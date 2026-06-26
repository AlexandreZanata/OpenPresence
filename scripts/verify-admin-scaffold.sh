#!/usr/bin/env bash
# Manual admin panel scaffold verification — build + dev HTTP smoke.
#
# Usage: ./scripts/verify-admin-scaffold.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN="${ROOT}/web/admin"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Admin scaffold verification ==="
echo

if ! command -v npm >/dev/null 2>&1; then
  fail "npm not installed"
  exit 1
fi
pass "npm installed"

for f in package.json vite.config.ts src/router.tsx src/routes/index.tsx .env.example; do
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
echo "--- dev server HTTP smoke (localhost:5174) ---"
if curl -sf http://localhost:5174/ 2>/dev/null | grep -q 'OpenPresence Admin'; then
  pass "dev server already serving landing page"
else
  echo "Starting vite dev in background..."
  (cd "$ADMIN" && npm run dev >/tmp/openpresence-admin-dev.log 2>&1 &)
  DEV_PID=$!
  trap 'kill "$DEV_PID" 2>/dev/null || true' EXIT
  for _ in $(seq 1 30); do
    if curl -sf http://localhost:5174/ 2>/dev/null | grep -q 'OpenPresence Admin'; then
      pass "GET / renders OpenPresence Admin"
      break
    fi
    sleep 1
  done
  if ! curl -sf http://localhost:5174/ 2>/dev/null | grep -q 'OpenPresence Admin'; then
    fail "dev server did not render landing page"
    tail -20 /tmp/openpresence-admin-dev.log 2>/dev/null || true
  fi
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Admin scaffold verification: ALL PASSED ==="
  exit 0
fi
echo "=== Admin scaffold verification: FAILED ==="
exit 1
