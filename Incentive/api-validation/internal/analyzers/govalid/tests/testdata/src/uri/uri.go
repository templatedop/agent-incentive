//go:generate ./uri.go
package uri

// Resource is a struct for testing uri validation
type Resource struct {
	// +govalid:uri
	Endpoint string `json:"endpoint"`

	// +govalid:uri
	Reference string `json:"reference"`

	// +govalid:uri
	Location string `json:"location"`
}
