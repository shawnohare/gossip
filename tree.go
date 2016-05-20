package gossip

// Tree is represents the parsed search tree of of the input Query.
type Tree struct {
	root *Node
	src  string
}

func (t *Tree) Root() *Node {
	if t == nil {
		return nil
	}
	return t.root
}

// String returns the source query that generated the tree.
func (t *Tree) String() string {
	if t == nil {
		return ""
	}
	return t.src
}

// leaves is a tail recursive function that computes the leaves of a tree.
func leaves(children []*Node, current []*Node) []*Node {
	if len(children) == 0 {
		return current
	}

	var nonLeaves []*Node
	for _, node := range children {
		if node.IsLeaf() {
			current = append(current, node)
		} else {
			nonLeaves = append(nonLeaves, node)
		}
	}
	return leaves(nonLeaves, current)
}

func (t *Tree) Leaves() []*Node {
	return leaves(t.Root().Children(), nil)
}

func (t *Tree) Height() int {
	var maxDepth int
	for _, node := range t.Leaves() {
		if d := node.Depth(); d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth
}
