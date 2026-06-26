#!/usr/bin/env bash
# Manual mobile KMP E2E verification — commonTest + offline sync with mock API.
#
# Usage: ./scripts/verify-mobile-e2e.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Mobile KMP E2E verification (UC-001 partial) ==="
echo

if [[ ! -x "$ROOT/gradlew" ]]; then
  fail "gradlew not found"
  exit 1
fi
pass "gradlew present"

echo
echo "--- base: verify-mobile.sh ---"
if "$ROOT/scripts/verify-mobile.sh"; then
  pass "verify-mobile.sh"
else
  fail "verify-mobile.sh"
fi

echo
echo "--- ./gradlew :mobile:shared:jvmTest --tests OfflinePunchSyncE2ETest ---"
if (cd "$ROOT" && ./gradlew :mobile:shared:jvmTest --no-daemon -q \
  --tests "com.openpresence.punch.e2e.OfflinePunchSyncE2ETest"); then
  pass "OfflinePunchSyncE2ETest"
else
  fail "OfflinePunchSyncE2ETest"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Mobile KMP E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== Mobile KMP E2E verification: FAILED ==="
  exit 1
fi
