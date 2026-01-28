-- ============================================================================
-- Migration: 002_add_commission_batches
-- Description: Add commission_batches table for batch processing tracking
-- Author: Claude Code Generator
-- Date: 2026-01-28
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Table: commission_batches (E-IC-011)
-- Description: Commission batch processing tracking with Temporal workflow integration
-- Business Rules: BR-IC-COM-001, BR-IC-COM-012 (6-hour SLA)
-- ----------------------------------------------------------------------------
CREATE TABLE commission_batches (
    batch_id VARCHAR(50) PRIMARY KEY,
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    year INT NOT NULL CHECK (year BETWEEN 2020 AND 2099),
    status VARCHAR(50) NOT NULL CHECK (status IN ('INITIATED', 'CALCULATING', 'TRIAL_GENERATED', 'COMPLETED', 'FAILED')),
    workflow_id VARCHAR(255) NOT NULL,
    workflow_state VARCHAR(50),
    total_policies INT DEFAULT 0 CHECK (total_policies >= 0),
    processed_records INT DEFAULT 0 CHECK (processed_records >= 0),
    failed_records INT DEFAULT 0 CHECK (failed_records >= 0),
    progress_percentage INT DEFAULT 0 CHECK (progress_percentage BETWEEN 0 AND 100),
    triggered_by VARCHAR(50) NOT NULL CHECK (triggered_by IN ('SYSTEM_SCHEDULER', 'MANUAL')),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    estimated_completion TIMESTAMP WITH TIME ZONE,
    sla_deadline TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_batch_month_year_unique UNIQUE (month, year),
    CONSTRAINT chk_batch_progress_valid CHECK (processed_records + failed_records <= total_policies),
    CONSTRAINT chk_batch_completion_status CHECK (
        (status != 'COMPLETED') OR (completed_at IS NOT NULL)
    ),
    CONSTRAINT chk_batch_failure_reason CHECK (
        (status != 'FAILED') OR (failure_reason IS NOT NULL)
    ),
    CONSTRAINT chk_sla_deadline_after_start CHECK (sla_deadline > started_at)
);

COMMENT ON TABLE commission_batches IS 'E-IC-011: Commission batch processing tracker (BR-IC-COM-001, BR-IC-COM-012)';
COMMENT ON COLUMN commission_batches.sla_deadline IS 'BR-IC-COM-012: 6-hour SLA deadline';
COMMENT ON COLUMN commission_batches.workflow_id IS 'Temporal workflow ID for tracking';
COMMENT ON COLUMN commission_batches.progress_percentage IS 'Calculated as (processed_records * 100) / total_policies';

-- Indexes for commission_batches
CREATE INDEX idx_batch_status ON commission_batches(status)
    WHERE status NOT IN ('COMPLETED', 'FAILED');

CREATE INDEX idx_batch_workflow ON commission_batches(workflow_id);

CREATE INDEX idx_batch_month_year ON commission_batches(year DESC, month DESC);

CREATE INDEX idx_batch_sla ON commission_batches(sla_deadline)
    WHERE status NOT IN ('COMPLETED', 'FAILED');

-- Trigger for updated_at
CREATE TRIGGER trg_commission_batches_updated_at
    BEFORE UPDATE ON commission_batches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
