# Tasks — Biometric service

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh owasp security biometric`
- [ ] Rust >= 1.78 and `protoc` installed
- [ ] ONNX models placed in `models/` (see `docs/BIOMETRICS.md`)

## Scaffold

- [ ] `cargo new` under `services/biometric/` (library + binary)
- [ ] Add dependencies: axum, tonic, prost, ort, ndarray, image, tokio, serde_json
- [ ] Define `proto/biometric.proto` from `docs/API-CONTRACT.md`
- [ ] `build.rs` + buf or tonic-build for codegen

## Core implementation

- [ ] `FaceProcessor` struct with sessions loaded once at startup
- [ ] `preprocess_for_liveness` — 80×80 BGR, normalize
- [ ] `preprocess_for_recognition` — 112×112 RGB, affine warp from landmarks
- [ ] `ensemble_liveness` — average softmax, threshold 0.80
- [ ] `cosine_similarity` for `[f32; 512]`
- [ ] gRPC `VerifyPunch` handler
- [ ] gRPC `EnrollFace` handler
- [ ] `/health/live` and `/health/ready` endpoints

## TDD — tests

- [ ] `test_liveness_real_face_above_threshold` (fixture or mock session)
- [ ] `test_liveness_printed_photo_rejected`
- [ ] `test_face_recognition_same_person` — similarity >= 0.75
- [ ] `test_face_recognition_different_persons` — similarity < 0.65
- [ ] `test_cosine_similarity_identical_vectors` → 1.0

## Validation

- [ ] `cd services/biometric && cargo test`
- [ ] `cargo clippy` — no warnings (or documented allows)
- [ ] Sessions not reloaded per request (code review)
- [ ] Add `models/` to `.gitignore` if not already

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
