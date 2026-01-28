package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestRequiredWithoutAllValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.RequiredWithoutAll
		expectError bool
	}{
		{
			name:        "valid - phone present, email optional",
			data:        test.RequiredWithoutAll{Phone: "555-1234", Fax: "", Email: ""},
			expectError: false,
		},
		{
			name:        "valid - both phone and fax present",
			data:        test.RequiredWithoutAll{Phone: "555-1234", Fax: "555-5678", Email: ""},
			expectError: false,
		},
		{
			name:        "invalid - both phone and fax missing, email required but missing",
			data:        test.RequiredWithoutAll{Phone: "", Fax: "", Email: ""},
			expectError: true,
		},
		{
			name:        "valid - both phone and fax missing but email present",
			data:        test.RequiredWithoutAll{Phone: "", Fax: "", Email: "test@example.com"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateRequiredWithoutAll(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
