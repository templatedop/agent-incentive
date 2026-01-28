package bootstrap

import (
	"agent-commission/handler"
	"agent-commission/repo/postgres"
	"agent-commission/workflows"
	"agent-commission/workflows/activities"

	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	"go.uber.org/fx"
)

// FxRepo module provides all repository implementations
// Repositories will be added as we implement each module
var FxRepo = fx.Module(
	"Repomodule",
	fx.Provide(
		// Phase 1: Agent Onboarding & Profile Management
		postgres.NewAgentProfileRepository,
		postgres.NewAgentHierarchyRepository,
		postgres.NewReferenceDataRepository,

		// Phase 2: License Management (to be added)
		// postgres.NewAgentLicenseRepository,

		// Phase 3: Commission Rate Configuration (to be added)
		// postgres.NewCommissionRateRepository,

		// Phase 4+: Additional repositories...
	),
)

// FxHandler module provides all HTTP handlers
// Each handler must be annotated to implement serverHandler.Handler interface
var FxHandler = fx.Module(
	"Handlermodule",
	fx.Provide(
		// Phase 1: Agent Onboarding & Profile Management
		fx.Annotate(
			handler.NewAgentOnboardingHandler,
			fx.As(new(serverHandler.Handler)),
			fx.ResultTags(serverHandler.ServerControllersGroupTag),
		),

		// Phase 2: License Management (to be added)
		// fx.Annotate(
		// 	handler.NewAgentLicenseHandler,
		// 	fx.As(new(serverHandler.Handler)),
		// 	fx.ResultTags(serverHandler.ServerControllersGroupTag),
		// ),

		// Phase 3: Commission Rate Configuration (to be added)
		// fx.Annotate(
		// 	handler.NewCommissionRateHandler,
		// 	fx.As(new(serverHandler.Handler)),
		// 	fx.ResultTags(serverHandler.ServerControllersGroupTag),
		// ),

		// Phase 4-10: Additional handlers...
	),
)

// FxWorkflow module provides Temporal workflow workers and activities
var FxWorkflow = fx.Module(
	"Workflowmodule",
	fx.Provide(
		// Temporal client
		workflows.NewTemporalClient,

		// Activities
		activities.NewAgentActivities,

		// Workers
		workflows.NewTemporalWorker,
	),
	fx.Invoke(
		// Start workers
		workflows.StartWorker,
	),
)
