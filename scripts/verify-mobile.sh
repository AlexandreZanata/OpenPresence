#!/usr/bin/env bash
# Manual KMP shared module verification — Gradle check + layer isolation.
#
# Usage: ./scripts/verify-mobile.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHARED="$ROOT/mobile/shared"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Mobile shared (KMP) verification ==="
echo "Module: mobile/shared"
echo

if [[ ! -x "$ROOT/gradlew" ]]; then
  fail "gradlew not found — run: gradle wrapper"
  exit 1
fi
pass "gradlew present"

echo
echo "--- ./gradlew :mobile:shared:jvmTest ---"
if (cd "$ROOT" && ./gradlew :mobile:shared:jvmTest --no-daemon -q); then
  pass "jvmTest"
else
  fail "jvmTest"
fi

echo
echo "--- ./gradlew :mobile:shared:check ---"
if (cd "$ROOT" && ./gradlew :mobile:shared:check --no-daemon -q); then
  pass "check"
else
  fail "check"
fi

echo
echo "--- layer isolation (no geofence math in mobile shared) ---"
if ! grep -rE 'haversine|ray.?cast|polygon' "$SHARED/src" 2>/dev/null; then
  pass "mobile shared has no server geofence algorithm"
else
  fail "mobile shared duplicates geofence rules"
fi

echo
echo "--- PunchState matches MOBILE-FLOWS ---"
REQUIRED_STATES=(
  Idle
  CheckingDevice
  WaitingLocation
  OpeningCamera
  DetectingFace
  CheckingLiveness
  Submitting
  Success
  Suspicious
  Error
  DeviceWarning
  OutOfGeofence
)
for state in "${REQUIRED_STATES[@]}"; do
  if grep -q "$state" "$SHARED/src/commonMain/kotlin/com/openpresence/punch/presentation/PunchState.kt"; then
    pass "PunchState.$state defined"
  else
    fail "PunchState.$state missing"
  fi
done

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Mobile verification: ALL PASSED ==="
  exit 0
else
  echo "=== Mobile verification: FAILED ==="
  exit 1
fi
