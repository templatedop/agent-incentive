package unit

import (
	"testing"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/test"
)

func TestExcludedWithoutValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.ExcludedWithout
		expectError bool
	}{
		{
			name:        "valid - premium true, free feature allowed",
			data:        test.ExcludedWithout{Premium: "premium", FreeFeature: "enabled"},
			expectError: false,
		},
		{
			name:        "valid - premium false, free feature empty",
			data:        test.ExcludedWithout{Premium: "", FreeFeature: ""},
			expectError: false,
		},
		{
			name:        "invalid - premium false but free feature present",
			data:        test.ExcludedWithout{Premium: "", FreeFeature: "enabled"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			govalidErr := test.ValidateExcludedWithout(&tt.data)
			govalidHasError := govalidErr != nil

			if govalidHasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (%v)", tt.expectError, govalidHasError, govalidErr)
			}
		})
	}
}
