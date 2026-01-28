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

		// Phase 4: Commission history
		postgres.NewCommissionHistoryRepository,

		// Phase 5: Clawback
		postgres.NewClawbackRepository,

		// Phase 6: Suspense
		postgres.NewSuspenseRepository,

		// Future phases
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

		// Phase 4: Commission History (1 API)
		handler.NewCommissionHistoryHandler,

		// Phase 5: Clawback (1 API)
		handler.NewClawbackHandler,

		// Phase 6: Suspense (2 APIs)
		handler.NewSuspenseHandler,

		// Phase 7: Workflow Management (8 APIs)
		handler.NewWorkflowManagementHandler,

		// Future phases
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
