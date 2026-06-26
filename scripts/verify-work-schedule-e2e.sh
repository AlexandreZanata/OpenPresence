#!/usr/bin/env bash
# Manual work schedule E2E verification — domain + CalculateDayAttendance integration (BR-030–034).
#
# Usage: ./scripts/verify-work-schedule-e2e.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/application/attendance/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Work schedule E2E verification (BR-030–034) ==="
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

if ! command -v docker >/dev/null 2>&1 || ! docker info >/dev/null 2>&1; then
  fail "docker required for CalculateDayAttendance integration tests"
  exit 1
fi
pass "docker available"

echo
echo "--- domain: verify-work-schedule.sh ---"
if "$ROOT/scripts/verify-work-schedule.sh"; then
  pass "verify-work-schedule.sh"
else
  fail "verify-work-schedule.sh"
fi

echo
echo "--- go test -tags=integration $PKG -run E2E ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PKG" -run E2E); then
  pass "CalculateDayAttendance E2E integration"
else
  fail "CalculateDayAttendance E2E integration"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Work schedule E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== Work schedule E2E verification: FAILED ==="
  exit 1
fi
