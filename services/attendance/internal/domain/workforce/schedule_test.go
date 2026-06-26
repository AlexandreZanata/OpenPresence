package workforce

import (
	"testing"
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

var loc = time.UTC

func dayAt(year int, month time.Month, day, hour, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}

func TestCalculateWorkedMinutes_BR030_PrivateOfficeWithBreak(t *testing.T) {
	d := dayAt(2026, 6, 26, 0, 0)
	punches := []PunchEvent{
		{Type: organization.PunchTypeClockIn, At: d.Add(9 * time.Hour)},
		{Type: organization.PunchTypeBreakStart, At: d.Add(12 * time.Hour)},
		{Type: organization.PunchTypeBreakEnd, At: d.Add(13 * time.Hour)},
		{Type: organization.PunchTypeClockOut, At: d.Add(18 * time.Hour)},
	}

	got := CalculateWorkedMinutes(punches)
	want := 8 * 60
	if got != want {
		t.Fatalf("BR-030: worked minutes %d, want %d (9h gross minus 1h break)", got, want)
	}
}

func TestCalculateLateness_BR031_TenMinutesInFiveTolerance(t *testing.T) {
	schedule := Standard8h()
	policy := organization.DefaultPolicy()
	policy.ToleranceMinutes = 5

	clockIn := dayAt(2026, 6, 26, 9, 10)
	got := CalculateLateness(clockIn, schedule, policy)
	want := 5
	if got != want {
		t.Fatalf("BR-031: lateness %d min, want %d", got, want)
	}
}

func TestCalculateOvertime_BR032_AfterEndPlusTolerance(t *testing.T) {
	schedule := Standard8h()
	policy := organization.DefaultPolicy()
	policy.ToleranceMinutes = 5

	clockOut := dayAt(2026, 6, 26, 18, 20)
	got := CalculateOvertime(clockOut, schedule, policy)
	want := 15
	if got != want {
		t.Fatalf("BR-032: overtime %d min, want %d", got, want)
	}
}

func TestCalculateOvertime_BR032_DisabledWhenPolicyOff(t *testing.T) {
	schedule := Standard8h()
	policy := organization.DefaultPolicy()
	policy.OvertimePolicy = organization.OvertimePolicyDisabled

	clockOut := dayAt(2026, 6, 26, 19, 0)
	if got := CalculateOvertime(clockOut, schedule, policy); got != 0 {
		t.Fatalf("BR-032: disabled overtime must be 0, got %d", got)
	}
}

func TestEvaluateWindows_BR033_Shift12x36CrossesMidnight(t *testing.T) {
	schedule := Shift12x36()
	day := dayAt(2026, 6, 26, 0, 0)
	windows := EvaluateWindows(day, schedule)

	if len(windows) != 1 {
		t.Fatalf("expected one 12×36 window, got %d", len(windows))
	}
	w := windows[0]
	if w.Start.Hour() != 19 || w.End.Hour() != 7 {
		t.Fatalf("12×36 window start/end hours got %v–%v", w.Start, w.End)
	}
	if !w.End.After(w.Start) {
		t.Fatal("resolved window end must be after start across midnight")
	}
	if w.End.Day() != day.Day()+1 {
		t.Fatal("night shift end must fall on next calendar day")
	}
}

func TestEvaluateWindows_BR033_SplitShiftIndependentWindows(t *testing.T) {
	schedule := SplitShift()
	day := dayAt(2026, 6, 26, 0, 0)
	windows := EvaluateWindows(day, schedule)

	if len(windows) != 2 {
		t.Fatalf("split shift expects 2 windows, got %d", len(windows))
	}
	if windows[0].End.Sub(windows[0].Start) != 4*time.Hour {
		t.Fatal("morning window must be 4 hours")
	}
	if windows[1].Start.Hour() != 14 {
		t.Fatal("afternoon window must start at 14:00")
	}
}

func TestUpdateTimeBank_BR034_AccumulatesWithTimeBankPolicy(t *testing.T) {
	got := UpdateTimeBank(120, 45, organization.OvertimePolicyTimeBank)
	if got != 165 {
		t.Fatalf("BR-034: balance %d, want 165", got)
	}
}

func TestUpdateTimeBank_BR034_StandardPolicyNoAccrual(t *testing.T) {
	got := UpdateTimeBank(120, 45, organization.OvertimePolicyStandard)
	if got != 120 {
		t.Fatalf("standard policy must not update bank, got %d", got)
	}
}

func TestCalculateWorkedMinutes_BR030_NursingNightShift(t *testing.T) {
	day := dayAt(2026, 6, 26, 0, 0)
	punches := []PunchEvent{
		{Type: organization.PunchTypeClockIn, At: day.Add(19 * time.Hour)},
		{Type: organization.PunchTypeBreakStart, At: day.Add(23 * time.Hour)},
		{Type: organization.PunchTypeBreakEnd, At: day.Add(23*time.Hour + 30*time.Minute)},
		{Type: organization.PunchTypeClockOut, At: day.Add(31 * time.Hour)},
	}

	got := CalculateWorkedMinutes(punches)
	want := 11*60 + 30
	if got != want {
		t.Fatalf("nursing 12×36 worked %d min, want %d", got, want)
	}
}

func TestWorkSchedule_AttachedToEmployeeByIDInFixtures(t *testing.T) {
	employeeSchedule := map[string]WorkSchedule{
		"emp-nursing": Shift12x36(),
		"emp-office":  Standard8h(),
	}
	if employeeSchedule["emp-nursing"].ScheduledStart != 19*time.Hour {
		t.Fatal("nursing employee must use 12×36 template")
	}
	if employeeSchedule["emp-office"].ScheduledEnd != 18*time.Hour {
		t.Fatal("office employee must use 8h template")
	}
}
