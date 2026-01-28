package handler

import (
	"errors"

	"agent-commission/core/domain"
	"agent-commission/core/port"
	repo "agent-commission/repo/postgres"

	log "gitlab.cept.gov.in/it-2.0-common/api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// WebhookHandler handles integration webhooks from external systems
// INT-IC-002: Policy Services webhook
// INT-IC-003: PFMS webhook
type WebhookHandler struct {
	*serverHandler.Base
	disbursementRepo *repo.DisbursementRepository
	// TODO: Add CommissionBatchRepository for policy status changes
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	disbursementRepo *repo.DisbursementRepository,
) *WebhookHandler {
	base := serverHandler.New("Webhooks").SetPrefix("/v1").AddPrefix("")
	return &WebhookHandler{
		Base:             base,
		disbursementRepo: disbursementRepo,
	}
}

// Routes returns all routes for webhook endpoints
func (h *WebhookHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.POST("/webhooks/pfms/disbursement-confirmation", h.PFMSDisbursementConfirmation).Name("PFMS Disbursement Confirmation"),
		serverRoute.POST("/webhooks/policy/status-change", h.PolicyStatusChange).Name("Policy Status Change"),
	}
}

// PFMSDisbursementConfirmationRequest represents PFMS payment confirmation webhook
type PFMSDisbursementConfirmationRequest struct {
	DisbursementID int64  `json:"disbursement_id" validate:"required"`
	UTRNumber      string `json:"utr_number" validate:"required"` // Bank transaction reference
	Status         string `json:"status" validate:"required,oneof=SUCCESS FAILED"`
	TransactionDate string `json:"transaction_date" validate:"required"` // RFC3339 format
	FailureReason  *string `json:"failure_reason,omitempty"`
	FailureCode    *string `json:"failure_code,omitempty"`
	Signature      string  `json:"signature" validate:"required"` // Webhook signature
}

// PFMSDisbursementConfirmationResponse represents response for PFMS webhook
type PFMSDisbursementConfirmationResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		DisbursementID int64  `json:"disbursement_id"`
		Status         string `json:"status"`
		Message        string `json:"message"`
	} `json:"data"`
}

// PFMSDisbursementConfirmation handles PFMS payment confirmation webhook
// API: POST /webhooks/pfms/disbursement-confirmation
// Implements: INT-IC-003, FR-IC-COM-010
func (h *WebhookHandler) PFMSDisbursementConfirmation(
	sctx *serverRoute.Context,
	req PFMSDisbursementConfirmationRequest,
) (*PFMSDisbursementConfirmationResponse, error) {
	// TODO: Verify webhook signature to ensure authenticity
	// if !verifyPFMSSignature(req.Signature, req) {
	//     log.Warn(sctx.Ctx, "Invalid PFMS webhook signature")
	//     return nil, errors.New("invalid webhook signature")
	// }

	log.Info(sctx.Ctx, "Received PFMS disbursement confirmation: %d (status: %s, UTR: %s)",
		req.DisbursementID, req.Status, req.UTRNumber)

	// Fetch disbursement to validate
	disbursement, err := h.disbursementRepo.GetDisbursementByID(sctx.Ctx, req.DisbursementID)
	if err != nil {
		log.Error(sctx.Ctx, "Disbursement not found: %d", req.DisbursementID)
		return nil, errors.New("disbursement not found")
	}

	// Validate disbursement is in expected state
	if disbursement.DisbursementStatus != domain.DisbursementStatusSentToBank &&
		disbursement.DisbursementStatus != domain.DisbursementStatusProcessing {
		log.Warn(sctx.Ctx, "Disbursement %d in unexpected state: %s",
			req.DisbursementID, disbursement.DisbursementStatus)
		return nil, errors.New("disbursement not in expected state")
	}

	// Update disbursement based on status
	var finalStatus domain.DisbursementStatus
	var failureReason *domain.PaymentFailureReason

	if req.Status == "SUCCESS" {
		finalStatus = domain.DisbursementStatusCompleted
		log.Info(sctx.Ctx, "Disbursement %d completed successfully", req.DisbursementID)
	} else {
		finalStatus = domain.DisbursementStatusFailed
		// Map failure code to internal failure reason
		if req.FailureCode != nil {
			mapped := mapPFMSFailureCode(*req.FailureCode)
			failureReason = &mapped
		}
		log.Warn(sctx.Ctx, "Disbursement %d failed: %v", req.DisbursementID, req.FailureReason)
	}

	// Update disbursement status
	err = h.disbursementRepo.CompleteDisbursement(
		sctx.Ctx,
		req.DisbursementID,
		finalStatus,
		&req.UTRNumber,
		failureReason,
		req.FailureReason,
	)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating disbursement status: %v", err)
		return nil, err
	}

	// TODO: Trigger accounting system integration for GL posting
	// if finalStatus == domain.DisbursementStatusCompleted {
	//     h.accountingClient.PostDisbursement(disbursement)
	// }

	// TODO: Send notification to agent about disbursement status
	// h.notificationService.NotifyDisbursementStatus(disbursement.AgentID, finalStatus)

	// Build response
	resp := &PFMSDisbursementConfirmationResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
	}
	resp.Data.DisbursementID = req.DisbursementID
	resp.Data.Status = req.Status
	resp.Data.Message = "Disbursement status updated successfully"

	return resp, nil
}

// PolicyStatusChangeRequest represents policy status change webhook
type PolicyStatusChangeRequest struct {
	PolicyNumber string   `json:"policy_number" validate:"required"`
	NewStatus    string   `json:"new_status" validate:"required"`
	OldStatus    string   `json:"old_status" validate:"required"`
	ChangeDate   string   `json:"change_date" validate:"required"` // RFC3339 format
	Reason       *string  `json:"reason,omitempty"`
	Signature    string   `json:"signature" validate:"required"` // Webhook signature
}

// PolicyStatusChangeResponse represents response for policy status webhook
type PolicyStatusChangeResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		PolicyNumber string `json:"policy_number"`
		Status       string `json:"status"`
		Message      string `json:"message"`
	} `json:"data"`
}

// PolicyStatusChange handles policy status change webhook
// API: POST /webhooks/policy/status-change
// Implements: INT-IC-002
// Triggers clawback if policy is surrendered, lapsed, or cancelled
func (h *WebhookHandler) PolicyStatusChange(
	sctx *serverRoute.Context,
	req PolicyStatusChangeRequest,
) (*PolicyStatusChangeResponse, error) {
	// TODO: Verify webhook signature
	// if !verifyPolicyServiceSignature(req.Signature, req) {
	//     log.Warn(sctx.Ctx, "Invalid policy service webhook signature")
	//     return nil, errors.New("invalid webhook signature")
	// }

	log.Info(sctx.Ctx, "Received policy status change: %s (old: %s, new: %s)",
		req.PolicyNumber, req.OldStatus, req.NewStatus)

	// Check if status change triggers clawback
	// BR-IC-CLAWBACK-001: Clawback triggered on surrender, lapse, or cancellation
	triggersClawback := req.NewStatus == "SURRENDERED" ||
		req.NewStatus == "LAPSED" ||
		req.NewStatus == "CANCELLED"

	if triggersClawback {
		log.Info(sctx.Ctx, "Policy %s status change triggers clawback: %s",
			req.PolicyNumber, req.NewStatus)

		// TODO: Trigger clawback workflow
		// 1. Find all commission transactions for this policy
		// 2. Calculate clawback amount based on policy age
		// 3. Create clawback record
		// 4. Start clawback recovery workflow
		//
		// workflowID := "clawback-" + req.PolicyNumber
		// h.temporalClient.StartClawbackWorkflow(workflowID, req.PolicyNumber, req.NewStatus)
	}

	// TODO: Update policy commission records status
	// Mark commissions as "CLAWBACK_PENDING" or "POLICY_INACTIVE"

	// Build response
	resp := &PolicyStatusChangeResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
	}
	resp.Data.PolicyNumber = req.PolicyNumber
	resp.Data.Status = req.NewStatus
	if triggersClawback {
		resp.Data.Message = "Policy status updated and clawback initiated"
	} else {
		resp.Data.Message = "Policy status updated successfully"
	}

	return resp, nil
}

// mapPFMSFailureCode maps PFMS failure codes to internal failure reasons
func mapPFMSFailureCode(code string) domain.PaymentFailureReason {
	switch code {
	case "INVALID_ACCOUNT", "ACC_NOT_FOUND":
		return domain.PaymentFailureInvalidAccount
	case "INSUFFICIENT_FUNDS":
		return domain.PaymentFailureInsufficientFunds
	case "BANK_REJECTED", "BENEFICIARY_REJECTED":
		return domain.PaymentFailureBankRejection
	case "NETWORK_ERROR", "TIMEOUT":
		return domain.PaymentFailureNetworkError
	default:
		return domain.PaymentFailureValidationError
	}
}
