//go:build integration

package punch_test

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

func startPostgres(t *testing.T) (admin, app *sqlx.DB) {
	t.Helper()
	ctx := context.Background()
	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("openpresence"),
		tcpostgres.WithUsername("openpresence"),
		tcpostgres.WithPassword("openpresence"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second),
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
		host, port.Port(),
	)
	appDB, err := sqlx.Connect("postgres", appConn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = appDB.Close() })
	return adminDB, appDB
}

func seedEmployees(t *testing.T, admin *sqlx.DB) (tenantID, otherTenantID, employeeID, otherEmployeeID uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	_, err := admin.ExecContext(ctx, `SET row_security = off`)
	require.NoError(t, err)

	err = admin.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-a') RETURNING id`).Scan(&tenantID)
	require.NoError(t, err)
	err = admin.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-b') RETURNING id`).Scan(&otherTenantID)
	require.NoError(t, err)
	err = admin.QueryRowContext(ctx, `
		INSERT INTO employees (tenant_id, registration, status)
		VALUES ($1, 'EMP-1', 'ACTIVE') RETURNING id`, tenantID).Scan(&employeeID)
	require.NoError(t, err)
	err = admin.QueryRowContext(ctx, `
		INSERT INTO employees (tenant_id, registration, status)
		VALUES ($1, 'EMP-B', 'ACTIVE') RETURNING id`, otherTenantID).Scan(&otherEmployeeID)
	require.NoError(t, err)
	return tenantID, otherTenantID, employeeID, otherEmployeeID
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
