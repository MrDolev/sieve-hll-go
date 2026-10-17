#!/usr/bin/env bash
#
# Emulates multiple clients hitting the HTTP API concurrently.
# Each "client" sends a different request (varying n) and a spoofed
# X-Forwarded-For header with a random IP taken from IP_CIDR, so you can
# check which classes of IPs the app filters (HTTP 204) and which ones it
# collects (HTTP 200).
#
# NOTE: X-Forwarded-For only changes what the app *sees/logs* if it
# trusts that header -- it does not change the real TCP source IP.
# For genuinely distinct source IPs, use docker-compose.clients.yml
# instead (real containers = real distinct IPs on the Docker network).
#
# Usage:
#   ./scripts/load_test.sh [BASE_URL] [CONCURRENCY] [N_MIN] [N_MAX]
#
# Env:
#   IP_CIDR   comma-separated IPv4 ranges the fake client IPs are drawn
#             from; each client picks one range at random (default 10.0.0.0/8)
#
# Examples:
#   ./scripts/load_test.sh
#   ./scripts/load_test.sh http://localhost:8081 20 10 500
#   IP_CIDR=192.168.1.0/24 ./scripts/load_test.sh
#   IP_CIDR=10.0.0.0/8,8.8.8.0/24 ./scripts/load_test.sh http://localhost:8081 50

set -euo pipefail

BASE_URL="${1:-http://localhost:8081}"
CONCURRENCY="${2:-10}"
N_MIN="${3:-10}"
N_MAX="${4:-200}"
IP_CIDR="${IP_CIDR:-10.0.0.0/8}"

for range in ${IP_CIDR//,/ }; do
  if ! [[ "$range" =~ ^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[12][0-9]|3[0-2]))?$ ]]; then
    echo "invalid IP_CIDR entry: ${range} (expected e.g. 192.168.1.0/24)" >&2
    exit 2
  fi
done

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT
STATUS_FILE="$WORK_DIR/status"
CODES_FILE="$WORK_DIR/codes"
: > "$CODES_FILE"
export STATUS_FILE CODES_FILE

echo "Target:      ${BASE_URL}"
echo "Concurrency: ${CONCURRENCY} simulated clients"
echo "n range:     ${N_MIN}-${N_MAX} (each client picks its own value)"
echo "IP ranges:   ${IP_CIDR}"
echo

# Prints a random IPv4 address inside the given CIDR.
random_ip() {
  local addr="${1%/*}" bits=32
  [[ "$1" == */* ]] && bits="${1#*/}"
  local IFS=.
  # shellcheck disable=SC2086
  set -- $addr
  local base=$(( ($1 << 24) | ($2 << 16) | ($3 << 8) | $4 ))
  local size=$(( 1 << (32 - bits) ))
  local network=$(( base & ~(size - 1) & 0xFFFFFFFF ))
  local rand=$(( (RANDOM << 17) | (RANDOM << 2) | (RANDOM & 3) ))
  local ip=$(( network + rand % size ))
  echo "$(( (ip >> 24) & 255 )).$(( (ip >> 16) & 255 )).$(( (ip >> 8) & 255 )).$(( ip & 255 ))"
}

run_client() {
  local id="$1"
  local n=$(( (RANDOM % (N_MAX - N_MIN + 1)) + N_MIN ))
  local ranges
  read -r -a ranges <<< "${IP_CIDR//,/ }"
  local range="${ranges[RANDOM % ${#ranges[@]}]}"
  local fake_ip
  fake_ip="$(random_ip "$range")"

  local result code
  if ! result=$(curl -s -o /dev/null -w "%{http_code} %{time_total}" \
      -H "X-Forwarded-For: ${fake_ip}" \
      "${BASE_URL}?n=${n}"); then
    result="ERROR:$?"
    printf '1' > "$STATUS_FILE"
  fi
  code="${result%% *}"

  local outcome
  case "$code" in
    200) outcome="collected" ;;
    204) outcome="filtered" ;;
    *)   outcome="other" ;;
  esac
  echo "$code" >> "$CODES_FILE"

  printf "client %-3s ip=%-15s from=%-18s n=%-4s -> %-16s %s\n" \
    "$id" "$fake_ip" "$range" "$n" "$result" "$outcome"
}

export -f run_client random_ip
export BASE_URL N_MIN N_MAX IP_CIDR

seq 1 "${CONCURRENCY}" | xargs -P "${CONCURRENCY}" -I{} bash -c 'run_client "$@"' _ {}

echo
echo "Summary:"
printf "  collected (200): %s\n" "$(grep -c '^200$' "$CODES_FILE" || true)"
printf "  filtered  (204): %s\n" "$(grep -c '^204$' "$CODES_FILE" || true)"
printf "  other:           %s\n" "$(grep -vc -e '^200$' -e '^204$' "$CODES_FILE" || true)"

echo
echo "Health check:"
curl -s -o /dev/null -w "HTTP %{http_code}\n" "${BASE_URL}"

if [ -s "$STATUS_FILE" ]; then
  exit 1
fi
