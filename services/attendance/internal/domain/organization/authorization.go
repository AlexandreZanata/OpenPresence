package organization

import "errors"

// Role is a named permission bundle (see docs/GLOSSARY.md).
type Role string

const (
	RoleSuperAdmin      Role = "SUPER_ADMIN"
	RoleOrgAdmin        Role = "ORG_ADMIN"
	RoleManager         Role = "MANAGER"
	RoleHRAnalyst       Role = "HR_ANALYST"
	RoleSecurityAnalyst Role = "SECURITY_ANALYST"
	RoleEmployee        Role = "EMPLOYEE"
	RoleAuditor         Role = "AUDITOR"
)

// ActorScope identifies who is performing an authorized action.
type ActorScope struct {
	Role              Role
	TenantID          string
	AssignedOrgNodeID string
	EmployeeID        string
}

// EmployeePlacementRef is the employee org assignment used for ABAC checks.
type EmployeePlacementRef struct {
	EmployeeID string
	TenantID   string
	OrgNodeID  string
}

var ErrCrossTenantAccess = errors.New("organization: cross-tenant access denied")

// IsDescendant reports whether nodeID is the same as or below ancestorID in the tree.
func (t *OrgTree) IsDescendant(ancestorID, nodeID string) (bool, error) {
	if ancestorID == nodeID {
		return true, nil
	}
	path, err := t.PathFromRoot(nodeID)
	if err != nil {
		return false, err
	}
	for _, id := range path {
		if id == ancestorID {
			return true, nil
		}
	}
	return false, nil
}

// CanAccessSubtree reports whether targetNodeID is within actorNodeID scope.
func CanAccessSubtree(actorNodeID, targetNodeID string, tree *OrgTree) (bool, error) {
	return tree.IsDescendant(actorNodeID, targetNodeID)
}

// CanApprovePunch enforces manager/org-admin subtree approval (docs/ORGANIZATION.md ABAC).
func CanApprovePunch(actor ActorScope, placement EmployeePlacementRef, tree *OrgTree) (bool, error) {
	if err := assertSameTenant(actor.TenantID, placement.TenantID, tree.TenantID); err != nil {
		return false, err
	}
	switch actor.Role {
	case RoleSuperAdmin, RoleHRAnalyst, RoleSecurityAnalyst:
		return true, nil
	case RoleOrgAdmin, RoleManager:
		return tree.IsDescendant(actor.AssignedOrgNodeID, placement.OrgNodeID)
	case RoleAuditor, RoleEmployee:
		return false, nil
	default:
		return false, nil
	}
}

// CanExportPayroll allows HR and auditors within tenant boundary only.
func CanExportPayroll(actor ActorScope, tenantID string) bool {
	if actor.TenantID != tenantID {
		return false
	}
	switch actor.Role {
	case RoleSuperAdmin, RoleHRAnalyst, RoleAuditor:
		return true
	default:
		return false
	}
}

// CanReadPunch allows tenant-scoped read for oversight roles and self-access.
func CanReadPunch(actor ActorScope, placement EmployeePlacementRef, tree *OrgTree) (bool, error) {
	if err := assertSameTenant(actor.TenantID, placement.TenantID, tree.TenantID); err != nil {
		return false, err
	}
	if actor.Role == RoleEmployee {
		return actor.EmployeeID == placement.EmployeeID, nil
	}
	switch actor.Role {
	case RoleSuperAdmin, RoleHRAnalyst, RoleSecurityAnalyst, RoleAuditor:
		return true, nil
	case RoleOrgAdmin, RoleManager:
		return tree.IsDescendant(actor.AssignedOrgNodeID, placement.OrgNodeID)
	default:
		return false, nil
	}
}

func assertSameTenant(actorTenant, placementTenant, treeTenant string) error {
	if actorTenant != placementTenant || actorTenant != treeTenant {
		return ErrCrossTenantAccess
	}
	return nil
}
