//go:build integration

package enrollment_test

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

	appenroll "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/enrollment"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/postgres"
)

func TestSaveFaceEmbedding_E2E_RLS_CrossTenantIsolation(t *testing.T) {
	env := newEnrollmentE2EEnv(t)
	ctx := context.Background()

	err := env.handler.Handle(ctx, appenroll.SaveFaceEmbeddingCommand{
		TenantID: env.tenantB, EmployeeID: env.employeeB,
	})
	require.NoError(t, err)

	leakCount, err := env.embedRepo.CountActive(ctx, env.tenantA, env.employeeB)
	require.NoError(t, err)
	require.Equal(t, 0, leakCount, "tenant A must not see tenant B embeddings")

	ownCount, err := env.embedRepo.CountActive(ctx, env.tenantB, env.employeeB)
	require.NoError(t, err)
	require.Equal(t, 1, ownCount, "tenant B must see its own embedding")
}

func TestSaveFaceEmbedding_E2E_RLS_CrossTenantEnrollDenied(t *testing.T) {
	env := newEnrollmentE2EEnv(t)
	ctx := context.Background()

	err := env.handler.Handle(ctx, appenroll.SaveFaceEmbeddingCommand{
		TenantID: env.tenantA, EmployeeID: env.employeeB,
	})
	require.ErrorIs(t, err, appenroll.ErrEmployeeNotFound)

	count, err := env.embedRepo.CountActive(ctx, env.tenantB, env.employeeB)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

type enrollmentE2EEnv struct {
	handler   appenroll.SaveFaceEmbeddingHandler
	embedRepo *postgres.FaceEmbeddingRepository
	tenantA   uuid.UUID
	tenantB   uuid.UUID
	employeeB uuid.UUID
}

func newEnrollmentE2EEnv(t *testing.T) enrollmentE2EEnv {
	t.Helper()
	adminDB, appDB := startPostgres(t)
	tenantA, tenantB, _, employeeB := seedEmployees(t, adminDB)
	empRepo := postgres.NewEmployeeRepository(appDB)
	embedRepo := postgres.NewFaceEmbeddingRepository(appDB)
	handler := appenroll.SaveFaceEmbeddingHandler{
		Employees:  postgresEmployeeAdapter{repo: empRepo},
		Embeddings: embedRepo,
	}
	return enrollmentE2EEnv{
		handler: handler, embedRepo: embedRepo,
		tenantA: tenantA, tenantB: tenantB, employeeB: employeeB,
	}
}

type postgresEmployeeAdapter struct {
	repo *postgres.EmployeeRepository
}

func (a postgresEmployeeAdapter) GetEmployee(
	ctx context.Context, tenantID, employeeID uuid.UUID,
) (*appenroll.TenantEmployee, error) {
	emp, err := a.repo.GetEmployee(ctx, tenantID, employeeID)
	if err != nil || emp == nil {
		return nil, err
	}
	return &appenroll.TenantEmployee{ID: emp.ID, TenantID: emp.TenantID}, nil
}

func seedEmployees(t *testing.T, admin *sqlx.DB) (tenantA, tenantB, employeeA, employeeB uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	_, err := admin.ExecContext(ctx, `SET row_security = off`)
	require.NoError(t, err)

	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-a') RETURNING id`).Scan(&tenantA))
	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-b') RETURNING id`).Scan(&tenantB))
	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO employees (tenant_id, registration, status)
		VALUES ($1, 'EMP-A', 'ACTIVE') RETURNING id`, tenantA).Scan(&employeeA))
	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO employees (tenant_id, registration, status)
		VALUES ($1, 'EMP-B', 'ACTIVE') RETURNING id`, tenantB).Scan(&employeeB))
	return tenantA, tenantB, employeeA, employeeB
}

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

	require.NoError(t, postgres.ApplyMigrations(adminDB.DB, migrationsDir(t)))

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

func migrationsDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	dir := filepath.Join(filepath.Dir(filename), "..", "..", "..", "migrations")
	abs, err := filepath.Abs(dir)
	require.NoError(t, err)
	return abs
}
