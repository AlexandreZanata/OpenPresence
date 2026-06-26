# Domain Events

Past-tense names. Immutable once published. Handlers in Application/Infrastructure layers only.

| Event | Raised when | Payload (summary) |
|-------|-------------|-----------------|
| `EmployeeEnrolled` | Biometric enrollment completes | `employeeId`, `tenantId`, `enrolledBy`, `embeddingCount` |
| `EmployeeDeactivated` | Status → INACTIVE/SUSPENDED | `employeeId`, `tenantId`, `reason` |
| `PunchRecorded` | Punch persisted (any status) | `punchRecordId`, `employeeId`, `type`, `status`, `punchedAt` |
| `PunchApproved` | Manager approves SUSPICIOUS | `punchRecordId`, `approvedBy` |
| `PunchRejected` | Manager or system rejects | `punchRecordId`, `reason`, `rejectedBy` |
| `FraudDetected` | FraudFlag with HIGH/CRITICAL | `punchRecordId`, `fraudType`, `severity` |
| `DeviceLocked` | BR-013 lockout triggered | `deviceId`, `employeeId`, `until` |
| `GeofenceCreated` | New zone active | `geofenceZoneId`, `tenantId`, `orgNodeId` |
| `PolicyUpdated` | AttendancePolicy changed on node | `orgNodeId`, `tenantId`, `changedFields` |
| `OfflinePunchDiscarded` | Sync past TTL | `localPunchId`, `employeeId`, `deviceTime` |

## NATS subjects (convention)

```
openpresence.{tenant_id}.attendance.punch_recorded
openpresence.{tenant_id}.fraud.detected
openpresence.{tenant_id}.audit.action
```

Consumers: audit writer, notification service, reporting projections.
