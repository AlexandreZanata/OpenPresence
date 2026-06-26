#!/usr/bin/env bash
# Verify ONNX models exist and match models/MANIFEST.json checksums.
#
# Usage: ./scripts/verify-models.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODELS_DIR="${ONNX_MODELS_PATH:-$ROOT/models}"
MANIFEST="$ROOT/models/MANIFEST.json"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1" >&2; FAIL=1; }

echo "=== Model verification ==="
echo "Manifest: $MANIFEST"
echo "Models:   $MODELS_DIR"
echo

if [[ ! -f "$MANIFEST" ]]; then
  fail "MANIFEST.json missing"
  exit 1
fi
pass "MANIFEST.json present"

if ! command -v python3 >/dev/null 2>&1; then
  fail "python3 required"
  exit 1
fi

export ROOT MODELS_DIR
while IFS=$'\t' read -r file sha; do
  dest="$MODELS_DIR/$file"
  if [[ ! -f "$dest" ]]; then
    fail "missing file: $file (run ./scripts/download-models.sh)"
    continue
  fi
  if echo "$sha  $dest" | sha256sum -c --status 2>/dev/null; then
    pass "checksum OK: $file"
  else
    fail "checksum mismatch: $file"
  fi
done < <(python3 - <<'PY'
import json
import os
from pathlib import Path

manifest = Path(os.environ["ROOT"]) / "models" / "MANIFEST.json"
data = json.loads(manifest.read_text())
for item in data["models"]:
    print(f"{item['file']}\t{item['sha256']}")
PY
)

echo
if [[ "$FAIL" -eq 0 ]]; then
  echo "=== Model verification: ALL PASSED ==="
  exit 0
else
  echo "=== Model verification: FAILED ==="
  exit 1
fi
