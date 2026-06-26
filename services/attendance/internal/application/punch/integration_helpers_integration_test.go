//go:build integration

package punch_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/postgres"
)

const integrationBaseTime = "2026-06-26T09:00:00Z"

type integrationEnv struct {
	handler          apppunch.SubmitPunchHandler
	punchRepo        *postgres.PunchRepository
	appDB            *sqlx.DB
	tenantID         uuid.UUID
	otherTenantID    uuid.UUID
	employeeID       uuid.UUID
	otherEmployeeID  uuid.UUID
	serverTime       time.Time
	integrationOpts  integrationOpts
}

type integrationOpts struct {
	biometric apppunch.BiometricClient
	clock     func() time.Time
	zones     []geofence.GeofenceZone
	lockout   *fraud.DeviceLockoutTracker
}

func defaultIntegrationOpts() integrationOpts {
	serverTime := mustParseTime(integrationBaseTime)
	return integrationOpts{
		biometric: configurableBiometricClient{
			IsLive: true, LivenessScore: 0.95,
			IsRecognized: true, RecognitionConfidence: 0.90,
		},
		clock: func() time.Time { return serverTime },
	}
}

type configurableBiometricClient struct {
	IsLive                bool
	LivenessScore         float64
	IsRecognized          bool
	RecognitionConfidence float64
}

func (c configurableBiometricClient) VerifyPunch(
	_ context.Context, _, _ uuid.UUID, _ []byte,
) (*apppunch.BiometricVerifyResult, error) {
	return &apppunch.BiometricVerifyResult{
		IsLive: c.IsLive, LivenessScore: c.LivenessScore,
		IsRecognized: c.IsRecognized, RecognitionConfidence: c.RecognitionConfidence,
		EmbeddingHash: "e2e-hash",
	}, nil
}

func newIntegrationEnv(t *testing.T) integrationEnv {
	t.Helper()
	return newIntegrationEnvWithOpts(t, defaultIntegrationOpts())
}

func newIntegrationEnvWithZones(t *testing.T, zones []geofence.GeofenceZone) integrationEnv {
	t.Helper()
	opts := defaultIntegrationOpts()
	opts.zones = zones
	return newIntegrationEnvWithOpts(t, opts)
}

func newIntegrationEnvWithOpts(t *testing.T, opts integrationOpts) integrationEnv {
	t.Helper()
	adminDB, appDB := startPostgres(t)
	tenantID, otherTenantID, employeeID, otherEmployeeID := seedEmployees(t, adminDB)
	serverTime := opts.clock()

	empRepo := postgres.NewEmployeeRepository(appDB)
	punchRepo := postgres.NewPunchRepository(appDB)

	handler := buildSubmitPunchHandler(
		empRepo, punchRepo, tenantID, employeeID, serverTime, opts,
	)

	return integrationEnv{
		handler: handler, punchRepo: punchRepo, appDB: appDB,
		tenantID: tenantID, otherTenantID: otherTenantID,
		employeeID: employeeID, otherEmployeeID: otherEmployeeID,
		serverTime: serverTime, integrationOpts: opts,
	}
}

func buildSubmitPunchHandler(
	empRepo *postgres.EmployeeRepository,
	punchRepo *postgres.PunchRepository,
	tenantID, employeeID uuid.UUID,
	serverTime time.Time,
	opts integrationOpts,
) apppunch.SubmitPunchHandler {
	zones := opts.zones
	if zones == nil {
		zones = []geofence.GeofenceZone{testZone()}
	}
	return apppunch.SubmitPunchHandler{
		Employees: employeeReaderAdapter{repo: empRepo},
		Placements: &stubPlacementReader{placement: &workforce.EmployeePlacement{
			ID: "pl-1", EmployeeID: employeeID.String(), TenantID: tenantID.String(),
			OrgNodeID: "site-1", Type: workforce.PlacementTypePrimary,
			ValidFrom: serverTime.Add(-time.Hour),
		}},
		Policies:  &stubPolicyReader{policy: organization.DefaultPolicy()},
		Geofences: &stubGeofenceReader{zones: zones},
		Biometric: opts.biometric,
		Punches:   punchRepo,
		Validator: domainpunch.PunchValidator{},
		Fraud:     fraud.FraudEvaluator{},
		Lockout:   opts.lockout,
		Clock:     opts.clock,
	}
}

func (env integrationEnv) handlerForTenant(tenantID, employeeID uuid.UUID) apppunch.SubmitPunchHandler {
	empRepo := postgres.NewEmployeeRepository(env.appDB)
	punchRepo := postgres.NewPunchRepository(env.appDB)
	return buildSubmitPunchHandler(
		empRepo, punchRepo, tenantID, employeeID, env.serverTime, env.integrationOpts,
	)
}

func validPunchCmd(env integrationEnv, mutate func(*apppunch.SubmitPunchCommand)) apppunch.SubmitPunchCommand {
	cmd := apppunch.SubmitPunchCommand{
		TenantID:   env.tenantID,
		EmployeeID: env.employeeID,
		Type:       domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10,
		},
		DeviceTime: env.serverTime,
		FrameJPEG:  []byte("e2e-frame"),
	}
	if mutate != nil {
		mutate(&cmd)
	}
	return cmd
}

func insideLocation() domainpunch.GpsCoordinate {
	return domainpunch.GpsCoordinate{Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10}
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

func mustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}
