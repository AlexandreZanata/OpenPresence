# Biometrics — Open Source Stack

Biometric data is **sensitive** (LGPD Art. 11). Processing stays self-hosted; no cloud biometric APIs in v1.

## Why self-hosted

- Regulatory: biometric data must not leave tenant infrastructure without explicit consent.
- Cost: cloud APIs can be ~20× more expensive at scale.
- Offline: on-device recognition with server confirmation.

## Commercially free components

### Liveness — MiniFASNet V2 + V1SE (Apache 2.0)

| Property | Value |
|----------|-------|
| Repository | [Silent-Face-Anti-Spoofing](https://github.com/minivision-ai/Silent-Face-Anti-Spoofing) |
| Format | ONNX (~4 MB ensemble) |
| Accuracy | ~98% (CelebA-Spoof) |
| Mobile latency | < 50 ms CPU (Cortex-A55) |
| Classes | live, print-attack, replay-attack |

Ensemble pipeline:

```
frame (80×80 BGR crop)
  → MiniFASNetV2 @ scale 2.7  → softmax[3]
  → MiniFASNetV1SE @ scale 4.0 → softmax[3]
  → average → argmax → {REAL | SPOOF}
```

**Known limitation:** printed photos and screen replay. High-quality 3D masks and deepfakes are out of scope; combine with challenge-response for high-security deployments.

### Face recognition — AuraFace (MIT)

| Property | Value |
|----------|-------|
| Repository | [auraface-project/auraface](https://github.com/auraface-project/auraface) |
| Architecture | ArcFace-style, 512-dim embedding |
| Benchmark | LFW 99.82% |

**Why not InsightFace buffalo_l?** InsightFace code is MIT but pretrained weights are **non-commercial only**. AuraFace publishes MIT-licensed weights with equivalent performance.

### Face detection — RetinaFace (ONNX, MIT)

Locates and aligns face before embedding. ~15 ms CPU, ~3 ms GPU.

## Rust biometric pipeline

```
CameraFrame (JPEG/WebP)
  ↓ RetinaFace → BoundingBox + 5-point landmarks
  ↓ Affine warp → 112×112 RGB (AuraFace)
  ↓ Affine warp → 80×80 BGR (MiniFASNet)
  ↓ [parallel]
      MiniFASNet ensemble → liveness_score (0.0–1.0)
      AuraFace → embedding Vec<f32; 512>
  ↓ cosine_similarity(embedding, registered_embeddings[])
  ↓ PunchResult { matched, liveness, confidence }
```

**Note:** pgvector search runs in Go (Attendance Service). Biometric Service returns embedding; Go queries the database.

## Development without ONNX weights

`services/biometric` runs in **stub mode** when `ONNX_MODELS_PATH` is unset (`BIOMETRIC_USE_STUB=true`):

- Pipeline math and preprocess functions are exercised with real JPEG decode
- gRPC `VerifyPunch` / `EnrollFace` return deterministic embeddings for integration tests
- Full ONNX inference: `cargo build --features onnx` with models in `models/`

```bash
./scripts/verify-biometric.sh
```

## Model files (not in git)

Weights are **not** committed. Download with checksum verification:

```bash
./scripts/download-models.sh
./scripts/verify-models.sh
```

Manifest: `models/MANIFEST.json` (URLs, SHA-256, license notes).

Expected layout after download:

```
models/
├── MANIFEST.json          # committed — URLs + hashes
├── auraface.onnx          # gitignored — recognition (AuraFace glintr100)
├── retinaface.onnx        # gitignored — detection (SCRFD from AuraFace bundle)
├── minifasnet_v2.onnx     # gitignored — liveness scale 2.7
└── minifasnet_v1se.onnx   # gitignored — liveness scale 4.0
```

Set `ONNX_MODELS_PATH=./models` and `BIOMETRIC_USE_STUB=false` for ONNX inference (task 13).

## License summary

| Component | License | Commercial |
|-----------|---------|------------|
| AuraFace | MIT | Yes |
| RetinaFace | MIT | Yes |
| MiniFASNet V2/V1SE | Apache 2.0 | Yes |
| ONNX Runtime | MIT | Yes |
| InsightFace weights (buffalo_l) | Non-commercial | No — use AuraFace |

See [ADR-002](adr/ADR-002-biometric-stack.md).
