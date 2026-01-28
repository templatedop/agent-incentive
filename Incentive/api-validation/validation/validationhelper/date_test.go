package validationhelper

import "testing"

func TestIsValidDateDDMMYY(t *testing.T) {
	cases := []struct {
		in string
		ok bool
	}{
		{"01/01/00", true},
		{"31/12/99", true},
		{"00/01/00", false},
		{"32/01/00", false},
		{"01/00/00", false},
		{"01/13/00", false},
		{"1/01/00", false},
		{"01/1/00", false},
		{"01/01/2000", false},
		{"", false},
	}
	for _, c := range cases {
		if got := IsValidDateDDMMYY(c.in); got != c.ok {
			t.Errorf("IsValidDateDDMMYY(%q)=%v want %v", c.in, got, c.ok)
		}
	}
}
