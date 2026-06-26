package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

// PunchRepository persists punch records under tenant RLS.
type PunchRepository struct {
	db *sqlx.DB
}

// NewPunchRepository returns a tenant-scoped punch repository.
func NewPunchRepository(db *sqlx.DB) *PunchRepository {
	return &PunchRepository{db: db}
}

type punchRow struct {
	ID         uuid.UUID `db:"id"`
	EmployeeID uuid.UUID `db:"employee_id"`
	TenantID   uuid.UUID `db:"tenant_id"`
	PunchType  string    `db:"punch_type"`
	PunchedAt  time.Time `db:"punched_at"`
	Status     string    `db:"status"`
}

// RecentPunches returns the most recent punches for an employee (newest first).
func (r *PunchRepository) RecentPunches(
	ctx context.Context,
	tenantID, employeeID uuid.UUID,
	limit int,
) ([]domainpunch.PunchRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	var rows []punchRow
	err := WithTenant(ctx, r.db, tenantID, func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, &rows, `
			SELECT id, tenant_id, employee_id, punch_type, punched_at, status
			FROM punch_records
			WHERE employee_id = $1
			ORDER BY punched_at DESC
			LIMIT $2`, employeeID, limit)
	})
	if err != nil {
		return nil, fmt.Errorf("recent punches: %w", err)
	}
	return rowsToDomain(rows), nil
}

// Save inserts a punch record under tenant scope.
func (r *PunchRepository) Save(
	ctx context.Context,
	tenantID uuid.UUID,
	record domainpunch.PunchRecord,
) error {
	id, err := uuid.Parse(record.ID)
	if err != nil {
		id = uuid.New()
	}
	employeeID, err := uuid.Parse(record.EmployeeID)
	if err != nil {
		return fmt.Errorf("parse employee id: %w", err)
	}
	return WithTenant(ctx, r.db, tenantID, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO punch_records (id, tenant_id, employee_id, punch_type, punched_at, status, sync_status)
			VALUES ($1, $2, $3, $4, $5, $6, 'SYNCED')`,
			id, tenantID, employeeID, string(record.Type), record.PunchedAt, string(record.Status))
		return err
	})
}

// CountByStatus counts punch rows for an employee with the given status.
func (r *PunchRepository) CountByStatus(
	ctx context.Context,
	tenantID, employeeID uuid.UUID,
	status domainpunch.PunchStatus,
) (int, error) {
	var count int
	err := WithTenant(ctx, r.db, tenantID, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &count, `
			SELECT count(*) FROM punch_records
			WHERE employee_id = $1 AND status = $2`, employeeID, string(status))
	})
	return count, err
}

func rowsToDomain(rows []punchRow) []domainpunch.PunchRecord {
	out := make([]domainpunch.PunchRecord, len(rows))
	for i, row := range rows {
		out[i] = domainpunch.PunchRecord{
			ID:         row.ID.String(),
			EmployeeID: row.EmployeeID.String(),
			TenantID:   row.TenantID.String(),
			PunchedAt:  row.PunchedAt,
			Type:       domainpunch.PunchType(row.PunchType),
			Status:     domainpunch.PunchStatus(row.Status),
		}
	}
	return out
}
