package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		`some words +[]`,     // last leaf is empty
		`"some words +[]" +`, // last leaf is empty
		`  `,                 // empty
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
