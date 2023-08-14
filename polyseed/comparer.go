package polyseed

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	// NumPrefixSymbols is the maximum number of utf-8 symbols that are compared
	// when matching a seed word with a polyseed Lang entry when the language's
	// HasPrefix field set to true.
	NumPrefixSymbols = 4
)

// prefix returns the first UTF-8 symbols of the passed word as a string. If the
// word is less than 4 symbols, the word is returned unchanged.
func prefix(word string) string {
	// Accented characters can be spread across more than one rune, so we have
	// to normalize.
	var p []byte
	var symbolIter norm.Iter
	symbolIter.InitString(norm.NFD, word)

	i := 0
	for !symbolIter.Done() {
		p = append(p, symbolIter.Next()...)
		i++
		if i >= NumPrefixSymbols {
			break
		}
	}

	return string(p)
}

// removeAccents replaces symbols with accents with their equivalent symbols without
// the accent. Example: "peñón" => "penon".
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	noAccents, _, err := transform.String(t, s)
	if err != nil {
		return s // doesn't appear to be reachable even with bad UTF-8 input
	}
	return noAccents
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
