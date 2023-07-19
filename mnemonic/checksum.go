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
func (w *WordList) GetChecksumWord(mnemonic []string) (string, error) {
	if len(mnemonic) == 0 {
		// API misuse, so panic
		panic("no seeds to compute checksum from")
	}

	hash := crc32.NewIEEE()
	for _, word := range mnemonic {
		if !w.HasWord(word) {
			return "", errSeedNotInList
		}
		_, _ = hash.Write([]byte(prefix(word, w.PrefixSz)))
	}
	sum := hash.Sum32()
	idx := sum % uint32(len(mnemonic))
	return mnemonic[idx], nil
}

// HasWord returns whether the passed word is in the current language's
// wordlist. Only the first N symbols are checked, where N is PrefixSz of the
// given language's settings.
func (w *WordList) HasWord(word string) bool {
	_, ok := w.PrefixMap[prefix(word, w.PrefixSz)]
	return ok
}
