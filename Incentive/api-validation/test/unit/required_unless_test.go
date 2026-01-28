package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestRequiredUnlessValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.RequiredUnless
		expectError bool
	}{
		{
			name:        "valid - status inactive, field can be empty",
			data:        test.RequiredUnless{Status: "inactive", ActiveField: ""},
			expectError: false,
		},
		{
			name:        "valid - status active and field has value",
			data:        test.RequiredUnless{Status: "active", ActiveField: "value"},
			expectError: false,
		},
		{
			name:        "invalid - status active but field empty",
			data:        test.RequiredUnless{Status: "active", ActiveField: ""},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateRequiredUnless(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
