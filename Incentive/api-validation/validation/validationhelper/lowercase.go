package validationhelper

import (
	"unicode"
)

// IsLowercase validates if all characters in the string are lowercase.
// Empty strings return true. Numbers and symbols are ignored.
func IsLowercase(s string) bool {
	if len(s) == 0 {
		return true
	}

	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}
