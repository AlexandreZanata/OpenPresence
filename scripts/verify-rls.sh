#!/usr/bin/env bash
# Manual PostgreSQL RLS verification — integration tests via testcontainers.
#
# Usage: ./scripts/verify-rls.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/infrastructure/postgres/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== PostgreSQL RLS verification ==="
echo "Package: $PKG"
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

if ! command -v docker >/dev/null 2>&1; then
  fail "docker not installed (required for testcontainers)"
  exit 1
fi
if ! docker info >/dev/null 2>&1; then
  fail "docker daemon not running"
  exit 1
fi
pass "docker available"

echo
echo "--- migrations present ---"
for f in 001_create_tenants.sql 002_create_employees.sql 003_create_punch_records.sql \
         004_create_face_embeddings.sql 005_enable_rls.sql 006_create_app_role.sql; do
  if [[ -f "$ATTENDANCE/migrations/$f" ]]; then
    pass "migration $f"
  else
    fail "missing migration $f"
  fi
done

echo
echo "--- go test ./... (unit) ---"
if (cd "$ATTENDANCE" && go test ./...); then
  pass "unit tests"
else
  fail "unit tests"
fi

echo
echo "--- go test -tags=integration $PKG ---"
if (cd "$ATTENDANCE" && go test -tags=integration -v -count=1 $PKG); then
  pass "RLS integration tests"
else
  fail "RLS integration tests"
fi

echo
echo "--- go vet $PKG ---"
if (cd "$ATTENDANCE" && go vet $PKG); then
  pass "go vet postgres"
else
  fail "go vet postgres"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== RLS verification: ALL PASSED ==="
  exit 0
else
  echo "=== RLS verification: FAILED ==="
  exit 1
fi
