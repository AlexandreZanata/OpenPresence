package enrollment

import (
	"context"

	"github.com/google/uuid"
)

// TenantEmployee is a minimal employee row for enrollment orchestration.
type TenantEmployee struct {
	ID       uuid.UUID
	TenantID uuid.UUID
}

// EmployeeReader loads tenant-scoped employees.
type EmployeeReader interface {
	GetEmployee(ctx context.Context, tenantID, employeeID uuid.UUID) (*TenantEmployee, error)
}

// FaceEmbeddingWriter persists embeddings after biometric enrollment.
type FaceEmbeddingWriter interface {
	Save(ctx context.Context, tenantID, employeeID uuid.UUID) error
}
