#!/usr/bin/env bash
# Manual authorization E2E verification — unit ABAC + Postgres integration (BR-040–043).
#
# Usage: ./scripts/verify-authorization-e2e.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/application/authorization/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Authorization E2E verification (ABAC) ==="
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

if ! command -v docker >/dev/null 2>&1 || ! docker info >/dev/null 2>&1; then
  fail "docker required for authorization integration tests"
  exit 1
fi
pass "docker available"

echo
echo "--- domain: verify-authorization.sh ---"
if "$ROOT/scripts/verify-authorization.sh"; then
  pass "verify-authorization.sh"
else
  fail "verify-authorization.sh"
fi

echo
echo "--- go test -tags=integration $PKG -run E2E ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PKG" -run E2E); then
  pass "authorization E2E integration"
else
  fail "authorization E2E integration"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Authorization E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== Authorization E2E verification: FAILED ==="
  exit 1
fi
