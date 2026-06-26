package organization

// PathFromRoot returns node IDs from top-level ancestor down to nodeID (inclusive).
func (t *OrgTree) PathFromRoot(nodeID string) ([]string, error) {
	chain, err := t.pathToRoot(nodeID)
	if err != nil {
		return nil, err
	}
	for left, right := 0, len(chain)-1; left < right; left, right = left+1, right-1 {
		chain[left], chain[right] = chain[right], chain[left]
	}
	return chain, nil
}

// ResolveEffectivePolicy walks root → node and merges tenant root with node overrides.
func (t *OrgTree) ResolveEffectivePolicy(
	nodeID string,
	root AttendancePolicy,
	overrides map[string]PolicyOverride,
) (AttendancePolicy, error) {
	path, err := t.PathFromRoot(nodeID)
	if err != nil {
		return AttendancePolicy{}, err
	}

	layers := make([]PolicyOverride, 0, len(path))
	for _, id := range path {
		if override, ok := overrides[id]; ok {
			layers = append(layers, override)
		}
	}
	return EffectivePolicy(root, layers), nil
}

func (t *OrgTree) pathToRoot(nodeID string) ([]string, error) {
	chain := make([]string, 0, len(t.nodes))
	current := nodeID
	for {
		node, ok := t.nodes[current]
		if !ok {
			return nil, ErrNodeNotFound
		}
		chain = append(chain, current)
		if node.ParentID == "" {
			return chain, nil
		}
		current = node.ParentID
	}
}
