API_DIR := apps/api
COMPOSE_FILE := deploy/docker-compose.yml

.PHONY: test build docker-build docker-up docker-down migrate-up migrate-down

test:
	cd $(API_DIR) && go test ./...

build:
	cd $(API_DIR) && go build ./...

docker-build:
	docker build -f apps/api/Dockerfile -t novascans-api:local .

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d --build

docker-down:
	docker compose -f $(COMPOSE_FILE) down --remove-orphans

migrate-up:
	migrate -path apps/api/migrations -database "$$DB_MAIN_DSN" up

migrate-down:
	migrate -path apps/api/migrations -database "$$DB_MAIN_DSN" down 1
