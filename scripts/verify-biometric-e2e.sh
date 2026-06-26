#!/usr/bin/env bash
# Manual biometric stack E2E — Rust service + Attendance gRPC client (BR-002, BR-010).
#
# Usage:
#   ./scripts/verify-biometric-e2e.sh
#   ONNX_MODELS_PATH=./models ./scripts/verify-biometric-e2e.sh   # also runs ONNX verify-biometric

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ATTENDANCE="$ROOT/services/attendance"
PUNCH_PKG="./internal/application/punch/..."
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Biometric stack E2E verification (BR-002, BR-010) ==="
echo

if ! command -v go >/dev/null 2>&1; then
  fail "go not installed"
  exit 1
fi
pass "go installed: $(go version)"

if ! command -v cargo >/dev/null 2>&1; then
  fail "cargo not installed (required for biometric-server subprocess)"
  exit 1
fi
pass "cargo installed: $(cargo --version)"

if ! command -v docker >/dev/null 2>&1 || ! docker info >/dev/null 2>&1; then
  fail "docker required for SubmitPunch + Postgres integration"
  exit 1
fi
pass "docker available"

echo
echo "--- stub: verify-biometric.sh ---"
if "$ROOT/scripts/verify-biometric.sh"; then
  pass "verify-biometric.sh (stub)"
else
  fail "verify-biometric.sh (stub)"
fi

echo
echo "--- go test -tags=integration $PUNCH_PKG -run BiometricGrpc ---"
if (cd "$ATTENDANCE" && go test -tags=integration -count=1 "$PUNCH_PKG" -run BiometricGrpc); then
  pass "Attendance → biometric gRPC integration"
else
  fail "Attendance → biometric gRPC integration"
fi

if [[ -f "${ROOT}/models/auraface.onnx" ]]; then
  echo
  echo "--- ONNX: verify-biometric.sh ---"
  if ONNX_MODELS_PATH="${ROOT}/models" "$ROOT/scripts/verify-biometric.sh"; then
    pass "verify-biometric.sh (ONNX)"
  else
    fail "verify-biometric.sh (ONNX)"
  fi
else
  echo
  echo "SKIP: ONNX verify — models/auraface.onnx not present (run ./scripts/download-models.sh)"
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Biometric stack E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== Biometric stack E2E verification: FAILED ==="
  exit 1
fi
