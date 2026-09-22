# Sistema de Controle de Chamados Internos

Aplicação full stack para que funcionários abram chamados internos de suporte ("meu computador travou",
"preciso de uma cadeira nova"...) e a equipe de suporte acompanhe, atenda e distribua esse trabalho de
forma equilibrada — desafio técnico da [Codificar](https://codificar.com.br).

- **Backend**: Go, Clean Architecture + SOLID, MongoDB, Swagger, testes de unidade com 90%+ de cobertura na
  camada de negócio. Documentação completa em [`backend/README.md`](backend/README.md).
- **Frontend**: Vue 3 + TypeScript + Tailwind CSS, consumindo a API via Axios. Documentação completa em
  [`frontend/README.md`](frontend/README.md).

Este README cobre apenas **como subir tudo com um único comando**. Decisões de arquitetura, justificativas
técnicas e detalhes de cada camada estão nos READMEs de `backend/` e `frontend/` linkados acima.

## Como executar (Docker Compose)

Pré-requisito: [Docker](https://www.docker.com/) com Docker Compose.

Na raiz do repositório:

```bash
docker compose up -d --build
```

Isso sobe três serviços:

| Serviço | URL | Descrição |
|---|---|---|
| `frontend` | http://localhost:5173 | Interface Vue (build de produção servido via Nginx) |
| `backend` | http://localhost:8080 | API REST em Go |
| `backend` (Swagger) | http://localhost:8080/swagger/index.html | Documentação interativa de todos os endpoints |
| `mongo` | `localhost:27017` (uso interno) | Banco de dados, com dados persistidos em volume Docker |

Na primeira subida, o backend popula automaticamente 3 responsáveis padrão (Ana Souza, Bruno Lima, Carla
Mendes) para que a aplicação já esteja pronta para uso — não é necessário nenhum passo manual de seed.

Para acompanhar os logs:

```bash
docker compose logs -f
```

Para derrubar tudo (mantendo os dados do Mongo no volume):

```bash
docker compose down
```

Para derrubar tudo e apagar também os dados do Mongo:

```bash
docker compose down -v
```

## Rodando cada parte isoladamente (sem Docker Compose)

Útil durante o desenvolvimento, quando você quer hot-reload em uma das duas pontas. Veja o passo a passo
detalhado em cada README:

- [`backend/README.md`](backend/README.md#como-rodar) — API em Go + MongoDB.
- [`frontend/README.md`](frontend/README.md#como-rodar) — interface Vue, apontando para a API via
  `VITE_API_URL`.

## Estrutura do repositório

```
.
├── backend/           # API REST em Go (Clean Architecture, MongoDB, Swagger, testes)
├── frontend/           # Interface em Vue 3 + TypeScript + Tailwind
└── docker-compose.yml   # Orquestra os três serviços (frontend, backend, mongo) juntos
```

## Resumo das decisões principais

Justificativas completas estão nos READMEs de cada camada; o resumo:

- **Go no backend, com [chi](https://github.com/go-chi/chi) como roteador** — `chi` é 100% compatível com
  `net/http` (não é um framework como Gin/Echo), só torna path params, rotas aninhadas e middlewares mais
  ergonômicos que a stdlib pura, sem mudar a forma como os handlers são escritos.
- **Vue + Tailwind no frontend, sem Vue Router/Pinia** — a aplicação é uma tela só (lista + modal), então
  roteamento e uma store global adicionariam complexidade sem necessidade real.
- **MongoDB** — solicitado no desafio; o chamado é um documento praticamente plano, o que combina bem com
  um modelo de documento.
- **"Em aberto" para fins de distribuição automática** = status `aberto` ou `em_andamento`. A distribuição
  automática sempre escolhe o responsável ativo com menos chamados nesses dois status. Detalhes em
  [`backend/README.md`](backend/README.md#decisões-de-modelagem).
- **CORS totalmente aberto** no backend (todas as origens/métodos/headers), como solicitado.
