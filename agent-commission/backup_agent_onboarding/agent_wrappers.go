package response

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	"fmt"
	"time"
)

// ========================================
// Response Wrappers with Status
// ========================================

// AgentOnboardingResponse wraps agent onboarding result
type AgentOnboardingResponse struct {
	port.StatusCodeAndMessage
	Data       *AgentProfileResponse `json:"data"`
	WorkflowID string                `json:"workflow_id"`
}

// AgentProfileDetailResponse wraps single agent profile
type AgentProfileDetailResponse struct {
	port.StatusCodeAndMessage
	Data *AgentProfileDetailData `json:"data"`
}

// AgentProfileDetailData contains full agent profile details
type AgentProfileDetailData struct {
	AgentProfileID    int64              `json:"agent_profile_id"`
	AgentCode         string             `json:"agent_code"`
	AgentType         domain.AgentType   `json:"agent_type"`
	PersonType        string             `json:"person_type"`
	EmployeeID        string             `json:"employee_id,omitempty"`
	Salutation        string             `json:"salutation"`
	FirstName         string             `json:"first_name"`
	MiddleName        string             `json:"middle_name,omitempty"`
	LastName          string             `json:"last_name"`
	FullName          string             `json:"full_name"`
	Gender            string             `json:"gender,omitempty"`
	DateOfBirth       *time.Time         `json:"date_of_birth,omitempty"`
	PAN               string             `json:"pan"`
	AadhaarNumber     string             `json:"aadhaar_number,omitempty"`
	BankAccountNumber string             `json:"bank_account_number,omitempty"`
	BankName          string             `json:"bank_name,omitempty"`
	BankBranch        string             `json:"bank_branch,omitempty"`
	IFSCCode          string             `json:"ifsc_code,omitempty"`
	AccountHolderName string             `json:"account_holder_name,omitempty"`
	POSBAccountNumber string             `json:"posb_account_number,omitempty"`
	POSBBranch        string             `json:"posb_branch,omitempty"`
	CircleID          int64              `json:"circle_id"`
	CircleName        string             `json:"circle_name"`
	DivisionID        int64              `json:"division_id"`
	DivisionName      string             `json:"division_name"`
	Status            domain.AgentStatus `json:"status"`
	JoiningDate       time.Time          `json:"joining_date"`
	TerminationDate   *time.Time         `json:"termination_date,omitempty"`
	TerminationReason string             `json:"termination_reason,omitempty"`
	Remarks           string             `json:"remarks,omitempty"`
	Addresses         []AgentAddressResponse `json:"addresses,omitempty"`
	Contacts          []AgentContactResponse `json:"contacts,omitempty"`
	Emails            []AgentEmailResponse   `json:"emails,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         *time.Time         `json:"updated_at,omitempty"`
}

// AgentListResponse wraps agent list with pagination
type AgentListResponse struct {
	port.StatusCodeAndMessage
	port.MetaDataResponse
	Data []AgentProfileResponse `json:"data"`
}

// AgentStatusUpdateResponse wraps status update result
type AgentStatusUpdateResponse struct {
	port.StatusCodeAndMessage
	Data *AgentStatusUpdateData `json:"data"`
}

// AgentStatusUpdateData contains status update details
type AgentStatusUpdateData struct {
	AgentID   int64  `json:"agent_id"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

// CoordinatorAssignmentResponse wraps coordinator assignment result
type CoordinatorAssignmentResponse struct {
	port.StatusCodeAndMessage
	Data *CoordinatorAssignmentData `json:"data"`
}

// CoordinatorAssignmentData contains coordinator assignment details
type CoordinatorAssignmentData struct {
	AgentID         int64     `json:"agent_id"`
	CoordinatorID   int64     `json:"coordinator_id"`
	CoordinatorCode string    `json:"coordinator_code"`
	EffectiveFrom   time.Time `json:"effective_from"`
}

// ========================================
// Helper Functions
// ========================================

// ToAgentProfileResponses converts domain agents to response DTOs
func ToAgentProfileResponses(agents []domain.AgentProfile) []AgentProfileResponse {
	responses := make([]AgentProfileResponse, len(agents))
	for i, agent := range agents {
		responses[i] = ToAgentProfileResponse(agent)
	}
	return responses
}

// ToAgentProfileResponse converts a single domain agent to response DTO
func ToAgentProfileResponse(agent domain.AgentProfile) AgentProfileResponse {
	return AgentProfileResponse{
		AgentProfileID: agent.AgentProfileID,
		AgentCode:      agent.AgentCode,
		AgentType:      agent.AgentType,
		FullName:       agent.GetFullName(),
		FirstName:      agent.FirstName,
		MiddleName:     agent.MiddleName,
		PAN:            agent.PAN,
		Status:         agent.Status,
		CircleName:     agent.CircleName,
		DivisionName:   agent.DivisionName,
		JoiningDate:    agent.JoiningDate,
		CreatedAt:      agent.CreatedAt,
	}
}

// ToAddressResponses converts domain addresses to response DTOs
func ToAddressResponses(addresses []domain.AgentAddress) []AgentAddressResponse {
	responses := make([]AgentAddressResponse, len(addresses))
	for i, addr := range addresses {
		responses[i] = AgentAddressResponse{
			AgentAddressID:   addr.AgentAddressID,
			AddressType:      addr.AddressType,
			AddressLine1:     addr.AddressLine1,
			AddressLine2:     &addr.AddressLine2,
			AddressLine3:     &addr.AddressLine3,
			Landmark:         &addr.Landmark,
			City:             addr.City,
			State:            addr.State,
			Pincode:          addr.Pincode,
			Country:          addr.Country,
			IsPrimary:        addr.IsPrimary,
			FormattedAddress: formatAddress(addr),
		}
	}
	return responses
}

// ToContactResponses converts domain contacts to response DTOs
func ToContactResponses(contacts []domain.AgentContact) []AgentContactResponse {
	responses := make([]AgentContactResponse, len(contacts))
	for i, contact := range contacts {
		responses[i] = AgentContactResponse{
			AgentContactID:    contact.AgentContactID,
			ContactType:       contact.ContactType,
			ContactNumber:     contact.ContactNumber,
			STDCode:           &contact.STDCode,
			Extension:         &contact.Extension,
			IsPrimary:         contact.IsPrimary,
			IsWhatsAppEnabled: contact.IsWhatsAppEnabled,
			FormattedNumber:   formatContactNumber(contact),
		}
	}
	return responses
}

// ToEmailResponses converts domain emails to response DTOs
func ToEmailResponses(emails []domain.AgentEmail) []AgentEmailResponse {
	responses := make([]AgentEmailResponse, len(emails))
	for i, email := range emails {
		responses[i] = AgentEmailResponse{
			AgentEmailID: email.AgentEmailID,
			EmailType:    email.EmailType,
			EmailAddress: email.EmailAddress,
			IsPrimary:    email.IsPrimary,
			IsVerified:   email.IsVerified,
			VerifiedAt:   email.VerifiedAt,
		}
	}
	return responses
}

// formatAddress formats address for display
func formatAddress(addr domain.AgentAddress) string {
	parts := []string{addr.AddressLine1}
	if addr.AddressLine2 != "" {
		parts = append(parts, addr.AddressLine2)
	}
	if addr.AddressLine3 != "" {
		parts = append(parts, addr.AddressLine3)
	}
	if addr.Landmark != "" {
		parts = append(parts, addr.Landmark)
	}
	parts = append(parts, addr.City, addr.State, addr.Pincode, addr.Country)

	formatted := ""
	for i, part := range parts {
		if i > 0 {
			formatted += ", "
		}
		formatted += part
	}
	return formatted
}

// formatContactNumber formats contact number for display
func formatContactNumber(contact domain.AgentContact) string {
	formatted := contact.ContactNumber
	if contact.STDCode != "" {
		formatted = fmt.Sprintf("+%s %s", contact.STDCode, contact.ContactNumber)
	}
	if contact.Extension != "" {
		formatted = fmt.Sprintf("%s ext %s", formatted, contact.Extension)
	}
	return formatted
}
