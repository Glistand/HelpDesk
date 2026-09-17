.PHONY: up up-ghcr up-server down logs ps reset nats-streams env web web-build proto build-libs build-services test-go e2e e2e-phase2 e2e-phase3 e2e-phase4 e2e-phase5 e2e-phase6 load-phase6 chaos-phase6

COMPOSE_FILE := deploy/compose/docker-compose.yml
COMPOSE_GHCR_FILE := deploy/compose/docker-compose.ghcr.yml
COMPOSE_SERVER_FILE := deploy/server/docker-compose.yml
SERVER_ENV_FILE := deploy/server/.env
ENV_FILE := .env

up: env
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d --build

# Pull prebuilt images from ghcr.io (see docker-compose.ghcr.yml).
# IMAGE_TAG=sha-abc1234 GHCR_OWNER=glistand make up-ghcr
up-ghcr: env
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) -f $(COMPOSE_GHCR_FILE) pull
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) -f $(COMPOSE_GHCR_FILE) up -d --no-build

# Server-oriented stack: GHCR images, only WEB_PORT + GATEWAY_PORT published.
up-server:
	@test -f $(SERVER_ENV_FILE) || cp deploy/server/.env.example $(SERVER_ENV_FILE)
	docker compose --env-file $(SERVER_ENV_FILE) -f $(COMPOSE_SERVER_FILE) pull
	docker compose --env-file $(SERVER_ENV_FILE) -f $(COMPOSE_SERVER_FILE) up -d

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
		github.com/Glistand/HelpDesk/libs/otelkit/... \
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
		github.com/Glistand/HelpDesk/libs/otelkit/... \
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

e2e-phase5:
	./scripts/e2e-phase5.sh

e2e-phase6:
	./scripts/e2e-phase6.sh

load-phase6:
	./scripts/load-phase6.sh

chaos-phase6:
	./scripts/chaos-phase6.sh

e2e: e2e-phase2 e2e-phase3 e2e-phase4 e2e-phase5 e2e-phase6
