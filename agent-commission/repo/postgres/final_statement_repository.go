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

// FinalStatementRepository handles final statement data operations
// Implements: FR-IC-COM-008 (Final statement generation)
// Business Rule: BR-IC-COM-007 (Final statement generation after trial approval)
type FinalStatementRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewFinalStatementRepository creates a new final statement repository
func NewFinalStatementRepository(db *dblib.DB, cfg *config.Config) *FinalStatementRepository {
	return &FinalStatementRepository{db: db, cfg: cfg}
}

const (
	finalStatementTable = "final_statements"
)

// CreateFinalStatement creates a new final statement
// Implements: FR-IC-COM-008, BR-IC-COM-007
func (r *FinalStatementRepository) CreateFinalStatement(
	ctx context.Context,
	stmt *domain.FinalStatement,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	q := dblib.Psql.Insert(finalStatementTable).Columns(
		"statement_number",
		"trial_statement_id",
		"agent_id",
		"statement_date",
		"total_gross_commission",
		"total_tds",
		"total_net_commission",
		"statement_status",
		"workflow_state",
	).Values(
		stmt.StatementNumber,
		stmt.TrialStatementID,
		stmt.AgentID,
		stmt.StatementDate,
		stmt.TotalGrossCommission,
		stmt.TotalTDS,
		stmt.TotalNetCommission,
		stmt.StatementStatus,
		stmt.WorkflowState,
	).Suffix("RETURNING final_statement_id")

	scanFn := func(row pgx.Row) (int64, error) {
		var id int64
		err := row.Scan(&id)
		return id, err
	}

	id, err := dblib.InsertReturning(ctx, r.db, q, scanFn)
	if err != nil {
		return err
	}
	stmt.FinalStatementID = id
	return nil
}

// GetFinalStatementByID retrieves a final statement by ID
func (r *FinalStatementRepository) GetFinalStatementByID(
	ctx context.Context,
	statementID int64,
) (domain.FinalStatement, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"final_statement_id",
		"statement_number",
		"trial_statement_id",
		"agent_id",
		"statement_date",
		"total_gross_commission",
		"total_tds",
		"total_net_commission",
		"statement_status",
		"pdf_path",
		"workflow_state",
		"workflow_state_updated_at",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).From(finalStatementTable).Where(sq.And{
		sq.Eq{"final_statement_id": statementID},
		sq.Eq{"deleted_at": nil},
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.FinalStatement])
}

// SearchFinalStatements searches for final statements with filters
// Uses batch optimization to fetch count + results in single round trip
func (r *FinalStatementRepository) SearchFinalStatements(
	ctx context.Context,
	filter *domain.FinalStatementFilter,
) ([]domain.FinalStatement, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Build base query with filters
	baseQuery := dblib.Psql.Select().From(finalStatementTable).Where(sq.Eq{"deleted_at": nil})

	// Apply filters
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
		"final_statement_id",
		"statement_number",
		"trial_statement_id",
		"agent_id",
		"statement_date",
		"total_gross_commission",
		"total_tds",
		"total_net_commission",
		"statement_status",
		"pdf_path",
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
	var statements []domain.FinalStatement

	// Queue count query
	scanCountFn := func(row pgx.Row) (int64, error) {
		var count int64
		err := row.Scan(&count)
		return count, err
	}
	dblib.QueueReturnRow(batch, countQuery, scanCountFn, &totalCount)

	// Queue results query
	scanFn := pgx.RowToStructByName[domain.FinalStatement]
	dblib.QueueReturn(batch, resultsQuery, scanFn, &statements)

	// Execute batch - QueueReturn functions already populate the pointers
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// QueueReturnRow and QueueReturn have already populated totalCount and statements
	// No need for manual processing
	return statements, totalCount, nil
}

// UpdateStatementStatus updates the final statement status
func (r *FinalStatementRepository) UpdateStatementStatus(
	ctx context.Context,
	statementID int64,
	status domain.FinalStatementStatus,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	q := dblib.Psql.Update(finalStatementTable).
		Set("statement_status", status).
		Set("updated_at", time.Now()).
		Set("version", sq.Expr("version + 1")).
		Where(sq.And{
			sq.Eq{"final_statement_id": statementID},
			sq.Eq{"deleted_at": nil},
		})

	return dblib.Update(ctx, r.db, q)
}

// GenerateFinalStatementNumber generates a unique statement number
// Format: FS_YYYYMMDD_NNNNNNN
func GenerateFinalStatementNumber() string {
	now := time.Now()
	return fmt.Sprintf("FS_%s_%07d", now.Format("20060102"), now.UnixNano()%10000000)
}
