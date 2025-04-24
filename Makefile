
GOOSE_MIGRATION_DIR ?= sql/schema
GOOSE_DBSTRING = "postgres://dev:OEOOMYb3rpZv3xr2aUOzmC9b135@viaduct.proxy.rlwy.net:21434/fortune_vault"

APP_NAME = server
dev:
	swag init -g cmd/server/main.go && go run ./cmd/${APP_NAME}

run:
	docker compose up -d && go run ./cmd/${APP_NAME}

build:
	go build -o $(APP_NAME) cmd/${APP_NAME}/main.go

upse:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up
downse:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) down
resetse:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) reset
