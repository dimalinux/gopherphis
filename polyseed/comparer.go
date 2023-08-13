package polyseed

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	// NumRunePrefix is the maximum number of utf-8 symbols that are compared
	// when matching a seed word with a polyseed Lang entry when the language's
	// HasPrefix field set to true.
	NumRunePrefix = 4
)

// prefix returns the first 4 runes (UTF-8 symbols) of the passed word as a
// string. If the word is less than 4 runes, the word is returned unchanged.
func prefix(word string) string {
	wordSymbols := []rune(word)
	if len(wordSymbols) > NumRunePrefix {
		word = string(wordSymbols[:NumRunePrefix])
	}
	return word
}

// removeAccents replaces symbols with accents with their equivalent symbols without
// the accent. Example: "peñón" => "penon".
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	if err != nil {
		panic(err) // not reachable
	}
	return output
}

// comparePrefix compares two strings ignoring any values after the 4th UTF-8
// symbol. "abcde" and "abcdf" are treated as equivalent, because the divergence
// only happens on the 5th symbol.
func comparePrefix(key, elem string) int {
	return strings.Compare(prefix(key), prefix(elem))
}

// compareNoAccent compares 2 strings treating accented characters as identical
// to their non-accented counterparts.
func compareNoAccent(key, elem string) int {
	return strings.Compare(removeAccents(key), removeAccents(elem))
}

// comparePrefix compares two strings ignoring any values after the 4th UTF-8
// symbol and also treating any accented characters as identical to their
// non-accented counterparts.
func comparePrefixNoAccent(key, elem string) int {
	return comparePrefix(removeAccents(key), removeAccents(elem))
}
