package fraud

import "time"

// FraudType classifies integrity violations (see docs/GLOSSARY.md).
type FraudType string

const (
	FraudTypeMockGPS            FraudType = "MOCK_GPS"
	FraudTypeClockManipulation  FraudType = "CLOCK_MANIPULATION"
	FraudTypeImpossibleSpeed    FraudType = "IMPOSSIBLE_SPEED"
	FraudTypeLivenessFailed     FraudType = "LIVENESS_FAILED"
	FraudTypeFaceNotRecognized  FraudType = "FACE_NOT_RECOGNIZED"
	FraudTypeOutOfGeofence      FraudType = "OUT_OF_GEOFENCE"
	FraudTypeDuplicatePunch     FraudType = "DUPLICATE_PUNCH"
	FraudTypeDeviceRooted       FraudType = "DEVICE_ROOTED"
	FraudTypeVPNDetected         FraudType = "VPN_DETECTED"
	FraudTypeGPSLowAccuracy     FraudType = "GPS_LOW_ACCURACY"
)

// FraudSeverity ranks anomaly impact (see docs/GLOSSARY.md).
type FraudSeverity string

const (
	FraudSeverityLow      FraudSeverity = "LOW"
	FraudSeverityMedium   FraudSeverity = "MEDIUM"
	FraudSeverityHigh     FraudSeverity = "HIGH"
	FraudSeverityCritical FraudSeverity = "CRITICAL"
)

// FraudFlag is a detected anomaly attached to a punch attempt.
type FraudFlag struct {
	Type       FraudType
	Severity   FraudSeverity
	DetectedAt time.Time
	Metadata   map[string]string
}

// GpsSnapshot captures punch-time location for fraud checks.
type GpsSnapshot struct {
	Latitude  float64
	Longitude float64
	Accuracy  float64
	IsMocked  bool
}

// BiometricSnapshot captures biometric outcomes for fraud checks.
type BiometricSnapshot struct {
	LivenessScore         float64
	RecognitionConfidence float64
	IsLive                bool
	IsRecognized          bool
}

// DeviceIntegrityReport is the mobile device payload at punch time.
type DeviceIntegrityReport struct {
	IsRooted  bool
	VPNActive bool
}

// PriorPunch is minimal history for speed and duplicate detection.
type PriorPunch struct {
	PunchedAt time.Time
	Location  GpsSnapshot
	Status    string
}

const (
	clockSkewSuspiciousSeconds = 300
	clockSkewCriticalSeconds   = 1800
	impossibleSpeedKmh         = 600.0
	duplicateWindowSeconds     = 60
)
