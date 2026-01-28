package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestNeValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Ne
		expectError bool
	}{
		{
			name:        "valid - both fields different",
			data:        test.Ne{Role: "user", Score: 10},
			expectError: false,
		},
		{
			name:        "invalid - role equals admin",
			data:        test.Ne{Role: "admin", Score: 10},
			expectError: true,
		},
		{
			name:        "invalid - score equals 0",
			data:        test.Ne{Role: "user", Score: 0},
			expectError: true,
		},
		{
			name:        "invalid - both match forbidden values",
			data:        test.Ne{Role: "admin", Score: 0},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateNe(&tt.data)
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
