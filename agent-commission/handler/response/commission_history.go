package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"time"
)

// CommissionHistorySummary represents a commission transaction in history
type CommissionHistorySummary struct {
	CommissionID       int64   `json:"commission_id"`
	AgentID            string  `json:"agent_id"`
	PolicyNumber       string  `json:"policy_number"`
	CommissionType     string  `json:"commission_type"`      // FIRST_YEAR, RENEWAL, BONUS
	ProductType        string  `json:"product_type"`         // PLI, RPLI
	AnnualisedPremium  float64 `json:"annualised_premium"`
	RatePercentage     float64 `json:"rate_percentage"`
	GrossCommission    float64 `json:"gross_commission"`
	TDSAmount          float64 `json:"tds_amount"`
	NetCommission      float64 `json:"net_commission"`
	CommissionDate     string  `json:"commission_date"`
	CommissionStatus   string  `json:"commission_status"`
	TrialStatementID   *int64  `json:"trial_statement_id,omitempty"`
	FinalStatementID   *int64  `json:"final_statement_id,omitempty"`
	DisbursementID     *string `json:"disbursement_id,omitempty"`
	BatchID            *string `json:"batch_id,omitempty"`
}

// NewCommissionHistorySummary converts domain.CommissionTransaction to summary
func NewCommissionHistorySummary(txn domain.CommissionTransaction) CommissionHistorySummary {
	return CommissionHistorySummary{
		CommissionID:      txn.CommissionID,
		AgentID:           txn.AgentID,
		PolicyNumber:      txn.PolicyNumber,
		CommissionType:    string(txn.CommissionType),
		ProductType:       string(txn.ProductType),
		AnnualisedPremium: txn.AnnualisedPremium,
		RatePercentage:    txn.RatePercentage,
		GrossCommission:   txn.GrossCommission,
		TDSAmount:         txn.TDSAmount,
		NetCommission:     txn.NetCommission,
		CommissionDate:    txn.CommissionDate.Format(time.RFC3339),
		CommissionStatus:  string(txn.CommissionStatus),
		TrialStatementID:  txn.TrialStatementID,
		FinalStatementID:  txn.FinalStatementID,
		DisbursementID:    txn.DisbursementID,
		BatchID:           txn.BatchID,
	}
}

// NewCommissionHistorySummaries converts a slice of transactions to summaries
func NewCommissionHistorySummaries(transactions []domain.CommissionTransaction) []CommissionHistorySummary {
	summaries := make([]CommissionHistorySummary, 0, len(transactions))
	for _, txn := range transactions {
		summaries = append(summaries, NewCommissionHistorySummary(txn))
	}
	return summaries
}

// CommissionHistoryResponse represents response for GET /commissions/history
type CommissionHistoryResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Transactions []CommissionHistorySummary `json:"transactions"`
		Pagination   PaginationMetadata         `json:"pagination"`
	} `json:"data"`
}
