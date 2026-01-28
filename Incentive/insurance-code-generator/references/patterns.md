# Code Generation Patterns

This file contains complete code generation patterns with full implementations for insurance microservices.

**Reference**: Always check template.md (Sections 15-18) before using these patterns.

---

## Pattern 1: Workflow State Struct

**Purpose**: Cache data fetched once and reused across multiple activities to reduce database calls.

**Reference**: template.md Section 18

### Structure

```go
type {Resource}WorkflowState struct {
    // ===== Workflow Metadata =====
    {Resource}ID     string
    WorkflowStage    string
    StartedAt        time.Time
    
    // ===== Cached Data (fetched once in FetchInitialData) =====
    // Primary entity
    {Resource}Data   *domain.{Resource}
    
    // Related entities
    PolicyData       *domain.Policy
    CustomerData     *domain.Customer
    AgentData        *domain.Agent
    
    // Lookup data
    RiderData        []domain.Rider
    NominationData   []domain.Nomination
    
    // ===== Calculated Data (computed in activities) =====
    CalculationResult *domain.Calculation
    ValidationResult  *domain.ValidationResult
    ApprovalStatus    string
    
    // ===== Workflow Control =====
    RequiresApproval  bool
    ApprovalLevel     int
    RetryCount        int
}
```

### Example: Death Claim Workflow State

```go
type DeathClaimWorkflowState struct {
    // Metadata
    ClaimID          string
    WorkflowStage    string
    StartedAt        time.Time
    
    // Cached data (fetched once)
    ClaimData        *domain.Claim
    PolicyData       *domain.Policy
    DeceasedData     *domain.Customer
    BeneficiaryData  []domain.Beneficiary
    CoverageData     []domain.Coverage
    
    // Calculated data
    PayoutCalculation *domain.PayoutCalculation
    EligibilityResult *domain.EligibilityResult
    DocumentStatus    map[string]string
    
    // Workflow control
    RequiresInvestigation bool
    InvestigationDays     int
    ApprovalLevel         int
}
```

---

## Pattern 2: Initial Fetch with pgx.Batch

**Purpose**: Fetch all required data in a single database round trip using pgx.Batch.

**Reference**: template.md Section 16

### Implementation

```go
func (a *{Resource}Activities) FetchInitialDataActivity(
    ctx context.Context,
    input {Resource}Input,
) (*InitialDataResult, error) {
    batch := &pgx.Batch{}

    // ===== Query 1: Primary entity (single row) =====
    query1 := dblib.Psql.
        Select("*").
        From("{resource}_table").
        Where(sq.Eq{"{resource}_id": input.{Resource}ID})
    
    var resourceData domain.{Resource}
    dblib.QueueReturnRow(
        batch, 
        query1, 
        pgx.RowToStructByNameLax[domain.{Resource}], 
        &resourceData,
    )

    // ===== Query 2: Related entity (single row) =====
    query2 := dblib.Psql.
        Select("*").
        From("policies").
        Where(sq.Eq{"policy_id": input.PolicyID})
    
    var policyData domain.Policy
    dblib.QueueReturnRow(
        batch, 
        query2, 
        pgx.RowToStructByNameLax[domain.Policy], 
        &policyData,
    )

    // ===== Query 3: Related collection (multiple rows) =====
    query3 := dblib.Psql.
        Select("*").
        From("riders").
        Where(sq.Eq{"policy_id": input.PolicyID})
    
    var ridersData []domain.Rider
    dblib.QueueReturn(
        batch, 
        query3, 
        pgx.RowToStructByNameLax[domain.Rider],
    )

    // ===== Query 4: Another collection =====
    query4 := dblib.Psql.
        Select("*").
        From("nominations").
        Where(sq.Eq{"policy_id": input.PolicyID})
    
    var nominationsData []domain.Nomination
    dblib.QueueReturn(
        batch, 
        query4, 
        pgx.RowToStructByNameLax[domain.Nomination],
    )

    // ===== Execute all queries in one round trip =====
    results := a.db.SendBatch(ctx, batch)
    defer results.Close()

    // Scan results in order
    if err := results.QueryRow().Scan(&resourceData); err != nil {
        return nil, fmt.Errorf("failed to fetch resource: %w", err)
    }

    if err := results.QueryRow().Scan(&policyData); err != nil {
        return nil, fmt.Errorf("failed to fetch policy: %w", err)
    }

    rows, err := results.Query()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch riders: %w", err)
    }
    ridersData, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Rider])
    if err != nil {
        return nil, fmt.Errorf("failed to collect riders: %w", err)
    }

    rows, err = results.Query()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch nominations: %w", err)
    }
    nominationsData, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Nomination])
    if err != nil {
        return nil, fmt.Errorf("failed to collect nominations: %w", err)
    }

    if err := results.Close(); err != nil {
        return nil, fmt.Errorf("failed to close batch: %w", err)
    }

    // ===== Return all data =====
    return &InitialDataResult{
        ResourceData:    &resourceData,
        PolicyData:      &policyData,
        RidersData:      ridersData,
        NominationsData: nominationsData,
    }, nil
}
```

### Example: Death Claim Initial Fetch

```go
func (a *ClaimActivities) FetchDeathClaimInitialDataActivity(
    ctx context.Context,
    input DeathClaimInput,
) (*DeathClaimInitialData, error) {
    batch := &pgx.Batch{}

    // Claim details
    claimQuery := dblib.Psql.
        Select("*").
        From("claims").
        Where(sq.Eq{"claim_id": input.ClaimID})
    var claimData domain.Claim
    dblib.QueueReturnRow(batch, claimQuery, pgx.RowToStructByNameLax[domain.Claim], &claimData)

    // Policy details
    policyQuery := dblib.Psql.
        Select("*").
        From("policies").
        Where(sq.Eq{"policy_id": input.PolicyID})
    var policyData domain.Policy
    dblib.QueueReturnRow(batch, policyQuery, pgx.RowToStructByNameLax[domain.Policy], &policyData)

    // Deceased customer
    customerQuery := dblib.Psql.
        Select("*").
        From("customers").
        Where(sq.Eq{"customer_id": input.CustomerID})
    var customerData domain.Customer
    dblib.QueueReturnRow(batch, customerQuery, pgx.RowToStructByNameLax[domain.Customer], &customerData)

    // Beneficiaries
    beneficiaryQuery := dblib.Psql.
        Select("*").
        From("beneficiaries").
        Where(sq.Eq{"policy_id": input.PolicyID})
    var beneficiaries []domain.Beneficiary
    dblib.QueueReturn(batch, beneficiaryQuery, pgx.RowToStructByNameLax[domain.Beneficiary])

    // Coverage details
    coverageQuery := dblib.Psql.
        Select("*").
        From("policy_coverage").
        Where(sq.Eq{"policy_id": input.PolicyID})
    var coverages []domain.Coverage
    dblib.QueueReturn(batch, coverageQuery, pgx.RowToStructByNameLax[domain.Coverage])

    // Execute batch
    err := a.db.SendBatch(ctx, batch).Close()
    if err != nil {
        return nil, fmt.Errorf("batch fetch failed: %w", err)
    }

    return &DeathClaimInitialData{
        Claim:         &claimData,
        Policy:        &policyData,
        Customer:      &customerData,
        Beneficiaries: beneficiaries,
        Coverages:     coverages,
    }, nil
}
```

---

## Pattern 3: Activities Use State (No DB Calls)

**Purpose**: Activities perform pure calculations/logic using workflow state without additional database queries.

**Reference**: template.md Section 18

### Implementation

```go
func (a *{Resource}Activities) Calculate{Something}Activity(
    ctx context.Context,
    state *{Resource}WorkflowState,
) (*domain.{CalculationResult}, error) {
    // ✅ NO database call - all data from workflow state
    
    // Validation
    if state.PolicyData == nil {
        return nil, fmt.Errorf("policy data not available in state")
    }
    
    // Pure calculation using state data
    result := &domain.{CalculationResult}{
        Amount:      calculateAmount(state.PolicyData, state.RiderData),
        Breakdown:   buildBreakdown(state),
        AppliedRate: state.PolicyData.PremiumRate,
        Timestamp:   time.Now(),
    }
    
    return result, nil
}

func calculateAmount(policy *domain.Policy, riders []domain.Rider) float64 {
    baseAmount := policy.SumAssured
    
    for _, rider := range riders {
        if rider.Active {
            baseAmount += rider.CoverageAmount
        }
    }
    
    return baseAmount
}
```

### Example: Death Claim Payout Calculation

```go
func (a *ClaimActivities) CalculateDeathPayoutActivity(
    ctx context.Context,
    state *DeathClaimWorkflowState,
) (*domain.PayoutCalculation, error) {
    // ✅ NO database call - uses state
    
    if state.PolicyData == nil || state.CoverageData == nil {
        return nil, fmt.Errorf("required data not in state")
    }
    
    // Calculate base payout
    basePayout := state.PolicyData.SumAssured
    
    // Add rider benefits
    var riderPayout float64
    for _, coverage := range state.CoverageData {
        if coverage.CoverageType == "RIDER" && coverage.Active {
            riderPayout += coverage.CoverageAmount
        }
    }
    
    // Apply deductions
    deductions := calculateDeductions(state.PolicyData)
    
    // Final calculation
    totalPayout := basePayout + riderPayout - deductions
    
    return &domain.PayoutCalculation{
        BasePayout:   basePayout,
        RiderPayout:  riderPayout,
        Deductions:   deductions,
        TotalPayout:  totalPayout,
        CalculatedAt: time.Now(),
        CalculatedBy: "SYSTEM",
    }, nil
}

func calculateDeductions(policy *domain.Policy) float64 {
    var deductions float64
    
    // Outstanding loan
    if policy.LoanOutstanding > 0 {
        deductions += policy.LoanOutstanding
    }
    
    // Unpaid premiums
    if policy.UnpaidPremiums > 0 {
        deductions += policy.UnpaidPremiums
    }
    
    return deductions
}
```

### Example: Eligibility Validation

```go
func (a *ClaimActivities) ValidateEligibilityActivity(
    ctx context.Context,
    state *DeathClaimWorkflowState,
) (*domain.EligibilityResult, error) {
    // ✅ NO database call - validation using state
    
    result := &domain.EligibilityResult{
        Eligible:      true,
        Reasons:       []string{},
        CheckedAt:     time.Now(),
    }
    
    // Check policy status
    if state.PolicyData.Status != "ACTIVE" {
        result.Eligible = false
        result.Reasons = append(result.Reasons, "Policy is not active")
    }
    
    // Check waiting period
    policyAge := time.Since(state.PolicyData.InceptionDate).Hours() / 24
    if policyAge < 365 { // 1 year waiting period
        result.Eligible = false
        result.Reasons = append(result.Reasons, "Policy within waiting period")
    }
    
    // Check exclusions
    for _, coverage := range state.CoverageData {
        if coverage.HasExclusion && coverage.ExclusionType == "DEATH" {
            result.Eligible = false
            result.Reasons = append(result.Reasons, fmt.Sprintf("Exclusion: %s", coverage.ExclusionReason))
        }
    }
    
    return result, nil
}
```

---

## Pattern 4: Handler Starting Workflow

**Purpose**: HTTP handler initiates a Temporal workflow and returns workflow ID for async processing.

**Reference**: template.md Section 17

### Implementation

```go
func (h *{Resource}Handler) {Operation}(c echo.Context) error {
    // ===== 1. Parse and validate request =====
    var req request.{Request}
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request format",
            "details": err.Error(),
        })
    }
    
    // Validate request
    if err := req.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Validation failed",
            "details": err.Error(),
        })
    }
    
    // ===== 2. Prepare workflow input =====
    workflowInput := workflows.{Resource}WorkflowInput{
        {Resource}ID: req.{Resource}ID,
        UserID:       c.Get("user_id").(string),
        RequestData:  req,
        Timestamp:    time.Now(),
    }
    
    // ===== 3. Start workflow =====
    workflowID := fmt.Sprintf("{resource}-%s-%d", req.{Resource}ID, time.Now().Unix())
    
    workflowRun, err := h.temporalClient.ExecuteWorkflow(
        c.Request().Context(),
        client.StartWorkflowOptions{
            ID:                    workflowID,
            TaskQueue:             "{resource}-queue",
            WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
        },
        workflows.{Resource}Workflow,
        workflowInput,
    )
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to start workflow",
            "details": err.Error(),
        })
    }
    
    // ===== 4. Return workflow reference =====
    return c.JSON(http.StatusAccepted, map[string]interface{}{
        "workflow_id": workflowRun.GetID(),
        "run_id":      workflowRun.GetRunID(),
        "message":     "{Resource} processing started",
    })
}
```

### Example: Death Claim Registration Handler

```go
func (h *ClaimHandler) RegisterDeathClaim(c echo.Context) error {
    // Parse request
    var req request.RegisterDeathClaimRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request",
        })
    }
    
    // Validate
    if req.PolicyID == "" || req.DeceasedCustomerID == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Policy ID and Customer ID are required",
        })
    }
    
    // Prepare workflow input
    workflowInput := workflows.DeathClaimWorkflowInput{
        ClaimID:           uuid.New().String(),
        PolicyID:          req.PolicyID,
        DeceasedID:        req.DeceasedCustomerID,
        DeathDate:         req.DeathDate,
        DeathCertificate:  req.DeathCertificateNumber,
        ClaimantID:        req.ClaimantID,
        InitiatedBy:       c.Get("user_id").(string),
        InitiatedAt:       time.Now(),
    }
    
    // Start workflow
    workflowID := fmt.Sprintf("death-claim-%s", workflowInput.ClaimID)
    
    workflowRun, err := h.temporalClient.ExecuteWorkflow(
        c.Request().Context(),
        client.StartWorkflowOptions{
            ID:        workflowID,
            TaskQueue: "death-claim-queue",
            WorkflowExecutionTimeout: 30 * 24 * time.Hour, // 30 days
        },
        workflows.DeathClaimWorkflow,
        workflowInput,
    )
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to start death claim workflow",
        })
    }
    
    // Return workflow reference
    return c.JSON(http.StatusAccepted, map[string]interface{}{
        "claim_id":    workflowInput.ClaimID,
        "workflow_id": workflowRun.GetID(),
        "status":      "PROCESSING",
        "message":     "Death claim registration initiated",
    })
}
```

---

## Pattern 5: Complete Workflow with State

**Purpose**: Full workflow implementation using state optimization.

**Reference**: template.md Section 17-18

### Implementation

```go
func {Resource}Workflow(
    ctx workflow.Context,
    input {Resource}WorkflowInput,
) error {
    // ===== 1. Initialize workflow state =====
    state := &{Resource}WorkflowState{
        {Resource}ID:  input.{Resource}ID,
        WorkflowStage: "INITIALIZED",
        StartedAt:     workflow.Now(ctx),
    }
    
    // ===== 2. Fetch initial data (ONE database call) =====
    var initialData InitialDataResult
    err := workflow.ExecuteActivity(
        ctx,
        activities.FetchInitialDataActivity,
        input,
    ).Get(ctx, &initialData)
    if err != nil {
        return fmt.Errorf("failed to fetch initial data: %w", err)
    }
    
    // ===== 3. Populate workflow state =====
    state.{Resource}Data = initialData.ResourceData
    state.PolicyData = initialData.PolicyData
    state.RiderData = initialData.RidersData
    state.WorkflowStage = "DATA_LOADED"
    
    // ===== 4. Execute activities using state =====
    
    // Activity 1: Validation (uses state)
    var validationResult domain.ValidationResult
    err = workflow.ExecuteActivity(
        ctx,
        activities.Validate{Resource}Activity,
        state,
    ).Get(ctx, &validationResult)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    state.ValidationResult = &validationResult
    state.WorkflowStage = "VALIDATED"
    
    // Activity 2: Calculation (uses state)
    var calculationResult domain.Calculation
    err = workflow.ExecuteActivity(
        ctx,
        activities.Calculate{Resource}Activity,
        state,
    ).Get(ctx, &calculationResult)
    if err != nil {
        return fmt.Errorf("calculation failed: %w", err)
    }
    state.CalculationResult = &calculationResult
    state.WorkflowStage = "CALCULATED"
    
    // Activity 3: Conditional approval (uses state)
    if state.RequiresApproval {
        var approvalSignal ApprovalSignal
        selector := workflow.NewSelector(ctx)
        
        approvalChan := workflow.GetSignalChannel(ctx, "approval")
        selector.AddReceive(approvalChan, func(c workflow.ReceiveChannel, more bool) {
            c.Receive(ctx, &approvalSignal)
        })
        
        selector.Select(ctx) // Wait for approval
        
        state.ApprovalStatus = approvalSignal.Status
        state.WorkflowStage = "APPROVED"
    }
    
    // ===== 5. Final activity: Save results =====
    err = workflow.ExecuteActivity(
        ctx,
        activities.Save{Resource}ResultActivity,
        state,
    ).Get(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to save result: %w", err)
    }
    
    state.WorkflowStage = "COMPLETED"
    return nil
}
```

### Example: Complete Death Claim Workflow

```go
func DeathClaimWorkflow(
    ctx workflow.Context,
    input DeathClaimWorkflowInput,
) error {
    // Initialize state
    state := &DeathClaimWorkflowState{
        ClaimID:       input.ClaimID,
        WorkflowStage: "INITIALIZED",
        StartedAt:     workflow.Now(ctx),
    }
    
    // Fetch all data in one batch
    var initialData DeathClaimInitialData
    err := workflow.ExecuteActivity(
        workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
            StartToCloseTimeout: 30 * time.Second,
        }),
        activities.FetchDeathClaimInitialDataActivity,
        input,
    ).Get(ctx, &initialData)
    if err != nil {
        return fmt.Errorf("data fetch failed: %w", err)
    }
    
    // Populate state
    state.ClaimData = initialData.Claim
    state.PolicyData = initialData.Policy
    state.DeceasedData = initialData.Customer
    state.BeneficiaryData = initialData.Beneficiaries
    state.CoverageData = initialData.Coverages
    state.WorkflowStage = "DATA_LOADED"
    
    // Validate eligibility (uses state)
    var eligibility domain.EligibilityResult
    err = workflow.ExecuteActivity(ctx, activities.ValidateEligibilityActivity, state).Get(ctx, &eligibility)
    if err != nil {
        return err
    }
    if !eligibility.Eligible {
        state.WorkflowStage = "REJECTED"
        return fmt.Errorf("claim not eligible: %v", eligibility.Reasons)
    }
    state.EligibilityResult = &eligibility
    state.WorkflowStage = "ELIGIBLE"
    
    // Calculate payout (uses state)
    var payout domain.PayoutCalculation
    err = workflow.ExecuteActivity(ctx, activities.CalculateDeathPayoutActivity, state).Get(ctx, &payout)
    if err != nil {
        return err
    }
    state.PayoutCalculation = &payout
    state.WorkflowStage = "CALCULATED"
    
    // Investigation period (21 days)
    if payout.TotalPayout > 1000000 { // High value claims
        state.RequiresInvestigation = true
        state.InvestigationDays = 21
        
        err = workflow.Sleep(ctx, 21*24*time.Hour)
        if err != nil {
            return err
        }
        state.WorkflowStage = "INVESTIGATION_COMPLETE"
    }
    
    // Approval workflow
    if payout.TotalPayout > 500000 {
        var approvalSignal ApprovalSignal
        approvalChan := workflow.GetSignalChannel(ctx, "approval-signal")
        
        selector := workflow.NewSelector(ctx)
        selector.AddReceive(approvalChan, func(c workflow.ReceiveChannel, more bool) {
            c.Receive(ctx, &approvalSignal)
        })
        selector.Select(ctx)
        
        if approvalSignal.Status != "APPROVED" {
            state.WorkflowStage = "REJECTED"
            return fmt.Errorf("claim rejected in approval")
        }
        state.ApprovalLevel = approvalSignal.ApproverLevel
    }
    
    state.WorkflowStage = "APPROVED"
    
    // Process payout
    err = workflow.ExecuteActivity(ctx, activities.ProcessPayoutActivity, state).Get(ctx, nil)
    if err != nil {
        return err
    }
    
    state.WorkflowStage = "COMPLETED"
    return nil
}
```

---

## Pattern 6: Repository with pgx.Batch

**Purpose**: Repository methods that use pgx.Batch for multiple database operations.

**Reference**: template.md Section 16

### Implementation

```go
func (r *{Resource}Repository) Create{Resource}WithRelated(
    ctx context.Context,
    resource *domain.{Resource},
    related []*domain.Related,
) error {
    batch := &pgx.Batch{}
    
    // Insert main resource
    insertResource := dblib.Psql.
        Insert("{resource}_table").
        Columns("column1", "column2", "created_at").
        Values(resource.Field1, resource.Field2, time.Now()).
        Suffix("RETURNING id")
    
    var resourceID int64
    dblib.QueueReturnRow(batch, insertResource, pgx.RowTo[int64], &resourceID)
    
    // Insert related records
    for _, rel := range related {
        insertRelated := dblib.Psql.
            Insert("related_table").
            Columns("resource_id", "field1", "field2").
            Values(resourceID, rel.Field1, rel.Field2)
        
        dblib.QueueExecRow(batch, insertRelated)
    }
    
    // Execute batch
    results := r.db.SendBatch(ctx, batch)
    defer results.Close()
    
    // Get resource ID
    if err := results.QueryRow().Scan(&resourceID); err != nil {
        return fmt.Errorf("failed to insert resource: %w", err)
    }
    
    // Execute remaining inserts
    for range related {
        if _, err := results.Exec(); err != nil {
            return fmt.Errorf("failed to insert related: %w", err)
        }
    }
    
    return results.Close()
}
```

---

## Version

**Version**: 1.0.0  
**Last Updated**: 2026-01-20  
**Template.md Sections**: 15-18
