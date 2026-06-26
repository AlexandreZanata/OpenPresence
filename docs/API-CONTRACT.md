# API Contract — OpenPresence

> Version all public APIs from day one. OpenAPI specs live in `docs/api/` (to be generated).

## Base URL

```
https://api.{tenant-host}/v1
```

## Authentication

| Method | Mechanism |
|--------|-----------|
| Bearer JWT | `Authorization: Bearer <token>` (15 min expiry) |
| Refresh | `POST /v1/auth/refresh` (7 day refresh token) |
| Tenant | Resolved from JWT claims + `X-Tenant-Id` validation |

## Error format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable safe message",
    "correlationId": "uuid"
  }
}
```

## Pagination

| Param | Default | Max |
|-------|---------|-----|
| `page` | 1 | — |
| `pageSize` | 20 | 50 |

---

## Auth Service

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/auth/login` | Registration ID + password → JWT + refresh |
| POST | `/v1/auth/refresh` | Refresh token → new JWT |
| POST | `/v1/auth/device/register` | Register device token |
| DELETE | `/v1/auth/device/{id}` | Revoke device |

## Attendance Service

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/v1/attendance/punch` | EMPLOYEE | Register punch |
| GET | `/v1/attendance/today` | EMPLOYEE | Today's punches |
| GET | `/v1/attendance/history` | EMPLOYEE | Paginated history |
| GET | `/v1/attendance/report` | MANAGER+ | Period report |
| GET | `/v1/attendance/suspicious` | MANAGER+ | Pending review |
| PATCH | `/v1/attendance/{id}/approve` | MANAGER+ | Approve suspicious |
| PATCH | `/v1/attendance/{id}/reject` | MANAGER+ | Reject with reason |
| PATCH | `/v1/attendance/{id}/adjust` | HR_ANALYST | Manual adjustment (audited) |

### POST `/v1/attendance/punch`

**Request (allow-listed fields only):**

```json
{
  "punchType": "CLOCK_IN",
  "deviceTime": "2026-06-26T08:01:00Z",
  "location": {
    "latitude": -12.5458,
    "longitude": -55.7061,
    "accuracy": 12.5,
    "altitude": 365.0,
    "provider": "FUSED",
    "isMocked": false
  },
  "frameBase64": "...",
  "deviceIntegrityReport": {
    "isRooted": false,
    "isVpnActive": false,
    "isDeveloperOptionsEnabled": false
  }
}
```

**Response 201:**

```json
{
  "id": "uuid",
  "status": "VALID",
  "punchedAt": "2026-06-26T08:01:02Z",
  "type": "CLOCK_IN",
  "fraudFlags": []
}
```

**Idempotency:** `Idempotency-Key` header required for offline bulk sync.

## Organization Service

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/org/nodes` | Create org node |
| GET | `/v1/org/tree` | Full tenant tree |
| PATCH | `/v1/org/nodes/{id}/policy` | Update attendance policy |
| POST | `/v1/org/geofences` | Create geofence |
| PATCH | `/v1/org/geofences/{id}` | Update geofence |
| DELETE | `/v1/org/geofences/{id}` | Deactivate geofence |

## Workforce Service

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/employees` | Create employee |
| PATCH | `/v1/employees/{id}` | Update employee |
| POST | `/v1/employees/{id}/enroll` | Start biometric enrollment |
| GET | `/v1/employees/{id}/profile` | Profile + biometric status |
| PATCH | `/v1/employees/{id}/status` | ACTIVE \| INACTIVE \| SUSPENDED |

## Biometric Service (internal gRPC)

```protobuf
service BiometricService {
  rpc VerifyPunch(VerifyPunchRequest) returns (VerifyPunchResponse);
  rpc EnrollFace(EnrollFaceRequest) returns (EnrollFaceResponse);
  rpc DeleteProfile(DeleteProfileRequest) returns (DeleteProfileResponse);
}

message VerifyPunchRequest {
  bytes frame_jpeg = 1;
  string employee_id = 2;
  string tenant_id = 3;
}

message VerifyPunchResponse {
  bool is_live = 1;
  float liveness_score = 2;
  bool is_recognized = 3;
  float recognition_confidence = 4;
  string matched_employee_id = 5;
  repeated FraudFlag fraud_flags = 6;
  bytes embedding = 7;  // internal only — never exposed on REST
}

message EnrollFaceRequest {
  bytes frame_jpeg = 1;
  string employee_id = 2;
  string tenant_id = 3;
  string angle = 4;  // FRONTAL | LEFT_15 | RIGHT_15
}

message EnrollFaceResponse {
  bool is_live = 1;
  float liveness_score = 2;
  float quality_score = 3;
  bytes embedding = 4;
  repeated FraudFlag fraud_flags = 5;
}

message DeleteProfileRequest {
  string employee_id = 1;
  string tenant_id = 2;
}

message DeleteProfileResponse {
  bool success = 1;
}
```

Proto source: `services/biometric/proto/biometric.proto`

**Security:** REST responses expose `faceEmbeddingHash` only, never raw embedding (see [SECURITY.md](SECURITY.md)).
