ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	@go build -o bin/app ./cmd/app

genapi:
	npx @redocly/openapi-cli@latest bundle internal/api/spec/main.yaml > internal/api/spec/api.yaml
	oapi-codegen --config=internal/api/spec/oapi-config.yaml internal/api/spec/api.yaml

geniss:
	oapi-codegen --config=internal/infrastructure/issclient/spec/oapi-config.yaml internal/infrastructure/issclient/spec/main.yaml

lint:
	golangci-lint run

run: build
	@./bin/app

test:
	@go test -v ./...

migrate_up:
	goose -dir db/migrations postgres "postgresql://postgres:$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" up

migrate_down:
	goose -dir db/migrations postgres "postgresql://postgres:$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" down

migrate_reset:
	goose -dir db/migrations postgres "postgresql://postgres:$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" reset
