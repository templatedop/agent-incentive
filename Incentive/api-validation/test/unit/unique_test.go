package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestUniqueValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Unique
		expectError bool
	}{
		{
			name:        "valid - all unique strings",
			data:        test.Unique{Tags: []string{"go", "rust", "python"}, IDs: []int{1, 2, 3}},
			expectError: false,
		},
		{
			name:        "valid - all unique integers",
			data:        test.Unique{Tags: []string{"a", "b", "c"}, IDs: []int{10, 20, 30}},
			expectError: false,
		},
		{
			name:        "invalid - duplicate strings",
			data:        test.Unique{Tags: []string{"go", "rust", "go"}, IDs: []int{1, 2, 3}},
			expectError: true,
		},
		{
			name:        "invalid - duplicate integers",
			data:        test.Unique{Tags: []string{"a", "b", "c"}, IDs: []int{1, 2, 1}},
			expectError: true,
		},
		{
			name:        "valid - empty slices",
			data:        test.Unique{Tags: []string{}, IDs: []int{}},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateUnique(&tt.data)
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
