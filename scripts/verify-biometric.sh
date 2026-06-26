#!/usr/bin/env bash
# Manual biometric service verification — unit tests + live server health + gRPC smoke.
#
# Usage: ./scripts/verify-biometric.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SVC="$ROOT/services/biometric"
FAIL=0
GRPC_PORT="${BIOMETRIC_GRPC_PORT:-19090}"
HTTP_PORT="${BIOMETRIC_HTTP_PORT:-19091}"
SERVER_PID=""

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

echo "=== Biometric service verification ==="
echo "Service: $SVC"
echo

if ! command -v cargo >/dev/null 2>&1; then
  fail "cargo not installed"
  exit 1
fi
pass "cargo installed: $(cargo --version)"

echo
echo "--- cargo test ---"
if (cd "$SVC" && cargo test --quiet); then
  pass "cargo test"
else
  fail "cargo test"
fi

echo
echo "--- cargo clippy ---"
if (cd "$SVC" && cargo clippy --quiet -- -D warnings); then
  pass "cargo clippy"
else
  fail "cargo clippy"
fi

echo
echo "--- cargo build biometric-server ---"
if (cd "$SVC" && cargo build --quiet --bin biometric-server); then
  pass "cargo build --bin biometric-server"
else
  fail "cargo build"
fi

echo
echo "--- live server: health + gRPC (stub processor) ---"
export BIOMETRIC_USE_STUB=true
export BIOMETRIC_GRPC_ADDR="127.0.0.1:${GRPC_PORT}"
export BIOMETRIC_HTTP_ADDR="127.0.0.1:${HTTP_PORT}"
export RUST_LOG=info

(cd "$SVC" && cargo run --quiet --bin biometric-server) &
SERVER_PID=$!

for _ in $(seq 1 40); do
  if curl -sf "http://127.0.0.1:${HTTP_PORT}/health/live" >/dev/null 2>&1; then
    break
  fi
  sleep 0.25
done

if curl -sf "http://127.0.0.1:${HTTP_PORT}/health/live" | grep -q '"status":"ok"'; then
  pass "GET /health/live"
else
  fail "GET /health/live"
fi

if curl -sf "http://127.0.0.1:${HTTP_PORT}/health/ready" | grep -q '"status":"ready"'; then
  pass "GET /health/ready"
else
  fail "GET /health/ready"
fi

if command -v grpcurl >/dev/null 2>&1; then
  if grpcurl -plaintext "127.0.0.1:${GRPC_PORT}" list 2>/dev/null | grep -q BiometricService; then
    pass "grpcurl list — gRPC server listening"
  else
    fail "grpcurl could not list gRPC services"
  fi
else
  if (echo >/dev/tcp/127.0.0.1/"${GRPC_PORT}") 2>/dev/null; then
    pass "TCP connect gRPC port ${GRPC_PORT}"
  else
    fail "gRPC port ${GRPC_PORT} not accepting connections"
  fi
fi

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Biometric verification: ALL PASSED ==="
  exit 0
else
  echo "=== Biometric verification: FAILED ==="
  exit 1
fi
