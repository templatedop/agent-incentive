//go:generate ./excluded_unless.go
package excluded_unless

// Order is a struct for testing excluded_unless validation
type Order struct {
	DeliveryMethod string `json:"delivery_method"`

	// +govalid:excluded_unless=DeliveryMethod home_delivery
	DeliveryAddress string `json:"delivery_address"`

	PaymentType string `json:"payment_type"`

	// +govalid:excluded_unless=PaymentType invoice
	InvoiceNumber string `json:"invoice_number"`
}
