package unit

import (
	"testing"
	"time"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestMinDurationValidation(t *testing.T) {
	// Note: go-playground/validator doesn't have minduration validator,
	// so we test only govalid behavior

	tests := []struct {
		name        string
		data        test.MinDuration
		expectError bool
	}{
		{
			name:        "valid - above minimum",
			data:        test.MinDuration{Timeout: 2 * time.Hour},
			expectError: false,
		},
		{
			name:        "valid - exactly at minimum",
			data:        test.MinDuration{Timeout: 1 * time.Hour},
			expectError: false,
		},
		{
			name:        "invalid - below minimum",
			data:        test.MinDuration{Timeout: 30 * time.Minute},
			expectError: true,
		},
		{
			name:        "invalid - zero",
			data:        test.MinDuration{Timeout: 0},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateMinDuration(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
