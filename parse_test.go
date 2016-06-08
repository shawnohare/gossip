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
		out *Node
	}{
		{"w", &Node{Phrase: "w", Verb: Should}},
		//
		{
			`+"machine learning"`,
			&Node{Verb: Must, Phrase: "machine learning"},
		},
		//
		{
			"x y",
			&Node{
				Verb: Should,
				Children: []*Node{
					{Phrase: "x", Verb: Should},
					{Phrase: "y", Verb: Should},
				},
			},
		},
		//
		{
			"x,y",
			&Node{
				Verb: Should,
				Children: []*Node{
					{Phrase: "x", Verb: Should},
					{Phrase: "y", Verb: Should},
				},
			},
		},
		//
		{
			"x,+y",
			&Node{
				Verb: Should,
				Children: []*Node{
					{Phrase: "x", Verb: Should},
					{Phrase: "y", Verb: Must},
				},
			},
		},
		//
		{
			"x, +y",
			&Node{
				Verb: Should,
				Children: []*Node{
					{Phrase: "x", Verb: Should},
					{Phrase: "y", Verb: Must},
				},
			},
		},
		//
		{
			"-x +y",
			&Node{
				Verb: Should,
				Children: []*Node{
					{Phrase: "x", Verb: MustNot},
					{Phrase: "y", Verb: Must},
				},
			},
		},
		//
		{
			"x +[+y -z]",
			&Node{
				Verb: Should,
				Children: []*Node{
					{Phrase: "x", Verb: Should},
					// Subquery +[+y - z]
					{
						Verb: Must,
						Children: []*Node{
							{Phrase: "y", Verb: Must},
							{Phrase: "z", Verb: MustNot},
						},
					},
				},
			},
		},
		//
		{
			"x, +[+y, -z]",
			&Node{
				Verb: Should,
				Children: []*Node{
					{Phrase: "x", Verb: Should},
					// Subquery +[+y - z]
					{
						Verb: Must,
						Children: []*Node{
							{Phrase: "y", Verb: Must},
							{Phrase: "z", Verb: MustNot},
						},
					},
				},
			},
		},
		//
		{
			`+"phrase one", [+"phrase 2", -z]`,
			&Node{
				Verb: Should,
				Children: []*Node{
					{
						Verb:   Must,
						Phrase: "phrase one"},
					// Subquery +[+y - z]
					{
						Verb: Should,
						Children: []*Node{
							{Phrase: "phrase 2", Verb: Must},
							{Phrase: "z", Verb: MustNot},
						},
					},
				},
			},
		},
		//
		{
			`+"phrase one" [+"phrase 2" -z]`,
			&Node{
				Verb: Should,
				Children: []*Node{
					{
						Verb:   Must,
						Phrase: "phrase one"},
					// Subquery +[+y - z]
					{
						Verb: Should,
						Children: []*Node{
							{Phrase: "phrase 2", Verb: Must},
							{Phrase: "z", Verb: MustNot},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		tree, err := Parse(tt.in)
		msg := fmt.Sprintf(
			"Fails test case (%d)\ninput: %s\ntree: %#v\ntree.String(): %s\nchildren %#v",
			i, tt.in, tree, tree, tree.Root().GetChildren(),
		)
		assert.NoError(t, err, msg)
		assert.True(t, tt.out.Tree().Equals(tree.Tree()), msg)
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
		`,`,                  // empty
		`,,,`,                // empty
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

func ExampleParse() {
	// Example search that should include the phrases
	// data science, machine learning and math,
	// but must require the word statistics and not the word hype.
	search := `"data science" "machine learning" +statistics -hype`
	tree, _ := Parse(search)
	fmt.Println(tree)
	// Output: ~[~"data science", ~"machine learning", +"statistics", -"hype"]
}

func ExampleLeafNode() {
	search := `"data science" "machine learning" +statistics -hype`
	node, _ := Parse(search)

	var children []*Node
	children = node.GetChildren()
	leaf := children[2]
	fmt.Printf("%s contain %s", leaf.Verb.Pretty(), leaf.Phrase)
	// Output: must contain statistics
}
