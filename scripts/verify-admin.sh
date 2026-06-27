#!/usr/bin/env bash
# Manual admin shell verification — build + API health smoke.
#
# Usage: ./scripts/verify-admin.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN="${ROOT}/web/admin"
FAIL=0
DEV_PID=""

API_BASE="${VITE_API_BASE_URL:-http://127.0.0.1:8088}"

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

echo "=== Admin shell verification ==="
echo

for f in \
  src/components/AdminShell.tsx \
  src/components/admin-shell.css \
  src/lib/api/health.ts \
  src/routes/_authenticated/dashboard.tsx
do
  if [[ -f "$ADMIN/$f" ]]; then
    pass "exists: web/admin/$f"
  else
    fail "missing: web/admin/$f"
  fi
done

if grep -q "AdminShell" "$ADMIN/src/routes/_authenticated.tsx" 2>/dev/null; then
  pass "_authenticated wraps AdminShell"
else
  fail "_authenticated missing AdminShell"
fi

if grep -q "useQuery" "$ADMIN/src/routes/_authenticated/dashboard.tsx" 2>/dev/null; then
  pass "dashboard uses TanStack Query for health"
else
  fail "dashboard missing useQuery health check"
fi

echo
echo "--- npm run build ---"
if (cd "$ADMIN" && npm run build); then
  pass "npm run build"
else
  fail "npm run build"
fi

echo
echo "--- API health smoke (${API_BASE}) ---"
if RESP="$(curl -sf "${API_BASE}/health/live" 2>/dev/null)" && echo "$RESP" | grep -q '"status":"ok"'; then
  pass "GET /health/live → ok"
else
  fail "GET /health/live failed (start backend: ./scripts/dev-backend.sh start)"
fi

echo
echo "--- HTTP smoke /login ---"
ensure_dev_server
if curl -sf http://localhost:5174/login 2>/dev/null | grep -q 'Registration ID'; then
  pass "GET /login renders form"
else
  fail "GET /login did not render form"
fi

echo
echo "--- Manual browser checks (required) ---"
echo "  1. Sign in at http://localhost:5174/login (admin/admin)"
echo "  2. Shell visible: sidebar (Dashboard, Employees, Settings) + header with user"
echo "  3. Dashboard card shows API status Online when backend is up"
echo "  4. Sign out → returns to login"

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Admin shell verification: ALL PASSED ==="
  exit 0
fi
echo "=== Admin shell verification: FAILED ==="
exit 1
