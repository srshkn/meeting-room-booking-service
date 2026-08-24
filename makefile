# mrb-service/makefile

ifneq (,$(wildcard .env))
	include .env
	export
endif

.DEFAULT_GOAL := help

# -------------------------------------------------------------------------
# CONSTANTS

# Environment
ENV_FILE=.env
ENV_EXAMPLE=.env.example

# API
V1API=api/v1/api.yml

# Migrations
MIGRATION_DB_HOST ?= localhost
MIGRATIONS_DIR := db/migrations
MIGRATION_DB_PORT ?= 5432
DATABASE_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(MIGRATION_DB_HOST):$(MIGRATION_DB_PORT)/$(DB_NAME)?sslmode=disable

# Compose
BASE=compose.yml
DEV=infra/compose/dev.yml

# JWT
PRIVATE_KEY=./secrets/private.pem
PUBLIC_KEY=./secrets/public.pem

# -------------------------------------------------------------------------
# PHONY

# Base
.PHONY: help

# Generation
.PHONY: env-gen v1api-gen jwt-keys-gen

# Migrations
.PHONY: migrate-create migrate-up migrate-down migrate-version

# Compose
.PHONY: compose-dev compose-clean

# -------------------------------------------------------------------------
# BASE

help:
	@echo "meeting-room-booking-service"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Generation:"
	@echo "  env-gen              Create .env from .env.example"
	@echo "  v1api-gen            Generate Go code from OpenAPI specification"
	@echo "  sql-gen              Generate Go code from SQL queries"
	@echo "  jwt-keys-gen              Generate JWT RSA keys"
	@echo ""
	@echo "Migrations:"
	@echo "  migrate-create       Create a new migration (name=<migration_name>)"
	@echo "  migrate-up           Apply all pending migrations"
	@echo "  migrate-down         Roll back the last migration"
	@echo "  migrate-version      Show current migration version"
	@echo ""
	@echo "Compose:"
	@echo "  compose-dev          Start development environment with Docker Compose"
	@echo "  compose-clean        Stop and remove development environment and volumes"

# -------------------------------------------------------------------------
# GENERATION

env-gen:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Creating .env from .env.example"; \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
	fi

v1api-gen:
	go tool oapi-codegen -config ./api/v1/configs/server.yml ./$(V1API)
	go tool oapi-codegen -config ./api/v1/configs/models.yml ./$(V1API)
	go tool oapi-codegen -config ./api/v1/configs/spec.yml ./$(V1API)

sql-gen:
	sqlc generate

jwt-keys-gen:
	@if [ -f $(PRIVATE_KEY) ] && [ -f $(PUBLIC_KEY) ]; then \
		echo "JWT keys already exist. Skipping generation."; \
	else \
		echo "Generating JWT RS256 keys..."; \
		mkdir -p secrets; \
		openssl genrsa -out $(PRIVATE_KEY) 2048; \
		openssl rsa -in $(PRIVATE_KEY) -pubout -out $(PUBLIC_KEY); \
		echo "JWT keys generated in ./secrets"; \
	fi

# -------------------------------------------------------------------------
# MIGRATIONS

migrate-create:
	@test -n "$(name)" || \
		(echo "Usage: make migrate-create name=add_user_status" && exit 1)
	migrate create \
		-ext sql \
		-dir $(MIGRATIONS_DIR) \
		-seq \
		$(name)

migrate-up:
	migrate \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		up

migrate-down:
	migrate \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		down 1

migrate-version:
	migrate \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		version

# -------------------------------------------------------------------------
# COMPOSE

compose-dev:
	docker compose --env-file ./$(ENV_FILE) \
		-f $(BASE) \
		-f $(DEV) \
		up --build

compose-clean:
	docker compose -f $(BASE) -f $(DEV) down -v --remove-orphans
	docker builder prune -af
