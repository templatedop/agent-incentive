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

// CommissionRateHandler handles commission rate lookup operations
// Implements: FR-IC-COM-001 (Commission rate table management)
type CommissionRateHandler struct {
	*serverHandler.Base
	repo *postgres.CommissionRateRepository
}

// NewCommissionRateHandler creates a new commission rate handler
func NewCommissionRateHandler(repo *postgres.CommissionRateRepository) *CommissionRateHandler {
	base := serverHandler.New("CommissionRates").SetPrefix("/v1").AddPrefix("")
	return &CommissionRateHandler{Base: base, repo: repo}
}

// Routes returns all routes for commission rate endpoints
func (h *CommissionRateHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.GET("/lookup/commission-rates", h.GetCommissionRates).Name("Get Commission Rates"),
	}
}

// GetCommissionRatesRequest represents query parameters for GET /lookup/commission-rates
type GetCommissionRatesRequest struct {
	ProductType string `form:"product_type" validate:"omitempty,oneof=PLI RPLI"`
	AgentType   string `form:"agent_type" validate:"omitempty"`
	PlanCode    string `form:"plan_code" validate:"omitempty"`
}

// GetCommissionRates retrieves commission rates based on filter criteria
// Implements: FR-IC-COM-001, BR-IC-COM-006
// API: GET /lookup/commission-rates
//
// Query Parameters:
//   - product_type (optional): Filter by product type (PLI, RPLI)
//   - agent_type (optional): Filter by agent type
//   - plan_code (optional): Filter by product plan code
//
// Returns:
//   - 200: Commission rates retrieved successfully
//   - 401: Unauthorized
func (h *CommissionRateHandler) GetCommissionRates(
	sctx *serverRoute.Context,
	req GetCommissionRatesRequest,
) (*resp.CommissionRatesListResponse, error) {
	// Build filter from query parameters
	filter := &domain.CommissionRateFilter{}

	if req.ProductType != "" {
		productType := domain.ProductType(req.ProductType)
		filter.ProductType = &productType
	}

	if req.AgentType != "" {
		filter.AgentType = &req.AgentType
	}

	if req.PlanCode != "" {
		filter.PlanCode = &req.PlanCode
	}

	// Fetch commission rates from repository
	rates, err := h.repo.GetCommissionRates(sctx.Ctx, filter)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching commission rates: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Successfully retrieved %d commission rates", len(rates))

	// Build response
	response := &resp.CommissionRatesListResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}
	response.Data.Rates = resp.NewCommissionRatesResponse(rates)

	return response, nil
}
