# Use Case: UC-002 — Biometric Enrollment

## Metadata

| Field | Value |
|-------|-------|
| ID | UC-002 |
| Actor | ORG_ADMIN or authorized administrator |
| Status | Accepted |

## Preconditions

- Employee record exists (ACTIVE or pending activation)
- Administrator authenticated with enrollment permission
- Enrollment performed in controlled environment (not remote self-service)

## Main flow

1. Administrator opens workforce enrollment for employee
2. For each required angle (FRONTAL, LEFT_15, RIGHT_15):
   - Capture frame
   - On-device liveness ≥ 0.85 (BR-002)
   - Image quality ≥ 0.7 (BR-003)
3. Submit frames to `POST /v1/employees/{id}/enroll`
4. Biometric Service extracts embeddings
5. Attendance/Workforce stores embeddings in `face_embeddings` (max 5 active, BR-005)
6. `EmployeeEnrolled` event published
7. Employee mobile app can punch after sync

## Alternate flows

### AF-1: Liveness failure

- **When:** score < 0.85 on capture
- **Then:** Reject capture; prompt retry

### AF-2: Maximum embeddings

- **When:** 5 active embeddings exist
- **Then:** Rotate oldest or reject per tenant policy (BR-005)

## Business rules applied

| Rule ID | Description |
|---------|-------------|
| BR-001 | Three angles required |
| BR-002 | Liveness on enrollment |
| BR-003 | Image quality |
| BR-004 | Administrator only |

## Domain events raised

| Event | When |
|-------|------|
| `EmployeeEnrolled` | After step 6 |

## Authorization

| Role | Allowed |
|------|---------|
| ORG_ADMIN | Yes (within org subtree) |
| SUPER_ADMIN | Yes |
| EMPLOYEE | No |

## Out of scope

- Employee self-enrollment via mobile
