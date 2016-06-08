package gossip

import (
	"strings"
	"unicode/utf8"
)

// Reserved runes that have special meaning in search queries.
const (
	Space        rune = 0x00000020
	Quote        rune = 0x00000022
	Plus         rune = 0x0000002b
	Comma        rune = 0x0000002c
	Minus        rune = 0x0000002d
	LeftBracket  rune = 0x0000005b
	RightBracket rune = 0x0000005d
	Tilde        rune = 0x0000007e
	// LeftParen    rune = 0x00000028
	// RightParen   rune = 0x00000029
	// At           rune = 0x00000040
	// Escape       rune = 0x0000005c // reverse solidus, \
)

// Reserved rune aliases.
const (
	SubqueryStart rune = LeftBracket
	SubqueryEnd   rune = RightBracket
	PhraseDelim   rune = Quote
)

var reservedRuneLookup map[rune]struct{} = map[rune]struct{}{
	Space:        struct{}{},
	Comma:        struct{}{},
	Quote:        struct{}{},
	Plus:         struct{}{},
	Minus:        struct{}{},
	Tilde:        struct{}{},
	LeftBracket:  struct{}{},
	RightBracket: struct{}{},
}

func IsReserved(r rune) bool {
	_, ok := reservedRuneLookup[r]
	return ok
}

// IsSeparator states if the input denotes a query object separator.
// An object can be words, phrases, or subqueries.
// Currently, a " " and "," are considered equivalent separators.
func IsSeparator(r rune) bool {
	return r == Space || r == Comma
}

// IsSubqueryStarts states if the input denotes the start of a nested subquery.
func IsSubqueryStart(r rune) bool {
	return r == SubqueryStart
}

// IsSubqueryEnd states if the input denotes the end of a nested subquery.
func IsSubqueryEnd(r rune) bool {
	return r == SubqueryEnd
}

// IsPhraseDelim states if the input indicates the start of a phrase literal.
func IsPhraseDelim(r rune) bool {
	return r == Quote
}

// nextReserved reports the next reserved rune and its index.
// If no reserved runes are found, the returned rune is a utf8.RuneError
// with index -1.
func NextReserved(s string) (rune, int) {
	for i, r := range s {
		if IsReserved(r) {
			return r, i
		}
	}
	return utf8.RuneError, -1
}

// IsValidPair states whether the ordered pair of reserved runes is a valid combination.
// This is useful for basic look behind checks.
// If the prev rune is utf8.RuneError, then the current rune is assumed
// to be the initial rune in a string.
//
// Valid combinations are:
//    +-  "  [  ]  _,
// +-  x  o  o  x  x
//  "  x  ?  x  o  o
//  [  o  o  o  o  o
//  ]  x  x  x  o  o
// _,  o  o  o  o  o
//  r  x  x  x  o  o
func IsPairValid(prev rune, curr rune) bool {
	if prev == curr && prev == utf8.RuneError {
		return false
	}

	// Only first or last should be set at a single time.
	var (
		first bool
		last  bool
		p     rune = prev
		c     rune = curr
	)

	if prev == utf8.RuneError {
		first = true
	}
	if curr == utf8.RuneError {
		last = true
		c = prev
	}
	lit := IsPhraseDelim(prev)

	var ok bool
	switch {
	case IsRuneVerb(c):
		// Fail if last or second condition not met.
		ok = !last && (first || lit || IsSubqueryStart(p) || IsSeparator(p))

	case IsPhraseDelim(c):
		ok = last || first || (IsReserved(p) && !IsSubqueryEnd(p))

	case IsSubqueryStart(c):
		// Fail if last or second condition not met.
		ok = !last && (first || lit || IsReserved(p) && !IsPhraseDelim(p) && !IsSubqueryEnd(p))

	case IsSubqueryEnd(c):
		// Fail if first or previous is a verb.
		ok = last || (!first && !IsRuneVerb(p))

	case IsSeparator(c):
		// Previous cannot be a verb.
		ok = !IsRuneVerb(p)

	case !IsReserved(c):
		// Previous cannot be a subquery.
		ok = first || !IsSubqueryEnd(p)
	}
	return ok
}

// IsValidTriple states whether the ordered triple of runes represents
// a valid sequence in a search query.
// If the prev rune is utf8.RuneError, then the current rune is assumed
// to be initial.  If the next rune is a utf8.RuneError, then the current
// rune is considered to be terminal.  When both the previous and next
// runes are utf8.RuneError, then the current rune is represents a
// singleton query.
//
// The Validity of a triple depends primarily on the first two elements.
func IsTripleValid(prev rune, curr rune, next rune) bool {
	if prev == next && prev == utf8.RuneError {
		if curr == PhraseDelim {
			return false
		}

	}
	return IsPairValid(prev, curr) && IsPairValid(curr, next)
}

// checkReserved determiens if a reserved rune is in a valid sequence.
// Matrix of acceptable (prev, curr) rune pairs. Current control rune on top.
// Previous rune on left.  Here _ is a space and r any non-reserved rune.
//    +-  "  [  ]  _,
// +-  x  o  o  x  x
//  "  x  ?  x  o  o
//  [  o  o  o  ?  o
//  ]  x  x  x  o  ?
// _,  o  o  o  o  o
//  r  x  x  x  o  o
//
// Most reserved runes can be the initial but not terminal rune in the string.
func checkReserved(s string, r rune, loc int, width int) bool {
	var (
		prev      rune = utf8.RuneError
		next      rune = utf8.RuneError
		nextIndex int
	)

	// loc or width is out of bounds.
	if loc < 0 || len(s) <= loc || width < 0 {
		return false
	}

	// Decode previous rune or remember that it does not exist.
	prev, w := utf8.DecodeLastRuneInString(s[:loc])
	if prev == utf8.RuneError && w == 1 {
		return false
	}
	nextIndex = loc + width

	// Decode next rune, or remember that it does not exist.
	if nextIndex < len(s) {
		tmp, w := utf8.DecodeRuneInString(s[nextIndex:])
		if tmp == utf8.RuneError && w == 1 {
			return false
		}
		next = tmp
	}

	return IsTripleValid(prev, r, next)
}

// that is not contained in a phrase literal.  For instance,
// if s = `[word1 "phrase with brackets[]" word2]` and r = `]`, then
// the index returned is len(s)-1.
func indexNonPhraseRune(s string, r rune) int {
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
	return -1
}
