package validationhelper

import "strconv"

// IsValidLongitude validates if a string represents a valid longitude (-180 to 180).
func IsValidLongitude(s string) bool {
	lon, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	return lon >= -180 && lon <= 180
}
