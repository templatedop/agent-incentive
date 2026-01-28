//go:generate ./longitude.go
package longitude

// Location is a struct for testing longitude validation
type Location struct {
	// +govalid:longitude
	Lon string `json:"lon"`

	// +govalid:longitude
	StartLon string `json:"start_lon"`

	// +govalid:longitude
	EndLon string `json:"end_lon"`
}
