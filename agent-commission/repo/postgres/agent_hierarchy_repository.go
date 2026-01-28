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

// AgentHierarchyRepository handles agent hierarchy relationships
// BR-IC-AH-001: Advisors MUST be linked to existing Advisor Coordinator
type AgentHierarchyRepository struct {
	db  *dblib.DB
	cfg *config.Config
}

// NewAgentHierarchyRepository creates a new agent hierarchy repository
func NewAgentHierarchyRepository(db *dblib.DB, cfg *config.Config) *AgentHierarchyRepository {
	return &AgentHierarchyRepository{db: db, cfg: cfg}
}

// CreateHierarchyRelationship creates a new agent-coordinator relationship
// BR-IC-AH-001: Link Advisor to Coordinator
func (r *AgentHierarchyRepository) CreateHierarchyRelationship(
	ctx context.Context,
	agentID int64,
	agentCode string,
	coordinatorID int64,
	coordinatorCode string,
	effectiveFrom time.Time,
	createdBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	ins := dblib.Psql.Insert(agentHierarchyTable).
		Columns(
			"agent_id", "agent_code", "coordinator_id", "coordinator_code",
			"hierarchy_level", "effective_from_date", "is_active",
			"created_at", "created_by",
		).
		Values(
			agentID, agentCode, coordinatorID, coordinatorCode,
			1, // Direct report
			effectiveFrom, true,
			sq.Expr("NOW()"), createdBy,
		)

	_, err := dblib.Insert(ctx, r.db, ins)
	return err
}

// GetActiveCoordinatorForAgent retrieves the currently active coordinator for an agent
// BR-IC-AH-001: Get agent's current coordinator
func (r *AgentHierarchyRepository) GetActiveCoordinatorForAgent(ctx context.Context, agentID int64) (*domain.AgentHierarchy, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutLow"))
	defer cancel()

	q := dblib.Psql.Select(
		"agent_hierarchy_id", "agent_id", "agent_code",
		"coordinator_id", "coordinator_code", "hierarchy_level",
		"effective_from_date", "effective_to_date", "is_active",
		"created_at", "created_by", "updated_at", "updated_by",
	).From(agentHierarchyTable).Where(sq.And{
		sq.Eq{"agent_id": agentID},
		sq.Eq{"is_active": true},
		sq.Eq{"effective_to_date": nil}, // Currently active
	})

	return dblib.SelectOne(ctx, r.db, q, pgx.RowToStructByName[domain.AgentHierarchy])
}

// GetSubordinatesForCoordinator retrieves all active subordinates for a coordinator
// Used for coordinator dashboard and reporting
func (r *AgentHierarchyRepository) GetSubordinatesForCoordinator(ctx context.Context, coordinatorID int64) ([]domain.AgentHierarchy, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	q := dblib.Psql.Select(
		"agent_hierarchy_id", "agent_id", "agent_code",
		"coordinator_id", "coordinator_code", "hierarchy_level",
		"effective_from_date", "effective_to_date", "is_active",
		"created_at", "created_by",
	).From(agentHierarchyTable).Where(sq.And{
		sq.Eq{"coordinator_id": coordinatorID},
		sq.Eq{"is_active": true},
		sq.Eq{"effective_to_date": nil},
	}).OrderBy("effective_from_date DESC")

	return dblib.SelectRows(ctx, r.db, q, pgx.RowToStructByName[domain.AgentHierarchy])
}

// UpdateHierarchyRelationship updates an existing hierarchy relationship (e.g., coordinator change)
// Deactivates old relationship and creates new one
func (r *AgentHierarchyRepository) UpdateHierarchyRelationship(
	ctx context.Context,
	agentID int64,
	newCoordinatorID int64,
	newCoordinatorCode string,
	effectiveFrom time.Time,
	updatedBy string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.GetDuration("db.QueryTimeoutMed"))
	defer cancel()

	return r.db.WithTx(ctx, func(tx pgx.Tx) error {
		// Step 1: Deactivate old relationship
		upd := dblib.Psql.Update(agentHierarchyTable).
			Set("is_active", false).
			Set("effective_to_date", effectiveFrom).
			Set("updated_at", sq.Expr("NOW()")).
			Set("updated_by", updatedBy).
			Where(sq.And{
				sq.Eq{"agent_id": agentID},
				sq.Eq{"is_active": true},
			})

		if err := dblib.TxExec(ctx, tx, upd); err != nil {
			return err
		}

		// Step 2: Get agent code
		var agentCode string
		agentQuery := dblib.Psql.Select("agent_code").
			From(agentProfileTable).
			Where(sq.Eq{"agent_profile_id": agentID})

		sql, args, err := agentQuery.ToSql()
		if err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, sql, args...).Scan(&agentCode); err != nil {
			return err
		}

		// Step 3: Create new relationship
		ins := dblib.Psql.Insert(agentHierarchyTable).
			Columns(
				"agent_id", "agent_code", "coordinator_id", "coordinator_code",
				"hierarchy_level", "effective_from_date", "is_active",
				"created_at", "created_by",
			).
			Values(
				agentID, agentCode, newCoordinatorID, newCoordinatorCode,
				1, effectiveFrom, true,
				sq.Expr("NOW()"), updatedBy,
			)

		return dblib.TxExec(ctx, tx, ins)
	})
}
