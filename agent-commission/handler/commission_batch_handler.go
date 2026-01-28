package handler

import (
	"errors"
	"time"

	"agent-commission/core/domain"
	"agent-commission/core/port"
	resp "agent-commission/handler/response"
	"agent-commission/repo/postgres"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// CommissionBatchHandler handles commission batch processing operations
// Implements: FR-IC-COM-002 (Commission batch processing)
type CommissionBatchHandler struct {
	*serverHandler.Base
	batchRepo *postgres.CommissionBatchRepository
}

// NewCommissionBatchHandler creates a new commission batch handler
func NewCommissionBatchHandler(batchRepo *postgres.CommissionBatchRepository) *CommissionBatchHandler {
	base := serverHandler.New("CommissionBatch").SetPrefix("/v1").AddPrefix("")
	return &CommissionBatchHandler{Base: base, batchRepo: batchRepo}
}

// Routes returns all routes for commission batch endpoints
func (h *CommissionBatchHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.POST("/commissions/batches/start", h.StartCommissionBatch).Name("Start Commission Batch"),
		serverRoute.GET("/commissions/batches/:batchId/status", h.GetCommissionBatchStatus).Name("Get Commission Batch Status"),
	}
}

// StartCommissionBatchRequest represents the request for POST /commissions/batches/start
type StartCommissionBatchRequest struct {
	Month         int    `json:"month" validate:"required,min=1,max=12"`
	Year          int    `json:"year" validate:"required,min=2020,max=2099"`
	TriggeredBy   string `json:"triggered_by" validate:"omitempty,oneof=SYSTEM_SCHEDULER MANUAL"`
	ManualTrigger bool   `json:"manual_trigger"`
}

// StartCommissionBatch initiates a monthly commission calculation batch
// Implements: FR-IC-COM-002, BR-IC-COM-001, BR-IC-COM-012 (6-hour SLA)
// API: POST /commissions/batches/start
//
// Request Body:
//   - month (required): Month (1-12)
//   - year (required): Year (2020-2099)
//   - triggered_by (optional): SYSTEM_SCHEDULER or MANUAL
//   - manual_trigger (optional): Boolean flag
//
// Returns:
//   - 200: Batch initiated successfully
//   - 400: Bad request - batch already running for this month
//   - 401: Unauthorized
//   - 403: Forbidden
//
// Business Rules:
//   - BR-IC-COM-001: Commission calculation rules
//   - BR-IC-COM-012: 6-hour SLA timeout
//   - Only one active batch per month/year
func (h *CommissionBatchHandler) StartCommissionBatch(
	sctx *serverRoute.Context,
	req StartCommissionBatchRequest,
) (*resp.StartCommissionBatchResponse, error) {
	// Check if batch already exists for this month/year
	existingBatch, err := h.batchRepo.GetBatchByMonthYear(sctx.Ctx, req.Month, req.Year)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error(sctx.Ctx, "Error checking existing batch: %v", err)
		return nil, err
	}

	// If batch exists and is not completed/failed, return error
	if existingBatch != nil {
		if existingBatch.Status != domain.BatchStatusCompleted && existingBatch.Status != domain.BatchStatusFailed {
			log.Warn(sctx.Ctx, "Batch already running for %d-%d: %s", req.Year, req.Month, existingBatch.BatchID)
			return nil, errors.New("commission batch already running for this month")
		}
	}

	// Generate batch ID
	batchID := postgres.GenerateBatchID(req.Month, req.Year)

	// Determine triggered_by
	triggeredBy := "MANUAL"
	if req.TriggeredBy != "" {
		triggeredBy = req.TriggeredBy
	}

	// Create batch
	now := time.Now()
	batch := &domain.CommissionBatch{
		BatchID:            batchID,
		Month:              req.Month,
		Year:               req.Year,
		Status:             domain.BatchStatusInitiated,
		WorkflowID:         batchID, // TODO: Replace with Temporal workflow ID
		WorkflowState:      "INITIATED",
		TotalPolicies:      0,
		ProcessedRecords:   0,
		FailedRecords:      0,
		ProgressPercentage: 0,
		TriggeredBy:        triggeredBy,
		StartedAt:          now,
		SLADeadline:        now.Add(6 * time.Hour), // BR-IC-COM-012: 6-hour SLA
		EstimatedCompletion: nil,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Save batch to database
	err = h.batchRepo.CreateBatch(sctx.Ctx, batch)
	if err != nil {
		log.Error(sctx.Ctx, "Error creating commission batch: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Commission batch initiated: %s for %d-%02d", batchID, req.Year, req.Month)

	// TODO: Start Temporal workflow here
	// workflowID, err := h.temporalClient.StartCommissionCalculationWorkflow(ctx, batch)

	// Build response
	response := &resp.StartCommissionBatchResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Data:                 resp.NewCommissionBatchResponse(*batch),
	}

	return response, nil
}

// BatchIDUri represents the URI parameter for batch ID
type BatchIDUri struct {
	BatchID string `uri:"batchId" validate:"required"`
}

// GetCommissionBatchStatus retrieves current status and progress of commission batch
// Implements: WF-IC-COM-001, WF-TEMPORAL-IC-001
// API: GET /commissions/batches/{batchId}/status
//
// URI Parameters:
//   - batchId (required): Batch ID
//
// Returns:
//   - 200: Batch status retrieved
//   - 401: Unauthorized
//   - 404: Batch not found
func (h *CommissionBatchHandler) GetCommissionBatchStatus(
	sctx *serverRoute.Context,
	uri BatchIDUri,
) (*resp.GetCommissionBatchStatusResponse, error) {
	// Fetch batch from database
	batch, err := h.batchRepo.GetBatchByID(sctx.Ctx, uri.BatchID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn(sctx.Ctx, "Batch not found: %s", uri.BatchID)
			return nil, errors.New("batch not found")
		}
		log.Error(sctx.Ctx, "Error fetching batch status: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Batch status retrieved: %s - %s", uri.BatchID, batch.Status)

	// TODO: Query Temporal workflow for real-time status
	// workflowStatus, err := h.temporalClient.GetWorkflowStatus(ctx, batch.WorkflowID)

	// Build response
	response := &resp.GetCommissionBatchStatusResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data:                 resp.NewCommissionBatchStatusResponse(*batch),
	}

	return response, nil
}
