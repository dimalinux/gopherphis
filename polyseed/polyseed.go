// Package polyseed implements Monero's newer mnemonic seed scheme that only
// requires 16 seeds and records the wallet's approximate birthday as part of
// the seed phrase.
package polyseed

import (
	"crypto/rand"
)

const (
	// NumSeedWords is the number of seed words is a polyseed mnemonic phrase.
	// Each word contains 11 bits of information as the word lists are 2048 in
	// size. 16 words * 11 bits/word = 176 bits: 150 bits of secret entropy, 11
	// bits of checksum, 10 bits for the wallet birthday, and 5 feature bits.
	NumSeedWords = 16

	// NumSecretBits is the number of entropy bits stored in the seed phrase
	// which are stretched to the 256-bit key.
	NumSecretBits = 150

	// kdfNumIterations is the number of iterations used when
	kdfNumIterations = 10000
)

// Coin is set to zero for Monero, but we are not trying to prevent
// this library from being used by other coins, so it is supported.
// The maximum supported Coin value is 2047.
type Coin uint16

// Constants for the known coin values
const (
	MoneroCoin = Coin(0)
	AEONCoin   = Coin(1)
	// If this repo adds a new coin, we'll add the matching value here.
	// https://github.com/tevador/polyseed
)

// defaultCoin is a variable that can be set to use this polyseed package with
// another coin besides Monero. This parameter is package wide, so you can only
// have one coin at a time.
var defaultCoin = MoneroCoin

// SetPackageCoin allows you to retarget the package for another coin other
// than Monero.
func SetPackageCoin(coin Coin) {
	if coin >= 2048 {
		panic("invalid coin")
	}
	defaultCoin = coin
}

// CreateNewSeedPhrase creates polyseed mnemonic using random values for the 150
// secret bits. Setting feature bits is not supported, but will be when/if a use
// case appears.
func CreateNewSeedPhrase(lang *Lang) ([]string, error) {
	seed := &SeedData{
		birthday: birthdayNow(),
		features: 0,
	}

	if _, err := rand.Read(seed.secret[:]); err != nil {
		return nil, err
	}
	seed.secret[numKeyEntropyBytes-1] &= clearBitMask

	// encode polynomial
	poly := dataToPoly(seed)

	return lang.getWords(poly.coeff[:]), nil
}

// CreateSeedData initializes a SeedData object with the passed seed phrase
// and returns it.
func CreateSeedData(seedWords []string) (*SeedData, error) {
	if len(seedWords) != NumSeedWords {
		return nil, ErrNumWords
	}

	indexes, err := getIndexes(seedWords)
	if err != nil {
		return nil, err
	}

	poly := &poly{}
	copy(poly.coeff[:], indexes)
	clear(indexes)

	// Finalize the polynomial. The coin value needs to be xor'ed before
	// checksum validation.
	poly.coeff[numChecksumWords] ^= uint16(defaultCoin)

	if !poly.ValidChecksum() {
		return nil, ErrChecksum
	}

	seed := poly.ToSeedData()

	// Encrypted is the only feature we support at the current time.
	if seed.features & ^uint(encryptedMask) != 0 {
		return nil, ErrUnsupported
	}

	return seed, nil
}
