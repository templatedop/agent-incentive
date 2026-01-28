package validationhelper

import "strconv"

// IsValidLatitude validates if a string represents a valid latitude (-90 to 90).
func IsValidLatitude(s string) bool {
	lat, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	return lat >= -90 && lat <= 90
}
