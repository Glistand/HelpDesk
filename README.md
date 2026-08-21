# Helpdesk Event Hub

Event-driven helpdesk для практики микросервисной архитектуры: **Go**-сервисы общаются через **Apache Kafka**, UI на **Next.js**.

Клиент или агент создаёт тикет → система назначает исполнителя, считает SLA, шлёт уведомления и при просрочке эскалирует — всё через события, а не через один god-service.

## Цели

- Bounded contexts и database-per-service
- Choreography через Kafka (fan-out, retries, DLQ)
- Transactional outbox и идемпотентные consumers
- SLA и escalation как отдельные сервисы
- BFF / API Gateway для фронтенда
- Локальный стенд на Docker Compose

## Стек

| Слой | Технология |
|------|------------|
| Frontend | Next.js *(стек UI уточняется)* |
| Edge | Go API Gateway / BFF |
| Services | Go microservices |
| Event bus | Apache Kafka |
| Storage | PostgreSQL (per service), Redis (SLA timers) |
| Search | Meilisearch *(после MVP)* |
| Local | Docker Compose |

## Архитектура (сервисы)

| Сервис | Роль |
|--------|------|
| `api-gateway` | Единая точка входа для Next.js, JWT, routing |
| `auth-service` | Пользователи, роли, токены |
| `ticket-service` | CRUD тикетов, статусы, outbox → Kafka |
| `assignment-service` | Авто-назначение агента / очереди |
| `sla-service` | Политики SLA, таймеры, `sla.breached` |
| `escalation-service` | Повышение приоритета / смена очереди |
| `notification-service` | Email / webhook / in-app (mock на старте) |
| `audit-service` | Append-only timeline событий |
| `search-service` | Индексация и поиск *(опционально)* |

### Happy path

```text
Portal → Gateway → ticket-service
                 → Kafka: ticket.created
                         → assignment → ticket.assigned
                         → sla (timers)
                         → notification + audit
         … SLA breach → escalation → ticket.escalated → notify L2
```

## MVP

1. Создать тикет → событие `ticket.created`
2. Авто-назначение → `ticket.assigned`
3. Mock-уведомление + запись в audit timeline
4. SLA breach → эскалация → уведомление L2
5. Список тикетов и лента событий в UI

**Позже:** реальные каналы уведомлений, вложения, multi-tenant, AI triage, Kubernetes.

## Структура репозитория (план)

```text
HelpDesk/
├── apps/
│   └── web/                 # Next.js
├── services/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── ticket-service/
│   ├── assignment-service/
│   ├── sla-service/
│   ├── escalation-service/
│   ├── notification-service/
│   ├── audit-service/
│   └── search-service/
├── libs/
│   └── eventkit/            # envelope, logging, health
├── deploy/
│   └── compose/             # Kafka, Postgres, Redis, …
└── README.md
```

Репозиторий сейчас в стартовом состоянии: зафиксированы цели и ignore-правила, код сервисов появится по roadmap.

## Быстрый старт

> Появится после Фазы 0 (Compose + каркас сервисов).

```bash
# планируется
docker compose -f deploy/compose/docker-compose.yml up -d
```

## Roadmap (кратко)

| Фаза | Что делаем |
|------|------------|
| 0 | Compose: Kafka, Postgres, Redis; общий Go-каркас |
| 1 | `ticket-service` + outbox + thin gateway + auth |
| 2 | assignment, notification, audit |
| 3 | SLA + escalation + DLQ |
| 4 | search + BFF aggregation |
| 5 | Next.js MVP |
| 6 | tracing, load/chaos, опционально K8s |

## Принципы

- Sync для команд пользователя (REST), async для побочных эффектов (Kafka)
- Partition key = `ticket_id` — порядок событий одного тикета
- Consumers идемпотентны по `event_id`
- Нет shared DB между сервисами

## Лицензия

TBD
