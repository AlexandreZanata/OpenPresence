package organization

import (
	"errors"
	"testing"
)

const testTenant = "tenant-001"

// Municipality fixture (docs/ORGANIZATION.md):
// Health Secretariat → Municipal Hospital → Nursing department.
func TestBuildTree_MunicipalityHealthSecretariat(t *testing.T) {
	nodes := []OrgNode{
		{ID: "health-sec", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "Health Secretariat"},
		{ID: "hospital", TenantID: testTenant, ParentID: "health-sec", Type: OrgNodeTypeLocation, Name: "Municipal Hospital"},
		{ID: "nursing", TenantID: testTenant, ParentID: "hospital", Type: OrgNodeTypeDepartment, Name: "Nursing"},
	}

	tree, err := BuildTree(testTenant, nodes)
	if err != nil {
		t.Fatalf("expected valid municipality tree, got %v", err)
	}
	if len(tree.Nodes()) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(tree.Nodes()))
	}
	nursing, ok := tree.Node("nursing")
	if !ok || nursing.Type != OrgNodeTypeDepartment {
		t.Fatal("nursing department must exist under hospital LOCATION")
	}
}

// Private company fixture (docs/ORGANIZATION.md):
// HQ São Paulo → Sales → Inside Sales team.
func TestBuildTree_PrivateCompanyHQSales(t *testing.T) {
	nodes := []OrgNode{
		{ID: "hq-sp", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "HQ São Paulo"},
		{ID: "sales", TenantID: testTenant, ParentID: "hq-sp", Type: OrgNodeTypeDepartment, Name: "Sales"},
		{ID: "inside-sales", TenantID: testTenant, ParentID: "sales", Type: OrgNodeTypeTeam, Name: "Inside Sales"},
	}

	tree, err := BuildTree(testTenant, nodes)
	if err != nil {
		t.Fatalf("expected valid private company tree, got %v", err)
	}
	team, ok := tree.Node("inside-sales")
	if !ok || team.ParentID != "sales" {
		t.Fatal("Inside Sales TEAM must hang under Sales DEPARTMENT")
	}
}

func TestBuildTree_WorkSiteUnderDivisionAndDepartment(t *testing.T) {
	nodes := []OrgNode{
		{ID: "projects", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "Projects"},
		{ID: "site-a", TenantID: testTenant, ParentID: "projects", Type: OrgNodeTypeWorkSite, Name: "Highway BR-163 Site"},
		{ID: "ops", TenantID: testTenant, ParentID: "projects", Type: OrgNodeTypeDepartment, Name: "Operations"},
		{ID: "site-b", TenantID: testTenant, ParentID: "ops", Type: OrgNodeTypeWorkSite, Name: "Temporary Yard"},
	}

	if _, err := BuildTree(testTenant, nodes); err != nil {
		t.Fatalf("WORK_SITE must attach under DIVISION or DEPARTMENT: %v", err)
	}
}

func TestBuildTree_MultipleTopLevelDivisions(t *testing.T) {
	nodes := []OrgNode{
		{ID: "health", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "Health Secretariat"},
		{ID: "education", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "Education Secretariat"},
	}

	tree, err := BuildTree(testTenant, nodes)
	if err != nil {
		t.Fatalf("municipality may have multiple secretariat DIVISION roots: %v", err)
	}
	if len(tree.Nodes()) != 2 {
		t.Fatalf("expected 2 top-level divisions, got %d", len(tree.Nodes()))
	}
}

func TestBuildTree_RejectsCycle(t *testing.T) {
	nodes := []OrgNode{
		{ID: "a", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "A"},
		{ID: "b", TenantID: testTenant, ParentID: "a", Type: OrgNodeTypeDepartment, Name: "B"},
	}
	nodes[0].ParentID = "b"

	_, err := BuildTree(testTenant, nodes)
	if !errors.Is(err, ErrCycle) {
		t.Fatalf("expected ErrCycle for A→B→A, got %v", err)
	}
}

func TestBuildTree_RejectsOrphan(t *testing.T) {
	nodes := []OrgNode{
		{ID: "root", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "Root"},
		{ID: "lost", TenantID: testTenant, ParentID: "missing", Type: OrgNodeTypeDepartment, Name: "Orphan"},
	}

	_, err := BuildTree(testTenant, nodes)
	if !errors.Is(err, ErrOrphanNode) {
		t.Fatalf("expected ErrOrphanNode, got %v", err)
	}
}

func TestBuildTree_RejectsTeamUnderDivision(t *testing.T) {
	nodes := []OrgNode{
		{ID: "hq", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "HQ"},
		{ID: "team", TenantID: testTenant, ParentID: "hq", Type: OrgNodeTypeTeam, Name: "Direct Team"},
	}

	_, err := BuildTree(testTenant, nodes)
	if !errors.Is(err, ErrInvalidChildType) {
		t.Fatalf("TEAM must not hang directly under DIVISION: %v", err)
	}
}

func TestBuildTree_RejectsTeamAsRoot(t *testing.T) {
	nodes := []OrgNode{
		{ID: "team-root", TenantID: testTenant, Type: OrgNodeTypeTeam, Name: "Invalid Root"},
	}

	_, err := BuildTree(testTenant, nodes)
	if !errors.Is(err, ErrInvalidRootType) {
		t.Fatalf("top-level node must be DIVISION: %v", err)
	}
}

func TestBuildTree_RejectsWorkSiteUnderTeam(t *testing.T) {
	nodes := []OrgNode{
		{ID: "div", TenantID: testTenant, Type: OrgNodeTypeDivision, Name: "Div"},
		{ID: "dept", TenantID: testTenant, ParentID: "div", Type: OrgNodeTypeDepartment, Name: "Dept"},
		{ID: "team", TenantID: testTenant, ParentID: "dept", Type: OrgNodeTypeTeam, Name: "Team"},
		{ID: "site", TenantID: testTenant, ParentID: "team", Type: OrgNodeTypeWorkSite, Name: "Site"},
	}

	_, err := BuildTree(testTenant, nodes)
	if !errors.Is(err, ErrInvalidChildType) {
		t.Fatalf("WORK_SITE must not attach under TEAM: %v", err)
	}
}

func TestCanBeChild_PublicSectorMapping(t *testing.T) {
	if !CanBeChild("", OrgNodeTypeDivision, true) {
		t.Fatal("secretariat maps to top-level DIVISION")
	}
	if !CanBeChild(OrgNodeTypeDivision, OrgNodeTypeLocation, false) {
		t.Fatal("hospital/UBS maps to LOCATION under secretariat DIVISION")
	}
	if !CanBeChild(OrgNodeTypeLocation, OrgNodeTypeDepartment, false) {
		t.Fatal("department maps under LOCATION (e.g. Nursing)")
	}
}
