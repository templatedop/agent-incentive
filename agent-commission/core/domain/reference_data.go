package domain

import "time"

// Circle represents a geographic circle in the postal system
// BR-IC-AH-002: Advisor Coordinators MUST be assigned to Circle and Division
type Circle struct {
	CircleID   int64  `db:"circle_id" json:"circle_id"`
	CircleCode string `db:"circle_code" json:"circle_code"` // e.g., "NORTH", "SOUTH"
	CircleName string `db:"circle_name" json:"circle_name"` // e.g., "North Circle"
	RegionCode *string `db:"region_code" json:"region_code,omitempty"` // Parent region if applicable
	IsActive   bool   `db:"is_active" json:"is_active"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// Division represents a division within a circle
// BR-IC-AH-002: Advisor Coordinators MUST be assigned to Circle and Division
type Division struct {
	DivisionID   int64  `db:"division_id" json:"division_id"`
	DivisionCode string `db:"division_code" json:"division_code"` // e.g., "DEL-001"
	DivisionName string `db:"division_name" json:"division_name"` // e.g., "Delhi Division 1"
	CircleID     int64  `db:"circle_id" json:"circle_id"`         // Foreign key to circles
	CircleCode   string `db:"circle_code" json:"circle_code"`     // Denormalized
	CircleName   string `db:"circle_name" json:"circle_name"`     // Denormalized
	HeadOffice   *string `db:"head_office" json:"head_office,omitempty"`
	IsActive     bool   `db:"is_active" json:"is_active"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// ProductPlan represents insurance product plan configurations
// Used for commission rate lookup (BR-IC-COM-006)
type ProductPlan struct {
	ProductPlanID   int64  `db:"product_plan_id" json:"product_plan_id"`
	PlanCode        string `db:"plan_code" json:"plan_code"`         // e.g., "ENDOWMENT", "WHOLE_LIFE"
	PlanName        string `db:"plan_name" json:"plan_name"`         // e.g., "Endowment Plan"
	ProductType     string `db:"product_type" json:"product_type"`   // PLI or RPLI
	PlanDescription *string `db:"plan_description" json:"plan_description,omitempty"`
	IsActive        bool   `db:"is_active" json:"is_active"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
