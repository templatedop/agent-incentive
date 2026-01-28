package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludedWithoutAllValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.ExcludedWithoutAll
		expectError bool
	}{
		{
			name:        "valid - one feature present",
			data:        test.ExcludedWithoutAll{FeatureA: "enabled", FeatureB: "", ConflictingFeature: "enabled"},
			expectError: false,
		},
		{
			name:        "valid - both features missing, conflicting empty",
			data:        test.ExcludedWithoutAll{FeatureA: "", FeatureB: "", ConflictingFeature: ""},
			expectError: false,
		},
		{
			name:        "invalid - both features missing but conflicting present",
			data:        test.ExcludedWithoutAll{FeatureA: "", FeatureB: "", ConflictingFeature: "enabled"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			govalidErr := test.ValidateExcludedWithoutAll(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
