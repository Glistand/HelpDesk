.PHONY: up down logs ps reset nats-streams env web web-build proto build-libs build-services test-go

COMPOSE_FILE := deploy/compose/docker-compose.yml
ENV_FILE := .env

up: env
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d --build

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

proto:
	cd api/proto && buf generate

build-libs:
	go build github.com/Glistand/HelpDesk/libs/eventkit/... \
		github.com/Glistand/HelpDesk/libs/grpckit/... \
		github.com/Glistand/HelpDesk/api/gen/go/...

build-services: build-libs
	go build -o bin/auth-service ./services/auth-service/cmd/auth-service
	go build -o bin/ticket-service ./services/ticket-service/cmd/ticket-service
	go build -o bin/api-gateway ./services/api-gateway/cmd/api-gateway

test-go:
	go test github.com/Glistand/HelpDesk/libs/eventkit/... \
		github.com/Glistand/HelpDesk/libs/grpckit/... \
		github.com/Glistand/HelpDesk/services/ticket-service/... \
		github.com/Glistand/HelpDesk/services/auth-service/... \
		github.com/Glistand/HelpDesk/services/api-gateway/...
