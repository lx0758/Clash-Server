.PHONY: build run clean test

build:
	cd web && $(MAKE) build
	cd server && $(MAKE) build

dev:
	cd server && $(MAKE) dev &

clean:
	cd server && $(MAKE) clean
	cd web && $(MAKE) clean

.DEFAULT_GOAL := build
