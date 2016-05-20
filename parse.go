package gossip

import (
	"errors"
	"log"
	"strings"
	"unicode/utf8"
)

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
			// FIXME: remove comments
			// j := indexUnescapedRune(s[i:], quote)
			j := strings.Index(s[i:], `"`)
			// j := indexRune(s[i:], quote)
			if j == -1 {
				return nil, errors.New(ErrorUnpairedQuotation)
			}
			j += i // point j to loc in s of matched quotation mark

			// FIXME:
			log.Println("start:end of quotes", i-width, j)
			log.Println("Phrase for phrae query:", s[i:j])

			q := &Node{
				verb:   currVerb,
				phrase: s[i:j],
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
			j := indexNonPhraseRune(s[i:], lbracket)
			if j == -1 {
				return nil, errors.New(ErrorUnpairedBracket)
			}

			subTree, err := Parse(s[i : j+i])
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
			// No special marks detected.  Search for next space.
			r, j := nextReserved(s[i:])
			if j == -1 {
				j = len(s)
			} else {
				j += i
			}
			// Error on words that contain reserved characters such as `word1[word2`.
			if r != space {
				return nil, errors.New(ErrorUnexpectedReservedRune)
			}
			q := &Node{verb: currVerb, phrase: s[i:j]}
			root.AddChild(q)
			i = j + 1
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
