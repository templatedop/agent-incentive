//go:generate ./dive.go
package dive

// Address is a nested struct
type Address struct {
	// +govalid:required
	Street string `json:"street"`

	// +govalid:required
	City string `json:"city"`

	// +govalid:minlength=5
	// +govalid:maxlength=10
	ZipCode string `json:"zip_code"`
}

// Person is a struct for testing dive validation
type Person struct {
	// +govalid:required
	Name string `json:"name"`

	// +govalid:dive
	Addresses []Address `json:"addresses"`
}
