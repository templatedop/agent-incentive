package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestURIValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.URI
		expectError bool
	}{
		{
			name:        "valid - http URL",
			data:        test.URI{Address: "http://example.com"},
			expectError: false,
		},
		{
			name:        "valid - https URL",
			data:        test.URI{Address: "https://example.com/path"},
			expectError: false,
		},
		{
			name:        "valid - ftp URI",
			data:        test.URI{Address: "ftp://ftp.example.com"},
			expectError: false,
		},
		{
			name:        "valid - file URI",
			data:        test.URI{Address: "file:///path/to/file"},
			expectError: false,
		},
		{
			name:        "invalid - not a URI",
			data:        test.URI{Address: "not a uri"},
			expectError: true,
		},
		{
			name:        "invalid - missing scheme",
			data:        test.URI{Address: "example.com"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateURI(&tt.data)
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
