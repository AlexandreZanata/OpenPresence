#!/usr/bin/env bash
# Manual hierarchy authorization verification — real go test + coverage on disk.
#
# Usage: ./scripts/verify-authorization.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
ORG_PKG="./internal/domain/organization/..."
APP_PKG="./internal/application/authorization/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Hierarchy authorization verification ==="
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

echo
echo "--- go test -v -run 'Auth|Approve|Descendant|Export|Auditor|CrossTenant' $ORG_PKG ---"
if (cd "$ATTENDANCE" && go test -v -run 'Auth|Approve|Descendant|Export|Auditor|CrossTenant' $ORG_PKG); then
  pass "go test organization authorization"
else
  fail "go test organization authorization"
fi

echo
echo "--- go test -v $APP_PKG ---"
if (cd "$ATTENDANCE" && go test -v $APP_PKG); then
  pass "go test application authorization"
else
  fail "go test application authorization"
fi

echo
echo "--- go vet $ORG_PKG $APP_PKG ---"
if (cd "$ATTENDANCE" && go vet $ORG_PKG $APP_PKG); then
  pass "go vet authorization"
else
  fail "go vet authorization"
fi

echo
echo "--- layer isolation ---"
if ! grep -rE 'internal/(infrastructure|interfaces)' "$ATTENDANCE/internal/domain/organization" 2>/dev/null; then
  pass "organization domain has no outer imports"
else
  fail "organization domain imports outer layers"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Authorization verification: ALL PASSED ==="
  exit 0
else
  echo "=== Authorization verification: FAILED ==="
  exit 1
fi
