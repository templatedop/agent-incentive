package handler

import (
	"errors"
	"fmt"
	"time"

	"agent-commission/core/domain"
	"agent-commission/core/port"
	"agent-commission/handler/response"
	repo "agent-commission/repo/postgres"

	log "gitlab.cept.gov.in/it-2.0-common/api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// WorkflowManagementHandler handles workflow monitoring and control operations
// Temporal workflow inspection and SLA tracking
type WorkflowManagementHandler struct {
	*serverHandler.Base
	commissionBatchRepo *repo.CommissionBatchRepository
	// TODO: Add Temporal client when available
	// temporalClient      temporal.Client
}

// NewWorkflowManagementHandler creates a new workflow management handler
func NewWorkflowManagementHandler(
	commissionBatchRepo *repo.CommissionBatchRepository,
) *WorkflowManagementHandler {
	base := serverHandler.New("WorkflowManagement").SetPrefix("/v1").AddPrefix("")
	return &WorkflowManagementHandler{
		Base:                base,
		commissionBatchRepo: commissionBatchRepo,
	}
}

// Routes returns all routes for workflow management endpoints
func (h *WorkflowManagementHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.GET("/workflows/:workflowId/status", h.GetWorkflowStatus).Name("Get Workflow Status"),
		serverRoute.GET("/workflows/:workflowId/history", h.GetWorkflowHistory).Name("Get Workflow History"),
		serverRoute.POST("/workflows/:workflowId/cancel", h.CancelWorkflow).Name("Cancel Workflow"),
		serverRoute.POST("/workflows/:workflowId/retry", h.RetryWorkflow).Name("Retry Workflow"),
		serverRoute.POST("/workflows/:workflowId/query", h.QueryWorkflow).Name("Query Workflow State"),
		serverRoute.GET("/batches/:batchId/status", h.GetBatchStatus).Name("Get Batch Status"),
		serverRoute.GET("/batches/:batchId/progress", h.GetBatchProgress).Name("Get Batch Progress"),
		serverRoute.POST("/batches/:batchId/cancel", h.CancelBatch).Name("Cancel Batch"),
	}
}

// WorkflowIDUri represents the URI parameter for workflow ID
type WorkflowIDUri struct {
	WorkflowID string `uri:"workflowId" validate:"required"`
}

// Note: BatchIDUri is defined in commission_batch_handler.go

// QueryWorkflowRequest represents request to query workflow state
type QueryWorkflowRequest struct {
	WorkflowID string `uri:"workflowId" validate:"required"`
	QueryType  string `json:"query_type" validate:"required"`
}

// GetWorkflowStatus retrieves the current status of a workflow
// API: GET /workflows/{workflowId}/status
// Temporal Integration: temporal.Client.DescribeWorkflowExecution
func (h *WorkflowManagementHandler) GetWorkflowStatus(
	sctx *serverRoute.Context,
	uri WorkflowIDUri,
) (*response.WorkflowStatusResponse, error) {
	log.Info(sctx.Ctx, "Getting workflow status for: %s", uri.WorkflowID)

	// TODO: Replace with actual Temporal client call
	// workflowExec := h.temporalClient.GetWorkflow(sctx.Ctx, uri.WorkflowID, "")
	// description, err := workflowExec.Describe(sctx.Ctx)
	// if err != nil {
	//     log.Error(sctx.Ctx, "Error describing workflow: %v", err)
	//     return nil, err
	// }

	// Mock implementation for now
	workflow := domain.WorkflowInfo{
		WorkflowID:   uri.WorkflowID,
		WorkflowType: domain.WorkflowTypeCommissionBatch,
		Status:       domain.WorkflowStatusRunning,
		StartTime:    time.Now().Add(-1 * time.Hour),
		SLABreached:  false,
	}

	// Calculate SLA remaining
	slaDeadline := time.Now().Add(5 * time.Hour)
	workflow.SLADeadline = &slaDeadline
	workflow.CalculateSLARemaining()
	workflow.CalculateExecutionTime()

	log.Info(sctx.Ctx, "Workflow status: %s", workflow.Status)

	resp := response.NewWorkflowStatusResponse(workflow)
	return resp, nil
}

// GetWorkflowHistory retrieves the execution history of a workflow
// API: GET /workflows/{workflowId}/history
// Temporal Integration: temporal.Client.GetWorkflowHistory
func (h *WorkflowManagementHandler) GetWorkflowHistory(
	sctx *serverRoute.Context,
	uri WorkflowIDUri,
) (*response.WorkflowHistoryResponse, error) {
	log.Info(sctx.Ctx, "Getting workflow history for: %s", uri.WorkflowID)

	// TODO: Replace with actual Temporal client call
	// historyIter := h.temporalClient.GetWorkflowHistory(sctx.Ctx, uri.WorkflowID, "", false, enumspb.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	// var events []domain.WorkflowHistoryEvent
	// for historyIter.HasNext() {
	//     event, err := historyIter.Next()
	//     if err != nil {
	//         return nil, err
	//     }
	//     events = append(events, mapToHistoryEvent(event))
	// }

	// Mock implementation for now
	history := domain.WorkflowHistory{
		WorkflowID: uri.WorkflowID,
		Events: []domain.WorkflowHistoryEvent{
			{
				EventID:   1,
				EventType: "WorkflowExecutionStarted",
				EventTime: time.Now().Add(-1 * time.Hour),
			},
			{
				EventID:      2,
				EventType:    "ActivityTaskScheduled",
				EventTime:    time.Now().Add(-55 * time.Minute),
				ActivityName: ptrString("CalculateCommissions"),
			},
			{
				EventID:   3,
				EventType: "ActivityTaskCompleted",
				EventTime: time.Now().Add(-50 * time.Minute),
				Result:    ptrString("Success"),
			},
		},
		TotalEvents: 3,
	}

	log.Info(sctx.Ctx, "Found %d history events", history.TotalEvents)

	resp := response.NewWorkflowHistoryResponse(history)
	return resp, nil
}

// CancelWorkflow cancels a running workflow
// API: POST /workflows/{workflowId}/cancel
// Temporal Integration: temporal.Client.CancelWorkflow
func (h *WorkflowManagementHandler) CancelWorkflow(
	sctx *serverRoute.Context,
	uri WorkflowIDUri,
) (*response.WorkflowCancelResponse, error) {
	log.Info(sctx.Ctx, "Canceling workflow: %s", uri.WorkflowID)

	// TODO: Replace with actual Temporal client call
	// err := h.temporalClient.CancelWorkflow(sctx.Ctx, uri.WorkflowID, "")
	// if err != nil {
	//     log.Error(sctx.Ctx, "Error canceling workflow: %v", err)
	//     return nil, err
	// }

	// Mock implementation for now
	resp := &response.WorkflowCancelResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
	}
	resp.Data.WorkflowID = uri.WorkflowID
	resp.Data.Message = "Workflow cancellation initiated successfully"
	resp.Data.CanceledAt = time.Now().Format(time.RFC3339)

	log.Info(sctx.Ctx, "Workflow cancellation initiated: %s", uri.WorkflowID)

	return resp, nil
}

// RetryWorkflow retries a failed workflow
// API: POST /workflows/{workflowId}/retry
// Creates a new workflow execution with the same input
func (h *WorkflowManagementHandler) RetryWorkflow(
	sctx *serverRoute.Context,
	uri WorkflowIDUri,
) (*response.WorkflowRetryResponse, error) {
	log.Info(sctx.Ctx, "Retrying workflow: %s", uri.WorkflowID)

	// TODO: Replace with actual Temporal client call
	// First, get the original workflow's input
	// workflowExec := h.temporalClient.GetWorkflow(sctx.Ctx, uri.WorkflowID, "")
	// description, err := workflowExec.Describe(sctx.Ctx)
	// if err != nil {
	//     return nil, err
	// }
	//
	// Check if workflow is in a failed state
	// if description.WorkflowExecutionInfo.Status != enumspb.WORKFLOW_EXECUTION_STATUS_FAILED {
	//     return nil, errors.New("workflow is not in a failed state and cannot be retried")
	// }
	//
	// Start a new workflow with the same input
	// newWorkflowID := uri.WorkflowID + "-retry-" + time.Now().Format("20060102150405")
	// workflowOptions := client.StartWorkflowOptions{
	//     ID:        newWorkflowID,
	//     TaskQueue: description.WorkflowExecutionInfo.TaskQueue,
	// }
	// run, err := h.temporalClient.ExecuteWorkflow(sctx.Ctx, workflowOptions, workflowType, input)
	// if err != nil {
	//     return nil, err
	// }

	// Mock implementation for now
	newWorkflowID := uri.WorkflowID + "-retry-" + time.Now().Format("20060102150405")

	resp := &response.WorkflowRetryResponse{
		StatusCodeAndMessage: port.CreateSuccess,
	}
	resp.Data.OriginalWorkflowID = uri.WorkflowID
	resp.Data.NewWorkflowID = newWorkflowID
	resp.Data.Message = "Workflow retry initiated successfully"
	resp.Data.RetriedAt = time.Now().Format(time.RFC3339)

	log.Info(sctx.Ctx, "Workflow retry initiated: %s -> %s", uri.WorkflowID, newWorkflowID)

	return resp, nil
}

// QueryWorkflow queries the state of a running workflow
// API: POST /workflows/{workflowId}/query
// Temporal Integration: temporal.Client.QueryWorkflow
func (h *WorkflowManagementHandler) QueryWorkflow(
	sctx *serverRoute.Context,
	req QueryWorkflowRequest,
) (*response.WorkflowQueryResponse, error) {
	log.Info(sctx.Ctx, "Querying workflow %s with query type: %s", req.WorkflowID, req.QueryType)

	// TODO: Replace with actual Temporal client call
	// workflowExec := h.temporalClient.GetWorkflow(sctx.Ctx, req.WorkflowID, "")
	// var result map[string]any
	// err := workflowExec.Query(sctx.Ctx, req.QueryType, &result)
	// if err != nil {
	//     log.Error(sctx.Ctx, "Error querying workflow: %v", err)
	//     return nil, err
	// }

	// Mock implementation for now
	mockResult := map[string]any{
		"query_type":         req.QueryType,
		"current_state":      "PROCESSING",
		"processed_records":  250,
		"total_records":      1000,
		"progress_percent":   25.0,
		"last_activity":      "CalculateCommissions",
		"last_activity_time": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
	}

	resp := &response.WorkflowQueryResponse{
		StatusCodeAndMessage: port.RetrieveSuccess,
	}
	resp.Data.WorkflowID = req.WorkflowID
	resp.Data.QueryType = req.QueryType
	resp.Data.Result = mockResult
	resp.Data.ExecutedAt = time.Now().Format(time.RFC3339)

	log.Info(sctx.Ctx, "Workflow query executed successfully")

	return resp, nil
}

// GetBatchStatus retrieves the status of a commission batch
// API: GET /batches/{batchId}/status
// Queries batch status from database
func (h *WorkflowManagementHandler) GetBatchStatus(
	sctx *serverRoute.Context,
	uri BatchIDUri,
) (*response.BatchStatusResponse, error) {
	log.Info(sctx.Ctx, "Getting batch status for: %s", uri.BatchID)

	// Get batch from database
	batch, err := h.commissionBatchRepo.GetBatchByID(sctx.Ctx, uri.BatchID)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching batch: %v", err)
		return nil, errors.New("batch not found")
	}

	resp := &response.BatchStatusResponse{
		StatusCodeAndMessage: port.RetrieveSuccess,
	}
	resp.Data.BatchID = batch.BatchID
	resp.Data.Status = string(batch.Status)
	resp.Data.WorkflowID = batch.WorkflowID
	resp.Data.StartTime = batch.StartTime.Format(time.RFC3339)
	resp.Data.SLABreached = batch.SLABreached

	if batch.EndTime != nil {
		endTimeStr := batch.EndTime.Format(time.RFC3339)
		resp.Data.EndTime = &endTimeStr
	}

	if batch.SLADeadline != nil {
		deadlineStr := batch.SLADeadline.Format(time.RFC3339)
		resp.Data.SLADeadline = &deadlineStr
	}

	log.Info(sctx.Ctx, "Batch status: %s", batch.Status)

	return resp, nil
}

// GetBatchProgress retrieves the detailed progress of a commission batch
// API: GET /batches/{batchId}/progress
// Combines database batch info with workflow query for real-time progress
func (h *WorkflowManagementHandler) GetBatchProgress(
	sctx *serverRoute.Context,
	uri BatchIDUri,
) (*response.BatchProgressResponse, error) {
	log.Info(sctx.Ctx, "Getting batch progress for: %s", uri.BatchID)

	// Get batch from database
	batch, err := h.commissionBatchRepo.GetBatchByID(sctx.Ctx, uri.BatchID)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching batch: %v", err)
		return nil, errors.New("batch not found")
	}

	// TODO: Query workflow for real-time progress if batch is running
	// if batch.WorkflowID != nil && batch.Status == domain.BatchStatusProcessing {
	//     workflowExec := h.temporalClient.GetWorkflow(sctx.Ctx, *batch.WorkflowID, "")
	//     var progressData map[string]any
	//     err := workflowExec.Query(sctx.Ctx, "GetProgress", &progressData)
	//     if err == nil {
	//         batch.ProcessedCount = int(progressData["processed"].(float64))
	//     }
	// }

	// Calculate progress
	progress := domain.BatchProgress{
		BatchID:           batch.BatchID,
		TotalPolicies:     batch.TotalCount,
		ProcessedPolicies: batch.ProcessedCount,
		FailedPolicies:    batch.FailedCount,
		Status:            string(batch.Status),
		StartTime:         batch.StartTime,
	}

	progress.CalculateProgressPercent()

	// Estimate end time if still processing
	if batch.Status == domain.BatchStatusProcessing {
		progress.EstimateEndTime()
		currentActivity := "Processing commissions"
		progress.CurrentActivity = &currentActivity
	}

	if batch.FailureReason != nil {
		progress.ErrorMessage = batch.FailureReason
	}

	log.Info(sctx.Ctx, "Batch progress: %.2f%% (%d/%d)", progress.ProgressPercent, progress.ProcessedPolicies, progress.TotalPolicies)

	resp := response.NewBatchProgressResponse(progress)
	return resp, nil
}

// CancelBatch cancels a running commission batch
// API: POST /batches/{batchId}/cancel
// Cancels the associated Temporal workflow
func (h *WorkflowManagementHandler) CancelBatch(
	sctx *serverRoute.Context,
	uri BatchIDUri,
) (*response.BatchCancelResponse, error) {
	log.Info(sctx.Ctx, "Canceling batch: %s", uri.BatchID)

	// Get batch from database
	batch, err := h.commissionBatchRepo.GetBatchByID(sctx.Ctx, uri.BatchID)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching batch: %v", err)
		return nil, errors.New("batch not found")
	}

	// Check if batch can be canceled
	if batch.Status != domain.BatchStatusProcessing {
		log.Warn(sctx.Ctx, "Batch is not in processing state: %s", batch.Status)
		return nil, fmt.Errorf("batch cannot be canceled in %s state", batch.Status)
	}

	if batch.WorkflowID == nil {
		log.Warn(sctx.Ctx, "Batch has no associated workflow")
		return nil, errors.New("batch has no associated workflow")
	}

	// TODO: Cancel the Temporal workflow
	// err = h.temporalClient.CancelWorkflow(sctx.Ctx, *batch.WorkflowID, "")
	// if err != nil {
	//     log.Error(sctx.Ctx, "Error canceling workflow: %v", err)
	//     return nil, err
	// }

	// Update batch status
	err = h.commissionBatchRepo.UpdateBatchStatus(sctx.Ctx, batch.BatchID, domain.BatchStatusCancelled, nil)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating batch status: %v", err)
		return nil, err
	}

	resp := &response.BatchCancelResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
	}
	resp.Data.BatchID = batch.BatchID
	resp.Data.WorkflowID = *batch.WorkflowID
	resp.Data.Message = "Batch cancellation initiated successfully"
	resp.Data.CanceledAt = time.Now().Format(time.RFC3339)

	log.Info(sctx.Ctx, "Batch cancellation initiated: %s", batch.BatchID)

	return resp, nil
}
