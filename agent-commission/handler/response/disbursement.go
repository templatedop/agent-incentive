package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"fmt"
	"time"
)

// DisbursementSummary represents a disbursement in list responses
type DisbursementSummary struct {
	DisbursementID     int64  `json:"disbursement_id"`
	FinalStatementID   int64  `json:"final_statement_id"`
	AgentID            string `json:"agent_id"`
	AgentCode          string `json:"agent_code"`
	DisbursementMode   string `json:"disbursement_mode"`   // CHEQUE or EFT
	DisbursementStatus string `json:"disbursement_status"` // Current status
	TotalNetCommission float64 `json:"total_net_commission"` // Amount disbursed
	ChequeNumber       *string `json:"cheque_number,omitempty"`
	UTRNumber          *string `json:"utr_number,omitempty"` // For EFT
	InitiatedAt        string  `json:"initiated_at"`
	CompletedAt        *string `json:"completed_at,omitempty"`
	SLADeadline        string  `json:"sla_deadline"`
	SLABreached        bool    `json:"sla_breached"`
	FailureReason      *string `json:"failure_reason,omitempty"`
}

// NewDisbursementSummary converts domain.Disbursement to DisbursementSummary
func NewDisbursementSummary(disbursement domain.Disbursement) DisbursementSummary {
	summary := DisbursementSummary{
		DisbursementID:     disbursement.DisbursementID,
		FinalStatementID:   disbursement.FinalStatementID,
		AgentID:            disbursement.AgentID,
		AgentCode:          disbursement.AgentCode,
		DisbursementMode:   string(disbursement.DisbursementMode),
		DisbursementStatus: string(disbursement.DisbursementStatus),
		TotalNetCommission: disbursement.TotalNetCommission,
		ChequeNumber:       disbursement.ChequeNumber,
		UTRNumber:          disbursement.UTRNumber,
		InitiatedAt:        disbursement.InitiatedAt.Format(time.RFC3339),
		SLADeadline:        disbursement.SLADeadline.Format(time.RFC3339),
		SLABreached:        disbursement.SLABreached,
	}

	if disbursement.CompletedAt != nil {
		completedStr := disbursement.CompletedAt.Format(time.RFC3339)
		summary.CompletedAt = &completedStr
	}

	if disbursement.FailureReason != nil {
		failureStr := string(*disbursement.FailureReason)
		summary.FailureReason = &failureStr
	}

	return summary
}

// NewDisbursementSummaries converts a slice of disbursements to summaries
func NewDisbursementSummaries(disbursements []domain.Disbursement) []DisbursementSummary {
	summaries := make([]DisbursementSummary, 0, len(disbursements))
	for _, disbursement := range disbursements {
		summaries = append(summaries, NewDisbursementSummary(disbursement))
	}
	return summaries
}

// DisbursementDetailResponse represents detailed disbursement information
type DisbursementDetailResponse struct {
	DisbursementID       int64    `json:"disbursement_id"`
	FinalStatementID     int64    `json:"final_statement_id"`
	AgentID              string   `json:"agent_id"`
	AgentCode            string   `json:"agent_code"`
	DisbursementMode     string   `json:"disbursement_mode"`
	DisbursementStatus   string   `json:"disbursement_status"`
	TotalGrossCommission float64  `json:"total_gross_commission"`
	TotalTDS             float64  `json:"total_tds"`
	TotalNetCommission   float64  `json:"total_net_commission"`
	ChequeNumber         *string  `json:"cheque_number,omitempty"`
	ChequeDate           *string  `json:"cheque_date,omitempty"`
	BankAccountNumber    *string  `json:"bank_account_number,omitempty"`
	BankName             *string  `json:"bank_name,omitempty"`
	BankBranch           *string  `json:"bank_branch,omitempty"`
	IFSCCode             *string  `json:"ifsc_code,omitempty"`
	AccountHolderName    *string  `json:"account_holder_name,omitempty"`
	UTRNumber            *string  `json:"utr_number,omitempty"`
	InitiatedAt          string   `json:"initiated_at"`
	SLADeadline          string   `json:"sla_deadline"`
	ProcessedAt          *string  `json:"processed_at,omitempty"`
	CompletedAt          *string  `json:"completed_at,omitempty"`
	SLABreached          bool     `json:"sla_breached"`
	TimeRemaining        *string  `json:"time_remaining,omitempty"` // Human-readable format
	FailureReason        *string  `json:"failure_reason,omitempty"`
	FailureDetails       *string  `json:"failure_details,omitempty"`
	RetryCount           int      `json:"retry_count"`
	WorkflowID           *string  `json:"workflow_id,omitempty"`
	WorkflowState        *string  `json:"workflow_state,omitempty"`
	VoucherNumber        *string  `json:"voucher_number,omitempty"`
	PostedToGL           bool     `json:"posted_to_gl"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
}

// NewDisbursementDetailResponse converts domain.Disbursement to detailed response
func NewDisbursementDetailResponse(disbursement domain.Disbursement) DisbursementDetailResponse {
	detail := DisbursementDetailResponse{
		DisbursementID:       disbursement.DisbursementID,
		FinalStatementID:     disbursement.FinalStatementID,
		AgentID:              disbursement.AgentID,
		AgentCode:            disbursement.AgentCode,
		DisbursementMode:     string(disbursement.DisbursementMode),
		DisbursementStatus:   string(disbursement.DisbursementStatus),
		TotalGrossCommission: disbursement.TotalGrossCommission,
		TotalTDS:             disbursement.TotalTDS,
		TotalNetCommission:   disbursement.TotalNetCommission,
		ChequeNumber:         disbursement.ChequeNumber,
		BankAccountNumber:    disbursement.BankAccountNumber,
		BankName:             disbursement.BankName,
		BankBranch:           disbursement.BankBranch,
		IFSCCode:             disbursement.IFSCCode,
		AccountHolderName:    disbursement.AccountHolderName,
		UTRNumber:            disbursement.UTRNumber,
		InitiatedAt:          disbursement.InitiatedAt.Format(time.RFC3339),
		SLADeadline:          disbursement.SLADeadline.Format(time.RFC3339),
		SLABreached:          disbursement.SLABreached,
		RetryCount:           disbursement.RetryCount,
		WorkflowID:           disbursement.WorkflowID,
		WorkflowState:        disbursement.WorkflowState,
		VoucherNumber:        disbursement.VoucherNumber,
		PostedToGL:           disbursement.PostedToGL,
		CreatedAt:            disbursement.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            disbursement.UpdatedAt.Format(time.RFC3339),
	}

	if disbursement.ChequeDate != nil {
		chequeDateStr := disbursement.ChequeDate.Format("2006-01-02")
		detail.ChequeDate = &chequeDateStr
	}

	if disbursement.ProcessedAt != nil {
		processedStr := disbursement.ProcessedAt.Format(time.RFC3339)
		detail.ProcessedAt = &processedStr
	}

	if disbursement.CompletedAt != nil {
		completedStr := disbursement.CompletedAt.Format(time.RFC3339)
		detail.CompletedAt = &completedStr
	}

	if disbursement.FailureReason != nil {
		failureStr := string(*disbursement.FailureReason)
		detail.FailureReason = &failureStr
	}

	detail.FailureDetails = disbursement.FailureDetails

	// Calculate time remaining for SLA
	if !disbursement.IsCompleted() && !disbursement.SLABreached {
		timeRemaining := time.Until(disbursement.SLADeadline)
		if timeRemaining > 0 {
			hours := int(timeRemaining.Hours())
			minutes := int(timeRemaining.Minutes()) % 60
			timeRemainingStr := fmt.Sprintf("%d hours %d minutes", hours, minutes)
			detail.TimeRemaining = &timeRemainingStr
		}
	}

	return detail
}

// CreateDisbursementResponse represents response for POST /commissions/disbursements
type CreateDisbursementResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		DisbursementID int64   `json:"disbursement_id"`
		AgentID        string  `json:"agent_id"`
		Amount         float64 `json:"amount"`
		Mode           string  `json:"mode"`
		Status         string  `json:"status"`
		SLADeadline    string  `json:"sla_deadline"`
		WorkflowID     *string `json:"workflow_id,omitempty"`
	} `json:"data"`
}

// DisbursementStatusResponse represents response for GET /commissions/disbursements/{id}/status
type DisbursementStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      DisbursementDetailResponse `json:"data"`
}

// DisbursementsListResponse represents response for GET /commissions/disbursements
type DisbursementsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Disbursements []DisbursementSummary `json:"disbursements"`
		Pagination    PaginationMetadata    `json:"pagination"`
	} `json:"data"`
}
