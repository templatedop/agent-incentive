package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"time"
)

// ExportJobResponse represents the response for export job creation
type ExportJobResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		ExportID     int64  `json:"export_id"`
		ExportType   string `json:"export_type"`
		ExportFormat string `json:"export_format"`
		ExportStatus string `json:"export_status"`
		RequestedAt  string `json:"requested_at"`
		Message      string `json:"message"`
	} `json:"data"`
}

// ExportJobStatusResponse represents the response for export job status query
type ExportJobStatusResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		ExportID       int64   `json:"export_id"`
		ExportType     string  `json:"export_type"`
		ExportFormat   string  `json:"export_format"`
		ExportStatus   string  `json:"export_status"`
		RequestedBy    string  `json:"requested_by"`
		RequestedAt    string  `json:"requested_at"`
		StartedAt      *string `json:"started_at,omitempty"`
		CompletedAt    *string `json:"completed_at,omitempty"`
		ExpiresAt      *string `json:"expires_at,omitempty"`
		FileName       *string `json:"file_name,omitempty"`
		FileSize       *int64  `json:"file_size_bytes,omitempty"`
		DownloadURL    *string `json:"download_url,omitempty"`
		RecordCount    *int    `json:"record_count,omitempty"`
		ErrorMessage   *string `json:"error_message,omitempty"`
		ProgressPercent float64 `json:"progress_percent"`
		IsExpired      bool    `json:"is_expired"`
	} `json:"data"`
}

// CommissionSummaryResponse represents the response for commission summary report
type CommissionSummaryResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		ReportID      int64  `json:"report_id"`
		ReportPeriod  string `json:"report_period"`
		PeriodStart   string `json:"period_start"`
		PeriodEnd     string `json:"period_end"`
		GeneratedAt   string `json:"generated_at"`
		Summary       CommissionSummaryData `json:"summary"`
		ReportFileURL *string `json:"report_file_url,omitempty"`
	} `json:"data"`
}

// CommissionSummaryData represents aggregated commission metrics
type CommissionSummaryData struct {
	TotalAgents          int     `json:"total_agents"`
	TotalPolicies        int     `json:"total_policies"`
	TotalCommissionGross float64 `json:"total_commission_gross"`
	TotalCommissionNet   float64 `json:"total_commission_net"`
	TotalTDS             float64 `json:"total_tds"`
	TotalClawback        float64 `json:"total_clawback"`
	TotalSuspense        float64 `json:"total_suspense"`
	// Breakdown by type
	ByType CommissionByType `json:"by_type"`
}

// CommissionByType represents commission breakdown by type
type CommissionByType struct {
	FirstYear float64 `json:"first_year_commission"`
	Renewal   float64 `json:"renewal_commission"`
	Bonus     float64 `json:"bonus_commission"`
}

// ExportJobListResponse represents list of export jobs
type ExportJobListResponse struct {
	port.StatusCodeAndMessage `json:",inline"`
	Data                      struct {
		ExportJobs []ExportJobSummary   `json:"export_jobs"`
		Pagination PaginationMetadata   `json:"pagination"`
	} `json:"data"`
}

// ExportJobSummary represents a summary of export job
type ExportJobSummary struct {
	ExportID     int64   `json:"export_id"`
	ExportType   string  `json:"export_type"`
	ExportFormat string  `json:"export_format"`
	ExportStatus string  `json:"export_status"`
	RequestedAt  string  `json:"requested_at"`
	CompletedAt  *string `json:"completed_at,omitempty"`
	RecordCount  *int    `json:"record_count,omitempty"`
	IsExpired    bool    `json:"is_expired"`
}

// NewExportJobResponse creates an export job response
func NewExportJobResponse(job domain.ExportJob) *ExportJobResponse {
	resp := &ExportJobResponse{
		StatusCodeAndMessage: port.CreateSuccess,
	}

	resp.Data.ExportID = job.ExportID
	resp.Data.ExportType = string(job.ExportType)
	resp.Data.ExportFormat = string(job.ExportFormat)
	resp.Data.ExportStatus = string(job.ExportStatus)
	resp.Data.RequestedAt = job.RequestedAt.Format(time.RFC3339)
	resp.Data.Message = "Export job created successfully. Check status using export_id."

	return resp
}

// NewExportJobStatusResponse creates an export job status response
func NewExportJobStatusResponse(job domain.ExportJob) *ExportJobStatusResponse {
	resp := &ExportJobStatusResponse{
		StatusCodeAndMessage: port.RetrieveSuccess,
	}

	resp.Data.ExportID = job.ExportID
	resp.Data.ExportType = string(job.ExportType)
	resp.Data.ExportFormat = string(job.ExportFormat)
	resp.Data.ExportStatus = string(job.ExportStatus)
	resp.Data.RequestedBy = job.RequestedBy
	resp.Data.RequestedAt = job.RequestedAt.Format(time.RFC3339)
	resp.Data.FileName = job.FileName
	resp.Data.FileSize = job.FileSize
	resp.Data.DownloadURL = job.DownloadURL
	resp.Data.RecordCount = job.RecordCount
	resp.Data.ErrorMessage = job.ErrorMessage
	resp.Data.IsExpired = job.IsExpired()

	if job.StartedAt != nil {
		startedAtStr := job.StartedAt.Format(time.RFC3339)
		resp.Data.StartedAt = &startedAtStr
	}

	if job.CompletedAt != nil {
		completedAtStr := job.CompletedAt.Format(time.RFC3339)
		resp.Data.CompletedAt = &completedAtStr
	}

	if job.ExpiresAt != nil {
		expiresAtStr := job.ExpiresAt.Format(time.RFC3339)
		resp.Data.ExpiresAt = &expiresAtStr
	}

	// Calculate progress
	switch job.ExportStatus {
	case domain.ExportStatusCompleted:
		resp.Data.ProgressPercent = 100.0
	case domain.ExportStatusProcessing:
		// Estimate based on time elapsed (mock)
		if job.StartedAt != nil {
			elapsed := time.Since(*job.StartedAt).Minutes()
			progress := elapsed / 5.0 * 100 // Assume 5 minutes for completion
			if progress > 95 {
				progress = 95
			}
			resp.Data.ProgressPercent = progress
		} else {
			resp.Data.ProgressPercent = 10.0
		}
	case domain.ExportStatusPending:
		resp.Data.ProgressPercent = 0.0
	case domain.ExportStatusFailed:
		resp.Data.ProgressPercent = 0.0
	}

	return resp
}

// NewCommissionSummaryResponse creates a commission summary response
func NewCommissionSummaryResponse(report domain.CommissionSummaryReport) *CommissionSummaryResponse {
	resp := &CommissionSummaryResponse{
		StatusCodeAndMessage: port.RetrieveSuccess,
	}

	resp.Data.ReportID = report.ReportID
	resp.Data.ReportPeriod = report.ReportPeriod
	resp.Data.PeriodStart = report.PeriodStart.Format(time.RFC3339)
	resp.Data.PeriodEnd = report.PeriodEnd.Format(time.RFC3339)
	resp.Data.GeneratedAt = report.GeneratedAt.Format(time.RFC3339)
	resp.Data.ReportFileURL = report.ReportFileURL

	resp.Data.Summary = CommissionSummaryData{
		TotalAgents:          report.TotalAgents,
		TotalPolicies:        report.TotalPolicies,
		TotalCommissionGross: report.TotalCommissionGross,
		TotalCommissionNet:   report.TotalCommissionNet,
		TotalTDS:             report.TotalTDS,
		TotalClawback:        report.TotalClawback,
		TotalSuspense:        report.TotalSuspense,
		ByType: CommissionByType{
			FirstYear: report.FirstYearCommission,
			Renewal:   report.RenewalCommission,
			Bonus:     report.BonusCommission,
		},
	}

	return resp
}

// NewExportJobSummary creates an export job summary
func NewExportJobSummary(job domain.ExportJob) ExportJobSummary {
	summary := ExportJobSummary{
		ExportID:     job.ExportID,
		ExportType:   string(job.ExportType),
		ExportFormat: string(job.ExportFormat),
		ExportStatus: string(job.ExportStatus),
		RequestedAt:  job.RequestedAt.Format(time.RFC3339),
		RecordCount:  job.RecordCount,
		IsExpired:    job.IsExpired(),
	}

	if job.CompletedAt != nil {
		completedAtStr := job.CompletedAt.Format(time.RFC3339)
		summary.CompletedAt = &completedAtStr
	}

	return summary
}

// NewExportJobSummaries creates a list of export job summaries
func NewExportJobSummaries(jobs []domain.ExportJob) []ExportJobSummary {
	summaries := make([]ExportJobSummary, len(jobs))
	for i, job := range jobs {
		summaries[i] = NewExportJobSummary(job)
	}
	return summaries
}
