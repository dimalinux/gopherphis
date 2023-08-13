package polyseed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_comparePrefix(t *testing.T) {
	// equal, shorter than prefix
	require.Zero(t, comparePrefix("a", "a"))

	// equal, length of prefix
	require.Zero(t, comparePrefix("abcd", "abcd"))

	// equal, symbols after prefix ignored
	require.Zero(t, comparePrefix("aaaab", "aaaac"))

	// less than zero
	require.Equal(t, -1, comparePrefix("aaab", "aaac"))

	// greater than zero, with multi-byte character set
	require.Equal(t, 1, comparePrefix("ааав", "аааб"))
}

func Test_removeAccents(t *testing.T) {
	// Spanish
	require.Equal(t, "penon", removeAccents("peñón"))

	// French
	require.Equal(t, "eleve", removeAccents("élève"))

	// Russian
	require.Equal(t, "орел", removeAccents("орёл"))

	// Russian е is not equal to the French e
	require.NotEqual(t, removeAccents("é"), removeAccents("ё"))
}

func Test_compareNoAccent(t *testing.T) {
	// shorter words ordered first
	require.True(t, compareNoAccent("pez", "pezuña") < 0)
	require.True(t, compareNoAccent("pezuña", "pez") > 0)

	// Words with accents removed are equal
	require.Zero(t, compareNoAccent("eleve", removeAccents("élève")))

	// Normally, á would come after a, but rábano comes before rabia when
	// accents are normalized.
	require.True(t, compareNoAccent("rábano", "rabia") < 0)
}

func Test_comparePrefixNoAccent(t *testing.T) {
	// equal, shorter than prefix
	require.Zero(t, comparePrefixNoAccent("è", "é"))

	// equal, symbols after prefix ignored
	require.Zero(t, comparePrefixNoAccent("́áááá", "aaaab"))

	// less than and greater than
	require.Equal(t, -1, comparePrefixNoAccent("áááá", "aaab"))
	require.Equal(t, 1, comparePrefixNoAccent("áááb", "aaaa"))
}
