# Biometric Service (Rust)

gRPC service for face liveness ensemble and embedding extraction. Internal-only — not exposed on public REST.

## Stack

- **tonic** — gRPC (`VerifyPunch`, `EnrollFace`, `DeleteProfile`)
- **axum** — HTTP `/health/live`, `/health/ready`
- **ONNX Runtime** (`ort`, optional feature `onnx`) — production inference when models are present

## Modes

| Mode | When | Env |
|------|------|-----|
| **Stub** (default) | No `ONNX_MODELS_PATH` | `BIOMETRIC_USE_STUB=true` |
| **ONNX** | Models on disk | `ONNX_MODELS_PATH=./models` + `--features onnx` |

Stub mode uses deterministic pipeline math for local dev and CI without downloading AuraFace / MiniFASNet weights.

Download production weights:

```bash
./scripts/download-models.sh
./scripts/verify-models.sh
```

See `models/MANIFEST.json` and `docs/BIOMETRICS.md`.

## Commands

```bash
cargo test
cargo test --features onnx          # includes ONNX unit test (ignored without models)
cargo clippy -- -D warnings
cargo run --bin biometric-server    # stub when BIOMETRIC_USE_STUB=true or no ONNX_MODELS_PATH
cargo run --features onnx --bin biometric-server   # ONNX when ONNX_MODELS_PATH is set
```

From repo root:

```bash
./scripts/verify-biometric.sh
ONNX_MODELS_PATH=./models ./scripts/verify-biometric.sh
```

## Ports (default)

| Protocol | Address | Purpose |
|----------|---------|---------|
| gRPC | `0.0.0.0:9090` | `BIOMETRIC_GRPC_ADDR` |
| HTTP | `0.0.0.0:9091` | `BIOMETRIC_HTTP_ADDR` |

## Proto

`proto/biometric.proto` — package `openpresence.biometric.v1`. See `docs/API-CONTRACT.md`.

## Related docs

- `docs/BIOMETRICS.md` — model licensing and pipeline
- `docs/adr/ADR-002-biometric-stack.md`
- `docs/AGENT-IMPLEMENTATION-GUIDE.md` — Task 01
