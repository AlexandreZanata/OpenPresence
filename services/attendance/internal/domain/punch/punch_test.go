package punch

import (
	"errors"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)
	validator = PunchValidator{}
)

func validInput() PunchValidationInput {
	return PunchValidationInput{
		ID:         "punch-1",
		EmployeeID: "emp-1",
		TenantID:   "tenant-1",
		Type:       PunchTypeClockIn,
		ServerTime: baseTime,
		DeviceTime: baseTime,
		Location: GpsCoordinate{
			Latitude: -12.5458, Longitude: -55.7061, Accuracy: 5, Provider: "GPS",
		},
		Biometric: BiometricResult{
			LivenessScore: 0.92, RecognitionConfidence: 0.88,
			FaceEmbeddingHash: "abc123", IsLive: true, IsRecognized: true,
		},
		InsideGeofence:    true,
		MatchedGeofenceID: "zone-hq",
	}
}

func TestValidator_BR010_ValidPunch(t *testing.T) {
	result := validator.Validate(validInput())
	if result.Status != PunchStatusValid {
		t.Fatalf("BR-010: expected VALID, got %s reasons=%v", result.Status, result.Reasons)
	}
}

func TestValidator_BR010_RejectLiveness(t *testing.T) {
	input := validInput()
	input.Biometric.LivenessScore = 0.70
	input.Biometric.IsLive = false

	result := validator.Validate(input)
	if result.Status != PunchStatusRejected {
		t.Fatal("low liveness must reject")
	}
	if !containsReason(result.Reasons, ReasonLivenessFailed) {
		t.Fatal("expected LIVENESS_FAILED reason")
	}
}

func TestValidator_BR010_RejectMockGPS(t *testing.T) {
	input := validInput()
	input.Location.IsMocked = true

	result := validator.Validate(input)
	if !containsReason(result.Reasons, ReasonMockGPS) {
		t.Fatal("mock GPS must reject")
	}
}

func TestValidator_BR010_RejectOutOfGeofence(t *testing.T) {
	input := validInput()
	input.InsideGeofence = false

	result := validator.Validate(input)
	if !containsReason(result.Reasons, ReasonOutOfGeofence) {
		t.Fatal("outside geofence must reject")
	}
}

func TestValidator_BR014_InvalidSequence(t *testing.T) {
	input := validInput()
	input.Type = PunchTypeClockOut
	input.RecentPunches = []PunchRecord{}

	result := validator.Validate(input)
	if result.Status != PunchStatusRejected {
		t.Fatal("BR-014: CLOCK_OUT without CLOCK_IN must reject")
	}
	if !containsReason(result.Reasons, ReasonInvalidSequence) {
		t.Fatal("expected INVALID_SEQUENCE")
	}
}

func TestValidator_BR014_ValidSequence(t *testing.T) {
	history := []PunchRecord{
		{Type: PunchTypeClockIn, Status: PunchStatusValid, PunchedAt: baseTime},
	}
	input := validInput()
	input.Type = PunchTypeBreakStart
	input.ServerTime = baseTime.Add(2 * time.Hour)
	input.DeviceTime = input.ServerTime
	input.RecentPunches = history

	result := validator.Validate(input)
	if result.Status != PunchStatusValid {
		t.Fatalf("BR-014: valid BREAK_START after CLOCK_IN, got %s", result.Status)
	}
}

func TestValidator_BR015_ServerTimeOfficial(t *testing.T) {
	input := validInput()
	input.ServerTime = baseTime
	input.DeviceTime = baseTime.Add(2 * time.Minute)

	result := validator.Validate(input)
	if !result.Record.PunchedAt.Equal(input.ServerTime) {
		t.Fatal("BR-015: punchedAt must equal server time")
	}
	if !result.Record.DeviceTime.Equal(input.DeviceTime) {
		t.Fatal("BR-015: deviceTime must be stored separately for audit")
	}
}

func TestValidator_AntiDuplicateWithin60Seconds(t *testing.T) {
	input := validInput()
	input.Type = PunchTypeBreakStart
	input.ServerTime = baseTime.Add(30 * time.Second)
	input.DeviceTime = input.ServerTime
	input.RecentPunches = []PunchRecord{
		{Type: PunchTypeClockIn, Status: PunchStatusValid, PunchedAt: baseTime},
	}

	result := validator.Validate(input)
	if result.Status != PunchStatusRejected {
		t.Fatal("duplicate within 60s must reject")
	}
	if !containsReason(result.Reasons, ReasonDuplicatePunch) {
		t.Fatal("expected DUPLICATE_PUNCH")
	}
}

func TestValidator_ClockManipulationOver300Seconds(t *testing.T) {
	input := validInput()
	input.DeviceTime = baseTime.Add(-10 * time.Minute)

	result := validator.Validate(input)
	if result.Status != PunchStatusRejected {
		t.Fatal("clock skew > 300s must reject")
	}
	if !containsReason(result.Reasons, ReasonClockManipulation) {
		t.Fatal("expected CLOCK_MANIPULATION")
	}
}

func TestValidator_BR011_OfflineSyncExpired(t *testing.T) {
	queued := baseTime.Add(-9 * time.Hour)
	input := validInput()
	input.IsOfflineSync = true
	input.OfflineQueuedAt = &queued
	input.OfflineSyncMaxAge = 8 * time.Hour

	result := validator.Validate(input)
	if result.Status != PunchStatusDiscarded {
		t.Fatalf("BR-011: expired offline punch must be DISCARDED, got %s", result.Status)
	}
	if !containsReason(result.Reasons, ReasonOfflineSyncExpired) {
		t.Fatal("expected OFFLINE_SYNC_EXPIRED")
	}
}

func TestValidator_BR011_OfflineSyncWithinTTL(t *testing.T) {
	queued := baseTime.Add(-7 * time.Hour)
	input := validInput()
	input.IsOfflineSync = true
	input.OfflineQueuedAt = &queued
	input.OfflineSyncMaxAge = 8 * time.Hour

	result := validator.Validate(input)
	if result.Status != PunchStatusValid {
		t.Fatalf("BR-011: offline within TTL must validate, got %s", result.Status)
	}
}

func TestTransition_PendingToValid(t *testing.T) {
	next, err := Transition(PunchStatusPending, PunchStatusValid)
	if err != nil || next != PunchStatusValid {
		t.Fatalf("pending→valid allowed, got %v err=%v", next, err)
	}
}

func TestTransition_ValidIsTerminal(t *testing.T) {
	_, err := Transition(PunchStatusValid, PunchStatusRejected)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatal("VALID must be terminal")
	}
}

func containsReason(reasons []ValidationReason, target ValidationReason) bool {
	for _, r := range reasons {
		if r == target {
			return true
		}
	}
	return false
}
