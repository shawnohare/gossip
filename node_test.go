package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case %d", i)
		actual := tt.n0.Equals(tt.n1)
		assert.Equal(t, tt.out, actual, msg)
	}
}
