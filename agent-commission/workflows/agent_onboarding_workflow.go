package workflows

import (
	"time"

	"agent-commission/core/domain"

	"go.temporal.io/sdk/workflow"
)

// AgentOnboardingWorkflow orchestrates the agent onboarding process
// WF-IC-ONB-001: New Agent Onboarding Process
// FR-IC-PROF-001: Agent profile creation
// FR-IC-PROF-002: Complete agent onboarding
// Handles different agent types: Advisor, Advisor Coordinator, Departmental Employee, Field Officer
func AgentOnboardingWorkflow(ctx workflow.Context, input AgentOnboardingInput) (*AgentOnboardingResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting Agent Onboarding Workflow", "AgentType", input.AgentType, "RequestedBy", input.RequestedBy)

	// Workflow options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &workflow.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result AgentOnboardingResult
	result.WorkflowID = workflow.GetInfo(ctx).WorkflowExecution.RunID

	// Step 1: Validate basic input
	logger.Info("Step 1: Validating basic input")
	var validationResult ValidationResult
	err := workflow.ExecuteActivity(ctx, ValidateAgentInputActivity, input).Get(ctx, &validationResult)
	if err != nil {
		logger.Error("Input validation failed", "Error", err)
		result.Status = "FAILED"
		result.ErrorMessage = err.Error()
		return &result, err
	}

	// Step 2: Check PAN uniqueness (VR-IC-PROF-002)
	if input.PAN != "" {
		logger.Info("Step 2: Checking PAN uniqueness", "PAN", input.PAN)
		var panExists bool
		err = workflow.ExecuteActivity(ctx, CheckPANExistsActivity, input.PAN, input.AgentProfileID).Get(ctx, &panExists)
		if err != nil {
			logger.Error("PAN check failed", "Error", err)
			result.Status = "FAILED"
			result.ErrorMessage = err.Error()
			return &result, err
		}
		if panExists {
			logger.Error("PAN already exists", "PAN", input.PAN)
			result.Status = "FAILED"
			result.ErrorMessage = "PAN already exists in the system"
			return &result, workflow.NewApplicationError("PAN already exists", "IC-ERR-002", nil)
		}
	}

	// Step 3: Handle Departmental Employee HRMS integration (BR-IC-AH-003)
	if input.AgentType == domain.AgentTypeDepartmentalEmployee && input.EmployeeID != "" {
		logger.Info("Step 3: Fetching employee data from HRMS", "EmployeeID", input.EmployeeID)
		var hrmsData HRMSEmployeeData
		err = workflow.ExecuteActivity(ctx, FetchEmployeeFromHRMSActivity, input.EmployeeID).Get(ctx, &hrmsData)
		if err != nil {
			logger.Error("HRMS fetch failed", "Error", err)
			result.Status = "FAILED"
			result.ErrorMessage = "Employee not found in HRMS (IC-ERR-011)"
			return &result, err
		}

		// Merge HRMS data with input
		input = mergeHRMSData(input, hrmsData)
		logger.Info("HRMS data merged successfully")
	}

	// Step 4: Validate Coordinator assignment for Advisors (BR-IC-AH-001)
	if input.AgentType == domain.AgentTypeAdvisor {
		if input.CoordinatorID == 0 {
			logger.Error("Advisor must be assigned to a coordinator")
			result.Status = "FAILED"
			result.ErrorMessage = "Advisor must be assigned to a coordinator (BR-IC-AH-001)"
			return &result, workflow.NewApplicationError("Missing coordinator", "IC-ERR-003", nil)
		}

		logger.Info("Step 4: Validating coordinator", "CoordinatorID", input.CoordinatorID)
		var coordinatorValid bool
		err = workflow.ExecuteActivity(ctx, ValidateCoordinatorActivity, input.CoordinatorID).Get(ctx, &coordinatorValid)
		if err != nil || !coordinatorValid {
			logger.Error("Invalid coordinator", "Error", err)
			result.Status = "FAILED"
			result.ErrorMessage = "Invalid or inactive coordinator"
			return &result, err
		}
	}

	// Step 5: Generate Agent Code
	logger.Info("Step 5: Generating agent code", "AgentType", input.AgentType)
	var agentCode string
	err = workflow.ExecuteActivity(ctx, GenerateAgentCodeActivity, input.AgentType).Get(ctx, &agentCode)
	if err != nil {
		logger.Error("Agent code generation failed", "Error", err)
		result.Status = "FAILED"
		result.ErrorMessage = err.Error()
		return &result, err
	}
	input.AgentCode = agentCode
	logger.Info("Agent code generated", "AgentCode", agentCode)

	// Step 6: Create Agent Profile with all relations (using batch)
	logger.Info("Step 6: Creating agent profile with relations")
	var profileCreationResult ProfileCreationResult
	err = workflow.ExecuteActivity(ctx, CreateAgentProfileActivity, input).Get(ctx, &profileCreationResult)
	if err != nil {
		logger.Error("Profile creation failed", "Error", err)
		result.Status = "FAILED"
		result.ErrorMessage = err.Error()
		return &result, err
	}

	// Step 7: Create hierarchy relationship if Advisor
	if input.AgentType == domain.AgentTypeAdvisor && input.CoordinatorID > 0 {
		logger.Info("Step 7: Creating hierarchy relationship", "AgentID", profileCreationResult.AgentProfileID, "CoordinatorID", input.CoordinatorID)
		err = workflow.ExecuteActivity(ctx, CreateHierarchyActivity, profileCreationResult.AgentProfileID, agentCode, input.CoordinatorID, input.RequestedBy).Get(ctx, nil)
		if err != nil {
			logger.Error("Hierarchy creation failed", "Error", err)
			// Note: Profile is already created, but we should mark this as partial success
			result.Status = "PARTIAL_SUCCESS"
			result.ErrorMessage = "Profile created but hierarchy assignment failed"
			result.AgentProfileID = profileCreationResult.AgentProfileID
			result.AgentCode = agentCode
			return &result, nil
		}
	}

	// Step 8: Send notification (optional, fire-and-forget)
	_ = workflow.ExecuteActivity(ctx, SendOnboardingNotificationActivity, SendNotificationInput{
		AgentProfileID: profileCreationResult.AgentProfileID,
		AgentCode:      agentCode,
		Email:          input.PrimaryEmail,
		AgentType:      input.AgentType,
	}).Get(ctx, nil)

	// Success
	logger.Info("Agent onboarding completed successfully", "AgentCode", agentCode, "AgentProfileID", profileCreationResult.AgentProfileID)
	result.Status = "SUCCESS"
	result.AgentProfileID = profileCreationResult.AgentProfileID
	result.AgentCode = agentCode
	result.CompletedAt = workflow.Now(ctx)

	return &result, nil
}

// mergeHRMSData merges HRMS employee data with onboarding input
func mergeHRMSData(input AgentOnboardingInput, hrmsData HRMSEmployeeData) AgentOnboardingInput {
	if hrmsData.FirstName != "" {
		input.FirstName = hrmsData.FirstName
	}
	if hrmsData.MiddleName != "" {
		input.MiddleName = hrmsData.MiddleName
	}
	if hrmsData.LastName != "" {
		input.LastName = hrmsData.LastName
	}
	if hrmsData.DateOfBirth != nil {
		input.DateOfBirth = hrmsData.DateOfBirth
	}
	if hrmsData.Gender != "" {
		input.Gender = hrmsData.Gender
	}
	if hrmsData.Email != "" {
		input.PrimaryEmail = hrmsData.Email
	}
	if hrmsData.ContactNumber != "" {
		input.PrimaryContact = hrmsData.ContactNumber
	}
	if hrmsData.CircleID > 0 {
		input.CircleID = hrmsData.CircleID
		input.CircleName = hrmsData.CircleName
	}
	if hrmsData.DivisionID > 0 {
		input.DivisionID = hrmsData.DivisionID
		input.DivisionName = hrmsData.DivisionName
	}
	return input
}
