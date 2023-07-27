package cryptonote

import (
	"encoding/hex"

	ed25519 "filippo.io/edwards25519"
)

const privateKeySize = 32

type PrivateKey struct {
	key *ed25519.Scalar
}

// Hex formats the key as a hex string
func (k *PrivateKey) Hex() string {
	return hex.EncodeToString(k.key.Bytes())
}

// String formats the key as a 0x-prefixed hex string
func (k *PrivateKey) String() string {
	return "0x" + k.Hex()
}

// Public returns the PublicKey corresponding to this PrivateKey.
func (k *PrivateKey) Public() *PublicKey {
	pk := new(ed25519.Point).ScalarBaseMult(k.key)
	return &PublicKey{
		key: pk,
	}
}
