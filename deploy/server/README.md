# Server deploy (GHCR)

Минимальный стек для теста на VPS: образы из `ghcr.io`, наружу только UI и API.

## Порты

| Хост | Сервис | Зачем |
|------|--------|--------|
| `WEB_PORT` (80) | `web` | агентская консоль |
| `GATEWAY_PORT` (8080) | `api-gateway` | REST / e2e / curl |

Postgres, Redis, NATS, Meilisearch, gRPC, Jaeger/OTLP — **только** во внутренней сети Docker.

## Запуск

```bash
# на сервере: клон репо (нужны init-скрипты в deploy/compose/init)
git clone git@github.com:Glistand/HelpDesk.git && cd HelpDesk

cp deploy/server/.env.example deploy/server/.env
# отредактируйте JWT_SECRET, POSTGRES_PASSWORD, MEILI_MASTER_KEY

# если пакеты GHCR private:
echo "$GITHUB_TOKEN" | docker login ghcr.io -u YOUR_GITHUB_USER --password-stdin

docker compose --env-file deploy/server/.env -f deploy/server/docker-compose.yml pull
docker compose --env-file deploy/server/.env -f deploy/server/docker-compose.yml up -d
docker compose --env-file deploy/server/.env -f deploy/server/docker-compose.yml ps
```

Проверка:

```bash
curl -s http://127.0.0.1:8080/healthz
# UI: http://SERVER_IP/
# login: agent@helpdesk.local / password
```

Опционально с хоста (как локальный e2e):

```bash
GATEWAY_URL=http://127.0.0.1:8080 make e2e-phase4
WEB_URL=http://127.0.0.1 make e2e-phase5   # если WEB_PORT=80
```

Jaeger UI без публикации порта:

```bash
docker compose --env-file deploy/server/.env -f deploy/server/docker-compose.yml port jaeger 16686
# или временный SSH-туннель:
ssh -L 16686:127.0.0.1:16686 user@SERVER
# внутри: docker compose exec не нужен — проще:
docker run --rm --network helpdesk_default curlimages/curl -s http://jaeger:16686/
```

Остановка:

```bash
docker compose --env-file deploy/server/.env -f deploy/server/docker-compose.yml down
# + volumes:
docker compose --env-file deploy/server/.env -f deploy/server/docker-compose.yml down -v
```
