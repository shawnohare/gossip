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
		currVerb int   // modal verb to apply to children
		i        int   // current index in input string
		tree     *Tree = NewTree()
		curr     *Node = tree.Root()
	)

	for i < len(s) {
		r, width := utf8.DecodeRuneInString(s[i:]) // Get next rune.

		switch {
		// When we see a quotation mark, search for the next occurrence and
		// create a child
		case r == Quote:
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
			curr.AddChild(q)

			// Update state.
			currVerb = 0
			i = j + width // advance head past matching quotation mark

		case r == Plus || r == Minus:
			// Update state.  If we already remember a verb, the query is malformed.
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorMalformedQuery)
			}
			currVerb = int(r)
			i += width

		// Replace the current node with a new child subquery node.
		case r == LeftBracket:
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorMalformedQuery)
			}
			curr = curr.AddChild(&Node{verb: currVerb})
			i += width
			currVerb = 0

		case r == RightBracket:
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorMalformedQuery)
			}
			curr = curr.Parent()
			if curr == nil {
				return nil, errors.New(ErrorUnpairedBracket)
			}

		case r == Space:
			if !checkReserved(s, r, i, width) {
				return nil, errors.New(ErrorMalformedQuery)
			}
			i += width

		default:
			// No reserved runes detected.  Search for next space.
			_, j := nextReserved(s[i:])
			if j == -1 {
				j = len(s)
			} else {
				j += i
			}
			_ = curr.AddChild(&Node{verb: currVerb, phrase: s[i:j]})
			i = j
		}
	}

	// Collapse unnecessary hierarchy, and do a basic sanity check.
	root := tree.Root()
	if len(root.children) == 1 && root.children[0].IsLeaf() {
		root = root.children[0]
		root.parent = nil
	}
	if !root.IsValid() {
		return nil, errors.New(ErrorMalformedQuery)
	}
	tree.root = root

	return tree, nil
}

// FIXME: delete, this is defunct.
// func Parse(s string) (*Tree, error) {
// 	var (
// 		currVerb int
// 		i        int             // current index in input string
// 		root     = &Node{src: s} // each element of the input is a child.
// 	)

// 	for i < len(s) {
// 		r, width := utf8.DecodeRuneInString(s[i:]) // Get next rune.

// 		switch {
// 		case r == quote:
// 			i += width
// 			// Create a leaf query consisting the substring between the matched
// 			// quotation and the next unescaped quotation mark.
// 			if i >= len(s)-1 {
// 				return nil, errors.New(ErrorUnpairedQuotation)
// 			}
// 			// FIXME: remove comments
// 			// j := indexUnescapedRune(s[i:], quote)
// 			j := strings.Index(s[i:], `"`)
// 			// j := indexRune(s[i:], quote)
// 			if j == -1 {
// 				return nil, errors.New(ErrorUnpairedQuotation)
// 			}
// 			j += i // point j to loc in s of matched quotation mark

// 			// FIXME:
// 			log.Println("start:end of quotes", i-width, j)
// 			log.Println("Phrase for phrae query:", s[i:j])

// 			q := &Node{
// 				verb:   currVerb,
// 				phrase: s[i:j],
// 			}
// 			root.AddChild(q)

// 			// Update state.
// 			currVerb = 0
// 			i = j + width // advance head past matching quotation mark

// 		case r == plus || r == minus:
// 			// Update state.  If we already remember a verb, the query is malformed.
// 			if currVerb != 0 {
// 				return nil, errors.New(ErrorMalformedQuery + "Unexpected + or -.")
// 			}
// 			currVerb = int(r)
// 			i += width

// 		case r == lbracket:
// 			// Create a nested subquery by recursively calling function on the
// 			// the substring
// 			// TODO replace with custom index func.
// 			j := indexNonPhraseRune(s[i:], lbracket)
// 			if j == -1 {
// 				return nil, errors.New(ErrorUnpairedBracket)
// 			}

// 			subTree, err := Parse(s[i : j+i])
// 			if err != nil {
// 				return nil, err
// 			}

// 			// TODO deal with incongruent verbs?
// 			sub := subTree.Root()
// 			if sub.Verb() != Should && sub.Verb() != currVerb {
// 				return nil, errors.New(ErrorMalformedQuery + "Unexpected combination of verbs.")
// 			}
// 			sub.verb = currVerb
// 			root.AddChild(sub)

// 			currVerb = 0
// 			i += j + width // advance past matching ], assumes width [ = width ]

// 		case r == space:
// 			i += width

// 		default:
// 			// No special marks detected.  Search for next space.
// 			r, j := nextReserved(s[i:])
// 			if j == -1 {
// 				j = len(s)
// 			} else {
// 				j += i
// 			}
// 			// Error on words that contain reserved characters such as `word1[word2`.
// 			if r != space {
// 				return nil, errors.New(ErrorUnexpectedReservedRune)
// 			}
// 			q := &Node{verb: currVerb, phrase: s[i:j]}
// 			root.AddChild(q)
// 			i = j + 1
// 		}
// 	}

// 	// TODO:
// 	// Collapse unnecessary hierarchy.
// 	if len(root.children) == 1 && root.children[0].IsLeaf() {
// 		root = root.children[0]
// 		root.parent = nil
// 		root.src = s
// 	}

// 	if !root.IsValid() {
// 		return nil, errors.New(ErrorMalformedQuery)
// 	}

// 	tree := &Tree{
// 		root: root,
// 		src:  s,
// 	}

// 	return tree, nil
// }
