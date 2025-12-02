SHELL := /bin/bash

.PHONY: all build test deps deps-cleancache

GOCMD=go
BUILD_DIR=build
BINARY_DIR=$(BUILD_DIR)/bin
CODE_COVERAGE=code-coverage

all: test build

${BINARY_DIR}:
	mkdir -p $(BINARY_DIR)

local-build: ${BINARY_DIR}
	$(GOCMD) build -o $(BINARY_DIR) -v ./cmd/api

run:
	$(GOCMD) run ./cmd/api

test:
	$(GOCMD) test ./... -cover

test-coverage:
	$(GOCMD) test ./... -coverprofile=$(CODE_COVERAGE).out
	$(GOCMD) tool cover -html=$(CODE_COVERAGE).out

deps:
	$(GOCMD) get -u -t -d -v ./...
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor

cleancache:
	$(GOCMD) clean -modcache

wire:
	cd internal/di && wire

up:
	@docker-compose up -d

down:
	@docker-compose down -v

build:
	@docker-compose down -v && @docker-compose up -d --build

air:
	air
