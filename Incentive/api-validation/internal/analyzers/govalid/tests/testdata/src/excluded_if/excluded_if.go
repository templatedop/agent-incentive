//go:generate ./excluded_if.go
package excluded_if

// Account is a struct for testing excluded_if validation
type Account struct {
	Type string `json:"type"`

	// +govalid:excluded_if=Type guest
	Password string `json:"password"`

	Plan string `json:"plan"`

	// +govalid:excluded_if=Plan free
	CreditCard string `json:"credit_card"`
}
