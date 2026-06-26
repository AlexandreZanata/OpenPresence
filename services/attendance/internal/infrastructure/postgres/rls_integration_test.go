//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/postgres"
)

type testDBs struct {
	admin *sqlx.DB
	app   *sqlx.DB
}

func TestMigrations_ApplyOnEmptyDB(t *testing.T) {
	dbs := newTestDBs(t)

	var tableCount int
	err := dbs.admin.QueryRow(`
		SELECT count(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name IN ('employees', 'punch_records', 'face_embeddings')
	`).Scan(&tableCount)
	require.NoError(t, err)
	require.Equal(t, 3, tableCount)
}

func TestRLS_TenantCannotReadOtherTenantEmployee(t *testing.T) {
	dbs := newTestDBs(t)
	ctx := context.Background()
	repo := postgres.NewEmployeeRepository(dbs.app)

	tenantA, tenantB, employeeA, employeeB := seedTenantsAndEmployees(t, dbs.admin)

	emp, err := repo.GetEmployee(ctx, tenantA, employeeB)
	require.NoError(t, err)
	require.Nil(t, emp, "tenant A must not see tenant B employee")

	emp, err = repo.GetEmployee(ctx, tenantB, employeeA)
	require.NoError(t, err)
	require.Nil(t, emp, "tenant B must not see tenant A employee")

	emp, err = repo.GetEmployee(ctx, tenantA, employeeA)
	require.NoError(t, err)
	require.NotNil(t, emp)
	require.Equal(t, "EMP-A", emp.Registration)

	emp, err = repo.GetEmployee(ctx, tenantB, employeeB)
	require.NoError(t, err)
	require.NotNil(t, emp)
	require.Equal(t, "EMP-B", emp.Registration)
}

func TestRLS_PunchRecordsAndEmbeddingsIsolated(t *testing.T) {
	dbs := newTestDBs(t)
	ctx := context.Background()

	tenantA, tenantB, employeeA, employeeB := seedTenantsAndEmployees(t, dbs.admin)
	punchA, punchB := seedPunchRecords(t, dbs.admin, tenantA, employeeA, tenantB, employeeB)
	embedA, embedB := seedFaceEmbeddings(t, dbs.admin, tenantA, employeeA, tenantB, employeeB)

	var punchCount int
	err := postgres.WithTenant(ctx, dbs.app, tenantA, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &punchCount, `SELECT count(*) FROM punch_records`)
	})
	require.NoError(t, err)
	require.Equal(t, 1, punchCount)

	var embedCount int
	err = postgres.WithTenant(ctx, dbs.app, tenantA, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &embedCount, `SELECT count(*) FROM face_embeddings`)
	})
	require.NoError(t, err)
	require.Equal(t, 1, embedCount)

	var crossPunchCount int
	err = postgres.WithTenant(ctx, dbs.app, tenantA, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &crossPunchCount, `SELECT count(*) FROM punch_records WHERE id = $1`, punchB)
	})
	require.NoError(t, err)
	require.Equal(t, 0, crossPunchCount)

	var crossEmbedCount int
	err = postgres.WithTenant(ctx, dbs.app, tenantB, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &crossEmbedCount, `SELECT count(*) FROM face_embeddings WHERE id = $1`, embedA)
	})
	require.NoError(t, err)
	require.Equal(t, 0, crossEmbedCount)

	_ = punchA
	_ = embedB
}

func newTestDBs(t *testing.T) testDBs {
	t.Helper()

	ctx := context.Background()
	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("openpresence"),
		tcpostgres.WithUsername("openpresence"),
		tcpostgres.WithPassword("openpresence"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	adminConn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	adminDB, err := sqlx.Connect("postgres", adminConn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = adminDB.Close() })

	err = postgres.ApplyMigrations(adminDB.DB, migrationsDir(t))
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	appConn := fmt.Sprintf(
		"postgres://attendance_app:attendance_app@%s:%s/openpresence?sslmode=disable",
		host,
		port.Port(),
	)
	appDB, err := sqlx.Connect("postgres", appConn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = appDB.Close() })

	return testDBs{admin: adminDB, app: appDB}
}

func migrationsDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	dir := filepath.Join(filepath.Dir(filename), "..", "..", "..", "migrations")
	abs, err := filepath.Abs(dir)
	require.NoError(t, err)
	return abs
}

func seedTenantsAndEmployees(t *testing.T, db *sqlx.DB) (uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `SET row_security = off`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `SET row_security = on`)
	})

	var tenantA, tenantB uuid.UUID
	err = db.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-a') RETURNING id`).Scan(&tenantA)
	require.NoError(t, err)
	err = db.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-b') RETURNING id`).Scan(&tenantB)
	require.NoError(t, err)

	var employeeA, employeeB uuid.UUID
	err = db.QueryRowContext(ctx, `
		INSERT INTO employees (tenant_id, registration) VALUES ($1, 'EMP-A') RETURNING id`,
		tenantA).Scan(&employeeA)
	require.NoError(t, err)
	err = db.QueryRowContext(ctx, `
		INSERT INTO employees (tenant_id, registration) VALUES ($1, 'EMP-B') RETURNING id`,
		tenantB).Scan(&employeeB)
	require.NoError(t, err)

	return tenantA, tenantB, employeeA, employeeB
}

func seedPunchRecords(t *testing.T, db *sqlx.DB, tenantA, employeeA, tenantB, employeeB uuid.UUID) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	var punchA, punchB uuid.UUID
	err := db.QueryRowContext(ctx, `
		INSERT INTO punch_records (tenant_id, employee_id, punch_type)
		VALUES ($1, $2, 'CLOCK_IN') RETURNING id`, tenantA, employeeA).Scan(&punchA)
	require.NoError(t, err)
	err = db.QueryRowContext(ctx, `
		INSERT INTO punch_records (tenant_id, employee_id, punch_type)
		VALUES ($1, $2, 'CLOCK_IN') RETURNING id`, tenantB, employeeB).Scan(&punchB)
	require.NoError(t, err)
	return punchA, punchB
}

func seedFaceEmbeddings(t *testing.T, db *sqlx.DB, tenantA, employeeA, tenantB, employeeB uuid.UUID) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	var embedA, embedB uuid.UUID
	err := db.QueryRowContext(ctx, `
		INSERT INTO face_embeddings (tenant_id, employee_id) VALUES ($1, $2) RETURNING id`,
		tenantA, employeeA).Scan(&embedA)
	require.NoError(t, err)
	err = db.QueryRowContext(ctx, `
		INSERT INTO face_embeddings (tenant_id, employee_id) VALUES ($1, $2) RETURNING id`,
		tenantB, employeeB).Scan(&embedB)
	require.NoError(t, err)
	return embedA, embedB
}
