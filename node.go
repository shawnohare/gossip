package gossip

import "errors"

// Node in a parsed search tree.  It contains pointers to its parent node,
// if any, and all of its children.
type Node struct {
	parent   *Node
	children []*Node
	verb     int    // modal verb of the query: must (not), should.
	phrase   string // phrase literal if this query is a leaf.
}

// Parent of the instance.  Is nil if the node is a root.
func (n *Node) Parent() *Node {
	if n == nil {
		return nil
	}
	return n.parent
}

// Children of the instance.
func (n *Node) Children() []*Node {
	if n == nil {
		return nil
	}
	return n.children
}

func (n *Node) Phrase() string {
	if n.IsValid() {
		return n.phrase
	}
	return ""
}

// Verb returns the integer code the the modal verb attached to the node.
// If the node is
func (n *Node) Verb() int {
	if n.IsValid() {
		return n.verb
	}
	return VerbError
}

// SetParent sets the node's parent and returns the instance.
func (n *Node) SetParent(parent *Node) *Node {
	if n == nil {
		n = NewNode()
	}
	n.parent = parent
	return n
}

// SetPhrase sets the node's phrase and returns the instance.
func (n *Node) SetPhrase(phrase string) *Node {
	if n == nil {
		n = NewNode()
	}
	n.phrase = phrase
	return n
}

// SetVerb sets the node's modal verb and returns the instance.
// the instance.
func (n *Node) SetVerb(verb int) *Node {
	if n == nil {
		n = NewNode()
	}
	n.verb = verb
	return n
}

// IsLeaf reports whether the node is a leaf node.  Nil instances are
// considered leaves.
func (n *Node) IsLeaf() bool {
	return len(n.Children()) == 0
}

// IsValid reports whether the node represents a semantically valid node in
// a parsed search tree. Only the instance is expected.
func (n *Node) IsValid() bool {
	if n == nil {
		return false
	}

	if !n.IsLeaf() && n.phrase != "" {
		return false
	}
	if n.IsLeaf() && n.phrase == "" {
		return false
	}
	// Root should not have a verb.
	if n.Parent == nil && !n.IsLeaf() && n.verb != 0 {
		return false
	}
	return true
}

// depth is a tail recursive node depth function.
func depth(node *Node, k int) int {
	if node.Parent() == nil {
		return k
	}
	return depth(node.Parent(), k+1)
}

// Depth reports the node's depth in the tree of which it a member.
func (n *Node) Depth() int {
	return depth(n, 0)
}

// AddChild specifies that the input should be a child of the instance.
// This appends a new node to the instance's children, sets the instance
// to be the parent of that child, and returns the child. A nil input
// is ignored.
func (n *Node) AddChild(child *Node) *Node {
	if child == nil {
		return nil
	}

	if n == nil {
		n = NewNode()
	}
	child.parent = n
	n.children = append(n.children, child)
	return child
}

// NewChildren creates a child node.
func (n *Node) NewChild() *Node {
	return n.AddChild(NewNode())
}

// NewSibling creates a sibling node of the instance.  The sibling and the
// instance have the same parent, which also contains both nodes as children.
func (n *Node) NewSibling() (*Node, error) {
	if n == nil || n.parent == nil {
		return nil, errors.New("Cannot create sibling Query without a parent.")
	}
	return n.parent.NewChild(), nil
}

// Equals reports whether the instance and input define semantically
// equal parsed subtrees.
func (n *Node) Equals(m *Node) bool {
	if !n.IsValid() || !m.IsValid() {
		return false
	}

	if len(n.children) != len(m.children) {
		return false
	}

	if n.IsLeaf() && m.IsLeaf() {
		return n.phrase == m.phrase && n.verb == m.verb
	}

	ok := true
	for i, ni := range n.children {
		if !ni.Equals(m.children[i]) {
			return false
		}
	}
	return ok
}

func NewNode() *Node {
	return new(Node)
}
