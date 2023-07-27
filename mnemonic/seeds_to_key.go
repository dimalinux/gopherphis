package mnemonic

import (
	"encoding/binary"
	"errors"
	"fmt"

	"ekyu.moe/cryptonight"
)

// CreateKeyFromSeedsWithoutChecksum creates and returns a 32-byte key from the
// 24 passed seeds.
func (wl *WordList) CreateKeyFromSeedsWithoutChecksum(mnemonic []string) ([]byte, error) {

	if len(mnemonic) != 24 {
		return nil, fmt.Errorf("expected 24 seeds, found %d", len(mnemonic))
	}

	var key [32]byte

	for i := 0; i < 8; i++ {
		// w1, w2, and w3 are all uint32
		w1 := wl.FindIndex32(mnemonic[i*3+0])
		w2 := wl.FindIndex32(mnemonic[i*3+1])
		w3 := wl.FindIndex32(mnemonic[i*3+2])

		x := w1 +
			WordListSize*(((WordListSize-w1)+w2)%WordListSize) +
			WordListSize*WordListSize*(((WordListSize-w2)+w3)%WordListSize)

		// 3 seeds represent more than 2^32 unique values (1626^3 > 2^32).
		// Invalid sequences will overflow the 32-bits and trigger this error.
		// I think it would be better to use uint64 and then check if the result
		// is more than 2^32-1, but this version is matching what the C++ does.
		if x%WordListSize != w1 {
			return nil, fmt.Errorf("invalid seed sequence starting on the %d seed",
				i*3+1) // using human-based index instead of zero-based
		}
		binary.LittleEndian.PutUint32(key[i*4:i*4+4], x)
	}

	return key[:], nil
}

// CreateSeedsWithoutChecksumFromKey creates a 24-seed mnemonic for the given
// 32-byte key without the 25th checksum seed.
func (wl *WordList) CreateSeedsWithoutChecksumFromKey(key []byte) []string {
	if len(key) != 32 {
		panic(fmt.Sprintf("expected key of 32 bytes but found %d", len(key)))
	}
	// return the 24 seeds that make-up the key without the checksum seed
	seeds := make([]string, 0, 24)

	// We break the 32-byte key into 8 chunks of 4 bytes. Each 4-byte
	// section is represented by 3 seeds.
	for i := 0; i < 8; i++ {
		word := key[4*i : 4*i+4]
		x := binary.LittleEndian.Uint32(word)
		w1 := x % WordListSize
		w2 := ((x / WordListSize) + w1) % WordListSize
		w3 := ((x / WordListSize / WordListSize) + w2) % WordListSize
		seeds = append(seeds, wl.Entries[w1], wl.Entries[w2], wl.Entries[w3])
	}
	return seeds
}

// CreateSeedsFromKey returns the 25 seeds for the passed 32-byte key. Note that
// this seed list is for direct key creation without a password. The original
// seed list that created the key, if it used a password, will be different that
// the returned seed list.
func (wl *WordList) CreateSeedsFromKey(key []byte) []string {
	seeds := wl.CreateSeedsWithoutChecksumFromKey(key)
	// This can only error if the seeds we just generated above are not in the
	// wordlist, so we can safely panic here.
	checkSumSeed, err := wl.GetChecksumWord(seeds)
	if err != nil {
		panic(err)
	}
	return append(seeds, checkSumSeed)
}

// CreateKeyFromSeeds creates a key from the passed mnemonic seeds. This method
// requires a checksum seed and should consist of 25 or 13 seeds.
func (wl *WordList) CreateKeyFromSeeds(seeds []string) ([]byte, error) {
	if len(seeds) != 25 && len(seeds) != 13 {
		return nil, fmt.Errorf("expected 25 or 13 seeds, but found %d", len(seeds))
	}

	keySeeds := seeds[:len(seeds)-1]
	checkSum := seeds[len(seeds)-1]

	expectedCheckSum, err := wl.GetChecksumWord(keySeeds)
	if err != nil {
		return nil, err
	}

	if checkSum != expectedCheckSum {
		return nil, errors.New("checksum seed does not match")
	}

	if len(keySeeds) == 12 {
		keySeeds = append(keySeeds, keySeeds...)
	}

	key, err := wl.CreateKeyFromSeedsWithoutChecksum(keySeeds)
	if err != nil {
		return nil, err
	}

	return key[:], nil
}

// CreateKeyFromSeedsAndPassword returns the 32 byte key for the given seeds and
// optional password.
func (wl *WordList) CreateKeyFromSeedsAndPassword(mnemonic []string, password string) ([]byte, error) {
	key, err := wl.CreateKeyFromSeeds(mnemonic)
	if err != nil {
		return nil, err
	}

	if len(password) > 0 {
		hash := cryptonight.Sum([]byte(password), 0)
		key = scSub(key, hash)
	}

	return key, nil
}

// CreateKeyFromSeedsAndPassword auto-detects the seed language and returns the
// 32-byte key for the given seeds and optional password.
func CreateKeyFromSeedsAndPassword(mnemonic []string, password string) ([]byte, error) {
	wl, err := FindLanguage(mnemonic)
	if err != nil {
		return nil, err
	}

	return wl.CreateKeyFromSeedsAndPassword(mnemonic, password)
}

// CreateKeyFromSeeds auto-detects the seed language and returns the
// 32-byte key for the given seeds.
func CreateKeyFromSeeds(seeds []string) ([]byte, error) {
	return CreateKeyFromSeedsAndPassword(seeds, "")
}
