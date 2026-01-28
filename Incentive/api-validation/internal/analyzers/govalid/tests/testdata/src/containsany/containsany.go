//go:generate ./containsany.go
package containsany

// Password is a struct for testing containsany validation
type Password struct {
	// +govalid:containsany=!@#$%
	SpecialChars string `json:"special_chars"`

	// +govalid:containsany=0123456789
	HasDigit string `json:"has_digit"`

	// +govalid:containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ
	HasUppercase string `json:"has_uppercase"`
}
