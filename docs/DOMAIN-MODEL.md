# Domain Model

Bounded contexts, aggregates, and value objects. Code names must match [GLOSSARY.md](GLOSSARY.md).

## Context map

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  Organization    │────▶│   Workforce      │────▶│   Attendance     │
│  Context         │     │   Context        │     │   Context        │
│  Tenant          │     │  Employee        │     │  PunchRecord     │
│  OrgNode         │     │  WorkSchedule    │     │  GeofenceZone    │
│  AttendancePolicy│     │  BiometricProfile│     │  FraudAttempt    │
└──────────────────┘     └──────────────────┘     └──────────────────┘
         │                        │                        │
         └──────────────┬─────────┘                        │
                        │                                  │
              ┌─────────▼──────────┐           ┌───────────▼──────────┐
              │ Identity & Access  │           │ Reporting & Audit      │
              │ User, Role, Perm   │           │ AuditLog, Reports      │
              └────────────────────┘           └────────────────────────┘
```

## Attendance context

### Aggregate: PunchRecord (root)

| Field | Type | Notes |
|-------|------|-------|
| `id` | PunchRecordId | UUID v7 |
| `employeeId` | EmployeeId | Reference |
| `tenantId` | TenantId | Multi-tenant |
| `punchedAt` | Timestamp | **Server time** — official record |
| `deviceTime` | Timestamp | Audit only — clock manipulation detection |
| `location` | GpsCoordinate | VO |
| `geofenceId` | GeofenceZoneId | Matched zone |
| `biometricResult` | BiometricResult | VO |
| `fraudFlags` | FraudFlag[] | VO collection |
| `status` | PunchStatus | VALID \| SUSPICIOUS \| REJECTED |
| `type` | PunchType | CLOCK_IN \| CLOCK_OUT \| BREAK_START \| BREAK_END |

### Value objects

**BiometricResult:** `livenessScore`, `recognitionConfidence`, `faceEmbeddingHash` (SHA-256, never raw embedding), `isLive`, `isRecognized`

**GpsCoordinate:** `latitude`, `longitude`, `accuracy` (m), `altitude`, `provider` (GPS \| NETWORK \| FUSED), `isMocked`

**FraudFlag:** `type` (FraudType), `severity` (LOW \| MEDIUM \| HIGH \| CRITICAL), `detectedAt`, `metadata`

**FraudType enum:** `MOCK_GPS`, `CLOCK_MANIPULATION`, `IMPOSSIBLE_SPEED`, `LIVENESS_FAILED`, `FACE_NOT_RECOGNIZED`, `OUT_OF_GEOFENCE`, `DUPLICATE_PUNCH`, `DEVICE_ROOTED`, `VPN_DETECTED`, `GPS_LOW_ACCURACY`

## Workforce context

### Aggregate: Employee (root)

| Field | Type |
|-------|------|
| `id` | EmployeeId |
| `tenantId` | TenantId |
| `registration` | string (employee ID / tax ID) |
| `name` | PersonName (VO) |
| `biometricProfile` | BiometricProfile (entity) |
| `workSchedule` | WorkSchedule (entity) |
| `assignedLocations` | GeofenceZoneId[] |
| `status` | ACTIVE \| INACTIVE \| SUSPENDED |
| `hierarchy` | OrganizationNode ref |

### Entity: BiometricProfile

`faceEmbeddings[]` (max 5 active), `enrolledAt`, `lastUpdatedAt`, `enrolledBy` (UserId)

### Value object: FaceEmbedding

`vector` (512-dim), `capturedAt`, `quality` (0–1), `angle` (FRONTAL \| LEFT_15 \| RIGHT_15)

## Organization context

### Entity: OrganizationNode (tree)

`type`: COMPANY \| DIVISION \| DEPARTMENT \| SECTION \| TEAM \| LOCATION \| WORK_SITE

Inherits `AttendancePolicy` from parent (overridable). Linked `GeofenceZone`s.

### Entity: GeofenceZone

`type`: CIRCLE \| POLYGON; optional `validFrom` / `validUntil` for temporary sites.

### Value object: AttendancePolicy

`workdayDuration`, `toleranceMinutes`, `allowedPunchTypes`, `geofenceRequired`, `biometricRequired`, `offlineSyncMaxAge`, `overtimePolicy`

## State machines

### PunchStatus

```
PENDING (offline) ──sync──▶ SYNCING ──success──▶ VALID
                              ├──expired──▶ DISCARDED
                              └──fraud──▶ SUSPICIOUS
VALID ◀──approve── SUSPICIOUS ──reject──▶ REJECTED
```

### PunchType sequence (per workday)

```
CLOCK_IN → BREAK_START → BREAK_END → CLOCK_OUT
```

Invalid sequence → REJECTED (see BR-014).
