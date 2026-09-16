# ticket-service

CRUD тикетов, transactional outbox → NATS JetStream. Exposes gRPC on `:50052`.

```bash
TICKET_DATABASE_URL=postgres://helpdesk:helpdesk@localhost:5432/ticket?sslmode=disable \
NATS_URL=nats://localhost:4222 \
go run ./cmd/ticket-service
```
