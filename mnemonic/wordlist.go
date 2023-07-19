// Package mnemonic is code for seed and address handing before Seraphis.
package mnemonic

import (
	"fmt"
)

const (
	// WordListSize is the number of seed words in the lists. All languages use
	// the same number of seed words.
	WordListSize = 1626
)

// WordList is the seed word list configuration for a language
type WordList struct {
	Name        string
	EnglishName string
	PrefixSz    uint
	Entries     []string
	PrefixMap   map[string]uint32 // map from word prefixes to their indices
}

// prefix is a port of Monero's utf8prefix method (much simpler in Go).
//
//	https://github.com/monero-project/monero/blob/v0.18.2.2/src/mnemonics/language_base.h#L52-L75
func prefix(s string, prefixLen uint) string {
	sr := []rune(s)
	if uint(len(sr)) < prefixLen {
		return s
	}
	return string(sr[:prefixLen])
}

func newWordList(name string, englishName string, prefixSz uint, entries []string) *WordList {
	if len(entries) != WordListSize {
		// this function is designed to be called with static initialization data, so panic'ing
		// is appropriate.
		panic(fmt.Sprintf("initialization error: WordList hash size %d, expected %d", len(entries), WordListSize))
	}

	pm := make(map[string]uint32)
	for i := 0; i < WordListSize; i++ {
		pm[prefix(entries[i], prefixSz)] = uint32(i)
	}

	return &WordList{
		Name:        name,
		EnglishName: englishName,
		PrefixSz:    prefixSz,
		Entries:     entries,
		PrefixMap:   pm,
	}
}

// WordLists is an array of all available seed word configs
var WordLists = []*WordList{
	EnglishWordList,
	RussianWordList,
	SpanishWordList,
	PortugueseWordList,
	GermanWordList,
	FrenchWordList,
	ItalianWordList,
	DutchWordList,
	JapaneseWordList,
	ChineseSimplifiedWordList,
}

// FindIndex32 returns the index of the passed word. This method should
// only be used when you are certain that the word exists in the list.
func (w *WordList) FindIndex32(word string) uint32 {
	p := prefix(word, w.PrefixSz)
	i, ok := w.PrefixMap[p]
	if !ok {
		panic("word not found")
	}

	return i
}
