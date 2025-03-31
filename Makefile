APP_NAME = server
dev:
	swag init -g cmd/server/main.go && go run ./cmd/${APP_NAME}

run:
	docker compose up -d && go run ./cmd/${APP_NAME}

build:
	go build -o $(APP_NAME) cmd/${APP_NAME}/main.go

