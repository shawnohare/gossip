package gossip

import (
	"strings"
	"unicode/utf8"
)

// Runes that have some special meaning.  Currently, Plus (+), Minus (i)
const (
	Space        rune = 0x00000020 // reserved
	Quote        rune = 0x00000022 // reserved
	Plus         rune = 0x0000002b // reserved
	Minus        rune = 0x0000002d // reserved
	LeftBracket  rune = 0x0000005b // reserved
	RightBracket rune = 0x0000005d // reserved
	LeftParen    rune = 0x00000028
	RightParen   rune = 0x00000029
	At           rune = 0x00000040
	Escape       rune = 0x0000005c // reverse solidus, \
)

var empty struct{}

var reservedRuneLookup map[rune]struct{} = map[rune]struct{}{
	Space:        empty,
	Quote:        empty,
	Plus:         empty,
	Minus:        empty,
	LeftBracket:  empty,
	RightBracket: empty,
}

func isRuneReserved(r rune) bool {
	_, ok := reservedRuneLookup[r]
	return ok
}

// nextReserved reports the next reserved rune and its index.
// If no reserved runes are found, the returned rune is a utf8.RuneError
// with index -1.
func nextReserved(s string) (rune, int) {
	for i, r := range s {
		if isRuneReserved(r) {
			return r, i
		}
	}
	return utf8.RuneError, -1
}

// checkReserved determiens if a reserved rune is in a valid sequence.
// Matrix of acceptable (prev, curr) rune pairs. Current control rune on top.
// Previous rune on left.  Here _ is a space and r any non-reserved rune.
//    +-  "  [  ]  _
// +-  x  o  o  x  x
//  "  x  ?  x  o  o
//  [  o  o  o  ?  o
//  ]  x  x  x  o  ?
//  _  o  o  o  o  o
//  r  x  x  x  o  o
//
// Most reserved runes can be the initial but not terminal rune in the string.
func checkReserved(s string, r rune, loc int, width int) bool {
	if loc < 0 || len(s) <= loc {
		return false
	}

	prev, w := utf8.DecodeLastRuneInString(s[:loc])
	first := w == 0
	last := loc+width >= len(s)

	var ok bool
	switch r {
	case Plus, Minus:
		ok = !last && (first || prev == LeftBracket || prev == Space)
	case Quote:
		ok = !last && (first || (isRuneReserved(prev) && prev != RightBracket))
	case LeftBracket:
		ok = !last && (first || (isRuneReserved(prev) && prev != Quote && prev != RightBracket))
	case RightBracket:
		ok = !first && prev != utf8.RuneError && prev != Plus && prev != Minus
	case Space:
		ok = prev != Plus && prev != Minus
	}

	return ok

}

// that is not contained in a phrase literal.  For instance,
// if s = `[word1 "phrase with brackets[]" word2]` and r = `]`, then
// the index returned is len(s)-1.
func indexNonPhraseRune(s string, r rune) int {
	inLiteral := false
	var i int
	for i < len(s) {
		ri, w := utf8.DecodeRuneInString(s[i:]) // Get next rune
		if ri == r {
			return i
		}

		i += w

		// If we see a quotation, find the matching mark and resume the search.
		if ri == Quote && i < len(s) {
			j := strings.IndexRune(s[i:], Quote)
			if j == -1 {
				return -1
			}
			i += j
		}

	}
	for i, ri := range s {
		// Only match if we are not in a phrase literal.
		if !inLiteral && ri == r {
			return i
		}
		// Toggle inLiteral when we see quotation marks.
		if ri == Quote {
			inLiteral = !inLiteral
		}
	}

	return -1
}

// Index reports the index of the next control rune or -1 if the rune is not
// present in the input string.  Runes inside phrase literals are ignored.
// ocurrence of r outside of a phrase literal.  For example,
//   Index(`"+" +`, Plus) = 4.
// func Index(s string, r rune) int {
// 	// TODO
// }
