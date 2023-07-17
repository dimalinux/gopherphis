package util

// ReverseSlice returns a slice of the same type with the elements in the
// reverse order. A new copy is returned and the backing array of the input will
// not be modified.
func ReverseSlice[T comparable](s []T) []T {
	var r []T
	for i := len(s) - 1; i >= 0; i-- {
		r = append(r, s[i])
	}
	return r
}
