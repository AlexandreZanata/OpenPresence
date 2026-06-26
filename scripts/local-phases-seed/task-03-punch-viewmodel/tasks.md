# Tasks — PunchViewModel

## Preparation

- [x] Read [README.md](README.md) and [official_source.md](official_source.md)
- [x] Run `./agent-harness/resolve-rules.sh domain mobile state tdd`
- [x] Kotlin 2.0+ and KMP plugin available

## Scaffold

- [x] Create `mobile/shared/` KMP module (commonMain, androidMain, iosMain)
- [x] Package `com.openpresence.punch`
- [x] Add Koin, Ktor, Coroutines dependencies

## Domain & ports

- [x] Define `PunchType`, `PunchErrorCode` enums (match glossary)
- [x] Define port interfaces: `PunchRepository`, `BiometricProcessor`, `GeofenceValidator`, `DeviceIntegrityChecker`

## ViewModel

- [x] Implement sealed class `PunchState` per `docs/MOBILE-FLOWS.md`
- [x] Implement `PunchViewModel.startPunch(type: PunchType)`
- [x] Flow: device check → location → geofence → camera → liveness → submit → result
- [x] Implement `handleOfflinePunch` — PENDING status stub

## DI

- [x] Koin `punchModule` with ViewModel and factory bindings

## Tests

- [x] Test `Idle` → `CheckingDevice` on `startPunch`
- [x] Test `OutOfGeofence` when validator returns false
- [x] Test `Success` when repository returns VALID
- [x] Test `Error` on network failure

## Validation

- [x] `./gradlew :shared:check` (or project equivalent) passes
- [x] No business logic duplicated from server domain (geofence rules stay server-side for authority)

## Completion

- [x] All steps above marked `[x]`
- [x] Update `.local/phases/README.md` active task
