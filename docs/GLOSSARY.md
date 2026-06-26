# Domain Glossary

> Ubiquitous language. Code, APIs, docs, and agents MUST use these terms exactly.
> Do not invent domain terms without adding an entry here.

---

## Tenant

**Definition:** Root organization (private company or municipality). All data is isolated by `tenant_id`.
**Not the same as:** OrgNode (tenant is the billing/isolation boundary).
**Code name:** `Tenant`

---

## OrgNode

**Definition:** Node in the organizational tree (company, division, department, team, location, work site).
**Enum values:** `COMPANY`, `DIVISION`, `DEPARTMENT`, `SECTION`, `TEAM`, `LOCATION`, `WORK_SITE`
**Code name:** `OrganizationNode`

---

## Employee

**Definition:** Person registered to punch time for a tenant. Has biometric profile and work schedule.
**Enum values (status):** `ACTIVE`, `INACTIVE`, `SUSPENDED`
**Code name:** `Employee`

---

## PunchRecord

**Definition:** Immutable attendance mark (clock in/out, break). Official time is server `punchedAt`.
**Enum values (status):** `VALID`, `SUSPICIOUS`, `REJECTED`, `PENDING`, `DISCARDED`
**Enum values (type):** `CLOCK_IN`, `CLOCK_OUT`, `BREAK_START`, `BREAK_END`
**Code name:** `PunchRecord`

---

## GeofenceZone

**Definition:** Geographic boundary where punch is allowed.
**Enum values (type):** `CIRCLE`, `POLYGON`
**Code name:** `GeofenceZone`

---

## AttendancePolicy

**Definition:** Rules for work duration, tolerance, biometric/geofence requirements, offline sync TTL.
Inherited down the org tree with override at each node.
**Code name:** `AttendancePolicy`

---

## BiometricProfile

**Definition:** Employee's enrolled face embeddings (up to 5 active). Enrollment by administrator only.
**Code name:** `BiometricProfile`

---

## FaceEmbedding

**Definition:** 512-dimensional AuraFace vector plus metadata (angle, quality, capturedAt).
**Enum values (angle):** `FRONTAL`, `LEFT_15`, `RIGHT_15`
**Code name:** `FaceEmbedding`

---

## BiometricResult

**Definition:** Outcome of liveness + recognition for one punch. Stores hash of embedding, never raw vector in API responses.
**Code name:** `BiometricResult`

---

## FraudFlag

**Definition:** Detected anomaly attached to a punch or attempt.
**Enum values (severity):** `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`
**Code name:** `FraudFlag`

---

## FraudType

**Definition:** Category of fraud or integrity violation.
**Enum values:** `MOCK_GPS`, `CLOCK_MANIPULATION`, `IMPOSSIBLE_SPEED`, `LIVENESS_FAILED`, `FACE_NOT_RECOGNIZED`, `OUT_OF_GEOFENCE`, `DUPLICATE_PUNCH`, `DEVICE_ROOTED`, `VPN_DETECTED`, `GPS_LOW_ACCURACY`
**Code name:** `FraudType`

---

## GpsCoordinate

**Definition:** Location snapshot at punch time with accuracy and mock detection flag.
**Enum values (provider):** `GPS`, `NETWORK`, `FUSED`
**Code name:** `GpsCoordinate`

---

## WorkSchedule

**Definition:** Planned work windows, scheduled start/end, tolerance for an employee.
**Code name:** `WorkSchedule`

---

## User

**Definition:** Account that authenticates to the system (admin, manager, employee app user).
**Code name:** `User`

---

## Role

**Definition:** Named permission bundle: `SUPER_ADMIN`, `ORG_ADMIN`, `MANAGER`, `HR_ANALYST`, `SECURITY_ANALYST`, `EMPLOYEE`, `AUDITOR`
**Code name:** `Role`

---

## AuditLog

**Definition:** Immutable record of who changed what, when, from which IP.
**Code name:** `AuditLog`

---

## DeviceIntegrityReport

**Definition:** Mobile payload: root, VPN, developer options, mock GPS flags at punch time.
**Code name:** `DeviceIntegrityReport`
