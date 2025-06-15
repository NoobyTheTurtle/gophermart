GOPHERMART_DIR = ./cmd/gophermart
ACCRUAL_DIR = ./cmd/accrual
GOPHERMART_BIN = $(GOPHERMART_DIR)/gophermart
ACCRUAL_BIN = $(ACCRUAL_DIR)/accrual_darwin_arm64

DATABASE_DSN ?= postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable

.DEFAULT_GOAL := help

.PHONY: generate
generate:
	@echo "Running go generate..."
	@go generate ./...

.PHONY: generate-mocks
generate-mocks:
	@echo "Regenerating all mocks..."
	@find ./internal -name "interfaces.go" | while read file; do \
		dir=$$(dirname "$$file"); \
		pkg=$$(basename "$$dir"); \
		echo "Generating mock for $$file -> $$dir/mocks.go"; \
		mockgen -source=$$file -destination=$$dir/mocks.go -package=$$pkg; \
	done
	@echo "Mocks successfully regenerated"

.PHONY: build-gophermart
build-gophermart:
	@echo "Building gophermart..."
	@go build -o $(GOPHERMART_BIN) $(GOPHERMART_DIR)
	@echo "Binary created at: $(GOPHERMART_BIN)"

.PHONY: run-gophermart
run-gophermart: build-gophermart
	@echo "Running gophermart..."
	@DATABASE_DSN="$(DATABASE_DSN)" $(GOPHERMART_BIN)

.PHONY: run-accrual
run-accrual:
	@echo "Running accrual system..."
	@$(ACCRUAL_BIN)

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying up modules..."
	@go mod tidy

.PHONY: mod-download
mod-download:
	@echo "Downloading modules..."
	@go mod download

.PHONY: postgres
postgres:
	@echo "Starting PostgreSQL in Docker..."
	@docker run --name gophermart-postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_DB=gophermart \
		-p 5432:5432 \
		-v $(shell pwd)/tmp/postgres-data:/var/lib/postgresql/data \
		-d \
		postgres:17-alpine

.PHONY: postgres-stop
postgres-stop:
	@echo "Stopping PostgreSQL Docker container..."
	@docker stop gophermart-postgres || true
	@docker rm gophermart-postgres || true

.PHONY: postgres-logs
postgres-logs:
	@echo "Showing PostgreSQL logs..."
	@docker logs -f gophermart-postgres

.PHONY: postgres-connect
postgres-connect:
	@echo "Connecting to PostgreSQL..."
	@docker exec -it gophermart-postgres psql -U postgres -d gophermart

.PHONY: migrate-up
migrate-up:
	@echo "Running database migrations up..."
	@goose -dir migrations postgres "$(DATABASE_DSN)" up

.PHONY: migrate-down
migrate-down:
	@echo "Running database migrations down..."
	@goose -dir migrations postgres "$(DATABASE_DSN)" down

.PHONY: migrate-status
migrate-status:
	@echo "Checking migration status..."
	@goose -dir migrations postgres "$(DATABASE_DSN)" status

.PHONY: migrate-create
migrate-create:
	@echo "Creating new migration (usage: make migrate-create NAME=your_migration_name)..."
	@if [ -z "$(NAME)" ]; then echo "Error: NAME parameter is required"; exit 1; fi
	@goose -dir migrations create $(NAME) sql

.PHONY: migrate-reset
migrate-reset:
	@echo "Resetting database..."
	@goose -dir migrations postgres "$(DATABASE_DSN)" reset

.PHONY: deps
deps:
	@echo "Installing development dependencies..."
	@go install go.uber.org/mock/mockgen@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: generate-swagger
generate-swagger:
	@echo "Generating swagger documentation..."
	@swag init -g cmd/gophermart/main.go -o docs --parseDependency --parseInternal

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make generate            - Run go generate"
	@echo "  make generate-mocks      - Regenerate all mocks"
	@echo "  make build-gophermart    - Build gophermart service"
	@echo "  make run-gophermart      - Run gophermart service"
	@echo "  make run-accrual         - Run accrual system"
	@echo "  make fmt                 - Format code"
	@echo "  make mod-tidy            - Tidy up modules"
	@echo "  make mod-download        - Download modules"
	@echo "  make postgres            - Start PostgreSQL in Docker"
	@echo "  make postgres-stop       - Stop and remove PostgreSQL Docker container"
	@echo "  make postgres-logs       - Show PostgreSQL logs"
	@echo "  make postgres-connect    - Connect to PostgreSQL"
	@echo "  make migrate-up          - Run database migrations up"
	@echo "  make migrate-down        - Run database migrations down"
	@echo "  make migrate-status      - Check migration status"
	@echo "  make migrate-create      - Create new migration (NAME=your_migration_name)"
	@echo "  make migrate-reset       - Reset database (down all + up all)"
	@echo "  make deps                - Install development dependencies"
	@echo "  make generate-swagger    - Generate swagger documentation"