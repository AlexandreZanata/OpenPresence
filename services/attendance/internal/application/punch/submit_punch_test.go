package punch_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
)

func TestSubmitPunchHandler_HappyPath_VALID(t *testing.T) {
	tenantID := uuid.New()
	employeeID := uuid.New()
	now := time.Date(2026, 6, 26, 8, 0, 0, 0, time.UTC)

	handler, _ := newTestHandler(t, tenantID, employeeID, testZone(), now)

	result, err := handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID:   tenantID,
		EmployeeID: employeeID,
		Type:       domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude:  -23.5505,
			Longitude: -46.6333,
			Accuracy:  10,
		},
		DeviceTime: now,
		FrameJPEG:  []byte("frame-bytes"),
	})
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
	require.Equal(t, domainpunch.PunchTypeClockIn, result.Record.Type)
}

func TestSubmitPunchHandler_OutOfGeofence_REJECTED(t *testing.T) {
	tenantID := uuid.New()
	employeeID := uuid.New()
	now := time.Date(2026, 6, 26, 8, 0, 0, 0, time.UTC)

	handler, repo := newTestHandler(t, tenantID, employeeID, testZone(), now)

	result, err := handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID:   tenantID,
		EmployeeID: employeeID,
		Type:       domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude:  -22.0,
			Longitude: -43.0,
			Accuracy:  10,
		},
		DeviceTime: now,
		FrameJPEG:  []byte("frame"),
	})
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonOutOfGeofence)

	validCount, err := repo.CountByStatus(context.Background(), tenantID, employeeID, domainpunch.PunchStatusValid)
	require.NoError(t, err)
	require.Equal(t, 0, validCount)
}

func TestSubmitPunchHandler_InvalidSequence_REJECTED(t *testing.T) {
	tenantID := uuid.New()
	employeeID := uuid.New()
	now := time.Date(2026, 6, 26, 8, 0, 0, 0, time.UTC)

	handler, repo := newTestHandler(t, tenantID, employeeID, testZone(), now)

	_, err := handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID: tenantID, EmployeeID: employeeID,
		Type: domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10,
		},
		DeviceTime: now, FrameJPEG: []byte("first"),
	})
	require.NoError(t, err)

	result, err := handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID: tenantID, EmployeeID: employeeID,
		Type: domainpunch.PunchTypeClockIn,
		Location: domainpunch.GpsCoordinate{
			Latitude: -23.5505, Longitude: -46.6333, Accuracy: 10,
		},
		DeviceTime: now.Add(time.Minute), FrameJPEG: []byte("second"),
	})
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonInvalidSequence)

	validCount, err := repo.CountByStatus(context.Background(), tenantID, employeeID, domainpunch.PunchStatusValid)
	require.NoError(t, err)
	require.Equal(t, 1, validCount)
}

func TestSubmitPunchHandler_EmployeeNotFound(t *testing.T) {
	tenantID := uuid.New()
	knownEmployee := uuid.New()
	handler, _ := newTestHandler(t, tenantID, knownEmployee, testZone(), time.Now().UTC())
	handler.Employees = &stubEmployeeReader{emp: nil}

	_, err := handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID:   tenantID,
		EmployeeID: uuid.New(),
		Type:       domainpunch.PunchTypeClockIn,
		Location:   domainpunch.GpsCoordinate{Latitude: -23.55, Longitude: -46.63},
		DeviceTime: time.Now().UTC(),
		FrameJPEG:  []byte("x"),
	})
	require.ErrorIs(t, err, apppunch.ErrEmployeeNotFound)
}

func newTestHandler(
	t *testing.T,
	tenantID, employeeID uuid.UUID,
	zone geofence.GeofenceZone,
	now time.Time,
) (apppunch.SubmitPunchHandler, *memoryPunchRepo) {
	t.Helper()
	repo := &memoryPunchRepo{}
	return apppunch.SubmitPunchHandler{
		Employees: &stubEmployeeReader{
			emp: &apppunch.Employee{ID: employeeID, TenantID: tenantID, Status: "ACTIVE"},
		},
		Placements: &stubPlacementReader{placement: &workforce.EmployeePlacement{
			ID: "pl-1", EmployeeID: employeeID.String(), TenantID: tenantID.String(),
			OrgNodeID: "site-1", Type: workforce.PlacementTypePrimary, ValidFrom: now.Add(-time.Hour),
		}},
		Policies:  &stubPolicyReader{policy: organization.DefaultPolicy()},
		Geofences: &stubGeofenceReader{zones: []geofence.GeofenceZone{zone}},
		Biometric: stubBiometricClient{},
		Punches:   repo,
		Validator: domainpunch.PunchValidator{},
		Fraud:     fraud.FraudEvaluator{},
		Clock:     func() time.Time { return now },
	}, repo
}

func testZone() geofence.GeofenceZone {
	center := geofence.GpsCoordinate{Latitude: -23.5505, Longitude: -46.6333}
	return geofence.GeofenceZone{
		ID: "zone-1", Type: geofence.GeofenceTypeCircle,
		Center: &center, RadiusMeters: 500, AllowedDeviation: 50,
	}
}

type stubEmployeeReader struct {
	emp *apppunch.Employee
}

func (s *stubEmployeeReader) GetEmployee(_ context.Context, _, _ uuid.UUID) (*apppunch.Employee, error) {
	return s.emp, nil
}

type stubPlacementReader struct {
	placement *workforce.EmployeePlacement
}

func (s *stubPlacementReader) ActivePrimaryPlacement(
	_ context.Context, _, _ uuid.UUID, _ time.Time,
) (*workforce.EmployeePlacement, error) {
	return s.placement, nil
}

type stubPolicyReader struct {
	policy organization.AttendancePolicy
}

func (s *stubPolicyReader) EffectivePolicy(
	_ context.Context, _ uuid.UUID, _ string,
) (organization.AttendancePolicy, error) {
	return s.policy, nil
}

type stubGeofenceReader struct {
	zones []geofence.GeofenceZone
}

func (s *stubGeofenceReader) ZonesForOrgPath(
	_ context.Context, _ uuid.UUID, _ []string,
) ([]geofence.GeofenceZone, error) {
	return s.zones, nil
}

type stubBiometricClient struct{}

func (stubBiometricClient) VerifyPunch(
	_ context.Context, _, _ uuid.UUID, _ []byte,
) (*apppunch.BiometricVerifyResult, error) {
	return &apppunch.BiometricVerifyResult{
		IsLive: true, LivenessScore: 0.95,
		IsRecognized: true, RecognitionConfidence: 0.90,
		EmbeddingHash: "abc123",
	}, nil
}

type memoryPunchRepo struct {
	mu      sync.Mutex
	records []domainpunch.PunchRecord
}

func (m *memoryPunchRepo) RecentPunches(
	_ context.Context, _, _ uuid.UUID, limit int,
) ([]domainpunch.PunchRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if limit <= 0 || limit > len(m.records) {
		limit = len(m.records)
	}
	out := make([]domainpunch.PunchRecord, limit)
	copy(out, m.records[len(m.records)-limit:])
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

func (m *memoryPunchRepo) Save(_ context.Context, _ uuid.UUID, record domainpunch.PunchRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, record)
	return nil
}

func (m *memoryPunchRepo) CountByStatus(
	_ context.Context, _, _ uuid.UUID, status domainpunch.PunchStatus,
) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for _, r := range m.records {
		if r.Status == status {
			count++
		}
	}
	return count, nil
}
