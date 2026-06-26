# Use Case: UC-003 — Approve Suspicious Punch

## Metadata

| Field | Value |
|-------|-------|
| ID | UC-003 |
| Actor | MANAGER |
| Status | Accepted |

## Preconditions

- PunchRecord exists with status SUSPICIOUS
- Employee belongs to manager's org subtree (ABAC)
- Manager authenticated

## Main flow

1. Manager opens suspicious punch list (`GET /v1/attendance/suspicious`)
2. Reviews punch details, fraud flags, GPS map, biometric scores (no raw embedding)
3. Manager approves via `PATCH /v1/attendance/{id}/approve`
4. Status changes SUSPICIOUS → VALID
5. `PunchApproved` event published
6. Employee notified; punch counts for payroll

## Alternate flows

### AF-1: Reject

- **When:** Manager determines fraud or invalid punch
- **Then:** `PATCH .../reject` with reason → REJECTED; `PunchRejected` event

## Business rules applied

| Rule ID | Description |
|---------|-------------|
| BR-012 | Suspicious punch workflow |

## Domain events raised

| Event | When |
|-------|------|
| `PunchApproved` | Step 5 |
| `PunchRejected` | AF-1 |

## Authorization

| Role | Allowed |
|------|---------|
| MANAGER | Yes (subtree only) |
| ORG_ADMIN | Yes |
| EMPLOYEE | No |

## Out of scope

- Automatic approval without human review
