#!/usr/bin/env bash
#
# Build the sieve-go Docker image.
#
# Usage:
#   ./build/build.sh source   # compile from source inside Docker (default, CI-safe)
#   ./build/build.sh local    # reuse local binary from `make bin` (fast iteration)
#
set -euo pipefail

MODE="${1:-source}"
IMAGE_NAME="${IMAGE_NAME:-sieve-go}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

# Run from repo root regardless of where the script is invoked from.
cd "$(dirname "$0")/.."

case "$MODE" in
  source)
    echo "==> Building image from source (build/Dockerfile)"
    docker build -f build/Dockerfile -t "${IMAGE_NAME}:${IMAGE_TAG}" .
    ;;
  local)
    echo "==> Compiling binary locally (make bin)"
    make bin
    echo "==> Building image from prebuilt binary (build/Dockerfile.local)"
    docker build -f build/Dockerfile.local -t "${IMAGE_NAME}:${IMAGE_TAG}" .
    ;;
  *)
    echo "Unknown mode: ${MODE} (expected 'source' or 'local')" >&2
    exit 1
    ;;
esac

echo "==> Built ${IMAGE_NAME}:${IMAGE_TAG}"
