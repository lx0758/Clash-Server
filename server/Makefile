.PHONY: build run clean test dev

build:
	go build -o bin/clash-server ./cmd/server/main.go

dev:
	go run ./cmd/server/main.go

clean:
	rm -f bin

test:
	go test ./...

.DEFAULT_GOAL := build