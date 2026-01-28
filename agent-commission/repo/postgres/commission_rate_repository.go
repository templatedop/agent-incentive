package postgres

import (
	"context"

	"agent-commission/core/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/api-db"
)

// CommissionRateRepository handles commission rate data operations
// Implements: FR-IC-COM-001 (Commission rate table management)
// Business Rule: BR-IC-COM-006 (Rate structure)
type CommissionRateRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewCommissionRateRepository creates a new commission rate repository
func NewCommissionRateRepository(db *dblib.DB, cfg *config.Config) *CommissionRateRepository {
	return &CommissionRateRepository{db: db, cfg: cfg}
}

const (
	commissionRateTable = "commission_rates"
)

// GetCommissionRates retrieves commission rates based on filter criteria
// Implements: FR-IC-COM-001, BR-IC-COM-006
// API: GET /lookup/commission-rates
//
// Filters:
//   - productType: Filter by product type (PLI, RPLI)
//   - agentType: Filter by agent type
//   - planCode: Filter by product plan code
//
// Returns only active rates if no specific filters provided
func (r *CommissionRateRepository) GetCommissionRates(
	ctx context.Context,
	filter *domain.CommissionRateFilter,
) ([]domain.CommissionRate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"rate_id",
		"rate_percentage",
		"policy_duration_months",
		"product_type",
		"product_plan_code",
		"agent_type",
		"policy_term_years",
		"effective_from",
		"effective_to",
		"is_active",
		"created_at",
		"updated_at",
	).From(commissionRateTable).Where(sq.Eq{"is_active": true})

	// Apply filters if provided
	if filter != nil {
		if filter.ProductType != nil {
			q = q.Where(sq.Eq{"product_type": string(*filter.ProductType)})
		}
		if filter.AgentType != nil && *filter.AgentType != "" {
			q = q.Where(sq.Eq{"agent_type": *filter.AgentType})
		}
		if filter.PlanCode != nil && *filter.PlanCode != "" {
			q = q.Where(sq.Eq{"product_plan_code": *filter.PlanCode})
		}
		// Filter by date range if AsOfDate provided
		if filter.AsOfDate != nil {
			q = q.Where(sq.LtOrEq{"effective_from": *filter.AsOfDate})
			q = q.Where(sq.Or{
				sq.Eq{"effective_to": nil},
				sq.GtOrEq{"effective_to": *filter.AsOfDate},
			})
		}
	}

	// Order by product type, agent type, and plan code for consistent results
	q = q.OrderBy(
		"product_type ASC",
		"agent_type ASC",
		"product_plan_code ASC",
		"policy_term_years ASC",
		"effective_from DESC",
	)

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.CommissionRate])
}

// GetCommissionRateByID retrieves a specific commission rate by ID
// Used for rate lookup in commission calculation
func (r *CommissionRateRepository) GetCommissionRateByID(ctx context.Context, rateID int64) (domain.CommissionRate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"rate_id",
		"rate_percentage",
		"policy_duration_months",
		"product_type",
		"product_plan_code",
		"agent_type",
		"policy_term_years",
		"effective_from",
		"effective_to",
		"is_active",
		"created_at",
		"updated_at",
	).From(commissionRateTable).Where(sq.Eq{"rate_id": rateID})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.CommissionRate])
}

// GetApplicableRate finds the applicable commission rate for a specific policy
// Implements: BR-IC-COM-001 (Commission calculation rules)
//
// Selection logic:
//   - Match: product_type, agent_type, plan_code, policy_term_years
//   - Active and within effective date range
//   - Returns most recent rate if multiple matches
func (r *CommissionRateRepository) GetApplicableRate(
	ctx context.Context,
	productType domain.ProductType,
	agentType string,
	planCode string,
	policyTermYears int,
) (domain.CommissionRate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"rate_id",
		"rate_percentage",
		"policy_duration_months",
		"product_type",
		"product_plan_code",
		"agent_type",
		"policy_term_years",
		"effective_from",
		"effective_to",
		"is_active",
		"created_at",
		"updated_at",
	).From(commissionRateTable).Where(sq.And{
		sq.Eq{"product_type": string(productType)},
		sq.Eq{"agent_type": agentType},
		sq.Eq{"product_plan_code": planCode},
		sq.Eq{"policy_term_years": policyTermYears},
		sq.Eq{"is_active": true},
		sq.LtOrEq{"effective_from": "NOW()"},
		sq.Or{
			sq.Eq{"effective_to": nil},
			sq.GtOrEq{"effective_to": "NOW()"},
		},
	}).OrderBy("effective_from DESC").Limit(1)

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.CommissionRate])
}
