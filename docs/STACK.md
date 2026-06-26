# Technology Stack

Definitive stack for OpenPresence. Changes require an ADR.

## Backend

| Layer | Technology | Rationale |
|-------|------------|-----------|
| API Gateway / Auth | Go + [Fiber v3](https://github.com/gofiber/fiber) | Sub-ms latency, high throughput, low memory |
| Biometric Service | Rust + [Axum](https://github.com/tokio-rs/axum) | Memory safety, ONNX Runtime for embeddings |
| Face embedding DB | [pgvector](https://github.com/pgvector/pgvector) + PostgreSQL 16 | Cosine similarity in-database |
| Cache / Session | Redis 7 (Valkey) | Token TTL, rate limiting, punch deduplication |
| Message Queue | NATS JetStream | Punch events, offline sync, async audit |
| Time-series audit | TimescaleDB (PG extension) | Immutable punch and fraud logs |
| Observability | OpenTelemetry → Grafana + Loki + Tempo | Distributed traces (Go + Rust) |

## Mobile

| Layer | Technology |
|-------|------------|
| UI | Kotlin Multiplatform + [Compose Multiplatform](https://www.jetbrains.com/compose-multiplatform/) |
| On-device face detection | ONNX Runtime Mobile |
| Liveness | MiniFASNetV2 + MiniFASNetV1SE (ONNX, Apache 2.0) |
| Face recognition | AuraFace / OpenCV SFace (commercial-friendly weights) |
| GPS anti-spoof | Native Android/iOS APIs + server validation |
| Offline storage | SQLDelight (KMP) |
| DI | Koin Multiplatform |
| Network | Ktor Client |
| Local crypto | Kotlinx Serialization + AES-256-GCM for biometric cache |

## DevOps

```
Docker Compose (dev) → GitHub Actions CI → Helm + K8s (prod)
PostgreSQL 16 + pgvector + TimescaleDB
Redis 7 (Valkey)
NATS JetStream
Prometheus + Grafana
```

## Required tooling (development)

| Tool | Minimum version |
|------|-----------------|
| Go | 1.22+ |
| Rust | 1.78 stable |
| Kotlin | 2.0+ |
| Docker | 25+ |
| sqlx-cli | latest |
| protoc | 3.21+ |
| buf | 1.30+ |

See [BIOMETRICS.md](BIOMETRICS.md) for model licensing details.
