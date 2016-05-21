package gossip

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestNextReservedRune(t *testing.T) {
	tests := []struct {
		in string
		r  rune
		i  int
	}{
		{"", utf8.RuneError, -1},
		{"anything", utf8.RuneError, -1},
		{"日本語", utf8.RuneError, -1},
		{"[", lbracket, 0},
		{"0]2", rbracket, 1},
		{`0"`, quote, 1},
		{`0+"567+"`, plus, 1},
		{`0-"567+"`, minus, 1},
		{`0123 "`, space, 4},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		r, i := nextReserved(tt.in)
		assert.Equal(t, tt.r, r, msg)
		assert.Equal(t, tt.i, i, msg)
	}

}

func TestLookBehindCheck(t *testing.T) {
	tests := []struct {
		before  string
		current rune
		out     bool
	}{
		// current = plus
		{"+", plus, false},
		{"-", plus, false},
		{`"`, plus, false},
		{"[", plus, true},
		{"]", plus, false},
		{" ", plus, true},
		{"+d", plus, false},
		{"", plus, true},
		// current = minus
		{"+", minus, false},
		{"-", minus, false},
		{`"`, minus, false},
		{"[", minus, true},
		{"]", minus, false},
		{" ", minus, true},
		{"+d", minus, false},
		{"", minus, true},
		// current = quote
		{"+", quote, true},
		{"-", quote, true},
		{`"`, quote, true},
		{"[", quote, true},
		{"]", quote, false},
		{" ", quote, true},
		{"+d", quote, false},
		{"", quote, true},
		// current = left bracket
		{"+", lbracket, true},
		{"-", lbracket, true},
		{`"`, lbracket, false},
		{"[", lbracket, true},
		{"]", lbracket, false},
		{" ", lbracket, true},
		{"+d", lbracket, false},
		{"", lbracket, true},
		// current = rbracket
		{"+", rbracket, false},
		{"-", rbracket, false},
		{`"`, rbracket, true},
		{"[", rbracket, true},
		{"]", rbracket, true},
		{" ", rbracket, true},
		{"+d", rbracket, true},
		{"", rbracket, true},
		// current = space
		{"+", space, false},
		{"-", space, false},
		{`"`, space, true},
		{"[", space, true},
		{"]", space, true},
		{" ", space, true},
		{"+d", space, true},
		{"", space, true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		actual := lookBehindCheck(tt.before, tt.current)
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestIndexNonPhraseRune(t *testing.T) {
	r0, _ := utf8.DecodeRuneInString("")
	r1, _ := utf8.DecodeRuneInString("語")
	tests := []struct {
		s   string
		r   rune
		out int
	}{
		{"", 0, -1},
		{"anything", r0, -1},
		{"", quote, -1},
		{"[", rbracket, -1},
		{`"\"`, escape, -1},
		{`012\"`, escape, 3}, // 5
		{`0123"567+"`, plus, -1},
		{`0123"`, quote, 4},
		{`0123"567+"`, plus, -1},
		{`w "[]"`, lbracket, -1},
		{"日本語", r1, 6}, // 10
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		assert.Equal(t, tt.out, indexNonPhraseRune(tt.s, tt.r), msg)
	}

}

// FIXME: remove
// func TestRemoveEscapes(t *testing.T) {
// 	tests := []struct {
// 		in  string
// 		out string
// 	}{
// 		{"", ""},
// 		{"a phrase", "a phrase"},
// 		{`\`, ``},
// 		{`\\`, `\`},
// 		{`\"a phrase\"`, `"a phrase"`},
// 	}

// 	for i, tt := range tests {
// 		msg := fmt.Sprintf("Fails test case %d", i)
// 		assert.Equal(t, tt.out, RemoveEscapes(tt.in), msg)
// 	}
// }
