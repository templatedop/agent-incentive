// Package main is the entry point for the govalid command line tool.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gostaticanalysis/codegen/singlegenerator"

	// Local package import
	govalid_pkg "gitlab.cept.gov.in/it-2.0-common/n-api-validation"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/govalid"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/markers"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/registry"
)

func main() {
	// Parse version flag
	var version bool

	flag.BoolVar(&version, "version", false, "print version information")
	flag.Parse()

	if version {
		fmt.Printf("govalid version %s\n", govalid_pkg.Version)
		os.Exit(0)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

// run initializes the analyzers and starts the unit checker.
func run() error {
	registry := registry.NewRegistry(
		registry.AddAnalyzers(markers.Initializer()),
		registry.AddGenerators(govalid.Initializer()),
	)

	if err := registry.Init(nil); err != nil {
		return fmt.Errorf("failed to initialize analyzers: %w", err)
	}

	govalid, err := registry.Generator(govalid.Name)
	if err != nil {
		return fmt.Errorf("failed to get govalid generator: %w", err)
	}

	singlegenerator.Main(govalid)

	return nil
}
