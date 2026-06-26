package punch

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

func punchTypeAllowed(policy organization.AttendancePolicy, t domainpunch.PunchType) bool {
	for _, allowed := range policy.AllowedPunchTypes {
		if organization.PunchType(t) == allowed {
			return true
		}
	}
	return false
}

func checkGeofence(
	checker geofence.GeofenceChecker,
	loc domainpunch.GpsCoordinate,
	zones []geofence.GeofenceZone,
) (inside bool, matchedID string) {
	coord := geofence.GpsCoordinate{Latitude: loc.Latitude, Longitude: loc.Longitude}
	for _, zone := range zones {
		if checker.IsInsideZone(coord, zone) {
			return true, zone.ID
		}
	}
	return false, ""
}

func mergeStatus(validated, fraudRecommended domainpunch.PunchStatus) domainpunch.PunchStatus {
	if validated != domainpunch.PunchStatusValid {
		return validated
	}
	if fraudRecommended == domainpunch.PunchStatusRejected {
		return domainpunch.PunchStatusRejected
	}
	if fraudRecommended == domainpunch.PunchStatusSuspicious {
		return domainpunch.PunchStatusSuspicious
	}
	return domainpunch.PunchStatusValid
}

func buildFraudInput(
	cmd SubmitPunchCommand,
	serverTime time.Time,
	inside bool,
	biometric domainpunch.BiometricResult,
	recent []domainpunch.PunchRecord,
	zones []geofence.GeofenceZone,
) fraud.EvaluateInput {
	deviation := maxAllowedDeviation(zones)
	return fraud.EvaluateInput{
		ServerTime: serverTime,
		DeviceTime: cmd.DeviceTime,
		Location: fraud.GpsSnapshot{
			Latitude:  cmd.Location.Latitude,
			Longitude: cmd.Location.Longitude,
			Accuracy:  cmd.Location.Accuracy,
			IsMocked:  cmd.Location.IsMocked,
		},
		Biometric: fraud.BiometricSnapshot{
			LivenessScore:         biometric.LivenessScore,
			RecognitionConfidence: biometric.RecognitionConfidence,
			IsLive:                biometric.IsLive,
			IsRecognized:          biometric.IsRecognized,
		},
		Device:           cmd.DeviceReport,
		InsideGeofence:   inside,
		AllowedDeviation: deviation,
		LastPunch:        priorPunch(recent),
		RecentAttempts:   priorAttempts(recent),
	}
}

func maxAllowedDeviation(zones []geofence.GeofenceZone) float64 {
	var maxDev float64
	for _, zone := range zones {
		if zone.AllowedDeviation > maxDev {
			maxDev = zone.AllowedDeviation
		}
	}
	return maxDev
}

func priorPunch(recent []domainpunch.PunchRecord) *fraud.PriorPunch {
	if len(recent) == 0 {
		return nil
	}
	last := recent[0]
	return &fraud.PriorPunch{
		PunchedAt: last.PunchedAt,
		Location: fraud.GpsSnapshot{
			Latitude:  last.Location.Latitude,
			Longitude: last.Location.Longitude,
		},
		Status: string(last.Status),
	}
}

func priorAttempts(recent []domainpunch.PunchRecord) []fraud.PriorPunch {
	out := make([]fraud.PriorPunch, 0, len(recent))
	for _, r := range recent {
		out = append(out, fraud.PriorPunch{
			PunchedAt: r.PunchedAt,
			Location: fraud.GpsSnapshot{
				Latitude:  r.Location.Latitude,
				Longitude: r.Location.Longitude,
			},
			Status: string(r.Status),
		})
	}
	return out
}

func hashFrame(frame []byte) string {
	sum := sha256.Sum256(frame)
	return hex.EncodeToString(sum[:16])
}

func reversePunchHistory(recent []domainpunch.PunchRecord) []domainpunch.PunchRecord {
	if len(recent) <= 1 {
		return recent
	}
	out := make([]domainpunch.PunchRecord, len(recent))
	for i, r := range recent {
		out[len(recent)-1-i] = r
	}
	return out
}
