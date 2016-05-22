package gossip

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Parse converts a raw text search into a structured query term tree.
// The produced tree is semantically somorphic to input query, but not identical.
// It is a poset morphism, in that the order that phrase leaves have the
// same total ordering as induced by the words in the raw text.
// In particular, non-branching paths are collapsed. For instance, the query
// `[[golang]]` is semantically identical to `golang`, and this function
// returns the height 0 tree for the later.
//
// Semantically empty search phrases will yield a parse error.
func Parse(s string) (*Tree, error) {
	var (
		currVerb rune  = Should // modal verb to apply to children
		i        int            // current index in input string
		tree     *Tree = NewTree()
		curr     *Node = tree.Root()
	)

	if s == "" {
		return nil, errors.New(ErrorEmptyQuery)
	}

	for i < len(s) {
		r, width := utf8.DecodeRuneInString(s[i:]) // Get next rune.

		switch {
		// When we see a quotation mark, search for the next occurrence and
		// create a child
		case IsPhraseDelim(r):
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorMalformedQuery)
			}

			// Create a leaf query consisting the substring between the matched
			// quotation and the next unescaped quotation mark.
			i += width
			j := strings.IndexRune(s[i:], Quote)
			if j == -1 {
				return nil, errors.New(ErrorUnpairedQuotation)
			}
			j += i // point j to loc in s of matched quotation mark

			q := &Node{
				verb:   currVerb,
				phrase: s[i:j],
			}
			if !q.IsValid() {
				return nil, errors.New(ErrorEmptyQuery)
			}
			curr.AddChild(q)

			// Update state.
			currVerb = Should
			i = j + width // advance head past matching quotation mark

		case IsVerb(r):
			// Update state.  If we already remember a verb, the query is malformed.
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorVerbSequence)
			}
			currVerb = r
			i += width

		// Replace the current node with a new child subquery node.
		case IsSubqueryStart(r):
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorMalformedQuery)
			}
			curr = curr.AddChild(&Node{verb: currVerb})
			i += width
			currVerb = Should

		case IsSubqueryEnd(r):
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorMalformedQuery)
			}
			if !curr.IsValid() {
				return nil, errors.New(ErrorMalformedQuery)
			}
			curr = curr.Parent()
			i += width

		case IsSeparator(r):
			// Bad separators are currently detected by other tests.
			// if !checkReserved(s, r, i, width) {
			// 	return nil, errors.New(ErrorMalformedQuery)
			// }
			i += width

		default:
			// No reserved runes detected.  Search for next space.
			_, j := NextReserved(s[i:])
			if j == -1 {
				j = len(s)
			} else {
				j += i
			}

			// This will add the node with phrase xyz for the bad
			// query xyz+, but the error will be caught in the next check.
			_ = curr.AddChild(&Node{verb: currVerb, phrase: s[i:j]})

			i = j
			currVerb = Should
		}
	}

	// Collapse unnecessary hierarchy, and do a basic sanity check.
	root := tree.Root()
	if len(root.children) == 1 && root.children[0].IsLeaf() {
		root = root.children[0]
		root.parent = nil
	}

	// Node checks are cheap.  Catches queries like "  ".
	if !root.IsValid() {
		return nil, errors.New(ErrorMalformedQuery)
	}
	tree.root = root

	return tree, nil
}
