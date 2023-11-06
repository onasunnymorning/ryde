package ryde

import (
	"strings"
)

// Removes \n (newlines)
// Standardizes multiple spaces or tabs (\t) to one space
// Removes leading and trailing spaces
func StandardizeString(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(strings.Replace(strings.ReplaceAll(s, "\n", ""), "\t", " ", -1)), " "))
}

// Runs StandardizeString on all elements of a slice of strings
func StandardizeStringSlice(s []string) []string {
	for i, v := range s {
		s[i] = StandardizeString(v)
	}
	return s
}
