# auth-service

JWT login / validate. Exposes gRPC on `:50051`.

Dev seed (password for all: `password`):

| Email | Role | ID | Name |
|-------|------|-----|------|
| `agent@helpdesk.local` | agent (L1) | a-1 | Алексей К. |
| `maria@helpdesk.local` | agent (L1) | a-2 | Марина С. |
| `ivan@helpdesk.local` | agent (L2) | a-3 | Денис В. |
| `admin@helpdesk.local` | admin | a-admin | Admin |
| `anna@company.local` | requester | r-1 | Анна П. |
| `peter@company.local` | requester | r-2 | Пётр В. |

```bash
JWT_SECRET=helpdesk-dev-secret-change-me go run ./cmd/auth-service
```
