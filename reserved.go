package gossip

import "unicode/utf8"

// Define modal verbs.
const (
	VerbError int = -2
	Should    int = 0
	MustNot   int = int(minus)
	Must      int = int(plus)
)

// Define special characters used when parsing a raw search query.
const (
	space    rune = 0x00000020
	quote    rune = 0x00000022
	at       rune = 0x00000040
	plus     rune = 0x0000002b
	minus    rune = 0x0000002d
	lparen   rune = 0x00000028
	rparen   rune = 0x00000029
	lbracket rune = 0x0000005b
	rbracket rune = 0x0000005d
	escape   rune = 0x0000005c // reverse solidus, \
)

var empty struct{}

var reservedRuneLookup map[rune]struct{} = map[rune]struct{}{
	space:    empty,
	quote:    empty,
	plus:     empty,
	minus:    empty,
	lbracket: empty,
	rbracket: empty,
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

// lookBehindCheck inspects the provided string and compares the last
// rune to the input rune, to determine whether the pair is valid.
// The current rune is assumed to be outside of a phrase literal.
func lookBehindCheck(before string, current rune) bool {
	prev, w := utf8.DecodeLastRuneInString(before)
	// Matrix of acceptable combinations. Current control rune on top.
	// Previous rune on left.  Here _ is a space and r any non-reserved rune.
	//    +-  "  [  ]  _
	// +-  x  o  o  x  x
	//  "  x  ?  x  o  o
	//  [  o  o  o  ?  o
	//  ]  x  x  x  o  ?
	//  _  o  o  o  o  o
	//  r  x  x  x  o  o

	// Valid the current rune is the initial or preceeded by space.
	// This will make multiple spaces acceptable inputs.
	if w == 0 {
		return true
	}

	// TODO finish all cases.  Also, need to consider exactly when this func
	// is called. Presumably any time outside of a phrase literal.
	var ok bool
	switch current {
	case plus, minus:
		ok = prev == lbracket || prev == space
	case quote:
		ok = isRuneReserved(prev) && prev != rbracket
	case lbracket:
		ok = isRuneReserved(prev) && prev != quote && prev != rbracket
	case rbracket, space:
		ok = prev != plus && prev != minus
	}

	return ok

}

// indexNonPhraseRune reports the index of the next rune matching the input
// that is not contained in a phrase literal.  For instance,
// if s = `[word1 "phrase with brackets[]" word2]` and r = `]`, then
// the index returned is len(s)-1.
func indexNonPhraseRune(s string, r rune) int {
	inLiteral := false
	for i, ri := range s {
		// Only match if we are not in a phrase literal.
		if !inLiteral && ri == r {
			return i
		}
		// Toggle inLiteral when we see quotation marks.
		if ri == quote {
			inLiteral = !inLiteral
		}
	}

	return -1
}

// FIXME: remove?
// RemoveEscapes replaces occurrences of \R with R, where R is any rune.
// func RemoveEscapes(s string) string {
// 	var (
// 		s0   []rune
// 		prev rune
// 	)

// 	for _, r := range s {
// 		if r != escape || prev == escape {
// 			s0 = append(s0, r)
// 		}
// 		prev = r
// 	}

// 	return string(s0)
// }
