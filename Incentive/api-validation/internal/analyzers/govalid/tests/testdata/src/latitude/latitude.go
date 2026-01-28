//go:generate ./latitude.go
package latitude

// Location is a struct for testing latitude validation
type Location struct {
	// +govalid:latitude
	Lat string `json:"lat"`

	// +govalid:latitude
	StartLat string `json:"start_lat"`

	// +govalid:latitude
	EndLat string `json:"end_lat"`
}
