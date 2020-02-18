package dynalist

type Node interface {
	NodeID() string
	NodeChildren() []string
}

type NodeTree struct {
	byID  map[string]Node
	roots []Node
}

func NewNodeTree(files []Node) *NodeTree {
	tree := &NodeTree{
		byID: make(map[string]Node),
	}

	for _, f := range files {
		tree.byID[f.NodeID()] = f
	}

	tree.roots = tree.ScanRoots()

	return tree
}

func (tree *NodeTree) ScanRoots() []Node {
	candidates := make(map[string]Node)
	for id, f := range tree.byID {
		candidates[id] = f
	}
	for _, f := range tree.byID {
		for _, child := range f.NodeChildren() {
			delete(candidates, child)
		}
	}

	roots := make([]Node, 0, len(candidates))
	for _, f := range candidates {
		roots = append(roots, f)
	}

	return roots
}

func (tree *NodeTree) Get(id string) Node {
	return tree.byID[id]
}

type NodeWalkFunc func(parents []Node, f Node) error

func nodeWalk(tree *NodeTree, walkFn NodeWalkFunc, parents []Node, f Node) error {
	if err := walkFn(parents, f); err != nil {
		return err
	}

	pparents := append(parents, f)
	for _, id := range f.NodeChildren() {
		child := tree.Get(id)
		if err := nodeWalk(tree, walkFn, pparents, child); err != nil {
			return err
		}
	}

	return nil
}

func (tree *NodeTree) Walk(walkFn NodeWalkFunc) error {
	for _, f := range tree.roots {
		parents := make([]Node, 0)
		if err := nodeWalk(tree, walkFn, parents, f); err != nil {
			return err
		}
	}

	return nil
}
