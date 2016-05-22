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
		out int
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

func TestNodeSetParent(t *testing.T) {
	var n *Node
	p := &Node{}
	n = n.SetParent(p)
	assert.Equal(t, p, n.parent)
}

func TestSetVerb(t *testing.T) {
	var n *Node
	n = n.SetVerb(-10)
	assert.Equal(t, -10, n.verb)
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
		{&Node{}, &Node{}, false},
		{&Node{}, &Node{phrase: "x"}, false},
		{&Node{phrase: "x"}, &Node{phrase: "y"}, false},
		{&Node{verb: Must, phrase: "x"}, &Node{phrase: "x"}, false},
		{&Node{verb: Must, phrase: "x"}, &Node{verb: MustNot, phrase: "x"}, false},
		{&Node{phrase: "x"}, &Node{phrase: "x"}, true},
		// Basic test with children.
		{
			&Node{
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
			},
			&Node{
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
				children: []*Node{
					&Node{
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
				children: []*Node{
					&Node{
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
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x1",
					},
				},
			},
			&Node{
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
				children: []*Node{
					&Node{
						verb:   Must,
						phrase: "x",
					},
				},
			},
			&Node{
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
		{&Node{phrase: "x"}, "x"},
		{&Node{children: []*Node{&Node{phrase: "x"}}}, ""},
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
