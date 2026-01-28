package activities

import (
	"context"
	"fmt"
	"time"

	"agent-commission/core/domain"
	"agent-commission/repo/postgres"
	"agent-commission/workflows"

	"go.temporal.io/sdk/activity"
)

// AgentActivities contains all activities for agent onboarding workflow
type AgentActivities struct {
	agentRepo     *postgres.AgentProfileRepository
	hierarchyRepo *postgres.AgentHierarchyRepository
	refDataRepo   *postgres.ReferenceDataRepository
}

// NewAgentActivities creates a new instance of AgentActivities
func NewAgentActivities(
	agentRepo *postgres.AgentProfileRepository,
	hierarchyRepo *postgres.AgentHierarchyRepository,
	refDataRepo *postgres.ReferenceDataRepository,
) *AgentActivities {
	return &AgentActivities{
		agentRepo:     agentRepo,
		hierarchyRepo: hierarchyRepo,
		refDataRepo:   refDataRepo,
	}
}

// ValidateAgentInputActivity validates the agent onboarding input
// VR-IC-PROF-003: Mandatory fields validation
// VR-IC-PROF-004: Address validation
// VR-IC-PROF-005: Contact validation
func (a *AgentActivities) ValidateAgentInputActivity(ctx context.Context, input workflows.AgentOnboardingInput) (workflows.ValidationResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Validating agent input", "AgentType", input.AgentType)

	var errors []string

	// VR-IC-PROF-003: Mandatory fields
	if input.FirstName == "" {
		errors = append(errors, "First name is required")
	}
	if input.LastName == "" {
		errors = append(errors, "Last name is required")
	}
	if input.PAN == "" {
		errors = append(errors, "PAN is required")
	}
	if len(input.PAN) != 10 {
		errors = append(errors, "PAN must be 10 characters")
	}

	// VR-IC-PROF-004: At least one address required
	if len(input.Addresses) == 0 {
		errors = append(errors, "At least one address is required")
	}

	// VR-IC-PROF-005: At least one contact required
	if len(input.Contacts) == 0 {
		errors = append(errors, "At least one contact number is required")
	}

	// Circle and Division required
	if input.CircleID == 0 {
		errors = append(errors, "Circle is required")
	}
	if input.DivisionID == 0 {
		errors = append(errors, "Division is required")
	}

	if len(errors) > 0 {
		return workflows.ValidationResult{
			IsValid:      false,
			Errors:       errors,
			ErrorMessage: "Validation failed: " + errors[0],
		}, fmt.Errorf("validation failed: %v", errors)
	}

	return workflows.ValidationResult{
		IsValid: true,
	}, nil
}

// CheckPANExistsActivity checks if PAN already exists in the system
// VR-IC-PROF-002: PAN uniqueness validation
func (a *AgentActivities) CheckPANExistsActivity(ctx context.Context, pan string, excludeAgentID int64) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Checking PAN existence", "PAN", pan)

	var excludePtr *int64
	if excludeAgentID > 0 {
		excludePtr = &excludeAgentID
	}

	exists, err := a.agentRepo.CheckPANExists(ctx, pan, excludePtr)
	if err != nil {
		logger.Error("Error checking PAN existence", "Error", err)
		return false, err
	}

	logger.Info("PAN check completed", "Exists", exists)
	return exists, nil
}

// FetchEmployeeFromHRMSActivity fetches employee data from HRMS system
// BR-IC-AH-003: Departmental employees must be fetched from HRMS
// TODO: Implement actual HRMS integration
func (a *AgentActivities) FetchEmployeeFromHRMSActivity(ctx context.Context, employeeID string) (workflows.HRMSEmployeeData, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Fetching employee from HRMS", "EmployeeID", employeeID)

	// TODO: Implement actual HRMS API call
	// For now, return mock data or error
	// In production, this would call the HRMS service

	// Mock implementation - replace with actual HRMS integration
	logger.Warn("HRMS integration not implemented, using mock data")

	// Return error if employee not found
	return workflows.HRMSEmployeeData{}, fmt.Errorf("employee not found in HRMS: IC-ERR-011")
}

// ValidateCoordinatorActivity validates that the coordinator exists and is active
// BR-IC-AH-001: Advisors must be linked to existing Advisor Coordinator
func (a *AgentActivities) ValidateCoordinatorActivity(ctx context.Context, coordinatorID int64) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Validating coordinator", "CoordinatorID", coordinatorID)

	coordinator, err := a.agentRepo.GetCoordinatorByID(ctx, coordinatorID)
	if err != nil {
		logger.Error("Coordinator validation failed", "Error", err)
		return false, err
	}

	if coordinator == nil {
		logger.Error("Coordinator not found")
		return false, fmt.Errorf("coordinator not found or inactive")
	}

	logger.Info("Coordinator validated successfully", "CoordinatorCode", coordinator.AgentCode)
	return true, nil
}

// GenerateAgentCodeActivity generates a unique agent code
// FR-IC-PROF-001: Auto-generate agent code
func (a *AgentActivities) GenerateAgentCodeActivity(ctx context.Context, agentType domain.AgentType) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Generating agent code", "AgentType", agentType)

	agentCode, err := a.agentRepo.GenerateAgentCode(ctx, agentType)
	if err != nil {
		logger.Error("Agent code generation failed", "Error", err)
		return "", err
	}

	logger.Info("Agent code generated successfully", "AgentCode", agentCode)
	return agentCode, nil
}

// CreateAgentProfileActivity creates agent profile with all related entities using batch
// FR-IC-PROF-001: Agent profile creation
// FR-IC-PROF-002: Complete agent onboarding
// Uses pgx.Batch for optimal performance
func (a *AgentActivities) CreateAgentProfileActivity(ctx context.Context, input workflows.AgentOnboardingInput) (workflows.ProfileCreationResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Creating agent profile", "AgentCode", input.AgentCode, "AgentType", input.AgentType)

	// Convert workflow input to domain models
	profile := &domain.AgentProfile{
		AgentCode:         input.AgentCode,
		AgentType:         input.AgentType,
		PersonType:        input.PersonType,
		EmployeeID:        input.EmployeeID,
		Salutation:        input.Salutation,
		FirstName:         input.FirstName,
		MiddleName:        input.MiddleName,
		LastName:          input.LastName,
		Gender:            input.Gender,
		DateOfBirth:       input.DateOfBirth,
		PAN:               input.PAN,
		BankAccountNumber: input.BankAccountNumber,
		BankName:          input.BankName,
		BankBranch:        input.BankBranch,
		IFSCCode:          input.IFSCCode,
		AccountHolderName: input.AccountHolderName,
		POSBAccountNumber: input.POSBAccountNumber,
		POSBBranch:        input.POSBBranch,
		CircleID:          input.CircleID,
		CircleName:        input.CircleName,
		DivisionID:        input.DivisionID,
		DivisionName:      input.DivisionName,
		Status:            domain.AgentStatusActive,
		JoiningDate:       input.JoiningDate,
		Remarks:           input.Remarks,
		CreatedBy:         input.RequestedBy,
	}

	// Convert addresses
	addresses := make([]domain.AgentAddress, len(input.Addresses))
	for i, addr := range input.Addresses {
		addresses[i] = domain.AgentAddress{
			AddressType:  addr.AddressType,
			AddressLine1: addr.AddressLine1,
			AddressLine2: addr.AddressLine2,
			AddressLine3: addr.AddressLine3,
			Landmark:     addr.Landmark,
			City:         addr.City,
			State:        addr.State,
			Pincode:      addr.Pincode,
			Country:      addr.Country,
			IsPrimary:    addr.IsPrimary,
		}
	}

	// Convert contacts
	contacts := make([]domain.AgentContact, len(input.Contacts))
	for i, contact := range input.Contacts {
		contacts[i] = domain.AgentContact{
			ContactType:       contact.ContactType,
			ContactNumber:     contact.ContactNumber,
			STDCode:           contact.STDCode,
			Extension:         contact.Extension,
			IsPrimary:         contact.IsPrimary,
			IsWhatsAppEnabled: contact.IsWhatsAppEnabled,
		}
	}

	// Convert emails
	emails := make([]domain.AgentEmail, len(input.Emails))
	for i, email := range input.Emails {
		emails[i] = domain.AgentEmail{
			EmailType:    email.EmailType,
			EmailAddress: email.EmailAddress,
			IsPrimary:    email.IsPrimary,
			IsVerified:   email.IsVerified,
		}
	}

	// Create profile with all relations (uses batch - single DB round trip)
	createdProfile, err := a.agentRepo.CreateAgentProfileWithRelations(ctx, profile, addresses, contacts, emails)
	if err != nil {
		logger.Error("Profile creation failed", "Error", err)
		return workflows.ProfileCreationResult{}, err
	}

	logger.Info("Agent profile created successfully", "AgentProfileID", createdProfile.AgentProfileID, "AgentCode", createdProfile.AgentCode)

	return workflows.ProfileCreationResult{
		AgentProfileID: createdProfile.AgentProfileID,
		AgentCode:      createdProfile.AgentCode,
		CreatedAt:      createdProfile.CreatedAt,
	}, nil
}

// CreateHierarchyActivity creates agent-coordinator hierarchy relationship
// BR-IC-AH-001: Advisors must be linked to existing Advisor Coordinator
func (a *AgentActivities) CreateHierarchyActivity(ctx context.Context, agentID int64, agentCode string, coordinatorID int64, createdBy string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Creating hierarchy relationship", "AgentID", agentID, "CoordinatorID", coordinatorID)

	// Get coordinator details
	coordinator, err := a.agentRepo.GetCoordinatorByID(ctx, coordinatorID)
	if err != nil {
		logger.Error("Failed to get coordinator", "Error", err)
		return err
	}

	// Create hierarchy relationship
	err = a.hierarchyRepo.CreateHierarchyRelationship(ctx, agentID, agentCode, coordinatorID, coordinator.AgentCode, time.Now(), createdBy)
	if err != nil {
		logger.Error("Hierarchy creation failed", "Error", err)
		return err
	}

	logger.Info("Hierarchy relationship created successfully")
	return nil
}

// SendOnboardingNotificationActivity sends onboarding notification to the agent
// FR-IC-PROF-003: Send welcome notification
// TODO: Implement actual email/SMS integration
func (a *AgentActivities) SendOnboardingNotificationActivity(ctx context.Context, input workflows.SendNotificationInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending onboarding notification", "AgentCode", input.AgentCode, "Email", input.Email)

	// TODO: Implement actual email/SMS service integration
	// For now, just log the notification
	logger.Info("Notification sent successfully (mock)", "AgentCode", input.AgentCode)

	return nil
}
