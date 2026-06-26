#!/usr/bin/env bash
# Manual scaffold verification — run after task-00 or when changing repo layout.
#
# Usage: ./scripts/verify-scaffold.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== OpenPresence scaffold verification ==="
echo "Root: $ROOT"
echo

# --- filesystem layout ---
REQUIRED_DIRS=(
  "$ATTENDANCE/internal/domain"
  "$ATTENDANCE/internal/application"
  "$ATTENDANCE/internal/infrastructure"
  "$ATTENDANCE/internal/interfaces"
  "$ROOT/infra"
  "$ROOT/infra/k8s"
  "$ROOT/infra/terraform"
)

for dir in "${REQUIRED_DIRS[@]}"; do
  if [[ -d "$dir" ]]; then
    pass "directory exists: ${dir#$ROOT/}"
  else
    fail "missing directory: ${dir#$ROOT/}"
  fi
done

for file in "$ROOT/go.work" "$ATTENDANCE/go.mod" "$ROOT/infra/docker-compose.yml"; do
  if [[ -f "$file" ]]; then
    pass "file exists: ${file#$ROOT/}"
  else
    fail "missing file: ${file#$ROOT/}"
  fi
done

# --- Go toolchain ---
if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
else
  GO_VER="$(go version)"
  pass "go installed: $GO_VER"
  GO_MINOR="$(go version | sed -nE 's/.*go1\.([0-9]+).*/\1/p')"
  if [[ -n "$GO_MINOR" ]] && [[ "$GO_MINOR" -ge 22 ]]; then
    pass "go version >= 1.22"
  else
    fail "go version must be >= 1.22 (got: $GO_VER)"
  fi
fi

# --- build & test (real execution) ---
echo
echo "--- go build ./... (services/attendance) ---"
if (cd "$ATTENDANCE" && go build ./...); then
  pass "go build ./..."
else
  fail "go build ./..."
fi

echo
echo "--- go test -v ./... (services/attendance) ---"
if (cd "$ATTENDANCE" && go test -v ./...); then
  pass "go test ./..."
else
  fail "go test ./..."
fi

echo
echo "--- go vet ./... (services/attendance) ---"
if (cd "$ATTENDANCE" && go vet ./...); then
  pass "go vet ./..."
else
  fail "go vet ./..."
fi

echo
echo "--- go work sync (repo root) ---"
if (cd "$ROOT" && go work sync 2>/dev/null || true); then
  pass "go work present at root"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Scaffold verification: ALL PASSED ==="
  exit 0
else
  echo "=== Scaffold verification: FAILED ==="
  exit 1
fi
