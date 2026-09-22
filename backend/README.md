# Backend — Sistema de Controle de Chamados Internos

API REST em **Go** para o desafio técnico da Codificar: um sistema simples para que funcionários abram
chamados internos ("meu computador travou", "preciso de uma cadeira nova"...) e a equipe de suporte
acompanhe, atenda e distribua esse trabalho de forma equilibrada.

## Sumário

- [Stack e por quê](#stack-e-por-quê)
- [Arquitetura](#arquitetura)
- [Decisões de modelagem](#decisões-de-modelagem)
- [Como rodar](#como-rodar)
- [Variáveis de ambiente](#variáveis-de-ambiente)
- [Documentação da API (Swagger)](#documentação-da-api-swagger)
- [Testes e cobertura](#testes-e-cobertura)
- [Estrutura de pastas](#estrutura-de-pastas)
- [Trade-offs e próximos passos](#trade-offs-e-próximos-passos)

## Stack e por quê

| Escolha | Motivo |
|---|---|
| **Go + [chi](https://github.com/go-chi/chi)** | `chi` é um roteador HTTP leve, 100% compatível com `net/http` (`http.Handler`/`http.HandlerFunc` puros, sem tipos próprios de request/response) — não é um framework web como Gin ou Echo. Foi escolhido por dar path params, roteamento aninhado (`Route`) e composição de middlewares mais ergonômicos que o `http.ServeMux` da stdlib, mantendo os handlers e o restante da aplicação exatamente como já eram: funções `func(http.ResponseWriter, *http.Request)` comuns. |
| **MongoDB** | Solicitado no desafio. O documento de chamado é praticamente plano (sem relacionamentos complexos), o que combina bem com um modelo de documento; e como o Go driver oficial (`mongo-driver/v2`) é bem tipado, a camada de infraestrutura fica simples de isolar atrás de interfaces. |
| **Clean Architecture (domain → usecase → delivery/infrastructure)** | Mantém a regra de negócio (ex.: "atribuir ao responsável com menos chamados em aberto") isolada de HTTP e de MongoDB. Isso permite testar 100% da lógica com mocks, sem banco de dados, e trocar qualquer camada externa sem tocar nas regras. |
| **OpenAPI escrito à mão + Swagger UI embutido** | Em vez de depender de geração de código (`swag init`) como passo extra de build, o `docs/openapi.json` é a fonte da verdade, versionado e revisável em PR como qualquer outro arquivo, e é embutido no binário via `go:embed`. |
| **`testify` (assert/require/mock)** | Reduz boilerplate de asserções e de mocks de interface, mantendo os testes tabulares idiomáticos em Go. |
| **CORS totalmente aberto** | Conforme solicitado — como o desafio não define domínios de produção, liberar todas as origens/métodos/headers evita fricção ao rodar o frontend em portas/hosts diferentes (Vite dev server, Docker, etc.). |

## Arquitetura

```
cmd/api            → composição (main.go): lê config, conecta no Mongo, injeta dependências, sobe o servidor
internal/domain     → entidades e regras invariantes (Ticket, Responsible, Priority, Status) + interfaces
                       de repositório (portas). Não importa nada de fora do próprio domínio.
internal/usecase    → regras de aplicação (abrir, editar, listar, atribuir chamado). Depende só das
                       interfaces do domain — testável com mocks, sem banco.
internal/delivery/http → adaptador HTTP: handlers, DTOs de request/response, middlewares (CORS, logging)
                       e o router. Depende das interfaces do usecase, nunca de MongoDB diretamente.
internal/infrastructure → adaptadores concretos: conexão Mongo, repositórios Mongo, config via env,
                       seed dos responsáveis padrão.
internal/mocks      → mocks (testify) das interfaces de domain/usecase, usados nos testes.
docs/               → especificação OpenAPI (openapi.json) embutida no binário via go:embed.
```

A regra de dependência é sempre "de fora para dentro": `delivery` e `infrastructure` dependem de
`usecase`/`domain` através de interfaces; nunca o contrário. Isso é o que permite os testes de unidade do
`usecase` e do `delivery/http` rodarem sem MongoDB algum — eles recebem um **mock** da interface de
repositório/usecase no lugar da implementação real.

### SOLID na prática

- **S**RP — cada arquivo tem uma responsabilidade (ex.: `ticket.go` só valida e muta o agregado; o
  handler só traduz HTTP ↔ usecase; o repositório só traduz domínio ↔ Mongo).
- **O**CP — novas regras de distribuição ou de filtro se adicionam sem alterar o contrato das interfaces.
- **L**SP — qualquer implementação de `domain.TicketRepository` (Mongo real ou mock) pode substituir a
  outra sem quebrar o `usecase`.
- **I**SP — interfaces pequenas e específicas (`TicketRepository`, `ResponsibleRepository`,
  `TicketUseCase`, `ResponsibleUseCase`), sem métodos que um consumidor não usa.
- **D**IP — `usecase` e `delivery` dependem de abstrações (`interface`) definidas no próprio consumidor
  (`domain` define a interface que `infrastructure` implementa), não de implementações concretas.

## Decisões de modelagem

- **Status do chamado**: `aberto`, `em_andamento`, `resolvido`, `fechado`.
- **"Em aberto" para fins de distribuição de carga** (item 4.3 do desafio): um chamado conta como carga
  em aberto de um responsável enquanto seu status é `aberto` **ou** `em_andamento`. Assim que o chamado
  vira `resolvido` ou `fechado`, ele deixa de contar — a lógica está em `domain.Status.IsOpen()` e é usada
  tanto na distribuição automática quanto na listagem/dashboard de responsáveis.
- **Distribuição automática**: ao abrir um chamado sem `responsible_id`, ou ao chamar
  `PATCH /tickets/{id}/assign` sem corpo, o sistema busca todos os responsáveis ativos, conta os chamados
  em aberto de cada um (`CountOpenByResponsible`) e escolhe o de menor carga; em empate, sorteia entre os
  empatados (`math/rand/v2`) em vez de sempre escolher o primeiro da lista — do contrário, quem viesse
  depois na listagem praticamente nunca seria escolhido. Responsáveis com `active=false` são ignorados na
  distribuição automática, mas continuam
  atribuíveis manualmente (útil para alguém temporariamente afastado, sem apagar o histórico dele).
- **Responsáveis**: 3 responsáveis padrão (Ana Souza, Bruno Lima, Carla Mendes) são criados
  automaticamente na primeira subida da aplicação (`internal/infrastructure/seed`), para que a aplicação já
  tenha quem selecionar. A partir daí, o cadastro é gerenciado por completo pela API
  (`POST`/`PUT`/`DELETE /api/v1/responsibles`): criar, editar nome/e-mail/indicador de ativo e remover.
- **IDs como UUID v7 (pacote `uuid` da stdlib, Go 1.27+)**: `Ticket` e `Responsible` geram seu próprio
  `ID` no construtor do domínio (`uuid.NewV7().String()`) em vez de depender do `ObjectID` do Mongo —
  desacopla o identificador de negócio do banco escolhido, e o v7 preserva a ordenação temporal que o
  `ObjectID` já dava. Documentos gravados antes dessa mudança (com `_id` ainda como `ObjectID`) continuam
  funcionando: o cliente Mongo decodifica esses IDs legados como string hex (`ObjectIDAsHexString`, em
  `database.Connect`), e uma migração roda a cada boot (`database.MigrateLegacyObjectIDs`) convertendo o
  `_id` desses documentos para o mesmo formato string — necessário porque o Mongo não permite alterar o
  `_id` de um documento já existente in-place.
- **Edição de chamado** (`PUT /tickets/{id}`) substitui título, descrição, prioridade e status em uma
  única chamada — mais simples de usar no frontend do que múltiplos endpoints PATCH parciais, e ainda
  assim mantém a responsabilidade de reatribuir separada (`PATCH /tickets/{id}/assign`), já que atribuir
  responsável é uma decisão diferente de editar o conteúdo do chamado.
- **Listagem**: suporta filtro por status, prioridade, responsável, busca textual (título/descrição,
  case-insensitive) e ordenação por data de abertura ou por prioridade, com paginação (`page`/`page_size`,
  máximo de 100 itens por página) — o suficiente para o dia a dia da pessoa de suporte acompanhar a fila
  sem virar uma tela sobrecarregada.

## Como rodar

### Opção 1 — Docker Compose (recomendado)

O `docker-compose.yml` que sobe backend + frontend + MongoDB juntos fica na **raiz do repositório**
(fora de `backend/` e `frontend/`). Veja o [README raiz](../README.md) para o passo a passo completo.

Resumo, a partir da raiz do repositório:

```bash
docker compose up -d --build
```

A API sobe em `http://localhost:8080`.

### Opção 2 — Rodando o backend isoladamente

Pré-requisitos: Go 1.27+ e um MongoDB acessível (local, Docker ou Atlas).

```bash
cd backend
cp .env.example .env        # ajuste MONGO_URI se necessário (ex.: mongodb://localhost:27017)
go mod download

# subindo um MongoDB local rapidamente, se não tiver um:
docker run -d --name mongo-chamados -p 27017:27017 mongo:7

go run ./cmd/api
```

A API sobe em `http://localhost:8080` (ou na porta definida em `PORT`). Na primeira subida com o banco
vazio, os 3 responsáveis padrão são criados automaticamente.

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `PORT` | `8080` | Porta HTTP da API. |
| `MONGO_URI` | `mongodb://localhost:27017` | String de conexão do MongoDB. |
| `MONGO_DATABASE` | `controle_chamados` | Nome do banco de dados usado pela aplicação. |

Veja `.env.example`.

## Documentação da API (Swagger)

Com o servidor rodando, a documentação interativa (Swagger UI) fica em:

```
http://localhost:8080/swagger/index.html
```

O JSON puro da especificação OpenAPI 3.0 fica em `http://localhost:8080/swagger/doc.json` (fonte:
`docs/openapi.json`, versionado no repositório).

### Endpoints

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Verifica disponibilidade da API. |
| `POST` | `/api/v1/tickets` | Abre um novo chamado (`responsible_id` omitido = distribuição automática). |
| `GET` | `/api/v1/tickets` | Lista chamados com filtros, busca, ordenação e paginação. |
| `GET` | `/api/v1/tickets/{id}` | Busca um chamado por id. |
| `PUT` | `/api/v1/tickets/{id}` | Edita título, descrição, prioridade e status de um chamado. |
| `PATCH` | `/api/v1/tickets/{id}/assign` | Reatribui o responsável (manual ou automático). |
| `GET` | `/api/v1/responsibles` | Lista responsáveis disponíveis com a contagem de chamados em aberto de cada um. |
| `POST` | `/api/v1/responsibles` | Cadastra um novo responsável. |
| `PUT` | `/api/v1/responsibles/{id}` | Edita nome, e-mail e o indicador de ativo de um responsável. |
| `DELETE` | `/api/v1/responsibles/{id}` | Remove um responsável. |

## Testes e cobertura

```bash
make test          # todos os testes de unidade (equivalente a: go test ./...)
make test-cover     # cobertura da camada de negócio, com relatório HTML em coverage.html
make test-integration  # testes de integração dos repositórios Mongo (requer um MongoDB acessível)
```

Estratégia de testes:

- **`domain`, `usecase`, `delivery/http` (handlers, DTOs, middlewares, router), `infrastructure/config`,
  `infrastructure/seed`** são cobertos por **testes de unidade puros**, sem tocar em rede/banco — os
  repositórios são substituídos por mocks (`internal/mocks`, feitos com `testify/mock`). Rodando
  `make test-cover`, a cobertura combinada dessas camadas fica em **~96%** (meta do desafio: 90%+).
- **`infrastructure/repository` e `infrastructure/database`** são os adaptadores que conversam
  diretamente com o driver do MongoDB. A parte pura deles (mapeamento BSON ↔ domínio, montagem de filtro/
  ordenação) tem testes de unidade normais; as operações que exigem um MongoDB real (`Create`, `Update`,
  `FindByID`, `FindAll`, `CountOpenByResponsible`, ...) têm **testes de integração** próprios, atrás da
  build tag `integration`, seguindo a prática usual em Go de não misturar teste de unidade (rápido,
  isolado) com teste de integração (precisa de infraestrutura real). Rode-os com `make test-integration`
  apontando `MONGO_URI` para um Mongo disponível.

## Estrutura de pastas

```
backend/
├── cmd/api/main.go
├── docs/
│   ├── docs.go            # go:embed do openapi.json
│   └── openapi.json
├── internal/
│   ├── domain/             # entidades, regras invariantes, interfaces de repositório
│   ├── usecase/             # regras de aplicação (orquestra domain + repositórios)
│   ├── delivery/http/
│   │   ├── dto/             # request/response da API
│   │   ├── handler/         # handlers HTTP
│   │   ├── middleware/      # CORS, logging estruturado
│   │   └── router.go
│   ├── infrastructure/
│   │   ├── config/          # leitura de variáveis de ambiente
│   │   ├── database/        # conexão com o MongoDB
│   │   ├── repository/      # implementação Mongo das interfaces de domain
│   │   └── seed/            # popula os responsáveis padrão no primeiro boot
│   └── mocks/                # mocks de teste (testify) das interfaces de domain/usecase
├── Dockerfile
├── Makefile
└── go.mod
```

## Trade-offs e próximos passos

Dado o escopo do desafio ("nível básico", ~8h), ficaram de fora conscientemente:

- **Autenticação/autorização** — o desafio não pede login; qualquer pessoa com acesso à API pode abrir e
  editar chamados, como no cenário descrito (um time pequeno usando a ferramenta internamente).
- **Cache/rate limiting** — desnecessário para o volume de uma equipe interna pequena.
- **Testes de carga** — fora do escopo de um MVP interno.

Um próximo passo natural, se o produto crescesse, seria adicionar autenticação simples (ex.: sessão por
responsável logado) e notificações quando um chamado é atribuído.
