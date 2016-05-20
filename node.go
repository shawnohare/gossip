package gossip

import "errors"

// Node in a parsed search tree.  It contains pointers to its parent node,
// if any, and all of its children.
type Node struct {
	parent   *Node
	children []*Node
	verb     int    // modal verb of the query: must (not), should.
	phrase   string // phrase literal if this query is a leaf.
	src      string
}

// depth is a tail recursive node depth function.
func depth(node *Node, k int) int {
	if node.Parent() == nil {
		return k
	}
	return depth(node.Parent(), k+1)
}

// IsValid reports whether the node represents a semantically valid node in
// a parsed search tree.
func (n *Node) IsValid() bool {
	if n == nil {
		return false
	}
	if n.IsLeaf() && n.phrase == "" {
		return false
	}
	return true
}

// Depth reports the node's depth in the tree of which it a member.
func (n *Node) Depth() int {
	return depth(n, 0)
}

// IsLeaf reports whether the node is a leaf node.  Nil instances are
// considered leaves.
func (n *Node) IsLeaf() bool {
	return len(n.Children()) == 0
}

// AddChild adds the input as a child of the instance.
func (n *Node) AddChild(child *Node) {
	if child == nil {
		return
	}

	if n == nil {
		n = new(Node)
	}
	child.parent = n
	n.children = append(n.children, child)
}

// NewChildren creates a child node.
func (n *Node) NewChild() *Node {
	if n == nil {
		n = new(Node)
	}
	child := new(Node)
	n.AddChild(child)
	return child
}

// NewSibling creates a sibling node of the instance.  The sibling and the
// instance have the same parent, which also contains both nodes as children.
func (n *Node) NewSibling() (*Node, error) {
	if n == nil || n.parent == nil {
		return nil, errors.New("Cannot create sibling Query without a parent.")
	}
	return n.parent.NewChild(), nil
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

// Verb returns the integer code the the modal verb attached to the node.
// If the node is
func (n *Node) Verb() int {
	if n.IsValid() {
		return n.verb
	}
	return VerbError
}
