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

// SuspenseRepository handles suspense account data operations
// FR-IC-COM-013: Suspense account management
type SuspenseRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewSuspenseRepository creates a new suspense repository
func NewSuspenseRepository(db *dblib.DB, cfg *config.Config) *SuspenseRepository {
	return &SuspenseRepository{db: db, cfg: cfg}
}

const (
	suspenseAccountTable     = "suspense_accounts"
	suspenseTransactionTable = "suspense_transactions"
)

// CreateSuspenseAccount creates a new suspense account entry
// Implements: FR-IC-COM-013, BR-IC-SUS-001
func (r *SuspenseRepository) CreateSuspenseAccount(
	ctx context.Context,
	suspense *domain.SuspenseAccount,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Auto-determine priority if not set
	if suspense.Priority == "" {
		suspense.Priority = domain.DeterminePriority(suspense.SuspenseAmount, suspense.SuspenseReason)
	}

	// Calculate resolution deadline
	resolutionDeadline := domain.CalculateResolutionDeadline(suspense.SuspenseDate, suspense.Priority)
	suspense.ResolutionDeadline = &resolutionDeadline

	// Set initial status
	if suspense.SuspenseStatus == "" {
		suspense.SuspenseStatus = domain.SuspenseStatusOpen
	}

	// Set timestamps
	now := time.Now()
	suspense.CreatedAt = now
	suspense.UpdatedAt = now
	suspense.SuspenseDate = now
	suspense.Version = 1

	// Calculate initial aging
	suspense.UpdateAgingDays()

	// Build insert query
	q := dblib.Psql.Insert(suspenseAccountTable).Columns(
		"agent_id",
		"policy_number",
		"commission_id",
		"disbursement_id",
		"suspense_amount",
		"suspense_reason",
		"suspense_status",
		"suspense_date",
		"aging_days",
		"resolution_deadline",
		"priority",
		"assigned_to",
		"notes",
		"workflow_id",
		"workflow_state",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"version",
	).Values(
		suspense.AgentID,
		suspense.PolicyNumber,
		suspense.CommissionID,
		suspense.DisbursementID,
		suspense.SuspenseAmount,
		string(suspense.SuspenseReason),
		string(suspense.SuspenseStatus),
		suspense.SuspenseDate,
		suspense.AgingDays,
		suspense.ResolutionDeadline,
		suspense.Priority,
		suspense.AssignedTo,
		suspense.Notes,
		suspense.WorkflowID,
		suspense.WorkflowState,
		suspense.CreatedAt,
		suspense.UpdatedAt,
		suspense.CreatedBy,
		suspense.UpdatedBy,
		suspense.Version,
	).Suffix("RETURNING suspense_id")

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

	suspense.SuspenseID = id

	// Create audit transaction
	return r.createSuspenseTransaction(ctx, suspense.SuspenseID, "CREATED", nil, string(suspense.SuspenseStatus), &suspense.SuspenseAmount, nil, suspense.CreatedBy)
}

// GetSuspenseAccountByID retrieves a suspense account by its ID
func (r *SuspenseRepository) GetSuspenseAccountByID(
	ctx context.Context,
	suspenseID int64,
) (domain.SuspenseAccount, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutShort"))
	defer cancel()

	q := dblib.Psql.Select(
		"suspense_id",
		"agent_id",
		"policy_number",
		"commission_id",
		"disbursement_id",
		"suspense_amount",
		"suspense_reason",
		"suspense_status",
		"suspense_date",
		"aging_days",
		"resolution_deadline",
		"resolved_date",
		"resolution_method",
		"resolved_amount",
		"write_off_date",
		"write_off_reason",
		"assigned_to",
		"priority",
		"notes",
		"workflow_id",
		"workflow_state",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).From(suspenseAccountTable).Where(sq.Eq{"suspense_id": suspenseID, "deleted_at": nil})

	scanFn := pgx.RowToStructByName[domain.SuspenseAccount]
	suspense, err := dblib.SelectOne(ctx, r.db, q, scanFn)
	if err != nil {
		return suspense, err
	}

	// Update aging days
	suspense.UpdateAgingDays()
	return suspense, nil
}

// SearchSuspenseAccounts searches suspense accounts with filters
// Uses batch optimization to fetch count + results in single round trip
func (r *SuspenseRepository) SearchSuspenseAccounts(
	ctx context.Context,
	filter *domain.SuspenseFilter,
) ([]domain.SuspenseAccount, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Build base query with filters
	baseQuery := dblib.Psql.Select().From(suspenseAccountTable).Where(sq.Eq{"deleted_at": nil})

	// Apply filters
	if filter.AgentID != nil && *filter.AgentID != "" {
		baseQuery = baseQuery.Where(sq.Eq{"agent_id": *filter.AgentID})
	}
	if filter.PolicyNumber != nil && *filter.PolicyNumber != "" {
		baseQuery = baseQuery.Where(sq.Eq{"policy_number": *filter.PolicyNumber})
	}
	if filter.SuspenseStatus != nil {
		baseQuery = baseQuery.Where(sq.Eq{"suspense_status": string(*filter.SuspenseStatus)})
	}
	if filter.SuspenseReason != nil {
		baseQuery = baseQuery.Where(sq.Eq{"suspense_reason": string(*filter.SuspenseReason)})
	}
	if filter.FromDate != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"suspense_date": *filter.FromDate})
	}
	if filter.ToDate != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"suspense_date": *filter.ToDate})
	}
	if filter.MinAmount != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"suspense_amount": *filter.MinAmount})
	}
	if filter.MaxAmount != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"suspense_amount": *filter.MaxAmount})
	}
	if filter.MinAgingDays != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"aging_days": *filter.MinAgingDays})
	}
	if filter.MaxAgingDays != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"aging_days": *filter.MaxAgingDays})
	}
	if filter.Priority != nil && *filter.Priority != "" {
		baseQuery = baseQuery.Where(sq.Eq{"priority": *filter.Priority})
	}
	if filter.AssignedTo != nil && *filter.AssignedTo != "" {
		baseQuery = baseQuery.Where(sq.Eq{"assigned_to": *filter.AssignedTo})
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
		"suspense_id",
		"agent_id",
		"policy_number",
		"commission_id",
		"disbursement_id",
		"suspense_amount",
		"suspense_reason",
		"suspense_status",
		"suspense_date",
		"aging_days",
		"resolution_deadline",
		"resolved_date",
		"resolution_method",
		"resolved_amount",
		"write_off_date",
		"write_off_reason",
		"assigned_to",
		"priority",
		"notes",
		"workflow_id",
		"workflow_state",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).OrderBy("suspense_date DESC").Limit(uint64(limit)).Offset(uint64(offset))

	// Use batch to execute both queries in one round trip
	batch := &pgx.Batch{}
	var totalCount int64
	var suspenseAccounts []domain.SuspenseAccount

	// Queue count query
	scanCountFn := func(row pgx.CollectableRow) (int64, error) {
		var count int64
		err := row.Scan(&count)
		return count, err
	}
	dblib.QueueReturnRow(batch, countQuery, scanCountFn, &totalCount)

	// Queue results query
	scanFn := pgx.RowToStructByName[domain.SuspenseAccount]
	dblib.QueueReturn(batch, resultsQuery, scanFn, &suspenseAccounts)

	// Execute batch
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// Update aging days for all results
	for i := range suspenseAccounts {
		suspenseAccounts[i].UpdateAgingDays()
	}

	return suspenseAccounts, totalCount, nil
}

// GetAgingReport generates aging report for suspense accounts
// BR-IC-SUS-002: Aging buckets analysis
func (r *SuspenseRepository) GetAgingReport(ctx context.Context) (*domain.SuspenseAgingReport, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Query for aging buckets
	agingQuery := `
		SELECT
			CASE
				WHEN aging_days <= 30 THEN '0-30 days'
				WHEN aging_days <= 60 THEN '31-60 days'
				WHEN aging_days <= 90 THEN '61-90 days'
				WHEN aging_days <= 180 THEN '91-180 days'
				ELSE '180+ days'
			END as aging_bucket,
			COUNT(*) as count,
			SUM(suspense_amount) as total_amount,
			MIN(aging_days) as min_aging_days,
			MAX(aging_days) as max_aging_days
		FROM suspense_accounts
		WHERE deleted_at IS NULL AND suspense_status = 'OPEN'
		GROUP BY aging_bucket
		ORDER BY MIN(aging_days)
	`

	rows, err := r.db.Pool.Query(ctx, agingQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buckets []domain.SuspenseAgingBucket
	var totalEntries int
	var totalAmount float64

	for rows.Next() {
		var bucket domain.SuspenseAgingBucket
		err := rows.Scan(&bucket.AgingBucket, &bucket.Count, &bucket.TotalAmount, &bucket.MinAgingDays, &bucket.MaxAgingDays)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
		totalEntries += bucket.Count
		totalAmount += bucket.TotalAmount
	}

	// Query for overdue entries
	overdueQuery := `
		SELECT COUNT(*), COALESCE(SUM(suspense_amount), 0)
		FROM suspense_accounts
		WHERE deleted_at IS NULL
		  AND suspense_status = 'OPEN'
		  AND resolution_deadline < NOW()
	`

	var overdueCount int
	var overdueAmount float64
	err = r.db.Pool.QueryRow(ctx, overdueQuery).Scan(&overdueCount, &overdueAmount)
	if err != nil {
		return nil, err
	}

	// Build report
	report := &domain.SuspenseAgingReport{
		ReportDate:    time.Now(),
		TotalEntries:  totalEntries,
		TotalAmount:   totalAmount,
		Buckets:       buckets,
		OverdueCount:  overdueCount,
		OverdueAmount: overdueAmount,
		ByReason:      make(map[domain.SuspenseReason]domain.SuspenseAgingStats),
		ByPriority:    make(map[string]domain.SuspenseAgingStats),
	}

	return report, nil
}

// ResolveSuspenseAccount resolves a suspense account entry
func (r *SuspenseRepository) ResolveSuspenseAccount(
	ctx context.Context,
	suspenseID int64,
	resolvedAmount float64,
	resolutionMethod string,
	resolvedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	now := time.Now()

	q := dblib.Psql.Update(suspenseAccountTable).
		Set("suspense_status", string(domain.SuspenseStatusResolved)).
		Set("resolved_date", now).
		Set("resolved_amount", resolvedAmount).
		Set("resolution_method", resolutionMethod).
		Set("updated_at", now).
		Set("updated_by", resolvedBy).
		Set("version", sq.Expr("version + 1")).
		Where(sq.Eq{"suspense_id": suspenseID, "deleted_at": nil})

	_, err := dblib.Update(ctx, r.db, q)
	if err != nil {
		return err
	}

	// Create audit transaction
	remarks := resolutionMethod
	return r.createSuspenseTransaction(ctx, suspenseID, "RESOLVED", ptrString(string(domain.SuspenseStatusOpen)), string(domain.SuspenseStatusResolved), &resolvedAmount, &remarks, resolvedBy)
}

// WriteOffSuspenseAccount writes off a suspense account entry
func (r *SuspenseRepository) WriteOffSuspenseAccount(
	ctx context.Context,
	suspenseID int64,
	writeOffReason string,
	writtenOffBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	now := time.Now()

	q := dblib.Psql.Update(suspenseAccountTable).
		Set("suspense_status", string(domain.SuspenseStatusWriteOff)).
		Set("write_off_date", now).
		Set("write_off_reason", writeOffReason).
		Set("updated_at", now).
		Set("updated_by", writtenOffBy).
		Set("version", sq.Expr("version + 1")).
		Where(sq.Eq{"suspense_id": suspenseID, "deleted_at": nil})

	_, err := dblib.Update(ctx, r.db, q)
	if err != nil {
		return err
	}

	// Create audit transaction
	return r.createSuspenseTransaction(ctx, suspenseID, "WRITE_OFF", ptrString(string(domain.SuspenseStatusOpen)), string(domain.SuspenseStatusWriteOff), nil, &writeOffReason, writtenOffBy)
}

// createSuspenseTransaction creates an audit transaction
func (r *SuspenseRepository) createSuspenseTransaction(
	ctx context.Context,
	suspenseID int64,
	transactionType string,
	oldStatus *string,
	newStatus string,
	amount *float64,
	remarks *string,
	performedBy string,
) error {
	now := time.Now()

	q := dblib.Psql.Insert(suspenseTransactionTable).Columns(
		"suspense_id",
		"transaction_type",
		"old_status",
		"new_status",
		"amount",
		"remarks",
		"performed_by",
		"transaction_date",
		"created_at",
	).Values(
		suspenseID,
		transactionType,
		oldStatus,
		newStatus,
		amount,
		remarks,
		performedBy,
		now,
		now,
	)

	_, err := dblib.Insert(ctx, r.db, q)
	return err
}

// Helper function to create string pointer
func ptrString(s string) *string {
	return &s
}
