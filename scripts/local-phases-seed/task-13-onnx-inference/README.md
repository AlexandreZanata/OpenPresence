# Task 13 — ONNX inference pipeline (Rust)

**Status:** pending  
**Phase ID:** task-13-onnx-inference

## Goal

Replace **stub-only** biometric processing with real **ONNX Runtime** inference when models are present: RetinaFace detection, MiniFASNet liveness ensemble, AuraFace 512-dim embedding. Image preprocessing and thresholds per `docs/BIOMETRICS.md`.

## Scope

**In scope:**

- `services/biometric/` — `ort` feature, session load at startup
- Wire existing preprocess + `ensemble_liveness` + `cosine_similarity` to real tensors
- Graceful fallback to stub when `BIOMETRIC_USE_STUB=true` or models missing
- Integration tests with sample images (committed small fixtures, not full models in CI)

**Out of scope:**

- pgvector search in Go (stays in attendance service)
- GPU provider tuning

## Acceptance

- With models: `verify-biometric.sh` exercises real inference path
- Without models: stub path unchanged
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
