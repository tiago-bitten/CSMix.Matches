.PHONY: run build test lint tidy

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

lint:
	gofmt -l .
	go vet ./...

tidy:
	go mod tidy
