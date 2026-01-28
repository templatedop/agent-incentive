package postgres

import (
	"context"
	"time"

	"agent-commission/core/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/api-db"
)

// ClawbackRepository handles clawback data operations
// FR-IC-COM-012: Clawback management
type ClawbackRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewClawbackRepository creates a new clawback repository
func NewClawbackRepository(db *dblib.DB, cfg *config.Config) *ClawbackRepository {
	return &ClawbackRepository{db: db, cfg: cfg}
}

const (
	clawbackTable         = "clawbacks"
	clawbackRecoveryTable = "clawback_recoveries"
)

// CreateClawback creates a new clawback record
// Implements: FR-IC-COM-012, BR-IC-CLAWBACK-001
// Calculates clawback amount based on policy age and graduated percentage
func (r *ClawbackRepository) CreateClawback(
	ctx context.Context,
	clawback *domain.Clawback,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Calculate policy age in months
	policyAgeMonths := int(time.Since(clawback.PolicyInceptionDate).Hours() / 24 / 30)
	clawback.PolicyAge = policyAgeMonths

	// Calculate clawback percentage based on policy age
	clawback.ClawbackPercentage = domain.CalculateClawbackPercentage(policyAgeMonths)

	// Calculate clawback amount
	clawback.ClawbackAmount = clawback.OriginalCommission * (clawback.ClawbackPercentage / 100)
	clawback.PendingAmount = clawback.ClawbackAmount
	clawback.RecoveredAmount = 0

	// Set initial status
	if clawback.ClawbackStatus == "" {
		clawback.ClawbackStatus = domain.ClawbackStatusPending
	}

	// Set timestamps
	now := time.Now()
	clawback.CreatedAt = now
	clawback.UpdatedAt = now
	clawback.TriggerDate = now
	clawback.Version = 1

	// Build insert query
	q := dblib.Psql.Insert(clawbackTable).Columns(
		"policy_number",
		"agent_id",
		"original_commission",
		"clawback_amount",
		"recovered_amount",
		"pending_amount",
		"clawback_percentage",
		"clawback_reason",
		"clawback_status",
		"policy_age",
		"trigger_date",
		"policy_inception_date",
		"policy_end_date",
		"recovery_schedule",
		"installment_months",
		"notes",
		"workflow_id",
		"workflow_state",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"version",
	).Values(
		clawback.PolicyNumber,
		clawback.AgentID,
		clawback.OriginalCommission,
		clawback.ClawbackAmount,
		clawback.RecoveredAmount,
		clawback.PendingAmount,
		clawback.ClawbackPercentage,
		string(clawback.ClawbackReason),
		string(clawback.ClawbackStatus),
		clawback.PolicyAge,
		clawback.TriggerDate,
		clawback.PolicyInceptionDate,
		clawback.PolicyEndDate,
		clawback.RecoverySchedule,
		clawback.InstallmentMonths,
		clawback.Notes,
		clawback.WorkflowID,
		clawback.WorkflowState,
		clawback.CreatedAt,
		clawback.UpdatedAt,
		clawback.CreatedBy,
		clawback.UpdatedBy,
		clawback.Version,
	).Suffix("RETURNING clawback_id")

	// Execute with InsertReturning
	scanFn := func(row pgx.CollectableRow) (int64, error) {
		var id int64
		err := row.Scan(&id)
		return id, err
	}

	id, err := dblib.InsertReturning(ctx, r.db, q, scanFn)
	if err != nil {
		return err
	}

	clawback.ClawbackID = id
	return nil
}

// GetClawbackByID retrieves a clawback by its ID
func (r *ClawbackRepository) GetClawbackByID(
	ctx context.Context,
	clawbackID int64,
) (domain.Clawback, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutShort"))
	defer cancel()

	q := dblib.Psql.Select(
		"clawback_id",
		"policy_number",
		"agent_id",
		"original_commission",
		"clawback_amount",
		"recovered_amount",
		"pending_amount",
		"clawback_percentage",
		"clawback_reason",
		"clawback_status",
		"policy_age",
		"trigger_date",
		"policy_inception_date",
		"policy_end_date",
		"recovery_start_date",
		"recovery_end_date",
		"recovery_schedule",
		"installment_months",
		"approved_by",
		"approved_at",
		"waived_by",
		"waived_at",
		"waiver_reason",
		"workflow_id",
		"workflow_state",
		"notes",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).From(clawbackTable).Where(sq.Eq{"clawback_id": clawbackID, "deleted_at": nil})

	scanFn := pgx.RowToStructByName[domain.Clawback]
	return dblib.SelectOne(ctx, r.db, q, scanFn)
}

// GetClawbacksByPolicyNumber retrieves all clawbacks for a policy
func (r *ClawbackRepository) GetClawbacksByPolicyNumber(
	ctx context.Context,
	policyNumber string,
) ([]domain.Clawback, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	q := dblib.Psql.Select(
		"clawback_id",
		"policy_number",
		"agent_id",
		"original_commission",
		"clawback_amount",
		"recovered_amount",
		"pending_amount",
		"clawback_percentage",
		"clawback_reason",
		"clawback_status",
		"policy_age",
		"trigger_date",
		"policy_inception_date",
		"policy_end_date",
		"recovery_start_date",
		"recovery_end_date",
		"recovery_schedule",
		"installment_months",
		"approved_by",
		"approved_at",
		"waived_by",
		"waived_at",
		"waiver_reason",
		"workflow_id",
		"workflow_state",
		"notes",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).From(clawbackTable).Where(sq.Eq{"policy_number": policyNumber, "deleted_at": nil}).
		OrderBy("trigger_date DESC")

	scanFn := pgx.RowToStructByName[domain.Clawback]
	return dblib.SelectRows(ctx, r.db, q, scanFn)
}

// SearchClawbacks searches clawbacks with filters
// Uses batch optimization to fetch count + results in single round trip
func (r *ClawbackRepository) SearchClawbacks(
	ctx context.Context,
	filter *domain.ClawbackFilter,
) ([]domain.Clawback, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Build base query with filters
	baseQuery := dblib.Psql.Select().From(clawbackTable).Where(sq.Eq{"deleted_at": nil})

	// Apply filters
	if filter.AgentID != nil && *filter.AgentID != "" {
		baseQuery = baseQuery.Where(sq.Eq{"agent_id": *filter.AgentID})
	}
	if filter.PolicyNumber != nil && *filter.PolicyNumber != "" {
		baseQuery = baseQuery.Where(sq.Eq{"policy_number": *filter.PolicyNumber})
	}
	if filter.ClawbackStatus != nil {
		baseQuery = baseQuery.Where(sq.Eq{"clawback_status": string(*filter.ClawbackStatus)})
	}
	if filter.ClawbackReason != nil {
		baseQuery = baseQuery.Where(sq.Eq{"clawback_reason": string(*filter.ClawbackReason)})
	}
	if filter.FromDate != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"trigger_date": *filter.FromDate})
	}
	if filter.ToDate != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"trigger_date": *filter.ToDate})
	}
	if filter.MinAmount != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"clawback_amount": *filter.MinAmount})
	}
	if filter.MaxAmount != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"clawback_amount": *filter.MaxAmount})
	}

	// Count query
	countQuery := baseQuery.Columns("COUNT(*)")

	// Results query with pagination
	offset := 0
	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Page > 0 {
		offset = (filter.Page - 1) * limit
	}

	resultsQuery := baseQuery.Columns(
		"clawback_id",
		"policy_number",
		"agent_id",
		"original_commission",
		"clawback_amount",
		"recovered_amount",
		"pending_amount",
		"clawback_percentage",
		"clawback_reason",
		"clawback_status",
		"policy_age",
		"trigger_date",
		"policy_inception_date",
		"policy_end_date",
		"recovery_start_date",
		"recovery_end_date",
		"recovery_schedule",
		"installment_months",
		"approved_by",
		"approved_at",
		"waived_by",
		"waived_at",
		"waiver_reason",
		"workflow_id",
		"workflow_state",
		"notes",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).OrderBy("trigger_date DESC").Limit(uint64(limit)).Offset(uint64(offset))

	// Use batch to execute both queries in one round trip
	batch := &pgx.Batch{}
	var totalCount int64
	var clawbacks []domain.Clawback

	// Queue count query
	scanCountFn := func(row pgx.CollectableRow) (int64, error) {
		var count int64
		err := row.Scan(&count)
		return count, err
	}
	dblib.QueueReturnRow(batch, countQuery, scanCountFn, &totalCount)

	// Queue results query
	scanFn := pgx.RowToStructByName[domain.Clawback]
	dblib.QueueReturn(batch, resultsQuery, scanFn, &clawbacks)

	// Execute batch
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	return clawbacks, totalCount, nil
}

// UpdateClawbackStatus updates the status of a clawback
func (r *ClawbackRepository) UpdateClawbackStatus(
	ctx context.Context,
	clawbackID int64,
	status domain.ClawbackStatus,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutShort"))
	defer cancel()

	q := dblib.Psql.Update(clawbackTable).
		Set("clawback_status", string(status)).
		Set("updated_at", time.Now()).
		Set("updated_by", updatedBy).
		Set("version", sq.Expr("version + 1")).
		Where(sq.Eq{"clawback_id": clawbackID, "deleted_at": nil})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// RecordRecovery records a clawback recovery transaction
// Uses batch for atomic execution of insert + update
func (r *ClawbackRepository) RecordRecovery(
	ctx context.Context,
	recovery *domain.ClawbackRecovery,
	updateClawback bool,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Set timestamps
	now := time.Now()
	recovery.CreatedAt = now
	recovery.UpdatedAt = now
	recovery.Version = 1

	// Create batch for atomic execution
	batch := &pgx.Batch{}
	var recoveryID int64

	// Insert recovery record
	insertQ := dblib.Psql.Insert(clawbackRecoveryTable).Columns(
		"clawback_id",
		"installment_number",
		"scheduled_amount",
		"recovered_amount",
		"recovery_date",
		"recovery_method",
		"statement_id",
		"disbursement_id",
		"transaction_ref",
		"recovery_status",
		"failure_reason",
		"retry_count",
		"next_retry_date",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"version",
	).Values(
		recovery.ClawbackID,
		recovery.InstallmentNumber,
		recovery.ScheduledAmount,
		recovery.RecoveredAmount,
		recovery.RecoveryDate,
		recovery.RecoveryMethod,
		recovery.StatementID,
		recovery.DisbursementID,
		recovery.TransactionRef,
		recovery.RecoveryStatus,
		recovery.FailureReason,
		recovery.RetryCount,
		recovery.NextRetryDate,
		recovery.CreatedAt,
		recovery.UpdatedAt,
		recovery.CreatedBy,
		recovery.UpdatedBy,
		recovery.Version,
	).Suffix("RETURNING recovery_id")

	scanFn := func(row pgx.CollectableRow) (int64, error) {
		var id int64
		err := row.Scan(&id)
		return id, err
	}

	// Queue insert with RETURNING
	dblib.QueueReturnRow(batch, insertQ, scanFn, &recoveryID)

	// Update clawback amounts if requested
	if updateClawback && recovery.RecoveryStatus == "COMPLETED" {
		updateQ := dblib.Psql.Update(clawbackTable).
			Set("recovered_amount", sq.Expr("recovered_amount + ?", recovery.RecoveredAmount)).
			Set("pending_amount", sq.Expr("pending_amount - ?", recovery.RecoveredAmount)).
			Set("updated_at", now).
			Set("version", sq.Expr("version + 1")).
			Where(sq.Eq{"clawback_id": recovery.ClawbackID})

		// Queue update
		dblib.QueueExecRow(batch, updateQ)
	}

	// Execute batch (implicit transaction)
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// Recovery ID is auto-populated by QueueReturnRow
	recovery.RecoveryID = recoveryID

	return nil
}

// GetRecoveriesByClawbackID retrieves all recovery transactions for a clawback
func (r *ClawbackRepository) GetRecoveriesByClawbackID(
	ctx context.Context,
	clawbackID int64,
) ([]domain.ClawbackRecovery, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	q := dblib.Psql.Select(
		"recovery_id",
		"clawback_id",
		"installment_number",
		"scheduled_amount",
		"recovered_amount",
		"recovery_date",
		"recovery_method",
		"statement_id",
		"disbursement_id",
		"transaction_ref",
		"recovery_status",
		"failure_reason",
		"retry_count",
		"next_retry_date",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"version",
	).From(clawbackRecoveryTable).Where(sq.Eq{"clawback_id": clawbackID}).
		OrderBy("installment_number ASC")

	scanFn := pgx.RowToStructByName[domain.ClawbackRecovery]
	return dblib.SelectRows(ctx, r.db, q, scanFn)
}
