# Module-Specific Examples

This file provides complete examples of code generation for different insurance modules, showing the analysis, decision-making, and generated files for each scenario.

---

## Example 1: Death Claim Registration (Claims Module)

### Context

**Endpoint**: `POST /claims/death/register`

**Module**: Claims Processing

**Business Requirements**:
- FR-CLM-DC-001: Register death claim with policy validation
- BR-CLM-DC-001: Policy must be active and in-force
- BR-CLM-DC-002: Death must occur after policy inception + waiting period
- BR-CLM-DC-003: Exclude policies with suicide clause within 1 year
- BR-CLM-DC-004: Calculate payout = Sum Assured + Riders - Deductions
- VR-CLM-DC-001: Death certificate mandatory
- WF-CLM-DC-001: 7-step workflow with 21-day investigation for high-value claims

**Swagger Extract**:
```yaml
/claims/death/register:
  post:
    summary: Register death claim
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            properties:
              policy_id: string
              deceased_customer_id: string
              death_date: string (date)
              death_certificate_number: string
              claimant_id: string
              relationship_to_deceased: string
```

**Database Tables**:
- `claims` (claim_id, policy_id, claim_type, status, amount)
- `policies` (policy_id, sum_assured, status, inception_date)
- `customers` (customer_id, name, date_of_birth)
- `beneficiaries` (beneficiary_id, policy_id, customer_id, percentage)
- `policy_coverage` (coverage_id, policy_id, coverage_type, amount)

### Analysis & Decisions

#### Decision Matrix

| Decision Point | Evaluation | Choice | Reason |
|----------------|------------|--------|---------|
| **Temporal Workflow?** | 21-day investigation + multi-step + approvals | ✅ YES | Long-running process with human approvals |
| **Workflow State?** | Policy/customer/beneficiary data used in 4+ activities | ✅ YES | Avoid 8 separate DB calls |
| **pgx.Batch?** | Need policy + customer + beneficiaries + coverage | ✅ YES | 4 queries → 1 round trip |
| **SQL Type** | Simple lookups by ID | Squirrel | No complex joins needed |

#### Optimization Impact

**Without Optimization**:
- Activity 1 (Validate): Query policy + customer = 2 DB calls
- Activity 2 (Calculate): Query policy + coverage = 2 DB calls
- Activity 3 (Check Exclusions): Query policy + coverage = 2 DB calls
- Activity 4 (Verify Beneficiary): Query beneficiaries = 1 DB call
- Activity 5 (Save): Insert claim = 1 DB call
- **Total: 8 DB calls**

**With Optimization**:
- Initial Fetch (pgx.Batch): All data in 1 call
- Activities 1-4: Use workflow state = 0 DB calls
- Activity 5 (Save): Insert claim = 1 DB call
- **Total: 2 DB calls (75% reduction)**

### Generated Files

#### 1. Domain Model: `core/domain/claim.go`

```go
package domain

import "time"

type Claim struct {
    ClaimID             string    `db:"claim_id"`
    PolicyID            string    `db:"policy_id"`
    ClaimType           string    `db:"claim_type"`
    ClaimDate           time.Time `db:"claim_date"`
    ClaimAmount         float64   `db:"claim_amount"`
    Status              string    `db:"status"`
    DeathCertificateNum string    `db:"death_certificate_number"`
    ClaimantID          string    `db:"claimant_id"`
    CreatedAt           time.Time `db:"created_at"`
    CreatedBy           string    `db:"created_by"`
    UpdatedAt           time.Time `db:"updated_at"`
    UpdatedBy           string    `db:"updated_by"`
}

type PayoutCalculation struct {
    BasePayout   float64
    RiderPayout  float64
    Deductions   float64
    TotalPayout  float64
    CalculatedAt time.Time
    CalculatedBy string
}

type EligibilityResult struct {
    Eligible  bool
    Reasons   []string
    CheckedAt time.Time
}
```

#### 2. Request DTO: `handler/request.go` (append)

```go
type RegisterDeathClaimRequest struct {
    PolicyID              string    `json:"policy_id" validate:"required"`
    DeceasedCustomerID    string    `json:"deceased_customer_id" validate:"required"`
    DeathDate             time.Time `json:"death_date" validate:"required"`
    DeathCertificateNum   string    `json:"death_certificate_number" validate:"required"`
    ClaimantID            string    `json:"claimant_id" validate:"required"`
    RelationshipToDeceased string   `json:"relationship_to_deceased" validate:"required"`
}

func (r *RegisterDeathClaimRequest) Validate() error {
    if r.DeathDate.After(time.Now()) {
        return errors.New("death date cannot be in future")
    }
    return nil
}
```

#### 3. Response DTO: `handler/response/claim.go`

```go
package response

type DeathClaimResponse struct {
    ClaimID    string `json:"claim_id"`
    WorkflowID string `json:"workflow_id"`
    Status     string `json:"status"`
    Message    string `json:"message"`
}
```

#### 4. Handler: `handler/claim_handler.go`

```go
package handler

func (h *ClaimHandler) RegisterDeathClaim(c echo.Context) error {
    var req request.RegisterDeathClaimRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }
    
    if err := req.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }
    
    workflowInput := workflows.DeathClaimWorkflowInput{
        ClaimID:           uuid.New().String(),
        PolicyID:          req.PolicyID,
        DeceasedID:        req.DeceasedCustomerID,
        DeathDate:         req.DeathDate,
        DeathCertificate:  req.DeathCertificateNum,
        ClaimantID:        req.ClaimantID,
    }
    
    workflowRun, err := h.temporalClient.ExecuteWorkflow(
        c.Request().Context(),
        client.StartWorkflowOptions{
            ID:        "death-claim-" + workflowInput.ClaimID,
            TaskQueue: "death-claim-queue",
        },
        workflows.DeathClaimWorkflow,
        workflowInput,
    )
    
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    
    return c.JSON(http.StatusAccepted, response.DeathClaimResponse{
        ClaimID:    workflowInput.ClaimID,
        WorkflowID: workflowRun.GetID(),
        Status:     "PROCESSING",
        Message:    "Death claim registered successfully",
    })
}
```

#### 5. Workflow: `workflows/death_claim_workflow.go`

```go
package workflows

type DeathClaimWorkflowState struct {
    ClaimID           string
    WorkflowStage     string
    
    // Cached data
    ClaimData         *domain.Claim
    PolicyData        *domain.Policy
    DeceasedData      *domain.Customer
    BeneficiaryData   []domain.Beneficiary
    CoverageData      []domain.Coverage
    
    // Calculated
    PayoutCalculation *domain.PayoutCalculation
    EligibilityResult *domain.EligibilityResult
}

func DeathClaimWorkflow(ctx workflow.Context, input DeathClaimWorkflowInput) error {
    state := &DeathClaimWorkflowState{
        ClaimID:       input.ClaimID,
        WorkflowStage: "INITIALIZED",
    }
    
    // Fetch all data in ONE batch
    var initialData DeathClaimInitialData
    err := workflow.ExecuteActivity(ctx, activities.FetchDeathClaimInitialDataActivity, input).Get(ctx, &initialData)
    if err != nil {
        return err
    }
    
    state.ClaimData = initialData.Claim
    state.PolicyData = initialData.Policy
    state.DeceasedData = initialData.Customer
    state.BeneficiaryData = initialData.Beneficiaries
    state.CoverageData = initialData.Coverages
    state.WorkflowStage = "DATA_LOADED"
    
    // Validate (uses state)
    var eligibility domain.EligibilityResult
    err = workflow.ExecuteActivity(ctx, activities.ValidateEligibilityActivity, state).Get(ctx, &eligibility)
    if err != nil || !eligibility.Eligible {
        return fmt.Errorf("claim not eligible")
    }
    state.EligibilityResult = &eligibility
    
    // Calculate (uses state)
    var payout domain.PayoutCalculation
    err = workflow.ExecuteActivity(ctx, activities.CalculateDeathPayoutActivity, state).Get(ctx, &payout)
    if err != nil {
        return err
    }
    state.PayoutCalculation = &payout
    
    // Investigation if high value
    if payout.TotalPayout > 1000000 {
        workflow.Sleep(ctx, 21*24*time.Hour)
    }
    
    // Approval
    if payout.TotalPayout > 500000 {
        var approvalSignal ApprovalSignal
        workflow.GetSignalChannel(ctx, "approval").Receive(ctx, &approvalSignal)
        if approvalSignal.Status != "APPROVED" {
            return fmt.Errorf("claim rejected")
        }
    }
    
    // Process payout
    err = workflow.ExecuteActivity(ctx, activities.ProcessPayoutActivity, state).Get(ctx, nil)
    return err
}
```

#### 6. Activities: `workflows/activities/claim_activities.go`

See references/patterns.md for complete implementation.

#### 7. Repository: `repo/postgres/claim_repository.go`

```go
package postgres

func (r *ClaimRepository) SaveClaim(ctx context.Context, claim *domain.Claim) error {
    query := dblib.Psql.
        Insert("claims").
        Columns("claim_id", "policy_id", "claim_type", "claim_amount", "status").
        Values(claim.ClaimID, claim.PolicyID, claim.ClaimType, claim.ClaimAmount, claim.Status)
    
    _, err := r.db.Exec(ctx, query)
    return err
}
```

---

## Example 2: Policy Endorsement (Policy Module)

### Context

**Endpoint**: `POST /policies/endorsement`

**Module**: Policy Issuance

**Requirements**:
- FR-POL-END-001: Update policy details without changing core terms
- BR-POL-END-002: Only active policies can be endorsed
- VR-POL-END-001: Validate endorsement type

**Database**: `policies`, `policy_riders`, `policy_loans`

### Analysis & Decisions

| Decision Point | Evaluation | Choice | Reason |
|----------------|------------|--------|---------|
| **Temporal?** | Quick update (< 1 minute) | ❌ NO | Simple synchronous operation |
| **pgx.Batch?** | Update policy + riders + loans | ✅ YES | 3 tables = 1 batch |
| **SQL Type** | Simple updates by ID | Squirrel | Standard UPDATE queries |

**Optimization**: 3 separate queries → 1 pgx.Batch call (67% reduction)

### Generated Files

1. `core/domain/endorsement.go`
2. `handler/request.go` (append EndorsementRequest)
3. `handler/response/endorsement.go`
4. `handler/policy_handler.go` (direct repository call, no Temporal)
5. `repo/postgres/policy_repository.go` (with pgx.Batch)
6. `bootstrap/bootstrapper.go` (update registration)

**Handler Pattern** (No Temporal):
```go
func (h *PolicyHandler) CreateEndorsement(c echo.Context) error {
    var req request.EndorsementRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }
    
    // Direct repository call - no workflow
    endorsement, err := h.policyRepo.CreateEndorsement(c.Request().Context(), &req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    
    return c.JSON(http.StatusOK, response.EndorsementResponse{
        EndorsementID: endorsement.ID,
        Status:        "COMPLETED",
    })
}
```

---

## Example 3: Agent Commission Calculation (Agent Module)

### Context

**Endpoint**: `POST /agents/commission/calculate`

**Module**: Agent Management

**Requirements**:
- FR-AGT-COM-001: Calculate commission based on hierarchy
- BR-AGT-COM-003: Multi-level commission split
- Database: Recursive agent hierarchy CTE needed

### Analysis & Decisions

| Decision Point | Evaluation | Choice | Reason |
|----------------|------------|--------|---------|
| **Temporal?** | Fast calculation | ❌ NO | Completes in seconds |
| **pgx.Batch?** | Agent + hierarchy + policies | ✅ YES | Multiple queries |
| **SQL Type** | Recursive CTE for hierarchy | Raw SQL | Squirrel doesn't support CTE |

**Key Pattern**: Recursive CTE for agent hierarchy

### Generated Files

1. `core/domain/commission.go`
2. `handler/agent_handler.go` (direct call)
3. `repo/postgres/agent_repository.go` (with recursive CTE)
4. `bootstrap/bootstrapper.go`

**Repository Pattern** (Raw SQL with CTE):
```go
func (r *AgentRepository) CalculateCommission(ctx context.Context, policyID string) (*domain.Commission, error) {
    // Recursive CTE for agent hierarchy
    query := `
        WITH RECURSIVE agent_hierarchy AS (
            SELECT agent_id, parent_agent_id, level, commission_rate
            FROM agents
            WHERE agent_id = (SELECT agent_id FROM policies WHERE policy_id = $1)
            
            UNION ALL
            
            SELECT a.agent_id, a.parent_agent_id, ah.level + 1, a.commission_rate
            FROM agents a
            JOIN agent_hierarchy ah ON a.agent_id = ah.parent_agent_id
            WHERE ah.level < 5
        )
        SELECT 
            agent_id,
            level,
            commission_rate,
            (SELECT premium FROM policies WHERE policy_id = $1) * commission_rate / 100 as commission_amount
        FROM agent_hierarchy
    `
    
    rows, err := r.db.Query(ctx, query, policyID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // Process commission split across hierarchy
    return processCommissionSplit(rows)
}
```

---

## Example 4: Address Change (Non-Financial Services Module)

### Context

**Endpoint**: `POST /services/address-change`

**Module**: Non-Financial Services

**Requirements**:
- FR-NFR-ADDR-001: Allow address updates
- BR-NFR-ADDR-002: Permanent address requires approval
- Workflow: Conditional (approval if permanent)

### Analysis & Decisions

| Decision Point | Evaluation | Choice | Reason |
|----------------|------------|--------|---------|
| **Temporal?** | If permanent (approval needed) | ✅ CONDITIONAL | Approval workflow |
| **Workflow State?** | Policy + customer data reused | ✅ YES | 3 activities use same data |
| **pgx.Batch?** | Fetch policy + customer + communications | ✅ YES | 3 queries |

**Conditional Logic**: Temporary = direct update, Permanent = approval workflow

### Generated Files

1. `core/domain/address_change.go`
2. `workflows/address_change_workflow.go` (conditional logic)
3. `workflows/activities/address_change_activities.go`
4. `handler/nfr_handler.go`
5. `repo/postgres/nfr_repository.go`

**Workflow Pattern** (Conditional):
```go
func AddressChangeWorkflow(ctx workflow.Context, input AddressChangeInput) error {
    // Fetch data with pgx.Batch
    var initialData AddressChangeInitialData
    workflow.ExecuteActivity(ctx, activities.FetchAddressDataActivity, input).Get(ctx, &initialData)
    
    state := &AddressChangeWorkflowState{
        PolicyData:   initialData.Policy,
        CustomerData: initialData.Customer,
    }
    
    // Conditional: Permanent address requires approval
    if input.AddressType == "PERMANENT" {
        var approvalSignal ApprovalSignal
        workflow.GetSignalChannel(ctx, "approval").Receive(ctx, &approvalSignal)
        
        if approvalSignal.Status != "APPROVED" {
            return fmt.Errorf("address change rejected")
        }
    }
    
    // Update address
    workflow.ExecuteActivity(ctx, activities.UpdateAddressActivity, state).Get(ctx, nil)
    return nil
}
```

---

## Example 5: GL Integration Batch (Accounting Module)

### Context

**Endpoint**: `POST /accounting/gl-integration/batch`

**Module**: Accounting

**Requirements**:
- FR-ACC-GL-001: Integrate transactions to GL
- BR-ACC-GL-002: Batch processing with retries
- Workflow: Batch process with error handling

### Analysis & Decisions

| Decision Point | Evaluation | Choice | Reason |
|----------------|------------|--------|---------|
| **Temporal?** | Batch processing with retries | ✅ YES | Long-running batch |
| **Raw SQL?** | Complex INSERT...SELECT for GL entries | ✅ YES | Bulk operations |
| **pgx.Batch?** | Bulk inserts | ✅ YES | Many inserts |

### Generated Files

1. `core/domain/gl_entry.go`
2. `workflows/gl_integration_workflow.go`
3. `workflows/activities/gl_activities.go`
4. `handler/accounting_handler.go`
5. `repo/postgres/accounting_repository.go`

**Repository Pattern** (Bulk INSERT...SELECT):
```go
func (r *AccountingRepository) CreateGLBatch(ctx context.Context, batchID string) error {
    query := `
        INSERT INTO gl_entries (
            entry_id, account_code, debit, credit, transaction_id, batch_id
        )
        SELECT 
            gen_random_uuid(),
            CASE 
                WHEN t.transaction_type = 'PREMIUM' THEN '1001'
                WHEN t.transaction_type = 'CLAIM' THEN '2001'
            END,
            CASE WHEN t.amount > 0 THEN t.amount ELSE 0 END,
            CASE WHEN t.amount < 0 THEN ABS(t.amount) ELSE 0 END,
            t.transaction_id,
            $1
        FROM transactions t
        WHERE t.gl_status = 'PENDING'
          AND t.created_at >= NOW() - INTERVAL '1 day'
    `
    
    _, err := r.db.Exec(ctx, query, batchID)
    return err
}
```

---

## Example 6: Customer Portal Profile Update (Portal/Billing Module)

### Context

**Endpoint**: `POST /portal/customer/profile-update`

**Module**: Customer Portal

**Requirements**:
- FR-POR-PROF-001: Update customer profile
- Quick synchronous update

### Analysis & Decisions

| Decision Point | Evaluation | Choice | Reason |
|----------------|------------|--------|---------|
| **Temporal?** | Quick update | ❌ NO | < 1 second |
| **pgx.Batch?** | Update customer + preferences + history | ✅ YES | 3 tables |
| **SQL Type** | Simple updates | Squirrel | Standard UPDATEs |

### Generated Files

1. `handler/portal_handler.go` (direct call)
2. `repo/postgres/customer_repository.go` (pgx.Batch for updates)

---

## Example 7: Aadhaar Verification (KYC/BCP Module)

### Context

**Endpoint**: `POST /kyc/aadhaar/verify`

**Module**: KYC

**Requirements**:
- FR-KYC-AAD-001: Verify Aadhaar with UIDAI
- External API call with retries

### Analysis & Decisions

| Decision Point | Evaluation | Choice | Reason |
|----------------|------------|--------|---------|
| **Temporal?** | External API with retries | ✅ YES | Retry policies needed |
| **Workflow State?** | Single external call | ❌ NO | No data reuse |

### Generated Files

1. `workflows/kyc_workflow.go`
2. `workflows/activities/kyc_activities.go` (external API call)
3. `handler/kyc_handler.go`

**Activity Pattern** (External API):
```go
func (a *KYCActivities) VerifyAadhaarActivity(ctx context.Context, input AadhaarInput) (*AadhaarResult, error) {
    // Call external UIDAI API
    client := &http.Client{Timeout: 10 * time.Second}
    
    resp, err := client.Post(
        "https://uidai.gov.in/verify",
        "application/json",
        bytes.NewBuffer([]byte(input.AadhaarNumber)),
    )
    
    if err != nil {
        return nil, err // Temporal will retry
    }
    
    defer resp.Body.Close()
    
    var result AadhaarResult
    json.NewDecoder(resp.Body).Decode(&result)
    
    return &result, nil
}
```

---

## Summary Table

| Example | Module | Temporal | State | Batch | SQL Type | Files |
|---------|--------|----------|-------|-------|----------|-------|
| Death Claim | Claims | ✅ | ✅ | ✅ | Squirrel | 8 |
| Policy Endorsement | Policy | ❌ | ❌ | ✅ | Squirrel | 6 |
| Agent Commission | Agent | ❌ | ❌ | ✅ | Raw SQL | 4 |
| Address Change | NFR | ✅ | ✅ | ✅ | Squirrel | 7 |
| GL Integration | Accounting | ✅ | ❌ | ✅ | Raw SQL | 7 |
| Profile Update | Portal | ❌ | ❌ | ✅ | Squirrel | 2 |
| Aadhaar Verify | KYC | ✅ | ❌ | ❌ | Squirrel | 3 |

---

## Version

**Version**: 1.0.0  
**Last Updated**: 2026-01-20
