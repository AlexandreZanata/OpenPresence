#!/usr/bin/env bash
# Regenerate Go gRPC stubs from services/biometric/proto/biometric.proto
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="$ROOT/services/attendance/internal/infrastructure/biometric/pb"
PROTO="$ROOT/services/biometric/proto/biometric.proto"

if ! command -v protoc >/dev/null 2>&1; then
  echo "protoc not found — install protobuf-compiler or use a release from github.com/protocolbuffers/protobuf" >&2
  exit 1
fi

mkdir -p "$OUT"
protoc \
  --go_out="$OUT" --go_opt=paths=source_relative \
  --go-grpc_out="$OUT" --go-grpc_opt=paths=source_relative \
  -I "$ROOT/services/biometric/proto" \
  "$PROTO"

echo "Generated: $OUT"
