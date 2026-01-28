package domain

import (
	"time"
)

// ProductType represents the type of insurance product
type ProductType string

const (
	ProductTypePLI  ProductType = "PLI"
	ProductTypeRPLI ProductType = "RPLI"
)

// CommissionRate represents a commission rate configuration
// Implements: BR-IC-COM-006 (Commission rate table structure)
// Requirement: FR-IC-COM-001 (Commission rate table management)
type CommissionRate struct {
	RateID               int64       `json:"rate_id" db:"rate_id"`
	RatePercentage       float64     `json:"rate_percentage" db:"rate_percentage"`
	PolicyDurationMonths int         `json:"policy_duration_months" db:"policy_duration_months"`
	ProductType          ProductType `json:"product_type" db:"product_type"`
	PlanCode             string      `json:"plan_code" db:"product_plan_code"` // Maps to product_plan_code in DB
	AgentType            string      `json:"agent_type" db:"agent_type"`
	PolicyTermYears      int         `json:"policy_term_years" db:"policy_term_years"`
	EffectiveFrom        time.Time   `json:"effective_from" db:"effective_from"`
	EffectiveTo          *time.Time  `json:"effective_to,omitempty" db:"effective_to"`
	IsActive             bool        `json:"is_active" db:"is_active"`
	CreatedAt            time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at" db:"updated_at"`
}

// IsActiveOn checks if the commission rate is active on a specific date
func (cr *CommissionRate) IsActiveOn(asOfDate time.Time) bool {
	if !cr.IsActive {
		return false
	}
	if asOfDate.Before(cr.EffectiveFrom) {
		return false
	}
	if cr.EffectiveTo != nil && asOfDate.After(*cr.EffectiveTo) {
		return false
	}
	return true
}

// CommissionRateHistory represents the change history of commission rates
// Used for audit trail and historical analysis
type CommissionRateHistory struct {
	HistoryID        string      `json:"history_id" db:"history_id"`
	RateID           string      `json:"rate_id" db:"rate_id"`
	OldPercentage    *float64    `json:"old_percentage,omitempty" db:"old_percentage"`
	NewPercentage    float64     `json:"new_percentage" db:"new_percentage"`
	ChangeReason     string      `json:"change_reason" db:"change_reason"`
	ChangedBy        string      `json:"changed_by" db:"changed_by"`
	EffectiveFrom    time.Time   `json:"effective_from" db:"effective_from"`
	OldEffectiveTo   *time.Time  `json:"old_effective_to,omitempty" db:"old_effective_to"`
	ProductType      ProductType `json:"product_type" db:"product_type"`
	PlanCode         string      `json:"plan_code" db:"plan_code"`
	AgentType        string      `json:"agent_type" db:"agent_type"`
	PolicyTermYears  int         `json:"policy_term_years" db:"policy_term_years"`
	ChangedAt        time.Time   `json:"changed_at" db:"changed_at"`
}

// CommissionRateFilter represents filter criteria for querying commission rates
type CommissionRateFilter struct {
	ProductType *ProductType
	AgentType   *string
	PlanCode    *string
	AsOfDate    *time.Time // Filter for rates active on this date
}
