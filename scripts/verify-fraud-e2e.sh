#!/usr/bin/env bash
# Manual fraud E2E verification — domain tests + SubmitPunch integration (BR-012–013).
#
# Usage: ./scripts/verify-fraud-e2e.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/application/punch/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Fraud E2E verification (BR-012–013) ==="
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
echo "--- domain: verify-fraud.sh ---"
if "$ROOT/scripts/verify-fraud.sh"; then
  pass "verify-fraud.sh"
else
  fail "verify-fraud.sh"
fi

echo
echo "--- go test -tags=integration $PKG -run E2E_Fraud ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PKG" -run E2E_Fraud); then
  pass "SubmitPunch fraud E2E integration"
else
  fail "SubmitPunch fraud E2E integration"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Fraud E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== Fraud E2E verification: FAILED ==="
  exit 1
fi
