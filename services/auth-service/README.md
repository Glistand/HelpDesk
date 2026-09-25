# auth-service

JWT login / validate. Exposes gRPC on `:50051`.

Dev seed (password for all: `password`):

| Email | Role | ID | Name |
|-------|------|-----|------|
| `agent@helpdesk.local` | agent (L1) | a-1 | Алексей К. |
| `admin@helpdesk.local` | admin | a-admin | Admin |

```bash
JWT_SECRET=helpdesk-dev-secret-change-me go run ./cmd/auth-service
```
