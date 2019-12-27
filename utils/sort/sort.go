package sort

import "sort"

// LexicographicalOrder sorts a string slice lexicographically.
func LexicographicalOrder(names []string) {
	sort.Strings(names)
}
