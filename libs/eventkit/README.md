# eventkit

Общая библиотека для работы с NATS JetStream в helpdesk-сервисах.

- envelope событий (`event_id`, `correlation_id`, `occurred_at`, payload)
- publish / subscribe helpers
- идемпотентный consumer wrapper
- health checks

Используется всеми сервисами monorepo [HelpDesk](https://github.com/Glistand/HelpDesk).
