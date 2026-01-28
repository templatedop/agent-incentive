package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludedWithValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.ExcludedWith
		expectError bool
	}{
		{
			name:        "valid - guest mode off, admin panel allowed",
			data:        test.ExcludedWith{GuestMode: "", AdminPanel: "enabled"},
			expectError: false,
		},
		{
			name:        "valid - guest mode on, admin panel empty",
			data:        test.ExcludedWith{GuestMode: "on", AdminPanel: ""},
			expectError: false,
		},
		{
			name:        "invalid - guest mode on but admin panel present",
			data:        test.ExcludedWith{GuestMode: "on", AdminPanel: "enabled"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			govalidErr := test.ValidateExcludedWith(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
