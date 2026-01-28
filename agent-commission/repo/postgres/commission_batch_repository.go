package postgres

import (
	"context"
	"fmt"
	"time"

	"agent-commission/core/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/api-db"
)

// CommissionBatchRepository handles commission batch data operations
// Implements: FR-IC-COM-002 (Commission batch processing)
// Business Rule: BR-IC-COM-001 (Commission calculation), BR-IC-COM-012 (6-hour SLA)
type CommissionBatchRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewCommissionBatchRepository creates a new commission batch repository
func NewCommissionBatchRepository(db *dblib.DB, cfg *config.Config) *CommissionBatchRepository {
	return &CommissionBatchRepository{db: db, cfg: cfg}
}

const (
	commissionBatchTable = "commission_batches"
)

// CreateBatch creates a new commission batch
// Implements: FR-IC-COM-002, BR-IC-COM-012 (6-hour SLA)
func (r *CommissionBatchRepository) CreateBatch(
	ctx context.Context,
	batch *domain.CommissionBatch,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Calculate SLA deadline (6 hours from start)
	batch.SLADeadline = batch.StartedAt.Add(6 * time.Hour)

	q := dblib.Psql.Insert(commissionBatchTable).Columns(
		"batch_id",
		"month",
		"year",
		"status",
		"workflow_id",
		"workflow_state",
		"total_policies",
		"processed_records",
		"failed_records",
		"progress_percentage",
		"triggered_by",
		"started_at",
		"estimated_completion",
		"sla_deadline",
	).Values(
		batch.BatchID,
		batch.Month,
		batch.Year,
		batch.Status,
		batch.WorkflowID,
		batch.WorkflowState,
		batch.TotalPolicies,
		batch.ProcessedRecords,
		batch.FailedRecords,
		batch.ProgressPercentage,
		batch.TriggeredBy,
		batch.StartedAt,
		batch.EstimatedCompletion,
		batch.SLADeadline,
	)

	_, err := dblib.Insert(ctx, r.db, q)
	return err
}

// GetBatchByID retrieves a commission batch by batch ID
func (r *CommissionBatchRepository) GetBatchByID(
	ctx context.Context,
	batchID string,
) (domain.CommissionBatch, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"batch_id",
		"month",
		"year",
		"status",
		"workflow_id",
		"workflow_state",
		"total_policies",
		"processed_records",
		"failed_records",
		"progress_percentage",
		"triggered_by",
		"started_at",
		"estimated_completion",
		"sla_deadline",
		"completed_at",
		"failure_reason",
		"created_at",
		"updated_at",
	).From(commissionBatchTable).Where(sq.Eq{"batch_id": batchID})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.CommissionBatch])
}

// GetBatchByMonthYear retrieves a commission batch by month and year
// Used to check if batch already exists for the month
func (r *CommissionBatchRepository) GetBatchByMonthYear(
	ctx context.Context,
	month, year int,
) (domain.CommissionBatch, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"batch_id",
		"month",
		"year",
		"status",
		"workflow_id",
		"workflow_state",
		"total_policies",
		"processed_records",
		"failed_records",
		"progress_percentage",
		"triggered_by",
		"started_at",
		"estimated_completion",
		"sla_deadline",
		"completed_at",
		"failure_reason",
		"created_at",
		"updated_at",
	).From(commissionBatchTable).Where(sq.And{
		sq.Eq{"month": month},
		sq.Eq{"year": year},
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.CommissionBatch])
}

// UpdateBatchProgress updates the batch progress
// Uses batch optimization to update multiple fields in one query
func (r *CommissionBatchRepository) UpdateBatchProgress(
	ctx context.Context,
	batchID string,
	totalPolicies, processedRecords, failedRecords int,
	workflowState string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Calculate progress percentage
	progressPercentage := 0
	if totalPolicies > 0 {
		progressPercentage = (processedRecords * 100) / totalPolicies
	}

	q := dblib.Psql.Update(commissionBatchTable).
		Set("total_policies", totalPolicies).
		Set("processed_records", processedRecords).
		Set("failed_records", failedRecords).
		Set("progress_percentage", progressPercentage).
		Set("workflow_state", workflowState).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"batch_id": batchID})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// UpdateBatchStatus updates the batch status
func (r *CommissionBatchRepository) UpdateBatchStatus(
	ctx context.Context,
	batchID string,
	status domain.BatchStatus,
	workflowState string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	q := dblib.Psql.Update(commissionBatchTable).
		Set("status", status).
		Set("workflow_state", workflowState).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"batch_id": batchID})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// CompleteBatch marks the batch as completed
func (r *CommissionBatchRepository) CompleteBatch(
	ctx context.Context,
	batchID string,
	status domain.BatchStatus,
	failureReason *string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	now := time.Now()

	q := dblib.Psql.Update(commissionBatchTable).
		Set("status", status).
		Set("completed_at", now).
		Set("updated_at", now)

	if failureReason != nil {
		q = q.Set("failure_reason", *failureReason)
	}

	q = q.Where(sq.Eq{"batch_id": batchID})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// GenerateBatchID generates a unique batch ID for a given month/year
// Format: BATCH_YYYYMM_NNN
func GenerateBatchID(month, year int) string {
	return fmt.Sprintf("BATCH_%04d%02d_001", year, month)
}
