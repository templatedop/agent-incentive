//go:generate ./boolean.go
package boolean

// Settings is a struct for testing boolean validation
type Settings struct {
	// +govalid:boolean
	Enabled string `json:"enabled"`

	// +govalid:boolean
	Active string `json:"active"`

	// +govalid:boolean
	FlagValue string `json:"flag_value"`
}
