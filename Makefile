include .env
export

BINARY_NAME   := server
MAIN_PATH     := .server
MIGRATION_DIR := internal/db/migration
SQLC_CONFIG   := sqlc.yaml

DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: nuke
nuke:
	docker compose down -v

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: ps
ps:
	docker compose ps

.PHONY: psql
psql:
	docker compose exec database psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_DSN)" up

.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$(DB_DSN)" down 1

.PHONY: migrate-down-all
migrate-down-all:
	migrate -path $(MIGRATION_DIR) -database "$(DB_DSN)" down -all

.PHONY: migrate-version
migrate-version:
	migrate -path $(MIGRATION_DIR) -database "$(DB_DSN)" version

.PHONY: migrate-force
migrate-force:
	migrate -path $(MIGRATION_DIR) -database "$(DB_DSN)" force $(v)

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Ошибка: укажи name=... например: make migrate-create name=add_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)
