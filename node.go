package gossip

import (
	"fmt"
	"strings"
)

// Node in a parsed search tree.  It contains pointers to its parent node,
// if any, and all of its children.  Generally it is expected that the
// Node getter and setter methods are used to access the exported fields.
// These fields are exported primarily so that an instance can be marshalled
// into a JSON string.
type Node struct {
	Parent   *Node   `json:"-"`
	Children []*Node `json:"children"`
	Verb     Verb    `json:"verb"`   // Modal verb of the query: must (not), should.
	Phrase   string  `json:"phrase"` // Phrase literal if this query is a leaf.
}

// IsLeaf reports whether the node is a leaf.
// Nil instances are not considered leaves.
func (n *Node) IsLeaf() bool {
	if n == nil {
		return false
	}
	return len(n.Children) == 0
}

// IsValid recursively determines if every node is valid.
// Conditions that lead to an invalid node are:
// - The instance is nil.
// - The instance is its own parent or contains itself as a child.
// - The instance's verb is not one of the constants Must, Should, MustNot.
// - The instance is a leaf with an empty phrase.
// - The instance is a non-leaf but contains a phrase.
// - Any child is invalid.
func (n *Node) IsValid() bool {
	if n == nil {
		return false
	}

	if n.Parent == n {
		return false
	}

	if !n.Verb.IsValid() {
		return false
	}

	if n.IsLeaf() && n.Phrase == "" {
		return false
	}

	// Preliminary non-leaf check.
	if !n.IsLeaf() && n.Phrase != "" {
		return false
	}

	for _, child := range n.GetChildren() {
		if !child.IsValid() || child == n || child.Parent != n {
			return false
		}
	}
	return true
}

// GetParent is a helper get method defined for all instances.
func (n *Node) GetParent() *Node {
	if n == nil {
		return nil
	}
	return n.Parent
}

// GetChildren is a helper get method defined for all instances.
func (n *Node) GetChildren() []*Node {
	if n == nil {
		return nil
	}
	return n.Children
}

// Tree sets the Parent field of the instance's descents to the appropriate
// node.  This effectively creates a subtree with the instance as a root,
// although the instance's Parent is not removed.
func (n *Node) Tree() *Node {
	if n == nil || n.IsLeaf() {
		return n
	}

	for _, child := range n.Children {
		_ = child.SetParent(n).Tree()
	}

	return n
}

// GetPhrase is a helper get method defined for all instances.
func (n *Node) GetPhrase() string {
	if n == nil {
		return ""
	}
	return n.Phrase
}

// GetVerb returns the integer code the the modal verb associated to the node.
func (n *Node) GetVerb() Verb {
	if n == nil {
		return VerbError
	}
	return n.Verb
}

// SetParent sets the node's parent and returns the instance.  If the
// parent input is the node itself, nothing is set.  This will not add the
// instance to the parent node's children, however.  For this functionality,
// see the AddChild function.
func (n *Node) SetParent(parent *Node) *Node {
	if n == nil {
		n = NewNode()
	}
	if n != parent {
		n.Parent = parent
	}
	return n
}

// SetPhrase sets the node's phrase and returns the instance.
func (n *Node) SetPhrase(phrase string) *Node {
	if n == nil {
		n = NewNode()
	}
	n.Phrase = phrase
	return n
}

// SetVerb sets the node's Verb field to the input verb, regardless
// of the input's semantic validity.
func (n *Node) SetVerb(verb Verb) *Node {
	if n == nil {
		n = NewNode()
	}
	n.Verb = verb
	return n
}

// IsRoot reports if the non-nil node is the root of a tree.
func (n *Node) IsRoot() bool {
	if n == nil {
		return false
	}
	return n.GetParent() == nil
}

// depth is a tail recursive node depth function.
func depth(node *Node, k int) int {
	if node.GetParent() == nil {
		return k
	}
	return depth(node.GetParent(), k+1)
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

	child.Parent = n
	n.Children = append(n.Children, child)
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

	if len(n.Children) != len(m.Children) {
		return false
	}

	if n.IsLeaf() && m.IsLeaf() {
		return n.Phrase == m.Phrase && n.Verb == m.Verb
	}

	for i, ni := range n.Children {
		if !ni.Equals(m.Children[i]) {
			return false
		}
	}
	return true
}

// Root returns the root node of the tree containing the instance.
func (n *Node) Root() *Node {

	tmp := n
	for tmp.GetParent() != nil {
		tmp = tmp.GetParent()
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
			newRemaining = append(newRemaining, node.GetChildren()...)
		}
	}
	return leaves(newRemaining, current)
}

// Leaves returns all external nodes of the subtree defined by the node.
func (n *Node) Leaves() []*Node {
	if n.IsLeaf() {
		return []*Node{n}
	}
	return leaves(n.GetChildren(), nil)
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
	if n.IsLeaf() {
		if n.IsValid() {
			return fmt.Sprintf("%s\"%s\"", n.Verb, n.Phrase)
		}
		return ""
	}

	// Otherwise convert each child to a string.  The result might
	// look something like [+w0 +"phrase1" -[...]]
	strs := make([]string, len(n.Children))
	for i, child := range n.Children {
		substring := child.String()
		if substring == "" {
			return ""
		}
		strs[i] = substring
	}

	return fmt.Sprintf("%s[%s]", n.Verb, strings.Join(strs, ", "))
}

// NewNode produces a leaf node with the default modal verb of Should,
// an empty phrase.  This node is not semantically valid until a phrase
// is set.
func NewNode() *Node {
	return &Node{
		Verb:   Should,
		Phrase: "",
	}
}

// New is an alias for NewNode.
// func New() *Node {
// 	return NewNode()
// }
