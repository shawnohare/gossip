package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeMethodsForNil(t *testing.T) {
	var n *Node
	assert.False(t, n.Equals(n))
	assert.False(t, n.IsValid())
	assert.Equal(t, 0, n.Height())
	assert.Equal(t, "", n.String())
	assert.Nil(t, n.Root())
	assert.Len(t, n.Leaves(), 0)
	assert.False(t, n.IsLeaf())
	assert.Equal(t, VerbError, n.GetVerb())
	assert.Len(t, n.GetChildren(), 0)
	assert.Equal(t, "", n.GetPhrase())
	assert.Equal(t, 0, n.Depth())
	assert.Nil(t, n.GetParent())
	assert.NotNil(t, n.AddChild(&Node{}))
	assert.Nil(t, n.AddChild(nil))
	assert.Nil(t, n.NewChild())
	assert.NotNil(t, n.SetParent(nil))
	assert.NotNil(t, n.SetParent(NewNode()))

}

func TestNodeAddChildForNil(t *testing.T) {
	var n *Node
	n = n.AddChild(NewNode())
	assert.NotNil(t, n)
	assert.Len(t, n.GetChildren(), 1)
}

func TestNodeAddChildSelf(t *testing.T) {
	m := NewNode()
	m.AddChild(m)
	assert.Len(t, m.GetChildren(), 0)
}

func TestNodeSetParentSelf(t *testing.T) {
	m := NewNode()
	m.SetParent(m)
	assert.Nil(t, m.GetParent())
}

func TestNodeIsLeafAfterAdds(t *testing.T) {
	tests := []struct {
		in  *Node
		out bool
	}{
		{nil, false},
		{NewNode(), true},
		{NewNode().NewChild().GetParent(), false},
		{NewNode().AddChild(NewNode()).GetParent(), false},
		// AddChild(nil) is nil, so its parent is also nil, hence not a leaf.
		{NewNode().AddChild(nil).GetParent(), false},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) with children %#v fails", i, tt.in.GetChildren())
		actual := tt.in.IsLeaf()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodeVerb(t *testing.T) {
	tests := []struct {
		in  *Node
		out Verb
	}{
		{nil, VerbError},
		{&Node{}, 0},
		{&Node{Verb: Must}, Must},
		{&Node{Verb: Should}, Should},
		{&Node{Verb: MustNot}, MustNot},
		{&Node{Verb: VerbError}, VerbError},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.GetVerb()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodeSetParent(t *testing.T) {
	var n *Node
	p := &Node{}
	n = n.SetParent(p)
	assert.Equal(t, p, n.Parent)
}

func TestSetVerb(t *testing.T) {
	var n *Node
	n = n.SetVerb(Must)
	assert.Equal(t, Must, n.Verb)
}

func TestSetVerbInvalidInput(t *testing.T) {
	var n *Node
	n = n.SetVerb(-10)
	assert.Equal(t, Verb(-10), n.Verb)
}

func TestSetPhrase(t *testing.T) {
	var n *Node
	n = n.SetPhrase("0")
	assert.Equal(t, "0", n.Phrase)
}

func TestNodeIsValid(t *testing.T) {
	n := NewNode()
	n.Parent = n

	tests := []struct {
		in  *Node
		out bool
	}{
		{nil, false},
		{&Node{}, false},
		{n, false}, // parent is self
		{
			&Node{
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x1",
					},
				},
				Phrase: "this phrase invalidates the node",
			},
			false,
		},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case %d", i)
		actual := tt.in.IsValid()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodeTree(t *testing.T) {
	r := NewNode()
	c0 := NewNode()
	c00 := NewNode()
	c0.Children = []*Node{c00}
	r.Children = []*Node{c0}
	r.Tree()
	assert.Equal(t, c00.Parent, c0)
	assert.Equal(t, c0.Parent, r)
}

func TestNodeEquals(t *testing.T) {
	tests := []struct {
		n0  *Node
		n1  *Node
		out bool
	}{
		{nil, nil, false},
		{nil, &Node{}, false},
		{nil, &Node{Phrase: "x"}, false},
		{NewNode(), NewNode(), false},
		{NewNode(), &Node{Phrase: "x", Verb: Should}, false},
		{&Node{Phrase: "x", Verb: Should}, &Node{Phrase: "y", Verb: Should}, false},
		{&Node{Verb: Must, Phrase: "x"}, &Node{Verb: Should, Phrase: "x"}, false},
		{&Node{Verb: Must, Phrase: "x"}, &Node{Verb: MustNot, Phrase: "x"}, false},
		{&Node{Phrase: "x", Verb: Should}, &Node{Phrase: "x", Verb: Should}, true},
		// 9. Basic test with children.
		{
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x1",
					},
				},
			},
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x1",
					},
				},
			},
			true,
		},
		// 10. Trees of height 2.
		{
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb: Should,
						Children: []*Node{
							&Node{
								Verb:   Must,
								Phrase: "x1",
							},
						},
					},
				},
			},
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb: Should,
						Children: []*Node{
							&Node{
								Verb:   Must,
								Phrase: "x1",
							},
						},
					},
				},
			},
			true,
		},
		// At least one node is invalid.
		{
			&Node{
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x1",
					},
				},
				Phrase: "this phrase invalidates the node",
			},
			&Node{
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x1",
					},
				},
			},
			false,
		},
		// Children are different lengths.
		{
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x1",
					},
				},
			},
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x1",
					},
					&Node{
						Verb:   Must,
						Phrase: "x2",
					},
				},
			},
			false,
		},
		// Children are different.
		{
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "x",
					},
				},
			},
			&Node{
				Verb: Should,
				Children: []*Node{
					&Node{
						Verb:   Must,
						Phrase: "y",
					},
				},
			},
			false,
		},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case %d", i)
		actual := tt.n0.Tree().Equals(tt.n1.Tree())
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodePhrase(t *testing.T) {
	tests := []struct {
		in  *Node
		out string
	}{
		{nil, ""},
		{&Node{}, ""},
		{&Node{Parent: &Node{}}, ""},
		{&Node{Phrase: "x", Verb: Should}, "x"},
		{&Node{Children: []*Node{&Node{Phrase: "x", Verb: Should}}}, ""},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.GetPhrase()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodeDepth(t *testing.T) {
	tests := []struct {
		in  *Node
		out int
	}{
		{nil, 0},
		{&Node{}, 0},
		{&Node{Parent: &Node{}}, 1},
		{&Node{Parent: &Node{Parent: &Node{}}}, 2},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.Depth()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestTreeIsValid(t *testing.T) {
	h1 := NewNode()
	h1.NewChild().SetPhrase("test")

	h2 := NewNode()
	h2.NewChild().NewChild().SetPhrase("test")

	h3 := NewNode()
	h3.NewChild().SetPhrase("test").NewChild().SetPhrase("test")

	tests := []struct {
		in  *Node
		out bool
	}{
		{nil, false},
		{NewNode(), false},
		{h1, true},
		{h2, true},
		{h3, false},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.IsValid()
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestNodeLeavesExact(t *testing.T) {
	leaves := []*Node{
		NewNode(),
		NewNode(), NewNode(),
		NewNode(), NewNode(), NewNode(), NewNode(),
	}
	// leaves looks like
	//        0
	//     0     0
	//    0 0   0 0
	leaves[0].AddChild(leaves[1])
	leaves[0].AddChild(leaves[2])
	leaves[1].AddChild(leaves[3])
	leaves[1].AddChild(leaves[4])
	leaves[2].AddChild(leaves[5])
	leaves[2].AddChild(leaves[6])

	n := NewNode()

	tests := []struct {
		in  *Node
		out []*Node
	}{
		{nil, nil},
		{n, []*Node{n}},
		{leaves[0], leaves[3:]},
		{leaves[1], leaves[3:5]},
		{leaves[2], leaves[5:]},
	}

	for i, tt := range tests {
		actual := tt.in.Leaves()
		msg := fmt.Sprintf("Test case (%d) fails", i)
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestNodeString(t *testing.T) {
	h1 := NewNode()
	h1.NewChild()

	h2 := NewNode()
	h2.NewChild().SetPhrase("x")
	h2.NewChild().SetPhrase("y")

	h3 := NewNode()
	c1 := h3.NewChild()
	c1.NewChild().SetPhrase("x")
	c1.NewChild().SetPhrase("y")
	c2 := h3.NewChild()
	c2.NewChild().SetPhrase("v").SetVerb(Must)
	c2.NewChild().SetPhrase("w").SetVerb(MustNot)

	tests := []struct {
		in  *Node
		out string
	}{
		{nil, ""},
		{NewNode(), ""},
		{h1, ""},
		{h2, `~[~"x", ~"y"]`},
		{h3, `~[~[~"x", ~"y"], ~[+"v", -"w"]]`},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.String()
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestNodeHeight(t *testing.T) {
	h1 := NewNode()
	h1.NewChild()

	h2 := NewNode()
	h2.NewChild().NewChild()
	h2.NewChild()

	h3 := NewNode()
	h3.NewChild()
	h3.NewChild().NewChild()
	h3.NewChild().NewChild().NewChild()

	tests := []struct {
		in  *Node
		out int
	}{
		{nil, 0},
		{new(Node), 0},
		{NewNode(), 0},
		{h1, 1},
		{h2, 2},
		{h3, 3},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.Height()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodeRoot(t *testing.T) {
	n := NewNode()
	m := n.NewChild()
	tests := []struct {
		in  *Node
		out *Node
	}{
		{nil, nil},
		{n, n},
		{m, n},
	}

	for i, tt := range tests {
		actual := tt.in.Root()
		msg := fmt.Sprintf("Test case (%d) fails", i)
		assert.Equal(t, tt.out, actual, msg)
	}

}
func TestNodeIsRoot(t *testing.T) {
	tests := []struct {
		in  *Node
		out bool
	}{
		{nil, false},
		{NewNode(), true},
		{NewNode().NewChild(), false},
	}

	for i, tt := range tests {
		actual := tt.in.IsRoot()
		msg := fmt.Sprintf("Test case (%d) fails", i)
		assert.Equal(t, tt.out, actual, msg)
	}
}
