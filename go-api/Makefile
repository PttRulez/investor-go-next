ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	@go build -o bin/app ./cmd/app

genapi:
	npx @redocly/openapi-cli@latest bundle spec/api/main.yaml > spec/api/api.yaml
	oapi-codegen --config=spec/api/oapi-config.yaml spec/api/api.yaml

geniss:
	oapi-codegen --config=spec/iss/oapi-config.yaml spec/iss/main.yaml

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
