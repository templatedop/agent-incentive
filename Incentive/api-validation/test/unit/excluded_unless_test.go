package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludedUnlessValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.ExcludedUnless
		expectError bool
	}{
		{
			name:        "valid - status active, field can have value",
			data:        test.ExcludedUnless{Status: "active", InactiveField: "value"},
			expectError: false,
		},
		{
			name:        "valid - status inactive and field empty",
			data:        test.ExcludedUnless{Status: "inactive", InactiveField: ""},
			expectError: false,
		},
		{
			name:        "invalid - status inactive but field has value",
			data:        test.ExcludedUnless{Status: "inactive", InactiveField: "value"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			govalidErr := test.ValidateExcludedUnless(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
