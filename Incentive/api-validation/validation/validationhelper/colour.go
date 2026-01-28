package validationhelper

import (
	"regexp"
	"strings"
)

var (
	hexColorRegex = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	rgbRegex      = regexp.MustCompile(`^rgba?\s*\(\s*\d+\s*,\s*\d+\s*,\s*\d+\s*(,\s*[\d.]+\s*)?\)$`)
	hslRegex      = regexp.MustCompile(`^hsla?\s*\(\s*\d+\s*,\s*\d+%\s*,\s*\d+%\s*(,\s*[\d.]+\s*)?\)$`)
)

// IsValidColour validates if a string represents a valid color.
// Supports: hex (#RGB, #RRGGBB, #RRGGBBAA), rgb/rgba, hsl/hsla, and named colors.
func IsValidColour(s string) bool {
	if len(s) == 0 {
		return false
	}

	s = strings.TrimSpace(strings.ToLower(s))

	// Check hex color
	if hexColorRegex.MatchString(s) {
		return true
	}

	// Check RGB/RGBA
	if rgbRegex.MatchString(s) {
		return true
	}

	// Check HSL/HSLA
	if hslRegex.MatchString(s) {
		return true
	}

	// Check named colors (common ones)
	return isNamedColor(s)
}

func isNamedColor(s string) bool {
	namedColors := map[string]bool{
		"black": true, "white": true, "red": true, "green": true, "blue": true,
		"yellow": true, "cyan": true, "magenta": true, "gray": true, "grey": true,
		"orange": true, "purple": true, "pink": true, "brown": true, "lime": true,
		"navy": true, "teal": true, "aqua": true, "maroon": true, "olive": true,
		"silver": true, "fuchsia": true, "transparent": true,
	}
	return namedColors[s]
}
