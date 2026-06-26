package authorization

import (
	"context"

	"github.com/google/uuid"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

// TenantEmployee is a minimal employee row for authorization orchestration.
type TenantEmployee struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Registration string
}

// EmployeeReader loads tenant-scoped employees.
type EmployeeReader interface {
	GetEmployee(ctx context.Context, tenantID, employeeID uuid.UUID) (*TenantEmployee, error)
}

// PlacementReader resolves org placement for ABAC checks.
type PlacementReader interface {
	EmployeePlacement(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
	) (organization.EmployeePlacementRef, error)
}
