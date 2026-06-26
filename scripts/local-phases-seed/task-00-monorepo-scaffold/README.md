# Task 00 — Monorepo scaffold

**Status:** done  
**Phase ID:** task-00-monorepo-scaffold

## Goal

Create the minimal monorepo layout so domain code can live under `services/attendance/` per [ADR-003](../../../docs/adr/ADR-003-monorepo-structure.md).

## Scope

**In scope:**

- `services/attendance/` Go module with `internal/domain/` layout
- Root `go.work` (optional) or documented single-module start
- Placeholder `infra/docker-compose.yml` skeleton (no full stack yet)

**Out of scope:**

- HTTP handlers, database, biometric service
- CI pipelines (later phase)

## Acceptance

- `go test ./...` runs (even if only empty packages)
- Folder layout matches [ARCHITECTURE.md](../../../docs/ARCHITECTURE.md)

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
