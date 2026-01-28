package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"fmt"
	"time"
)

// CommissionBatchResponse represents the commission batch data in API responses
type CommissionBatchResponse struct {
	BatchID             string  `json:"batch_id"`
	Status              string  `json:"status"`
	WorkflowID          string  `json:"workflow_id"`
	WorkflowState       string  `json:"workflow_state,omitempty"`
	EstimatedCompletion *string `json:"estimated_completion,omitempty"` // ISO 8601
	SLADeadline         string  `json:"sla_deadline"`                   // ISO 8601
	CreatedAt           string  `json:"created_at"`                     // ISO 8601
}

// NewCommissionBatchResponse converts domain.CommissionBatch to CommissionBatchResponse
func NewCommissionBatchResponse(batch domain.CommissionBatch) CommissionBatchResponse {
	resp := CommissionBatchResponse{
		BatchID:       batch.BatchID,
		Status:        string(batch.Status),
		WorkflowID:    batch.WorkflowID,
		WorkflowState: batch.WorkflowState,
		SLADeadline:   batch.SLADeadline.Format(time.RFC3339),
		CreatedAt:     batch.CreatedAt.Format(time.RFC3339),
	}

	if batch.EstimatedCompletion != nil {
		estCompletion := batch.EstimatedCompletion.Format(time.RFC3339)
		resp.EstimatedCompletion = &estCompletion
	}

	return resp
}

// StartCommissionBatchResponse represents the response for POST /commissions/batches/start
type StartCommissionBatchResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      CommissionBatchResponse `json:"data"`
}

// SLATracking represents SLA tracking information
type SLATracking struct {
	SLADeadline   string `json:"sla_deadline"`
	TimeRemaining string `json:"time_remaining"` // Human-readable format
	IsBreached    bool   `json:"is_breached"`
}

// BatchProgress represents batch processing progress
type BatchProgress struct {
	TotalPolicies      int `json:"total_policies"`
	ProcessedRecords   int `json:"processed_records"`
	FailedRecords      int `json:"failed_records"`
	ProgressPercentage int `json:"progress_percentage"`
}

// CommissionBatchStatusResponse represents detailed batch status
type CommissionBatchStatusResponse struct {
	BatchID       string        `json:"batch_id"`
	Status        string        `json:"status"`
	WorkflowState string        `json:"workflow_state,omitempty"`
	SLATracking   SLATracking   `json:"sla_tracking"`
	Progress      BatchProgress `json:"progress"`
	StartedAt     string        `json:"started_at"`      // ISO 8601
	CompletedAt   *string       `json:"completed_at,omitempty"` // ISO 8601
}

// NewCommissionBatchStatusResponse converts domain.CommissionBatch to detailed status response
func NewCommissionBatchStatusResponse(batch domain.CommissionBatch) CommissionBatchStatusResponse {
	// Calculate time remaining
	timeRemaining := time.Until(batch.SLADeadline)
	timeRemainingStr := formatDuration(timeRemaining)
	if timeRemaining < 0 {
		timeRemainingStr = "SLA breached"
	}

	slaTracking := SLATracking{
		SLADeadline:   batch.SLADeadline.Format(time.RFC3339),
		TimeRemaining: timeRemainingStr,
		IsBreached:    batch.IsSLABreached(),
	}

	progress := BatchProgress{
		TotalPolicies:      batch.TotalPolicies,
		ProcessedRecords:   batch.ProcessedRecords,
		FailedRecords:      batch.FailedRecords,
		ProgressPercentage: batch.ProgressPercentage,
	}

	resp := CommissionBatchStatusResponse{
		BatchID:       batch.BatchID,
		Status:        string(batch.Status),
		WorkflowState: batch.WorkflowState,
		SLATracking:   slaTracking,
		Progress:      progress,
		StartedAt:     batch.StartedAt.Format(time.RFC3339),
	}

	if batch.CompletedAt != nil {
		completedAt := batch.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &completedAt
	}

	return resp
}

// GetCommissionBatchStatusResponse represents the response for GET /commissions/batches/{batchId}/status
type GetCommissionBatchStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      CommissionBatchStatusResponse `json:"data"`
}

// formatDuration formats a duration to human-readable format
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "0h 0m"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	return fmt.Sprintf("%dh %dm", hours, minutes)
}
