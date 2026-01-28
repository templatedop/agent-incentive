package markers_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/markers"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()

	initializer := markers.Initializer()

	a, err := initializer.Init(nil)
	if err != nil {
		t.Fatalf("failed to initialize analyzer: %v", err)
	}

	analysistest.Run(t, testdata, a, "a")
}
