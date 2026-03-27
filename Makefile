.PHONY: build run clean test

VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")

build:
	cd web && VITE_APP_VERSION=$(VERSION) $(MAKE) build
	cd server && VERSION=$(VERSION) $(MAKE) build

dev:
	cd server && $(MAKE) dev &

clean:
	cd server && $(MAKE) clean
	cd web && $(MAKE) clean

.DEFAULT_GOAL := build
