package cryptonote

import (
	"encoding/hex"

	ed25519 "filippo.io/edwards25519"
)

// PublicKey represents a monero public spend, view or subaddress key.
type PublicKey struct {
	key *ed25519.Point
}

// Bytes returns the canonical 32-byte, little-endian encoding of PublicKey.
func (k *PublicKey) Bytes() []byte {
	return k.key.Bytes()
}

// Hex formats the key as a hex string
func (k *PublicKey) Hex() string {
	return hex.EncodeToString(k.key.Bytes())
}

// String formats the key as a 0x-prefixed hex string
func (k *PublicKey) String() string {
	return "0x" + k.Hex()
}
