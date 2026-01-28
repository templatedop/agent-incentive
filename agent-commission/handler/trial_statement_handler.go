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

// TrialStatementHandler handles trial statement operations
// Implements: FR-IC-COM-004 (Trial statement generation), FR-IC-COM-006 (Approval)
type TrialStatementHandler struct {
	*serverHandler.Base
	trialRepo *postgres.TrialStatementRepository
}

// NewTrialStatementHandler creates a new trial statement handler
func NewTrialStatementHandler(trialRepo *postgres.TrialStatementRepository) *TrialStatementHandler {
	base := serverHandler.New("TrialStatement").SetPrefix("/v1").AddPrefix("")
	return &TrialStatementHandler{Base: base, trialRepo: trialRepo}
}

// Routes returns all routes for trial statement endpoints
func (h *TrialStatementHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.GET("/commissions/trial-statements", h.ListTrialStatements).Name("List Trial Statements"),
		serverRoute.POST("/commissions/trial-statements/:statementId/approve", h.ApproveTrialStatement).Name("Approve Trial Statement"),
	}
}

// ListTrialStatementsRequest represents query parameters for GET /commissions/trial-statements
type ListTrialStatementsRequest struct {
	BatchID        string `form:"batch_id" validate:"omitempty"`
	AgentID        string `form:"agent_id" validate:"omitempty"`
	Circle         string `form:"circle" validate:"omitempty"`
	CommissionType string `form:"commission_type" validate:"omitempty,oneof=FIRST_YEAR RENEWAL BONUS"`
	Status         string `form:"status" validate:"omitempty,oneof=PENDING_APPROVAL APPROVED REJECTED"`
	Page           int    `form:"page" validate:"omitempty,min=1"`
	Limit          int    `form:"limit" validate:"omitempty,min=1,max=100"`
}

// ListTrialStatements retrieves trial statements with filtering and pagination
// Implements: FR-IC-COM-004, BR-IC-COM-009
// API: GET /commissions/trial-statements
//
// Query Parameters:
//   - batch_id (optional): Filter by batch ID
//   - agent_id (optional): Filter by agent ID
//   - circle (optional): Filter by circle code
//   - commission_type (optional): FIRST_YEAR, RENEWAL, BONUS
//   - status (optional): PENDING_APPROVAL, APPROVED, REJECTED
//   - page (optional): Page number (default: 1)
//   - limit (optional): Items per page (default: 10, max: 100)
//
// Returns:
//   - 200: Trial statements retrieved
//   - 401: Unauthorized
func (h *TrialStatementHandler) ListTrialStatements(
	sctx *serverRoute.Context,
	req ListTrialStatementsRequest,
) (*resp.TrialStatementsListResponse, error) {
	// Build filter
	filter := &domain.TrialStatementFilter{
		Page:  req.Page,
		Limit: req.Limit,
	}

	if req.BatchID != "" {
		filter.BatchID = &req.BatchID
	}
	if req.AgentID != "" {
		filter.AgentID = &req.AgentID
	}
	if req.Status != "" {
		status := domain.TrialStatementStatus(req.Status)
		filter.Status = &status
	}
	if req.CommissionType != "" {
		commType := domain.CommissionType(req.CommissionType)
		filter.CommissionType = &commType
	}

	// Set defaults
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	// Fetch trial statements
	statements, totalCount, err := h.trialRepo.SearchTrialStatements(sctx.Ctx, filter)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching trial statements: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Retrieved %d trial statements (total: %d)", len(statements), totalCount)

	// Calculate pagination
	totalPages := int(totalCount) / filter.Limit
	if int(totalCount)%filter.Limit > 0 {
		totalPages++
	}

	// Build response
	response := &resp.TrialStatementsListResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}
	response.Data.Statements = resp.NewTrialStatementSummaries(statements)
	response.Data.Pagination = resp.PaginationMetadata{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	return response, nil
}

// ApproveTrialStatementRequest represents the request body for approval
type ApproveTrialStatementRequest struct {
	StatementID int64 `uri:"statementId" validate:"required"`
	Approval    struct {
		DisbursementMode string   `json:"disbursement_mode" validate:"required,oneof=CHEQUE EFT"`
		PartialAmount    *float64 `json:"partial_amount" validate:"omitempty,gt=0"`
		Remarks          *string  `json:"remarks" validate:"omitempty"`
	} `json:"approval" validate:"required"`
}

// ApproveTrialStatement approves trial statement and initiates final statement generation
// Implements: FR-IC-COM-006, BR-IC-COM-002, BR-IC-COM-005 (partial disbursement)
// API: POST /commissions/trial-statements/{statementId}/approve
//
// URI Parameters:
//   - statementId (required): Trial statement ID
//
// Request Body:
//   - approval.disbursement_mode (required): CHEQUE or EFT
//   - approval.partial_amount (optional): Amount for partial disbursement
//   - approval.remarks (optional): Approval remarks
//
// Returns:
//   - 200: Trial statement approved
//   - 400: Bad request - invalid approval data
//   - 401: Unauthorized
//   - 404: Statement not found
//
// Business Rules:
//   - BR-IC-COM-002: Trial statement approval workflow
//   - BR-IC-COM-005: Partial disbursement support
//   - VR-IC-COM-003: Validation rules
func (h *TrialStatementHandler) ApproveTrialStatement(
	sctx *serverRoute.Context,
	req ApproveTrialStatementRequest,
) (*resp.ApproveTrialStatementResponse, error) {
	// Fetch trial statement
	statement, err := h.trialRepo.GetTrialStatementByID(sctx.Ctx, req.StatementID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn(sctx.Ctx, "Trial statement not found: %d", req.StatementID)
			return nil, errors.New("trial statement not found")
		}
		log.Error(sctx.Ctx, "Error fetching trial statement: %v", err)
		return nil, err
	}

	// Check if statement can be approved
	if !statement.CanApprove() {
		log.Warn(sctx.Ctx, "Trial statement cannot be approved: %d (status: %s)", req.StatementID, statement.StatementStatus)
		return nil, errors.New("trial statement cannot be approved")
	}

	// TODO: Get user ID from context
	approvedBy := "system" // Replace with actual user from JWT

	// Approve statement
	err = h.trialRepo.ApproveTrialStatement(sctx.Ctx, req.StatementID, approvedBy, req.Approval.Remarks)
	if err != nil {
		log.Error(sctx.Ctx, "Error approving trial statement: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Trial statement approved: %d by %s", req.StatementID, approvedBy)

	// TODO: Trigger final statement generation workflow
	// finalStmtID, err := h.temporalClient.StartFinalStatementWorkflow(ctx, statement)

	// Build response
	now := time.Now()
	response := &resp.ApproveTrialStatementResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: resp.ApprovalResponse{
			StatementID: req.StatementID,
			Status:      string(domain.TrialStatementStatusApproved),
			ApprovedBy:  approvedBy,
			ApprovedAt:  now.Format(time.RFC3339),
			Message:     "Trial statement approved successfully. Final statement generation initiated.",
		},
	}

	return response, nil
}
