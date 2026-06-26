package attendance

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
)

// CalculateDayAttendanceCommand requests time accounting for one calendar day.
type CalculateDayAttendanceCommand struct {
	TenantID             uuid.UUID
	EmployeeID           uuid.UUID
	Day                  time.Time
	PriorTimeBankMinutes int
}

// DayAttendanceResult is BR-030–034 accounting output for a day.
type DayAttendanceResult struct {
	WorkedMinutes   int
	LatenessMinutes int
	OvertimeMinutes int
	Windows         []workforce.ResolvedWindow
	TimeBankBalance int
}

// CalculateDayAttendanceHandler computes worked time, lateness, overtime, and time bank.
type CalculateDayAttendanceHandler struct {
	Punches   PunchDayReader
	Schedules ScheduleResolver
	Policies  PolicyResolver
}

// Handle loads punches from storage and applies workforce domain rules.
func (h CalculateDayAttendanceHandler) Handle(
	ctx context.Context,
	cmd CalculateDayAttendanceCommand,
) (*DayAttendanceResult, error) {
	punches, err := h.Punches.PunchesForDay(ctx, cmd.TenantID, cmd.EmployeeID, cmd.Day)
	if err != nil {
		return nil, fmt.Errorf("load day punches: %w", err)
	}

	schedule, err := h.Schedules.WorkSchedule(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("load schedule: %w", err)
	}

	policy, err := h.Policies.EffectivePolicy(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("load policy: %w", err)
	}

	events := toPunchEvents(punches)
	worked := workforce.CalculateWorkedMinutes(events)

	var lateness, overtime int
	if clockIn := firstClockIn(events); clockIn != nil {
		lateness = workforce.CalculateLateness(*clockIn, schedule, policy)
	}
	if clockOut := lastClockOut(events); clockOut != nil {
		overtime = workforce.CalculateOvertime(*clockOut, schedule, policy)
	}

	day := dayStart(cmd.Day)
	windows := workforce.EvaluateWindows(day, schedule)
	timeBank := workforce.UpdateTimeBank(cmd.PriorTimeBankMinutes, overtime, policy.OvertimePolicy)

	return &DayAttendanceResult{
		WorkedMinutes:   worked,
		LatenessMinutes: lateness,
		OvertimeMinutes: overtime,
		Windows:         windows,
		TimeBankBalance: timeBank,
	}, nil
}

func toPunchEvents(punches []domainpunch.PunchRecord) []workforce.PunchEvent {
	valid := make([]domainpunch.PunchRecord, 0, len(punches))
	for _, p := range punches {
		if p.Status == domainpunch.PunchStatusValid {
			valid = append(valid, p)
		}
	}
	sort.Slice(valid, func(i, j int) bool {
		return valid[i].PunchedAt.Before(valid[j].PunchedAt)
	})
	out := make([]workforce.PunchEvent, len(valid))
	for i, p := range valid {
		out[i] = workforce.PunchEvent{
			Type: organization.PunchType(p.Type),
			At:   p.PunchedAt,
		}
	}
	return out
}

func firstClockIn(events []workforce.PunchEvent) *time.Time {
	for _, e := range events {
		if e.Type == organization.PunchTypeClockIn {
			t := e.At
			return &t
		}
	}
	return nil
}

func lastClockOut(events []workforce.PunchEvent) *time.Time {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Type == organization.PunchTypeClockOut {
			t := events[i].At
			return &t
		}
	}
	return nil
}

func dayStart(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
}
