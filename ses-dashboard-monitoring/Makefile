# SES Dashboard Monitoring - Backend Makefile

.PHONY: build run test clean docker-build docker-run swagger deps migrate-up migrate-down migrate-create load-env

# Load environment variables from root .env file (generated from config.yaml)
include ../.env
export

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=ses-monitoring
BINARY_UNIX=$(BINARY_NAME)_unix

# Database URL for migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/api

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/api
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Generate swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/api

# Docker build
docker-build:
	docker build -t ses-monitoring-backend .

# Docker run
docker-run:
	docker run -p 8080:8080 ses-monitoring-backend

# Development server with hot reload
dev:
	air -c .air.toml

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Security check
sec:
	gosec ./...

# All checks (format, lint, test)
check: fmt lint test

# Install migrate tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Generate .env from config.yaml
generate-env:
	../generate-env.sh

# Run migrations up
migrate-up:
	go run cmd/migrate/main.go -action=up -host=$(DB_HOST) -port=$(DB_PORT) -user=$(DB_USER) -password=$(DB_PASSWORD) -dbname=$(DB_NAME)

# Run migrations down
migrate-down:
	go run cmd/migrate/main.go -action=down -host=$(DB_HOST) -port=$(DB_PORT) -user=$(DB_USER) -password=$(DB_PASSWORD) -dbname=$(DB_NAME)

# Create new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/infrastructure/database/migration $$name

# Force migration version
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path internal/infrastructure/database/migration -database "$(DB_URL)" force $$version

# Check migration version
migrate-version:
	migrate -path internal/infrastructure/database/migration -database "$(DB_URL)" version

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  swagger       - Generate swagger documentation"
	@echo "  dev           - Run development server with hot reload"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  check         - Run all checks (format, lint, test)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo ""
	@echo "Migration commands:"
	@echo "  install-migrate - Install golang-migrate tool"
	@echo "  migrate-up      - Run migrations up"
	@echo "  migrate-down    - Run migrations down"
	@echo "  migrate-create  - Create new migration"
	@echo "  migrate-force   - Force migration version"
	@echo "  migrate-version - Check current migration version"