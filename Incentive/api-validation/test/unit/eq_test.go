package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestEqValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Eq
		expectError bool
	}{
		{
			name:        "valid - both fields match",
			data:        test.Eq{Status: "active", Count: 100},
			expectError: false,
		},
		{
			name:        "invalid - status mismatch",
			data:        test.Eq{Status: "inactive", Count: 100},
			expectError: true,
		},
		{
			name:        "invalid - count mismatch",
			data:        test.Eq{Status: "active", Count: 99},
			expectError: true,
		},
		{
			name:        "invalid - both mismatch",
			data:        test.Eq{Status: "pending", Count: 50},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateEq(&tt.data)
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
