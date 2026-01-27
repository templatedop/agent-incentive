-- ============================================================================
-- Agent Incentive, Commission and Producer Management Database Schema
-- Database: incentive_db
-- PostgreSQL Version: 16
-- Schema Version: 1.0.0
-- Date: 2026-01-27
-- ============================================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================================
-- ENUM TYPES
-- ============================================================================

-- Agent Types (BR-IC-AH-001, BR-IC-AH-003, BR-IC-AH-004)
CREATE TYPE agent_type_enum AS ENUM (
    'ADVISOR',
    'ADVISOR_COORDINATOR',
    'DEPARTMENTAL_EMPLOYEE',
    'FIELD_OFFICER'
);

-- Person Type
CREATE TYPE person_type_enum AS ENUM (
    'INDIVIDUAL',
    'CORPORATE_GROUP'
);

-- Gender
CREATE TYPE gender_enum AS ENUM (
    'MALE',
    'FEMALE',
    'OTHER'
);

-- Marital Status
CREATE TYPE marital_status_enum AS ENUM (
    'SINGLE',
    'MARRIED',
    'DIVORCED',
    'WIDOWED'
);

-- Agent Status (BR-IC-PROF-004)
CREATE TYPE agent_status_enum AS ENUM (
    'ACTIVE',
    'SUSPENDED',
    'TERMINATED',
    'EXPIRED',
    'DEACTIVATED'
);

-- Commission Status (BR-IC-COM-002)
CREATE TYPE commission_status_enum AS ENUM (
    'CALCULATED',
    'TRIAL_PENDING',
    'TRIAL_APPROVED',
    'FINALIZED',
    'READY_FOR_DISBURSEMENT',
    'DISBURSED',
    'CANCELLED',
    'HELD_IN_SUSPENSE'
);

-- Commission Types
CREATE TYPE commission_type_enum AS ENUM (
    'FIRST_YEAR',
    'RENEWAL',
    'BONUS'
);

-- Product Type
CREATE TYPE product_type_enum AS ENUM (
    'PLI',
    'RPLI'
);

-- Address Type
CREATE TYPE address_type_enum AS ENUM (
    'OFFICIAL',
    'PERMANENT',
    'COMMUNICATION'
);

-- Contact Type
CREATE TYPE contact_type_enum AS ENUM (
    'OFFICIAL_LANDLINE',
    'RESIDENTIAL_LANDLINE',
    'MOBILE'
);

-- Email Type
CREATE TYPE email_type_enum AS ENUM (
    'OFFICIAL',
    'PERMANENT',
    'COMMUNICATION'
);

-- Payment Mode (BR-IC-COM-008)
CREATE TYPE payment_mode_enum AS ENUM (
    'CHEQUE',
    'EFT'
);

-- License Status
CREATE TYPE license_status_enum AS ENUM (
    'ACTIVE',
    'EXPIRED',
    'RENEWED',
    'CANCELLED'
);

-- Resident Status
CREATE TYPE resident_status_enum AS ENUM (
    'RESIDENT',
    'NON_RESIDENT'
);

-- Trial Statement Status
CREATE TYPE trial_statement_status_enum AS ENUM (
    'PENDING_APPROVAL',
    'APPROVED',
    'REJECTED'
);

-- Final Statement Status
CREATE TYPE final_statement_status_enum AS ENUM (
    'FINALIZED',
    'READY_FOR_DISBURSEMENT',
    'DISBURSEMENT_QUEUED',
    'DISBURSED'
);

-- Disbursement Status
CREATE TYPE disbursement_status_enum AS ENUM (
    'PENDING',
    'PROCESSING',
    'COMPLETED',
    'FAILED'
);

-- Commission Status for Agent (BR-IC-SUSPENSE-001)
CREATE TYPE commission_status_agent_enum AS ENUM (
    'ACTIVE',
    'SUSPENDED',
    'SUSPENDED_PENDING_CLAWBACK',
    'SUSPENDED_PENDING_INVESTIGATION'
);

-- Clawback Recovery Status
CREATE TYPE clawback_recovery_status_enum AS ENUM (
    'PENDING',
    'PARTIALLY_RECOVERED',
    'FULLY_RECOVERED',
    'WRITTEN_OFF'
);

-- Suspense Status
CREATE TYPE suspense_status_enum AS ENUM (
    'SUSPENDED',
    'RELEASED',
    'FORFEITED'
);

-- ============================================================================
-- TABLES
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Table: agent_profiles (E-IC-001)
-- Description: Master table for agent profile information
-- Business Rules: BR-IC-AH-001, BR-IC-AH-003, BR-IC-AH-004, BR-IC-PROF-004
-- ----------------------------------------------------------------------------
CREATE TABLE agent_profiles (
    agent_id VARCHAR(20) PRIMARY KEY,
    agent_type agent_type_enum NOT NULL,
    person_type person_type_enum NOT NULL,
    advisor_coordinator_id VARCHAR(20),
    profile_type VARCHAR(50),
    office_type VARCHAR(50),
    office_code VARCHAR(20),
    advisor_sub_type VARCHAR(50),
    effective_date DATE NOT NULL,
    distribution_channel VARCHAR(50),
    product_class VARCHAR(50),

    -- Circle and Division assignment (BR-IC-AH-002)
    circle_id VARCHAR(20),
    circle_name VARCHAR(100),
    division_id VARCHAR(20),
    division_name VARCHAR(100),

    title VARCHAR(10),
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    gender gender_enum NOT NULL,
    date_of_birth DATE NOT NULL,
    category VARCHAR(50),
    marital_status marital_status_enum,
    aadhar_number VARCHAR(12),
    pan VARCHAR(10) UNIQUE NOT NULL,
    designation VARCHAR(100),
    service_number VARCHAR(50),
    professional_title VARCHAR(100),
    status agent_status_enum NOT NULL DEFAULT 'ACTIVE',
    status_reason VARCHAR(255),
    status_date DATE,
    status_effective_date DATE,
    termination_date DATE,
    advisor_undergoing_training BOOLEAN DEFAULT FALSE,
    preferred_payment_mode payment_mode_enum,

    -- Clawback and Suspense tracking (BR-IC-SUSPENSE-001)
    clawback_pending_amount DECIMAL(15,2) DEFAULT 0,
    suspense_amount DECIMAL(15,2) DEFAULT 0,
    commission_status commission_status_agent_enum DEFAULT 'ACTIVE',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- Flexible metadata and full-text search
    metadata JSONB,
    search_vector tsvector,

    -- Constraints
    CONSTRAINT chk_agent_coordinator_self
        CHECK (agent_id != advisor_coordinator_id OR advisor_coordinator_id IS NULL),
    CONSTRAINT chk_dob_valid
        CHECK (date_of_birth < CURRENT_DATE),
    CONSTRAINT chk_age_minimum
        CHECK (EXTRACT(YEAR FROM AGE(date_of_birth)) >= 18),
    CONSTRAINT chk_pan_format
        CHECK (pan ~* '^[A-Z]{5}[0-9]{4}[A-Z]{1}$'),
    CONSTRAINT chk_aadhar_length
        CHECK (aadhar_number IS NULL OR aadhar_number ~ '^[0-9]{12}$'),
    CONSTRAINT chk_coordinator_required_for_advisor
        CHECK (agent_type != 'ADVISOR' OR advisor_coordinator_id IS NOT NULL),
    CONSTRAINT chk_circle_required_for_coordinator
        CHECK (agent_type != 'ADVISOR_COORDINATOR' OR circle_id IS NOT NULL),
    CONSTRAINT chk_clawback_amount_positive
        CHECK (clawback_pending_amount >= 0),
    CONSTRAINT chk_suspense_amount_positive
        CHECK (suspense_amount >= 0)
);

-- Comments for agent_profiles
COMMENT ON TABLE agent_profiles IS 'E-IC-001: Master agent profile table with all agent-related information';
COMMENT ON COLUMN agent_profiles.agent_id IS 'Unique agent code - primary key (VR-IC-PROF-002)';
COMMENT ON COLUMN agent_profiles.pan IS 'PAN number - must be unique (VR-IC-PROF-002)';
COMMENT ON COLUMN agent_profiles.date_of_birth IS 'Date of birth - must be 18+ (VR-IC-PROF-003)';
COMMENT ON COLUMN agent_profiles.advisor_coordinator_id IS 'Link to advisor coordinator (BR-IC-AH-001)';
COMMENT ON COLUMN agent_profiles.circle_id IS 'Circle assignment for coordinator (BR-IC-AH-002)';
COMMENT ON COLUMN agent_profiles.division_id IS 'Division assignment for coordinator (BR-IC-AH-002)';
COMMENT ON COLUMN agent_profiles.status IS 'Agent status (BR-IC-PROF-004)';
COMMENT ON COLUMN agent_profiles.clawback_pending_amount IS 'Total clawback pending recovery (BR-IC-SUSPENSE-001)';
COMMENT ON COLUMN agent_profiles.suspense_amount IS 'Total suspense amount held (BR-IC-SUSPENSE-001)';
COMMENT ON COLUMN agent_profiles.commission_status IS 'Commission payment status (BR-IC-SUSPENSE-001)';

-- ----------------------------------------------------------------------------
-- Table: agent_addresses (E-IC-002)
-- Description: Agent address information with historical tracking
-- ----------------------------------------------------------------------------
CREATE TABLE agent_addresses (
    address_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id) ON DELETE CASCADE,
    address_type address_type_enum NOT NULL,
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    village VARCHAR(100),
    taluka VARCHAR(100),
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100),
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) DEFAULT 'India',
    pin_code VARCHAR(10) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    effective_from DATE NOT NULL,
    effective_to DATE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Flexible metadata
    metadata JSONB,

    -- Constraints
    CONSTRAINT chk_pin_code_format CHECK (pin_code ~ '^[0-9]{6}$'),
    CONSTRAINT chk_address_date_range CHECK (effective_to IS NULL OR effective_to >= effective_from)
);

COMMENT ON TABLE agent_addresses IS 'E-IC-002: Agent address information with type-based classification';

-- ----------------------------------------------------------------------------
-- Table: agent_contacts (E-IC-003)
-- Description: Agent contact information (phone numbers)
-- ----------------------------------------------------------------------------
CREATE TABLE agent_contacts (
    contact_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id) ON DELETE CASCADE,
    contact_type contact_type_enum NOT NULL,
    contact_value VARCHAR(15) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_phone_format CHECK (contact_value ~ '^[0-9]{10}$'),
    CONSTRAINT chk_one_primary_contact_per_type
        CHECK (NOT is_primary OR NOT EXISTS (
            SELECT 1 FROM agent_contacts ac2
            WHERE ac2.agent_id = agent_contacts.agent_id
            AND ac2.contact_type = agent_contacts.contact_type
            AND ac2.is_primary = TRUE
            AND ac2.contact_id != agent_contacts.contact_id
        ))
);

COMMENT ON TABLE agent_contacts IS 'E-IC-003: Agent contact phone numbers';

-- ----------------------------------------------------------------------------
-- Table: agent_emails (E-IC-004)
-- Description: Agent email addresses
-- ----------------------------------------------------------------------------
CREATE TABLE agent_emails (
    email_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id) ON DELETE CASCADE,
    email_type email_type_enum NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT chk_email_format CHECK (email_address ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT chk_one_primary_email_per_type
        CHECK (NOT is_primary OR NOT EXISTS (
            SELECT 1 FROM agent_emails ae2
            WHERE ae2.agent_id = agent_emails.agent_id
            AND ae2.email_type = agent_emails.email_type
            AND ae2.is_primary = TRUE
            AND ae2.email_id != agent_emails.email_id
        ))
);

COMMENT ON TABLE agent_emails IS 'E-IC-004: Agent email addresses';

-- ----------------------------------------------------------------------------
-- Table: agent_bank_accounts (E-IC-005)
-- Description: Agent bank account details for commission disbursement
-- Business Rules: VR-IC-COM-004
-- ----------------------------------------------------------------------------
CREATE TABLE agent_bank_accounts (
    bank_account_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id) ON DELETE CASCADE,
    bank_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(30) NOT NULL,
    ifsc_code VARCHAR(11) NOT NULL,
    account_type VARCHAR(50),
    is_primary BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Constraints (VR-IC-COM-004)
    CONSTRAINT chk_ifsc_format CHECK (ifsc_code ~* '^[A-Z]{4}0[A-Z0-9]{6}$'),
    CONSTRAINT chk_one_primary_account CHECK (NOT is_primary OR NOT EXISTS (
        SELECT 1 FROM agent_bank_accounts aba2
        WHERE aba2.agent_id = agent_bank_accounts.agent_id
        AND aba2.is_primary = TRUE
        AND aba2.bank_account_id != agent_bank_accounts.bank_account_id
    ))
);

COMMENT ON TABLE agent_bank_accounts IS 'E-IC-005: Agent bank accounts for EFT disbursement (VR-IC-COM-004)';

-- ----------------------------------------------------------------------------
-- Table: agent_licenses (E-IC-006)
-- Description: Agent license information with renewal tracking
-- Business Rules: BR-IC-LIC-003, BR-IC-LIC-004
-- ----------------------------------------------------------------------------
CREATE TABLE agent_licenses (
    license_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id) ON DELETE CASCADE,
    license_line VARCHAR(50) NOT NULL,
    license_type VARCHAR(50) NOT NULL,
    license_number VARCHAR(50) NOT NULL,
    resident_status resident_status_enum NOT NULL,
    license_date DATE NOT NULL,
    renewal_date DATE NOT NULL,
    authority_date DATE,
    license_status license_status_enum NOT NULL DEFAULT 'ACTIVE',
    is_active BOOLEAN DEFAULT TRUE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Flexible metadata
    metadata JSONB,

    -- Constraints
    CONSTRAINT chk_renewal_date_after_license CHECK (renewal_date > license_date)
);

COMMENT ON TABLE agent_licenses IS 'E-IC-006: Agent license information with renewal tracking (BR-IC-LIC-004)';
COMMENT ON COLUMN agent_licenses.renewal_date IS 'License renewal date for reminder batch (BR-IC-LIC-004)';

-- ----------------------------------------------------------------------------
-- Table: commission_rates (E-IC-007)
-- Description: Commission rate configuration table
-- Business Rules: BR-IC-COM-006
-- ----------------------------------------------------------------------------
CREATE TABLE commission_rates (
    rate_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    rate_percentage DECIMAL(5,2) NOT NULL,
    policy_duration_months INT NOT NULL,
    product_type product_type_enum NOT NULL,
    product_plan_code VARCHAR(50) NOT NULL,
    agent_type VARCHAR(50) NOT NULL,
    policy_term_years INT NOT NULL,
    effective_from DATE NOT NULL,
    effective_to DATE,
    is_active BOOLEAN DEFAULT TRUE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_rate_percentage_positive CHECK (rate_percentage > 0),
    CONSTRAINT chk_rate_percentage_max CHECK (rate_percentage <= 100),
    CONSTRAINT chk_duration_positive CHECK (policy_duration_months > 0),
    CONSTRAINT chk_term_positive CHECK (policy_term_years > 0),
    CONSTRAINT chk_effective_date_range CHECK (effective_to IS NULL OR effective_to > effective_from)
);

COMMENT ON TABLE commission_rates IS 'E-IC-007: Commission rate table (BR-IC-COM-006)';

-- ----------------------------------------------------------------------------
-- Table: commission_transactions (E-IC-008)
-- Description: Commission calculation and transaction records
-- Business Rules: BR-IC-COM-001, BR-IC-COM-002, BR-IC-COM-003, BR-IC-COM-006
-- ----------------------------------------------------------------------------
CREATE TABLE commission_transactions (
    commission_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id),
    policy_number VARCHAR(50) NOT NULL,
    commission_type commission_type_enum NOT NULL,
    product_type product_type_enum NOT NULL,
    annualised_premium DECIMAL(15,2) NOT NULL,
    rate_percentage DECIMAL(5,2) NOT NULL,
    gross_commission DECIMAL(15,2) NOT NULL,
    tds_rate DECIMAL(5,2) DEFAULT 0,
    tds_amount DECIMAL(15,2) DEFAULT 0,
    net_commission DECIMAL(15,2) NOT NULL,
    commission_date DATE NOT NULL,
    commission_status commission_status_enum NOT NULL DEFAULT 'CALCULATED',
    trial_statement_id BIGINT,
    final_statement_id BIGINT,
    disbursement_id VARCHAR(50),
    batch_id VARCHAR(50),

    -- Workflow state tracking
    workflow_state VARCHAR(50),
    workflow_state_updated_at TIMESTAMP WITH TIME ZONE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- Flexible metadata and full-text search
    metadata JSONB,
    search_vector tsvector,

    -- Constraints
    CONSTRAINT chk_annualised_premium_positive CHECK (annualised_premium > 0),
    CONSTRAINT chk_rate_positive CHECK (rate_percentage > 0),
    CONSTRAINT chk_gross_commission_positive CHECK (gross_commission >= 0),
    CONSTRAINT chk_tds_rate_valid CHECK (tds_rate >= 0 AND tds_rate <= 100),
    CONSTRAINT chk_tds_amount_valid CHECK (tds_amount >= 0),
    CONSTRAINT chk_net_commission_calculation CHECK (net_commission = gross_commission - tds_amount),
    CONSTRAINT chk_disbursement_after_trial CHECK (
        commission_status NOT IN ('FINALIZED', 'READY_FOR_DISBURSEMENT', 'DISBURSED')
        OR trial_statement_id IS NOT NULL
    )
);

COMMENT ON TABLE commission_transactions IS 'E-IC-008: Commission calculation records (BR-IC-COM-001, BR-IC-COM-003)';
COMMENT ON COLUMN commission_transactions.commission_status IS 'Commission workflow status (BR-IC-COM-002)';
COMMENT ON COLUMN commission_transactions.tds_amount IS 'TDS deduction amount (BR-IC-COM-003)';

-- ----------------------------------------------------------------------------
-- Table: trial_statements (E-IC-009)
-- Description: Trial statement for finance review and approval
-- Business Rules: BR-IC-COM-002
-- ----------------------------------------------------------------------------
CREATE TABLE trial_statements (
    trial_statement_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    statement_number VARCHAR(50) UNIQUE NOT NULL,
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id),
    statement_date DATE NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    total_policies INT DEFAULT 0,
    total_gross_commission DECIMAL(15,2) DEFAULT 0,
    total_tds DECIMAL(15,2) DEFAULT 0,
    total_net_commission DECIMAL(15,2) DEFAULT 0,
    statement_status trial_statement_status_enum NOT NULL DEFAULT 'PENDING_APPROVAL',
    approved_by VARCHAR(50),
    approved_at TIMESTAMP WITH TIME ZONE,
    approval_remarks TEXT,
    processing_unit VARCHAR(50),
    batch_id VARCHAR(50),

    -- Workflow state tracking
    workflow_state VARCHAR(50),
    workflow_state_updated_at TIMESTAMP WITH TIME ZONE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- Flexible metadata and full-text search
    metadata JSONB,
    search_vector tsvector,

    -- Constraints
    CONSTRAINT chk_date_range_valid CHECK (to_date >= from_date),
    CONSTRAINT chk_totals_positive CHECK (
        total_policies >= 0 AND
        total_gross_commission >= 0 AND
        total_tds >= 0 AND
        total_net_commission >= 0
    ),
    CONSTRAINT chk_approval_consistency CHECK (
        (statement_status != 'APPROVED') OR
        (approved_by IS NOT NULL AND approved_at IS NOT NULL)
    )
);

COMMENT ON TABLE trial_statements IS 'E-IC-009: Trial statements for finance approval (BR-IC-COM-002)';

-- ----------------------------------------------------------------------------
-- Table: final_statements (E-IC-010)
-- Description: Final commission statement after trial approval
-- Business Rules: BR-IC-COM-007
-- ----------------------------------------------------------------------------
CREATE TABLE final_statements (
    final_statement_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    statement_number VARCHAR(50) UNIQUE NOT NULL,
    trial_statement_id BIGINT NOT NULL REFERENCES trial_statements(trial_statement_id),
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id),
    statement_date DATE NOT NULL,
    total_gross_commission DECIMAL(15,2) DEFAULT 0,
    total_tds DECIMAL(15,2) DEFAULT 0,
    total_net_commission DECIMAL(15,2) DEFAULT 0,
    statement_status final_statement_status_enum NOT NULL DEFAULT 'FINALIZED',
    pdf_path VARCHAR(255),

    -- Workflow state tracking
    workflow_state VARCHAR(50),
    workflow_state_updated_at TIMESTAMP WITH TIME ZONE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- Flexible metadata and full-text search
    metadata JSONB,
    search_vector tsvector,

    -- Constraints
    CONSTRAINT chk_totals_final_positive CHECK (
        total_gross_commission >= 0 AND
        total_tds >= 0 AND
        total_net_commission >= 0
    )
);

COMMENT ON TABLE final_statements IS 'E-IC-010: Final commission statements (BR-IC-COM-007)';

-- ----------------------------------------------------------------------------
-- Table: disbursements (E-IC-011)
-- Description: Commission disbursement tracking
-- Business Rules: BR-IC-COM-008, BR-IC-COM-011
-- ----------------------------------------------------------------------------
CREATE TABLE disbursements (
    disbursement_id VARCHAR(50) PRIMARY KEY,
    final_statement_id BIGINT NOT NULL REFERENCES final_statements(final_statement_id),
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id),
    payment_mode payment_mode_enum NOT NULL,
    disbursement_amount DECIMAL(15,2) NOT NULL,
    disbursement_status disbursement_status_enum NOT NULL DEFAULT 'PENDING',
    cheque_number VARCHAR(50),
    bank_name VARCHAR(255),
    ifsc_code VARCHAR(11),
    account_number VARCHAR(30),
    payment_date DATE,
    payment_reference VARCHAR(100),
    pfms_transaction_id VARCHAR(100),
    trial_approval_date DATE,
    sla_breach BOOLEAN DEFAULT FALSE,
    penalty_amount DECIMAL(15,2) DEFAULT 0,
    remarks TEXT,

    -- Partial disbursement tracking (BR-IC-COM-005)
    disbursement_mode VARCHAR(20) DEFAULT 'FULL',
    disbursement_percentage DECIMAL(5,2) DEFAULT 100,
    pending_amount DECIMAL(15,2) DEFAULT 0,

    -- Payment retry tracking (BR-IC-SUSPENSE-002)
    retry_count INT DEFAULT 0,
    retry_status VARCHAR(50),
    last_retry_date DATE,
    next_retry_date DATE,
    max_retries INT DEFAULT 3,

    -- Workflow state tracking
    workflow_state VARCHAR(50),
    workflow_state_updated_at TIMESTAMP WITH TIME ZONE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- Flexible metadata
    metadata JSONB,

    -- Constraints (BR-IC-COM-011)
    CONSTRAINT chk_disbursement_amount_positive CHECK (disbursement_amount > 0),
    CONSTRAINT chk_penalty_amount_positive CHECK (penalty_amount >= 0),
    CONSTRAINT chk_eft_requires_bank_details CHECK (
        payment_mode != 'EFT' OR
        (ifsc_code IS NOT NULL AND account_number IS NOT NULL)
    ),
    CONSTRAINT chk_cheque_requires_cheque_number CHECK (
        payment_mode != 'CHEQUE' OR
        cheque_number IS NOT NULL
    ),
    CONSTRAINT chk_disbursement_mode_valid CHECK (disbursement_mode IN ('FULL', 'PARTIAL')),
    CONSTRAINT chk_disbursement_percentage_valid CHECK (
        disbursement_mode = 'FULL' AND disbursement_percentage = 100 OR
        disbursement_mode = 'PARTIAL' AND disbursement_percentage BETWEEN 0 AND 100
    ),
    CONSTRAINT chk_retry_count_max CHECK (retry_count <= max_retries)
);

COMMENT ON TABLE disbursements IS 'E-IC-011: Commission disbursement tracking (BR-IC-COM-008, BR-IC-COM-011)';
COMMENT ON COLUMN disbursements.trial_approval_date IS 'For SLA calculation (BR-IC-COM-011)';
COMMENT ON COLUMN disbursements.sla_breach IS 'SLA breach flag (10 working days) (BR-IC-COM-011)';
COMMENT ON COLUMN disbursements.disbursement_mode IS 'FULL or PARTIAL disbursement (BR-IC-COM-005)';
COMMENT ON COLUMN disbursements.disbursement_percentage IS 'Percentage for partial disbursement (BR-IC-COM-005)';
COMMENT ON COLUMN disbursements.pending_amount IS 'Remaining unpaid amount for partial disbursement (BR-IC-COM-005)';
COMMENT ON COLUMN disbursements.retry_count IS 'Number of retry attempts for failed payments (BR-IC-SUSPENSE-002)';
COMMENT ON COLUMN disbursements.retry_status IS 'Current retry status (BR-IC-SUSPENSE-002)';

-- ----------------------------------------------------------------------------
-- Table: commission_history (E-IC-012)
-- Description: Historical commission transaction records
-- Business Rules: BR-IC-COM-009
-- ----------------------------------------------------------------------------
CREATE TABLE commission_history (
    history_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    commission_id BIGINT NOT NULL REFERENCES commission_transactions(commission_id),
    agent_id VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id),
    policy_number VARCHAR(50) NOT NULL,
    commission_type commission_type_enum NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    tds_amount DECIMAL(15,2) DEFAULT 0,
    net_amount DECIMAL(15,2) NOT NULL,
    product_type product_type_enum NOT NULL,
    transaction_date DATE NOT NULL,
    processed_date DATE NOT NULL,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Flexible metadata
    metadata JSONB,

    -- Constraints
    CONSTRAINT chk_history_amount_positive CHECK (amount > 0),
    CONSTRAINT chk_history_net_calculation CHECK (net_amount = amount - tds_amount)
);

COMMENT ON TABLE commission_history IS 'E-IC-012: Commission history for search and inquiry (BR-IC-COM-009)';

-- ----------------------------------------------------------------------------
-- Table: commission_clawbacks (E-IC-013)
-- Description: Commission clawback tracking for lapsed policies
-- Business Rules: BR-IC-CLAWBACK-001
-- ----------------------------------------------------------------------------
CREATE TABLE commission_clawbacks (
    clawback_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    policy_number VARCHAR(50) NOT NULL,
    agent_code VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id),
    policy_issue_date DATE NOT NULL,
    lapse_date DATE NOT NULL,
    months_active INT NOT NULL,
    original_commission DECIMAL(15,2) NOT NULL,
    clawback_percentage DECIMAL(5,2) NOT NULL,
    clawback_amount DECIMAL(15,2) NOT NULL,
    recovery_status clawback_recovery_status_enum NOT NULL DEFAULT 'PENDING',
    recovered_amount DECIMAL(15,2) DEFAULT 0,
    created_date DATE NOT NULL,
    clawback_reason VARCHAR(255),

    -- Workflow state tracking
    workflow_state VARCHAR(50),
    workflow_state_updated_at TIMESTAMP WITH TIME ZONE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- Flexible metadata
    metadata JSONB,

    -- Constraints
    CONSTRAINT chk_months_active_positive CHECK (months_active > 0),
    CONSTRAINT chk_original_commission_positive CHECK (original_commission > 0),
    CONSTRAINT chk_clawback_percentage_valid CHECK (
        clawback_percentage IN (25, 50, 75, 100)
    ),
    CONSTRAINT chk_clawback_amount_positive CHECK (clawback_amount >= 0),
    CONSTRAINT chk_recovered_amount_valid CHECK (
        recovered_amount >= 0 AND recovered_amount <= clawback_amount
    ),
    CONSTRAINT chk_lapse_after_issue CHECK (lapse_date > policy_issue_date)
);

COMMENT ON TABLE commission_clawbacks IS 'E-IC-013: Commission clawback records (BR-IC-CLAWBACK-001)';

-- ----------------------------------------------------------------------------
-- Table: commission_suspense (E-IC-014)
-- Description: Commission suspense for disputed/investigation policies
-- Business Rules: BR-IC-SUSPENSE-001
-- ----------------------------------------------------------------------------
CREATE TABLE commission_suspense (
    suspense_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    policy_number VARCHAR(50) NOT NULL,
    agent_code VARCHAR(20) NOT NULL REFERENCES agent_profiles(agent_id),
    commission_amount DECIMAL(15,2) NOT NULL,
    suspense_reason VARCHAR(100) NOT NULL,
    investigation_type VARCHAR(50),
    suspense_date DATE NOT NULL,
    expected_resolution_date DATE,
    status suspense_status_enum NOT NULL DEFAULT 'SUSPENDED',
    release_date DATE,
    release_reason TEXT,
    forfeiture_date DATE,
    forfeiture_reason TEXT,
    investigation_reference VARCHAR(100),

    -- Workflow state tracking
    workflow_state VARCHAR(50),
    workflow_state_updated_at TIMESTAMP WITH TIME ZONE,

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(50),
    updated_by VARCHAR(50),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1,

    -- Flexible metadata
    metadata JSONB,

    -- Constraints
    CONSTRAINT chk_commission_amount_positive CHECK (commission_amount > 0),
    CONSTRAINT chk_expected_resolution_after_suspense CHECK (
        expected_resolution_date IS NULL OR expected_resolution_date >= suspense_date
    ),
    CONSTRAINT chk_release_after_suspense CHECK (
        release_date IS NULL OR release_date >= suspense_date
    ),
    CONSTRAINT chk_forfeiture_after_suspense CHECK (
        forfeiture_date IS NULL OR forfeiture_date >= suspense_date
    )
);

COMMENT ON TABLE commission_suspense IS 'E-IC-014: Commission suspense for disputed policies (BR-IC-SUSPENSE-001)';

-- ============================================================================
-- INDEXES
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Indexes for agent_profiles
-- ----------------------------------------------------------------------------
CREATE INDEX idx_agent_pan ON agent_profiles(pan) WHERE deleted_at IS NULL;
CREATE INDEX idx_agent_status ON agent_profiles(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_agent_coordinator ON agent_profiles(advisor_coordinator_id)
    WHERE advisor_coordinator_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_agent_name ON agent_profiles(first_name, last_name)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_agent_type ON agent_profiles(agent_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_agent_commission_status ON agent_profiles(commission_status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_agent_clawback_pending ON agent_profiles(clawback_pending_amount)
    WHERE clawback_pending_amount > 0 AND deleted_at IS NULL;
CREATE INDEX idx_agent_suspense_amount ON agent_profiles(suspense_amount)
    WHERE suspense_amount > 0 AND deleted_at IS NULL;
CREATE INDEX idx_agent_metadata ON agent_profiles USING gin(metadata);
CREATE INDEX idx_agent_search ON agent_profiles USING gin(search_vector);

-- ----------------------------------------------------------------------------
-- Indexes for agent_addresses
-- ----------------------------------------------------------------------------
CREATE INDEX idx_address_agent ON agent_addresses(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_address_type ON agent_addresses(address_type)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_address_active ON agent_addresses(is_active)
    WHERE is_active = TRUE AND deleted_at IS NULL;

-- ----------------------------------------------------------------------------
-- Indexes for agent_contacts
-- ----------------------------------------------------------------------------
CREATE INDEX idx_contact_agent ON agent_contacts(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_contact_type ON agent_contacts(contact_type)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_contact_value ON agent_contacts(contact_value)
    WHERE is_active = TRUE AND deleted_at IS NULL;

-- ----------------------------------------------------------------------------
-- Indexes for agent_emails
-- ----------------------------------------------------------------------------
CREATE INDEX idx_email_agent ON agent_emails(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_email_type ON agent_emails(email_type)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_email_address ON agent_emails(email_address)
    WHERE is_active = TRUE AND deleted_at IS NULL;

-- ----------------------------------------------------------------------------
-- Indexes for agent_bank_accounts
-- ----------------------------------------------------------------------------
CREATE INDEX idx_bank_agent ON agent_bank_accounts(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_bank_primary ON agent_bank_accounts(is_primary)
    WHERE is_primary = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_bank_active ON agent_bank_accounts(is_active)
    WHERE is_active = TRUE AND deleted_at IS NULL;

-- ----------------------------------------------------------------------------
-- Indexes for agent_licenses (BR-IC-LIC-004)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_license_agent ON agent_licenses(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_license_number ON agent_licenses(license_number)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_license_renewal_date ON agent_licenses(renewal_date)
    WHERE license_status = 'ACTIVE' AND deleted_at IS NULL;
CREATE INDEX idx_license_status ON agent_licenses(license_status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_license_expiry_reminders ON agent_licenses(renewal_date)
    WHERE license_status = 'ACTIVE'
    AND renewal_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '90 days'
    AND deleted_at IS NULL;

-- ----------------------------------------------------------------------------
-- Indexes for commission_rates (BR-IC-COM-006)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_rate_lookup ON commission_rates(
    product_type, product_plan_code, agent_type, policy_term_years, policy_duration_months
) WHERE is_active = TRUE;
CREATE INDEX idx_rate_active ON commission_rates(is_active);

-- ----------------------------------------------------------------------------
-- Indexes for commission_transactions
-- ----------------------------------------------------------------------------
CREATE INDEX idx_commission_agent ON commission_transactions(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_commission_policy ON commission_transactions(policy_number)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_commission_status ON commission_transactions(commission_status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_commission_date ON commission_transactions(commission_date)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_commission_batch ON commission_transactions(batch_id)
    WHERE batch_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_commission_trial ON commission_transactions(trial_statement_id)
    WHERE trial_statement_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_commission_type ON commission_transactions(commission_type)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_commission_workflow_state ON commission_transactions(workflow_state)
    WHERE workflow_state IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_commission_metadata ON commission_transactions USING gin(metadata);
CREATE INDEX idx_commission_search ON commission_transactions USING gin(search_vector);

-- ----------------------------------------------------------------------------
-- Indexes for trial_statements
-- ----------------------------------------------------------------------------
CREATE INDEX idx_trial_agent ON trial_statements(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_trial_status ON trial_statements(statement_status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_trial_date ON trial_statements(statement_date)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_trial_number ON trial_statements(statement_number)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_trial_batch ON trial_statements(batch_id)
    WHERE batch_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_trial_approval_pending ON trial_statements(agent_id, statement_date)
    WHERE statement_status = 'PENDING_APPROVAL' AND deleted_at IS NULL;
CREATE INDEX idx_trial_workflow_state ON trial_statements(workflow_state)
    WHERE workflow_state IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_trial_metadata ON trial_statements USING gin(metadata);
CREATE INDEX idx_trial_search ON trial_statements USING gin(search_vector);

-- ----------------------------------------------------------------------------
-- Indexes for final_statements
-- ----------------------------------------------------------------------------
CREATE INDEX idx_final_agent ON final_statements(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_final_trial ON final_statements(trial_statement_id);
CREATE INDEX idx_final_status ON final_statements(statement_status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_final_number ON final_statements(statement_number)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_final_workflow_state ON final_statements(workflow_state)
    WHERE workflow_state IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_final_metadata ON final_statements USING gin(metadata);
CREATE INDEX idx_final_search ON final_statements USING gin(search_vector);

-- ----------------------------------------------------------------------------
-- Indexes for disbursements (BR-IC-COM-011)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_disb_agent ON disbursements(agent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_disb_final ON disbursements(final_statement_id);
CREATE INDEX idx_disb_status ON disbursements(disbursement_status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_disb_trial_approval_date ON disbursements(trial_approval_date)
    WHERE trial_approval_date IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_disb_sla_breach ON disbursements(sla_breach, trial_approval_date)
    WHERE sla_breach = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_disb_payment_mode ON disbursements(payment_mode)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_disb_pfms_transaction ON disbursements(pfms_transaction_id)
    WHERE pfms_transaction_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_disb_workflow_state ON disbursements(workflow_state)
    WHERE workflow_state IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_disb_metadata ON disbursements USING gin(metadata);

-- ----------------------------------------------------------------------------
-- Indexes for commission_history (BR-IC-COM-009)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_hist_agent ON commission_history(agent_id);
CREATE INDEX idx_hist_policy ON commission_history(policy_number);
CREATE INDEX idx_hist_date ON commission_history(transaction_date);
CREATE INDEX idx_hist_type ON commission_history(commission_type);
CREATE INDEX idx_hist_product_type ON commission_history(product_type);
CREATE INDEX idx_hist_agent_date ON commission_history(agent_id, transaction_date);
CREATE INDEX idx_hist_metadata ON commission_history USING gin(metadata);

-- ----------------------------------------------------------------------------
-- Indexes for commission_clawbacks
-- ----------------------------------------------------------------------------
CREATE INDEX idx_clawback_agent ON commission_clawbacks(agent_code)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_clawback_policy ON commission_clawbacks(policy_number)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_clawback_status ON commission_clawbacks(recovery_status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_clawback_date ON commission_clawbacks(created_date)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_clawback_pending ON commission_clawbacks(agent_code, recovery_status)
    WHERE recovery_status IN ('PENDING', 'PARTIALLY_RECOVERED') AND deleted_at IS NULL;
CREATE INDEX idx_clawback_workflow_state ON commission_clawbacks(workflow_state)
    WHERE workflow_state IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_clawback_metadata ON commission_clawbacks USING gin(metadata);

-- ----------------------------------------------------------------------------
-- Indexes for commission_suspense (BR-IC-SUSPENSE-001)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_suspense_agent ON commission_suspense(agent_code)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_suspense_policy ON commission_suspense(policy_number)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_suspense_status ON commission_suspense(status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_suspense_date ON commission_suspense(suspense_date)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_suspense_pending ON commission_suspense(agent_code, status)
    WHERE status = 'SUSPENDED' AND deleted_at IS NULL;
CREATE INDEX idx_suspense_expected_resolution ON commission_suspense(expected_resolution_date)
    WHERE status = 'SUSPENDED' AND deleted_at IS NULL;
CREATE INDEX idx_suspense_workflow_state ON commission_suspense(workflow_state)
    WHERE workflow_state IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_suspense_metadata ON commission_suspense USING gin(metadata);

-- ----------------------------------------------------------------------------
-- Indexes for circle/division (BR-IC-AH-002)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_agent_circle ON agent_profiles(circle_id)
    WHERE circle_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_agent_division ON agent_profiles(division_id)
    WHERE division_id IS NOT NULL AND deleted_at IS NULL;

-- ----------------------------------------------------------------------------
-- Indexes for partial disbursement (BR-IC-COM-005)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_disb_mode ON disbursements(disbursement_mode)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_disb_partial ON disbursements(disbursement_mode, pending_amount)
    WHERE disbursement_mode = 'PARTIAL' AND pending_amount > 0 AND deleted_at IS NULL;

-- ----------------------------------------------------------------------------
-- Indexes for payment retry (BR-IC-SUSPENSE-002)
-- ----------------------------------------------------------------------------
CREATE INDEX idx_disb_retry_count ON disbursements(retry_count)
    WHERE retry_count > 0 AND deleted_at IS NULL;
CREATE INDEX idx_disb_retry_status ON disbursements(retry_status)
    WHERE retry_status IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_disb_next_retry ON disbursements(next_retry_date)
    WHERE retry_status IS NOT NULL AND next_retry_date IS NOT NULL AND deleted_at IS NULL;

-- ============================================================================
-- FOREIGN KEY CONSTRAINTS (ADDITIONAL)
-- ============================================================================

-- Foreign key from commission_transactions to trial_statements
ALTER TABLE commission_transactions
    ADD CONSTRAINT fk_commission_trial_statement
    FOREIGN KEY (trial_statement_id)
    REFERENCES trial_statements(trial_statement_id);

-- Foreign key from commission_transactions to final_statements
ALTER TABLE commission_transactions
    ADD CONSTRAINT fk_commission_final_statement
    FOREIGN KEY (final_statement_id)
    REFERENCES final_statements(final_statement_id);

-- Foreign key from commission_transactions to disbursements
ALTER TABLE commission_transactions
    ADD CONSTRAINT fk_commission_disbursement
    FOREIGN KEY (disbursement_id)
    REFERENCES disbursements(disbursement_id);

-- ============================================================================
-- VIEWS
-- ============================================================================

-- ----------------------------------------------------------------------------
-- View: v_active_agents
-- Description: Active agents with license and commission status
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_active_agents AS
SELECT
    ap.agent_id,
    ap.agent_type,
    ap.first_name,
    ap.last_name,
    ap.pan,
    ap.status,
    ap.commission_status,
    ap.clawback_pending_amount,
    ap.suspense_amount,
    ap.circle_id,
    ap.circle_name,
    ap.division_id,
    ap.division_name,
    COUNT(DISTINCT al.license_id) FILTER (
        WHERE al.license_status = 'ACTIVE' AND al.is_active = TRUE
    ) as active_licenses_count,
    MIN(al.renewal_date) FILTER (
        WHERE al.license_status = 'ACTIVE' AND al.is_active = TRUE
    ) as next_renewal_date,
    COUNT(ct.commission_id) FILTER (
        WHERE ct.commission_status = 'DISBURSED'
        AND ct.commission_date >= CURRENT_DATE - INTERVAL '12 months'
    ) as commissions_last_12_months,
    COALESCE(SUM(ct.net_commission) FILTER (
        WHERE ct.commission_status = 'DISBURSED'
        AND ct.commission_date >= CURRENT_DATE - INTERVAL '12 months'
    ), 0) as total_commission_12_months
FROM agent_profiles ap
LEFT JOIN agent_licenses al ON ap.agent_id = al.agent_id AND al.deleted_at IS NULL
LEFT JOIN commission_transactions ct ON ap.agent_id = ct.agent_id AND ct.deleted_at IS NULL
WHERE ap.status = 'ACTIVE'
    AND ap.deleted_at IS NULL
GROUP BY ap.agent_id, ap.agent_type, ap.first_name, ap.last_name, ap.pan,
         ap.status, ap.commission_status, ap.clawback_pending_amount, ap.suspense_amount,
         ap.circle_id, ap.circle_name, ap.division_id, ap.division_name;

COMMENT ON VIEW v_active_agents IS 'Active agents with license and commission summary';

-- ----------------------------------------------------------------------------
-- View: v_trial_approval_queue
-- Description: Trial statements pending approval
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_trial_approval_queue AS
SELECT
    ts.trial_statement_id,
    ts.statement_number,
    ts.agent_id,
    ap.first_name,
    ap.last_name,
    ts.statement_date,
    ts.from_date,
    ts.to_date,
    ts.total_policies,
    ts.total_gross_commission,
    ts.total_net_commission,
    ts.statement_status,
    EXTRACT(DAY FROM (CURRENT_DATE - ts.statement_date)) as days_pending,
    ts.batch_id
FROM trial_statements ts
JOIN agent_profiles ap ON ts.agent_id = ap.agent_id
WHERE ts.statement_status = 'PENDING_APPROVAL'
    AND ts.deleted_at IS NULL
    AND ap.deleted_at IS NULL
ORDER BY ts.statement_date ASC;

COMMENT ON VIEW v_trial_approval_queue IS 'Trial statements pending finance approval';

-- ----------------------------------------------------------------------------
-- View: v_disbursement_sla_tracking
-- Description: Disbursements with SLA tracking (BR-IC-COM-011)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_disbursement_sla_tracking AS
SELECT
    d.disbursement_id,
    d.agent_id,
    ap.first_name,
    ap.last_name,
    d.payment_mode,
    d.disbursement_amount,
    d.disbursement_status,
    d.trial_approval_date,
    d.payment_date,
    d.sla_breach,
    d.penalty_amount,
    EXTRACT(DAY FROM (COALESCE(d.payment_date, CURRENT_DATE) - d.trial_approval_date)) as days_since_approval,
    CASE
        WHEN d.payment_date IS NOT NULL THEN
            EXTRACT(DAY FROM (d.payment_date - d.trial_approval_date))
        ELSE
            EXTRACT(DAY FROM (CURRENT_DATE - d.trial_approval_date))
    END as processing_days,
    CASE
        WHEN d.payment_date IS NOT NULL AND
             EXTRACT(DAY FROM (d.payment_date - d.trial_approval_date)) > 10 THEN TRUE
        WHEN d.payment_date IS NULL AND
             EXTRACT(DAY FROM (CURRENT_DATE - d.trial_approval_date)) > 10 THEN TRUE
        ELSE FALSE
    END as sla_breached_calculated,
    fs.statement_number
FROM disbursements d
JOIN agent_profiles ap ON d.agent_id = ap.agent_id
JOIN final_statements fs ON d.final_statement_id = fs.final_statement_id
WHERE d.deleted_at IS NULL
    AND ap.deleted_at IS NULL
    AND d.trial_approval_date IS NOT NULL
ORDER BY d.trial_approval_date DESC;

COMMENT ON VIEW v_disbursement_sla_tracking IS 'Disbursement SLA tracking (10 working days) (BR-IC-COM-011)';

-- ----------------------------------------------------------------------------
-- View: v_license_expiry_report
-- Description: License expiry tracking (BR-IC-LIC-004)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_license_expiry_report AS
SELECT
    al.license_id,
    al.agent_id,
    ap.first_name,
    ap.last_name,
    ap.agent_type,
    al.license_number,
    al.license_type,
    al.license_line,
    al.renewal_date,
    al.license_status,
    EXTRACT(DAY FROM (al.renewal_date - CURRENT_DATE)) as days_until_expiry,
    CASE
        WHEN EXTRACT(DAY FROM (al.renewal_date - CURRENT_DATE)) <= 0 THEN 'EXPIRED'
        WHEN EXTRACT(DAY FROM (al.renewal_date - CURRENT_DATE)) <= 7 THEN 'T-7_DAYS'
        WHEN EXTRACT(DAY FROM (al.renewal_date - CURRENT_DATE)) <= 15 THEN 'T-15_DAYS'
        WHEN EXTRACT(DAY FROM (al.renewal_date - CURRENT_DATE)) <= 30 THEN 'T-30_DAYS'
        ELSE 'UPCOMING'
    END as urgency_level,
    ap.status as agent_status,
    ap.commission_status
FROM agent_licenses al
JOIN agent_profiles ap ON al.agent_id = ap.agent_id
WHERE al.is_active = TRUE
    AND al.deleted_at IS NULL
    AND ap.deleted_at IS NULL
    AND al.renewal_date <= CURRENT_DATE + INTERVAL '90 days'
ORDER BY al.renewal_date ASC;

COMMENT ON VIEW v_license_expiry_report IS 'License expiry tracking for renewal reminders (BR-IC-LIC-004)';

-- ----------------------------------------------------------------------------
-- View: v_sla_breach_report
-- Description: SLA breach summary
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_sla_breach_report AS
SELECT
    d.disbursement_id,
    d.agent_id,
    ap.first_name,
    ap.last_name,
    d.payment_mode,
    d.disbursement_amount,
    d.trial_approval_date,
    d.payment_date,
    d.sla_breach,
    d.penalty_amount,
    EXTRACT(DAY FROM (COALESCE(d.payment_date, CURRENT_DATE) - d.trial_approval_date)) as days_taken,
    fs.statement_number,
    ts.from_date,
    ts.to_date
FROM disbursements d
JOIN agent_profiles ap ON d.agent_id = ap.agent_id
JOIN final_statements fs ON d.final_statement_id = fs.final_statement_id
JOIN trial_statements ts ON fs.trial_statement_id = ts.trial_statement_id
WHERE d.deleted_at IS NULL
    AND ap.deleted_at IS NULL
    AND d.sla_breach = TRUE
ORDER BY d.trial_approval_date DESC;

COMMENT ON VIEW v_sla_breach_report IS 'SLA breach report for penalty calculation';

-- ----------------------------------------------------------------------------
-- View: v_commission_summary_by_agent
-- Description: Commission summary by agent
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_commission_summary_by_agent AS
SELECT
    ct.agent_id,
    ap.first_name,
    ap.last_name,
    ap.agent_type,
    ct.commission_type,
    DATE_TRUNC('month', ct.commission_date) as commission_month,
    COUNT(*) as transaction_count,
    SUM(ct.gross_commission) as total_gross_commission,
    SUM(ct.tds_amount) as total_tds,
    SUM(ct.net_commission) as total_net_commission,
    COUNT(DISTINCT ct.policy_number) as unique_policies
FROM commission_transactions ct
JOIN agent_profiles ap ON ct.agent_id = ap.agent_id
WHERE ct.commission_status = 'DISBURSED'
    AND ct.deleted_at IS NULL
    AND ap.deleted_at IS NULL
GROUP BY ct.agent_id, ap.first_name, ap.last_name, ap.agent_type,
         ct.commission_type, DATE_TRUNC('month', ct.commission_date)
ORDER BY commission_month DESC, total_net_commission DESC;

COMMENT ON VIEW v_commission_summary_by_agent IS 'Monthly commission summary by agent';

-- ----------------------------------------------------------------------------
-- View: v_clawback_pending_recovery
-- Description: Clawbacks pending recovery
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_clawback_pending_recovery AS
SELECT
    cc.clawback_id,
    cc.agent_code,
    ap.first_name,
    ap.last_name,
    cc.policy_number,
    cc.policy_issue_date,
    cc.lapse_date,
    cc.months_active,
    cc.original_commission,
    cc.clawback_percentage,
    cc.clawback_amount,
    cc.recovered_amount,
    (cc.clawback_amount - cc.recovered_amount) as pending_recovery,
    cc.recovery_status,
    cc.created_date,
    EXTRACT(DAY FROM (CURRENT_DATE - cc.created_date)) as days_since_created
FROM commission_clawbacks cc
JOIN agent_profiles ap ON cc.agent_code = ap.agent_id
WHERE cc.recovery_status IN ('PENDING', 'PARTIALLY_RECOVERED')
    AND cc.deleted_at IS NULL
    AND ap.deleted_at IS NULL
ORDER BY cc.created_date ASC;

COMMENT ON VIEW v_clawback_pending_recovery IS 'Clawbacks pending recovery';

-- ----------------------------------------------------------------------------
-- View: v_suspense_pending_release
-- Description: Suspense entries pending resolution
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_suspense_pending_release AS
SELECT
    cs.suspense_id,
    cs.agent_code,
    ap.first_name,
    ap.last_name,
    cs.policy_number,
    cs.commission_amount,
    cs.suspense_reason,
    cs.investigation_type,
    cs.suspense_date,
    cs.expected_resolution_date,
    cs.status,
    cs.investigation_reference,
    EXTRACT(DAY FROM (CURRENT_DATE - cs.suspense_date)) as days_in_suspense,
    CASE
        WHEN cs.expected_resolution_date < CURRENT_DATE THEN TRUE
        ELSE FALSE
    END as overdue
FROM commission_suspense cs
JOIN agent_profiles ap ON cs.agent_code = ap.agent_id
WHERE cs.status = 'SUSPENDED'
    AND cs.deleted_at IS NULL
    AND ap.deleted_at IS NULL
ORDER BY cs.suspense_date ASC;

COMMENT ON VIEW v_suspense_pending_release IS 'Commission suspense pending resolution (BR-IC-SUSPENSE-001)';

-- ----------------------------------------------------------------------------
-- View: v_partial_disbursements
-- Description: Partial disbursements with pending amounts (BR-IC-COM-005)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_partial_disbursements AS
SELECT
    d.disbursement_id,
    d.agent_id,
    ap.first_name,
    ap.last_name,
    d.final_statement_id,
    fs.statement_number,
    d.disbursement_amount,
    d.disbursement_mode,
    d.disbursement_percentage,
    d.pending_amount,
    d.payment_mode,
    d.disbursement_status,
    d.created_at
FROM disbursements d
JOIN agent_profiles ap ON d.agent_id = ap.agent_id
JOIN final_statements fs ON d.final_statement_id = fs.final_statement_id
WHERE d.disbursement_mode = 'PARTIAL'
    AND d.pending_amount > 0
    AND d.deleted_at IS NULL
    AND ap.deleted_at IS NULL
ORDER BY d.created_at DESC;

COMMENT ON VIEW v_partial_disbursements IS 'Partial disbursements with pending amounts (BR-IC-COM-005)';

-- ----------------------------------------------------------------------------
-- View: v_payment_retries
-- Description: Failed payments awaiting retry (BR-IC-SUSPENSE-002)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_payment_retries AS
SELECT
    d.disbursement_id,
    d.agent_id,
    ap.first_name,
    ap.last_name,
    d.payment_mode,
    d.disbursement_amount,
    d.disbursement_status,
    d.retry_count,
    d.max_retries,
    d.retry_status,
    d.last_retry_date,
    d.next_retry_date,
    d.pfms_transaction_id,
    CASE
        WHEN d.retry_count >= d.max_retries THEN 'MAX_RETRIES_REACHED'
        WHEN d.next_retry_date <= CURRENT_DATE THEN 'READY_FOR_RETRY'
        ELSE 'PENDING'
    END as retry_action,
    EXTRACT(DAY FROM (CURRENT_DATE - d.created_at)) as days_since_creation
FROM disbursements d
JOIN agent_profiles ap ON d.agent_id = ap.agent_id
WHERE d.retry_count > 0
    AND d.disbursement_status IN ('FAILED', 'PROCESSING')
    AND d.deleted_at IS NULL
    AND ap.deleted_at IS NULL
ORDER BY d.next_retry_date ASC NULLS LAST, d.retry_count DESC;

COMMENT ON VIEW v_payment_retries IS 'Failed payments awaiting retry (BR-IC-SUSPENSE-002)';

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Function: update_updated_at_column()
-- Description: Automatically update updated_at timestamp and version
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    IF OLD IS NOT NULL THEN
        NEW.version = OLD.version + 1;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_updated_at_column() IS 'Auto-update updated_at and version on row modification';

-- ----------------------------------------------------------------------------
-- Function: calculate_penalty_interest()
-- Description: Calculate penalty interest for SLA breach (BR-IC-COM-011)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION calculate_penalty_interest(
    p_disbursement_id VARCHAR
)
RETURNS DECIMAL(15,2) AS $$
DECLARE
    v_trial_approval_date DATE;
    v_payment_date DATE;
    v_days_overdue INT;
    v_disbursement_amount DECIMAL(15,2);
    v_penalty_rate DECIMAL(5,2) := 8.0; -- 8% annual rate
    v_penalty_amount DECIMAL(15,2);
BEGIN
    SELECT trial_approval_date, payment_date, disbursement_amount
    INTO v_trial_approval_date, v_payment_date, v_disbursement_amount
    FROM disbursements
    WHERE disbursement_id = p_disbursement_id;

    IF NOT FOUND THEN
        RETURN 0;
    END IF;

    -- Calculate days overdue (excluding weekends and holidays)
    v_days_overdue := EXTRACT(DAY FROM (
        COALESCE(v_payment_date, CURRENT_DATE) - v_trial_approval_date
    ))::INT - 10; -- 10 working days SLA

    IF v_days_overdue <= 0 THEN
        RETURN 0;
    END IF;

    -- Calculate penalty: (Amount * Rate * Days) / (365 * 100)
    v_penalty_amount := (
        v_disbursement_amount * v_penalty_rate * v_days_overdue
    ) / (365 * 100);

    RETURN ROUND(v_penalty_amount, 2);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION calculate_penalty_interest(VARCHAR) IS 'Calculate penalty for SLA breach (8% annual) (BR-IC-COM-011)';

-- ----------------------------------------------------------------------------
-- Function: validate_workflow_transition()
-- Description: Validate workflow state transitions
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION validate_workflow_transition(
    p_table_name TEXT,
    p_record_id BIGINT,
    p_new_state TEXT
)
RETURNS BOOLEAN AS $$
DECLARE
    v_current_state TEXT;
    v_allowed_transitions TEXT[];
BEGIN
    -- Get current state
    EXECUTE format(
        'SELECT workflow_state FROM %I WHERE %I = $1',
        p_table_name,
        case
            when p_table_name = 'commission_transactions' then 'commission_id'
            when p_table_name = 'trial_statements' then 'trial_statement_id'
            when p_table_name = 'final_statements' then 'final_statement_id'
            when p_table_name = 'disbursements' then 'disbursement_id::bigint'
            else 'id'
        end
    ) INTO v_current_state USING p_record_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Record not found in %', p_table_name;
    END IF;

    -- Define allowed transitions for each table
    IF p_table_name = 'commission_transactions' THEN
        v_allowed_transitions := ARRAY[
            'CALCULATED:TRIAL_PENDING',
            'TRIAL_PENDING:TRIAL_APPROVED',
            'TRIAL_APPROVED:FINALIZED',
            'FINALIZED:READY_FOR_DISBURSEMENT',
            'READY_FOR_DISBURSEMENT:DISBURSED'
        ];
    ELSIF p_table_name = 'trial_statements' THEN
        v_allowed_transitions := ARRAY[
            'PENDING_APPROVAL:APPROVED',
            'PENDING_APPROVAL:REJECTED'
        ];
    ELSIF p_table_name = 'final_statements' THEN
        v_allowed_transitions := ARRAY[
            'FINALIZED:READY_FOR_DISBURSEMENT',
            'READY_FOR_DISBURSEMENT:DISBURSEMENT_QUEUED',
            'DISBURSEMENT_QUEUED:DISBURSED'
        ];
    ELSIF p_table_name = 'disbursements' THEN
        v_allowed_transitions := ARRAY[
            'PENDING:PROCESSING',
            'PROCESSING:COMPLETED',
            'PROCESSING:FAILED'
        ];
    ELSE
        RETURN TRUE;
    END IF;

    -- Check if transition is allowed
    IF v_current_state IS NULL THEN
        RETURN TRUE; -- No current state, allow any state
    END IF;

    IF p_new_state = v_current_state THEN
        RETURN TRUE; -- Same state, allow
    END IF;

    IF (v_current_state || ':' || p_new_time) = ANY(v_allowed_transitions) THEN
        RETURN TRUE;
    END IF;

    RAISE EXCEPTION 'Invalid workflow transition from % to % in table %',
                   v_current_state, p_new_state, p_table_name;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION validate_workflow_transition(TEXT, BIGINT, TEXT) IS 'Validate workflow state transitions';

-- ----------------------------------------------------------------------------
-- Function: update_agent_commission_status()
-- Description: Update agent commission status based on clawback/suspense
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_agent_commission_status()
RETURNS TRIGGER AS $$
BEGIN
    -- After insert or update on clawbacks
    IF TG_TABLE_NAME = 'commission_clawbacks' THEN
        IF TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND NEW.recovery_status IN ('PENDING', 'PARTIALLY_RECOVERED')) THEN
            UPDATE agent_profiles
            SET clawback_pending_amount = (
                SELECT COALESCE(SUM(clawback_amount - recovered_amount), 0)
                FROM commission_clawbacks
                WHERE agent_code = NEW.agent_code
                    AND recovery_status IN ('PENDING', 'PARTIALLY_RECOVERED')
                    AND deleted_at IS NULL
            ),
            commission_status = CASE
                WHEN EXISTS (
                    SELECT 1 FROM commission_clawbacks
                    WHERE agent_code = NEW.agent_code
                    AND recovery_status IN ('PENDING', 'PARTIALLY_RECOVERED')
                    AND deleted_at IS NULL
                ) THEN 'SUSPENDED_PENDING_CLAWBACK'
                ELSE 'ACTIVE'
            END
            WHERE agent_id = NEW.agent_code;
        END IF;
    END IF;

    -- After insert or update on suspense
    IF TG_TABLE_NAME = 'commission_suspense' THEN
        IF TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND NEW.status = 'SUSPENDED') THEN
            UPDATE agent_profiles
            SET suspense_amount = (
                SELECT COALESCE(SUM(commission_amount), 0)
                FROM commission_suspense
                WHERE agent_code = NEW.agent_code
                    AND status = 'SUSPENDED'
                    AND deleted_at IS NULL
            ),
            commission_status = CASE
                WHEN EXISTS (
                    SELECT 1 FROM commission_suspense
                    WHERE agent_code = NEW.agent_code
                    AND status = 'SUSPENDED'
                    AND deleted_at IS NULL
                ) THEN 'SUSPENDED_PENDING_INVESTIGATION'
                WHEN EXISTS (
                    SELECT 1 FROM commission_clawbacks
                    WHERE agent_code = NEW.agent_code
                    AND recovery_status IN ('PENDING', 'PARTIALLY_RECOVERED')
                    AND deleted_at IS NULL
                ) THEN 'SUSPENDED_PENDING_CLAWBACK'
                ELSE 'ACTIVE'
            END
            WHERE agent_id = NEW.agent_code;
        ELSIF TG_OP = 'UPDATE' AND NEW.status IN ('RELEASED', 'FORFEITED') THEN
            -- Recalculate suspense amount and potentially update commission status
            UPDATE agent_profiles
            SET suspense_amount = (
                SELECT COALESCE(SUM(commission_amount), 0)
                FROM commission_suspense
                WHERE agent_code = NEW.agent_code
                    AND status = 'SUSPENDED'
                    AND deleted_at IS NULL
            ),
            commission_status = CASE
                WHEN EXISTS (
                    SELECT 1 FROM commission_suspense
                    WHERE agent_code = NEW.agent_code
                    AND status = 'SUSPENDED'
                    AND deleted_at IS NULL
                ) THEN 'SUSPENDED_PENDING_INVESTIGATION'
                WHEN EXISTS (
                    SELECT 1 FROM commission_clawbacks
                    WHERE agent_code = NEW.agent_code
                    AND recovery_status IN ('PENDING', 'PARTIALLY_RECOVERED')
                    AND deleted_at IS NULL
                ) THEN 'SUSPENDED_PENDING_CLAWBACK'
                ELSE 'ACTIVE'
            END
            WHERE agent_id = NEW.agent_code;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_agent_commission_status() IS 'Update agent commission status based on clawback/suspense (BR-IC-SUSPENSE-001)';

-- ----------------------------------------------------------------------------
-- Function: check_license_expiry()
-- Description: Check for license expiry and deactivate agents (BR-IC-LIC-003)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION check_license_expiry()
RETURNS TABLE(
    agent_id VARCHAR,
    license_id BIGINT,
    license_number VARCHAR,
    renewal_date DATE,
    days_expired INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        al.agent_id,
        al.license_id,
        al.license_number,
        al.renewal_date,
        EXTRACT(DAY FROM (CURRENT_DATE - al.renewal_date))::INT as days_expired
    FROM agent_licenses al
    JOIN agent_profiles ap ON al.agent_id = ap.agent_id
    WHERE CURRENT_DATE > al.renewal_date
        AND al.license_status != 'RENEWED'
        AND al.is_active = TRUE
        AND al.deleted_at IS NULL
        AND ap.deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION check_license_expiry() IS 'Check for expired licenses (BR-IC-LIC-003)';

-- ----------------------------------------------------------------------------
-- Function: calculate_clawback_percentage()
-- Description: Calculate clawback percentage based on months active (BR-IC-CLAWBACK-001)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION calculate_clawback_percentage(p_months_active INT)
RETURNS DECIMAL(5,2) AS $$
BEGIN
    IF p_months_active < 6 THEN
        RETURN 100.0;
    ELSIF p_months_active < 12 THEN
        RETURN 75.0;
    ELSIF p_months_active < 18 THEN
        RETURN 50.0;
    ELSIF p_months_active < 24 THEN
        RETURN 25.0;
    ELSE
        RETURN 0.0;
    END IF;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION calculate_clawback_percentage(INT) IS 'Calculate clawback % based on months active (BR-IC-CLAWBACK-001)';

-- ----------------------------------------------------------------------------
-- Function: schedule_payment_retry()
-- Description: Schedule next retry attempt for failed payments (BR-IC-SUSPENSE-002)
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION schedule_payment_retry(
    p_disbursement_id VARCHAR,
    p_failure_reason VARCHAR
)
RETURNS BOOLEAN AS $$
DECLARE
    v_retry_count INT;
    v_max_retries INT;
    v_next_retry_date DATE;
    v_retry_status VARCHAR(50);
BEGIN
    -- Get current retry info
    SELECT retry_count, max_retries
    INTO v_retry_count, v_max_retries
    FROM disbursements
    WHERE disbursement_id = p_disbursement_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Disbursement % not found', p_disbursement_id;
    END IF;

    -- Check if max retries reached
    IF v_retry_count >= v_max_retries THEN
        -- Mark as final failure
        UPDATE disbursements
        SET retry_status = 'MAX_RETRIES_REACHED',
            disbursement_status = 'FAILED',
            next_retry_date = NULL
        WHERE disbursement_id = p_disbursement_id;

        -- Notify for manual intervention
        RAISE NOTICE 'Max retries reached for disbursement %, manual intervention required', p_disbursement_id;
        RETURN FALSE;
    END IF;

    -- Calculate next retry date (exponential backoff: 1 day, 3 days, 7 days)
    v_retry_count := v_retry_count + 1;

    CASE v_retry_count
        WHEN 1 THEN v_next_retry_date := CURRENT_DATE + INTERVAL '1 day';
        WHEN 2 THEN v_next_retry_date := CURRENT_DATE + INTERVAL '3 days';
        WHEN 3 THEN v_next_retry_date := CURRENT_DATE + INTERVAL '7 days';
        ELSE v_next_retry_date := CURRENT_DATE + INTERVAL '7 days';
    END CASE;

    v_retry_status := 'RETRY_' || v_retry_count;

    -- Update disbursement
    UPDATE disbursements
    SET retry_count = v_retry_count,
        retry_status = v_retry_status,
        last_retry_date = CURRENT_DATE,
        next_retry_date = v_next_retry_date,
        disbursement_status = 'PROCESSING'
    WHERE disbursement_id = p_disbursement_id;

    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION schedule_payment_retry(VARCHAR, VARCHAR) IS 'Schedule next payment retry with exponential backoff (BR-IC-SUSPENSE-002)';

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Triggers for update_updated_at_column()
-- ----------------------------------------------------------------------------
CREATE TRIGGER trg_agent_profiles_updated_at
    BEFORE UPDATE ON agent_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_agent_addresses_updated_at
    BEFORE UPDATE ON agent_addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_agent_contacts_updated_at
    BEFORE UPDATE ON agent_contacts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_agent_emails_updated_at
    BEFORE UPDATE ON agent_emails
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_agent_bank_accounts_updated_at
    BEFORE UPDATE ON agent_bank_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_agent_licenses_updated_at
    BEFORE UPDATE ON agent_licenses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_commission_rates_updated_at
    BEFORE UPDATE ON commission_rates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_commission_transactions_updated_at
    BEFORE UPDATE ON commission_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_trial_statements_updated_at
    BEFORE UPDATE ON trial_statements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_final_statements_updated_at
    BEFORE UPDATE ON final_statements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_disbursements_updated_at
    BEFORE UPDATE ON disbursements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_commission_clawbacks_updated_at
    BEFORE UPDATE ON commission_clawbacks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_commission_suspense_updated_at
    BEFORE UPDATE ON commission_suspense
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ----------------------------------------------------------------------------
-- Trigger: update_agent_commission_status for clawbacks
-- ----------------------------------------------------------------------------
CREATE TRIGGER trg_clawback_agent_status
    AFTER INSERT OR UPDATE ON commission_clawbacks
    FOR EACH ROW EXECUTE FUNCTION update_agent_commission_status();

-- ----------------------------------------------------------------------------
-- Trigger: update_agent_commission_status for suspense
-- ----------------------------------------------------------------------------
CREATE TRIGGER trg_suspense_agent_status
    AFTER INSERT OR UPDATE ON commission_suspense
    FOR EACH ROW EXECUTE FUNCTION update_agent_commission_status();

-- ----------------------------------------------------------------------------
-- Function: update_search_vector()
-- Description: Update full-text search vector
-- ----------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    -- For agent_profiles
    IF TG_TABLE_NAME = 'agent_profiles' THEN
        NEW.search_vector := to_tsvector('english',
            COALESCE(NEW.first_name, '') || ' ' ||
            COALESCE(NEW.middle_name, '') || ' ' ||
            COALESCE(NEW.last_name, '') || ' ' ||
            COALESCE(NEW.pan, '') || ' ' ||
            COALESCE(NEW.agent_id, '')
        );
    END IF;

    -- For commission_transactions
    IF TG_TABLE_NAME = 'commission_transactions' THEN
        NEW.search_vector := to_tsvector('english',
            COALESCE(NEW.agent_id, '') || ' ' ||
            COALESCE(NEW.policy_number, '') || ' ' ||
            COALESCE(NEW.batch_id, '')
        );
    END IF;

    -- For trial_statements
    IF TG_TABLE_NAME = 'trial_statements' THEN
        NEW.search_vector := to_tsvector('english',
            COALESCE(NEW.statement_number, '') || ' ' ||
            COALESCE(NEW.agent_id, '') || ' ' ||
            COALESCE(NEW.batch_id, '')
        );
    END IF;

    -- For final_statements
    IF TG_TABLE_NAME = 'final_statements' THEN
        NEW.search_vector := to_tsvector('english',
            COALESCE(NEW.statement_number, '') || ' ' ||
            COALESCE(NEW.agent_id, '')
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_search_vector() IS 'Update full-text search vector';

-- ----------------------------------------------------------------------------
-- Triggers for update_search_vector()
-- ----------------------------------------------------------------------------
CREATE TRIGGER trg_agent_profiles_search_vector
    BEFORE INSERT OR UPDATE ON agent_profiles
    FOR EACH ROW EXECUTE FUNCTION update_search_vector();

CREATE TRIGGER trg_commission_transactions_search_vector
    BEFORE INSERT OR UPDATE ON commission_transactions
    FOR EACH ROW EXECUTE FUNCTION update_search_vector();

CREATE TRIGGER trg_trial_statements_search_vector
    BEFORE INSERT OR UPDATE ON trial_statements
    FOR EACH ROW EXECUTE FUNCTION update_search_vector();

CREATE TRIGGER trg_final_statements_search_vector
    BEFORE INSERT OR UPDATE ON final_statements
    FOR EACH ROW EXECUTE FUNCTION update_search_vector();

-- ============================================================================
-- SAMPLE DATA (OPTIONAL - FOR TESTING)
-- ============================================================================

-- Insert sample commission rates (BR-IC-COM-006)
-- Uncomment to insert sample data

/*
INSERT INTO commission_rates (
    rate_percentage, policy_duration_months, product_type,
    product_plan_code, agent_type, policy_term_years, effective_from
) VALUES
    (5.0, 12, 'PLI', 'ENDOWMENT', 'DIRECT', 15, CURRENT_DATE),
    (4.5, 12, 'RPLI', 'WHOLE_LIFE', 'FIELD_OFFICER', 20, CURRENT_DATE),
    (6.0, 12, 'PLI', 'MONEY_BACK', 'ADVISOR', 10, CURRENT_DATE),
    (3.5, 24, 'PLI', 'ENDOWMENT', 'DIRECT', 15, CURRENT_DATE),
    (2.5, 24, 'RPLI', 'WHOLE_LIFE', 'FIELD_OFFICER', 20, CURRENT_DATE);
*/

-- ============================================================================
-- END OF SCHEMA
-- ============================================================================
