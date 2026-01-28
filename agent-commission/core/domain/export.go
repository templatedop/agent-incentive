package domain

import (
	"time"
)

// ExportStatus represents the status of an export job
type ExportStatus string

const (
	ExportStatusPending    ExportStatus = "PENDING"    // Export job queued
	ExportStatusProcessing ExportStatus = "PROCESSING" // Export being generated
	ExportStatusCompleted  ExportStatus = "COMPLETED"  // Export ready for download
	ExportStatusFailed     ExportStatus = "FAILED"     // Export failed
	ExportStatusExpired    ExportStatus = "EXPIRED"    // Download link expired
)

// ExportFormat represents the format of export file
type ExportFormat string

const (
	ExportFormatExcel ExportFormat = "EXCEL" // .xlsx
	ExportFormatPDF   ExportFormat = "PDF"   // .pdf
	ExportFormatCSV   ExportFormat = "CSV"   // .csv
)

// ExportType represents the type of data being exported
type ExportType string

const (
	ExportTypeTrialStatements  ExportType = "TRIAL_STATEMENTS"
	ExportTypeFinalStatements  ExportType = "FINAL_STATEMENTS"
	ExportTypeCommissionHistory ExportType = "COMMISSION_HISTORY"
	ExportTypeClawback         ExportType = "CLAWBACK_REPORT"
	ExportTypeSuspense         ExportType = "SUSPENSE_REPORT"
	ExportTypeCommissionSummary ExportType = "COMMISSION_SUMMARY"
)

// ExportJob represents an export job
// FR-IC-COM-010: Excel/PDF export
type ExportJob struct {
	ExportID     int64        `db:"export_id"`
	ExportType   ExportType   `db:"export_type"`
	ExportFormat ExportFormat `db:"export_format"`
	ExportStatus ExportStatus `db:"export_status"`
	RequestedBy  string       `db:"requested_by"`
	RequestedAt  time.Time    `db:"requested_at"`
	StartedAt    *time.Time   `db:"started_at"`
	CompletedAt  *time.Time   `db:"completed_at"`
	ExpiresAt    *time.Time   `db:"expires_at"` // Download link expiry (7 days)
	// File details
	FileName     *string `db:"file_name"`      // Generated file name
	FileSize     *int64  `db:"file_size"`      // File size in bytes
	FilePath     *string `db:"file_path"`      // S3/storage path
	DownloadURL  *string `db:"download_url"`   // Signed URL for download
	RecordCount  *int    `db:"record_count"`   // Number of records exported
	// Filter criteria (JSON stored as string)
	FilterCriteria *string `db:"filter_criteria"` // JSON string of filters applied
	// Error handling
	ErrorMessage *string `db:"error_message"`
	RetryCount   int     `db:"retry_count"`
	// Workflow tracking
	WorkflowID    *string `db:"workflow_id"`    // Temporal workflow ID for export job
	WorkflowState *string `db:"workflow_state"` // Workflow state (INITIATED, RUNNING, etc.)
	// Metadata
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// CommissionSummaryReport represents aggregated commission data
// FR-IC-COM-012: Commission summary reports
type CommissionSummaryReport struct {
	ReportID      int64     `db:"report_id"`
	ReportPeriod  string    `db:"report_period"` // MONTHLY, QUARTERLY, YEARLY
	PeriodStart   time.Time `db:"period_start"`
	PeriodEnd     time.Time `db:"period_end"`
	GeneratedAt   time.Time `db:"generated_at"`
	GeneratedBy   string    `db:"generated_by"`
	// Summary metrics
	TotalAgents          int     `db:"total_agents"`
	TotalPolicies        int     `db:"total_policies"`
	TotalCommissionGross float64 `db:"total_commission_gross"`
	TotalCommissionNet   float64 `db:"total_commission_net"`
	TotalTDS             float64 `db:"total_tds"`
	TotalClawback        float64 `db:"total_clawback"`
	TotalSuspense        float64 `db:"total_suspense"`
	// By type breakdown
	FirstYearCommission float64 `db:"first_year_commission"`
	RenewalCommission   float64 `db:"renewal_commission"`
	BonusCommission     float64 `db:"bonus_commission"`
	// By product breakdown (JSON)
	ByProduct *string `db:"by_product"` // JSON string
	// By agent breakdown (JSON)
	TopAgents *string `db:"top_agents"` // JSON string
	// Report file
	ReportFileURL *string `db:"report_file_url"`
}

// IsCompleted checks if export job is completed
func (e *ExportJob) IsCompleted() bool {
	return e.ExportStatus == ExportStatusCompleted
}

// IsFailed checks if export job has failed
func (e *ExportJob) IsFailed() bool {
	return e.ExportStatus == ExportStatusFailed
}

// IsProcessing checks if export job is processing
func (e *ExportJob) IsProcessing() bool {
	return e.ExportStatus == ExportStatusProcessing
}

// CanRetry checks if export job can be retried
func (e *ExportJob) CanRetry() bool {
	return e.ExportStatus == ExportStatusFailed && e.RetryCount < 3
}

// IsExpired checks if download link has expired
func (e *ExportJob) IsExpired() bool {
	if e.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*e.ExpiresAt)
}

// MarkAsProcessing marks the export job as processing
func (e *ExportJob) MarkAsProcessing() {
	e.ExportStatus = ExportStatusProcessing
	now := time.Now()
	e.StartedAt = &now
	e.UpdatedAt = now
}

// MarkAsCompleted marks the export job as completed
func (e *ExportJob) MarkAsCompleted(fileName string, fileSize int64, filePath string, downloadURL string, recordCount int) {
	e.ExportStatus = ExportStatusCompleted
	now := time.Now()
	e.CompletedAt = &now
	e.UpdatedAt = now

	e.FileName = &fileName
	e.FileSize = &fileSize
	e.FilePath = &filePath
	e.DownloadURL = &downloadURL
	e.RecordCount = &recordCount

	// Set expiry to 7 days from now
	expiresAt := now.AddDate(0, 0, 7)
	e.ExpiresAt = &expiresAt
}

// MarkAsFailed marks the export job as failed
func (e *ExportJob) MarkAsFailed(errorMessage string) {
	e.ExportStatus = ExportStatusFailed
	e.ErrorMessage = &errorMessage
	e.RetryCount++
	e.UpdatedAt = time.Now()
}

// ExportFilter represents common export filters
type ExportFilter struct {
	FromDate     *time.Time
	ToDate       *time.Time
	AgentID      *string
	PolicyNumber *string
	Status       *string
}

// ExportJobFilter represents search criteria for export jobs
type ExportJobFilter struct {
	ExportType   *ExportType
	ExportStatus *ExportStatus
	RequestedBy  *string
	FromDate     *time.Time
	ToDate       *time.Time
	Page         int
	Limit        int
}
