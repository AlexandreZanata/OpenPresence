# Task 01 — Biometric service (Rust)

**Status:** done  
**Phase ID:** task-01-biometric-service

## Goal

gRPC server in `services/biometric/` for face detection, liveness ensemble, and embedding extraction using ONNX Runtime. Implements AGENT guide Task 01.

## Scope

**In scope:**

- Rust crate with tonic gRPC
- `FaceProcessor` with ONNX sessions at startup (`Arc`)
- `VerifyPunch` and `EnrollFace` RPCs (proto defined)
- Unit tests for preprocess, liveness ensemble, cosine similarity

**Out of scope:**

- pgvector search (stays in Go attendance service)
- Model weight download CI (document manual `models/` setup)

## Acceptance

- `cargo test` passes in `services/biometric/`
- gRPC server starts and responds to health check
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
