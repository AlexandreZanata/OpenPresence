package enrollment

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrEmployeeNotFound = errors.New("enrollment: employee not found")

// SaveFaceEmbeddingCommand stores an embedding after successful EnrollFace.
type SaveFaceEmbeddingCommand struct {
	TenantID   uuid.UUID
	EmployeeID uuid.UUID
}

// SaveFaceEmbeddingHandler persists embeddings only for visible tenant employees.
type SaveFaceEmbeddingHandler struct {
	Employees  EmployeeReader
	Embeddings FaceEmbeddingWriter
}

// Handle saves the embedding when the employee exists under the tenant.
func (h SaveFaceEmbeddingHandler) Handle(
	ctx context.Context,
	cmd SaveFaceEmbeddingCommand,
) error {
	emp, err := h.Employees.GetEmployee(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		return err
	}
	if emp == nil {
		return ErrEmployeeNotFound
	}
	return h.Embeddings.Save(ctx, cmd.TenantID, cmd.EmployeeID)
}
