package attendance

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/workforce"
)

// PunchDayReader loads VALID punches for a calendar day.
type PunchDayReader interface {
	PunchesForDay(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
		day time.Time,
	) ([]domainpunch.PunchRecord, error)
}

// ScheduleResolver returns the employee work schedule template.
type ScheduleResolver interface {
	WorkSchedule(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
	) (workforce.WorkSchedule, error)
}

// PolicyResolver returns attendance policy for time accounting.
type PolicyResolver interface {
	EffectivePolicy(
		ctx context.Context,
		tenantID, employeeID uuid.UUID,
	) (organization.AttendancePolicy, error)
}
