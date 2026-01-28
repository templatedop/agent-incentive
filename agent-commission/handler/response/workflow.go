package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"time"
)

// WorkflowStatusResponse represents the response for workflow status query
type WorkflowStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		WorkflowID     string  `json:"workflow_id"`
		WorkflowType   string  `json:"workflow_type"`
		Status         string  `json:"status"`
		StartTime      string  `json:"start_time"`
		CloseTime      *string `json:"close_time,omitempty"`
		ExecutionTime  *int64  `json:"execution_time_ms,omitempty"`
		ErrorMessage   *string `json:"error_message,omitempty"`
		SLADeadline    *string `json:"sla_deadline,omitempty"`
		SLARemaining   *int64  `json:"sla_remaining_seconds,omitempty"`
		SLABreached    bool    `json:"sla_breached"`
		BatchID        *string `json:"batch_id,omitempty"`
		PolicyNumber   *string `json:"policy_number,omitempty"`
		AgentID        *string `json:"agent_id,omitempty"`
		Memo           map[string]any `json:"memo,omitempty"`
	} `json:"data"`
}

// WorkflowHistoryResponse represents the response for workflow history query
type WorkflowHistoryResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		WorkflowID  string                    `json:"workflow_id"`
		Events      []WorkflowHistoryEventDTO `json:"events"`
		TotalEvents int                       `json:"total_events"`
	} `json:"data"`
}

// WorkflowHistoryEventDTO represents a workflow history event
type WorkflowHistoryEventDTO struct {
	EventID      int64          `json:"event_id"`
	EventType    string         `json:"event_type"`
	EventTime    string         `json:"event_time"`
	ActivityName *string        `json:"activity_name,omitempty"`
	Input        *string        `json:"input,omitempty"`
	Result       *string        `json:"result,omitempty"`
	ErrorMessage *string        `json:"error_message,omitempty"`
	Details      map[string]any `json:"details,omitempty"`
}

// WorkflowCancelResponse represents the response for cancel operation
type WorkflowCancelResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		WorkflowID string `json:"workflow_id"`
		Message    string `json:"message"`
		CanceledAt string `json:"canceled_at"`
	} `json:"data"`
}

// WorkflowRetryResponse represents the response for retry operation
type WorkflowRetryResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		OriginalWorkflowID string `json:"original_workflow_id"`
		NewWorkflowID      string `json:"new_workflow_id"`
		Message            string `json:"message"`
		RetriedAt          string `json:"retried_at"`
	} `json:"data"`
}

// WorkflowQueryResponse represents the response for workflow query operation
type WorkflowQueryResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		WorkflowID string         `json:"workflow_id"`
		QueryType  string         `json:"query_type"`
		Result     map[string]any `json:"result"`
		ExecutedAt string         `json:"executed_at"`
	} `json:"data"`
}

// BatchStatusResponse represents the response for batch status query
type BatchStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		BatchID     string  `json:"batch_id"`
		Status      string  `json:"status"`
		WorkflowID  *string `json:"workflow_id,omitempty"`
		StartTime   string  `json:"start_time"`
		EndTime     *string `json:"end_time,omitempty"`
		SLADeadline *string `json:"sla_deadline,omitempty"`
		SLABreached bool    `json:"sla_breached"`
	} `json:"data"`
}

// BatchProgressResponse represents the response for batch progress query
type BatchProgressResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		BatchID           string  `json:"batch_id"`
		TotalPolicies     int     `json:"total_policies"`
		ProcessedPolicies int     `json:"processed_policies"`
		FailedPolicies    int     `json:"failed_policies"`
		ProgressPercent   float64 `json:"progress_percent"`
		Status            string  `json:"status"`
		StartTime         string  `json:"start_time"`
		EstimatedEndTime  *string `json:"estimated_end_time,omitempty"`
		CurrentActivity   *string `json:"current_activity,omitempty"`
		ErrorMessage      *string `json:"error_message,omitempty"`
	} `json:"data"`
}

// BatchCancelResponse represents the response for batch cancel operation
type BatchCancelResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		BatchID    string `json:"batch_id"`
		WorkflowID string `json:"workflow_id"`
		Message    string `json:"message"`
		CanceledAt string `json:"canceled_at"`
	} `json:"data"`
}

// NewWorkflowStatusResponse creates a workflow status response
func NewWorkflowStatusResponse(workflow domain.WorkflowInfo) *WorkflowStatusResponse {
	resp := &WorkflowStatusResponse{
		StatusCodeAndMessage: port.RetrieveSuccess,
	}

	resp.Data.WorkflowID = workflow.WorkflowID
	resp.Data.WorkflowType = string(workflow.WorkflowType)
	resp.Data.Status = string(workflow.Status)
	resp.Data.StartTime = workflow.StartTime.Format(time.RFC3339)
	resp.Data.ExecutionTime = workflow.ExecutionTime
	resp.Data.ErrorMessage = workflow.ErrorMessage
	resp.Data.SLARemaining = workflow.SLARemaining
	resp.Data.SLABreached = workflow.SLABreached
	resp.Data.BatchID = workflow.BatchID
	resp.Data.PolicyNumber = workflow.PolicyNumber
	resp.Data.AgentID = workflow.AgentID
	resp.Data.Memo = workflow.Memo

	if workflow.CloseTime != nil {
		closeTimeStr := workflow.CloseTime.Format(time.RFC3339)
		resp.Data.CloseTime = &closeTimeStr
	}

	if workflow.SLADeadline != nil {
		deadlineStr := workflow.SLADeadline.Format(time.RFC3339)
		resp.Data.SLADeadline = &deadlineStr
	}

	return resp
}

// NewWorkflowHistoryResponse creates a workflow history response
func NewWorkflowHistoryResponse(history domain.WorkflowHistory) *WorkflowHistoryResponse {
	resp := &WorkflowHistoryResponse{
		StatusCodeAndMessage: port.ListSuccess,
	}

	resp.Data.WorkflowID = history.WorkflowID
	resp.Data.TotalEvents = history.TotalEvents
	resp.Data.Events = make([]WorkflowHistoryEventDTO, len(history.Events))

	for i, event := range history.Events {
		resp.Data.Events[i] = WorkflowHistoryEventDTO{
			EventID:      event.EventID,
			EventType:    event.EventType,
			EventTime:    event.EventTime.Format(time.RFC3339),
			ActivityName: event.ActivityName,
			Input:        event.Input,
			Result:       event.Result,
			ErrorMessage: event.ErrorMessage,
			Details:      event.Details,
		}
	}

	return resp
}

// NewBatchProgressResponse creates a batch progress response
func NewBatchProgressResponse(progress domain.BatchProgress) *BatchProgressResponse {
	resp := &BatchProgressResponse{
		StatusCodeAndMessage: port.RetrieveSuccess,
	}

	resp.Data.BatchID = progress.BatchID
	resp.Data.TotalPolicies = progress.TotalPolicies
	resp.Data.ProcessedPolicies = progress.ProcessedPolicies
	resp.Data.FailedPolicies = progress.FailedPolicies
	resp.Data.ProgressPercent = progress.ProgressPercent
	resp.Data.Status = progress.Status
	resp.Data.StartTime = progress.StartTime.Format(time.RFC3339)
	resp.Data.CurrentActivity = progress.CurrentActivity
	resp.Data.ErrorMessage = progress.ErrorMessage

	if progress.EstimatedEndTime != nil {
		endTimeStr := progress.EstimatedEndTime.Format(time.RFC3339)
		resp.Data.EstimatedEndTime = &endTimeStr
	}

	return resp
}
