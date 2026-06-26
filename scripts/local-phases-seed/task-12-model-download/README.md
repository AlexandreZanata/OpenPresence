# Task 12 — Biometric model download (Rust/ops)

**Status:** pending  
**Phase ID:** task-12-model-download

## Goal

Provide **reproducible download** of ONNX models (RetinaFace, MiniFASNet ensemble, AuraFace) into `models/` with checksum verification. Enable local and CI setup for real image inference without committing weights to git.

## Scope

**In scope:**

- `scripts/download-models.sh` — fetch, verify SHA-256, layout per `docs/BIOMETRICS.md`
- `models/.gitkeep` or manifest `models/MANIFEST.json` with URLs and hashes
- Update `.env.example` (`ONNX_MODELS_PATH`, `BIOMETRIC_USE_STUB`)
- Manual verify script `scripts/verify-models.sh`

**Out of scope:**

- ONNX Runtime inference (task 13)
- Mobile on-device model bundling

## Acceptance

- Script downloads all required models on clean machine
- Checksum mismatch fails loudly
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
