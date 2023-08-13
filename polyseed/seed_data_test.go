package polyseed

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeedData_Crypt(t *testing.T) {
	const password = "암호123"

	phrase, err := CreateNewSeedPhrase(KoreanLang)
	require.NoError(t, err)

	sd, err := CreateSeedData(phrase)
	require.NoError(t, err)
	require.False(t, sd.IsEncrypted())
	key1 := hex.EncodeToString(sd.KeyGen())

	// Calling Crypt once creates a new key, but calling Crypt twice produces
	// the identity function back to the original key.

	sd.Crypt(password)
	require.True(t, sd.IsEncrypted())
	key2 := hex.EncodeToString(sd.KeyGen())

	sd.Crypt(password)
	require.False(t, sd.IsEncrypted())
	key3 := hex.EncodeToString(sd.KeyGen())

	require.NotEqual(t, key1, key2)
	require.Equal(t, key1, key3)
}
