//go:generate ./lowercase.go
package lowercase

// User is a struct for testing lowercase validation
type User struct {
	// +govalid:lowercase
	Username string `json:"username"`

	// +govalid:lowercase
	Email string `json:"email"`

	// +govalid:lowercase
	Slug string `json:"slug"`
}
