APP_NAME := carbon-emission-management

.PHONY: build run test clean docker-up docker-down deps lint fmt

build:
	go build -o ./bin/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./... -v

clean:
	rm -rf ./bin

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

deps:
	go mod tidy

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .
