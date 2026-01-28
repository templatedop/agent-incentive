package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestOneOfValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.OneOf
		expectError bool
	}{
		{
			name:        "valid - red and 1",
			data:        test.OneOf{Color: "red", Level: 1},
			expectError: false,
		},
		{
			name:        "valid - green and 2",
			data:        test.OneOf{Color: "green", Level: 2},
			expectError: false,
		},
		{
			name:        "valid - blue and 3",
			data:        test.OneOf{Color: "blue", Level: 3},
			expectError: false,
		},
		{
			name:        "invalid - color not in list",
			data:        test.OneOf{Color: "yellow", Level: 1},
			expectError: true,
		},
		{
			name:        "invalid - level not in list",
			data:        test.OneOf{Color: "red", Level: 4},
			expectError: true,
		},
		{
			name:        "invalid - both not in list",
			data:        test.OneOf{Color: "purple", Level: 5},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateOneOf(&tt.data)
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
