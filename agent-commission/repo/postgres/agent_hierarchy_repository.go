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
// Optimization: Uses batch instead of transaction (batch provides implicit transaction)
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

	// Create batch - batch provides implicit transaction
	batch := &pgx.Batch{}

	// Step 1: Queue deactivation of old relationship
	upd := dblib.Psql.Update(agentHierarchyTable).
		Set("is_active", false).
		Set("effective_to_date", effectiveFrom).
		Set("updated_at", sq.Expr("NOW()")).
		Set("updated_by", updatedBy).
		Where(sq.And{
			sq.Eq{"agent_id": agentID},
			sq.Eq{"is_active": true},
		})

	if err := dblib.QueueExecRow(batch, upd); err != nil {
		return err
	}

	// Step 2: Queue creation of new relationship (using subquery for agent_code)
	ins := dblib.Psql.Insert(agentHierarchyTable).
		Columns(
			"agent_id", "agent_code", "coordinator_id", "coordinator_code",
			"hierarchy_level", "effective_from_date", "is_active",
			"created_at", "created_by",
		).
		Values(
			agentID,
			sq.Expr("(SELECT agent_code FROM "+agentProfileTable+" WHERE agent_profile_id = ?)", agentID),
			newCoordinatorID, newCoordinatorCode,
			1, effectiveFrom, true,
			sq.Expr("NOW()"), updatedBy,
		)

	if err := dblib.QueueExecRow(batch, ins); err != nil {
		return err
	}

	// Execute batch (implicit transaction)
	batchResults := r.db.Pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// Process results
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchResults.Exec(); err != nil {
			return fmt.Errorf("batch operation %d failed: %w", i, err)
		}
	}

	return nil
}
