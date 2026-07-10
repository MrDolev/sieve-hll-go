#!/usr/bin/env bash
#
# Runs inside a client container. Fires REQUEST_COUNT requests at
# TARGET_URL/sieve with a random n in [N_MIN, N_MAX], waiting DELAY_MS
# between each. This container's own IP (visible to the app/Redis via
# Docker networking) is a real, distinct source IP -- no spoofing needed.

set -euo pipefail

TARGET_URL="${TARGET_URL:-http://app:8081}"
N_MIN="${N_MIN:-10}"
N_MAX="${N_MAX:-100}"
REQUEST_COUNT="${REQUEST_COUNT:-20}"
DELAY_MS="${DELAY_MS:-200}"

my_ip="$(hostname -i 2>/dev/null || echo unknown)"
echo "[$(hostname)] source ip=${my_ip} target=${TARGET_URL} requests=${REQUEST_COUNT}"

for i in $(seq 1 "${REQUEST_COUNT}"); do
  n=$(( (RANDOM % (N_MAX - N_MIN + 1)) + N_MIN ))
  code=$(curl -s -o /dev/null -w "%{http_code}" "${TARGET_URL}/sieve?n=${n}")
  echo "[$(hostname)] req=${i}/${REQUEST_COUNT} n=${n} -> HTTP ${code}"
  sleep "$(awk -v ms="${DELAY_MS}" 'BEGIN{print ms/1000}')"
done

echo "[$(hostname)] done"
