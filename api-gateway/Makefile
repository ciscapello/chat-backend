include .env
export $(shell sed 's/=.*//' .env)

DB_CONNECTION=postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)"

MIGRATIONS_DIR=migrations

.PHONY: run
run:
	go run ./cmd/api-gateway/main.go

migrate-up:
	goose -dir $(MIGRATIONS_DIR) $(DB_CONNECTION) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) $(DB_CONNECTION) down
