package workforce

import "time"

// WorkWindow is one independent attendance window (split shift or 12×36 segment).
type WorkWindow struct {
	Start time.Duration
	End   time.Duration // when End <= Start, window crosses midnight into next day
}

// WorkSchedule defines planned work for an employee (see docs/GLOSSARY.md).
type WorkSchedule struct {
	ScheduledStart   time.Duration
	ScheduledEnd     time.Duration
	Windows          []WorkWindow
	ToleranceMinutes int
}

// ResolvedWindow is an absolute time span on a calendar day.
type ResolvedWindow struct {
	Start time.Time
	End   time.Time
}

// Standard8h returns a private-office schedule 09:00–18:00 with 10m tolerance.
func Standard8h() WorkSchedule {
	return WorkSchedule{
		ScheduledStart:   9 * time.Hour,
		ScheduledEnd:     18 * time.Hour,
		ToleranceMinutes: 10,
	}
}

// Shift12x36 returns a public-health night shift 19:00–07:00 (crosses midnight).
func Shift12x36() WorkSchedule {
	return WorkSchedule{
		ScheduledStart: 19 * time.Hour,
		ScheduledEnd:   7 * time.Hour,
		Windows: []WorkWindow{
			{Start: 19 * time.Hour, End: 7 * time.Hour},
		},
		ToleranceMinutes: 30,
	}
}

// SplitShift returns morning and afternoon windows with a midday gap.
func SplitShift() WorkSchedule {
	return WorkSchedule{
		ScheduledStart: 8 * time.Hour,
		ScheduledEnd:   18 * time.Hour,
		Windows: []WorkWindow{
			{Start: 8 * time.Hour, End: 12 * time.Hour},
			{Start: 14 * time.Hour, End: 18 * time.Hour},
		},
		ToleranceMinutes: 10,
	}
}

func scheduleWindows(schedule WorkSchedule) []WorkWindow {
	if len(schedule.Windows) > 0 {
		return schedule.Windows
	}
	return []WorkWindow{{Start: schedule.ScheduledStart, End: schedule.ScheduledEnd}}
}

func dayStart(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
}

func resolveWindow(day time.Time, window WorkWindow) ResolvedWindow {
	start := dayStart(day).Add(window.Start)
	end := dayStart(day).Add(window.End)
	if window.End <= window.Start {
		end = end.Add(24 * time.Hour)
	}
	return ResolvedWindow{Start: start, End: end}
}
