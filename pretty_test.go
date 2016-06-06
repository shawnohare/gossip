package gossip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodePretty(t *testing.T) {
	tests := []struct {
		in  *Node
		out *Pretty
	}{
		// Fails
		{new(Node), nil},
		{
			&Node{verb: Must, phrase: "test phrase"},
			&Pretty{Verb: VerbStringHuman(Must), Phrase: "test phrase"},
		},
		{
			&Node{
				verb: Should,
				children: []*Node{
					&Node{verb: MustNot, phrase: "c0 phrase"},
					&Node{verb: Should, phrase: "c1 phrase"},
				},
			},
			&Pretty{
				Verb: VerbStringHuman(Should),
				Children: []*Pretty{
					&Pretty{Verb: VerbStringHuman(MustNot), Phrase: "c0 phrase"},
					&Pretty{Verb: VerbStringHuman(Should), Phrase: "c1 phrase"},
				},
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.out, tt.in.Pretty())
	}
}
