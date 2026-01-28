//go:generate ./iscolour.go
package iscolour

// Theme is a struct for testing iscolour validation
type Theme struct {
	// +govalid:iscolour
	Primary string `json:"primary"`

	// +govalid:iscolour
	Secondary string `json:"secondary"`

	// +govalid:iscolour
	Background string `json:"background"`
}
