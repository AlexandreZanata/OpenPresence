# Business Rules

All rules use **GIVEN / WHEN / THEN**. IDs are stable references for tests and use cases.

## Enrollment (BR-001 – BR-006)

### BR-001 — Minimum enrollment angles

- **GIVEN** an administrator starts biometric enrollment for an employee
- **WHEN** fewer than 3 face captures are submitted
- **THEN** enrollment is rejected; required angles: FRONTAL, LEFT_15, RIGHT_15

### BR-002 — Liveness on enrollment

- **GIVEN** a face capture during enrollment
- **WHEN** `liveness_score < 0.85`
- **THEN** the capture is rejected

### BR-003 — Image quality on enrollment

- **GIVEN** a face capture during enrollment
- **WHEN** image quality (sharpness + lighting) `< 0.7`
- **THEN** the capture is rejected; no embedding stored

### BR-004 — Administrator-only enrollment

- **GIVEN** an enrollment request
- **WHEN** the actor is not an authenticated administrator
- **THEN** enrollment is rejected; employees cannot self-enroll remotely

### BR-005 — Maximum embeddings

- **GIVEN** an employee biometric profile
- **WHEN** storing a new embedding would exceed 5 active embeddings
- **THEN** oldest active embedding is rotated out (or enrollment rejected per policy)

### BR-006 — Deactivation retention

- **GIVEN** an employee is deactivated
- **WHEN** status changes to INACTIVE or SUSPENDED
- **THEN** embeddings are marked INACTIVE, not deleted (LGPD Art. 16 audit retention)

## Punch (BR-010 – BR-015)

### BR-010 — Valid punch criteria (all required)

- **GIVEN** a punch submission
- **WHEN** all conditions hold simultaneously:
  - `liveness_score >= 0.80`
  - `recognition_confidence >= 0.75`
  - `gps.isMocked == false`
  - coordinate inside assigned geofence
  - `|device_time - server_time| <= 300` seconds
  - no valid punch in last 60 seconds (anti-duplicate)
- **THEN** `PunchRecord.status = VALID`

### BR-011 — Offline punch TTL

- **GIVEN** an offline punch without connectivity
- **WHEN** sync occurs within `AttendancePolicy.offlineSyncMaxAge` (default 8h)
- **THEN** punch is processed; otherwise DISCARDED with audit log

### BR-012 — Suspicious punch

- **GIVEN** a punch with any FraudFlag
- **WHEN** severity does not mandate auto-reject
- **THEN** record is persisted with `status = SUSPICIOUS`; manager review required

### BR-013 — Device lockout

- **GIVEN** 3 consecutive REJECTED attempts within 10 minutes
- **WHEN** the third rejection is recorded
- **THEN** device blocked 30 minutes; manager alert sent

### BR-014 — Punch type sequence

- **GIVEN** an employee's punch history for the workday
- **WHEN** next punch type breaks order CLOCK_IN → BREAK_START → BREAK_END → CLOCK_OUT
- **THEN** punch is REJECTED

### BR-015 — Official timestamp

- **GIVEN** any accepted punch
- **WHEN** persisted
- **THEN** `punchedAt` uses server time; `deviceTime` stored for audit only

## Geofence (BR-020 – BR-024)

### BR-020 — Circle geofence

- **GIVEN** zone type CIRCLE
- **WHEN** `haversine(coordinate, center) <= radius + allowedDeviation`
- **THEN** inside zone

### BR-021 — Polygon geofence

- **GIVEN** zone type POLYGON
- **WHEN** point-in-polygon (ray casting) with `allowedDeviation` buffer
- **THEN** inside zone

### BR-022 — Low GPS accuracy

- **GIVEN** `gps.accuracy > allowedDeviation * 2`
- **WHEN** punch otherwise valid
- **THEN** accept with FraudFlag `GPS_LOW_ACCURACY`, severity LOW

### BR-023 — Multiple zones

- **GIVEN** employee assigned to multiple geofences
- **WHEN** inside any one assigned zone
- **THEN** geofence check passes

### BR-024 — Temporary zones

- **GIVEN** geofence with `validFrom` / `validUntil`
- **WHEN** punch outside validity window
- **THEN** zone is ignored for matching

## Work schedule (BR-030 – BR-034)

### BR-030 — Worked time calculation

- **GIVEN** completed punch pairs for a day
- **WHEN** calculating total worked time
- **THEN** `Σ(CLOCK_OUT - CLOCK_IN) - Σ(BREAK_END - BREAK_START)`

### BR-031 — Lateness

- **GIVEN** employee WorkSchedule with `scheduled_start`
- **WHEN** CLOCK_IN after scheduled start + tolerance
- **THEN** lateness minutes recorded

### BR-032 — Overtime

- **GIVEN** CLOCK_OUT after `scheduled_end + toleranceMinutes`
- **WHEN** policy allows overtime
- **THEN** overtime minutes recorded

### BR-033 — Split shifts

- **GIVEN** policy with multiple work windows per day
- **WHEN** calculating attendance
- **THEN** each window evaluated independently (12×36, split shifts supported)

### BR-034 — Time bank

- **GIVEN** organizational overtime policy
- **WHEN** period closes (week / biweek / month per policy)
- **THEN** time bank balance updated cumulatively
