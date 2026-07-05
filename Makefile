APP_NAME := sieve-go
BIN_DIR  := bin
CMD_DIR  := ./cmd
MAIN_PKG := $(CMD_DIR)/main.go

GO            := go
GOFLAGS       :=
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint
GOSEC         := $(shell go env GOPATH)/bin/gosec

.PHONY: all build bin test lint security validate clean run tools

all: validate build

## Compile the code (check it builds, no binary output)
build:
	$(GO) build $(GOFLAGS) -o /dev/null ./...

## Generate the binary into bin/
bin:
	mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PKG)

## Run tests
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

## Run lint, security, and tests together (use before build/release)
validate: lint security test

## Run the compiled binary
run: bin
	./$(BIN_DIR)/$(APP_NAME)

## Remove build artifacts
clean:
	rm -rf $(BIN_DIR)
