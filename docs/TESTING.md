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

## Domain tests — Fraud detection (Go)

Implemented in `services/attendance/internal/domain/fraud/`:

| Test | Rule |
|------|------|
| `TestEvaluate_AllGlossaryFraudTypes/*` | each `FraudType` in glossary |
| `TestEvaluate_MockGPS_HighSeverity` | MOCK_GPS → HIGH |
| `TestEvaluate_ClockManipulation_MediumOver300s` | clock skew > 300s |
| `TestEvaluate_ClockManipulation_CriticalOver30Min` | CRITICAL → REJECTED |
| `TestEvaluate_ImpossibleSpeed_Critical` | > 600 km/h → CRITICAL |
| `TestEvaluate_DuplicatePunch_Within60s` | duplicate within 60s |
| `TestEvaluate_GPSLowAccuracy_LowAcceptWithFlag` | BR-022 accept with flag |
| `TestEvaluate_BR012_SuspiciousWhenNotCritical` | BR-012 non-critical → SUSPICIOUS |
| `TestDeviceLockoutTracker_BR013_ThreeRejectsInTenMinutes` | BR-013 device lockout |

Manual verification:

```bash
./scripts/verify-fraud.sh
./scripts/verify-fraud-e2e.sh
```

## Fraud E2E tests — SubmitPunch integration (Go)

Implemented in `services/attendance/internal/application/punch/submit_punch_fraud_e2e_integration_test.go`:

| Test | Rule |
|------|------|
| `TestSubmitPunch_E2E_Fraud_BR012_VPN_SUSPICIOUSInDB` | BR-012 VPN → SUSPICIOUS persisted |
| `TestSubmitPunch_E2E_Fraud_BR012_CriticalClock_REJECTED` | BR-012 critical fraud → REJECTED |
| `TestSubmitPunch_E2E_Fraud_BR013_DeviceLockoutAfterThreeRejects` | BR-013 lockout after 3 rejects |

`SubmitPunchHandler` wires `DeviceLockoutTracker` via optional `Lockout` field and `DeviceID` on the command.

## Domain tests — Hierarchy authorization (Go)

Implemented in `services/attendance/internal/domain/organization/` and `internal/application/authorization/`:

| Test | Expectation |
|------|-------------|
| `TestCanApprovePunch_PublicHealthManagerApprovesNurse` | Health manager approves nursing subtree |
| `TestCanApprovePunch_PublicHealthManagerRejectsEducation` | Cross-secretariat denied |
| `TestCanApprovePunch_PrivateSalesManagerApprovesInsideSales` | Sales manager approves team |
| `TestCanApprovePunch_PrivateSalesManagerRejectsIT` | Sibling department denied |
| `TestCanReadPunch_AuditorAllowedWriteDenied` | AUDITOR read-only |
| `TestCanApprovePunch_CrossTenantDenied` | Cross-tenant always denied |
| `TestPunchAuthorizationService_*` | Application layer with mocked `OrgTreeReader` |

Manual verification:

```bash
./scripts/verify-authorization.sh
```

## Authorization E2E tests — ApprovePunch integration (Go)

`AuthorizePunchApprovalHandler` wires Postgres `EmployeeRepository`, placement resolution, and `PunchAuthorizationService` ABAC checks.

Implemented in `services/attendance/internal/application/authorization/approve_punch_integration_test.go`:

| Test | Expectation |
|------|-------------|
| `TestPunchAuthorization_E2E_ManagerApprovesNurseInSubtree` | Health manager approves nurse in nursing subtree |
| `TestPunchAuthorization_E2E_ManagerRejectsEducationEmployee` | Cross-secretariat approval denied |
| `TestPunchAuthorization_E2E_AuditorWriteDenied` | Auditor read allowed, write denied |
| `TestPunchAuthorization_E2E_CrossTenantActorDenied` | Cross-tenant actor rejected |

Manual verification (domain unit + Postgres integration):

```bash
./scripts/verify-authorization-e2e.sh
```

## Model download (ONNX)

| Script | Expectation |
|--------|-------------|
| `./scripts/download-models.sh` | Fetches all models from `models/MANIFEST.json` |
| `./scripts/verify-models.sh` | All files present + SHA-256 match manifest |

Models are gitignored; manifest is committed at `models/MANIFEST.json`.

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

## Geofence E2E tests — SubmitPunch integration (Go)

Implemented in `services/attendance/internal/application/punch/submit_punch_geofence_e2e_integration_test.go`:

| Test | Rule |
|------|------|
| `TestSubmitPunch_E2E_Geofence_BR020_CircleInside_VALID` | BR-020 circle inside |
| `TestSubmitPunch_E2E_Geofence_BR020_CircleOutside_REJECTED` | BR-020 circle outside |
| `TestSubmitPunch_E2E_Geofence_BR021_PolygonInside_VALID` | BR-021 polygon inside |
| `TestSubmitPunch_E2E_Geofence_BR022_LowAccuracyFlag_VALID` | BR-022 GPS_LOW_ACCURACY flag, still VALID |
| `TestSubmitPunch_E2E_Geofence_BR023_AnyAssignedZone_VALID` | BR-023 match any assigned zone |
| `TestSubmitPunch_E2E_Geofence_BR024_ExpiredZoneIgnored_REJECTED` | BR-024 expired zone ignored |

Manual verification (domain + Postgres integration):

```bash
./scripts/verify-geofence-e2e.sh
```

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
./scripts/verify-work-schedule-e2e.sh
```

## Work schedule E2E tests — CalculateDayAttendance (Go)

Application use case: `internal/application/attendance/CalculateDayAttendanceHandler` loads VALID punches from Postgres and applies BR-030–034.

| Test | Rule |
|------|------|
| `TestCalculateDay_E2E_BR030_WorkedMinutesFromDB` | BR-030 worked minutes from seeded punches |
| `TestCalculateDay_E2E_BR031_LatenessFromDB` | BR-031 lateness after tolerance |
| `TestCalculateDay_E2E_BR032_OvertimeFromDB` | BR-032 overtime when policy allows |
| `TestCalculateDay_E2E_BR033_Shift12x36Windows` | BR-033 12×36 window resolution |
| `TestCalculateDay_E2E_BR034_TimeBankFromDB` | BR-034 cumulative time bank |

## Biometric tests (Rust)

Implemented in `services/biometric/`:

| Test | Expectation |
|------|-------------|
| `test_liveness_real_face_above_threshold` | is_live = true |
| `test_liveness_printed_photo_rejected` | is_live = false |
| `test_face_recognition_same_person` | similarity >= 0.75 |
| `test_face_recognition_different_persons` | similarity < 0.65 |
| `test_cosine_similarity_identical_vectors` | similarity = 1.0 |
| `real_inference_returns_512_embedding` | ONNX — 512-dim embedding (ignored unless models present) |
| `grpc_verify_punch_and_enroll_roundtrip` | integration — live gRPC server |

Manual verification (starts server, curls health, checks gRPC port):

```bash
# Stub mode (CI default — no models required)
./scripts/verify-biometric.sh

# ONNX mode (after ./scripts/download-models.sh)
ONNX_MODELS_PATH=./models ./scripts/verify-biometric.sh
```

Build with ONNX Runtime: `cargo test --features onnx` in `services/biometric/`.

## Enrollment E2E tests (Rust gRPC)

Implemented in `services/biometric/tests/enrollment_e2e.rs`:

| Test | Rule |
|------|------|
| `enrollment_e2e_br001_three_angles_success` | BR-001 FRONTAL, LEFT_15, RIGHT_15 |
| `enrollment_e2e_br002_liveness_fail_rejected` | BR-002 liveness < 0.85 → no embedding |
| `enrollment_e2e_br003_low_quality_rejected` | BR-003 quality < 0.7 → no embedding |

Manual verification (integration tests + live server; grpcurl optional):

```bash
./scripts/verify-enrollment.sh
ONNX_MODELS_PATH=./models ./scripts/verify-enrollment.sh
```

Fixtures for grpcurl: `services/biometric/tests/fixtures/*.jpg`.

## Application tests — SubmitPunch (Go)

Implemented in `services/attendance/internal/application/punch/`:

| Test | Expectation |
|------|-------------|
| `TestSubmitPunchHandler_HappyPath_VALID` | unit — VALID punch |
| `TestSubmitPunchHandler_OutOfGeofence_REJECTED` | unit — REJECTED, no VALID row |
| `TestSubmitPunchHandler_InvalidSequence_REJECTED` | unit — sequence BR-014 |
| `TestSubmitPunch_Integration_HappyPath_VALIDInDB` | integration — VALID in Postgres + RLS |
| `TestSubmitPunch_Integration_CrossTenant_Rejected` | integration — employee not visible |
| `TestSubmitPunch_Integration_InvalidSequence_REJECTED` | integration — one VALID only |
| `TestSubmitPunch_Integration_OutOfGeofence_REJECTED` | integration — geofence rejection |
| `TestSubmitPunch_E2E_BR010_LowLiveness_REJECTED` | E2E integration — liveness fail |
| `TestSubmitPunch_E2E_BR010_MockGPS_REJECTED` | E2E integration — mock GPS |
| `TestSubmitPunch_E2E_BR010_ClockSkew_REJECTED` | E2E integration — clock skew > 300s |
| `TestSubmitPunch_E2E_BR010_Duplicate_REJECTED` | E2E integration — duplicate within 60s |
| `TestSubmitPunch_E2E_BR011_OfflineExpired_DISCARDED` | E2E integration — offline TTL expired |
| `TestSubmitPunch_E2E_BR011_OfflineWithinTTL_VALID` | E2E integration — offline within 8h |
| `TestSubmitPunch_E2E_BR014_InvalidSequence_REJECTED` | E2E integration — sequence BR-014 |
| `TestSubmitPunch_E2E_BR015_ServerTimeOfficial` | E2E integration — punchedAt = server time |

Manual verification:

```bash
./scripts/verify-punch-usecase.sh
go test -tags=integration ./services/attendance/internal/application/punch/... -run E2E
```

Integration tests require Docker (testcontainers Postgres 16).

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

## Master business-rules runner (local)

After `./scripts/setup-local-ai.sh`, run all repo verification scripts in one pass:

```bash
./.local/scripts/verify-all-business-rules.sh           # unit + integration (Docker)
./.local/scripts/verify-all-business-rules.sh --quick   # domain/unit only (no Docker)
./.local/scripts/verify-all-business-rules.sh --e2e     # full + per-phase E2E run.sh stubs
```

Coverage matrix (local): `.local/phases/e2e-testing/BUSINESS-RULES-COVERAGE.md`.

Prerequisites: Go ≥ 1.22, Rust/Cargo (biometric), Docker (RLS + SubmitPunch integration), Gradle/JDK (mobile). ONNX models optional via `./scripts/download-models.sh`.
