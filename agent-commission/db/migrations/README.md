# Database Migrations

This directory contains database migration scripts for the Agent Commission Management System.

## Migration Tool

We use [golang-migrate/migrate](https://github.com/golang-migrate/migrate) for managing database migrations.

### Installation

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Migration Files

Migrations follow the naming convention:
```
{version}_{description}.up.sql   # Forward migration
{version}_{description}.down.sql # Rollback migration
```

### Current Migrations

| Version | Description | Status |
|---------|-------------|--------|
| 001 | Initial schema (enums, tables, indexes) | Ready |

## Running Migrations

### Apply All Migrations (Up)

```bash
# Using migrate CLI
migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/agent_commission?sslmode=disable" up

# Or using docker-compose
docker-compose exec postgres psql -U postgres -d agent_commission -f /migrations/001_initial_schema.up.sql
```

### Rollback Last Migration (Down)

```bash
migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/agent_commission?sslmode=disable" down 1
```

### Rollback All Migrations

```bash
migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/agent_commission?sslmode=disable" down
```

### Check Migration Status

```bash
migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/agent_commission?sslmode=disable" version
```

## Database Connection String Format

```
postgresql://[user]:[password]@[host]:[port]/[database]?sslmode=[mode]
```

### Environment-specific

**Development:**
```
postgresql://postgres:postgres@localhost:5432/agent_commission_dev?sslmode=disable
```

**Test:**
```
postgresql://postgres:postgres@localhost:5432/agent_commission_test?sslmode=disable
```

**Production:**
```
postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:5432/agent_commission_prod?sslmode=require
```

## Schema Overview

### Tables Created (14 total)

1. **ref_circles** - Circle master data
2. **ref_divisions** - Division master data  
3. **ref_product_plans** - Product plan configurations
4. **agent_profiles** - Agent master records (FR-IC-PROF-001)
5. **agent_addresses** - Agent address details
6. **agent_contacts** - Agent contact numbers
7. **agent_emails** - Agent email addresses
8. **agent_hierarchy** - Agent-Coordinator relationships (BR-IC-AH-001)
9. **agent_licenses** - License management (FR-IC-LIC-001)
10. **license_renewal_reminders** - Renewal reminder tracking
11. **commission_rate_config** - Commission rate configuration (FR-IC-RATE-001)
12. **commission_records** - Individual commission records (FR-IC-COM-002)
13. **commission_batch_log** - Batch processing logs
14. **commission_trial_statements** - Trial statements (FR-IC-COM-003)
15. **commission_final_statements** - Final statements (FR-IC-COM-007)
16. **commission_disbursements** - Disbursement records (FR-IC-DIS-001)
17. **disbursement_cheque_details** - Cheque payment details
18. **disbursement_eft_details** - EFT payment details
19. **commission_clawback** - Clawback management (FR-IC-CLAW-001)
20. **commission_suspense_accounts** - Suspense account management (FR-IC-SUSP-001)

### Enum Types Created

- agent_type_enum
- person_type_enum
- gender_enum
- marital_status_enum
- agent_status_enum
- commission_status_enum
- commission_type_enum
- product_type_enum
- address_type_enum
- contact_type_enum
- email_type_enum
- payment_mode_enum
- disbursement_status_enum
- license_status_enum
- reminder_status_enum
- clawback_status_enum
- suspense_status_enum

## Development Workflow

1. **Create New Migration:**
   ```bash
   migrate create -ext sql -dir db/migrations -seq migration_name
   ```

2. **Test Migration:**
   ```bash
   # Apply
   migrate -path db/migrations -database "postgresql://..." up
   
   # Verify
   psql -U postgres -d agent_commission -c "\dt"
   
   # Rollback
   migrate -path db/migrations -database "postgresql://..." down 1
   ```

3. **Commit to Git:**
   ```bash
   git add db/migrations/
   git commit -m "Add migration: description"
   ```

## Troubleshooting

### Migration Version Conflict
```bash
# Force version
migrate -path db/migrations -database "postgresql://..." force VERSION
```

### Reset Database (DANGER: Deletes all data)
```bash
# Drop all migrations
migrate -path db/migrations -database "postgresql://..." down

# Reapply all migrations
migrate -path db/migrations -database "postgresql://..." up
```

### View Migration Status in Database
```sql
SELECT * FROM schema_migrations;
```

## Best Practices

1. **Never modify existing migrations** - Create a new migration instead
2. **Always test both up and down migrations**
3. **Keep migrations small and focused**
4. **Add comments with FR/BR/VR traceability**
5. **Test migrations on a copy of production data before deploying**

---

**Last Updated:** 2026-01-28
**Schema Version:** 1.0.0
