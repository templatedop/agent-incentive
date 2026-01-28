package validationhelper

import "regexp"

var ddmmyy = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/([0-9]{2})$`)

// IsValidDateDDMMYY reports whether s matches dd/mm/yy.
func IsValidDateDDMMYY(s string) bool {
	return ddmmyy.MatchString(s)
}
