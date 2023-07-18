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
}

func TestJamtisKeys(t *testing.T) {
	masterKey, err := hex.DecodeString(tc.masterKey)
	require.NoError(t, err)

	viewBalanceKey, err := genViewBalanceKey(masterKey)
	require.NoError(t, err)
	require.Equal(t, tc.viewBalanceKey, hex.EncodeToString(viewBalanceKey))

	unlockAmountsKey, err := genUnlockAmountsKey(viewBalanceKey)
	require.NoError(t, err)
	require.Equal(t, tc.unlockAmountsKey, hex.EncodeToString(unlockAmountsKey))

	findReceivedKey, err := genFindReceivedKey(viewBalanceKey)
	require.NoError(t, err)
	require.Equal(t, tc.findReceivedKey, hex.EncodeToString(findReceivedKey))

	genAddressSecret, err := genGenAddressSecret(viewBalanceKey)
	require.NoError(t, err)
	require.Equal(t, tc.generateAddressSecret, hex.EncodeToString(genAddressSecret))

	genCipherTagSecret, err := genCipherTagSecret(genAddressSecret)
	require.NoError(t, err)
	require.Equal(t, tc.cipherTagSecret, hex.EncodeToString(genCipherTagSecret))

	spendKeyBase, err := genSeraphisSpendKey(viewBalanceKey, masterKey)
	require.NoError(t, err)
	require.Equal(t, tc.jamtisSpendKeyBase, hex.EncodeToString(spendKeyBase))

	unlockAmountsPubKey := genUnlockAmountsPubKey(unlockAmountsKey)
	require.Equal(t, tc.unlockAmountsPubKey, hex.EncodeToString(unlockAmountsPubKey))

	findReceivedPubKey := genFindReceivedPubKey(findReceivedKey, unlockAmountsPubKey)
	require.Equal(t, tc.findReceivedPubKey, hex.EncodeToString(findReceivedPubKey))
}
