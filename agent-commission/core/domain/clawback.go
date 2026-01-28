package domain

import (
	"time"
)

// ClawbackStatus represents the status of a clawback
type ClawbackStatus string

const (
	ClawbackStatusPending    ClawbackStatus = "PENDING"     // Clawback initiated, not yet started recovery
	ClawbackStatusInProgress ClawbackStatus = "IN_PROGRESS" // Recovery in progress
	ClawbackStatusCompleted  ClawbackStatus = "COMPLETED"   // Fully recovered
	ClawbackStatusPartial    ClawbackStatus = "PARTIAL"     // Partially recovered, policy term ended
	ClawbackStatusWaived     ClawbackStatus = "WAIVED"      // Clawback waived by management
	ClawbackStatusWriteOff   ClawbackStatus = "WRITE_OFF"   // Written off as bad debt
)

// ClawbackReason represents why a clawback was triggered
type ClawbackReason string

const (
	ClawbackReasonPolicySurrendered ClawbackReason = "POLICY_SURRENDERED"
	ClawbackReasonPolicyLapsed      ClawbackReason = "POLICY_LAPSED"
	ClawbackReasonPolicyCancelled   ClawbackReason = "POLICY_CANCELLED"
	ClawbackReasonFraud             ClawbackReason = "FRAUD_DETECTED"
	ClawbackReasonChargebackRequest ClawbackReason = "CHARGEBACK_REQUEST"
)

// Clawback represents a commission clawback record
// FR-IC-COM-012: Clawback triggers and calculation
// BR-IC-CLAWBACK-001: Surrender within N years triggers graduated clawback
type Clawback struct {
	ClawbackID          int64          `db:"clawback_id"`
	PolicyNumber        string         `db:"policy_number"`
	AgentID             string         `db:"agent_id"`
	OriginalCommission  float64        `db:"original_commission"` // Total commission paid for this policy
	ClawbackAmount      float64        `db:"clawback_amount"`     // Total amount to recover
	RecoveredAmount     float64        `db:"recovered_amount"`    // Amount recovered so far
	PendingAmount       float64        `db:"pending_amount"`      // Remaining to recover
	ClawbackPercentage  float64        `db:"clawback_percentage"` // 0-100% based on policy age
	ClawbackReason      ClawbackReason `db:"clawback_reason"`
	ClawbackStatus      ClawbackStatus `db:"clawback_status"`
	PolicyAge           int            `db:"policy_age"`            // Age in months at clawback time
	TriggerDate         time.Time      `db:"trigger_date"`          // When clawback was triggered
	PolicyInceptionDate time.Time      `db:"policy_inception_date"` // Original policy start date
	PolicyEndDate       *time.Time     `db:"policy_end_date"`       // When policy terminated
	RecoveryStartDate   *time.Time     `db:"recovery_start_date"`   // When recovery began
	RecoveryEndDate     *time.Time     `db:"recovery_end_date"`     // When recovery completed
	RecoverySchedule    string         `db:"recovery_schedule"`     // IMMEDIATE, INSTALLMENT
	InstallmentMonths   *int           `db:"installment_months"`    // Number of months for recovery
	ApprovedBy          *string        `db:"approved_by"`           // User who approved clawback
	ApprovedAt          *time.Time     `db:"approved_at"`
	WaivedBy            *string        `db:"waived_by"`   // User who waived clawback
	WaivedAt            *time.Time     `db:"waived_at"`
	WaiverReason        *string        `db:"waiver_reason"`
	WorkflowID          *string        `db:"workflow_id"` // Temporal workflow tracking
	WorkflowState       *string        `db:"workflow_state"`
	Notes               *string        `db:"notes"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
	CreatedBy           string         `db:"created_by"`
	UpdatedBy           string         `db:"updated_by"`
	DeletedAt           *time.Time     `db:"deleted_at"`
	Version             int            `db:"version"` // Optimistic locking
}

// ClawbackRecovery represents individual recovery transactions
// Tracks installment-based recovery of clawback amounts
type ClawbackRecovery struct {
	RecoveryID         int64      `db:"recovery_id"`
	ClawbackID         int64      `db:"clawback_id"`
	InstallmentNumber  int        `db:"installment_number"`   // 1, 2, 3... for installment recovery
	ScheduledAmount    float64    `db:"scheduled_amount"`     // Amount scheduled for this installment
	RecoveredAmount    float64    `db:"recovered_amount"`     // Actual amount recovered
	RecoveryDate       time.Time  `db:"recovery_date"`        // When recovery happened
	RecoveryMethod     string     `db:"recovery_method"`      // DEDUCTION, CHEQUE_RETURN, etc.
	StatementID        *int64     `db:"statement_id"`         // Which statement this was deducted from
	DisbursementID     *int64     `db:"disbursement_id"`      // Which disbursement was adjusted
	TransactionRef     *string    `db:"transaction_ref"`      // External reference
	RecoveryStatus     string     `db:"recovery_status"`      // SCHEDULED, COMPLETED, FAILED
	FailureReason      *string    `db:"failure_reason"`
	RetryCount         int        `db:"retry_count"`
	NextRetryDate      *time.Time `db:"next_retry_date"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
	CreatedBy          string     `db:"created_by"`
	UpdatedBy          string     `db:"updated_by"`
	Version            int        `db:"version"`
}

// CalculateClawbackPercentage calculates the clawback percentage based on policy age
// BR-IC-CLAWBACK-002: Graduated clawback percentage
// Year 1: 100%, Year 2: 75%, Year 3: 50%, Year 4: 25%, Year 5+: 0%
func CalculateClawbackPercentage(policyAgeMonths int) float64 {
	switch {
	case policyAgeMonths < 12:
		return 100.0 // First year
	case policyAgeMonths < 24:
		return 75.0 // Second year
	case policyAgeMonths < 36:
		return 50.0 // Third year
	case policyAgeMonths < 48:
		return 25.0 // Fourth year
	default:
		return 0.0 // Fifth year onwards - no clawback
	}
}

// IsFullyRecovered checks if the clawback has been fully recovered
func (c *Clawback) IsFullyRecovered() bool {
	return c.ClawbackStatus == ClawbackStatusCompleted ||
		c.PendingAmount <= 0.01 // Use small epsilon for float comparison
}

// CanRecover checks if the clawback can still be recovered
func (c *Clawback) CanRecover() bool {
	return c.ClawbackStatus == ClawbackStatusPending ||
		c.ClawbackStatus == ClawbackStatusInProgress
}

// IsActive checks if the clawback is in an active state
func (c *Clawback) IsActive() bool {
	return c.ClawbackStatus != ClawbackStatusCompleted &&
		c.ClawbackStatus != ClawbackStatusWaived &&
		c.ClawbackStatus != ClawbackStatusWriteOff
}

// RecoveryProgress returns the recovery completion percentage (0-100)
func (c *Clawback) RecoveryProgress() float64 {
	if c.ClawbackAmount == 0 {
		return 0
	}
	return (c.RecoveredAmount / c.ClawbackAmount) * 100
}

// UpdateRecoveredAmount updates the recovered and pending amounts
func (c *Clawback) UpdateRecoveredAmount(amount float64) {
	c.RecoveredAmount += amount
	c.PendingAmount = c.ClawbackAmount - c.RecoveredAmount

	// Update status if fully recovered
	if c.PendingAmount <= 0.01 {
		c.ClawbackStatus = ClawbackStatusCompleted
		now := time.Now()
		c.RecoveryEndDate = &now
	} else if c.ClawbackStatus == ClawbackStatusPending {
		// Move to in-progress on first recovery
		c.ClawbackStatus = ClawbackStatusInProgress
		if c.RecoveryStartDate == nil {
			now := time.Now()
			c.RecoveryStartDate = &now
		}
	}
}

// ClawbackFilter represents search criteria for clawbacks
type ClawbackFilter struct {
	AgentID        *string
	PolicyNumber   *string
	ClawbackStatus *ClawbackStatus
	ClawbackReason *ClawbackReason
	FromDate       *time.Time
	ToDate         *time.Time
	MinAmount      *float64
	MaxAmount      *float64
	Page           int
	Limit          int
}
