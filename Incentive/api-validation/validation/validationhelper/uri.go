package validationhelper

import "net/url"

// IsValidURI validates if a string is a valid URI (more permissive than URL).
// Accepts any valid URI scheme, not just HTTP/HTTPS.
func IsValidURI(s string) bool {
	if len(s) == 0 {
		return false
	}

	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	// URI must have a scheme
	return u.Scheme != ""
}
