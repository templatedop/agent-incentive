package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestLongitudeValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.Longitude
		expectError bool
	}{
		{
			name:        "valid - zero",
			data:        test.Longitude{Lon: "0"},
			expectError: false,
		},
		{
			name:        "valid - positive",
			data:        test.Longitude{Lon: "120.5"},
			expectError: false,
		},
		{
			name:        "valid - negative",
			data:        test.Longitude{Lon: "-120.5"},
			expectError: false,
		},
		{
			name:        "valid - max longitude",
			data:        test.Longitude{Lon: "180"},
			expectError: false,
		},
		{
			name:        "valid - min longitude",
			data:        test.Longitude{Lon: "-180"},
			expectError: false,
		},
		{
			name:        "invalid - above max",
			data:        test.Longitude{Lon: "180.1"},
			expectError: true,
		},
		{
			name:        "invalid - below min",
			data:        test.Longitude{Lon: "-180.1"},
			expectError: true,
		},
		{
			name:        "invalid - not a number",
			data:        test.Longitude{Lon: "xyz"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateLongitude(&tt.data)
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
