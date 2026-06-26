package punch

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
)

// Employee is a tenant-scoped employee row visible to the use case.
type Employee struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Registration string
	Status       string
}

// BiometricVerifyResult is the outcome of a biometric gRPC VerifyPunch call.
type BiometricVerifyResult struct {
	IsLive                bool
	LivenessScore         float64
	RecognitionConfidence float64
	IsRecognized          bool
	EmbeddingHash         string
}

// EmployeeReader loads employees under tenant RLS.
type EmployeeReader interface {
	GetEmployee(ctx context.Context, tenantID, employeeID uuid.UUID) (*Employee, error)
}

// PlacementReader returns the active primary placement for an employee.
type PlacementReader interface {
	ActivePrimaryPlacement(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
		at time.Time,
	) (*workforce.EmployeePlacement, error)
}

// PolicyResolver returns the effective attendance policy for an org node.
type PolicyResolver interface {
	EffectivePolicy(
		ctx context.Context,
		tenantID uuid.UUID,
		orgNodeID string,
	) (organization.AttendancePolicy, error)
}

// GeofenceResolver returns geofence zones attached to org nodes on the placement path.
type GeofenceResolver interface {
	ZonesForOrgPath(
		ctx context.Context,
		tenantID uuid.UUID,
		orgNodeIDs []string,
	) ([]geofence.GeofenceZone, error)
}

// BiometricClient verifies a punch frame against enrolled biometrics.
type BiometricClient interface {
	VerifyPunch(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
		frameJPEG []byte,
	) (*BiometricVerifyResult, error)
}

// PunchRepository persists and loads punch history under tenant scope.
type PunchRepository interface {
	RecentPunches(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
		limit int,
	) ([]domainpunch.PunchRecord, error)
	Save(ctx context.Context, tenantID uuid.UUID, record domainpunch.PunchRecord) error
	CountByStatus(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
		status domainpunch.PunchStatus,
	) (int, error)
}
