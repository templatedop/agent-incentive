# Agent Commission Management System - Project Context

## Current Status

**Phase**: Phase 0 - Planning Complete
**Current Task**: Awaiting approval to start implementation
**Last Updated**: 2026-01-28
**Progress**: 0% (0/105 APIs completed)

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
- **Plan**: `/home/user/agent-incentive/plan.md`
- **Requirements**: `/home/user/agent-incentive/Incentive/analysis/IC_Incentive_Commission_Producer_Management_Analysis.md`
- **SRS**: `/home/user/agent-incentive/Incentive/srs/Agent_SRS_Incentive-Commission-and-Producer-Management.md`
- **Swagger**: `/home/user/agent-incentive/Incentive/swagger/swagger.yaml`
- **Database Schema**: `/home/user/agent-incentive/Incentive/db/schema.sql`
- **Template**: `/home/user/agent-incentive/Incentive/template.md`
- **DB Library README**: `/home/user/agent-incentive/Incentive/db-README.md`

### Skill
- **Name**: insurance-code-generator
- **Location**: `/home/user/agent-incentive/Incentive/insurance-code-generator/`
- **Status**: ✅ Loaded and ready

---

## Implementation Progress

### Phase 0: Project Foundation & Setup [NOT STARTED]
- [ ] 0.1 Create project directory structure
- [ ] 0.2 Initialize Go modules
- [ ] 0.3 Setup configuration files
- [ ] 0.4 Create bootstrap/bootstrapper.go
- [ ] 0.5 Setup database migrations
- [ ] 0.6 Create main.go
- [ ] 0.7 Setup Swagger
- [ ] 0.8 Create project files
- [ ] 0.9 Verify build

### Phase 1: Agent Onboarding & Profile Management [NOT STARTED]
**Module 1.1: Agent Onboarding** (6 APIs)
- [ ] 1.1.1 Read requirements
- [ ] 1.1.2 Pattern decisions
- [ ] 1.1.3 Domain models
- [ ] 1.1.4 DTOs
- [ ] 1.1.5 Repository
- [ ] 1.1.6 Temporal workflow
- [ ] 1.1.7 Handler
- [ ] 1.1.8 Unit tests

**Module 1.2: Agent Profile Management** (6 APIs)
- [ ] Not started

**Module 1.3: Agent Termination** (3 APIs)
- [ ] Not started

### Phase 2: License Management [NOT STARTED]
- APIs: 10 endpoints
- Status: Not started

### Phase 3: Commission Rate Configuration [NOT STARTED]
- APIs: 6 endpoints
- Status: Not started

### Phase 4: Commission Calculation & Processing [NOT STARTED]
- APIs: 18 endpoints
- Status: Not started

### Phase 5: Commission Disbursement [NOT STARTED]
- APIs: 12 endpoints
- Status: Not started

### Phase 6: Commission History & Search [NOT STARTED]
- APIs: 8 endpoints
- Status: Not started

### Phase 7: Commission Clawback [NOT STARTED]
- APIs: 7 endpoints
- Status: Not started

### Phase 8: Suspense Account Management [NOT STARTED]
- APIs: 11 endpoints
- Status: Not started

### Phase 9: Lookup & Reference Data [NOT STARTED]
- APIs: 10 endpoints
- Status: Not started

### Phase 10: Export, Reports & Webhooks [NOT STARTED]
- APIs: 10 endpoints
- Status: Not started

### Phase 11: Integration Testing [NOT STARTED]
- Status: Not started

### Phase 12: Documentation & Deployment [NOT STARTED]
- Status: Not started

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
