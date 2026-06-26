package workforce

import (
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

// PlacementService manages employee placements against a tenant org tree.
type PlacementService struct {
	tree       *organization.OrgTree
	placements []EmployeePlacement
}

// NewPlacementService creates a service bound to a validated org tree.
func NewPlacementService(tree *organization.OrgTree) *PlacementService {
	return &PlacementService{tree: tree}
}

// Assign registers a new placement after validating tenant and overlap rules.
func (s *PlacementService) Assign(placement EmployeePlacement) error {
	if err := s.validatePlacement(placement); err != nil {
		return err
	}
	if err := s.checkOverlapRules(placement); err != nil {
		return err
	}
	s.placements = append(s.placements, placement)
	return nil
}

// EndPlacement closes an open placement at the given instant.
func (s *PlacementService) EndPlacement(placementID string, at time.Time) error {
	for i := range s.placements {
		if s.placements[i].ID != placementID {
			continue
		}
		if s.placements[i].ValidUntil != nil {
			return ErrPlacementNotFound
		}
		s.placements[i].ValidUntil = &at
		return nil
	}
	return ErrPlacementNotFound
}

// Transfer ends the active primary placement and opens a new primary at newOrgNodeID.
func (s *PlacementService) Transfer(employeeID, newOrgNodeID, newPlacementID string, at time.Time) error {
	active, err := s.ActivePrimaryAt(employeeID, at)
	if err != nil {
		return err
	}
	if active == nil {
		return ErrPlacementNotFound
	}
	if err := s.EndPlacement(active.ID, at); err != nil {
		return err
	}
	return s.Assign(EmployeePlacement{
		ID:         newPlacementID,
		EmployeeID: employeeID,
		TenantID:   active.TenantID,
		OrgNodeID:  newOrgNodeID,
		Type:       PlacementTypePrimary,
		ValidFrom:  at,
	})
}

// ActivePlacementsAt returns all placements active for the employee at the instant.
func (s *PlacementService) ActivePlacementsAt(employeeID string, at time.Time) []EmployeePlacement {
	out := make([]EmployeePlacement, 0)
	for _, p := range s.placements {
		if p.EmployeeID != employeeID || !p.IsActiveAt(at) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// ActivePrimaryAt returns the active primary placement, if any.
func (s *PlacementService) ActivePrimaryAt(employeeID string, at time.Time) (*EmployeePlacement, error) {
	var found *EmployeePlacement
	for i := range s.placements {
		p := &s.placements[i]
		if p.EmployeeID != employeeID || p.Type != PlacementTypePrimary || !p.IsActiveAt(at) {
			continue
		}
		if found != nil {
			return nil, ErrDuplicatePrimary
		}
		found = p
	}
	return found, nil
}

// EffectiveOrgPath returns root → node path for geofence inheritance (BR-023).
func (s *PlacementService) EffectiveOrgPath(employeeID string, at time.Time) ([]string, error) {
	primary, err := s.ActivePrimaryAt(employeeID, at)
	if err != nil {
		return nil, err
	}
	if primary == nil {
		return nil, ErrPlacementNotFound
	}
	return s.tree.PathFromRoot(primary.OrgNodeID)
}

func (s *PlacementService) validatePlacement(p EmployeePlacement) error {
	if p.EmployeeID == "" {
		return ErrEmptyEmployeeID
	}
	if p.OrgNodeID == "" {
		return ErrEmptyOrgNodeID
	}
	if p.TenantID == "" {
		return ErrEmptyTenantID
	}
	if p.Type != PlacementTypePrimary && p.Type != PlacementTypeSecondary {
		return ErrInvalidPlacementType
	}
	if s.tree.TenantID != p.TenantID {
		return ErrTenantMismatch
	}
	if _, ok := s.tree.Node(p.OrgNodeID); !ok {
		return ErrNodeNotInTenant
	}
	return nil
}

func (s *PlacementService) checkOverlapRules(incoming EmployeePlacement) error {
	for _, existing := range s.placements {
		if existing.EmployeeID != incoming.EmployeeID {
			continue
		}
		if !periodsOverlap(existing.ValidFrom, existing.ValidUntil, incoming.ValidFrom, incoming.ValidUntil) {
			continue
		}
		if incoming.Type == PlacementTypePrimary && existing.Type == PlacementTypePrimary {
			return ErrDuplicatePrimary
		}
		if incoming.Type == PlacementTypeSecondary && existing.Type == PlacementTypeSecondary {
			return ErrOverlappingSecondary
		}
	}
	return nil
}
