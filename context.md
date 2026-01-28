# Agent Commission Management System - Project Context

## Current Status

**Phase**: Phase 1 - Foundation & Reference Data - COMPLETE ✅
**Current Module**: Ready for Phase 2 - Commission Processing Core
**Last Updated**: 2026-01-28
**Progress**: 10% (3/30 APIs completed) | Phase 1: 100% (3/3 APIs)

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

### Phase 2: Commission Processing Core [NOT STARTED]
- Module 2.1: Commission Batch Processing (5 APIs)
- Module 2.2: Trial Statement Management (2 APIs)
- Module 2.3: Final Statement Generation (1 API)

### Phase 3: Commission Disbursement [NOT STARTED]
- Module 3.1: Disbursement Processing (2 APIs)
- Module 3.2: Integration Webhooks (2 APIs)

### Phase 4: Commission History & Inquiry [NOT STARTED]
- Module 4.1: Commission History (1 API)

### Phase 5: Clawback Management [NOT STARTED]
- Module 5.1: Commission Clawback (1 API)

### Phase 6: Suspense Account Management [NOT STARTED]
- Module 6.1: Suspense Accounts (2 APIs)

### Phase 7: Workflow Management [NOT STARTED]
- Module 7.1: Workflow Status & Control (8 APIs)

### Phase 8: Export & Reporting [NOT STARTED]
- Module 8.1: Export APIs (6 APIs)

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
