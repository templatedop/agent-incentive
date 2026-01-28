package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestRequiredWithValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.RequiredWith
		expectError bool
	}{
		{
			name:        "valid - both fields have values",
			data:        test.RequiredWith{Email: "test@example.com", EmailConfirmation: "test@example.com"},
			expectError: false,
		},
		{
			name:        "valid - both fields empty",
			data:        test.RequiredWith{Email: "", EmailConfirmation: ""},
			expectError: false,
		},
		{
			name:        "invalid - email present but confirmation missing",
			data:        test.RequiredWith{Email: "test@example.com", EmailConfirmation: ""},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateRequiredWith(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
