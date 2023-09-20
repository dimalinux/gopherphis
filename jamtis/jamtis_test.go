package jamtis

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	masterKey             string
	viewBalanceKey        string
	unlockAmountsKey      string
	findReceivedKey       string
	generateAddressSecret string
	cipherTagSecret       string
	jamtisSpendKeyBase    string
	unlockAmountsPubKey   string
	findReceivedPubKey    string
	addressK1             string
	addressK2             string
	addressK3             string
	addressTag            string
}

var tc = testCase{
	masterKey:             "56dca2096a788c5b7c1f2bbd34a7ce1e25eb51c723121b9874d33d7c7fe1b407",
	viewBalanceKey:        "50b19cc007c9cdcdf3c1aefecbb167da67bfd0309cc54fcc2ebb0cbc19f74a19",
	unlockAmountsKey:      "2013823116e9dfdda33eab61fa3db11786323cf66d6a2376a0e9391ee4dd6705",
	findReceivedKey:       "2824824113258b645c97e82aa41d9a1194c64cb53774610659a586a30c41b24b",
	generateAddressSecret: "15178ff2626099384607e078da6d5fb28f365c816797b621aba1b80a28461ca9",
	cipherTagSecret:       "d7afd2426d21bb6282b4dcc647328b513376f820e44f816686a3c9e1ca381673",
	jamtisSpendKeyBase:    "bd829d74e26dd91a8b16f774b8799c8742bbeb9d25ff7e1a57449dc9e5411d79",
	unlockAmountsPubKey:   "437e2a1e5896afd6df54041ff5d9a18fe25814027ccc4a800493910a3f731907",
	findReceivedPubKey:    "37c2d38f79503b9cd4f2685d4aa567fb42e470552866fee1bdf2bb4693054057",
	addressK1:             "2a253091a8005d6ccb3326bb92b0b04eb3e6f0028abe1a66d242e2c5683efe47",
	addressK2:             "dfd0a4f919cef620405ebbd27efe084b7bfb149a01365a13c14a2a3d0ffb4c75",
	addressK3:             "5f5174d730d818f30de3a814c1148b488ed3addb2547f73f0e744ef3e7878246",
	addressTag:            "a435a7bf8076247e6976d48f32e003adcaa9",
}

func TestJamtisKeys(t *testing.T) {
	masterKey, err := hex.DecodeString(tc.masterKey)
	require.NoError(t, err)

	viewBalanceKey, err := GenViewBalancePrivKey(masterKey)
	require.NoError(t, err)
	require.Equal(t, tc.viewBalanceKey, hex.EncodeToString(viewBalanceKey))

	unlockAmountsPrivKey, err := GenUnlockAmountsPrivKey(viewBalanceKey)
	require.NoError(t, err)
	require.Equal(t, tc.unlockAmountsKey, hex.EncodeToString(unlockAmountsPrivKey))

	findReceivedPrivKey, err := GenFindReceivedPrivKey(viewBalanceKey)
	require.NoError(t, err)
	require.Equal(t, tc.findReceivedKey, hex.EncodeToString(findReceivedPrivKey))

	genAddressSecret, err := GenGenAddressSecret(viewBalanceKey)
	require.NoError(t, err)
	require.Equal(t, tc.generateAddressSecret, hex.EncodeToString(genAddressSecret))

	cipherTagSecret, err := GenCipherTagSecret(genAddressSecret)
	require.NoError(t, err)
	require.Equal(t, tc.cipherTagSecret, hex.EncodeToString(cipherTagSecret))

	spendKeyBase, err := GenSeraphisSpendKey(viewBalanceKey, masterKey)
	require.NoError(t, err)
	require.Equal(t, tc.jamtisSpendKeyBase, hex.EncodeToString(spendKeyBase))

	unlockAmountsPubKey := GenUnlockAmountsPubKey(unlockAmountsPrivKey)
	require.Equal(t, tc.unlockAmountsPubKey, hex.EncodeToString(unlockAmountsPubKey))

	findReceivedPubKey := GenFindReceivedPubKey(findReceivedPrivKey, unlockAmountsPubKey)
	require.Equal(t, tc.findReceivedPubKey, hex.EncodeToString(findReceivedPubKey))

	j := [AddressIndexLen]byte{1}
	address, err := GenJamtisAddressV1(spendKeyBase, unlockAmountsPubKey, findReceivedPubKey, genAddressSecret, j[:])
	require.NoError(t, err)
	require.Equal(t, tc.addressK1, hex.EncodeToString(address.K1[:]))
	require.Equal(t, tc.addressK2, hex.EncodeToString(address.K2[:]))
	require.Equal(t, tc.addressK3, hex.EncodeToString(address.K3[:]))
	require.Equal(t, tc.addressTag, hex.EncodeToString(address.Tag))
}
