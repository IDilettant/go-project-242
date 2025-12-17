test:
	go test -v ./... -race

lint:
	golangci-lint run ./...

build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size

.PHONY: test lint build
