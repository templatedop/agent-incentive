package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"fmt"
)

// AdvisorCoordinatorSummary represents coordinator information for lookup
type AdvisorCoordinatorSummary struct {
	CoordinatorID string `json:"coordinator_id"`
	AgentCode     string `json:"agent_code"`
	Name          string `json:"name"`
	Circle        string `json:"circle"`
	CircleName    string `json:"circle_name"`
	Division      string `json:"division"`
	DivisionName  string `json:"division_name"`
}

// NewAdvisorCoordinatorSummary converts domain.AgentProfile to AdvisorCoordinatorSummary
func NewAdvisorCoordinatorSummary(profile domain.AgentProfile) AdvisorCoordinatorSummary {
	// Combine first, middle, last name
	name := profile.FirstName
	if profile.MiddleName != nil && *profile.MiddleName != "" {
		name += " " + *profile.MiddleName
	}
	name += " " + profile.LastName

	return AdvisorCoordinatorSummary{
		CoordinatorID: fmt.Sprintf("%d", profile.AgentProfileID),
		AgentCode:     profile.AgentCode,
		Name:          name,
		Circle:        fmt.Sprintf("%d", profile.CircleID),
		CircleName:    profile.CircleName,
		Division:      fmt.Sprintf("%d", profile.DivisionID),
		DivisionName:  profile.DivisionName,
	}
}

// NewAdvisorCoordinatorSummaries converts a slice of agent profiles to coordinator summaries
func NewAdvisorCoordinatorSummaries(profiles []domain.AgentProfile) []AdvisorCoordinatorSummary {
	resp := make([]AdvisorCoordinatorSummary, 0, len(profiles))
	for _, profile := range profiles {
		resp = append(resp, NewAdvisorCoordinatorSummary(profile))
	}
	return resp
}

// AdvisorCoordinatorsResponse represents the response for GET /lookup/advisor-coordinators
type AdvisorCoordinatorsResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Coordinators []AdvisorCoordinatorSummary `json:"coordinators"`
	} `json:"data"`
}

// CircleSummary represents circle information for lookup
type CircleSummary struct {
	CircleCode string `json:"circle_code"`
	CircleName string `json:"circle_name"`
}

// NewCircleSummary converts domain.Circle to CircleSummary
func NewCircleSummary(circle domain.Circle) CircleSummary {
	return CircleSummary{
		CircleCode: circle.CircleCode,
		CircleName: circle.CircleName,
	}
}

// NewCircleSummaries converts a slice of circles to circle summaries
func NewCircleSummaries(circles []domain.Circle) []CircleSummary {
	resp := make([]CircleSummary, 0, len(circles))
	for _, circle := range circles {
		resp = append(resp, NewCircleSummary(circle))
	}
	return resp
}

// CirclesResponse represents the response for GET /lookup/circles
type CirclesResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		Circles []CircleSummary `json:"circles"`
	} `json:"data"`
}
