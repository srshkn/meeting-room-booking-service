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
