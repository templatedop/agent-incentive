package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"time"
)

// TrialStatementSummary represents trial statement summary data
type TrialStatementSummary struct {
	StatementID      int64   `json:"statement_id"`
	AgentCode        string  `json:"agent_code"`
	AgentName        string  `json:"agent_name"`
	Circle           string  `json:"circle,omitempty"`
	Division         string  `json:"division,omitempty"`
	GrossCommission  float64 `json:"gross_commission"`
	TDSAmount        float64 `json:"tds_amount"`
	NetAmount        float64 `json:"net_amount"`
	PolicyCount      int     `json:"policy_count"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"` // ISO 8601
}

// NewTrialStatementSummary converts domain.TrialStatement to TrialStatementSummary
func NewTrialStatementSummary(stmt domain.TrialStatement) TrialStatementSummary {
	return TrialStatementSummary{
		StatementID:     stmt.TrialStatementID,
		AgentCode:       stmt.AgentID,
		AgentName:       "", // TODO: Join with agent profile
		GrossCommission: stmt.TotalGrossCommission,
		TDSAmount:       stmt.TotalTDS,
		NetAmount:       stmt.TotalNetCommission,
		PolicyCount:     stmt.TotalPolicies,
		Status:          string(stmt.StatementStatus),
		CreatedAt:       stmt.CreatedAt.Format(time.RFC3339),
	}
}

// NewTrialStatementSummaries converts a slice of domain.TrialStatement to []TrialStatementSummary
func NewTrialStatementSummaries(stmts []domain.TrialStatement) []TrialStatementSummary {
	resp := make([]TrialStatementSummary, 0, len(stmts))
	for _, stmt := range stmts {
		resp = append(resp, NewTrialStatementSummary(stmt))
	}
	return resp
}

// PaginationMetadata represents pagination information
type PaginationMetadata struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalCount int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
}

// TrialStatementsListResponse represents the response for GET /commissions/trial-statements
type TrialStatementsListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Statements []TrialStatementSummary `json:"statements"`
		Pagination PaginationMetadata      `json:"pagination"`
	} `json:"data"`
}

// ApprovalResponse represents the response for trial statement approval
type ApprovalResponse struct {
	StatementID int64  `json:"statement_id"`
	Status      string `json:"status"`
	ApprovedBy  string `json:"approved_by"`
	ApprovedAt  string `json:"approved_at"` // ISO 8601
	Message     string `json:"message"`
}

// ApproveTrialStatementResponse represents the response for POST /commissions/trial-statements/{statementId}/approve
type ApproveTrialStatementResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      ApprovalResponse `json:"data"`
}
