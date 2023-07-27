package mnemonic

import (
	"errors"
	"hash/crc32"
)

var (
	errSeedNotInList = errors.New("seed not in list")
)

// GetChecksumWord determines the correct language for the passed mnemonic
// string and returns the expected checksum word for the list. The checksum
// calculation is language specific, as only the language specific prefix length
// number symbols are used in its creation. See the WordList member function
// with this same name for an additional documentation.
func GetChecksumWord(mnemonic []string) (string, *WordList, error) {
	for _, l := range WordLists {
		w, err := l.GetChecksumWord(mnemonic)
		if err != nil {
			if errors.Is(err, errSeedNotInList) {
				continue
			}
			return "", nil, err
		}
		return w, l, nil
	}
	return "", nil, errors.New("could not find language with all mnemonic seeds")
}

// GetChecksumWord returns the checksum word for a given mnemonic. While you
// should probably be passing a 24 seed slice (or in some cases 12), the code
// only requires that you pass at least one seed. The checksum word is always
// one of the passed mnemonic seeds.
func (wl *WordList) GetChecksumWord(mnemonic []string) (string, error) {
	if len(mnemonic) == 0 {
		// API misuse, so panic
		panic("no seeds to compute checksum from")
	}

	hash := crc32.NewIEEE()
	for _, word := range mnemonic {
		if !wl.HasWord(word) {
			return "", errSeedNotInList
		}
		_, _ = hash.Write([]byte(prefix(word, wl.PrefixSz)))
	}
	sum := hash.Sum32()
	idx := sum % uint32(len(mnemonic))
	return mnemonic[idx], nil
}

// HasWord returns whether the passed seed is in the current language's
// wordlist. Only the first N symbols are checked, where N is PrefixSz of the
// given language's settings.
func (wl *WordList) HasWord(seed string) bool {
	_, ok := wl.PrefixMap[prefix(seed, wl.PrefixSz)]
	return ok
}

// HasWords returns whether all the passed words are in the current language's
// wordlist, based on the prefix N symbols, as well as the number of words that
// were an exact match.
func (wl *WordList) HasWords(seeds []string) (bool, int) {
	exactCount := 0
	for _, s := range seeds {
		index, ok := wl.PrefixMap[prefix(s, wl.PrefixSz)]
		if !ok {
			return false, -1
		}
		if s == wl.Entries[index] {
			exactCount++
		}
	}

	return true, exactCount
}
