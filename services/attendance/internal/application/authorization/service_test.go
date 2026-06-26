package authorization_test

import (
	"errors"
	"testing"

	appauth "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/authorization"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

const tenant = "city-example"

type stubTreeReader struct {
	trees map[string]*organization.OrgTree
	err   error
}

func (s stubTreeReader) Tree(tenantID string) (*organization.OrgTree, error) {
	if s.err != nil {
		return nil, s.err
	}
	tree, ok := s.trees[tenantID]
	if !ok {
		return nil, organization.ErrEmptyTree
	}
	return tree, nil
}

func municipalTree(t *testing.T) *organization.OrgTree {
	t.Helper()
	nodes := []organization.OrgNode{
		{ID: "health", TenantID: tenant, Type: organization.OrgNodeTypeDivision, Name: "Health"},
		{ID: "nursing", TenantID: tenant, ParentID: "health", Type: organization.OrgNodeTypeDepartment, Name: "Nursing"},
		{ID: "education", TenantID: tenant, Type: organization.OrgNodeTypeDivision, Name: "Education"},
		{ID: "school-a", TenantID: tenant, ParentID: "education", Type: organization.OrgNodeTypeLocation, Name: "School A"},
	}
	tree, err := organization.BuildTree(tenant, nodes)
	if err != nil {
		t.Fatalf("build tree: %v", err)
	}
	return tree
}

func TestPunchAuthorizationService_ApproveHealthManager(t *testing.T) {
	svc := appauth.NewPunchAuthorizationService(stubTreeReader{trees: map[string]*organization.OrgTree{tenant: municipalTree(t)}})
	actor := organization.ActorScope{Role: organization.RoleManager, TenantID: tenant, AssignedOrgNodeID: "health"}
	placement := organization.EmployeePlacementRef{EmployeeID: "nurse-1", TenantID: tenant, OrgNodeID: "nursing"}

	ok, err := svc.ApprovePunch(actor, placement)
	if err != nil || !ok {
		t.Fatal("service must allow health manager approval")
	}
}

func TestPunchAuthorizationService_RejectCrossSecretariat(t *testing.T) {
	svc := appauth.NewPunchAuthorizationService(stubTreeReader{trees: map[string]*organization.OrgTree{tenant: municipalTree(t)}})
	actor := organization.ActorScope{Role: organization.RoleManager, TenantID: tenant, AssignedOrgNodeID: "health"}
	placement := organization.EmployeePlacementRef{EmployeeID: "teacher-1", TenantID: tenant, OrgNodeID: "school-a"}

	ok, err := svc.ApprovePunch(actor, placement)
	if err != nil || ok {
		t.Fatal("service must deny cross-secretariat approval")
	}
}

func TestPunchAuthorizationService_AuditorReadNoWrite(t *testing.T) {
	svc := appauth.NewPunchAuthorizationService(stubTreeReader{trees: map[string]*organization.OrgTree{tenant: municipalTree(t)}})
	actor := organization.ActorScope{Role: organization.RoleAuditor, TenantID: tenant}
	placement := organization.EmployeePlacementRef{EmployeeID: "nurse-1", TenantID: tenant, OrgNodeID: "nursing"}

	canRead, err := svc.ReadPunch(actor, placement)
	if err != nil || !canRead {
		t.Fatal("auditor read must be allowed")
	}
	if svc.WritePunch(actor) {
		t.Fatal("auditor write must be denied")
	}
}

func TestPunchAuthorizationService_CrossTenantDenied(t *testing.T) {
	svc := appauth.NewPunchAuthorizationService(stubTreeReader{trees: map[string]*organization.OrgTree{tenant: municipalTree(t)}})
	actor := organization.ActorScope{Role: organization.RoleSuperAdmin, TenantID: "other-tenant"}
	placement := organization.EmployeePlacementRef{EmployeeID: "nurse-1", TenantID: tenant, OrgNodeID: "nursing"}

	_, err := svc.ApprovePunch(actor, placement)
	if !errors.Is(err, organization.ErrCrossTenantAccess) {
		t.Fatalf("expected cross-tenant error, got %v", err)
	}
}
