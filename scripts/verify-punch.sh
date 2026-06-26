#!/usr/bin/env bash
# Manual PunchRecord domain verification — real go test + coverage on disk.
#
# Usage: ./scripts/verify-punch.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/domain/punch/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== PunchRecord validation verification ==="
echo "Package: $PKG"
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

echo
echo "--- go test -v $PKG ---"
if (cd "$ATTENDANCE" && go test -v $PKG); then
  pass "go test punch domain"
else
  fail "go test punch domain"
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
  pass "go vet punch"
else
  fail "go vet punch"
fi

echo
echo "--- layer isolation ---"
if ! grep -rE 'internal/(infrastructure|interfaces|application)' "$ATTENDANCE/internal/domain/punch" 2>/dev/null; then
  pass "punch package has no outer layer imports"
else
  fail "punch imports outer layers"
fi

echo
echo "--- geofence regression ---"
if "$ROOT/scripts/verify-geofence.sh"; then
  pass "geofence regression"
else
  fail "geofence regression"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Punch verification: ALL PASSED ==="
  exit 0
else
  echo "=== Punch verification: FAILED ==="
  exit 1
fi
