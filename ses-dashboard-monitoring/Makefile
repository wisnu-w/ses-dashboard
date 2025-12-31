# ========================
# CONFIG
# ========================
APP_ENV ?= local
#CONFIG_FILE := config/config.$(APP_ENV).yaml
CONFIG_FILE := config/config.yaml
MIGRATION_PATH := internal/infrastructure/database/migration
SHELL := bash

# ========================
# DB CONFIG (yq REQUIRED)
# ========================

DB_URL := postgres://ses_user:password123!@localhost:5432/ses_dashboard?sslmode=disable

# ========================
# MIGRATION COMMANDS
# ========================

migrate-up:
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) up

migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) down 1

migrate-drop:
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) drop -f

migrate-force:
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) force $(version)

migrate-version:
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) version

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

.PHONY: migrate-up migrate-down migrate-drop migrate-force migrate-version migrate-create
