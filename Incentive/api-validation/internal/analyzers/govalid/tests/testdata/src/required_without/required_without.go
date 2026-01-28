//go:generate ./required_without.go
package required_without

// Contact is a struct for testing required_without validation
type Contact struct {
	Phone string `json:"phone"`

	// +govalid:required_without=Phone
	Email string `json:"email"`

	HomeAddress string `json:"home_address"`

	// +govalid:required_without=HomeAddress
	OfficeAddress string `json:"office_address"`
}
