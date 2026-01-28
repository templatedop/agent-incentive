package postgres

import (
	"context"

	"agent-commission/core/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/api-db"
)

// CommissionHistoryRepository handles commission history queries
// FR-IC-COM-011: Commission history inquiry
type CommissionHistoryRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewCommissionHistoryRepository creates a new commission history repository
func NewCommissionHistoryRepository(db *dblib.DB, cfg *config.Config) *CommissionHistoryRepository {
	return &CommissionHistoryRepository{db: db, cfg: cfg}
}

const (
	commissionTransactionTable = "commission_transactions"
)

// SearchCommissionHistory searches commission transactions with filters
// Implements: FR-IC-COM-011, BR-IC-COM-009
// Uses batch optimization to fetch count + results in single round trip
func (r *CommissionHistoryRepository) SearchCommissionHistory(
	ctx context.Context,
	filter *domain.CommissionHistoryFilter,
) ([]domain.CommissionTransaction, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Build base query with filters
	baseQuery := dblib.Psql.Select().From(commissionTransactionTable).Where(sq.Eq{"deleted_at": nil})

	// Apply filters
	if filter.AgentID != nil && *filter.AgentID != "" {
		baseQuery = baseQuery.Where(sq.Eq{"agent_id": *filter.AgentID})
	}
	if filter.PolicyNumber != nil && *filter.PolicyNumber != "" {
		baseQuery = baseQuery.Where(sq.Eq{"policy_number": *filter.PolicyNumber})
	}
	if filter.CommissionType != nil {
		baseQuery = baseQuery.Where(sq.Eq{"commission_type": string(*filter.CommissionType)})
	}
	if filter.CommissionStatus != nil {
		baseQuery = baseQuery.Where(sq.Eq{"commission_status": string(*filter.CommissionStatus)})
	}
	if filter.FromDate != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"commission_date": *filter.FromDate})
	}
	if filter.ToDate != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"commission_date": *filter.ToDate})
	}
	if filter.BatchID != nil && *filter.BatchID != "" {
		baseQuery = baseQuery.Where(sq.Eq{"batch_id": *filter.BatchID})
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
		"commission_id",
		"agent_id",
		"policy_number",
		"commission_type",
		"product_type",
		"annualised_premium",
		"rate_percentage",
		"gross_commission",
		"tds_rate",
		"tds_amount",
		"net_commission",
		"commission_date",
		"commission_status",
		"trial_statement_id",
		"final_statement_id",
		"disbursement_id",
		"batch_id",
		"workflow_state",
		"workflow_state_updated_at",
		"created_at",
		"updated_at",
		"created_by",
		"updated_by",
		"deleted_at",
		"version",
	).OrderBy("commission_date DESC").Limit(uint64(limit)).Offset(uint64(offset))

	// Use batch to execute both queries in one round trip
	batch := &pgx.Batch{}
	var totalCount int64
	var transactions []domain.CommissionTransaction

	// Queue count query
	scanCountFn := func(row pgx.CollectableRow) (int64, error) {
		var count int64
		err := row.Scan(&count)
		return count, err
	}
	dblib.QueueReturnRow(batch, countQuery, scanCountFn, &totalCount)

	// Queue results query
	scanFn := pgx.RowToStructByName[domain.CommissionTransaction]
	dblib.QueueReturn(batch, resultsQuery, scanFn, &transactions)

	// Execute batch - QueueReturn functions already populate the pointers
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// QueueReturnRow and QueueReturn have already populated totalCount and transactions
	// No need for manual processing
	return transactions, totalCount, nil
}
