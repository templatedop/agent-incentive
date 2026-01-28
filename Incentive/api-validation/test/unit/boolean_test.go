package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestBooleanValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Boolean
		expectError bool
	}{
		{
			name:        "valid - true",
			data:        test.Boolean{Flag: "true"},
			expectError: false,
		},
		{
			name:        "valid - false",
			data:        test.Boolean{Flag: "false"},
			expectError: false,
		},
		{
			name:        "valid - 1",
			data:        test.Boolean{Flag: "1"},
			expectError: false,
		},
		{
			name:        "valid - 0",
			data:        test.Boolean{Flag: "0"},
			expectError: false,
		},
		{
			name:        "valid - yes",
			data:        test.Boolean{Flag: "yes"},
			expectError: false,
		},
		{
			name:        "valid - no",
			data:        test.Boolean{Flag: "no"},
			expectError: false,
		},
		{
			name:        "invalid - random string",
			data:        test.Boolean{Flag: "maybe"},
			expectError: true,
		},
		{
			name:        "invalid - number",
			data:        test.Boolean{Flag: "42"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateBoolean(&tt.data)
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
