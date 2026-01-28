//go:generate ./excludesall.go
package excludesall

// SafeInput is a struct for testing excludesall validation
type SafeInput struct {
	// +govalid:excludesall=<>
	NoHTML string `json:"no_html"`

	// +govalid:excludesall='\"
	NoQuotes string `json:"no_quotes"`

	// +govalid:excludesall=;|&
	NoShellChars string `json:"no_shell_chars"`
}
