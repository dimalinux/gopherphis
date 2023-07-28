// Package cryptonote is for libraries to manage the keys and addresses
// used before Jamtis.
package cryptonote

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	ed25519 "filippo.io/edwards25519"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/dimalinux/gopherphis/mcrypto"
)

const (
	// KeySize is the size, in bytes, of both public and private keys
	// used in cryptonote.
	KeySize = 32
)

var (
	errInvalidInput = errors.New("input is not 32 bytes")
)

// PrivateKeyPair represents a monero private spend and view key.
type PrivateKeyPair struct {
	sk *PrivateSpendKey
	vk *PrivateViewKey
}

// NewPrivateKeyPairFromBytes returns a new PrivateKeyPair given the canonical byte representation of
// a private spend and view key.
func NewPrivateKeyPairFromBytes(skBytes, vkBytes []byte) (*PrivateKeyPair, error) {
	if len(skBytes) != KeySize || len(vkBytes) != KeySize {
		return nil, errInvalidInput
	}

	sk, err := ed25519.NewScalar().SetCanonicalBytes(skBytes)
	if err != nil {
		return nil, err
	}

	vk, err := ed25519.NewScalar().SetCanonicalBytes(vkBytes)
	if err != nil {
		return nil, err
	}

	return &PrivateKeyPair{
		sk: &PrivateSpendKey{key: sk},
		vk: &PrivateViewKey{key: vk},
	}, nil
}

// SpendKeyBytes returns the canonical byte encoding of the private spend key.
func (kp *PrivateKeyPair) SpendKeyBytes() []byte {
	return kp.sk.key.Bytes()
}

// PublicKeyPair returns the PublicKeyPair corresponding to the PrivateKeyPair
func (kp *PrivateKeyPair) PublicKeyPair() *PublicKeyPair {
	return &PublicKeyPair{
		isSubAddress: false,
		sk:           kp.sk.Public(),
		vk:           kp.vk.Public(),
	}
}

// subAddressSecretKey creates and returns the private key used when generating the
// public keys for a subaddress.
func (kp *PrivateKeyPair) subAddressSecret(accountIndex uint32, subAddrIndex uint32) *ed25519.Scalar {
	if accountIndex == 0 && subAddrIndex == 0 {
		panic("accountIndex=0, subAddrIndex=0 is not a subaddress")
	}

	const prefix = "SubAddr\000"
	const hashSize = len(prefix) + 32 + 2*4
	b := make([]byte, 0, hashSize)
	b = append(b, []byte(prefix)...)
	b = append(b, kp.PrivateViewKey().Bytes()...)
	b = binary.LittleEndian.AppendUint32(b, accountIndex)
	b = binary.LittleEndian.AppendUint32(b, subAddrIndex)

	h := ethcrypto.Keccak256(b)
	h = mcrypto.ScReduce32(h)
	s, err := ed25519.NewScalar().SetCanonicalBytes(h)
	if err != nil {
		panic("ed25519 error: setting scalar failed")
	}
	return s
}

// SubAddrPubKeyPair returns the PublicKeyPair of the requested subaddress.
func (kp *PrivateKeyPair) SubAddrPubKeyPair(accountIndex uint32, subAddrIndex uint32) *PublicKeyPair {

	if accountIndex == 0 && subAddrIndex == 0 {
		// It's a primary key pair, not a subaddress
		return kp.PublicKeyPair()
	}

	subAddrSecret := kp.subAddressSecret(accountIndex, subAddrIndex)
	subAddrSecretPub := new(ed25519.Point).ScalarBaseMult(subAddrSecret)

	spendKeyPub := new(ed25519.Point).Add(kp.sk.Public().key, subAddrSecretPub)
	viewKeyPub := new(ed25519.Point).ScalarMult(kp.vk.key, spendKeyPub)

	return &PublicKeyPair{
		isSubAddress: true,
		sk:           &PublicKey{key: spendKeyPub},
		vk:           &PublicKey{key: viewKeyPub},
	}
}

// SpendKey returns the key pair's private spend key
func (kp *PrivateKeyPair) SpendKey() *PrivateSpendKey {
	return kp.sk
}

// PrivateViewKey returns the key pair's private view key
func (kp *PrivateKeyPair) PrivateViewKey() *PrivateViewKey {
	return kp.vk
}

// PrivateSpendKey represents a monero private spend key
type PrivateSpendKey struct {
	key *ed25519.Scalar
}

// NewPrivateSpendKey returns a new PrivateSpendKey from the given canonically-encoded scalar.
func NewPrivateSpendKey(b []byte) (*PrivateSpendKey, error) {
	if len(b) != KeySize {
		return nil, errInvalidInput
	}

	sk, err := ed25519.NewScalar().SetCanonicalBytes(b)
	if err != nil {
		return nil, err
	}

	return &PrivateSpendKey{
		key: sk,
	}, nil
}

// Public returns the public key corresponding to the private key.
func (k *PrivateSpendKey) Public() *PublicKey {
	pk := new(ed25519.Point).ScalarBaseMult(k.key)
	return &PublicKey{
		key: pk,
	}
}

// Hex formats the key as a hex string
func (k *PrivateSpendKey) Hex() string {
	return hex.EncodeToString(k.key.Bytes())
}

// String formats the key as a 0x-prefixed hex string
func (k *PrivateSpendKey) String() string {
	return "0x" + k.Hex()
}

// AsPrivateKeyPair returns the PrivateSpendKey as a PrivateKeyPair.
func (k *PrivateSpendKey) AsPrivateKeyPair() (*PrivateKeyPair, error) {
	vk, err := k.PrivateViewKey()
	if err != nil {
		return nil, err
	}

	return &PrivateKeyPair{
		sk: k,
		vk: vk,
	}, nil
}

// PrivateViewKey returns the private view key using the standard algorithm from
// the PrivateSpendKey. View keys do not not have to be derived from the
// spend key, but by doing it this way, you preserve compatibility with
// most wallets, some of which require it to be this way.
func (k *PrivateSpendKey) PrivateViewKey() (*PrivateViewKey, error) {
	h := ethcrypto.Keccak256(k.key.Bytes())
	// We can't use SetBytesWithClamping below, which would do the sc_reduce32 computation
	// for us, because standard monero wallets do not modify the first and last byte when
	// calculating the view key.
	vkBytes := mcrypto.ScReduce32(h[:])
	vk, err := ed25519.NewScalar().SetCanonicalBytes(vkBytes[:])
	if err != nil {
		return nil, err
	}

	return &PrivateViewKey{
		key: vk,
	}, nil
}

// Bytes returns the PrivateSpendKey as canonical bytes
func (k *PrivateSpendKey) Bytes() []byte {
	return k.key.Bytes()
}

// PrivateViewKey represents a monero private view key.
type PrivateViewKey struct {
	key *ed25519.Scalar
}

// Public returns the PublicKey corresponding to this PrivateViewKey.
func (k *PrivateViewKey) Public() *PublicKey {
	pk := new(ed25519.Point).ScalarBaseMult(k.key)
	return &PublicKey{
		key: pk,
	}
}

// Bytes returns the canonical 32-byte little-endian encoding of PrivateViewKey.
func (k *PrivateViewKey) Bytes() []byte {
	return k.key.Bytes()
}

// Hex formats the key as a hex string
func (k *PrivateViewKey) Hex() string {
	return hex.EncodeToString(k.key.Bytes())
}

// String formats the key as a 0x-prefixed hex string
func (k *PrivateViewKey) String() string {
	return "0x" + k.Hex()
}

// PublicKeyPair contains a public SpendKey and ViewKey
type PublicKeyPair struct {
	isSubAddress bool
	sk           *PublicKey
	vk           *PublicKey
}

// SpendKey returns the key pair's spend key.
func (kp *PublicKeyPair) SpendKey() *PublicKey {
	return kp.sk
}

// ViewKey returns the key pair's view key.
func (kp *PublicKeyPair) ViewKey() *PublicKey {
	return kp.vk
}

// GenerateKeys generates a private spend key and view key
func GenerateKeys() (*PrivateKeyPair, error) {
	var seed [32]byte
	_, err := rand.Read(seed[:])
	if err != nil {
		return nil, err
	}

	// we hash the seed for compatibility w/ the ed25519 stdlib
	h := sha512.Sum512(seed[:])

	s, err := ed25519.NewScalar().SetBytesWithClamping(h[:32])
	if err != nil {
		return nil, fmt.Errorf("failed to set bytes: %w", err)
	}

	sk := &PrivateSpendKey{key: s}

	return sk.AsPrivateKeyPair()
}
