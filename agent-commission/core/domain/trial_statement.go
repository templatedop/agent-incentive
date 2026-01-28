package domain

import (
	"time"
)

// TrialStatementStatus represents the status of a trial statement
type TrialStatementStatus string

const (
	TrialStatementStatusPendingApproval TrialStatementStatus = "PENDING_APPROVAL"
	TrialStatementStatusApproved        TrialStatementStatus = "APPROVED"
	TrialStatementStatusRejected        TrialStatementStatus = "REJECTED"
	TrialStatementStatusCorrectionNeeded TrialStatementStatus = "CORRECTION_NEEDED"
)

// TrialStatement represents a trial commission statement for finance review
// Implements: FR-IC-COM-004 (Trial statement generation)
// Business Rule: BR-IC-COM-002 (Trial statement approval), BR-IC-COM-009 (7-day SLA)
type TrialStatement struct {
	TrialStatementID       int64                `json:"trial_statement_id" db:"trial_statement_id"`
	StatementNumber        string               `json:"statement_number" db:"statement_number"`
	AgentID                string               `json:"agent_id" db:"agent_id"`
	StatementDate          time.Time            `json:"statement_date" db:"statement_date"`
	FromDate               time.Time            `json:"from_date" db:"from_date"`
	ToDate                 time.Time            `json:"to_date" db:"to_date"`
	TotalPolicies          int                  `json:"total_policies" db:"total_policies"`
	TotalGrossCommission   float64              `json:"total_gross_commission" db:"total_gross_commission"`
	TotalTDS               float64              `json:"total_tds" db:"total_tds"`
	TotalNetCommission     float64              `json:"total_net_commission" db:"total_net_commission"`
	StatementStatus        TrialStatementStatus `json:"statement_status" db:"statement_status"`
	ApprovedBy             *string              `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt             *time.Time           `json:"approved_at,omitempty" db:"approved_at"`
	ApprovalRemarks        *string              `json:"approval_remarks,omitempty" db:"approval_remarks"`
	ProcessingUnit         *string              `json:"processing_unit,omitempty" db:"processing_unit"`
	BatchID                *string              `json:"batch_id,omitempty" db:"batch_id"`
	WorkflowState          *string              `json:"workflow_state,omitempty" db:"workflow_state"`
	WorkflowStateUpdatedAt *time.Time           `json:"workflow_state_updated_at,omitempty" db:"workflow_state_updated_at"`
	CreatedAt              time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at" db:"updated_at"`
	CreatedBy              *string              `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy              *string              `json:"updated_by,omitempty" db:"updated_by"`
	DeletedAt              *time.Time           `json:"deleted_at,omitempty" db:"deleted_at"`
	Version                int                  `json:"version" db:"version"`
}

// CanApprove checks if the trial statement can be approved
func (ts *TrialStatement) CanApprove() bool {
	return ts.StatementStatus == TrialStatementStatusPendingApproval && ts.DeletedAt == nil
}

// Approve marks the trial statement as approved
func (ts *TrialStatement) Approve(approvedBy string, remarks *string) {
	now := time.Now()
	ts.StatementStatus = TrialStatementStatusApproved
	ts.ApprovedBy = &approvedBy
	ts.ApprovedAt = &now
	ts.ApprovalRemarks = remarks
	ts.UpdatedAt = now
}

// Reject marks the trial statement as rejected
func (ts *TrialStatement) Reject(rejectedBy string, remarks string) {
	now := time.Now()
	ts.StatementStatus = TrialStatementStatusRejected
	ts.ApprovedBy = &rejectedBy
	ts.ApprovedAt = &now
	ts.ApprovalRemarks = &remarks
	ts.UpdatedAt = now
}

// TrialStatementFilter represents filter criteria for querying trial statements
type TrialStatementFilter struct {
	BatchID        *string
	AgentID        *string
	Circle         *string
	CommissionType *CommissionType
	Status         *TrialStatementStatus
	Page           int
	Limit          int
}
