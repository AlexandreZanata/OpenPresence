package fraud

import (
	"strconv"
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

// EvaluateInput bundles signals for fraud evaluation.
type EvaluateInput struct {
	ServerTime       time.Time
	DeviceTime       time.Time
	Location         GpsSnapshot
	Biometric        BiometricSnapshot
	Device           DeviceIntegrityReport
	InsideGeofence   bool
	AllowedDeviation float64
	LastPunch        *PriorPunch
	RecentAttempts   []PriorPunch
}

// EvaluationResult holds detected flags and recommended punch status.
type EvaluationResult struct {
	Flags             []FraudFlag
	RecommendedStatus punch.PunchStatus
}

// FraudEvaluator classifies anomalies and recommends punch status (BR-012).
type FraudEvaluator struct{}

// Evaluate detects fraud flags and derives recommended status.
func (FraudEvaluator) Evaluate(input EvaluateInput) EvaluationResult {
	flags := collectFlags(input)
	return EvaluationResult{
		Flags:             flags,
		RecommendedStatus: recommendStatus(flags),
	}
}

func collectFlags(input EvaluateInput) []FraudFlag {
	at := input.ServerTime
	var flags []FraudFlag

	if input.Location.IsMocked {
		flags = append(flags, newFlag(FraudTypeMockGPS, FraudSeverityHigh, at, nil))
	}
	if skew := clockSkewSeconds(input.ServerTime, input.DeviceTime); skew > clockSkewSuspiciousSeconds {
		severity := FraudSeverityMedium
		if skew > clockSkewCriticalSeconds {
			severity = FraudSeverityCritical
		}
		flags = append(flags, newFlag(FraudTypeClockManipulation, severity, at, map[string]string{
			"skew_seconds": strconv.FormatInt(skew, 10),
		}))
	}
	if flag := detectImpossibleSpeed(input); flag != nil {
		flags = append(flags, *flag)
	}
	if isDuplicate(input.RecentAttempts, at) {
		flags = append(flags, newFlag(FraudTypeDuplicatePunch, FraudSeverityMedium, at, nil))
	}
	if input.AllowedDeviation > 0 && input.Location.Accuracy > input.AllowedDeviation*2 {
		flags = append(flags, newFlag(FraudTypeGPSLowAccuracy, FraudSeverityLow, at, map[string]string{
			"accuracy_m": strconv.FormatFloat(input.Location.Accuracy, 'f', 1, 64),
		}))
	}
	if !input.InsideGeofence {
		flags = append(flags, newFlag(FraudTypeOutOfGeofence, FraudSeverityHigh, at, nil))
	}
	if !input.Biometric.IsLive || input.Biometric.LivenessScore < 0.80 {
		flags = append(flags, newFlag(FraudTypeLivenessFailed, FraudSeverityHigh, at, nil))
	}
	if !input.Biometric.IsRecognized || input.Biometric.RecognitionConfidence < 0.75 {
		flags = append(flags, newFlag(FraudTypeFaceNotRecognized, FraudSeverityHigh, at, nil))
	}
	if input.Device.IsRooted {
		flags = append(flags, newFlag(FraudTypeDeviceRooted, FraudSeverityMedium, at, nil))
	}
	if input.Device.VPNActive {
		flags = append(flags, newFlag(FraudTypeVPNDetected, FraudSeverityLow, at, nil))
	}
	return flags
}

func detectImpossibleSpeed(input EvaluateInput) *FraudFlag {
	if input.LastPunch == nil {
		return nil
	}
	elapsed := input.ServerTime.Sub(input.LastPunch.PunchedAt)
	if elapsed <= 0 {
		return nil
	}
	from := geofence.GpsCoordinate{
		Latitude:  input.LastPunch.Location.Latitude,
		Longitude: input.LastPunch.Location.Longitude,
	}
	to := geofence.GpsCoordinate{Latitude: input.Location.Latitude, Longitude: input.Location.Longitude}
	distanceM := geofence.HaversineDistance(from, to)
	speedKmh := (distanceM / 1000) / elapsed.Hours()
	if speedKmh <= impossibleSpeedKmh {
		return nil
	}
	flag := newFlag(FraudTypeImpossibleSpeed, FraudSeverityCritical, input.ServerTime, map[string]string{
		"speed_kmh": strconv.FormatFloat(speedKmh, 'f', 1, 64),
	})
	return &flag
}

func recommendStatus(flags []FraudFlag) punch.PunchStatus {
	if len(flags) == 0 {
		return punch.PunchStatusValid
	}
	for _, flag := range flags {
		if flag.Severity == FraudSeverityCritical {
			return punch.PunchStatusRejected
		}
	}
	if onlyGPSLowAccuracy(flags) {
		return punch.PunchStatusValid
	}
	return punch.PunchStatusSuspicious
}

func onlyGPSLowAccuracy(flags []FraudFlag) bool {
	if len(flags) != 1 {
		return false
	}
	return flags[0].Type == FraudTypeGPSLowAccuracy && flags[0].Severity == FraudSeverityLow
}

func isDuplicate(recent []PriorPunch, at time.Time) bool {
	window := time.Duration(duplicateWindowSeconds) * time.Second
	for _, prior := range recent {
		if prior.Status != string(punch.PunchStatusValid) && prior.Status != string(punch.PunchStatusSuspicious) {
			continue
		}
		delta := at.Sub(prior.PunchedAt)
		if delta >= 0 && delta < window {
			return true
		}
	}
	return false
}

func newFlag(t FraudType, s FraudSeverity, at time.Time, meta map[string]string) FraudFlag {
	return FraudFlag{Type: t, Severity: s, DetectedAt: at, Metadata: meta}
}

func clockSkewSeconds(server, device time.Time) int64 {
	skew := server.Sub(device)
	if skew < 0 {
		skew = -skew
	}
	return int64(skew.Seconds())
}
