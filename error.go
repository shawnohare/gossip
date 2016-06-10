package gossip

// Define some common error codes.
const (
	ErrorMalformedQuery         = "gossip: Search query is malformed. "
	ErrorUnpairedQuotation      = ErrorMalformedQuery + "Unpaired quotation mark."
	ErrorUnpairedBracket        = ErrorMalformedQuery + "Unpaired bracket."
	ErrorUnexpectedReservedRune = ErrorMalformedQuery + "Unexpected reserved rune."
	ErrorEmptyQuery             = ErrorMalformedQuery + "Semantically empty."
	ErrorVerbSequence           = ErrorMalformedQuery + "Unexpected verb sequence."
	ErrorVerbString             = "gossip: Verb string is not recognized."
)
