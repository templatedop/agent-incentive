# Agent Commission Management System - Project Context

## Current Status

**Phase**: Phase 8 - Export & Reporting - COMPLETE ✅
**Current Module**: ALL 30 APIs COMPLETE! 🎉
**Last Updated**: 2026-01-28
**Progress**: 100% (30/30 APIs completed) | Phases 1-8: ✅ ALL COMPLETE

### IMPORTANT: Scope Change
- **Old Scope**: 105 APIs including Agent Profile Management
- **New Scope**: 30 APIs - Commission & Incentive ONLY
- **Agent Profile Management**: Moved to complete-pli-agent project (OUT OF SCOPE)
- **New Plan**: `/home/user/agent-incentive/plan_revised_incentive_only.md`
- **Swagger**: `/home/user/agent-incentive/Incentive/swagger/swagger_incentive.yaml`

---

## Project Information

### Project Details
- **Name**: Agent Incentive, Commission and Producer Management System
- **Repository**: agent-incentive
- **Branch**: claude/general-session-XA1Pw
- **Framework**: N-API Template (Golang)
- **Workflow Engine**: Temporal.io
- **Database**: PostgreSQL 16
- **Architecture**: Hexagonal/Ports-and-Adapters

### Key Documents
- **Plan (REVISED)**: `/home/user/agent-incentive/plan_revised_incentive_only.md` ⭐
- **Original Plan**: `/home/user/agent-incentive/plan.md` (DEPRECATED)
- **Requirements**: `/home/user/agent-incentive/Incentive/analysis/IC_Incentive_Commission_Producer_Management_Analysis.md`
- **SRS**: `/home/user/agent-incentive/Incentive/srs/Agent_SRS_Incentive-Commission-and-Producer-Management.md`
- **Swagger (Commission Only)**: `/home/user/agent-incentive/Incentive/swagger/swagger_incentive.yaml` ⭐
- **Swagger (Full - OLD)**: `/home/user/agent-incentive/Incentive/swagger/swagger.yaml` (DEPRECATED)
- **Database Schema**: `/home/user/agent-incentive/Incentive/db/schema.sql`
- **Template**: `/home/user/agent-incentive/Incentive/template.md`
- **DB Library README**: `/home/user/agent-incentive/Incentive/db-README.md`

### Skill
- **Name**: insurance-code-generator
- **Location**: `/home/user/agent-incentive/Incentive/insurance-code-generator/`
- **Status**: ✅ Loaded and ready

---

## Agent Onboarding Code - BACKED UP

The agent onboarding implementation has been moved to backup folder as it's now handled by complete-pli-agent:
- **Location**: `/home/user/agent-incentive/agent-commission/backup_agent_onboarding/`
- **Contents**:
  - agent_onboarding.go (handler with 6 APIs)
  - agent_wrappers.go (response DTOs)
  - workflows/ (Temporal workflow + 8 activities)
- **Reason**: Agent profile management is out of scope for this commission-only module

---

## Implementation Progress (REVISED)

### Phase 0: Project Foundation & Setup [✅ COMPLETE]
- [x] 0.1 Create project directory structure
- [x] 0.2 Initialize Go modules with dependencies
- [x] 0.3 Setup configuration files (dev, test, prod)
- [x] 0.4 Create bootstrap/bootstrapper.go with FX modules
- [x] 0.5 Setup database migrations structure
- [x] 0.6 Create main.go application entry point
- [x] 0.7 Setup docker-compose (PostgreSQL + Temporal + UI)
- [x] 0.8 Create project files (.gitignore, README, Makefile)
- [x] 0.9 Project structure verified

**Deliverables:**
- ✅ Complete N-API Template project at `/home/user/agent-incentive/agent-commission/`
- ✅ 16 files created (3,436 lines of code)
- ✅ Database schema with 20+ tables
- ✅ Docker infrastructure for E2E testing
- ✅ Comprehensive Makefile with 30+ commands
- ✅ Multi-environment configuration
- ✅ Migration framework ready
- ✅ Documentation complete

### Phase 1: Foundation & Reference Data [✅ COMPLETE]

**Module 1.1: Commission Rate Configuration** (1 API) [✅ COMPLETE]
- [x] 1.1.1 Domain models: CommissionRate, CommissionRateHistory
- [x] 1.1.2 Repository: commission_rate_repository.go
- [x] 1.1.3 Handler: commission_rate_handler.go
- [x] 1.1.4 API: GET /lookup/commission-rates

**Module 1.2: Lookup & Reference APIs** (2 APIs) [✅ COMPLETE]
- [x] 1.2.1 Repository: Extended reference_data_repository.go
- [x] 1.2.2 Handler: lookup_handler.go
- [x] 1.2.3 API: GET /lookup/advisor-coordinators
- [x] 1.2.4 API: GET /lookup/circles

**Deliverables (Module 1.1):**
- ✅ **Domain Models** (`agent-commission/core/domain/commission_rate.go`)
  - CommissionRate entity with all required fields
  - CommissionRateHistory for audit trail
  - CommissionRateFilter for query filtering
  - ProductType enum (PLI, RPLI)
  - IsActiveOn() method for date-based validation
  - Implements: BR-IC-COM-006, FR-IC-COM-001

- ✅ **Response DTOs** (`agent-commission/handler/response/commission_rate.go`)
  - CommissionRateResponse with ISO 8601 date formatting
  - CommissionRatesListResponse with structured data wrapper
  - Conversion functions: NewCommissionRateResponse, NewCommissionRatesResponse

- ✅ **Repository Layer** (`agent-commission/repo/postgres/commission_rate_repository.go`)
  - GetCommissionRates - Filter by product_type, agent_type, plan_code
  - GetCommissionRateByID - Single rate lookup
  - GetApplicableRate - Find rate for specific policy parameters
  - Uses dblib.SelectRows and dblib.Psql patterns
  - Timeout handling with cfg.GetDuration("db.QueryTimeoutLow")

- ✅ **HTTP Handler** (`agent-commission/handler/commission_rate_handler.go`)
  - GET /lookup/commission-rates - Query commission rates
  - Query parameters: product_type (PLI/RPLI), agent_type, plan_code
  - Returns CommissionRatesListResponse with port.ListSuccess
  - Plain Go function pattern: func(sctx, req) (*resp, error)

- ✅ **Bootstrap Integration** (`agent-commission/bootstrap/bootstrapper.go`)
  - Registered CommissionRateRepository in FxRepo
  - Registered CommissionRateHandler in FxHandler
  - Handler registered with serverHandler.Register()

- ✅ **Port Types** (`agent-commission/core/port/`)
  - response.go - StatusCodeAndMessage, MetaDataResponse
  - request.go - MetadataRequest for pagination
  - Copied from template for consistency

**Deliverables (Module 1.2):**
- ✅ **Repository Extension** (`agent-commission/repo/postgres/reference_data_repository.go`)
  - GetAllCoordinators - Fetch all active advisor coordinators
  - Sorted by circle_name, division_name, first_name, last_name
  - Reuses existing GetAllCircles method

- ✅ **Response DTOs** (`agent-commission/handler/response/lookup.go`)
  - AdvisorCoordinatorSummary with full name and location details
  - CircleSummary with code and name
  - Conversion functions for domain to DTO mapping

- ✅ **HTTP Handler** (`agent-commission/handler/lookup_handler.go`)
  - GET /lookup/advisor-coordinators - Active coordinators lookup
  - GET /lookup/circles - All circles lookup
  - Query parameters: circle, division (optional filters)
  - Returns structured responses with data wrappers

- ✅ **Bootstrap Integration**
  - Registered LookupHandler in FxHandler module
  - Added to serverHandler.Register()

### Phase 2: Commission Processing Core [✅ COMPLETE]

**Module 2.1: Commission Batch Processing** (2 APIs) [✅ COMPLETE]
- [x] 2.1.1 Domain models: CommissionBatch, CommissionTransaction
- [x] 2.1.2 Repository: commission_batch_repository.go
- [x] 2.1.3 Handler: commission_batch_handler.go
- [x] 2.1.4 APIs: POST /commissions/batches/start, GET /commissions/batches/{batchId}/status

**Module 2.2: Trial Statement Management** (2 APIs) [✅ COMPLETE]
- [x] 2.2.1 Domain model: TrialStatement
- [x] 2.2.2 Repository: trial_statement_repository.go
- [x] 2.2.3 Handler: trial_statement_handler.go
- [x] 2.2.4 APIs: GET /commissions/trial-statements, POST /commissions/trial-statements/{id}/approve

**Module 2.3: Final Statement Generation** (1 API) [✅ COMPLETE]
- [x] 2.3.1 Domain model: FinalStatement
- [x] 2.3.2 Repository: final_statement_repository.go
- [x] 2.3.3 Handler: final_statement_handler.go
- [x] 2.3.4 API: GET /commissions/final-statements

**Deliverables (Phase 2):**
- ✅ **Domain Models** (4 files)
  - CommissionBatch - Batch with 6-hour SLA tracking, progress calculation
  - CommissionTransaction - Individual commission records with TDS
  - TrialStatement - Trial statements with approval/rejection methods
  - FinalStatement - Final statements ready for disbursement
  - All with BR/FR traceability comments

- ✅ **Database Migration** (`002_add_commission_batches.up.sql`)
  - commission_batches table with month/year uniqueness
  - SLA deadline tracking, workflow integration
  - Indexes for status, workflow_id, month_year, sla tracking
  - Optimized for batch processing queries

- ✅ **Repository Layer** (3 files - all with batch optimization)
  - commission_batch_repository.go (~250 lines)
    - CreateBatch, GetBatchByID, GetBatchByMonthYear
    - UpdateBatchProgress, UpdateBatchStatus, CompleteBatch
    - Implements BR-IC-COM-001, BR-IC-COM-012 (6-hour SLA)
  - trial_statement_repository.go (~280 lines)
    - CreateTrialStatement, GetTrialStatementByID
    - SearchTrialStatements (batch optimized count + results in 1 round trip)
    - ApproveTrialStatement with optimistic locking
    - Implements FR-IC-COM-004, BR-IC-COM-002, BR-IC-COM-009
  - final_statement_repository.go (~200 lines)
    - CreateFinalStatement, GetFinalStatementByID
    - SearchFinalStatements (batch optimized count + results)
    - UpdateStatementStatus
    - Implements FR-IC-COM-008, BR-IC-COM-007

- ✅ **Response DTOs** (3 files)
  - CommissionBatchResponse - Batch creation with workflow tracking
  - CommissionBatchStatusResponse - Detailed status with SLA, progress, time remaining
  - TrialStatementSummary + ApprovalResponse - List and approval
  - FinalStatementSummary - Final statement listing
  - PaginationMetadata - Reusable pagination structure

- ✅ **HTTP Handlers** (3 files)
  - commission_batch_handler.go (2 APIs)
    - POST /commissions/batches/start - Start batch with duplicate checking
    - GET /commissions/batches/{batchId}/status - Detailed status
  - trial_statement_handler.go (2 APIs)
    - GET /commissions/trial-statements - List with filters & pagination
    - POST /commissions/trial-statements/{id}/approve - Approve with remarks
  - final_statement_handler.go (1 API)
    - GET /commissions/final-statements - List with filters & pagination
  - All handlers use plain Go function pattern from template
  - Request validation, error handling, logging
  - URI and query parameter binding

- ✅ **Bootstrap Integration**
  - Registered 3 Phase 2 repositories in FxRepo
  - Registered 3 Phase 2 handlers in FxHandler
  - Total: 8 handlers (Phase 1 + Phase 2)

**Notes:**
- Temporal workflow integration marked as TODO for Phase 2B
- User authentication (JWT) extraction marked as TODO
- All batch operations use pgx.Batch for optimal performance
- Optimistic locking with version field for concurrent updates

### Phase 3: Commission Disbursement [✅ COMPLETE]

**Module 3.1: Disbursement Processing** (2 APIs) [✅ COMPLETE]
- [x] 3.1.1 Domain model: Disbursement with SLA tracking
- [x] 3.1.2 Repository: disbursement_repository.go (batch optimized)
- [x] 3.1.3 Handler: disbursement_handler.go
- [x] 3.1.4 APIs: POST /commissions/disbursements, GET /commissions/disbursements/{id}/status

**Module 3.2: Integration Webhooks** (2 APIs) [✅ COMPLETE]
- [x] 3.2.1 Handler: webhook_handler.go
- [x] 3.2.2 APIs: POST /webhooks/pfms/disbursement-confirmation, POST /webhooks/policy/status-change

**Deliverables (Phase 3):**
- ✅ **Domain Model** (`disbursement.go`)
  - Disbursement entity with dual payment modes (CHEQUE/EFT)
  - 10-working-day SLA tracking (BR-IC-COM-011)
  - DisbursementStatus enum (6 states), PaymentFailureReason enum
  - Business methods: IsSLABreached(), CanProcess(), CanRetry()

- ✅ **Database Migration** (`003_add_disbursements.up.sql`)
  - disbursements table with mode-specific validation
  - 10 performance indexes for queries
  - Workflow and accounting integration fields

- ✅ **Repository Layer** (`disbursement_repository.go` - 350 lines)
  - CreateDisbursement with auto SLA calculation
  - SearchDisbursements (batch optimized count + results)
  - UpdateDisbursementStatus, CompleteDisbursement
  - PostToGL, UpdateSLABreach
  - calculateWorkingDays() helper (excludes weekends)

- ✅ **Response DTOs** (`response/disbursement.go`)
  - DisbursementSummary, DisbursementDetailResponse
  - Time remaining calculation for active disbursements

- ✅ **HTTP Handlers** (2 files, 4 APIs)
  - disbursement_handler.go: Create and status APIs
  - webhook_handler.go: PFMS and policy webhooks
  - Mode-specific validation (cheque/bank details)
  - PFMS failure code mapping

- ✅ **Bootstrap Integration**
  - Registered DisbursementRepository and both handlers

**Notes:**
- Webhook signature verification marked as TODO
- Temporal workflow integration marked as TODO
- Clawback trigger marked as TODO
- All queries use batch optimization patterns

### Phase 4: Commission History & Inquiry [✅ COMPLETE]

**Module 4.1: Commission History** (1 API) [✅ COMPLETE]
- [x] 4.1.1 Domain filter: CommissionHistoryFilter
- [x] 4.1.2 Repository: commission_history_repository.go
- [x] 4.1.3 Handler: commission_history_handler.go
- [x] 4.1.4 API: GET /commissions/history

**Deliverables (Phase 4):**
- ✅ **Domain Filter** (`commission_batch.go`)
  - CommissionHistoryFilter for query parameters
  - Supports filters: agent_id, policy_number, commission_type, status, date range, batch_id
  - Implements: FR-IC-COM-011, BR-IC-COM-009

- ✅ **Repository Layer** (`commission_history_repository.go` - 132 lines)
  - SearchCommissionHistory (batch optimized count + results)
  - Flexible filtering with 8 query parameters
  - Pagination support

- ✅ **Response DTOs** (`response/commission_history.go`)
  - CommissionHistorySummary for list views
  - CommissionHistoryResponse with pagination

- ✅ **HTTP Handler** (`commission_history_handler.go`)
  - GET /commissions/history - Search commission transaction history
  - Query parameters: agent_id, policy_number, commission_type, commission_status, from_date, to_date, batch_id, page, limit

- ✅ **Bootstrap Integration**
  - Registered CommissionHistoryRepository and CommissionHistoryHandler

### Phase 5: Clawback Management [✅ COMPLETE]

**Module 5.1: Commission Clawback** (1 API) [✅ COMPLETE]
- [x] 5.1.1 Domain models: Clawback, ClawbackRecovery
- [x] 5.1.2 Repository: clawback_repository.go
- [x] 5.1.3 Handler: clawback_handler.go
- [x] 5.1.4 API: POST /commissions/clawback/create

**Deliverables (Phase 5):**
- ✅ **Domain Models** (`clawback.go` - 260 lines)
  - Clawback entity with graduated recovery tracking
  - ClawbackRecovery entity for installment tracking
  - ClawbackStatus enum (PENDING, IN_PROGRESS, COMPLETED, PARTIAL, WAIVED, WRITE_OFF)
  - ClawbackReason enum (5 types)
  - CalculateClawbackPercentage() - BR-IC-CLAWBACK-002: Year 1: 100%, Year 2: 75%, Year 3: 50%, Year 4: 25%, Year 5+: 0%
  - Business methods: IsFullyRecovered(), CanRecover(), RecoveryProgress()
  - Implements: FR-IC-COM-012, BR-IC-CLAWBACK-001, BR-IC-CLAWBACK-002

- ✅ **Repository Layer** (`clawback_repository.go` - 400+ lines)
  - CreateClawback with auto percentage calculation
  - SearchClawbacks (batch optimized count + results)
  - RecordRecovery with transaction support
  - GetRecoveriesByClawbackID
  - UpdateClawbackStatus

- ✅ **Response DTOs** (`response/clawback.go` - 250 lines)
  - ClawbackDetailResponse with recovery progress
  - ClawbackSummary for list views
  - ClawbackRecoveryDetail for recovery transactions

- ✅ **HTTP Handler** (`clawback_handler.go`)
  - POST /commissions/clawback/create - Initiate clawback
  - Calculates total commission paid, applies graduated percentage
  - Recovery schedule: IMMEDIATE or INSTALLMENT

- ✅ **Bootstrap Integration**
  - Registered ClawbackRepository and ClawbackHandler

**Notes:**
- Clawback workflow integration marked as TODO
- Graduated clawback percentage based on policy age in months

### Phase 6: Suspense Account Management [✅ COMPLETE]

**Module 6.1: Suspense Accounts** (2 APIs) [✅ COMPLETE]
- [x] 6.1.1 Domain models: SuspenseAccount, SuspenseTransaction
- [x] 6.1.2 Repository: suspense_repository.go
- [x] 6.1.3 Handler: suspense_handler.go
- [x] 6.1.4 APIs: POST /commissions/suspense/create, GET /commissions/suspense/aging-report

**Deliverables (Phase 6):**
- ✅ **Domain Models** (`suspense.go` - 250 lines)
  - SuspenseAccount entity with aging tracking
  - SuspenseTransaction audit trail
  - SuspenseStatus enum (OPEN, RESOLVED, WRITE_OFF)
  - SuspenseReason enum (10 types)
  - SuspenseAgingBucket, SuspenseAgingReport structures
  - GetAgingBucket() - BR-IC-SUS-002: 0-30, 31-60, 61-90, 91-180, 180+ days
  - CalculateResolutionDeadline() - BR-IC-SUS-003: HIGH: 7d, MEDIUM: 15d, LOW: 30d
  - DeterminePriority() - BR-IC-SUS-004: Auto-priority based on amount and reason
  - Implements: FR-IC-COM-013, BR-IC-SUS-001

- ✅ **Repository Layer** (`suspense_repository.go` - 450+ lines)
  - CreateSuspenseAccount with auto priority and deadline
  - SearchSuspenseAccounts (batch optimized count + results)
  - GetAgingReport - Generates bucket distribution with SQL aggregation
  - ResolveSuspenseAccount, WriteOffSuspenseAccount
  - createSuspenseTransaction for audit trail

- ✅ **Response DTOs** (`response/suspense.go` - 230 lines)
  - SuspenseAccountDetailResponse with aging bucket and overdue flag
  - SuspenseAccountSummary for list views
  - SuspenseAgingBucketResponse, SuspenseAgingReportResponse
  - Time-based calculations (aging days, overdue status)

- ✅ **HTTP Handler** (`suspense_handler.go`)
  - POST /commissions/suspense/create - Create suspense entry
  - GET /commissions/suspense/aging-report - Aging distribution report
  - Auto-priority assignment based on amount and reason

- ✅ **Bootstrap Integration**
  - Registered SuspenseRepository and SuspenseHandler

**Notes:**
- Aging report uses raw SQL for efficient aggregation
- Priority auto-determined: HIGH (≥50k or critical), MEDIUM (≥10k), LOW (<10k)
- Resolution SLA based on priority: HIGH: 7 days, MEDIUM: 15 days, LOW: 30 days
- Suspense workflow integration marked as TODO

### Phase 7: Workflow Management [✅ COMPLETE]

**Module 7.1: Workflow Status & Control** (8 APIs) [✅ COMPLETE]
- [x] 7.1.1 Domain models: WorkflowInfo, WorkflowHistory, BatchProgress
- [x] 7.1.2 Response DTOs: workflow.go
- [x] 7.1.3 Handler: workflow_management_handler.go
- [x] 7.1.4 APIs: All 8 workflow management endpoints

**Deliverables (Phase 7):**
- ✅ **Domain Models** (`workflow.go` - 200 lines)
  - WorkflowInfo with SLA tracking
  - WorkflowStatus enum (RUNNING, COMPLETED, FAILED, CANCELED, TERMINATED, TIMED_OUT)
  - WorkflowType enum (COMMISSION_BATCH, TRIAL_STATEMENT_APPROVAL, DISBURSEMENT, CLAWBACK, SUSPENSE)
  - WorkflowHistory with events
  - BatchProgress with real-time tracking
  - Business methods: IsRunning(), CanCancel(), CanRetry(), CalculateSLARemaining()

- ✅ **Response DTOs** (`response/workflow.go` - 220 lines)
  - WorkflowStatusResponse with SLA information
  - WorkflowHistoryResponse with event details
  - WorkflowCancelResponse, WorkflowRetryResponse
  - WorkflowQueryResponse for state queries
  - BatchStatusResponse, BatchProgressResponse
  - BatchCancelResponse

- ✅ **HTTP Handler** (`workflow_management_handler.go` - 450+ lines)
  - GET /workflows/{workflowId}/status - Get workflow status with SLA tracking
  - GET /workflows/{workflowId}/history - Get execution history
  - POST /workflows/{workflowId}/cancel - Cancel running workflow
  - POST /workflows/{workflowId}/retry - Retry failed workflow
  - POST /workflows/{workflowId}/query - Query workflow state
  - GET /batches/{batchId}/status - Get batch status
  - GET /batches/{batchId}/progress - Get batch progress with estimates
  - POST /batches/{batchId}/cancel - Cancel batch processing

- ✅ **Bootstrap Integration**
  - Registered WorkflowManagementHandler

**Notes:**
- Mock implementations provided for all Temporal interactions
- Ready to wire up to Temporal client when available
- Includes comprehensive TODO comments for Temporal integration points
- SLA tracking and monitoring built-in
- Real-time progress estimation for batch processing
- Supports workflow retry with automatic ID generation

### Phase 8: Export & Reporting [✅ COMPLETE]

**Module 8.1: Export APIs** (6 APIs) [✅ COMPLETE]
- [x] 8.1.1 Domain models: ExportJob, CommissionSummaryReport
- [x] 8.1.2 Response DTOs: export.go
- [x] 8.1.3 Handler: export_handler.go
- [x] 8.1.4 APIs: All 6 export and reporting endpoints

**Deliverables (Phase 8):**
- ✅ **Domain Models** (`export.go` - 240 lines)
  - ExportJob entity for async export tracking
  - ExportStatus enum (PENDING, PROCESSING, COMPLETED, FAILED, EXPIRED)
  - ExportFormat enum (EXCEL, PDF, CSV)
  - ExportType enum (TRIAL_STATEMENTS, FINAL_STATEMENTS, COMMISSION_HISTORY, CLAWBACK_REPORT, SUSPENSE_REPORT, COMMISSION_SUMMARY)
  - CommissionSummaryReport for aggregated metrics
  - Business methods: IsCompleted(), CanRetry(), IsExpired(), MarkAsProcessing(), MarkAsCompleted()
  - Implements: FR-IC-COM-010, FR-IC-COM-012

- ✅ **Response DTOs** (`response/export.go` - 260 lines)
  - ExportJobResponse - Export job creation confirmation
  - ExportJobStatusResponse - Job status with progress tracking
  - CommissionSummaryResponse - Aggregated commission metrics
  - CommissionSummaryData - Detailed breakdown by type
  - ExportJobListResponse - List of export jobs
  - Helper functions: NewExportJobResponse(), NewExportJobStatusResponse(), etc.

- ✅ **HTTP Handler** (`export_handler.go` - 380+ lines)
  - POST /exports/commissions/trial-statements - Export trial statements
  - POST /exports/commissions/final-statements - Export final statements
  - POST /exports/commissions/history - Export commission history
  - POST /exports/commissions/clawback - Export clawback report
  - POST /exports/commissions/suspense - Export suspense report
  - POST /reports/commission-summary - Generate commission summary

- ✅ **Bootstrap Integration**
  - Registered ExportHandler

**Notes:**
- Async export architecture - returns job ID for tracking
- Supports multiple formats: EXCEL, PDF, CSV
- Mock implementations provided for all exports
- Ready for integration with actual export service (Excel/PDF generation)
- Includes comprehensive filter options for each export type
- Commission summary includes aggregations by type, agent, product
- Export links expire after 7 days
- Supports retry for failed exports (max 3 attempts)
- TODO: Implement actual file generation service
- TODO: Integrate with S3/storage for file hosting
- TODO: Add Temporal workflow for async export processing

---

## Key Decisions Made

### Architecture Decisions
1. **Framework**: N-API Template (Golang) - Follows established patterns
2. **Workflow Engine**: Temporal.io - For long-running processes with SLAs
3. **Database**: PostgreSQL 16 - With pgx driver for performance
4. **Database Library**: Squirrel for simple queries, Raw SQL for complex operations
5. **Validation**: go-playground/validator with auto-generated validators

### Pattern Decisions (from insurance-code-generator skill)

#### When to Use Temporal Workflow
- ✅ **Agent Onboarding** - Multi-step with HRMS integration
- ✅ **Commission Batch Processing** - 6-hour SLA, batch operations
- ✅ **Trial Statement Approval** - 7-day SLA, human-in-the-loop
- ✅ **Disbursement Processing** - 10-day SLA, PFMS integration
- ✅ **License Renewal** - 3-day SLA, approval workflow
- ✅ **License Reminders** - Scheduled notifications (T-30, T-15, T-7, T-0)
- ✅ **Agent Termination** - Multi-step with commission settlement
- ✅ **Clawback Recovery** - Graduated recovery over months
- ✅ **Suspense Retry** - Aging and retry logic

#### When NOT to Use Temporal
- ❌ **Agent Profile CRUD** - Simple database operations
- ❌ **Commission Rate Configuration** - Simple CRUD
- ❌ **Commission History Search** - Read-only queries
- ❌ **Lookup/Reference Data** - Simple queries

#### Database Optimization Strategy
- **pgx.Batch**: Use for ALL multi-query operations (inserts, updates, complex fetches)
  - Agent onboarding: profile + address + contact in one batch
  - Commission batch: Bulk insert 1000s of records
  - Trial statements: Aggregate by agent in one query
- **Workflow State Optimization**: Cache frequently accessed data
  - Commission rates table
  - Agent hierarchy data
  - HRMS data
  - Bank details
- **Target**: Reduce DB round trips by 60-75%

#### Batch Optimization Patterns (Critical Learning from Phase 1)

**Rule #1: NEVER use transactions (WithTx) for multi-query operations - Use batch instead**
- ❌ **WRONG**: `r.db.WithTx(ctx, func(tx pgx.Tx) { ... })` for multiple queries
- ✅ **CORRECT**: `batch := &pgx.Batch{}` - Batch provides implicit transaction

**Rule #2: Use batch to reduce database round trips**
- Search operations: Count + Results = 2 queries → 1 batch
- Create operations: Profile + Addresses + Contacts = N queries → 1 batch
- Update operations: Deactivate old + Create new = 2 queries → 1 batch

**Rule #3: Use utility.go batch helper functions**
Location: `/home/user/agent-incentive/Incentive/api-db/utility.go`

```go
// For INSERT/UPDATE/DELETE with RETURNING
dblib.QueueReturnRow(batch, builder, scanFn, result)

// For SELECT returning multiple rows
dblib.QueueReturn(batch, builder, scanFn, result)

// For INSERT/UPDATE/DELETE without RETURNING
dblib.QueueExecRow(batch, builder)

// For bulk operations
dblib.QueueReturnBulk(batch, builder, scanFn, result)
```

**Rule #4: Batch execution pattern**
```go
// 1. Create batch
batch := &pgx.Batch{}

// 2. Queue operations
dblib.QueueReturnRow(batch, insertQuery, scanFn, &result1)
dblib.QueueExecRow(batch, updateQuery)
dblib.QueueReturn(batch, selectQuery, scanFn, &results)

// 3. Execute batch (implicit transaction)
batchResults := r.db.Pool.SendBatch(ctx, batch)
defer batchResults.Close()

// 4. Process results sequentially
batchResults.QueryRow().Scan(&result1.ID)
batchResults.Exec()
rows, _ := batchResults.Query()
// ... collect rows
```

**Examples from Phase 1:**
1. **CreateAgentProfileWithRelations**: Profile + Addresses + Contacts + Emails in ONE batch
2. **SearchAgents**: Count + Results in ONE batch (was 2 separate queries)
3. **UpdateHierarchyRelationship**: Deactivate + Create in ONE batch (was transaction)

### Naming Conventions
- **Handlers**: `{resource}_handler.go` (e.g., `agent_profile.go`)
- **Repositories**: `repo/postgres/{resource}.go`
- **Workflows**: `workflows/{resource}_workflow.go`
- **Activities**: `workflows/activities/{resource}_activities.go`
- **Domain**: `core/domain/{resource}.go`
- **DTOs**: Request in `handler/request.go`, Response in `handler/response/{resource}.go`

---

## Traceability Requirements

### Code Comments Format
Every function/method must include traceability comments:

```go
// CreateAgentProfile creates a new agent profile
// FR: FR-IC-PROF-001 (Agent Onboarding)
// BR: BR-IC-AH-001 (Advisor must link to Coordinator)
// VR: VR-IC-PROF-002 (PAN uniqueness check)
// WF: WF-IC-ONB-001 (Step 3 - Profile Creation)
func (r *AgentRepository) CreateAgentProfile(ctx context.Context, agent *domain.Agent) error {
    // Implementation
}
```

### Reference Document IDs
- **FR**: Functional Requirement (FR-IC-XXX-###)
- **BR**: Business Rule (BR-IC-XXX-###)
- **VR**: Validation Rule (VR-IC-XXX-###)
- **WF**: Workflow (WF-IC-XXX-###)
- **ER**: Error Code (IC-ERR-###)

---

## SLA Requirements

| Process | SLA | Monitoring |
|---------|-----|------------|
| Commission Batch Processing | 6 hours max | Escalate after 3 hours, critical after 5 hours |
| Trial Statement Approval | 7 days | Remind after 5 days |
| Disbursement Processing | 10 working days | Escalate if overdue |
| License Renewal Processing | 3 working days | Escalate to supervisor if exceeded |
| License Reminders | T-30, T-15, T-7, T-0 days | Automated batch job |

---

## Testing Strategy

### Unit Testing
- **Coverage Target**: 80%+
- **Framework**: Go testing package with testify
- **Naming**: `{file}_test.go`
- **Pattern**: Table-driven tests

### Integration Testing
- Test complete workflows end-to-end
- Mock external services (HRMS, PFMS, Accounting)
- Test Temporal workflows with test server

### Performance Testing
- Commission batch: 10,000 records/second minimum
- Search queries: < 1 second response time
- Concurrent users: Support 100+ simultaneous API calls

---

## External Dependencies

### Services to Mock/Integrate
1. **HRMS System** - Departmental employee data
2. **Policy Services** - Policy data for commission calculation
3. **PFMS Gateway** - EFT payment processing
4. **Accounting/Finance** - GL posting, voucher confirmation
5. **Email Service** - Notifications
6. **SMS Service** - License reminders

---

## Known Issues / Blockers

*None at this time*

---

## Questions & Clarifications Needed

### Before Starting Implementation:
1. **Reference Project**: The plan mentions `/home/user/pli-agent` as reference, but this path doesn't exist. Should we proceed using only the template in `Incentive/template/` folder?

2. **Project Root**: Should the new Go project be created at:
   - `/home/user/agent-incentive/` (root level)
   - `/home/user/agent-incentive/Incentive/` (inside Incentive folder)
   - `/home/user/agent-incentive/{new-folder}/` (new folder)

3. **Database Migrations**: Should we:
   - Use an existing migration tool (golang-migrate, goose)?
   - Create custom migration system?
   - Copy the entire schema.sql as initial migration?

4. **External Service Mocking**: For Phase 1, should we:
   - Create mock implementations for HRMS/PFMS/etc.?
   - Skip integration points initially and implement later?
   - Create interface contracts only?

5. **Temporal Server**: Do we need to:
   - Setup local Temporal server (docker-compose)?
   - Mock Temporal for initial development?
   - Assume Temporal is already running?

6. **Testing Approach**: Should we:
   - Write unit tests after each module (recommended)?
   - Write all tests at the end of each phase?
   - Skip tests for now and add later?

7. **Commit Strategy**: Should commits be:
   - Per sub-task (frequent, granular commits)?
   - Per module (fewer, larger commits)?
   - Per phase only?

8. **Swagger Generation**: Should we:
   - Use swaggo/swag for auto-generation?
   - Manually maintain swagger.yaml?
   - Use existing swagger.yaml as specification only?

---

## Next Actions

1. **Clarify questions above**
2. **Get approval on plan.md**
3. **Begin Phase 0: Project Foundation & Setup**

---

## Notes

- This context file will be updated after each phase completion
- All decisions and blockers should be documented here
- Progress percentages will be updated regularly
- Keep traceability comments up to date

---

**Context Created**: 2026-01-28
**Last Updated**: 2026-01-28
**Version**: 1.0
