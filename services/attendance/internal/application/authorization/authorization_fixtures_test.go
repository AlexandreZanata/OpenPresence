package authorization_test

import (
	"testing"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

func municipalTree(t *testing.T, tenantID string) *organization.OrgTree {
	t.Helper()
	nodes := []organization.OrgNode{
		{ID: "health", TenantID: tenantID, Type: organization.OrgNodeTypeDivision, Name: "Health"},
		{ID: "nursing", TenantID: tenantID, ParentID: "health", Type: organization.OrgNodeTypeDepartment, Name: "Nursing"},
		{ID: "education", TenantID: tenantID, Type: organization.OrgNodeTypeDivision, Name: "Education"},
		{ID: "school-a", TenantID: tenantID, ParentID: "education", Type: organization.OrgNodeTypeLocation, Name: "School A"},
	}
	tree, err := organization.BuildTree(tenantID, nodes)
	if err != nil {
		t.Fatalf("build tree: %v", err)
	}
	return tree
}
