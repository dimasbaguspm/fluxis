.PHONY: init dev build run down logs sqlc

init:
	go mod download
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


build:
	docker compose -f infra/docker-compose.yaml build

dev:
	docker compose -f infra/docker-compose.yaml up --build

down:
	docker compose -f infra/docker-compose.yaml down

logs:
	docker compose -f infra/docker-compose.yaml logs -f app

sqlc:
	sqlc generate