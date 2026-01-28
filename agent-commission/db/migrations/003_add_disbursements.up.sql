-- Migration: Add disbursements table
-- Version: 003
-- Description: Create disbursements table for commission payment processing
-- FR-IC-COM-009: Cheque disbursement
-- FR-IC-COM-010: EFT disbursement via PFMS
-- BR-IC-COM-011: 10-working-day SLA for disbursement

-- Create disbursements table
CREATE TABLE IF NOT EXISTS disbursements (
    disbursement_id BIGSERIAL PRIMARY KEY,

    -- Reference Information
    final_statement_id BIGINT NOT NULL,
    agent_id VARCHAR(50) NOT NULL,
    agent_code VARCHAR(50) NOT NULL,

    -- Disbursement Details
    disbursement_mode VARCHAR(20) NOT NULL CHECK (disbursement_mode IN ('CHEQUE', 'EFT')),
    disbursement_status VARCHAR(20) NOT NULL CHECK (disbursement_status IN ('PENDING', 'PROCESSING', 'SENT_TO_BANK', 'COMPLETED', 'FAILED', 'CANCELLED')),

    -- Financial Information
    total_gross_commission DECIMAL(15, 2) NOT NULL,
    total_tds DECIMAL(15, 2) NOT NULL DEFAULT 0,
    total_net_commission DECIMAL(15, 2) NOT NULL,

    -- Cheque Details (for CHEQUE mode)
    cheque_number VARCHAR(50),
    cheque_date DATE,

    -- EFT Details (for EFT mode)
    bank_account_number VARCHAR(50),
    bank_name VARCHAR(255),
    bank_branch VARCHAR(255),
    ifsc_code VARCHAR(11),
    account_holder_name VARCHAR(255),
    utr_number VARCHAR(50),  -- PFMS/Bank transaction reference

    -- SLA Tracking
    initiated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sla_deadline TIMESTAMP WITH TIME ZONE NOT NULL,  -- 10 working days
    processed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    sla_breached BOOLEAN NOT NULL DEFAULT FALSE,
    breach_reason_id BIGINT,

    -- Failure Tracking
    failure_reason VARCHAR(50) CHECK (failure_reason IN ('INVALID_ACCOUNT', 'INSUFFICIENT_FUNDS', 'BANK_REJECTION', 'NETWORK_ERROR', 'VALIDATION_ERROR')),
    failure_details TEXT,
    retry_count INT NOT NULL DEFAULT 0,

    -- Workflow Integration
    workflow_id VARCHAR(255),
    workflow_run_id VARCHAR(255),
    workflow_state VARCHAR(100),
    workflow_state_updated_at TIMESTAMP WITH TIME ZONE,

    -- Accounting Integration
    voucher_number VARCHAR(50),
    posted_to_gl BOOLEAN NOT NULL DEFAULT FALSE,
    gl_posted_at TIMESTAMP WITH TIME ZONE,

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INT NOT NULL DEFAULT 1,  -- Optimistic locking

    -- Constraints
    CONSTRAINT fk_disbursement_final_statement FOREIGN KEY (final_statement_id)
        REFERENCES final_statements(final_statement_id) ON DELETE RESTRICT,
    CONSTRAINT chk_disbursement_mode_details CHECK (
        (disbursement_mode = 'CHEQUE' AND cheque_number IS NOT NULL) OR
        (disbursement_mode = 'EFT' AND bank_account_number IS NOT NULL AND ifsc_code IS NOT NULL)
    )
);

-- Indexes for performance optimization

-- Primary queries: Search by agent, status, date range
CREATE INDEX idx_disbursements_agent_id ON disbursements(agent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_disbursements_status ON disbursements(disbursement_status) WHERE deleted_at IS NULL;
CREATE INDEX idx_disbursements_mode ON disbursements(disbursement_mode) WHERE deleted_at IS NULL;

-- SLA monitoring queries
CREATE INDEX idx_disbursements_sla_deadline ON disbursements(sla_deadline)
    WHERE disbursement_status NOT IN ('COMPLETED', 'CANCELLED') AND deleted_at IS NULL;
CREATE INDEX idx_disbursements_sla_breached ON disbursements(sla_breached)
    WHERE sla_breached = TRUE AND deleted_at IS NULL;

-- Date range queries
CREATE INDEX idx_disbursements_initiated_at ON disbursements(initiated_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_disbursements_completed_at ON disbursements(completed_at DESC) WHERE deleted_at IS NULL;

-- Workflow tracking
CREATE INDEX idx_disbursements_workflow_id ON disbursements(workflow_id) WHERE workflow_id IS NOT NULL;

-- Failed disbursements for retry
CREATE INDEX idx_disbursements_failed_retry ON disbursements(disbursement_status, retry_count)
    WHERE disbursement_status = 'FAILED' AND deleted_at IS NULL;

-- Accounting integration
CREATE INDEX idx_disbursements_gl_pending ON disbursements(posted_to_gl)
    WHERE posted_to_gl = FALSE AND disbursement_status = 'COMPLETED' AND deleted_at IS NULL;

-- Composite index for common filter combinations
CREATE INDEX idx_disbursements_agent_status_date ON disbursements(agent_id, disbursement_status, initiated_at DESC)
    WHERE deleted_at IS NULL;

-- Comments
COMMENT ON TABLE disbursements IS 'Commission disbursement records with SLA tracking and workflow integration';
COMMENT ON COLUMN disbursements.disbursement_mode IS 'Payment mode: CHEQUE or EFT (via PFMS)';
COMMENT ON COLUMN disbursements.sla_deadline IS 'BR-IC-COM-011: 10 working days from initiated_at';
COMMENT ON COLUMN disbursements.workflow_id IS 'Temporal workflow ID for WF-IC-DISB-001';
COMMENT ON COLUMN disbursements.version IS 'Optimistic locking version for concurrent updates';
