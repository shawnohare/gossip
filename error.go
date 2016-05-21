package gossip

// Define some common error codes.
const (
	ErrorMalformedQuery         = "gossip: Search query is malformed. "
	ErrorUnpairedQuotation      = ErrorMalformedQuery + "Unpaired quotation mark."
	ErrorUnpairedBracket        = ErrorMalformedQuery + "Unpaired bracket."
	ErrorVerbSequence           = ErrorMalformedQuery + "Unexpected verb sequence."
	ErrorUnexpectedReservedRune = ErrorMalformedQuery + "Unexpected reserved rune."
	ErrorEmptyQuery             = ErrorMalformedQuery + "Semantically empty."
)
