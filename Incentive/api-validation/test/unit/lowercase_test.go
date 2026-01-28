package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestLowercaseValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Lowercase
		expectError bool
	}{
		{
			name:        "valid - all lowercase",
			data:        test.Lowercase{Username: "alice"},
			expectError: false,
		},
		{
			name:        "valid - lowercase with numbers",
			data:        test.Lowercase{Username: "user123"},
			expectError: false,
		},
		{
			name:        "invalid - has uppercase",
			data:        test.Lowercase{Username: "Alice"},
			expectError: true,
		},
		{
			name:        "invalid - all uppercase",
			data:        test.Lowercase{Username: "ADMIN"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateLowercase(&tt.data)
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
