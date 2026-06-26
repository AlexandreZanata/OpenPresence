#!/usr/bin/env bash
# Download ONNX biometric models from models/MANIFEST.json with SHA-256 verification.
#
# Usage: ./scripts/download-models.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODELS_DIR="${ONNX_MODELS_PATH:-$ROOT/models}"
MANIFEST="$ROOT/models/MANIFEST.json"

if [[ ! -f "$MANIFEST" ]]; then
  echo "FAIL: manifest not found at $MANIFEST" >&2
  exit 1
fi

if ! command -v curl >/dev/null 2>&1; then
  echo "FAIL: curl is required" >&2
  exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
  echo "FAIL: python3 is required to parse MANIFEST.json" >&2
  exit 1
fi

mkdir -p "$MODELS_DIR"

echo "=== OpenPresence model download ==="
echo "Manifest: $MANIFEST"
echo "Target:   $MODELS_DIR"
echo

download_one() {
  local file="$1" url="$2" sha="$3"
  local dest="$MODELS_DIR/$file"
  local tmp="$dest.part"

  if [[ -f "$dest" ]]; then
    if echo "$sha  $dest" | sha256sum -c --status 2>/dev/null; then
      echo "SKIP: $file (checksum OK)"
      return 0
    fi
    echo "WARN: $file exists but checksum mismatch — re-downloading"
    rm -f "$dest"
  fi

  echo "FETCH: $file"
  if ! curl -fsSL --retry 3 --retry-delay 2 --connect-timeout 30 -o "$tmp" "$url"; then
    rm -f "$tmp"
    echo "FAIL: download failed for $file" >&2
    exit 1
  fi

  if [[ ! -s "$tmp" ]]; then
    rm -f "$tmp"
    echo "FAIL: empty download for $file" >&2
    exit 1
  fi

  if ! echo "$sha  $tmp" | sha256sum -c --status; then
    rm -f "$tmp"
    echo "FAIL: checksum mismatch for $file" >&2
    exit 1
  fi

  mv "$tmp" "$dest"
  echo "OK:   $file"
}

export ROOT MODELS_DIR
while IFS=$'\t' read -r file url sha; do
  download_one "$file" "$url" "$sha"
done < <(python3 - <<'PY'
import json
import os
from pathlib import Path

manifest = Path(os.environ["ROOT"]) / "models" / "MANIFEST.json"
data = json.loads(manifest.read_text())
for item in data["models"]:
    print(f"{item['file']}\t{item['url']}\t{item['sha256']}")
PY
)

echo
echo "=== Model download complete ==="
