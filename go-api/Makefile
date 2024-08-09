ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	@go build -o bin/app ./cmd/app

run: build
	@./bin/app

test:
	@go test -v ./...

migrate_up:
	goose -dir db/migrations postgres "postgresql://postgres:$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" up

migrate_down:
	goose -dir db/migrations postgres "postgresql://postgres:$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" down




