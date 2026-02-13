APP_NAME := git-sync-status
CMD_PATH := ./cmd/git-sync-status

.PHONY: build run test lint vet race tidy

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

race:
	go test -race ./...

tidy:
	go mod tidy
