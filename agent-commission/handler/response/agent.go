package response

import (
	"agent-commission/core/domain"
	"time"
)

// ========================================
// Agent Onboarding Responses
// ========================================

// OnboardingSessionResponse represents an onboarding session
// POST /agents/new/init response
type OnboardingSessionResponse struct {
	SessionID       string                  `json:"session_id"`
	WorkflowState   string                  `json:"workflow_state"` // "INITIATED", "TYPE_SELECTED", "PROFILE_ENTERED", etc.
	LookupData      *LookupData            `json:"lookup_data,omitempty"`
	ValidationRules *ValidationRules       `json:"validation_rules,omitempty"`
	ExpiresAt       time.Time              `json:"expires_at"` // Session expiration
}

// LookupData provides reference data for dropdowns
type LookupData struct {
	AgentTypes     []string               `json:"agent_types"`
	PersonTypes    []string               `json:"person_types"`
	Genders        []string               `json:"genders"`
	MaritalStatuses []string              `json:"marital_statuses"`
	AddressTypes   []string               `json:"address_types"`
	ContactTypes   []string               `json:"contact_types"`
	EmailTypes     []string               `json:"email_types"`
	Circles        []CircleResponse       `json:"circles"`
	Divisions      []DivisionResponse     `json:"divisions"`
}

// ValidationRules provides validation rules for client-side validation
type ValidationRules struct {
	PANLength          int    `json:"pan_length"`           // 10
	PANPattern         string `json:"pan_pattern"`          // "^[A-Z]{5}[0-9]{4}[A-Z]{1}$"
	AadhaarLength      int    `json:"aadhaar_length"`       // 12
	PincodeLength      int    `json:"pincode_length"`       // 6
	MobileLength       int    `json:"mobile_length"`        // 10
	IFSCLength         int    `json:"ifsc_length"`          // 11
	MinAge             int    `json:"min_age"`              // 18
	MaxAge             int    `json:"max_age"`              // 65
}

// AgentTypeSelectionResponse responds to agent type selection
// POST /agents/new/{sessionId}/select-type response
type AgentTypeSelectionResponse struct {
	SessionID         string                  `json:"session_id"`
	AgentType         domain.AgentType        `json:"agent_type"`
	WorkflowState     string                  `json:"workflow_state"`
	FormConfiguration *FormConfiguration      `json:"form_configuration,omitempty"`
}

// FormConfiguration provides dynamic form configuration based on agent type
type FormConfiguration struct {
	RequiresCoordinator     bool     `json:"requires_coordinator"`      // true for ADVISOR (BR-IC-AH-001)
	RequiresCircleAssignment bool    `json:"requires_circle_assignment"` // true for COORDINATOR (BR-IC-AH-002)
	AllowsEmployeeIDLookup  bool     `json:"allows_employee_id_lookup"` // true for DEPT_EMPLOYEE (BR-IC-AH-003)
	MandatoryFields         []string `json:"mandatory_fields"`
	OptionalFields          []string `json:"optional_fields"`
	HiddenFields            []string `json:"hidden_fields"`
}

// AgentProfileValidationResponse responds to validation request
// POST /agents/validate response
type AgentProfileValidationResponse struct {
	ValidationResult string            `json:"validation_result"` // "VALID" or "INVALID"
	Errors           []ValidationError `json:"errors,omitempty"`
	Warnings         []ValidationWarning `json:"warnings,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`    // e.g., "IC-ERR-002", "IC-ERR-003"
	Message string `json:"message"`
}

// ValidationWarning represents a validation warning (non-blocking)
type ValidationWarning struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// AgentProfileCreatedResponse responds to successful agent creation
// POST /agents/new/{sessionId}/submit response
type AgentProfileCreatedResponse struct {
	AgentID        int64            `json:"agent_id"`
	AgentCode      string           `json:"agent_code"`
	FullName       string           `json:"full_name"`
	AgentType      domain.AgentType `json:"agent_type"`
	Status         domain.AgentStatus `json:"status"`
	WorkflowID     string           `json:"workflow_id"` // Temporal workflow ID
	WorkflowRunID  string           `json:"workflow_run_id"`
	Message        string           `json:"message"`
	CreatedAt      time.Time        `json:"created_at"`
}

// ========================================
// Agent Profile Management Responses
// ========================================

// AgentProfileResponse represents a complete agent profile
// GET /agents/{agentId} response
type AgentProfileResponse struct {
	AgentProfileID    int64            `json:"agent_profile_id"`
	AgentCode         string           `json:"agent_code"`
	AgentType         domain.AgentType `json:"agent_type"`
	PersonType        domain.PersonType `json:"person_type"`
	EmployeeID        *string          `json:"employee_id,omitempty"`

	// Personal Information
	Salutation    string                 `json:"salutation"`
	FirstName     string                 `json:"first_name"`
	MiddleName    *string                `json:"middle_name,omitempty"`
	LastName      string                 `json:"last_name"`
	FullName      string                 `json:"full_name"`
	Gender        *domain.Gender         `json:"gender,omitempty"`
	DateOfBirth   *time.Time             `json:"date_of_birth,omitempty"`
	Age           *int                   `json:"age,omitempty"`
	MaritalStatus *domain.MaritalStatus `json:"marital_status,omitempty"`

	// Tax & Identification
	PAN           string  `json:"pan"`
	AadhaarNumber *string `json:"aadhaar_number,omitempty"`

	// Bank Details
	BankAccountNumber *string `json:"bank_account_number,omitempty"`
	BankName          *string `json:"bank_name,omitempty"`
	BankBranch        *string `json:"bank_branch,omitempty"`
	IFSCCode          *string `json:"ifsc_code,omitempty"`
	AccountHolderName *string `json:"account_holder_name,omitempty"`
	POSBAccountNumber *string `json:"posb_account_number,omitempty"`
	POSBBranch        *string `json:"posb_branch,omitempty"`

	// Geographic Assignment
	CircleID     *int64  `json:"circle_id,omitempty"`
	CircleName   *string `json:"circle_name,omitempty"`
	DivisionID   *int64  `json:"division_id,omitempty"`
	DivisionName *string `json:"division_name,omitempty"`

	// Coordinator (if agent is an Advisor)
	CoordinatorID   *int64  `json:"coordinator_id,omitempty"`
	CoordinatorCode *string `json:"coordinator_code,omitempty"`
	CoordinatorName *string `json:"coordinator_name,omitempty"`

	// Status & Lifecycle
	Status            domain.AgentStatus `json:"status"`
	JoiningDate       time.Time          `json:"joining_date"`
	TerminationDate   *time.Time         `json:"termination_date,omitempty"`
	TerminationReason *string            `json:"termination_reason,omitempty"`

	// Related Data
	Addresses []AgentAddressResponse `json:"addresses,omitempty"`
	Contacts  []AgentContactResponse `json:"contacts,omitempty"`
	Emails    []AgentEmailResponse   `json:"emails,omitempty"`

	// Remarks
	Remarks *string `json:"remarks,omitempty"`

	// Audit
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
}

// AgentSearchResponse represents paginated search results
// POST /agents/search response
type AgentSearchResponse struct {
	Agents     []AgentSummaryResponse `json:"agents"`
	TotalCount int                    `json:"total_count"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

// AgentSummaryResponse represents a summary view of an agent (for lists)
type AgentSummaryResponse struct {
	AgentProfileID  int64            `json:"agent_profile_id"`
	AgentCode       string           `json:"agent_code"`
	FullName        string           `json:"full_name"`
	AgentType       domain.AgentType `json:"agent_type"`
	Status          domain.AgentStatus `json:"status"`
	PAN             string           `json:"pan"`
	Mobile          *string          `json:"mobile,omitempty"`
	Email           *string          `json:"email,omitempty"`
	CoordinatorName *string          `json:"coordinator_name,omitempty"`
	CircleName      *string          `json:"circle_name,omitempty"`
	DivisionName    *string          `json:"division_name,omitempty"`
	JoiningDate     time.Time        `json:"joining_date"`
}

// AgentAddressResponse represents an agent's address
type AgentAddressResponse struct {
	AgentAddressID int64              `json:"agent_address_id"`
	AddressType    domain.AddressType `json:"address_type"`
	AddressLine1   string             `json:"address_line_1"`
	AddressLine2   *string            `json:"address_line_2,omitempty"`
	AddressLine3   *string            `json:"address_line_3,omitempty"`
	Landmark       *string            `json:"landmark,omitempty"`
	City           string             `json:"city"`
	State          string             `json:"state"`
	Pincode        string             `json:"pincode"`
	Country        string             `json:"country"`
	IsPrimary      bool               `json:"is_primary"`
	FormattedAddress string           `json:"formatted_address"`
}

// AgentContactResponse represents an agent's contact
type AgentContactResponse struct {
	AgentContactID    int64              `json:"agent_contact_id"`
	ContactType       domain.ContactType `json:"contact_type"`
	ContactNumber     string             `json:"contact_number"`
	STDCode           *string            `json:"std_code,omitempty"`
	Extension         *string            `json:"extension,omitempty"`
	IsPrimary         bool               `json:"is_primary"`
	IsWhatsAppEnabled bool               `json:"is_whatsapp_enabled"`
	FormattedNumber   string             `json:"formatted_number"`
}

// AgentEmailResponse represents an agent's email
type AgentEmailResponse struct {
	AgentEmailID int64            `json:"agent_email_id"`
	EmailType    domain.EmailType `json:"email_type"`
	EmailAddress string           `json:"email_address"`
	IsPrimary    bool             `json:"is_primary"`
	IsVerified   bool             `json:"is_verified"`
	VerifiedAt   *time.Time       `json:"verified_at,omitempty"`
}

// ========================================
// Reference Data Responses
// ========================================

// CircleResponse represents a circle
type CircleResponse struct {
	CircleID   int64  `json:"circle_id"`
	CircleCode string `json:"circle_code"`
	CircleName string `json:"circle_name"`
	IsActive   bool   `json:"is_active"`
}

// DivisionResponse represents a division
type DivisionResponse struct {
	DivisionID   int64  `json:"division_id"`
	DivisionCode string `json:"division_code"`
	DivisionName string `json:"division_name"`
	CircleID     int64  `json:"circle_id"`
	CircleName   string `json:"circle_name"`
	IsActive     bool   `json:"is_active"`
}

// CoordinatorResponse represents a coordinator (for selection)
type CoordinatorResponse struct {
	CoordinatorID   int64  `json:"coordinator_id"`
	CoordinatorCode string `json:"coordinator_code"`
	CoordinatorName string `json:"coordinator_name"`
	CircleID        int64  `json:"circle_id"`
	CircleName      string `json:"circle_name"`
	DivisionID      int64  `json:"division_id"`
	DivisionName    string `json:"division_name"`
	IsActive        bool   `json:"is_active"`
}
