-- Rollback: Remove disbursements table
-- Version: 003

-- Drop indexes first
DROP INDEX IF EXISTS idx_disbursements_agent_status_date;
DROP INDEX IF EXISTS idx_disbursements_gl_pending;
DROP INDEX IF EXISTS idx_disbursements_failed_retry;
DROP INDEX IF EXISTS idx_disbursements_workflow_id;
DROP INDEX IF EXISTS idx_disbursements_completed_at;
DROP INDEX IF EXISTS idx_disbursements_initiated_at;
DROP INDEX IF EXISTS idx_disbursements_sla_breached;
DROP INDEX IF EXISTS idx_disbursements_sla_deadline;
DROP INDEX IF EXISTS idx_disbursements_mode;
DROP INDEX IF EXISTS idx_disbursements_status;
DROP INDEX IF EXISTS idx_disbursements_agent_id;

-- Drop table
DROP TABLE IF EXISTS disbursements;
