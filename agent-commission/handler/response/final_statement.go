package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"time"
)

// FinalStatementSummary represents final statement summary data
type FinalStatementSummary struct {
	FinalStatementID  int64   `json:"final_statement_id"`
	TrialStatementID  int64   `json:"trial_statement_id"`
	AgentCode         string  `json:"agent_code"`
	NetPayable        float64 `json:"net_payable"`
	Status            string  `json:"status"`
	GeneratedAt       string  `json:"generated_at"` // ISO 8601
}

// NewFinalStatementSummary converts domain.FinalStatement to FinalStatementSummary
func NewFinalStatementSummary(stmt domain.FinalStatement) FinalStatementSummary {
	return FinalStatementSummary{
		FinalStatementID: stmt.FinalStatementID,
		TrialStatementID: stmt.TrialStatementID,
		AgentCode:        stmt.AgentID,
		NetPayable:       stmt.TotalNetCommission,
		Status:           string(stmt.StatementStatus),
		GeneratedAt:      stmt.CreatedAt.Format(time.RFC3339),
	}
}

// NewFinalStatementSummaries converts a slice of domain.FinalStatement to []FinalStatementSummary
func NewFinalStatementSummaries(stmts []domain.FinalStatement) []FinalStatementSummary {
	resp := make([]FinalStatementSummary, 0, len(stmts))
	for _, stmt := range stmts {
		resp = append(resp, NewFinalStatementSummary(stmt))
	}
	return resp
}

// FinalStatementsListResponse represents the response for GET /commissions/final-statements
type FinalStatementsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Statements []FinalStatementSummary `json:"statements"`
		Pagination PaginationMetadata      `json:"pagination"`
	} `json:"data"`
}
