#!/usr/bin/env bash
# Manual enrollment E2E verification — gRPC EnrollFace (BR-001–003).
#
# Usage:
#   ./scripts/verify-enrollment.sh              # stub mode (default)
#   ONNX_MODELS_PATH=./models ./scripts/verify-enrollment.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SVC="$ROOT/services/biometric"
FIXTURE_DIR="$SVC/tests/fixtures"
FAIL=0
GRPC_PORT="${BIOMETRIC_GRPC_PORT:-19092}"
HTTP_PORT="${BIOMETRIC_HTTP_PORT:-19093}"
SERVER_PID=""
CARGO_EXTRA=()

if [[ -n "${ONNX_MODELS_PATH:-}" ]] && [[ "${BIOMETRIC_USE_STUB:-}" != "true" ]]; then
  CARGO_EXTRA=(--features onnx)
  echo "Mode: ONNX"
else
  export BIOMETRIC_USE_STUB=true
  echo "Mode: stub"
fi

cleanup() {
  if [[ -n "$SERVER_PID" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

if [[ -f "$HOME/.cargo/env" ]]; then
  # shellcheck source=/dev/null
  . "$HOME/.cargo/env"
fi

echo "=== Enrollment E2E verification (BR-001–003) ==="
echo "Service: $SVC"
echo

if ! command -v cargo >/dev/null 2>&1; then
  fail "cargo not installed"
  exit 1
fi
pass "cargo installed: $(cargo --version)"

echo
echo "--- cargo test enrollment_e2e ---"
if (cd "$SVC" && cargo test --quiet "${CARGO_EXTRA[@]}" enrollment_e2e -- --test-threads=1); then
  pass "enrollment_e2e integration tests"
else
  fail "enrollment_e2e integration tests"
fi

echo
echo "--- fixture JPEGs for grpcurl ---"
if [[ ! -f "$FIXTURE_DIR/valid_128.jpg" ]]; then
  (cd "$SVC" && cargo test --quiet write_enrollment_fixture_jpegs -- --ignored --test-threads=1)
fi
for f in valid_128.jpg low_liveness_128.jpg low_quality_32.jpg; do
  if [[ -f "$FIXTURE_DIR/$f" ]]; then
    pass "fixture $f"
  else
    fail "missing fixture $f"
  fi
done

echo
echo "--- live server: grpcurl EnrollFace smoke ---"
export BIOMETRIC_GRPC_ADDR="127.0.0.1:${GRPC_PORT}"
export BIOMETRIC_HTTP_ADDR="127.0.0.1:${HTTP_PORT}"
export RUST_LOG=warn

(cd "$SVC" && cargo run --quiet --bin biometric-server "${CARGO_EXTRA[@]}") &
SERVER_PID=$!

for _ in $(seq 1 60); do
  if curl -sf "http://127.0.0.1:${HTTP_PORT}/health/live" >/dev/null 2>&1; then
    break
  fi
  sleep 0.25
done

if ! command -v grpcurl >/dev/null 2>&1; then
  echo "SKIP: grpcurl not installed — integration tests above are sufficient"
else
  proto="$SVC/proto/biometric.proto"
  b64_jpeg="$(base64 -w0 "$FIXTURE_DIR/valid_128.jpg")"
  for angle in FRONTAL LEFT_15 RIGHT_15; do
    if grpcurl -plaintext -import-path "$SVC/proto" -proto biometric.proto \
      -d "{\"frame_jpeg\":\"${b64_jpeg}\",\"employee_id\":\"emp-manual\",\"tenant_id\":\"tenant-manual\",\"angle\":\"${angle}\"}" \
      "127.0.0.1:${GRPC_PORT}" openpresence.biometric.v1.BiometricService/EnrollFace \
      | grep -q '"isLive": true'; then
      pass "grpcurl EnrollFace $angle"
    else
      fail "grpcurl EnrollFace $angle"
    fi
  done

  b64_bad="$(base64 -w0 "$FIXTURE_DIR/low_liveness_128.jpg")"
  if grpcurl -plaintext -import-path "$SVC/proto" -proto biometric.proto \
    -d "{\"frame_jpeg\":\"${b64_bad}\",\"employee_id\":\"emp-manual\",\"tenant_id\":\"tenant-manual\",\"angle\":\"FRONTAL\"}" \
    "127.0.0.1:${GRPC_PORT}" openpresence.biometric.v1.BiometricService/EnrollFace \
    | grep -q '"isLive": false'; then
    pass "grpcurl EnrollFace liveness reject (BR-002)"
  else
    fail "grpcurl EnrollFace liveness reject (BR-002)"
  fi
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Enrollment E2E verification: ALL PASSED ==="
  exit 0
else
  echo "=== Enrollment E2E verification: FAILED ==="
  exit 1
fi
