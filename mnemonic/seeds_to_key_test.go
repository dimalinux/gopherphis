package mnemonic

import (
	"encoding/hex"
	"fmt"
	"testing"

	"ekyu.moe/cryptonight"
	"github.com/stretchr/testify/require"
)

func Test_createPrivateSpendKeyFromSeeds_knownKey(t *testing.T) {
	keyHex := "025d2ad614953e6ecae5a8f557a06ed099a66b337f0d09605aee42385747620b"
	mnemonic := []string{
		"veteran", "weekday", "soil", "husband", "wiring", "idols", "roared", "olympics",
		"needed", "roster", "highway", "demonstrate", "lunar", "stacking", "actress", "onboard",
		"afield", "huge", "scrub", "sieve", "zeal", "buffet", "haunted", "industrial",
		"husband",
	}

	newKey, err := EnglishWordList.CreateKeyFromSeeds(mnemonic)
	require.NoError(t, err)
	require.Equal(t, keyHex, hex.EncodeToString(newKey))
}

func Test_createPrivateSpendKeyFromSeeds_knownAddress(t *testing.T) {
	expectedAddress := "44mTQkfkgg7UjMTjJuGT9kVhsp6vf4NKHdBJwHxWPVjsBzsEN1KVWtA2hEEvK3JpAE1ZStqksrypG1bAcNnH7hXEL5W88M4"
	mnemonic := []string{
		"wedge", "mundane", "shocking", "muffin", "ritual", "gnaw", "tumbling", "yearbook",
		"truth", "flying", "ponies", "obvious", "menu", "edited", "gauze", "sequence",
		"bugs", "ongoing", "iguana", "emulate", "aimless", "hawk", "getting", "gossip", "menu",
	}

	spendKeyBytes, err := EnglishWordList.CreateKeyFromSeeds(mnemonic)
	require.NoError(t, err)

	spendKey, err := NewPrivateSpendKey(spendKeyBytes)
	require.NoError(t, err)

	key, err := spendKey.AsPrivateKeyPair()
	require.NoError(t, err)

	address := key.PublicKeyPair().Address(Mainnet).String()
	require.Equal(t, expectedAddress, address)
}

func Test_CreateKeyFromSeedsWithoutChecksum(t *testing.T) {
	wl := EnglishWordList
	for i := 0; i < WordListSize-24; i++ {
		seeds1 := wl.Entries[i : i+24]

		key, err := wl.CreateKeyFromSeedsWithoutChecksum(seeds1)
		require.NoError(t, err)

		t.Logf("Key is %x", key)
		seeds2 := wl.CreateSeedsWithoutChecksumFromKey(key)
		t.Logf("Seeds are: %s", seeds2)

		require.EqualValues(t, seeds1, seeds2)
	}
}

func TestWordList_CreateKeyFromSeedsAndPassword(t *testing.T) {
	wl := EnglishWordList
	for i, tc := range seedsWithPasswordTests {
		failMsg := fmt.Sprintf("case %d failed", i)
		key, err := wl.CreateKeyFromSeedsAndPassword(tc.mnemonic, tc.password)
		require.NoError(t, err, failMsg)
		require.Equal(t, tc.expectedKeyHex, hex.EncodeToString(key), failMsg)
		if i > 0 && i%100 == 0 {
			t.Logf("case %d succeeded", i)
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
