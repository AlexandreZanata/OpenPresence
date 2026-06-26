package workforce

import (
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

// PunchEvent is a minimal punch input for time accounting (BR-030–034).
type PunchEvent struct {
	Type organization.PunchType
	At   time.Time
}

// CalculateWorkedMinutes implements BR-030: in/out pairs minus breaks.
func CalculateWorkedMinutes(punches []PunchEvent) int {
	var worked int
	var breaks int
	var clockIn *time.Time
	var breakStart *time.Time

	for _, punch := range punches {
		switch punch.Type {
		case organization.PunchTypeClockIn:
			t := punch.At
			clockIn = &t
		case organization.PunchTypeClockOut:
			if clockIn != nil {
				worked += minutesBetween(*clockIn, punch.At)
				clockIn = nil
			}
		case organization.PunchTypeBreakStart:
			t := punch.At
			breakStart = &t
		case organization.PunchTypeBreakEnd:
			if breakStart != nil {
				breaks += minutesBetween(*breakStart, punch.At)
				breakStart = nil
			}
		}
	}
	return worked - breaks
}

// CalculateLateness implements BR-031 using schedule window start + policy tolerance.
func CalculateLateness(clockIn time.Time, schedule WorkSchedule, policy organization.AttendancePolicy) int {
	window := nearestWindowStart(clockIn, schedule)
	if window == nil {
		return 0
	}
	allowed := window.Start.Add(time.Duration(policy.ToleranceMinutes) * time.Minute)
	if !clockIn.After(allowed) {
		return 0
	}
	return minutesBetween(allowed, clockIn)
}

// CalculateOvertime implements BR-032 when policy allows overtime.
func CalculateOvertime(clockOut time.Time, schedule WorkSchedule, policy organization.AttendancePolicy) int {
	if policy.OvertimePolicy == organization.OvertimePolicyDisabled {
		return 0
	}
	window := nearestWindowEnd(clockOut, schedule)
	if window == nil {
		return 0
	}
	allowed := window.End.Add(time.Duration(policy.ToleranceMinutes) * time.Minute)
	if !clockOut.After(allowed) {
		return 0
	}
	return minutesBetween(allowed, clockOut)
}

// EvaluateWindows implements BR-033 — each configured window on the given day.
func EvaluateWindows(day time.Time, schedule WorkSchedule) []ResolvedWindow {
	windows := scheduleWindows(schedule)
	out := make([]ResolvedWindow, 0, len(windows))
	for _, w := range windows {
		out = append(out, resolveWindow(day, w))
	}
	return out
}

// UpdateTimeBank implements BR-034 cumulative balance when policy uses time bank.
func UpdateTimeBank(balanceMinutes, overtimeMinutes int, policy organization.OvertimePolicy) int {
	if policy != organization.OvertimePolicyTimeBank {
		return balanceMinutes
	}
	return balanceMinutes + overtimeMinutes
}

func nearestWindowStart(at time.Time, schedule WorkSchedule) *ResolvedWindow {
	day := dayStart(at)
	candidates := EvaluateWindows(day, schedule)
	if at.Before(day) {
		candidates = append(EvaluateWindows(day.Add(-24*time.Hour), schedule), candidates...)
	}
	var best *ResolvedWindow
	for i := range candidates {
		w := &candidates[i]
		if at.Before(w.Start) {
			continue
		}
		if best == nil || w.Start.After(best.Start) {
			best = w
		}
	}
	return best
}

func nearestWindowEnd(at time.Time, schedule WorkSchedule) *ResolvedWindow {
	day := dayStart(at)
	candidates := EvaluateWindows(day, schedule)
	if at.Hour() < 12 {
		prev := EvaluateWindows(day.Add(-24*time.Hour), schedule)
		candidates = append(prev, candidates...)
	}
	var best *ResolvedWindow
	for i := range candidates {
		w := &candidates[i]
		if at.Before(w.End) {
			continue
		}
		if best == nil || w.End.Before(best.End) {
			best = w
		}
	}
	if best == nil && len(candidates) > 0 {
		return &candidates[len(candidates)-1]
	}
	return best
}

func minutesBetween(from, to time.Time) int {
	return int(to.Sub(from).Minutes())
}
