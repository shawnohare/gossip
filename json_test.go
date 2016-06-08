package gossip

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSONBad(t *testing.T) {
	var (
		n   *Node
		err error
	)

	n, err = UnmarshalJSON(nil)
	assert.Error(t, err)
	assert.Nil(t, n)

	n, err = UnmarshalJSON([]byte{})
	assert.Error(t, err)
	assert.Nil(t, n)

	n, err = UnmarshalJSON([]byte("garbage"))
	assert.Error(t, err)
	assert.Nil(t, n)
}

func TestUnmarshalJSON(t *testing.T) {
	n0 := NewNode()
	data, err := json.Marshal(n0)
	assert.NoError(t, err)
	n1, err := UnmarshalJSON(data)
	assert.Equal(t, n0, n1)
}
