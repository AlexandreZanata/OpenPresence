package fraud

import "time"

const (
	lockoutRejectThreshold = 3
	lockoutWindowMinutes   = 10
	lockoutDurationMinutes = 30
)

// DeviceLockoutTracker enforces BR-013 consecutive rejection lockout.
type DeviceLockoutTracker struct {
	rejects map[string][]time.Time
	locked  map[string]time.Time
}

// NewDeviceLockoutTracker creates an empty tracker.
func NewDeviceLockoutTracker() *DeviceLockoutTracker {
	return &DeviceLockoutTracker{
		rejects: make(map[string][]time.Time),
		locked:  make(map[string]time.Time),
	}
}

// RecordRejected registers a rejected attempt and applies lockout when threshold is met.
func (t *DeviceLockoutTracker) RecordRejected(deviceID string, at time.Time) {
	if t.IsLocked(deviceID, at) {
		return
	}
	windowStart := at.Add(-time.Duration(lockoutWindowMinutes) * time.Minute)
	recent := filterSince(t.rejects[deviceID], windowStart)
	recent = append(recent, at)
	t.rejects[deviceID] = recent

	if len(recent) >= lockoutRejectThreshold {
		t.locked[deviceID] = at.Add(time.Duration(lockoutDurationMinutes) * time.Minute)
	}
}

// IsLocked reports whether the device is blocked at the given instant.
func (t *DeviceLockoutTracker) IsLocked(deviceID string, at time.Time) bool {
	until, ok := t.locked[deviceID]
	if !ok {
		return false
	}
	return at.Before(until)
}

// LockedUntil returns lock expiry when device is locked.
func (t *DeviceLockoutTracker) LockedUntil(deviceID string) (time.Time, bool) {
	until, ok := t.locked[deviceID]
	return until, ok
}

func filterSince(times []time.Time, since time.Time) []time.Time {
	out := make([]time.Time, 0, len(times))
	for _, ts := range times {
		if !ts.Before(since) {
			out = append(out, ts)
		}
	}
	return out
}
