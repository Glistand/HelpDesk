.PHONY: up down logs ps reset nats-streams env web web-build

COMPOSE_FILE := deploy/compose/docker-compose.yml
ENV_FILE := .env

up: env
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d

down:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down

logs:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) logs -f

ps:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) ps

reset:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down -v

nats-streams:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) run --rm nats-init

env:
	@test -f $(ENV_FILE) || cp .env.example $(ENV_FILE)

WEB_PORT ?= 3001

web:
	cd apps/web && npm run dev -- --port $(WEB_PORT)

web-build:
	cd apps/web && npm run build
