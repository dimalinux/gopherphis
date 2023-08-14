package polyseed

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that ordered languages are ordered using their selected compare algorithm.
func TestLangOrdered(t *testing.T) {
	for _, lang := range Languages() {
		if !lang.IsSorted {
			continue
		}
		cmp := lang.getComparer()
		for i := 1; i < LangSize; i++ {
			require.Equal(t, -1, cmp(lang.Words[i-1], lang.Words[i]))
		}
	}
}

func TestPrefixesUnique(t *testing.T) {
	for _, lang := range Languages() {
		if !lang.HasPrefix {
			continue
		}
		var seenMap = make(map[string]bool)
		for i := 0; i < LangSize; i++ {
			wordPrefix := prefix(lang.Words[i])
			require.False(t, seenMap[wordPrefix])
			seenMap[wordPrefix] = true
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

	// test reversing back to words
	words := EnglishLang.getWords(indexes)
	require.EqualValues(t, seeds, words)
}

func TestLang_getIndex(t *testing.T) {
	indexes := []uint16{0, LangSize/2 - 1, LangSize / 2, LangSize/2 + 1, LangSize - 1}
	for _, lang := range Languages() {
		cmp := lang.getComparer()
		for _, idx := range indexes {
			word := lang.Words[idx]
			foundIdx, found := lang.getIndex(word, cmp)
			require.True(t, found)
			require.Equal(t, idx, foundIdx)
		}
		_, found := lang.getIndex("xxxxxx", cmp)
		require.False(t, found)
	}
}
