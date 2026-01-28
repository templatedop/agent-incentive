package validationhelper

import "unicode"

// IsAlphanum validates if a string contains only alphanumeric characters.
// Returns false for empty strings.
func IsAlphanum(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
