package handler

import (
	"fmt"

	"agent-commission/core/domain"
	"agent-commission/core/port"
	"agent-commission/handler/response"
	repo "agent-commission/repo/postgres"

	log "gitlab.cept.gov.in/it-2.0-common/api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// SuspenseHandler handles suspense account operations
// FR-IC-COM-013: Suspense account management
// BR-IC-SUS-001: Commission moved to suspense if agent details incomplete/invalid
type SuspenseHandler struct {
	*serverHandler.Base
	suspenseRepo *repo.SuspenseRepository
}

// NewSuspenseHandler creates a new suspense handler
func NewSuspenseHandler(
	suspenseRepo *repo.SuspenseRepository,
) *SuspenseHandler {
	base := serverHandler.New("Suspense").SetPrefix("/v1").AddPrefix("")
	return &SuspenseHandler{
		Base:         base,
		suspenseRepo: suspenseRepo,
	}
}

// Routes returns all routes for suspense endpoints
func (h *SuspenseHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.POST("/commissions/suspense/create", h.CreateSuspenseAccount).Name("Create Suspense Account"),
		serverRoute.GET("/commissions/suspense/aging-report", h.GetAgingReport).Name("Get Suspense Aging Report"),
	}
}

// CreateSuspenseAccountRequest represents request to create a suspense account entry
type CreateSuspenseAccountRequest struct {
	AgentID        *string                `json:"agent_id,omitempty"`        // May be null if agent not found
	PolicyNumber   *string                `json:"policy_number,omitempty"`   // Source policy
	CommissionID   *int64                 `json:"commission_id,omitempty"`   // Source commission transaction
	DisbursementID *int64                 `json:"disbursement_id,omitempty"` // Failed disbursement
	SuspenseAmount float64                `json:"suspense_amount" validate:"required,gt=0"`
	SuspenseReason domain.SuspenseReason  `json:"suspense_reason" validate:"required"`
	Priority       *string                `json:"priority,omitempty"` // HIGH, MEDIUM, LOW (auto-calculated if not provided)
	AssignedTo     *string                `json:"assigned_to,omitempty"`
	Notes          *string                `json:"notes,omitempty"`
}

// CreateSuspenseAccount creates a new suspense account entry
// API: POST /commissions/suspense/create
// Implements: FR-IC-COM-013, BR-IC-SUS-001
// Used when commission cannot be paid due to incomplete/invalid agent details
func (h *SuspenseHandler) CreateSuspenseAccount(
	sctx *serverRoute.Context,
	req CreateSuspenseAccountRequest,
) (*response.CreateSuspenseAccountResponse, error) {
	log.Info(sctx.Ctx, "Creating suspense account entry: amount=%.2f, reason=%s",
		req.SuspenseAmount, req.SuspenseReason)

	// Create suspense account record
	suspense := &domain.SuspenseAccount{
		AgentID:        req.AgentID,
		PolicyNumber:   req.PolicyNumber,
		CommissionID:   req.CommissionID,
		DisbursementID: req.DisbursementID,
		SuspenseAmount: req.SuspenseAmount,
		SuspenseReason: req.SuspenseReason,
		AssignedTo:     req.AssignedTo,
		Notes:          req.Notes,
		CreatedBy:      "system", // TODO: Get from JWT
		UpdatedBy:      "system",
	}

	// Set priority if provided, otherwise it will be auto-determined
	if req.Priority != nil {
		suspense.Priority = *req.Priority
	}

	// Repository will auto-calculate:
	// - Priority (if not set)
	// - Resolution deadline based on priority
	// - Initial aging days
	err := h.suspenseRepo.CreateSuspenseAccount(sctx.Ctx, suspense)
	if err != nil {
		log.Error(sctx.Ctx, "Error creating suspense account: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Suspense account created successfully: %d (priority: %s, deadline: %v)",
		suspense.SuspenseID,
		suspense.Priority,
		suspense.ResolutionDeadline,
	)

	// Start Temporal workflow for suspense resolution
	// Workflow will handle retry logic and aging tracking
	agentIDStr := "UNKNOWN"
	if suspense.AgentID != nil {
		agentIDStr = *suspense.AgentID
	}
	workflowID := fmt.Sprintf("suspense-%d-%s-%s", suspense.SuspenseID, agentIDStr, suspense.SuspenseReason)
	suspense.WorkflowID = &workflowID
	suspense.WorkflowState = ptrString("INITIATED")

	// TODO: Uncomment when Temporal client is available
	// err = h.temporalClient.ExecuteWorkflow(sctx.Ctx, client.StartWorkflowOptions{
	//     ID:        workflowID,
	//     TaskQueue: "suspense-resolution-queue",
	// }, "SuspenseResolutionWorkflow", suspense)
	// if err != nil {
	//     log.Error(sctx.Ctx, "Error starting suspense workflow: %v", err)
	//     return nil, err
	// }

	log.Info(sctx.Ctx, "Suspense workflow ID: %s", workflowID)

	// TODO: Send notification to assigned user
	// if suspense.AssignedTo != nil {
	//     h.notificationService.NotifySuspenseAssignment(*suspense.AssignedTo, suspense)
	// }

	// Build response
	resp := &response.CreateSuspenseAccountResponse{
		StatusCodeAndMessage: port.CreateSuccess,
	}
	resp.Data.SuspenseID = suspense.SuspenseID
	resp.Data.AgentID = suspense.AgentID
	resp.Data.PolicyNumber = suspense.PolicyNumber
	resp.Data.SuspenseAmount = suspense.SuspenseAmount
	resp.Data.SuspenseReason = string(suspense.SuspenseReason)
	resp.Data.SuspenseStatus = string(suspense.SuspenseStatus)
	resp.Data.Priority = suspense.Priority
	if suspense.ResolutionDeadline != nil {
		resp.Data.ResolutionDeadline = suspense.ResolutionDeadline.Format("2006-01-02")
	}
	resp.Data.WorkflowID = suspense.WorkflowID

	return resp, nil
}

// GetAgingReport generates aging report for suspense accounts
// API: GET /commissions/suspense/aging-report
// Implements: FR-IC-COM-013, BR-IC-SUS-002
// Shows distribution of suspense entries across aging buckets
func (h *SuspenseHandler) GetAgingReport(
	sctx *serverRoute.Context,
	_ struct{}, // No request parameters
) (*response.SuspenseAgingReportResponse, error) {
	log.Info(sctx.Ctx, "Generating suspense aging report")

	// Get aging report from repository
	report, err := h.suspenseRepo.GetAgingReport(sctx.Ctx)
	if err != nil {
		log.Error(sctx.Ctx, "Error generating aging report: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Aging report generated: %d entries, total amount: %.2f",
		report.TotalEntries, report.TotalAmount)

	// Build response
	resp := &response.SuspenseAgingReportResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}
	resp.Data.ReportDate = report.ReportDate.Format("2006-01-02T15:04:05Z07:00")
	resp.Data.TotalEntries = report.TotalEntries
	resp.Data.TotalAmount = report.TotalAmount
	resp.Data.Buckets = response.NewSuspenseAgingBucketResponses(report.Buckets)
	resp.Data.OverdueCount = report.OverdueCount
	resp.Data.OverdueAmount = report.OverdueAmount

	return resp, nil
}
