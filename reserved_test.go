package gossip

import (
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestIsVerbType(t *testing.T) {
	assert.True(t, IsMust(Must))
	assert.False(t, IsMust(Should))
	assert.False(t, IsMust(0))

	assert.True(t, IsMustNot(MustNot))
	assert.False(t, IsMustNot(Should))
	assert.False(t, IsMustNot(0))

	assert.True(t, IsShould(Should))
	assert.False(t, IsShould(Must))
	assert.False(t, IsShould(0))
}

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
		{`0+"567+"`, Must, 1},
		{`0-"567+"`, MustNot, 1},
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
		{Must, false},
		{MustNot, false},
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
		{Must, Must, false},
		{MustNot, Must, false},
		{PhraseDelim, Must, true},
		{SubqueryStart, Must, true},
		{SubqueryEnd, Must, false},
		{Space, Must, true},
		{e, Must, true},
		{a, Must, false},
		// current = minus
		{Must, MustNot, false},
		{MustNot, MustNot, false},
		{PhraseDelim, MustNot, true},
		{SubqueryStart, MustNot, true},
		{SubqueryEnd, MustNot, false},
		{Space, MustNot, true},
		{e, MustNot, true},
		{a, MustNot, false},
		// current = quote
		{Must, PhraseDelim, true},
		{MustNot, PhraseDelim, true},
		{PhraseDelim, PhraseDelim, true},
		{SubqueryStart, PhraseDelim, true},
		{SubqueryEnd, PhraseDelim, false},
		{Space, PhraseDelim, true},
		{e, PhraseDelim, true},
		{a, PhraseDelim, false},
		// current = subquery start
		{Must, SubqueryStart, true},
		{MustNot, SubqueryStart, true},
		{PhraseDelim, SubqueryStart, true},
		{SubqueryStart, SubqueryStart, true},
		{SubqueryEnd, SubqueryStart, false},
		{Space, SubqueryStart, true},
		{e, SubqueryStart, true},
		{a, SubqueryStart, false},
		// current =Subuery end
		{Must, SubqueryEnd, false},
		{MustNot, SubqueryEnd, false},
		{PhraseDelim, SubqueryEnd, true},
		{SubqueryStart, SubqueryEnd, true},
		{SubqueryEnd, SubqueryEnd, true},
		{Space, SubqueryEnd, true},
		{e, SubqueryEnd, false},
		{a, SubqueryEnd, true},
		// current = space
		{Must, Space, false},
		{MustNot, Space, false},
		{PhraseDelim, Space, true},
		{SubqueryStart, Space, true},
		{SubqueryEnd, Space, true},
		{Space, Space, true},
		{e, Space, true},
		{a, Space, true},
		// non-reserved
		{Must, a, true},
		{MustNot, a, true},
		{PhraseDelim, a, true},
		{SubqueryStart, a, true},
		{SubqueryEnd, a, false},
		{Space, a, true},
		{e, a, true},
		{a, a, true},
		// Last
		{Must, e, false},
		{MustNot, e, false},
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
		{Must, Must, a, false},
		{MustNot, Must, a, false},
		{Should, Must, a, false},
		{SubqueryEnd, Must, a, false},
		{Space, Must, Space, false},
		{Space, Must, Comma, false},
		{Space, Must, e, false},
		{e, Must, e, false},
		{Space, Must, Must, false},
		{Space, Must, MustNot, false},
		{Space, Must, Should, false},
		{PhraseDelim, Must, a, true},
		{SubqueryStart, Must, a, true},
		{e, Must, a, true},
		{Space, Must, a, true},
		{Comma, Must, a, true},
		// Subquery starts
		{SubqueryEnd, SubqueryStart, a, false},
		{a, SubqueryStart, e, false},
		{e, SubqueryStart, e, false},
		{e, SubqueryStart, e, false},
		{Must, SubqueryStart, a, true},
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
		{`0123"567+"`, Must, -1},
		{`0123"`, PhraseDelim, 4},
		{`0123"567+"`, Must, -1},
		{`w "[]"`, SubqueryStart, -1},
		{"日本語", r1, 6}, // 10
		{`日本語"+`, Must, -1},
		{`+`, Must, 0},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %#v", i, tt)
		assert.Equal(t, tt.out, indexNonPhraseRune(tt.s, tt.r), msg)
	}

}

func TestVerbStringForTree(t *testing.T) {
	tests := []struct {
		in  rune
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
		in  rune
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
