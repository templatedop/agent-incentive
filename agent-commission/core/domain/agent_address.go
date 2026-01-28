package domain

import "time"

// AddressType represents the type of address
type AddressType string

const (
	AddressTypeOfficial      AddressType = "OFFICIAL"
	AddressTypePermanent     AddressType = "PERMANENT"
	AddressTypeCommunication AddressType = "COMMUNICATION"
)

// AgentAddress represents an agent's address
// An agent can have multiple addresses (official, permanent, communication)
// FR-IC-PROF-001: Agent profile includes address details
type AgentAddress struct {
	// Primary Key
	AgentAddressID int64 `db:"agent_address_id" json:"agent_address_id"`

	// Foreign Key
	AgentProfileID int64 `db:"agent_profile_id" json:"agent_profile_id"`

	// Address Type
	AddressType AddressType `db:"address_type" json:"address_type"`

	// Address Components
	AddressLine1 string  `db:"address_line_1" json:"address_line_1"` // Mandatory
	AddressLine2 *string `db:"address_line_2" json:"address_line_2,omitempty"`
	AddressLine3 *string `db:"address_line_3" json:"address_line_3,omitempty"`
	Landmark     *string `db:"landmark" json:"landmark,omitempty"`
	City         string  `db:"city" json:"city"`         // Mandatory
	State        string  `db:"state" json:"state"`       // Mandatory
	Pincode      string  `db:"pincode" json:"pincode"`   // Mandatory, 6 digits
	Country      string  `db:"country" json:"country"`   // Default: "India"

	// Flags
	IsPrimary bool `db:"is_primary" json:"is_primary"` // One primary address per type

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
}

// FormattedAddress returns a formatted single-line address
func (a *AgentAddress) FormattedAddress() string {
	address := a.AddressLine1
	if a.AddressLine2 != nil && *a.AddressLine2 != "" {
		address += ", " + *a.AddressLine2
	}
	if a.AddressLine3 != nil && *a.AddressLine3 != "" {
		address += ", " + *a.AddressLine3
	}
	if a.Landmark != nil && *a.Landmark != "" {
		address += ", " + *a.Landmark
	}
	address += ", " + a.City + ", " + a.State + " - " + a.Pincode
	return address
}
