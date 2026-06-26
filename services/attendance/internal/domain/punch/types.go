package punch

import (
	"errors"
	"time"
)

// PunchType is the punch action (see docs/GLOSSARY.md).
type PunchType string

const (
	PunchTypeClockIn    PunchType = "CLOCK_IN"
	PunchTypeClockOut   PunchType = "CLOCK_OUT"
	PunchTypeBreakStart PunchType = "BREAK_START"
	PunchTypeBreakEnd   PunchType = "BREAK_END"
)

// PunchStatus is the punch lifecycle state (see docs/DOMAIN-MODEL.md).
type PunchStatus string

const (
	PunchStatusValid      PunchStatus = "VALID"
	PunchStatusSuspicious PunchStatus = "SUSPICIOUS"
	PunchStatusRejected   PunchStatus = "REJECTED"
	PunchStatusPending    PunchStatus = "PENDING"
	PunchStatusDiscarded  PunchStatus = "DISCARDED"
)

// ValidationReason explains a rejected or discarded punch.
type ValidationReason string

const (
	ReasonLivenessFailed       ValidationReason = "LIVENESS_FAILED"
	ReasonFaceNotRecognized    ValidationReason = "FACE_NOT_RECOGNIZED"
	ReasonMockGPS              ValidationReason = "MOCK_GPS"
	ReasonOutOfGeofence        ValidationReason = "OUT_OF_GEOFENCE"
	ReasonClockManipulation    ValidationReason = "CLOCK_MANIPULATION"
	ReasonDuplicatePunch       ValidationReason = "DUPLICATE_PUNCH"
	ReasonInvalidSequence      ValidationReason = "INVALID_SEQUENCE"
	ReasonOfflineSyncExpired   ValidationReason = "OFFLINE_SYNC_EXPIRED"
)

// BiometricResult is the outcome of liveness and recognition for one punch.
type BiometricResult struct {
	LivenessScore         float64
	RecognitionConfidence float64
	FaceEmbeddingHash     string
	IsLive                bool
	IsRecognized          bool
}

// GpsCoordinate is a location snapshot at punch time.
type GpsCoordinate struct {
	Latitude  float64
	Longitude float64
	Accuracy  float64
	Altitude  float64
	Provider  string
	IsMocked  bool
}

// PunchRecord is the punch aggregate root.
type PunchRecord struct {
	ID              string
	EmployeeID      string
	TenantID        string
	PunchedAt       time.Time
	DeviceTime      time.Time
	Location        GpsCoordinate
	GeofenceID      string
	BiometricResult BiometricResult
	Status          PunchStatus
	Type            PunchType
}

const (
	minLivenessScore         = 0.80
	minRecognitionConfidence = 0.75
	maxClockSkewSeconds      = 300
	duplicateWindowSeconds   = 60
)

var ErrInvalidTransition = errors.New("punch: invalid status transition")
