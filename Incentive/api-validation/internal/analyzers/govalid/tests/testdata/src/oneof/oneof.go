//go:generate ./oneof.go
package oneof

// Priority is a struct for testing oneof validation
type Priority struct {
	// +govalid:oneof=low medium high
	Level string `json:"level"`

	// +govalid:oneof=draft published archived
	Status string `json:"status"`

	// +govalid:oneof=red green blue
	Color string `json:"color"`
}
