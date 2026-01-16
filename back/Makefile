.PHONY: help
help: ## Show this help message
	@echo "LauraTech Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ========================================
# Development
# ========================================

.PHONY: dev
dev: ## Run development server with hot reload (requires air)
	air

.PHONY: run
run: ## Run application
	go run cmd/api/main.go

.PHONY: build
build: ## Build application binary
	go build -o bin/api cmd/api/main.go

# ========================================
# Dependencies
# ========================================

.PHONY: deps
deps: ## Install Go dependencies
	go mod download
	go mod tidy

.PHONY: install-tools
install-tools: ## Install development tools (migrate, sqlc, air)
	@echo "Installing golang-migrate..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing sqlc..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Installing air (hot reload)..."
	go install github.com/cosmtrek/air@latest
	@echo "✅ All tools installed"

# ========================================
# Database Migrations
# ========================================

.PHONY: migrate-up
migrate-up: ## Run all up migrations
	migrate -path db/migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Run all down migrations
	migrate -path db/migrations -database "$(DATABASE_URL)" down

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=add_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=add_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(name)

.PHONY: migrate-force
migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	migrate -path db/migrations -database "$(DATABASE_URL)" force $(version)

.PHONY: migrate-version
migrate-version: ## Show current migration version
	migrate -path db/migrations -database "$(DATABASE_URL)" version

# ========================================
# sqlc (Type-safe SQL)
# ========================================

.PHONY: sqlc
sqlc: ## Generate Go code from SQL queries
	cd db && sqlc generate

.PHONY: sqlc-verify
sqlc-verify: ## Verify sqlc configuration
	cd db && sqlc verify

# ========================================
# Testing
# ========================================

.PHONY: test
test: ## Run all tests
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	go test -v -race -short ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	go test -v -race -tags=integration ./tests/integration/...

.PHONY: test-coverage
test-coverage: test ## Show test coverage in browser
	go tool cover -html=coverage.out

# ========================================
# Code Quality
# ========================================

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint)
	golangci-lint run

# ========================================
# Docker
# ========================================

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t lauratech-backend:latest .

.PHONY: docker-up
docker-up: ## Start Docker Compose services
	docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop Docker Compose services
	docker-compose down

.PHONY: docker-logs
docker-logs: ## Show Docker Compose logs
	docker-compose logs -f

# ========================================
# Utilities
# ========================================

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

.PHONY: env
env: ## Create .env file from example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ .env file created from .env.example"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi
