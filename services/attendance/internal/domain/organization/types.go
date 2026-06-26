package organization

import "errors"

// OrgNodeType classifies nodes in the tenant org tree (see docs/GLOSSARY.md).
type OrgNodeType string

const (
	OrgNodeTypeDivision   OrgNodeType = "DIVISION"
	OrgNodeTypeDepartment OrgNodeType = "DEPARTMENT"
	OrgNodeTypeSection    OrgNodeType = "SECTION"
	OrgNodeTypeTeam       OrgNodeType = "TEAM"
	OrgNodeTypeLocation   OrgNodeType = "LOCATION"
	OrgNodeTypeWorkSite   OrgNodeType = "WORK_SITE"
)

// OrgNode is one node in the tenant-scoped organization tree.
// Public sector: secretariat → DIVISION, hospital/UBS → LOCATION.
type OrgNode struct {
	ID       string
	TenantID string
	ParentID string
	Type     OrgNodeType
	Name     string
	Code     string
}

var (
	ErrEmptyTenant       = errors.New("organization: tenant id is required")
	ErrDuplicateNode     = errors.New("organization: duplicate node id")
	ErrOrphanNode        = errors.New("organization: parent not found in tree")
	ErrMultipleTenants   = errors.New("organization: all nodes must share tenant")
	ErrInvalidRootType   = errors.New("organization: top-level node must be DIVISION")
	ErrInvalidChildType  = errors.New("organization: invalid child type for parent")
	ErrCycle             = errors.New("organization: cycle detected in tree")
	ErrEmptyTree         = errors.New("organization: tree has no nodes")
	ErrInvalidNode       = errors.New("organization: node id and name are required")
	ErrNodeNotFound      = errors.New("organization: node not found in tree")
)
