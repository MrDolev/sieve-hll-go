# syntax=docker/dockerfile:1
#
# Build context must be the repo root, e.g.:
#   docker build -f build/Dockerfile -t sieve-hll-go:latest .
#
## ---- Build stage: compile the binary ----------------------------------
FROM golang:1.24.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/sieve-hll-go ./cmd/main.go

## ---- Runtime stage: minimal image, binary only ------------------------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /out/sieve-hll-go /app/sieve-hll-go

USER appuser

ENV HTTP_ADDR=:8080
EXPOSE 8080

ENTRYPOINT ["/app/sieve-hll-go"]
