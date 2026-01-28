-- ============================================================================
-- Agent Incentive Database Schema Rollback
-- Drops all tables and types in reverse order
-- ============================================================================

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS commission_suspense_accounts CASCADE;
DROP TABLE IF EXISTS commission_clawback CASCADE;
DROP TABLE IF EXISTS disbursement_eft_details CASCADE;
DROP TABLE IF EXISTS disbursement_cheque_details CASCADE;
DROP TABLE IF EXISTS commission_disbursements CASCADE;
DROP TABLE IF EXISTS commission_final_statements CASCADE;
DROP TABLE IF EXISTS commission_trial_statements CASCADE;
DROP TABLE IF EXISTS commission_batch_log CASCADE;
DROP TABLE IF EXISTS commission_records CASCADE;
DROP TABLE IF EXISTS commission_rate_config CASCADE;
DROP TABLE IF EXISTS license_renewal_reminders CASCADE;
DROP TABLE IF EXISTS agent_licenses CASCADE;
DROP TABLE IF EXISTS agent_hierarchy CASCADE;
DROP TABLE IF EXISTS agent_emails CASCADE;
DROP TABLE IF EXISTS agent_contacts CASCADE;
DROP TABLE IF EXISTS agent_addresses CASCADE;
DROP TABLE IF EXISTS agent_profiles CASCADE;
DROP TABLE IF EXISTS ref_product_plans CASCADE;
DROP TABLE IF EXISTS ref_divisions CASCADE;
DROP TABLE IF EXISTS ref_circles CASCADE;

-- Drop enum types
DROP TYPE IF EXISTS email_type_enum CASCADE;
DROP TYPE IF EXISTS contact_type_enum CASCADE;
DROP TYPE IF EXISTS address_type_enum CASCADE;
DROP TYPE IF EXISTS product_type_enum CASCADE;
DROP TYPE IF EXISTS commission_type_enum CASCADE;
DROP TYPE IF EXISTS commission_status_enum CASCADE;
DROP TYPE IF EXISTS agent_status_enum CASCADE;
DROP TYPE IF EXISTS marital_status_enum CASCADE;
DROP TYPE IF EXISTS gender_enum CASCADE;
DROP TYPE IF EXISTS person_type_enum CASCADE;
DROP TYPE IF EXISTS agent_type_enum CASCADE;
DROP TYPE IF EXISTS payment_mode_enum CASCADE;
DROP TYPE IF EXISTS disbursement_status_enum CASCADE;
DROP TYPE IF EXISTS suspense_status_enum CASCADE;
DROP TYPE IF EXISTS clawback_status_enum CASCADE;
DROP TYPE IF EXISTS license_status_enum CASCADE;
DROP TYPE IF EXISTS reminder_status_enum CASCADE;

-- Drop extensions
DROP EXTENSION IF EXISTS "pg_trgm";
DROP EXTENSION IF EXISTS "pgcrypto";
DROP EXTENSION IF EXISTS "uuid-ossp";
