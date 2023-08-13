package polyseed

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that ordered languages are ordered using their selected compare algorithm.
func TestLangOrdered(t *testing.T) {
	for _, lang := range Languages() {
		cmp := lang.getComparer()
		for i := 1; i < LangSize; i++ {
			if lang.IsSorted {
				require.Equal(t, -1, cmp(lang.Words[i-1], lang.Words[i]))
			}
		}
	}
}

func TestLang_getIndexes(t *testing.T) {
	expectedIndices := []uint16{691, 1962, 1644, 430, 1966, 128, 1602, 1632, 540, 1462, 1594, 1467, 1296, 987, 140, 225}
	seedStr := "filter vocal snow cupboard volume avoid sign slot drum replace shrug resist pear kiwi bag bring"
	seeds := strings.Split(seedStr, " ")
	indexes, err := getIndexes(seeds)
	require.NoError(t, err)
	require.EqualValues(t, expectedIndices, indexes)
}
