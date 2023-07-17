package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverseSlice(t *testing.T) {
	in := []byte{0xa, 0xb, 0xc}
	expected := []byte{0xc, 0xb, 0xa}
	require.Equal(t, expected, ReverseSlice(in))
	require.Equal(t, []byte{0xa, 0xb, 0xc}, in) // backing array of original slice is unmodified
}
