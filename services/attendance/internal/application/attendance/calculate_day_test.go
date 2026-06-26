package attendance_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	appattendance "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/attendance"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
)

func TestCalculateDayAttendance_BR030_WorkedMinutes(t *testing.T) {
	day := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	handler := newTestHandler(t, workforce.Standard8h(), organization.DefaultPolicy(), []domainpunch.PunchRecord{
		validPunch(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour)),
		validPunch(domainpunch.PunchTypeBreakStart, day.Add(12*time.Hour)),
		validPunch(domainpunch.PunchTypeBreakEnd, day.Add(13*time.Hour)),
		validPunch(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour)),
	})

	result, err := handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: uuid.New(), EmployeeID: uuid.New(), Day: day,
	})
	require.NoError(t, err)
	require.Equal(t, 8*60, result.WorkedMinutes)
}

func TestCalculateDayAttendance_BR031_Lateness(t *testing.T) {
	day := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	policy := organization.DefaultPolicy()
	policy.ToleranceMinutes = 5
	handler := newTestHandler(t, workforce.Standard8h(), policy, []domainpunch.PunchRecord{
		validPunch(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour+10*time.Minute)),
		validPunch(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour)),
	})

	result, err := handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: uuid.New(), EmployeeID: uuid.New(), Day: day,
	})
	require.NoError(t, err)
	require.Equal(t, 5, result.LatenessMinutes)
}

func TestCalculateDayAttendance_BR032_Overtime(t *testing.T) {
	day := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	policy := organization.DefaultPolicy()
	policy.ToleranceMinutes = 5
	handler := newTestHandler(t, workforce.Standard8h(), policy, []domainpunch.PunchRecord{
		validPunch(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour)),
		validPunch(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour+20*time.Minute)),
	})

	result, err := handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: uuid.New(), EmployeeID: uuid.New(), Day: day,
	})
	require.NoError(t, err)
	require.Equal(t, 15, result.OvertimeMinutes)
}

func TestCalculateDayAttendance_BR033_Shift12x36Windows(t *testing.T) {
	day := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	handler := newTestHandler(t, workforce.Shift12x36(), organization.DefaultPolicy(), nil)

	result, err := handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: uuid.New(), EmployeeID: uuid.New(), Day: day,
	})
	require.NoError(t, err)
	require.Len(t, result.Windows, 1)
	require.Equal(t, 19, result.Windows[0].Start.Hour())
	require.Equal(t, 7, result.Windows[0].End.Hour())
}

func TestCalculateDayAttendance_BR034_TimeBank(t *testing.T) {
	day := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	policy := organization.DefaultPolicy()
	policy.OvertimePolicy = organization.OvertimePolicyTimeBank
	policy.ToleranceMinutes = 5
	handler := newTestHandler(t, workforce.Standard8h(), policy, []domainpunch.PunchRecord{
		validPunch(domainpunch.PunchTypeClockIn, day.Add(9*time.Hour)),
		validPunch(domainpunch.PunchTypeClockOut, day.Add(18*time.Hour+20*time.Minute)),
	})

	result, err := handler.Handle(context.Background(), appattendance.CalculateDayAttendanceCommand{
		TenantID: uuid.New(), EmployeeID: uuid.New(), Day: day,
		PriorTimeBankMinutes: 120,
	})
	require.NoError(t, err)
	require.Equal(t, 15, result.OvertimeMinutes)
	require.Equal(t, 135, result.TimeBankBalance)
}

func validPunch(t domainpunch.PunchType, at time.Time) domainpunch.PunchRecord {
	return domainpunch.PunchRecord{Type: t, PunchedAt: at, Status: domainpunch.PunchStatusValid}
}

func newTestHandler(
	t *testing.T,
	schedule workforce.WorkSchedule,
	policy organization.AttendancePolicy,
	punches []domainpunch.PunchRecord,
) appattendance.CalculateDayAttendanceHandler {
	t.Helper()
	return appattendance.CalculateDayAttendanceHandler{
		Punches:   &stubPunchReader{punches: punches},
		Schedules: &stubScheduleResolver{schedule: schedule},
		Policies:  &stubPolicyResolver{policy: policy},
	}
}

type stubPunchReader struct {
	punches []domainpunch.PunchRecord
}

func (s *stubPunchReader) PunchesForDay(
	_ context.Context, _, _ uuid.UUID, _ time.Time,
) ([]domainpunch.PunchRecord, error) {
	return s.punches, nil
}

type stubScheduleResolver struct {
	schedule workforce.WorkSchedule
}

func (s *stubScheduleResolver) WorkSchedule(
	_ context.Context, _, _ uuid.UUID,
) (workforce.WorkSchedule, error) {
	return s.schedule, nil
}

type stubPolicyResolver struct {
	policy organization.AttendancePolicy
}

func (s *stubPolicyResolver) EffectivePolicy(
	_ context.Context, _, _ uuid.UUID,
) (organization.AttendancePolicy, error) {
	return s.policy, nil
}
