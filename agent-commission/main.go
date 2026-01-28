package main

import (
	"context"

	"agent-commission/bootstrap"

	bootstrapper "gitlab.cept.gov.in/it-2.0-common/n-api-bootstrapper"
)

// main is the application entry point
// It initializes and starts the server with all dependencies using Uber FX
func main() {
	app := bootstrapper.New().Options(
		bootstrap.FxHandler,  // Register all HTTP handlers
		bootstrap.FxRepo,     // Register all repositories
		bootstrap.FxWorkflow, // Temporal workflow workers
	)
	app.WithContext(context.Background()).Run()
}

// The bootstrapper automatically handles:
// - Configuration loading from configs/ directory
// - Database connection pool (PostgreSQL via pgx)
// - HTTP server initialization (Gin framework)
// - Graceful shutdown and signal handling
// - Dependency injection via Uber FX
// - Health check endpoints
// - Metrics and observability
// - Request logging and tracing
