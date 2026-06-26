package organization

import "time"

// PunchType is a allowed punch action (see docs/GLOSSARY.md).
type PunchType string

const (
	PunchTypeClockIn    PunchType = "CLOCK_IN"
	PunchTypeClockOut   PunchType = "CLOCK_OUT"
	PunchTypeBreakStart PunchType = "BREAK_START"
	PunchTypeBreakEnd   PunchType = "BREAK_END"
)

// OvertimePolicy controls overtime handling (see docs/BUSINESS-RULES.md BR-032, BR-034).
type OvertimePolicy string

const (
	OvertimePolicyDisabled OvertimePolicy = "DISABLED"
	OvertimePolicyStandard OvertimePolicy = "STANDARD"
	OvertimePolicyTimeBank OvertimePolicy = "TIME_BANK"
)

// AttendancePolicy is the effective rules for punch validation on an org node.
type AttendancePolicy struct {
	WorkdayDuration   time.Duration
	ToleranceMinutes  int
	AllowedPunchTypes []PunchType
	GeofenceRequired  bool
	BiometricRequired bool
	OfflineSyncMaxAge time.Duration
	OvertimePolicy    OvertimePolicy
}

// PolicyOverride holds optional field overrides. Nil pointer fields inherit from parent.
// Non-nil pointers apply explicit values, including false and zero durations.
type PolicyOverride struct {
	WorkdayDuration   *time.Duration
	ToleranceMinutes  *int
	AllowedPunchTypes *[]PunchType
	GeofenceRequired  *bool
	BiometricRequired *bool
	OfflineSyncMaxAge *time.Duration
	OvertimePolicy    *OvertimePolicy
}

// DefaultPolicy returns baseline tenant policy. OfflineSyncMaxAge follows BR-011 (8h).
func DefaultPolicy() AttendancePolicy {
	return AttendancePolicy{
		WorkdayDuration:   8 * time.Hour,
		ToleranceMinutes:  10,
		AllowedPunchTypes: allPunchTypes(),
		GeofenceRequired:  true,
		BiometricRequired: true,
		OfflineSyncMaxAge: 8 * time.Hour,
		OvertimePolicy:    OvertimePolicyStandard,
	}
}

// PublicSectorPreset targets municipalities: strict geofence, 12×36 shift tolerance.
func PublicSectorPreset() AttendancePolicy {
	return AttendancePolicy{
		WorkdayDuration:   12 * time.Hour,
		ToleranceMinutes:  30,
		AllowedPunchTypes: allPunchTypes(),
		GeofenceRequired:  true,
		BiometricRequired: true,
		OfflineSyncMaxAge: 8 * time.Hour,
		OvertimePolicy:    OvertimePolicyTimeBank,
	}
}

// PrivateSectorPreset targets companies: flexible tolerance windows and standard overtime.
func PrivateSectorPreset() AttendancePolicy {
	return AttendancePolicy{
		WorkdayDuration:   8 * time.Hour,
		ToleranceMinutes:  15,
		AllowedPunchTypes: allPunchTypes(),
		GeofenceRequired:  true,
		BiometricRequired: true,
		OfflineSyncMaxAge: 8 * time.Hour,
		OvertimePolicy:    OvertimePolicyStandard,
	}
}

// MergePolicy applies override on top of parent. Nil override fields inherit parent values.
func MergePolicy(parent AttendancePolicy, override PolicyOverride) AttendancePolicy {
	merged := parent
	if override.WorkdayDuration != nil {
		merged.WorkdayDuration = *override.WorkdayDuration
	}
	if override.ToleranceMinutes != nil {
		merged.ToleranceMinutes = *override.ToleranceMinutes
	}
	if override.AllowedPunchTypes != nil {
		merged.AllowedPunchTypes = append([]PunchType(nil), (*override.AllowedPunchTypes)...)
	}
	if override.GeofenceRequired != nil {
		merged.GeofenceRequired = *override.GeofenceRequired
	}
	if override.BiometricRequired != nil {
		merged.BiometricRequired = *override.BiometricRequired
	}
	if override.OfflineSyncMaxAge != nil {
		merged.OfflineSyncMaxAge = *override.OfflineSyncMaxAge
	}
	if override.OvertimePolicy != nil {
		merged.OvertimePolicy = *override.OvertimePolicy
	}
	return merged
}

// EffectivePolicy merges root policy with overrides along ancestors root → node (in order).
func EffectivePolicy(root AttendancePolicy, overrides []PolicyOverride) AttendancePolicy {
	policy := root
	for _, override := range overrides {
		policy = MergePolicy(policy, override)
	}
	return policy
}

func allPunchTypes() []PunchType {
	return []PunchType{
		PunchTypeClockIn,
		PunchTypeClockOut,
		PunchTypeBreakStart,
		PunchTypeBreakEnd,
	}
}
