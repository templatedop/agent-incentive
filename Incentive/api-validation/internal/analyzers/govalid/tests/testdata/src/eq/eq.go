//go:generate ./eq.go
package eq

// Status is a struct for testing eq validation
type Status struct {
	// +govalid:eq=active
	State string `json:"state"`

	// +govalid:eq=100
	Count int `json:"count"`

	// +govalid:eq=3.14
	Value float64 `json:"value"`
}
