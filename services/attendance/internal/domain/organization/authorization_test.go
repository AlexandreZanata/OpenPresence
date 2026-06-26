package organization

import "testing"

const authTenant = "city-example"

func municipalTree(t *testing.T) *OrgTree {
	t.Helper()
	nodes := []OrgNode{
		{ID: "health", TenantID: authTenant, Type: OrgNodeTypeDivision, Name: "Health Secretariat"},
		{ID: "hospital", TenantID: authTenant, ParentID: "health", Type: OrgNodeTypeLocation, Name: "Municipal Hospital"},
		{ID: "nursing", TenantID: authTenant, ParentID: "hospital", Type: OrgNodeTypeDepartment, Name: "Nursing"},
		{ID: "education", TenantID: authTenant, Type: OrgNodeTypeDivision, Name: "Education Secretariat"},
		{ID: "school-a", TenantID: authTenant, ParentID: "education", Type: OrgNodeTypeLocation, Name: "School A"},
	}
	tree, err := BuildTree(authTenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}
	return tree
}

func privateTree(t *testing.T) *OrgTree {
	t.Helper()
	nodes := []OrgNode{
		{ID: "hq", TenantID: authTenant, Type: OrgNodeTypeDivision, Name: "HQ"},
		{ID: "sales", TenantID: authTenant, ParentID: "hq", Type: OrgNodeTypeDepartment, Name: "Sales"},
		{ID: "inside", TenantID: authTenant, ParentID: "sales", Type: OrgNodeTypeTeam, Name: "Inside Sales"},
		{ID: "it", TenantID: authTenant, ParentID: "hq", Type: OrgNodeTypeDepartment, Name: "IT"},
	}
	tree, err := BuildTree(authTenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}
	return tree
}

func TestIsDescendant_NursingUnderHealth(t *testing.T) {
	tree := municipalTree(t)
	ok, err := tree.IsDescendant("health", "nursing")
	if err != nil || !ok {
		t.Fatal("nursing must be descendant of health secretariat")
	}
}

func TestCanApprovePunch_PublicHealthManagerApprovesNurse(t *testing.T) {
	tree := municipalTree(t)
	actor := ActorScope{Role: RoleManager, TenantID: authTenant, AssignedOrgNodeID: "health"}
	placement := EmployeePlacementRef{EmployeeID: "nurse-1", TenantID: authTenant, OrgNodeID: "nursing"}

	ok, err := CanApprovePunch(actor, placement, tree)
	if err != nil || !ok {
		t.Fatal("health manager must approve hospital nurse punch")
	}
}

func TestCanApprovePunch_PublicHealthManagerRejectsEducation(t *testing.T) {
	tree := municipalTree(t)
	actor := ActorScope{Role: RoleManager, TenantID: authTenant, AssignedOrgNodeID: "health"}
	placement := EmployeePlacementRef{EmployeeID: "teacher-1", TenantID: authTenant, OrgNodeID: "school-a"}

	ok, err := CanApprovePunch(actor, placement, tree)
	if err != nil || ok {
		t.Fatal("health manager must not approve education employee")
	}
}

func TestCanApprovePunch_PrivateSalesManagerApprovesInsideSales(t *testing.T) {
	tree := privateTree(t)
	actor := ActorScope{Role: RoleManager, TenantID: authTenant, AssignedOrgNodeID: "sales"}
	placement := EmployeePlacementRef{EmployeeID: "rep-1", TenantID: authTenant, OrgNodeID: "inside"}

	ok, err := CanApprovePunch(actor, placement, tree)
	if err != nil || !ok {
		t.Fatal("sales manager must approve inside sales punch")
	}
}

func TestCanApprovePunch_PrivateSalesManagerRejectsIT(t *testing.T) {
	tree := privateTree(t)
	actor := ActorScope{Role: RoleManager, TenantID: authTenant, AssignedOrgNodeID: "sales"}
	placement := EmployeePlacementRef{EmployeeID: "dev-1", TenantID: authTenant, OrgNodeID: "it"}

	ok, err := CanApprovePunch(actor, placement, tree)
	if err != nil || ok {
		t.Fatal("sales manager must not approve IT employee")
	}
}

func TestCanApprovePunch_OrgAdminSubtreeOnly(t *testing.T) {
	tree := privateTree(t)
	actor := ActorScope{Role: RoleOrgAdmin, TenantID: authTenant, AssignedOrgNodeID: "hq"}
	placement := EmployeePlacementRef{EmployeeID: "rep-1", TenantID: authTenant, OrgNodeID: "inside"}

	ok, err := CanApprovePunch(actor, placement, tree)
	if err != nil || !ok {
		t.Fatal("ORG_ADMIN at HQ must approve subtree employee")
	}
}

func TestCanReadPunch_AuditorAllowedWriteDenied(t *testing.T) {
	tree := municipalTree(t)
	actor := ActorScope{Role: RoleAuditor, TenantID: authTenant}
	placement := EmployeePlacementRef{EmployeeID: "nurse-1", TenantID: authTenant, OrgNodeID: "nursing"}

	canRead, err := CanReadPunch(actor, placement, tree)
	if err != nil || !canRead {
		t.Fatal("AUDITOR must read tenant punches")
	}
	canApprove, _ := CanApprovePunch(actor, placement, tree)
	if canApprove {
		t.Fatal("AUDITOR must not approve punches")
	}
}

func TestCanExportPayroll_HRWithinTenant(t *testing.T) {
	actor := ActorScope{Role: RoleHRAnalyst, TenantID: authTenant}
	if !CanExportPayroll(actor, authTenant) {
		t.Fatal("HR_ANALYST must export within tenant")
	}
	if CanExportPayroll(actor, "other-tenant") {
		t.Fatal("HR_ANALYST must not export cross-tenant")
	}
}

func TestCanApprovePunch_CrossTenantDenied(t *testing.T) {
	tree := municipalTree(t)
	actor := ActorScope{Role: RoleSuperAdmin, TenantID: "other-tenant"}
	placement := EmployeePlacementRef{EmployeeID: "nurse-1", TenantID: authTenant, OrgNodeID: "nursing"}

	_, err := CanApprovePunch(actor, placement, tree)
	if err != ErrCrossTenantAccess {
		t.Fatalf("cross-tenant must be denied, got %v", err)
	}
}
