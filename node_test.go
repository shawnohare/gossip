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
	assert.Equal(t, VerbError, n.Verb())
	assert.Len(t, n.Children(), 0)
	assert.Equal(t, "", n.Phrase())
	assert.Equal(t, 0, n.Depth())
	assert.Nil(t, n.Parent())
	assert.Nil(t, n.AddChild(&Node{}))
	assert.Nil(t, n.AddChild(nil))
	assert.Nil(t, n.NewChild())
}

func TestNodeIsLeafAfterAdds(t *testing.T) {
	tests := []struct {
		in  *Node
		out bool
	}{
		{nil, false},
		{NewNode(), true},
		{NewNode().NewChild().Parent(), false},
		{NewNode().AddChild(NewNode()).Parent(), false},
		// AddChild(nil) is nil, so its parent is also nil, hence not a leaf.
		{NewNode().AddChild(nil).Parent(), false},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) with children %#v fails", i, tt.in.Children())
		actual := tt.in.IsLeaf()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodeVerb(t *testing.T) {
	tests := []struct {
		in  *Node
		out rune
	}{
		{nil, VerbError},
		{&Node{}, VerbError},
		{&Node{verb: Must}, VerbError},
		{&Node{verb: Must, phrase: "a"}, Must},
		{&Node{verb: Should, phrase: "a"}, Should},
		{&Node{verb: MustNot, phrase: "a"}, MustNot},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.Verb()
		assert.Equal(t, tt.out, actual, msg)
	}
}

func TestNodeVerbString(t *testing.T) {
	tests := []struct {
		in *Node
	}{
		{nil},
		{&Node{}},
		{&Node{verb: Must}},
		{&Node{verb: Must, phrase: "a"}},
		{&Node{verb: Should, phrase: "a"}},
		{&Node{verb: MustNot, phrase: "a"}},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.VerbString()
		assert.Equal(t, VerbStringHuman(tt.in.Verb()), actual, msg)
	}
}

func TestNodeSetParent(t *testing.T) {
	var n *Node
	p := &Node{}
	n = n.SetParent(p)
	assert.Equal(t, p, n.parent)
}

func TestSetVerb(t *testing.T) {
	var n *Node
	n = n.SetVerb(Must)
	assert.Equal(t, Must, n.verb)
}

func TestSetVerbInvalidInput(t *testing.T) {
	var n *Node
	n = n.SetVerb(-10)
	assert.Equal(t, Should, n.verb)
}

func TestSetPhrase(t *testing.T) {
	var n *Node
	n = n.SetPhrase("0")
	assert.Equal(t, "0", n.phrase)
}

func TestNodeIsValid(t *testing.T) {
	tests := []struct {
		in  *Node
		out bool
	}{
		{nil, false},
		{&Node{}, false},
		{
			&Node{
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
				phrase: "this phrase invalidates the node",
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

func TestNodeEquals(t *testing.T) {
	tests := []struct {
		n0  *Node
		n1  *Node
		out bool
	}{
		{nil, nil, false},
		{nil, &Node{}, false},
		{nil, &Node{phrase: "x"}, false},
		{NewNode(), NewNode(), false},
		{NewNode(), &Node{phrase: "x", verb: Should}, false},
		{&Node{phrase: "x", verb: Should}, &Node{phrase: "y", verb: Should}, false},
		{&Node{verb: Must, phrase: "x"}, &Node{verb: Should, phrase: "x"}, false},
		{&Node{verb: Must, phrase: "x"}, &Node{verb: MustNot, phrase: "x"}, false},
		{&Node{phrase: "x", verb: Should}, &Node{phrase: "x", verb: Should}, true},
		// Basic test with children.
		{
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
			},
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
			},
			true,
		},
		// 10. Trees of height 2.
		{
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb: Should,
						children: []*Node{
							&Node{
								verb:   Must,
								phrase: "x1",
							},
						},
					},
				},
			},
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb: Should,
						children: []*Node{
							&Node{
								verb:   Must,
								phrase: "x1",
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
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
				phrase: "this phrase invalidates the node",
			},
			&Node{
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
			},
			false,
		},
		// Children are different lengths.
		{
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
			},
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
					&Node{
						verb:   Must,
						phrase: "x2",
					},
				},
			},
			false,
		},
		// Children are different.
		{
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x",
					},
				},
			},
			&Node{
				verb: Should,
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "y",
					},
				},
			},
			false,
		},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case %d", i)
		actual := tt.n0.Equals(tt.n1)
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
		{&Node{parent: &Node{}}, ""},
		{&Node{phrase: "x", verb: Should}, "x"},
		{&Node{children: []*Node{&Node{phrase: "x", verb: Should}}}, ""},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.Phrase()
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
		{&Node{parent: &Node{}}, 1},
		{&Node{parent: &Node{parent: &Node{}}}, 2},
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
		{h2, `["x", "y"]`},
		{h3, `[["x", "y"], [+"v", -"w"]]`},
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
