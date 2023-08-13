package polyseed

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateNewSeedPhrase(t *testing.T) {
	seedPhrase, err := CreateNewSeedPhrase(EnglishLang)
	require.NoError(t, err)
	require.Len(t, seedPhrase, NumSeedWords)

	sd, err := CreateSeedData(seedPhrase)
	require.NoError(t, err)
	require.Zero(t, sd.features)
	t.Logf("New phrase birthdate: %s", time.Unix(sd.BirthDate(), 0))
}

// TODO: Rename test
func TestCreateSeedData(t *testing.T) {
	const (
		password            = "пароль123"
		expectedKeyNoPass   = "3b02f326adfd6f20e8599ead565390685add9551dd4c421743dbd5f40028ee28"
		expectedKeyWithPass = "9caf676adee9b0463e4b6b8db69cec6e7b37b6b257195b7741cf70ea7c0e3d45"
	)

	seedStr := "filter vocal snow cupboard volume avoid sign slot drum replace shrug resist pear kiwi bag bring"
	seeds := strings.Split(seedStr, " ")
	seedData, err := CreateSeedData(seeds)
	require.NoError(t, err)

	// test unencrypted key generation
	key := hex.EncodeToString(seedData.KeyGen())
	require.Equal(t, expectedKeyNoPass, key)

	// encrypt
	require.False(t, seedData.IsEncrypted())
	seedData.Crypt(password)
	require.True(t, seedData.IsEncrypted())

	// test generated key after encryption
	key = hex.EncodeToString(seedData.KeyGen())
	require.Equal(t, expectedKeyWithPass, key)
}
