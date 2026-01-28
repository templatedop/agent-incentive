package domain

import (
	"time"
)

// WorkflowStatus represents the status of a Temporal workflow
type WorkflowStatus string

const (
	WorkflowStatusRunning   WorkflowStatus = "RUNNING"
	WorkflowStatusCompleted WorkflowStatus = "COMPLETED"
	WorkflowStatusFailed    WorkflowStatus = "FAILED"
	WorkflowStatusCanceled  WorkflowStatus = "CANCELED"
	WorkflowStatusTerminated WorkflowStatus = "TERMINATED"
	WorkflowStatusTimedOut   WorkflowStatus = "TIMED_OUT"
)

// WorkflowType represents different types of workflows in the system
type WorkflowType string

const (
	WorkflowTypeCommissionBatch      WorkflowType = "COMMISSION_BATCH"
	WorkflowTypeTrialStatementApproval WorkflowType = "TRIAL_STATEMENT_APPROVAL"
	WorkflowTypeDisbursement         WorkflowType = "DISBURSEMENT"
	WorkflowTypeClawback             WorkflowType = "CLAWBACK"
	WorkflowTypeSuspense             WorkflowType = "SUSPENSE"
)

// WorkflowInfo represents the current state of a workflow
type WorkflowInfo struct {
	WorkflowID    string         `json:"workflow_id"`
	WorkflowType  WorkflowType   `json:"workflow_type"`
	Status        WorkflowStatus `json:"status"`
	StartTime     time.Time      `json:"start_time"`
	CloseTime     *time.Time     `json:"close_time,omitempty"`
	ExecutionTime *int64         `json:"execution_time_ms,omitempty"` // Duration in milliseconds
	ErrorMessage  *string        `json:"error_message,omitempty"`
	Memo          map[string]any `json:"memo,omitempty"` // Workflow memo (metadata)
	// SLA tracking
	SLADeadline    *time.Time `json:"sla_deadline,omitempty"`
	SLARemaining   *int64     `json:"sla_remaining_seconds,omitempty"`
	SLABreached    bool       `json:"sla_breached"`
	// Batch-specific fields
	BatchID        *string `json:"batch_id,omitempty"`
	PolicyNumber   *string `json:"policy_number,omitempty"`
	AgentID        *string `json:"agent_id,omitempty"`
}

// WorkflowHistoryEvent represents a single event in workflow execution history
type WorkflowHistoryEvent struct {
	EventID      int64      `json:"event_id"`
	EventType    string     `json:"event_type"`    // WorkflowExecutionStarted, ActivityTaskScheduled, etc.
	EventTime    time.Time  `json:"event_time"`
	ActivityName *string    `json:"activity_name,omitempty"`
	Input        *string    `json:"input,omitempty"`
	Result       *string    `json:"result,omitempty"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	Details      map[string]any `json:"details,omitempty"`
}

// WorkflowHistory represents the full execution history of a workflow
type WorkflowHistory struct {
	WorkflowID string                  `json:"workflow_id"`
	Events     []WorkflowHistoryEvent  `json:"events"`
	TotalEvents int                    `json:"total_events"`
}

// WorkflowQueryResult represents the result of a workflow query
type WorkflowQueryResult struct {
	WorkflowID  string         `json:"workflow_id"`
	QueryType   string         `json:"query_type"`
	Result      map[string]any `json:"result"`
	ExecutedAt  time.Time      `json:"executed_at"`
}

// BatchProgress represents the progress of a commission batch workflow
type BatchProgress struct {
	BatchID           string    `json:"batch_id"`
	TotalPolicies     int       `json:"total_policies"`
	ProcessedPolicies int       `json:"processed_policies"`
	FailedPolicies    int       `json:"failed_policies"`
	ProgressPercent   float64   `json:"progress_percent"`
	Status            string    `json:"status"`
	StartTime         time.Time `json:"start_time"`
	EstimatedEndTime  *time.Time `json:"estimated_end_time,omitempty"`
	CurrentActivity   *string   `json:"current_activity,omitempty"`
	ErrorMessage      *string   `json:"error_message,omitempty"`
}

// IsRunning checks if the workflow is currently running
func (w *WorkflowInfo) IsRunning() bool {
	return w.Status == WorkflowStatusRunning
}

// IsCompleted checks if the workflow has completed successfully
func (w *WorkflowInfo) IsCompleted() bool {
	return w.Status == WorkflowStatusCompleted
}

// IsFailed checks if the workflow has failed
func (w *WorkflowInfo) IsFailed() bool {
	return w.Status == WorkflowStatusFailed
}

// CanRetry checks if the workflow can be retried
func (w *WorkflowInfo) CanRetry() bool {
	return w.Status == WorkflowStatusFailed || w.Status == WorkflowStatusTimedOut
}

// CanCancel checks if the workflow can be canceled
func (w *WorkflowInfo) CanCancel() bool {
	return w.Status == WorkflowStatusRunning
}

// CalculateSLARemaining calculates remaining SLA time in seconds
func (w *WorkflowInfo) CalculateSLARemaining() {
	if w.SLADeadline != nil {
		remaining := time.Until(*w.SLADeadline).Seconds()
		if remaining < 0 {
			w.SLABreached = true
			remaining = 0
		}
		remainingInt := int64(remaining)
		w.SLARemaining = &remainingInt
	}
}

// CalculateExecutionTime calculates workflow execution duration
func (w *WorkflowInfo) CalculateExecutionTime() {
	if w.CloseTime != nil {
		duration := w.CloseTime.Sub(w.StartTime).Milliseconds()
		w.ExecutionTime = &duration
	} else if w.IsRunning() {
		duration := time.Since(w.StartTime).Milliseconds()
		w.ExecutionTime = &duration
	}
}

// CalculateProgressPercent calculates batch progress percentage
func (b *BatchProgress) CalculateProgressPercent() {
	if b.TotalPolicies > 0 {
		b.ProgressPercent = (float64(b.ProcessedPolicies) / float64(b.TotalPolicies)) * 100
	}
}

// EstimateEndTime estimates when the batch will complete
func (b *BatchProgress) EstimateEndTime() {
	if b.ProcessedPolicies > 0 && b.TotalPolicies > b.ProcessedPolicies {
		elapsed := time.Since(b.StartTime)
		avgTimePerPolicy := elapsed / time.Duration(b.ProcessedPolicies)
		remaining := b.TotalPolicies - b.ProcessedPolicies
		estimatedRemaining := avgTimePerPolicy * time.Duration(remaining)
		estimatedEnd := time.Now().Add(estimatedRemaining)
		b.EstimatedEndTime = &estimatedEnd
	}
}
