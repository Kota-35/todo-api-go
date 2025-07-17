.PHONY: build run down db-up migrate

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

down:
	docker compose -f docker/docker-compose.yml down

db-up:
	docker compose -f docker/docker-compose.yml up -d db

db-push:
	go run github.com/steebchen/prisma-client-go db push

migrate:
	go run github.com/steebchen/prisma-client-go migrate dev
