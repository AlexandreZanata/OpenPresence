# Mobile Flows

Kotlin Multiplatform app. States and transitions must match [DOMAIN-MODEL.md](DOMAIN-MODEL.md).

## Device onboarding

1. Install app
2. Enter organization invite code (QR or alphanumeric)
3. Employee auth (registration ID + temporary password)
4. Initial sync: employee data, policies, geofences
5. Configure push notifications
6. Ready to punch

## Punch happy path

```
Open app
  → Device integrity check (root, VPN, dev options)
  → Capture location (FusedLocationProvider)
  → Geofence validation
  → Open front camera (CameraX / AVFoundation)
  → On-device RetinaFace (~30 fps)
  → On-device MiniFASNet liveness (auto-capture at score >= 0.80)
  → HTTPS/mTLS submit: {frame, gps, device_integrity, device_time, punch_type}
  → Server: policy + geofence + biometric gRPC
  → Response: VALID | SUSPICIOUS | REJECTED
  → Offline: SQLDelight queue, sync every 30s when online
```

## PunchViewModel states

```
Idle → CheckingDevice → WaitingLocation → OpeningCamera
  → DetectingFace → CheckingLiveness → Submitting
  → Success | Suspicious | Error

Branches:
  DeviceWarning(flags)
  OutOfGeofence(distance)
```

## Offline sync

| State | Transition |
|-------|------------|
| OFFLINE_PENDING | → SYNCING on connectivity |
| SYNCING | → VALID on success |
| SYNCING | → DISCARDED if past offlineSyncMaxAge |
| SYNCING | → SUSPICIOUS on fraud flags |

## Implementation reference

See [AGENT-IMPLEMENTATION-GUIDE.md](AGENT-IMPLEMENTATION-GUIDE.md) — Task 03 (PunchViewModel).

Use case: [use-cases/UC-001-punch-clock-in.md](use-cases/UC-001-punch-clock-in.md)
