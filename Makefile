.PHONY: up down logs ps reset nats-streams env

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
