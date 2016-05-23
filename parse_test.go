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
		{"w", &Node{phrase: "w", verb: Should}},
		//
		{
			`+"machine learning"`,
			&Node{verb: Must, phrase: "machine learning"},
		},
		//
		{
			"x y",
			&Node{
				verb: Should,
				children: []*Node{
					{phrase: "x", verb: Should},
					{phrase: "y", verb: Should},
				},
			},
		},
		//
		{
			"x,y",
			&Node{
				verb: Should,
				children: []*Node{
					{phrase: "x", verb: Should},
					{phrase: "y", verb: Should},
				},
			},
		},
		//
		{
			"x,+y",
			&Node{
				verb: Should,
				children: []*Node{
					{phrase: "x", verb: Should},
					{phrase: "y", verb: Must},
				},
			},
		},
		//
		{
			"x, +y",
			&Node{
				verb: Should,
				children: []*Node{
					{phrase: "x", verb: Should},
					{phrase: "y", verb: Must},
				},
			},
		},
		//
		{
			"-x +y",
			&Node{
				verb: Should,
				children: []*Node{
					{phrase: "x", verb: MustNot},
					{phrase: "y", verb: Must},
				},
			},
		},
		//
		{
			"x +[+y -z]",
			&Node{
				verb: Should,
				children: []*Node{
					{phrase: "x", verb: Should},
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
		//
		{
			"x, +[+y, -z]",
			&Node{
				verb: Should,
				children: []*Node{
					{phrase: "x", verb: Should},
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
		//
		{
			`+"phrase one", [+"phrase 2", -z]`,
			&Node{
				verb: Should,
				children: []*Node{
					{
						verb:   Must,
						phrase: "phrase one"},
					// Subquery +[+y - z]
					{
						verb: Should,
						children: []*Node{
							{phrase: "phrase 2", verb: Must},
							{phrase: "z", verb: MustNot},
						},
					},
				},
			},
		},
		//
		{
			`+"phrase one" [+"phrase 2" -z]`,
			&Node{
				verb: Should,
				children: []*Node{
					{
						verb:   Must,
						phrase: "phrase one"},
					// Subquery +[+y - z]
					{
						verb: Should,
						children: []*Node{
							{phrase: "phrase 2", verb: Must},
							{phrase: "z", verb: MustNot},
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
			i, tt.in, tree, tree, tree.Root().Children(),
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
	// Output: ["data science", "machine learning", +"statistics", -"hype"]
}

func ExampleLeafNode() {
	search := `"data science" "machine learning" +statistics -hype`
	tree, _ := Parse(search)

	var children []*Node
	children = tree.Children()
	leaf := children[2]
	fmt.Printf("%s contain %s", leaf.VerbString(), leaf.Phrase())
	// Output: must contain statistics
}
