# Use Case: UC-001 — Clock In (Happy Path)

## Metadata

| Field | Value |
|-------|-------|
| ID | UC-001 |
| Actor | Employee (mobile app) |
| Status | Accepted |

## Preconditions

- Employee status ACTIVE
- Device registered and not locked
- Biometric profile enrolled (≥ 3 angles)
- Employee inside assigned geofence
- Network available (or offline policy allows queue)

## Main flow

1. Employee opens app and selects CLOCK_IN
2. App runs device integrity checks
3. App captures GPS; validates geofence
4. App opens camera; detects face (RetinaFace on-device)
5. Liveness score ≥ 0.80; frame captured
6. App submits punch to `POST /v1/attendance/punch`
7. Attendance Service validates BR-010 criteria
8. Biometric Service confirms liveness + recognition via gRPC
9. PunchRecord persisted with status VALID, server `punchedAt`
10. App displays success with official time

## Alternate flows

### AF-1: Out of geofence

- **When:** GPS outside all assigned zones
- **Then:** App shows distance message; no submission (BR-023)

### AF-2: Offline punch

- **When:** No network at step 6
- **Then:** Save to SQLDelight PENDING; sync when online (BR-011)

## Business rules applied

| Rule ID | Description |
|---------|-------------|
| BR-010 | Valid punch criteria |
| BR-014 | CLOCK_IN sequence |
| BR-015 | Server timestamp official |

## Domain events raised

| Event | When |
|-------|------|
| `PunchRecorded` | After step 9 |

## Authorization

| Role | Allowed |
|------|---------|
| EMPLOYEE | Yes (own punch only) |

## Out of scope

- Manager approval (not needed for VALID punch)
- Break punches (see UC-003)
