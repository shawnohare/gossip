package gossip

import (
	"errors"
)

type Verb rune

// Modal verbs.
const (
	VerbError Verb = Verb(-2)
	Should    Verb = Verb(Tilde)
	Not       Verb = Verb(Minus)
	Must      Verb = Verb(Plus)
)

// Modal verbs as their literal string representation.
const (
	VerbErrorString string = "_error"
	ShouldString    string = "~"
	NotString       string = "-"
	MustString      string = "+"
)

// Human readable modal verbs.
const (
	VerbErrorStringPretty string = "_error"
	ShouldStringPretty    string = "should"
	NotStringPretty       string = "not"
	MustStringPretty      string = "must"
)

var verbStringsForHumans = map[rune]string{
	rune(Must):      MustStringPretty,
	rune(Not):       NotStringPretty,
	rune(Should):    ShouldStringPretty,
	rune(VerbError): VerbErrorStringPretty,
}

var verbStrings = map[rune]string{
	rune(Must):      MustString,
	rune(Not):       NotString,
	rune(Should):    ShouldString,
	rune(VerbError): VerbErrorString,
}

var verbStringLookup = map[string]Verb{
	ShouldString:       Should,
	ShouldStringPretty: Should,
	NotString:          Not,
	NotStringPretty:    Not,
	MustString:         Must,
	MustStringPretty:   Must,
}

func (v Verb) String() string {
	if vs, ok := verbStrings[rune(v)]; ok {
		return vs
	}
	return VerbErrorString
}

func (v Verb) Pretty() string {
	if vs, ok := verbStringsForHumans[rune(v)]; ok {
		return vs
	}
	return VerbErrorStringPretty
}

// IsValid reports whether the instance is a valid modal verb.
func (v Verb) IsValid() bool {
	return v == Should || v == Must || v == Not
}

// IsMust reports whether the instance is the modal verb "must".
func (v Verb) IsMust() bool {
	return v == Must
}

// IsMustNot reports whether the instance is the modal verb "must not".
func (v Verb) IsMustNot() bool {
	return v == Not
}

// IsShould reports whether the instance is the modal verb "should".
func (v Verb) IsShould() bool {
	return v == Should
}

// IsRuneVerb states if the input represents a modal verb such as "must".
func IsRuneVerb(r rune) bool {
	v := Verb(r)
	return v == Should || v == Must || v == Not
}

// IsRuneMust states if the input represents the modal verb "must".
func IsRuneMust(r rune) bool {
	return Verb(r) == Must
}

// IsRuneMustNot states if the input represents the modal verb "must not".
func IsRuneMustNot(r rune) bool {
	return Verb(r) == Not
}

// IsRuneShould states if the input represents the modal verb "should".
func IsRuneShould(r rune) bool {
	return Verb(r) == Should
}

// ParseVerbString converts a string representation of a verb into the
// appropriate Verb instance. If the input string is not a valid Verb,
// the VerbError Verb instance is returned, along with an error.
func ParseVerbString(verb string) (Verb, error) {
	if v, ok := verbStringLookup[verb]; ok {
		return v, nil
	}
	return VerbError, errors.New(ErrorVerbString)
}
