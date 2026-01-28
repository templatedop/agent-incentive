package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestIsColourValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.IsColour
		expectError bool
	}{
		{
			name:        "valid - hex 3 digits",
			data:        test.IsColour{Color: "#FFF"},
			expectError: false,
		},
		{
			name:        "valid - hex 6 digits",
			data:        test.IsColour{Color: "#FF5733"},
			expectError: false,
		},
		{
			name:        "valid - hex 8 digits (with alpha)",
			data:        test.IsColour{Color: "#FF5733AA"},
			expectError: false,
		},
		{
			name:        "valid - rgb",
			data:        test.IsColour{Color: "rgb(255, 87, 51)"},
			expectError: false,
		},
		{
			name:        "valid - rgba",
			data:        test.IsColour{Color: "rgba(255, 87, 51, 0.5)"},
			expectError: false,
		},
		{
			name:        "valid - hsl",
			data:        test.IsColour{Color: "hsl(12, 100%, 60%)"},
			expectError: false,
		},
		{
			name:        "valid - hsla",
			data:        test.IsColour{Color: "hsla(12, 100%, 60%, 0.8)"},
			expectError: false,
		},
		{
			name:        "valid - named color",
			data:        test.IsColour{Color: "red"},
			expectError: false,
		},
		{
			name:        "invalid - not a color",
			data:        test.IsColour{Color: "notacolor"},
			expectError: true,
		},
		{
			name:        "invalid - invalid hex",
			data:        test.IsColour{Color: "#GGGGGG"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateIsColour(&tt.data)
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
