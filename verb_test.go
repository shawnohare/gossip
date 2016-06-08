package gossip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerbIsX(t *testing.T) {
	assert.True(t, Must.IsValid())
	assert.True(t, Must.IsMust())
	assert.False(t, Must.IsMustNot())
	assert.False(t, Must.IsShould())

	assert.True(t, Should.IsValid())
	assert.False(t, Should.IsMust())
	assert.False(t, Should.IsMustNot())
	assert.True(t, Should.IsShould())

	assert.True(t, MustNot.IsValid())
	assert.False(t, MustNot.IsMust())
	assert.True(t, MustNot.IsMustNot())
	assert.False(t, MustNot.IsShould())

	assert.False(t, Verb(-93).IsValid())
	assert.False(t, Verb(-93).IsMust())
	assert.False(t, Verb(-93).IsMustNot())
	assert.False(t, Verb(-93).IsShould())
}

func TestIsVerb(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{-92, false},
		{0, false},
		{rune(VerbError), false},
		{rune(Must), true},
		{rune(Should), true},
		{rune(MustNot), true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, IsRuneVerb(tt.in), msg)
	}
}

func TestIsMust(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{-92, false},
		{0, false},
		{rune(VerbError), false},
		{rune(Must), true},
		{rune(Should), false},
		{rune(MustNot), false},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, IsRuneMust(tt.in), msg)
	}
}

func TestIsMustNot(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{-92, false},
		{0, false},
		{rune(VerbError), false},
		{rune(Must), false},
		{rune(Should), false},
		{rune(MustNot), true},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, IsRuneMustNot(tt.in), msg)
	}
}

func TestIsShould(t *testing.T) {
	tests := []struct {
		in  rune
		out bool
	}{
		{-92, false},
		{0, false},
		{rune(VerbError), false},
		{rune(Must), false},
		{rune(Should), true},
		{rune(MustNot), false},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, IsRuneShould(tt.in), msg)
	}
}

func TestVerbString(t *testing.T) {
	tests := []struct {
		in  Verb
		out string
	}{
		{Must, MustString},
		{Should, ShouldString},
		{MustNot, MustNotString},
		{VerbError, VerbErrorString},
		{999, VerbErrorString},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, tt.in.String(), msg)
	}

}

func TestVerbPretty(t *testing.T) {
	tests := []struct {
		in  Verb
		out string
	}{
		{Must, MustPretty},
		{Should, ShouldPretty},
		{MustNot, MustNotPretty},
		{VerbError, VerbErrorPretty},
		{999, VerbErrorPretty},
	}

	for i, tt := range tests {
		msg := fmt.Sprintf("Fails test case (%d) %s", i, tt.in)
		assert.Equal(t, tt.out, tt.in.Pretty(), msg)
	}

}
