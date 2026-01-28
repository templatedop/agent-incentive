package validationhelper

import (
	"strings"
)

// IsValidFQDN validates if a string is a Fully Qualified Domain Name.
// FQDN must have at least one dot, valid characters, and proper structure.
func IsValidFQDN(s string) bool {
	if len(s) == 0 || len(s) > 255 {
		return false
	}

	// Remove trailing dot if present (valid in FQDNs)
	s = strings.TrimSuffix(s, ".")

	// Must contain at least one dot
	if !strings.Contains(s, ".") {
		return false
	}

	labels := strings.Split(s, ".")
	if len(labels) < 2 {
		return false
	}

	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}

		// Label must start and end with alphanumeric
		if !isAlphanumericByte(label[0]) || !isAlphanumericByte(label[len(label)-1]) {
			return false
		}

		// Check all characters are valid (alphanumeric or hyphen)
		for i := 0; i < len(label); i++ {
			c := label[i]
			if !isAlphanumericByte(c) && c != '-' {
				return false
			}
		}
	}

	return true
}

func isAlphanumericByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
