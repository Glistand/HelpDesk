.PHONY: up down logs ps reset nats-streams env web web-build proto build-libs build-services test-go e2e e2e-phase2 e2e-phase3 e2e-phase4

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
	go build -o bin/assignment-service ./services/assignment-service/cmd/assignment-service
	go build -o bin/notification-service ./services/notification-service/cmd/notification-service
	go build -o bin/audit-service ./services/audit-service/cmd/audit-service
	go build -o bin/sla-service ./services/sla-service/cmd/sla-service
	go build -o bin/escalation-service ./services/escalation-service/cmd/escalation-service
	go build -o bin/search-service ./services/search-service/cmd/search-service
	go build -o bin/api-gateway ./services/api-gateway/cmd/api-gateway

test-go:
	go test github.com/Glistand/HelpDesk/libs/eventkit/... \
		github.com/Glistand/HelpDesk/libs/grpckit/... \
		github.com/Glistand/HelpDesk/services/ticket-service/... \
		github.com/Glistand/HelpDesk/services/auth-service/... \
		github.com/Glistand/HelpDesk/services/assignment-service/... \
		github.com/Glistand/HelpDesk/services/notification-service/... \
		github.com/Glistand/HelpDesk/services/audit-service/... \
		github.com/Glistand/HelpDesk/services/sla-service/... \
		github.com/Glistand/HelpDesk/services/escalation-service/... \
		github.com/Glistand/HelpDesk/services/search-service/... \
		github.com/Glistand/HelpDesk/services/api-gateway/...

e2e-phase2:
	./scripts/e2e-phase2.sh

e2e-phase3:
	./scripts/e2e-phase3.sh

e2e-phase4:
	./scripts/e2e-phase4.sh

e2e: e2e-phase2 e2e-phase3 e2e-phase4
