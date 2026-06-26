# Tasks — PunchViewModel

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh domain mobile state tdd`
- [ ] Kotlin 2.0+ and KMP plugin available

## Scaffold

- [ ] Create `mobile/shared/` KMP module (commonMain, androidMain, iosMain)
- [ ] Package `com.openpresence.punch`
- [ ] Add Koin, Ktor, Coroutines dependencies

## Domain & ports

- [ ] Define `PunchType`, `PunchErrorCode` enums (match glossary)
- [ ] Define port interfaces: `PunchRepository`, `BiometricProcessor`, `GeofenceValidator`, `DeviceIntegrityChecker`

## ViewModel

- [ ] Implement sealed class `PunchState` per `docs/MOBILE-FLOWS.md`
- [ ] Implement `PunchViewModel.startPunch(type: PunchType)`
- [ ] Flow: device check → location → geofence → camera → liveness → submit → result
- [ ] Implement `handleOfflinePunch` — PENDING status stub

## DI

- [ ] Koin `punchModule` with ViewModel and factory bindings

## Tests

- [ ] Test `Idle` → `CheckingDevice` on `startPunch`
- [ ] Test `OutOfGeofence` when validator returns false
- [ ] Test `Success` when repository returns VALID
- [ ] Test `Error` on network failure

## Validation

- [ ] `./gradlew :shared:check` (or project equivalent) passes
- [ ] No business logic duplicated from server domain (geofence rules stay server-side for authority)

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
