package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestIsDefaultValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.IsDefault
		expectError bool
	}{
		{
			name:        "valid - all default values",
			data:        test.IsDefault{OptionalField: "", OptionalNumber: 0},
			expectError: false,
		},
		{
			name:        "invalid - string has value",
			data:        test.IsDefault{OptionalField: "value", OptionalNumber: 0},
			expectError: true,
		},
		{
			name:        "invalid - number has value",
			data:        test.IsDefault{OptionalField: "", OptionalNumber: 42},
			expectError: true,
		},
		{
			name:        "invalid - both have values",
			data:        test.IsDefault{OptionalField: "test", OptionalNumber: 1},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateIsDefault(&tt.data)
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
