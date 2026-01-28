# Agent Incentive & Commission Management - REVISED Implementation Plan

**Project**: Commission & Incentive Module (ONLY)
**Scope**: Based on `swagger_incentive.yaml`
**Total APIs**: 30 endpoints
**Duration**: 30 working days
**Agent Profile Management**: OUT OF SCOPE (handled by complete-pli-agent)

---

## 🎯 Scope Clarification

### ✅ IN SCOPE (This Project):
- Commission calculation & processing
- Trial statement approval workflow
- Commission disbursement (Cheque/EFT)
- Commission history & inquiry
- Clawback management
- Suspense account management
- Commission rate configuration
- Lookup APIs (read-only reference data)
- Integration webhooks (PFMS, Policy Services)
- Workflow management APIs
- Export & reporting

### ❌ OUT OF SCOPE (Handled by complete-pli-agent):
- Agent onboarding
- Agent profile CRUD
- Agent search & filtering
- Agent status management
- Agent termination
- License management (full module)
- Agent hierarchy management

---

## 📋 Revised Module Breakdown

### **Phase 1: Foundation & Reference Data** (Days 1-3)

#### Module 1.1: Commission Rate Configuration (3 APIs)
- [ ] 1.1.1 Domain models: CommissionRate, CommissionRateHistory
- [ ] 1.1.2 Repository: commission_rate_repository.go
- [ ] 1.1.3 Handler: commission_rate_handler.go
  - GET /lookup/commission-rates - Get commission rates
  - POST /commissions/rates - Create/Update rate (admin)
  - GET /commissions/rates/history - Rate change history

**Key Requirements**:
- FR-IC-COM-001: Commission rate table management
- BR-IC-COM-006: Rate structure (product type, agent type, plan code)
- VR-IC-COM-001: Rate validation rules

#### Module 1.2: Lookup & Reference APIs (2 APIs)
- [ ] 1.2.1 Repository: reference_data_repository.go (extend existing)
- [ ] 1.2.2 Handler: lookup_handler.go
  - GET /lookup/advisor-coordinators - Active coordinators (from agent profile DB)
  - GET /lookup/circles - All circles

**Integration**: Read-only access to agent profile database (complete-pli-agent)

---

### **Phase 2: Commission Processing Core** (Days 4-10)

#### Module 2.1: Commission Batch Processing (5 APIs)
- [ ] 2.1.1 Domain models: CommissionBatch, CommissionCalculation, PolicyCommission
- [ ] 2.1.2 Repository: commission_batch_repository.go, commission_calculation_repository.go
- [ ] 2.1.3 Temporal Workflow: commission_batch_workflow.go
  - WF-IC-COM-001: Monthly Commission Processing
  - WF-TEMPORAL-IC-001: Commission Batch Workflow
  - Activities: FetchPolicies, CalculateCommission, ApplyTDS, GenerateTrialStatements
- [ ] 2.1.4 Handler: commission_batch_handler.go
  - POST /commissions/batches/start - Start commission batch
  - GET /commissions/batches/{batchId}/status - Get batch status

**Key Requirements**:
- FR-IC-COM-002: Monthly commission calculation
- BR-IC-COM-001: Commission calculation rules
- BR-IC-COM-003: TDS calculation based on PAN
- BR-IC-COM-012: 6-hour SLA for batch completion
- INT-IC-002: Policy Services integration

**Database Optimization**:
- Use pgx.Batch for bulk commission inserts (1000+ records per batch)
- Batch policy data fetching from Policy Services

#### Module 2.2: Trial Statement Management (2 APIs)
- [ ] 2.2.1 Domain models: TrialStatement, TrialStatementDetail
- [ ] 2.2.2 Repository: trial_statement_repository.go
- [ ] 2.2.3 Temporal Workflow: trial_approval_workflow.go
  - WF-IC-TRIAL-001: Trial Statement Approval (7-day SLA)
- [ ] 2.2.4 Handler: trial_statement_handler.go
  - GET /commissions/trial-statements - List trial statements
  - POST /commissions/trial-statements/{statementId}/approve - Approve trial

**Key Requirements**:
- FR-IC-COM-004: Trial statement generation
- FR-IC-COM-006: Trial statement approval
- BR-IC-COM-002: Approval workflow (Finance Head for > ₹10L)
- BR-IC-COM-005: Partial disbursement support
- BR-IC-COM-009: 7-day approval SLA

#### Module 2.3: Final Statement Generation (1 API)
- [ ] 2.3.1 Domain models: FinalStatement
- [ ] 2.3.2 Repository: final_statement_repository.go
- [ ] 2.3.3 Handler: final_statement_handler.go
  - GET /commissions/final-statements - List final statements

**Key Requirements**:
- FR-IC-COM-008: Final statement generation after approval

---

### **Phase 3: Commission Disbursement** (Days 11-15)

#### Module 3.1: Disbursement Processing (2 APIs)
- [ ] 3.1.1 Domain models: Disbursement, DisbursementDetail
- [ ] 3.1.2 Repository: disbursement_repository.go
- [ ] 3.1.3 Temporal Workflow: disbursement_workflow.go
  - WF-IC-DISB-001: Commission Disbursement Workflow
  - WF-TEMPORAL-IC-003: Disbursement SLA Monitoring
  - Activities: ValidateBankDetails, ProcessCheque, ProcessEFT, PostAccounting
- [ ] 3.1.4 Handler: disbursement_handler.go
  - POST /commissions/disbursements - Process disbursement
  - GET /commissions/disbursements/{disbursementId}/status - Get status

**Key Requirements**:
- FR-IC-COM-009: Cheque disbursement
- FR-IC-COM-010: EFT disbursement via PFMS
- BR-IC-COM-008: Disbursement mode workflow (Cheque/EFT)
- BR-IC-COM-011: 10-working-day SLA
- VR-IC-COM-004: Disbursement validation
- INT-IC-003: PFMS/Bank integration
- INT-IC-004: Accounting system integration

**Database Optimization**:
- Batch disbursement status updates
- Use Temporal for retry logic on PFMS failures

#### Module 3.2: Integration Webhooks (2 APIs)
- [ ] 3.2.1 Handler: webhook_handler.go
  - POST /webhooks/pfms/disbursement-confirmation - PFMS callback
  - POST /webhooks/policy/status-change - Policy status callback

**Key Requirements**:
- INT-IC-002: Policy Services webhook
- INT-IC-003: PFMS webhook
- Webhook signature verification
- Idempotency handling

---

### **Phase 4: Commission History & Inquiry** (Days 16-18)

#### Module 4.1: Commission History (1 API)
- [ ] 4.1.1 Domain models: CommissionTransaction
- [ ] 4.1.2 Repository: commission_history_repository.go (with batch-optimized search)
- [ ] 4.1.3 Handler: commission_history_handler.go
  - GET /commissions/history - Search commission history

**Key Requirements**:
- FR-IC-COM-011: Commission history inquiry
- BR-IC-COM-009: History filtering and pagination

**Database Optimization**:
- Use batch for count + results (like SearchAgents pattern)
- Index on agent_id, policy_number, date fields

---

### **Phase 5: Clawback Management** (Days 19-21)

#### Module 5.1: Commission Clawback (1 API)
- [ ] 5.1.1 Domain models: Clawback, ClawbackRecovery
- [ ] 5.1.2 Repository: clawback_repository.go
- [ ] 5.1.3 Temporal Workflow: clawback_workflow.go
  - WF-IC-CLAWBACK-001: Commission Clawback Process
  - WF-TEMPORAL-IC-005: Clawback Recovery Workflow
  - Activities: CalculateClawback, ScheduleRecovery, PostRecovery
- [ ] 5.1.4 Handler: clawback_handler.go
  - POST /commissions/clawback/create - Create clawback

**Key Requirements**:
- FR-IC-COM-011: Clawback initiation
- BR-IC-CLAWBACK-001: Graduated recovery (100%, 75%, 50% over months)
- Integration with Policy status webhook

---

### **Phase 6: Suspense Account Management** (Days 22-24)

#### Module 6.1: Suspense Accounts (2 APIs)
- [ ] 6.1.1 Domain models: SuspenseAccount, SuspenseTransaction
- [ ] 6.1.2 Repository: suspense_repository.go
- [ ] 6.1.3 Temporal Workflow: suspense_workflow.go
  - WF-IC-SUSPENSE-001: Suspense Management Workflow
  - WF-TEMPORAL-IC-006: Suspense Retry & Aging
  - Activities: CreateSuspense, RetryPayment, GenerateAgingReport
- [ ] 6.1.4 Handler: suspense_handler.go
  - POST /commissions/suspense/create - Create suspense entry
  - GET /commissions/suspense/aging-report - Get aging report

**Key Requirements**:
- FR-IC-COM-013: Suspense account creation
- FR-IC-COM-014: Suspense resolution
- FR-IC-COM-018: Aging report
- BR-IC-SUSPENSE-001: Suspense reasons (investigation, payment failed, disputed)
- BR-IC-SUSPENSE-002: Retry with exponential backoff
- BR-IC-SUSPENSE-003: Age buckets (0-30, 31-60, 61-90, 90+ days)

---

### **Phase 7: Workflow Management** (Days 25-26)

#### Module 7.1: Workflow Status & Control (8 APIs)
- [ ] 7.1.1 Handler: workflow_management_handler.go
  - GET /workflows/{workflowId}/status - Get workflow status
  - GET /workflows/{workflowId}/history - Get execution history
  - POST /workflows/{workflowId}/cancel - Cancel workflow
  - POST /workflows/{workflowId}/retry - Retry failed workflow
  - POST /workflows/{workflowId}/query - Query workflow state
  - GET /batches/{batchId}/status - Get batch status
  - GET /batches/{batchId}/progress - Get batch progress
  - POST /batches/{batchId}/cancel - Cancel batch

**Key Requirements**:
- Temporal workflow inspection APIs
- SLA tracking and monitoring
- Batch progress tracking

---

### **Phase 8: Export & Reporting** (Days 27-29)

#### Module 8.1: Export APIs (6 APIs)
- [ ] 8.1.1 Domain models: Export, ExportJob
- [ ] 8.1.2 Repository: export_repository.go
- [ ] 8.1.3 Export Service: excel_export_service.go, pdf_export_service.go
- [ ] 8.1.4 Handler: export_handler.go
  - POST /exports/commissions/trial-statements - Export trial statements
  - POST /exports/commissions/final-statements - Export final statements
  - POST /exports/commissions/history - Export commission history
  - POST /exports/commissions/clawback - Export clawback report
  - POST /exports/commissions/suspense - Export suspense report
  - POST /reports/commission-summary - Generate summary report

**Key Requirements**:
- FR-IC-COM-010: Excel/PDF export
- FR-IC-COM-012: Commission summary reports
- BR-IC-COM-010: Report generation rules
- Async export with download URLs
- Scheduled report generation

---

### **Phase 9: Integration & Testing** (Day 30)

#### Module 9.1: Integration Layer
- [ ] 9.1.1 Policy Service client (HTTP/gRPC)
- [ ] 9.1.2 PFMS/Bank client
- [ ] 9.1.3 Accounting system client
- [ ] 9.1.4 Webhook signature verification

#### Module 9.2: End-to-End Testing
- [ ] 9.2.1 Commission batch E2E test
- [ ] 9.2.2 Disbursement workflow E2E test
- [ ] 9.2.3 Clawback workflow E2E test
- [ ] 9.2.4 Load test: 10,000 policies commission calculation

---

## 🗑️ Cleanup Required

### What to Remove from Current Codebase:
1. **agent-commission/handler/agent_onboarding.go** - ❌ Delete (handled by complete-pli-agent)
2. **agent-commission/workflows/agent_onboarding_workflow.go** - ❌ Delete
3. **agent-commission/workflows/activities/agent_activities.go** - ❌ Delete
4. **agent-commission/handler/response/agent_wrappers.go** - ❌ Delete
5. **bootstrap/bootstrapper.go** - ⚠️ Remove AgentOnboardingHandler registration

### What to Keep:
1. **agent-commission/core/domain/** - ✅ Keep (can reuse AgentProfile, Address for read-only lookups)
2. **agent-commission/repo/postgres/agent_profile_repository.go** - ✅ Keep ONLY read methods (GetCoordinatorByID, GetAgentByID for lookups)
3. **agent-commission/repo/postgres/reference_data_repository.go** - ✅ Keep (needed for circles, divisions)
4. **Database schema** - ✅ Keep (commission module needs agent reference data)
5. **Batch optimization patterns** - ✅ Keep and apply to all commission operations

---

## 📊 API Summary

| Module | APIs | Temporal Workflows | Days |
|--------|------|-------------------|------|
| Commission Rate Configuration | 3 | 0 | 2 |
| Lookup & Reference | 2 | 0 | 1 |
| Commission Batch Processing | 5 | 2 | 5 |
| Trial & Final Statements | 3 | 1 | 2 |
| Disbursement | 2 | 2 | 4 |
| Integration Webhooks | 2 | 0 | 1 |
| Commission History | 1 | 0 | 2 |
| Clawback | 1 | 2 | 3 |
| Suspense Accounts | 2 | 2 | 3 |
| Workflow Management | 8 | 0 | 2 |
| Export & Reports | 6 | 0 | 3 |
| Integration & Testing | - | - | 2 |
| **TOTAL** | **30** | **9** | **30** |

---

## 🎯 Key Technical Patterns

### 1. Batch Optimization (Applied Throughout)
```go
// Pattern: Multi-query operations use pgx.Batch
batch := &pgx.Batch{}
dblib.QueueReturnRow(batch, query1, scanFn, &result1)
dblib.QueueReturn(batch, query2, scanFn, &result2)
batchResults := r.db.Pool.SendBatch(ctx, batch)
defer batchResults.Close()
```

### 2. Temporal Workflows
- Commission batch processing (6-hour SLA)
- Trial approval workflow (7-day SLA)
- Disbursement SLA monitoring (10-day SLA)
- Clawback graduated recovery
- Suspense retry with backoff

### 3. Integration Patterns
- Webhook signature verification
- Idempotency handling
- Retry with exponential backoff
- Circuit breaker for external services

### 4. Handler Pattern (from template)
```go
func (h *Handler) HandlerFunc(sctx *serverRoute.Context, req RequestType) (*ResponseType, error) {
    // Plain Go function with sctx.Ctx for context
}
```

---

## 🔗 Dependencies

### External Systems:
1. **complete-pli-agent** (Agent Profile DB) - Read-only access for:
   - Agent lookups (coordinator assignment)
   - Circle/Division reference data
   - Agent status verification

2. **Policy Services** - Integration for:
   - Policy data for commission calculation
   - Policy status change webhooks

3. **PFMS/Bank Gateway** - Integration for:
   - EFT payment processing
   - Payment confirmation webhooks

4. **Accounting System** - Integration for:
   - GL posting
   - Voucher confirmation

### Shared Database Tables:
- Read-only access to: `agent_profiles`, `circles`, `divisions`, `agent_hierarchy`
- Own tables: All commission-related tables (commission_batches, trial_statements, disbursements, clawback, suspense)

---

## ✅ Next Steps

1. **Cleanup Phase** (1-2 hours):
   - Remove agent onboarding handler, workflows, activities
   - Update bootstrap configuration
   - Commit cleanup changes

2. **Start Phase 1** (Commission Rate Configuration):
   - Create commission rate domain models
   - Implement commission rate repository with batch optimization
   - Build commission rate handler with 3 APIs

3. **Parallel Track**:
   - Define database schema for commission tables
   - Setup integration client stubs (Policy Services, PFMS)

---

**Plan Created**: 2026-01-28
**Based On**: `Incentive/swagger/swagger_incentive.yaml`
**Focus**: Commission & Incentive Management ONLY
