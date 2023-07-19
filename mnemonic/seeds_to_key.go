package mnemonic

import (
	"encoding/binary"
	"fmt"

	"ekyu.moe/cryptonight"
)

// CreateKeyFromSeedsWithoutChecksum creates and returns a spend key from the 24
// passed seeds.
func (w *WordList) CreateKeyFromSeedsWithoutChecksum(mnemonic []string) ([]byte, error) {

	if len(mnemonic) != 24 {
		return nil, fmt.Errorf("expected 24 seeds, found %d", len(mnemonic))
	}

	var privateSpendKey [32]byte

	for i := 0; i < 8; i++ {
		// w1, w2, and w3 are all uint32
		w1 := w.FindIndex32(mnemonic[i*3+0])
		w2 := w.FindIndex32(mnemonic[i*3+1])
		w3 := w.FindIndex32(mnemonic[i*3+2])

		x := w1 +
			WordListSize*(((WordListSize-w1)+w2)%WordListSize) +
			WordListSize*WordListSize*(((WordListSize-w2)+w3)%WordListSize)

		// 3 seeds represent more than 2^32 unique values (1626^3 > 2^32).
		// Invalid sequences will overflow the 32-bits and trigger this error.
		// I think it would be better to use uint64 and then check if the result
		// is more than 2^32-1, but this version is matching what the C++ does.
		if x%WordListSize != w1 {
			return nil, fmt.Errorf("%s %s %s is not a valid seed sequence",
				mnemonic[i*3+0], mnemonic[i*3+1], mnemonic[i*3+2])
		}
		binary.LittleEndian.PutUint32(privateSpendKey[i*4:i*4+4], x)
	}

	return privateSpendKey[:], nil
}

// CreateSeedsWithoutChecksumFromKey creates a 24-seed mnemonic for the given
// 32-byte key without the 25th checksum seed.
func (w *WordList) CreateSeedsWithoutChecksumFromKey(key []byte) []string {
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
		seeds = append(seeds, w.Entries[w1], w.Entries[w2], w.Entries[w3])
	}
	return seeds
}

// CreateSeedsFromKey returns the 25 seeds for the passed 32-byte key. Note that
// this seed list is for direct key creation without a password. The original
// seed list that created the key, if it used a password, will be different that
// the returned seed list.
func (w *WordList) CreateSeedsFromKey(key []byte) []string {
	seeds := w.CreateSeedsWithoutChecksumFromKey(key)
	// This can only error if the seeds we just generated above are not in the
	// wordlist, so we can safely panic here.
	checkSumSeed, err := w.GetChecksumWord(seeds)
	if err != nil {
		panic(err)
	}
	return append(seeds, checkSumSeed)
}

// CreateKeyFromSeeds creates a key from the passed mnemonic seeds. This method
// requires a checksum seed and should consist of 25 or 13 seeds.
func (w *WordList) CreateKeyFromSeeds(mnemonic []string) ([]byte, error) {

	if len(mnemonic) != 25 && len(mnemonic) != 13 {
		return nil, fmt.Errorf("expected 25 or 13 seeds, but found %d", len(mnemonic))
	}

	keySeeds := mnemonic[:len(mnemonic)-1]
	checkSum := mnemonic[len(mnemonic)-1]

	expectedCheckSum, err := w.GetChecksumWord(keySeeds)
	if err != nil {
		return nil, err
	}

	if checkSum != expectedCheckSum {
		return nil, fmt.Errorf("expected %q as checksum but found %q", expectedCheckSum, checkSum)
	}

	if len(keySeeds) == 12 {
		keySeeds = append(keySeeds, keySeeds...)
	}

	privateSpendKey, err := w.CreateKeyFromSeedsWithoutChecksum(keySeeds)
	if err != nil {
		return nil, err
	}

	return privateSpendKey[:], nil
}

// CreateKeyFromSeedsAndPassword returns the 32 byte key for the given seeds and
// optional password.
func (w *WordList) CreateKeyFromSeedsAndPassword(mnemonic []string, password string) ([]byte, error) {
	key, err := w.CreateKeyFromSeeds(mnemonic)
	if err != nil {
		return nil, err
	}

	hash := cryptonight.Sum([]byte(password), 0)

	return scSub(key, hash), nil
}
