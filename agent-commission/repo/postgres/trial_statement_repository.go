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

// TrialStatementRepository handles trial statement data operations
// Implements: FR-IC-COM-004 (Trial statement generation)
// Business Rule: BR-IC-COM-002 (Trial statement approval), BR-IC-COM-009 (7-day SLA)
type TrialStatementRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewTrialStatementRepository creates a new trial statement repository
func NewTrialStatementRepository(db *dblib.DB, cfg *config.Config) *TrialStatementRepository {
	return &TrialStatementRepository{db: db, cfg: cfg}
}

const (
	trialStatementTable = "trial_statements"
)

// CreateTrialStatement creates a new trial statement
// Implements: FR-IC-COM-004
func (r *TrialStatementRepository) CreateTrialStatement(
	ctx context.Context,
	stmt *domain.TrialStatement,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	q := dblib.Psql.Insert(trialStatementTable).Columns(
		"statement_number",
		"agent_id",
		"statement_date",
		"from_date",
		"to_date",
		"total_policies",
		"total_gross_commission",
		"total_tds",
		"total_net_commission",
		"statement_status",
		"processing_unit",
		"batch_id",
		"workflow_state",
	).Values(
		stmt.StatementNumber,
		stmt.AgentID,
		stmt.StatementDate,
		stmt.FromDate,
		stmt.ToDate,
		stmt.TotalPolicies,
		stmt.TotalGrossCommission,
		stmt.TotalTDS,
		stmt.TotalNetCommission,
		stmt.StatementStatus,
		stmt.ProcessingUnit,
		stmt.BatchID,
		stmt.WorkflowState,
	).Suffix("RETURNING trial_statement_id")

	scanFn := func(row pgx.Row) (int64, error) {
		var id int64
		err := row.Scan(&id)
		return id, err
	}

	return dblib.InsertReturning(ctx, r.db, q, scanFn, &stmt.TrialStatementID)
}

// GetTrialStatementByID retrieves a trial statement by ID
func (r *TrialStatementRepository) GetTrialStatementByID(
	ctx context.Context,
	statementID int64,
) (domain.TrialStatement, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"trial_statement_id",
		"statement_number",
		"agent_id",
		"statement_date",
		"from_date",
		"to_date",
		"total_policies",
		"total_gross_commission",
		"total_tds",
		"total_net_commission",
		"statement_status",
		"approved_by",
		"approved_at",
		"approval_remarks",
		"processing_unit",
		"batch_id",
		"workflow_state",
		"workflow_state_updated_at",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).From(trialStatementTable).Where(sq.And{
		sq.Eq{"trial_statement_id": statementID},
		sq.Eq{"deleted_at": nil},
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.TrialStatement])
}

// SearchTrialStatements searches for trial statements with filters
// Implements: FR-IC-COM-004, BR-IC-COM-009
// Uses batch optimization to fetch count + results in single round trip
func (r *TrialStatementRepository) SearchTrialStatements(
	ctx context.Context,
	filter *domain.TrialStatementFilter,
) ([]domain.TrialStatement, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Build base query with filters
	baseQuery := dblib.Psql.Select().From(trialStatementTable).Where(sq.Eq{"deleted_at": nil})

	// Apply filters
	if filter.BatchID != nil && *filter.BatchID != "" {
		baseQuery = baseQuery.Where(sq.Eq{"batch_id": *filter.BatchID})
	}
	if filter.AgentID != nil && *filter.AgentID != "" {
		baseQuery = baseQuery.Where(sq.Eq{"agent_id": *filter.AgentID})
	}
	if filter.Status != nil {
		baseQuery = baseQuery.Where(sq.Eq{"statement_status": string(*filter.Status)})
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
		"trial_statement_id",
		"statement_number",
		"agent_id",
		"statement_date",
		"from_date",
		"to_date",
		"total_policies",
		"total_gross_commission",
		"total_tds",
		"total_net_commission",
		"statement_status",
		"approved_by",
		"approved_at",
		"approval_remarks",
		"processing_unit",
		"batch_id",
		"workflow_state",
		"workflow_state_updated_at",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).OrderBy("created_at DESC").Limit(uint64(limit)).Offset(uint64(offset))

	// Use batch to execute both queries in one round trip
	batch := &pgx.Batch{}
	var totalCount int64
	var statements []domain.TrialStatement

	// Queue count query
	scanCountFn := func(row pgx.Row) (int64, error) {
		var count int64
		err := row.Scan(&count)
		return count, err
	}
	dblib.QueueReturnRow(batch, countQuery, scanCountFn, &totalCount)

	// Queue results query
	scanFn := pgx.RowToStructByName[domain.TrialStatement]
	dblib.QueueReturn(batch, resultsQuery, scanFn, &statements)

	// Execute batch - QueueReturn functions already populate the pointers
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// QueueReturnRow and QueueReturn have already populated totalCount and statements
	// No need for manual processing
	return statements, totalCount, nil
}

// ApproveTrialStatement approves a trial statement
// Implements: FR-IC-COM-006, BR-IC-COM-002
func (r *TrialStatementRepository) ApproveTrialStatement(
	ctx context.Context,
	statementID int64,
	approvedBy string,
	remarks *string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	now := time.Now()

	q := dblib.Psql.Update(trialStatementTable).
		Set("statement_status", domain.TrialStatementStatusApproved).
		Set("approved_by", approvedBy).
		Set("approved_at", now).
		Set("approval_remarks", remarks).
		Set("updated_at", now).
		Set("version", sq.Expr("version + 1")).
		Where(sq.And{
			sq.Eq{"trial_statement_id": statementID},
			sq.Eq{"statement_status": domain.TrialStatementStatusPendingApproval},
			sq.Eq{"deleted_at": nil},
		})

	return dblib.Update(ctx, r.db, q)
}

// GenerateStatementNumber generates a unique statement number
// Format: TS_YYYYMMDD_NNNNNNN
func GenerateStatementNumber() string {
	now := time.Now()
	return fmt.Sprintf("TS_%s_%07d", now.Format("20060102"), now.UnixNano()%10000000)
}
