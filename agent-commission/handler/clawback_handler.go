package handler

import (
	"errors"
	"fmt"
	"time"

	"agent-commission/core/domain"
	"agent-commission/core/port"
	"agent-commission/handler/response"
	repo "agent-commission/repo/postgres"

	"github.com/jackc/pgx/v5"
	log "gitlab.cept.gov.in/it-2.0-common/api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// ClawbackHandler handles commission clawback operations
// FR-IC-COM-012: Clawback management
// BR-IC-CLAWBACK-001: Graduated clawback based on policy age
type ClawbackHandler struct {
	*serverHandler.Base
	clawbackRepo         *repo.ClawbackRepository
	commissionBatchRepo  *repo.CommissionBatchRepository
}

// NewClawbackHandler creates a new clawback handler
func NewClawbackHandler(
	clawbackRepo *repo.ClawbackRepository,
	commissionBatchRepo *repo.CommissionBatchRepository,
) *ClawbackHandler {
	base := serverHandler.New("Clawback").SetPrefix("/v1").AddPrefix("")
	return &ClawbackHandler{
		Base:                base,
		clawbackRepo:        clawbackRepo,
		commissionBatchRepo: commissionBatchRepo,
	}
}

// Routes returns all routes for clawback endpoints
func (h *ClawbackHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.POST("/commissions/clawback/create", h.CreateClawback).Name("Create Clawback"),
	}
}

// CreateClawbackRequest represents request to create a clawback
type CreateClawbackRequest struct {
	PolicyNumber        string                 `json:"policy_number" validate:"required"`
	AgentID             string                 `json:"agent_id" validate:"required"`
	ClawbackReason      domain.ClawbackReason  `json:"clawback_reason" validate:"required"`
	PolicyInceptionDate string                 `json:"policy_inception_date" validate:"required"` // RFC3339 format
	PolicyEndDate       *string                `json:"policy_end_date,omitempty"`                 // RFC3339 format
	RecoverySchedule    string                 `json:"recovery_schedule" validate:"required,oneof=IMMEDIATE INSTALLMENT"`
	InstallmentMonths   *int                   `json:"installment_months,omitempty"` // Required if INSTALLMENT
	Notes               *string                `json:"notes,omitempty"`
}

// CreateClawback initiates a commission clawback
// API: POST /commissions/clawback/create
// Implements: FR-IC-COM-012, BR-IC-CLAWBACK-001, BR-IC-CLAWBACK-002
// Triggered when policy is surrendered, lapsed, or cancelled
func (h *ClawbackHandler) CreateClawback(
	sctx *serverRoute.Context,
	req CreateClawbackRequest,
) (*response.CreateClawbackResponse, error) {
	// Validate recovery schedule
	if req.RecoverySchedule == "INSTALLMENT" && req.InstallmentMonths == nil {
		log.Warn(sctx.Ctx, "Installment months required for INSTALLMENT recovery schedule")
		return nil, errors.New("installment_months required for INSTALLMENT recovery schedule")
	}

	// Parse dates
	policyInceptionDate, err := time.Parse(time.RFC3339, req.PolicyInceptionDate)
	if err != nil {
		log.Warn(sctx.Ctx, "Invalid policy inception date format: %v", err)
		return nil, errors.New("invalid policy_inception_date format, expected RFC3339")
	}

	var policyEndDate *time.Time
	if req.PolicyEndDate != nil {
		parsedEndDate, err := time.Parse(time.RFC3339, *req.PolicyEndDate)
		if err != nil {
			log.Warn(sctx.Ctx, "Invalid policy end date format: %v", err)
			return nil, errors.New("invalid policy_end_date format, expected RFC3339")
		}
		policyEndDate = &parsedEndDate
	}

	// Calculate total commission paid for this policy
	// In real implementation, this should fetch from commission_transactions table
	// For now, we'll need to get it from the request or calculate it
	commissionTransactions, err := h.commissionBatchRepo.GetCommissionsByPolicyNumber(sctx.Ctx, req.PolicyNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn(sctx.Ctx, "No commission transactions found for policy: %s", req.PolicyNumber)
			return nil, errors.New("no commission transactions found for this policy")
		}
		log.Error(sctx.Ctx, "Error fetching commission transactions: %v", err)
		return nil, err
	}

	// Calculate total commission paid (sum of net commissions)
	var totalCommissionPaid float64
	for _, txn := range commissionTransactions {
		totalCommissionPaid += txn.NetCommission
	}

	if totalCommissionPaid == 0 {
		log.Warn(sctx.Ctx, "No commission paid for policy %s", req.PolicyNumber)
		return nil, errors.New("no commission has been paid for this policy")
	}

	log.Info(sctx.Ctx, "Total commission paid for policy %s: %.2f", req.PolicyNumber, totalCommissionPaid)

	// Create clawback record
	clawback := &domain.Clawback{
		PolicyNumber:        req.PolicyNumber,
		AgentID:             req.AgentID,
		OriginalCommission:  totalCommissionPaid,
		ClawbackReason:      req.ClawbackReason,
		PolicyInceptionDate: policyInceptionDate,
		PolicyEndDate:       policyEndDate,
		RecoverySchedule:    req.RecoverySchedule,
		InstallmentMonths:   req.InstallmentMonths,
		Notes:               req.Notes,
		CreatedBy:           "system", // TODO: Get from JWT
		UpdatedBy:           "system",
	}

	// Repository will calculate policy age, clawback percentage, and clawback amount
	err = h.clawbackRepo.CreateClawback(sctx.Ctx, clawback)
	if err != nil {
		log.Error(sctx.Ctx, "Error creating clawback: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Clawback created successfully: %d for policy %s (%.0f%% of %.2f = %.2f)",
		clawback.ClawbackID,
		clawback.PolicyNumber,
		clawback.ClawbackPercentage,
		clawback.OriginalCommission,
		clawback.ClawbackAmount,
	)

	// Start Temporal workflow for clawback recovery
	// Workflow will handle graduated recovery over time
	workflowID := fmt.Sprintf("clawback-%d-%s", clawback.ClawbackID, clawback.PolicyNumber)
	clawback.WorkflowID = &workflowID
	clawback.WorkflowState = ptrString("INITIATED")

	// TODO: Uncomment when Temporal client is available
	// err = h.temporalClient.ExecuteWorkflow(sctx.Ctx, client.StartWorkflowOptions{
	//     ID:        workflowID,
	//     TaskQueue: "clawback-recovery-queue",
	// }, "ClawbackRecoveryWorkflow", clawback)
	// if err != nil {
	//     log.Error(sctx.Ctx, "Error starting clawback workflow: %v", err)
	//     return nil, err
	// }

	log.Info(sctx.Ctx, "Clawback workflow ID: %s", workflowID)

	// Build response
	resp := &response.CreateClawbackResponse{
		StatusCodeAndMessage: port.CreateSuccess,
	}
	resp.Data.ClawbackID = clawback.ClawbackID
	resp.Data.PolicyNumber = clawback.PolicyNumber
	resp.Data.AgentID = clawback.AgentID
	resp.Data.ClawbackAmount = clawback.ClawbackAmount
	resp.Data.ClawbackPercentage = clawback.ClawbackPercentage
	resp.Data.RecoverySchedule = clawback.RecoverySchedule
	resp.Data.Status = string(clawback.ClawbackStatus)
	resp.Data.TriggerDate = clawback.TriggerDate.Format(time.RFC3339)
	resp.Data.WorkflowID = clawback.WorkflowID

	return resp, nil
}

// ptrString returns a pointer to the given string
func ptrString(s string) *string {
	return &s
}
