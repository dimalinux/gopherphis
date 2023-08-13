package polyseed

import (
	"slices"
	"strings"
)

const (
	// LangSize is the number of words in the seed dictionary of each language.
	LangSize = 2048
)

// Lang holds the individual language configurations for polyseed's word lists.
type Lang struct {
	Name        string
	EnglishName string
	Separator   string
	IsSorted    bool
	HasPrefix   bool
	HasAccents  bool
	Compose     bool
	Words       [LangSize]string
}

// Languages is an array of all polyseed supported languages.
func Languages() []*Lang {
	return []*Lang{
		EnglishLang,
		JapaneseLang,
		KoreanLang,
		SpanishLang,
		FrenchLang,
		ItalianLang,
		CzechLang,
		PortugueseLang,
		ChineseSimpleLang,
		ChineseLang,
	}
}

type seedCompare func(a, b string) int

func (l *Lang) getComparer() seedCompare {
	if l.HasPrefix {
		if l.HasAccents {
			return comparePrefixNoAccent
		}
		return comparePrefix
	}

	if l.HasAccents {
		return compareNoAccent
	}

	return strings.Compare
}

// getIndex returns the index of the word if found. The 2nd return value
// indicates if the word was found.
func (l *Lang) getIndex(word string, cmp seedCompare) (uint16, bool) {
	word = strings.ToLower(word)

	if l.IsSorted {
		i, found := slices.BinarySearchFunc(l.Words[:], word, cmp)
		return uint16(i), found
	}

	// linear search, since language is not sorted
	for i := uint16(0); i < LangSize; i++ {
		if cmp(word, l.Words[i]) == 0 {
			return i, true
		}
	}

	return 0, false
}

func (l *Lang) getIndexes(words []string) []uint16 {
	cmp := l.getComparer()
	indexes := make([]uint16, 0, len(words))
	for _, w := range words {
		index, found := l.getIndex(w, cmp)
		if !found {
			return nil
		}
		indexes = append(indexes, index)
	}
	return indexes
}

func (l *Lang) getWords(indexes []uint16) []string {
	words := make([]string, 0, len(indexes))

	for _, idx := range indexes {
		words = append(words, l.Words[idx])
	}

	return words
}

func getIndexes(words []string) ([]uint16, error) {
	for _, l := range Languages() {
		indexes := l.getIndexes(words)
		if indexes != nil {
			return indexes, nil
		}
	}
	return nil, ErrLang
}
