# api-gateway

HTTP/REST BFF → gRPC (`auth-service`, `ticket-service`). Listens on `:8080`.

| Method | Path | Auth |
|--------|------|------|
| GET | `/healthz` | no |
| POST | `/auth/login` | no |
| GET/POST | `/tickets` | Bearer |
| GET | `/tickets/{id}` | Bearer |
| PATCH | `/tickets/{id}/status` | Bearer |

```bash
AUTH_GRPC_ADDR=localhost:50051 TICKET_GRPC_ADDR=localhost:50052 \
go run ./cmd/api-gateway
```
