package gossip

// Define modal verbs.
const (
	VerbError int = -2
	VerbUnknown
	Should  int = 0
	MustNot int = int(Minus)
	Must    int = int(Plus)
)

var verbStringsForHumans = map[int]string{
	Must:      "must",
	MustNot:   "must not",
	Should:    "should",
	VerbError: "_error_",
}

var verbStrings = map[int]string{
	Must:      "+",
	MustNot:   "-",
	Should:    "",
	VerbError: "_error_",
}

// VerbString is used for constructing Tree strings. If the verb is unknown,
// an empty string is returned.
func VerbString(verb int) string {
	return verbStrings[verb]
}

// VerbStringHuman converts an integer modal verb code to
// its human readable counterpart. If the verb is unknown, an empty string
// is returned.
func VerbStringHuman(verb int) string {
	return verbStringsForHumans[verb]
}
