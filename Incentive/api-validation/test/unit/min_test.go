package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestMinValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Min
		expectError bool
	}{
		{
			name:        "valid - above minimum",
			data:        test.Min{Age: 20},
			expectError: false,
		},
		{
			name:        "limit minus one",
			data:        test.Min{Age: 9},
			expectError: true,
		},
		{
			name:        "exactly at limit",
			data:        test.Min{Age: 10},
			expectError: false,
		},
		{
			name:        "limit plus one",
			data:        test.Min{Age: 11},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateMin(&tt.data)
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
