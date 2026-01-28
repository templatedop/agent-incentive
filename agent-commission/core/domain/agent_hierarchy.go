package domain

import "time"

// AgentHierarchy represents the relationship between agents and their coordinators
// BR-IC-AH-001: Advisors MUST be linked to existing Advisor Coordinator
// This table establishes the agent reporting hierarchy
type AgentHierarchy struct {
	// Primary Key
	AgentHierarchyID int64 `db:"agent_hierarchy_id" json:"agent_hierarchy_id"`

	// Agent (the subordinate)
	AgentID   int64  `db:"agent_id" json:"agent_id"`     // Foreign key to agent_profiles
	AgentCode string `db:"agent_code" json:"agent_code"` // Denormalized for performance

	// Coordinator (the supervisor)
	CoordinatorID   int64  `db:"coordinator_id" json:"coordinator_id"`     // Foreign key to agent_profiles
	CoordinatorCode string `db:"coordinator_code" json:"coordinator_code"` // Denormalized

	// Hierarchy Level
	HierarchyLevel int `db:"hierarchy_level" json:"hierarchy_level"` // 1=Direct report, 2=Skip level, etc.

	// Effective Dates
	EffectiveFromDate time.Time  `db:"effective_from_date" json:"effective_from_date"`
	EffectiveToDate   *time.Time `db:"effective_to_date" json:"effective_to_date,omitempty"` // NULL means currently active

	// Flags
	IsActive bool `db:"is_active" json:"is_active"` // Derived from effective dates

	// Remarks (for transfers, changes)
	Remarks *string `db:"remarks" json:"remarks,omitempty"`

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
}

// IsCurrentlyActive returns true if this hierarchy relationship is currently active
func (h *AgentHierarchy) IsCurrentlyActive() bool {
	now := time.Now()
	return h.EffectiveFromDate.Before(now) &&
		(h.EffectiveToDate == nil || h.EffectiveToDate.After(now))
}

// AgentWithCoordinator represents an agent profile with coordinator details
// Used for query results that join agent_profiles with agent_hierarchy
type AgentWithCoordinator struct {
	Agent              AgentProfile
	CoordinatorID      *int64
	CoordinatorCode    *string
	CoordinatorName    *string
	HierarchyLevel     *int
	EffectiveFromDate  *time.Time
}
