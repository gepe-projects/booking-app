# Makefile untuk Database Migrations

# Database configuration
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= booking_app
DB_USER ?= root
DB_PASS ?= root
DB_SSL_MODE ?= disable

# Migration directory
MIGRATION_DIR = pkg/database/migrations

# Construct database URL
DB_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

# Include environment variables if .env exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: migrate-create migrate-up migrate-down migrate-force migrate-version migrate-drop

# Create new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

# Apply all pending migrations
migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

# Apply specific number of migrations
migrate-up-steps:
	@if [ -z "$(steps)" ]; then \
		echo "Usage: make migrate-up-steps steps=<number>"; \
		exit 1; \
	fi
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up $(steps)

# Rollback last migration
migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1

# Rollback specific number of migrations
migrate-down-steps:
	@if [ -z "$(steps)" ]; then \
		echo "Usage: make migrate-down-steps steps=<number>"; \
		exit 1; \
	fi
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down $(steps)

# Check current migration version
migrate-version:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version

# Help command
migrate-help:
	@echo "Database Migration Commands:"
	@echo "  make migrate-create name=<name>      Create new migration"
	@echo "  make migrate-up                      Apply all pending migrations"
	@echo "  make migrate-up-steps steps=<n>      Apply specific number of migrations"
	@echo "  make migrate-down                    Rollback last migration"
	@echo "  make migrate-down-steps steps=<n>    Rollback specific number of migrations"
	@echo "  make migrate-version                 Check current version"