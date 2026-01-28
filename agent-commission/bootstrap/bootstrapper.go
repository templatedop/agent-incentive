package bootstrap

import (
	"agent-commission/handler"
	"agent-commission/repo/postgres"

	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	"go.uber.org/fx"
)

// FxRepo module provides all repository implementations
// Repositories will be added as we implement each module
var FxRepo = fx.Module(
	"Repomodule",
	fx.Provide(
		// Agent repositories will be added here
		// Example: postgres.NewAgentProfileRepository,
		// Example: postgres.NewAgentLicenseRepository,
		// Example: postgres.NewCommissionRepository,
	),
)

// FxHandler module provides all HTTP handlers
// Each handler must be annotated to implement serverHandler.Handler interface
var FxHandler = fx.Module(
	"Handlermodule",
	fx.Provide(
		// Handlers will be added as we implement each module
		// Phase 1: Agent Onboarding & Profile Management
		// fx.Annotate(
		// 	handler.NewAgentOnboardingHandler,
		// 	fx.As(new(serverHandler.Handler)),
		// 	fx.ResultTags(serverHandler.ServerControllersGroupTag),
		// ),
		
		// Phase 2: License Management
		// fx.Annotate(
		// 	handler.NewAgentLicenseHandler,
		// 	fx.As(new(serverHandler.Handler)),
		// 	fx.ResultTags(serverHandler.ServerControllersGroupTag),
		// ),
		
		// Phase 3: Commission Rate Configuration
		// fx.Annotate(
		// 	handler.NewCommissionRateHandler,
		// 	fx.As(new(serverHandler.Handler)),
		// 	fx.ResultTags(serverHandler.ServerControllersGroupTag),
		// ),
		
		// Phase 4-10: Additional handlers...
	),
)

// FxWorkflow module provides Temporal workflow workers (to be added)
// This will be populated when we implement Temporal workflows
var FxWorkflow = fx.Module(
	"Workflowmodule",
	fx.Provide(
		// Temporal workflow workers will be added here
		// Example: workflows.NewAgentOnboardingWorker,
		// Example: workflows.NewCommissionBatchWorker,
	),
)
