package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
)

// CommissionRateResponse represents the commission rate data in API responses
type CommissionRateResponse struct {
	RateID               int64   `json:"rate_id"`
	RatePercentage       float64 `json:"rate_percentage"`
	PolicyDurationMonths int     `json:"policy_duration_months"`
	ProductType          string  `json:"product_type"`
	PlanCode             string  `json:"plan_code"`
	AgentType            string  `json:"agent_type"`
	PolicyTermYears      int     `json:"policy_term_years"`
	EffectiveFrom        string  `json:"effective_from"` // ISO 8601 date format
	EffectiveTo          *string `json:"effective_to,omitempty"`
}

// NewCommissionRateResponse converts domain.CommissionRate to CommissionRateResponse
func NewCommissionRateResponse(rate domain.CommissionRate) CommissionRateResponse {
	resp := CommissionRateResponse{
		RateID:               rate.RateID,
		RatePercentage:       rate.RatePercentage,
		PolicyDurationMonths: rate.PolicyDurationMonths,
		ProductType:          string(rate.ProductType),
		PlanCode:             rate.PlanCode,
		AgentType:            rate.AgentType,
		PolicyTermYears:      rate.PolicyTermYears,
		EffectiveFrom:        rate.EffectiveFrom.Format("2006-01-02"),
	}

	if rate.EffectiveTo != nil {
		effectiveTo := rate.EffectiveTo.Format("2006-01-02")
		resp.EffectiveTo = &effectiveTo
	}

	return resp
}

// NewCommissionRatesResponse converts a slice of domain.CommissionRate to []CommissionRateResponse
func NewCommissionRatesResponse(rates []domain.CommissionRate) []CommissionRateResponse {
	resp := make([]CommissionRateResponse, 0, len(rates))
	for _, rate := range rates {
		resp = append(resp, NewCommissionRateResponse(rate))
	}
	return resp
}

// CommissionRatesListResponse represents the response for GET /lookup/commission-rates
// Implements: FR-IC-COM-001, BR-IC-COM-006
type CommissionRatesListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Rates []CommissionRateResponse `json:"rates"`
	} `json:"data"`
}
