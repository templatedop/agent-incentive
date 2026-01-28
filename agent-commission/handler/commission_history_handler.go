package handler

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"agent-commission/handler/response"
	repo "agent-commission/repo/postgres"

	log "gitlab.cept.gov.in/it-2.0-common/api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// CommissionHistoryHandler handles commission history inquiry operations
// FR-IC-COM-011: Commission history inquiry
type CommissionHistoryHandler struct {
	*serverHandler.Base
	repo *repo.CommissionHistoryRepository
}

// NewCommissionHistoryHandler creates a new commission history handler
func NewCommissionHistoryHandler(repo *repo.CommissionHistoryRepository) *CommissionHistoryHandler {
	base := serverHandler.New("CommissionHistory").SetPrefix("/v1").AddPrefix("")
	return &CommissionHistoryHandler{Base: base, repo: repo}
}

// Routes returns all routes for commission history endpoints
func (h *CommissionHistoryHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.GET("/commissions/history", h.SearchCommissionHistory).Name("Search Commission History"),
	}
}

// SearchCommissionHistoryRequest represents query parameters for history search
type SearchCommissionHistoryRequest struct {
	AgentID          *string `query:"agent_id"`
	PolicyNumber     *string `query:"policy_number"`
	CommissionType   *string `query:"commission_type"`   // FIRST_YEAR, RENEWAL, BONUS
	CommissionStatus *string `query:"commission_status"` // CALCULATED, DISBURSED, etc.
	FromDate         *string `query:"from_date"`         // RFC3339 format
	ToDate           *string `query:"to_date"`           // RFC3339 format
	BatchID          *string `query:"batch_id"`
	Page             int     `query:"page"`
	Limit            int     `query:"limit"`
}

// SearchCommissionHistory searches commission transaction history
// API: GET /commissions/history
// Implements: FR-IC-COM-011, BR-IC-COM-009
func (h *CommissionHistoryHandler) SearchCommissionHistory(
	sctx *serverRoute.Context,
	req SearchCommissionHistoryRequest,
) (*response.CommissionHistoryResponse, error) {
	// Build filter from request
	filter := &domain.CommissionHistoryFilter{
		AgentID:        req.AgentID,
		PolicyNumber:   req.PolicyNumber,
		BatchID:        req.BatchID,
		Page:           req.Page,
		Limit:          req.Limit,
	}

	// Parse enum types
	if req.CommissionType != nil && *req.CommissionType != "" {
		commType := domain.CommissionType(*req.CommissionType)
		filter.CommissionType = &commType
	}
	if req.CommissionStatus != nil && *req.CommissionStatus != "" {
		commStatus := domain.CommissionStatus(*req.CommissionStatus)
		filter.CommissionStatus = &commStatus
	}

	// TODO: Parse date filters
	// if req.FromDate != nil {
	//     fromDate, _ := time.Parse(time.RFC3339, *req.FromDate)
	//     filter.FromDate = &fromDate
	// }
	// if req.ToDate != nil {
	//     toDate, _ := time.Parse(time.RFC3339, *req.ToDate)
	//     filter.ToDate = &toDate
	// }

	// Search commission history
	transactions, totalCount, err := h.repo.SearchCommissionHistory(sctx.Ctx, filter)
	if err != nil {
		log.Error(sctx.Ctx, "Error searching commission history: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Found %d commission transactions (total: %d)", len(transactions), totalCount)

	// Build response
	resp := &response.CommissionHistoryResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}
	resp.Data.Transactions = response.NewCommissionHistorySummaries(transactions)
	resp.Data.Pagination = response.PaginationMetadata{
		CurrentPage: filter.Page,
		PageSize:    filter.Limit,
		TotalCount:  totalCount,
		TotalPages:  int((totalCount + int64(filter.Limit) - 1) / int64(filter.Limit)),
	}

	return resp, nil
}
