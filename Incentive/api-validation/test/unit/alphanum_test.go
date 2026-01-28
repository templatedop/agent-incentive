package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestAlphanumValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Alphanum
		expectError bool
	}{
		{
			name:        "valid - letters only",
			data:        test.Alphanum{Code: "ABC"},
			expectError: false,
		},
		{
			name:        "valid - numbers only",
			data:        test.Alphanum{Code: "123"},
			expectError: false,
		},
		{
			name:        "valid - mixed alphanumeric",
			data:        test.Alphanum{Code: "ABC123"},
			expectError: false,
		},
		{
			name:        "invalid - has special characters",
			data:        test.Alphanum{Code: "ABC-123"},
			expectError: true,
		},
		{
			name:        "invalid - has space",
			data:        test.Alphanum{Code: "ABC 123"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateAlphanum(&tt.data)
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
