package gossip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeJSON(t *testing.T) {
	tests := []struct {
		in  *Node
		out *JSON
	}{
		// Fails
		{new(Node), nil},
		{
			&Node{verb: Must, phrase: "test phrase"},
			&JSON{Verb: VerbStringHuman(Must), Phrase: "test phrase"},
		},
		{
			&Node{
				verb: Should,
				children: []*Node{
					&Node{verb: MustNot, phrase: "c0 phrase"},
					&Node{verb: Should, phrase: "c1 phrase"},
				},
			},
			&JSON{
				Verb: VerbStringHuman(Should),
				Children: []*JSON{
					&JSON{Verb: VerbStringHuman(MustNot), Phrase: "c0 phrase"},
					&JSON{Verb: VerbStringHuman(Should), Phrase: "c1 phrase"},
				},
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.out, tt.in.JSON())
	}
}
