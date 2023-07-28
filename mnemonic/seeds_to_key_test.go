package mnemonic

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"ekyu.moe/cryptonight"
	"github.com/stretchr/testify/require"
)

// Simple, stand-alone test if you want to test something smaller
// than our large regression test set.
func Test_createKeyFromSeeds_knownKey(t *testing.T) {
	const expectedKeyHex = "025d2ad614953e6ecae5a8f557a06ed099a66b337f0d09605aee42385747620b"
	seeds := []string{
		"veteran", "weekday", "soil", "husband", "wiring", "idols", "roared", "olympics",
		"needed", "roster", "highway", "demonstrate", "lunar", "stacking", "actress", "onboard",
		"afield", "huge", "scrub", "sieve", "zeal", "buffet", "haunted", "industrial",
		"husband",
	}

	wl, err := FindLanguage(seeds)
	require.NoError(t, err)
	require.Equal(t, "English", wl.Name)

	newKey, err := wl.CreateKeyFromSeeds(seeds)
	require.NoError(t, err)
	require.Equal(t, expectedKeyHex, hex.EncodeToString(newKey))
}

// This test creates a seed from the wordlist entries, then
// reverses the wordlist entries back from the keys and
// verifies that they match.
func Test_CreateKeyFromSeedsWithoutChecksum(t *testing.T) {
	wl := EnglishWordList
	// There are more seed combinations than the 256 bits in the key. Not all
	// seed sequences will form a key, but these sequences in forward order do.
	for i := 0; i < WordListSize-24; i++ {
		seeds1 := wl.Entries[i : i+24]

		key, err := wl.CreateKeyFromSeedsWithoutChecksum(seeds1)
		require.NoError(t, err)

		// verify that reversing the seeds from the key gives back
		// the original set of seeds.
		seeds2 := wl.CreateSeedsWithoutChecksumFromKey(key)
		require.EqualValues(t, seeds1, seeds2)
	}
}

// This test verifies that all words in every wordlist have the correct index.
// See the README-tests.md for a description of how the test data was generated.
func Test_CreateKeyFromSeeds(t *testing.T) {
	data, err := os.ReadFile("testdata/test_all_seeds_all_langs.json")
	require.NoError(t, err)

	type testCaseSecret struct {
		Secret    string              `json:"secret"`
		LangSeeds map[string][]string `json:"langSeeds"`
	}

	var secrets []*testCaseSecret
	err = json.Unmarshal(data, &secrets)
	require.NoError(t, err)

	for _, tc := range secrets {
		for _, wl := range WordLists {
			secret, err := CreateKeyFromSeeds(tc.LangSeeds[wl.Name])
			require.NoError(t, err)
			require.Equal(t, tc.Secret, hex.EncodeToString(secret))
		}
	}
}

func TestWordList_CreateKeyFromSeedsAndPassword(t *testing.T) {
	data, err := os.ReadFile("testdata/test_seeds_with_passwords.json")
	require.NoError(t, err)

	type testCaseSeeds struct {
		Seeds        []string `json:"seeds"`
		SeedPassword string   `json:"seedPassword"`
		Key          string   `json:"key"`
	}

	var testCases []*testCaseSeeds
	err = json.Unmarshal(data, &testCases)
	require.NoError(t, err)

	for i, tc := range testCases {
		secret, err := CreateKeyFromSeedsAndPassword(tc.Seeds, tc.SeedPassword)
		require.NoError(t, err)
		require.Equal(t, tc.Key, hex.EncodeToString(secret))
		if (i+1)%100 == 0 {
			t.Logf("%d subtests completed", i+1)
		}
	}
}

// Sanity check of the 3rd party library
func TestMoneroHash(t *testing.T) {
	const expectedHash = "bbec2cacf69866a8e740380fe7b818fc78f8571221742d729d9d02d7f8989b87"
	p, err := hex.DecodeString("63617665617420656d70746f72")
	require.NoError(t, err)
	hash := cryptonight.Sum(p, 0)
	require.Len(t, hash, 32)
	require.Equal(t, expectedHash, hex.EncodeToString(hash))
}
