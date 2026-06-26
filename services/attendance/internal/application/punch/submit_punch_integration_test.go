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

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/postgres"
)

func TestSubmitPunch_Integration_HappyPath_VALIDInDB(t *testing.T) {
	env := newIntegrationEnv(t)
	now := time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)

	result, err := env.handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID:   env.tenantID,
		EmployeeID: env.employeeID,
		Type:       domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10,
		},
		DeviceTime: now,
		FrameJPEG:  []byte("integration-frame"),
	})
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)

	count, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestSubmitPunch_Integration_CrossTenant_Rejected(t *testing.T) {
	env := newIntegrationEnv(t)
	now := time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)

	_, err := env.handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID:   env.otherTenantID,
		EmployeeID: env.employeeID,
		Type:       domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10,
		},
		DeviceTime: now,
		FrameJPEG:  []byte("frame"),
	})
	require.ErrorIs(t, err, apppunch.ErrEmployeeNotFound)

	count, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestSubmitPunch_Integration_InvalidSequence_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)
	now := time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)

	_, err := env.handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID: env.tenantID, EmployeeID: env.employeeID,
		Type: domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10,
		},
		DeviceTime: now, FrameJPEG: []byte("first"),
	})
	require.NoError(t, err)

	result, err := env.handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID: env.tenantID, EmployeeID: env.employeeID,
		Type: domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10,
		},
		DeviceTime: now.Add(2 * time.Minute), FrameJPEG: []byte("second"),
	})
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)

	validCount, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, validCount)
}

func TestSubmitPunch_Integration_OutOfGeofence_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)
	now := time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)

	result, err := env.handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID:   env.tenantID,
		EmployeeID: env.employeeID,
		Type:       domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -22.0, Longitude: -43.0, Accuracy: 10,
		},
		DeviceTime: now,
		FrameJPEG:  []byte("far-away"),
	})
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)

	validCount, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, validCount)
}

type integrationEnv struct {
	handler         apppunch.SubmitPunchHandler
	punchRepo       *postgres.PunchRepository
	tenantID        uuid.UUID
	otherTenantID   uuid.UUID
	employeeID      uuid.UUID
}

func newIntegrationEnv(t *testing.T) integrationEnv {
	t.Helper()
	ctx := context.Background()
	adminDB, appDB := startPostgres(t)

	tenantID, otherTenantID, employeeID := seedEmployee(t, adminDB)
	zone := testZone()
	now := time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)

	empRepo := postgres.NewEmployeeRepository(appDB)
	punchRepo := postgres.NewPunchRepository(appDB)

	handler := apppunch.SubmitPunchHandler{
		Employees:  employeeReaderAdapter{repo: empRepo},
		Placements: &stubPlacementReader{placement: &workforce.EmployeePlacement{
			ID: "pl-1", EmployeeID: employeeID.String(), TenantID: tenantID.String(),
			OrgNodeID: "site-1", Type: workforce.PlacementTypePrimary, ValidFrom: now.Add(-time.Hour),
		}},
		Policies:  &stubPolicyReader{policy: organization.DefaultPolicy()},
		Geofences: &stubGeofenceReader{zones: []geofence.GeofenceZone{zone}},
		Biometric: stubBiometricClient{},
		Punches:   punchRepo,
		Validator: domainpunch.PunchValidator{},
		Fraud:     fraud.FraudEvaluator{},
		Clock:     func() time.Time { return now },
	}

	_ = ctx
	return integrationEnv{
		handler: handler, punchRepo: punchRepo,
		tenantID: tenantID, otherTenantID: otherTenantID, employeeID: employeeID,
	}
}

type employeeReaderAdapter struct {
	repo *postgres.EmployeeRepository
}

func (a employeeReaderAdapter) GetEmployee(
	ctx context.Context, tenantID, employeeID uuid.UUID,
) (*apppunch.Employee, error) {
	emp, err := a.repo.GetEmployee(ctx, tenantID, employeeID)
	if err != nil || emp == nil {
		return nil, err
	}
	return &apppunch.Employee{
		ID: emp.ID, TenantID: emp.TenantID,
		Registration: emp.Registration, Status: emp.Status,
	}, nil
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

func seedEmployee(t *testing.T, admin *sqlx.DB) (tenantID, otherTenantID, employeeID uuid.UUID) {
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
	return tenantID, otherTenantID, employeeID
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
