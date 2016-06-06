package gossip

import "strings"

// Node in a parsed search tree.  It contains pointers to its parent node,
// if any, and all of its children.
type Node struct {
	parent   *Node
	children []*Node
	verb     rune   // modal verb of the query: must (not), should.
	phrase   string // phrase literal if this query is a leaf.
}

// IsLeaf reports whether the node is a leaf.
// Nil instances are not considered leaves.
func (n *Node) IsLeaf() bool {
	if n == nil {
		return false
	}
	return len(n.Children()) == 0
}

// IsValid recursively determines if every node is valid.
func (n *Node) IsValid() bool {
	if n == nil {
		return false
	}

	if n.parent == n {
		return false
	}

	if !IsVerb(n.verb) {
		return false
	}

	if n.IsLeaf() {
		return n.phrase != "" && IsVerb(n.verb)
	}

	// Preliminary non-leaf check.
	if !n.IsLeaf() && n.phrase != "" {
		return false
	}
	// Check the children of a valid non-leaf.
	for _, child := range n.Children() {
		if !child.IsValid() {
			return false
		}
	}
	return true
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
func (n *Node) Verb() rune {
	if n.IsValid() {
		return n.verb
	}
	return VerbError
}

// VerbString returns a human readable version of the node's modal verb.
func (n *Node) VerbString() string {
	return VerbStringHuman(n.Verb())
}

// SetParent sets the node's parent and returns the instance.  If the
// parent input is the node itself, nothing is set.
func (n *Node) SetParent(parent *Node) *Node {
	if n == nil {
		n = NewNode()
	}
	if n != parent {
		n.parent = parent
	}
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

// SetVerb checks whether the input is a valid modal verb representation and
// if so sets the node's verb attribute and returns the instance.
// For example:
//   var n *node
//   n = n.SetVerb(Must)
func (n *Node) SetVerb(verb rune) *Node {
	if n == nil {
		n = NewNode()
	}
	if IsVerb(verb) {
		n.verb = verb
	}
	return n
}

// IsRoot reports if the node is the root of a tree.
func (n *Node) IsRoot() bool {
	if n == nil {
		return false
	}
	return n.Parent() == nil
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
		return n
	}
	if n == nil {
		n = NewNode()
	}

	// Do not add the instance to itself.
	if n == child {
		return n
	}

	child.parent = n
	n.children = append(n.children, child)
	return n
}

// NewChildren creates and returns a new child node, provided the instance
// is not nil.
func (n *Node) NewChild() *Node {
	if n == nil {
		return nil
	}
	c := NewNode()
	n.AddChild(c)
	return c
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

// Root returns the root node of the tree containing the instance.
func (n *Node) Root() *Node {

	tmp := n
	for tmp.Parent() != nil {
		tmp = tmp.Parent()
	}

	return tmp
}

// leaves is a tail recursive function that computes the leaves of a tree.
func leaves(remaining []*Node, current []*Node) []*Node {
	if len(remaining) == 0 {
		return current
	}

	var newRemaining []*Node
	for _, node := range remaining {
		if node.IsLeaf() {
			current = append(current, node)
		} else {
			newRemaining = append(newRemaining, node.Children()...)
		}
	}
	return leaves(newRemaining, current)
}

// Leaves returns all external nodes of the subtree defined by the node.
func (n *Node) Leaves() []*Node {
	if n.IsLeaf() {
		return []*Node{n}
	}
	return leaves(n.Children(), nil)
}

// Height returns the height of the node.  This is the max depth of
// all the node's leaves.
func (n *Node) Height() int {
	var maxDepth int
	for _, node := range n.Leaves() {
		if d := node.Depth(); d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}

	// Just return the phrase if the root is a leaf.
	vs := VerbString(n.verb)
	if n.IsLeaf() {
		if n.IsValid() {
			return vs + `"` + n.Phrase() + `"`
		}
		return ""
	}

	// Otherwise convert each child to a string.  The result might
	// look something like [+w0 +"phrase1" -[...]]
	strs := make([]string, len(n.children))
	for i, child := range n.children {
		substring := child.String()
		if substring == "" {
			return ""
		}
		strs[i] = substring
	}

	s := vs + "[" + strings.Join(strs, ", ") + "]"
	return s
}

func (n *Node) JSON() *JSON {
	if !n.IsValid() {
		return nil
	}

	jnode := &JSON{
		Verb:   n.VerbString(),
		Phrase: n.phrase,
	}

	if n.IsLeaf() {
		return jnode
	}

	jnode.Children = make([]*JSON, len(n.children))
	for i, child := range n.children {
		jnode.Children[i] = child.JSON()
	}
	return jnode
}

// NewNode produces a leaf node with the default modal verb of Should,
// an empty phrase.  This node is not semantically valid until a phrase
// is set.
func NewNode() *Node {
	return &Node{
		verb:   Should,
		phrase: "",
	}
}

// New is an alias for NewNode.
// func New() *Node {
// 	return NewNode()
// }
