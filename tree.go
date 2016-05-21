package gossip

// Tree is represents the parsed search tree of some input query.
// It is effectively a root node together with some additional
// methods.
type Tree struct {
	root *Node
}

// Root returns the tree's root node.
func (t *Tree) Root() *Node {
	if t == nil {
		return nil
	}
	return t.root
}

// SetRoot sets the tree's root
// Caution: any existing root will be overwritten.
func (t *Tree) SetRoot(root *Node) {
	if t == nil {
		t = NewTree()
	}
	t.root = root
}

// IsValid recursively determines if every node is valid.
func (t *Tree) IsValid() bool {
	r := t.Root()
	if r == nil {
		return false
	}

	if r.IsLeaf() {
		return r.IsValid()
	}

	for _, child := range r.Children() {
		if !NewTreeFromRoot(child).IsValid() {
			return false
		}
	}
	return true
}

// String returns the source query that generated the tree.
func (t *Tree) String() string {
	// if t == nil {
	// 	return ""
	// }
	// return t.src

	if t == nil || t.root == nil {
		return ""
	}

	vs := verbString(t.root.verb)

	// Just return the phrase if the root is a leaf.
	if t.root.IsLeaf() {
		return vs + t.root.phrase
	}

	// Otherwise convert each child to a string.  The result might
	// look something like [+w0 +"phrase1" -[...]]
	var s string
	for _, child := range t.root.children {
		subtree := NewTreeFromRoot(child)
		s += subtree.String() + " "
	}
	s = vs + "[" + s + "]"

	return s
}

// Equals reports whether the trees are semantically equivalent.
// In particular, it is a poset isomorphism with respect to the
// total ordering on phrases induced by the source query.
func (t *Tree) Equals(r *Tree) bool {
	return t.Root().Equals(r.Root())
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

func NewTree() *Tree {
	return &Tree{
		root: NewNode(),
	}
}

func NewTreeFromRoot(root *Node) *Tree {
	return &Tree{
		root: root,
	}
}
