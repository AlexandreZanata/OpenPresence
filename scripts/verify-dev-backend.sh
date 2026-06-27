#!/usr/bin/env bash
# Manual dev backend verification — Postgres + biometric gRPC + attendance HTTP.
#
# Usage: ./scripts/verify-dev-backend.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE_URL="${ATTENDANCE_HTTP_URL:-http://127.0.0.1:8088}"
BIOMETRIC_PORT="${BIOMETRIC_GRPC_ADDR:-127.0.0.1:9090}"
BIOMETRIC_PORT="${BIOMETRIC_PORT##*:}"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

is_listening() {
  (echo >/dev/tcp/127.0.0.1/"$1") 2>/dev/null
}

echo "=== Dev backend verification ==="
echo

for cmd in docker go cargo curl; do
  if command -v "$cmd" >/dev/null 2>&1; then
    pass "$cmd available"
  else
    fail "$cmd missing"
  fi
done

echo
echo "--- ensure stack is up ---"
if ! curl -sf "${ATTENDANCE_URL}/health/live" >/dev/null 2>&1; then
  echo "Starting dev backend..."
  "${ROOT}/scripts/dev-backend.sh" start
else
  pass "attendance already responding"
fi

echo
echo "--- ./scripts/dev-backend.sh status ---"
STATUS_OUT="$("${ROOT}/scripts/dev-backend.sh" status 2>&1)" || true
echo "$STATUS_OUT"
echo "$STATUS_OUT" | grep -q "RUNNING: attendance" && pass "attendance RUNNING" || fail "attendance not RUNNING"
echo "$STATUS_OUT" | grep -q "RUNNING: biometric" && pass "biometric RUNNING" || fail "biometric not RUNNING"
echo "$STATUS_OUT" | grep -q "Up" && pass "postgres container up" || fail "postgres not up"

echo
echo "--- curl health/live ---"
RESP="$(curl -sf "${ATTENDANCE_URL}/health/live")"
if echo "$RESP" | grep -q '"status":"ok"'; then
  pass "GET /health/live → ok"
else
  fail "unexpected health response: $RESP"
fi

echo
echo "--- CORS preflight (admin dev origin) ---"
CORS_HEADERS="$(curl -sI -X OPTIONS "${ATTENDANCE_URL}/health/live" \
  -H "Origin: http://localhost:5174" \
  -H "Access-Control-Request-Method: GET" 2>/dev/null | tr -d '\r')"
if echo "$CORS_HEADERS" | grep -qi "access-control-allow-origin: http://localhost:5174"; then
  pass "OPTIONS /health/live allows admin origin"
else
  fail "CORS headers missing (restart backend: ./scripts/dev-backend.sh stop && start)"
fi

echo
echo "--- biometric gRPC port ${BIOMETRIC_PORT} ---"
if is_listening "$BIOMETRIC_PORT"; then
  pass "biometric port listening"
else
  fail "biometric port not listening"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Dev backend verification: ALL PASSED ==="
  exit 0
fi
echo "=== Dev backend verification: FAILED ==="
exit 1
