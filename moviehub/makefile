# =============================================================================
# GinMovieAPI — Makefile
# Usage: make <target>
# =============================================================================

.PHONY: help build up down restart logs ps clean \
        run-movie run-user run-gateway \
        test test-cover lint fmt \
        db-up db-shell-mongo db-shell-postgres db-shell-redis \
        seed docker-prune

# Default target
.DEFAULT_GOAL := help

# ── Variables ─────────────────────────────────────────────────────────────────
COMPOSE    := docker compose
GO         := go
GOFLAGS    := -v
BINARY_DIR := ./bin

# ── Help ──────────────────────────────────────────────────────────────────────
help: ## Show this help message
	@echo "GinMovieAPI — Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ── Docker Compose (Dev) ──────────────────────────────────────────────────────
build: ## Build all Docker images (dev)
	$(COMPOSE) build

up: ## Start all services (detached)
	$(COMPOSE) up -d

down: ## Stop and remove containers (keep volumes)
	$(COMPOSE) down

restart: down up ## Full restart of all services

# ── Docker Compose (Production) ───────────────────────────────────────────────
prod-build: ## Fast production build with BuildKit caching
	DOCKER_BUILDKIT=1 $(COMPOSE) -f docker-compose.prod.yml build

prod-up: ## Start production services in detached mode
	DOCKER_BUILDKIT=1 $(COMPOSE) -f docker-compose.prod.yml up -d --build

prod-down: ## Stop production containers
	$(COMPOSE) -f docker-compose.prod.yml down

prod-restart: prod-down prod-up ## Full restart of production services

prod-logs: ## Tail production logs
	$(COMPOSE) -f docker-compose.prod.yml logs -f

logs: ## Tail logs for all services
	$(COMPOSE) logs -f

logs-movie: ## Tail movie-service logs
	$(COMPOSE) logs -f movie-service

logs-user: ## Tail user-service logs
	$(COMPOSE) logs -f user-service

logs-gateway: ## Tail api-gateway logs
	$(COMPOSE) logs -f api-gateway

ps: ## Show running containers
	$(COMPOSE) ps

clean: ## Stop containers and remove volumes (destroys data)
	$(COMPOSE) down -v --remove-orphans

# ── Infrastructure only ───────────────────────────────────────────────────────
db-up: ## Start only infrastructure (mongo, postgres, redis)
	$(COMPOSE) up -d redisdb postgresdb mongodb

# ── Native Run (for development with hot-reload) ──────────────────────────────
run-movie: ## Run movie-service natively (requires db-up)
	$(GO) run $(GOFLAGS) ./cmd/movie

run-user: ## Run user-service natively (requires db-up)
	$(GO) run $(GOFLAGS) ./cmd/user

run-gateway: ## Run api-gateway natively
	$(GO) run $(GOFLAGS) ./cmd/gateway

# ── Hot Reload with Air ───────────────────────────────────────────────────────
air-movie: ## Hot reload movie-service with Air
	air -c .air.movie.toml

air-user: ## Hot reload user-service with Air
	air -c .air.user.toml

air-gateway: ## Hot reload api-gateway with Air
	air -c .air.gateway.toml

# ── Build Binaries ────────────────────────────────────────────────────────────
build-movie: ## Build movie-service binary
	mkdir -p $(BINARY_DIR)
	CGO_ENABLED=0 GOOS=linux $(GO) build -ldflags="-s -w" -o $(BINARY_DIR)/movie-service ./cmd/movie

build-user: ## Build user-service binary
	mkdir -p $(BINARY_DIR)
	CGO_ENABLED=0 GOOS=linux $(GO) build -ldflags="-s -w" -o $(BINARY_DIR)/user-service ./cmd/user

build-gateway: ## Build api-gateway binary
	mkdir -p $(BINARY_DIR)
	CGO_ENABLED=0 GOOS=linux $(GO) build -ldflags="-s -w" -o $(BINARY_DIR)/api-gateway ./cmd/gateway

build-all: build-movie build-user build-gateway ## Build all binaries

# ── Testing ───────────────────────────────────────────────────────────────────
test: ## Run all tests
	$(GO) test $(GOFLAGS) ./...

test-cover: ## Run tests with coverage report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# ── Code Quality ──────────────────────────────────────────────────────────────
lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format all Go files
	$(GO) fmt ./...
	goimports -w .

vet: ## Run go vet
	$(GO) vet ./...

# ── Database Shells ───────────────────────────────────────────────────────────
db-shell-mongo: ## Open MongoDB shell
	$(COMPOSE) exec mongodb mongosh -u admin -p password123 --authenticationDatabase admin

db-shell-postgres: ## Open PostgreSQL shell
	$(COMPOSE) exec postgresdb psql -U app_user -d app_db

db-shell-redis: ## Open Redis CLI
	$(COMPOSE) exec redisdb redis-cli

# ── Seeding ───────────────────────────────────────────────────────────────────
seed: ## Run movie seed script
	bash ./seed_movies.sh

# ── Maintenance ───────────────────────────────────────────────────────────────
tidy: ## Tidy Go modules
	$(GO) mod tidy
	$(GO) mod verify

docker-prune: ## Remove unused Docker resources
	docker system prune -f
	docker volume prune -f

# ── Health checks ─────────────────────────────────────────────────────────────
health: ## Check health of all services
	@echo "==> api-gateway"
	@curl -sf http://localhost:8080/health && echo "  OK" || echo "  FAIL"
	@echo "==> movie-service"
	@curl -sf http://localhost:8082/health && echo "  OK" || echo "  FAIL"
	@echo "==> user-service"
	@curl -sf http://localhost:8083/health && echo "  OK" || echo "  FAIL"
	@echo "==> Redis"
	@$(COMPOSE) exec -T redisdb redis-cli ping | grep -q PONG && echo "  OK" || echo "  FAIL"
	@echo "==> PostgreSQL"
	@$(COMPOSE) exec -T postgresdb pg_isready -U app_user -d app_db && echo "  OK" || echo "  FAIL"
