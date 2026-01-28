package workflows

import (
	"time"

	"agent-commission/core/domain"
)

// AgentOnboardingInput represents the input for agent onboarding workflow
type AgentOnboardingInput struct {
	// Agent Profile Fields
	AgentCode   string           `json:"agent_code,omitempty"`
	AgentType   domain.AgentType `json:"agent_type"`
	PersonType  string           `json:"person_type"`
	EmployeeID  string           `json:"employee_id,omitempty"`
	Salutation  string           `json:"salutation,omitempty"`
	FirstName   string           `json:"first_name"`
	MiddleName  string           `json:"middle_name,omitempty"`
	LastName    string           `json:"last_name"`
	Gender      string           `json:"gender,omitempty"`
	DateOfBirth *time.Time       `json:"date_of_birth,omitempty"`
	PAN         string           `json:"pan"`

	// Hierarchy
	CoordinatorID int64 `json:"coordinator_id,omitempty"`

	// Circle/Division
	CircleID     int64  `json:"circle_id"`
	CircleName   string `json:"circle_name"`
	DivisionID   int64  `json:"division_id"`
	DivisionName string `json:"division_name"`

	// Contact Information
	PrimaryEmail   string `json:"primary_email,omitempty"`
	PrimaryContact string `json:"primary_contact,omitempty"`

	// Addresses
	Addresses []AddressInput `json:"addresses"`

	// Contacts
	Contacts []ContactInput `json:"contacts"`

	// Emails
	Emails []EmailInput `json:"emails"`

	// Bank Details
	BankAccountNumber  string `json:"bank_account_number,omitempty"`
	BankName           string `json:"bank_name,omitempty"`
	BankBranch         string `json:"bank_branch,omitempty"`
	IFSCCode           string `json:"ifsc_code,omitempty"`
	AccountHolderName  string `json:"account_holder_name,omitempty"`
	POSBAccountNumber  string `json:"posb_account_number,omitempty"`
	POSBBranch         string `json:"posb_branch,omitempty"`

	// Metadata
	AgentProfileID int64     `json:"agent_profile_id,omitempty"` // For updates
	JoiningDate    time.Time `json:"joining_date"`
	Remarks        string    `json:"remarks,omitempty"`
	RequestedBy    string    `json:"requested_by"`
}

// AddressInput represents address input for workflow
type AddressInput struct {
	AddressType   string `json:"address_type"`
	AddressLine1  string `json:"address_line_1"`
	AddressLine2  string `json:"address_line_2,omitempty"`
	AddressLine3  string `json:"address_line_3,omitempty"`
	Landmark      string `json:"landmark,omitempty"`
	City          string `json:"city"`
	State         string `json:"state"`
	Pincode       string `json:"pincode"`
	Country       string `json:"country"`
	IsPrimary     bool   `json:"is_primary"`
}

// ContactInput represents contact input for workflow
type ContactInput struct {
	ContactType        string `json:"contact_type"`
	ContactNumber      string `json:"contact_number"`
	STDCode            string `json:"std_code,omitempty"`
	Extension          string `json:"extension,omitempty"`
	IsPrimary          bool   `json:"is_primary"`
	IsWhatsAppEnabled  bool   `json:"is_whatsapp_enabled"`
}

// EmailInput represents email input for workflow
type EmailInput struct {
	EmailType    string `json:"email_type"`
	EmailAddress string `json:"email_address"`
	IsPrimary    bool   `json:"is_primary"`
	IsVerified   bool   `json:"is_verified"`
}

// AgentOnboardingResult represents the output of agent onboarding workflow
type AgentOnboardingResult struct {
	WorkflowID     string    `json:"workflow_id"`
	Status         string    `json:"status"` // SUCCESS, FAILED, PARTIAL_SUCCESS
	AgentProfileID int64     `json:"agent_profile_id,omitempty"`
	AgentCode      string    `json:"agent_code,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CompletedAt    time.Time `json:"completed_at,omitempty"`
}

// ValidationResult represents the result of input validation
type ValidationResult struct {
	IsValid      bool     `json:"is_valid"`
	Errors       []string `json:"errors,omitempty"`
	ErrorMessage string   `json:"error_message,omitempty"`
}

// HRMSEmployeeData represents employee data fetched from HRMS
type HRMSEmployeeData struct {
	EmployeeID    string     `json:"employee_id"`
	FirstName     string     `json:"first_name"`
	MiddleName    string     `json:"middle_name"`
	LastName      string     `json:"last_name"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Gender        string     `json:"gender"`
	Email         string     `json:"email"`
	ContactNumber string     `json:"contact_number"`
	CircleID      int64      `json:"circle_id"`
	CircleName    string     `json:"circle_name"`
	DivisionID    int64      `json:"division_id"`
	DivisionName  string     `json:"division_name"`
	Designation   string     `json:"designation"`
	JoiningDate   time.Time  `json:"joining_date"`
}

// ProfileCreationResult represents the result of profile creation
type ProfileCreationResult struct {
	AgentProfileID int64     `json:"agent_profile_id"`
	AgentCode      string    `json:"agent_code"`
	CreatedAt      time.Time `json:"created_at"`
}

// SendNotificationInput represents input for sending onboarding notification
type SendNotificationInput struct {
	AgentProfileID int64            `json:"agent_profile_id"`
	AgentCode      string           `json:"agent_code"`
	Email          string           `json:"email"`
	AgentType      domain.AgentType `json:"agent_type"`
}
