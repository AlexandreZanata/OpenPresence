package authorization

import (
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

// OrgTreeReader loads tenant org trees for authorization checks.
type OrgTreeReader interface {
	Tree(tenantID string) (*organization.OrgTree, error)
}

// PunchAuthorizationService gates punch operations with ABAC rules.
type PunchAuthorizationService struct {
	trees OrgTreeReader
}

// NewPunchAuthorizationService wires the org tree port.
func NewPunchAuthorizationService(trees OrgTreeReader) *PunchAuthorizationService {
	return &PunchAuthorizationService{trees: trees}
}

// ApprovePunch checks whether actor may approve a suspicious punch for the employee.
func (s *PunchAuthorizationService) ApprovePunch(actor organization.ActorScope, placement organization.EmployeePlacementRef) (bool, error) {
	if actor.TenantID != placement.TenantID {
		return false, organization.ErrCrossTenantAccess
	}
	tree, err := s.trees.Tree(placement.TenantID)
	if err != nil {
		return false, err
	}
	return organization.CanApprovePunch(actor, placement, tree)
}

// ExportPayroll checks tenant-scoped payroll export permission.
func (s *PunchAuthorizationService) ExportPayroll(actor organization.ActorScope, tenantID string) bool {
	return organization.CanExportPayroll(actor, tenantID)
}

// ReadPunch checks read access to an employee punch record.
func (s *PunchAuthorizationService) ReadPunch(actor organization.ActorScope, placement organization.EmployeePlacementRef) (bool, error) {
	if actor.TenantID != placement.TenantID {
		return false, organization.ErrCrossTenantAccess
	}
	tree, err := s.trees.Tree(placement.TenantID)
	if err != nil {
		return false, err
	}
	return organization.CanReadPunch(actor, placement, tree)
}

// WritePunch checks mutating punch operations (auditors denied).
func (s *PunchAuthorizationService) WritePunch(actor organization.ActorScope) bool {
	switch actor.Role {
	case organization.RoleAuditor, organization.RoleEmployee:
		return false
	case organization.RoleSuperAdmin, organization.RoleOrgAdmin, organization.RoleManager,
		organization.RoleHRAnalyst, organization.RoleSecurityAnalyst:
		return actor.TenantID != ""
	default:
		return false
	}
}
