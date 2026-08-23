# Helpdesk Event Hub

Event-driven helpdesk для практики микросервисной архитектуры: **Go**-сервисы, **gRPC** между ними, **Apache Kafka** для событий, UI на **Next.js**.

Клиент или агент создаёт тикет → система назначает исполнителя, считает SLA, шлёт уведомления и при просрочке эскалирует. Синхронные вызовы — по gRPC, побочные эффекты — через Kafka.

## Цели

- Bounded contexts и database-per-service
- **gRPC** как единственный sync-транспорт между сервисами
- Choreography через Kafka (fan-out, retries, DLQ)
- Transactional outbox и идемпотентные consumers
- SLA и escalation как отдельные сервисы
- BFF / API Gateway (HTTP снаружи → gRPC внутрь)
- Локальный стенд на Docker Compose

## Стек

| Слой | Технология |
|------|------------|
| Frontend | Next.js *(стек UI уточняется)* |
| Edge | Go API Gateway / BFF — **HTTP/REST** для браузера |
| Inter-service | **gRPC** (+ protobuf) |
| Event bus | Apache Kafka |
| Storage | PostgreSQL (per service), Redis (SLA timers) |
| Search | Meilisearch *(после MVP)* |
| Local | Docker Compose |

## Транспорт: что чем ходит

| Связь | Протокол | Зачем |
|-------|----------|--------|
| Next.js → `api-gateway` | HTTP/REST (JSON) | удобно для браузера / BFF |
| `api-gateway` → сервисы | **gRPC** | типизированные контракты, низкая латентность |
| сервис → сервис (sync) | **gRPC** | запросы данных / команды, когда нужен ответ сейчас |
| сервис → сервис (async) | **Kafka** | assign, SLA, notify, audit, search — choreography |

Правило: если вызывающему нужен **ответ в этом же запросе** — gRPC. Если это **побочный эффект / реакция на факт** — Kafka.

## Архитектура (сервисы)

| Сервис | Роль | Exposes |
|--------|------|---------|
| `api-gateway` | JWT, routing, aggregation для UI | HTTP |
| `auth-service` | Пользователи, роли, токены | gRPC |
| `ticket-service` | CRUD тикетов, статусы, outbox → Kafka | gRPC |
| `assignment-service` | Авто-назначение агента / очереди | gRPC + Kafka consumer |
| `sla-service` | Политики SLA, таймеры, `sla.breached` | gRPC + Kafka consumer |
| `escalation-service` | Повышение приоритета / смена очереди | gRPC + Kafka consumer |
| `notification-service` | Email / webhook / in-app (mock) | gRPC + Kafka consumer |
| `audit-service` | Append-only timeline | gRPC + Kafka consumer |
| `search-service` | Индексация и поиск *(опционально)* | gRPC + Kafka consumer |

### Happy path

```text
Portal --HTTP--> Gateway --gRPC--> ticket-service
                              └─→ Kafka: ticket.created
                                      → assignment → ticket.assigned
                                      → sla (timers)
                                      → notification + audit
              … SLA breach → escalation → ticket.escalated → notify L2

Gateway --gRPC--> audit-service / sla-service   # BFF: склеить карточку тикета
```

## MVP

1. Создать тикет (HTTP → gRPC) → событие `ticket.created`
2. Авто-назначение → `ticket.assigned`
3. Mock-уведомление + запись в audit timeline
4. SLA breach → эскалация → уведомление L2
5. Список тикетов и лента событий в UI (Gateway агрегирует по gRPC)

**Позже:** реальные каналы уведомлений, вложения, multi-tenant, AI triage, Kubernetes.

## Структура репозитория (план)

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
│   ├── eventkit/            # Kafka envelope, logging, health
│   └── grpckit/             # interceptors, metadata, errors
├── deploy/
│   └── compose/             # Kafka, Postgres, Redis, …
└── README.md
```

Репозиторий сейчас в стартовом состоянии: зафиксированы цели и ignore-правила, код сервисов появится по roadmap.

## Быстрый старт

> Появится после Фазы 0 (Compose + каркас сервисов + proto).

```bash
# планируется
docker compose -f deploy/compose/docker-compose.yml up -d
```

## Roadmap (кратко)

| Фаза | Что делаем |
|------|------------|
| 0 | Compose; `api/proto`; общий Go-каркас (`grpckit`, `eventkit`) |
| 1 | `ticket-service` (gRPC) + outbox + gateway (HTTP→gRPC) + auth |
| 2 | assignment, notification, audit |
| 3 | SLA + escalation + DLQ |
| 4 | search + BFF aggregation по gRPC |
| 5 | Next.js MVP |
| 6 | tracing (gRPC + Kafka), load/chaos, опционально K8s |

## Принципы

- Снаружи (браузер) — HTTP/REST; внутри — **только gRPC** для sync
- Async side-effects — Kafka; partition key = `ticket_id`
- Consumers идемпотентны по `event_id`
- Контракты версионируются через protobuf (`api/proto`)
- Нет shared DB между сервисами
- gRPC metadata: `authorization`, `x-correlation-id`, `x-request-id`

## Лицензия

TBD
