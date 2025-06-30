.PHONY: build run dev

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

dev:
	docker compose -f docker/docker-compose.yml up -d db && make run

down:
	docker compose -f docker/docker-compose.yml down