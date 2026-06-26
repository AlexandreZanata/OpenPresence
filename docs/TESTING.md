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

| Test | Rule |
|------|------|
| `TestPunchRecord_ValidPunch` | BR-010 all criteria → VALID |
| `TestPunchRecord_RejectMockGPS` | mocked GPS → SUSPICIOUS + MOCK_GPS |
| `TestPunchRecord_RejectOutOfGeofence` | outside zone → REJECTED |
| `TestPunchRecord_ClockManipulationDetection` | 10min delta → CLOCK_MANIPULATION |
| `TestPunchRecord_ImpossibleSpeed` | 500km in 30s → CRITICAL |
| `TestPunchRecord_OfflineSync_Expired` | past TTL → DISCARDED + audit |

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
