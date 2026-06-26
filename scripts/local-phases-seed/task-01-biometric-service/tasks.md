# Tasks — Biometric service

## Preparation

- [x] Read [README.md](README.md) and [official_source.md](official_source.md)
- [x] Run `./agent-harness/resolve-rules.sh owasp security biometric`
- [x] Rust >= 1.78 and `protoc` installed (via protoc-bin-vendored in build)
- [x] ONNX models optional — stub mode for dev/CI

## Scaffold

- [x] `cargo new` under `services/biometric/` (library + binary)
- [x] Add dependencies: axum, tonic, prost, image, tokio, serde_json (+ ort optional)
- [x] Define `proto/biometric.proto` from `docs/API-CONTRACT.md`
- [x] `build.rs` + protoc-bin-vendored + tonic-build

## Core implementation

- [x] `StubProcessor` / `BiometricProcessor` trait (ONNX via `--features onnx`)
- [x] `preprocess_for_liveness` — 80×80 BGR, normalize
- [x] `preprocess_for_recognition` — 112×112 RGB normalize
- [x] `ensemble_liveness` — average softmax, threshold 0.80
- [x] `cosine_similarity` for 512-dim vectors
- [x] gRPC `VerifyPunch` handler
- [x] gRPC `EnrollFace` handler
- [x] `/health/live` and `/health/ready` endpoints

## TDD — tests

- [x] `test_liveness_real_face_above_threshold`
- [x] `test_liveness_printed_photo_rejected`
- [x] `test_face_recognition_same_person`
- [x] `test_face_recognition_different_persons`
- [x] `test_cosine_similarity_identical_vectors`
- [x] `grpc_verify_punch_and_enroll_roundtrip` (integration)

## Validation

- [x] `cd services/biometric && cargo test`
- [x] `cargo clippy -- -D warnings`
- [x] `./scripts/verify-biometric.sh` — live server + curl health
- [x] `models/` in `.gitignore`

## Completion

- [x] All steps above marked `[x]`
- [x] Update `.local/phases/README.md` active task
