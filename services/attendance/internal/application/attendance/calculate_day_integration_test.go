//go:build integration

package attendance_test

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

	appattendance "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/attendance"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/postgres"
)

const testDay = "2026-06-26T00:00:00Z"

func TestCalculateDay_E2E_BR030_WorkedMinutesFromDB(t *testing.T) {
	day := mustDay(testDay)
	env := newAttendanceE2EEnv(t, workforce.Standard8h(), organization.DefaultPolicy())
	seedValidPunches(t, env, []domainpunch.PunchRecord{
		punchAt(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour)),
		punchAt(domainpunch.PunchTypeBreakStart, day.Add(12*time.Hour)),
		punchAt(domainpunch.PunchTypeBreakEnd, day.Add(13*time.Hour)),
		punchAt(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour)),
	})

	result, err := env.handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: env.tenantID, EmployeeID: env.employeeID, Day: day,
	})
	require.NoError(t, err)
	require.Equal(t, 8*60, result.WorkedMinutes)
}

func TestCalculateDay_E2E_BR031_LatenessFromDB(t *testing.T) {
	day := mustDay(testDay)
	policy := organization.DefaultPolicy()
	policy.ToleranceMinutes = 5
	env := newAttendanceE2EEnv(t, workforce.Standard8h(), policy)
	seedValidPunches(t, env, []domainpunch.PunchRecord{
		punchAt(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour+10*time.Minute)),
		punchAt(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour)),
	})

	result, err := env.handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: env.tenantID, EmployeeID: env.employeeID, Day: day,
	})
	require.NoError(t, err)
	require.Equal(t, 5, result.LatenessMinutes)
}

func TestCalculateDay_E2E_BR032_OvertimeFromDB(t *testing.T) {
	day := mustDay(testDay)
	policy := organization.DefaultPolicy()
	policy.ToleranceMinutes = 5
	env := newAttendanceE2EEnv(t, workforce.Standard8h(), policy)
	seedValidPunches(t, env, []domainpunch.PunchRecord{
		punchAt(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour)),
		punchAt(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour+20*time.Minute)),
	})

	result, err := env.handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: env.tenantID, EmployeeID: env.employeeID, Day: day,
	})
	require.NoError(t, err)
	require.Equal(t, 15, result.OvertimeMinutes)
}

func TestCalculateDay_E2E_BR033_Shift12x36Windows(t *testing.T) {
	day := mustDay(testDay)
	env := newAttendanceE2EEnv(t, workforce.Shift12x36(), organization.DefaultPolicy())

	result, err := env.handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: env.tenantID, EmployeeID: env.employeeID, Day: day,
	})
	require.NoError(t, err)
	require.Len(t, result.Windows, 1)
	require.Equal(t, 19, result.Windows[0].Start.Hour())
}

func TestCalculateDay_E2E_BR034_TimeBankFromDB(t *testing.T) {
	day := mustDay(testDay)
	policy := organization.DefaultPolicy()
	policy.OvertimePolicy = organization.OvertimePolicyTimeBank
	policy.ToleranceMinutes = 5
	env := newAttendanceE2EEnv(t, workforce.Standard8h(), policy)
	seedValidPunches(t, env, []domainpunch.PunchRecord{
		punchAt(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour)),
		punchAt(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour+20*time.Minute)),
	})

	result, err := env.handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: env.tenantID, EmployeeID: env.employeeID, Day: day,
		PriorTimeBankMinutes: 120,
	})
	require.NoError(t, err)
	require.Equal(t, 135, result.TimeBankBalance)
}

type attendanceE2EEnv struct {
	handler    appattendance.CalculateDayAttendanceHandler
	punchRepo  *postgres.PunchRepository
	tenantID   uuid.UUID
	employeeID uuid.UUID
}

func newAttendanceE2EEnv(
	t *testing.T,
	schedule workforce.WorkSchedule,
	policy organization.AttendancePolicy,
) attendanceE2EEnv {
	t.Helper()
	adminDB, appDB := startPostgres(t)
	tenantID, _, employeeID := seedEmployee(t, adminDB)
	punchRepo := postgres.NewPunchRepository(appDB)

	handler := appattendance.CalculateDayAttendanceHandler{
		Punches: punchRepo,
		Schedules: &fixedScheduleResolver{schedule: schedule},
		Policies:  &fixedPolicyResolver{policy: policy},
	}
	return attendanceE2EEnv{
		handler: handler, punchRepo: punchRepo,
		tenantID: tenantID, employeeID: employeeID,
	}
}

func seedValidPunches(t *testing.T, env attendanceE2EEnv, punches []domainpunch.PunchRecord) {
	t.Helper()
	ctx := context.Background()
	for i, p := range punches {
		p.ID = uuid.New().String()
		p.EmployeeID = env.employeeID.String()
		p.TenantID = env.tenantID.String()
		p.Status = domainpunch.PunchStatusValid
		require.NoError(t, env.punchRepo.Save(ctx, env.tenantID, p), "punch %d", i)
	}
}

func punchAt(t domainpunch.PunchType, at time.Time) domainpunch.PunchRecord {
	return domainpunch.PunchRecord{Type: t, PunchedAt: at}
}

func mustDay(value string) time.Time {
	day, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return day.UTC()
}

type fixedScheduleResolver struct {
	schedule workforce.WorkSchedule
}

func (f *fixedScheduleResolver) WorkSchedule(
	_ context.Context, _, _ uuid.UUID,
) (workforce.WorkSchedule, error) {
	return f.schedule, nil
}

type fixedPolicyResolver struct {
	policy organization.AttendancePolicy
}

func (f *fixedPolicyResolver) EffectivePolicy(
	_ context.Context, _, _ uuid.UUID,
) (organization.AttendancePolicy, error) {
	return f.policy, nil
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
