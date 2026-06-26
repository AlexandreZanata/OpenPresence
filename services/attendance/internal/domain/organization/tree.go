package organization

// OrgTree is a validated tenant-scoped organization tree.
type OrgTree struct {
	TenantID string
	nodes    map[string]OrgNode
}

// BuildTree validates and constructs an OrgTree from a flat node list.
func BuildTree(tenantID string, nodes []OrgNode) (*OrgTree, error) {
	if tenantID == "" {
		return nil, ErrEmptyTenant
	}
	if len(nodes) == 0 {
		return nil, ErrEmptyTree
	}

	tree := &OrgTree{TenantID: tenantID, nodes: make(map[string]OrgNode, len(nodes))}
	for _, node := range nodes {
		if err := tree.insert(node); err != nil {
			return nil, err
		}
	}
	if err := tree.validateStructure(); err != nil {
		return nil, err
	}
	return tree, nil
}

// Node returns a copy of the node with the given id.
func (t *OrgTree) Node(id string) (OrgNode, bool) {
	node, ok := t.nodes[id]
	return node, ok
}

// Nodes returns all nodes in arbitrary order.
func (t *OrgTree) Nodes() []OrgNode {
	out := make([]OrgNode, 0, len(t.nodes))
	for _, node := range t.nodes {
		out = append(out, node)
	}
	return out
}

func (t *OrgTree) insert(node OrgNode) error {
	if node.ID == "" || node.Name == "" {
		return ErrInvalidNode
	}
	if node.TenantID != t.TenantID {
		return ErrMultipleTenants
	}
	if _, exists := t.nodes[node.ID]; exists {
		return ErrDuplicateNode
	}
	t.nodes[node.ID] = node
	return nil
}

func (t *OrgTree) validateStructure() error {
	for id := range t.nodes {
		if err := t.detectCycle(id); err != nil {
			return err
		}
	}
	for _, node := range t.nodes {
		if err := t.validateParentLink(node); err != nil {
			return err
		}
	}
	return nil
}

func (t *OrgTree) validateParentLink(node OrgNode) error {
	if node.ParentID == "" {
		if !CanBeChild("", node.Type, true) {
			return ErrInvalidRootType
		}
		return nil
	}
	parent, ok := t.nodes[node.ParentID]
	if !ok {
		return ErrOrphanNode
	}
	if !CanBeChild(parent.Type, node.Type, false) {
		return ErrInvalidChildType
	}
	return nil
}

func (t *OrgTree) detectCycle(startID string) error {
	visited := make(map[string]struct{}, len(t.nodes))
	current := startID
	for step := 0; step <= len(t.nodes); step++ {
		node, ok := t.nodes[current]
		if !ok {
			return nil
		}
		if _, seen := visited[current]; seen {
			return ErrCycle
		}
		visited[current] = struct{}{}
		if node.ParentID == "" {
			return nil
		}
		current = node.ParentID
	}
	return ErrCycle
}
