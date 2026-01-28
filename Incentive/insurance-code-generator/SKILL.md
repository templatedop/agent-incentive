---
name: insurance-code-generator
description: "Generate production-ready Golang code for insurance microservices across ALL domains (Claims, Policy, Agent, Accounting, Portals, KYC, etc.) using Temporal workflows, PostgreSQL, and pgx optimizations. Use when: (1) generating code from Swagger/OpenAPI specs, (2) implementing Temporal workflows for insurance processes, (3) applying database optimizations with pgx.Batch and workflow state management, (4) ensuring consistency across insurance microservices. Handles requirements analysis, API design, database schema integration, and applies patterns from template.md for code generation."
compatibility:
  products: ["claude-code", "claude-desktop", "claude-web"]
---

# Insurance Code Generator (Universal - All Modules)

Generates production-ready Golang code for **ANY insurance API endpoint** across **ALL business domains** by analyzing multiple sources and intelligently applying patterns from template.md.

## Supported Modules (ALL)

✅ **Claims** - Death, Maturity, Survival Benefit, Freelook
✅ **Policy Issuance** - New business, Endorsements, Renewals
✅ **Agent Management** - Commission, Hierarchy, Portal, Producer
✅ **Non-Financial Services** - Address change, Nomination, Name change
✅ **Accounting** - PLI/RPLI calculations, GL integration
✅ **Portals/Apps/Billing** - Customer portal, Mobile apps, Collections
✅ **BCP/KYC/IVRS** - KYC, Aadhaar, PAN, e-KYC
✅ **Any Other Module** - The patterns are universal

## Required Files

### 1. template.md (in project root)
Location: `/mnt/project/template.md`
Contains: All code patterns, database optimization, Temporal workflows, pgx.Batch patterns

### 2. Requirements File (user provides)
Example: `Phase3_Claims_Analysis.md`
Contains: Business rules (BR-*), Functional requirements (FR-*), Validation rules (VR-*), Workflows (WF-*)

### 3. Swagger/OpenAPI (user provides)  
Example: `claims_api_swagger_complete.yaml`
Contains: API endpoints, request/response schemas

### 4. Database DDL (user provides)
Example: `claims_database_schema.sql`
Contains: Table definitions, constraints, indexes

## Three-Phase Methodology

### Phase 1: Multi-Source Analysis

**CRITICAL**: Always analyze ALL sources before making decisions.

#### Step 1.1: Read Requirements
Extract from requirements file:
- Business Rules (BR-*) referenced in Swagger
- Functional Requirements (FR-*)
- Validation Rules (VR-*)
- Workflow Definitions (WF-*)
- SLA requirements

#### Step 1.2: Parse Swagger
Extract:
- HTTP method and path
- Request schema (fields, types, required)
- Response schema
- Related requirements (from **Related**: in description)

#### Step 1.3: Analyze DDL
Extract:
- Primary table and related tables
- Column names, types, constraints
- Indexes (especially partial indexes)
- Partitioning strategy
- Audit column patterns

### Phase 2: Pattern Decisions

**Reference template.md before making ANY decisions.**

#### Decision 1: Use Temporal Workflow?

Read: template.md Section 17

Evaluate:
├─ Process duration > 1 minute? → Use Temporal
├─ Multi-step with waits? → Use Temporal
├─ Human approvals needed? → Use Temporal
├─ Retry/compensation logic? → Use Temporal
└─ Simple CRUD (< 1 minute)? → Skip Temporal

#### Decision 2: Use Workflow State Optimization?

Read: template.md Section 18

IF using Temporal:
  IF multiple activities need same data:
    ✅ Use Workflow State
    - FetchInitialData (pgx.Batch)
    - Store in state
    - Activities use state (NO extra DB calls)

#### Decision 3: Use pgx.Batch?

Read: template.md Section 16

Count queries:
├─ 1 query? → No batch needed
└─ 2+ queries? → Use pgx.Batch
   - QueueReturnRow (single row)
   - QueueReturn (multiple rows)
   - QueueExecRow (no return)

#### Decision 4: Squirrel vs Raw SQL?

Read: template.md Section 15

Per query:
├─ Simple SELECT/INSERT/UPDATE? → Squirrel
├─ INSERT...SELECT? → Raw SQL
├─ UPDATE...FROM? → Raw SQL
├─ WITH (CTE)? → Raw SQL
└─ Dynamic filters? → Squirrel

**For detailed decision criteria and thresholds**, see [references/decisions.md](references/decisions.md)

### Phase 3: Code Generation

**Follow template.md patterns exactly. Never invent patterns not in template.md.**

#### Files Generated

**Always**:
1. `core/domain/{resource}.go` - Domain models
2. `handler/request.go` - Request DTOs
3. `handler/response/{resource}.go` - Response DTOs
4. `repo/postgres/{resource}_repository.go` - Repository

**If Temporal**:
5. `handler/{resource}_handler.go` - Starts workflow
6. `workflows/{resource}_workflow.go` - Workflow with state
7. `workflows/activities/{resource}_activities.go` - Activities
8. `bootstrap/bootstrapper.go` - Registration

**If No Temporal**:
5. `handler/{resource}_handler.go` - Direct repository call
6. `bootstrap/bootstrapper.go` - Registration

**For complete code generation patterns and examples**, see [references/patterns.md](references/patterns.md)

## Code Generation Patterns

### Pattern Overview

The skill uses four main patterns:

1. **Workflow State Struct** - Cache data fetched once and reused across activities
2. **Initial Fetch with pgx.Batch** - Batch multiple queries into one DB round trip
3. **Activities Use State** - Pure calculation/logic using cached workflow state
4. **Handler Starting Workflow** - Initiate Temporal workflow from HTTP handler

**For detailed pattern implementations with complete code examples**, see [references/patterns.md](references/patterns.md)

## Examples by Module

The skill handles endpoints across all insurance modules:

- **Claims Module**: Death claim registration with 21-day workflow
- **Policy Module**: Policy endorsement with batch updates
- **Agent Module**: Commission calculation with hierarchy queries
- **Non-Financial Services**: Address change with conditional approval
- **Accounting Module**: GL integration with batch processing
- **Portal/Billing Module**: Customer profile updates
- **KYC Module**: Aadhaar verification with external API

**For complete examples with analysis, decisions, and generated files**, see [references/examples.md](references/examples.md)

## Quality Checklist

Before delivering code:

**Pattern Compliance**:
- [ ] Follows template.md patterns exactly
- [ ] Uses dblib library (not raw pgx)
- [ ] Squirrel for simple, Raw SQL for complex
- [ ] pgx.Batch used correctly
- [ ] Workflow state when appropriate

**Completeness**:
- [ ] All Swagger fields mapped
- [ ] All business rules implemented
- [ ] Validation rules included
- [ ] Audit columns populated

**Database**:
- [ ] Column names match DDL
- [ ] Data types correct
- [ ] Constraints respected
- [ ] Indexes utilized

**Temporal** (if applicable):
- [ ] Workflow state struct defined
- [ ] Initial fetch uses pgx.Batch
- [ ] Subsequent activities use state
- [ ] Signal/Query handlers included
- [ ] Retry policies configured

## Quick Usage Guide

### Simple CRUD (No Temporal)
```
User: "Generate code for GET /policies/{id}"
→ Handler → Repository (direct call)
```

### Complex Workflow (With Temporal)
```
User: "Generate code for POST /claims/death/register"
→ Analyze requirements (21-day SLA, approvals)
→ Decide: Temporal + Workflow State + pgx.Batch
→ Generate 8 files with optimizations
```

### Batch Processing
```
User: "Generate code for POST /accounting/gl-integration/batch"
→ Decide: Temporal + Raw SQL + pgx.Batch
→ Generate workflow with bulk operations
```

## Error Handling

**template.md not found**:
```
Solution: Verify template.md is in /mnt/project/ directory
```

**Missing business rules**:
```
Solution: Generate with available info, add TODO comments
```

**Swagger schema incomplete**:
```
Solution: Generate with available fields, warn user
```

## Reference Files

- **[references/decisions.md](references/decisions.md)** - Complete decision criteria tables and thresholds
- **[references/patterns.md](references/patterns.md)** - Full code generation patterns with implementations
- **[references/examples.md](references/examples.md)** - Detailed examples for each insurance module

## Version

**Skill**: 1.0.0  
**Template.md**: Requires v1.0+ (Sections 15-18)  
**Tech**: Go 1.25, Temporal 1.25, PostgreSQL 16  
**Framework**: n-api-template

---

**Note**: Always read template.md sections before generating code. Never invent patterns not in template.md. Reference the files in the `references/` directory for detailed examples, patterns, and decision criteria.
