# ADR-003: Monorepo Structure

**Status:** Accepted  
**Date:** 2026-06-26  
**Deciders:** Tech lead

## Context

Multiple services (Go + Rust), mobile KMP, shared docs, and infra need coordinated versioning and atomic changes across API contracts.

## Decision

Single monorepo:

```
services/{api-gateway,attendance,organization,workforce,biometric}
mobile/{shared,androidApp,iosApp}
infra/
docs/
models/   # gitignored ONNX weights
```

## Consequences

### Positive

- Single PR can update API + mobile + docs
- Shared CI and commit conventions
- Easier local docker-compose for full stack

### Negative

- Larger clone size over time
- CI must path-filter per service

## Alternatives considered

| Option | Rejected because |
|--------|------------------|
| Polyrepo per service | Contract drift; harder atomic changes |
| Mobile separate repo | Duplicated API contract maintenance |

See [ARCHITECTURE.md](../ARCHITECTURE.md).
