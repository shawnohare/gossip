package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeMethodsForNil(t *testing.T) {
	var x *Tree
	assert.False(t, x.Equals(x))
	assert.False(t, x.IsValid())
	assert.Equal(t, 0, x.Height())
	assert.Equal(t, "", x.String())
	assert.Nil(t, x.Root())
	assert.Len(t, x.Leaves(), 0)

	r := &Node{}
	y := x.SetRoot(r)
	assert.Equal(t, r, y.Root())
}

func TestTreeSetRoot(t *testing.T) {
	var x *Tree
	r := &Node{}
	y := x.SetRoot(r)
	assert.Equal(t, r, y.Root())

	y = NewTree()
	y.SetRoot(r)
	assert.Equal(t, r, y.Root())
}

func TestTreeIsValid(t *testing.T) {
	h1 := NewTree()
	h1.root.NewChild().SetPhrase("test")

	h2 := NewTree()
	h2.root.NewChild().NewChild().SetPhrase("test")

	h3 := NewTree()
	h3.root.NewChild().SetPhrase("test").NewChild().SetPhrase("test")

	tests := []struct {
		in  *Tree
		out bool
	}{
		{nil, false},
		{NewTree(), false},
		{h1, true},
		{h2, true},
		{h3, false},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.IsValid()
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestTreeLeavesExact(t *testing.T) {
	leaves := []*Node{NewNode(), NewNode(), NewNode(), NewNode()}
	leaves[0].AddChild(leaves[1])
	leaves[1].AddChild(leaves[2])
	leaves[1].AddChild(leaves[3])

	nt := NewTree()

	tests := []struct {
		in  *Tree
		out []*Node
	}{
		{nil, nil},
		{nt, []*Node{nt.Root()}},
		{NewTreeFromRoot(leaves[0]), leaves[2:]},
		{NewTreeFromRoot(leaves[1]), leaves[2:]},
		{NewTreeFromRoot(leaves[2]), leaves[2:3]},
	}

	for i, tt := range tests {
		actual := tt.in.Leaves()
		msg := fmt.Sprintf("Test case (%d) fails", i)
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestTreeString(t *testing.T) {
	h1 := NewTree()
	h1.root.NewChild()

	h2 := NewTree()
	h2.root.NewChild().SetPhrase("x")
	h2.root.NewChild().SetPhrase("y")

	h3 := NewTree()
	c1 := h3.root.NewChild()
	c1.NewChild().SetPhrase("x")
	c1.NewChild().SetPhrase("y")
	c2 := h3.root.NewChild()
	c2.NewChild().SetPhrase("v").SetVerb(Must)
	c2.NewChild().SetPhrase("w").SetVerb(MustNot)

	tests := []struct {
		in  *Tree
		out string
	}{
		{nil, ""},
		{NewTree(), ""},
		{h1, ""},
		{h2, "[x y]"},
		{h3, "[[x y] [+v -w]]"},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.String()
		assert.Equal(t, tt.out, actual, msg)
	}

}

func TestTreeHeight(t *testing.T) {
	h1 := NewTree()
	h1.root.NewChild()

	h2 := NewTree()
	h2.root.NewChild().NewChild()
	h2.root.NewChild()

	h3 := NewTree()
	h3.root.NewChild()
	h3.root.NewChild().NewChild()
	h3.root.NewChild().NewChild().NewChild()

	tests := []struct {
		in  *Tree
		out int
	}{
		{nil, 0},
		{new(Tree), 0},
		{h1, 1},
		{h2, 2},
		{h3, 3},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Test case (%d) %#v fails", i, tt)
		actual := tt.in.Height()
		assert.Equal(t, tt.out, actual, msg)
	}
}
