package domain

import (
	"time"
)

// AgentType represents the type of agent
// BR-IC-AH-001: Advisors must link to Coordinator
// BR-IC-AH-002: Coordinators must assign to Circle/Division
// BR-IC-AH-003: Departmental Employee from HRMS
// BR-IC-AH-004: Field Officer onboarding modes
type AgentType string

const (
	AgentTypeAdvisor              AgentType = "ADVISOR"
	AgentTypeAdvisorCoordinator   AgentType = "ADVISOR_COORDINATOR"
	AgentTypeDepartmentalEmployee AgentType = "DEPARTMENTAL_EMPLOYEE"
	AgentTypeFieldOfficer         AgentType = "FIELD_OFFICER"
)

// PersonType represents whether agent is individual or corporate
type PersonType string

const (
	PersonTypeIndividual     PersonType = "INDIVIDUAL"
	PersonTypeCorporateGroup PersonType = "CORPORATE_GROUP"
)

// Gender enumeration
type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

// MaritalStatus enumeration
type MaritalStatus string

const (
	MaritalStatusSingle   MaritalStatus = "SINGLE"
	MaritalStatusMarried  MaritalStatus = "MARRIED"
	MaritalStatusDivorced MaritalStatus = "DIVORCED"
	MaritalStatusWidowed  MaritalStatus = "WIDOWED"
)

// AgentStatus represents the current status of an agent
// BR-IC-PROF-004: Status lifecycle management
type AgentStatus string

const (
	AgentStatusActive      AgentStatus = "ACTIVE"
	AgentStatusSuspended   AgentStatus = "SUSPENDED"
	AgentStatusTerminated  AgentStatus = "TERMINATED"
	AgentStatusExpired     AgentStatus = "EXPIRED"
	AgentStatusDeactivated AgentStatus = "DEACTIVATED"
)

// AgentProfile represents the core agent entity
// FR-IC-PROF-001: Agent onboarding
// FR-IC-PROF-002: Agent profile creation
// VR-IC-PROF-001 to VR-IC-PROF-005: Validation rules
type AgentProfile struct {
	// Primary Key
	AgentProfileID int64 `db:"agent_profile_id" json:"agent_profile_id"`

	// Agent Identification
	AgentCode     string    `db:"agent_code" json:"agent_code"`               // Auto-generated unique code
	AgentType     AgentType `db:"agent_type" json:"agent_type"`               // BR-IC-AH-001 to BR-IC-AH-004
	PersonType    PersonType `db:"person_type" json:"person_type"`            // Individual or Corporate
	EmployeeID    *string   `db:"employee_id" json:"employee_id,omitempty"`   // For Dept Employee/Field Officer (BR-IC-AH-003)

	// Personal Information
	Salutation    string        `db:"salutation" json:"salutation"`
	FirstName     string        `db:"first_name" json:"first_name"`             // VR-IC-PROF-003: Mandatory
	MiddleName    *string       `db:"middle_name" json:"middle_name,omitempty"`
	LastName      string        `db:"last_name" json:"last_name"`               // VR-IC-PROF-003: Mandatory
	Gender        *Gender       `db:"gender" json:"gender,omitempty"`
	DateOfBirth   *time.Time    `db:"date_of_birth" json:"date_of_birth,omitempty"` // VR-IC-PROF-004: DOB validation
	MaritalStatus *MaritalStatus `db:"marital_status" json:"marital_status,omitempty"`

	// Tax & Identification
	PAN           string    `db:"pan" json:"pan"`                               // VR-IC-PROF-002: Unique, 10 chars
	AadhaarNumber *string   `db:"aadhaar_number" json:"aadhaar_number,omitempty"`

	// Bank Details
	BankAccountNumber *string `db:"bank_account_number" json:"bank_account_number,omitempty"`
	BankName          *string `db:"bank_name" json:"bank_name,omitempty"`
	BankBranch        *string `db:"bank_branch" json:"bank_branch,omitempty"`
	IFSCCode          *string `db:"ifsc_code" json:"ifsc_code,omitempty"`
	AccountHolderName *string `db:"account_holder_name" json:"account_holder_name,omitempty"`

	// POSB Account (Alternative)
	POSBAccountNumber *string `db:"posb_account_number" json:"posb_account_number,omitempty"`
	POSBBranch        *string `db:"posb_branch" json:"posb_branch,omitempty"`

	// Geographic Assignment (for Coordinators - BR-IC-AH-002)
	CircleID    *int64  `db:"circle_id" json:"circle_id,omitempty"`
	CircleName  *string `db:"circle_name" json:"circle_name,omitempty"`       // Denormalized for query performance
	DivisionID  *int64  `db:"division_id" json:"division_id,omitempty"`
	DivisionName *string `db:"division_name" json:"division_name,omitempty"`  // Denormalized

	// Status & Lifecycle
	Status            AgentStatus `db:"status" json:"status"`                   // BR-IC-PROF-004
	JoiningDate       time.Time   `db:"joining_date" json:"joining_date"`
	TerminationDate   *time.Time  `db:"termination_date" json:"termination_date,omitempty"`
	TerminationReason *string     `db:"termination_reason" json:"termination_reason,omitempty"`

	// Remarks & Notes
	Remarks *string `db:"remarks" json:"remarks,omitempty"`

	// Audit Fields
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Soft delete
}

// IsAdvisor returns true if agent is an Advisor
func (a *AgentProfile) IsAdvisor() bool {
	return a.AgentType == AgentTypeAdvisor
}

// IsCoordinator returns true if agent is an Advisor Coordinator
func (a *AgentProfile) IsCoordinator() bool {
	return a.AgentType == AgentTypeAdvisorCoordinator
}

// IsDepartmentalEmployee returns true if agent is a Departmental Employee
func (a *AgentProfile) IsDepartmentalEmployee() bool {
	return a.AgentType == AgentTypeDepartmentalEmployee
}

// IsFieldOfficer returns true if agent is a Field Officer
func (a *AgentProfile) IsFieldOfficer() bool {
	return a.AgentType == AgentTypeFieldOfficer
}

// IsActive returns true if agent status is ACTIVE
func (a *AgentProfile) IsActive() bool {
	return a.Status == AgentStatusActive
}

// RequiresCoordinator returns true if this agent type requires coordinator assignment
// BR-IC-AH-001: Advisors MUST be linked to existing Advisor Coordinator
func (a *AgentProfile) RequiresCoordinator() bool {
	return a.AgentType == AgentTypeAdvisor
}

// RequiresCircleAssignment returns true if this agent type requires circle/division assignment
// BR-IC-AH-002: Advisor Coordinators MUST be assigned to Circle and Division
func (a *AgentProfile) RequiresCircleAssignment() bool {
	return a.AgentType == AgentTypeAdvisorCoordinator
}

// FullName returns the complete name of the agent
func (a *AgentProfile) FullName() string {
	name := a.FirstName
	if a.MiddleName != nil && *a.MiddleName != "" {
		name += " " + *a.MiddleName
	}
	name += " " + a.LastName
	return name
}
