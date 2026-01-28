package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestLatitudeValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Latitude
		expectError bool
	}{
		{
			name:        "valid - zero",
			data:        test.Latitude{Lat: "0"},
			expectError: false,
		},
		{
			name:        "valid - positive",
			data:        test.Latitude{Lat: "45.5"},
			expectError: false,
		},
		{
			name:        "valid - negative",
			data:        test.Latitude{Lat: "-45.5"},
			expectError: false,
		},
		{
			name:        "valid - max latitude",
			data:        test.Latitude{Lat: "90"},
			expectError: false,
		},
		{
			name:        "valid - min latitude",
			data:        test.Latitude{Lat: "-90"},
			expectError: false,
		},
		{
			name:        "invalid - above max",
			data:        test.Latitude{Lat: "90.1"},
			expectError: true,
		},
		{
			name:        "invalid - below min",
			data:        test.Latitude{Lat: "-90.1"},
			expectError: true,
		},
		{
			name:        "invalid - not a number",
			data:        test.Latitude{Lat: "abc"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateLatitude(&tt.data)
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
