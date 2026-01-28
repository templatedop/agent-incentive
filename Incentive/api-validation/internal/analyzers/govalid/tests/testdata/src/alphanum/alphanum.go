//go:generate ./alphanum.go
package alphanum

// Identifier is a struct for testing alphanum validation
type Identifier struct {
	// +govalid:alphanum
	Code string `json:"code"`

	// +govalid:alphanum
	SKU string `json:"sku"`

	// +govalid:alphanum
	Token string `json:"token"`
}
