.PHONY: build run down db-up

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

down:
	docker compose -f docker/docker-compose.yml down

db-up:
	docker compose -f docker/docker-compose.yml up -d db