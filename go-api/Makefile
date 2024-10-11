ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	@go build -o bin/app ./cmd/app

genapi:
	npx @redocly/openapi-cli@latest bundle spec/http-server/main.yaml > spec/http-server/generated.yaml
	oapi-codegen --config=spec/http-server/oapi-config.yaml spec/http-server/generated.yaml

geniss:
	oapi-codegen --config=spec/iss-client/oapi-config.yaml spec/iss-client/main.yaml

lint:
	golangci-lint run

run: build
	@./bin/app

test:
	@go test -v ./...

migrate_up:
	goose -dir db/migrations postgres "postgresql://$(PG_USERNAME):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" up

migrate_down:
	goose -dir db/migrations postgres "postgresql://$(PG_USERNAME):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" down

migrate_down_all:
	goose -dir db/migrations postgres "postgresql://$(PG_USERNAME):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" down-to 0

migrate_reset:
	goose -dir db/migrations postgres "postgresql://$(PG_USERNAME):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable" reset

proto:
	protoc --go_out=. --go-grpc_out=. proto/investor.proto

.PHONY: proto
