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

// DisbursementRepository handles disbursement data operations
// FR-IC-COM-009: Cheque disbursement
// FR-IC-COM-010: EFT disbursement via PFMS
type DisbursementRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewDisbursementRepository creates a new disbursement repository
func NewDisbursementRepository(db *dblib.DB, cfg *config.Config) *DisbursementRepository {
	return &DisbursementRepository{db: db, cfg: cfg}
}

const (
	disbursementTable = "disbursements"
)

// CreateDisbursement creates a new disbursement record
// Implements: FR-IC-COM-009, FR-IC-COM-010
// VR-IC-COM-004: Disbursement validation
func (r *DisbursementRepository) CreateDisbursement(
	ctx context.Context,
	disbursement *domain.Disbursement,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Calculate SLA deadline (10 working days from initiated_at)
	disbursement.SLADeadline = calculateWorkingDays(disbursement.InitiatedAt, 10)

	q := dblib.Psql.Insert(disbursementTable).Columns(
		"final_statement_id",
		"agent_id",
		"agent_code",
		"disbursement_mode",
		"disbursement_status",
		"total_gross_commission",
		"total_tds",
		"total_net_commission",
		"cheque_number",
		"cheque_date",
		"bank_account_number",
		"bank_name",
		"bank_branch",
		"ifsc_code",
		"account_holder_name",
		"initiated_at",
		"sla_deadline",
		"workflow_id",
		"workflow_state",
		"created_by",
		"updated_by",
	).Values(
		disbursement.FinalStatementID,
		disbursement.AgentID,
		disbursement.AgentCode,
		disbursement.DisbursementMode,
		disbursement.DisbursementStatus,
		disbursement.TotalGrossCommission,
		disbursement.TotalTDS,
		disbursement.TotalNetCommission,
		disbursement.ChequeNumber,
		disbursement.ChequeDate,
		disbursement.BankAccountNumber,
		disbursement.BankName,
		disbursement.BankBranch,
		disbursement.IFSCCode,
		disbursement.AccountHolderName,
		disbursement.InitiatedAt,
		disbursement.SLADeadline,
		disbursement.WorkflowID,
		disbursement.WorkflowState,
		disbursement.CreatedBy,
		disbursement.UpdatedBy,
	).Suffix("RETURNING disbursement_id")

	scanFn := func(row pgx.CollectableRow) (int64, error) {
		var id int64
		err := row.Scan(&id)
		return id, err
	}

	id, err := dblib.InsertReturning(ctx, r.db, q, scanFn)
	if err != nil {
		return err
	}
	disbursement.DisbursementID = id
	return nil
}

// GetDisbursementByID retrieves a disbursement by ID
func (r *DisbursementRepository) GetDisbursementByID(
	ctx context.Context,
	disbursementID int64,
) (domain.Disbursement, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"disbursement_id", "final_statement_id", "agent_id", "agent_code",
		"disbursement_mode", "disbursement_status",
		"total_gross_commission", "total_tds", "total_net_commission",
		"cheque_number", "cheque_date",
		"bank_account_number", "bank_name", "bank_branch", "ifsc_code", "account_holder_name", "utr_number",
		"initiated_at", "sla_deadline", "processed_at", "completed_at", "sla_breached", "breach_reason_id",
		"failure_reason", "failure_details", "retry_count",
		"workflow_id", "workflow_run_id", "workflow_state", "workflow_state_updated_at",
		"voucher_number", "posted_to_gl", "gl_posted_at",
		"created_at", "created_by", "updated_at", "updated_by", "deleted_at", "version",
	).From(disbursementTable).Where(sq.And{
		sq.Eq{"disbursement_id": disbursementID},
		sq.Eq{"deleted_at": nil},
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.Disbursement])
}

// SearchDisbursements searches for disbursements with filters
// Uses batch optimization to fetch count + results in single round trip
func (r *DisbursementRepository) SearchDisbursements(
	ctx context.Context,
	filter *domain.DisbursementFilter,
) ([]domain.Disbursement, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	// Build base query with filters
	baseQuery := dblib.Psql.Select().From(disbursementTable).Where(sq.Eq{"deleted_at": nil})

	// Apply filters
	if filter.AgentID != nil && *filter.AgentID != "" {
		baseQuery = baseQuery.Where(sq.Eq{"agent_id": *filter.AgentID})
	}
	if filter.DisbursementMode != nil {
		baseQuery = baseQuery.Where(sq.Eq{"disbursement_mode": string(*filter.DisbursementMode)})
	}
	if filter.DisbursementStatus != nil {
		baseQuery = baseQuery.Where(sq.Eq{"disbursement_status": string(*filter.DisbursementStatus)})
	}
	if filter.FromDate != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"initiated_at": *filter.FromDate})
	}
	if filter.ToDate != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"initiated_at": *filter.ToDate})
	}
	if filter.SLABreached != nil && *filter.SLABreached {
		baseQuery = baseQuery.Where(sq.Eq{"sla_breached": true})
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
		"disbursement_id", "final_statement_id", "agent_id", "agent_code",
		"disbursement_mode", "disbursement_status",
		"total_gross_commission", "total_tds", "total_net_commission",
		"cheque_number", "cheque_date",
		"bank_account_number", "bank_name", "bank_branch", "ifsc_code", "account_holder_name", "utr_number",
		"initiated_at", "sla_deadline", "processed_at", "completed_at", "sla_breached", "breach_reason_id",
		"failure_reason", "failure_details", "retry_count",
		"workflow_id", "workflow_run_id", "workflow_state", "workflow_state_updated_at",
		"voucher_number", "posted_to_gl", "gl_posted_at",
		"created_at", "created_by", "updated_at", "updated_by", "deleted_at", "version",
	).OrderBy("initiated_at DESC").Limit(uint64(limit)).Offset(uint64(offset))

	// Use batch to execute both queries in one round trip
	batch := &pgx.Batch{}
	var totalCount int64
	var disbursements []domain.Disbursement

	// Queue count query
	scanCountFn := func(row pgx.CollectableRow) (int64, error) {
		var count int64
		err := row.Scan(&count)
		return count, err
	}
	dblib.QueueReturnRow(batch, countQuery, scanCountFn, &totalCount)

	// Queue results query
	scanFn := pgx.RowToStructByName[domain.Disbursement]
	dblib.QueueReturn(batch, resultsQuery, scanFn, &disbursements)

	// Execute batch - QueueReturn functions already populate the pointers
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// QueueReturnRow and QueueReturn have already populated totalCount and disbursements
	// No need for manual processing
	return disbursements, totalCount, nil
}

// UpdateDisbursementStatus updates the disbursement status
func (r *DisbursementRepository) UpdateDisbursementStatus(
	ctx context.Context,
	disbursementID int64,
	status domain.DisbursementStatus,
	workflowState *string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	q := dblib.Psql.Update(disbursementTable).
		Set("disbursement_status", status).
		Set("workflow_state", workflowState).
		Set("workflow_state_updated_at", time.Now()).
		Set("updated_at", time.Now()).
		Set("version", sq.Expr("version + 1")).
		Where(sq.And{
			sq.Eq{"disbursement_id": disbursementID},
			sq.Eq{"deleted_at": nil},
		})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// CompleteDisbursement marks a disbursement as completed or failed
func (r *DisbursementRepository) CompleteDisbursement(
	ctx context.Context,
	disbursementID int64,
	status domain.DisbursementStatus,
	utrNumber *string,
	failureReason *domain.PaymentFailureReason,
	failureDetails *string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMedium"))
	defer cancel()

	now := time.Now()

	q := dblib.Psql.Update(disbursementTable).
		Set("disbursement_status", status).
		Set("completed_at", now).
		Set("updated_at", now).
		Set("version", sq.Expr("version + 1"))

	if utrNumber != nil {
		q = q.Set("utr_number", *utrNumber)
	}
	if failureReason != nil {
		q = q.Set("failure_reason", *failureReason)
	}
	if failureDetails != nil {
		q = q.Set("failure_details", *failureDetails)
	}

	q = q.Where(sq.And{
		sq.Eq{"disbursement_id": disbursementID},
		sq.Eq{"deleted_at": nil},
	})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// IncrementRetryCount increments the retry count for a failed disbursement
func (r *DisbursementRepository) IncrementRetryCount(
	ctx context.Context,
	disbursementID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Update(disbursementTable).
		Set("retry_count", sq.Expr("retry_count + 1")).
		Set("updated_at", time.Now()).
		Where(sq.And{
			sq.Eq{"disbursement_id": disbursementID},
			sq.Eq{"deleted_at": nil},
		})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// PostToGL marks the disbursement as posted to General Ledger
// INT-IC-004: Accounting system integration
func (r *DisbursementRepository) PostToGL(
	ctx context.Context,
	disbursementID int64,
	voucherNumber string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	now := time.Now()

	q := dblib.Psql.Update(disbursementTable).
		Set("voucher_number", voucherNumber).
		Set("posted_to_gl", true).
		Set("gl_posted_at", now).
		Set("updated_at", now).
		Where(sq.And{
			sq.Eq{"disbursement_id": disbursementID},
			sq.Eq{"deleted_at": nil},
		})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// UpdateSLABreach marks a disbursement as SLA breached
// BR-IC-COM-011: 10-working-day SLA tracking
func (r *DisbursementRepository) UpdateSLABreach(
	ctx context.Context,
	disbursementID int64,
	breachReasonID *int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Update(disbursementTable).
		Set("sla_breached", true).
		Set("breach_reason_id", breachReasonID).
		Set("updated_at", time.Now()).
		Where(sq.And{
			sq.Eq{"disbursement_id": disbursementID},
			sq.Eq{"deleted_at": nil},
		})

	_, err := dblib.Update(ctx, r.db, q)
	return err
}

// calculateWorkingDays calculates a date N working days from the start date
// Excludes weekends (Saturday, Sunday) - can be extended to exclude holidays
func calculateWorkingDays(start time.Time, workingDays int) time.Time {
	days := 0
	current := start

	for days < workingDays {
		current = current.AddDate(0, 0, 1)
		// Skip weekends
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			days++
		}
	}

	return current
}
