package punch

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

// SubmitPunchCommand carries a punch submission request from the delivery layer.
type SubmitPunchCommand struct {
	TenantID      uuid.UUID
	EmployeeID    uuid.UUID
	Type          domainpunch.PunchType
	Location      domainpunch.GpsCoordinate
	DeviceTime    time.Time
	FrameJPEG     []byte
	DeviceReport  fraud.DeviceIntegrityReport
	DeviceID      string
	IsOfflineSync bool
	OfflineQueued *time.Time
}

// SubmitPunchResult is the outcome of punch orchestration.
type SubmitPunchResult struct {
	Record     domainpunch.PunchRecord
	FraudFlags []fraud.FraudFlag
	Reasons    []domainpunch.ValidationReason
}

// SubmitPunchHandler orchestrates placement, policy, geofence, biometric, validation, and persistence.
type SubmitPunchHandler struct {
	Employees  EmployeeReader
	Placements PlacementReader
	Policies   PolicyResolver
	Geofences  GeofenceResolver
	Biometric  BiometricClient
	Punches    PunchRepository
	Checker    geofence.GeofenceChecker
	Validator  domainpunch.PunchValidator
	Fraud      fraud.FraudEvaluator
	Lockout    *fraud.DeviceLockoutTracker
	Clock      func() time.Time
}

// Handle runs the SubmitPunch use case end-to-end.
func (h SubmitPunchHandler) Handle(ctx context.Context, cmd SubmitPunchCommand) (*SubmitPunchResult, error) {
	serverTime := h.now()

	if h.isDeviceLocked(cmd.DeviceID, serverTime) {
		return nil, ErrDeviceLocked
	}

	emp, err := h.Employees.GetEmployee(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("load employee: %w", err)
	}
	if emp == nil {
		return nil, ErrEmployeeNotFound
	}
	if emp.Status != "" && emp.Status != "ACTIVE" {
		return nil, ErrEmployeeInactive
	}

	placement, err := h.Placements.ActivePrimaryPlacement(ctx, cmd.TenantID, cmd.EmployeeID, serverTime)
	if err != nil {
		return nil, fmt.Errorf("load placement: %w", err)
	}
	if placement == nil {
		return nil, ErrNoActivePlacement
	}

	policy, err := h.Policies.EffectivePolicy(ctx, cmd.TenantID, placement.OrgNodeID)
	if err != nil {
		return nil, fmt.Errorf("resolve policy: %w", err)
	}
	if !punchTypeAllowed(policy, cmd.Type) {
		return nil, ErrPunchTypeNotAllowed
	}

	zones, err := h.Geofences.ZonesForOrgPath(ctx, cmd.TenantID, []string{placement.OrgNodeID})
	if err != nil {
		return nil, fmt.Errorf("resolve geofences: %w", err)
	}

	inside, matchedID := checkGeofence(h.checker(), cmd.Location, zones)

	biometric, err := h.resolveBiometric(ctx, cmd, policy)
	if err != nil {
		return nil, err
	}

	recent, err := h.Punches.RecentPunches(ctx, cmd.TenantID, cmd.EmployeeID, 20)
	if err != nil {
		return nil, fmt.Errorf("load recent punches: %w", err)
	}
	chronological := reversePunchHistory(recent)

	punchID := uuid.New().String()
	validation := h.validator().Validate(domainpunch.PunchValidationInput{
		ID:                punchID,
		EmployeeID:        cmd.EmployeeID.String(),
		TenantID:          cmd.TenantID.String(),
		Type:              cmd.Type,
		ServerTime:        serverTime,
		DeviceTime:        cmd.DeviceTime,
		Location:          cmd.Location,
		Biometric:         biometric,
		InsideGeofence:    inside,
		MatchedGeofenceID: matchedID,
		RecentPunches:     chronological,
		IsOfflineSync:     cmd.IsOfflineSync,
		OfflineQueuedAt:   cmd.OfflineQueued,
		OfflineSyncMaxAge: policy.OfflineSyncMaxAge,
	})

	fraudResult := h.fraudEval().Evaluate(buildFraudInput(cmd, serverTime, inside, biometric, chronological, zones))

	record := validation.Record
	record.Status = mergeStatus(validation.Status, fraudResult.RecommendedStatus)

	if err := h.Punches.Save(ctx, cmd.TenantID, record); err != nil {
		return nil, fmt.Errorf("persist punch: %w", err)
	}

	h.recordDeviceRejection(cmd.DeviceID, serverTime, record.Status)

	return &SubmitPunchResult{
		Record:     record,
		FraudFlags: fraudResult.Flags,
		Reasons:    validation.Reasons,
	}, nil
}

func (h SubmitPunchHandler) now() time.Time {
	if h.Clock != nil {
		return h.Clock()
	}
	return time.Now().UTC()
}

func (h SubmitPunchHandler) checker() geofence.GeofenceChecker {
	if h.Checker != nil {
		return h.Checker
	}
	return geofence.NewChecker()
}

func (h SubmitPunchHandler) validator() domainpunch.PunchValidator {
	return h.Validator
}

func (h SubmitPunchHandler) fraudEval() fraud.FraudEvaluator {
	return h.Fraud
}

func (h SubmitPunchHandler) isDeviceLocked(deviceID string, at time.Time) bool {
	if h.Lockout == nil || deviceID == "" {
		return false
	}
	return h.Lockout.IsLocked(deviceID, at)
}

func (h SubmitPunchHandler) recordDeviceRejection(deviceID string, at time.Time, status domainpunch.PunchStatus) {
	if h.Lockout == nil || deviceID == "" || status != domainpunch.PunchStatusRejected {
		return
	}
	h.Lockout.RecordRejected(deviceID, at)
}

func (h SubmitPunchHandler) resolveBiometric(
	ctx context.Context,
	cmd SubmitPunchCommand,
	policy organization.AttendancePolicy,
) (domainpunch.BiometricResult, error) {
	if !policy.BiometricRequired {
		return domainpunch.BiometricResult{
			LivenessScore:         1.0,
			RecognitionConfidence: 1.0,
			IsLive:                true,
			IsRecognized:          true,
		}, nil
	}
	if h.Biometric == nil {
		return domainpunch.BiometricResult{}, ErrBiometricRequired
	}
	result, err := h.Biometric.VerifyPunch(ctx, cmd.TenantID, cmd.EmployeeID, cmd.FrameJPEG)
	if err != nil {
		return domainpunch.BiometricResult{}, fmt.Errorf("biometric verify: %w", err)
	}
	hash := result.EmbeddingHash
	if hash == "" {
		hash = hashFrame(cmd.FrameJPEG)
	}
	return domainpunch.BiometricResult{
		LivenessScore:         result.LivenessScore,
		RecognitionConfidence: result.RecognitionConfidence,
		IsLive:                result.IsLive,
		IsRecognized:          result.IsRecognized,
		FaceEmbeddingHash:     hash,
	}, nil
}
