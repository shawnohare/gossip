package gossip

import (
	"errors"
	"strings"
	// "regexp"
	"unicode/utf8"
)

const (
	ErrorMalformedQuery    = "Search query is malformed: "
	ErrorUnpairedQuotation = ErrorMalformedQuery + "Unpaired quotation mark."
	ErrorUnpairedBracket   = ErrorMalformedQuery + "Unpaired bracket."
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

// Define modal verbs.
const (
	VerbError int = -2
	Should    int = 0
	MustNot   int = int(minus)
	Must      int = int(plus)
)

// lookBehindCheck inspects the provided string and compares the last
// rune to the input rune, to determine whether the pair is valid.
// The current rune is assumed to be marked as a control.  For instance,
// it is the the left quotation mark in a pair, the left bracket (or +, -)
// outside of a string literal, etc.
func lookBehindCheck(before string, current rune) bool {
	prev, w := utf8.DecodeLastRuneInString(before)

	// Valid the current rune is the initial
	if w == 0 {
		return true
	}
	if prev == current {
		return false
	}

	if prev == space ||
		prev == lbracket ||
		prev == quote ||
		prev == minus ||
		prev == plus {
		return true
	}
	return false

}

// indexUnescapedRune reports the index of the next rune matching the input
// that is not preceeded by an unescaped escape character.
//
// For instance, if r is the quotation mark rune, and
// s := `phrase with double escape \\"`, then the f(s, r) = len(s)-1,
// where f is this function.
func indexUnescapedRune(s string, r rune) int {
	for i, ri := range s {
		if ri == r {
			b, w := utf8.DecodeLastRuneInString(s[:i])
			if b != escape {
				return i
			}

			// If the escape is itself escaped, report the current match.
			j := i - w
			if j >= 0 {
				b, _ := utf8.DecodeLastRuneInString(s[:j])
				if b == escape {
					return i
				}
			}
		}
	}

	return -1
}

// RemoveEscapes replaces occurrences of \R with R, where R is any rune.
func RemoveEscapes(s string) string {
	var (
		s0   []rune
		prev rune
	)

	for _, r := range s {
		if r != escape || prev == escape {
			s0 = append(s0, r)
		}
		prev = r
	}

	return string(s0)
}

// Parse recursively converts a raw text search into a
// structured query tree.  The produced tree, represented by the returned
// query, is semantically isomorphic to input query but not identical.  In
// particular, non-branching paths are collapsed. For instance, the query
// `[[golang]]` is semantically identical to `golang`, and this function
// returns the single node tree for the later.
func Parse(s string) (*Tree, error) {
	var (
		currVerb int
		i        int             // current index in input string
		root     = &Node{src: s} // each element of the input is a child.
	)

	for i < len(s) {
		r, width := utf8.DecodeRuneInString(s[i:]) // Get next rune.

		switch {
		case r == quote:
			i += width
			// Create a leaf query consisting the substring between the matched
			// quotation and the next unescaped quotation mark.
			if i >= len(s)-1 {
				return nil, errors.New(ErrorUnpairedQuotation)
			}
			j := indexUnescapedRune(s[i:], quote)
			if j == -1 {
				return nil, errors.New(ErrorUnpairedQuotation)
			}
			j += i // point j to loc in s of matched quotation mark

			q := &Node{
				verb:   currVerb,
				phrase: RemoveEscapes(s[i:j]),
			}
			root.AddChild(q)

			// Update state.
			currVerb = 0
			i = j + width // advance head past matching quotation mark

		case r == plus || r == minus:
			// Update state.  If we already remember a verb, the query is malformed.
			if currVerb != 0 {
				return nil, errors.New(ErrorMalformedQuery + "Unexpected + or -.")
			}
			currVerb = int(r)
			i += width

		case r == lbracket:
			// Create a nested subquery by recursively calling function on the
			// the substring
			// TODO replace with custom index func.
			j := indexUnescapedRune(s[i:], lbracket)
			if j == -1 {
				return nil, errors.New(ErrorUnpairedBracket)
			}

			subTree, err := Parse(s[i:j])
			if err != nil {
				return nil, err
			}

			// TODO deal with incongruent verbs?
			sub := subTree.Root()
			if sub.Verb() != Should && sub.Verb() != currVerb {
				return nil, errors.New(ErrorMalformedQuery + "Unexpected combination of verbs.")
			}
			sub.verb = currVerb
			root.AddChild(sub)

			currVerb = 0
			i += j + width // advance past matching ], assumes width [ = width ]

		case r == space:
			i += width

		default:
			// No special marks detected.  Search for next space or end of query.
			j := strings.Index(s[i:], " ")
			if j == -1 {
				j = len(s)
			}
			q := &Node{verb: currVerb, phrase: s[i:j]}
			root.AddChild(q)
			i += j + 1 // space has width 1
		}
	}

	// TODO:
	// Collapse unnecessary hierarchy.
	if len(root.children) == 1 && root.children[0].IsLeaf() {
		root = root.children[0]
		root.parent = nil
		root.src = s
	}

	if !root.IsValid() {
		return nil, errors.New(ErrorMalformedQuery)
	}

	tree := &Tree{
		root: root,
		src:  s,
	}

	return tree, nil
}
