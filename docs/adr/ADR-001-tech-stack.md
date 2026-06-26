# ADR-001: Backend and Mobile Technology Stack

**Status:** Accepted  
**Date:** 2026-06-26  
**Deciders:** Product / Tech lead

## Context

OpenPresence requires low-latency punch processing with on-device biometrics, self-hosted sensitive data (LGPD), and cross-platform mobile (Android + iOS).

## Decision

- **Backend services:** Go + Fiber v3 for API/auth/attendance/org/workforce; Rust + Axum for biometric ONNX pipeline.
- **Mobile:** Kotlin Multiplatform + Compose Multiplatform.
- **Data:** PostgreSQL 16 + pgvector + TimescaleDB, Redis (Valkey), NATS JetStream.

## Consequences

### Positive

- Go: proven throughput for gateway and CRUD services
- Rust: memory safety for biometric inference
- KMP: shared domain and use cases across mobile platforms

### Negative

- Two backend languages increase hiring and CI complexity
- KMP iOS toolchain adds setup friction

## Alternatives considered

| Option | Rejected because |
|--------|------------------|
| Single language (Go only) | ONNX/biometric ergonomics weaker than Rust |
| Flutter mobile | Team Kotlin expertise; Compose Multiplatform maturity |
| Cloud biometric APIs | LGPD and offline requirements |
