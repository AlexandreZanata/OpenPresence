#!/usr/bin/env bash
# Manual geofence E2E verification — domain tests + SubmitPunch integration (BR-020–024).
#
# Usage: ./scripts/verify-geofence-e2e.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/application/punch/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Geofence E2E verification (BR-020–024) ==="
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

if ! command -v docker >/dev/null 2>&1 || ! docker info >/dev/null 2>&1; then
  fail "docker required for SubmitPunch integration tests"
  exit 1
fi
pass "docker available"

echo
echo "--- domain: verify-geofence.sh ---"
if "$ROOT/scripts/verify-geofence.sh"; then
  pass "verify-geofence.sh"
else
  fail "verify-geofence.sh"
fi

echo
echo "--- domain: verify-workforce-placement.sh ---"
if "$ROOT/scripts/verify-workforce-placement.sh"; then
  pass "verify-workforce-placement.sh"
else
  fail "verify-workforce-placement.sh"
fi

echo
echo "--- go test -tags=integration $PKG -run E2E_Geofence ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PKG" -run E2E_Geofence); then
  pass "SubmitPunch geofence E2E integration"
else
  fail "SubmitPunch geofence E2E integration"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Geofence E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== Geofence E2E verification: FAILED ==="
  exit 1
fi
