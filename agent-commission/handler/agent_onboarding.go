package handler

import (
	"agent-commission/core/domain"
	"agent-commission/core/port"
	resp "agent-commission/handler/response"
	"agent-commission/repo/postgres"
	"agent-commission/workflows"

	"go.temporal.io/sdk/client"

	log "gitlab.cept.gov.in/it-2.0-common/n-api-log"
	serverHandler "gitlab.cept.gov.in/it-2.0-common/n-api-server/handler"
	serverRoute "gitlab.cept.gov.in/it-2.0-common/n-api-server/route"
)

// AgentOnboardingHandler handles agent onboarding API endpoints
type AgentOnboardingHandler struct {
	*serverHandler.Base
	agentRepo     *postgres.AgentProfileRepository
	hierarchyRepo *postgres.AgentHierarchyRepository
	refDataRepo   *postgres.ReferenceDataRepository
	temporalClient client.Client
}

// NewAgentOnboardingHandler creates a new agent onboarding handler
func NewAgentOnboardingHandler(
	agentRepo *postgres.AgentProfileRepository,
	hierarchyRepo *postgres.AgentHierarchyRepository,
	refDataRepo *postgres.ReferenceDataRepository,
	temporalClient client.Client,
) *AgentOnboardingHandler {
	base := serverHandler.New("Agent Onboarding").SetPrefix("/v1/agents").AddPrefix("")
	return &AgentOnboardingHandler{
		Base:           base,
		agentRepo:      agentRepo,
		hierarchyRepo:  hierarchyRepo,
		refDataRepo:    refDataRepo,
		temporalClient: temporalClient,
	}
}

// Routes defines all routes for agent onboarding
func (h *AgentOnboardingHandler) Routes() []serverRoute.Route {
	return []serverRoute.Route{
		serverRoute.POST("/onboard", h.OnboardAgent).Name("Onboard New Agent"),
		serverRoute.GET("/:id", h.GetAgentProfile).Name("Get Agent Profile"),
		serverRoute.GET("", h.SearchAgents).Name("Search Agents"),
		serverRoute.PUT("/:id", h.UpdateAgentProfile).Name("Update Agent Profile"),
		serverRoute.PUT("/:id/status", h.UpdateAgentStatus).Name("Update Agent Status"),
		serverRoute.PUT("/:id/coordinator", h.AssignCoordinator).Name("Assign/Change Coordinator"),
	}
}

// OnboardAgent initiates agent onboarding workflow
// FR-IC-PROF-001: Agent profile creation
// FR-IC-PROF-002: Complete agent onboarding
// WF-IC-ONB-001: Agent Onboarding Workflow
func (h *AgentOnboardingHandler) OnboardAgent(sctx *serverRoute.Context, req CreateAgentProfileRequest) (*resp.AgentOnboardingResponse, error) {
	log.Info(sctx.Ctx, "Onboarding new agent", "AgentType", req.AgentType, "FirstName", req.FirstName, "LastName", req.LastName)

	// Prepare workflow input
	workflowInput := workflows.AgentOnboardingInput{
		AgentType:         req.AgentType,
		PersonType:        req.PersonType,
		EmployeeID:        req.EmployeeID,
		Salutation:        req.Salutation,
		FirstName:         req.FirstName,
		MiddleName:        req.MiddleName,
		LastName:          req.LastName,
		Gender:            req.Gender,
		DateOfBirth:       req.DateOfBirth,
		PAN:               req.PAN,
		CoordinatorID:     req.CoordinatorID,
		CircleID:          req.CircleID,
		CircleName:        req.CircleName,
		DivisionID:        req.DivisionID,
		DivisionName:      req.DivisionName,
		BankAccountNumber: req.BankAccountNumber,
		BankName:          req.BankName,
		BankBranch:        req.BankBranch,
		IFSCCode:          req.IFSCCode,
		AccountHolderName: req.AccountHolderName,
		POSBAccountNumber: req.POSBAccountNumber,
		POSBBranch:        req.POSBBranch,
		JoiningDate:       req.JoiningDate,
		Remarks:           req.Remarks,
		RequestedBy:       sctx.Ctx.Value("user_id").(string), // TODO: Get from auth context
	}

	// Convert addresses
	workflowInput.Addresses = make([]workflows.AddressInput, len(req.Addresses))
	for i, addr := range req.Addresses {
		workflowInput.Addresses[i] = workflows.AddressInput{
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
	workflowInput.Contacts = make([]workflows.ContactInput, len(req.Contacts))
	for i, contact := range req.Contacts {
		workflowInput.Contacts[i] = workflows.ContactInput{
			ContactType:       contact.ContactType,
			ContactNumber:     contact.ContactNumber,
			STDCode:           contact.STDCode,
			Extension:         contact.Extension,
			IsPrimary:         contact.IsPrimary,
			IsWhatsAppEnabled: contact.IsWhatsAppEnabled,
		}
	}

	// Convert emails
	workflowInput.Emails = make([]workflows.EmailInput, len(req.Emails))
	for i, email := range req.Emails {
		workflowInput.Emails[i] = workflows.EmailInput{
			EmailType:    email.EmailType,
			EmailAddress: email.EmailAddress,
			IsPrimary:    email.IsPrimary,
			IsVerified:   false, // Email verification handled separately
		}
	}

	// Start Temporal workflow
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: workflows.TaskQueueAgentOnboarding,
	}

	we, err := h.temporalClient.ExecuteWorkflow(sctx.Ctx, workflowOptions, workflows.AgentOnboardingWorkflow, workflowInput)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to start onboarding workflow", "Error", err)
		return nil, err
	}

	// Wait for workflow completion (synchronous for now)
	var result workflows.AgentOnboardingResult
	err = we.Get(sctx.Ctx, &result)
	if err != nil {
		log.Error(sctx.Ctx, "Workflow execution failed", "Error", err, "WorkflowID", we.GetID())
		return nil, err
	}

	if result.Status != "SUCCESS" {
		log.Error(sctx.Ctx, "Agent onboarding failed", "Status", result.Status, "Error", result.ErrorMessage)
		return nil, &port.Error{
			Code:    "ONBOARDING_FAILED",
			Message: result.ErrorMessage,
		}
	}

	// Fetch created profile
	profile, addresses, contacts, emails, err := h.agentRepo.GetAgentByID(sctx.Ctx, result.AgentProfileID)
	if err != nil {
		log.Error(sctx.Ctx, "Failed to fetch created profile", "Error", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Agent onboarded successfully", "AgentCode", result.AgentCode, "AgentProfileID", result.AgentProfileID)

	r := &resp.AgentOnboardingResponse{
		StatusCodeAndMessage: port.CreateSuccess,
		Data: &resp.AgentProfileResponse{
			AgentProfileID: profile.AgentProfileID,
			AgentCode:      profile.AgentCode,
			AgentType:      profile.AgentType,
			FullName:       profile.GetFullName(),
			FirstName:      profile.FirstName,
			MiddleName:     profile.MiddleName,
			LastName:       profile.LastName,
			PAN:            profile.PAN,
			Status:         profile.Status,
			CircleName:     profile.CircleName,
			DivisionName:   profile.DivisionName,
			JoiningDate:    profile.JoiningDate,
			Addresses:      resp.ToAddressResponses(addresses),
			Contacts:       resp.ToContactResponses(contacts),
			Emails:         resp.ToEmailResponses(emails),
			CreatedAt:      profile.CreatedAt,
		},
		WorkflowID: result.WorkflowID,
	}
	return r, nil
}

// GetAgentProfile retrieves agent profile by ID
// FR-IC-PROF-005: Get agent profile details
func (h *AgentOnboardingHandler) GetAgentProfile(sctx *serverRoute.Context, req AgentIDRequest) (*resp.AgentProfileDetailResponse, error) {
	log.Info(sctx.Ctx, "Fetching agent profile", "AgentID", req.ID)

	profile, addresses, contacts, emails, err := h.agentRepo.GetAgentByID(sctx.Ctx, req.ID)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching agent profile", "Error", err)
		return nil, err
	}

	r := &resp.AgentProfileDetailResponse{
		StatusCodeAndMessage: port.FetchSuccess,
		Data: &resp.AgentProfileDetailData{
			AgentProfileID:    profile.AgentProfileID,
			AgentCode:         profile.AgentCode,
			AgentType:         profile.AgentType,
			PersonType:        profile.PersonType,
			EmployeeID:        profile.EmployeeID,
			Salutation:        profile.Salutation,
			FirstName:         profile.FirstName,
			MiddleName:        profile.MiddleName,
			LastName:          profile.LastName,
			FullName:          profile.GetFullName(),
			Gender:            profile.Gender,
			DateOfBirth:       profile.DateOfBirth,
			PAN:               profile.PAN,
			AadhaarNumber:     profile.AadhaarNumber,
			BankAccountNumber: profile.BankAccountNumber,
			BankName:          profile.BankName,
			BankBranch:        profile.BankBranch,
			IFSCCode:          profile.IFSCCode,
			AccountHolderName: profile.AccountHolderName,
			POSBAccountNumber: profile.POSBAccountNumber,
			POSBBranch:        profile.POSBBranch,
			CircleID:          profile.CircleID,
			CircleName:        profile.CircleName,
			DivisionID:        profile.DivisionID,
			DivisionName:      profile.DivisionName,
			Status:            profile.Status,
			JoiningDate:       profile.JoiningDate,
			TerminationDate:   profile.TerminationDate,
			TerminationReason: profile.TerminationReason,
			Remarks:           profile.Remarks,
			Addresses:         resp.ToAddressResponses(addresses),
			Contacts:          resp.ToContactResponses(contacts),
			Emails:            resp.ToEmailResponses(emails),
			CreatedAt:         profile.CreatedAt,
			UpdatedAt:         profile.UpdatedAt,
		},
	}
	return r, nil
}

// SearchAgents searches for agents with filters
// FR-IC-PROF-004: Agent search functionality
func (h *AgentOnboardingHandler) SearchAgents(sctx *serverRoute.Context, req AgentSearchRequest) (*resp.AgentListResponse, error) {
	log.Info(sctx.Ctx, "Searching agents", "Filters", req)

	// Build filters map
	filters := make(map[string]interface{})
	if req.AgentCode != "" {
		filters["agent_code"] = req.AgentCode
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.PAN != "" {
		filters["pan"] = req.PAN
	}
	if req.AgentType != "" {
		filters["agent_type"] = req.AgentType
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.CircleID > 0 {
		filters["circle_id"] = req.CircleID
	}
	if req.DivisionID > 0 {
		filters["division_id"] = req.DivisionID
	}

	// Default pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	agents, totalCount, err := h.agentRepo.SearchAgents(sctx.Ctx, filters, page, pageSize)
	if err != nil {
		log.Error(sctx.Ctx, "Error searching agents", "Error", err)
		return nil, err
	}

	data := resp.ToAgentProfileResponses(agents)
	md := port.NewMetaDataResponse(uint64((page-1)*pageSize), uint64(pageSize), uint64(totalCount))

	r := &resp.AgentListResponse{
		StatusCodeAndMessage: port.ListSuccess,
		MetaDataResponse:     md,
		Data:                 data,
	}
	return r, nil
}

// UpdateAgentProfile updates agent profile details
// FR-IC-PROF-006: Agent profile update
func (h *AgentOnboardingHandler) UpdateAgentProfile(sctx *serverRoute.Context, req UpdateAgentProfileRequest) (*resp.AgentProfileDetailResponse, error) {
	log.Info(sctx.Ctx, "Updating agent profile", "AgentID", req.ID)

	// Build updates map
	updates := make(map[string]interface{})
	if req.Salutation != "" {
		updates["salutation"] = req.Salutation
	}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.MiddleName != "" {
		updates["middle_name"] = req.MiddleName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = req.DateOfBirth
	}
	if req.BankAccountNumber != "" {
		updates["bank_account_number"] = req.BankAccountNumber
	}
	if req.BankName != "" {
		updates["bank_name"] = req.BankName
	}
	if req.BankBranch != "" {
		updates["bank_branch"] = req.BankBranch
	}
	if req.IFSCCode != "" {
		updates["ifsc_code"] = req.IFSCCode
	}
	if req.AccountHolderName != "" {
		updates["account_holder_name"] = req.AccountHolderName
	}
	if req.Remarks != "" {
		updates["remarks"] = req.Remarks
	}

	updatedBy := sctx.Ctx.Value("user_id").(string) // TODO: Get from auth context
	err := h.agentRepo.UpdateAgentProfile(sctx.Ctx, req.ID, updates, updatedBy)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating agent profile", "Error", err)
		return nil, err
	}

	// Fetch updated profile
	profile, addresses, contacts, emails, err := h.agentRepo.GetAgentByID(sctx.Ctx, req.ID)
	if err != nil {
		log.Error(sctx.Ctx, "Error fetching updated profile", "Error", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Agent profile updated successfully", "AgentCode", profile.AgentCode)

	r := &resp.AgentProfileDetailResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: &resp.AgentProfileDetailData{
			AgentProfileID:    profile.AgentProfileID,
			AgentCode:         profile.AgentCode,
			AgentType:         profile.AgentType,
			FullName:          profile.GetFullName(),
			FirstName:         profile.FirstName,
			MiddleName:        profile.MiddleName,
			LastName:          profile.LastName,
			PAN:               profile.PAN,
			Status:            profile.Status,
			BankAccountNumber: profile.BankAccountNumber,
			BankName:          profile.BankName,
			Addresses:         resp.ToAddressResponses(addresses),
			Contacts:          resp.ToContactResponses(contacts),
			Emails:            resp.ToEmailResponses(emails),
			UpdatedAt:         profile.UpdatedAt,
		},
	}
	return r, nil
}

// UpdateAgentStatus updates agent status
// BR-IC-PROF-004: Status lifecycle management
func (h *AgentOnboardingHandler) UpdateAgentStatus(sctx *serverRoute.Context, req UpdateAgentStatusRequest) (*resp.AgentStatusUpdateResponse, error) {
	log.Info(sctx.Ctx, "Updating agent status", "AgentID", req.ID, "NewStatus", req.Status)

	updatedBy := sctx.Ctx.Value("user_id").(string) // TODO: Get from auth context
	err := h.agentRepo.UpdateAgentStatus(sctx.Ctx, req.ID, domain.AgentStatus(req.Status), req.Reason, updatedBy)
	if err != nil {
		log.Error(sctx.Ctx, "Error updating agent status", "Error", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Agent status updated successfully", "AgentID", req.ID, "Status", req.Status)

	r := &resp.AgentStatusUpdateResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: &resp.AgentStatusUpdateData{
			AgentID:   req.ID,
			Status:    req.Status,
			UpdatedAt: sctx.Ctx.Value("timestamp").(string), // TODO: Get actual timestamp
		},
	}
	return r, nil
}

// AssignCoordinator assigns or changes coordinator for an advisor
// BR-IC-AH-001: Advisors must be linked to existing Advisor Coordinator
func (h *AgentOnboardingHandler) AssignCoordinator(sctx *serverRoute.Context, req AssignCoordinatorRequest) (*resp.CoordinatorAssignmentResponse, error) {
	log.Info(sctx.Ctx, "Assigning coordinator", "AgentID", req.ID, "CoordinatorID", req.CoordinatorID)

	// Validate coordinator exists
	coordinator, err := h.agentRepo.GetCoordinatorByID(sctx.Ctx, req.CoordinatorID)
	if err != nil || coordinator == nil {
		log.Error(sctx.Ctx, "Invalid coordinator", "Error", err)
		return nil, &port.Error{
			Code:    "INVALID_COORDINATOR",
			Message: "Coordinator not found or inactive",
		}
	}

	updatedBy := sctx.Ctx.Value("user_id").(string) // TODO: Get from auth context
	err = h.hierarchyRepo.UpdateHierarchyRelationship(sctx.Ctx, req.ID, req.CoordinatorID, coordinator.AgentCode, req.EffectiveFrom, updatedBy)
	if err != nil {
		log.Error(sctx.Ctx, "Error assigning coordinator", "Error", err)
		return nil, err
	}

	log.Info(sctx.Ctx, "Coordinator assigned successfully", "AgentID", req.ID, "CoordinatorCode", coordinator.AgentCode)

	r := &resp.CoordinatorAssignmentResponse{
		StatusCodeAndMessage: port.UpdateSuccess,
		Data: &resp.CoordinatorAssignmentData{
			AgentID:         req.ID,
			CoordinatorID:   req.CoordinatorID,
			CoordinatorCode: coordinator.AgentCode,
			EffectiveFrom:   req.EffectiveFrom,
		},
	}
	return r, nil
}
