//go:generate ./required_with.go
package required_with

// Registration is a struct for testing required_with validation
type Registration struct {
	Email string `json:"email"`

	// +govalid:required_with=Email
	EmailConfirmation string `json:"email_confirmation"`

	Phone string `json:"phone"`

	// +govalid:required_with=Phone
	PhoneConfirmation string `json:"phone_confirmation"`
}
