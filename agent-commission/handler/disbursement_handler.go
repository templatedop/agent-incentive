package handler

import (
	"errors"
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

// DisbursementHandler handles commission disbursement operations
// FR-IC-COM-009: Cheque disbursement
// FR-IC-COM-010: EFT disbursement via PFMS
type DisbursementHandler struct {
	*serverHandler.Base
	disbursementRepo *repo.DisbursementRepository
	finalStmtRepo    *repo.FinalStatementRepository
}

// NewDisbursementHandler creates a new disbursement handler
func NewDisbursementHandler(
	disbursementRepo *repo.DisbursementRepository,
	finalStmtRepo *repo.FinalStatementRepository,
) *DisbursementHandler {
	base := serverHandler.New("Disbursements").SetPrefix("/v1").AddPrefix("")
	return &DisbursementHandler{
		Base:             base,
		disbursementRepo: disbursementRepo,
		finalStmtRepo:    finalStmtRepo,
	}
}

// Routes returns all routes for disbursement endpoints
func (h *DisbursementHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.POST("/commissions/disbursements", h.CreateDisbursement).Name("Create Disbursement"),
		serverRoute.GET("/commissions/disbursements/:disbursementId/status", h.GetDisbursementStatus).Name("Get Disbursement Status"),
	}
}

// CreateDisbursementRequest represents request to create a disbursement
type CreateDisbursementRequest struct {
	FinalStatementID int64                    `json:"final_statement_id" validate:"required"`
	DisbursementMode domain.DisbursementMode  `json:"disbursement_mode" validate:"required,oneof=CHEQUE EFT"`

	// Cheque details (required if mode is CHEQUE)
	ChequeDetails *struct {
		ChequeNumber string `json:"cheque_number" validate:"required"`
		ChequeDate   string `json:"cheque_date" validate:"required"` // Format: YYYY-MM-DD
	} `json:"cheque_details,omitempty"`

	// Bank details (required if mode is EFT)
	BankDetails *struct {
		AccountNumber    string `json:"account_number" validate:"required"`
		AccountHolder    string `json:"account_holder" validate:"required"`
		BankName         string `json:"bank_name" validate:"required"`
		BankBranch       string `json:"bank_branch" validate:"required"`
		IFSCCode         string `json:"ifsc_code" validate:"required,len=11"`
	} `json:"bank_details,omitempty"`

	Remarks *string `json:"remarks,omitempty"`
}

// CreateDisbursement processes commission disbursement
// API: POST /commissions/disbursements
// Implements: FR-IC-COM-009, FR-IC-COM-010, VR-IC-COM-004
// BR-IC-COM-008: Disbursement mode workflow (Cheque/EFT)
// BR-IC-COM-011: 10-working-day SLA
func (h *DisbursementHandler) CreateDisbursement(
	sctx *serverRoute.Context,
	req CreateDisbursementRequest,
) (*response.CreateDisbursementResponse, error) {
	// VR-IC-COM-004: Validate disbursement mode and details
	if req.DisbursementMode == domain.DisbursementModeCheque && req.ChequeDetails == nil {
		log.Warn(sctx.Ctx, "Cheque details required for CHEQUE mode")
		return nil, errors.New("cheque details required for CHEQUE disbursement mode")
	}

	if req.DisbursementMode == domain.DisbursementModeEFT && req.BankDetails == nil {
		log.Warn(sctx.Ctx, "Bank details required for EFT mode")
		return nil, errors.New("bank details required for EFT disbursement mode")
	}

	// Fetch final statement to validate and get agent details
	finalStatement, err := h.finalStmtRepo.GetFinalStatementByID(sctx.Ctx, req.FinalStatementID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn(sctx.Ctx, "Final statement not found: %d", req.FinalStatementID)
			return nil, errors.New("final statement not found")
		}
		log.Error(sctx.Ctx, "Error fetching final statement: %v", err)
		return nil, err
	}

	// Validate final statement status
	if finalStatement.StatementStatus != domain.FinalStatementStatusReadyForDisbursement {
		log.Warn(sctx.Ctx, "Final statement not ready for disbursement: %d (status: %s)",
			req.FinalStatementID, finalStatement.StatementStatus)
		return nil, errors.New("final statement is not ready for disbursement")
	}

	// Create disbursement record
	disbursement := &domain.Disbursement{
		FinalStatementID:     req.FinalStatementID,
		AgentID:              finalStatement.AgentID,
		AgentCode:            "", // TODO: Get from agent profile lookup
		DisbursementMode:     req.DisbursementMode,
		DisbursementStatus:   domain.DisbursementStatusPending,
		TotalGrossCommission: finalStatement.TotalGrossCommission,
		TotalTDS:             finalStatement.TotalTDS,
		TotalNetCommission:   finalStatement.TotalNetCommission,
		InitiatedAt:          time.Now(),
		RetryCount:           0,
		PostedToGL:           false,
		SLABreached:          false,
		CreatedBy:            "system", // TODO: Get from JWT
		UpdatedBy:            "system",
	}

	// Set mode-specific details
	if req.DisbursementMode == domain.DisbursementModeCheque {
		disbursement.ChequeNumber = &req.ChequeDetails.ChequeNumber
		chequeDate, _ := time.Parse("2006-01-02", req.ChequeDetails.ChequeDate)
		disbursement.ChequeDate = &chequeDate
	} else if req.DisbursementMode == domain.DisbursementModeEFT {
		disbursement.BankAccountNumber = &req.BankDetails.AccountNumber
		disbursement.AccountHolderName = &req.BankDetails.AccountHolder
		disbursement.BankName = &req.BankDetails.BankName
		disbursement.BankBranch = &req.BankDetails.BankBranch
		disbursement.IFSCCode = &req.BankDetails.IFSCCode
	}

	// TODO: Start Temporal workflow for disbursement processing
	// workflowID := "disbursement-" + fmt.Sprintf("%d", req.FinalStatementID)
	// disbursement.WorkflowID = &workflowID

	// Create disbursement
	err = h.disbursementRepo.CreateDisbursement(sctx.Ctx, disbursement)
	if err != nil {
		log.Error(sctx.Ctx, "Error creating disbursement: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Disbursement created successfully: %d for agent %s",
		disbursement.DisbursementID, disbursement.AgentID)

	// Build response
	resp := &response.CreateDisbursementResponse{
		StatusCodeAndMessage: port.CreateSuccess,
	}
	resp.Data.DisbursementID = disbursement.DisbursementID
	resp.Data.AgentID = disbursement.AgentID
	resp.Data.Amount = disbursement.TotalNetCommission
	resp.Data.Mode = string(disbursement.DisbursementMode)
	resp.Data.Status = string(disbursement.DisbursementStatus)
	resp.Data.SLADeadline = disbursement.SLADeadline.Format(time.RFC3339)
	resp.Data.WorkflowID = disbursement.WorkflowID

	return resp, nil
}

// DisbursementIDUri represents the URI parameter for disbursement ID
type DisbursementIDUri struct {
	DisbursementID int64 `uri:"disbursementId" validate:"required"`
}

// GetDisbursementStatus retrieves the status of a disbursement
// API: GET /commissions/disbursements/{disbursementId}/status
// Implements: FR-IC-COM-010
func (h *DisbursementHandler) GetDisbursementStatus(
	sctx *serverRoute.Context,
	uri DisbursementIDUri,
) (*response.DisbursementStatusResponse, error) {
	// Fetch disbursement
	disbursement, err := h.disbursementRepo.GetDisbursementByID(sctx.Ctx, uri.DisbursementID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn(sctx.Ctx, "Disbursement not found: %d", uri.DisbursementID)
			return nil, errors.New("disbursement not found")
		}
		log.Error(sctx.Ctx, "Error fetching disbursement: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Fetched disbursement status: %d (status: %s)",
		uri.DisbursementID, disbursement.DisbursementStatus)

	// Build response
	resp := &response.DisbursementStatusResponse{
		StatusCodeAndMessage: port.ListSuccess,
		Data:                 response.NewDisbursementDetailResponse(disbursement),
	}

	return resp, nil
}
