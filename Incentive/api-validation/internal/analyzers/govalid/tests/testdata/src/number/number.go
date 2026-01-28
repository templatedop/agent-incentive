//go:generate ./number.go
package number

// Measurement is a struct for testing number validation
type Measurement struct {
	// +govalid:number
	Value string `json:"value"`

	// +govalid:number
	Quantity string `json:"quantity"`

	// +govalid:number
	Amount string `json:"amount"`
}
