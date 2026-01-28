//go:generate ./excluded_with.go
package excluded_with

// Preference is a struct for testing excluded_with validation
type Preference struct {
	AutoSave bool `json:"auto_save"`

	// +govalid:excluded_with=AutoSave
	ManualSaveButton string `json:"manual_save_button"`

	DarkMode bool `json:"dark_mode"`

	// +govalid:excluded_with=DarkMode
	LightTheme string `json:"light_theme"`
}
