API_DIR := apps/api
COMPOSE_FILE := deploy/docker-compose.yml
APP_VERSION ?= $(shell tr -d '\r\n' < VERSION)

.PHONY: test build docker-build docker-up docker-down migrate-up migrate-down

test:
	cd $(API_DIR) && APP_VERSION=$(APP_VERSION) go test ./...

build:
	cd $(API_DIR) && APP_VERSION=$(APP_VERSION) go build ./...

docker-build:
	docker build -f apps/api/Dockerfile -t novascans-api:$(APP_VERSION) .

docker-up:
	APP_VERSION=$(APP_VERSION) docker compose -f $(COMPOSE_FILE) up -d --build

docker-down:
	docker compose -f $(COMPOSE_FILE) down --remove-orphans

migrate-up:
	migrate -path apps/api/migrations -database "$$DB_MAIN_DSN" up

migrate-down:
	migrate -path apps/api/migrations -database "$$DB_MAIN_DSN" down 1
