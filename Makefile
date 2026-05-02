# cooksense-backend — developer Makefile.
# Canonical UX for local dev. See specifications/SPEC-MAKE.md.

-include .env
export

.PHONY: help up down migrate seed run build test lint clean

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Start Postgres and wait until healthy.
	docker compose up -d postgres
	@docker compose wait postgres 2>/dev/null || \
	  ( i=0; \
	    until docker compose exec -T postgres pg_isready -q; do \
	      i=$$((i+2)); \
	      if [ $$i -ge 30 ]; then echo "postgres not healthy after 30s" >&2; exit 1; fi; \
	      sleep 2; \
	    done )

down: ## Stop Postgres (data volume preserved).
	docker compose down

migrate: up ## Apply database migrations.
	go run ./cmd/cooksense-server migrate

seed: up ## Load curated YAML data.
	go run ./cmd/cooksense-server seed

run: up migrate ## Start the server with modd hot-reload.
	modd

build: ## Build bin/cooksense-server.
	@mkdir -p bin
	go build -buildvcs=false -o bin/cooksense-server ./cmd/cooksense-server

test: ## Run go test ./....
	go test ./...

lint: ## Run go vet (and golangci-lint if available).
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
	  golangci-lint run ./...; \
	else \
	  echo "golangci-lint not on PATH, skipping (install: https://golangci-lint.run/welcome/install/)"; \
	fi

clean: ## Remove build artefacts. Set CLEAN_VOLUMES=1 to also drop Postgres data.
	rm -rf bin
ifeq ($(CLEAN_VOLUMES),1)
	docker compose down -v
endif
