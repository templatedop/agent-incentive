package handler

import (
	"agent-commission/core/domain"
	"time"
)

// ========================================
// Agent Onboarding Requests
// ========================================

// InitAgentOnboardingRequest initializes a new onboarding session
// POST /agents/new/init
// FR-IC-PROF-001: Agent onboarding initialization
// WF-IC-ONB-001: Step 1
type InitAgentOnboardingRequest struct {
	AgentType *domain.AgentType `json:"agent_type,omitempty"` // Optional pre-selection
}

// SelectAgentTypeRequest selects the agent type for onboarding
// POST /agents/new/{sessionId}/select-type
// FR-IC-PROF-001: Agent type selection
// BR-IC-AH-001: Advisor requires Coordinator
// VR-IC-PROF-001: Agent type validation
// WF-IC-ONB-001: Step 2
type SelectAgentTypeRequest struct {
	AgentType domain.AgentType `json:"agent_type" validate:"required,oneof=ADVISOR ADVISOR_COORDINATOR DEPARTMENTAL_EMPLOYEE FIELD_OFFICER"`
}

// CreateAgentProfileRequest creates a new agent profile
// POST /agents/new/{sessionId}/profile-details
// FR-IC-PROF-001, FR-IC-PROF-002: Agent profile creation
// VR-IC-PROF-001 to VR-IC-PROF-005: Validation rules
// WF-IC-ONB-001: Step 3
type CreateAgentProfileRequest struct {
	// Agent Type
	AgentType  domain.AgentType  `json:"agent_type" validate:"required"`
	PersonType domain.PersonType `json:"person_type" validate:"required,oneof=INDIVIDUAL CORPORATE_GROUP"`

	// Employee ID (for Dept Employee/Field Officer - BR-IC-AH-003, BR-IC-AH-004)
	EmployeeID *string `json:"employee_id,omitempty"`

	// Personal Information
	Salutation    string                 `json:"salutation" validate:"required,oneof=Mr Mrs Ms Dr"`
	FirstName     string                 `json:"first_name" validate:"required,min=2,max=100"` // VR-IC-PROF-003
	MiddleName    *string                `json:"middle_name,omitempty" validate:"omitempty,max=100"`
	LastName      string                 `json:"last_name" validate:"required,min=2,max=100"`   // VR-IC-PROF-003
	Gender        *domain.Gender         `json:"gender,omitempty" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	DateOfBirth   *time.Time             `json:"date_of_birth,omitempty"` // VR-IC-PROF-004
	MaritalStatus *domain.MaritalStatus `json:"marital_status,omitempty" validate:"omitempty,oneof=SINGLE MARRIED DIVORCED WIDOWED"`

	// Tax & Identification
	PAN           string  `json:"pan" validate:"required,len=10,alphanum"` // VR-IC-PROF-002: Unique, 10 chars
	AadhaarNumber *string `json:"aadhaar_number,omitempty" validate:"omitempty,len=12,numeric"`

	// Bank Details
	BankAccountNumber *string `json:"bank_account_number,omitempty" validate:"omitempty,min=9,max=18"`
	BankName          *string `json:"bank_name,omitempty" validate:"omitempty,max=200"`
	BankBranch        *string `json:"bank_branch,omitempty" validate:"omitempty,max=200"`
	IFSCCode          *string `json:"ifsc_code,omitempty" validate:"omitempty,len=11,alphanum"`
	AccountHolderName *string `json:"account_holder_name,omitempty" validate:"omitempty,max=200"`

	// POSB Account (Alternative)
	POSBAccountNumber *string `json:"posb_account_number,omitempty" validate:"omitempty,max=50"`
	POSBBranch        *string `json:"posb_branch,omitempty" validate:"omitempty,max=200"`

	// Geographic Assignment (for Coordinators - BR-IC-AH-002)
	CircleID   *int64 `json:"circle_id,omitempty"`   // Mandatory for ADVISOR_COORDINATOR
	DivisionID *int64 `json:"division_id,omitempty"` // Mandatory for ADVISOR_COORDINATOR

	// Addresses
	Addresses []CreateAddressRequest `json:"addresses" validate:"required,min=1,dive"`

	// Contacts
	Contacts []CreateContactRequest `json:"contacts" validate:"required,min=1,dive"`

	// Emails
	Emails []CreateEmailRequest `json:"emails" validate:"required,min=1,dive"`

	// Remarks
	Remarks *string `json:"remarks,omitempty"`
}

// CreateAddressRequest represents an address to be created
type CreateAddressRequest struct {
	AddressType  domain.AddressType `json:"address_type" validate:"required,oneof=OFFICIAL PERMANENT COMMUNICATION"`
	AddressLine1 string             `json:"address_line_1" validate:"required,max=200"`
	AddressLine2 *string            `json:"address_line_2,omitempty" validate:"omitempty,max=200"`
	AddressLine3 *string            `json:"address_line_3,omitempty" validate:"omitempty,max=200"`
	Landmark     *string            `json:"landmark,omitempty" validate:"omitempty,max=200"`
	City         string             `json:"city" validate:"required,max=100"`
	State        string             `json:"state" validate:"required,max=100"`
	Pincode      string             `json:"pincode" validate:"required,len=6,numeric"`
	Country      string             `json:"country" validate:"required,max=100"`
	IsPrimary    bool               `json:"is_primary"`
}

// CreateContactRequest represents a contact number to be created
type CreateContactRequest struct {
	ContactType       domain.ContactType `json:"contact_type" validate:"required,oneof=OFFICIAL_LANDLINE RESIDENTIAL_LANDLINE MOBILE"`
	ContactNumber     string             `json:"contact_number" validate:"required"`
	STDCode           *string            `json:"std_code,omitempty" validate:"omitempty,max=10"`
	Extension         *string            `json:"extension,omitempty" validate:"omitempty,max=10"`
	IsPrimary         bool               `json:"is_primary"`
	IsWhatsAppEnabled bool               `json:"is_whatsapp_enabled"`
}

// CreateEmailRequest represents an email address to be created
type CreateEmailRequest struct {
	EmailType    domain.EmailType `json:"email_type" validate:"required,oneof=OFFICIAL PERMANENT COMMUNICATION"`
	EmailAddress string           `json:"email_address" validate:"required,email"`
	IsPrimary    bool             `json:"is_primary"`
}

// SelectCoordinatorRequest assigns a coordinator to an advisor
// POST /agents/new/{sessionId}/select-coordinator
// BR-IC-AH-001: Advisor must link to Coordinator
// WF-IC-ONB-001: Step 4
type SelectCoordinatorRequest struct {
	CoordinatorID int64 `json:"coordinator_id" validate:"required,gt=0"`
}

// SubmitAgentOnboardingRequest finalizes and submits the onboarding
// POST /agents/new/{sessionId}/submit
// FR-IC-PROF-001: Final submission
// WF-IC-ONB-001: Step 6 - Start Temporal workflow
type SubmitAgentOnboardingRequest struct {
	// Additional submission metadata
	SubmittedBy string  `json:"submitted_by" validate:"required"`
	Comments    *string `json:"comments,omitempty"`
}

// AgentProfileValidationRequest validates agent profile data
// POST /agents/validate
// VR-IC-PROF-002 to VR-IC-PROF-007: All validation rules
// WF-IC-ONB-001: Step 5
type AgentProfileValidationRequest struct {
	PAN           string              `json:"pan" validate:"required,len=10,alphanum"`
	AgentType     domain.AgentType    `json:"agent_type" validate:"required"`
	CoordinatorID *int64              `json:"coordinator_id,omitempty"` // Required if AgentType=ADVISOR
	CircleID      *int64              `json:"circle_id,omitempty"`      // Required if AgentType=COORDINATOR
	DivisionID    *int64              `json:"division_id,omitempty"`    // Required if AgentType=COORDINATOR
	EmployeeID    *string             `json:"employee_id,omitempty"`    // Required if AgentType=DEPT_EMPLOYEE
}

// ========================================
// Agent Profile Management Requests (Module 1.2)
// ========================================

// AgentSearchRequest searches for agents
// POST /agents/search
// FR-IC-PROF-004: Agent search
type AgentSearchRequest struct {
	AgentID    *string            `json:"agent_id,omitempty"`
	AgentCode  *string            `json:"agent_code,omitempty"`
	Name       *string            `json:"name,omitempty"`
	PAN        *string            `json:"pan,omitempty"`
	Mobile     *string            `json:"mobile,omitempty"`
	AgentType  *domain.AgentType  `json:"agent_type,omitempty"`
	Status     *domain.AgentStatus `json:"status,omitempty"`
	CircleID   *int64             `json:"circle_id,omitempty"`
	DivisionID *int64             `json:"division_id,omitempty"`

	// Pagination
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// UpdateAgentProfileRequest updates an existing agent profile
// PUT /agents/{agentId}
// FR-IC-PROF-006: Agent profile update
// BR-IC-PROF-003: Audit trail required
type UpdateAgentProfileRequest struct {
	// Allow updating specific fields only
	FirstName     *string                `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	MiddleName    *string                `json:"middle_name,omitempty" validate:"omitempty,max=100"`
	LastName      *string                `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Gender        *domain.Gender         `json:"gender,omitempty"`
	DateOfBirth   *time.Time             `json:"date_of_birth,omitempty"`
	MaritalStatus *domain.MaritalStatus `json:"marital_status,omitempty"`
	BankAccountNumber *string            `json:"bank_account_number,omitempty"`
	BankName          *string            `json:"bank_name,omitempty"`
	BankBranch        *string            `json:"bank_branch,omitempty"`
	IFSCCode          *string            `json:"ifsc_code,omitempty"`
	AccountHolderName *string            `json:"account_holder_name,omitempty"`
	POSBAccountNumber *string            `json:"posb_account_number,omitempty"`
	POSBBranch        *string            `json:"posb_branch,omitempty"`
	Remarks           *string            `json:"remarks,omitempty"`

	// Audit
	UpdatedBy string `json:"updated_by" validate:"required"`
	Reason    string `json:"reason" validate:"required,max=500"` // BR-IC-PROF-003: Audit trail
}

// UpdateAgentStatusRequest updates agent status
// PUT /agents/{agentId}/status
// BR-IC-PROF-004: Status management
type UpdateAgentStatusRequest struct {
	Status domain.AgentStatus `json:"status" validate:"required,oneof=ACTIVE SUSPENDED TERMINATED EXPIRED DEACTIVATED"`
	Reason string             `json:"reason" validate:"required,max=500"`
	UpdatedBy string          `json:"updated_by" validate:"required"`
}
