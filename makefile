# mrb-service/makefile

ifneq (,$(wildcard .env))
	include .env
	export
endif

.DEFAULT_GOAL := help

# -------------------------------------------------------------------------
# PHONY

.PHONY: help

help:
	@echo "meeting-room-booking-service"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Generation:"
	@echo "  api-v1gen            Generate Go code from OpenAPI specification"


api-v1gen:
	go tool oapi-codegen -config ./api/v1/configs/server.yml ./api/v1/api.yml
	go tool oapi-codegen -config ./api/v1/configs/models.yml ./api/v1/api.yml
	go tool oapi-codegen -config ./api/v1/configs/spec.yml ./api/v1/api.yml