package postgres

import (
	"context"

	"agent-commission/core/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	dblib "gitlab.cept.gov.in/it-2.0-common/api-db"
)

// ReferenceDataRepository handles reference/master data operations
// BR-IC-AH-002: Coordinator circle and division assignment
type ReferenceDataRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewReferenceDataRepository creates a new reference data repository
func NewReferenceDataRepository(db *dblib.DB, cfg *config.Config) *ReferenceDataRepository {
	return &ReferenceDataRepository{db: db, cfg: cfg}
}

const (
	circleTable      = "ref_circles"
	divisionTable    = "ref_divisions"
	productPlanTable = "ref_product_plans"
)

// GetAllCircles retrieves all active circles
// BR-IC-AH-002: Used for coordinator circle assignment
func (r *ReferenceDataRepository) GetAllCircles(ctx context.Context) ([]domain.Circle, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"circle_id", "circle_code", "circle_name", "region_code", "is_active", "created_at", "updated_at",
	).From(circleTable).Where(sq.Eq{"is_active": true}).OrderBy("circle_name ASC")

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.Circle])
}

// GetCircleByID retrieves a circle by ID
// BR-IC-AH-002: Validate circle during coordinator assignment
func (r *ReferenceDataRepository) GetCircleByID(ctx context.Context, circleID int64) (*domain.Circle, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"circle_id", "circle_code", "circle_name", "region_code", "is_active", "created_at", "updated_at",
	).From(circleTable).Where(sq.And{
		sq.Eq{"circle_id": circleID},
		sq.Eq{"is_active": true},
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.Circle])
}

// GetDivisionsByCircle retrieves all divisions for a given circle
// BR-IC-AH-002: Used for coordinator division assignment
func (r *ReferenceDataRepository) GetDivisionsByCircle(ctx context.Context, circleID int64) ([]domain.Division, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"division_id", "division_code", "division_name", "circle_id", "circle_code", "circle_name",
		"head_office", "is_active", "created_at", "updated_at",
	).From(divisionTable).Where(sq.And{
		sq.Eq{"circle_id": circleID},
		sq.Eq{"is_active": true},
	}).OrderBy("division_name ASC")

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.Division])
}

// GetDivisionByID retrieves a division by ID
// BR-IC-AH-002: Validate division during coordinator assignment
func (r *ReferenceDataRepository) GetDivisionByID(ctx context.Context, divisionID int64) (*domain.Division, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"division_id", "division_code", "division_name", "circle_id", "circle_code", "circle_name",
		"head_office", "is_active", "created_at", "updated_at",
	).From(divisionTable).Where(sq.And{
		sq.Eq{"division_id": divisionID},
		sq.Eq{"is_active": true},
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.Division])
}

// GetAllDivisions retrieves all active divisions
func (r *ReferenceDataRepository) GetAllDivisions(ctx context.Context) ([]domain.Division, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"division_id", "division_code", "division_name", "circle_id", "circle_code", "circle_name",
		"head_office", "is_active", "created_at", "updated_at",
	).From(divisionTable).Where(sq.Eq{"is_active": true}).OrderBy("circle_name ASC", "division_name ASC")

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.Division])
}

// GetProductPlans retrieves all active product plans
// Used for commission rate configuration and lookup
func (r *ReferenceDataRepository) GetProductPlans(ctx context.Context, productType *string) ([]domain.ProductPlan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"product_plan_id", "plan_code", "plan_name", "product_type", "plan_description",
		"is_active", "created_at", "updated_at",
	).From(productPlanTable).Where(sq.Eq{"is_active": true})

	if productType != nil && *productType != "" {
		q = q.Where(sq.Eq{"product_type": *productType})
	}

	q = q.OrderBy("product_type ASC", "plan_name ASC")

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.ProductPlan])
}

// GetProductPlanByCode retrieves a product plan by code
func (r *ReferenceDataRepository) GetProductPlanByCode(ctx context.Context, planCode string) (*domain.ProductPlan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"product_plan_id", "plan_code", "plan_name", "product_type", "plan_description",
		"is_active", "created_at", "updated_at",
	).From(productPlanTable).Where(sq.And{
		sq.Eq{"plan_code": planCode},
		sq.Eq{"is_active": true},
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.ProductPlan])
}

// GetCoordinatorsByCircle retrieves all active coordinators in a circle
// Used for agent onboarding - coordinator selection
func (r *ReferenceDataRepository) GetCoordinatorsByCircle(ctx context.Context, circleID int64) ([]domain.AgentProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"agent_profile_id", "agent_code", "first_name", "middle_name", "last_name",
		"circle_id", "circle_name", "division_id", "division_name", "status",
	).From(agentProfileTable).Where(sq.And{
		sq.Eq{"agent_type": domain.AgentTypeAdvisorCoordinator},
		sq.Eq{"circle_id": circleID},
		sq.Eq{"status": domain.AgentStatusActive},
		sq.Eq{"deleted_at": nil},
	}).OrderBy("first_name ASC", "last_name ASC")

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.AgentProfile])
}

// GetCoordinatorsByDivision retrieves all active coordinators in a division
func (r *ReferenceDataRepository) GetCoordinatorsByDivision(ctx context.Context, divisionID int64) ([]domain.AgentProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"agent_profile_id", "agent_code", "first_name", "middle_name", "last_name",
		"circle_id", "circle_name", "division_id", "division_name", "status",
	).From(agentProfileTable).Where(sq.And{
		sq.Eq{"agent_type": domain.AgentTypeAdvisorCoordinator},
		sq.Eq{"division_id": divisionID},
		sq.Eq{"status": domain.AgentStatusActive},
		sq.Eq{"deleted_at": nil},
	}).OrderBy("first_name ASC", "last_name ASC")

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.AgentProfile])
}
