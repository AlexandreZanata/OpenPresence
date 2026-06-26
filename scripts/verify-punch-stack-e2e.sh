#!/usr/bin/env bash
# Manual SubmitPunch full-stack E2E — Postgres + RLS + real biometric gRPC (BR-010, BR-014, BR-023).
#
# Usage: ./scripts/verify-punch-stack-e2e.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/application/punch/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== SubmitPunch full-stack E2E verification (BR-010, BR-014, BR-023) ==="
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
echo "--- domain: verify-punch.sh ---"
if "$ROOT/scripts/verify-punch.sh"; then
  pass "verify-punch.sh"
else
  fail "verify-punch.sh"
fi

echo
echo "--- use case: verify-punch-usecase.sh ---"
if "$ROOT/scripts/verify-punch-usecase.sh"; then
  pass "verify-punch-usecase.sh"
else
  fail "verify-punch-usecase.sh"
fi

echo
echo "--- go test -tags=integration $PKG -run E2E_Stack ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PKG" -run E2E_Stack); then
  pass "SubmitPunch full-stack E2E integration"
else
  fail "SubmitPunch full-stack E2E integration"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== SubmitPunch full-stack E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== SubmitPunch full-stack E2E verification: FAILED ==="
  exit 1
fi
