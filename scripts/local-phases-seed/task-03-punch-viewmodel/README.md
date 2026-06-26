# Task 03 — PunchViewModel (Kotlin Multiplatform)

**Status:** pending  
**Phase ID:** task-03-punch-viewmodel

## Goal

Implement `PunchViewModel` with sealed `PunchState`, full punch happy path, and offline queue per `docs/MOBILE-FLOWS.md` and UC-001.

## Scope

**In scope:**

- KMP `shared` module structure
- ViewModel state machine
- Koin DI module stubs
- Unit tests for state transitions (commonTest)

**Out of scope:**

- Full CameraX / ONNX on-device integration (can mock `BiometricProcessor`)
- Production SQLDelight schema (stub repository OK for first pass)

## Acceptance

- State flow matches `docs/MOBILE-FLOWS.md`
- UC-001 main flow covered by tests or documented manual checklist
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
