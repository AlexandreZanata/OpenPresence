package workforce

import (
	"errors"
	"testing"
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
)

const testTenant = "tenant-workforce"

var (
	jan1  = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	jun1  = time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	jul1  = time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	aug15 = time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC)
)

func municipalHealthTree(t *testing.T) *organization.OrgTree {
	t.Helper()
	nodes := []organization.OrgNode{
		{ID: "health", TenantID: testTenant, Type: organization.OrgNodeTypeDivision, Name: "Health Secretariat"},
		{ID: "hospital", TenantID: testTenant, ParentID: "health", Type: organization.OrgNodeTypeLocation, Name: "Municipal Hospital"},
		{ID: "nursing", TenantID: testTenant, ParentID: "hospital", Type: organization.OrgNodeTypeDepartment, Name: "Nursing"},
	}
	tree, err := organization.BuildTree(testTenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}
	return tree
}

func municipalEducationAndHealthTree(t *testing.T) *organization.OrgTree {
	t.Helper()
	nodes := []organization.OrgNode{
		{ID: "education", TenantID: testTenant, Type: organization.OrgNodeTypeDivision, Name: "Education Secretariat"},
		{ID: "health", TenantID: testTenant, Type: organization.OrgNodeTypeDivision, Name: "Health Secretariat"},
		{ID: "school-a", TenantID: testTenant, ParentID: "education", Type: organization.OrgNodeTypeLocation, Name: "School A"},
		{ID: "nursing", TenantID: testTenant, ParentID: "health", Type: organization.OrgNodeTypeDepartment, Name: "Nursing"},
	}
	tree, err := organization.BuildTree(testTenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}
	return tree
}

func privateSalesTree(t *testing.T) *organization.OrgTree {
	t.Helper()
	nodes := []organization.OrgNode{
		{ID: "hq", TenantID: testTenant, Type: organization.OrgNodeTypeDivision, Name: "HQ São Paulo"},
		{ID: "sales", TenantID: testTenant, ParentID: "hq", Type: organization.OrgNodeTypeDepartment, Name: "Sales"},
		{ID: "inside", TenantID: testTenant, ParentID: "sales", Type: organization.OrgNodeTypeTeam, Name: "Inside Sales"},
		{ID: "field", TenantID: testTenant, ParentID: "sales", Type: organization.OrgNodeTypeTeam, Name: "Field Sales"},
	}
	tree, err := organization.BuildTree(testTenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}
	return tree
}

func TestAssign_PublicPrimaryAtMunicipalHospitalNursing(t *testing.T) {
	svc := NewPlacementService(municipalHealthTree(t))
	placement := EmployeePlacement{
		ID: "pl-1", EmployeeID: "emp-servidor", TenantID: testTenant,
		OrgNodeID: "nursing", Type: PlacementTypePrimary, ValidFrom: jan1,
	}
	if err := svc.Assign(placement); err != nil {
		t.Fatalf("assign primary at nursing failed: %v", err)
	}

	active := svc.ActivePlacementsAt("emp-servidor", jun1)
	if len(active) != 1 || active[0].OrgNodeID != "nursing" {
		t.Fatal("employee must be primarily placed at Municipal Hospital / Nursing")
	}

	path, err := svc.EffectiveOrgPath("emp-servidor", jun1)
	if err != nil {
		t.Fatalf("org path failed: %v", err)
	}
	want := []string{"health", "hospital", "nursing"}
	for i := range want {
		if path[i] != want[i] {
			t.Fatalf("path[%d]=%s want %s", i, path[i], want[i])
		}
	}
}

func TestTransfer_PublicEducationToHealthClosesOldPlacement(t *testing.T) {
	svc := NewPlacementService(municipalEducationAndHealthTree(t))
	if err := svc.Assign(EmployeePlacement{
		ID: "pl-edu", EmployeeID: "emp-transfer", TenantID: testTenant,
		OrgNodeID: "school-a", Type: PlacementTypePrimary, ValidFrom: jan1,
	}); err != nil {
		t.Fatalf("initial assign failed: %v", err)
	}

	if err := svc.Transfer("emp-transfer", "nursing", "pl-health", jul1); err != nil {
		t.Fatalf("transfer failed: %v", err)
	}

	before := svc.ActivePlacementsAt("emp-transfer", jun1)
	if len(before) != 1 || before[0].OrgNodeID != "school-a" {
		t.Fatal("before transfer employee must remain at Education")
	}

	after := svc.ActivePlacementsAt("emp-transfer", aug15)
	if len(after) != 1 || after[0].OrgNodeID != "nursing" {
		t.Fatal("after transfer employee primary must be at Health / Nursing")
	}

	primary, err := svc.ActivePrimaryAt("emp-transfer", aug15)
	if err != nil || primary == nil || primary.ID != "pl-health" {
		t.Fatal("active primary must be the new health placement")
	}
}

func TestAssign_PrivatePrimarySalesSecondaryFieldTeam(t *testing.T) {
	svc := NewPlacementService(privateSalesTree(t))
	if err := svc.Assign(EmployeePlacement{
		ID: "pl-primary", EmployeeID: "emp-sales", TenantID: testTenant,
		OrgNodeID: "inside", Type: PlacementTypePrimary, ValidFrom: jan1,
	}); err != nil {
		t.Fatalf("primary assign failed: %v", err)
	}
	if err := svc.Assign(EmployeePlacement{
		ID: "pl-secondary", EmployeeID: "emp-sales", TenantID: testTenant,
		OrgNodeID: "field", Type: PlacementTypeSecondary, ValidFrom: jan1,
	}); err != nil {
		t.Fatalf("secondary assign failed: %v", err)
	}

	active := svc.ActivePlacementsAt("emp-sales", jun1)
	if len(active) != 2 {
		t.Fatalf("expected primary + secondary active, got %d", len(active))
	}
}

func TestAssign_RejectsSecondActivePrimary(t *testing.T) {
	svc := NewPlacementService(privateSalesTree(t))
	if err := svc.Assign(EmployeePlacement{
		ID: "pl-a", EmployeeID: "emp-dup", TenantID: testTenant,
		OrgNodeID: "inside", Type: PlacementTypePrimary, ValidFrom: jan1,
	}); err != nil {
		t.Fatalf("first primary failed: %v", err)
	}

	err := svc.Assign(EmployeePlacement{
		ID: "pl-b", EmployeeID: "emp-dup", TenantID: testTenant,
		OrgNodeID: "field", Type: PlacementTypePrimary, ValidFrom: jun1,
	})
	if !errors.Is(err, ErrDuplicatePrimary) {
		t.Fatalf("expected ErrDuplicatePrimary, got %v", err)
	}
}

func TestAssign_RejectsOrgNodeOutsideTenant(t *testing.T) {
	svc := NewPlacementService(municipalHealthTree(t))
	err := svc.Assign(EmployeePlacement{
		ID: "pl-bad", EmployeeID: "emp-x", TenantID: testTenant,
		OrgNodeID: "unknown-node", Type: PlacementTypePrimary, ValidFrom: jan1,
	})
	if !errors.Is(err, ErrNodeNotInTenant) {
		t.Fatalf("expected ErrNodeNotInTenant, got %v", err)
	}
}

func TestAssign_RejectsOverlappingSecondary(t *testing.T) {
	svc := NewPlacementService(privateSalesTree(t))
	base := EmployeePlacement{
		EmployeeID: "emp-sec", TenantID: testTenant,
		OrgNodeID: "field", Type: PlacementTypeSecondary, ValidFrom: jan1,
	}
	if err := svc.Assign(EmployeePlacement{ID: "sec-1", ValidFrom: jan1, EmployeeID: base.EmployeeID, TenantID: base.TenantID, OrgNodeID: base.OrgNodeID, Type: base.Type}); err != nil {
		t.Fatalf("first secondary failed: %v", err)
	}
	err := svc.Assign(EmployeePlacement{
		ID: "sec-2", EmployeeID: base.EmployeeID, TenantID: base.TenantID,
		OrgNodeID: "field", Type: PlacementTypeSecondary, ValidFrom: jun1,
	})
	if !errors.Is(err, ErrOverlappingSecondary) {
		t.Fatalf("expected ErrOverlappingSecondary, got %v", err)
	}
}

func TestActivePrimaryAt_NoDuplicateActivePrimary(t *testing.T) {
	svc := NewPlacementService(privateSalesTree(t))
	// Simulate bad state shouldn't happen via Assign, but query must detect double primary.
	svc.placements = []EmployeePlacement{
		{ID: "a", EmployeeID: "emp-bad", TenantID: testTenant, OrgNodeID: "inside", Type: PlacementTypePrimary, ValidFrom: jan1},
		{ID: "b", EmployeeID: "emp-bad", TenantID: testTenant, OrgNodeID: "field", Type: PlacementTypePrimary, ValidFrom: jan1},
	}
	_, err := svc.ActivePrimaryAt("emp-bad", jun1)
	if !errors.Is(err, ErrDuplicatePrimary) {
		t.Fatalf("expected ErrDuplicatePrimary from query, got %v", err)
	}
}
