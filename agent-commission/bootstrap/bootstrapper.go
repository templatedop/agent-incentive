package bootstrap

import (
	"agent-commission/handler"
	"agent-commission/repo/postgres"

	"go.uber.org/fx"
)

// FxRepo module provides all repository implementations
var FxRepo = fx.Module(
	"Repomodule",
	fx.Provide(
		// Phase 1: Reference data repositories
		postgres.NewReferenceDataRepository,
		postgres.NewCommissionRateRepository,

		// Phase 2: Commission processing repositories
		postgres.NewCommissionBatchRepository,
		postgres.NewTrialStatementRepository,
		postgres.NewFinalStatementRepository,

		// Phase 3: Disbursement repositories
		postgres.NewDisbursementRepository,

		// Future phases
		// Phase 5: postgres.NewClawbackRepository,
		// Phase 6: postgres.NewSuspenseRepository,
		// etc.
	),
)

// FxHandler module provides all HTTP handlers
var FxHandler = fx.Module(
	"Handlermodule",
	fx.Provide(
		// Phase 1: Commission Rate Configuration & Lookup (3 APIs)
		handler.NewCommissionRateHandler,
		handler.NewLookupHandler,

		// Phase 2: Commission Processing (5 APIs)
		handler.NewCommissionBatchHandler,
		handler.NewTrialStatementHandler,
		handler.NewFinalStatementHandler,

		// Phase 3: Disbursement & Webhooks (4 APIs)
		handler.NewDisbursementHandler,
		handler.NewWebhookHandler,

		// Future phases
		// Phase 4: Commission History
		// Phase 5: Clawback
		// Phase 6: Suspense
		// etc.
	),
)

// FxWorkflow module provides Temporal workflow workers and activities
// Will be enabled when we implement commission workflows
var FxWorkflow = fx.Module(
	"Workflowmodule",
	fx.Provide(
	// Temporal workflows will be added here
	// workflows.NewTemporalClient,
	// activities.NewCommissionActivities,
	// workflows.NewCommissionWorker,
	),
)
