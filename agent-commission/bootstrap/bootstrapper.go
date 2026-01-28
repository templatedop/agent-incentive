package bootstrap

import (
	"agent-commission/handler"
	"agent-commission/repo/postgres"

	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	"go.uber.org/fx"
)

// FxRepo module provides all repository implementations
var FxRepo = fx.Module(
	"Repomodule",
	fx.Provide(
		// Reference data repositories (read-only, for lookups)
		postgres.NewReferenceDataRepository,

		// Commission module repositories
		postgres.NewCommissionRateRepository,
		// postgres.NewCommissionBatchRepository,
		// postgres.NewDisbursementRepository,
		// etc.
	),
)

// FxHandler module provides all HTTP handlers
var FxHandler = fx.Module(
	"Handlermodule",
	fx.Provide(
		// Phase 1: Commission Rate Configuration & Lookup
		handler.NewCommissionRateHandler,
		handler.NewLookupHandler,
		// Phase 2: Commission Batch Processing
		// Phase 3: Disbursement
		// etc.
	),
	fx.Invoke(func(
		commissionRateHandler *handler.CommissionRateHandler,
		lookupHandler *handler.LookupHandler,
	) {
		serverHandler.Register(commissionRateHandler)
		serverHandler.Register(lookupHandler)
	}),
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
