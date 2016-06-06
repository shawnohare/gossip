package gossip

// JSON is a JSON representation of a Node.
type JSON struct {
	Verb     string  `json:"verb,omitempty"`
	Children []*JSON `json:"children,omitempty"`
	Phrase   string  `json:"phrase,omitempty"`
}
