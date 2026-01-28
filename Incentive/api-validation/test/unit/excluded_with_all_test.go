package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludedWithAllValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.ExcludedWithAll
		expectError bool
	}{
		{
			name:        "valid - only one condition true",
			data:        test.ExcludedWithAll{ReadOnly: "true", Archived: "", EditButton: "enabled"},
			expectError: false,
		},
		{
			name:        "valid - both conditions true, edit button empty",
			data:        test.ExcludedWithAll{ReadOnly: "true", Archived: "true", EditButton: ""},
			expectError: false,
		},
		{
			name:        "invalid - both conditions true but edit button present",
			data:        test.ExcludedWithAll{ReadOnly: "true", Archived: "true", EditButton: "enabled"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			govalidErr := test.ValidateExcludedWithAll(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
