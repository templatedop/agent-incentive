package handler

import (
	"agent-commission/core/port"
	resp "agent-commission/handler/response"
	"agent-commission/repo/postgres"

	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// LookupHandler handles lookup and reference data operations
// Implements: FR-IC-PROF-001 (Agent profile management)
type LookupHandler struct {
	*serverHandler.Base
	refRepo *postgres.ReferenceDataRepository
}

// NewLookupHandler creates a new lookup handler
func NewLookupHandler(refRepo *postgres.ReferenceDataRepository) *LookupHandler {
	base := serverHandler.New("Lookup").SetPrefix("/v1").AddPrefix("")
	return &LookupHandler{Base: base, refRepo: refRepo}
}

// Routes returns all routes for lookup endpoints
func (h *LookupHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.GET("/lookup/advisor-coordinators", h.GetAdvisorCoordinators).Name("Get Advisor Coordinators"),
		serverRoute.GET("/lookup/circles", h.GetCircles).Name("Get Circles"),
	}
}

// GetAdvisorCoordinatorsRequest represents query parameters for GET /lookup/advisor-coordinators
type GetAdvisorCoordinatorsRequest struct {
	Circle   string `form:"circle" validate:"omitempty"`
	Division string `form:"division" validate:"omitempty"`
}

// GetAdvisorCoordinators retrieves active advisor coordinators for dropdown/selection
// Implements: FR-IC-PROF-001, BR-IC-AH-001
// API: GET /lookup/advisor-coordinators
//
// Query Parameters:
//   - circle (optional): Filter by circle
//   - division (optional): Filter by division
//
// Returns:
//   - 200: Coordinators retrieved successfully
//   - 401: Unauthorized
func (h *LookupHandler) GetAdvisorCoordinators(
	sctx *serverRoute.Context,
	req GetAdvisorCoordinatorsRequest,
) (*resp.AdvisorCoordinatorsResponse, error) {
	var coordinators []resp.AdvisorCoordinatorSummary
	var err error

	// TODO: Handle circle and division filters when needed
	// For now, return all coordinators
	profiles, err := h.refRepo.GetAllCoordinators(sctx.Ctx)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching advisor coordinators: %v", err)
		return nil, err
	}

	coordinators = resp.NewAdvisorCoordinatorSummaries(profiles)

	log.Info(sctx.Ctx, "Successfully retrieved %d advisor coordinators", len(coordinators))

	// Build response
	response := &resp.AdvisorCoordinatorsResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}
	response.Data.Coordinators = coordinators

	return response, nil
}

// GetCircles retrieves all circles for reference data
// API: GET /lookup/circles
//
// Returns:
//   - 200: Circles retrieved successfully
//   - 401: Unauthorized
func (h *LookupHandler) GetCircles(
	sctx *serverRoute.Context,
	_ struct{},
) (*resp.CirclesResponse, error) {
	circles, err := h.refRepo.GetAllCircles(sctx.Ctx)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching circles: %v", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Successfully retrieved %d circles", len(circles))

	// Build response
	response := &resp.CirclesResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}
	response.Data.Circles = resp.NewCircleSummaries(circles)

	return response, nil
}
