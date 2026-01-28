//go:generate ./excludes.go
package excludes

// Content is a struct for testing excludes validation
type Content struct {
	// +govalid:excludes=spam
	Text string `json:"text"`

	// +govalid:excludes=admin
	Username string `json:"username"`

	// +govalid:excludes=test
	Email string `json:"email"`
}
