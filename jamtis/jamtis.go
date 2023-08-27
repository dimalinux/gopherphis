// Package jamtis is for libraries to manage the keys and addresses
// used when Monero switches to the Seraphis protocol. See here
// for details:
// https://gist.github.com/tevador/50160d160d24cfc6c52ae02eb3d17024?permalink_comment_id=4240591#47-wallet-public-keys
package jamtis

import (
	ed25519 "filippo.io/edwards25519"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/twofish"

	"github.com/dimalinux/gopherphis/mcrypto"
)

const (
	prefix                               = "monero"
	hashKeyJamtisIndexExtensionGenerator = "jamtis_index_extension_generator"
	hashKeyJamtisSpendKeyExtensionG      = "jamtis_spendkey_extension_g"
	hashKeyJamtisSpendKeyExtensionX      = "jamtis_spendkey_extension_x"
	hashKeyJamtisSpendKeyExtensionU      = "jamtis_spendkey_extension_u"
)

type Address struct {
	K1  []byte
	K2  []byte
	K3  []byte
	Tag []byte
}

func keyDerive1(key []byte, name string) ([]byte, error) {
	const prefix = "monero"

	h, err := blake2b.New(32, key)
	if err != nil {
		return nil, err
	}

	_, err = h.Write([]byte(prefix + name))
	if err != nil {
		return nil, err
	}

	h32 := h.Sum(nil)
	h32[0] &= 0xf8
	h32[31] &= 0x7f

	return h32, nil
}

func secretDerive(key []byte, name string) ([]byte, error) {
	const prefix = "monero"

	h, err := blake2b.New(32, key)
	if err != nil {
		return nil, err
	}

	_, err = h.Write([]byte(prefix + name))
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func genViewBalanceKey(masterKey []byte) ([]byte, error) {
	return keyDerive1(masterKey, "jamtis_view_balance_key")
}

func genUnlockAmountsKey(viewBalanceKey []byte) ([]byte, error) {
	return keyDerive1(viewBalanceKey, "jamtis_unlock_amounts_key")
}

func genFindReceivedKey(viewBalanceKey []byte) ([]byte, error) {
	return keyDerive1(viewBalanceKey, "jamtis_find_received_key")
}

func genGenAddressSecret(viewBalanceKey []byte) ([]byte, error) {
	return secretDerive(viewBalanceKey, "jamtis_generate_address_secret")
}

func genCipherTagSecret(genAddressSecret []byte) ([]byte, error) {
	return secretDerive(genAddressSecret, "jamtis_cipher_tag_secret")
}

func genUnlockAmountsPubKey(unlockAmountsPrivKey []byte) []byte {
	var pubKey [32]byte
	x25519ScalarMult(pubKey[:], unlockAmountsPrivKey, curve25519.Basepoint)
	return pubKey[:]
}

func genFindReceivedPubKey(findReceivedPrivKey []byte, unlockAmountsPubKey []byte) []byte {
	var pubKey [32]byte
	x25519ScalarMult(pubKey[:], findReceivedPrivKey, unlockAmountsPubKey)
	return pubKey[:]
}

func genJamtisAddressV1(
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

	cipherTagSecret, err := genCipherTagSecret(generateAddressPrivKey)
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
	addressExtensionKeyU, err := genJamtisSpendKeyExtension(hashKeyJamtisSpendKeyExtensionU, spendPubKey, generateAddress, j)
	if err != nil {
		return nil, err
	}

	//k^j_x
	addressExtensionKeyX, err := genJamtisSpendKeyExtension(hashKeyJamtisSpendKeyExtensionX, spendPubKey, generateAddress, j)
	if err != nil {
		return nil, err
	}

	//k^j_g
	addressExtensionKeyG, err := genJamtisSpendKeyExtension(hashKeyJamtisSpendKeyExtensionG, spendPubKey, generateAddress, j)
	if err != nil {
		return nil, err
	}

	extendedSeraphisSpendKey, err := new(ed25519.Point).SetBytes(spendPubKey)
	if err != nil {
		return nil, err
	}

	//k^j_u U + K_s
	extendSeraphisSpendKey(addressExtensionKeyU, getUPoint(), extendedSeraphisSpendKey)

	//k^j_x X + (k^j_u U + K_s)
	extendSeraphisSpendKey(addressExtensionKeyX, getXPoint(), extendedSeraphisSpendKey)

	//k^j_g G + (k^j_x X + k^j_u U + K_s)
	extendedSeraphisSpendKey.Add(extendedSeraphisSpendKey, new(ed25519.Point).ScalarBaseMult(addressExtensionKeyG))

	return extendedSeraphisSpendKey.Bytes(), nil
}

func extendSeraphisSpendKey(
	addressExtensionKey *ed25519.Scalar,
	generatorPt *ed25519.Point,
	addressSpendKeyInOut *ed25519.Point,
) {
	extenderKey := new(ed25519.Point).ScalarMult(addressExtensionKey, generatorPt)
	addressSpendKeyInOut.Add(addressSpendKeyInOut, extenderKey)
}

func blake2bHash(optionalKey []byte, hashSize int, inputs ...any) ([]byte, error) {
	h, err := blake2b.New(hashSize, optionalKey)
	if err != nil {
		return nil, err
	}

	for _, input := range inputs {
		var inputBytes []byte
		switch arg := input.(type) {
		case string:
			inputBytes = []byte(arg)
		case []byte:
			inputBytes = arg
		case *ed25519.Scalar:
			inputBytes = arg.Bytes()
		case *ed25519.Point:
			inputBytes = arg.Bytes()
		default:
			panic("invalid hash input type")
		}

		_, err = h.Write(inputBytes)
		if err != nil {
			return nil, err
		}
	}

	return h.Sum(nil), nil
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
