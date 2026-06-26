package app

import (
	"context"
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

// PunchStack wires SubmitPunch with Postgres and in-memory policy ports.
type PunchStack struct {
	Handler   apppunch.SubmitPunchHandler
	PunchRepo *postgres.PunchRepository
}

// PunchStackConfig configures stub ports for the punch use case.
type PunchStackConfig struct {
	DB        *sqlx.DB
	Biometric apppunch.BiometricClient
	Clock     func() time.Time
	Zones     []geofence.GeofenceZone
}

// NewPunchStack builds the SubmitPunch handler graph for HTTP and E2E tests.
func NewPunchStack(cfg PunchStackConfig) PunchStack {
	clock := cfg.Clock
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	zones := cfg.Zones
	if zones == nil {
		zones = defaultZones()
	}
	now := clock()
	empRepo := postgres.NewEmployeeRepository(cfg.DB)
	punchRepo := postgres.NewPunchRepository(cfg.DB)

	return PunchStack{
		Handler: apppunch.SubmitPunchHandler{
			Employees:  employeeAdapter{repo: empRepo},
			Placements: stubPlacement{now: now},
			Policies:   stubPolicy{},
			Geofences:  stubGeofence{zones: zones},
			Biometric:  cfg.Biometric,
			Punches:    punchRepo,
			Validator:  domainpunch.PunchValidator{},
			Fraud:      fraud.FraudEvaluator{},
			Clock:      clock,
		},
		PunchRepo: punchRepo,
	}
}

func defaultZones() []geofence.GeofenceZone {
	center := geofence.GpsCoordinate{Latitude: -23.5505, Longitude: -46.6333}
	return []geofence.GeofenceZone{{
		ID: "site-1", Type: geofence.GeofenceTypeCircle,
		Center: &center, RadiusMeters: 500, AllowedDeviation: 50,
	}}
}

type employeeAdapter struct {
	repo *postgres.EmployeeRepository
}

func (a employeeAdapter) GetEmployee(
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

type stubPlacement struct {
	now time.Time
}

func (s stubPlacement) ActivePrimaryPlacement(
	_ context.Context, tenantID, employeeID uuid.UUID, _ time.Time,
) (*workforce.EmployeePlacement, error) {
	return &workforce.EmployeePlacement{
		ID: "pl-1", EmployeeID: employeeID.String(), TenantID: tenantID.String(),
		OrgNodeID: "site-1", Type: workforce.PlacementTypePrimary,
		ValidFrom: s.now.Add(-time.Hour),
	}, nil
}

type stubPolicy struct{}

func (stubPolicy) EffectivePolicy(
	_ context.Context, _ uuid.UUID, _ string,
) (organization.AttendancePolicy, error) {
	return organization.DefaultPolicy(), nil
}

type stubGeofence struct {
	zones []geofence.GeofenceZone
}

func (s stubGeofence) ZonesForOrgPath(
	_ context.Context, _ uuid.UUID, _ []string,
) ([]geofence.GeofenceZone, error) {
	return s.zones, nil
}
