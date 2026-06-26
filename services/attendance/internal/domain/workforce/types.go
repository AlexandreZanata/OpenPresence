package workforce

import (
	"errors"
	"time"
)

// PlacementType classifies employee assignment to an org node.
type PlacementType string

const (
	PlacementTypePrimary   PlacementType = "PRIMARY"
	PlacementTypeSecondary PlacementType = "SECONDARY"
)

// EmployeePlacement assigns an employee to an org node for a time window.
type EmployeePlacement struct {
	ID         string
	EmployeeID string
	TenantID   string
	OrgNodeID  string
	Type       PlacementType
	ValidFrom  time.Time
	ValidUntil *time.Time
}

var (
	ErrEmptyEmployeeID      = errors.New("workforce: employee id is required")
	ErrEmptyOrgNodeID       = errors.New("workforce: org node id is required")
	ErrEmptyTenantID        = errors.New("workforce: tenant id is required")
	ErrNodeNotInTenant      = errors.New("workforce: org node not found in tenant tree")
	ErrTenantMismatch       = errors.New("workforce: placement tenant does not match org tree")
	ErrDuplicatePrimary     = errors.New("workforce: active primary placement already exists")
	ErrOverlappingSecondary = errors.New("workforce: overlapping secondary placement")
	ErrPlacementNotFound    = errors.New("workforce: placement not found")
	ErrInvalidPlacementType = errors.New("workforce: invalid placement type")
)

// IsActiveAt reports whether the placement covers the given instant.
func (p EmployeePlacement) IsActiveAt(at time.Time) bool {
	if at.Before(p.ValidFrom) {
		return false
	}
	if p.ValidUntil != nil && !at.Before(*p.ValidUntil) {
		return false
	}
	return true
}

func periodsOverlap(aFrom time.Time, aUntil *time.Time, bFrom time.Time, bUntil *time.Time) bool {
	aEnd := farFuture
	if aUntil != nil {
		aEnd = *aUntil
	}
	bEnd := farFuture
	if bUntil != nil {
		bEnd = *bUntil
	}
	return aFrom.Before(bEnd) && bFrom.Before(aEnd)
}

var farFuture = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
