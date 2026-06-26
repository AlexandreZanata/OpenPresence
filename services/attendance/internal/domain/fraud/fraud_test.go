package fraud

import (
	"testing"
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

var (
	evaluator = FraudEvaluator{}
	baseTime  = time.Date(2026, 6, 26, 9, 0, 0, 0, time.UTC)
)

func cleanInput() EvaluateInput {
	return EvaluateInput{
		ServerTime: baseTime,
		DeviceTime: baseTime,
		Location: GpsSnapshot{
			Latitude: -12.5458, Longitude: -55.7061, Accuracy: 5,
		},
		Biometric: BiometricSnapshot{
			LivenessScore: 0.9, RecognitionConfidence: 0.85, IsLive: true, IsRecognized: true,
		},
		InsideGeofence:   true,
		AllowedDeviation: 10,
	}
}

func TestEvaluate_MockGPS_HighSeverity(t *testing.T) {
	input := cleanInput()
	input.Location.IsMocked = true

	result := evaluator.Evaluate(input)
	flag := findFlag(result.Flags, FraudTypeMockGPS)
	if flag == nil || flag.Severity != FraudSeverityHigh {
		t.Fatal("MOCK_GPS must be HIGH severity")
	}
}

func TestEvaluate_ClockManipulation_MediumOver300s(t *testing.T) {
	input := cleanInput()
	input.DeviceTime = baseTime.Add(-6 * time.Minute)

	result := evaluator.Evaluate(input)
	flag := findFlag(result.Flags, FraudTypeClockManipulation)
	if flag == nil || flag.Severity != FraudSeverityMedium {
		t.Fatal("clock skew > 300s must be MEDIUM")
	}
}

func TestEvaluate_ClockManipulation_CriticalOver30Min(t *testing.T) {
	input := cleanInput()
	input.DeviceTime = baseTime.Add(-31 * time.Minute)

	result := evaluator.Evaluate(input)
	if result.RecommendedStatus != punch.PunchStatusRejected {
		t.Fatal("BR-012: CRITICAL clock manipulation must reject")
	}
}

func TestEvaluate_ImpossibleSpeed_Critical(t *testing.T) {
	input := cleanInput()
	input.ServerTime = baseTime.Add(30 * time.Second)
	input.Location = GpsSnapshot{Latitude: 0, Longitude: 0, Accuracy: 5}
	input.LastPunch = &PriorPunch{
		PunchedAt: baseTime,
		Location:  GpsSnapshot{Latitude: 10, Longitude: 10, Accuracy: 5},
		Status:    string(punch.PunchStatusValid),
	}

	result := evaluator.Evaluate(input)
	flag := findFlag(result.Flags, FraudTypeImpossibleSpeed)
	if flag == nil || flag.Severity != FraudSeverityCritical {
		t.Fatal("IMPOSSIBLE_SPEED must be CRITICAL")
	}
}

func TestEvaluate_DuplicatePunch_Within60s(t *testing.T) {
	input := cleanInput()
	input.RecentAttempts = []PriorPunch{
		{PunchedAt: baseTime.Add(-30 * time.Second), Status: string(punch.PunchStatusValid)},
	}

	result := evaluator.Evaluate(input)
	if findFlag(result.Flags, FraudTypeDuplicatePunch) == nil {
		t.Fatal("DUPLICATE_PUNCH within 60s must be flagged")
	}
}

func TestEvaluate_GPSLowAccuracy_LowAcceptWithFlag(t *testing.T) {
	input := cleanInput()
	input.Location.Accuracy = 25
	input.AllowedDeviation = 10

	result := evaluator.Evaluate(input)
	flag := findFlag(result.Flags, FraudTypeGPSLowAccuracy)
	if flag == nil || flag.Severity != FraudSeverityLow {
		t.Fatal("GPS_LOW_ACCURACY must be LOW severity")
	}
	if result.RecommendedStatus != punch.PunchStatusValid {
		t.Fatal("BR-022: low accuracy alone should keep VALID with flag")
	}
}

func TestEvaluate_BR012_SuspiciousWhenNotCritical(t *testing.T) {
	input := cleanInput()
	input.Device.VPNActive = true

	result := evaluator.Evaluate(input)
	if result.RecommendedStatus != punch.PunchStatusSuspicious {
		t.Fatalf("BR-012: non-critical flags → SUSPICIOUS, got %s", result.RecommendedStatus)
	}
}

func TestEvaluate_AllGlossaryFraudTypes(t *testing.T) {
	cases := []struct {
		name     string
		mutate   func(*EvaluateInput)
		wantType FraudType
	}{
		{"MOCK_GPS", func(i *EvaluateInput) { i.Location.IsMocked = true }, FraudTypeMockGPS},
		{"CLOCK_MANIPULATION", func(i *EvaluateInput) { i.DeviceTime = i.ServerTime.Add(-10 * time.Minute) }, FraudTypeClockManipulation},
		{"IMPOSSIBLE_SPEED", func(i *EvaluateInput) {
			i.ServerTime = baseTime.Add(30 * time.Second)
			i.Location = GpsSnapshot{Latitude: 0, Longitude: 0}
			i.LastPunch = &PriorPunch{PunchedAt: baseTime, Location: GpsSnapshot{Latitude: 10, Longitude: 10}, Status: string(punch.PunchStatusValid)}
		}, FraudTypeImpossibleSpeed},
		{"LIVENESS_FAILED", func(i *EvaluateInput) { i.Biometric.IsLive = false }, FraudTypeLivenessFailed},
		{"FACE_NOT_RECOGNIZED", func(i *EvaluateInput) { i.Biometric.IsRecognized = false }, FraudTypeFaceNotRecognized},
		{"OUT_OF_GEOFENCE", func(i *EvaluateInput) { i.InsideGeofence = false }, FraudTypeOutOfGeofence},
		{"DUPLICATE_PUNCH", func(i *EvaluateInput) {
			i.RecentAttempts = []PriorPunch{{PunchedAt: baseTime.Add(-20 * time.Second), Status: string(punch.PunchStatusValid)}}
		}, FraudTypeDuplicatePunch},
		{"DEVICE_ROOTED", func(i *EvaluateInput) { i.Device.IsRooted = true }, FraudTypeDeviceRooted},
		{"VPN_DETECTED", func(i *EvaluateInput) { i.Device.VPNActive = true }, FraudTypeVPNDetected},
		{"GPS_LOW_ACCURACY", func(i *EvaluateInput) { i.Location.Accuracy = 30; i.AllowedDeviation = 10 }, FraudTypeGPSLowAccuracy},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := cleanInput()
			tc.mutate(&input)
			result := evaluator.Evaluate(input)
			if findFlag(result.Flags, tc.wantType) == nil {
				t.Fatalf("expected flag %s", tc.wantType)
			}
		})
	}
}

func TestDeviceLockoutTracker_BR013_ThreeRejectsInTenMinutes(t *testing.T) {
	tracker := NewDeviceLockoutTracker()
	device := "device-1"
	t1 := baseTime
	t2 := t1.Add(3 * time.Minute)
	t3 := t2.Add(3 * time.Minute)

	tracker.RecordRejected(device, t1)
	tracker.RecordRejected(device, t2)
	if tracker.IsLocked(device, t2) {
		t.Fatal("lockout must not trigger before third reject")
	}
	tracker.RecordRejected(device, t3)
	if !tracker.IsLocked(device, t3) {
		t.Fatal("BR-013: third reject within 10 minutes must lock device")
	}
	until, ok := tracker.LockedUntil(device)
	if !ok || !until.Equal(t3.Add(30*time.Minute)) {
		t.Fatalf("lockout must last 30 minutes, until=%v", until)
	}
}

func findFlag(flags []FraudFlag, target FraudType) *FraudFlag {
	for i := range flags {
		if flags[i].Type == target {
			return &flags[i]
		}
	}
	return nil
}
