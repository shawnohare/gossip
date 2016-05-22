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
		{"[", LeftBracket, 0},
		{"0]2", RightBracket, 1},
		{`0"`, Quote, 1},
		{`0+"567+"`, Plus, 1},
		{`0-"567+"`, Minus, 1},
		{`0123 "`, Space, 4},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		r, i := nextReserved(tt.in)
		assert.Equal(t, tt.r, r, msg)
		assert.Equal(t, tt.i, i, msg)
	}

}

// Checks when the reserved rune is the last element of the search.
func TestCheckReservedRuneWhenLast(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{Plus, false},
		{Minus, false},
		{Quote, false},
		{LeftBracket, false},
		{RightBracket, true},
		{Space, true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		// Effectively replace the _ with the current rune.
		actual := checkReserved(" _", tt.in, 1, 1)
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestCheckReservedOutOfBounds(t *testing.T) {
	assert.Equal(t, false, checkReserved("", 0, 1, 1))
	assert.Equal(t, false, checkReserved("", -1, 1, 1))

}

// Checks when the reserved rune is the only element of the search string.
func TestCheckReservedSingleton(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{Plus, false},
		{Minus, false},
		{Quote, false},
		{LeftBracket, false},
		{RightBracket, false},
		{Space, true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		// Effectively replace the _ with the current rune.
		actual := checkReserved("_", tt.in, 0, 1)
		assert.Equal(t, tt.out, actual, msg)
	}

}

// These tests assume the string to parse looks like before + rune + *
// where * is a non-empty string.  In particular the rune to check is
// never last.
func TestCheckReservedLookBehindNotLast(t *testing.T) {
	tests := []struct {
		before  string
		current rune
		out     bool
	}{
		// current = plus
		{"+", Plus, false},
		{"-", Plus, false},
		{`"`, Plus, false},
		{"[", Plus, true},
		{"]", Plus, false},
		{" ", Plus, true},
		// current = minus
		{"+", Minus, false},
		{"-", Minus, false},
		{`"`, Minus, false},
		{"[", Minus, true},
		{"]", Minus, false},
		{" ", Minus, true},
		// current = quote
		{"+", Quote, true},
		{"-", Quote, true},
		{`"`, Quote, true},
		{"[", Quote, true},
		{"]", Quote, false},
		{" ", Quote, true},
		// current = left bracket
		{"+", LeftBracket, true},
		{"-", LeftBracket, true},
		{`"`, LeftBracket, false},
		{"[", LeftBracket, true},
		{"]", LeftBracket, false},
		{" ", LeftBracket, true},
		// current = rbracket
		{"+", RightBracket, false},
		{"-", RightBracket, false},
		{`"`, RightBracket, true},
		{"[", RightBracket, true},
		{"]", RightBracket, true},
		{" ", RightBracket, true},
		// current = space
		{"+", Space, false},
		{"-", Space, false},
		{`"`, Space, true},
		{"[", Space, true},
		{"]", Space, true},
		{" ", Space, true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		// Effectively replace the first _ with the current rune.
		actual := checkReserved(tt.before+"__", tt.current, 1, 1)
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
		{"", Quote, -1},
		{"[", RightBracket, -1},
		{`"\"`, Escape, -1},
		{`012\"`, Escape, 3}, // 5
		{`0123"567+"`, Plus, -1},
		{`0123"`, Quote, 4},
		{`0123"567+"`, Plus, -1},
		{`w "[]"`, LeftBracket, -1},
		{"日本語", r1, 6}, // 10
		{`日本語"+`, Plus, -1},
		{`+`, Plus, 0},
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
