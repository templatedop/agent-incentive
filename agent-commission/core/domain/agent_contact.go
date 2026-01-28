package domain

import "time"

// ContactType represents the type of contact number
type ContactType string

const (
	ContactTypeOfficialLandline     ContactType = "OFFICIAL_LANDLINE"
	ContactTypeResidentialLandline  ContactType = "RESIDENTIAL_LANDLINE"
	ContactTypeMobile               ContactType = "MOBILE"
)

// AgentContact represents an agent's contact numbers
// An agent can have multiple contact numbers
// FR-IC-PROF-001: Agent profile includes contact details
type AgentContact struct {
	// Primary Key
	AgentContactID int64 `db:"agent_contact_id" json:"agent_contact_id"`

	// Foreign Key
	AgentProfileID int64 `db:"agent_profile_id" json:"agent_profile_id"`

	// Contact Type
	ContactType ContactType `db:"contact_type" json:"contact_type"`

	// Contact Details
	ContactNumber string  `db:"contact_number" json:"contact_number"` // Mandatory, 10 digits for mobile
	STDCode       *string `db:"std_code" json:"std_code,omitempty"`   // For landline numbers
	Extension     *string `db:"extension" json:"extension,omitempty"` // For office extensions

	// Flags
	IsPrimary       bool `db:"is_primary" json:"is_primary"`             // One primary contact per type
	IsWhatsAppEnabled bool `db:"is_whatsapp_enabled" json:"is_whatsapp_enabled"` // For mobile only

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
}

// FormattedNumber returns a formatted contact number
func (c *AgentContact) FormattedNumber() string {
	if c.STDCode != nil && *c.STDCode != "" {
		return *c.STDCode + "-" + c.ContactNumber
	}
	return c.ContactNumber
}

// IsMobile returns true if this is a mobile number
func (c *AgentContact) IsMobile() bool {
	return c.ContactType == ContactTypeMobile
}
