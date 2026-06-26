#!/usr/bin/env bash
# UC-001 full E2E — POST /v1/attendance/punch → attendance → biometric gRPC → VALID.
#
# Usage:
#   ./scripts/verify-uc001-e2e.sh           # integration tests (Postgres + live biometric)
#   ./scripts/verify-uc001-e2e.sh --curl  # optional HTTP smoke via docker compose stack

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/interfaces/httpapi/..."
COMPOSE_FILE="$ROOT/infra/docker-compose.e2e.yml"
FAIL=0

TENANT_ID="11111111-1111-1111-1111-111111111111"
EMPLOYEE_ID="22222222-2222-2222-2222-222222222222"
BASE_URL="${ATTENDANCE_HTTP_URL:-http://127.0.0.1:8088}"

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== UC-001 full clock-in E2E verification ==="
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

if ! command -v cargo >/dev/null 2>&1; then
  fail "cargo not installed (biometric-server subprocess)"
  exit 1
fi
pass "cargo installed: $(cargo --version)"

if ! command -v docker >/dev/null 2>&1 || ! docker info >/dev/null 2>&1; then
  fail "docker required for Postgres integration"
  exit 1
fi
pass "docker available"

echo
echo "--- go test -tags=integration $PKG -run UC001 ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PKG" -run UC001); then
  pass "Punch API UC-001 integration tests"
else
  fail "Punch API UC-001 integration tests"
fi

if [[ "${1:-}" == "--curl" ]]; then
  echo
  echo "--- HTTP smoke: docker compose + curl POST punch ---"
  if ! command -v curl >/dev/null 2>&1 && ! command -v http >/dev/null 2>&1; then
    fail "curl or httpie required for --curl"
  else
    FIXTURE="$ROOT/services/biometric/tests/fixtures/valid_128.jpg"
    if [[ ! -f "$FIXTURE" ]]; then
      fail "fixture missing: $FIXTURE"
    else
      docker compose -f "$COMPOSE_FILE" up -d --wait postgres attendance biometric 2>/dev/null || \
        docker compose -f "$COMPOSE_FILE" up -d postgres attendance biometric

      for _ in $(seq 1 90); do
        if curl -sf "$BASE_URL/health/live" >/dev/null 2>&1; then
          break
        fi
        sleep 2
      done
      if ! curl -sf "$BASE_URL/health/live" >/dev/null 2>&1; then
        fail "attendance HTTP not ready at $BASE_URL"
      else
        pass "attendance HTTP ready"
        docker compose -f "$COMPOSE_FILE" exec -T postgres \
          psql -U openpresence -d openpresence -v ON_ERROR_STOP=1 \
          -f - < "$ROOT/infra/e2e/seed-uc001.sql" >/dev/null

        FRAME_B64="$(base64 -w0 "$FIXTURE")"
        BODY=$(cat <<EOF
{
  "punchType": "CLOCK_IN",
  "deviceTime": "2026-06-26T09:00:00Z",
  "location": {"latitude": -23.5505, "longitude": -46.6333, "accuracy": 10.0, "isMocked": false},
  "frameBase64": "$FRAME_B64",
  "deviceIntegrityReport": {"isRooted": false, "isVpnActive": false, "isDeveloperOptionsEnabled": false},
  "offlineSync": false
}
EOF
)
        TOKEN="e2e.${TENANT_ID}.${EMPLOYEE_ID}"
        if command -v http >/dev/null 2>&1; then
          RESP=$(http --check-status --ignore-stdin --print=b POST \
            "$BASE_URL/v1/attendance/punch" \
            "Authorization:Bearer $TOKEN" \
            Content-Type:application/json <<<"$BODY" 2>/dev/null || true)
        else
          RESP=$(curl -sf -X POST "$BASE_URL/v1/attendance/punch" \
            -H "Authorization: Bearer $TOKEN" \
            -H "Content-Type: application/json" \
            -d "$BODY")
        fi
        if echo "$RESP" | grep -q '"status":"VALID"'; then
          pass "curl POST punch returned VALID"
        else
          fail "curl POST punch unexpected response: $RESP"
        fi
      fi
      docker compose -f "$COMPOSE_FILE" down >/dev/null 2>&1 || true
    fi
  fi
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== UC-001 E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== UC-001 E2E verification: FAILED ==="
  exit 1
fi
