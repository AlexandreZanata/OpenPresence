package organization

// allowedChildren maps parent type to permitted child types.
// Root (empty parent) accepts only DIVISION — see docs/ORGANIZATION.md.
var allowedChildren = map[OrgNodeType][]OrgNodeType{
	OrgNodeTypeDivision:   {OrgNodeTypeDepartment, OrgNodeTypeLocation, OrgNodeTypeWorkSite},
	OrgNodeTypeDepartment: {OrgNodeTypeSection, OrgNodeTypeTeam, OrgNodeTypeWorkSite},
	OrgNodeTypeSection:    {OrgNodeTypeTeam},
	OrgNodeTypeLocation:   {OrgNodeTypeDepartment},
}

var allowedRootTypes = []OrgNodeType{OrgNodeTypeDivision}

// CanBeChild reports whether childType is valid under parentType.
// When parentType is empty, child is a direct child of the tenant root.
func CanBeChild(parentType OrgNodeType, childType OrgNodeType, isRoot bool) bool {
	if isRoot {
		return containsType(allowedRootTypes, childType)
	}
	allowed, ok := allowedChildren[parentType]
	if !ok {
		return false
	}
	return containsType(allowed, childType)
}

func containsType(types []OrgNodeType, target OrgNodeType) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}
