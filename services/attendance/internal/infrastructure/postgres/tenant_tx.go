package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// WithTenant runs fn inside a transaction scoped to tenantID via SET LOCAL app.tenant_id.
func WithTenant(ctx context.Context, db *sqlx.DB, tenantID uuid.UUID, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(
		ctx,
		"SELECT set_config('app.tenant_id', $1, true)",
		tenantID.String(),
	); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
