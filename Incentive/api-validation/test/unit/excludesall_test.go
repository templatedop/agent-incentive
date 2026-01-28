package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludesAllValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.ExcludesAll
		expectError bool
	}{
		{
			name:        "valid - no forbidden characters",
			data:        test.ExcludesAll{Comment: "This is a safe comment"},
			expectError: false,
		},
		{
			name:        "invalid - contains <",
			data:        test.ExcludesAll{Comment: "Hello <world"},
			expectError: true,
		},
		{
			name:        "invalid - contains >",
			data:        test.ExcludesAll{Comment: "Hello world>"},
			expectError: true,
		},
		{
			name:        "invalid - contains both < and >",
			data:        test.ExcludesAll{Comment: "<script>alert()</script>"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateExcludesAll(&tt.data)
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
