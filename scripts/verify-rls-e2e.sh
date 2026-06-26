#!/usr/bin/env bash
# Manual RLS E2E verification — postgres RLS + SubmitPunch + enrollment isolation.
#
# Usage: ./scripts/verify-rls-e2e.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PUNCH_PKG="./internal/application/punch/..."
ENROLL_PKG="./internal/application/enrollment/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== RLS E2E verification (multi-tenant isolation) ==="
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

if ! command -v docker >/dev/null 2>&1 || ! docker info >/dev/null 2>&1; then
  fail "docker required for RLS integration tests"
  exit 1
fi
pass "docker available"

echo
echo "--- base: verify-rls.sh ---"
if "$ROOT/scripts/verify-rls.sh"; then
  pass "verify-rls.sh"
else
  fail "verify-rls.sh"
fi

echo
echo "--- go test -tags=integration $PUNCH_PKG -run E2E_RLS ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PUNCH_PKG" -run E2E_RLS); then
  pass "SubmitPunch RLS E2E integration"
else
  fail "SubmitPunch RLS E2E integration"
fi

echo
echo "--- go test -tags=integration $ENROLL_PKG -run E2E_RLS ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$ENROLL_PKG" -run E2E_RLS); then
  pass "enrollment RLS E2E integration"
else
  fail "enrollment RLS E2E integration"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== RLS E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== RLS E2E verification: FAILED ==="
  exit 1
fi
