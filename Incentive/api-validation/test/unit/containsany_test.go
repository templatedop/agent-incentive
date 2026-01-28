package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestContainsAnyValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.ContainsAny
		expectError bool
	}{
		{
			name:        "valid - contains !",
			data:        test.ContainsAny{Password: "pass!word"},
			expectError: false,
		},
		{
			name:        "valid - contains @",
			data:        test.ContainsAny{Password: "p@ssword"},
			expectError: false,
		},
		{
			name:        "valid - contains #",
			data:        test.ContainsAny{Password: "pa#sword"},
			expectError: false,
		},
		{
			name:        "valid - contains $",
			data:        test.ContainsAny{Password: "pass$word"},
			expectError: false,
		},
		{
			name:        "invalid - no special characters",
			data:        test.ContainsAny{Password: "password"},
			expectError: true,
		},
		{
			name:        "invalid - different special characters",
			data:        test.ContainsAny{Password: "pass-word"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateContainsAny(&tt.data)
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
