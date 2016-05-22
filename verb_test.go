package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerbStringForTree(t *testing.T) {
	tests := []struct {
		in  int
		out string
	}{
		{Must, verbStrings[Must]},
		{Should, verbStrings[Should]},
		{MustNot, verbStrings[MustNot]},
		{VerbError, verbStrings[VerbError]},
		{999, ""},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, VerbString(tt.in), msg)
	}

}

func TestVerbStringForHumans(t *testing.T) {
	tests := []struct {
		in  int
		out string
	}{
		{Must, verbStringsForHumans[Must]},
		{Should, verbStringsForHumans[Should]},
		{MustNot, verbStringsForHumans[MustNot]},
		{VerbError, verbStringsForHumans[VerbError]},
		{999, ""},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, VerbStringHuman(tt.in), msg)
	}

}
