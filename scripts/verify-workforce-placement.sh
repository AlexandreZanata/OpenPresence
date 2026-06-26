#!/usr/bin/env bash
# Manual workforce placement domain verification — real go test + coverage on disk.
#
# Usage: ./scripts/verify-workforce-placement.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/domain/workforce/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Workforce placement verification ==="
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
  pass "go test workforce placement"
else
  fail "go test workforce placement"
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
  pass "go vet workforce"
else
  fail "go vet workforce"
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
  echo "=== Workforce placement verification: ALL PASSED ==="
  exit 0
else
  echo "=== Workforce placement verification: FAILED ==="
  exit 1
fi
