package domain

import (
	"time"
)

// BatchStatus represents the status of a commission batch
type BatchStatus string

const (
	BatchStatusInitiated      BatchStatus = "INITIATED"
	BatchStatusCalculating    BatchStatus = "CALCULATING"
	BatchStatusTrialGenerated BatchStatus = "TRIAL_GENERATED"
	BatchStatusCompleted      BatchStatus = "COMPLETED"
	BatchStatusFailed         BatchStatus = "FAILED"
)

// CommissionBatch represents a monthly commission calculation batch
// Implements: FR-IC-COM-002 (Commission batch processing)
// Business Rule: BR-IC-COM-001 (Commission calculation), BR-IC-COM-012 (6-hour SLA)
type CommissionBatch struct {
	BatchID             string      `json:"batch_id" db:"batch_id"`
	Month               int         `json:"month" db:"month"`
	Year                int         `json:"year" db:"year"`
	Status              BatchStatus `json:"status" db:"status"`
	WorkflowID          string      `json:"workflow_id" db:"workflow_id"`
	WorkflowState       string      `json:"workflow_state" db:"workflow_state"`
	TotalPolicies       int         `json:"total_policies" db:"total_policies"`
	ProcessedRecords    int         `json:"processed_records" db:"processed_records"`
	FailedRecords       int         `json:"failed_records" db:"failed_records"`
	ProgressPercentage  int         `json:"progress_percentage" db:"progress_percentage"`
	TriggeredBy         string      `json:"triggered_by" db:"triggered_by"` // SYSTEM_SCHEDULER or MANUAL
	StartedAt           time.Time   `json:"started_at" db:"started_at"`
	EstimatedCompletion *time.Time  `json:"estimated_completion,omitempty" db:"estimated_completion"`
	SLADeadline         time.Time   `json:"sla_deadline" db:"sla_deadline"` // 6 hours from start
	CompletedAt         *time.Time  `json:"completed_at,omitempty" db:"completed_at"`
	FailureReason       *string     `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt           time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at" db:"updated_at"`
}

// CalculateProgress updates the progress percentage
func (cb *CommissionBatch) CalculateProgress() {
	if cb.TotalPolicies > 0 {
		cb.ProgressPercentage = (cb.ProcessedRecords * 100) / cb.TotalPolicies
	}
}

// IsSLABreached checks if the 6-hour SLA deadline is breached
func (cb *CommissionBatch) IsSLABreached() bool {
	return time.Now().After(cb.SLADeadline)
}

// CommissionType represents the type of commission
type CommissionType string

const (
	CommissionTypeFirstYear CommissionType = "FIRST_YEAR"
	CommissionTypeRenewal   CommissionType = "RENEWAL"
	CommissionTypeBonus     CommissionType = "BONUS"
)

// CommissionStatus represents the commission lifecycle status
type CommissionStatus string

const (
	CommissionStatusCalculated          CommissionStatus = "CALCULATED"
	CommissionStatusTrialGenerated      CommissionStatus = "TRIAL_GENERATED"
	CommissionStatusFinalized           CommissionStatus = "FINALIZED"
	CommissionStatusReadyForDisbursement CommissionStatus = "READY_FOR_DISBURSEMENT"
	CommissionStatusDisbursed           CommissionStatus = "DISBURSED"
	CommissionStatusCancelled           CommissionStatus = "CANCELLED"
)

// CommissionTransaction represents individual commission calculation
// Implements: BR-IC-COM-001, BR-IC-COM-002, BR-IC-COM-003 (TDS), BR-IC-COM-006 (Rate)
type CommissionTransaction struct {
	CommissionID         int64            `json:"commission_id" db:"commission_id"`
	AgentID              string           `json:"agent_id" db:"agent_id"`
	PolicyNumber         string           `json:"policy_number" db:"policy_number"`
	CommissionType       CommissionType   `json:"commission_type" db:"commission_type"`
	ProductType          ProductType      `json:"product_type" db:"product_type"`
	AnnualisedPremium    float64          `json:"annualised_premium" db:"annualised_premium"`
	RatePercentage       float64          `json:"rate_percentage" db:"rate_percentage"`
	GrossCommission      float64          `json:"gross_commission" db:"gross_commission"`
	TDSRate              float64          `json:"tds_rate" db:"tds_rate"`
	TDSAmount            float64          `json:"tds_amount" db:"tds_amount"`
	NetCommission        float64          `json:"net_commission" db:"net_commission"`
	CommissionDate       time.Time        `json:"commission_date" db:"commission_date"`
	CommissionStatus     CommissionStatus `json:"commission_status" db:"commission_status"`
	TrialStatementID     *int64           `json:"trial_statement_id,omitempty" db:"trial_statement_id"`
	FinalStatementID     *int64           `json:"final_statement_id,omitempty" db:"final_statement_id"`
	DisbursementID       *string          `json:"disbursement_id,omitempty" db:"disbursement_id"`
	BatchID              *string          `json:"batch_id,omitempty" db:"batch_id"`
	WorkflowState        *string          `json:"workflow_state,omitempty" db:"workflow_state"`
	WorkflowStateUpdated *time.Time       `json:"workflow_state_updated_at,omitempty" db:"workflow_state_updated_at"`
	CreatedAt            time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at" db:"updated_at"`
	CreatedBy            *string          `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy            *string          `json:"updated_by,omitempty" db:"updated_by"`
	DeletedAt            *time.Time       `json:"deleted_at,omitempty" db:"deleted_at"`
	Version              int              `json:"version" db:"version"`
}

// CalculateNetCommission computes net commission after TDS
// BR-IC-COM-003: TDS deduction based on PAN availability
func (ct *CommissionTransaction) CalculateNetCommission() {
	ct.TDSAmount = (ct.GrossCommission * ct.TDSRate) / 100
	ct.NetCommission = ct.GrossCommission - ct.TDSAmount
}
