package domain

import (
	"time"
)

// DisbursementMode represents the mode of disbursement
// BR-IC-COM-008: Disbursement mode workflow (Cheque/EFT)
type DisbursementMode string

const (
	DisbursementModeCheque DisbursementMode = "CHEQUE"
	DisbursementModeEFT    DisbursementMode = "EFT"
)

// DisbursementStatus represents the current status of a disbursement
// FR-IC-COM-009, FR-IC-COM-010: Disbursement lifecycle
type DisbursementStatus string

const (
	DisbursementStatusPending    DisbursementStatus = "PENDING"     // Created, awaiting processing
	DisbursementStatusProcessing DisbursementStatus = "PROCESSING"  // Being processed
	DisbursementStatusSentToBank DisbursementStatus = "SENT_TO_BANK" // EFT sent to PFMS/Bank
	DisbursementStatusCompleted  DisbursementStatus = "COMPLETED"   // Successfully disbursed
	DisbursementStatusFailed     DisbursementStatus = "FAILED"      // Disbursement failed
	DisbursementStatusCancelled  DisbursementStatus = "CANCELLED"   // Manually cancelled
)

// PaymentFailureReason represents reasons for payment failure
type PaymentFailureReason string

const (
	PaymentFailureInvalidAccount     PaymentFailureReason = "INVALID_ACCOUNT"
	PaymentFailureInsufficientFunds  PaymentFailureReason = "INSUFFICIENT_FUNDS"
	PaymentFailureBankRejection      PaymentFailureReason = "BANK_REJECTION"
	PaymentFailureNetworkError       PaymentFailureReason = "NETWORK_ERROR"
	PaymentFailureValidationError    PaymentFailureReason = "VALIDATION_ERROR"
)

// Disbursement represents a commission disbursement to an agent
// FR-IC-COM-009: Cheque disbursement
// FR-IC-COM-010: EFT disbursement via PFMS
// BR-IC-COM-011: 10-working-day SLA for disbursement
type Disbursement struct {
	// Primary Key
	DisbursementID int64 `db:"disbursement_id" json:"disbursement_id"`

	// Reference Information
	FinalStatementID int64  `db:"final_statement_id" json:"final_statement_id"` // Link to final statement
	AgentID          string `db:"agent_id" json:"agent_id"`
	AgentCode        string `db:"agent_code" json:"agent_code"`

	// Disbursement Details
	DisbursementMode   DisbursementMode   `db:"disbursement_mode" json:"disbursement_mode"`     // CHEQUE or EFT
	DisbursementStatus DisbursementStatus `db:"disbursement_status" json:"disbursement_status"` // Current status

	// Financial Information
	TotalGrossCommission float64 `db:"total_gross_commission" json:"total_gross_commission"`
	TotalTDS             float64 `db:"total_tds" json:"total_tds"`
	TotalNetCommission   float64 `db:"total_net_commission" json:"total_net_commission"` // Amount to disburse

	// Cheque Details (for CHEQUE mode)
	ChequeNumber *string    `db:"cheque_number" json:"cheque_number,omitempty"`
	ChequeDate   *time.Time `db:"cheque_date" json:"cheque_date,omitempty"`

	// EFT Details (for EFT mode)
	BankAccountNumber *string `db:"bank_account_number" json:"bank_account_number,omitempty"`
	BankName          *string `db:"bank_name" json:"bank_name,omitempty"`
	BankBranch        *string `db:"bank_branch" json:"bank_branch,omitempty"`
	IFSCCode          *string `db:"ifsc_code" json:"ifsc_code,omitempty"`
	AccountHolderName *string `db:"account_holder_name" json:"account_holder_name,omitempty"`
	UTRNumber         *string `db:"utr_number" json:"utr_number,omitempty"` // PFMS transaction reference

	// SLA Tracking - BR-IC-COM-011: 10 working days
	InitiatedAt     time.Time  `db:"initiated_at" json:"initiated_at"`
	SLADeadline     time.Time  `db:"sla_deadline" json:"sla_deadline"`           // 10 working days from initiated_at
	ProcessedAt     *time.Time `db:"processed_at" json:"processed_at,omitempty"` // When processing started
	CompletedAt     *time.Time `db:"completed_at" json:"completed_at,omitempty"` // When completed/failed
	SLABreached     bool       `db:"sla_breached" json:"sla_breached"`
	BreachReasonID  *int64     `db:"breach_reason_id" json:"breach_reason_id,omitempty"`

	// Failure Tracking
	FailureReason *PaymentFailureReason `db:"failure_reason" json:"failure_reason,omitempty"`
	FailureDetails *string              `db:"failure_details" json:"failure_details,omitempty"`
	RetryCount     int                  `db:"retry_count" json:"retry_count"` // Number of retry attempts

	// Workflow Integration
	WorkflowID              *string    `db:"workflow_id" json:"workflow_id,omitempty"`
	WorkflowRunID           *string    `db:"workflow_run_id" json:"workflow_run_id,omitempty"`
	WorkflowState           *string    `db:"workflow_state" json:"workflow_state,omitempty"`
	WorkflowStateUpdatedAt  *time.Time `db:"workflow_state_updated_at" json:"workflow_state_updated_at,omitempty"`

	// Accounting Integration
	VoucherNumber *string    `db:"voucher_number" json:"voucher_number,omitempty"` // GL voucher number
	PostedToGL    bool       `db:"posted_to_gl" json:"posted_to_gl"`
	GLPostedAt    *time.Time `db:"gl_posted_at" json:"gl_posted_at,omitempty"`

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	UpdatedBy string     `db:"updated_by" json:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Version   int        `db:"version" json:"version"` // Optimistic locking
}

// IsSLABreached checks if the disbursement has breached SLA
// BR-IC-COM-011: 10 working days SLA
func (d *Disbursement) IsSLABreached() bool {
	return time.Now().After(d.SLADeadline) && d.DisbursementStatus != DisbursementStatusCompleted
}

// CanProcess checks if disbursement can be processed
// VR-IC-COM-004: Disbursement validation
func (d *Disbursement) CanProcess() bool {
	return d.DisbursementStatus == DisbursementStatusPending && d.DeletedAt == nil
}

// CanRetry checks if failed disbursement can be retried
func (d *Disbursement) CanRetry() bool {
	return d.DisbursementStatus == DisbursementStatusFailed && d.RetryCount < 3 && d.DeletedAt == nil
}

// IsCompleted checks if disbursement is successfully completed
func (d *Disbursement) IsCompleted() bool {
	return d.DisbursementStatus == DisbursementStatusCompleted
}

// DisbursementFilter represents filter criteria for searching disbursements
type DisbursementFilter struct {
	AgentID            *string             `json:"agent_id" query:"agent_id"`
	DisbursementMode   *DisbursementMode   `json:"disbursement_mode" query:"disbursement_mode"`
	DisbursementStatus *DisbursementStatus `json:"disbursement_status" query:"disbursement_status"`
	FromDate           *time.Time          `json:"from_date" query:"from_date"`
	ToDate             *time.Time          `json:"to_date" query:"to_date"`
	SLABreached        *bool               `json:"sla_breached" query:"sla_breached"`
	Page               int                 `json:"page" query:"page"`
	Limit              int                 `json:"limit" query:"limit"`
}
