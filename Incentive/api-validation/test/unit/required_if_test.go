package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestRequiredIfValidation(t *testing.T) {
	// Note: go-playground/validator's required_if behavior may differ

	tests := []struct {
		name        string
		data        test.RequiredIf
		expectError bool
	}{
		{
			name:        "valid - status active and field has value",
			data:        test.RequiredIf{Status: "active", ActiveField: "value"},
			expectError: false,
		},
		{
			name:        "valid - status not active, field can be empty",
			data:        test.RequiredIf{Status: "inactive", ActiveField: ""},
			expectError: false,
		},
		{
			name:        "invalid - status active but field empty",
			data:        test.RequiredIf{Status: "active", ActiveField: ""},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateRequiredIf(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
