# auth-service

JWT login / validate. Exposes gRPC on `:50051`.

Dev seed:
- `agent@helpdesk.local` / `password` (agent)
- `admin@helpdesk.local` / `password` (admin)

```bash
JWT_SECRET=helpdesk-dev-secret-change-me go run ./cmd/auth-service
```
