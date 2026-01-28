package domain

import "time"

// EmailType represents the type of email address
type EmailType string

const (
	EmailTypeOfficial      EmailType = "OFFICIAL"
	EmailTypePermanent     EmailType = "PERMANENT"
	EmailTypeCommunication EmailType = "COMMUNICATION"
)

// AgentEmail represents an agent's email addresses
// An agent can have multiple email addresses
// FR-IC-PROF-001: Agent profile includes email details
type AgentEmail struct {
	// Primary Key
	AgentEmailID int64 `db:"agent_email_id" json:"agent_email_id"`

	// Foreign Key
	AgentProfileID int64 `db:"agent_profile_id" json:"agent_profile_id"`

	// Email Type
	EmailType EmailType `db:"email_type" json:"email_type"`

	// Email Details
	EmailAddress string `db:"email_address" json:"email_address"` // Mandatory, validated format

	// Flags
	IsPrimary  bool `db:"is_primary" json:"is_primary"`   // One primary email per type
	IsVerified bool `db:"is_verified" json:"is_verified"` // Email verification status

	// Verification
	VerifiedAt *time.Time `db:"verified_at" json:"verified_at,omitempty"`

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
}
