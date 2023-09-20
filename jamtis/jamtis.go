// Package jamtis is for libraries to manage the keys and addresses
// used when Monero switches to the Seraphis protocol. See here
// for details:
// https://gist.github.com/tevador/50160d160d24cfc6c52ae02eb3d17024?permalink_comment_id=4240591#47-wallet-public-keys
package jamtis

import (
	ed25519 "filippo.io/edwards25519"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/twofish" //nolint:staticcheck

	"github.com/dimalinux/gopherphis/mcrypto"
)

const (
	// AddressIndexLen is the length of the jamtis sub-address identifier in
	// bytes. This index is is generated randomly and not a sequential counter.
	AddressIndexLen = 16

	prefix                               = "monero"
	hashKeyJamtisIndexExtensionGenerator = "jamtis_index_extension_generator"
	hashKeyJamtisSpendKeyExtensionG      = "jamtis_spendkey_extension_g"
	hashKeyJamtisSpendKeyExtensionX      = "jamtis_spendkey_extension_x"
	hashKeyJamtisSpendKeyExtensionU      = "jamtis_spendkey_extension_u"
)

// Address holds the binary fields that make up a Jamtis user address (a
// destination for funds). In the C++ code, this struct is called
// JamtisDestinationV1.
type Address struct {
	K1  []byte // Spend Public Key (32 bytes) - k^j_g G + k^j_x X + k^j_u U + K_s
	K2  []byte // View Public Key (32 bytes) - xk^j_a xK_fr
	K3  []byte // DH Base key (32 bytes) - xk^j_a xK_ua
	Tag []byte // Address tag (18 bytes, does not include tag hint)
}

func keyDerive1(key []byte, name string) ([]byte, error) {
	const prefix = "monero"

	h32, err := blake2bHash(key, 32, []byte(prefix+name))
	if err != nil {
		return nil, err
	}

	h32[0] &= 0xf8
	h32[31] &= 0x7f

	return h32, nil
}

func secretDerive(key []byte, name string) ([]byte, error) {
	const prefix = "monero"
	return blake2bHash(key, 32, []byte(prefix+name))
}

// GenViewBalancePrivKey generates the view balance private key which, in turn,
// is a master key for generating the find-received, unlock-amounts and
// generate-addresses keys.
func GenViewBalancePrivKey(masterPrivKey []byte) ([]byte, error) {
	return keyDerive1(masterPrivKey, "jamtis_view_balance_key")
}

// GenUnlockAmountsPrivKey generates the private key used to encrypt/decrypt the
// amount of a transaction.
func GenUnlockAmountsPrivKey(viewBalancePrivKey []byte) ([]byte, error) {
	return keyDerive1(viewBalancePrivKey, "jamtis_unlock_amounts_key")
}

// GenUnlockAmountsPubKey generates the public key from the unlock amounts
// private key.
func GenUnlockAmountsPubKey(unlockAmountsPrivKey []byte) []byte {
	var pubKey [32]byte
	x25519ScalarMult(pubKey[:], unlockAmountsPrivKey, curve25519.Basepoint)
	return pubKey[:]
}

// GenFindReceivedPrivKey generates the private key used to scan for received
// transaction outputs.
func GenFindReceivedPrivKey(viewBalancePrivKey []byte) ([]byte, error) {
	return keyDerive1(viewBalancePrivKey, "jamtis_find_received_key")
}

// GenFindReceivedPubKey generates the public key from the find received private
// key.
func GenFindReceivedPubKey(findReceivedPrivKey []byte, unlockAmountsPubKey []byte) []byte {
	var pubKey [32]byte
	x25519ScalarMult(pubKey[:], findReceivedPrivKey, unlockAmountsPubKey)
	return pubKey[:]
}

// GenGenAddressSecret generates the secret key used to generate a new address.
// (Note: All addresses in Jamtis are subaddresses.)
func GenGenAddressSecret(viewBalanceKey []byte) ([]byte, error) {
	return secretDerive(viewBalanceKey, "jamtis_generate_address_secret")
}

// GenCipherTagSecret generates the cipher-tag secret which is used to
// encrypt/create the address tag.
func GenCipherTagSecret(genAddressSecret []byte) ([]byte, error) {
	return secretDerive(genAddressSecret, "jamtis_cipher_tag_secret")
}

// GenJamtisAddressV1 generates the Jamtis address data for the passed address
// index.
func GenJamtisAddressV1(
	spendPubKey []byte,
	unlockAmountsPubKey []byte,
	findReceivedPubKey []byte,
	generateAddressPrivKey []byte,
	addressIndex []byte,
) (*Address, error) {

	K1, err := genJamtisAddressSpendKey(spendPubKey, generateAddressPrivKey, addressIndex)
	if err != nil {
		return nil, err
	}

	addressPrivKey, err := genJamtisAddressPrivKey(spendPubKey, generateAddressPrivKey, addressIndex)
	if err != nil {
		return nil, err
	}

	K2 := make([]byte, 32)
	x25519ScalarMult(K2, addressPrivKey, findReceivedPubKey)

	K3 := make([]byte, 32)
	x25519ScalarMult(K3, addressPrivKey, unlockAmountsPubKey)

	cipherTagSecret, err := GenCipherTagSecret(generateAddressPrivKey)
	if err != nil {
		return nil, err
	}

	cipher, err := twofish.NewCipher(cipherTagSecret)
	if err != nil {
		return nil, err
	}

	encryptedAddressIndexAndHint := make([]byte, cipher.BlockSize()+2)
	encryptedAddressIndex := encryptedAddressIndexAndHint[0:cipher.BlockSize()]
	//hint := encryptedAddressIndexAndHint[cipher.BlockSize():]
	// Encrypt the data.
	cipher.Encrypt(encryptedAddressIndex, addressIndex)

	hint, err := genAddressTagHint(cipherTagSecret, encryptedAddressIndex)
	if err != nil {
		return nil, err
	}
	copy(encryptedAddressIndexAndHint[cipher.BlockSize():], hint)

	a := &Address{
		K1:  K1,
		K2:  K2,
		K3:  K3,
		Tag: encryptedAddressIndexAndHint,
	}

	return a, nil
}

func genAddressTagHint(cipherKey []byte, encryptedAddressIndex []byte) ([]byte, error) {
	// assemble hash contents: prefix || 'domain-sep' || k || cipher[k](j)
	const domainSeparator = "jamtis_address_tag_hint"
	return blake2bHash(nil, 2, prefix, domainSeparator, cipherKey, encryptedAddressIndex)
}

func genJamtisAddressSpendKey(spendPubKey []byte, generateAddress []byte, j []byte) ([]byte, error) {

	// K_1 = k^j_g G + k^j_x X + k^j_u U + K_s

	//k^j_u
	addrExtKeyU, err := genJamtisSpendKeyExtension(hashKeyJamtisSpendKeyExtensionU, spendPubKey, generateAddress, j)
	if err != nil {
		return nil, err
	}

	//k^j_x
	addrExtKeyX, err :=
		genJamtisSpendKeyExtension(hashKeyJamtisSpendKeyExtensionX, spendPubKey, generateAddress, j)
	if err != nil {
		return nil, err
	}

	//k^j_g
	addrExtKeyG, err := genJamtisSpendKeyExtension(hashKeyJamtisSpendKeyExtensionG, spendPubKey, generateAddress, j)
	if err != nil {
		return nil, err
	}

	extendedSeraphisSpendKey, err := new(ed25519.Point).SetBytes(spendPubKey)
	if err != nil {
		return nil, err
	}

	//k^j_u U + K_s
	extendSeraphisSpendKey(extendedSeraphisSpendKey, addrExtKeyU, getUPoint())

	//k^j_x X + (k^j_u U + K_s)
	extendSeraphisSpendKey(extendedSeraphisSpendKey, addrExtKeyX, getXPoint())

	//k^j_g G + (k^j_x X + k^j_u U + K_s)
	extendedSeraphisSpendKey.Add(extendedSeraphisSpendKey, new(ed25519.Point).ScalarBaseMult(addrExtKeyG))

	return extendedSeraphisSpendKey.Bytes(), nil
}

func extendSeraphisSpendKey(
	addrSpendKeyInOut *ed25519.Point,
	addrExtKey *ed25519.Scalar,
	generatorPt *ed25519.Point,
) {
	extenderKey := new(ed25519.Point).ScalarMult(addrExtKey, generatorPt)
	addrSpendKeyInOut.Add(addrSpendKeyInOut, extenderKey)
}

func genJamtisSpendKeyExtension(
	domainSeparator string,
	spendPubKey []byte,
	generatorAddress []byte,
	j []byte, // 16-byte address index
) (*ed25519.Scalar, error) {

	// s^j_gen
	generator, err := blake2bHash(generatorAddress, 32, prefix, hashKeyJamtisIndexExtensionGenerator, j)
	if err != nil {
		return nil, err
	}

	// k^j_?
	extensionOut, err := blake2bHash(nil, 64, prefix, domainSeparator, spendPubKey, j, generator)
	if err != nil {
		return nil, err
	}

	return new(ed25519.Scalar).SetCanonicalBytes(mcrypto.ScReduce32(extensionOut))
}

func genJamtisAddressPrivKey(
	spendPubKey []byte,
	generateAddressPrivKey []byte,
	addressIndex []byte,
) ([]byte, error) {
	generator, err := genJamtisIndexExtensionGenerator(generateAddressPrivKey, addressIndex)
	if err != nil {
		return nil, err
	}

	// xk^j_a = H_n_x25519(K_s, j, H_32[s_ga](j))
	const domainSeparator = "jamtis_address_privkey"
	addressPrivKeyOut, err := blake2bHash(nil, 32, prefix, domainSeparator, spendPubKey, addressIndex, generator)
	if err != nil {
		return nil, err
	}

	addressPrivKeyOut[0] &= 255 - 7
	addressPrivKeyOut[31] &= 127

	return addressPrivKeyOut, nil
}

func genJamtisIndexExtensionGenerator(
	generateAddressSecret []byte,
	addressIndex []byte,
) ([]byte, error) {
	const domainSeparator = "jamtis_index_extension_generator"
	return blake2bHash(generateAddressSecret, 32, prefix, domainSeparator, addressIndex)
}
