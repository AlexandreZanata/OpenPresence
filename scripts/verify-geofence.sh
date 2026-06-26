#!/usr/bin/env bash
# Manual geofence domain verification — real go test + coverage on disk.
#
# Usage: ./scripts/verify-geofence.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/domain/geofence/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Geofence engine verification ==="
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
  pass "go test geofence"
else
  fail "go test geofence"
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
  pass "go vet geofence"
else
  fail "go vet geofence"
fi

echo
echo "--- layer isolation (no outer imports in geofence) ---"
if ! grep -rE 'internal/(infrastructure|interfaces|application)' "$ATTENDANCE/internal/domain/geofence" 2>/dev/null; then
  pass "geofence package has no outer layer imports"
else
  fail "geofence imports outer layers"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Geofence verification: ALL PASSED ==="
  exit 0
else
  echo "=== Geofence verification: FAILED ==="
  exit 1
fi
