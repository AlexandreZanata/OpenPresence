# Official sources — Biometric service

## Repository documentation

| Document | Path |
|----------|------|
| Biometrics stack | `docs/BIOMETRICS.md` |
| ADR-002 | `docs/adr/ADR-002-biometric-stack.md` |
| API gRPC contract | `docs/API-CONTRACT.md` (Biometric Service section) |
| Testing (Rust tests) | `docs/TESTING.md` |
| Security (no raw embeddings in REST) | `docs/SECURITY.md` |
| Agent guide Task 01 | `docs/AGENT-IMPLEMENTATION-GUIDE.md` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh owasp security biometric onnx
```

| Rule file | Why |
|-----------|-----|
| `agent-rules/03-security/secrets-and-credentials.md` | No leaked embeddings |
| `agent-rules/07-data-management/pii-and-data-retention.md` | Biometric = sensitive |
| `agent-rules/04-testing/tdd.md` | Unit tests per pipeline stage |

## Business rules

| ID | Summary |
|----|---------|
| BR-002 | Liveness >= 0.85 on enrollment |
| BR-010 | Liveness >= 0.80 on punch |

## External references

| Topic | URL |
|-------|-----|
| AuraFace | https://github.com/auraface-project/auraface |
| MiniFASNet (Silent-Face-Anti-Spoofing) | https://github.com/minivision-ai/Silent-Face-Anti-Spoofing |
| ort (ONNX Runtime Rust) | https://github.com/pykeio/ort |
| tonic gRPC | https://github.com/hyperium/tonic |

## Glossary terms

- `BiometricResult`, `FaceEmbedding`, `BiometricProfile`
