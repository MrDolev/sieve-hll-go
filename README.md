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
- `.env.example` provides baseline runtime values, but Docker Compose also applies service-specific overrides for container networking.
- When running via Docker Compose, `docker-compose.yml` and `docker-compose.stage.yaml` load `.env.example` and `.env` via `env_file`, and the app service additionally sets `REDIS_ADDR=redis:6379` so it can resolve the Redis service by container name.
- Edit the `.env` used by the compose command or set environment variables in your CI/deployment pipeline.

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

Access the Redis container CLI with:

```bash
make local-redis-cli
```

Run a simulated local load test against the app on port 8081:

```bash
make local-load-test
```

For Docker Compose, access the Redis container CLI with:

```bash
make compose-redis-cli
```

For the stage stack, use:

```bash
make stage-redis-cli
```

If you want to hit the server directly with curl:

```bash
curl -v "http://localhost:8081/?n=123"
```

For the Docker Compose dev stack, use:

```bash
make compose-up
make compose-redis-cli
make compose-load-test
```

For the stage Docker Compose stack, use:

```bash
make stage-up
make stage-logs
make stage-redis-cli
make stage-load-test
```

The stage stack exposes the app on host port `8082`.

Stop the local Redis instance with:

```bash
make local-down
```

For a full containerized local stack:

```bash
make dev-up
```

Alternatively, use the dedicated Docker Compose workflow:

```bash
make compose-build
make compose-up
```

Watch logs with:

```bash
make compose-logs
```

Tear it down with:

```bash
make compose-down
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

### Setup & Configuration
- `make env` - bootstrap `.env` from `.env.example`
- `make tools` - install build/lint/security tools

### Build & Compile
- `make build` - compile the project (no binary output)
- `make bin` - build the executable into `bin/`
- `make run` - build and run the binary directly

### Testing & Validation
- `make test` - run Go unit tests
- `make test-race` - run Go race condition detectors
- `make lint` - run static analysis with `golangci-lint`
- `make security` - run `gosec` and `go vet`
- `make validate` - run lint, security, and unit tests together
- `make integration-test` - run integration tests against a transient Redis instance
- `make integration-test-container` - run integration tests inside Docker
- `make clean` - remove build artifacts

### Local Development (Host Go Toolchain)
- `make local-up` / `make local-down` - start/stop Redis for local development
- `make local-run` - build and run the app against local Redis
- `make local-test` - run unit tests locally
- `make local-logs` - tail Redis logs
- `make local-redis-cli` - run Redis CLI against the local Redis-only stack
- `make local-load-test` - run simulated load test against localhost:8081

### Docker Compose (Dev Stack)
- `make compose-build` - build the app image for Docker Compose deployment
- `make compose-up` - build and deploy the local Docker Compose stack (app + redis on port 8081)
- `make compose-down` - stop the local Docker Compose stack
- `make compose-logs` - tail app + redis logs for the local Compose stack
- `make compose-redis-cli` - run Redis CLI against the local Compose stack
- `make compose-load-test` - run the load test against the local Compose stack (localhost:8081)
- `make compose-deploy` - alias for `make compose-up`

### Docker Dev Stack (Alternative)
- `make dev-build` - build the app image for dev (from source)
- `make dev-up` / `make dev-down` - start/stop the full Docker-based dev stack
- `make dev-logs` - tail dev stack logs

### Docker Stage Stack (Isolated)
- `make stage-build` - build the app image for stage
- `make stage-up` / `make stage-down` - start/stop the stage stack (app + redis on port 8082)
- `make stage-logs` - tail logs for the stage stack
- `make stage-redis-cli` - run Redis CLI against the stage stack
- `make stage-load-test` - run the load test against the stage stack (localhost:8082)

### Multi-Client Simulation
- `make load-test` - curl-based load test against localhost:8081 (works with any running instance)
- `make clients-up` / `make clients-down` - multi-container client simulation (requires dev-up first)

For a complete list of all targets and their descriptions, see the `Makefile` comments.

## Development & AI Assistance

Parts of this project have been developed with the assistance of Artificial Intelligence (AI) tools.
