# sieve-go

sieve-go is a lightweight, high-throughput Go service that receives client events from Sieve in real time. Requests are buffered in memory and processed in batches by a dedicated worker goroutine, which estimates unique IP cardinality using Redis HyperLogLog. The service exposes a simple HTTP API for querying aggregated results, with a focus on low latency and minimal Redis overhead.

## Project goal

The primary goal of this project is to provide a performant backend component that:

- Receives client events over HTTP and buffers them in a channel for asynchronous processing.
- Processes events in batches, reducing the number of writes to Redis.
- Estimates the number of unique IPs per day using Redis HyperLogLog (`PFADD`), trading exact counts for speed and low memory usage.
- Exposes a simple HTTP API for reading this aggregated data with low latency.

## Setup

1. Clone the repository:

```bash
   git clone https://github.com/MrDolev/sieve-go.git
   cd sieve-go
```

1. Create an environment file from the example if needed:

```bash
   make env
```

   This copies `.env.example` to `.env` without overwriting an existing `.env`.
3. Edit `.env` only when you need to override the default Redis or server settings.

## Configuring .env

The application reads runtime configuration from a `.env` file (created with `make env`). The following variables are available in `.env.example` and can be adjusted for your environment:

- `REDIS_ADDR` (default: `localhost:6379`) — Redis connection address (host:port).
- `REDIS_PASSWORD` (default: empty) — Password for Redis, if required.
- `REDIS_DB` (default: `0`) — Redis logical DB index to use.
- `REDIS_PROTOCOL` (default: `2`) — Protocol/client settings used by the Redis client (kept as provided).
- `SERVER_PORT` (default: `8081`) — TCP port the HTTP server listens on.
- `SERVER_PATH` (default: `/`) — Base HTTP path prefix for the server endpoints.

Example `.env` for a local run:

```env
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_PROTOCOL=2

SERVER_PORT=8081
SERVER_PATH=/
```

Notes:

- `make env` will not overwrite an existing `.env` file — it only bootstraps one from `.env.example` if needed.
- When running via Docker Compose the compose files supply the environment to containers via `env_file` entries; edit the `.env` used by the compose command or set environment variables in your CI/deployment pipeline.

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

## Development & AI Assistance

Parts of this project have been developed with the assistance of Artificial Intelligence (AI) tools.
