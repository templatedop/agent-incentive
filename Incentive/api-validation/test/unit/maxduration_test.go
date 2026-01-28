package unit

import (
	"testing"
	"time"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestMaxDurationValidation(t *testing.T) {
	// Note: go-playground/validator doesn't have maxduration validator,
	// so we test only govalid behavior

	tests := []struct {
		name        string
		data        test.MaxDuration
		expectError bool
	}{
		{
			name:        "valid - below maximum",
			data:        test.MaxDuration{Interval: 12 * time.Hour},
			expectError: false,
		},
		{
			name:        "valid - exactly at maximum",
			data:        test.MaxDuration{Interval: 24 * time.Hour},
			expectError: false,
		},
		{
			name:        "invalid - above maximum",
			data:        test.MaxDuration{Interval: 48 * time.Hour},
			expectError: true,
		},
		{
			name:        "valid - zero",
			data:        test.MaxDuration{Interval: 0},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			govalidErr := test.ValidateMaxDuration(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
