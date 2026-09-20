# conversation-service

Support chat threads + OpenRouter bot replies. gRPC `:50058`.

```bash
CONVERSATION_DATABASE_URL=postgres://helpdesk:helpdesk@localhost:5432/conversation?sslmode=disable \
NATS_URL=nats://localhost:4222 \
OPENROUTER_API_KEY=sk-or-... \
OPENROUTER_MODEL=inclusionai/ling-3.0-flash-vl:free \
go run ./cmd/conversation-service
```
