package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestRequiredWithAllValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.RequiredWithAll
		expectError bool
	}{
		{
			name:        "valid - all fields present",
			data:        test.RequiredWithAll{FirstName: "John", LastName: "Doe", FullName: "John Doe"},
			expectError: false,
		},
		{
			name:        "valid - only first name present",
			data:        test.RequiredWithAll{FirstName: "John", LastName: "", FullName: ""},
			expectError: false,
		},
		{
			name:        "valid - all fields empty",
			data:        test.RequiredWithAll{FirstName: "", LastName: "", FullName: ""},
			expectError: false,
		},
		{
			name:        "invalid - first and last present but full name missing",
			data:        test.RequiredWithAll{FirstName: "John", LastName: "Doe", FullName: ""},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateRequiredWithAll(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
