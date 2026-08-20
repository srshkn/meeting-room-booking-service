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

# Compose
BASE=compose.yml
DEV=infra/compose/dev.yml

# -------------------------------------------------------------------------
# PHONY

# Base
.PHONY: help

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
