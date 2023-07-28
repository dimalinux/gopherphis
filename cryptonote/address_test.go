package cryptonote

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dimalinux/gopherphis/mnemonic"
)

func TestNewAddress(t *testing.T) {
	const addrStr = "42ey1afDFnn4886T7196doS9GPMzexD9gXpsZJDwVjeRVdFCSoHnv7KPbBeGpzJBzHRCAs9UxqeoyFQMYbqSWYTfJJQAWDm"
	addr, err := NewAddress(addrStr, Mainnet)
	require.NoError(t, err)
	require.Equal(t, addrStr, addr.String())
}

func TestNewAddress_fail(t *testing.T) {
	_, err := NewAddress("fake", Mainnet)
	require.ErrorIs(t, err, errInvalidAddressLength)
}

func TestValidateAddress(t *testing.T) {
	kp, err := GenerateKeys()
	require.NoError(t, err)
	pubKeys := kp.PublicKeyPair()

	// mainnet address checks
	addr := pubKeys.Address(Mainnet)
	require.NoError(t, addr.ValidateNet(Mainnet))
	require.ErrorIs(t, addr.ValidateNet(Stagenet), errInvalidPrefixGotMainnet)

	// stagenet address checks
	addr = pubKeys.Address(Stagenet)
	require.NoError(t, addr.ValidateNet(Stagenet))
	require.ErrorIs(t, addr.ValidateNet(Mainnet), errInvalidPrefixGotStagenet)

	// testnet address check
	const testnetAddress = "9ujeXrjzf7bfeK3KZdCqnYaMwZVFuXemPU8Ubw335rj2FN1CdMiWNyFV3ksEfMFvRp9L9qum5UxkP5rN9aLcPxbH1au4WAB" //nolint:lll
	require.NoError(t, addr.UnmarshalText([]byte(testnetAddress)))
	require.ErrorIs(t, addr.ValidateNet(Mainnet), errInvalidPrefixGotTestnet)

	// uninitialized address validation
	addr = new(Address) // empty
	require.ErrorIs(t, addr.ValidateNet(Mainnet), errAddressNotInitialized)
}

func TestValidateAddress_loop(t *testing.T) {
	// Tests our address encoding/decoding with randomised data
	for i := 0; i < 1000; i++ {
		kp, err := GenerateKeys() // create random key
		require.NoError(t, err)
		// Generate the address, convert it to its base58 string form,
		// then convert the base58 form back into a new address, then
		// verify that the bytes of the 2 addresses are identical.
		addr1 := kp.PublicKeyPair().Address(Mainnet)
		addr2, err := NewAddress(addr1.String(), Mainnet)
		require.NoError(t, err)
		require.Equal(t, addr1.String(), addr2.String())
	}
}

func TestAddress_Equal(t *testing.T) {
	kp, err := GenerateKeys() // create random key
	require.NoError(t, err)
	pubKeys := kp.PublicKeyPair()

	addr1 := pubKeys.Address(Mainnet)
	addr2 := pubKeys.Address(Mainnet)
	addr3 := pubKeys.Address(Stagenet)

	require.False(t, addr1.Equal(nil))
	require.True(t, addr1.Equal(addr1)) // identity

	require.False(t, addr1 == addr2)    // the pointers are unique,
	require.True(t, addr1.Equal(addr2)) // but the values are the same

	require.False(t, addr1.Equal(addr3)) // same keys, but different network
}

func Test_createPrivateSpendKeyFromSeeds_knownAddress(t *testing.T) {
	expectedAddress := "44mTQkfkgg7UjMTjJuGT9kVhsp6vf4NKHdBJwHxWPVjsBzsEN1KVWtA2hEEvK3JpAE1ZStqksrypG1bAcNnH7hXEL5W88M4"
	seeds := []string{
		"wedge", "mundane", "shocking", "muffin", "ritual", "gnaw", "tumbling", "yearbook",
		"truth", "flying", "ponies", "obvious", "menu", "edited", "gauze", "sequence",
		"bugs", "ongoing", "iguana", "emulate", "aimless", "hawk", "getting", "gossip", "menu",
	}

	spendKeyBytes, err := mnemonic.EnglishWordList.CreateKeyFromSeeds(seeds)
	require.NoError(t, err)

	spendKey, err := NewPrivateSpendKey(spendKeyBytes)
	require.NoError(t, err)

	key, err := spendKey.AsPrivateKeyPair()
	require.NoError(t, err)

	address := key.PublicKeyPair().Address(Mainnet).String()
	require.Equal(t, expectedAddress, address)
}

const (
	testCaseAccounts         = 3
	testCaseAccountAddresses = 3
)

type SeedTestCase struct {
	Language         string   `json:"language"`
	Seeds            []string `json:"seeds"`
	Password         string   `json:"seedPassword"`
	PrivateSpendKey  string   `json:"privateSpendKey"`
	PrivateViewKey   string   `json:"privateViewKey"`
	AccountAddresses [testCaseAccounts]struct {
		AccountIndex uint32                           `json:"accountIndex"` // informative only, values are in order
		Addresses    [testCaseAccountAddresses]string `json:"addresses"`
	} `json:"accountAddresses"`
}

func TestWordList_CreateKeysAndAddressesFromSeeds(t *testing.T) {
	data, err := os.ReadFile("testdata/address_tests.json")
	require.NoError(t, err)

	var testCases []*SeedTestCase
	err = json.Unmarshal(data, &testCases)
	require.NoError(t, err)

	for i, tc := range testCases {
		failMsg := fmt.Sprintf("case %d failed", i)
		key, err := mnemonic.CreateKeyFromSeedsAndPassword(tc.Seeds, tc.Password)
		require.NoError(t, err, failMsg)

		spendKey, err := NewPrivateSpendKey(key)
		require.NoError(t, err)
		//require.Equal(t, tc.PrivateSpendKey, hex.EncodeToString(spendKey.Bytes()))

		viewKey, err := spendKey.PrivateViewKey()
		require.NoError(t, err)
		require.Equal(t, tc.PrivateViewKey, hex.EncodeToString(viewKey.Bytes()))

		keyPair, err := spendKey.AsPrivateKeyPair()
		require.NoError(t, err)
		primaryAddr := keyPair.PublicKeyPair().Address(Mainnet).String()
		require.Equal(t, tc.AccountAddresses[0].Addresses[0], primaryAddr)

		for j := uint32(0); j < testCaseAccounts; j++ {
			for k := uint32(0); k < testCaseAccountAddresses; k++ {
				pubKeyPair := keyPair.SubAddrPubKeyPair(j, k)
				subAddr := pubKeyPair.Address(Mainnet).String()
				require.Equal(t, tc.AccountAddresses[j].Addresses[k], subAddr, failMsg)
			}
		}

		if (i+1)%100 == 0 {
			t.Logf("%d subtests completed", i+1)
		}
	}
}
