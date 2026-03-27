.PHONY: build run clean test dev

VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")

build:
	go build -ldflags "-X clash-server/pkg/version.Version=$(VERSION)" -o bin/clash-server ./cmd/server/main.go

dev:
	go run ./cmd/server/main.go

clean:
	rm -f bin

test:
	go test ./...

.DEFAULT_GOAL := build
