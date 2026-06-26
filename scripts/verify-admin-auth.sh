#!/usr/bin/env bash
# Manual admin auth core verification — mock login + storage files + build.
#
# Usage: ./scripts/verify-admin-auth.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN="${ROOT}/web/admin"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Admin auth core verification ==="
echo

for f in \
  src/lib/auth/types.ts \
  src/lib/auth/storage.ts \
  src/lib/auth/AuthProvider.tsx \
  src/lib/auth/login-api.ts \
  src/lib/auth/dev-mock.ts \
  src/lib/api/client.ts \
  src/components/AuthRouterProvider.tsx \
  src/client.tsx
do
  if [[ -f "$ADMIN/$f" ]]; then
    pass "exists: web/admin/$f"
  else
    fail "missing: web/admin/$f"
  fi
done

echo
echo "--- npm run auth-smoke ---"
if (cd "$ADMIN" && npm run auth-smoke); then
  pass "auth mock smoke"
else
  fail "auth mock smoke"
fi

echo
echo "--- npm run build ---"
if (cd "$ADMIN" && npm run build); then
  pass "npm run build"
else
  fail "npm run build"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Admin auth core verification: ALL PASSED ==="
  exit 0
fi
echo "=== Admin auth core verification: FAILED ==="
exit 1
