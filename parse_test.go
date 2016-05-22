package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePasses(t *testing.T) {
	// All the tests should pass.
	tests := []struct {
		in  string
		out *Tree
	}{
		{"w", &Tree{&Node{phrase: "w"}}},
		{
			`+"machine learning"`,
			&Tree{
				root: &Node{verb: Must, phrase: "machine learning"},
			},
		},
		// 1
		{
			"x y",
			&Tree{
				root: &Node{
					children: []*Node{
						{phrase: "x"},
						{phrase: "y"},
					},
				},
			},
		},
		// 2
		{
			"-x +y",
			&Tree{
				root: &Node{
					children: []*Node{
						{phrase: "x", verb: MustNot},
						{phrase: "y", verb: Must},
					},
				},
			},
		},
		// 3
		{
			"x +[+y -z]",
			&Tree{
				root: &Node{
					children: []*Node{
						{phrase: "x"},
						// Subquery +[+y - z]
						{
							verb: Must,
							children: []*Node{
								{phrase: "y", verb: Must},
								{phrase: "z", verb: MustNot},
							},
						},
					},
				},
			},
		},
		{
			`+"phrase one" [+"phrase 2" -z]`,
			&Tree{
				root: &Node{
					children: []*Node{
						{
							verb:   Must,
							phrase: "phrase one"},
						// Subquery +[+y - z]
						{
							children: []*Node{
								{phrase: "phrase 2", verb: Must},
								{phrase: "z", verb: MustNot},
							},
						},
					},
				},
			},
		},
		{
			`+"phrase one" [+"phrase 2" -z]`,
			&Tree{
				root: &Node{
					children: []*Node{
						{
							verb:   Must,
							phrase: "phrase one"},
						// Subquery +[+y - z]
						{
							children: []*Node{
								{phrase: "phrase 2", verb: Must},
								{phrase: "z", verb: MustNot},
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		tree, err := Parse(tt.in)
		msg := fmt.Sprintf(
			"Fails test case (%d)\ninput: %s\ntree: %#v\ntree.String(): %s",
			i, tt.in, tree, tree.String(),
		)
		assert.NoError(t, err, msg)
		assert.True(t, tt.out.Equals(tree), msg)
	}
}

func TestParseFailures(t *testing.T) {
	// All the tests should raise a parse error.
	tests := []string{
		"", // empty
		`\ "no closing quotation`,
		`+`,
		`-`,
		`]`,
		`[`,
		`"`,
		`++`,
		`+-`,
		`[+-]`,
		`[[]]`,
		`[`,
		`[]`,
		`+[]`,
		`-[]`,
		`-[+]`,
		`+-word`,
		`++[word]`,
		`some words +[`,
		`some words with empty phrase ""`,
		`some words +[]`,     // last leaf is empty
		`"some words +[]" +`, // last leaf is empty
		`  `,                 // empty
		` + ]`,               // empty
		`  ]]`,               // empty
		`+w + `,              // empty
		`+w+ `,               // empty
	}

	for i, tt := range tests {
		tree, err := Parse(tt)
		msg := fmt.Sprintf(
			"Fails test case (%d)\ninput: %s\ntree: %#v\ntree.String(): %s",
			i, tt, tree, tree.String(),
		)
		assert.Error(t, err, msg)
	}
}
