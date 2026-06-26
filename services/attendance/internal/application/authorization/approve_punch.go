package authorization

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

var (
	ErrEmployeeNotFound = errors.New("authorization: employee not found")
	ErrWriteDenied      = errors.New("authorization: write denied for role")
)

// AuthorizePunchApprovalCommand requests manager approval for a suspicious punch.
type AuthorizePunchApprovalCommand struct {
	Actor      organization.ActorScope
	TenantID   uuid.UUID
	EmployeeID uuid.UUID
}

// AuthorizePunchApprovalHandler wires employee lookup, placement, and ABAC rules.
type AuthorizePunchApprovalHandler struct {
	Employees  EmployeeReader
	Placements PlacementReader
	Auth       *PunchAuthorizationService
}

// Approve checks ABAC approval after resolving the employee placement from storage.
func (h AuthorizePunchApprovalHandler) Approve(
	ctx context.Context,
	cmd AuthorizePunchApprovalCommand,
) (bool, error) {
	emp, err := h.Employees.GetEmployee(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		return false, err
	}
	if emp == nil {
		return false, ErrEmployeeNotFound
	}

	placement, err := h.Placements.EmployeePlacement(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		return false, err
	}

	return h.Auth.ApprovePunch(cmd.Actor, placement)
}

// AuthorizeWrite rejects auditors and employees from mutating punches.
func (h AuthorizePunchApprovalHandler) AuthorizeWrite(actor organization.ActorScope) error {
	if !h.Auth.WritePunch(actor) {
		return ErrWriteDenied
	}
	return nil
}
