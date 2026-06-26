package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// FaceEmbeddingRepository persists enrollment embeddings under tenant RLS.
type FaceEmbeddingRepository struct {
	db *sqlx.DB
}

// NewFaceEmbeddingRepository returns a tenant-scoped embedding repository.
func NewFaceEmbeddingRepository(db *sqlx.DB) *FaceEmbeddingRepository {
	return &FaceEmbeddingRepository{db: db}
}

// Save inserts an active face embedding for an employee.
func (r *FaceEmbeddingRepository) Save(
	ctx context.Context,
	tenantID, employeeID uuid.UUID,
) error {
	return WithTenant(ctx, r.db, tenantID, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO face_embeddings (tenant_id, employee_id, active)
			VALUES ($1, $2, true)`, tenantID, employeeID)
		return err
	})
}

// CountActive returns active embeddings visible under tenant scope.
func (r *FaceEmbeddingRepository) CountActive(
	ctx context.Context,
	tenantID, employeeID uuid.UUID,
) (int, error) {
	var count int
	err := WithTenant(ctx, r.db, tenantID, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &count, `
			SELECT count(*) FROM face_embeddings
			WHERE employee_id = $1 AND active = true`, employeeID)
	})
	if err != nil {
		return 0, fmt.Errorf("count active embeddings: %w", err)
	}
	return count, nil
}
