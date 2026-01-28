package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"fmt"
	"time"
)

// CreateClawbackResponse represents the response for creating a clawback
type CreateClawbackResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		ClawbackID         int64   `json:"clawback_id"`
		PolicyNumber       string  `json:"policy_number"`
		AgentID            string  `json:"agent_id"`
		ClawbackAmount     float64 `json:"clawback_amount"`
		ClawbackPercentage float64 `json:"clawback_percentage"`
		RecoverySchedule   string  `json:"recovery_schedule"`
		Status             string  `json:"status"`
		TriggerDate        string  `json:"trigger_date"`
		WorkflowID         *string `json:"workflow_id,omitempty"`
	} `json:"data"`
}

// ClawbackDetailResponse represents detailed clawback information
type ClawbackDetailResponse struct {
	ClawbackID          int64    `json:"clawback_id"`
	PolicyNumber        string   `json:"policy_number"`
	AgentID             string   `json:"agent_id"`
	OriginalCommission  float64  `json:"original_commission"`
	ClawbackAmount      float64  `json:"clawback_amount"`
	RecoveredAmount     float64  `json:"recovered_amount"`
	PendingAmount       float64  `json:"pending_amount"`
	RecoveryProgress    float64  `json:"recovery_progress_percentage"` // 0-100%
	ClawbackPercentage  float64  `json:"clawback_percentage"`
	ClawbackReason      string   `json:"clawback_reason"`
	ClawbackStatus      string   `json:"clawback_status"`
	PolicyAge           int      `json:"policy_age_months"`
	TriggerDate         string   `json:"trigger_date"`
	PolicyInceptionDate string   `json:"policy_inception_date"`
	PolicyEndDate       *string  `json:"policy_end_date,omitempty"`
	RecoveryStartDate   *string  `json:"recovery_start_date,omitempty"`
	RecoveryEndDate     *string  `json:"recovery_end_date,omitempty"`
	RecoverySchedule    string   `json:"recovery_schedule"`
	InstallmentMonths   *int     `json:"installment_months,omitempty"`
	ApprovedBy          *string  `json:"approved_by,omitempty"`
	ApprovedAt          *string  `json:"approved_at,omitempty"`
	WaivedBy            *string  `json:"waived_by,omitempty"`
	WaivedAt            *string  `json:"waived_at,omitempty"`
	WaiverReason        *string  `json:"waiver_reason,omitempty"`
	WorkflowID          *string  `json:"workflow_id,omitempty"`
	WorkflowState       *string  `json:"workflow_state,omitempty"`
	Notes               *string  `json:"notes,omitempty"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	Version             int      `json:"version"`
}

// ClawbackSummary represents a summary of clawback for list views
type ClawbackSummary struct {
	ClawbackID         int64   `json:"clawback_id"`
	PolicyNumber       string  `json:"policy_number"`
	AgentID            string  `json:"agent_id"`
	ClawbackAmount     float64 `json:"clawback_amount"`
	RecoveredAmount    float64 `json:"recovered_amount"`
	PendingAmount      float64 `json:"pending_amount"`
	RecoveryProgress   float64 `json:"recovery_progress_percentage"`
	ClawbackPercentage float64 `json:"clawback_percentage"`
	ClawbackStatus     string  `json:"clawback_status"`
	TriggerDate        string  `json:"trigger_date"`
}

// ClawbackRecoveryDetail represents recovery transaction details
type ClawbackRecoveryDetail struct {
	RecoveryID        int64   `json:"recovery_id"`
	ClawbackID        int64   `json:"clawback_id"`
	InstallmentNumber int     `json:"installment_number"`
	ScheduledAmount   float64 `json:"scheduled_amount"`
	RecoveredAmount   float64 `json:"recovered_amount"`
	RecoveryDate      string  `json:"recovery_date"`
	RecoveryMethod    string  `json:"recovery_method"`
	RecoveryStatus    string  `json:"recovery_status"`
	StatementID       *int64  `json:"statement_id,omitempty"`
	DisbursementID    *int64  `json:"disbursement_id,omitempty"`
	TransactionRef    *string `json:"transaction_ref,omitempty"`
	FailureReason     *string `json:"failure_reason,omitempty"`
}

// ClawbackSearchResponse represents search results for clawbacks
type ClawbackSearchResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Clawbacks  []ClawbackSummary    `json:"clawbacks"`
		Pagination PaginationMetadata   `json:"pagination"`
	} `json:"data"`
}

// ClawbackDetailWithRecoveriesResponse includes recovery history
type ClawbackDetailWithRecoveriesResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Clawback   ClawbackDetailResponse   `json:"clawback"`
		Recoveries []ClawbackRecoveryDetail `json:"recoveries"`
	} `json:"data"`
}

// NewClawbackDetailResponse creates a detailed clawback response
func NewClawbackDetailResponse(clawback domain.Clawback) ClawbackDetailResponse {
	detail := ClawbackDetailResponse{
		ClawbackID:          clawback.ClawbackID,
		PolicyNumber:        clawback.PolicyNumber,
		AgentID:             clawback.AgentID,
		OriginalCommission:  clawback.OriginalCommission,
		ClawbackAmount:      clawback.ClawbackAmount,
		RecoveredAmount:     clawback.RecoveredAmount,
		PendingAmount:       clawback.PendingAmount,
		RecoveryProgress:    clawback.RecoveryProgress(),
		ClawbackPercentage:  clawback.ClawbackPercentage,
		ClawbackReason:      string(clawback.ClawbackReason),
		ClawbackStatus:      string(clawback.ClawbackStatus),
		PolicyAge:           clawback.PolicyAge,
		TriggerDate:         clawback.TriggerDate.Format(time.RFC3339),
		PolicyInceptionDate: clawback.PolicyInceptionDate.Format(time.RFC3339),
		RecoverySchedule:    clawback.RecoverySchedule,
		InstallmentMonths:   clawback.InstallmentMonths,
		ApprovedBy:          clawback.ApprovedBy,
		WaivedBy:            clawback.WaivedBy,
		WaiverReason:        clawback.WaiverReason,
		WorkflowID:          clawback.WorkflowID,
		WorkflowState:       clawback.WorkflowState,
		Notes:               clawback.Notes,
		CreatedAt:           clawback.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           clawback.UpdatedAt.Format(time.RFC3339),
		Version:             clawback.Version,
	}

	if clawback.PolicyEndDate != nil {
		policyEndDateStr := clawback.PolicyEndDate.Format(time.RFC3339)
		detail.PolicyEndDate = &policyEndDateStr
	}
	if clawback.RecoveryStartDate != nil {
		recoveryStartDateStr := clawback.RecoveryStartDate.Format(time.RFC3339)
		detail.RecoveryStartDate = &recoveryStartDateStr
	}
	if clawback.RecoveryEndDate != nil {
		recoveryEndDateStr := clawback.RecoveryEndDate.Format(time.RFC3339)
		detail.RecoveryEndDate = &recoveryEndDateStr
	}
	if clawback.ApprovedAt != nil {
		approvedAtStr := clawback.ApprovedAt.Format(time.RFC3339)
		detail.ApprovedAt = &approvedAtStr
	}
	if clawback.WaivedAt != nil {
		waivedAtStr := clawback.WaivedAt.Format(time.RFC3339)
		detail.WaivedAt = &waivedAtStr
	}

	return detail
}

// NewClawbackSummary creates a summary clawback response
func NewClawbackSummary(clawback domain.Clawback) ClawbackSummary {
	return ClawbackSummary{
		ClawbackID:         clawback.ClawbackID,
		PolicyNumber:       clawback.PolicyNumber,
		AgentID:            clawback.AgentID,
		ClawbackAmount:     clawback.ClawbackAmount,
		RecoveredAmount:    clawback.RecoveredAmount,
		PendingAmount:      clawback.PendingAmount,
		RecoveryProgress:   clawback.RecoveryProgress(),
		ClawbackPercentage: clawback.ClawbackPercentage,
		ClawbackStatus:     string(clawback.ClawbackStatus),
		TriggerDate:        clawback.TriggerDate.Format(time.RFC3339),
	}
}

// NewClawbackSummaries creates a list of clawback summaries
func NewClawbackSummaries(clawbacks []domain.Clawback) []ClawbackSummary {
	summaries := make([]ClawbackSummary, len(clawbacks))
	for i, clawback := range clawbacks {
		summaries[i] = NewClawbackSummary(clawback)
	}
	return summaries
}

// NewClawbackRecoveryDetail creates a recovery detail response
func NewClawbackRecoveryDetail(recovery domain.ClawbackRecovery) ClawbackRecoveryDetail {
	detail := ClawbackRecoveryDetail{
		RecoveryID:        recovery.RecoveryID,
		ClawbackID:        recovery.ClawbackID,
		InstallmentNumber: recovery.InstallmentNumber,
		ScheduledAmount:   recovery.ScheduledAmount,
		RecoveredAmount:   recovery.RecoveredAmount,
		RecoveryDate:      recovery.RecoveryDate.Format(time.RFC3339),
		RecoveryMethod:    recovery.RecoveryMethod,
		RecoveryStatus:    recovery.RecoveryStatus,
		StatementID:       recovery.StatementID,
		DisbursementID:    recovery.DisbursementID,
		TransactionRef:    recovery.TransactionRef,
		FailureReason:     recovery.FailureReason,
	}
	return detail
}

// NewClawbackRecoveryDetails creates a list of recovery details
func NewClawbackRecoveryDetails(recoveries []domain.ClawbackRecovery) []ClawbackRecoveryDetail {
	details := make([]ClawbackRecoveryDetail, len(recoveries))
	for i, recovery := range recoveries {
		details[i] = NewClawbackRecoveryDetail(recovery)
	}
	return details
}

// FormatClawbackPercentage formats the clawback percentage with explanation
func FormatClawbackPercentage(policyAgeMonths int, percentage float64) string {
	year := (policyAgeMonths / 12) + 1
	return fmt.Sprintf("%.0f%% (Policy Year %d)", percentage, year)
}
