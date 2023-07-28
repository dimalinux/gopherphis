package mnemonic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The prefixes in this test's seed list match both French and English, but
// French will get picked even though English is is tested first, as the exact
// matches are higher for French.
func TestFindLanguage_ambiguousSeeds(t *testing.T) {
	seeds := []string{
		"anneau",
		"annoncer",
		"apercevoir",
		"apparence",
		"appel",
		"apporter",
		"apprendre",
		"appuyer",
		"arbre",
		"arcade",
		"arceau",
		"arche",
		"ardeur",
		"argent",
		"argile",
		"aride",
		"arme",
		"armure",
		"arracher",
		"arriver",
		"article",
		"asile",
		"aspect",
		"assaut",
		"arcade",
	}

	// omitted exact match count is zero
	hasPrefixes, exactMatchCount := EnglishWordList.HasWords(seeds)
	require.True(t, hasPrefixes)
	require.Zero(t, exactMatchCount)

	// omitted exact match count is 25
	hasPrefixes, exactMatchCount = FrenchWordList.HasWords(seeds)
	require.True(t, hasPrefixes)
	require.Equal(t, 25, exactMatchCount)

	wl, err := FindLanguage(seeds)
	require.NoError(t, err)
	require.True(t, wl == FrenchWordList)
}
