package validationhelper

// IsNumber validates if a string contains only numeric characters.
// Accepts optional leading sign (+/-) and decimal point.
func IsNumber(s string) bool {
	if len(s) == 0 {
		return false
	}

	start := 0
	// Allow leading sign
	if s[0] == '+' || s[0] == '-' {
		if len(s) == 1 {
			return false
		}
		start = 1
	}

	hasDigit := false
	hasDot := false

	for i := start; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			hasDigit = true
		} else if c == '.' {
			if hasDot {
				return false // Multiple dots
			}
			hasDot = true
		} else {
			return false
		}
	}

	return hasDigit
}
