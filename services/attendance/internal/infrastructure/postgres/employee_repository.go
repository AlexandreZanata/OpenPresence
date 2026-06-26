package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Employee is a minimal tenant-scoped row for RLS demonstrations.
type Employee struct {
	ID           uuid.UUID `db:"id"`
	TenantID     uuid.UUID `db:"tenant_id"`
	Registration string    `db:"registration"`
	Status       string    `db:"status"`
}

// EmployeeRepository loads employees through tenant-scoped transactions.
type EmployeeRepository struct {
	db *sqlx.DB
}

// NewEmployeeRepository returns a repository backed by db.
func NewEmployeeRepository(db *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// GetEmployee returns the employee when visible under tenantID, or nil when RLS filters the row.
func (r *EmployeeRepository) GetEmployee(
	ctx context.Context,
	tenantID, employeeID uuid.UUID,
) (*Employee, error) {
	var emp Employee
	err := WithTenant(ctx, r.db, tenantID, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &emp, `
			SELECT id, tenant_id, registration, status
			FROM employees
			WHERE id = $1`, employeeID)
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get employee: %w", err)
	}
	return &emp, nil
}
