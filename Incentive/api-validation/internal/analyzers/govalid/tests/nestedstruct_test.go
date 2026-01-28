package tests

import (
	"testing"

	"github.com/gostaticanalysis/codegen/codegentest"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/govalid"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/markers"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/registry"
)

func TestNestedStruct(t *testing.T) {
	registry := registry.NewRegistry(
		registry.AddAnalyzers(markers.Initializer()),
		registry.AddGenerators(govalid.Initializer()),
	)

	if err := registry.Init(nil); err != nil {
		t.Fatalf("failed to initialize analyzers: %v", err)
	}

	govalid, err := registry.Generator(govalid.Name)
	if err != nil {
		t.Fatalf("failed to get govalid generator: %v", err)
	}

	testCases := []string{
		"nestedstruct/nop",
		"nestedstruct/inside",
		"nestedstruct/partial",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			results := codegentest.Run(t, codegentest.TestData(), govalid, tc)
			codegentest.Golden(t, results, update)
		})
	}
}
