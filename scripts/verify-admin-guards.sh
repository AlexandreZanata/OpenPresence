#!/usr/bin/env bash
# Manual admin router guards verification — build + HTTP smoke on protected routes.
#
# Usage: ./scripts/verify-admin-guards.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN="${ROOT}/web/admin"
FAIL=0
DEV_PID=""

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

cleanup() {
  if [[ -n "$DEV_PID" ]]; then
    kill "$DEV_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

ensure_dev_server() {
  if curl -sf http://localhost:5174/login 2>/dev/null | grep -q 'Registration ID'; then
    pass "dev server already running"
    return
  fi
  echo "Starting vite dev in background..."
  (cd "$ADMIN" && npm run dev >/tmp/openpresence-admin-dev.log 2>&1 &)
  DEV_PID=$!
  for _ in $(seq 1 30); do
    if curl -sf http://localhost:5174/login 2>/dev/null | grep -q 'Registration ID'; then
      pass "dev server ready"
      return
    fi
    sleep 1
  done
  fail "dev server did not start"
  tail -20 /tmp/openpresence-admin-dev.log 2>/dev/null || true
}

echo "=== Admin router guards verification ==="
echo

for f in \
  src/routes/_authenticated.tsx \
  src/routes/_authenticated/dashboard.tsx
do
  if [[ -f "$ADMIN/$f" ]]; then
    pass "exists: web/admin/$f"
  else
    fail "missing: web/admin/$f"
  fi
done

if [[ -f "$ADMIN/src/routes/dashboard.tsx" ]]; then
  fail "legacy route must be removed: web/admin/src/routes/dashboard.tsx"
else
  pass "legacy dashboard.tsx removed"
fi

if grep -q "beforeLoad" "$ADMIN/src/routes/_authenticated.tsx" 2>/dev/null; then
  pass "_authenticated beforeLoad guard"
else
  fail "_authenticated missing beforeLoad"
fi

if grep -q "redirect" "$ADMIN/src/routes/_authenticated.tsx" 2>/dev/null; then
  pass "_authenticated redirect to login"
else
  fail "_authenticated missing redirect"
fi

echo
echo "--- npm run build ---"
if (cd "$ADMIN" && npm run build); then
  pass "npm run build"
else
  fail "npm run build"
fi

if grep -q "'/_authenticated'" "$ADMIN/src/routeTree.gen.ts" 2>/dev/null; then
  pass "routeTree includes _authenticated layout"
else
  fail "routeTree missing _authenticated"
fi

echo
echo "--- HTTP smoke (unauthenticated) ---"
ensure_dev_server

LOGIN_URL="$(curl -sI http://localhost:5174/dashboard 2>/dev/null | tr -d '\r' | grep -i '^location:' | awk '{print $2}' || true)"
if [[ "$LOGIN_URL" == *"/login"* ]] && [[ "$LOGIN_URL" == *"redirect"* ]]; then
  pass "GET /dashboard redirects to /login?redirect=..."
elif curl -sf "http://localhost:5174/dashboard" 2>/dev/null | grep -q 'Registration ID\|Sign in'; then
  pass "GET /dashboard shows login (client guard)"
else
  fail "GET /dashboard did not redirect to login (got: ${LOGIN_URL:-none})"
fi

ROOT_URL="$(curl -sI http://localhost:5174/ 2>/dev/null | tr -d '\r' | grep -i '^location:' | awk '{print $2}' || true)"
if [[ "$ROOT_URL" == *"/login"* ]] || [[ "$ROOT_URL" == *"/dashboard"* ]]; then
  pass "GET / redirects to login or dashboard"
else
  fail "GET / did not redirect (got: ${ROOT_URL:-none})"
fi

echo
echo "--- Manual browser checks (required) ---"
echo "  1. Open http://localhost:5174/dashboard — should land on login with ?redirect="
echo "  2. Sign in (admin/admin) — should reach dashboard"
echo "  3. Refresh /dashboard — session persists"
echo "  4. Sign out — /dashboard sends you back to login"

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Admin router guards verification: ALL PASSED ==="
  exit 0
fi
echo "=== Admin router guards verification: FAILED ==="
exit 1
