package gossip

type Verb rune

// Modal verbs.
const (
	VerbError Verb = Verb(-2)
	Should    Verb = Verb(Tilde)
	MustNot   Verb = Verb(Minus)
	Must      Verb = Verb(Plus)
)

// Modal verbs as their literal string representation.
const (
	VerbErrorString string = "_error"
	ShouldString    string = "~"
	MustNotString   string = "-"
	MustString      string = "+"
)

// Human readable modal verbs.
const (
	VerbErrorPretty string = "_error"
	ShouldPretty    string = "should"
	MustNotPretty   string = "must not"
	MustPretty      string = "must"
)

var verbStringsForHumans = map[rune]string{
	rune(Must):      MustPretty,
	rune(MustNot):   MustNotPretty,
	rune(Should):    ShouldPretty,
	rune(VerbError): VerbErrorPretty,
}

var verbStrings = map[rune]string{
	rune(Must):      MustString,
	rune(MustNot):   MustNotString,
	rune(Should):    ShouldString,
	rune(VerbError): VerbErrorString,
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
	return VerbErrorPretty
}

// IsValid reports whether the instance is a valid modal verb.
func (v Verb) IsValid() bool {
	return v == Should || v == Must || v == MustNot
}

// IsMust reports whether the instance is the modal verb "must".
func (v Verb) IsMust() bool {
	return v == Must
}

// IsMustNot reports whether the instance is the modal verb "must not".
func (v Verb) IsMustNot() bool {
	return v == MustNot
}

// IsShould reports whether the instance is the modal verb "should".
func (v Verb) IsShould() bool {
	return v == Should
}

// IsRuneVerb states if the input represents a modal verb such as "must".
func IsRuneVerb(r rune) bool {
	v := Verb(r)
	return v == Should || v == Must || v == MustNot
}

// IsRuneMust states if the input represents the modal verb "must".
func IsRuneMust(r rune) bool {
	return Verb(r) == Must
}

// IsRuneMustNot states if the input represents the modal verb "must not".
func IsRuneMustNot(r rune) bool {
	return Verb(r) == MustNot
}

// IsRuneShould states if the input represents the modal verb "should".
func IsRuneShould(r rune) bool {
	return Verb(r) == Should
}
