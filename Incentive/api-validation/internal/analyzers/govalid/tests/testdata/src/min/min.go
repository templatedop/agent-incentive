//go:generate ./min.go
package min

// Product is a struct for testing min validation
type Product struct {
	// +govalid:min=10
	Price int `json:"price"`

	// +govalid:min=0
	Quantity int `json:"quantity"`

	// +govalid:min=18
	Age int `json:"age"`
}
