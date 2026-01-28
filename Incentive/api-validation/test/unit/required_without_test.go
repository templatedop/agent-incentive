package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestRequiredWithoutValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.RequiredWithout
		expectError bool
	}{
		{
			name:        "valid - phone present, email optional",
			data:        test.RequiredWithout{Phone: "555-1234", Email: ""},
			expectError: false,
		},
		{
			name:        "valid - both present",
			data:        test.RequiredWithout{Phone: "555-1234", Email: "test@example.com"},
			expectError: false,
		},
		{
			name:        "invalid - phone missing, email required but missing",
			data:        test.RequiredWithout{Phone: "", Email: ""},
			expectError: true,
		},
		{
			name:        "valid - phone missing but email present",
			data:        test.RequiredWithout{Phone: "", Email: "test@example.com"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateRequiredWithout(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
