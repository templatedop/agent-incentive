package handler

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	resp "agent-commission/handler/response"
	"agent-commission/repo/postgres"

	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// FinalStatementHandler handles final statement operations
// Implements: FR-IC-COM-008 (Final statement generation)
type FinalStatementHandler struct {
	*serverHandler.Base
	finalRepo *postgres.FinalStatementRepository
}

// NewFinalStatementHandler creates a new final statement handler
func NewFinalStatementHandler(finalRepo *postgres.FinalStatementRepository) *FinalStatementHandler {
	base := serverHandler.New("FinalStatement").SetPrefix("/v1").AddPrefix("")
	return &FinalStatementHandler{Base: base, finalRepo: finalRepo}
}

// Routes returns all routes for final statement endpoints
func (h *FinalStatementHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.GET("/commissions/final-statements", h.ListFinalStatements).Name("List Final Statements"),
	}
}

// ListFinalStatementsRequest represents query parameters for GET /commissions/final-statements
type ListFinalStatementsRequest struct {
	BatchID string `form:"batch_id" validate:"omitempty"`
	AgentID string `form:"agent_id" validate:"omitempty"`
	Status  string `form:"status" validate:"omitempty,oneof=READY_FOR_DISBURSEMENT DISBURSED"`
	Page    int    `form:"page" validate:"omitempty,min=1"`
	Limit   int    `form:"limit" validate:"omitempty,min=1,max=100"`
}

// ListFinalStatements retrieves final statements ready for disbursement
// Implements: FR-IC-COM-008, WF-IC-COM-001 (Step 6)
// API: GET /commissions/final-statements
//
// Query Parameters:
//   - batch_id (optional): Filter by batch ID
//   - agent_id (optional): Filter by agent ID
//   - status (optional): READY_FOR_DISBURSEMENT, DISBURSED
//   - page (optional): Page number (default: 1)
//   - limit (optional): Items per page (default: 10, max: 100)
//
// Returns:
//   - 200: Final statements retrieved
//   - 401: Unauthorized
func (h *FinalStatementHandler) ListFinalStatements(
	sctx *serverRoute.Context,
	req ListFinalStatementsRequest,
) (*resp.FinalStatementsListResponse, error) {
	// Build filter
	filter := &domain.FinalStatementFilter{
		Page:  req.Page,
		Limit: req.Limit,
	}

	if req.BatchID != "" {
		filter.BatchID = &req.BatchID
	}
	if req.AgentID != "" {
		filter.AgentID = &req.AgentID
	}
	if req.Status != "" {
		status := domain.FinalStatementStatus(req.Status)
		filter.Status = &status
	}

	// Set defaults
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	// Fetch final statements
	statements, totalCount, err := h.finalRepo.SearchFinalStatements(sctx.Ctx, filter)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching final statements: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Retrieved %d final statements (total: %d)", len(statements), totalCount)

	// Calculate pagination
	totalPages := int(totalCount) / filter.Limit
	if int(totalCount)%filter.Limit > 0 {
		totalPages++
	}

	// Build response
	response := &resp.FinalStatementsListResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}
	response.Data.Statements = resp.NewFinalStatementSummaries(statements)
	response.Data.Pagination = resp.PaginationMetadata{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}

	return response, nil
}
