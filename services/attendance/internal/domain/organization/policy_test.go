package organization

import (
	"testing"
	"time"
)

func TestMergePolicy_InheritsUnsetOverrideFields(t *testing.T) {
	parent := DefaultPolicy()
	override := PolicyOverride{ToleranceMinutes: intPtr(25)}

	merged := MergePolicy(parent, override)
	if merged.ToleranceMinutes != 25 {
		t.Fatalf("expected tolerance override 25, got %d", merged.ToleranceMinutes)
	}
	if merged.OfflineSyncMaxAge != 8*time.Hour {
		t.Fatal("unset fields must inherit parent offlineSyncMaxAge (BR-011)")
	}
}

func TestEffectivePolicy_RootDivisionDepartmentChain(t *testing.T) {
	root := DefaultPolicy()
	overrides := []PolicyOverride{
		{ToleranceMinutes: intPtr(20)},
		{GeofenceRequired: boolPtr(false)},
	}

	effective := EffectivePolicy(root, overrides)
	if effective.ToleranceMinutes != 20 {
		t.Fatal("division tolerance override must apply")
	}
	if effective.GeofenceRequired {
		t.Fatal("department geofence override must apply")
	}
	if effective.OfflineSyncMaxAge != 8*time.Hour {
		t.Fatal("department must inherit root offline TTL")
	}
}

func TestResolveEffectivePolicy_TenantDivisionDepartment(t *testing.T) {
	const tenant = "tenant-policy"
	nodes := []OrgNode{
		{ID: "div", TenantID: tenant, Type: OrgNodeTypeDivision, Name: "Division"},
		{ID: "dept", TenantID: tenant, ParentID: "div", Type: OrgNodeTypeDepartment, Name: "Department"},
	}
	tree, err := BuildTree(tenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}

	overrides := map[string]PolicyOverride{
		"div":  {ToleranceMinutes: intPtr(20)},
		"dept": {BiometricRequired: boolPtr(false)},
	}

	effective, err := tree.ResolveEffectivePolicy("dept", DefaultPolicy(), overrides)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if effective.ToleranceMinutes != 20 {
		t.Fatal("department must inherit division tolerance")
	}
	if effective.BiometricRequired {
		t.Fatal("department local override must disable biometric")
	}
}

func TestResolveEffectivePolicy_PublicSecretariatHospital(t *testing.T) {
	const tenant = "city-example"
	nodes := []OrgNode{
		{ID: "health", TenantID: tenant, Type: OrgNodeTypeDivision, Name: "Health Secretariat"},
		{ID: "hospital", TenantID: tenant, ParentID: "health", Type: OrgNodeTypeLocation, Name: "Municipal Hospital"},
		{ID: "nursing", TenantID: tenant, ParentID: "hospital", Type: OrgNodeTypeDepartment, Name: "Nursing"},
	}
	tree, err := BuildTree(tenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}

	root := PublicSectorPreset()
	overrides := map[string]PolicyOverride{
		"health":   {OfflineSyncMaxAge: durationPtr(8 * time.Hour)},
		"hospital": {ToleranceMinutes: intPtr(45)},
	}

	nursing, err := tree.ResolveEffectivePolicy("nursing", root, overrides)
	if err != nil {
		t.Fatalf("resolve nursing failed: %v", err)
	}
	if nursing.OfflineSyncMaxAge != 8*time.Hour {
		t.Fatal("public secretariat keeps BR-011 offline TTL at 8h")
	}
	if nursing.ToleranceMinutes != 45 {
		t.Fatal("hospital tolerance override must flow to nursing")
	}
	if !nursing.GeofenceRequired {
		t.Fatal("public preset keeps strict geofence")
	}
}

func TestResolveEffectivePolicy_PrivateHQDisablesOvertimeBranchReenables(t *testing.T) {
	const tenant = "acme"
	nodes := []OrgNode{
		{ID: "hq", TenantID: tenant, Type: OrgNodeTypeDivision, Name: "HQ São Paulo"},
		{ID: "branch", TenantID: tenant, Type: OrgNodeTypeDivision, Name: "Cuiabá Branch"},
	}
	tree, err := BuildTree(tenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}

	root := PrivateSectorPreset()
	disabled := OvertimePolicyDisabled
	standard := OvertimePolicyStandard
	overrides := map[string]PolicyOverride{
		"hq":     {OvertimePolicy: &disabled},
		"branch": {OvertimePolicy: &standard},
	}

	hq, err := tree.ResolveEffectivePolicy("hq", root, overrides)
	if err != nil {
		t.Fatalf("resolve hq failed: %v", err)
	}
	if hq.OvertimePolicy != OvertimePolicyDisabled {
		t.Fatal("HQ must disable overtime")
	}

	branch, err := tree.ResolveEffectivePolicy("branch", root, overrides)
	if err != nil {
		t.Fatalf("resolve branch failed: %v", err)
	}
	if branch.OvertimePolicy != OvertimePolicyStandard {
		t.Fatal("branch must re-enable standard overtime")
	}
}

func TestPublicSectorPreset_12x36ShiftDefaults(t *testing.T) {
	preset := PublicSectorPreset()
	if preset.WorkdayDuration != 12*time.Hour {
		t.Fatalf("public preset expects 12h workday, got %v", preset.WorkdayDuration)
	}
	if preset.ToleranceMinutes != 30 {
		t.Fatalf("public preset expects 30m tolerance, got %d", preset.ToleranceMinutes)
	}
}

func TestPathFromRoot_OrderRootToLeaf(t *testing.T) {
	const tenant = "path-tenant"
	nodes := []OrgNode{
		{ID: "a", TenantID: tenant, Type: OrgNodeTypeDivision, Name: "A"},
		{ID: "b", TenantID: tenant, ParentID: "a", Type: OrgNodeTypeDepartment, Name: "B"},
		{ID: "c", TenantID: tenant, ParentID: "b", Type: OrgNodeTypeTeam, Name: "C"},
	}
	tree, err := BuildTree(tenant, nodes)
	if err != nil {
		t.Fatalf("tree build failed: %v", err)
	}

	path, err := tree.PathFromRoot("c")
	if err != nil {
		t.Fatalf("path failed: %v", err)
	}
	want := []string{"a", "b", "c"}
	if len(path) != len(want) {
		t.Fatalf("path len %d, want %d", len(path), len(want))
	}
	for i := range want {
		if path[i] != want[i] {
			t.Fatalf("path[%d]=%s, want %s", i, path[i], want[i])
		}
	}
}

func intPtr(v int) *int                       { return &v }
func boolPtr(v bool) *bool                    { return &v }
func durationPtr(v time.Duration) *time.Duration { return &v }
