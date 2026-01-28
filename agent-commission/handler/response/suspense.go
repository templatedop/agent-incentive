package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"time"
)

// CreateSuspenseAccountResponse represents the response for creating a suspense account
type CreateSuspenseAccountResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		SuspenseID         int64   `json:"suspense_id"`
		AgentID            *string `json:"agent_id,omitempty"`
		PolicyNumber       *string `json:"policy_number,omitempty"`
		SuspenseAmount     float64 `json:"suspense_amount"`
		SuspenseReason     string  `json:"suspense_reason"`
		SuspenseStatus     string  `json:"suspense_status"`
		Priority           string  `json:"priority"`
		ResolutionDeadline string  `json:"resolution_deadline"`
		WorkflowID         *string `json:"workflow_id,omitempty"`
	} `json:"data"`
}

// SuspenseAccountDetailResponse represents detailed suspense account information
type SuspenseAccountDetailResponse struct {
	SuspenseID         int64   `json:"suspense_id"`
	AgentID            *string `json:"agent_id,omitempty"`
	PolicyNumber       *string `json:"policy_number,omitempty"`
	CommissionID       *int64  `json:"commission_id,omitempty"`
	DisbursementID     *int64  `json:"disbursement_id,omitempty"`
	SuspenseAmount     float64 `json:"suspense_amount"`
	SuspenseReason     string  `json:"suspense_reason"`
	SuspenseStatus     string  `json:"suspense_status"`
	SuspenseDate       string  `json:"suspense_date"`
	AgingDays          int     `json:"aging_days"`
	AgingBucket        string  `json:"aging_bucket"`
	ResolutionDeadline *string `json:"resolution_deadline,omitempty"`
	IsOverdue          bool    `json:"is_overdue"`
	ResolvedDate       *string `json:"resolved_date,omitempty"`
	ResolutionMethod   *string `json:"resolution_method,omitempty"`
	ResolvedAmount     *float64 `json:"resolved_amount,omitempty"`
	WriteOffDate       *string `json:"write_off_date,omitempty"`
	WriteOffReason     *string `json:"write_off_reason,omitempty"`
	AssignedTo         *string `json:"assigned_to,omitempty"`
	Priority           string  `json:"priority"`
	Notes              *string `json:"notes,omitempty"`
	WorkflowID         *string `json:"workflow_id,omitempty"`
	WorkflowState      *string `json:"workflow_state,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	Version            int     `json:"version"`
}

// SuspenseAccountSummary represents a summary for list views
type SuspenseAccountSummary struct {
	SuspenseID         int64   `json:"suspense_id"`
	AgentID            *string `json:"agent_id,omitempty"`
	PolicyNumber       *string `json:"policy_number,omitempty"`
	SuspenseAmount     float64 `json:"suspense_amount"`
	SuspenseReason     string  `json:"suspense_reason"`
	SuspenseStatus     string  `json:"suspense_status"`
	AgingDays          int     `json:"aging_days"`
	AgingBucket        string  `json:"aging_bucket"`
	Priority           string  `json:"priority"`
	IsOverdue          bool    `json:"is_overdue"`
	ResolutionDeadline *string `json:"resolution_deadline,omitempty"`
}

// SuspenseAgingBucketResponse represents aging bucket data
type SuspenseAgingBucketResponse struct {
	AgingBucket  string  `json:"aging_bucket"`
	Count        int     `json:"count"`
	TotalAmount  float64 `json:"total_amount"`
	MinAgingDays int     `json:"min_aging_days"`
	MaxAgingDays int     `json:"max_aging_days"`
}

// SuspenseAgingStatsResponse represents statistics
type SuspenseAgingStatsResponse struct {
	Count       int     `json:"count"`
	TotalAmount float64 `json:"total_amount"`
	AvgAgeDays  float64 `json:"avg_age_days"`
}

// SuspenseAgingReportResponse represents the full aging report
type SuspenseAgingReportResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		ReportDate    string                        `json:"report_date"`
		TotalEntries  int                           `json:"total_entries"`
		TotalAmount   float64                       `json:"total_amount"`
		Buckets       []SuspenseAgingBucketResponse `json:"buckets"`
		OverdueCount  int                           `json:"overdue_count"`
		OverdueAmount float64                       `json:"overdue_amount"`
	} `json:"data"`
}

// SuspenseSearchResponse represents search results
type SuspenseSearchResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		SuspenseAccounts []SuspenseAccountSummary `json:"suspense_accounts"`
		Pagination       PaginationMetadata       `json:"pagination"`
	} `json:"data"`
}

// NewSuspenseAccountDetailResponse creates a detailed suspense account response
func NewSuspenseAccountDetailResponse(suspense domain.SuspenseAccount) SuspenseAccountDetailResponse {
	detail := SuspenseAccountDetailResponse{
		SuspenseID:     suspense.SuspenseID,
		AgentID:        suspense.AgentID,
		PolicyNumber:   suspense.PolicyNumber,
		CommissionID:   suspense.CommissionID,
		DisbursementID: suspense.DisbursementID,
		SuspenseAmount: suspense.SuspenseAmount,
		SuspenseReason: string(suspense.SuspenseReason),
		SuspenseStatus: string(suspense.SuspenseStatus),
		SuspenseDate:   suspense.SuspenseDate.Format(time.RFC3339),
		AgingDays:      suspense.AgingDays,
		AgingBucket:    suspense.GetAgingBucket(),
		IsOverdue:      suspense.IsOverdue(),
		ResolvedAmount: suspense.ResolvedAmount,
		ResolutionMethod: suspense.ResolutionMethod,
		WriteOffReason: suspense.WriteOffReason,
		AssignedTo:     suspense.AssignedTo,
		Priority:       suspense.Priority,
		Notes:          suspense.Notes,
		WorkflowID:     suspense.WorkflowID,
		WorkflowState:  suspense.WorkflowState,
		CreatedAt:      suspense.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      suspense.UpdatedAt.Format(time.RFC3339),
		Version:        suspense.Version,
	}

	if suspense.ResolutionDeadline != nil {
		deadlineStr := suspense.ResolutionDeadline.Format(time.RFC3339)
		detail.ResolutionDeadline = &deadlineStr
	}
	if suspense.ResolvedDate != nil {
		resolvedDateStr := suspense.ResolvedDate.Format(time.RFC3339)
		detail.ResolvedDate = &resolvedDateStr
	}
	if suspense.WriteOffDate != nil {
		writeOffDateStr := suspense.WriteOffDate.Format(time.RFC3339)
		detail.WriteOffDate = &writeOffDateStr
	}

	return detail
}

// NewSuspenseAccountSummary creates a summary suspense account response
func NewSuspenseAccountSummary(suspense domain.SuspenseAccount) SuspenseAccountSummary {
	summary := SuspenseAccountSummary{
		SuspenseID:     suspense.SuspenseID,
		AgentID:        suspense.AgentID,
		PolicyNumber:   suspense.PolicyNumber,
		SuspenseAmount: suspense.SuspenseAmount,
		SuspenseReason: string(suspense.SuspenseReason),
		SuspenseStatus: string(suspense.SuspenseStatus),
		AgingDays:      suspense.AgingDays,
		AgingBucket:    suspense.GetAgingBucket(),
		Priority:       suspense.Priority,
		IsOverdue:      suspense.IsOverdue(),
	}

	if suspense.ResolutionDeadline != nil {
		deadlineStr := suspense.ResolutionDeadline.Format(time.RFC3339)
		summary.ResolutionDeadline = &deadlineStr
	}

	return summary
}

// NewSuspenseAccountSummaries creates a list of suspense account summaries
func NewSuspenseAccountSummaries(suspenseAccounts []domain.SuspenseAccount) []SuspenseAccountSummary {
	summaries := make([]SuspenseAccountSummary, len(suspenseAccounts))
	for i, suspense := range suspenseAccounts {
		summaries[i] = NewSuspenseAccountSummary(suspense)
	}
	return summaries
}

// NewSuspenseAgingBucketResponse creates aging bucket response
func NewSuspenseAgingBucketResponse(bucket domain.SuspenseAgingBucket) SuspenseAgingBucketResponse {
	return SuspenseAgingBucketResponse{
		AgingBucket:  bucket.AgingBucket,
		Count:        bucket.Count,
		TotalAmount:  bucket.TotalAmount,
		MinAgingDays: bucket.MinAgingDays,
		MaxAgingDays: bucket.MaxAgingDays,
	}
}

// NewSuspenseAgingBucketResponses creates a list of aging bucket responses
func NewSuspenseAgingBucketResponses(buckets []domain.SuspenseAgingBucket) []SuspenseAgingBucketResponse {
	responses := make([]SuspenseAgingBucketResponse, len(buckets))
	for i, bucket := range buckets {
		responses[i] = NewSuspenseAgingBucketResponse(bucket)
	}
	return responses
}
