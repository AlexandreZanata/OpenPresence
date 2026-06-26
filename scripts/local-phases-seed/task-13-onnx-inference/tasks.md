# Tasks — ONNX inference pipeline

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh biometric onnx security tdd`
- [ ] Task 12 models downloadable (`./scripts/download-models.sh`)

## ONNX integration

- [ ] Enable `ort` in `Cargo.toml` with `onnx` feature flag
- [ ] `FaceProcessor::from_models(path)` loads RetinaFace, MiniFASNet×2, AuraFace once (`Arc`)
- [ ] `detect_face` → bounding box + landmarks
- [ ] `liveness_score` → ensemble ≥ 0.80 threshold
- [ ] `embed` → `Vec<f32>` length 512

## Image processing

- [ ] JPEG/WebP decode → preprocess 80×80 BGR (liveness), 112×112 RGB (recognition)
- [ ] Align/warp using 5-point landmarks
- [ ] Unit tests with fixed tensor fixtures (no ONNX in unit tests if slow)

## Stub fallback

- [ ] `BIOMETRIC_USE_STUB=true` or missing model → existing stub behavior
- [ ] Log clear message at startup which mode is active

## Validation

- [ ] `cargo test` (unit)
- [ ] `./scripts/verify-biometric.sh` with stub
- [ ] Manual: `./scripts/verify-biometric.sh` with `ONNX_MODELS_PATH=./models` after download
- [ ] Update `docs/TESTING.md` biometric section

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
