#!/usr/bin/env bash
# Manual AttendancePolicy domain verification — real go test + coverage on disk.
#
# Usage: ./scripts/verify-attendance-policy.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PKG="./internal/domain/organization/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== AttendancePolicy inheritance verification ==="
echo "Package: $PKG"
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

echo
echo "--- go test -v -run 'Policy|Merge|Effective|PathFromRoot|Preset' $PKG ---"
if (cd "$ATTENDANCE" && go test -v -run 'Policy|Merge|Effective|PathFromRoot|Preset' $PKG); then
  pass "go test attendance policy"
else
  fail "go test attendance policy"
fi

echo
echo "--- go test -cover $PKG ---"
COVER_OUTPUT="$(cd "$ATTENDANCE" && go test -cover $PKG 2>&1)"
echo "$COVER_OUTPUT"
if echo "$COVER_OUTPUT" | grep -q "ok"; then
  pass "coverage run"
else
  fail "coverage run"
fi

echo
echo "--- go vet $PKG ---"
if (cd "$ATTENDANCE" && go vet $PKG); then
  pass "go vet organization policy"
else
  fail "go vet organization policy"
fi

echo
echo "--- layer isolation ---"
if ! grep -rE 'internal/(infrastructure|interfaces|application)' "$ATTENDANCE/internal/domain/organization" 2>/dev/null; then
  pass "organization package has no outer layer imports"
else
  fail "organization imports outer layers"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== AttendancePolicy verification: ALL PASSED ==="
  exit 0
else
  echo "=== AttendancePolicy verification: FAILED ==="
  exit 1
fi
