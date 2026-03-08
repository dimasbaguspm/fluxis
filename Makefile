.PHONY: init dev build run down logs sqlc swagger apitest

init:
	go mod download
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

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

swagger:
	swag init --generalInfo cmd/fluxis/main.go --outputTypes json --output ./api

apitest:
	go test -v -count=1 -timeout=120s ./cmd/apitest/...