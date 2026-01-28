package domain

import (
	"time"
)

// FinalStatementStatus represents the status of a final statement
type FinalStatementStatus string

const (
	FinalStatementStatusFinalized           FinalStatementStatus = "FINALIZED"
	FinalStatementStatusReadyForDisbursement FinalStatementStatus = "READY_FOR_DISBURSEMENT"
	FinalStatementStatusDisbursed           FinalStatementStatus = "DISBURSED"
)

// FinalStatement represents the final commission statement after trial approval
// Implements: FR-IC-COM-008 (Final statement generation)
// Business Rule: BR-IC-COM-007 (Final statement generation after trial approval)
type FinalStatement struct {
	FinalStatementID       int64                `json:"final_statement_id" db:"final_statement_id"`
	StatementNumber        string               `json:"statement_number" db:"statement_number"`
	TrialStatementID       int64                `json:"trial_statement_id" db:"trial_statement_id"`
	AgentID                string               `json:"agent_id" db:"agent_id"`
	StatementDate          time.Time            `json:"statement_date" db:"statement_date"`
	TotalGrossCommission   float64              `json:"total_gross_commission" db:"total_gross_commission"`
	TotalTDS               float64              `json:"total_tds" db:"total_tds"`
	TotalNetCommission     float64              `json:"total_net_commission" db:"total_net_commission"`
	StatementStatus        FinalStatementStatus `json:"statement_status" db:"statement_status"`
	PDFPath                *string              `json:"pdf_path,omitempty" db:"pdf_path"`
	WorkflowState          *string              `json:"workflow_state,omitempty" db:"workflow_state"`
	WorkflowStateUpdatedAt *time.Time           `json:"workflow_state_updated_at,omitempty" db:"workflow_state_updated_at"`
	CreatedAt              time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at" db:"updated_at"`
	CreatedBy              *string              `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy              *string              `json:"updated_by,omitempty" db:"updated_by"`
	DeletedAt              *time.Time           `json:"deleted_at,omitempty" db:"deleted_at"`
	Version                int                  `json:"version" db:"version"`
}

// MarkReadyForDisbursement updates the status to ready for disbursement
func (fs *FinalStatement) MarkReadyForDisbursement() {
	fs.StatementStatus = FinalStatementStatusReadyForDisbursement
	fs.UpdatedAt = time.Now()
}

// MarkDisbursed updates the status to disbursed
func (fs *FinalStatement) MarkDisbursed() {
	fs.StatementStatus = FinalStatementStatusDisbursed
	fs.UpdatedAt = time.Now()
}

// FinalStatementFilter represents filter criteria for querying final statements
type FinalStatementFilter struct {
	BatchID *string
	AgentID *string
	Status  *FinalStatementStatus
	Page    int
	Limit   int
}
