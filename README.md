# sieve-go

High-efficiency Golang & Redis backend for Sieve real-time chat IP tracking.

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/MrDolev/sieve-go.git
   cd sieve-go
   ```
2. Create an environment file from the example if needed:
   ```bash
   make env
   ```
   This copies `.env.example` to `.env` without overwriting an existing `.env`.
3. Edit `.env` only when you need to override the default Redis or server settings.

## Installation

The project is a standard Go module. To build the binary locally:

```bash
make bin
```

This creates the app binary at `bin/sieve-go`.

If you prefer to install dependencies and tools first, run:

```bash
make tools
```

## Running locally

For a fast local development workflow that uses the host Go toolchain with Redis:

```bash
make local-up
make local-run
```

Stop the local Redis instance with:

```bash
make local-down
```

For a full containerized local stack:

```bash
make dev-up
```

Tear it down with:

```bash
make dev-down
```

## Testing

Run unit tests:

```bash
make test
```

Run the full validation suite (lint, security, and unit tests):

```bash
make validate
```

Run integration tests with a temporary Redis container:

```bash
make integration-test
```

## Makefile reference

Common targets:

- `make env` - bootstrap `.env` from `.env.example`
- `make build` - compile the project (no binary output)
- `make bin` - build the executable into `bin/`
- `make test` - run Go unit tests
- `make test-race` - run Go race condition detectors
- `make lint` - run static analysis with `golangci-lint`
- `make security` - run `gosec` and `go vet`
- `make validate` - run lint, security, and unit tests together
- `make clean` - remove build artifacts
- `make local-up` / `make local-down` - start/stop Redis for local development
- `make local-run` - build and run the app against local Redis
- `make dev-up` / `make dev-down` - start/stop the full Docker-based dev stack
- `make integration-test` - run integration tests against a transient Redis instance

For more commands, inspect the `Makefile`.

## 🤖 Development & AI Assistance

Parts of this project have been developed with the assistance of Artificial Intelligence (AI) tools.
