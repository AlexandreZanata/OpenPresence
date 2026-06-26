# Mobile Shared (Kotlin Multiplatform)

Shared punch flow for Android and iOS (`com.openpresence.punch`).

## Modules

| Package | Role |
|---------|------|
| `domain` | PunchType, PunchStatus, models |
| `ports` | Repository and platform port interfaces |
| `presentation` | `PunchViewModel`, `PunchState` |
| `di` | Koin `punchModule` |

Geofence **rules** stay server-side (`services/attendance`); mobile only calls `GeofenceValidator` port.

## Commands

```bash
./gradlew :mobile:shared:check
./gradlew :mobile:shared:jvmTest
```

From repo root:

```bash
./scripts/verify-mobile.sh
```

## Related docs

- `docs/MOBILE-FLOWS.md`
- `docs/use-cases/UC-001-punch-clock-in.md`
