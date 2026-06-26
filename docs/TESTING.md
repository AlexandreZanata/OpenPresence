# Testing Strategy

TDD pyramid: **75% unit / 20% integration / 5% E2E**. Domain layer ≥ **90% coverage**.

## Test layout

```
services/attendance/
  internal/domain/...     *_test.go     # unit — no DB
services/biometric/
  src/pipeline/           #[cfg(test)]  # unit — math, preprocess
  tests/                  grpc_integration.rs
  internal/application/   *_test.go     # unit with mocked ports
  tests/integration/      *_test.go     # real DB (testcontainers)
mobile/shared/
  commonTest/                          # KMP shared tests
```

## Domain tests — PunchRecord (Go)

Implemented in `services/attendance/internal/domain/punch/`:

| Test | Rule |
|------|------|
| `TestValidator_BR010_ValidPunch` | BR-010 all criteria → VALID |
| `TestValidator_BR010_RejectLiveness` | liveness below threshold → REJECTED |
| `TestValidator_BR010_RejectMockGPS` | mocked GPS → REJECTED |
| `TestValidator_BR010_RejectOutOfGeofence` | outside zone → REJECTED |
| `TestValidator_BR014_InvalidSequence` | BR-014 invalid sequence → REJECTED |
| `TestValidator_BR014_ValidSequence` | BR-014 CLOCK_IN → BREAK_START |
| `TestValidator_BR015_ServerTimeOfficial` | BR-015 punchedAt = server time |
| `TestValidator_AntiDuplicateWithin60Seconds` | duplicate within 60s → REJECTED |
| `TestValidator_ClockManipulationOver300Seconds` | \|device−server\| > 300s → REJECTED |
| `TestValidator_BR011_OfflineSyncExpired` | past offline TTL → DISCARDED |
| `TestValidator_BR011_OfflineSyncWithinTTL` | BR-011 sync within TTL → VALID |
| `TestTransition_PendingToValid` | PunchStatus state machine |
| `TestTransition_ValidIsTerminal` | VALID is terminal |

Manual verification:

```bash
./scripts/verify-punch.sh
```

## Domain tests — Geofence (Go)

| Test | Case |
|------|------|
| `TestIsInsideCircle_CenterPoint` | true |
| `TestIsInsideCircle_OnBoundary` | true |
| `TestIsInsideCircle_JustOutside` | false |
| `TestIsInsideCircle_WithDeviation` | true within buffer |
| `TestIsInsidePolygon_InsideConvex` | true |
| `TestIsInsidePolygon_ConcaveShape` | true |
| `TestHaversineDistance_KnownPair` | delta < 1m |

Write geofence tests **before** implementation (TDD). See AGENT Task 02.

```bash
./scripts/verify-geofence.sh
```

Implemented in `services/attendance/internal/domain/geofence/`.

## Domain tests — WorkSchedule (Go)

Implemented in `services/attendance/internal/domain/workforce/`:

| Test | Rule |
|------|------|
| `TestCalculateWorkedMinutes_BR030_PrivateOfficeWithBreak` | BR-030 in/out minus breaks |
| `TestCalculateWorkedMinutes_BR030_NursingNightShift` | BR-030 12×36 night shift |
| `TestCalculateLateness_BR031_TenMinutesInFiveTolerance` | BR-031 lateness after tolerance |
| `TestCalculateOvertime_BR032_AfterEndPlusTolerance` | BR-032 overtime after end + tolerance |
| `TestCalculateOvertime_BR032_DisabledWhenPolicyOff` | BR-032 no overtime when disabled |
| `TestEvaluateWindows_BR033_Shift12x36CrossesMidnight` | BR-033 night window crosses midnight |
| `TestEvaluateWindows_BR033_SplitShiftIndependentWindows` | BR-033 split shift windows |
| `TestUpdateTimeBank_BR034_AccumulatesWithTimeBankPolicy` | BR-034 cumulative time bank |
| `TestUpdateTimeBank_BR034_StandardPolicyNoAccrual` | BR-034 standard policy skips bank |

Manual verification:

```bash
./scripts/verify-work-schedule.sh
```

## Biometric tests (Rust)

Implemented in `services/biometric/`:

| Test | Expectation |
|------|-------------|
| `test_liveness_real_face_above_threshold` | is_live = true |
| `test_liveness_printed_photo_rejected` | is_live = false |
| `test_face_recognition_same_person` | similarity >= 0.75 |
| `test_face_recognition_different_persons` | similarity < 0.65 |
| `test_cosine_similarity_identical_vectors` | similarity = 1.0 |
| `grpc_verify_punch_and_enroll_roundtrip` | integration — live gRPC server |

Manual verification (starts server, curls health, checks gRPC port):

```bash
./scripts/verify-biometric.sh
```

## Integration tests — Punch API

| Test | Scope |
|------|-------|
| `TestPunchAPI_FullValidFlow` | E2E with mocked biometric gRPC |
| `TestPunchAPI_Unauthorized` | 401 without JWT |
| `TestPunchAPI_RateLimit` | 429 after burst |
| `TestPunchAPI_OfflineSync_BulkPunches` | 50 offline punches sync |

## PostgreSQL RLS tests (Go)

Implemented in `services/attendance/internal/infrastructure/postgres/` (`//go:build integration`):

| Test | Expectation |
|------|-------------|
| `TestMigrations_ApplyOnEmptyDB` | migrations 001–006 on empty Postgres 16 |
| `TestRLS_TenantCannotReadOtherTenantEmployee` | cross-tenant `GetEmployee` returns nil |
| `TestRLS_PunchRecordsAndEmbeddingsIsolated` | tenant-scoped row counts on related tables |

Uses testcontainers (`attendance_app` role, no `BYPASSRLS`). Manual verification:

```bash
./scripts/verify-rls.sh
```

## Mobile tests (KMP)

Implemented in `mobile/shared/src/commonTest/`:

| Test | Expectation |
|------|-------------|
| `startPunch_emitsCheckingDeviceFirst` | `CheckingDevice` before device port runs |
| `startPunch_outOfGeofence_whenValidatorFails` | `OutOfGeofence` with distance |
| `startPunch_success_whenRepositoryReturnsValid` | `Success` + `VALID` |
| `startPunch_error_onNetworkFailure` | `Error` + `NETWORK` |
| `handleOfflinePunch_queuesPendingResult` | offline queue + `PENDING` |

Manual verification (Gradle check + PunchState contract):

```bash
./scripts/verify-mobile.sh
```

## CI gates

- All unit tests pass
- Domain coverage ≥ 90%
- `golangci-lint`, `cargo clippy`, `ktlint` (when configured)
- No commit without running affected test suite
