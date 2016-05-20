package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseErrors(t *testing.T) {
	// All the tests should raise a parse error.
	tests := []string{
		"",
		`\ "no closing quotation`,
		`+"`,
		`[`,
		`+`,
		`+-word`,
		`++[word]`,
		`-`,
		`[]`,
		`+[]`,
		`-[]`,
		`-[+]`,
	}

	for i, tt := range tests {
		tree, err := Parse(tt)
		msg := fmt.Sprintf("Fails test case (%d) %s: tree = %#v", i, tt, tree)

		assert.Error(t, err, msg)
	}
}
