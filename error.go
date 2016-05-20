package gossip

// Define some common error codes.
const (
	ErrorMalformedQuery         = "Search query is malformed: "
	ErrorUnpairedQuotation      = ErrorMalformedQuery + "Unpaired quotation mark."
	ErrorUnpairedBracket        = ErrorMalformedQuery + "Unpaired bracket."
	ErrorVerbSequence           = ErrorMalformedQuery + "Unexpected verb sequence."
	ErrorUnexpectedReservedRune = ErrorMalformedQuery + "Unexpected reserved rune."
)
