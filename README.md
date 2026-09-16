# Helpdesk Event Hub

Event-driven helpdesk: **Go**-сервисы, **gRPC** между ними, **NATS JetStream** для событий, UI на **Next.js**.

Клиент или агент создаёт тикет → система назначает исполнителя, считает SLA, шлёт уведомления и при просрочке эскалирует. Синхронные вызовы — по gRPC, побочные эффекты — через NATS.

Репозиторий: [github.com/Glistand/HelpDesk](https://github.com/Glistand/HelpDesk) — monorepo, без отдельных GitLab-проектов.

## Цели

- Bounded contexts и database-per-service
- **gRPC** как единственный sync-транспорт между сервисами
- Choreography через NATS JetStream (fan-out, retries, DLQ)
- Transactional outbox и идемпотентные consumers
- SLA и escalation как отдельные сервисы
- BFF / API Gateway (HTTP снаружи → gRPC внутрь)
- Локальный и первый деплой — **Docker Compose**

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
| `api-gateway` | JWT, routing, aggregation для UI | HTTP |
| `auth-service` | Пользователи, роли, токены | gRPC |
| `ticket-service` | CRUD тикетов, статусы, outbox → NATS | gRPC |
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
│   └── web/                 # Next.js
├── api/
│   └── proto/               # .proto контракты + buf/protoc codegen
├── services/
│   ├── api-gateway/         # HTTP → gRPC clients
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
├── Makefile
└── README.md
```

## Быстрый старт

```bash
cp .env.example .env
make up
```

Поднимаются:

| Сервис | Порт | Назначение |
|--------|------|------------|
| PostgreSQL | 5432 | отдельные БД на сервис (`auth`, `ticket`, …) |
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
| `notification-service` | — | mock notify → `notification.sent` |
| `api-gateway` | 8080 | HTTP → gRPC |

JetStream streams создаются автоматически контейнером `nats-init`:

- `HELP_DESK_EVENTS` — subjects `helpdesk.ticket.>`, `helpdesk.sla.>`, …
- `HELP_DESK_DLQ` — subjects `helpdesk.dlq.>` (отдельный stream, без overlap)

Проверка:

```bash
make ps
curl http://localhost:8222/healthz
curl http://localhost:8080/healthz

# login (seed: agent@helpdesk.local / password)
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}' | jq -r .access_token)

curl -s -X POST http://localhost:8080/tickets \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"VPN down","description":"Cannot connect","priority":"high","category":"Network"}'

# Phase 2+3 happy path
make e2e          # assign/notify/audit + SLA breach → L2
make e2e-phase3   # только SLA/escalation
```

Остановка:

```bash
make down      # сохранить данные
make reset     # удалить volumes
```

## Go monorepo

Workspace: [`go.work`](go.work) включает libs, codegen и сервисы Фаз 0–3.

| Пакет | Назначение |
|-------|------------|
| [`api/proto`](api/proto) | `.proto` + `buf` codegen → [`api/gen/go`](api/gen/go) |
| [`libs/grpckit`](libs/grpckit) | gRPC interceptors, metadata, errors |
| [`libs/eventkit`](libs/eventkit) | NATS JetStream envelope / publish / subscribe / DLQ |
| [`services/auth-service`](services/auth-service) | Login / ValidateToken (JWT) |
| [`services/ticket-service`](services/ticket-service) | CRUD + transactional outbox |
| [`services/assignment-service`](services/assignment-service) | consume created → AssignTicket |
| [`services/notification-service`](services/notification-service) | mock notify + `notification.sent` |
| [`services/audit-service`](services/audit-service) | append-only timeline |
| [`services/sla-service`](services/sla-service) | Redis SLA timers |
| [`services/escalation-service`](services/escalation-service) | breach → L2 |
| [`services/api-gateway`](services/api-gateway) | HTTP BFF (`timeline`, `sla`) |

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
| 4 | search + BFF aggregation по gRPC |
| 5 | Next.js MVP (UI preview already in `apps/web`) |
| 6 | tracing (gRPC + NATS), load/chaos, hardening |

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
| `helpdesk.dlq.>` | любой consumer | DLQ stream, ручной replay |

## Лицензия

TBD
