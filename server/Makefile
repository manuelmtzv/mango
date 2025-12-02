.PHONY: run seed build test docker-up docker-down

run:
	go run ./cmd/api .

seed:
	go run ./cmd/seed

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
