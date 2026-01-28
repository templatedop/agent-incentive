//go:generate ./required_with_all.go
package required_with_all

// Person is a struct for testing required_with_all validation
type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	// +govalid:required_with_all=FirstName LastName
	FullName string `json:"full_name"`

	Street string `json:"street"`
	City   string `json:"city"`

	// +govalid:required_with_all=Street City
	ZipCode string `json:"zip_code"`
}
