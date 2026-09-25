# HelpDesk — Universal Support Desk

Встраиваемый **чат-виджет** для сайтов + **консоль support desk**. Один запуск обслуживает один проект: описание и инструкции бота задаются через `PROJECT_*` в `.env` конкретного deployment. Посетитель пишет в пузырь справа снизу → бот (OpenRouter) отвечает → handoff создаёт связанный bot-тикет → менеджер отвечает в UI.

Тикеты/SLA/escalation остаются в compose как legacy Event Hub, но primary UX — **беседы**.

Репозиторий: [github.com/Glistand/HelpDesk](https://github.com/Glistand/HelpDesk) — monorepo.

## Продукт (эта фаза)

| Компонент | Назначение |
|-----------|------------|
| `apps/widget` | FAB + чат; `<script src="…/widget.js" data-site-key="demo-site" data-gateway="http://localhost:8080">` |
| `conversation-service` | сессии/сообщения, OpenRouter, NATS `helpdesk.conversation.>` |
| `api-gateway` | публичный `/widget/*` + JWT `/conversations*` |
| `apps/web` | inbox бесед, тред, ответ агента |

Демо-страница виджета: `http://localhost:3001/widget-demo.html` (после `make web` или compose `web`).

## Цели (платформа)

- Bounded contexts и database-per-service
- **gRPC** как единственный sync-транспорт между сервисами
- Choreography через NATS JetStream (fan-out, retries, DLQ)
- Transactional outbox и идемпотентные consumers
- SLA и escalation как отдельные сервисы *(не в primary UX)*
- BFF / API Gateway (HTTP снаружи → gRPC внутрь)
- Локальный и первый деплой — **Docker Compose**

## Профиль проекта и роли

Каждый Compose-стек — отдельный проект, а не tenant внутри общей БД. Значения `PROJECT_*` из
`.env` являются конфигурацией конкретного deployment. Пример по умолчанию — SalonPro, платформа
записи и управления салонами красоты.

Роль `admin` видит всю очередь. Роль `agent` получает в API только
назначенные ей тикеты и чаты. Тикеты содержат источник `manager` или `bot` и могут хранить
связанный `conversation_id`.

## Стек

| Слой | Технология |
|------|------------|
| Frontend | Next.js *(стек UI уточняется)* |
| Edge | Go API Gateway / BFF — **HTTP/REST** для браузера |
| Inter-service | **gRPC** (+ protobuf) |
| Event bus | **NATS JetStream** |
| Storage | PostgreSQL (per service), Redis (SLA timers) |
| Search | Meilisearch *(после MVP)* |
| Deploy | **Docker Compose** |

## Транспорт: что чем ходит

| Связь | Протокол | Зачем |
|-------|----------|--------|
| Next.js → `api-gateway` | HTTP/REST (JSON) | удобно для браузера / BFF |
| `api-gateway` → сервисы | **gRPC** | типизированные контракты, низкая латентность |
| сервис → сервис (sync) | **gRPC** | запросы данных / команды, когда нужен ответ сейчас |
| сервис → сервис (async) | **NATS JetStream** | assign, SLA, notify, audit, search — choreography |

Правило: если вызывающему нужен **ответ в этом же запросе** — gRPC. Если это **побочный эффект / реакция на факт** — NATS.

## Архитектура (сервисы)

| Сервис | Роль | Exposes |
|--------|------|---------|
| `api-gateway` | JWT, routing, widget public API, aggregation | HTTP |
| `auth-service` | Пользователи, роли, токены | gRPC |
| `conversation-service` | Беседы, сообщения, OpenRouter bot | gRPC |
| `ticket-service` | CRUD тикетов, статусы, outbox → NATS *(legacy UX)* | gRPC |
| `assignment-service` | Авто-назначение агента / очереди | gRPC + NATS consumer |
| `sla-service` | Политики SLA, таймеры, `sla.breached` | gRPC + NATS consumer |
| `escalation-service` | Повышение приоритета / смена очереди | gRPC + NATS consumer |
| `notification-service` | Email / webhook / in-app (mock) | gRPC + NATS consumer |
| `audit-service` | Append-only timeline | gRPC + NATS consumer |
| `search-service` | Индексация и поиск *(опционально)* | gRPC + NATS consumer |

### Happy path

```text
Portal --HTTP--> Gateway --gRPC--> ticket-service
                              └─→ NATS: helpdesk.ticket.created
                                      → assignment → helpdesk.ticket.assigned
                                      → sla (timers)
                                      → notification + audit
              … SLA breach → escalation → helpdesk.ticket.escalated → notify L2

Gateway --gRPC--> audit-service / sla-service   # BFF: склеить карточку тикета
```

## MVP

1. Создать тикет (HTTP → gRPC) → событие `helpdesk.ticket.created`
2. Авто-назначение → `helpdesk.ticket.assigned`
3. Mock-уведомление + запись в audit timeline
4. SLA breach → эскалация → уведомление L2
5. Список тикетов и лента событий в UI (Gateway агрегирует по gRPC)

**Позже:** реальные каналы уведомлений, вложения, multi-tenant, AI triage.

## Структура репозитория

```text
HelpDesk/
├── apps/
│   ├── web/                 # Next.js agent console (conversations)
│   └── widget/              # embeddable chat FAB → builds to web/public/widget.js
├── api/
│   └── proto/               # .proto контракты + buf/protoc codegen
├── services/
│   ├── api-gateway/         # HTTP → gRPC (+ /widget public API)
│   ├── conversation-service/# chat + OpenRouter bot
│   ├── auth-service/
│   ├── ticket-service/
│   ├── assignment-service/
│   ├── sla-service/
│   ├── escalation-service/
│   ├── notification-service/
│   ├── audit-service/
│   └── search-service/
├── libs/
│   ├── eventkit/            # NATS envelope, logging, health
│   └── grpckit/             # interceptors, metadata, errors
├── deploy/
│   └── compose/             # Docker Compose: NATS, Postgres, Redis, …
├── scripts/
│   └── smoke-widget.sh      # widget → bot → handoff → agent
├── Makefile
└── README.md
```

## Быстрый старт

```bash
cp .env.example .env
# optional: set OPENROUTER_API_KEY for live bot replies (free model by default)
make up
```

Поднимаются:

| Сервис | Порт | Назначение |
|--------|------|------------|
| PostgreSQL | 5432 | отдельные БД на сервис (`auth`, `ticket`, `conversation`, …) |
| Redis | 6379 | SLA timers |
| NATS | 4222 | клиентский порт |
| NATS monitor | 8222 | health / metrics |
| Meilisearch | 7700 | поиск *(после MVP)* |
| `auth-service` | 50051 | gRPC JWT |
| `ticket-service` | 50052 | gRPC tickets + outbox → NATS |
| `assignment-service` | 50053 | auto-assign (round-robin L1) |
| `audit-service` | 50054 | timeline gRPC |
| `sla-service` | 50055 | Redis timers → `sla.warned` / `sla.breached` |
| `escalation-service` | 50056 | on breach → L2 assign + `ticket.escalated` |
| `search-service` | 50057 | Meilisearch index + gRPC search |
| `conversation-service` | 50058 | chat + OpenRouter |
| `notification-service` | — | mock notify → `notification.sent` |
| `api-gateway` | 8080 | HTTP → gRPC (widget + conversations + tickets) |
| `web` (Next.js) | 3001 | agent console → gateway (`GATEWAY_URL`) |
| Jaeger UI | 16686 | traces (OTLP `:4318`) |

JetStream streams создаются автоматически контейнером `nats-init`:

- `HELP_DESK_EVENTS` — subjects `helpdesk.ticket.>`, `helpdesk.conversation.>`, `helpdesk.sla.>`, …
- `HELP_DESK_DLQ` — subjects `helpdesk.dlq.>` (отдельный stream, без overlap)

Проверка:

```bash
make ps
curl http://localhost:8222/healthz
curl http://localhost:8080/healthz

# Support desk smoke: widget → bot → handoff → agent
./scripts/smoke-widget.sh
# make smoke-widget

# Existing Postgres volume without `conversation` DB:
# docker exec helpdesk-postgres psql -U helpdesk -d postgres -c 'CREATE DATABASE conversation;'

# login (seed: agent@helpdesk.local / password)
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}' | jq -r .access_token)

# Legacy tickets still work
curl -s -X POST http://localhost:8080/tickets \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"VPN down","description":"Cannot connect","priority":"high","category":"Network"}'

# Phase 2–6 happy path (tickets/SLA)
make e2e
make e2e-phase6

# UI (dev, без Docker)
cd apps/web && cp .env.example .env.local   # GATEWAY_URL=http://localhost:8080
make web                                    # http://localhost:3001
# widget demo: http://localhost:3001/widget-demo.html
# seed: agent@helpdesk.local / password

# rebuild embed script after editing apps/widget
cd apps/widget && npm install && npm run build
```

Остановка:

```bash
make down      # сохранить данные
make reset     # удалить volumes
```

## CI: Docker → GHCR

Workflow [`.github/workflows/docker-publish.yml`](.github/workflows/docker-publish.yml) на каждый push в `main` / тег `v*` / `workflow_dispatch` собирает все сервисы и `web`, пушит в GitHub Container Registry:

| Image | Пример |
|-------|--------|
| `ghcr.io/<owner>/helpdesk-api-gateway` | `:latest`, `:sha-<short>`, `:v1.2.3` |
| `ghcr.io/<owner>/helpdesk-web` | то же |
| … | все `*-service` из матрицы |

На PR образы только **собираются** (без push).

После первого успешного run:

1. GitHub → **Packages** у репозитория — появятся `helpdesk-*`
2. Для private package: Settings → Package → Manage Actions access (репо уже связано через `GITHUB_TOKEN`)
3. Локально (если пакеты private):

```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u YOUR_GITHUB_USER --password-stdin
GHCR_OWNER=glistand IMAGE_TAG=latest make up-ghcr
```

`make up-ghcr` тянет образы через [`deploy/compose/docker-compose.ghcr.yml`](deploy/compose/docker-compose.ghcr.yml) и поднимает стек без локальной сборки.

### Деплой на сервер

Два файла — [`deploy/server/docker-compose.yml`](deploy/server/docker-compose.yml) + [`deploy/server/.env`](deploy/server/.env). Репозиторий на сервере не нужен.

```bash
# скопируйте оба файла в каталог на сервере, затем:
docker compose --env-file .env up -d
# UI :80   API :8080
```

Подробности — [`deploy/server/README.md`](deploy/server/README.md).

## Go monorepo

Workspace: [`go.work`](go.work) включает libs, codegen и сервисы Фаз 0–4. UI — [`apps/web`](apps/web) (Next.js App Router).

| Пакет | Назначение |
|-------|------------|
| [`api/proto`](api/proto) | `.proto` + `buf` codegen → [`api/gen/go`](api/gen/go) |
| [`libs/grpckit`](libs/grpckit) | gRPC interceptors, metadata, errors, OTel stats |
| [`libs/eventkit`](libs/eventkit) | NATS JetStream envelope / publish / subscribe / DLQ + consume spans |
| [`libs/otelkit`](libs/otelkit) | OTLP/HTTP tracer bootstrap |
| [`services/auth-service`](services/auth-service) | Login / ValidateToken (JWT) |
| [`services/ticket-service`](services/ticket-service) | CRUD + transactional outbox |
| [`services/assignment-service`](services/assignment-service) | consume created → AssignTicket |
| [`services/notification-service`](services/notification-service) | mock notify + `notification.sent` |
| [`services/audit-service`](services/audit-service) | append-only timeline |
| [`services/sla-service`](services/sla-service) | Redis SLA timers |
| [`services/escalation-service`](services/escalation-service) | breach → L2 |
| [`services/search-service`](services/search-service) | Meilisearch indexer + search |
| [`services/conversation-service`](services/conversation-service) | conversations + OpenRouter |
| [`services/api-gateway`](services/api-gateway) | HTTP BFF (`/widget/*`, `/conversations`, tickets) |
| [`apps/web`](apps/web) | Next.js: conversations inbox/thread + legacy tickets |
| [`apps/widget`](apps/widget) | Embeddable support chat → `public/widget.js` |

```bash
make proto
make build-services
make test-go
make e2e
```

Dev seed user: `agent@helpdesk.local` / `password`.

Compose SLA defaults (override via `.env`): `SLA_FIRST_RESPONSE=5s`, `SLA_RESOLVE=30s`, `SLA_WARN_RATIO=0.5`.

## Roadmap

| Фаза | Что делаем |
|------|------------|
| 0 | Compose; `api/proto`; общий Go-каркас (`grpckit`, `eventkit`) — **done** |
| 1 | `ticket-service` (gRPC) + outbox + gateway (HTTP→gRPC) + auth — **done** |
| 2 | assignment, notification, audit — **done** |
| 3 | SLA + escalation + DLQ (`helpdesk.dlq.>`) — **done** |
| 4 | search + BFF aggregation по gRPC — **done** |
| 5 | Next.js MVP (live API, cookie auth, compose `web`) — **done** |
| 6 | tracing (gRPC + NATS + Jaeger), load/chaos, hardening — **done** |

## Принципы

- Снаружи (браузер) — HTTP/REST; внутри — **только gRPC** для sync
- Async side-effects — NATS JetStream; subject key включает `ticket_id` для ordering
- Consumers идемпотентны по `event_id` (NATS dedup window + app-level dedup)
- Контракты версионируются через protobuf (`api/proto`)
- Нет shared DB между сервисами
- gRPC metadata: `authorization`, `x-correlation-id`, `x-request-id`
- Деплой — Docker Compose; Kubernetes не в scope первой версии

## NATS: subjects (вместо Kafka topics)

| Subject | Кто публикует | Кто слушает |
|---------|---------------|-------------|
| `helpdesk.ticket.created` | ticket-service | assignment, sla, notification, audit, search |
| `helpdesk.ticket.assigned` | ticket-service (после AssignTicket) | sla, notification, audit |
| `helpdesk.ticket.updated` | ticket-service | audit, search |
| `helpdesk.notification.sent` | notification-service | audit |
| `helpdesk.sla.warned` | sla-service | notification |
| `helpdesk.sla.breached` | sla-service | escalation, notification |
| `helpdesk.ticket.escalated` | escalation-service | notification, audit |
| `helpdesk.conversation.created` | conversation-service | (observability / future) |
| `helpdesk.conversation.message` | conversation-service | (observability / future) |
| `helpdesk.conversation.handoff` | conversation-service | (observability / future) |
| `helpdesk.dlq.>` | любой consumer | DLQ stream, ручной replay |

## Лицензия

TBD
