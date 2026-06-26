#!/usr/bin/env bash
# Manual WorkSchedule domain verification — real go test + coverage on disk.
#
# Usage: ./scripts/verify-work-schedule.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/domain/workforce/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== WorkSchedule & time accounting verification ==="
echo "Package: $PKG"
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

echo
echo "--- go test -v -run 'BR03|Worked|Lateness|Overtime|Evaluate|TimeBank|WorkSchedule' $PKG ---"
if (cd "$ATTENDANCE" && go test -v -run 'BR03|Worked|Lateness|Overtime|Evaluate|TimeBank|WorkSchedule' $PKG); then
  pass "go test work schedule"
else
  fail "go test work schedule"
fi

echo
echo "--- go test -cover $PKG ---"
COVER_OUTPUT="$(cd "$ATTENDANCE" && go test -cover $PKG 2>&1)"
echo "$COVER_OUTPUT"
if echo "$COVER_OUTPUT" | grep -q "ok"; then
  pass "coverage run"
else
  fail "coverage run"
fi

echo
echo "--- go vet $PKG ---"
if (cd "$ATTENDANCE" && go vet $PKG); then
  pass "go vet workforce schedule"
else
  fail "go vet workforce schedule"
fi

echo
echo "--- layer isolation ---"
if ! grep -rE 'internal/(infrastructure|interfaces|application)' "$ATTENDANCE/internal/domain/workforce" 2>/dev/null; then
  pass "workforce package has no outer layer imports"
else
  fail "workforce imports outer layers"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== WorkSchedule verification: ALL PASSED ==="
  exit 0
else
  echo "=== WorkSchedule verification: FAILED ==="
  exit 1
fi
