package domain

import (
	"time"
)

// SuspenseStatus represents the status of a suspense account entry
type SuspenseStatus string

const (
	SuspenseStatusOpen     SuspenseStatus = "OPEN"     // Active suspense entry
	SuspenseStatusResolved SuspenseStatus = "RESOLVED" // Issue resolved, cleared
	SuspenseStatusWriteOff SuspenseStatus = "WRITE_OFF" // Written off
)

// SuspenseReason represents why commission was moved to suspense
type SuspenseReason string

const (
	SuspenseReasonAgentNotFound          SuspenseReason = "AGENT_NOT_FOUND"
	SuspenseReasonInvalidBankDetails     SuspenseReason = "INVALID_BANK_DETAILS"
	SuspenseReasonPaymentFailed          SuspenseReason = "PAYMENT_FAILED"
	SuspenseReasonDocumentationIncomplete SuspenseReason = "DOCUMENTATION_INCOMPLETE"
	SuspenseReasonLicenseExpired         SuspenseReason = "LICENSE_EXPIRED"
	SuspenseReasonLicenseSuspended       SuspenseReason = "LICENSE_SUSPENDED"
	SuspenseReasonKYCIncomplete          SuspenseReason = "KYC_INCOMPLETE"
	SuspenseReasonDuplicatePayment       SuspenseReason = "DUPLICATE_PAYMENT"
	SuspenseReasonDisputeUnderReview     SuspenseReason = "DISPUTE_UNDER_REVIEW"
	SuspenseReasonOther                  SuspenseReason = "OTHER"
)

// SuspenseAccount represents a commission suspense account entry
// FR-IC-COM-013: Suspense account management
// BR-IC-SUS-001: Commission moved to suspense if agent details incomplete/invalid
type SuspenseAccount struct {
	SuspenseID          int64          `db:"suspense_id"`
	AgentID             *string        `db:"agent_id"`             // May be null if agent not found
	PolicyNumber        *string        `db:"policy_number"`        // Source policy
	CommissionID        *int64         `db:"commission_id"`        // Source commission transaction
	DisbursementID      *int64         `db:"disbursement_id"`      // Failed disbursement
	SuspenseAmount      float64        `db:"suspense_amount"`      // Amount held in suspense
	SuspenseReason      SuspenseReason `db:"suspense_reason"`
	SuspenseStatus      SuspenseStatus `db:"suspense_status"`
	SuspenseDate        time.Time      `db:"suspense_date"`        // When moved to suspense
	AgingDays           int            `db:"aging_days"`           // Days in suspense (calculated)
	ResolutionDeadline  *time.Time     `db:"resolution_deadline"`  // Target resolution date
	ResolvedDate        *time.Time     `db:"resolved_date"`
	ResolutionMethod    *string        `db:"resolution_method"`    // How it was resolved
	ResolvedAmount      *float64       `db:"resolved_amount"`      // Amount paid out on resolution
	WriteOffDate        *time.Time     `db:"write_off_date"`
	WriteOffReason      *string        `db:"write_off_reason"`
	AssignedTo          *string        `db:"assigned_to"`          // User assigned to resolve
	Priority            string         `db:"priority"`             // HIGH, MEDIUM, LOW
	Notes               *string        `db:"notes"`
	WorkflowID          *string        `db:"workflow_id"`
	WorkflowState       *string        `db:"workflow_state"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
	CreatedBy           string         `db:"created_by"`
	UpdatedBy           string         `db:"updated_by"`
	DeletedAt           *time.Time     `db:"deleted_at"`
	Version             int            `db:"version"`
}

// SuspenseTransaction represents audit trail for suspense account operations
type SuspenseTransaction struct {
	TransactionID   int64      `db:"transaction_id"`
	SuspenseID      int64      `db:"suspense_id"`
	TransactionType string     `db:"transaction_type"` // CREATED, UPDATED, RESOLVED, WRITE_OFF
	OldStatus       *string    `db:"old_status"`
	NewStatus       string     `db:"new_status"`
	Amount          *float64   `db:"amount"`
	Remarks         *string    `db:"remarks"`
	PerformedBy     string     `db:"performed_by"`
	TransactionDate time.Time  `db:"transaction_date"`
	CreatedAt       time.Time  `db:"created_at"`
}

// IsOpen checks if the suspense entry is still open
func (s *SuspenseAccount) IsOpen() bool {
	return s.SuspenseStatus == SuspenseStatusOpen
}

// IsOverdue checks if the suspense entry has exceeded resolution deadline
func (s *SuspenseAccount) IsOverdue() bool {
	if s.ResolutionDeadline == nil {
		return false
	}
	return time.Now().After(*s.ResolutionDeadline) && s.IsOpen()
}

// UpdateAgingDays calculates and updates the aging days
func (s *SuspenseAccount) UpdateAgingDays() {
	s.AgingDays = int(time.Since(s.SuspenseDate).Hours() / 24)
}

// GetAgingBucket returns the aging bucket for this suspense entry
// BR-IC-SUS-002: Aging buckets: 0-30, 31-60, 61-90, 91-180, 180+ days
func (s *SuspenseAccount) GetAgingBucket() string {
	s.UpdateAgingDays()
	switch {
	case s.AgingDays <= 30:
		return "0-30 days"
	case s.AgingDays <= 60:
		return "31-60 days"
	case s.AgingDays <= 90:
		return "61-90 days"
	case s.AgingDays <= 180:
		return "91-180 days"
	default:
		return "180+ days"
	}
}

// CanResolve checks if the suspense entry can be resolved
func (s *SuspenseAccount) CanResolve() bool {
	return s.SuspenseStatus == SuspenseStatusOpen
}

// CanWriteOff checks if the suspense entry can be written off
func (s *SuspenseAccount) CanWriteOff() bool {
	return s.SuspenseStatus == SuspenseStatusOpen
}

// SuspenseFilter represents search criteria for suspense accounts
type SuspenseFilter struct {
	AgentID        *string
	PolicyNumber   *string
	SuspenseStatus *SuspenseStatus
	SuspenseReason *SuspenseReason
	MinAmount      *float64
	MaxAmount      *float64
	FromDate       *time.Time
	ToDate         *time.Time
	MinAgingDays   *int
	MaxAgingDays   *int
	Priority       *string
	AssignedTo     *string
	Page           int
	Limit          int
}

// SuspenseAgingBucket represents aging report bucket
type SuspenseAgingBucket struct {
	AgingBucket  string  // "0-30 days", "31-60 days", etc.
	Count        int     // Number of entries
	TotalAmount  float64 // Total amount in bucket
	MinAgingDays int     // Minimum days in this bucket
	MaxAgingDays int     // Maximum days in this bucket
}

// SuspenseAgingReport represents the full aging report
type SuspenseAgingReport struct {
	ReportDate      time.Time
	TotalEntries    int
	TotalAmount     float64
	Buckets         []SuspenseAgingBucket
	ByReason        map[SuspenseReason]SuspenseAgingStats
	ByPriority      map[string]SuspenseAgingStats
	OverdueCount    int
	OverdueAmount   float64
}

// SuspenseAgingStats represents statistics for aging analysis
type SuspenseAgingStats struct {
	Count       int
	TotalAmount float64
	AvgAgeDays  float64
}

// CalculateResolutionDeadline calculates the target resolution date
// BR-IC-SUS-003: Resolution SLA based on priority
// HIGH: 7 days, MEDIUM: 15 days, LOW: 30 days
func CalculateResolutionDeadline(suspenseDate time.Time, priority string) time.Time {
	var days int
	switch priority {
	case "HIGH":
		days = 7
	case "MEDIUM":
		days = 15
	case "LOW":
		days = 30
	default:
		days = 30 // Default to LOW priority
	}
	return suspenseDate.AddDate(0, 0, days)
}

// DeterminePriority determines priority based on amount and reason
// BR-IC-SUS-004: Priority rules
func DeterminePriority(amount float64, reason SuspenseReason) string {
	// High priority for large amounts or critical reasons
	if amount >= 50000 || reason == SuspenseReasonDuplicatePayment || reason == SuspenseReasonDisputeUnderReview {
		return "HIGH"
	}
	// Medium priority for moderate amounts
	if amount >= 10000 {
		return "MEDIUM"
	}
	// Low priority for small amounts
	return "LOW"
}
