package validationhelper

// IsValidBoolean validates if a string represents a valid boolean value.
// Accepts: "true", "false", "1", "0", "yes", "no", "on", "off" (case-insensitive).
func IsValidBoolean(s string) bool {
	switch s {
	case "true", "false", "TRUE", "FALSE", "True", "False",
		"1", "0",
		"yes", "no", "YES", "NO", "Yes", "No",
		"on", "off", "ON", "OFF", "On", "Off":
		return true
	}
	return false
}
