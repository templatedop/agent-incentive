//go:generate ./required_unless.go
package required_unless

// Payment is a struct for testing required_unless validation
type Payment struct {
	PaymentMethod string `json:"payment_method"`

	// +govalid:required_unless=PaymentMethod cash
	CardNumber string `json:"card_number"`

	AccountType string `json:"account_type"`

	// +govalid:required_unless=AccountType guest
	AccountID string `json:"account_id"`
}
