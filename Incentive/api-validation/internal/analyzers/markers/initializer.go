package markers

import (
	"sync"

	"golang.org/x/tools/go/analysis"

	// Local package import
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/registry"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/config"
)

// Initializer returns a new instance of the initializer for the markers analyzer.
func Initializer() registry.AnalyzerInitializer {
	return &initializer{}
}

// initializer is a struct that implements the registry.AnalyzerInitializer interface.
type initializer struct{}

var once sync.Once

// Init initializes the markers analyzer with the provided configuration.
func (i *initializer) Init(_ *config.GovalidConfig) (*analysis.Analyzer, error) {
	analyzer := newAnalyzer()

	once.Do(func() {
		Analyzer = analyzer
	})

	return Analyzer, nil
}

// Name returns the name of the markers analyzer.
func (i *initializer) Name() string {
	return Name
}
