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
		{"[", SubqueryStart, 0},
		{"0]2", SubqueryEnd, 1},
		{`0"`, PhraseDelim, 1},
		{`0+"567+"`, rune(Must), 1},
		{`0-"567+"`, rune(Not), 1},
		{`0123 "`, Space, 4},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		r, i := NextReserved(tt.in)
		assert.Equal(t, tt.r, r, msg)
		assert.Equal(t, tt.i, i, msg)
	}

}

func TestCheckReservedOutOfBounds(t *testing.T) {
	assert.Equal(t, false, checkReserved("", 0, 1, 1))
	assert.Equal(t, false, checkReserved("", -1, 1, 1))

}

func TestCheckReservedNonUTF8(t *testing.T) {
	const bad = "\xbd\xb2\xbd"
	var r rune = 1
	assert.Equal(t, false, checkReserved(bad, r, 0, 1))
	assert.Equal(t, false, checkReserved(bad, r, 1, 1))
}

// Checks when the reserved rune is the only element of the search string.
func TestCheckReservedSingleton(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{rune(Must), false},
		{rune(Not), false},
		{PhraseDelim, false},
		{SubqueryStart, false},
		{SubqueryEnd, false},
		{Space, true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		// Effectively replace the _ with the current rune.
		actual := IsTripleValid(utf8.RuneError, tt.in, utf8.RuneError)
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestIsValidPair(t *testing.T) {
	a, _ := utf8.DecodeRuneInString("a")
	e := utf8.RuneError

	tests := []struct {
		prev rune
		curr rune
		out  bool
	}{
		// current = plus
		{rune(Must), rune(Must), false},
		{rune(Not), rune(Must), false},
		{PhraseDelim, rune(Must), true},
		{SubqueryStart, rune(Must), true},
		{SubqueryEnd, rune(Must), false},
		{Space, rune(Must), true},
		{e, rune(Must), true},
		{a, rune(Must), false},
		// current = minus
		{rune(Must), rune(Not), false},
		{rune(Not), rune(Not), false},
		{PhraseDelim, rune(Not), true},
		{SubqueryStart, rune(Not), true},
		{SubqueryEnd, rune(Not), false},
		{Space, rune(Not), true},
		{e, rune(Not), true},
		{a, rune(Not), false},
		// current = quote
		{rune(Must), PhraseDelim, true},
		{rune(Not), PhraseDelim, true},
		{PhraseDelim, PhraseDelim, true},
		{SubqueryStart, PhraseDelim, true},
		{SubqueryEnd, PhraseDelim, false},
		{Space, PhraseDelim, true},
		{e, PhraseDelim, true},
		{a, PhraseDelim, false},
		// current = subquery start
		{rune(Must), SubqueryStart, true},
		{rune(Not), SubqueryStart, true},
		{PhraseDelim, SubqueryStart, true},
		{SubqueryStart, SubqueryStart, true},
		{SubqueryEnd, SubqueryStart, false},
		{Space, SubqueryStart, true},
		{e, SubqueryStart, true},
		{a, SubqueryStart, false},
		// current =Subuery end
		{rune(Must), SubqueryEnd, false},
		{rune(Not), SubqueryEnd, false},
		{PhraseDelim, SubqueryEnd, true},
		{SubqueryStart, SubqueryEnd, true},
		{SubqueryEnd, SubqueryEnd, true},
		{Space, SubqueryEnd, true},
		{e, SubqueryEnd, false},
		{a, SubqueryEnd, true},
		// current = space
		{rune(Must), Space, false},
		{rune(Not), Space, false},
		{PhraseDelim, Space, true},
		{SubqueryStart, Space, true},
		{SubqueryEnd, Space, true},
		{Space, Space, true},
		{e, Space, true},
		{a, Space, true},
		// non-reserved
		{rune(Must), a, true},
		{rune(Not), a, true},
		{PhraseDelim, a, true},
		{SubqueryStart, a, true},
		{SubqueryEnd, a, false},
		{Space, a, true},
		{e, a, true},
		{a, a, true},
		// Last
		{rune(Must), e, false},
		{rune(Not), e, false},
		{PhraseDelim, e, true},
		{SubqueryStart, e, false},
		{SubqueryEnd, e, true},
		{Space, e, true},
		{e, e, false},
		{a, e, true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		// Effectively replace the first _ with the current rune.
		actual := IsPairValid(tt.prev, tt.curr)
		assert.Equal(t, tt.out, actual, msg)
	}
}

// We do not tests all possible valid triples here, only ones that we
func TestValidTriple(t *testing.T) {
	a, _ := utf8.DecodeRuneInString("a")
	e := utf8.RuneError
	tests := []struct {
		prev rune
		curr rune
		next rune
		out  bool
	}{
		// Modal verbs
		{rune(Must), rune(Must), a, false},
		{rune(Not), rune(Must), a, false},
		{rune(Should), rune(Must), a, false},
		{SubqueryEnd, rune(Must), a, false},
		{Space, rune(Must), Space, false},
		{Space, rune(Must), Comma, false},
		{Space, rune(Must), e, false},
		{e, rune(Must), e, false},
		{Space, rune(Must), rune(Must), false},
		{Space, rune(Must), rune(Not), false},
		{Space, rune(Must), rune(Should), false},
		{PhraseDelim, rune(Must), a, true},
		{SubqueryStart, rune(Must), a, true},
		{e, rune(Must), a, true},
		{Space, rune(Must), a, true},
		{Comma, rune(Must), a, true},
		// Subquery starts
		{SubqueryEnd, SubqueryStart, a, false},
		{a, SubqueryStart, e, false},
		{e, SubqueryStart, e, false},
		{e, SubqueryStart, e, false},
		{rune(Must), SubqueryStart, a, true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		// Effectively replace the first _ with the current rune.
		actual := IsTripleValid(tt.prev, tt.curr, tt.next)
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
		{"", PhraseDelim, -1},
		{"[", SubqueryEnd, -1},
		// {`"\"`, Escape, -1},
		// {`012\"`, Escape, 3}, // 5
		{`0123"567+"`, rune(Must), -1},
		{`0123"`, PhraseDelim, 4},
		{`0123"567+"`, rune(Must), -1},
		{`w "[]"`, SubqueryStart, -1},
		{"日本語", r1, 6}, // 10
		{`日本語"+`, rune(Must), -1},
		{`+`, rune(Must), 0},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		assert.Equal(t, tt.out, indexNonPhraseRune(tt.s, tt.r), msg)
	}

}
