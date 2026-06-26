#!/usr/bin/env bash
# Manual SubmitPunch use case verification — unit + integration tests with testcontainers.
#
# Usage: ./scripts/verify-punch-usecase.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/application/punch/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== SubmitPunch use case verification ==="
echo "Package: $PKG"
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

echo
echo "--- go test (unit) $PKG ---"
if (cd "$ATTENDANCE" && go test -v "$PKG" -count=1); then
  pass "go test unit"
else
  fail "go test unit"
fi

echo
echo "--- go test -tags=integration $PKG ---"
if (cd "$ATTENDANCE" && go test -tags=integration -v "$PKG" -count=1); then
  pass "go test integration"
else
  fail "go test integration"
fi

echo
echo "--- go vet $PKG ---"
if (cd "$ATTENDANCE" && go vet "$PKG"); then
  pass "go vet"
else
  fail "go vet"
fi

echo
echo "--- layer isolation ---"
if ! grep -rE 'internal/(infrastructure|interfaces)' "$ATTENDANCE/internal/application/punch"/*.go 2>/dev/null | grep -v '_integration_test.go'; then
  pass "application punch has no infrastructure imports in production code"
else
  fail "application punch imports infrastructure in production code"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== SubmitPunch use case verification: ALL PASSED ==="
  exit 0
else
  echo "=== SubmitPunch use case verification: FAILED ==="
  exit 1
fi
