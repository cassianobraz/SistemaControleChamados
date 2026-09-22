# Frontend — Sistema de Controle de Chamados Internos

Interface em **Vue 3 + TypeScript + Tailwind CSS**, consumindo a [API do backend](../backend/README.md) via
Axios.

## Stack e por quê

| Escolha | Motivo |
|---|---|
| **Vue 3 + `<script setup lang="ts">` (Composition API)** | Padrão atual recomendado pelo próprio time do Vue; `<script setup>` reduz boilerplate e o TypeScript tipa props/emits/composables de ponta a ponta. |
| **Vite** | Build e dev server rápidos, HMR instantâneo — importante para produtividade num projeto pequeno. |
| **Tailwind CSS v4** | Permite montar a interface direto no markup, sem precisar criar e nomear uma folha de estilos própria para cada componente — atende ao pedido do desafio de "não reinventar a roda" no styleguide. |
| **Axios** | Cliente HTTP explícito (conforme pedido), com uma instância única (`src/api/http.ts`) e módulos por recurso (`tickets`, `responsibles`) em vez de chamadas soltas espalhadas pelos componentes. |
| **Sem Vue Router / Pinia** | A aplicação é uma única tela (lista de chamados + modal de criação/edição); adicionar roteamento ou uma store global só traria complexidade sem necessidade real. Estado é local aos composables (`useTickets`, `useResponsibles`), que é o suficiente para o escopo atual. |

## Arquitetura

```
src/
├── api/            → Axios: instância (http.ts) + funções tipadas por recurso (tickets.ts, responsibles.ts)
├── types/           → Interfaces TypeScript espelhando os DTOs do backend
├── composables/      → Estado + regras de tela (useTickets, useResponsibles) — a "camada de aplicação"
│                       do frontend: componentes não chamam a API diretamente, chamam o composable.
├── utils/            → Funções puras (formatação de data, labels/cores de status e prioridade)
├── components/
│   ├── ui/           → Primitivas reutilizáveis e sem regra de negócio (BaseButton, BaseBadge, BaseModal,
│   │                    PaginationControls)
│   └── layout/        → Casca da aplicação (AppHeader)
├── features/tickets/  → Componentes da feature de chamados:
│   ├── TicketsPage.vue          → Container: orquestra os composables e compõe os componentes abaixo
│   ├── TicketFilters.vue        → Filtros/busca/ordenação (v-model por campo)
│   ├── TicketList.vue           → Lista de chamados: tabela em telas md+, cartões empilhados no mobile
│   ├── TicketRow.vue            → Uma linha da tabela (badges, seleção de responsável, ações)
│   ├── TicketCard.vue           → Card equivalente à linha, usado no layout mobile
│   ├── TicketFormModal.vue      → Modal único para abrir e editar chamado
│   ├── ResponsibleSidebar.vue   → Painel de carga da equipe + ações de novo/editar/excluir responsável
│   └── ResponsibleFormModal.vue → Modal único para cadastrar e editar responsável
├── App.vue            → Raiz: só compõe AppHeader + TicketsPage
└── main.ts
```

`App.vue` e `TicketsPage.vue` seguem o princípio de "superfície de composição": eles não têm lógica própria
além de orquestrar composables e passar dados para componentes filhos via **props down / events up** — toda
mutação de estado (criar, editar, atribuir chamado) acontece nos composables, nunca diretamente nos
componentes de apresentação.

### Fluxo de dados

- `useTickets()` guarda a lista, filtros, paginação e expõe ações (`createTicket`, `updateTicket`,
  `assignTicket`) que já re-sincronizam a lista após cada mutação — a UI nunca fica com dado desatualizado.
- `useResponsibles()` guarda a lista de responsáveis com sua carga de chamados em aberto, usada tanto no
  seletor de filtros/formulário quanto no painel lateral, e expõe as ações de cadastro
  (`createResponsible`, `updateResponsible`, `deleteResponsible`) que já re-sincronizam a lista após cada
  mutação — mesmo padrão de `useTickets()`.
- Estado devolvido pelos composables é `readonly()`; a única forma de alterá-lo é pelas ações que o próprio
  composable expõe — evita que um componente filho mude a lista "por baixo dos panos".
- Reatribuição de responsável (`PATCH /tickets/{id}/assign`) é uma ação separada de editar o chamado
  (`PUT /tickets/{id}`), espelhando os dois endpoints do backend: a tabela permite tanto escolher um
  responsável manualmente (select por linha) quanto redistribuir automaticamente (botão "Auto"), e o modal
  de edição cuida de título/descrição/prioridade/status.

## Como rodar

### Opção 1 — Docker Compose (recomendado)

Veja o [README raiz](../README.md). A partir da raiz do repositório:

```bash
docker compose up -d --build
```

O frontend sobe em `http://localhost:5173` (build de produção servido via Nginx).

### Opção 2 — Rodando isoladamente

Pré-requisitos: Node.js 22+ e o [backend](../backend/README.md) rodando (local ou via Docker).

```bash
cd frontend
cp .env.example .env.local     # ajuste VITE_API_URL se a API não estiver em localhost:8080
npm install
npm run dev
```

A aplicação sobe em `http://localhost:5173`.

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `VITE_API_URL` | `http://localhost:8080` | URL base da API consumida pelo Axios. Em produção (Docker), é definida em tempo de build (ver `docker-compose.yml` na raiz). |

## Scripts

```bash
npm run dev       # servidor de desenvolvimento com HMR
npm run build      # type-check (vue-tsc) + build de produção em dist/
npm run preview     # serve o build de produção localmente, para conferência
```

## Decisões e trade-offs

- **CSS**: só os tokens de marca (`--color-brand-*`) foram promovidos a tema no `@theme` do Tailwind; o
  resto da interface é composto diretamente com utilitárias no markup, evitando uma camada de abstração de
  CSS desnecessária para o tamanho do projeto.
- **Sem testes automatizados no frontend**: o desafio pede cobertura de testes explicitamente para o
  backend; para manter o frontend simples e dentro do escopo, não foi adicionada uma suíte de testes aqui.
  A estrutura em composables (lógica isolada de componentes) foi pensada para que, se necessário, testes
  com Vitest + Vue Test Utils possam ser adicionados depois sem refatoração.
- **Paginação simples** (anterior/próxima) em vez de números de página, suficiente para o volume esperado
  de um time de suporte interno.
- **Layout responsivo com um componente de apresentação por breakpoint** (`TicketRow` para a tabela em
  telas médias/grandes, `TicketCard` para os cartões empilhados no mobile) em vez de uma única tabela com
  scroll horizontal — mantém a leitura confortável em telas pequenas sem duplicar estado, só marcação.
