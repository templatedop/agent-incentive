# Agent Commission Management System - Implementation Plan

## Project Overview

**Project Name**: Agent Incentive, Commission and Producer Management System
**Framework**: N-API Template (Golang)
**Workflow Engine**: Temporal.io
**Database**: PostgreSQL
**Total Endpoints**: 105 APIs
**Estimated Duration**: 8-10 weeks
**Team Size**: 1-2 developers

---

## Implementation Methodology

Following **insurance-code-generator skill** three-phase approach:

1. **Phase 1**: Multi-Source Analysis (Requirements, Swagger, Database)
2. **Phase 2**: Pattern Decisions (Temporal, Batch, Squirrel vs Raw SQL)
3. **Phase 3**: Code Generation (Domain, Handler, Repository, Workflows)

---

## Phase Breakdown

### PHASE 0: Project Foundation & Setup (Week 1: Days 1-2)
**Duration**: 2 days
**Goal**: Establish project structure and tooling

#### Tasks:
- [ ] **0.1** Create project directory structure following N-API Template
  - `main.go`, `go.mod`, `bootstrap/`, `core/`, `handler/`, `repo/`, `workflows/`
  - `configs/`, `db/migrations/`, `docs/`
- [ ] **0.2** Initialize Go modules with dependencies
  - pgx, squirrel, temporal, fx, validator, gin
  - Reference: `Incentive/template.md` Section 4
- [ ] **0.3** Setup configuration files (dev, sit, staging, prod)
  - Database connections, Temporal connection, JWT settings
  - Reference: `Incentive/template.md` Section 5
- [ ] **0.4** Create bootstrap/bootstrapper.go with FX modules
  - FxRepo, FxHandler, FxWorkflow modules
  - Reference: `Incentive/template.md` Section 3
- [ ] **0.5** Setup database migration structure
  - Create `db/migrations/` folder
  - Copy `Incentive/db/schema.sql` as base migration
- [ ] **0.6** Create main.go application entry point
  - Reference: `Incentive/template.md` Section 2
- [ ] **0.7** Setup Swagger documentation auto-generation
- [ ] **0.8** Create `.gitignore`, `README.md`, `Makefile`
- [ ] **0.9** Verify build and run hello-world endpoint

**Testing**: Build succeeds, server starts, health check endpoint responds

---

### PHASE 1: Agent Onboarding & Profile Management (Week 1-2: Days 3-10)
**Duration**: 8 days
**APIs**: 15 endpoints
**Goal**: Complete agent lifecycle from onboarding to termination

#### Module 1.1: Agent Onboarding (Days 3-5)
**APIs**: 6 endpoints
**Temporal Workflow**: WF-IC-ONB-001 (Agent Onboarding Workflow)

##### Tasks:
- [ ] **1.1.1** Read and analyze requirements
  - **FR**: FR-IC-PROF-001 to FR-IC-PROF-003
  - **BR**: BR-IC-AH-001 to BR-IC-AH-004
  - **VR**: VR-IC-PROF-001 to VR-IC-PROF-005
  - **WF**: WF-IC-ONB-001
  - **Swagger**: `/agents/new/*` endpoints (lines 95-300)
  - **DB**: `agent_profiles`, `agent_hierarchy`, `agent_addresses` tables

- [ ] **1.1.2** Make pattern decisions using insurance-code-generator skill
  - Use Temporal? → YES (multi-step approval, HRMS integration)
  - Use Workflow State? → YES (reduce HRMS/DB round trips)
  - Use pgx.Batch? → YES (insert profile + addresses + contacts in one trip)
  - Squirrel vs Raw SQL? → Squirrel for simple CRUD, Raw for hierarchy queries

- [ ] **1.1.3** Create domain models
  - `core/domain/agent_profile.go` - Agent, Coordinator, Employee, Field Officer
  - `core/domain/agent_address.go`, `core/domain/agent_contact.go`
  - Add comments: BR-IC-AH-001 (Advisor-Coordinator link), etc.

- [ ] **1.1.4** Create request/response DTOs
  - `handler/request.go`: `CreateAgentRequest`, `SelectAgentTypeRequest`
  - `handler/response/agent.go`: `AgentProfileResponse`, `OnboardingSessionResponse`
  - Add validator tags and generate validators

- [ ] **1.1.5** Create repository layer
  - `repo/postgres/agent_profile.go`:
    - `CreateAgentProfile(ctx, agent)` - Use pgx.Batch for profile+address+contact
    - `GetCoordinatorByID(ctx, id)` - For advisor validation
    - `CheckPANExists(ctx, pan)` - Duplicate PAN check (VR-IC-PROF-002)
    - `GetAgentByEmployeeID(ctx, empID)` - HRMS lookup
  - Add comments with BR/VR/FR traceability

- [ ] **1.1.6** Create Temporal workflow for agent onboarding
  - `workflows/agent_onboarding_workflow.go`:
    - `AgentOnboardingWorkflow(ctx, request)` - Orchestration
    - Workflow State: Cache coordinator data, HRMS data
  - `workflows/activities/agent_onboarding_activities.go`:
    - `ValidatePANActivity` - VR-IC-PROF-002
    - `FetchHRMSDataActivity` - BR-IC-AH-003
    - `CreateAgentProfileActivity` - Calls repository
    - `AssignCoordinatorActivity` - BR-IC-AH-001
    - `GenerateAgentCodeActivity` - FR-IC-PROF-001
  - Add comments with workflow step references

- [ ] **1.1.7** Create handler with routes
  - `handler/agent_onboarding.go`:
    - `POST /agents/new/init` - Initialize onboarding session
    - `POST /agents/new/{sessionId}/select-type` - Select agent type
    - `POST /agents/new/{sessionId}/profile-details` - Submit profile data
    - `POST /agents/new/{sessionId}/select-coordinator` - Link coordinator
    - `GET /agents/new/{sessionId}/preview` - Preview before save
    - `POST /agents/new/{sessionId}/submit` - Start Temporal workflow
  - Add route registration in `Routes()` method

- [ ] **1.1.8** Unit tests
  - `handler/agent_onboarding_test.go` - Handler tests
  - `repo/postgres/agent_profile_test.go` - Repository tests
  - `workflows/agent_onboarding_workflow_test.go` - Workflow tests
  - Coverage target: 80%+

**Traceability**:
- FR-IC-PROF-001, FR-IC-PROF-002, FR-IC-PROF-003
- BR-IC-AH-001 (Advisor-Coordinator), BR-IC-AH-002 (Coordinator-Circle), BR-IC-AH-003 (HRMS), BR-IC-AH-004 (Field Officer)
- VR-IC-PROF-001 to VR-IC-PROF-005
- WF-IC-ONB-001

---

#### Module 1.2: Agent Profile Management (Days 6-8)
**APIs**: 6 endpoints
**No Temporal Workflow** (simple CRUD operations)

##### Tasks:
- [ ] **1.2.1** Read and analyze requirements
  - **FR**: FR-IC-PROF-004 to FR-IC-PROF-006
  - **BR**: BR-IC-PROF-001 to BR-IC-PROF-004
  - **VR**: VR-IC-PROF-006 to VR-IC-PROF-010
  - **Swagger**: `/agents/search`, `/agents/{agentId}/*` endpoints

- [ ] **1.2.2** Pattern decisions
  - Use Temporal? → NO (simple CRUD)
  - Use pgx.Batch? → YES (fetch profile + addresses + licenses in one query)
  - Squirrel vs Raw SQL? → Squirrel for simple queries, Raw for complex search

- [ ] **1.2.3** Extend domain models
  - Update `core/domain/agent_profile.go` with search fields

- [ ] **1.2.4** Create request/response DTOs
  - `handler/request.go`: `AgentSearchRequest`, `UpdateAgentRequest`
  - `handler/response/agent.go`: `AgentSearchResponse`, `AgentDetailResponse`

- [ ] **1.2.5** Create repository methods
  - `repo/postgres/agent_profile.go`:
    - `SearchAgents(ctx, filters)` - Dynamic WHERE with Squirrel
    - `GetAgentByID(ctx, id)` - Use pgx.Batch for profile+addresses+licenses
    - `UpdateAgentProfile(ctx, agent)` - BR-IC-PROF-003 (audit trail)
    - `UpdateAgentName(ctx, id, name)` - VR-IC-PROF-006
    - `UpdateAgentStatus(ctx, id, status)` - BR-IC-PROF-004

- [ ] **1.2.6** Create handler with routes
  - `handler/agent_profile.go`:
    - `POST /agents/search` - FR-IC-PROF-004
    - `GET /agents/{agentId}` - FR-IC-PROF-005
    - `PUT /agents/{agentId}` - FR-IC-PROF-006
    - `PUT /agents/{agentId}/name` - VR-IC-PROF-006
    - `PUT /agents/{agentId}/status` - BR-IC-PROF-004
    - `GET /agents/{agentId}/audit-trail` - BR-IC-PROF-003

- [ ] **1.2.7** Unit tests
  - Test search with multiple filters
  - Test batch query performance
  - Test audit trail generation

**Traceability**:
- FR-IC-PROF-004, FR-IC-PROF-005, FR-IC-PROF-006
- BR-IC-PROF-001 to BR-IC-PROF-004
- VR-IC-PROF-006 to VR-IC-PROF-010

---

#### Module 1.3: Agent Termination (Days 9-10)
**APIs**: 3 endpoints
**Temporal Workflow**: WF-IC-TERM-001 (Agent Termination Workflow)

##### Tasks:
- [ ] **1.3.1** Read and analyze requirements
  - **FR**: FR-IC-PROF-007
  - **BR**: BR-IC-TERM-001, BR-IC-TERM-002
  - **VR**: VR-IC-TERM-001
  - **WF**: WF-IC-TERM-001
  - **Swagger**: `/agents/{agentId}/terminate` endpoints

- [ ] **1.3.2** Pattern decisions
  - Use Temporal? → YES (multi-step: stop commission, settle dues, deactivate)
  - Use Workflow State? → YES (cache agent, commission, license data)
  - Use pgx.Batch? → YES (update agent + licenses + hierarchy in one transaction)

- [ ] **1.3.3** Create Temporal workflow
  - `workflows/agent_termination_workflow.go`:
    - `AgentTerminationWorkflow(ctx, agentID, reason)`
  - `workflows/activities/agent_termination_activities.go`:
    - `CheckPendingCommissionsActivity` - BR-IC-TERM-001
    - `SettleCommissionsActivity` - BR-IC-TERM-002
    - `DeactivateLicensesActivity`
    - `UpdateAgentStatusActivity` - Status = TERMINATED
    - `NotifyStakeholdersActivity`

- [ ] **1.3.4** Create handler
  - `handler/agent_termination.go`:
    - `POST /agents/{agentId}/terminate` - Start workflow
    - `GET /agents/{agentId}/termination-status` - Check progress
    - `POST /agents/{agentId}/termination-cancel` - Cancel if needed

- [ ] **1.3.5** Unit tests
  - Test termination with pending commissions
  - Test termination with settled dues
  - Test workflow cancellation

**Traceability**:
- FR-IC-PROF-007
- BR-IC-TERM-001, BR-IC-TERM-002
- VR-IC-TERM-001
- WF-IC-TERM-001

---

### PHASE 2: License Management (Week 2-3: Days 11-15)
**Duration**: 5 days
**APIs**: 10 endpoints
**Goal**: Complete license lifecycle management

#### Module 2.1: License CRUD & Renewal (Days 11-13)
**APIs**: 6 endpoints
**Temporal Workflow**: WF-IC-LIC-002 (License Renewal Workflow)

##### Tasks:
- [ ] **2.1.1** Read and analyze requirements
  - **FR**: FR-IC-LIC-001 to FR-IC-LIC-003
  - **BR**: BR-IC-LIC-001 to BR-IC-LIC-003
  - **VR**: VR-IC-LIC-001 to VR-IC-LIC-003
  - **WF**: WF-IC-LIC-002
  - **Swagger**: `/agents/{agentId}/licenses/*` endpoints
  - **DB**: `agent_licenses`, `license_renewal_reminders` tables

- [ ] **2.1.2** Pattern decisions
  - Use Temporal? → YES (renewal workflow has approval + 3-day SLA)
  - Use Workflow State? → YES (cache agent data)
  - Use pgx.Batch? → YES (insert license + renewal reminder records)

- [ ] **2.1.3** Create domain models
  - `core/domain/agent_license.go` - License, RenewalReminder

- [ ] **2.1.4** Create repository
  - `repo/postgres/agent_license.go`:
    - `CreateLicense(ctx, license)` - Use pgx.Batch
    - `GetLicensesByAgentID(ctx, agentID)`
    - `UpdateLicenseStatus(ctx, id, status)`
    - `GetExpiringLicenses(ctx, daysUntilExpiry)` - For reminder batch
    - `RenewLicense(ctx, licenseID, renewalDate)` - BR-IC-LIC-001

- [ ] **2.1.5** Create Temporal workflow
  - `workflows/license_renewal_workflow.go`:
    - `LicenseRenewalWorkflow(ctx, licenseID)`
    - SLA: 3 working days
  - `workflows/activities/license_renewal_activities.go`:
    - `ValidateRenewalDocumentsActivity` - VR-IC-LIC-002
    - `UpdateLicenseActivity`
    - `NotifyAgentActivity`

- [ ] **2.1.6** Create handler
  - `handler/agent_license.go`:
    - `POST /agents/{agentId}/licenses` - Add new license
    - `GET /agents/{agentId}/licenses` - List all licenses
    - `GET /agents/{agentId}/licenses/{licenseId}` - Get details
    - `POST /agents/{agentId}/licenses/{licenseId}/renew` - Start renewal workflow
    - `PUT /agents/{agentId}/licenses/{licenseId}` - Update license
    - `DELETE /agents/{agentId}/licenses/{licenseId}` - Cancel license

- [ ] **2.1.7** Unit tests

**Traceability**:
- FR-IC-LIC-001, FR-IC-LIC-002, FR-IC-LIC-003
- BR-IC-LIC-001, BR-IC-LIC-002, BR-IC-LIC-003
- VR-IC-LIC-001, VR-IC-LIC-002, VR-IC-LIC-003
- WF-IC-LIC-002

---

#### Module 2.2: License Expiry & Reminders (Days 14-15)
**APIs**: 4 endpoints
**Temporal Workflow**: WF-IC-LIC-001 (License Reminder Workflow)

##### Tasks:
- [ ] **2.2.1** Read and analyze requirements
  - **FR**: FR-IC-LIC-004, FR-IC-LIC-005
  - **BR**: BR-IC-LIC-004, BR-IC-LIC-005
  - **WF**: WF-IC-LIC-001
  - **Swagger**: `/licenses/reminders/*` endpoints

- [ ] **2.2.2** Create Temporal workflow for reminders
  - `workflows/license_reminder_workflow.go`:
    - `LicenseReminderBatchWorkflow(ctx)` - Runs daily
    - Sends reminders at T-30, T-15, T-7, T-0 days
  - `workflows/activities/license_reminder_activities.go`:
    - `FetchExpiringLicensesActivity` - BR-IC-LIC-005
    - `SendReminderEmailActivity`
    - `SendReminderSMSActivity`
    - `LogReminderActivity`
    - `DeactivateExpiredAgentsActivity` - VR-IC-LIC-003

- [ ] **2.2.3** Create handler
  - `handler/license_reminder.go`:
    - `GET /licenses/reminders/upcoming` - View upcoming reminders
    - `POST /licenses/reminders/trigger` - Manual trigger (admin)
    - `GET /licenses/reminders/history` - View reminder history
    - `GET /licenses/expired` - List expired licenses

- [ ] **2.2.4** Unit tests
  - Test reminder scheduling
  - Test expiry deactivation logic

**Traceability**:
- FR-IC-LIC-004, FR-IC-LIC-005
- BR-IC-LIC-004, BR-IC-LIC-005
- VR-IC-LIC-003
- WF-IC-LIC-001

---

### PHASE 3: Commission Rate Configuration (Week 3: Days 16-17)
**Duration**: 2 days
**APIs**: 6 endpoints
**Goal**: Setup commission rate tables

#### Tasks:
- [ ] **3.1** Read and analyze requirements
  - **FR**: FR-IC-RATE-001 to FR-IC-RATE-003
  - **BR**: BR-IC-COM-006, BR-IC-COM-007
  - **Swagger**: `/commission-rates/*` endpoints
  - **DB**: `commission_rate_config` table

- [ ] **3.2** Pattern decisions
  - Use Temporal? → NO (simple CRUD)
  - Use pgx.Batch? → NO (single row operations)
  - Squirrel vs Raw SQL? → Squirrel

- [ ] **3.3** Create domain models
  - `core/domain/commission_rate.go` - RateConfig

- [ ] **3.4** Create repository
  - `repo/postgres/commission_rate.go`:
    - `CreateRate(ctx, rate)`
    - `GetRateByParams(ctx, agentType, product, commissionType, year)` - BR-IC-COM-006
    - `UpdateRate(ctx, rate)`
    - `GetAllRates(ctx, filters)`
    - `ArchiveRate(ctx, id)` - Soft delete

- [ ] **3.5** Create handler
  - `handler/commission_rate.go`:
    - `POST /commission-rates` - Create rate
    - `GET /commission-rates` - List rates with filters
    - `GET /commission-rates/{rateId}` - Get rate details
    - `PUT /commission-rates/{rateId}` - Update rate
    - `DELETE /commission-rates/{rateId}` - Archive rate
    - `GET /commission-rates/lookup` - Lookup rate by params

- [ ] **3.6** Unit tests
  - Test rate lookup logic
  - Test rate version management

**Traceability**:
- FR-IC-RATE-001, FR-IC-RATE-002, FR-IC-RATE-003
- BR-IC-COM-006, BR-IC-COM-007

---

### PHASE 4: Commission Calculation & Processing (Week 4-5: Days 18-30)
**Duration**: 13 days
**APIs**: 18 endpoints
**Goal**: Complete commission lifecycle from calculation to approval

#### Module 4.1: Commission Batch Calculation (Days 18-22)
**APIs**: 6 endpoints
**Temporal Workflow**: WF-IC-COM-001 (Commission Batch Processing)

##### Tasks:
- [ ] **4.1.1** Read and analyze requirements
  - **FR**: FR-IC-COM-001, FR-IC-COM-002
  - **BR**: BR-IC-COM-001, BR-IC-COM-006, BR-IC-COM-009, BR-IC-COM-012
  - **VR**: VR-IC-COM-001
  - **WF**: WF-IC-COM-001 (Steps 1-2)
  - **Swagger**: `/commission/batch/*` endpoints
  - **DB**: `commission_records`, `commission_batch_log` tables

- [ ] **4.1.2** Pattern decisions
  - Use Temporal? → YES (6-hour SLA, batch processing with retries)
  - Use Workflow State? → YES (cache rate table, policy data)
  - Use pgx.Batch? → YES (batch insert commission records - critical optimization)
  - Raw SQL for bulk operations

- [ ] **4.1.3** Create domain models
  - `core/domain/commission.go` - CommissionRecord, BatchLog
  - Add calculation formula: (Annualised Premium × Rate) / 100

- [ ] **4.1.4** Create repository
  - `repo/postgres/commission.go`:
    - `GetPoliciesForCommission(ctx, period)` - Raw SQL with JOIN
    - `BatchInsertCommissions(ctx, records)` - Use pgx.Batch (CRITICAL)
    - `UpdateCommissionStatus(ctx, id, status)`
    - `CreateBatchLog(ctx, log)`
    - `UpdateBatchProgress(ctx, batchID, processed, total)` - BR-IC-COM-012

- [ ] **4.1.5** Create Temporal workflow
  - `workflows/commission_batch_workflow.go`:
    - `CommissionBatchWorkflow(ctx, period)` - SLA: 6 hours
    - Workflow State: Cache all commission rates, agent hierarchy
  - `workflows/activities/commission_batch_activities.go`:
    - `FetchPoliciesActivity` - Get policies for period
    - `FetchCommissionRatesActivity` - Cache all rates
    - `CalculateCommissionsActivity` - Calculate in batches of 1000
    - `SaveCommissionsActivity` - Use pgx.Batch for bulk insert
    - `UpdateProgressActivity` - BR-IC-COM-012
    - `HandleFailedRecordsActivity` - Retry logic
  - Add comments: FR-IC-COM-002, BR-IC-COM-001

- [ ] **4.1.6** Create handler
  - `handler/commission_batch.go`:
    - `POST /commission/batch/trigger` - Start batch workflow
    - `GET /commission/batch/{batchId}/status` - Monitor progress
    - `GET /commission/batch/{batchId}/progress` - Detailed progress
    - `POST /commission/batch/{batchId}/cancel` - Cancel batch
    - `GET /commission/batch/history` - Past batches
    - `GET /commission/batch/{batchId}/failed-records` - Failed records

- [ ] **4.1.7** Unit tests
  - Test batch calculation logic
  - Test pgx.Batch performance (target: 10,000 records/second)
  - Test retry mechanism

**Traceability**:
- FR-IC-COM-001, FR-IC-COM-002
- BR-IC-COM-001, BR-IC-COM-006, BR-IC-COM-009, BR-IC-COM-012
- VR-IC-COM-001
- WF-IC-COM-001 (Steps 1-2)

---

#### Module 4.2: Trial Statement Generation & Approval (Days 23-26)
**APIs**: 8 endpoints
**Temporal Workflow**: WF-IC-COM-001 (Steps 3-5)

##### Tasks:
- [ ] **4.2.1** Read and analyze requirements
  - **FR**: FR-IC-COM-003, FR-IC-COM-004, FR-IC-COM-005, FR-IC-COM-006
  - **BR**: BR-IC-COM-002, BR-IC-COM-003, BR-IC-COM-004, BR-IC-COM-005
  - **VR**: VR-IC-COM-002, VR-IC-COM-003, VR-IC-COM-004, VR-IC-COM-005
  - **WF**: WF-IC-COM-001 (Steps 3-5)
  - **Swagger**: `/commission/trial-statements/*` endpoints
  - **DB**: `commission_trial_statements`, `commission_final_statements` tables

- [ ] **4.2.2** Pattern decisions
  - Use Temporal? → YES (7-day SLA for approval, human-in-the-loop)
  - Use Workflow State? → YES (cache trial data)
  - Use pgx.Batch? → YES (aggregate commissions by agent)
  - Squirrel for simple queries, Raw SQL for aggregations

- [ ] **4.2.3** Create domain models
  - `core/domain/commission_statement.go` - TrialStatement, FinalStatement
  - Add TDS calculation formula

- [ ] **4.2.4** Create repository
  - `repo/postgres/commission_statement.go`:
    - `GenerateTrialStatements(ctx, batchID)` - Raw SQL with GROUP BY
    - `GetTrialStatements(ctx, filters)` - Search with pagination
    - `GetTrialStatementByID(ctx, id)` - Use pgx.Batch for statement+records
    - `ApproveTrialStatement(ctx, id, approverID)` - BR-IC-COM-002
    - `RejectTrialStatement(ctx, id, reason)`
    - `PartialApproveStatement(ctx, id, percentage)` - BR-IC-COM-005

- [ ] **4.2.5** Extend Temporal workflow
  - Update `workflows/commission_batch_workflow.go`:
    - Add Step 3: Trial Statement Generation
    - Add Step 4: Wait for Finance Approval (signal pattern)
    - Add Step 5: Process Approval/Rejection
  - `workflows/activities/commission_statement_activities.go`:
    - `GenerateTrialStatementsActivity` - FR-IC-COM-003
    - `CalculateTDSActivity` - BR-IC-COM-003
    - `NotifyFinanceTeamActivity`
    - `WaitForApprovalActivity` - Signal wait
    - `ProcessApprovalActivity` - VR-IC-COM-005
    - `ProcessRejectionActivity`

- [ ] **4.2.6** Create handler
  - `handler/commission_trial_statement.go`:
    - `GET /commission/trial-statements` - List with filters
    - `GET /commission/trial-statements/{statementId}` - Get details
    - `GET /commission/trial-statements/{statementId}/records` - Commission records
    - `POST /commission/trial-statements/{statementId}/approve` - Approve
    - `POST /commission/trial-statements/{statementId}/reject` - Reject
    - `POST /commission/trial-statements/{statementId}/approve-partial` - Partial approve
    - `GET /commission/trial-statements/{statementId}/history` - Approval history
    - `POST /commission/trial-statements/bulk-approve` - Bulk approval

- [ ] **4.2.7** Unit tests
  - Test TDS calculation
  - Test partial approval logic
  - Test workflow signal handling

**Traceability**:
- FR-IC-COM-003, FR-IC-COM-004, FR-IC-COM-005, FR-IC-COM-006
- BR-IC-COM-002, BR-IC-COM-003, BR-IC-COM-004, BR-IC-COM-005
- VR-IC-COM-002, VR-IC-COM-003, VR-IC-COM-004, VR-IC-COM-005
- WF-IC-COM-001 (Steps 3-5)

---

#### Module 4.3: Final Statement Generation (Days 27-30)
**APIs**: 4 endpoints
**Temporal Workflow**: WF-IC-COM-001 (Steps 6-7)

##### Tasks:
- [ ] **4.3.1** Read and analyze requirements
  - **FR**: FR-IC-COM-007, FR-IC-COM-008
  - **BR**: BR-IC-COM-002 (locking)
  - **WF**: WF-IC-COM-001 (Steps 6-7)
  - **Swagger**: `/commission/final-statements/*` endpoints

- [ ] **4.3.2** Create repository methods
  - `repo/postgres/commission_statement.go`:
    - `GenerateFinalStatements(ctx, trialStatementID)` - Lock trial data
    - `GetFinalStatements(ctx, filters)`
    - `GetFinalStatementByID(ctx, id)`
    - `ExportFinalStatementPDF(ctx, id)` - Generate PDF

- [ ] **4.3.3** Extend Temporal workflow
  - Update `workflows/commission_batch_workflow.go`:
    - Add Step 6: Final Statement Generation
    - Add Step 7: Queue for Disbursement
  - `workflows/activities/commission_statement_activities.go`:
    - `LockTrialDataActivity` - BR-IC-COM-002
    - `GenerateFinalStatementsActivity`
    - `GeneratePDFActivity`
    - `QueueDisbursementActivity`

- [ ] **4.3.4** Create handler
  - `handler/commission_final_statement.go`:
    - `GET /commission/final-statements` - List with filters
    - `GET /commission/final-statements/{statementId}` - Get details
    - `GET /commission/final-statements/{statementId}/pdf` - Download PDF
    - `GET /commission/final-statements/{statementId}/records` - Records

- [ ] **4.3.5** Unit tests
  - Test locking mechanism
  - Test PDF generation

**Traceability**:
- FR-IC-COM-007, FR-IC-COM-008
- BR-IC-COM-002
- WF-IC-COM-001 (Steps 6-7)

---

### PHASE 5: Commission Disbursement (Week 6: Days 31-37)
**Duration**: 7 days
**APIs**: 12 endpoints
**Goal**: Process commission payments via Cheque and EFT

#### Module 5.1: Disbursement Processing (Days 31-34)
**APIs**: 8 endpoints
**Temporal Workflow**: WF-IC-DIS-001 (Disbursement Processing)

##### Tasks:
- [ ] **5.1.1** Read and analyze requirements
  - **FR**: FR-IC-DIS-001 to FR-IC-DIS-004
  - **BR**: BR-IC-COM-008, BR-IC-COM-010, BR-IC-COM-011
  - **VR**: VR-IC-DIS-001, VR-IC-DIS-002
  - **WF**: WF-IC-DIS-001
  - **Swagger**: `/commission/disbursements/*` endpoints
  - **DB**: `commission_disbursements`, `disbursement_cheques`, `disbursement_eft` tables

- [ ] **5.1.2** Pattern decisions
  - Use Temporal? → YES (10-day SLA, PFMS integration with retries)
  - Use Workflow State? → YES (cache agent bank details)
  - Use pgx.Batch? → YES (update multiple disbursement records)

- [ ] **5.1.3** Create domain models
  - `core/domain/disbursement.go` - Disbursement, ChequeDetails, EFTDetails

- [ ] **5.1.4** Create repository
  - `repo/postgres/disbursement.go`:
    - `CreateDisbursement(ctx, disbursement)`
    - `GetDisbursementsByAgentID(ctx, agentID)`
    - `UpdateDisbursementStatus(ctx, id, status)`
    - `GetPendingDisbursements(ctx)` - For batch processing
    - `RecordChequeDetails(ctx, details)`
    - `RecordEFTDetails(ctx, details)` - BR-IC-COM-008

- [ ] **5.1.5** Create Temporal workflow
  - `workflows/disbursement_workflow.go`:
    - `DisbursementWorkflow(ctx, finalStatementID)` - SLA: 10 working days
  - `workflows/activities/disbursement_activities.go`:
    - `ValidateBankDetailsActivity` - VR-IC-DIS-002 (for EFT)
    - `CreateDisbursementRecordsActivity`
    - `ProcessChequeActivity` - Generate cheque number
    - `ProcessEFTActivity` - Call PFMS API with retries
    - `UpdateAccountingActivity` - Post to GL
    - `NotifyAgentActivity` - Send disbursement confirmation

- [ ] **5.1.6** Create handler
  - `handler/disbursement.go`:
    - `POST /commission/disbursements/trigger` - Start disbursement workflow
    - `GET /commission/disbursements` - List disbursements
    - `GET /commission/disbursements/{disbursementId}` - Get details
    - `PUT /commission/disbursements/{disbursementId}/status` - Update status
    - `GET /commission/disbursements/agent/{agentId}` - Agent's disbursements
    - `POST /commission/disbursements/{disbursementId}/retry` - Retry failed
    - `GET /commission/disbursements/pending` - Pending disbursements
    - `POST /commission/disbursements/bulk-process` - Bulk processing

- [ ] **5.1.7** Unit tests
  - Test cheque generation
  - Test EFT processing
  - Test PFMS retry logic

**Traceability**:
- FR-IC-DIS-001, FR-IC-DIS-002, FR-IC-DIS-003, FR-IC-DIS-004
- BR-IC-COM-008, BR-IC-COM-010, BR-IC-COM-011
- VR-IC-DIS-001, VR-IC-DIS-002
- WF-IC-DIS-001

---

#### Module 5.2: Disbursement Tracking & SLA (Days 35-37)
**APIs**: 4 endpoints

##### Tasks:
- [ ] **5.2.1** Read and analyze requirements
  - **FR**: FR-IC-DIS-005
  - **BR**: BR-IC-COM-011 (10-day SLA)
  - **Swagger**: `/commission/disbursements/sla/*` endpoints

- [ ] **5.2.2** Create repository methods
  - `repo/postgres/disbursement.go`:
    - `GetSLAReport(ctx, filters)` - Raw SQL for SLA calculation
    - `GetOverdueDisbursements(ctx)` - BR-IC-COM-011

- [ ] **5.2.3** Create handler
  - `handler/disbursement_sla.go`:
    - `GET /commission/disbursements/sla/dashboard` - SLA dashboard
    - `GET /commission/disbursements/sla/overdue` - Overdue list
    - `GET /commission/disbursements/sla/report` - SLA report
    - `POST /commission/disbursements/sla/escalate` - Escalate overdue

- [ ] **5.2.4** Unit tests
  - Test SLA calculation logic
  - Test overdue detection

**Traceability**:
- FR-IC-DIS-005
- BR-IC-COM-011

---

### PHASE 6: Commission History & Search (Week 6-7: Days 38-40)
**Duration**: 3 days
**APIs**: 8 endpoints
**Goal**: Commission inquiry and reporting

#### Tasks:
- [ ] **6.1** Read and analyze requirements
  - **FR**: FR-IC-COM-009 to FR-IC-COM-012
  - **Swagger**: `/commission/history/*` endpoints

- [ ] **6.2** Pattern decisions
  - Use Temporal? → NO (read-only queries)
  - Use pgx.Batch? → YES (fetch commission + statement + disbursement)
  - Raw SQL for complex search queries

- [ ] **6.3** Create repository
  - `repo/postgres/commission_history.go`:
    - `SearchCommissionHistory(ctx, filters)` - Complex search with joins
    - `GetCommissionDetail(ctx, commissionID)` - Use pgx.Batch
    - `GetAgentCommissionSummary(ctx, agentID, period)` - Aggregations
    - `ExportCommissionHistory(ctx, filters)` - Export to Excel

- [ ] **6.4** Create handler
  - `handler/commission_history.go`:
    - `POST /commission/history/search` - Advanced search
    - `GET /commission/history/{commissionId}` - Get detail
    - `GET /commission/history/agent/{agentId}` - Agent history
    - `GET /commission/history/agent/{agentId}/summary` - Summary
    - `POST /commission/history/export` - Export to Excel/PDF
    - `GET /commission/history/statistics` - Statistics dashboard
    - `GET /commission/history/trends` - Trends analysis
    - `POST /commission/history/compare` - Period comparison

- [ ] **6.5** Unit tests
  - Test search with complex filters
  - Test export functionality
  - Test aggregation queries

**Traceability**:
- FR-IC-COM-009, FR-IC-COM-010, FR-IC-COM-011, FR-IC-COM-012

---

### PHASE 7: Commission Clawback (Week 7: Days 41-44)
**Duration**: 4 days
**APIs**: 7 endpoints
**Goal**: Manage commission recovery for policy cancellations

#### Tasks:
- [ ] **7.1** Read and analyze requirements
  - **FR**: FR-IC-CLAW-001 to FR-IC-CLAW-004
  - **BR**: BR-IC-CLAW-001 to BR-IC-CLAW-004
  - **VR**: VR-IC-CLAW-001, VR-IC-CLAW-002
  - **WF**: WF-IC-CLAW-001
  - **Swagger**: `/commission/clawback/*` endpoints
  - **DB**: `commission_clawback`, `clawback_recovery` tables

- [ ] **7.2** Pattern decisions
  - Use Temporal? → YES (graduated recovery over months)
  - Use Workflow State? → YES (cache agent, commission data)
  - Use pgx.Batch? → YES (update multiple records)

- [ ] **7.3** Create domain models
  - `core/domain/clawback.go` - Clawback, Recovery

- [ ] **7.4** Create repository
  - `repo/postgres/clawback.go`:
    - `CreateClawback(ctx, clawback)`
    - `GetClawbacksByAgentID(ctx, agentID)`
    - `RecordRecovery(ctx, recovery)` - BR-IC-CLAW-003
    - `GetPendingRecoveries(ctx)` - For batch processing

- [ ] **7.5** Create Temporal workflow
  - `workflows/clawback_workflow.go`:
    - `ClawbackWorkflow(ctx, policyID)` - Graduated recovery
  - `workflows/activities/clawback_activities.go`:
    - `CalculateClawbackActivity` - BR-IC-CLAW-001
    - `CreateRecoveryScheduleActivity` - BR-IC-CLAW-003
    - `DeductFromCommissionActivity`
    - `UpdateAccountingActivity`

- [ ] **7.6** Create handler
  - `handler/clawback.go`:
    - `POST /commission/clawback/trigger` - Trigger clawback
    - `GET /commission/clawback` - List clawbacks
    - `GET /commission/clawback/{clawbackId}` - Get details
    - `GET /commission/clawback/agent/{agentId}` - Agent clawbacks
    - `POST /commission/clawback/{clawbackId}/adjust` - Manual adjustment
    - `GET /commission/clawback/{clawbackId}/recovery-schedule` - Schedule
    - `POST /commission/clawback/bulk-process` - Bulk processing

- [ ] **7.7** Unit tests
  - Test clawback calculation
  - Test graduated recovery logic

**Traceability**:
- FR-IC-CLAW-001 to FR-IC-CLAW-004
- BR-IC-CLAW-001 to BR-IC-CLAW-004
- WF-IC-CLAW-001

---

### PHASE 8: Suspense Account Management (Week 8: Days 45-48)
**Duration**: 4 days
**APIs**: 11 endpoints
**Goal**: Handle disputed/held commissions

#### Tasks:
- [ ] **8.1** Read and analyze requirements
  - **FR**: FR-IC-SUSP-001 to FR-IC-SUSP-005
  - **BR**: BR-IC-SUSPENSE-001 to BR-IC-SUSPENSE-003
  - **VR**: VR-IC-SUSP-001
  - **WF**: WF-IC-SUSP-001
  - **Swagger**: `/commission/suspense/*` endpoints
  - **DB**: `commission_suspense_accounts` table

- [ ] **8.2** Pattern decisions
  - Use Temporal? → YES (aging and retry logic)
  - Use Workflow State? → YES
  - Use pgx.Batch? → YES

- [ ] **8.3** Create domain models
  - `core/domain/suspense.go` - SuspenseAccount

- [ ] **8.4** Create repository
  - `repo/postgres/suspense.go`:
    - `CreateSuspenseAccount(ctx, suspense)`
    - `GetSuspenseAccounts(ctx, filters)`
    - `UpdateSuspenseStatus(ctx, id, status)`
    - `GetAgingReport(ctx)` - BR-IC-SUSPENSE-003
    - `RetryDisbursement(ctx, id)`

- [ ] **8.5** Create Temporal workflow
  - `workflows/suspense_workflow.go`:
    - `SuspenseRetryWorkflow(ctx, suspenseID)` - Retry logic
  - `workflows/activities/suspense_activities.go`:
    - `ValidateDataActivity` - VR-IC-SUSP-001
    - `RetryDisbursementActivity`
    - `EscalateActivity` - After 30 days

- [ ] **8.6** Create handler
  - `handler/suspense.go`:
    - `POST /commission/suspense/create` - Create suspense
    - `GET /commission/suspense` - List suspense accounts
    - `GET /commission/suspense/{suspenseId}` - Get details
    - `PUT /commission/suspense/{suspenseId}/resolve` - Resolve
    - `POST /commission/suspense/{suspenseId}/retry` - Retry disbursement
    - `GET /commission/suspense/aging` - Aging report
    - `GET /commission/suspense/dashboard` - Dashboard
    - `POST /commission/suspense/{suspenseId}/escalate` - Escalate
    - `POST /commission/suspense/bulk-resolve` - Bulk resolve
    - `GET /commission/suspense/statistics` - Statistics
    - `POST /commission/suspense/export` - Export

- [ ] **8.7** Unit tests
  - Test retry logic
  - Test aging calculation

**Traceability**:
- FR-IC-SUSP-001 to FR-IC-SUSP-005
- BR-IC-SUSPENSE-001 to BR-IC-SUSPENSE-003
- WF-IC-SUSP-001

---

### PHASE 9: Lookup & Reference Data (Week 8: Days 49-50)
**Duration**: 2 days
**APIs**: 10 endpoints
**Goal**: Provide reference data for dropdowns and lookups

#### Tasks:
- [ ] **9.1** Read and analyze requirements
  - **Swagger**: `/lookup/*` endpoints
  - **DB**: `ref_circles`, `ref_divisions`, `ref_product_plans` tables

- [ ] **9.2** Create domain models
  - `core/domain/reference.go` - Circle, Division, ProductPlan

- [ ] **9.3** Create repository
  - `repo/postgres/reference.go`:
    - `GetCircles(ctx)`
    - `GetDivisionsByCircle(ctx, circleID)`
    - `GetCoordinatorsByCircle(ctx, circleID)`
    - `GetProductPlans(ctx, productType)`
    - `GetAgentTypes(ctx)`

- [ ] **9.4** Create handler
  - `handler/lookup.go`:
    - `GET /lookup/circles` - All circles
    - `GET /lookup/circles/{circleId}/divisions` - Divisions by circle
    - `GET /lookup/circles/{circleId}/coordinators` - Coordinators by circle
    - `GET /lookup/divisions/{divisionId}/coordinators` - Coordinators by division
    - `GET /lookup/product-plans` - Product plans
    - `GET /lookup/product-plans/{productType}` - Plans by product type
    - `GET /lookup/agent-types` - Agent types
    - `GET /lookup/commission-types` - Commission types
    - `GET /lookup/payment-modes` - Payment modes
    - `GET /lookup/all` - All lookups in one call

- [ ] **9.5** Unit tests

---

### PHASE 10: Export, Reports & Webhooks (Week 9: Days 51-53)
**Duration**: 3 days
**APIs**: 10 endpoints
**Goal**: Reporting and external integration

#### Module 10.1: Export & Reports (Days 51-52)
**APIs**: 6 endpoints

##### Tasks:
- [ ] **10.1.1** Read requirements
  - **Swagger**: `/export/*` endpoints

- [ ] **10.1.2** Create handler
  - `handler/export.go`:
    - `POST /export/agents` - Export agents to Excel
    - `POST /export/commissions` - Export commissions to Excel
    - `POST /export/trial-statements` - Export trial statements
    - `POST /export/disbursements` - Export disbursements
    - `POST /export/clawback` - Export clawback data
    - `GET /export/{jobId}/status` - Check export job status

- [ ] **10.1.3** Unit tests

---

#### Module 10.2: Integration Webhooks (Day 53)
**APIs**: 4 endpoints

##### Tasks:
- [ ] **10.2.1** Read requirements
  - **Swagger**: `/webhooks/*` endpoints

- [ ] **10.2.2** Create handler
  - `handler/webhooks.go`:
    - `POST /webhooks/pfms/payment-status` - PFMS payment callback
    - `POST /webhooks/accounting/voucher-confirmation` - Accounting confirmation
    - `POST /webhooks/policy/cancellation` - Policy cancellation trigger
    - `POST /webhooks/hrms/employee-update` - HRMS sync

- [ ] **10.2.3** Unit tests
  - Test webhook signature verification
  - Test retry handling

---

### PHASE 11: Integration Testing (Week 9-10: Days 54-58)
**Duration**: 5 days
**Goal**: End-to-end testing of all workflows

#### Tasks:
- [ ] **11.1** Create integration test suite
  - Test full agent onboarding flow
  - Test commission batch end-to-end (calculation → trial → final → disbursement)
  - Test license reminder workflow
  - Test clawback recovery workflow
  - Test suspense retry workflow

- [ ] **11.2** Performance testing
  - Test batch processing with 100,000 commission records
  - Test pgx.Batch performance
  - Test Temporal workflow scalability
  - Verify 6-hour SLA for commission batch

- [ ] **11.3** Load testing
  - Concurrent user testing
  - Database connection pool optimization
  - API rate limiting verification

- [ ] **11.4** Security testing
  - JWT authentication
  - Input validation
  - SQL injection prevention
  - XSS prevention

---

### PHASE 12: Documentation & Deployment (Week 10: Days 59-60)
**Duration**: 2 days
**Goal**: Production readiness

#### Tasks:
- [ ] **12.1** Documentation
  - Update README.md with setup instructions
  - Document API endpoints (Swagger auto-generated)
  - Create deployment guide
  - Create troubleshooting guide
  - Document temporal workflow monitoring

- [ ] **12.2** Database migrations
  - Finalize all migration scripts
  - Test migration rollback
  - Create seed data scripts

- [ ] **12.3** Configuration
  - Setup production configs
  - Configure Temporal workers
  - Configure database connection pools
  - Setup monitoring and alerting

- [ ] **12.4** Deployment preparation
  - Create Dockerfile
  - Create docker-compose.yml
  - Create Kubernetes manifests (if applicable)
  - Setup CI/CD pipeline

---

## Summary Statistics

| Metric | Count |
|--------|-------|
| **Total Phases** | 12 |
| **Total APIs** | 105 |
| **Temporal Workflows** | 8 |
| **Database Tables** | 25+ |
| **Estimated Duration** | 10 weeks |
| **Unit Test Files** | 30+ |

---

## Critical Success Factors

1. **Strict Adherence to Analysis Document** - All FR/BR/VR/WF must be traceable in code comments
2. **Database Optimization** - pgx.Batch usage mandatory for all multi-query operations
3. **Temporal Workflow State** - Cache frequently accessed data to reduce DB round trips
4. **SLA Compliance** - Monitor and enforce all SLAs (6-hour batch, 7-day approval, 10-day disbursement)
5. **Error Handling** - Implement all error codes (IC-ERR-001 to IC-ERR-041)
6. **Unit Testing** - 80%+ code coverage required
7. **Traceability** - Every handler/repository/workflow must have FR/BR/VR comments

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| **Scope Creep** | Strict phase boundaries, defer non-critical features |
| **Performance Issues** | Early performance testing in Phase 4 |
| **Integration Failures** | Mock external services, test independently |
| **Temporal Complexity** | Follow workflow state optimization patterns |
| **Database Bottlenecks** | Use pgx.Batch, optimize queries early |

---

## Next Steps

1. Review and approve this plan
2. Ask clarifying questions
3. Create `context.md` to track progress
4. Begin Phase 0: Project Foundation

---

**Plan Created**: 2026-01-28
**Last Updated**: 2026-01-28
**Version**: 1.0
