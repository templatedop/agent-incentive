package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestNumberValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Number
		expectError bool
	}{
		{
			name:        "valid - positive integer",
			data:        test.Number{NumericString: "123"},
			expectError: false,
		},
		{
			name:        "valid - negative integer",
			data:        test.Number{NumericString: "-456"},
			expectError: false,
		},
		{
			name:        "valid - decimal",
			data:        test.Number{NumericString: "123.45"},
			expectError: false,
		},
		{
			name:        "valid - negative decimal",
			data:        test.Number{NumericString: "-67.89"},
			expectError: false,
		},
		{
			name:        "valid - zero",
			data:        test.Number{NumericString: "0"},
			expectError: false,
		},
		{
			name:        "invalid - letters",
			data:        test.Number{NumericString: "abc"},
			expectError: true,
		},
		{
			name:        "invalid - mixed",
			data:        test.Number{NumericString: "12abc"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateNumber(&tt.data)
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
