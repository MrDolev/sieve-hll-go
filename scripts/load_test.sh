#!/usr/bin/env bash
#
# Emulates multiple clients hitting the HTTP API concurrently.
# Each "client" sends a different request (varying n) and, optionally,
# a spoofed X-Forwarded-For header so app-side logs show distinct IPs.
#
# NOTE: X-Forwarded-For only changes what the app *sees/logs* if it
# trusts that header -- it does not change the real TCP source IP.
# For genuinely distinct source IPs, use docker-compose.clients.yml
# instead (real containers = real distinct IPs on the Docker network).
#
# Usage:
#   ./scripts/load_test.sh [BASE_URL] [CONCURRENCY] [N_MIN] [N_MAX]
#
# Examples:
#   ./scripts/load_test.sh
#   ./scripts/load_test.sh http://localhost:8081 20 10 500

set -euo pipefail

BASE_URL="${1:-http://localhost:8081}"
CONCURRENCY="${2:-10}"
N_MIN="${3:-10}"
N_MAX="${4:-200}"

STATUS_FILE="$(mktemp)"
trap 'rm -f "$STATUS_FILE"' EXIT
export STATUS_FILE

echo "Target:      ${BASE_URL}"
echo "Concurrency: ${CONCURRENCY} simulated clients"
echo "n range:     ${N_MIN}-${N_MAX} (each client picks its own value)"
echo

run_client() {
  local id="$1"
  local n=$(( (RANDOM % (N_MAX - N_MIN + 1)) + N_MIN ))
  # Fake, distinct-looking source IP per client for X-Forwarded-For.
  local fake_ip
  fake_ip="10.$(( RANDOM % 256 )).$(( RANDOM % 256 )).$(( RANDOM % 256 ))"

  local result
  if ! result=$(curl -s -o /dev/null -w "%{http_code} %{time_total}" \
      -H "X-Forwarded-For: ${fake_ip}" \
      "${BASE_URL}?n=${n}"); then
    result="ERROR:$?"
    printf '1' > "$STATUS_FILE"
  fi

  printf "client %-3s ip=%-15s n=%-4s -> %s\n" "$id" "$fake_ip" "$n" "$result"
}

export -f run_client
export BASE_URL N_MIN N_MAX

seq 1 "${CONCURRENCY}" | xargs -P "${CONCURRENCY}" -I{} bash -c 'run_client "$@"' _ {}

echo
echo "Health check:"
curl -s -o /dev/null -w "HTTP %{http_code}\n" "${BASE_URL}"

if [ -s "$STATUS_FILE" ]; then
  exit 1
fi
