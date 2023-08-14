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

func changeDefaultCoinForTest(t *testing.T, coin Coin) {
	SetPackageCoin(coin)
	t.Cleanup(func() {
		SetPackageCoin(MoneroCoin)
	})
}

func TestCreateSeedData_AEON(t *testing.T) {
	const (
		seedStr             = "적성 큰딸 그토록 순수 매달 불꽃 점점 개성 상업 부장 놀이 편지 시각 발음 사탕 이념" //nolint:lll
		password            = "qwerty123"
		expectedKeyNoPass   = "140660311eb94ffbed063c796fa9c37ee5433dabda716ddee18ca957bfa32ab7"
		expectedKeyWithPass = "99d73e3628b8959dbb1536d516fcfe6722b43f666d0debd68ce34ec1c60953a4"
	)

	changeDefaultCoinForTest(t, AEONCoin)

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
