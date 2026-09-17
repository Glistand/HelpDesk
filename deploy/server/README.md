# Standalone server pack (2 files)

На сервер копируете **только**:

- `docker-compose.yml`
- `.env`

```bash
mkdir -p ~/helpdesk && cd ~/helpdesk
# scp оба файла сюда

# если GHCR private:
echo "$GITHUB_TOKEN" | docker login ghcr.io -u YOUR_USER --password-stdin

docker compose --env-file .env up -d
docker compose ps
curl -s http://127.0.0.1:8080/healthz
```

- UI: `http://SERVER/`
- API: `http://SERVER:8080`
- Логин: `agent@helpdesk.local` / `password`

Init Postgres/NATS встроен в compose (отдельные one-shot контейнеры) — репозиторий и папка `init/` не нужны.
