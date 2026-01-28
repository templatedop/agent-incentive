//go:generate ./date.go
package date

// Event is a struct for testing date validation
type Event struct {
	// +govalid:date
	StartDate string `json:"start_date"`

	// +govalid:date
	EndDate string `json:"end_date"`

	// +govalid:date
	BirthDate string `json:"birth_date"`
}
