package gossip

// Pretty is a printable representation of a Node that is JOSN encodable.
type Pretty struct {
	Verb     string    `json:"verb,omitempty"`
	Children []*Pretty `json:"children,omitempty"`
	Phrase   string    `json:"phrase,omitempty"`
}
