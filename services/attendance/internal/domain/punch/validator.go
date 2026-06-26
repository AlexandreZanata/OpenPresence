package punch

import "time"

// PunchValidationInput bundles all data required to validate a punch submission.
type PunchValidationInput struct {
	ID                string
	EmployeeID        string
	TenantID          string
	Type              PunchType
	ServerTime        time.Time
	DeviceTime        time.Time
	Location          GpsCoordinate
	Biometric         BiometricResult
	InsideGeofence    bool
	MatchedGeofenceID string
	RecentPunches     []PunchRecord
	IsOfflineSync     bool
	OfflineQueuedAt   *time.Time
	OfflineSyncMaxAge time.Duration
}

// ValidationResult is the outcome of punch validation.
type ValidationResult struct {
	Status  PunchStatus
	Record  PunchRecord
	Reasons []ValidationReason
}

// PunchValidator evaluates punch submissions against BR-010–BR-015.
type PunchValidator struct{}

// Validate applies business rules and returns status with a draft PunchRecord.
func (PunchValidator) Validate(input PunchValidationInput) ValidationResult {
	record := buildRecord(input)
	reasons := collectRejectReasons(input)

	if input.IsOfflineSync {
		if expired := offlineSyncExpired(input); expired != "" {
			record.Status = PunchStatusDiscarded
			return ValidationResult{Status: PunchStatusDiscarded, Record: record, Reasons: []ValidationReason{expired}}
		}
	}

	if !isValidSequence(input.RecentPunches, input.Type) {
		record.Status = PunchStatusRejected
		return ValidationResult{Status: PunchStatusRejected, Record: record, Reasons: []ValidationReason{ReasonInvalidSequence}}
	}

	if isDuplicate(input.RecentPunches, input.ServerTime) {
		record.Status = PunchStatusRejected
		return ValidationResult{Status: PunchStatusRejected, Record: record, Reasons: []ValidationReason{ReasonDuplicatePunch}}
	}

	if len(reasons) > 0 {
		record.Status = PunchStatusRejected
		return ValidationResult{Status: PunchStatusRejected, Record: record, Reasons: reasons}
	}

	record.Status = PunchStatusValid
	return ValidationResult{Status: PunchStatusValid, Record: record, Reasons: nil}
}

func buildRecord(input PunchValidationInput) PunchRecord {
	return PunchRecord{
		ID:              input.ID,
		EmployeeID:      input.EmployeeID,
		TenantID:        input.TenantID,
		PunchedAt:       input.ServerTime,
		DeviceTime:      input.DeviceTime,
		Location:        input.Location,
		GeofenceID:      input.MatchedGeofenceID,
		BiometricResult: input.Biometric,
		Type:            input.Type,
	}
}

func collectRejectReasons(input PunchValidationInput) []ValidationReason {
	var reasons []ValidationReason
	if input.Biometric.LivenessScore < minLivenessScore || !input.Biometric.IsLive {
		reasons = append(reasons, ReasonLivenessFailed)
	}
	if input.Biometric.RecognitionConfidence < minRecognitionConfidence || !input.Biometric.IsRecognized {
		reasons = append(reasons, ReasonFaceNotRecognized)
	}
	if input.Location.IsMocked {
		reasons = append(reasons, ReasonMockGPS)
	}
	if !input.InsideGeofence {
		reasons = append(reasons, ReasonOutOfGeofence)
	}
	if clockSkewSeconds(input.ServerTime, input.DeviceTime) > maxClockSkewSeconds {
		reasons = append(reasons, ReasonClockManipulation)
	}
	return reasons
}

func offlineSyncExpired(input PunchValidationInput) ValidationReason {
	if input.OfflineQueuedAt == nil || input.OfflineSyncMaxAge == 0 {
		return ReasonOfflineSyncExpired
	}
	deadline := input.OfflineQueuedAt.Add(input.OfflineSyncMaxAge)
	if input.ServerTime.After(deadline) {
		return ReasonOfflineSyncExpired
	}
	return ""
}

func isDuplicate(history []PunchRecord, at time.Time) bool {
	window := time.Duration(duplicateWindowSeconds) * time.Second
	for _, prior := range history {
		if prior.Status != PunchStatusValid && prior.Status != PunchStatusSuspicious {
			continue
		}
		delta := at.Sub(prior.PunchedAt)
		if delta >= 0 && delta < window {
			return true
		}
	}
	return false
}

func clockSkewSeconds(server, device time.Time) int64 {
	skew := server.Sub(device)
	if skew < 0 {
		skew = -skew
	}
	return int64(skew.Seconds())
}
