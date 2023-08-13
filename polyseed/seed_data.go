package polyseed

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

// Errors that API calls may return to polyseed library users.
var (
	ErrNumWords    = errors.New("wrong number of words in the phrase")
	ErrLang        = errors.New("unknown language or unsupported words")
	ErrChecksum    = errors.New("checksum mismatch")
	ErrUnsupported = errors.New("unsupported seed features")
	ErrFormat      = errors.New("invalid seed format")
)

const (
	// KeySizeBytes is the size in bytes of the key 256-bit key generated from
	// the seed words. The longer key size is derived from 150 bits of entropy.
	KeySizeBytes = 32

	// numKeyEntropyBytes is the number of bytes required to store the 150 bits
	// of key entropy bits. 19*8=152, leaving 2 extra clear bits in the final byte.
	numKeyEntropyBytes = 19

	// clearBitMask is a mask to keep the lower 6 entropy bits in the final byte of an
	// array of key entropy data, while ignoring the 2 unused clear bits.
	clearBitMask = 0x3F
)

// SeedData is structure for serialization/deserialization
type SeedData struct {
	birthday uint
	features uint
	// padded with zeroes for future compatibility with longer seeds
	secret   [KeySizeBytes]uint8
	checksum uint16
}

// Clear attempts to clear the secret and other fields of the SeedData. Due to
// limitations of Golang, this is just a best-effort attempt.
func (sd *SeedData) Clear() {
	sd.birthday = 0
	sd.features = 0
	sd.checksum = 0
	for i := 0; i < KeySizeBytes; i++ {
		sd.secret[i] = 0
	}
}

// BirthDate returns the wallet's birthday in Unix epoch time.
func (sd *SeedData) BirthDate() int64 {
	return birthdayDecode(sd.birthday)
}

// IsEncrypted returns whether the SeedData has been encrypted with a password.
func (sd *SeedData) IsEncrypted() bool {
	return sd.featureEnabled(encryptedMask)
}

func (sd *SeedData) featureEnabled(featureMask uint) bool {
	return (sd.features & featureMask) == featureMask
}

// KeyGen creates and returns the 32-byte key for the SeedData.
func (sd *SeedData) KeyGen() []byte {
	var salt [32]byte
	copy(salt[:], "POLYSEED key")
	salt[13] = 0xff
	salt[14] = 0xff
	salt[15] = 0xff
	le := binary.LittleEndian
	le.PutUint32(salt[16:], uint32(defaultCoin)) // domain separate by coin
	le.PutUint32(salt[20:], uint32(sd.birthday)) // domain separate by birthday
	le.PutUint32(salt[24:], uint32(sd.features)) // domain separate by features

	return pbkdf2.Key(sd.secret[:], salt[:], kdfNumIterations, KeySizeBytes, sha256.New)
}

// Crypt encrypts or decrypts the seed data with a password.
func (sd *SeedData) Crypt(password string) {
	if len(password) == 0 {
		return
	}

	var salt [16]byte
	copy(salt[:], "POLYSEED mask")
	salt[14] = 0xff
	salt[15] = 0xff

	// derive an encryption mask
	mask := pbkdf2.Key([]byte(password), salt[:], kdfNumIterations, 32, sha256.New)

	// apply mask
	for i := 0; i < numKeyEntropyBytes; i++ {
		sd.secret[i] ^= mask[i]
	}

	sd.secret[numKeyEntropyBytes-1] &= clearBitMask
	sd.features ^= encryptedMask // flip the encrypted bit

	// encode polynomial
	poly := dataToPoly(sd)
	sd.checksum = poly.coeff[0] // TODO: Where should this be set?!?
}
