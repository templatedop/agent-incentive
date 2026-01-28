package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludedIfValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.ExcludedIf
		expectError bool
	}{
		{
			name:        "valid - status inactive, field can have value",
			data:        test.ExcludedIf{Status: "active", InactiveField: "value"},
			expectError: false,
		},
		{
			name:        "valid - status inactive and field empty",
			data:        test.ExcludedIf{Status: "inactive", InactiveField: ""},
			expectError: false,
		},
		{
			name:        "invalid - status inactive but field has value",
			data:        test.ExcludedIf{Status: "inactive", InactiveField: "value"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			govalidErr := test.ValidateExcludedIf(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
