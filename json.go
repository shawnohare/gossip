package gossip

import "encoding/json"

// UnmarshalJSON converts a JSON encoded representation of a Node (or tree)
// into an instance where the Parent field is correctly set.
func UnmarshalJSON(data []byte) (*Node, error) {
	var root *Node
	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}
	return root.Tree(), nil
}
