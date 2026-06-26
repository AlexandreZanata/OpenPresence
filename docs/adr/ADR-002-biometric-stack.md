# ADR-002: Self-Hosted Biometric Stack (AuraFace + MiniFASNet)

**Status:** Accepted  
**Date:** 2026-06-26  
**Deciders:** Product / Tech lead

## Context

Face recognition requires commercially usable pretrained weights. InsightFace code is MIT but buffalo_l weights are non-commercial only.

## Decision

- **Recognition:** AuraFace (MIT weights, 512-dim ArcFace-style embeddings)
- **Liveness:** MiniFASNet V2 + V1SE ensemble (Apache 2.0)
- **Detection:** RetinaFace ONNX (MIT)
- **Runtime:** ONNX Runtime (Rust `ort`, mobile ONNX Runtime)
- **Storage:** pgvector cosine search in PostgreSQL

## Consequences

### Positive

- Full commercial use without extra license fees
- Self-hosted satisfies LGPD Art. 11
- Offline-capable on-device inference

### Negative

- 3D masks and deepfakes not fully covered
- Model files (~tens of MB) distributed out-of-band, not in git

## Alternatives considered

| Option | Rejected because |
|--------|------------------|
| InsightFace buffalo_l | Non-commercial weights |
| AWS Rekognition / Azure Face | Data residency, cost, offline |
| Single liveness model | Lower accuracy than ensemble |

See [BIOMETRICS.md](../BIOMETRICS.md).
