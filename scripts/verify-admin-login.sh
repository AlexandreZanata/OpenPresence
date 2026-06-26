#!/usr/bin/env bash
# Manual admin login form verification — build + HTTP smoke on /login.
#
# Usage: ./scripts/verify-admin-login.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN="${ROOT}/web/admin"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Admin login form verification ==="
echo

for f in \
  src/routes/login.tsx \
  src/routes/_authenticated/dashboard.tsx \
  src/components/login/LoginForm.tsx
do
  if [[ -f "$ADMIN/$f" ]]; then
    pass "exists: web/admin/$f"
  else
    fail "missing: web/admin/$f"
  fi
done

echo
echo "--- npm run build ---"
if (cd "$ADMIN" && npm run build); then
  pass "npm run build"
else
  fail "npm run build"
fi

echo
echo "--- HTTP smoke /login ---"
if curl -sf http://localhost:5174/login 2>/dev/null | grep -q 'Sign in'; then
  pass "dev server /login renders Sign in"
elif curl -sf http://localhost:5174/login 2>/dev/null | grep -q 'Registration ID'; then
  pass "dev server /login renders form"
else
  echo "Starting vite dev in background..."
  (cd "$ADMIN" && npm run dev >/tmp/openpresence-admin-dev.log 2>&1 &)
  DEV_PID=$!
  trap 'kill "$DEV_PID" 2>/dev/null || true' EXIT
  for _ in $(seq 1 30); do
    if curl -sf http://localhost:5174/login 2>/dev/null | grep -q 'Registration ID'; then
      pass "GET /login renders login form"
      break
    fi
    sleep 1
  done
  if ! curl -sf http://localhost:5174/login 2>/dev/null | grep -q 'Registration ID'; then
    fail "GET /login did not render form"
    tail -20 /tmp/openpresence-admin-dev.log 2>/dev/null || true
  fi
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Admin login form verification: ALL PASSED ==="
  exit 0
fi
echo "=== Admin login form verification: FAILED ==="
exit 1
