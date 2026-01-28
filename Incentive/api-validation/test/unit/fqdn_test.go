package unit

import (
	"testing"

	"github.com/go-playground/validator/v10"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestFQDNValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		data        test.FQDN
		expectError bool
	}{
		{
			name:        "valid - simple FQDN",
			data:        test.FQDN{Domain: "example.com"},
			expectError: false,
		},
		{
			name:        "valid - subdomain",
			data:        test.FQDN{Domain: "www.example.com"},
			expectError: false,
		},
		{
			name:        "valid - multiple subdomains",
			data:        test.FQDN{Domain: "api.v1.example.com"},
			expectError: false,
		},
		{
			name:        "valid - with trailing dot",
			data:        test.FQDN{Domain: "example.com."},
			expectError: false,
		},
		{
			name:        "invalid - no dot",
			data:        test.FQDN{Domain: "localhost"},
			expectError: true,
		},
		{
			name:        "invalid - starts with dot",
			data:        test.FQDN{Domain: ".example.com"},
			expectError: true,
		},
		{
			name:        "invalid - ends with hyphen",
			data:        test.FQDN{Domain: "example-.com"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateFQDN(&tt.data)
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
