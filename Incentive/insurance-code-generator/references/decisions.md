# Decision Reference Guide

This file provides detailed decision criteria and thresholds for choosing the right patterns when generating insurance microservice code.

## Decision 1: Temporal Workflow

### When to Use Temporal Workflow

| Indicator | Threshold | Use Temporal | Reason |
|-----------|-----------|--------------|---------|
| Process duration | > 1 minute | ✅ | Temporal handles long-running processes |
| Approvals required | Any | ✅ | Human-in-the-loop workflows need Temporal |
| Wait conditions | Any | ✅ | Timer/signal patterns for waits |
| SLA tracking | > 1 day | ✅ | Built-in visibility and monitoring |
| Compensation logic | Complex | ✅ | Saga pattern support |
| Retries needed | Custom policies | ✅ | Configurable retry strategies |
| External API calls | With retries | ✅ | Activity retry patterns |
| Multi-step coordination | 3+ steps | ✅ | Workflow orchestration |

### When NOT to Use Temporal

| Indicator | Threshold | Skip Temporal | Reason |
|-----------|-----------|---------------|---------|
| Process duration | < 1 minute | ✅ | Overhead not justified |
| Simple CRUD | Single operation | ✅ | Direct repository call sufficient |
| Synchronous response | Immediate | ✅ | No async benefit |
| No retry logic | Atomic operation | ✅ | Database transaction enough |

### Examples by Module

**Claims Module**:
- Death claim registration: ✅ Temporal (21-day investigation, approvals)
- Claim status check: ❌ No Temporal (simple query)

**Policy Module**:
- Policy issuance: ✅ Temporal (underwriting, approval workflow)
- Policy endorsement: ❌ No Temporal (quick update)

**Agent Module**:
- Agent onboarding: ✅ Temporal (multi-step verification, approvals)
- Commission calculation: ❌ No Temporal (fast calculation)

**Accounting Module**:
- GL batch integration: ✅ Temporal (batch processing with retries)
- Balance inquiry: ❌ No Temporal (simple query)

---

## Decision 2: Workflow State Optimization

### When to Use Workflow State

| Indicator | Threshold | Use State | Benefit |
|-----------|-----------|-----------|---------|
| Data reuse | 2+ activities | ✅ | Reduce DB calls |
| DB call reduction | 30%+ | ✅ | Performance gain |
| Complex data structure | Shared across activities | ✅ | Consistency |
| Immutable reference data | Used multiple times | ✅ | Cache efficiency |

### Pattern: Workflow State Struct

**Structure**:
```
WorkflowState:
├─ Metadata (IDs, stage)
├─ Cached Data (fetched once)
└─ Calculated Data (computed once)
```

**Optimization Impact**:
```
Without State: N activities × M queries = N×M DB calls
With State: 1 batch query + N activities = 1 DB call
Reduction: ((N×M) - 1) / (N×M) × 100%
```

### Examples

**Example 1: Death Claim (8 → 2 DB calls = 75% reduction)**
```
Activities needing policy data:
1. Validate eligibility (needs policy)
2. Calculate payout (needs policy + riders)
3. Check exclusions (needs policy + coverage)
4. Verify beneficiary (needs policy + nominations)

Without State: 4 queries × 2 tables = 8 DB calls
With State: 1 batch query (fetch all) = 1 DB call
Savings: 87.5%
```

**Example 2: Address Change (6 → 2 DB calls = 67% reduction)**
```
Activities:
1. Validate address (needs policy + customer)
2. Check communication preferences (needs customer)
3. Update address (writes to DB)

Without State: 3 queries
With State: 1 batch query + 1 update = 2 DB calls
Savings: 33%
```

---

## Decision 3: pgx.Batch

### When to Use pgx.Batch

| Scenario | Query Count | Use Batch | Method |
|----------|-------------|-----------|---------|
| Single query | 1 | ❌ | Direct query |
| Multiple reads | 2+ | ✅ | QueueReturn/QueueReturnRow |
| Multiple writes | 2+ | ✅ | QueueExecRow |
| Mixed read/write | 2+ | ✅ | Combined methods |
| Transaction required | Multiple ops | ✅ | Atomic batch |

### Batch Methods

**QueueReturnRow** - Single row result:
```go
var entity domain.Entity
dblib.QueueReturnRow(batch, query, pgx.RowToStructByNameLax[domain.Entity], &entity)
```

**QueueReturn** - Multiple rows result:
```go
var entities []domain.Entity
dblib.QueueReturn(batch, query, pgx.RowToStructByNameLax[domain.Entity])
```

**QueueExecRow** - No result (INSERT/UPDATE/DELETE):
```go
dblib.QueueExecRow(batch, query)
```

### Performance Impact

| Queries | Without Batch | With Batch | Reduction |
|---------|--------------|------------|-----------|
| 2 | 2 round trips | 1 round trip | 50% |
| 3 | 3 round trips | 1 round trip | 67% |
| 4 | 4 round trips | 1 round trip | 75% |
| 5 | 5 round trips | 1 round trip | 80% |

---

## Decision 4: Squirrel vs Raw SQL

### Query Type Matrix

| Query Pattern | Tool | Reason |
|--------------|------|---------|
| Simple SELECT | Squirrel | Type-safe, readable |
| Simple INSERT | Squirrel | Auto-generated values |
| Simple UPDATE | Squirrel | Safe parameter binding |
| Simple DELETE | Squirrel | WHERE clause safety |
| Dynamic filters | Squirrel | Conditional WHERE building |
| INSERT...SELECT | Raw SQL | Squirrel doesn't support |
| UPDATE...FROM | Raw SQL | Squirrel doesn't support |
| WITH (CTE) | Raw SQL | Complex query structure |
| Recursive CTE | Raw SQL | Advanced SQL feature |
| WINDOW functions | Raw SQL | Complex aggregations |
| Complex JOINs (4+) | Raw SQL | Query readability |
| Database-specific features | Raw SQL | PostgreSQL extensions |

### Examples

**Squirrel - Simple Queries**:
```go
// SELECT with WHERE
query := dblib.Psql.Select("*").
    From("policies").
    Where(sq.Eq{"policy_id": id})

// INSERT with returning
query := dblib.Psql.Insert("claims").
    Columns("policy_id", "claim_type", "amount").
    Values(policyID, claimType, amount).
    Suffix("RETURNING claim_id")

// UPDATE with WHERE
query := dblib.Psql.Update("policies").
    Set("status", "ACTIVE").
    Where(sq.Eq{"policy_id": id})

// Dynamic filters
query := dblib.Psql.Select("*").From("claims")
if status != "" {
    query = query.Where(sq.Eq{"status": status})
}
if fromDate != nil {
    query = query.Where(sq.GtOrEq{"created_at": fromDate})
}
```

**Raw SQL - Complex Queries**:
```go
// INSERT...SELECT (Squirrel doesn't support)
query := `
    INSERT INTO claim_history (claim_id, action, user_id)
    SELECT claim_id, 'APPROVED', $1
    FROM claims
    WHERE status = 'PENDING'
`

// UPDATE...FROM (Squirrel doesn't support)
query := `
    UPDATE policies p
    SET premium_status = 'OVERDUE'
    FROM premium_schedules ps
    WHERE p.policy_id = ps.policy_id
      AND ps.due_date < NOW()
      AND ps.payment_status = 'UNPAID'
`

// WITH (CTE)
query := `
    WITH agent_hierarchy AS (
        SELECT agent_id, parent_agent_id, level
        FROM agents
        WHERE agent_id = $1
        UNION ALL
        SELECT a.agent_id, a.parent_agent_id, ah.level + 1
        FROM agents a
        JOIN agent_hierarchy ah ON a.parent_agent_id = ah.agent_id
    )
    SELECT * FROM agent_hierarchy
`

// Complex aggregation with WINDOW functions
query := `
    SELECT 
        policy_id,
        premium_amount,
        SUM(premium_amount) OVER (PARTITION BY agent_id ORDER BY policy_date) as running_total,
        ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY policy_date DESC) as policy_rank
    FROM policies
    WHERE policy_date >= $1
`
```

---

## Decision Workflow

### Step-by-Step Decision Process

```
1. Analyze Requirements
   ├─ Check duration → > 1 min? → Consider Temporal
   ├─ Check approvals → Any? → Use Temporal
   └─ Check complexity → Multi-step? → Use Temporal

2. If Using Temporal
   ├─ Count data dependencies → 2+ activities? → Use Workflow State
   └─ Design state struct

3. Analyze Database Operations
   ├─ Count queries → 2+? → Use pgx.Batch
   └─ Choose batch methods (QueueReturn*, QueueExecRow)

4. For Each Query
   ├─ Simple CRUD? → Use Squirrel
   ├─ INSERT...SELECT? → Raw SQL
   ├─ UPDATE...FROM? → Raw SQL
   ├─ CTE needed? → Raw SQL
   └─ Complex join? → Raw SQL

5. Generate Code
   └─ Apply patterns from template.md
```

### Decision Matrix Cheat Sheet

| Scenario | Temporal | State | Batch | SQL Type |
|----------|----------|-------|-------|----------|
| Simple GET /resource/{id} | ❌ | ❌ | ❌ | Squirrel |
| POST with 3 related tables | ❌ | ❌ | ✅ | Squirrel |
| Long workflow (approval) | ✅ | ✅ | ✅ | Mixed |
| Batch processing | ✅ | ❌ | ✅ | Raw SQL |
| Complex hierarchy query | ❌ | ❌ | ❌ | Raw SQL |
| Multi-step calculation | ✅ | ✅ | ✅ | Squirrel |

---

## Module-Specific Decision Patterns

### Claims Module
- **Death Claims**: Temporal + State + Batch (long investigation)
- **Maturity Claims**: Temporal + State + Batch (calculation + approval)
- **Freelook Cancellation**: No Temporal (quick refund)

### Policy Module
- **New Policy Issuance**: Temporal + State + Batch (underwriting)
- **Endorsement**: No Temporal, Batch only (quick update)
- **Renewal**: Temporal (periodic process)

### Agent Module
- **Commission Calculation**: No Temporal, Raw SQL (complex CTE)
- **Agent Onboarding**: Temporal + State + Batch (verification workflow)
- **Hierarchy Management**: No Temporal, Raw SQL (recursive CTE)

### Accounting Module
- **GL Integration Batch**: Temporal + Batch + Raw SQL (bulk processing)
- **PLI/RPLI Calculation**: No Temporal, Raw SQL (complex formulas)
- **Balance Inquiry**: No Temporal, Squirrel (simple query)

### Non-Financial Services
- **Address Change**: Temporal (if permanent, needs approval)
- **Nomination Update**: No Temporal, Batch (multiple tables)
- **Communication Preference**: No Temporal, Squirrel (simple update)

### Portal/Billing
- **Profile Update**: No Temporal, Batch (customer + preferences)
- **Payment Processing**: Temporal + Batch (payment gateway workflow)
- **Statement Generation**: Temporal (batch report generation)

### KYC/BCP
- **Aadhaar Verification**: Temporal (external API with retries)
- **PAN Validation**: Temporal (external API with retries)
- **KYC Document Upload**: No Temporal, Batch (storage + metadata)

---

## Common Optimization Patterns

### Pattern 1: Fetch-Once-Use-Many
```
Scenario: Multiple activities need same policy data
Without: 4 activities × 1 query each = 4 DB calls
With: 1 batch fetch + state = 1 DB call
Savings: 75%
```

### Pattern 2: Batch Related Entities
```
Scenario: Fetch policy + riders + loans + nominations
Without: 4 separate queries = 4 round trips
With: 1 pgx.Batch = 1 round trip
Savings: 75%
```

### Pattern 3: Combine Temporal + State + Batch
```
Scenario: Complex workflow with multiple steps
FetchInitialData (pgx.Batch): policy + customer + agents
Store in workflow state
Activities: Use state for pure logic
Final step: Update with pgx.Batch

Total DB calls: 2 (fetch + update)
Without optimization: 10+ DB calls
Savings: 80%+
```

---

## Version

**Version**: 1.0.0  
**Last Updated**: 2026-01-20  
**Compatible with**: template.md v1.0+
