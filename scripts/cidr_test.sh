#!/usr/bin/env bash
#
# End-to-end check of the CIDR filter against a running instance.
#
# Sends one request per probe IP (spoofed via X-Forwarded-For) and compares
# the HTTP status with the one expected from the CIDR env var -- the same
# variable the app reads at startup:
#
#   204  ip inside at least one CIDR  -> filtered, not collected
#   200  ip outside every CIDR        -> collected
#   400  invalid ip                   -> rejected
#
# The expected result is computed here, independently from the Go code, so
# the app and the script must see the same CIDR value. `make cidr-test`
# loads it from .env.example/.env, exactly like `make local-run` and the
# docker compose stacks do. Restart the app after changing CIDR.
#
# Only IPv4 CIDR entries are used to compute expectations; IPv6 entries are
# reported and ignored (the default probes are all IPv4).
#
# Usage:
#   ./scripts/cidr_test.sh [BASE_URL]
#
# Env:
#   CIDR        comma-separated list, e.g. 10.0.0.0/8,192.168.0.0/16,8.8.8.8
#   PROBE_IPS   optional comma-separated list replacing the default probes
#
# Examples:
#   CIDR=10.0.0.0/8,172.16.0.0/12 ./scripts/cidr_test.sh http://localhost:8081
#   PROBE_IPS=10.1.2.3,8.8.8.8 make cidr-test

set -euo pipefail

BASE_URL="${1:-http://localhost:8081}"
CIDR="${CIDR:-}"

# "ip|group" -- boundaries are included so off-by-one mask errors show up.
DEFAULT_PROBES="
10.0.0.0|private class A (first)
10.42.7.1|private class A
10.255.255.255|private class A (last)
11.0.0.1|public, just after 10.0.0.0/8
172.15.255.255|public, just before 172.16.0.0/12
172.16.0.1|private class B (first)
172.20.10.5|private class B
172.31.255.255|private class B (last)
172.32.0.1|public, just after 172.16.0.0/12
192.167.255.255|public, just before 192.168.0.0/16
192.168.0.1|private class C
192.168.1.50|private class C
192.168.255.255|private class C (last)
192.169.0.1|public, just after 192.168.0.0/16
127.0.0.1|loopback
169.254.10.10|link-local
100.64.0.1|carrier-grade NAT
224.0.0.1|multicast
8.8.8.8|public
1.1.1.1|public
203.0.113.10|documentation (TEST-NET-3)
"

ip_to_int() {
  local IFS=.
  # shellcheck disable=SC2086
  set -- $1
  echo $(( ($1 << 24) | ($2 << 16) | ($3 << 8) | $4 ))
}

is_ipv4() {
  local re='^([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})$'
  [[ "$1" =~ $re ]] || return 1
  local octet
  for octet in "${BASH_REMATCH[@]:1}"; do
    (( 10#$octet <= 255 )) || return 1
  done
}

# Parsed IPv4 ranges as "network mask cidr" lines.
RANGES=""
IFS=',' read -r -a entries <<< "$CIDR"
for entry in ${entries[@]+"${entries[@]}"}; do
  entry="$(echo "$entry" | tr -d '[:space:]')"
  [ -z "$entry" ] && continue
  addr="${entry%/*}"
  bits=32
  [[ "$entry" == */* ]] && bits="${entry#*/}"
  if [[ "$addr" == *:* ]]; then
    echo "note: skipping IPv6 entry ${entry}"
    continue
  fi
  if ! is_ipv4 "$addr" || ! [[ "$bits" =~ ^[0-9]+$ ]] || (( bits > 32 )); then
    echo "invalid CIDR entry: ${entry} (the app would refuse to start)" >&2
    exit 2
  fi
  if (( bits == 0 )); then
    mask=0
  else
    mask=$(( (0xFFFFFFFF << (32 - bits)) & 0xFFFFFFFF ))
  fi
  network=$(( $(ip_to_int "$addr") & mask ))
  RANGES+="${network} ${mask} ${entry}"$'\n'
done

# Prints the matching CIDR entry, or nothing.
match() {
  local ip_int network mask entry
  ip_int=$(ip_to_int "$1")
  while read -r network mask entry; do
    [ -z "$network" ] && continue
    if (( (ip_int & mask) == network )); then
      echo "$entry"
      return
    fi
  done <<< "$RANGES"
}

if [ -n "${PROBE_IPS:-}" ]; then
  PROBES="$(echo "$PROBE_IPS" | tr ',' '\n' | sed 's/[[:space:]]//g; /^$/d; s/$/|custom/')"
else
  PROBES="$DEFAULT_PROBES"
fi
PROBES+=$'\n'"not-an-ip|invalid"$'\n'"999.1.1.1|invalid"

echo "Target: ${BASE_URL}"
echo "CIDR:   ${CIDR:-<empty, nothing filtered>}"
echo

if ! curl -s -o /dev/null --max-time 5 "${BASE_URL}"; then
  echo "cannot reach ${BASE_URL} -- is the app running?" >&2
  exit 2
fi

printf "%-6s %-17s %-38s %-8s %-8s %s\n" RESULT IP GROUP EXPECTED GOT MATCHED
failures=0
total=0
while IFS='|' read -r ip group; do
  [ -z "$ip" ] && continue
  total=$((total + 1))
  matched=""
  if ! is_ipv4 "$ip"; then
    expected=400
  else
    matched="$(match "$ip")"
    if [ -n "$matched" ]; then expected=204; else expected=200; fi
  fi

  got=$(curl -s -o /dev/null --max-time 5 -w "%{http_code}" \
    -H "X-Forwarded-For: ${ip}" "${BASE_URL}" || true)

  if [ "$got" = "$expected" ]; then
    result=ok
  else
    result=FAIL
    failures=$((failures + 1))
  fi
  printf "%-6s %-17s %-38s %-8s %-8s %s\n" "$result" "$ip" "$group" "$expected" "$got" "${matched:--}"
done <<< "$PROBES"

echo
echo "$((total - failures))/${total} passed"
[ "$failures" -eq 0 ]
