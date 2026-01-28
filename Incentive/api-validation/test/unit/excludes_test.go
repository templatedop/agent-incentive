package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludesValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Excludes
		expectError bool
	}{
		{
			name:        "valid - does not contain admin",
			data:        test.Excludes{Username: "alice"},
			expectError: false,
		},
		{
			name:        "valid - admin as part of larger word",
			data:        test.Excludes{Username: "administrator"},
			expectError: true, // excludes checks if substring exists
		},
		{
			name:        "invalid - exactly admin",
			data:        test.Excludes{Username: "admin"},
			expectError: true,
		},
		{
			name:        "invalid - contains admin substring",
			data:        test.Excludes{Username: "admin123"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateExcludes(&tt.data)
			govalidHasError := govalidErr != nil

			// Test go-playground/validator
			playgroundErr := validate.Struct(&tt.data)
			playgroundHasError := playgroundErr != nil

			// Compare results
			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
			if playgroundHasError != tt.expectError {
				t.Errorf("go-playground: expected error=%v, got error=%v (%v)", tt.expectError, playgroundHasError, playgroundErr)
			}
			if govalidHasError != playgroundHasError {
				t.Errorf("behavior mismatch: govalid=%v, playground=%v", govalidHasError, playgroundHasError)
			}
		})
	}
}
