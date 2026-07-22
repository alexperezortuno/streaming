APP_DIR = ./app
BINARY  = ../bin/streaming

.PHONY: help build run test test-race test-v vet lint tidy clean \
        dev up down logs ps graph all

help:           ## List all targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build:          ## Compile server binary
	cd $(APP_DIR) && go build -o $(BINARY) ./cmd/server

run:            ## Start the server (requires PostgreSQL)
	cd $(APP_DIR) && go run ./cmd/server

test:           ## Run all tests
	cd $(APP_DIR) && go test ./...

test-race:      ## Run tests with race detector
	cd $(APP_DIR) && go test -race ./...

test-v:         ## Run tests verbosely
	cd $(APP_DIR) && go test -v ./...

vet:            ## Run static analysis
	cd $(APP_DIR) && go vet ./...

lint:           ## Alias for vet
	cd $(APP_DIR) && go vet ./...

tidy:           ## Tidy Go module dependencies
	cd $(APP_DIR) && go mod tidy

clean:          ## Remove build artifacts
	rm -f $(APP_DIR)/$(BINARY)
	cd $(APP_DIR) && go clean

dev:            ## Start all services with Docker (foreground, rebuild)
	docker compose up --build

up:             ## Start all services in background
	docker compose up -d --build

down:           ## Stop all services
	docker compose down

logs:           ## Tail logs from all services
	docker compose logs -f

ps:             ## Show service status
	docker compose ps

graph:          ## Rebuild knowledge graph
	graphify .

all: vet test build  ## Run full pipeline (vet -> test -> build)
