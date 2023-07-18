// Package jamtis is for libraries to manage the keys and addresses
// used when Monero switches to the Seraphis protocol. See here
// for details:
// https://gist.github.com/tevador/50160d160d24cfc6c52ae02eb3d17024?permalink_comment_id=4240591#47-wallet-public-keys
package jamtis

import (
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/curve25519"
)

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
