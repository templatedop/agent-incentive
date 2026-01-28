//go:generate ./required_if.go
package required_if

// Form is a struct for testing required_if validation
type Form struct {
	Type string `json:"type"`

	// +govalid:required_if=Type business
	CompanyName string `json:"company_name"`

	Status string `json:"status"`

	// +govalid:required_if=Status active
	ActiveField string `json:"active_field"`
}
