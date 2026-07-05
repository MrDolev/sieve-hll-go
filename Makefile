APP_NAME := sieve-go
BIN_DIR  := bin
CMD_DIR  := ./cmd
MAIN_PKG := $(CMD_DIR)/main.go

GO            := go
GOFLAGS       :=
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint
GOSEC         := $(shell go env GOPATH)/bin/gosec

COMPOSE := docker compose

ENV_FILE    := .env
ENV_EXAMPLE := .env.example

# Load config into Make's own variables: .env.example first (defaults),
# then .env on top if it exists (real overrides win). Same precedence
# used by the Compose files' env_file layering below.
-include $(ENV_EXAMPLE)
-include $(ENV_FILE)
export

.PHONY: all build bin test lint security validate clean run tools env \
        local-up local-down local-logs local-run local-test \
        integration-up integration-down integration-test integration-test-container \
        dev-build dev-up dev-down dev-logs \
        stage-build stage-up stage-down stage-logs \
        clients-up clients-down load-test

all: validate build

## Create .env from .env.example if it doesn't exist yet (never overwrites)
env:
	@if [ ! -f $(ENV_FILE) ]; then \
		if [ ! -f $(ENV_EXAMPLE) ]; then \
			echo "$(ENV_EXAMPLE) not found, cannot bootstrap $(ENV_FILE)"; \
			exit 1; \
		fi; \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
		echo "Created $(ENV_FILE) from $(ENV_EXAMPLE) -- review/edit values before use."; \
	else \
		echo "$(ENV_FILE) already exists, leaving it untouched."; \
	fi

## ---- Build & test (no environment attached) ---------------------------

## Compile the code (check it builds, no binary output)
build:
	$(GO) build $(GOFLAGS) -o /dev/null ./...

## Generate the binary into bin/
bin:
	mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PKG)

## Run unit tests (no external dependencies required)
test:
	$(GO) test $(GOFLAGS) -v ./...

## Run static analysis / linting (golangci-lint)
lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "golangci-lint not found, installing..."; \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	$(GOLANGCI_LINT) run ./...

## Run security scan (gosec) + Go's built-in vulnerability check
security:
	@if [ ! -x "$(GOSEC)" ]; then \
		echo "gosec not found, installing..."; \
		$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	$(GOSEC) ./...
	$(GO) vet ./...

## Install lint/security tools upfront
tools:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest

## Run lint, security, and unit tests together (use before build/release)
validate: lint security test

## Remove build artifacts
clean:
	rm -rf $(BIN_DIR)

## Alias: build the binary and run it (uses whatever REDIS_ADDR is in .env/.env.example)
run: env bin
	./$(BIN_DIR)/$(APP_NAME)

## ================= LOCAL: Redis only, compile/run Go directly =================
## Use this while writing/testing code -- no app container, fastest loop.

LOCAL_COMPOSE := docker-compose.redis.yaml
LOCAL_PROJECT := sieve-local

## Start local Redis (reads .env.example, overlaid by .env if present)
local-up: env
	$(COMPOSE) -p $(LOCAL_PROJECT) -f $(LOCAL_COMPOSE) up -d

## Stop local Redis
local-down:
	$(COMPOSE) -p $(LOCAL_PROJECT) -f $(LOCAL_COMPOSE) down -v

## Tail local Redis logs
local-logs:
	$(COMPOSE) -p $(LOCAL_PROJECT) -f $(LOCAL_COMPOSE) logs -f

## Build the binary and run it against local Redis
local-run: local-up bin
	./$(BIN_DIR)/$(APP_NAME)

## Unit tests only, no containers required
local-test: test

## ================= INTEGRATION: Redis + go test -tags=integration =============

INTEGRATION_COMPOSE := docker-compose.test.yml
INTEGRATION_PROJECT := sieve-integration

## Start Redis for integration tests
integration-up:
	$(COMPOSE) -p $(INTEGRATION_PROJECT) -f $(INTEGRATION_COMPOSE) up -d redis-test

## Stop the integration Redis instance
integration-down:
	$(COMPOSE) -p $(INTEGRATION_PROJECT) -f $(INTEGRATION_COMPOSE) down -v

## Bring up Redis, run integration tests from the host against it, tear down
integration-test: integration-up
	@echo "Waiting for redis-test to be healthy..."
	@until [ "$$($(COMPOSE) -p $(INTEGRATION_PROJECT) -f $(INTEGRATION_COMPOSE) ps -q redis-test | xargs docker inspect -f '{{.State.Health.Status}}')" = "healthy" ]; do sleep 1; done
	REDIS_ADDR=localhost:6380 $(GO) test -tags=integration -v ./test/integration/... ; \
	status=$$?; \
	$(MAKE) integration-down; \
	exit $$status

## Run integration tests fully inside Docker (no local Go toolchain needed), then tear down
integration-test-container:
	$(COMPOSE) -p $(INTEGRATION_PROJECT) -f $(INTEGRATION_COMPOSE) run --rm test-runner ; \
	status=$$?; \
	$(COMPOSE) -p $(INTEGRATION_PROJECT) -f $(INTEGRATION_COMPOSE) down -v; \
	exit $$status

## ================= DEV: full local stack (app + redis) =========================

DEV_COMPOSE := docker-compose.yml
DEV_PROJECT := sieve-dev

## Build the app image for dev (compiles from source, uses build/Dockerfile)
dev-build:
	./build/build.sh source

## Start app + redis for local dev (reads .env.example, overlaid by .env)
dev-up: env
	$(COMPOSE) -p $(DEV_PROJECT) -f $(DEV_COMPOSE) up -d --build

## Stop the dev stack
dev-down:
	$(COMPOSE) -p $(DEV_PROJECT) -f $(DEV_COMPOSE) down -v

## Tail dev stack logs
dev-logs:
	$(COMPOSE) -p $(DEV_PROJECT) -f $(DEV_COMPOSE) logs -f

## ================= STAGE: isolated staging-like local run =======================
## Same image/services as dev, but separate ports/network/container names so
## it can run alongside dev-up without colliding. Adds an optional .env.stage
## layer on top of .env.example -> .env for stage-only overrides.

STAGE_COMPOSE := -f docker-compose.yml -f docker-compose.stage.yml
STAGE_PROJECT := sieve-stage

## Build the app image for stage (same Dockerfile as dev)
stage-build:
	./build/build.sh source

## Start app + redis for a staging-like local run
stage-up: env
	$(COMPOSE) -p $(STAGE_PROJECT) $(STAGE_COMPOSE) up -d --build

## Stop the stage stack
stage-down:
	$(COMPOSE) -p $(STAGE_PROJECT) $(STAGE_COMPOSE) down -v

## Tail stage stack logs
stage-logs:
	$(COMPOSE) -p $(STAGE_PROJECT) $(STAGE_COMPOSE) logs -f

## ---- Multi-client emulation (targets whatever stack is up, default dev) ------

## Quick host-based emulation: concurrent curl calls, spoofed X-Forwarded-For
load-test:
	./scripts/load_test.sh

## Real multi-container clients with distinct Docker network IPs
## (requires `make dev-up` first)
clients-up:
	$(COMPOSE) -f docker-compose.clients.yml up --build

clients-down:
	$(COMPOSE) -f docker-compose.clients.yml down -v
