

# Incentive, Commission and Producer Management - Detailed Requirements Analysis

## Document Control

| Attribute | Details |
|-----------|---------|
| **Module** | Incentive, Commission and Producer Management (IC) |
| **Phase** | Phase 4 - Agent Management |
| **Team** | Team 1 - Agent Management |
| **Analysis Date** | January 24, 2026 |
| **Source Documents** | 2 SRS Documents |
| **Complexity** | High |
| **Technology Stack** | Golang, Temporal.io, PostgreSQL, Kafka, React |

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Business Rules](#2-business-rules)
3. [Functional Requirements](#3-functional-requirements)
4. [Validation Rules](#4-validation-rules)
5. [Error Codes](#5-error-codes)
6. [Workflows](#6-workflows)
7. [Data Entities](#7-data-entities)
8. [Integration Points](#8-integration-points)
9. [Temporal Workflows](#9-temporal-workflows)
10. [Traceability Matrix](#10-traceability-matrix)
11. [Commission Rate Structure](#11-commission-rate-structure)
12. [Taxation and Compliance](#12-taxation-and-compliance)
13. [Report Specifications](#13-report-specifications)
14. [Exception Handling](#14-exception-handling)
15. [User Interface Requirements](#15-user-interface-requirements)
16. [Security and Access Control Details](#16-security-and-access-control-details)
17. [Performance SLAs](#17-performance-slas)
18. [Sample Calculations and Examples](#18-sample-calculations-and-examples)
19. [Glossary and Definitions](#19-glossary-and-definitions)
20. [Commission Clawback and Suspense Details](#20-commission-clawback-and-suspense-details)

---

## 1. Executive Summary

### 1.1 Purpose
This document provides comprehensive business requirements analysis for the **Incentive, Commission and Producer Management** module of the Postal Life Insurance (PLI) and Rural Postal Life Insurance (RPLI) system. This module is responsible for agent onboarding, commission processing, license management, and producer lifecycle operations.

### 1.2 Scope
The analysis covers the Incentive, Commission and Producer Management functionality including:
1. **Agent Onboarding** - Advisor, Advisor Coordinator, Departmental Employee, Field Officer
2. **Agent Profile Management** - Search, update, status management, termination
3. **License Management** - Issuance, renewal tracking, reminders, expiry handling
4. **Commission Rate Configuration** - Rate table management by product/agent type
5. **Commission Calculation & Processing** - Monthly batch calculation, annualised premium
6. **Trial & Final Statement Generation** - Review, approval, partial disbursement
7. **Commission Disbursement** - Cheque and EFT modes with PFMS integration
8. **Agent Termination** - Status changes, commission handling, deactivation

### 1.3 Key Statistics

| Metric | Count |
|--------|-------|
| **Business Rules** | 37 (30 + 7 clawback/suspense) |
| **Functional Requirements** | 44 (32 + 12 clawback/suspense) |
| **Validation Rules** | 45+ |
| **Workflows** | 10 (8 + 2 clawback/suspense) |
| **Temporal Workflows** | 6 (4 + 2 clawback/suspense with complete Go code) |
| **Commission Types** | 3 (First Year, Renewal, Bonus) |
| **Agent Types** | 4 (Advisor, Coordinator, Dept Employee, Field Officer) |
| **Payment Modes** | 2 (Cheque, EFT) |
| **Data Entities** | 14 (12 + 2 clawback/suspense) |
| **Integration Points** | 4 |
| **Error Codes** | 19 (16 + 3 clawback/suspense) |

### 1.4 Critical Dependencies

| Dependency | Purpose | Impact |
|------------|---------|--------|
| **Policy Services** | Policy issuance triggers commissions | Commission calculation cannot proceed without policy data |
| **Accounting/Finance** | Commission disbursement accounting & PFMS integration | Payment processing, ledger entries |
| **HRMS System** | Departmental employee data auto-population | Dept employee onboarding |
| **Banking/PFMS** | EFT payment processing | Electronic fund transfers |

### 1.5 SLA Requirements

| Process | SLA | Penalty/Escalation |
|---------|-----|-------------------|
| License Renewal Processing | 3 working days | Escalation to supervisor |
| Commission Batch Processing | 6 hours max | Escalation after 3 hours, critical after 5 hours |
| License Renewal Reminders | T-30, T-15, T-7, T-0 days | Automated notifications |

### 1.6 Key Business Rules Summary

#### Agent Hierarchy
- **Advisors MUST be linked to existing Advisor Coordinator**
- **Advisor Coordinators MUST be assigned to Circle and Division**

#### License Management
- **First renewal**: After 1 year
- **Subsequent renewals**: Every 3 years
- **Auto-deactivation**: If renewal date elapses without renewal
- **Reminders**: At 30, 15, 7, and 0 days before expiry

#### Commission Processing
- **Monthly batch calculation**: Will be triggered manually
- **Trial statement MUST be approved before disbursement**
- **TDS deducted** as applicable based on PAN
- **Premium** = Annual=1, Semi-annual=2, Quarterly=4, Monthly=12
- **Partial disbursement allowed** with percentage
- **Disbursement SLA**: 10 working days from trial approval

### 1.7 Document Structure

This analysis document follows the same structure as the Phase4_Agent_Management_Analysis.md but contains **only IC-specific content** extracted from both source documents:

1. **Agent_SRS_Incentive-Commission-and-Producer-Management.md** - Primary IC SRS
2. **Phase4_Agent_Management_Analysis.md** - Extracted IC-related sections

Content NOT included in this document (present in Phase4 but NOT IC):
- Medical Appointment System requirements
- Agent Portal authentication beyond basic onboarding
- Agent Goal Setting functionality
- Agent Letters and Reports generation
- Actuary Valuation Report details (beyond annualised premium reference)

---

## 2. Commission Types

### 2.1 First Year Commission
- Paid on new policies procured by the advisor in the past 12 months
- Typically higher rate to incentivize new business
- Calculated on premium collected
- Triggered by policy issuance from Policy Services

### 2.2 Renewal Commission
- Paid on policy premiums for subsequent years after completion of first 12 months of policy procurement
- Lower rate compared to first year
- Calculated on premium collected
- Triggered by premium collection events

---

## 4. Payment Modes

### 4.1 Cheque Payment
- **Immediate disbursement completion**
- Physical cheque generated
- Payment advice sent to agent
- No external system dependency
- Lower processing complexity

### 4.2 EFT (Electronic Funds Transfer)

- **Requires bank account details** (account number, IFSC code) or **POSB account details** (account number)
- **Payment file sent to PFMS/Bank**
- **2-3 days processing time**
- **Confirmation callback updates status**
- Integration with PFMS system required
- Higher processing complexity but better for bulk payments

---

## 5. Quick Reference Tables

### 5.1 Commission Processing Timeline

| Event | Timeline | Action |
|-------|----------|--------|
| Month End | Last day of month | Policies eligible for commission |
| Batch Calculation | First working day (within 6 hours) | Calculate commissions |
| Trial Statement | Same day or next | Generate for review |
| Finance Review | Within 7 days | Approve or reject |
| Final Statement | After approval | Lock and finalize |
| Disbursement | Within 10 working days | Process payments |
| SLA Breach | After 10 days | Escalate + penalty interest |

---

## Document Navigation

**Next Section**: [Business Rules](#2-business-rules)

This section contains 30+ detailed business rules organized by category:
- Agent Hierarchy Rules (4 rules)
- License Management Rules (5 rules)
- Commission Processing Rules (12 rules)
- Agent Profile Management Rules (4 rules)
- Data Validation Rules (3 rules)

---

**End of Section 1: Executive Summary**
# Incentive, Commission and Producer Management - Business Rules

## 2. Business Rules


### 2.3 Commission Processing Rules

#### BR-IC-COM-001: Monthly Commission Calculation
- **ID**: BR-IC-COM-001
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: Commission calculated monthly via batch jobs
- **Rule**: `Run commission_calculation_batch ON first_working_day_of_month FOR previous_month_policies`
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 3 (Business Requirements), FS_IC_024, Lines 218-220
- **Traceability**: FR-IC-COM-001
- **Impact**: Regular monthly commission processing
- **Example**: Batch job runs on February 1, 2026 (first working day) to calculate January 2026 commissions.

#### BR-IC-COM-002: Trial Statement Before Disbursement
- **ID**: BR-IC-COM-002
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: Trial statement must be generated and approved before final disbursement
- **Rule**:
  ```
  disbursement_allowed = TRUE ONLY IF trial_statement_status = 'APPROVED'

  IF trial_statement_status = 'PENDING' THEN
    block_disbursement()
    notify_user('Trial statement must be approved before disbursement')
  END
  ```
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 3 (Business Requirements), FS_IC_025-FS_IC_028, Lines 222-241; Section 4.4.3 (Commission Processing Steps), Lines 747-769
- **Traceability**: FR-IC-COM-003, FR-IC-COM-005
- **Impact**: Verification before payment
- **Example**: Agent AGT001 has ₹50,000 commission. Trial statement generated. Finance reviews and approves on 2026-02-05. Only then can disbursement proceed.

#### BR-IC-COM-003: TDS Deduction Requirement
- **ID**: BR-IC-COM-003
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: TDS must be deducted from commission payments as applicable
- **Rule**:
  ```
  IF agent.pan_available AND TDS_applicable THEN
    TDS_amount = gross_commission * TDS_rate
  ELSE
    TDS_amount = 0
  END

  net_payable = gross_commission - TDS_amount
  ```
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4 (Agent Commission Management), Lines 676-800; Section 4.4.3.4 (Manual Trial Statement Generation), Line 735
- **Traceability**: FR-IC-COM-007
- **Impact**: Tax compliance
- **Example**: Commission = ₹50,000, TDS rate = 5%, TDS = ₹2,500, Net payable = ₹47,500.

#### BR-IC-COM-005: Partial Disbursement Option
- **ID**: BR-IC-COM-005
- **Category**: Commission Disbursement
- **Priority**: MEDIUM
- **Description**: System must allow full or partial commission disbursement
- **Rule**:
  ```
  disbursement_mode IN ['FULL', 'PARTIAL']

  IF disbursement_mode = 'PARTIAL' THEN
    disbursement_percentage = user_input_percentage (0-100)
    disbursement_amount = gross_commission × (disbursement_percentage / 100)
    pending_amount = gross_commission - disbursement_amount
  ELSE
    disbursement_amount = gross_commission
    pending_amount = 0
  END
  ```
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3 (Commission Processing Steps), Lines 747-769
- **Traceability**: FR-IC-COM-006
- **Impact**: Flexible payment options
- **Example**: ₹50,000 commission. Finance approves 60% partial disbursement = ₹30,000 paid, ₹20,000 pending.

#### BR-IC-COM-006: Commission Rate Table Structure
- **ID**: BR-IC-COM-006
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: System must maintain a commission rate table with fields for Rate (%), Policy Duration (Months), Product Type (PLI, RPLI), Product Plan Code, Agent Type, and Policy Term (Years)
- **Formula**: `commission_rate = lookup(rate_table, product_type, plan_code, agent_type, policy_term, duration)`
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Lines 208-212, FS_IC_022; Section 4.4.1 (Commission Rate Table View Page)
- **Traceability**: FR-IC-COM-008
- **Rationale**: Foundation for all commission calculations across product lines
- **Impact**: Without this structure, commission calculations cannot be performed accurately
- **Example**:
  | Rate % | Duration (Months) | Product Type | Plan Code | Agent Type | Policy Term (Years) |
  |--------|------------------|--------------|-----------|-----------|---------------------|
  | 5.0    | 12               | PLI          | ENDOWMENT | DIRECT    | 15                  |
  | 4.5    | 24               | RPLI         | WHOLE_LIFE| FIELD_OFFICER| 20              |

#### BR-IC-COM-007: Final Statement Generation Batch
- **ID**: BR-IC-COM-007
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: System must support automatic execution of Final Incentive Statement Generation batch job after trial statement approval
- **Rule**:
  ```
  AFTER trial_statement_status = 'APPROVED' DO
    Run final_statement_batch()
    Lock trial_data()
    Calculate final_amounts()
    Generate final_statements()
    Set status = 'READY_FOR_DISBURSEMENT'
  END
  ```
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Lines 243-245, FS_IC_029; Section 4.4.3.6 (Final Incentive Statement Generation Batch Job)
- **Traceability**: FR-IC-COM-009
- **Rationale**: Locks approved trial data and generates final payment statements
- **Impact**: Critical for commission disbursement workflow
- **Example**: Trial statement approved at 10:00 AM. Final statement batch runs at 11:00 AM generating final PDFs.

#### BR-IC-COM-008: Disbursement Mode Workflow
- **ID**: BR-IC-COM-008
- **Category**: Commission Disbursement
- **Priority**: CRITICAL
- **Description**: System must support cheque and EFT payment modes with PFMS/Bank integration for EFT processing
- **Rule**:
  ```
  IF payment_mode = 'CHEQUE' THEN
    generate_cheque_payment_advice()
    mark_disbursement_complete_immediately()
    log_payment_details()
  ELSE IF payment_mode = 'EFT' THEN
    generate_payment_file(format='PFMS/BANK')
    send_to_PFMS_or_Bank()
    set_status = 'DISBURSEMENT_QUEUED'
    wait_for_confirmation_callback()
    ON confirmation:
      update_status = 'DISBURSED'
      send_confirmation_to_agent()
  END
  ```
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Lines 251-258, FS_IC_032; Section 4.4.3.9 (Automatic Disbursement Option), Lines 833-844
- **Traceability**: FR-IC-COM-010
- **Rationale**: Different processing workflows for different payment methods
- **Impact**: Essential for automated commission disbursement
- **Example**:
  - **Cheque**: Generated today, marked as disbursed immediately
  - **EFT**: File sent to PFMS, status = "QUEUED", callback updates to "DISBURSED" in 2-3 days

#### BR-IC-COM-009: Commission History Search
- **ID**: BR-IC-COM-009
- **Category**: Commission Processing
- **Priority**: HIGH
- **Description**: System must provide commission history search by Agent ID, Policy Number, Date Range, Product Type, and Commission Type (First Year, Renewal, Bonus)
- **Rule**: `commission_history = search(agent_id, policy_number, date_range, product_type, commission_type)`
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Lines 214-216, FS_IC_023; Section 4.4.2 (Commission History Search Page), Lines 646-656
- **Traceability**: FR-IC-COM-011
- **Rationale**: Transparency and audit trail for agents and administrators
- **Impact**: Required for dispute resolution and reconciliation
- **Example**: Agent searches all "First Year" commissions for "January 2026" to verify received amounts.

#### BR-IC-COM-010: Export Commission Statements
- **ID**: BR-IC-COM-010
- **Category**: Commission Processing
- **Priority**: MEDIUM
- **Description**: System must provide export functionality for trial and final statements in Excel/PDF formats
- **Rule**:
  ```
  export_statements(format IN ['EXCEL', 'PDF'], statement_ids)

  IF format = 'PDF' THEN
    generate_individual_statement_pdfs()
    OR generate_agent_wise_consolidated_pdf()
  ELSE IF format = 'EXCEL' THEN
    generate_excel_workbook(with_multiple_sheets = TRUE)
  END
  ```
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Lines 707-711, 796-800 (View Trial Statement Page, Final Statements Page)
- **Traceability**: FR-IC-COM-012
- **Rationale**: Agents need downloadable records for personal accounting
- **Impact**: User experience and record-keeping
- **Example**: Finance exports all approved trial statements to Excel for bulk review before approval.

#### BR-IC-COM-012: Commission Batch Processing 6-Hour Timeout
- **ID**: BR-IC-COM-012
- **Category**: Commission Processing
- **Priority**: HIGH
- **Description**: Commission calculation batch job must complete within 6 hours with monitoring and retry logic
- **Rule**:
  ```
  batch_start_time = now()
  batch_timeout = 6 hours

  -- Monitor batch progress every 30 minutes
  EVERY 30 minutes:
    batch_progress = (processed_records / total_records) * 100
    elapsed_time = now() - batch_start_time

    IF elapsed_time > 3 hours AND batch_progress < 50% THEN
      send_alert(message='Commission batch running slow - only ' + batch_progress + '% complete', urgency='MEDIUM')
    END

    IF elapsed_time > 5 hours AND batch_progress < 80% THEN
      send_alert(message='Commission batch critical - ' + batch_progress + '% complete with 1 hour remaining', urgency='HIGH')
    END

    IF elapsed_time > batch_timeout THEN
      mark_batch_failed(reason='Timeout after 6 hours')
      send_alert(urgency='CRITICAL', message='Commission batch timeout - failed after 6 hours')
      notify_support_team()
      trigger_manual_intervention()
    END
  ```

  -- Retry logic for failed records
  FOR EACH failed_record:
    retry_count = 0
    max_retries = 3

    WHILE retry_count < max_retries:
      retry_record_processing()
      retry_count++

      IF success THEN break
    END

    IF retry_count = max_retries AND NOT success THEN
      log_error(record_id, error_details)
      add_to_manual_review_queue()
    END
  END
  ```
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Lines 220-245 (Commission Calculation section)
- **Traceability**: FR-IC-COM-001
- **Impact**: Batch reliability, timely commission processing
- **Example**: Batch starts at 9:00 AM. At 12:30 PM (3.5 hours), only 40% complete → Alert sent. At 2:00 PM (5 hours), 75% complete → Critical alert.

---

### 2.7 Commission Suspense Rules

#### BR-IC-SUSPENSE-001: Commission Suspense for Disputed Policies
- **ID**: BR-IC-SUSPENSE-001
- **Category**: Commission Suspense
- **Priority**: CRITICAL
- **Description**: Commission must be held in suspense for policies under dispute or investigation
- **Rule**:
  ```
  -- Trigger when policy marked as disputed
  WHEN policy_status_changed_to('UNDER_INVESTIGATION'):
    policy_number = policy.policy_number
    investigation_reason = policy.investigation_reason

    -- Check if commission already paid
    commission_paid = SUM(
      SELECT gross_commission
      FROM commission_transactions
      WHERE policy_number = policy_number
      AND payment_status = 'DISBURSED'
    )

    IF commission_paid > 0 THEN
      -- Commission already paid - create clawback suspense
      CREATE commission_suspense {
        policy_number: policy_number
        agent_code: policy.agent_code
        commission_amount: commission_paid
        suspense_reason: 'POLICY_UNDER_INVESTIGATION'
        investigation_type: investigation_reason
        suspense_date: current_date
        status: 'SUSPENDED'
        investigation_reference: policy.investigation_reference
        expected_resolution_date: current_date + 30 days
      }

      -- Suspend future payments to agent
      agent.commission_status = 'SUSPENDED_PENDING_INVESTIGATION'
      agent.suspense_amount += commission_paid

    ELSE
      -- Commission not yet paid - mark for hold
      UPDATE commission_transactions
      SET payment_status = 'HELD_IN_SUSPENSE',
          suspense_reason = 'POLICY_UNDER_INVESTIGATION',
          suspense_date = current_date,
          hold_until = policy.expected_investigation_closure_date
      WHERE policy_number = policy_number
    END

    -- Notification
    notify_agent(
      agent_code=policy.agent_code,
      subject='Commission Suspended - Investigation',
      message='Commission for policy ' + policy_number +
              ' (₹' + commission_paid + ') held in suspense due to: ' + investigation_reason + '. ' +
              'You will be notified when investigation is complete.',
      priority='HIGH'
    )

    notify_finance_team(
      message='Commission suspense: Policy ' + policy_number +
              ', Reason: ' + investigation_reason
    )
  END

  -- Release suspense when investigation closed
  WHEN policy_status_changed_to('INVESTIGATION_CLOSED'):
    investigation_outcome = policy.investigation_outcome
    suspense_entry = get_suspense_entry(policy.policy_number)

    IF investigation_outcome == 'POLICY_GENUINE' THEN
      -- Release commission
      UPDATE commission_transactions
      SET payment_status = 'APPROVED_FOR_PAYMENT',
          suspense_released_date = current_date
      WHERE policy_number = policy.policy_number

      UPDATE commission_suspense
      SET status = 'RELEASED',
          release_date = current_date,
          release_reason = 'Investigation cleared - policy genuine'
      WHERE policy_number = policy.policy_number

      agent.commission_status = 'ACTIVE'
      agent.suspense_amount -= suspense_entry.commission_amount

      notify_agent(
        message='Commission suspense released for policy ' + policy.policy_number +
                '. Amount ₹' + suspense_entry.commission_amount + ' will be paid in next cycle.'
      )

    ELSE IF investigation_outcome == 'POLICY_FRAUDULENT' THEN
      -- Forfeit commission
      UPDATE commission_transactions
      SET payment_status = 'FORFEITED',
          forfeiture_date = current_date
      WHERE policy_number = policy.policy_number

      UPDATE commission_suspense
      SET status = 'FORFEITED',
          forfeiture_date = current_date,
          forfeiture_reason = 'Policy determined fraudulent'
      WHERE policy_number = policy.policy_number

      agent.suspense_amount -= suspense_entry.commission_amount

      -- Flag agent for review
      flag_agent_for_review(
        agent_code=policy.agent_code,
        reason='Fraudulent policy detected',
        severity='CRITICAL'
      )

      notify_agent(
        message='Commission for policy ' + policy.policy_number +
                ' has been forfeited due to fraud determination. Compliance review initiated.',
        priority='CRITICAL'
      )

      escalate_to_compliance_team(policy.agent_code)
    END
  END
  ```
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 320-360 (Suspense Account Management section, Suspense creation and investigation workflow)
- **Traceability**: FR-IC-COM-013, FR-IC-COM-014, FR-IC-COM-015
- **Impact**: Prevents payment of fraudulent commissions, protects company from financial loss, ₹5-10 Crore annual risk
- **Example**:
  - Policy under investigation for suspected fraud
  - Commission paid: ₹25,000
  - Suspense created on January 10, 2026
  - Investigation duration: 30 days
  - If cleared: Commission paid in February cycle
  - If fraud confirmed: Commission forfeited, agent flagged

#### BR-IC-SUSPENSE-002: Commission Payment Failure Suspense with Retry Logic
- **ID**: BR-IC-SUSPENSE-002
- **Category**: Commission Disbursement
- **Priority**: HIGH
- **Description**: Failed commission payments moved to suspense with 3 retry attempts before manual intervention
- **Rule**:
  ```
  WHEN commission_payment_failed(transaction_id, failure_reason):
    commission_txn = get_commission_transaction(transaction_id)

    commission_txn.payment_failure_count += 1
    commission_txn.last_failure_reason = failure_reason
    commission_txn.last_failure_date = current_date
    commission_txn.last_failure_time = current_time

    -- Retry logic with exponential backoff
    IF commission_txn.payment_failure_count <= 3 THEN
      -- Calculate retry delay: 2, 4, 8 hours
      retry_delay_hours = 2^(commission_txn.payment_failure_count)
      retry_time = current_time + retry_delay_hours hours

      schedule_retry(
        transaction_id=transaction_id,
        retry_attempt=commission_txn.payment_failure_count,
        retry_after=retry_time
      )

      commission_txn.payment_status = 'RETRY_SCHEDULED'
      commission_txn.next_retry_time = retry_time

      -- Log retry attempt
      CREATE payment_retry_log {
        transaction_id: transaction_id
        retry_attempt: commission_txn.payment_failure_count
        failure_reason: failure_reason
        retry_scheduled_time: retry_time
        status: 'SCHEDULED'
      }

      notify_agent(
        agent_code=commission_txn.agent_code,
        message='Commission payment temporarily failed: ' + failure_reason +
                '. Automatic retry scheduled in ' + retry_delay_hours + ' hours.',
        priority='MEDIUM'
      )

    ELSE
      -- Move to suspense after 3 failures
      CREATE commission_payment_suspense {
        transaction_id: transaction_id
        agent_code: commission_txn.agent_code
        amount: commission_txn.net_commission
        failure_reason: failure_reason
        failure_count: commission_txn.payment_failure_count
        suspense_date: current_date
        status: 'MANUAL_REVIEW_REQUIRED'
        bank_account_number: commission_txn.bank_account_number
        ifsc_code: commission_txn.ifsc_code
      }

      commission_txn.payment_status = 'FAILED_MOVED_TO_SUSPENSE'

      -- Create finance task
      CREATE finance_task {
        task_type: 'COMMISSION_PAYMENT_FAILURE_RESOLUTION'
        transaction_id: transaction_id
        agent_code: commission_txn.agent_code
        amount: commission_txn.net_commission
        failure_reason: failure_reason
        priority: 'HIGH'
        assigned_to: 'FINANCE_TEAM'
        due_date: current_date + 2 working days
        description: 'Resolve commission payment failure after 3 retry attempts. ' +
                     'Verify bank details and payment mode. Amount: ₹' + commission_txn.net_commission
      }

      notify_agent(
        agent_code=commission_txn.agent_code,
        subject='Commission Payment Failed - Manual Review',
        message='Commission payment of ₹' + commission_txn.net_commission +
                ' failed after 3 attempts. Reason: ' + failure_reason + '. ' +
                'Our finance team will contact you within 2 working days to resolve this.',
        priority='HIGH'
      )

      notify_finance_team(
        subject='Commission Payment Suspense - Manual Intervention Required',
        message='Transaction ID: ' + transaction_id +
                ', Agent: ' + commission_txn.agent_code +
                ', Amount: ₹' + commission_txn.net_commission +
                ', Failures: 3, Reason: ' + failure_reason,
        urgency='HIGH',
        assign_to='COMMISSION_RESOLUTION_TEAM'
      )
    END
  END
  ```
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 255-280 (Commission Disbursement section, Payment failure handling and retry logic)
- **Traceability**: FR-IC-COM-016, FR-IC-COM-017
- **Impact**: Ensures no commission amounts lost in failed transactions, maintains agent trust
- **Example**:
  - Payment attempt 1: Failed at 10:00 AM (invalid account), retry at 12:00 PM (2 hours)
  - Payment attempt 2: Failed at 12:00 PM (bank timeout), retry at 4:00 PM (4 hours)
  - Payment attempt 3: Failed at 4:00 PM (account frozen), retry at 12:00 AM next day (8 hours)
  - Payment attempt 4: Failed at 12:00 AM - moved to suspense, finance team notified
  - Resolution: Agent updates bank details, payment successful

#### BR-IC-SUSPENSE-004: Overpayment Recovery Suspense
- **ID**: BR-IC-SUSPENSE-004
- **Category**: Commission Adjustment
- **Priority**: MEDIUM
- **Description**: Commission overpayments identified and recovered through suspense mechanism
- **Rule**:
  ```
  -- Daily/monthly audit batch to detect overpayments
  FOR each agent IN agents_with_recent_payments:
    -- Detect overpayment scenarios
    overpayments = detect_overpayments(agent.agent_code)

    FOR each overpayment IN overpayments:
      CREATE commission_overpayment_suspense {
        agent_code: agent.agent_code
        overpayment_type: overpayment.type
        original_transaction_id: overpayment.transaction_id
        calculated_amount: overpayment.correct_amount
        paid_amount: overpayment.paid_amount
        overpaid_amount: overpayment.paid_amount - overpayment.correct_amount
        overpayment_reason: overpayment.reason
        detection_date: current_date
        recovery_status: 'PENDING_RECOVERY'
        recovery_method: determine_recovery_method(overpayment.overpaid_amount)
      }

      agent.overpayment_pending_recovery += overpayment.overpaid_amount

      notify_agent(
        agent_code=agent.agent_code,
        subject='Commission Overpayment Detected',
        message='An overpayment of ₹' + overpayment.overpaid_amount +
                ' has been identified for ' + overpayment.reason + '. ' +
                'This amount will be adjusted in your next commission payment. ' +
                'Reference: ' + overpayment.transaction_id,
        priority='MEDIUM'
      )
    END
  END

  -- Overpayment detection function
  FUNCTION detect_overpayments(agent_code):
    overpayments = []

    -- Type 1: Duplicate payments
    duplicate_payments = SELECT transaction_id, policy_number, amount, COUNT(*)
      FROM commission_transactions
      WHERE agent_code = agent_code
      AND payment_status = 'DISBURSED'
      AND payment_date >= current_date - 90 days
      GROUP BY policy_number, amount
      HAVING COUNT(*) > 1

    FOR each dup IN duplicate_payments:
      overpayments.add({
        type: 'DUPLICATE_PAYMENT',
        transaction_id: dup.transaction_id,
        correct_amount: dup.amount,
        paid_amount: dup.amount * 2,
        reason: 'Duplicate payment for policy ' + dup.policy_number
      })
    END

    -- Type 2: Incorrect rate application
    incorrect_rates = SELECT *
      FROM commission_transactions
      WHERE agent_code = agent_code
      AND payment_date >= current_date - 90 days
      AND commission_rate != expected_rate_for_policy_type

    FOR each txn IN incorrect_rates:
      correct_commission = txn.annualized_premium * correct_rate / 100
      IF txn.gross_commission > correct_commission THEN
        overpayments.add({
          type: 'INCORRECT_RATE',
          transaction_id: txn.transaction_id,
          correct_amount: correct_commission,
          paid_amount: txn.gross_commission,
          reason: 'Wrong commission rate applied'
        })
      END
    END

    -- Type 3: Policy cancelled but commission not reversed
    cancelled_policies = SELECT ct.*
      FROM commission_transactions ct
      JOIN policies p ON ct.policy_number = p.policy_number
      WHERE ct.agent_code = agent_code
      AND p.status IN ('CANCELLED', 'LAPSED')
      AND p.cancellation_within_free_look = TRUE
      AND ct.reversal_processed = FALSE

    FOR each txn IN cancelled_policies:
      overpayments.add({
        type: 'CANCELLED_POLICY',
        transaction_id: txn.transaction_id,
        correct_amount: 0,
        paid_amount: txn.gross_commission,
        reason: 'Policy cancelled within free-look period'
      })
    END

    RETURN overpayments
  END FUNCTION
  ```
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 320-360 (Commission Recovery section, Overpayment recovery and adjustment rules)
- **Traceability**: FR-IC-COM-016
- **Impact**: Recovers erroneous payments, maintains financial accuracy
- **Example**:
  - Overpayment detected: ₹10,000 (duplicate payment)
  - Next month commission: ₹30,000
  - Max recovery (50%): ₹15,000
  - Actual recovery: ₹10,000 (full overpayment)
  - Net payable: ₹20,000

#### BR-IC-SUSPENSE-005: Commission Hold for Investigation
- **ID**: BR-IC-SUSPENSE-005
- **Category**: Commission Suspense
- **Priority**: HIGH
- **Description**: All commission payments suspended when agent under investigation for compliance violations
- **Rule**:
  ```
  WHEN agent_investigation_initiated(agent_code, investigation_type, investigation_reason):
    agent = get_agent(agent_code)

    -- Suspend all pending and future commissions
    UPDATE commission_transactions
    SET payment_status = 'HELD_INVESTIGATION',
        hold_reason = 'AGENT_UNDER_INVESTIGATION',
        hold_date = current_date,
        investigation_reference = investigation.reference_number
    WHERE agent_code = agent_code
    AND payment_status IN ('CALCULATED', 'APPROVED', 'RETRY_SCHEDULED')

    -- Calculate total amount on hold
    total_held_amount = SUM(
      SELECT net_commission FROM commission_transactions
      WHERE agent_code = agent_code AND payment_status = 'HELD_INVESTIGATION'
    )

    CREATE agent_investigation_suspense {
      agent_code: agent_code
      investigation_type: investigation_type
      investigation_reason: investigation_reason
      investigation_reference: investigation.reference_number
      suspense_date: current_date
      total_amount_held: total_held_amount
      status: 'INVESTIGATION_IN_PROGRESS'
      expected_resolution_date: current_date + investigation.estimated_duration
    }

    agent.commission_status = 'SUSPENDED_UNDER_INVESTIGATION'
    agent.investigation_active = TRUE

    notify_agent(
      agent_code=agent_code,
      subject='Commission Payments Suspended - Investigation',
      message='All commission payments have been temporarily suspended due to ' +
              investigation_reason + '. Investigation reference: ' +
              investigation.reference_number + '. You will be contacted by the compliance team.',
      priority='CRITICAL'
    )

    notify_compliance_team(
      message='Agent ' + agent_code + ' commissions suspended. ' +
              'Amount on hold: ₹' + total_held_amount
    )
  END

  -- Investigation closure and commission release/forfeiture
  WHEN investigation_closed(agent_code, investigation_outcome, outcome_details):
    investigation_suspense = get_investigation_suspense(agent_code)
    agent = get_agent(agent_code)

    IF investigation_outcome == 'NO_VIOLATION_FOUND' THEN
      -- Release all held commissions
      UPDATE commission_transactions
      SET payment_status = 'APPROVED',
          hold_released_date = current_date
      WHERE agent_code = agent_code AND payment_status = 'HELD_INVESTIGATION'

      investigation_suspense.status = 'CLOSED_NO_VIOLATION'
      investigation_suspense.closure_date = current_date

      agent.commission_status = 'ACTIVE'
      agent.investigation_active = FALSE

      notify_agent(
        message='Investigation closed. No violations found. ' +
                'All held commissions (₹' + investigation_suspense.total_amount_held +
                ') will be processed in next payment cycle.',
        priority='HIGH'
      )

    ELSE IF investigation_outcome == 'MINOR_VIOLATION' THEN
      -- Release commissions with penalty deduction
      penalty_percentage = outcome_details.penalty_percentage
      penalty_amount = investigation_suspense.total_amount_held * (penalty_percentage / 100)

      UPDATE commission_transactions
      SET payment_status = 'APPROVED_WITH_PENALTY',
          penalty_deducted = (net_commission * penalty_percentage / 100),
          net_commission = net_commission * (1 - penalty_percentage / 100)
      WHERE agent_code = agent_code AND payment_status = 'HELD_INVESTIGATION'

      investigation_suspense.status = 'CLOSED_PENALTY_APPLIED'
      investigation_suspense.penalty_amount = penalty_amount

      agent.commission_status = 'ACTIVE_WITH_WARNING'

      notify_agent(
        message='Investigation closed. Minor violation found. ' +
                'Penalty of ' + penalty_percentage + '% (₹' + penalty_amount + ') applied. ' +
                'Remaining commissions will be paid. Warning issued.'
      )

    ELSE IF investigation_outcome == 'MAJOR_VIOLATION' THEN
      -- Forfeit all held commissions
      UPDATE commission_transactions
      SET payment_status = 'FORFEITED',
          forfeiture_reason = outcome_details.violation_details,
          forfeiture_date = current_date
      WHERE agent_code = agent_code AND payment_status = 'HELD_INVESTIGATION'

      investigation_suspense.status = 'CLOSED_COMMISSIONS_FORFEITED'
      investigation_suspense.forfeiture_amount = investigation_suspense.total_amount_held

      agent.commission_status = 'SUSPENDED'
      agent.status = 'SUSPENDED'

      notify_agent(
        message='Investigation closed. Major violation found. ' +
                'All held commissions (₹' + investigation_suspense.total_amount_held + ') forfeited. ' +
                'Account suspended. Disciplinary action pending.',
        priority='CRITICAL'
      )

      -- Initiate termination workflow if severe
      IF outcome_details.severity == 'SEVERE' THEN
        initiate_termination_workflow(agent_code, 'INVESTIGATION_VIOLATION')
      END
    END
  END
  ```
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 590-640 (Agent Termination section, Commission hold during investigation)
- **Traceability**: WF-IC-TERMINATION-001
- **Impact**: Protects company from paying fraudulent agents, enables investigation
- **Example**:
  - Investigation initiated for suspected fraud
  - Commissions on hold: ₹2,50,000
  - Investigation duration: 30 days
  - Outcome: No violation found
  - Action: ₹2,50,000 released in next cycle

---

**End of Section 2: Business Rules**

**Next Section**: [Functional Requirements](#3-functional-requirements)
# Incentive, Commission and Producer Management - Functional Requirements

## 3. Functional Requirements

### 3.3 Commission Processing Module

#### FR-IC-COM-001: Commission Rate Table Management
- **ID**: FR-IC-COM-001
- **Category**: Commission Configuration
- **Priority**: CRITICAL
- **Description**: System shall provide commission rate table view and configuration interface
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.1 (Commission Rate Table View Page), Lines 619-637
- **Traceability**: BR-IC-COM-006
- **Table Columns**:
  - Rate (%): Decimal
  - Policy Duration (Months): Integer
  - Product Type: Dropdown [PLI, RPLI]
  - Product Plan Code: Text
  - Agent Type: Dropdown [Direct Agent, Field Officer, etc.]
  - Policy Term (Years): Integer
- **Actions**: View only (configuration may be separate admin function)

#### FR-IC-COM-002: Commission Calculation Batch Job
- **ID**: FR-IC-COM-002
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: System shall support monthly commission calculation batch processing
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 3, FS_IC_024, Lines 218-220; Section 4.4.3.1 (Run Commission Calculation Batch Jobs), Lines 676-680
- **Traceability**: BR-IC-COM-001, BR-IC-COM-004
- **Trigger**: Scheduled (first working day of month) or Manual
- **Process**:
  1. Fetch all active policies for previous month
  2. Fetch agent details for each policy
  3. Lookup commission rate from rate table
  4. Calculate annualised premium
  5. Calculate commission amount = (Premium × Rate) / 100
  6. Store commission records with status = 'CALCULATED'
- **Monitoring**: Progress tracking, timeout after 6 hours, retry failed records

#### FR-IC-COM-003: Trial Statement Generation
- **ID**: FR-IC-COM-003
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: System shall generate trial commission statements for review and approval
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 3, FS_IC_025, Lines 222-224; Section 4.4.3.2 (Trial Statement Generation Batch Job), Lines 683-687
- **Traceability**: BR-IC-COM-002
- **Trigger**: After commission calculation batch completes
- **Process**:
  1. Group calculated commissions by agent
  2. Calculate agent-wise totals
  3. Apply TDS rules
  4. Generate trial statement for each agent
  5. Set status = 'PENDING_APPROVAL'
  6. Send notification to finance team

#### FR-IC-COM-004: Trial Statement View
- **ID**: FR-IC-COM-004
- **Category**: Commission Processing
- **Priority**: HIGH
- **Description**: System shall provide interface to view and filter trial statements
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.3 (View Trial Statement Page), Lines 689-711
- **Traceability**: BR-IC-COM-009
- **Filters**:
  - Agent ID
  - Policy Number
  - Circle
  - Commission Type
  - Statement Date Range
- **Table Columns**:
  - Agent ID
  - Agent Name
  - Policy Number
  - Commission Type
  - Calculated Amount
  - TDS Amount
  - Net Amount
  - Status (Pending/Approved)
  - Remarks
- **Actions**:
  - Export to Excel/PDF
  - Raise Correction
  - View Details

#### FR-IC-COM-005: Manual Trial Statement Generation
- **ID**: FR-IC-COM-005
- **Category**: Commission Processing
- **Priority**: MEDIUM
- **Description**: System shall allow manual trial statement generation with specific parameters
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.4 (Manual Trial Statement Generation Page), Lines 713-742
- **Traceability**: BR-IC-COM-002
- **Fields**:
  - Processing Unit: Dropdown [IT2.0, etc.]
  - Statement Format: Dropdown [Standard]
  - Max Statement Due Date: Date picker
  - Max Transaction Effective Date: Date picker
  - Max Process Date: Date picker
  - Statement Date: Date picker
  - Contract Holder: Text
  - Advisor Coordinator: Text
  - Carrier: Text
  - Tax Deduction (TDS %): Decimal
- **Actions**:
  - Generate Statement
  - Save Draft
  - Submit for Approval

#### FR-IC-COM-006: Trial Statement Approval
- **ID**: FR-IC-COM-006
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: System shall provide trial statement approval interface with full/partial disbursement options
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.5 (Approving Trial Statement Page), Lines 745-769
- **Traceability**: BR-IC-COM-002, BR-IC-COM-005
- **Table Columns**:
  - Agent ID
  - Policy Number
  - Commission Amount
  - Status
  - Remarks
  - Disbursement Option (Full/Partial)
  - Part Disbursement % (if Partial selected)
- **Actions**:
  - Apply Part Disbursement
  - Approve All Rows in P/U (Processing Unit)
  - Submit
- **On Approval**:
  - Update trial_statement_status = 'APPROVED'
  - Trigger final statement generation workflow

#### FR-IC-COM-007: Final Statement Generation
- **ID**: FR-IC-COM-007
- **Category**: Commission Processing
- **Priority**: CRITICAL
- **Description**: System shall generate final incentive statements after trial approval
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 3, FS_IC_029, Lines 243-245; Section 4.4.3.6 (Final Incentive Statement Generation Batch Job), Lines 771-777
- **Traceability**: BR-IC-COM-007
- **Trigger**: Scheduled or Manual after trial approval
- **Process**:
  1. Lock approved trial data (no further modifications)
  2. Generate final commission amounts
  3. Calculate final TDS
  4. Generate final statement PDFs
  5. Set status = 'READY_FOR_DISBURSEMENT'

#### FR-IC-COM-008: Final Statement View
- **ID**: FR-IC-COM-008
- **Category**: Commission Processing
- **Priority**: HIGH
- **Description**: System shall provide final statements viewing interface
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.7 (Final Statements Page), Lines 779-800
- **Traceability**: BR-IC-COM-007
- **Table Columns**:
  - Agent ID
  - Policy Number
  - Final Commission Amount
  - TDS Deducted
  - Net Payable
  - Payment Status (Pending/Queued/Disbursed)
- **Actions**:
  - View Statement (PDF)
  - Export PDF/Excel
  - Send to Disbursement

#### FR-IC-COM-009: Disbursement Details Entry
- **ID**: FR-IC-COM-009
- **Category**: Commission Disbursement
- **Priority**: CRITICAL
- **Description**: System shall provide disbursement details entry interface
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.8 (Disbursement Details Page), Lines 802-829
- **Traceability**: BR-IC-COM-008
- **Fields**:
  - Agent ID: Auto-populated
  - Payment Mode: Dropdown [Cheque, EFT]
  - **If Cheque**:
    - Cheque Number: Textbox
    - Bank Name: Textbox
    - Payment Date: Date picker
  - **If EFT**:
    - Bank Name or POSB: Dropdown/Textbox
    - IFSC Code: Textbox
    - Account Number: Textbox
  - Amount Paid: Decimal
  - Remarks: Textbox
- **Actions**:
  - Save
  - Submit
  - Generate Payment File (for EFT)

#### FR-IC-COM-010: Automatic Disbursement Processing
- **ID**: FR-IC-COM-010
- **Category**: Commission Disbursement
- **Priority**: CRITICAL
- **Description**: System shall support automatic disbursement for both cheque and EFT modes
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 3, FS_IC_032, Lines 251-258; Section 4.4.3.9 (Automatic Disbursement Option), Lines 833-844
- **Traceability**: BR-IC-COM-008
- **Process**:
  - **For Cheque**:
    - Mark disbursement as complete immediately
    - Generate payment advice
    - Send notification to agent
  - **For EFT**:
    - Generate payment file in PFMS/Bank format
    - Submit file to PFMS/Bank gateway
    - Set status = 'DISBURSEMENT_QUEUED'
    - Wait for confirmation callback
    - On confirmation: Set status = 'DISBURSED', notify agent

#### FR-IC-COM-011: Commission History Search
- **ID**: FR-IC-COM-011
- **Category**: Commission Reporting
- **Priority**: HIGH
- **Description**: System shall provide commission history search functionality
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 3, FS_IC_023, Lines 214-216; Section 4.4.2 (Commission History Search Page), Lines 642-667
- **Traceability**: BR-IC-COM-009
- **Search Filters**:
  - Agent ID
  - Policy Number
  - Date Range
  - Product Type
  - Commission Type (First Year, Renewal, Bonus)
- **Result Columns**:
  - Agent ID
  - Agent Name
  - Policy Number
  - Product Type
  - Commission Type
  - Amount
  - Status
  - Date Processed
- **Actions**:
  - Export to Excel/PDF
  - View Detailed Statement

#### FR-IC-COM-012: Statement Export
- **ID**: FR-IC-COM-012
- **Category**: Commission Reporting
- **Priority**: MEDIUM
- **Description**: System shall provide export of trial and final statements in multiple formats
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, referenced in Lines 707-711, 796-800
- **Traceability**: BR-IC-COM-010
- **Export Formats**:
  - Excel (with multiple sheets)
  - PDF (individual statements)
  - PDF (agent-wise consolidated)
- **Options**:
  - Export single statement
  - Export all statements in batch
  - Export filtered results

#### FR-IC-COM-011: Recovery from Future Commissions
- **ID**: FR-IC-COM-011
- **Category**: Commission Clawback
- **Priority**: CRITICAL
- **Description**: System shall recover clawback amount from agent's future commission payments
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 280-320
- **Traceability**: BR-IC-CLAWBACK-001
- **Functionality**:
  - Update agent clawback pending balance
  - Deduct from next commission payment
  - Track recovery progress
  - Update recovery status
- **Acceptance Criteria**:
  - Recovery initiated on next commission cycle
  - Partial recovery allowed (max 50% of current commission)
  - Recovery status accurately tracked
  - Balance zero when fully recovered

#### FR-IC-COM-012: Accounting Entry for Clawback
- **ID**: FR-IC-COM-012
- **Category**: Commission Clawback
- **Priority**: HIGH
- **Description**: System shall post accounting entries for clawback transactions
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 280-320
- **Traceability**: BR-IC-CLAWBACK-001
- **Functionality**:
  - Debit Commission Expense Reversal
  - Credit Agent Payable
  - Include reference to policy number
  - Pass effective date
- **Acceptance Criteria**:
  - Accounting entry posted on clawback creation
  - Double-entry bookkeeping maintained
  - Integration with Finance system verified

#### FR-IC-COM-013: Create Suspense Entry
- **ID**: FR-IC-COM-013
- **Category**: Commission Suspense
- **Priority**: CRITICAL
- **Description**: System shall create suspense entry for disputed/investigated policies
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 320-360
- **Traceability**: BR-IC-SUSPENSE-001, BR-IC-SUSPENSE-005
- **Functionality**:
  - Check if commission paid or pending
  - Create suspense record if paid
  - Mark as HELD_IN_SUSPENSE if pending
  - Update agent suspense balance
- **Acceptance Criteria**:
  - Suspense created within 1 hour of investigation trigger
  - Accurate amount tracked
  - Agent commission status updated

#### FR-IC-COM-014: Release Suspense on Clearance
- **ID**: FR-IC-COM-014
- **Category**: Commission Suspense
- **Priority**: HIGH
- **Description**: System shall release suspense amount when investigation clears policy
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 320-360
- **Traceability**: BR-IC-SUSPENSE-001, BR-IC-SUSPENSE-005
- **Functionality**:
  - Update suspense status to RELEASED
  - Release commission for payment
  - Update agent suspense balance
  - Send notification to agent
- **Acceptance Criteria**:
  - Release triggered within 1 hour of investigation closure
  - Payment processed in next cycle
  - Agent notified

#### FR-IC-COM-015: Forfeit Suspense on Fraud
- **ID**: FR-IC-COM-015
- **Category**: Commission Suspense
- **Priority**: CRITICAL
- **Description**: System shall forfeit suspense amount if policy determined fraudulent
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 320-360
- **Traceability**: BR-IC-SUSPENSE-001, BR-IC-SUSPENSE-005
- **Functionality**:
  - Update suspense status to FORFEITED
  - Update commission status to FORFEITED
  - Flag agent for compliance review
  - Escalate to compliance team
- **Acceptance Criteria**:
  - Forfeiture processed within 1 hour of fraud determination
  - Agent flagged for review
  - Compliance team notified

#### FR-IC-COM-016: Payment Retry Logic
- **ID**: FR-IC-COM-016
- **Category**: Commission Disbursement
- **Priority**: HIGH
- **Description**: System shall implement retry logic for failed commission payments
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 255-280
- **Traceability**: BR-IC-SUSPENSE-002
- **Functionality**:
  - Retry 3 times with exponential backoff (2, 4, 8 hours)
  - Log each retry attempt
  - Track failure count and reasons
  - Notify agent of retries
- **Acceptance Criteria**:
  - Retry scheduled immediately after failure
  - Correct delay applied (2^n hours)
  - All retries logged
  - Agent notified on each retry

#### FR-IC-COM-017: Manual Intervention Trigger
- **ID**: FR-IC-COM-017
- **Category**: Commission Disbursement
- **Priority**: MEDIUM
- **Description**: System shall trigger manual intervention after 3 failed payment retries
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 255-280
- **Traceability**: BR-IC-SUSPENSE-002
- **Functionality**:
  - Create suspense entry after 3 failures
  - Create finance task for resolution
  - Notify finance team
  - Notify agent of manual review required
- **Acceptance Criteria**:
  - Finance task created within 30 minutes of 3rd failure
  - Task assigned to correct team
  - 2-day SLA set for resolution
  - Agent notification sent

#### FR-IC-COM-018: Suspense Account Balance Tracking
- **ID**: FR-IC-COM-018
- **Category**: Commission Suspense
- **Priority**: HIGH
- **Description**: System shall track suspense account balances for each agent
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, Lines 320-360
- **Traceability**: BR-IC-SUSPENSE-001, BR-IC-SUSPENSE-003, BR-IC-SUSPENSE-004, BR-IC-SUSPENSE-005
- **Functionality**:
  - Maintain suspense_amount field in agent_profiles
  - Update on suspense creation
  - Update on suspense release/forfeiture
  - Display in agent profile
- **Acceptance Criteria**:
  - Balance accurately reflects all suspense transactions
  - Real-time updates
  - Audit trail maintained

#### FR-IC-COM-019: Suspense Aging Report
- **ID**: FR-IC-COM-019
- **Category**: Commission Reporting
- **Priority**: MEDIUM
- **Description**: System shall generate aging report for suspense accounts
- **Source**: `Agent_SRS_Incentive, Commission and Producer Management.md`, implied in suspense functionality
- **Traceability**: BR-IC-SUSPENSE-001
- **Functionality**:
  - Report suspense entries by age buckets (0-30, 31-60, 61-90, 90+ days)
  - Group by suspense reason
  - Show total amount by bucket
  - Export to Excel/PDF
- **Acceptance Criteria**:
  - Report available on demand
  - Accurate aging calculation
  - Drill-down to individual entries

---

**End of Section 3: Functional Requirements**

**Previous Section**: [Business Rules](#2-business-rules)
**Next Section**: [Validation Rules](#4-validation-rules)
# Incentive, Commission and Producer Management - Validation Rules

## 4. Validation Rules

### 4.1 Agent Profile Validation Rules

#### VR-IC-PROF-001: PAN Format Validation
- **ID**: VR-IC-PROF-001
- **Category**: PAN Validation
- **Priority**: CRITICAL
- **Rule**: `PAN must match pattern: [A-Z]{5}[0-9]{4}[A-Z]{1}`
- **Error Message**: "Please enter correct PAN"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.1.2, Error Messages 3-4
- **Traceability**: BR-IC-VAL-001
- **Test Cases**:
  - Valid: "ABCDE1234F" → Pass
  - Invalid: "12345ABCDE" → Fail
  - Invalid: "ABC12345FG" → Fail
  - Invalid: "ABCDE12345" → Fail (11 chars)

#### VR-IC-PROF-002: PAN Uniqueness Check
- **ID**: VR-IC-PROF-002
- **Category**: PAN Validation
- **Priority**: CRITICAL
- **Rule**: `PAN must not exist in agents table for any other agent_id`
- **Error Message**: "PAN number entered already exists for another advisor's profile and cannot be for this profile"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.1.2, Error Message 2
- **Traceability**: BR-IC-PROF-002
- **Test Cases**:
  - New PAN "XYZAB5678C" → Pass (not in database)
  - Existing PAN "ABCDE1234F" → Fail (already assigned to Agent AGT001)

#### VR-IC-PROF-003: First Name Mandatory
- **ID**: VR-IC-PROF-003
- **Category**: Personal Details Validation
- **Priority**: CRITICAL
- **Rule**: `first_name IS NOT NULL AND first_name != ''`
- **Error Message**: "Please enter a First name"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.1.2, Error Message 6
- **Traceability**: BR-IC-VAL-002
- **Test Cases**:
  - First Name = "Rajesh" → Pass
  - First Name = "" → Fail
  - First Name = NULL → Fail

#### VR-IC-PROF-004: Last Name Mandatory
- **ID**: VR-IC-PROF-004
- **Category**: Personal Details Validation
- **Priority**: CRITICAL
- **Rule**: `last_name IS NOT NULL AND last_name != ''`
- **Error Message**: "Please enter a Last name"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.1.2, Error Message 5
- **Traceability**: BR-IC-VAL-002
- **Test Cases**:
  - Last Name = "Kumar" → Pass
  - Last Name = "" → Fail

#### VR-IC-PROF-005: Date of Birth Validation
- **ID**: VR-IC-PROF-005
- **Category**: Personal Details Validation
- **Priority**: CRITICAL
- **Rule**: `date_of_birth IS NOT NULL AND date_of_birth < current_date - 18 years`
- **Error Message**: "Please enter a valid Date of Birth"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.1.2, Error Message 7
- **Traceability**: BR-IC-VAL-003
- **Test Cases**:
  - DOB = "1985-05-15" (Age > 18) → Pass
  - DOB = "2010-05-15" (Age < 18) → Fail
  - DOB = NULL → Fail
  - DOB = "2027-01-01" (Future date) → Fail

#### VR-IC-PROF-006: Mobile Number Format
- **ID**: VR-IC-PROF-006
- **Category**: Contact Validation
- **Priority**: HIGH
- **Rule**: `mobile_number MUST match pattern: [6-9][0-9]{9}`
- **Error Message**: "Please enter a valid 10-digit mobile number"
- **Source**: Standard Indian mobile validation (referenced in profile section)
- **Traceability**: BR-IC-PROF-001
- **Test Cases**:
  - "9876543210" → Pass
  - "8765432109" → Pass
  - "1234567890" → Fail (doesn't start with 6-9)
  - "987654321" → Fail (9 digits)
  - "98765432101" → Fail (11 digits)

#### VR-IC-PROF-007: Email Format Validation
- **ID**: VR-IC-PROF-007
- **Category**: Contact Validation
- **Priority**: MEDIUM
- **Rule**: `email MUST match standard email format regex`
- **Error Message**: "Please enter a valid email address"
- **Source**: Standard email validation (referenced in profile section)
- **Traceability**: FR-IC-ONB-002
- **Test Cases**:
  - "agent@example.com" → Pass
  - "rajesh.kumar@postal.gov.in" → Pass
  - "invalid-email" → Fail
  - "@example.com" → Fail

---

### 4.2 Commission Validation Rules

#### VR-IC-COM-001: Commission Rate Required
- **ID**: VR-IC-COM-001
- **Category**: Commission Configuration
- **Priority**: CRITICAL
- **Rule**: `Commission rate MUST exist in rate_table for given combination`
- **Error Message**: "Commission rate not found for Product Type: {type}, Plan: {plan}, Agent Type: {agent_type}"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.1
- **Traceability**: BR-IC-COM-006
- **Test Cases**:
  - Valid combination in rate table → Pass
  - Missing combination → Fail with error

#### VR-IC-COM-002: TDS Rate Validation
- **ID**: VR-IC-COM-002
- **Category**: Commission Calculation
- **Priority**: HIGH
- **Rule**: `TDS rate MUST be between 0% and 30%`
- **Error Message**: "TDS rate must be between 0% and 30%"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.4
- **Traceability**: BR-IC-COM-003
- **Test Cases**:
  - TDS = 5% → Pass
  - TDS = 0% → Pass (no TDS)
  - TDS = 30% → Pass (maximum)
  - TDS = 31% → Fail

#### VR-IC-COM-003: Disbursement Amount Validation
- **ID**: VR-IC-COM-003
- **Category**: Disbursement Validation
- **Priority**: HIGH
- **Rule**: `disbursement_amount MUST be <= net_commission_amount`
- **Error Message**: "Disbursement amount cannot exceed net commission amount"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.5
- **Traceability**: BR-IC-COM-005
- **Test Cases**:
  - Commission = 50,000, Disbursement = 50,000 → Pass (Full)
  - Commission = 50,000, Disbursement = 30,000 → Pass (Partial)
  - Commission = 50,000, Disbursement = 60,000 → Fail

#### VR-IC-COM-004: Bank Account Validation
- **ID**: VR-IC-COM-004
- **Category**: Disbursement Validation
- **Priority**: CRITICAL
- **Rule**: `For EFT: bank_account_number AND ifsc_code or posb_account_number ARE mandatory`
- **Error Message**: "Bank Account Number and IFSC Code or POSB Account Number are required for EFT payment"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3.8
- **Traceability**: FR-IC-COM-009
- **Test Cases**:
  - EFT with Account + IFSC → Pass
  - EFT with POSB Account → Pass
  - EFT with Account only → Fail
  - EFT with IFSC only → Fail
  - Cheque without Bank details → Pass (not required)

#### VR-IC-COM-005: Trial Statement Approval Validation
- **ID**: VR-IC-COM-005
- **Category**: Disbursement Validation
- **Priority**: HIGH
- **Rule**: `Final disbursement ONLY allowed if trial_statement_status = 'APPROVED'`
- **Error Message**: "Trial statement must be approved before disbursement"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.4.3
- **Traceability**: BR-IC-COM-002
- **Test Cases**:
  - Status = 'APPROVED' → Allow disbursement
  - Status = 'PENDING' → Block disbursement
  - Status = 'REJECTED' → Block disbursement

---

### 4.3 License Validation Rules

#### VR-IC-LIC-001: License Number Format
- **ID**: VR-IC-LIC-001
- **Category**: License Validation
- **Priority**: HIGH
- **Rule**: `license_number MUST be alphanumeric and 5-20 characters`
- **Error Message**: "License number must be 5-20 alphanumeric characters"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.2.3
- **Traceability**: FR-IC-PROF-004
- **Test Cases**:
  - "LIC12345" → Pass
  - "IRDA-12345-2026" → Pass
  - "1234" → Fail (too short)
  - "" → Fail (empty)

#### VR-IC-LIC-002: Renewal Date Future Validation
- **ID**: VR-IC-LIC-002
- **Category**: License Validation
- **Priority**: HIGH
- **Rule**: `For new license: renewal_date > license_date`
- **Error Message**: "Renewal date must be after license issue date"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.2.3
- **Traceability**: BR-IC-LIC-001
- **Test Cases**:
  - License: 2026-01-01, Renewal: 2027-01-01 → Pass
  - License: 2026-01-01, Renewal: 2025-01-01 → Fail

#### VR-IC-LIC-003: License Expiry Check
- **ID**: VR-IC-LIC-003
- **Category**: License Validation
- **Priority**: CRITICAL
- **Rule**: `Block agent operations IF current_date > renewal_date AND license_status != 'RENEWED'`
- **Error Message**: "License expired. Please renew to continue operations"
- **Source**: `Agent_SRS_Incentive-Commission-and-Producer-Management.md`, Section 4.2.3
- **Traceability**: BR-IC-LIC-003
- **Test Cases**:
  - Current: 2026-01-15, Renewal: 2027-01-01 → Active
  - Current: 2027-01-16, Renewal: 2027-01-15, Not Renewed → Deactivated

---

**End of Section 4: Validation Rules**

**Previous Section**: [Functional Requirements](#3-functional-requirements)
**Next Section**: [Error Codes](#5-error-codes)
# Incentive, Commission and Producer Management - Error Codes

## 5. Error Codes

| Error Code | Error Message | Category | Severity | Trigger Condition | Related BR |
|------------|---------------|----------|----------|-------------------|------------|
| **IC-ERR-001** | Please select a Profile Type | Agent Onboarding | ERROR | Profile Type not selected | BR-IC-VAL-002 |
| **IC-ERR-002** | PAN number entered already exists for another advisor's profile and cannot be for this profile | Agent Profile | ERROR | Duplicate PAN | BR-IC-PROF-002 |
| **IC-ERR-003** | Please enter a 10 digit Permanent Account Number (PAN) | Agent Profile | ERROR | Invalid PAN length | BR-IC-VAL-001 |
| **IC-ERR-004** | Please enter correct PAN | Agent Profile | ERROR | Invalid PAN format | BR-IC-VAL-001 |
| **IC-ERR-005** | Please enter a Last name | Agent Profile | ERROR | Last Name missing | BR-IC-VAL-002 |
| **IC-ERR-006** | Please enter a First name | Agent Profile | ERROR | First Name missing | BR-IC-VAL-002 |
| **IC-ERR-007** | Please enter a valid Date of Birth | Agent Profile | ERROR | Invalid/missing DOB | BR-IC-VAL-003 |
| **IC-ERR-008** | Your selected criteria did not return any rows. Please change your selections and try again | Agent Search | WARNING | No search results found | BR-IC-PROF-001 |
| **IC-ERR-009** | Coordinator ID is mandatory | Agent Onboarding | ERROR | Creating Advisor without Coordinator | BR-IC-AH-001 |
| **IC-ERR-010** | Circle assignment required | Agent Onboarding | ERROR | Creating Coordinator without Circle | BR-IC-AH-002 |
| **IC-ERR-011** | Employee ID not found in HRMS | Agent Onboarding | ERROR | Invalid Employee ID | BR-IC-AH-003 |
| **IC-ERR-012** | Trial statement must be approved before disbursement | Commission | ERROR | Attempting disbursement without approval | BR-IC-COM-002 |
| **IC-ERR-013** | Commission rate not found for given parameters | Commission | ERROR | Missing rate configuration | BR-IC-COM-006 |
| **IC-ERR-014** | Disbursement amount cannot exceed commission amount | Disbursement | ERROR | Invalid disbursement amount | BR-IC-COM-005 |
| **IC-ERR-015** | Bank Account Number and IFSC Code are required for EFT payment | Disbursement | ERROR | Missing bank details for EFT | BR-IC-COM-008 |
| **IC-ERR-016** | License has expired. Agent code deactivated | License | CRITICAL | Operating with expired license | BR-IC-LIC-003 |

### Error Handling Patterns

#### 1. Onboarding Errors (IC-ERR-001 to IC-ERR-011)
**Pattern**: Block save, display error message, allow correction

**User Experience**:
- Error displayed inline with field highlighting
- User can correct and resubmit
- No data loss

**Example**:
```
User Action: Click "Continue" without selecting Agent Type
System Response:
- Display IC-ERR-001: "Please select a Profile Type"
- Highlight Agent Type dropdown in red
- Keep all other entered data
- Allow user to select and retry
```

#### 2. Commission Processing Errors (IC-ERR-012 to IC-ERR-014)
**Pattern**: Block operation, notify user, log for audit

**User Experience**:
- Operation blocked with clear message
- Finance team notified for manual intervention
- Audit trail created

**Example**:
```
User Action: Attempt disbursement without trial approval
System Response:
- Display IC-ERR-012: "Trial statement must be approved before disbursement"
- Block disbursement button
- Send notification to Finance Team
- Log attempt in audit trail
```

#### 3. License Errors (IC-ERR-016)
**Pattern**: Auto-deactivate, notify all stakeholders

**User Experience**:
- Agent code deactivated automatically
- Notifications sent to Agent, Coordinator, Operations
- Manual reactivation required after renewal

**Example**:
```
Event: License expiry detected in batch job
System Response:
- Deactivate agent code
- Send email to Agent: "Your license has expired. Agent code deactivated."
- Send email to Coordinator: "Agent {name} deactivated due to license expiry."
- Update status to 'DEACTIVATED'
- Log event
```

### Error Code Ranges

| Range | Category | Usage |
|-------|----------|-------|
| IC-ERR-001 to 011 | Agent Onboarding & Profile | Data entry and validation errors |
| IC-ERR-012 to 015 | Commission Processing | Calculation and disbursement errors |
| IC-ERR-016 to 020 | License Management | License validation and expiry errors |
| IC-ERR-021 to 030 | Integration Errors | HRMS, PFMS, Policy Services errors |
| IC-ERR-031 to 040 | Batch Processing Errors | Commission batch, reminder batch errors |
| IC-ERR-041 to 050 | System Errors | Database, network, configuration errors |

### Integration Error Codes (Planned)

| Error Code | Error Message | Integration |
|------------|---------------|-------------|
| IC-ERR-021 | HRMS system unavailable. Please try again later | HRMS Integration |
| IC-ERR-022 | Employee ID not found in HRMS | HRMS Integration |
| IC-ERR-023 | PFMS gateway unavailable. Payment queued for retry | PFMS Integration |
| IC-ERR-024 | Policy service timeout. Commission calculation queued | Policy Services |
| IC-ERR-025 | Accounting service unavailable. Voucher queued | Accounting Integration |

---

**End of Section 5: Error Codes**

**Previous Section**: [Validation Rules](#4-validation-rules)
**Next Section**: [Workflows](#6-workflows)
# Incentive, Commission and Producer Management - Workflows

## 6. Workflows

### 6.1 Agent Onboarding Workflow

#### WF-IC-ONB-001: New Agent Onboarding Process

**Description**: End-to-end workflow for onboarding a new agent

**Participants**:
- Operations User
- Advisor Coordinator (for approval)
- HRMS System (for departmental employees)

**Preconditions**:
- Advisor Coordinator exists (for Advisor onboarding)
- HRMS system is accessible (for Departmental Employee)

**Process Flow**:

```
1. START: Operations user selects "New Agent Profile"

2. SELECT AGENT TYPE
   ├─ Advisor → Proceed to Step 3
   ├─ Advisor Coordinator → Proceed to Step 3
   ├─ Departmental Employee → Proceed to Step 4
   └─ Field Officer → Proceed to Step 5

3. ADVISOR ONBOARDING
   ├─ Enter Profile Details
   ├─ Select Advisor Coordinator (BR-IC-AH-001)
   ├─ Validate PAN uniqueness (VR-IC-PROF-002)
   ├─ Validate mandatory fields (VR-IC-PROF-003, VR-IC-PROF-004, VR-IC-PROF-005)
   ├─ Save Profile
   └─ Generate Agent Code → END

4. DEPARTMENTAL EMPLOYEE ONBOARDING
   ├─ Enter Employee ID
   ├─ Fetch from HRMS (BR-IC-AH-003)
   ├─ IF Employee Found:
   │  ├─ Auto-populate profile
   │  ├─ User completes remaining fields
   │  ├─ Validate and Save
   │  └─ Generate Agent Code → END
   └─ ELSE: Display error (IC-ERR-011)

5. FIELD OFFICER ONBOARDING (BR-IC-AH-004)
   ├─ Option A: Enter Employee ID → Go to Step 4
   └─ Option B: Manual Entry → Go to Step 3

6. END: Agent profile created successfully
```

**Error Handling**:
- IC-ERR-001: Profile Type not selected
- IC-ERR-002: Duplicate PAN
- IC-ERR-009: No Coordinator selected for Advisor
- IC-ERR-011: Employee ID not found

**Success Criteria**:
- Agent profile saved with unique Agent Code
- All mandatory fields validated
- Link to Advisor Coordinator established (for Advisors)

**Estimated Time**: 15-30 minutes per profile

---

### 6.2 Commission Processing Workflow

#### WF-IC-COM-001: Monthly Commission Processing

**Description**: Automated workflow for calculating and processing commissions

**Participants**:
- Batch Scheduler
- Commission Calculation Engine
- Finance Team
- Disbursement System
- PFMS/Bank Gateway

**Preconditions**:
- Commission rate table configured
- Active policies exist for the period
- Agent profiles valid

**Process Flow**:

```
1. START: Commission batch triggered (first working day of month)

2. COMMISSION CALCULATION BATCH (FR-IC-COM-002)
   ├─ Fetch all active policies for previous month
   ├─ FOR EACH policy:
   │  ├─ Fetch agent details
   │  ├─ Lookup commission rate (BR-IC-COM-006)
   │  ├─ Lookup premium
   │  ├─ Calculate commission = (Annualised Premium × Rate) / 100
   │  ├─ Store commission_record (status='CALCULATED')
   │  └─ Monitor progress (BR-IC-COM-012)
   ├─ Handle failed records (retry up to 3 times)
   └─ On completion → Trigger Step 3

3. TRIAL STATEMENT GENERATION (FR-IC-COM-003)
   ├─ Group commissions by agent
   ├─ Calculate agent-wise totals
   ├─ Apply TDS calculation (BR-IC-COM-003)
   ├─ Generate trial statements (status='PENDING_APPROVAL')
   └─ Notify Finance Team → Wait for approval

4. FINANCE REVIEW
   ├─ View Trial Statements (FR-IC-COM-004)
   ├─ Filter by Agent/Circle/Commission Type
   ├─ Raise corrections if needed
   └─ Approve Statements → Trigger Step 5

5. TRIAL STATEMENT APPROVAL (FR-IC-COM-006)
   ├─ For each approved statement:
   │  ├─ Select disbursement type (Full/Partial)
   │  ├─ IF Partial: Enter percentage (BR-IC-COM-005)
   │  └─ Submit approval
   ├─ Update status = 'APPROVED'
   └─ Trigger Step 6

6. FINAL STATEMENT GENERATION (FR-IC-COM-007)
   ├─ Lock approved trial data (BR-IC-COM-002)
   ├─ Generate final commission amounts
   ├─ Generate final TDS amounts
   ├─ Create final statement PDFs
   ├─ Set status = 'READY_FOR_DISBURSEMENT'
   └─ Notify Finance Team → Wait for disbursement

7. DISBURSEMENT PROCESSING (FR-IC-COM-009, FR-IC-COM-010)
   ├─ Enter disbursement details
   ├─ FOR EACH final statement:
   │  ├─ IF Payment Mode = Cheque:
   │  │  ├─ Enter Cheque Number, Bank Name, Payment Date
   │  │  ├─ Mark disbursed immediately
   │  │  └─ Send notification to agent
   │  └─ IF Payment Mode = EFT:
   │  │  ├─ Validate bank details (VR-IC-COM-004)
   │  │  ├─ Generate payment file
   │  │  ├─ Send to PFMS/Bank
   │  │  ├─ Set status = 'DISBURSEMENT_QUEUED'
   │  │  ├─ Wait for confirmation callback
   │  │  ├─ On confirmation: Mark 'DISBURSED', notify agent
   │  │  └─ Monitor SLA (BR-IC-COM-011)
   └─ Continue until all processed

8. END: Commission disbursement complete
```

**Error Handling**:
- IC-ERR-012: Disbursement without trial approval
- IC-ERR-013: Commission rate not found
- IC-ERR-014: Disbursement amount exceeds commission
- IC-ERR-015: Bank details missing for EFT

**Success Criteria**:
- All eligible agents receive commission statements
- Trial statements approved within SLA
- Disbursements completed within 10 working days (BR-IC-COM-011)
- TDS correctly deducted and reported

**SLA Summary**:
- Batch calculation: 6 hours
- Trial approval: 7 days
- Disbursement: 10 working days

---

### 6.3 License Renewal Workflow

#### WF-IC-LIC-001: License Renewal Process

**Description**: Workflow for managing license renewals and reminders

**Participants**:
- Batch Scheduler
- Agent
- Licensing Authority
- Operations Team

**Preconditions**:
- Agent has active license
- Renewal tracking enabled

**Process Flow**:

```
1. START: License with future renewal date in system

2. RENEWAL REMINDER SCHEDULE (BR-IC-LIC-004)
   ├─ T-30 days: Send first reminder
   ├─ T-15 days: Send second reminder
   ├─ T-7 days: Send third reminder
   └─ T-0 days: Send final reminder (expiry day)

3. AGENT SUBMITS RENEWAL REQUEST
   ├─ Upload renewal documents
   ├─ Enter new renewal date
   └─ Submit to Operations

4. OPERATIONS REVIEW (BR-IC-LIC-005)
   ├─ Validate documents complete
   ├─ Verify with licensing authority (if needed)
   ├─ IF documents valid AND no discrepancies:
   │  ├─ Auto-approve renewal
   │  ├─ Update license.renewal_date
   │  └─ Send confirmation to agent → END
   └─ ELSE: Manual review required

5. MANUAL REVIEW
   ├─ Operations reviews documents
   ├─ Request additional documents if needed
   ├─ Approve or Reject renewal
   └─ IF Approved: Go to Step 4 (Auto-approve path)

6. LICENSE EXPIRY CHECK (Daily Batch)
   ├─ FOR EACH license:
   │  ├─ IF current_date > renewal_date AND status != 'RENEWED':
   │  │  ├─ Deactivate agent code (BR-IC-LIC-003)
   │  │  ├─ Send expiry notification to agent
   │  │  └─ Notify Advisor Coordinator
   │  └─ ELSE: Continue monitoring
   └─ Repeat daily

7. END: License renewed or agent deactivated
```

**Error Handling**:
- Renewal documents incomplete
- Renewal date in past
- Agent operations blocked after expiry

**Success Criteria**:
- Agents receive timely reminders
- Renewals processed within 3 working days (BR-IC-LIC-005)
- No unlicensed agents can operate
- Auto-deactivation on expiry

---

### 6.4 Agent Termination Workflow

#### WF-IC-TERM-001: Agent Termination Process

**Description**: Workflow for terminating agent profiles

**Participants**:
- Operations User
- Advisor Coordinator
- Finance Team

**Preconditions**:
- Agent exists in system
- Termination reason documented

**Process Flow**:

```
1. START: Operations user initiates agent termination

2. SEARCH AGENT
   └─ Open Agent Profile Maintenance Page

3. INITIATE TERMINATION
   ├─ Select "Change Status To" = Terminated
   ├─ Click "Update"
   └─ Open Termination Page

4. ENTER TERMINATION DETAILS (FR-IC-PROF-005)
   ├─ Status Reason (Mandatory)
   ├─ Status Date
   ├─ Effective Date
   └─ Termination Date

5. VALIDATE INPUT
   ├─ All mandatory fields present (BR-IC-PROF-004)
   └─ Proceed to Step 6

6. PROCESS TERMINATION (BR-IC-PROF-004)
   ├─ Update agent.status = 'TERMINATED'
   ├─ Deactivate agent code
   ├─ Cancel all pending commissions
   ├─ Archive profile data
   ├─ Send notifications:
   │  ├─ To Agent
   │  ├─ To Advisor Coordinator
   │  └─ To Finance Team
   └─ Log audit trail

7. END: Agent termination complete
```

**Error Handling**:
- Missing termination reason
- Termination date before effective date

**Success Criteria**:
- Agent code deactivated
- Pending commissions cancelled
- All notifications sent
- Audit trail complete

---

### 6.5 Trial Statement Approval Workflow

#### WF-IC-TRIAL-001: Trial Statement Review and Approval

**Description**: Finance team workflow for reviewing and approving trial statements

**Participants**:
- Finance Manager
- Finance Head
- Operations Team

**Process Flow**:

```
1. START: Trial statements generated

2. STATEMENT DISTRIBUTION
   ├─ Email notifications sent to Finance Team
   └─ Statements available in portal

3. PRELIMINARY REVIEW
   ├─ Finance Manager reviews statements
   ├─ Filters by Circle, Agent Type, Amount Range
   ├─ Identifies anomalies or discrepancies
   └─ Two paths:

4. PATH A: RAISE CORRECTION
   ├─ Mark statement as "Needs Correction"
   ├─ Add correction comments
   ├─ Notify Operations Team
   ├─ Operations investigates and recalculates
   ├─ New trial statement generated
   └─ Return to Step 2

5. PATH B: APPROVE STATEMENT
   ├─ Finance Manager reviews statement details
   ├─ Select disbursement type:
   │  ├─ Full Disbursement (100%)
   │  └─ Partial Disbursement (enter %)
   ├─ Add approval remarks (optional)
   ├─ Submit approval
   ├─ Status changes to "APPROVED"
   └─ Trigger Final Statement Generation

6. ESCALATION (If needed)
   ├─ For amounts > ₹10 Lakhs: Finance Head approval
   ├─ For disputed amounts: Finance Committee review
   └─ Special cases: Director approval

7. END: Trial statement approved
```

**Decision Points**:
- Amount threshold for escalation: ₹10 Lakhs
- Correction vs Approval decision
- Full vs Partial disbursement

**Success Criteria**:
- All statements reviewed within 7 days
- Corrections addressed promptly
- Proper audit trail maintained

---

### 6.6 Disbursement Processing Workflow

#### WF-IC-DISB-001: Payment Disbursement

**Description**: End-to-end disbursement workflow for both cheque and EFT modes

**Process Flow**:

```
1. START: Final statements ready for disbursement

2. PAYMENT MODE SELECTION
   ├─ Cheque → Go to Step 3
   └─ EFT → Go to Step 4

3. CHEQUE DISBURSEMENT
   ├─ Generate cheque payment advice
   ├─ Print cheque details
   ├─ Record cheque number
   ├─ Mark disbursement as "COMPLETED"
   ├─ Send notification to agent
   └─ Update accounting → END

4. EFT DISBURSEMENT
   ├─ Validate bank details (account number, IFSC)
   ├─ Generate payment file (PFMS format)
   ├─ Upload to PFMS/Bank gateway
   ├─ Set status to "DISBURSEMENT_QUEUED"
   ├─ Wait for confirmation (2-3 days)
   ├─ On success:
   │  ├─ Update status to "DISBURSED"
   │  ├─ Send confirmation to agent
   │  └─ Update accounting
   ├─ On failure:
   │  ├─ Log error details
   │  ├─ Notify Finance Team
   │  ├─ Queue for retry
   │  └─ Manual intervention if 3 retries fail
   └─ END

5. SLA MONITORING
   ├─ Track time from trial approval to disbursement
   ├─ T-3 days: Send reminder if pending
   ├─ T+0 days: Escalate if overdue
   └─ Calculate penalty if breached (8% interest)
```

**Error Handling**:
- Invalid bank details
- PFMS gateway timeout
- Payment failure from bank
- Account closed/invalid

**Success Criteria**:
- Cheque: Same day completion
- EFT: 2-3 days completion
- SLA: 10 working days
- Failure rate < 1%

---

### 6.7 Commission History Inquiry Workflow

#### WF-IC-HIST-001: Agent Commission History

**Description**: Self-service workflow for agents to view their commission history

**Process Flow**:

```
1. START: Agent logs into Agent Portal

2. NAVIGATE TO COMMISSION HISTORY
   └─ Menu: My Commissions → Commission History

3. APPLY SEARCH FILTERS
   ├─ Date Range (default: last 12 months)
   ├─ Commission Type (First Year, Renewal, Bonus)
   ├─ Policy Number (optional)
   └─ Status (All, Paid, Pending)

4. VIEW RESULTS
   ├─ Table displays commission records
   ├─ Columns: Date, Policy, Type, Amount, TDS, Net, Status
   ├─ Sortable by any column
   └─ Pagination for large datasets

5. DRILL-DOWN DETAILS
   ├─ Click on record to view details
   ├─ Shows policy details
   ├─ Shows calculation breakdown
   ├─ Shows TDS deduction
   ├─ Shows payment details
   └─ Download PDF statement

6. EXPORT OPTIONS
   ├─ Export to Excel (with all fields)
   ├─ Download PDF statement
   └─ Print-friendly view

7. END: Agent completes inquiry
```

**Use Cases**:
- Agent verifying commission received
- Reconciling with bank statements
- Tax planning and TDS tracking
- Dispute resolution

**Performance Requirements**:
- Search response: < 2 seconds
- Export generation: < 5 seconds
- Support 10,000+ records per agent

---

### 6.8 License Renewal Submission Workflow

#### WF-IC-LIC-002: Agent License Renewal

**Description**: Agent-initiated workflow for submitting license renewal request

**Process Flow**:

```
1. START: Agent logs into Agent Portal

2. NAVIGATE TO LICENSES
   └─ Menu: My Profile → Licenses

3. VIEW LICENSE STATUS
   ├─ Display active licenses
   ├─ Show renewal due date
   ├─ Show expiry countdown
   └─ "Renew Now" button if within 60 days of expiry

4. INITIATE RENEWAL
   ├─ Click "Renew Now"
   ├─ System displays renewal checklist:
   │  ├─ Renewal application form
   │  ├─ Training completion certificate
   │  ├─ Previous license copy
   │  └─ Renewal fee payment receipt

5. UPLOAD DOCUMENTS
   ├─ Upload each document (PDF/JPG, max 5MB)
   ├─ System validates file type and size
   ├─ Enter new renewal date
   └─ Submit application

6. OPERATIONS REVIEW
   ├─ Operations receives notification
   ├─ Reviews uploaded documents
   ├─ Auto-approve if complete and valid
   └─ Manual review if discrepancies found

7. APPROVAL/REJECTION
   ├─ IF Approved:
   │  ├─ License updated with new renewal date
   │  ├─ Email confirmation to agent
   │  └─ Reminder schedule reset
   └─ IF Rejected:
   │  ├─ Email with reason for rejection
   │  ├─ Documents returned to agent
   │  └─ Agent can resubmit after corrections

8. END: Renewal processed
```

**Success Criteria**:
- Processing within 3 working days
- Auto-approval rate > 80%
- Rejection with clear reasons
- Document upload success rate > 95%

---

### 6.9 Commission Clawback Workflow

**Workflow ID**: `WF-IC-CLAWBACK-001`
**Name**: Commission Clawback on Policy Lapse
**Trigger**: Policy lapses (daily batch detection)
**Priority**: CRITICAL

**Workflow Diagram**:

```
START: Policy Lapse Detected
   │
   ├─ 1. DETECTION
   │  ├─ Daily batch job identifies lapsed policies
   │  ├─ Filter: lapse_date within last 24 hours
   │  └─ Check: clawback_processed = FALSE
   │
   ├─ 2. ELIGIBILITY CHECK
   │  ├─ Calculate months_active (issue_date to lapse_date)
   │  ├─ IF months_active >= 24: END (No clawback)
   │  └─ IF months_active < 24: Continue
   │
   ├─ 3. CALCULATE CLAWBACK PERCENTAGE
   │  ├─ < 6 months: 100%
   │  ├─ 6-12 months: 75%
   │  ├─ 12-18 months: 50%
   │  ├─ 18-24 months: 25%
   │  └─ ≥ 24 months: 0% (exit workflow)
   │
   ├─ 4. RETRIEVE ORIGINAL COMMISSION
   │  ├─ Query commission_transactions for first-year commission
   │  ├─ SUM all transactions where commission_year = 1
   │  └─ Get original_commission amount
   │
   ├─ 5. CALCULATE CLAWBACK AMOUNT
   │  └─ clawback_amount = original_commission × (clawback_percentage / 100)
   │
   ├─ 6. CREATE CLAWBACK ENTRY
   │  ├─ INSERT into commission_clawbacks table
   │  ├─ Set recovery_status = 'PENDING'
   │  └─ Log clawback_reason
   │
   ├─ 7. UPDATE AGENT BALANCE
   │  ├─ agent.clawback_pending_amount += clawback_amount
   │  └─ agent.commission_status = 'SUSPENDED_PENDING_CLAWBACK'
   │
   ├─ 8. POST ACCOUNTING ENTRY
   │  ├─ Debit: Commission Expense Reversal
   │  ├─ Credit: Agent Payable - {agent_code}
   │  └─ Reference: "Clawback - Policy {policy_number}"
   │
   ├─ 9. SEND NOTIFICATIONS
   │  ├─ Email to agent (HIGH priority)
   │  │  ├─ Subject: "Commission Clawback Notice"
   │  │  └─ Details: policy_number, amount, percentage, reason
   │  └─ Alert to Finance Team
   │
   ├─ 10. MARK POLICY PROCESSED
   │  ├─ policy.clawback_processed = TRUE
   │  └─ policy.clawback_amount = clawback_amount
   │
   └─ END: Clawback Initiated
```

**Success Criteria**:
- All lapsed policies detected within 24 hours
- Clawback calculation accurate (verified by audit)
- Agent notified within 1 hour of clawback creation
- Accounting entry posted successfully
- Recovery scheduled from next commission payment

**SLA**:
- Detection: Daily (within 24 hours of lapse)
- Processing: 2 hours from detection
- Notification: 1 hour from processing

---

### 6.10 Suspense Account Management Workflow

**Workflow ID**: `WF-IC-SUSPENSE-001`
**Name**: Commission Suspense for Disputed Policies
**Trigger**: Policy marked as UNDER_INVESTIGATION
**Priority**: CRITICAL

**Workflow Diagram**:

```
START: Policy Under Investigation
   │
   ├─ 1. TRIGGER DETECTION
   │  ├─ Policy status changed to 'UNDER_INVESTIGATION'
   │  ├─ Get investigation_reason and investigation_reference
   │  └─ Identify agent_code
   │
   ├─ 2. CHECK COMMISSION STATUS
   │  ├─ Query commission_transactions for this policy
   │  ├─ Check if payment_status = 'DISBURSED'
   │  └─ Get commission_paid amount
   │
   ├─ 3a. IF COMMISSION PAID
   │  ├─ CREATE SUSPENSE ENTRY
   │  │  ├─ INSERT into commission_suspense table
   │  │  ├─ Set status = 'SUSPENDED'
   │  │  ├─ Set suspense_reason = 'POLICY_UNDER_INVESTIGATION'
   │  │  └─ Set expected_resolution_date = current_date + 30 days
   │  │
   │  ├─ UPDATE AGENT STATUS
   │  │  ├─ agent.commission_status = 'SUSPENDED_PENDING_INVESTIGATION'
   │  │  └─ agent.suspense_amount += commission_paid
   │  │
   │  └─ Go to Step 5
   │
   ├─ 3b. IF COMMISSION NOT PAID
   │  ├─ HOLD PENDING COMMISSION
   │  │  ├─ UPDATE commission_transactions
   │  │  ├─ Set payment_status = 'HELD_IN_SUSPENSE'
   │  │  └─ Set hold_until = expected_investigation_closure_date
   │  │
   │  └─ Go to Step 5
   │
   ├─ 5. SEND NOTIFICATIONS
   │  ├─ Email to agent (HIGH priority)
   │  │  ├─ Subject: "Commission Suspended - Investigation"
   │  │  └─ Details: policy_number, amount, investigation_reason
   │  └─ Alert to Finance Team
   │
   ├─ 6. WAIT FOR INVESTIGATION OUTCOME
   │  ├─ Listen for policy status change to 'INVESTIGATION_CLOSED'
   │  ├─ Get investigation_outcome
   │  │  ├─ 'POLICY_GENUINE'
   │  │  ├─ 'POLICY_FRAUDULENT'
   │  │  └─ 'INCONCLUSIVE'
   │  └─ Maximum wait: 30 days (escalate if exceeded)
   │
   ├─ 7a. IF OUTCOME = POLICY_GENUINE
   │  ├─ RELEASE SUSPENSE
   │  │  ├─ UPDATE commission_suspense SET status = 'RELEASED'
   │  │  ├─ UPDATE commission_transactions SET payment_status = 'APPROVED_FOR_PAYMENT'
   │  │  ├─ agent.commission_status = 'ACTIVE'
   │  │  └─ agent.suspense_amount -= commission_amount
   │  │
   │  ├─ PROCESS PAYMENT
   │  │  └─ Include in next commission disbursement cycle
   │  │
   │  └─ NOTIFY AGENT
   │     └─ Email: "Commission suspense released - payment will be processed"
   │
   ├─ 7b. IF OUTCOME = POLICY_FRAUDULENT
   │  ├─ FORFEIT COMMISSION
   │  │  ├─ UPDATE commission_suspense SET status = 'FORFEITED'
   │  │  ├─ UPDATE commission_transactions SET payment_status = 'FORFEITED'
   │  │  └─ agent.suspense_amount -= commission_amount
   │  │
   │  ├─ FLAG AGENT FOR REVIEW
   │  │  ├─ flag_agent_for_review(agent_code, 'Fraudulent policy detected', 'CRITICAL')
   │  │  └─ Escalate to compliance team
   │  │
   │  └─ NOTIFY AGENT
   │     └─ Email (CRITICAL): "Commission forfeited - compliance review initiated"
   │
   └─ END: Suspense Resolved
```

**Success Criteria**:
- Suspense created within 1 hour of investigation trigger
- Accurate tracking of suspense amounts
- Release/forfeiture processed within 1 hour of investigation closure
- Agent notified at all critical stages
- Audit trail maintained

**SLA**:
- Suspense creation: 1 hour from trigger
- Investigation: Maximum 30 days (escalate at 20 days)
- Release/forfeiture: 1 hour from investigation closure

**Branch Workflows**:
- **Payment Failure Retry Workflow** (WF-IC-PAYMENT-RETRY-001): Handles failed EFT payments with 3 retry attempts
- **Overpayment Recovery Workflow** (BR-IC-SUSPENSE-004): Detects and recovers overpaid commissions

---

**End of Section 6: Workflows**

**Previous Section**: [Error Codes](#5-error-codes)
**Next Section**: [Temporal Workflows](#9-temporal-workflows)
# Incentive, Commission and Producer Management - Data Entities

## 7. Data Entities

### 7.1 Agent Profile Entity

**Table Name**: `agent_profiles`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **agent_id** | VARCHAR(20) | PK, NOT NULL | Unique Agent Code |
| **agent_type** | ENUM | NOT NULL | ADVISOR, ADVISOR_COORDINATOR, DEPARTMENTAL_EMPLOYEE, FIELD_OFFICER |
| **person_type** | ENUM | NOT NULL | INDIVIDUAL, CORPORATE_GROUP |
| **advisor_coordinator_id** | VARCHAR(20) | FK, NULLABLE | Link to Advisor Coordinator (NULL for Coordinators) |
| **profile_type** | VARCHAR(50) | NULLABLE | Profile classification |
| **office_type** | VARCHAR(50) | NULLABLE | Office type classification |
| **office_code** | VARCHAR(20) | NULLABLE | Affiliated office code |
| **advisor_sub_type** | VARCHAR(50) | NULLABLE | Advisor sub-category |
| **effective_date** | DATE | NOT NULL | Profile effective date |
| **distribution_channel** | VARCHAR(50) | NULLABLE | Distribution channel (e.g., India Post) |
| **product_class** | VARCHAR(50) | NULLABLE | PLI, RPLI |
| **title** | VARCHAR(10) | NULLABLE | Mr, Ms, Mrs, etc. |
| **first_name** | VARCHAR(100) | NOT NULL | First name |
| **middle_name** | VARCHAR(100) | NULLABLE | Middle name |
| **last_name** | VARCHAR(100) | NOT NULL | Last name |
| **gender** | ENUM | NOT NULL | MALE, FEMALE, OTHER |
| **date_of_birth** | DATE | NOT NULL | Date of birth (must be 18+) |
| **category** | VARCHAR(50) | NULLABLE | Category classification |
| **marital_status** | ENUM | NULLABLE | SINGLE, MARRIED, DIVORCED, WIDOWED |
| **aadhar_number** | VARCHAR(12) | NULLABLE, ENCRYPTED | Aadhar number |
| **pan** | VARCHAR(10) | UNIQUE, NOT NULL | PAN number |
| **designation** | VARCHAR(100) | NULLABLE | Designation/Rank |
| **service_number** | VARCHAR(50) | NULLABLE | Service number (for dept. employees) |
| **professional_title** | VARCHAR(100) | NULLABLE | Professional title |
| **status** | ENUM | NOT NULL | ACTIVE, SUSPENDED, TERMINATED, EXPIRED, DEACTIVATED |
| **status_reason** | VARCHAR(255) | NULLABLE | Reason for current status |
| **status_date** | DATE | NULLABLE | Date status changed |
| **effective_date** | DATE | NULLABLE | Status effective date |
| **termination_date** | DATE | NULLABLE | Termination date (if applicable) |
| **advisor_undergoing_training** | BOOLEAN | DEFAULT FALSE | Training flag |
| **preferred_payment_mode** | ENUM | NULLABLE | CHEQUE, EFT |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| **updated_at** | TIMESTAMP | ON UPDATE NOW() | Last update timestamp |
| **created_by** | VARCHAR(50) | NULLABLE | Created by user |
| **updated_by** | VARCHAR(50) | NULLABLE | Last updated by user |

**Indexes**:
- `idx_agent_pan` ON `pan`
- `idx_agent_status` ON `status`
- `idx_agent_coordinator` ON `advisor_coordinator_id`
- `idx_agent_name` ON `first_name`, `last_name`
- `idx_agent_type` ON `agent_type`

---

### 7.2 Agent Address Entity

**Table Name**: `agent_addresses`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **address_id** | BIGINT | PK, AUTO_INCREMENT | Unique Address ID |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **address_type** | ENUM | NOT NULL | OFFICIAL, PERMANENT, COMMUNICATION |
| **address_line1** | VARCHAR(255) | NOT NULL | Address line 1 |
| **address_line2** | VARCHAR(255) | NULLABLE | Address line 2 |
| **village** | VARCHAR(100) | NULLABLE | Village |
| **taluka** | VARCHAR(100) | NULLABLE | Taluka |
| **city** | VARCHAR(100) | NOT NULL | City |
| **district** | VARCHAR(100) | NULLABLE | District |
| **state** | VARCHAR(100) | NOT NULL | State |
| **country** | VARCHAR(100) | DEFAULT 'India' | Country |
| **pin_code** | VARCHAR(10) | NOT NULL | PIN code |
| **is_active** | BOOLEAN | DEFAULT TRUE | Active flag |
| **effective_from** | DATE | NOT NULL | Address effective from |
| **effective_to** | DATE | NULLABLE | Address effective to (NULL for current) |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_address_agent` ON `agent_id`
- `idx_address_type` ON `address_type`
- `idx_address_active` ON `is_active`

---

### 7.3 Agent Contact Entity

**Table Name**: `agent_contacts`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **contact_id** | BIGINT | PK, AUTO_INCREMENT | Unique Contact ID |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **contact_type** | ENUM | NOT NULL | OFFICIAL_LANDLINE, RESIDENTIAL_LANDLINE, MOBILE |
| **contact_value** | VARCHAR(15) | NOT NULL | Phone number |
| **is_primary** | BOOLEAN | DEFAULT FALSE | Primary contact flag |
| **is_active** | BOOLEAN | DEFAULT TRUE | Active flag |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_contact_agent` ON `agent_id`
- `idx_contact_type` ON `contact_type`
- `idx_contact_value` ON `contact_value`

---

### 7.4 Agent Email Entity

**Table Name**: `agent_emails`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **email_id** | BIGINT | PK, AUTO_INCREMENT | Unique Email ID |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **email_type** | ENUM | NOT NULL | OFFICIAL, PERMANENT, COMMUNICATION |
| **email_address** | VARCHAR(255) | NOT NULL | Email address |
| **is_primary** | BOOLEAN | DEFAULT FALSE | Primary email flag |
| **is_active** | BOOLEAN | DEFAULT TRUE | Active flag |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_email_agent` ON `agent_id`
- `idx_email_type` ON `email_type`
- `idx_email_address` ON `email_address`

---

### 7.5 Agent Bank Account Entity

**Table Name**: `agent_bank_accounts`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **bank_account_id** | BIGINT | PK, AUTO_INCREMENT | Unique Bank Account ID |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **bank_name** | VARCHAR(255) | NOT NULL | Bank name |
| **account_number** | VARCHAR(30) | NOT NULL, ENCRYPTED | Bank account number |
| **ifsc_code** | VARCHAR(11) | NOT NULL | IFSC code |
| **account_type** | VARCHAR(50) | NULLABLE | Savings, Current |
| **is_primary** | BOOLEAN | DEFAULT FALSE | Primary account for disbursement |
| **is_active** | BOOLEAN | DEFAULT TRUE | Active flag |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_bank_agent` ON `agent_id`
- `idx_bank_primary` ON `is_primary`

---

### 7.6 Agent License Entity

**Table Name**: `agent_licenses`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **license_id** | BIGINT | PK, AUTO_INCREMENT | Unique License ID |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **license_line** | VARCHAR(50) | NOT NULL | License line (e.g., Life) |
| **license_type** | VARCHAR(50) | NOT NULL | License type |
| **license_number** | VARCHAR(50) | NOT NULL | License number |
| **resident_status** | ENUM | NOT NULL | RESIDENT, NON_RESIDENT |
| **license_date** | DATE | NOT NULL | License issue date |
| **renewal_date** | DATE | NOT NULL | License renewal date |
| **authority_date** | DATE | NULLABLE | Authority date |
| **license_status** | ENUM | NOT NULL | ACTIVE, EXPIRED, RENEWED, CANCELLED |
| **is_active** | BOOLEAN | DEFAULT TRUE | Active flag |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| **updated_at** | TIMESTAMP | ON UPDATE NOW() | Last update timestamp |

**Indexes**:
- `idx_license_agent` ON `agent_id`
- `idx_license_number` ON `license_number`
- `idx_license_renewal_date` ON `renewal_date` (for reminder batch)
- `idx_license_status` ON `license_status`

---

### 7.7 Commission Rate Table Entity

**Table Name**: `commission_rates`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **rate_id** | BIGINT | PK, AUTO_INCREMENT | Unique Rate ID |
| **rate_percentage** | DECIMAL(5,2) | NOT NULL | Commission rate % |
| **policy_duration_months** | INT | NOT NULL | Policy duration in months |
| **product_type** | ENUM | NOT NULL | PLI, RPLI |
| **product_plan_code** | VARCHAR(50) | NOT NULL | Product plan code |
| **agent_type** | VARCHAR(50) | NOT NULL | Agent type |
| **policy_term_years** | INT | NOT NULL | Policy term in years |
| **effective_from** | DATE | NOT NULL | Rate effective from |
| **effective_to** | DATE | NULLABLE | Rate effective to (NULL for current) |
| **is_active** | BOOLEAN | DEFAULT TRUE | Active flag |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_rate_lookup` ON `product_type`, `product_plan_code`, `agent_type`, `policy_term_years`, `policy_duration_months`
- `idx_rate_active` ON `is_active`

---

### 7.8 Commission Transaction Entity

**Table Name**: `commission_transactions`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **commission_id** | BIGINT | PK, AUTO_INCREMENT | Unique Commission ID |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **policy_number** | VARCHAR(50) | FK, NOT NULL | Link to Policy |
| **commission_type** | ENUM | NOT NULL | FIRST_YEAR, RENEWAL, BONUS |
| **product_type** | ENUM | NOT NULL | PLI, RPLI |
| **annualised_premium** | DECIMAL(15,2) | NOT NULL | Annualised premium amount |
| **rate_percentage** | DECIMAL(5,2) | NOT NULL | Commission rate applied |
| **gross_commission** | DECIMAL(15,2) | NOT NULL | Gross commission amount |
| **tds_rate** | DECIMAL(5,2) | DEFAULT 0 | TDS rate % |
| **tds_amount** | DECIMAL(15,2) | DEFAULT 0 | TDS amount deducted |
| **net_commission** | DECIMAL(15,2) | NOT NULL | Net commission (gross - TDS) |
| **commission_date** | DATE | NOT NULL | Commission calculation date |
| **commission_status** | ENUM | NOT NULL | CALCULATED, TRIAL_PENDING, TRIAL_APPROVED, FINALIZED, READY_FOR_DISBURSEMENT, DISBURSED, CANCELLED |
| **trial_statement_id** | BIGINT | FK, NULLABLE | Link to Trial Statement |
| **final_statement_id** | BIGINT | FK, NULLABLE | Link to Final Statement |
| **disbursement_id** | BIGINT | FK, NULLABLE | Link to Disbursement |
| **batch_id** | VARCHAR(50) | NULLABLE | Batch job ID |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| **updated_at** | TIMESTAMP | ON UPDATE NOW() | Last update timestamp |

**Indexes**:
- `idx_commission_agent` ON `agent_id`
- `idx_commission_policy` ON `policy_number`
- `idx_commission_status` ON `commission_status`
- `idx_commission_date` ON `commission_date`
- `idx_commission_batch` ON `batch_id`
- `idx_commission_trial` ON `trial_statement_id`

---

### 7.9 Trial Statement Entity

**Table Name**: `trial_statements`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **trial_statement_id** | BIGINT | PK, AUTO_INCREMENT | Unique Trial Statement ID |
| **statement_number** | VARCHAR(50) | UNIQUE, NOT NULL | Statement number |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **statement_date** | DATE | NOT NULL | Statement generation date |
| **from_date** | DATE | NOT NULL | Period start date |
| **to_date** | DATE | NOT NULL | Period end date |
| **total_policies** | INT | DEFAULT 0 | Total policies in statement |
| **total_gross_commission** | DECIMAL(15,2) | DEFAULT 0 | Total gross commission |
| **total_tds** | DECIMAL(15,2) | DEFAULT 0 | Total TDS |
| **total_net_commission** | DECIMAL(15,2) | DEFAULT 0 | Total net commission |
| **statement_status** | ENUM | NOT NULL | PENDING_APPROVAL, APPROVED, REJECTED |
| **approved_by** | VARCHAR(50) | NULLABLE | Approved by user |
| **approved_at** | TIMESTAMP | NULLABLE | Approval timestamp |
| **approval_remarks** | TEXT | NULLABLE | Approval/rejection remarks |
| **processing_unit** | VARCHAR(50) | NULLABLE | Processing unit |
| **batch_id** | VARCHAR(50) | NULLABLE | Batch job ID |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| **updated_at** | TIMESTAMP | ON UPDATE NOW() | Last update timestamp |

**Indexes**:
- `idx_trial_agent` ON `agent_id`
- `idx_trial_status` ON `statement_status`
- `idx_trial_date` ON `statement_date`
- `idx_trial_number` ON `statement_number`

---

### 7.10 Final Statement Entity

**Table Name**: `final_statements`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **final_statement_id** | BIGINT | PK, AUTO_INCREMENT | Unique Final Statement ID |
| **statement_number** | VARCHAR(50) | UNIQUE, NOT NULL | Statement number |
| **trial_statement_id** | BIGINT | FK, NOT NULL | Link to Trial Statement |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **statement_date** | DATE | NOT NULL | Statement generation date |
| **total_gross_commission** | DECIMAL(15,2) | DEFAULT 0 | Total gross commission |
| **total_tds** | DECIMAL(15,2) | DEFAULT 0 | Total TDS |
| **total_net_commission** | DECIMAL(15,2) | DEFAULT 0 | Total net commission |
| **statement_status** | ENUM | NOT NULL | FINALIZED, READY_FOR_DISBURSEMENT, DISBURSEMENT_QUEUED, DISBURSED |
| **pdf_path** | VARCHAR(255) | NULLABLE | Path to generated PDF |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| **updated_at** | TIMESTAMP | ON UPDATE NOW() | Last update timestamp |

**Indexes**:
- `idx_final_agent` ON `agent_id`
- `idx_final_trial` ON `trial_statement_id`
- `idx_final_status` ON `statement_status`
- `idx_final_number` ON `statement_number`

---

### 7.11 Disbursement Entity

**Table Name**: `disbursements`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **disbursement_id** | VARCHAR(50) | PK, NOT NULL | Unique Disbursement ID (UUID) |
| **final_statement_id** | BIGINT | FK, NOT NULL | Link to Final Statement |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **payment_mode** | ENUM | NOT NULL | CHEQUE, EFT |
| **disbursement_amount** | DECIMAL(15,2) | NOT NULL | Disbursement amount |
| **disbursement_status** | ENUM | NOT NULL | PENDING, PROCESSING, COMPLETED, FAILED |
| **cheque_number** | VARCHAR(50) | NULLABLE | Cheque number (if cheque) |
| **bank_name** | VARCHAR(255) | NULLABLE | Bank name |
| **ifsc_code** | VARCHAR(11) | NULLABLE | IFSC code (if EFT) |
| **account_number** | VARCHAR(30) | NULLABLE, ENCRYPTED | Account number (if EFT) |
| **payment_date** | DATE | NULLABLE | Payment date |
| **payment_reference** | VARCHAR(100) | NULLABLE | Payment reference number |
| **pfms_transaction_id** | VARCHAR(100) | NULLABLE | PFMS transaction ID |
| **trial_approval_date** | DATE | NULLABLE | For SLA calculation |
| **sla_breach** | BOOLEAN | DEFAULT FALSE | SLA breached flag |
| **penalty_amount** | DECIMAL(15,2) | DEFAULT 0 | Penalty amount if breached |
| **remarks** | TEXT | NULLABLE | Disbursement remarks |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |
| **updated_at** | TIMESTAMP | ON UPDATE NOW() | Last update timestamp |
| **completed_at** | TIMESTAMP | NULLABLE | Completion timestamp |

**Indexes**:
- `idx_disb_agent` ON `agent_id`
- `idx_disb_final` ON `final_statement_id`
- `idx_disb_status` ON `disbursement_status`
- `idx_disb_date` ON `trial_approval_date` (for SLA monitoring)

---

### 7.12 Commission History Entity

**Table Name**: `commission_history`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **history_id** | BIGINT | PK, AUTO_INCREMENT | Unique History ID |
| **commission_id** | BIGINT | FK, NOT NULL | Link to Commission |
| **agent_id** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **policy_number** | VARCHAR(50) | FK, NOT NULL | Link to Policy |
| **commission_type** | ENUM | NOT NULL | FIRST_YEAR, RENEWAL, BONUS |
| **amount** | DECIMAL(15,2) | NOT NULL | Commission amount |
| **tds_amount** | DECIMAL(15,2) | DEFAULT 0 | TDS deducted |
| **net_amount** | DECIMAL(15,2) | NOT NULL | Net amount |
| **product_type** | ENUM | NOT NULL | PLI, RPLI |
| **transaction_date** | DATE | NOT NULL | Transaction date |
| **processed_date** | DATE | NOT NULL | Processing date |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_hist_agent` ON `agent_id`
- `idx_hist_policy` ON `policy_number`
- `idx_hist_date` ON `transaction_date`
- `idx_hist_type` ON `commission_type`

---

### 7.13 Commission Clawback Entity

**Table Name**: `commission_clawbacks`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **clawback_id** | BIGINT | PK, AUTO_INCREMENT | Unique Clawback ID |
| **policy_number** | VARCHAR(50) | FK, NOT NULL | Link to Policy |
| **agent_code** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **policy_issue_date** | DATE | NOT NULL | Policy issue date |
| **lapse_date** | DATE | NOT NULL | Policy lapse date |
| **months_active** | INT | NOT NULL | Months policy was active |
| **original_commission** | DECIMAL(15,2) | NOT NULL | First-year commission amount |
| **clawback_percentage** | DECIMAL(5,2) | NOT NULL | Clawback percentage (100, 75, 50, 25) |
| **clawback_amount** | DECIMAL(15,2) | NOT NULL | Amount to claw back |
| **recovery_status** | ENUM | NOT NULL, DEFAULT 'PENDING' | PENDING, PARTIALLY_RECOVERED, FULLY_RECOVERED, WRITTEN_OFF |
| **recovered_amount** | DECIMAL(15,2) | DEFAULT 0 | Amount recovered so far |
| **created_date** | DATE | NOT NULL | Clawback creation date |
| **clawback_reason** | VARCHAR(255) | | Reason for clawback |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_clawback_agent` ON `agent_code`
- `idx_clawback_policy` ON `policy_number`
- `idx_clawback_status` ON `recovery_status`
- `idx_clawback_date` ON `created_date`

**Foreign Keys**:
- `agent_code` REFERENCES `agent_profiles(agent_id)`
- `policy_number` REFERENCES `policies(policy_number)`

---

### 7.14 Commission Suspense Entity

**Table Name**: `commission_suspense`

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **suspense_id** | BIGINT | PK, AUTO_INCREMENT | Unique Suspense ID |
| **policy_number** | VARCHAR(50) | FK, NOT NULL | Link to Policy |
| **agent_code** | VARCHAR(20) | FK, NOT NULL | Link to Agent |
| **commission_amount** | DECIMAL(15,2) | NOT NULL | Commission amount held |
| **suspense_reason** | VARCHAR(100) | NOT NULL | Reason for suspense |
| **investigation_type** | VARCHAR(50) | | Type of investigation |
| **suspense_date** | DATE | NOT NULL | Suspense creation date |
| **expected_resolution_date** | DATE | | Expected resolution date |
| **status** | ENUM | NOT NULL, DEFAULT 'SUSPENDED' | SUSPENDED, RELEASED, FORFEITED |
| **release_date** | DATE | | Actual release date |
| **release_reason** | TEXT | | Reason for release |
| **forfeiture_date** | DATE | | Forfeiture date |
| **forfeiture_reason** | TEXT | | Reason for forfeiture |
| **investigation_reference** | VARCHAR(100) | | Investigation reference number |
| **created_at** | TIMESTAMP | DEFAULT NOW() | Record creation timestamp |

**Indexes**:
- `idx_suspense_agent` ON `agent_code`
- `idx_suspense_policy` ON `policy_number`
- `idx_suspense_status` ON `status`
- `idx_suspense_date` ON `suspense_date`

**Foreign Keys**:
- `agent_code` REFERENCES `agent_profiles(agent_id)`
- `policy_number` REFERENCES `policies(policy_number)`

---

### 7.15 Agent Profile Entity Updates

**Additional Fields Required for agent_profiles Table**:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| **clawback_pending_amount** | DECIMAL(15,2) | DEFAULT 0 | Total clawback pending recovery |
| **suspense_amount** | DECIMAL(15,2) | DEFAULT 0 | Total suspense amount held |
| **commission_status** | ENUM | DEFAULT 'ACTIVE' | ACTIVE, SUSPENDED, SUSPENDED_PENDING_CLAWBACK, SUSPENDED_PENDING_INVESTIGATION |

**ALTER Statements**:
```sql
ALTER TABLE agent_profiles
ADD COLUMN clawback_pending_amount DECIMAL(15,2) DEFAULT 0,
ADD COLUMN suspense_amount DECIMAL(15,2) DEFAULT 0,
ADD COLUMN commission_status ENUM('ACTIVE', 'SUSPENDED', 'SUSPENDED_PENDING_CLAWBACK', 'SUSPENDED_PENDING_INVESTIGATION') DEFAULT 'ACTIVE';

-- Create indexes for new fields
CREATE INDEX idx_agent_clawback ON agent_profiles(clawback_pending_amount);
CREATE INDEX idx_agent_suspense ON agent_profiles(suspense_amount);
CREATE INDEX idx_agent_comm_status ON agent_profiles(commission_status);
```

---

### Entity Relationships (ER Diagram)

```
agent_profiles (1) ----< (N) agent_addresses
agent_profiles (1) ----< (N) agent_contacts
agent_profiles (1) ----< (N) agent_emails
agent_profiles (1) ----< (N) agent_bank_accounts
agent_profiles (1) ----< (N) agent_licenses

agent_profiles (1) ----< (N) commission_transactions
agent_profiles (1) ----< (N) commission_clawbacks
agent_profiles (1) ----< (N) commission_suspense

commission_transactions (N) ----> (1) trial_statements
trial_statements (1) ----< (1) final_statements
final_statements (1) ----< (1) disbursements

commission_transactions (N) ----> (1) agent_profiles (agent_id)
commission_transactions (N) ----> (1) policies (policy_number)

commission_clawbacks (N) ----> (1) agent_profiles (agent_code)
commission_clawbacks (N) ----> (1) policies (policy_number)

commission_suspense (N) ----> (1) agent_profiles (agent_code)
commission_suspense (N) ----> (1) policies (policy_number)

commission_rates (lookup table) - referenced by commission_transactions
```

**New Relationships Added**:
- `agent_profiles` to `commission_clawbacks` (1:N) - One agent can have multiple clawbacks
- `agent_profiles` to `commission_suspense` (1:N) - One agent can have multiple suspense entries
- `commission_clawbacks` to `policies` (N:1) - Clawbacks linked to policies
- `commission_suspense` to `policies` (N:1) - Suspense entries linked to policies

---

**End of Section 7: Data Entities**

**Previous Section**: [Temporal Workflows](#9-temporal-workflows)
**Next Section**: [Integration Points](#8-integration-points)
# Incentive, Commission and Producer Management - Integration Points

## 8. Integration Points

### 8.1 HRMS Integration

**Integration**: HRMS System ↔ Agent Management Module

**Purpose**: Auto-populate Departmental Employee profiles

**Description**:
- Fetch employee details using Employee ID
- Retrieve: Name, DOB, Designation, Office, Contact details

**API Contract**:

```
GET /api/hrms/employee/{employee_id}

Request:
- employee_id: String (mandatory)

Response (200 OK):
{
  "status": "success",
  "employee": {
    "employee_id": "EMP12345",
    "first_name": "Rajesh",
    "last_name": "Kumar",
    "date_of_birth": "1985-05-15",
    "designation": "Postal Assistant",
    "office_code": "OFF001",
    "office_name": "New Delhi GPO",
    "mobile_number": "9876543210",
    "email_address": "rajesh.kumar@postal.gov.in"
  }
}

Response (404 Not Found):
{
  "status": "error",
  "message": "Employee ID not found"
}

Response (500 Internal Server Error):
{
  "status": "error",
  "message": "HRMS system unavailable"
}
```

**Error Handling**:
- 404: Employee ID not found → Display error IC-ERR-011
- 500: HRMS system unavailable → Show maintenance message, allow manual entry

**SLA**:
- Response time: < 2 seconds
- Availability: 99.5%

**Retry Strategy**:
- 3 attempts with exponential backoff (1s, 2s, 4s)
- Fallback to manual entry after failed retries

---

### 8.2 Policy Services Integration

**Integration**: Policy Management ↔ Commission Processing

**Purpose**: Fetch policy data for commission calculation

**Description**:
- Fetch active policies for commission period
- Retrieve policy details: Premium amount, mode, product type, plan, term

**API Contract**:

```
GET /api/policies/commission-eligible

Query Parameters:
- from_date: Date (mandatory) - Format: YYYY-MM-DD
- to_date: Date (mandatory) - Format: YYYY-MM-DD
- product_type: String (optional) - PLI, RPLI, ALL

Request Example:
GET /api/policies/commission-eligible?from_date=2026-01-01&to_date=2026-01-31&product_type=ALL

Response (200 OK):
{
  "status": "success",
  "policies": [
    {
      "policy_number": "PLI123456789",
      "agent_id": "AGT001",
      "product_type": "PLI",
      "plan_code": "ENDOWMENT",
      "policy_term_years": 15,
      "premium_amount": 1000.00,
      "premium_mode": "MONTHLY",
      "sum_assured": 200000.00,
      "issue_date": "2026-01-01",
      "policy_status": "ACTIVE"
    }
  ],
  "total_count": 1500,
  "page": 1,
  "page_size": 100
}

Response (400 Bad Request):
{
  "status": "error",
  "message": "Invalid date range"
}

Response (500 Internal Server Error):
{
  "status": "error",
  "message": "Policy service temporarily unavailable"
}
```

**Batch Processing**:
- Pagination: 1000 records per page
- Streaming support for large datasets
- Timeout: 5 seconds per 1000 records

**Error Handling**:
- 400: Invalid date range
- 500: Policy service unavailable → Retry with exponential backoff

**SLA**:
- Batch response time: < 5 seconds per 1000 policies
- Availability: 99.9%

**Retry Strategy**:
- Automatic retry up to 3 times
- Alert operations team if all retries fail

---

### 8.3 PFMS/Bank Integration

**Integration**: Disbursement Module ↔ PFMS/Bank Gateway

**Purpose**: Process EFT payments for commission disbursement

**Description**:
- Upload payment file for bulk disbursement
- Receive payment confirmation callback
- Update disbursement status

**API Contract (Outbound - Upload Payment File)**:

```
POST /api/pfms/payment-file

Request Headers:
- Content-Type: application/json
- Authorization: Bearer {access_token}
- X-Request-ID: {unique_request_id}

Request Body:
{
  "file_id": "PAY20260115001",
  "payment_date": "2026-01-15",
  "total_amount": 500000.00,
  "total_records": 50,
  "bank_code": "PUNB0112345",
  "debit_account": "1234567890",
  "payments": [
    {
      "reference_id": "DISB001",
      "beneficiary_name": "Rajesh Kumar",
      "account_number": "9876543210",
      "ifsc_code": "PUNB0112345",
      "amount": 50000.00,
      "remarks": "Commission Jan 2026"
    }
  ]
}

Response (200 OK - Accepted):
{
  "status": "accepted",
  "file_id": "PAY20260115001",
  "acknowledgement_number": "PFMS20260115001",
  "processing_status": "QUEUED",
  "estimated_processing_time": "2-3 working days"
}

Response (400 Bad Request):
{
  "status": "error",
  "message": "Invalid file format",
  "errors": [
    {
      "field": "payments[0].account_number",
      "error": "Invalid account number length"
    }
  ]
}

Response (401 Unauthorized):
{
  "status": "error",
  "message": "Invalid authentication credentials"
}

Response (503 Service Unavailable):
{
  "status": "error",
  "message": "PFMS gateway temporarily unavailable"
}
```

**API Contract (Inbound - Payment Confirmation Callback)**:

```
POST /api/disbursement/callback

Request (PFMS → Our System):
Headers:
- X-PFMS-Signature: {hmac_signature}
- X-Webhook-ID: {webhook_id}

Request Body:
{
  "acknowledgement_number": "PFMS20260115001",
  "file_id": "PAY20260115001",
  "processing_status": "COMPLETED",
  "processed_at": "2026-01-15T14:30:00Z",
  "transactions": [
    {
      "reference_id": "DISB001",
      "status": "SUCCESS",
      "transaction_id": "TXN123456789",
      "amount": 50000.00,
      "credited_at": "2026-01-15T14:35:00Z",
      "utr": "UTR123456789012"
    },
    {
      "reference_id": "DISB002",
      "status": "FAILED",
      "error_code": "ACCOUNT_CLOSED",
      "error_message": "Beneficiary account is closed"
    }
  ]
}

Response (200 OK):
{
  "status": "acknowledged",
  "processed_at": "2026-01-15T14:30:05Z"
}

Response (400 Bad Request):
{
  "status": "error",
  "message": "Invalid callback format"
}
```

**Security**:
- HMAC signature verification
- IP whitelisting for PFMS servers
- Retry handler for duplicate callbacks

**Error Handling**:
- 400: Invalid file format → Validate and retry
- 401: Authentication failure → Refresh token
- 503: PFMS service unavailable → Queue for retry (every 1 hour)

**SLA**:
- File upload response: < 10 seconds
- Payment confirmation: 2-3 working days
- Availability: 99%

**Retry Strategy**:
- Failed uploads: Retry every 1 hour for 24 hours
- Failed payments: Manual intervention after callback

---

### 8.4 Accounting Integration

**Integration**: Commission Disbursement ↔ Accounting Module

**Purpose**: Post commission disbursement accounting entries

**Description**:
- Post debit/credit entries for commission payments
- Create accounting voucher for each disbursement
- Track commission expense and liability

**API Contract**:

```
POST /api/accounting/voucher

Request Headers:
- Content-Type: application/json
- Authorization: Bearer {access_token}

Request Body:
{
  "voucher_type": "PAYMENT",
  "voucher_date": "2026-01-15",
  "reference_number": "DISB001",
  "reference_type": "COMMISSION_DISBURSEMENT",
  "amount": 50000.00,
  "narration": "Commission payment to Agent AGT001 for Jan 2026",
  "debit_account": {
    "code": "COMMISSION_EXPENSE",
    "name": "Commission Expense Account"
  },
  "credit_account": {
    "code": "BANK_ACCOUNT",
    "name": "PNB Operating Account"
  },
  "cost_center": "CIRCLE_DELHI",
  "agent_id": "AGT001",
  "disbursement_id": "DISB001",
  "batch_id": "BATCH_JAN2026"
}

Response (201 Created):
{
  "status": "success",
  "voucher_number": "JV20260115001",
  "voucher_date": "2026-01-15",
  "posting_status": "POSTED",
  "gl postings": [
    {
      "account": "COMMISSION_EXPENSE",
      "debit": 50000.00,
      "credit": 0
    },
    {
      "account": "BANK_ACCOUNT",
      "debit": 0,
      "credit": 50000.00
    }
  ],
  "created_at": "2026-01-15T10:30:00Z"
}

Response (400 Bad Request):
{
  "status": "error",
  "message": "Invalid voucher data",
  "errors": [
    {
      "field": "debit_account",
      "error": "Account code does not exist"
    }
  ]
}

Response (500 Internal Server Error):
{
  "status": "error",
  "message": "Accounting service error"
}
```

**Batch Vouchers**:
For bulk disbursements, multiple vouchers can be created in a single batch:

```
POST /api/accounting/vouchers/batch

Request Body:
{
  "batch_date": "2026-01-15",
  "vouchers": [...],  // Array of voucher objects
  "total_amount": 2500000.00,
  "total_count": 50
}

Response (201 Created):
{
  "status": "success",
  "batch_number": "JV_BATCH_20260115001",
  "voucher_numbers": ["JV20260115001", "JV20260115002", ...],
  "posting_status": "POSTED"
}
```

**Error Handling**:
- 400: Invalid voucher data → Validate and correct
- 500: Accounting service unavailable → Queue for retry

**SLA**:
- Response time: < 3 seconds
- Posting: Immediate (synchronous)
- Availability: 99.5%

**Retry Strategy**:
- Failed posts: Retry up to 3 times
- After 3 failures: Escalate to finance team

---

### 8.5 Notification Service Integration

**Integration**: All Modules → Notification Service

**Purpose**: Send email and SMS notifications

**API Contract**:

```
POST /api/notifications/send

Request Body:
{
  "notification_type": "EMAIL",
  "recipients": [
    {
      "type": "TO",
      "email": "agent@example.com",
      "name": "Rajesh Kumar"
    },
    {
      "type": "CC",
      "email": "coordinator@example.com"
    }
  ],
  "template_id": "COMMISSION_STATEMENT_READY",
  "template_data": {
    "agent_name": "Rajesh Kumar",
    "statement_date": "2026-01-15",
    "total_amount": "50000.00"
  },
  "priority": "NORMAL",
  "scheduled_at": null
}

Response (200 OK):
{
  "status": "queued",
  "notification_id": "NOTIF123456",
  "estimated_delivery": "2026-01-15T10:31:00Z"
}
```

**Notification Types**:
- EMAIL: Email notifications
- SMS: SMS notifications
- PUSH: Push notifications (for mobile app)
- WEBHOOK: Webhook callbacks

---

### Integration Summary Table

| Integration | Direction | Protocol | Authentication | Purpose |
|-------------|-----------|----------|----------------|----------|
| HRMS System | Outbound | REST API | API Key + OAuth | Fetch employee data |
| Policy Services | Outbound | REST API | Mutual TLS | Get policies for commission |
| PFMS/Bank | Outbound | REST API | Mutual TLS + HMAC | Submit EFT payments |
| PFMS/Bank | Inbound | Webhook | HMAC Signature | Payment confirmation |
| Accounting | Outbound | REST API | OAuth 2.0 | Post accounting entries |
| Notification | Outbound | REST API | API Key | Send emails/SMS |

---

### Integration Monitoring

**Health Check Endpoints**:
```
GET /api/health/integrations

Response:
{
  "status": "healthy",
  "integrations": {
    "hrms": {
      "status": "up",
      "last_check": "2026-01-15T10:00:00Z",
      "response_time_ms": 150
    },
    "policy_service": {
      "status": "up",
      "last_check": "2026-01-15T10:00:01Z",
      "response_time_ms": 80
    },
    "pfms": {
      "status": "degraded",
      "last_check": "2026-01-15T10:00:02Z",
      "response_time_ms": 5000
    },
    "accounting": {
      "status": "up",
      "last_check": "2026-01-15T10:00:03Z",
      "response_time_ms": 200
    }
  }
}
```

---

**End of Section 8: Integration Points**

**Previous Section**: [Data Entities](#7-data-entities)
**Next Section**: [Traceability Matrix](#10-traceability-matrix)
# Incentive, Commission and Producer Management - Temporal Workflows

## 9. Temporal Workflows

### 9.1 Monthly Commission Processing Workflow

**Workflow ID**: WF-TEMPORAL-IC-001
**Workflow Name**: MonthlyCommissionProcessing
**Description**: Orchestrates end-to-end commission calculation and disbursement

**Temporal Workflow (Go)**:

```go
package workflow

import (
    "time"
    "go.temporal.io/sdk/workflow"
)

// MonthlyCommissionProcessing orchestrates the complete commission lifecycle
func MonthlyCommissionProcessing(ctx workflow.Context, input CommissionInput) (ProcessingResult, error) {
    // Activity options with retry and timeout
    ao := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: 6 * time.Hour,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumAttempts:    3,
        },
    })

    // Step 1: Calculate Commissions
    var calcResult CalculationResult
    err := workflow.ExecuteActivity(ao, CalculateMonthlyCommissions, input.Month, input.Year).Get(ctx, &calcResult)
    if err != nil {
        return ProcessingResult{}, err
    }

    // Step 2: Generate Trial Statements
    var trialResult TrialGenerationResult
    err = workflow.ExecuteActivity(ao, GenerateTrialStatements, calcResult).Get(ctx, &trialResult)
    if err != nil {
        return ProcessingResult{}, err
    }

    // Step 3: Wait for Manual Approval (with 7-day timeout)
    var approvalSignal ApprovalSignal
    approvalChannel := workflow.GetSignalChannel(ctx, "approveTrialStatement")

    // Timer for timeout (7 days)
    timeoutFuture := workflow.NewTimer(ctx, 7*24*time.Hour)
    approvalFuture := workflow.NewFuture(ctx)

    // Signal handler
    err = workflow.SetSignalHandler(ctx, "approveTrialStatement", func(signal ApprovalSignal) error {
        approvalSignal = signal
        approvalFuture.Set(ctx, nil)
        return nil
    })

    // Wait for approval or timeout
    selector := workflow.NewSelector(ctx)

    selector.AddReceive(approvalChannel, func(c workflow.ReceiveChannel, more bool) {
        var signal ApprovalSignal
        c.Receive(ctx, &signal)
        approvalSignal = signal
    })

    selector.AddFuture(timeoutFuture, func(f workflow.Future, _ bool) {
        // Timeout - escalate
        workflow.ExecuteActivity(ao, EscalateApprovalTimeout, trialResult.StatementID)
    })

    selector.Select(ctx)

    // Step 4: Generate Final Statements
    var finalResult FinalGenerationResult
    err = workflow.ExecuteActivity(ao, GenerateFinalStatements, approvalSignal.StatementID).Get(ctx, &finalResult)
    if err != nil {
        return ProcessingResult{}, err
    }

    // Step 5: Process Disbursements
    var disbResult DisbursementResult
    err = workflow.ExecuteActivity(ao, ProcessDisbursements, finalResult).Get(ctx, &disbResult)
    if err != nil {
        return ProcessingResult{}, err
    }

    // Step 6: Monitor Disbursement SLA (10 working days)
    err = workflow.ExecuteActivity(ao, MonitorDisbursementSLA, disbResult.DisbursementID, 10*24*time.Hour).Get(ctx, nil)
    if err != nil {
        return ProcessingResult{}, err
    }

    return ProcessingResult{
        CommissionBatchID: calcResult.BatchID,
        TrialStatementID:  trialResult.StatementID,
        FinalStatementID:  finalResult.StatementID,
        DisbursementID:    disbResult.DisbursementID,
        TotalAmount:       finalResult.TotalAmount,
        AgentCount:        calcResult.AgentCount,
        CompletedAt:       time.Now(),
    }, nil
}

// Query handler for progress
func (w *Workflow) GetCommissionProgress(ctx workflow.Context) (float64, error) {
    // Return current progress percentage
    return 0.0, nil
}
```

**Activities**:

```go
package activities

import (
    "context"
    "time"
)

// CalculateMonthlyCommissions - Activity 1
func (a *Activities) CalculateMonthlyCommissions(ctx context.Context, month string, year int) (CalculationResult, error) {
    // 1. Fetch all active policies for the previous month
    policies, err := a.policyService.GetEligiblePolicies(ctx, month, year)
    if err != nil {
        return CalculationResult{}, err
    }

    totalRecords := len(policies)
    processedRecords := 0
    failedRecords := []string{}

    // 2. Process each policy
    for _, policy := range policies {
        // Heartbeat every 50 records
        if processedRecords%50 == 0 {
            if err := a.Heartbeat(ctx, processedRecords, totalRecords); err != nil {
                return CalculationResult{}, err
            }
        }

        // Calculate commission
        commission, err := a.calculateCommissionForPolicy(ctx, policy)
        if err != nil {
            failedRecords = append(failedRecords, policy.PolicyNumber)
            continue
        }

        // Store commission record
        err = a.commissionRepo.Save(ctx, commission)
        if err != nil {
            failedRecords = append(failedRecords, policy.PolicyNumber)
            continue
        }

        processedRecords++
    }

    // 3. Handle failed records (retry logic)
    for _, policyNum := range failedRecords {
        // Retry up to 3 times
        for i := 0; i < 3; i++ {
            // Retry processing
            if success := a.retryCommissionCalculation(ctx, policyNum); success {
                break
            }
        }
    }

    return CalculationResult{
        BatchID:          generateBatchID(),
        TotalPolicies:    totalRecords,
        ProcessedRecords: processedRecords,
        FailedRecords:    len(failedRecords),
        AgentCount:       a.getUniqueAgentCount(ctx),
        CompletedAt:      time.Now(),
    }, nil
}

// GenerateTrialStatements - Activity 2
func (a *Activities) GenerateTrialStatements(ctx context.Context, calcResult CalculationResult) (TrialGenerationResult, error) {
    // 1. Group commissions by agent
    agentCommissions := a.groupByAgent(ctx, calcResult.BatchID)

    // 2. For each agent, calculate totals and apply TDS
    var statementIDs []string

    for agentID, commissions := range agentCommissions {
        // Calculate totals
        totalGross := a.sumGrossCommission(commissions)
        totalTDS := a.calculateTDS(ctx, agentID, totalGross)
        totalNet := totalGross - totalTDS

        // Generate trial statement
        statement := &TrialStatement{
            StatementID:      generateStatementID(),
            AgentID:          agentID,
            StatementDate:    time.Now(),
            TotalGross:       totalGross,
            TotalTDS:         totalTDS,
            TotalNet:         totalNet,
            CommissionCount:  len(commissions),
            Status:           "PENDING_APPROVAL",
        }

        err := a.trialStatementRepo.Save(ctx, statement)
        if err != nil {
            return TrialGenerationResult{}, err
        }

        statementIDs = append(statementIDs, statement.StatementID)
    }

    // 3. Notify finance team
    a.notificationService.SendFinanceNotification(ctx, Notification{
        Type:    "TRIAL_STATEMENT_READY",
        Count:   len(statementIDs),
        Subject: "Trial Statements Ready for Review",
    })

    return TrialGenerationResult{
        StatementIDs: statementIDs,
        StatementID:  statementIDs[0], // Primary statement ID
        AgentCount:   len(agentCommissions),
        TotalAmount:  a.sumAllAgentCommissions(ctx, statementIDs),
        CreatedAt:    time.Now(),
    }, nil
}

// EscalateApprovalTimeout - Activity 3
func (a *Activities) EscalateApprovalTimeout(ctx context.Context, statementID string) error {
    // Send escalation to Finance Head
    return a.notificationService.SendEscalation(ctx, Escalation{
        Level:       "HIGH",
        Reason:      "Trial statement approval timeout - 7 days exceeded",
        StatementID: statementID,
        Recipients:  []string{"Finance Head", "Operations Head"},
    })
}

// GenerateFinalStatements - Activity 4
func (a *Activities) GenerateFinalStatements(ctx context.Context, trialStatementID string) (FinalGenerationResult, error) {
    // 1. Lock trial data
    err := a.trialStatementRepo.Lock(ctx, trialStatementID)
    if err != nil {
        return FinalGenerationResult{}, err
    }

    // 2. Get trial statement
    trial, err := a.trialStatementRepo.Get(ctx, trialStatementID)
    if err != nil {
        return FinalGenerationResult{}, err
    }

    // 3. Generate final statement
    final := &FinalStatement{
        StatementID:        generateFinalStatementID(),
        TrialStatementID:   trialStatementID,
        AgentID:            trial.AgentID,
        StatementDate:      time.Now(),
        TotalGross:         trial.TotalGross,
        TotalTDS:           trial.TotalTDS,
        TotalNet:           trial.TotalNet,
        Status:             "READY_FOR_DISBURSEMENT",
    }

    err = a.finalStatementRepo.Save(ctx, final)
    if err != nil {
        return FinalGenerationResult{}, err
    }

    // 4. Generate PDF
    pdfPath, err := a.pdfGenerator.GenerateFinalStatementPDF(ctx, final)
    if err != nil {
        return FinalGenerationResult{}, err
    }

    return FinalGenerationResult{
        StatementID:    final.StatementID,
        PDFPath:        pdfPath,
        TotalAmount:    final.TotalNet,
        AgentID:        final.AgentID,
        CreatedAt:      time.Now(),
    }, nil
}

// ProcessDisbursements - Activity 5
func (a *Activities) ProcessDisbursements(ctx context.Context, finalResult FinalGenerationResult) (DisbursementResult, error) {
    // Implementation depends on payment mode (Cheque/EFT)
    // This is a simplified version

    final, err := a.finalStatementRepo.Get(ctx, finalResult.StatementID)
    if err != nil {
        return DisbursementResult{}, err
    }

    // Get agent's preferred payment mode
    agent, err := a.agentRepo.Get(ctx, final.AgentID)
    if err != nil {
        return DisbursementResult{}, err
    }

    var disbursementID string

    if agent.PreferredPaymentMode == "CHEQUE" {
        // Process cheque disbursement
        disbursementID, err = a.processChequeDisbursement(ctx, final, agent)
    } else {
        // Process EFT disbursement
        disbursementID, err = a.processEFTDisbursement(ctx, final, agent)
    }

    if err != nil {
        return DisbursementResult{}, err
    }

    return DisbursementResult{
        DisbursementID: disbursementID,
        PaymentMode:     agent.PreferredPaymentMode,
        Amount:          final.TotalNet,
        AgentID:         final.AgentID,
    }, nil
}

// MonitorDisbursementSLA - Activity 6
func (a *Activities) MonitorDisbursementSLA(ctx context.Context, disbursementID string, slaDuration time.Duration) error {
    // Checkpoints for SLA monitoring
    checkpoints := []struct {
        offset  time.Duration
        urgency string
    }{
        {slaDuration - 3*24*time.Hour, "LOW"},     // T-3 days
        {slaDuration - 1*24*time.Hour, "MEDIUM"},  // T-1 day
        {slaDuration, "HIGH"},                      // T-0 (due date)
    }

    disbursement, err := a.disbursementRepo.Get(ctx, disbursementID)
    if err != nil {
        return err
    }

    for _, cp := range checkpoints {
        if time.Since(disbursement.CreatedAt) < cp.offset {
            continue
        }

        // Check if still pending
        if disbursement.Status != "COMPLETED" {
            // Send alert
            a.notificationService.SendSLAAlert(ctx, SLAAlert{
                DisbursementID: disbursementID,
                Urgency:        cp.urgency,
                DueDate:        disbursement.CreatedAt.Add(slaDuration),
                OverdueBy:      time.Since(disbursement.CreatedAt.Add(slaDuration)),
            })
        } else {
            return nil // Completed on time
        }
    }

    return nil
}
```

**Signals**:
```go
type ApprovalSignal struct {
    StatementID       string
    ApprovedBy        string
    ApprovalRemarks   string
    DisbursementType  string // "FULL" or "PARTIAL"
    PartialPercentage float64 // for partial disbursement
}
```

**Queries**:
```go
// Query: Get current processing status
func (w *Workflow) GetStatus(ctx workflow.Context) (string, error) {
    return w.currentStatus, nil
}

// Query: Get progress percentage
func (w *Workflow) GetProgress(ctx workflow.Context) (float64, error) {
    return (float64(w.completedSteps) / float64(w.totalSteps)) * 100, nil
}
```

---

### 9.2 License Renewal Reminder Workflow

**Workflow ID**: WF-TEMPORAL-IC-002
**Workflow Name**: LicenseRenewalReminder
**Description**: Manages license renewal reminders and expiry monitoring

**Temporal Workflow (Go)**:

```go
package workflow

import (
    "time"
    "go.temporal.io/sdk/workflow"
)

// LicenseRenewalReminder manages automated reminders and expiry
func LicenseRenewalReminder(ctx workflow.Context, licenseID string) error {
    // Get license details
    var license License
    ao := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: 5 * time.Minute,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumAttempts:    3,
        },
    })

    err := workflow.ExecuteActivity(ao, GetLicenseDetails, licenseID).Get(ctx, &license)
    if err != nil {
        return err
    }

    renewalDate := license.RenewalDate

    // Calculate reminder dates
    reminders := []struct {
        name     string
        offset   time.Duration
        daysLeft int
    }{
        {"T-30 days", 30 * 24 * time.Hour, 30},
        {"T-15 days", 15 * 24 * time.Hour, 15},
        {"T-7 days", 7 * 24 * time.Hour, 7},
        {"T-0 days", 0, 0},
    }

    // Send reminders at scheduled intervals
    for _, reminder := range reminders {
        waitDuration := renewalDate.Sub(time.Now()) - reminder.offset
        if waitDuration > 0 {
            // Wait until reminder time
            _ = workflow.NewTimer(ctx, waitDuration).Get(ctx, nil)

            // Send reminder
            err := workflow.ExecuteActivity(ao, SendLicenseReminder, licenseID, reminder.daysLeft).Get(ctx, nil)
            if err != nil {
                // Log error but continue
                workflow.GetLogger(ctx).Error("Failed to send reminder", "licenseID", licenseID, "error", err)
            }
        }
    }

    // Wait for renewal submission or expiry
    renewalSignal := workflow.NewSignalChannel(ctx, "licenseRenewalSubmitted")

    // Check daily after expiry
    for {
        // Wait for renewal signal or 1 day
        selector := workflow.NewSelector(ctx)

        selector.AddReceive(renewalSignal, func(c workflow.ReceiveChannel, more bool) {
            var renewal RenewalRequest
            c.Receive(ctx, &renewal)

            // Process renewal
            err := workflow.ExecuteActivity(ao, ProcessLicenseRenewal, renewal).Get(ctx, nil)
            if err != nil {
                workflow.GetLogger(ctx).Error("Failed to process renewal", "error", err)
            } else {
                // Renewal successful - end workflow
                return
            }
        })

        selector.AddFuture(workflow.NewTimer(ctx, 24*time.Hour), func(f workflow.Future, _ bool) {
            // Daily check - has license expired?
            var currentStatus LicenseStatus
            workflow.ExecuteActivity(ao, CheckLicenseStatus, licenseID).Get(ctx, &currentStatus)

            if currentStatus != "RENEWED" && time.Now().After(renewalDate) {
                // Deactivate agent
                workflow.ExecuteActivity(ao, DeactivateAgentOnLicenseExpiry, licenseID).Get(ctx, nil)

                // Send notifications
                workflow.ExecuteActivity(ao, SendLicenseExpiryNotifications, licenseID).Get(ctx, nil)

                // End workflow - agent deactivated
                return
            }
            // Continue monitoring
        })

        selector.Select(ctx)
    }

    return nil
}
```

**Activities**:

```go
package activities

// GetLicenseDetails retrieves license information
func (a *Activities) GetLicenseDetails(ctx context.Context, licenseID string) (License, error) {
    return a.licenseRepo.GetByID(ctx, licenseID)
}

// SendLicenseReminder sends reminder notification
func (a *Activities) SendLicenseReminder(ctx context.Context, licenseID string, daysLeft int) error {
    license, err := a.licenseRepo.GetByID(ctx, licenseID)
    if err != nil {
        return err
    }

    agent, err := a.agentRepo.Get(ctx, license.AgentID)
    if err != nil {
        return err
    }

    // Compose reminder email
    email := Email{
        To:       agent.Email,
        Subject:  fmt.Sprintf("License Renewal Reminder - %d days remaining", daysLeft),
        Template: "license_renewal_reminder",
        Data: map[string]interface{}{
            "AgentName":     agent.FirstName + " " + agent.LastName,
            "DaysLeft":      daysLeft,
            "RenewalDate":   license.RenewalDate,
            "LicenseNumber": license.LicenseNumber,
        },
    }

    // Send to agent
    err = a.emailService.Send(ctx, email)
    if err != nil {
        return err
    }

    // If urgent (<= 7 days), also notify coordinator
    if daysLeft <= 7 {
        coordinator, _ := a.agentRepo.Get(ctx, agent.AdvisorCoordinatorID)
        if coordinator != nil {
            email.To = coordinator.Email
            email.Subject = fmt.Sprintf("URGENT: Agent License Expiring Soon - %d days", daysLeft)
            a.emailService.Send(ctx, email)
        }
    }

    return nil
}

// ProcessLicenseRenewal processes renewal submission
func (a *Activities) ProcessLicenseRenewal(ctx context.Context, renewal RenewalRequest) error {
    // Validate documents
    if !a.validateRenewalDocuments(ctx, renewal) {
        return errors.New("incomplete documents")
    }

    // Update license
    license, err := a.licenseRepo.GetByID(ctx, renewal.LicenseID)
    if err != nil {
        return err
    }

    license.RenewalDate = renewal.NewRenewalDate
    license.Status = "RENEWED"

    err = a.licenseRepo.Update(ctx, license)
    if err != nil {
        return err
    }

    // Send confirmation
    a.sendRenewalConfirmation(ctx, license)

    return nil
}

// CheckLicenseStatus checks if license has been renewed
func (a *Activities) CheckLicenseStatus(ctx context.Context, licenseID string) (LicenseStatus, error) {
    license, err := a.licenseRepo.GetByID(ctx, licenseID)
    if err != nil {
        return "", err
    }
    return license.Status, nil
}

// DeactivateAgentOnLicenseExpiry deactivates agent when license expires
func (a *Activities) DeactivateAgentOnLicenseExpiry(ctx context.Context, licenseID string) error {
    license, err := a.licenseRepo.GetByID(ctx, licenseID)
    if err != nil {
        return err
    }

    // Update license status
    license.Status = "EXPIRED"
    err = a.licenseRepo.Update(ctx, license)
    if err != nil {
        return err
    }

    // Deactivate agent
    agent, err := a.agentRepo.Get(ctx, license.AgentID)
    if err != nil {
        return err
    }

    agent.Status = "DEACTIVATED"
    agent.DeactivationReason = "License expired"
    err = a.agentRepo.Update(ctx, agent)
    if err != nil {
        return err
    }

    return nil
}

// SendLicenseExpiryNotifications sends notifications after expiry
func (a *Activities) SendLicenseExpiryNotifications(ctx context.Context, licenseID string) error {
    license, _ := a.licenseRepo.GetByID(ctx, licenseID)
    agent, _ := a.agentRepo.Get(ctx, license.AgentID)

    // Notify agent
    a.emailService.Send(ctx, Email{
        To:      agent.Email,
        Subject: "License Expired - Agent Code Deactivated",
        Body:    "Your insurance license has expired. Your agent code has been deactivated. Please renew your license to reactivate.",
    })

    // Notify coordinator
    if agent.AdvisorCoordinatorID != "" {
        coordinator, _ := a.agentRepo.Get(ctx, agent.AdvisorCoordinatorID)
        a.emailService.Send(ctx, Email{
            To:      coordinator.Email,
            Subject: "Agent Deactivated - License Expired",
            Body:    fmt.Sprintf("Agent %s has been deactivated due to license expiry.", agent.FirstName+" "+agent.LastName),
        })
    }

    return nil
}
```

---

### 9.3 Commission Disbursement SLA Monitoring Workflow

**Workflow ID**: WF-TEMPORAL-IC-003
**Workflow Name**: DisbursementSLAMonitor
**Description**: Monitors disbursement SLA and handles escalations/penalties

**Temporal Workflow (Go)**:

```go
package workflow

import (
    "time"
    "go.temporal.io/sdk/workflow"
)

// DisbursementSLAMonitor monitors 10-day disbursement SLA
func DisbursementSLAMonitor(ctx workflow.Context, disbursementID string) error {
    ao := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: 1 * time.Minute,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumAttempts:    3,
        },
    })

    // Get disbursement details
    var disb Disbursement
    err := workflow.ExecuteActivity(ao, GetDisbursementDetails, disbursementID).Get(ctx, &disb)
    if err != nil {
        return err
    }

    slaDuration := 10 * 24 * time.Hour // 10 working days
    dueDate := disb.TrialApprovalDate.Add(slaDuration)

    // Set up SLA checkpoints
    checkpoints := []struct {
        name     string
        offset   time.Duration
        urgency  string
    }{
        {"T-3 days", 3 * 24 * time.Hour, "LOW"},
        {"T-1 day", 1 * 24 * time.Hour, "MEDIUM"},
        {"T-0 hours (due date)", 0, "HIGH"},
    }

    // Monitor checkpoints
    for _, cp := range checkpoints {
        waitTime := dueDate.Sub(time.Now()) - cp.offset
        if waitTime > 0 {
            // Wait until checkpoint
            _ = workflow.NewTimer(ctx, waitTime).Get(ctx, nil)

            // Check status
            var status DisbursementStatus
            workflow.ExecuteActivity(ao, GetDisbursementStatus, disbursementID).Get(ctx, &status)

            if status != "COMPLETED" {
                // Still pending - send alert
                workflow.ExecuteActivity(ao, SendSLAAlert, disbursementID, cp.urgency).Get(ctx, nil)
            } else {
                return nil // Disbursement completed
            }
        }
    }

    // SLA Breached - handle breach
    workflow.ExecuteActivity(ao, HandleSLABreach, disbursementID).Get(ctx, nil)

    // Wait for completion (with 5-day grace period)
    completionSignal := workflow.GetSignalChannel(ctx, "disbursementCompleted")
    graceTimer := workflow.NewTimer(ctx, 5*24*time.Hour)

    selector := workflow.NewSelector(ctx)

    selector.AddReceive(completionSignal, func(c workflow.ReceiveChannel, more bool) {
        var signal DisbursementCompletedSignal
        c.Receive(ctx, &signal)

        // Calculate and apply penalty
        workflow.ExecuteActivity(ao, CalculateAndApplyPenalty, disbursementID, signal.CompletedAt).Get(ctx, nil)
    })

    selector.AddFuture(graceTimer, func(f workflow.Future, _ bool) {
        // Still not completed after grace period - critical escalation
        workflow.ExecuteActivity(ao, CriticalEscalation, disbursementID).Get(ctx, nil)
    })

    selector.Select(ctx)

    return nil
}
```

**Activities**:

```go
package activities

// GetDisbursementDetails retrieves disbursement information
func (a *Activities) GetDisbursementDetails(ctx context.Context, disbursementID string) (Disbursement, error) {
    return a.disbursementRepo.Get(ctx, disbursementID)
}

// GetDisbursementStatus gets current disbursement status
func (a *Activities) GetDisbursementStatus(ctx context.Context, disbursementID string) (DisbursementStatus, error) {
    disb, err := a.disbursementRepo.Get(ctx, disbursementID)
    if err != nil {
        return "", err
    }
    return disb.Status, nil
}

// SendSLAAlert sends SLA warning alert
func (a *Activities) SendSLAAlert(ctx context.Context, disbursementID string, urgency string) error {
    disb, _ := a.disbursementRepo.Get(ctx, disbursementID)
    final, _ := a.finalStatementRepo.Get(ctx, disb.FinalStatementID)
    agent, _ := a.agentRepo.Get(ctx, disb.AgentID)

    alert := Alert{
        Type:        "SLA_WARNING",
        Urgency:     urgency,
        DisbursementID: disbursementID,
        AgentID:     agent.AgentID,
        AgentName:   agent.FirstName + " " + agent.LastName,
        Amount:      final.TotalNet,
        DueDate:     disb.TrialApprovalDate.Add(10 * 24 * time.Hour),
        Urgency:     urgency,
    }

    // Send to Finance Head
    a.notificationService.SendAlert(ctx, "Finance Head", alert)
    a.notificationService.SendAlert(ctx, "Operations Head", alert)

    return nil
}

// HandleSLABreach handles SLA breach
func (a *Activities) HandleSLABreach(ctx context.Context, disbursementID string) error {
    disb, _ := a.disbursementRepo.Get(ctx, disbursementID)

    // Log breach
    a.slaRepo.LogBreach(ctx, SLABreach{
        DisbursementID: disbursementID,
        BreachedAt:     time.Now(),
        OverdueBy:      time.Since(disb.TrialApprovalDate.Add(10 * 24 * time.Hour)),
    })

    // Notify stakeholders
    a.notificationService.SendSLABreachNotification(ctx, disbursementID)

    return nil
}

// CalculateAndApplyPenalty calculates penalty interest (8% per annum)
func (a *Activities) CalculateAndApplyPenalty(ctx context.Context, disbursementID string, completedAt time.Time) error {
    disb, _ := a.disbursementRepo.Get(ctx, disbursementID)
    final, _ := a.finalStatementRepo.Get(ctx, disb.FinalStatementID)

    // Calculate days delayed
    dueDate := disb.TrialApprovalDate.Add(10 * 24 * time.Hour)
    daysDelayed := int(completedAt.Sub(dueDate).Hours() / 24)

    if daysDelayed <= 0 {
        return nil // No penalty
    }

    // Calculate penalty: Amount × 8% × (days / 365)
    penaltyAmount := final.TotalNet * 0.08 * float64(daysDelayed) / 365.0

    // Apply penalty as credit to agent
    a.agentRepo.AddCredit(ctx, disb.AgentID, Credit{
        Type:        "PENALTY_INTEREST",
        Amount:      penaltyAmount,
        Reference:   disbursementID,
        Description: fmt.Sprintf("SLA breach penalty - %d days delayed", daysDelayed),
    })

    // Notify agent
    agent, _ := a.agentRepo.Get(ctx, disb.AgentID)
    a.emailService.Send(ctx, Email{
        To:      agent.Email,
        Subject: "SLA Breach Penalty Applied",
        Body:    fmt.Sprintf("A penalty of ₹%.2f has been credited to your account due to %d-day delay in commission payment.", penaltyAmount, daysDelayed),
    })

    return nil
}

// CriticalEscalation escalates to Director level
func (a *Activities) CriticalEscalation(ctx context.Context, disbursementID string) error {
    // Create critical incident
    incident := Incident{
        Type:        "CRITICAL_SLA_BREACH",
        Severity:    "CRITICAL",
        Reference:   disbursementID,
        CreatedAt:   time.Now(),
    }

    a.incidentRepo.Create(ctx, incident)

    // Notify Director
    a.notificationService.SendCriticalAlert(ctx, "Director", incident)

    return nil
}
```

**Signals**:
```go
type DisbursementCompletedSignal struct {
    DisbursementID string
    CompletedAt    time.Time
    PaymentMode    string
    Reference      string
}
```

---

### 9.4 Agent Onboarding Orchestration Workflow

**Workflow ID**: WF-TEMPORAL-IC-004
**Workflow Name**: AgentOnboardingOrchestration
**Description**: Orchestrates agent onboarding with parallel processing and validation

**Temporal Workflow (Go)**:

```go
package workflow

import (
    "time"
    "go.temporal.io/sdk/workflow"
)

// AgentOnboardingOrchestration orchestrates complete onboarding process
func AgentOnboardingOrchestration(ctx workflow.Context, request OnboardingRequest) (string, error) {
    ao := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Minute,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumAttempts:    3,
        },
    })

    // Step 1: Validate mandatory fields
    err := workflow.ExecuteActivity(ao, ValidateOnboardingRequest, request).Get(ctx, nil)
    if err != nil {
        return "", err
    }

    // Step 2: Fetch data from external systems (parallel)
    var futures []workflow.Future

    // For Departmental Employee - fetch HRMS data
    if request.AgentType == "DEPARTMENTAL_EMPLOYEE" {
        f := workflow.ExecuteActivity(ao, FetchHRMSData, request.EmployeeID)
        futures = append(futures, f)
    }

    // PAN Validation (parallel)
    f := workflow.ExecuteActivity(ao, ValidatePANUniqueness, request.PAN)
    futures = append(futures, f)

    // Coordinator Validation (for Advisors)
    if request.AgentType == "ADVISOR" {
        f := workflow.ExecuteActivity(ao, ValidateCoordinatorExists, request.CoordinatorID)
        futures = append(futures, f)
    }

    // Wait for all parallel activities
    for _, f := range futures {
        err := f.Get(ctx, nil)
        if err != nil {
            return "", err
        }
    }

    // Step 3: Generate Agent Code
    var agentCode string
    err = workflow.ExecuteActivity(ao, GenerateAgentCode, request).Get(ctx, &agentCode)
    if err != nil {
        return "", err
    }

    // Step 4: Create Profile (with retries)
    err = workflow.ExecuteActivity(ao, CreateAgentProfile, agentCode, request).Get(ctx, nil)
    if err != nil {
        return "", err
    }

    // Step 5: Create child records (parallel)
    for _, addr := range request.Addresses {
        workflow.ExecuteActivity(ao, CreateAgentAddress, agentCode, addr)
    }

    for _, contact := range request.Contacts {
        workflow.ExecuteActivity(ao, CreateAgentContact, agentCode, contact)
    }

    // Step 6: Send Notifications
    workflow.ExecuteActivity(ao, SendOnboardingNotifications, agentCode, request).Get(ctx, nil)

    return agentCode, nil
}
```

**Activities**:

```go
package activities

// ValidateOnboardingRequest validates all mandatory fields
func (a *Activities) ValidateOnboardingRequest(ctx context.Context, request OnboardingRequest) error {
    // Validate agent type
    if request.AgentType == "" {
        return errors.New("agent type is required")
    }

    // Validate mandatory fields based on agent type
    if request.AgentType == "DEPARTMENTAL_EMPLOYEE" && request.EmployeeID == "" {
        return errors.New("employee ID required for departmental employee")
    }

    if request.AgentType == "ADVISOR" && request.CoordinatorID == "" {
        return errors.New("coordinator ID required for advisor")
    }

    // Validate PAN format
    if !isValidPAN(request.PAN) {
        return errors.New("invalid PAN format")
    }

    // Validate names
    if request.FirstName == "" || request.LastName == "" {
        return errors.New("first name and last name are required")
    }

    // Validate DOB (must be 18+)
    age := time.Since(request.DateOfBirth).Hours() / 24 / 365
    if age < 18 {
        return errors.New("agent must be at least 18 years old")
    }

    return nil
}

// FetchHRMSData fetches employee data from HRMS
func (a *Activities) FetchHRMSData(ctx context.Context, employeeID string) (HRMSEmployeeData, error) {
    return a.hrmsClient.GetEmployee(ctx, employeeID)
}

// ValidatePANUniqueness checks if PAN already exists
func (a *Activities) ValidatePANUniqueness(ctx context.Context, pan string) (bool, error) {
    exists, err := a.agentRepo.PANExists(ctx, pan)
    if err != nil {
        return false, err
    }
    return !exists, nil // Valid if PAN doesn't exist
}

// ValidateCoordinatorExists validates coordinator is active
func (a *Activities) ValidateCoordinatorExists(ctx context.Context, coordinatorID string) (bool, error) {
    coordinator, err := a.agentRepo.Get(ctx, coordinatorID)
    if err != nil {
        return false, errors.New("coordinator not found")
    }

    if coordinator.Status != "ACTIVE" {
        return false, errors.New("coordinator is not active")
    }

    return true, nil
}

// GenerateAgentCode generates unique agent code
func (a *Activities) GenerateAgentCode(ctx context.Context, request OnboardingRequest) (string, error) {
    // Format: AGT{CircleCode}{SequenceNumber}
    circleCode := request.CircleCode
    sequence, err := a.agentRepo.GetNextSequence(ctx, circleCode)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("AGT%s%05d", circleCode, sequence), nil
}

// CreateAgentProfile creates agent profile
func (a *Activities) CreateAgentProfile(ctx context.Context, agentCode string, request OnboardingRequest) error {
    agent := &Agent{
        AgentID:              agentCode,
        AgentType:            request.AgentType,
        FirstName:            request.FirstName,
        LastName:             request.LastName,
        PAN:                  request.PAN,
        DateOfBirth:          request.DateOfBirth,
        AdvisorCoordinatorID: request.CoordinatorID,
        CircleID:             request.CircleID,
        DivisionID:           request.DivisionID,
        Status:               "ACTIVE",
        CreatedAt:            time.Now(),
    }

    return a.agentRepo.Create(ctx, agent)
}

// CreateAgentAddress creates address record
func (a *Activities) CreateAgentAddress(ctx context.Context, agentCode string, address Address) error {
    address.AgentID = agentCode
    return a.addressRepo.Create(ctx, address)
}

// CreateAgentContact creates contact record
func (a *Activities) CreateAgentContact(ctx context.Context, agentCode string, contact Contact) error {
    contact.AgentID = agentCode
    return a.contactRepo.Create(ctx, contact)
}

// SendOnboardingNotifications sends welcome notifications
func (a *Activities) SendOnboardingNotifications(ctx context.Context, agentCode string, request OnboardingRequest) error {
    agent, _ := a.agentRepo.Get(ctx, agentCode)

    // Send welcome email to agent
    a.emailService.Send(ctx, Email{
        To:      agent.Email,
        Subject: "Welcome to Postal Life Insurance",
        Template: "agent_welcome",
        Data: map[string]interface{}{
            "AgentName": agent.FirstName,
            "AgentCode": agentCode,
        },
    })

    // Notify coordinator
    if agent.AdvisorCoordinatorID != "" {
        coordinator, _ := a.agentRepo.Get(ctx, agent.AdvisorCoordinatorID)
        a.emailService.Send(ctx, Email{
            To:      coordinator.Email,
            Subject: "New Agent Onboarded",
            Body:    fmt.Sprintf("New agent %s (%s) has been onboarded under you.", agent.FirstName+" "+agent.LastName, agentCode),
        })
    }

    return nil
}
```

---


---

### 9.5 Commission Clawback Workflow

**Workflow ID**: `WF-TEMPORAL-IC-005`
**Name**: `CommissionClawbackWorkflow`
**Description**: Temporal workflow for processing commission clawbacks when policies lapse within 24 months

**Trigger**: Policy lapse detected (daily batch)

**Go Implementation**:

```go
package workflow

import (
    "time"
    "go.temporal.io/sdk/workflow"
)

// CommissionClawbackInput represents input to the clawback workflow
type CommissionClawbackInput struct {
    PolicyNumber     string
    AgentCode        string
    PolicyIssueDate  time.Time
    LapseDate        time.Time
}

// CommissionClawbackWorkflow processes clawbacks for lapsed policies
func CommissionClawbackWorkflow(ctx workflow.Context, input CommissionClawbackInput) error {
    // Step 1: Calculate months active
    var monthsActive int
    err := workflow.ExecuteActivity(ctx, a.CalculateMonthsActiveActivity, input.PolicyIssueDate, input.LapseDate).Get(ctx, &monthsActive)
    if err != nil {
        return err
    }

    // Step 2: Check eligibility (clawback only if < 24 months)
    if monthsActive >= 24 {
        // No clawback required - exit workflow
        return nil
    }

    // Step 3: Calculate clawback percentage
    var clawbackPercentage float64
    err = workflow.ExecuteActivity(ctx, a.CalculateClawbackPercentageActivity, monthsActive).Get(ctx, &clawbackPercentage)
    if err != nil {
        return err
    }

    // Step 4: Get original first-year commission
    var originalCommission float64
    err = workflow.ExecuteActivity(ctx, a.GetFirstYearCommissionActivity, input.PolicyNumber).Get(ctx, &originalCommission)
    if err != nil {
        return err
    }

    // Step 5: Calculate clawback amount
    clawbackAmount := originalCommission * (clawbackPercentage / 100.0)

    // Step 6: Create clawback entry
    var clawbackID int64
    err = workflow.ExecuteActivity(ctx, a.CreateClawbackEntryActivity, CreateClawbackEntryInput{
        PolicyNumber:        input.PolicyNumber,
        AgentCode:           input.AgentCode,
        PolicyIssueDate:     input.PolicyIssueDate,
        LapseDate:           input.LapseDate,
        MonthsActive:        monthsActive,
        OriginalCommission:  originalCommission,
        ClawbackPercentage:  clawbackPercentage,
        ClawbackAmount:      clawbackAmount,
    }).Get(ctx, &clawbackID)
    if err != nil {
        return err
    }

    // Step 7: Update agent clawback balance
    err = workflow.ExecuteActivity(ctx, a.UpdateAgentClawbackBalanceActivity, UpdateAgentBalanceInput{
        AgentCode:        input.AgentCode,
        ClawbackAmount:   clawbackAmount,
        CommissionStatus: "SUSPENDED_PENDING_CLAWBACK",
    }).Get(ctx, nil)
    if err != nil {
        return err
    }

    // Step 8: Post accounting entry
    err = workflow.ExecuteActivity(ctx, a.PostAccountingEntryActivity, AccountingEntryInput{
        DebitAccount:  "Commission Expense Reversal",
        CreditAccount: "Agent Payable - " + input.AgentCode,
        Amount:        clawbackAmount,
        Reference:     "Clawback - Policy " + input.PolicyNumber,
        EffectiveDate: time.Now(),
    }).Get(ctx, nil)
    if err != nil {
        return err
    }

    // Step 9: Send notifications (execute in parallel)
    var agentNotificationFuture workflow.Future
    var financeNotificationFuture workflow.Future

    // Notify agent
    agentNotificationFuture = workflow.ExecuteActivity(ctx, a.SendEmailActivity, SendEmailInput{
        To:      getAgentEmail(input.AgentCode),
        Subject: "Commission Clawback Notice",
        Body:    formatClawbackEmail(input.PolicyNumber, clawbackAmount, clawbackPercentage, monthsActive),
        Priority: "HIGH",
    })

    // Notify finance team
    financeNotificationFuture = workflow.ExecuteActivity(ctx, a.SendFinanceAlertActivity, SendFinanceAlertInput{
        Message:  formatFinanceClawbackAlert(input.PolicyNumber, input.AgentCode, clawbackAmount),
        Priority: "HIGH",
    })

    // Wait for both notifications
    if err := agentNotificationFuture.Get(ctx, nil); err != nil {
        return err
    }
    if err := financeNotificationFuture.Get(ctx, nil); err != nil {
        return err
    }

    // Step 10: Mark policy as processed
    err = workflow.ExecuteActivity(ctx, a.MarkPolicyClawbackProcessedActivity, MarkPolicyProcessedInput{
        PolicyNumber:     input.PolicyNumber,
        ClawbackAmount:   clawbackAmount,
        ClawbackProcessed: true,
    }).Get(ctx, nil)
    if err != nil {
        return err
    }

    // Step 11: Schedule recovery from future commissions
    // This is a long-running process - we'll listen for commission payment events
    childWorkflowOptions := workflow.ChildWorkflowOptions{
        WorkflowExecutionTimeout: 365 * 24 * time.Hour,
    }

    var recoveryResult interface{}
    err = workflow.ExecuteChildWorkflow(ctx, RecoverClawbackFromCommissionsWorkflow, childWorkflowOptions, RecoverClawbackInput{
        AgentCode:       input.AgentCode,
        ClawbackAmount:  clawbackAmount,
        ClawbackID:      clawbackID,
    }).Get(ctx, &recoveryResult)

    return err
}

// CalculateClawbackPercentageActivity calculates clawback percentage based on months active
func (a *Activities) CalculateClawbackPercentageActivity(ctx context.Context, monthsActive int) (float64, error) {
    if monthsActive < 6 {
        return 100.0, nil  // Full clawback
    } else if monthsActive < 12 {
        return 75.0, nil   // 75% clawback
    } else if monthsActive < 18 {
        return 50.0, nil   // 50% clawback
    } else if monthsActive < 24 {
        return 25.0, nil   // 25% clawback
    } else {
        return 0.0, nil    // No clawback
    }
}

// RecoverClawbackFromCommissionsWorkflow recovers clawback from future commission payments
func RecoverClawbackFromCommissionsWorkflow(ctx workflow.Context, input RecoverClawbackInput) error {
    remainingAmount := input.ClawbackAmount

    // Listen for commission payment events for this agent
    for remainingAmount > 0 {
        // Wait for next commission payment signal
        var signal CommissionPaymentSignal
        signalChannel := workflow.GetSignalChannel(ctx, "commission-payment-"+input.AgentCode)

        var more bool
        selector := workflow.NewSelector(ctx)

        selector.AddReceive(signalChannel, func(c workflow.ReceiveChannel, more bool) {
            c.Receive(ctx, &signal)
            more = more
        })

        selector.Select(ctx)

        if !more {
            break // Channel closed
        }

        // Calculate recovery (max 50% of current commission)
        maxRecovery := signal.CommissionAmount * 0.50
        actualRecovery := remainingAmount

        if maxRecovery < remainingAmount {
            actualRecovery = maxRecovery
        }

        // Execute recovery
        var newRemainingAmount float64
        err := workflow.ExecuteActivity(ctx, a.DeductFromCommissionActivity, DeductCommissionInput{
            AgentCode:       input.AgentCode,
            ClawbackID:      input.ClawbackID,
            CommissionID:    signal.CommissionID,
            RecoveryAmount:  actualRecovery,
        }).Get(ctx, &newRemainingAmount)

        if err != nil {
            return err
        }

        remainingAmount = newRemainingAmount

        // Update clawback status
        if remainingAmount == 0 {
            workflow.ExecuteActivity(ctx, a.UpdateClawbackStatusActivity, input.ClawbackID, "FULLY_RECOVERED")
        } else {
            workflow.ExecuteActivity(ctx, a.UpdateClawbackStatusActivity, input.ClawbackID, "PARTIALLY_RECOVERED")
        }

        // Notify agent of recovery
        workflow.ExecuteActivity(ctx, a.SendEmailActivity, SendEmailInput{
            To:      getAgentEmail(input.AgentCode),
            Subject: "Commission Recovery Deduction",
            Body:    formatRecoveryEmail(actualRecovery, remainingAmount),
            Priority: "MEDIUM",
        })
    }

    return nil
}
```

---

### 9.6 Suspense Management Workflow

**Workflow ID**: `WF-TEMPORAL-IC-006`
**Name**: `SuspenseManagementWorkflow`
**Description**: Temporal workflow for managing commission suspense accounts during policy investigations

**Trigger**: Policy marked as UNDER_INVESTIGATION

**Go Implementation**:

```go
package workflow

import (
    "time"
    "go.temporal.io/sdk/workflow"
)

// SuspenseManagementInput represents input to the suspense workflow
type SuspenseManagementInput struct {
    PolicyNumber           string
    AgentCode              string
    InvestigationReason    string
    InvestigationReference string
    ExpectedDuration       int // days
}

// InvestigationOutcome represents the result of an investigation
type InvestigationOutcome struct {
    Status      string // "POLICY_GENUINE", "POLICY_FRAUDULENT", "INCONCLUSIVE"
    Reason      string
    DeterminedBy string
    DeterminedAt time.Time
}

// SuspenseManagementWorkflow manages commission suspense during investigations
func SuspenseManagementWorkflow(ctx workflow.Context, input SuspenseManagementInput) error {
    // Step 1: Check commission payment status
    var commissionPaid bool
    var commissionAmount float64

    err := workflow.ExecuteActivity(ctx, a.CheckCommissionPaymentStatusActivity, input.PolicyNumber).Get(ctx, &struct {
        Paid   bool
        Amount float64
    }{commissionPaid, commissionAmount})

    if err != nil {
        return err
    }

    // Step 2: Create suspense or hold payment
    var suspenseID int64

    if commissionPaid {
        // Commission already paid - create suspense entry
        err = workflow.ExecuteActivity(ctx, a.CreateSuspenseEntryActivity, CreateSuspenseInput{
            PolicyNumber:           input.PolicyNumber,
            AgentCode:              input.AgentCode,
            CommissionAmount:       commissionAmount,
            SuspenseReason:         "POLICY_UNDER_INVESTIGATION",
            InvestigationType:      input.InvestigationReason,
            InvestigationReference: input.InvestigationReference,
            ExpectedResolutionDate: time.Now().AddDate(0, 0, input.ExpectedDuration),
            Status:                 "SUSPENDED",
        }).Get(ctx, &suspenseID)

        if err != nil {
            return err
        }

        // Update agent status
        err = workflow.ExecuteActivity(ctx, a.UpdateAgentCommissionStatusActivity, UpdateAgentStatusInput{
            AgentCode:        input.AgentCode,
            CommissionStatus: "SUSPENDED_PENDING_INVESTIGATION",
            SuspenseAmount:   commissionAmount,
        }).Get(ctx, nil)

        if err != nil {
            return err
        }

    } else {
        // Commission not yet paid - hold payment
        err = workflow.ExecuteActivity(ctx, a.HoldCommissionPaymentActivity, HoldCommissionInput{
            PolicyNumber: input.PolicyNumber,
            HoldReason:   "POLICY_UNDER_INVESTIGATION",
            HoldUntil:    time.Now().AddDate(0, 0, input.ExpectedDuration),
        }).Get(ctx, nil)

        if err != nil {
            return err
        }
    }

    // Step 3: Send notifications
    err = workflow.ExecuteActivity(ctx, a.SendEmailActivity, SendEmailInput{
        To:      getAgentEmail(input.AgentCode),
        Subject: "Commission Suspended - Investigation",
        Body:    formatSuspenseEmail(input.PolicyNumber, commissionAmount, input.InvestigationReason),
        Priority: "HIGH",
    }).Get(ctx, nil)

    if err != nil {
        return err
    }

    // Step 4: Wait for investigation outcome with SLA timer
    investigationFuture, _ := workflow.NewFuture(ctx)

    // Set up signal handler for investigation complete
    var outcome InvestigationOutcome
    signalChannel := workflow.GetSignalChannel(ctx, "investigation-complete-"+input.PolicyNumber)

    selector := workflow.NewSelector(ctx)

    // Signal handler
    selector.AddReceive(signalChannel, func(c workflow.ReceiveChannel, more bool) {
        c.Receive(ctx, &outcome)
        investigationFuture.Set(result, nil)
    })

    // SLA timer (escalate if exceeds expected duration)
    slaTimer := workflow.NewTimer(ctx, time.Duration(input.ExpectedDuration+5)*24*time.Hour)

    selector.AddFuture(slaTimer, func(f workflow.Future, _ bool) {
        // SLA breached - escalate
        workflow.ExecuteActivity(ctx, a.EscalateSuspenseSLAActivity, EscalationInput{
            SuspenseID:     suspenseID,
            PolicyNumber:   input.PolicyNumber,
            SLABreachedBy:  input.ExpectedDuration + 5,
        })
    })

    // Wait for either signal or timer
    selector.Select(ctx)

    // Step 5: Process investigation outcome
    if outcome.Status == "" {
        // No outcome yet (timer fired) - wait for outcome
        signalChannel.Receive(ctx, &outcome)
    }

    switch outcome.Status {
    case "POLICY_GENUINE":
        // Release suspense
        err = workflow.ExecuteActivity(ctx, a.ReleaseSuspenseActivity, ReleaseSuspenseInput{
            SuspenseID:    suspenseID,
            PolicyNumber:  input.PolicyNumber,
            ReleaseReason: "Investigation cleared - policy genuine",
            ReleaseDate:   time.Now(),
        }).Get(ctx, nil)

        if err != nil {
            return err
        }

        // Update agent status
        err = workflow.ExecuteActivity(ctx, a.UpdateAgentCommissionStatusActivity, UpdateAgentStatusInput{
            AgentCode:        input.AgentCode,
            CommissionStatus: "ACTIVE",
            SuspenseAmount:   -commissionAmount, // Deduct
        }).Get(ctx, nil)

        if err != nil {
            return err
        }

        // Notify agent
        workflow.ExecuteActivity(ctx, a.SendEmailActivity, SendEmailInput{
            To:      getAgentEmail(input.AgentCode),
            Subject: "Commission Suspense Released",
            Body:    formatSuspenseReleasedEmail(input.PolicyNumber, commissionAmount),
            Priority: "HIGH",
        })

    case "POLICY_FRAUDULENT":
        // Forfeit suspense
        err = workflow.ExecuteActivity(ctx, a.ForfeitSuspenseActivity, ForfeitSuspenseInput{
            SuspenseID:       suspenseID,
            PolicyNumber:     input.PolicyNumber,
            ForfeitureReason: "Policy determined fraudulent",
            ForfeitureDate:   time.Now(),
        }).Get(ctx, nil)

        if err != nil {
            return err
        }

        // Update agent status
        err = workflow.ExecuteActivity(ctx, a.UpdateAgentCommissionStatusActivity, UpdateAgentStatusInput{
            AgentCode:        input.AgentCode,
            CommissionStatus: "ACTIVE", // May change based on investigation
            SuspenseAmount:   -commissionAmount,
        }).Get(ctx, nil)

        if err != nil {
            return err
        }

        // Flag agent for review
        workflow.ExecuteActivity(ctx, a.FlagAgentForReviewActivity, FlagAgentInput{
            AgentCode: input.AgentCode,
            Reason:    "Fraudulent policy detected",
            Severity:  "CRITICAL",
        }).Get(ctx, nil)

        // Notify agent
        workflow.ExecuteActivity(ctx, a.SendEmailActivity, SendEmailInput{
            To:       getAgentEmail(input.AgentCode),
            Subject:  "Commission Forfeited",
            Body:     formatForfeitureEmail(input.PolicyNumber, commissionAmount, outcome.Reason),
            Priority: "CRITICAL",
        })

        // Escalate to compliance
        workflow.ExecuteActivity(ctx, a.EscalateToComplianceTeamActivity, ComplianceEscalationInput{
            AgentCode:     input.AgentCode,
            PolicyNumber:  input.PolicyNumber,
            SuspenseID:    suspenseID,
            ForfeitAmount: commissionAmount,
            FraudDetails:  outcome.Reason,
        })

    case "INCONCLUSIVE":
        // Extend investigation or partial release
        err = workflow.ExecuteActivity(ctx, a.ExtendSuspenseActivity, ExtendSuspenseInput{
            SuspenseID:            suspenseID,
            AdditionalDays:        30,
            ExtensionReason:       "Investigation inconclusive - extended for further review",
        }).Get(ctx, nil)

        if err != nil {
            return err
        }

        // Notify agent
        workflow.ExecuteActivity(ctx, a.SendEmailActivity, SendEmailInput{
            To:      getAgentEmail(input.AgentCode),
            Subject: "Investigation Extended",
            Body:    formatExtensionEmail(input.PolicyNumber, 30),
            Priority: "HIGH",
        })
    }

    return nil
}
```

---
**End of Section 9: Temporal Workflows**

**Previous Section**: [Workflows](#6-workflows)
**Next Section**: [Data Entities](#7-data-entities)
# Incentive, Commission and Producer Management - Traceability Matrix

## 10. Traceability Matrix

### 10.1 Business Rules to Functional Requirements

| Business Rule ID | Business Rule Name | Functional Requirements |
|------------------|-------------------|------------------------|
| **BR-IC-AH-001** | Advisor Linkage Requirement | FR-IC-ONB-001, FR-IC-ONB-003 |
| **BR-IC-AH-002** | Advisor Coordinator Assignment | FR-IC-ONB-001 |
| **BR-IC-AH-003** | Departmental Employee Auto-Population | FR-IC-ONB-004 |
| **BR-IC-AH-004** | Field Officer Onboarding Mode | FR-IC-ONB-005 |
| **BR-IC-LIC-001** | First License Renewal Period | FR-IC-PROF-004 |
| **BR-IC-LIC-002** | Subsequent License Renewal Period | FR-IC-PROF-004 |
| **BR-IC-LIC-003** | Auto-Deactivation on Expiry | FR-IC-PROF-004 |
| **BR-IC-LIC-004** | License Renewal Reminder Schedule | FR-IC-PROF-004 |
| **BR-IC-LIC-005** | License Renewal 3-Day SLA | FR-IC-PROF-004 |
| **BR-IC-COM-001** | Monthly Commission Calculation | FR-IC-COM-002 |
| **BR-IC-COM-002** | Trial Statement Before Disbursement | FR-IC-COM-003, FR-IC-COM-004, FR-IC-COM-006 |
| **BR-IC-COM-003** | TDS Deduction Requirement | FR-IC-COM-003, FR-IC-COM-007 |
| **BR-IC-COM-004** | Annualised Premium Calculation | FR-IC-COM-002 |
| **BR-IC-COM-005** | Partial Disbursement Option | FR-IC-COM-006 |
| **BR-IC-COM-006** | Commission Rate Table Structure | FR-IC-COM-001 |
| **BR-IC-COM-007** | Final Statement Generation Batch | FR-IC-COM-007 |
| **BR-IC-COM-008** | Disbursement Mode Workflow | FR-IC-COM-009, FR-IC-COM-010 |
| **BR-IC-COM-009** | Commission History Search | FR-IC-COM-011 |
| **BR-IC-COM-010** | Export Commission Statements | FR-IC-COM-012 |
| **BR-IC-COM-011** | Commission Disbursement 10-Day SLA | FR-IC-COM-010 |
| **BR-IC-COM-012** | Commission Batch 6-Hour Timeout | FR-IC-COM-002 |
| **BR-IC-CLAWBACK-001** | Commission Clawback on Policy Lapse | FR-IC-COM-008, FR-IC-COM-009, FR-IC-COM-010, FR-IC-COM-011, FR-IC-COM-012 |
| **BR-IC-CLAWBACK-002** | Commission Adjustment on Early Death Claims | FR-IC-COM-013 |
| **BR-IC-SUSPENSE-001** | Commission Suspense for Disputed Policies | FR-IC-COM-013, FR-IC-COM-014, FR-IC-COM-015 |
| **BR-IC-SUSPENSE-002** | Commission Payment Failure Suspense with Retry Logic | FR-IC-COM-016, FR-IC-COM-017 |
| **BR-IC-SUSPENSE-003** | Incentive Suspense for Persistency Shortfall | FR-IC-COM-018, FR-IC-COM-019 |
| **BR-IC-SUSPENSE-004** | Overpayment Recovery Suspense | FR-IC-COM-016 |
| **BR-IC-SUSPENSE-005** | Commission Hold for Investigation | FR-IC-COM-013, FR-IC-COM-014, FR-IC-COM-015 |
| **BR-IC-PROF-001** | Agent Search Functionality | FR-IC-PROF-001 |
| **BR-IC-PROF-002** | PAN Uniqueness Validation | FR-IC-ONB-002 |
| **BR-IC-PROF-003** | Agent Status Change with Reason | FR-IC-PROF-003 |
| **BR-IC-PROF-004** | Agent Termination Workflow | FR-IC-PROF-005 |

---

### 10.2 Functional Requirements to Validation Rules

| Functional Requirement ID | Functional Requirement Name | Validation Rules |
|---------------------------|----------------------------|------------------|
| **FR-IC-ONB-001** | New Profile Selection | VR-IC-PROF-003, VR-IC-PROF-004, VR-IC-PROF-005 |
| **FR-IC-ONB-002** | Profile Details Entry | VR-IC-PROF-001, VR-IC-PROF-002, VR-IC-PROF-003, VR-IC-PROF-004, VR-IC-PROF-005, VR-IC-PROF-006, VR-IC-PROF-007 |
| **FR-IC-COM-002** | Commission Calculation Batch | VR-IC-COM-001 |
| **FR-IC-COM-005** | Partial Disbursement | VR-IC-COM-003 |
| **FR-IC-COM-009** | Disbursement Details Entry | VR-IC-COM-004 |
| **FR-IC-COM-010** | Automatic Disbursement | VR-IC-COM-005 |
| **FR-IC-PROF-004** | License Management | VR-IC-LIC-001, VR-IC-LIC-002, VR-IC-LIC-003 |

---

### 10.3 Workflows to Temporal Workflows

| Workflow ID | Workflow Name | Temporal Workflow ID | Temporal Workflow Name |
|-------------|---------------|----------------------|------------------------|
| **WF-IC-COM-001** | Monthly Commission Processing | WF-TEMPORAL-IC-001 | MonthlyCommissionProcessing |
| **WF-IC-LIC-001** | License Renewal | WF-TEMPORAL-IC-002 | LicenseRenewalReminder |
| (SLA Monitoring) | Disbursement SLA Monitor | WF-TEMPORAL-IC-003 | DisbursementSLAMonitor |
| **WF-IC-ONB-001** | Agent Onboarding | WF-TEMPORAL-IC-004 | AgentOnboardingOrchestration |
| **WF-IC-CLAWBACK-001** | Commission Clawback Process | WF-TEMPORAL-IC-005 | CommissionClawbackWorkflow |
| **WF-IC-SUSPENSE-001** | Suspense Account Management | WF-TEMPORAL-IC-006 | SuspenseManagementWorkflow |

---

### 10.4 Error Codes to Business Rules

| Error Code | Error Message | Related Business Rules |
|------------|---------------|------------------------|
| **IC-ERR-001** | Please select a Profile Type | BR-IC-VAL-002 |
| **IC-ERR-002** | PAN already exists for another profile | BR-IC-PROF-002 |
| **IC-ERR-003** | Please enter a 10 digit PAN | BR-IC-VAL-001 |
| **IC-ERR-004** | Please enter correct PAN | BR-IC-VAL-001 |
| **IC-ERR-005** | Please enter a Last name | BR-IC-VAL-002 |
| **IC-ERR-006** | Please enter a First name | BR-IC-VAL-002 |
| **IC-ERR-007** | Please enter a valid Date of Birth | BR-IC-VAL-003 |
| **IC-ERR-008** | Your selected criteria did not return any rows | BR-IC-PROF-001 |
| **IC-ERR-009** | Coordinator ID is mandatory | BR-IC-AH-001 |
| **IC-ERR-010** | Circle assignment required | BR-IC-AH-002 |
| **IC-ERR-011** | Employee ID not found in HRMS | BR-IC-AH-003 |
| **IC-ERR-012** | Trial statement must be approved before disbursement | BR-IC-COM-002 |
| **IC-ERR-013** | Commission rate not found | BR-IC-COM-006 |
| **IC-ERR-014** | Disbursement amount exceeds commission | BR-IC-COM-005 |
| **IC-ERR-015** | Bank details required for EFT | BR-IC-COM-008 |
| **IC-ERR-016** | License expired. Agent code deactivated | BR-IC-LIC-003 |
| **IC-ERR-017** | Commission clawback initiated for policy {policy_number}. Amount: ₹{amount} | BR-IC-CLAWBACK-001 |
| **IC-ERR-018** | Commission held in suspense for policy {policy_number}. Reason: {reason} | BR-IC-SUSPENSE-001, BR-IC-SUSPENSE-005 |
| **IC-ERR-019** | Payment failed after 3 retry attempts. Finance team notified | BR-IC-SUSPENSE-002 |

---

### 10.5 SRS Requirements to Implementation Traceability

| SRS Requirement ID | SRS Requirement Title | Business Rule | Functional Requirement | Data Entity |
|-------------------|----------------------|---------------|------------------------|-------------|
| **FS_IC_001** | Agent linked to Coordinator | BR-IC-AH-001 | FR-IC-ONB-001, FR-IC-ONB-003 | agent_profiles |
| **FS_IC_002** | Coordinator Circle/Division assignment | BR-IC-AH-002 | FR-IC-ONB-001 | agent_profiles |
| **FS_IC_003** | Dept Employee HRMS auto-population | BR-IC-AH-003 | FR-IC-ONB-004 | agent_profiles |
| **FS_IC_004** | Field Officer onboarding modes | BR-IC-AH-004 | FR-IC-ONB-005 | agent_profiles |
| **FS_IC_005** | Agent search interface | BR-IC-PROF-001 | FR-IC-PROF-001 | agent_profiles |
| **FS_IC_006** | Agent profile dashboard | BR-IC-PROF-003 | FR-IC-PROF-002 | agent_profiles |
| **FS_IC_007** | Name update with validation | BR-IC-VAL-002 | FR-IC-PROF-002 | agent_profiles |
| **FS_IC_008** | PAN update with validation | BR-IC-PROF-002, BR-IC-VAL-001 | FR-IC-PROF-002 | agent_profiles |
| **FS_IC_009** | Status update with reason | BR-IC-PROF-003 | FR-IC-PROF-003 | agent_profiles |
| **FS_IC_010** | Personal information update | BR-IC-VAL-003 | FR-IC-PROF-002 | agent_profiles |
| **FS_IC_011** | Distribution channel details | - | FR-IC-PROF-002 | agent_profiles |
| **FS_IC_012** | External ID numbers | - | FR-IC-PROF-002 | agent_profiles |
| **FS_IC_013** | Product class information | - | FR-IC-PROF-002 | agent_profiles |
| **FS_IC_014** | Address management | - | FR-IC-ONB-002 | agent_addresses |
| **FS_IC_015** | Phone number management | VR-IC-PROF-006 | FR-IC-ONB-002 | agent_contacts |
| **FS_IC_016** | Email address management | VR-IC-PROF-007 | FR-IC-ONB-002 | agent_emails |
| **FS_IC_017** | Authority types | - | FR-IC-ONB-002 | agent_profiles |
| **FS_IC_018** | License renewal reminders | BR-IC-LIC-004 | FR-IC-PROF-004 | agent_licenses |
| **FS_IC_019** | License renewal periods | BR-IC-LIC-001, BR-IC-LIC-002 | FR-IC-PROF-004 | agent_licenses |
| **FS_IC_020** | Auto-deactivation on expiry | BR-IC-LIC-003 | FR-IC-PROF-004 | agent_licenses, agent_profiles |
| **FS_IC_021** | Agent termination | BR-IC-PROF-004 | FR-IC-PROF-005 | agent_profiles |
| **FS_IC_022** | Commission rate table | BR-IC-COM-006 | FR-IC-COM-001 | commission_rates |
| **FS_IC_023** | Commission history search | BR-IC-COM-009 | FR-IC-COM-011 | commission_history |
| **FS_IC_024** | Monthly commission calculation | BR-IC-COM-001, BR-IC-COM-012 | FR-IC-COM-002 | commission_transactions |
| **FS_IC_025** | Trial statement generation | BR-IC-COM-002 | FR-IC-COM-003 | trial_statements |
| **FS_IC_026** | Trial statement view | BR-IC-COM-009 | FR-IC-COM-004 | trial_statements |
| **FS_IC_027** | Manual trial statement | BR-IC-COM-002 | FR-IC-COM-005 | trial_statements |
| **FS_IC_028** | Trial statement approval | BR-IC-COM-002, BR-IC-COM-005 | FR-IC-COM-006 | trial_statements |
| **FS_IC_029** | Final statement generation | BR-IC-COM-007 | FR-IC-COM-007 | final_statements |
| **FS_IC_030** | Final statement view | BR-IC-COM-007 | FR-IC-COM-008 | final_statements |
| **FS_IC_031** | Disbursement details entry | BR-IC-COM-008 | FR-IC-COM-009 | disbursements |
| **FS_IC_032** | Automatic disbursement | BR-IC-COM-008, BR-IC-COM-011 | FR-IC-COM-010 | disbursements |

---

### 10.6 Test Case Coverage Matrix

| Test Case ID | Test Scenario | Covers Requirements | Covers BR | Covers VR |
|--------------|--------------|--------------------|-----------|-----------|
| **TC_001** | Create Advisor with Coordinator | FS_IC_001 | BR-IC-AH-001 | - |
| **TC_002** | Create Advisor without Coordinator | FS_IC_001 | BR-IC-AH-001 | - |
| **TC_003** | Create Coordinator with Circle | FS_IC_002 | BR-IC-AH-002 | - |
| **TC_004** | Create Coordinator without Circle | FS_IC_002 | BR-IC-AH-002 | - |
| **TC_005** | Dept Employee auto-population | FS_IC_003 | BR-IC-AH-003 | - |
| **TC_006** | Invalid Employee ID | FS_IC_003 | BR-IC-AH-003 | - |
| **TC_007** | Field Officer manual entry | FS_IC_004 | BR-IC-AH-004 | - |
| **TC_008** | Field Officer missing fields | FS_IC_004 | - | VR-IC-PROF-003 |
| **TC_009** | Agent search by ID | FS_IC_005 | BR-IC-PROF-001 | - |
| **TC_010** | Agent search invalid ID | FS_IC_005 | BR-IC-PROF-001 | - |
| **TC_011** | PAN update valid format | FS_IC_008 | BR-IC-PROF-002 | VR-IC-PROF-001 |
| **TC_012** | PAN update invalid format | FS_IC_008 | - | VR-IC-PROF-001 |
| **TC_013** | Status update with reason | FS_IC_009 | BR-IC-PROF-003 | - |
| **TC_014** | Status update without reason | FS_IC_009 | BR-IC-PROF-003 | - |
| **TC_015** | License reminder T-15 | FS_IC_018 | BR-IC-LIC-004 | - |
| **TC_016** | License expiry handling | FS_IC_020 | BR-IC-LIC-003 | - |
| **TC_017** | Commission rate setup | FS_IC_022 | BR-IC-COM-006 | - |
| **TC_018** | Commission rate missing product | FS_IC_022 | - | VR-IC-COM-001 |
| **TC_019** | Commission batch execution | FS_IC_024 | BR-IC-COM-001 | - |
| **TC_020** | Commission batch no policies | FS_IC_024 | - | - |
| **TC_021** | Approve trial statement | FS_IC_028 | BR-IC-COM-002 | - |
| **TC_022** | Approve without selection | FS_IC_028 | - | VR-IC-COM-005 |
| **TC_023** | EFT disbursement | FS_IC_032 | BR-IC-COM-008 | - |
| **TC_024** | Disbursement without bank details | FS_IC_032 | - | VR-IC-COM-004 |

---

### 10.7 SLA Compliance Matrix

| Process | SLA Requirement | Business Rule | Monitoring Workflow | Alert Escalation |
|---------|----------------|---------------|---------------------|------------------|
| Commission Disbursement | 10 working days | BR-IC-COM-011 | WF-TEMPORAL-IC-003 | Finance Head → Director |
| License Renewal Processing | 3 working days | BR-IC-LIC-005 | WF-TEMPORAL-IC-002 | Operations → Supervisor |
| Commission Batch Processing | 6 hours max | BR-IC-COM-012 | WF-TEMPORAL-IC-001 | 3hr: MEDIUM, 5hr: HIGH, 6hr: CRITICAL |
| License Renewal Reminders | T-30, T-15, T-7, T-0 | BR-IC-LIC-004 | WF-TEMPORAL-IC-002 | Automated notifications |

---

### 10.8 Database Entity Relationships

```
agent_profiles (1) ────< (N) agent_addresses
agent_profiles (1) ────< (N) agent_contacts
agent_profiles (1) ────< (N) agent_emails
agent_profiles (1) ────< (N) agent_bank_accounts
agent_profiles (1) ────< (N) agent_licenses

agent_profiles (1) ────< (N) commission_transactions
commission_transactions (N) ────> (1) trial_statements
trial_statements (1) ────< (1) final_statements
final_statements (1) ────< (1) disbursements

agent_profiles (1) ────< (N) commission_history
commission_transactions (N) ────> (1) commission_rates (lookup)
```

---

### 10.9 Security and Compliance Mapping

| Requirement | Security Control | Data Entities Affected | Compliance Standard |
|-------------|------------------|----------------------|---------------------|
| PAN Protection | Encryption at rest | agent_profiles | Tax Regulations |
| Bank Account Protection | Encryption at rest | agent_bank_accounts | Banking Regulations |
| TDS Compliance | Audit logging | commission_transactions, trial_statements, final_statements | Tax Laws |
| License Tracking | Immutable records | agent_licenses | IRDAI Regulations |
| Commission Audit Trail | Full audit logging | All commission tables | Financial Regulations |
| Data Privacy | Access controls | All personal data | Privacy Laws |
| Payment Security | HMAC signature verification | disbursements | PCI-DSS-like standards |

---

### 10.10 Feature Completeness Matrix

| Feature Module | Requirements Count | Functional Requirements | Workflows | Temporal Workflows | Status |
|---------------|-------------------|------------------------|-----------|-------------------|--------|
| Agent Onboarding | 5 | 5 | 1 | 1 | ✓ Complete |
| Agent Profile Management | 5 | 5 | 2 | 0 | ✓ Complete |
| License Management | 5 | 1 | 1 | 1 | ✓ Complete |
| Commission Configuration | 1 | 1 | 0 | 0 | ✓ Complete |
| Commission Processing | 12 | 12 | 1 | 2 | ✓ Complete |
| Commission Reporting | 2 | 2 | 1 | 0 | ✓ Complete |
| Disbursement | 2 | 2 | 1 | 1 | ✓ Complete |
| **TOTAL** | **32** | **32** | **8** | **4** | **✓ Complete** |

---
**End of Section 10: Traceability Matrix**

**Previous Section**: [Temporal Workflows](#9-temporal-workflows)
**Next Section**: [Commission Rate Structure](#11-commission-rate-structure)

---

# Incentive, Commission and Producer Management - Commission Rate Structure

## 11. Commission Rate Structure

### 11.1 PLI First Year Incentive Rates

The following commission rates apply to PLI policies based on policy type and term:

| Policy Type | Policy Term | Incentive Rate |
|-------------|-------------|----------------|
| Other than AEA | <= 15 years | 4% |
| Other than AEA | 15 < Term <= 25 years | 10% |
| Other than AEA | > 25 years | 20% |
| AEA | <= 15 years | 5% |
| AEA | > 15 years | 7% |

**Applicable From:** 01 April 2025

**Business Rule Reference:** BR-IC-COM-006

**Calculation Notes:**
- Rate is applied to First Year Premium
- Commission is calculated only after free-look period completion (15 days)
- No commission payable until policy is accepted and free-look completes

### 11.2 RPLI First Year Incentive

| Policy Type | Incentive Rate |
|-------------|----------------|
| RPLI (all policies) | 10% of first-year premium |

**Applicable From:** 01 April 2025

**Eligible Agents:**
- Departmental Employees (DE)
- Field Officers (FO)
- Direct Agents (DA)
- Gramin Dak Sevak (GDS)

### 11.3 Renewal Incentive Rates

Renewal commission is payable on subsequent year premiums after completion of first 12 months:

| Policy Procurement Period | Eligible Agents | Renewal Incentive Rate |
|---------------------------|-----------------|------------------------|
| 01.10.2009 - 31.03.2017 | Field Officers, Direct Agents | 1% of renewal premium |
| On or after 01.07.2020 | All Agent Types | 1% of renewal premium |
| PLI (Cash Policies) | All Agent Types | 1% of renewal premium |
| RPLI Policies | All Agent Types | 2.5% of renewal premium |

**Notes:**
- No renewal incentive for pay policies
- Agent must be active at time of policy procurement for eligibility
- Renewal commission continues as long as policy is in-force and premiums are paid

### 11.4 Monitoring Staff Procurement Incentive

Monitoring staff receive procurement incentive on new business premium for sales force reporting to them:

| Monitoring Role | Rate on New Business Premium | Reporting Staff |
|-----------------|------------------------------|-----------------|
| Development Officer | 0.8% | DA/FO under them |
| Sub-Divisional Head | 0.6% | GDS under them |
| Mail Overseer | 0.2% | GDS under them |
| Sub-Divisional Head | 0.8% | DE under them |
| ASP(HQ)/Office Superintendent | 0.8% | DE under them |
| Divisional Head | 0.2% | All sales force in division |

**Important Notes:**
- No renewal incentive is payable to monitoring staff
- Incentive is calculated on the actual premium collected (first year only)
- Multiple monitoring staff may receive incentive for the same policy based on hierarchy

### 11.5 Rate Table Configuration

The commission rate table must support the following fields:

| Field | Description | Data Type |
|-------|-------------|-----------|
| Rate (%) | Commission percentage | Decimal(5,2) |
| Policy Duration | Duration in months | Integer |
| Product Type | PLI/RPLI | Varchar(10) |
| Plan Code | Product plan code | Varchar(50) |
| Agent Type | DE/FO/DA/GDS/Coordinator | Varchar(20) |
| Policy Term | Term in years | Integer |
| Effective Date | Rate effective from | Date |
| Expiry Date | Rate effective to (nullable) | Date |

**Business Rule Reference:** BR-IC-COM-006

### 11.6 Commission Rate Examples

**Example 1: PLI Non-AEA Policy (Term 20 years)**
- Policy Type: PLI Whole Life
- Policy Term: 20 years
- First Year Premium: Rs. 12,000
- Agent Type: Direct Agent
- Incentive Rate: 10% (15 < Term <= 25 years)
- Gross Commission: 12,000 x 10% = Rs. 1,200

**Example 2: PLI AEA Policy (Term 10 years)**
- Policy Type: PLI AEA
- Policy Term: 10 years
- First Year Premium: Rs. 8,000
- Agent Type: Field Officer
- Incentive Rate: 5% (AEA, Term <= 15 years)
- Gross Commission: 8,000 x 5% = Rs. 400

**Example 3: RPLI Policy**
- Policy Type: RPLI
- First Year Premium: Rs. 5,000
- Agent Type: GDS
- Incentive Rate: 10%
- Gross Commission: 5,000 x 10% = Rs. 500
- Monitoring Staff Incentive:
  - Sub-Divisional Head: 5,000 x 0.6% = Rs. 30
  - Mail Overseer: 5,000 x 0.2% = Rs. 10
  - Divisional Head: 5,000 x 0.2% = Rs. 10

---
**End of Section 11: Commission Rate Structure**

**Previous Section**: [Traceability Matrix](#10-traceability-matrix)
**Next Section**: [Taxation and Compliance](#12-taxation-and-compliance)

---

# Incentive, Commission and Producer Management - Taxation and Compliance

## 12. Taxation and Compliance

### 12.1 Tax Deduction at Source (TDS)

**TDS Rate:** 2% of gross commission

**Applicability:**
- TDS is deducted from all commission payments where PAN is available
- If PAN is not available, higher TDS rate may apply as per Income Tax rules

**TDS Calculation:**
```
Gross Commission = Sum of all eligible commissions for the period
TDS Amount = Gross Commission x 2%
Net Payable = Gross Commission - TDS Amount
```

**Example:**
- Gross Commission: Rs. 50,000
- TDS @ 2%: Rs. 1,000
- Net Payable: Rs. 49,000

**Business Rule Reference:** BR-IC-COM-003

### 12.2 Goods and Services Tax (GST)

**GST Rate:** 18% (Reverse Charge Mechanism - RCM)

**Payment Responsibility:**
- GST is paid by the Department (not deducted from agent commission)
- Commission amount payable to agent is gross of GST
- Department pays GST under RCM and claims input tax credit

**GST Calculation:**
```
Commission Value = Gross Commission paid to agent
GST Liability = Commission Value x 18%
```

**Example:**
- Commission paid to agent: Rs. 50,000
- GST payable by Department: Rs. 9,000 (18% of Rs. 50,000)

### 12.3 Monthly Tax Liability Reporting

The system must generate the following reports monthly:

**TDS Liability Report:**
| Column | Description |
|--------|-------------|
| Agent Code | Unique agent identifier |
| Agent Name | Name of agent |
| PAN Number | PAN of agent |
| Gross Commission | Total commission before TDS |
| TDS Rate | TDS percentage (2%) |
| TDS Amount | T deducted at source |
| Net Payable | Amount payable to agent |

**GST Liability Report:**
| Column | Description |
|--------|-------------|
| Reporting Period | Month/Year |
| Total Commission Paid | Sum of all commissions |
| GST Rate | GST percentage (18%) |
| GST Liability | GST payable under RCM |
| Commission Category | First Year / Renewal |

### 12.4 Annual TDS Reconciliation

**Form 16A Generation:**
- System must generate Form 16A for each agent
- Form 16A contains TDS certificate for tax filing
- Must be generated annually for each agent who received commission

**Quarterly TDS Returns:**
- Form 24Q: Salary TDS (if applicable)
- Form 26Q: Non-salary TDS (commission payments)
- Due dates: Quarterly (July, October, January, May)

**Annual TDS Statement:**
| Column | Description |
|--------|-------------|
| Agent Code | Unique agent identifier |
| Agent Name | Name of agent |
| PAN Number | PAN of agent |
| Financial Year | April to March |
| Total Gross Commission | Sum of all commissions |
| Total TDS Deducted | Sum of all TDS |
| TDS Deposited | Amount deposited with tax department |

### 12.5 GST Compliance

**Monthly GST Returns:**
- GST-3B: Monthly self-assessment return
- Due date: 20th of following month
- Contains: outward supplies, input tax credit, tax payable

**GST Reconciliation:**
- Monthly reconciliation of GST paid vs. commission paid
- Annual reconciliation with books of accounts
- GST audit trail maintenance

**Compliance Requirements:**
| Requirement | Description | Frequency |
|-------------|-------------|-----------|
| GST Payment | Payment of GST under RCM | Monthly |
| GST Return Filing | GST-3B filing | Monthly |
| GST Reconciliation | Match books with GST returns | Monthly |
| GST Audit Trail | Maintain records | Ongoing |

### 12.6 Income Tax Compliance

**TDS Compliance Requirements:**
| Requirement | Description | Frequency |
|-------------|-------------|-----------|
| TDS Deduction | Deduct TDS from commission | Every payment |
| TDS Deposit | Deposit TDS with government | Monthly (7th of next month) |
| TDS Return | File quarterly returns | Quarterly |
| Form 16A | Issue TDS certificate to agents | Annual |
| TDS Reconciliation | Annual TDS reconciliation | Annual |

**Penalty for Non-Compliance:**
- Late TDS deposit: Interest @ 1% per month
- Late return filing: Late fee up to Rs. 5,000
- Non-deduction of TDS: Disallowed expense + penalty

---
**End of Section 12: Taxation and Compliance**

**Previous Section**: [Commission Rate Structure](#11-commission-rate-structure)
**Next Section**: [Report Specifications](#13-report-specifications)

---

# Incentive, Commission and Producer Management - Report Specifications

## 13. Report Specifications

### 13.1 Report Overview

The system must generate the following reports for management review and compliance:

| Report ID | Report Name | Frequency | Primary Users |
|-----------|-------------|-----------|---------------|
| RPT-IC-001 | Agent-wise Incentive Summary | Monthly | Finance, Division Heads |
| RPT-IC-002 | Circle/Division-wise Incentive Summary | Monthly | Circle Heads, Regional Heads |
| RPT-IC-003 | Monitoring Staff Procurement Incentive Register | Monthly | Monitoring Staff, Finance |
| RPT-IC-004 | TDS and GST Summary Report | Monthly | Finance, Tax Team |
| RPT-IC-005 | Pending Approval/Disbursement Report | Daily/Weekly | Finance, Operations |
| RPT-IC-006 | Policy-wise Incentive Report | On-demand | Agents, Auditors |
| RPT-IC-007 | Agent Category-wise Incentive Report | Monthly | Management |

### 13.2 RPT-IC-001: Agent-wise Incentive Summary

**Purpose:** Provides detailed commission summary for each agent

**Frequency:** Monthly

**Report Format:**

| Agent Code | Agent Name | Agent Type | Circle | Division | No. of Policies | Total Premium | Gross Commission | TDS | Net Payable | Status |
|------------|------------|------------|--------|----------|-----------------|---------------|------------------|-----|-------------|--------|
| DA001 | Rajesh Kumar | Direct Agent | Delhi | New Delhi | 5 | 45,000 | 5,200 | 104 | 5,096 | Approved |
| FO023 | Sunita Sharma | Field Officer | Mumbai | South Mumbai | 3 | 2,80,000 | 2,800 | 56 | 2,744 | Pending |

**Report Parameters:**
- Report Period (Month/Year)
- Circle (optional filter)
- Division (optional filter)
- Agent Type (optional filter)
- Status (All/Approved/Pending/Disbursed)

### 13.3 RPT-IC-002: Circle/Division-wise Incentive Summary

**Purpose:** Aggregated commission view by geographical hierarchy

**Frequency:** Monthly

**Report Format:**

| Circle | Division | Agent Count | Total Policies | Total Premium | Gross Commission | TDS | Net Payable |
|--------|----------|-------------|----------------|---------------|------------------|-----|-------------|
| Delhi | New Delhi | 45 | 234 | 45,50,000 | 4,55,000 | 9,100 | 4,45,900 |
| Delhi | Connaught Place | 32 | 189 | 32,40,000 | 3,24,000 | 6,480 | 3,17,520 |

**Drill-down Capability:**
- Click on Division to view Agent-wise details
- Click on Agent to view Policy-wise details

### 13.4 RPT-IC-003: Monitoring Staff Procurement Incentive Register

**Purpose:** Tracks procurement incentive for monitoring staff

**Frequency:** Monthly

**Report Format:**

| Staff Code | Staff Name | Role | Circle | Division | Reporting Agents | Policies Procured | Total Premium | Incentive Rate | Incentive Amount |
|------------|------------|------|--------|----------|------------------|-------------------|---------------|----------------|-----------------|
| DH001 | Ramesh Gupta | Divisional Head | Delhi | New Delhi | 45 | 234 | 45,50,000 | 0.2% | 9,100 |
| DO001 | Suresh Patil | Development Officer | Mumbai | South Mumbai | 12 | 78 | 12,40,000 | 0.8% | 9,920 |

### 13.5 RPT-IC-004: TDS and GST Summary Report

**Purpose:** Tax compliance and liability reporting

**Frequency:** Monthly

**Report Format:**

| Category | Gross Commission | TDS @ 2% | Net Payable | GST @ 18% (RCM) | Total Liability |
|----------|------------------|----------|-------------|-----------------|-----------------|
| Direct Agents | 85,000 | 1,700 | 83,300 | 15,300 | 17,000 |
| Departmental Employees | 35,000 | 700 | 34,300 | 6,300 | 7,000 |
| Field Officers | 5,000 | 100 | 4,900 | 900 | 1,000 |
| **TOTAL** | **1,25,000** | **2,500** | **1,22,500** | **22,500** | **25,000** |

### 13.6 RPT-IC-005: Pending Approval/Disbursement Report

**Purpose:** Track commissions pending approval or disbursement

**Frequency:** Daily/Weekly

**Report Format:**

| Trial Statement ID | Period | Agent Count | Total Amount | Pending Since | Age (Days) | Status | Action Required |
|--------------------|--------|-------------|--------------|---------------|------------|--------|-----------------|
| TS-2026-01-001 | Jan 2026 | 45 | 4,55,000 | 2026-02-01 | 3 | Pending Approval | Finance Review |
| TS-2026-01-002 | Jan 2026 | 32 | 3,24,000 | 2026-02-03 | 1 | Approved | Process Payment |

### 13.7 RPT-IC-006: Policy-wise Incentive Report

**Purpose:** Detailed commission at individual policy level

**Frequency:** On-demand

**Report Format:**

| Policy Number | Policy Type | Agent Code | Agent Name | Premium | Commission Rate | Commission Amount | TDS | Net Amount | Disbursement Date |
|---------------|-------------|------------|------------|---------|-----------------|-------------------|-----|------------|-------------------|
| PLI123456789 | Whole Life | DA001 | Rajesh Kumar | 12,000 | 10% | 1,200 | 24 | 1,176 | 2026-02-05 |
| RPLI987654321 | Whole Life | FO023 | Sunita Sharma | 5,000 | 10% | 500 | 10 | 490 | 2026-02-06 |

### 13.8 RPT-IC-007: Agent Category-wise Incentive Report

**Purpose:** Comparative analysis by agent category

**Frequency:** Monthly

**Report Format:**

| Agent Category | Agent Count | Policies | First Year Premium | FY Commission | Renewal Premium | Renewal Commission | Total Commission |
|----------------|-------------|----------|--------------------|---------------|-----------------|-------------------|------------------|
| Direct Agents | 120 | 890 | 45,00,000 | 4,50,000 | 12,50,000 | 1,25,000 | 5,75,000 |
| Field Officers | 85 | 654 | 32,40,000 | 3,24,000 | 9,80,000 | 98,000 | 4,22,000 |
| Dept Employees | 45 | 345 | 18,20,000 | 1,82,000 | 5,40,000 | 54,000 | 2,36,000 |

### 13.9 Report Export Formats

All reports support the following export formats:
- **Excel (.xlsx)** - For data analysis
- **PDF** - For printing and archival
- **CSV** - For data interchange

### 13.10 Report Distribution

**Report Distribution Schedule:**

| Report | Distribution | Recipients | Method |
|--------|--------------|------------|--------|
| Agent-wise Incentive Summary | Monthly (by 5th) | Division Heads, Circle Heads | Email |
| TDS/GST Summary | Monthly (by 5th) | Finance Head, Tax Team | Email |
| Pending Approval | Daily | Finance Team | Dashboard |
| Monitoring Staff Register | Monthly (by 5th) | All Monitoring Staff | Email |

---
**End of Section 13: Report Specifications**

**Previous Section**: [Taxation and Compliance](#12-taxation-and-compliance)
**Next Section**: [Exception Handling](#14-exception-handling)

---

# Incentive, Commission and Producer Management - Exception Handling

## 14. Exception Handling

### 14.1 Overview

The commission system must handle various exception scenarios that arise during processing. This section defines the exception handling logic for each category.

### 14.2 Data Validation Failures

#### 14.2.1 Missing Agent Details

**Scenario:** Commission calculation triggered but agent details incomplete

**Handling Logic:**
```
IF agent_details_missing THEN
  1. Flag policy for manual review
  2. Generate exception report for admin
  3. Hold commission until resolution
  4. Notify admin via dashboard alert
  5. Set policy.commission_status = 'HELD_MISSING_AGENT'
END
```

**Required Fields for Commission Processing:**
- Agent Code
- Agent Name
- PAN Number
- Bank Account Details (for EFT)
- Active Status

**SLA for Resolution:** To be decided by CEPT

#### 14.2.2 Invalid Policy Data

**Scenario:** Policy data fails validation rules

**Validation Checks:**
| Field | Validation Rule | Action on Failure |
|-------|-----------------|-------------------|
| Premium Amount | Must be > 0 | Reject calculation |
| Policy Term | Must be within defined range | Reject calculation |
| Acceptance Date | Must not be future | Reject calculation |
| Policy Status | Must be 'Accepted' | Hold calculation |
| Free-Look Status | Must be 'Completed' | Hold calculation |

**Error Handling:**
```
IF validation_fails THEN
  1. Log error with details
  2. Mark policy.commission_eligible = FALSE
  3. Set policy.rejection_reason = validation_error
  4. Notify source system
  5. Add to exception report
END
```

#### 14.2.3 Hierarchy Mismatch

**Scenario:** Agent not properly mapped to monitoring staff

**Handling Logic:**
```
IF monitoring_staff_not_found THEN
  1. Allocate agent commission (allow)
  2. Hold monitoring staff commission
  3. Generate exception report:
     - Agent Code
     - Missing Monitoring Staff Role
     - Expected Reporting Line
  4. Notify administrator
  5. Set status = 'MONITORING_STAFF_MISSING'
END
```

**Roles Validated:**
- Development Officer
- Sub-Divisional Head
- Mail Overseer
- ASP(HQ)/Office Superintendent
- Divisional Head

### 14.3 Premium Realization Issues

#### 14.3.1 Premium Not Realized

**Scenario:** Policy issued but premium not yet credited

**Handling Logic:**
```
IF premium_realization_date IS NULL THEN
  1. Commission calculation on HOLD
  2. Monitor for 90 days from policy acceptance
  3. IF NOT realized within 90 days THEN
     Mark policy.status = 'PREMIUM_PENDING'
     No commission payable until realization
  END
END
```

**Monitoring States:**
| State | Description | Action |
|-------|-------------|--------|
| PREMIUM_PENDING | Premium not realized | Monitor daily |
| PARTIALLY_REALIZED | Partial premium received | Calculate on realized portion |
| FULLY_REALIZED | Full premium received | Process commission |
| REALIZATION_FAILED | Premium bounced | Clawback if already paid |

#### 14.3.2 Partial Realization

**Scenario:** Only partial premium received

**Handling Logic:**
```
IF partial_premium_received THEN
  1. Calculate commission on realized portion
  2. Hold balance commission until full realization
  3. Track partial payments in commission_partial_ledger
  4. Cumulative commission = SUM of all partial commissions
END
```

**Partial Ledger Structure:**
| Field | Description |
|-------|-------------|
| Policy Number | Policy identifier |
| Installment Number | Which installment (1, 2, 3...) |
| Due Amount | Expected premium |
| Realized Amount | Actual received |
| Commission Calculated | Commission on this portion |
| Balance Pending | Remaining commission |

#### 14.3.3 Premium Bounce/Reversal

**Scenario:** Premium paid but later bounced or reversed

**Handling Logic:**
```
IF premium_bounced AND commission_already_paid THEN
  1. Create clawback entry immediately
  2. Calculate: clawback_amount = commission_paid
  3. Negative entry in next commission cycle
  4. Net payable = current_commission - clawback_amount
  5. IF insufficient_commission THEN
     Initiate recovery through establishment
  END
END
```

### 14.4 Integration Failures

#### 14.4.1 Finacle Integration Failure

**Scenario:** Payment file transfer to Finacle fails

**Retry Logic:**
```
ATTEMPT_COUNT = 0
MAX_ATTEMPTS = 3
RETRY_INTERVAL = 1 hour

WHILE ATTEMPT_COUNT < MAX_ATTEMPTS:
  TRY
    send_payment_file_to_finacle()
    IF success THEN
      mark_disbursement_status = 'SENT_TO_FINACLE'
      BREAK
    END
  CATCH error
    ATTEMPT_COUNT++
    LOG error_details
    IF ATTEMPT_COUNT >= MAX_ATTEMPTS THEN
      mark_disbursement_status = 'FAILED'
      move_to_manual_queue()
      send_email_alert('CEPT', 'Finance Team')
      generate_failure_report()
    END
    SLEEP(RETRY_INTERVAL)
  END
END
```

**Failure Report Contents:**
| Field | Description |
|-------|-------------|
| File ID | Unique identifier |
| Attempt Count | Number of retries |
| Error Code | Integration error code |
| Error Message | Detailed error |
| Timestamp | Failure time |
| Records Affected | Count of payments |

#### 14.4.2 Source System Data Feed Failure

**Scenario:** Policy data feed not received

**Fallback Logic:**
```
IF policy_feed_not_available THEN
  1. Use previous day's data for critical reports
  2. Flag calculations as 'PROVISIONAL'
  3. Display warning on dashboard
  4. ON feed_restored:
     Reconcile provisional vs actual
     Adjust if discrepancies found
  END
END
```

**Provisional Calculation Warning:**
> "Data from source system not available. Calculations are based on previous day's data and marked as PROVISIONAL. Reconciliation will be performed when feed is restored."

### 14.5 Calculation Errors

#### 14.5.1 Rate Configuration Error

**Scenario:** Commission rate not found in rate table

**Validation (Pre-processing):**
```
BEFORE batch_processing:
  FOR EACH policy_in_batch:
    rate = lookup_rate_table(policy.product_type, policy.term, policy.agent_type)
    IF rate IS NULL THEN
      halt_batch_processing()
      log_error('Rate not found', policy details)
      notify_admin('Rate configuration missing')
      RETURN ERROR
    END
  END
END
```

**Error Notification:**
```
TO: CEPT Configuration Team
SUBJECT: Commission Rate Missing
BODY:
Rate not found for:
- Product Type: {product_type}
- Policy Term: {policy_term}
- Agent Type: {agent_type}
Batch processing halted until rate is configured.
```

#### 14.5.2 Arithmetic Overflow/Underflow

**Scenario:** Commission amount exceeds system limits

**Boundary Checks:**
| Check | Min | Max | Action |
|-------|-----|-----|--------|
| Single Policy Commission | Re. 1 | Rs. 50,000 | Log exception |
| Agent Monthly Commission | Re. 1 | Rs. 5,00,000 | Log exception |
| Division Monthly Commission | Re. 1 | Rs. 50,00,000 | Log exception |

**Exception Handling:**
```
IF commission_amount < min_limit OR commission_amount > max_limit THEN
  1. Log exception with details
  2. Add to manual_review_queue
  3. Notify admin with policy details
  4. Continue processing other records
END
```

### 14.6 Exception Categories Summary

| Exception Type | Severity | Auto-Resolution | Manual Intervention | SLA |
|----------------|----------|-----------------|---------------------|-----|
| Missing Agent Details | High | No | Yes | To be defined |
| Invalid Policy Data | High | No | No | Reject |
| Hierarchy Mismatch | Medium | No | Yes | To be defined |
| Premium Not Realized | Medium | Yes | No | 90 days |
| Premium Bounce | High | Partial | Yes (recovery) | Immediate |
| Finacle Failure | Critical | Yes (retry) | Yes (after 3 attempts) | 1 hour |
| Rate Configuration Error | Critical | No | Yes | Immediate |
| Arithmetic Overflow | Medium | No | Yes | To be defined |

---
**End of Section 14: Exception Handling**

**Previous Section**: [Report Specifications](#13-report-specifications)
**Next Section**: [User Interface Requirements](#15-user-interface-requirements)

---

# Incentive, Commission and Producer Management - User Interface Requirements

## 15. User Interface Requirements

### 15.1 Dashboard - Management View

**Purpose:** Provide real-time visibility into commission processing status and key metrics

#### 15.1.1 Key Performance Indicators (KPIs)

| KPI | Description | Update Frequency |
|-----|-------------|------------------|
| Total Commission Payable | Sum of all approved commissions | Real-time |
| Total Policies Processed | Count of policies in current batch | Real-time |
| Total Agents Eligible | Count of agents with payable commission | Real-time |
| Average Commission per Agent | Total / Agent Count | Real-time |
| Exception Count | Pending approvals, errors | Real-time |
| Pending Approvals | Count awaiting approval | Hourly |
| Pending Disbursements | Count awaiting payment | Hourly |

#### 15.1.2 Filters

**Dashboard Filters:**
| Filter | Options | Description |
|--------|---------|-------------|
| Date Range | Month/Year selector | Select calculation period |
| Circle | Dropdown | Filter by circle |
| Division | Dropdown | Filter by division |
| Agent Type | Multi-select | DE, FO, DA, GDS, Coordinator |
| Policy Type | Multi-select | PLI, RPLI, AEA, Non-AEA |
| Commission Type | Multi-select | First Year, Renewal, Bonus |

#### 15.1.3 Visualizations

**Charts Required:**
1. **Bar Chart** - Circle-wise commission comparison
2. **Line Chart** - Trend analysis (last 6 months)
3. **Pie Chart** - Policy type distribution
4. **Data Table** - Drill-down capability with sorting

**Trend Analysis Metrics:**
- Commission by month (6 months)
- Policies by month (6 months)
- Agent performance trends
- Division-wise growth

**Top 10 Performing Agents Table:**
| Rank | Agent Code | Agent Name | Division | Policies | Commission |
|------|------------|------------|----------|----------|------------|
| 1 | DA001 | Rajesh Kumar | New Delhi | 45 | 85,000 |
| 2 | FO023 | Sunita Sharma | South Mumbai | 38 | 72,000 |

### 15.2 Commission Calculation Screen

**Purpose:** Interface for running monthly commission calculation batch

**User Role:** Commission Processing Officer

#### 15.2.1 Screen Elements

**Section 1: Period Selection**
```
┌─────────────────────────────────────────────┐
│ Calculation Period: [Month] [Year]          │
│ Scope: [All Circles] [Specific Circle]      │
│                                             │
│ [Calculate Commission] [View History]       │
└─────────────────────────────────────────────┘
```

**Section 2: Progress Indicator**
```
┌─────────────────────────────────────────────┐
│ Progress: ████████░░░░ 60%                 │
│                                             │
│ Policies Processed: 5,234 / 8,500          │
│ Agents Calculated: 145 / 200               │
│ Errors: 0                                  │
│                                             │
│ Estimated Time Remaining: 1 hr 15 min       │
└─────────────────────────────────────────────┘
```

**Section 3: Summary (After Completion)**
| Metric | Value |
|--------|-------|
| Total Policies | 8,500 |
| Eligible for Commission | 7,850 |
| Total Agents | 200 |
| Total Gross Commission | Rs. 45,50,000 |
| Total TDS | Rs. 91,000 |
| Total Net Payable | Rs. 44,59,000 |

#### 15.2.2 Validations

**Pre-calculation Checks:**
- Cannot process same period twice (unless cancelled)
- Cannot process future periods
- Must complete previous period before next
- Rate table must be configured

**Error Messages:**
| Condition | Message |
|-----------|---------|
| Period already processed | "Commission for January 2026 has already been processed. Cancel previous run first." |
| Future period selected | "Cannot process commission for future periods." |
| Previous period incomplete | "Complete December 2025 processing before starting January 2026." |
| Rates not configured | "Commission rates not configured for selected period. Configure rates first." |

### 15.3 Approval Workflow Screen

**Purpose:** Interface for Divisional Heads to approve trial statements

**User Role:** Approving Authority (Divisional Head)

#### 15.3.1 Pending Approvals Queue

**Table Layout:**
| Select | Trial ID | Period | Division | Agent Count | Total Amount | Pending Since | Actions |
|--------|----------|--------|----------|-------------|--------------|---------------|---------|
| [ ] | TS-2026-01-001 | Jan 2026 | New Delhi | 45 | 4,55,000 | 3 days | [View] |
| [ ] | TS-2026-01-002 | Jan 2026 | Connaught Place | 32 | 3,24,000 | 2 days | [View] |

**Actions:**
- [Approve Selected] - Bulk approve
- [Reject Selected] - Bulk reject with remarks
- [View Details] - View individual agent breakdown

#### 15.3.2 Detail View

**Agent-wise Breakdown:**
| Agent Code | Agent Name | Policies | Premium | Commission | TDS | Net Amount |
|------------|------------|----------|---------|------------|-----|------------|
| DA001 | Rajesh Kumar | 5 | 45,000 | 5,200 | 104 | 5,096 |
| DA002 | Amit Verma | 3 | 28,000 | 2,800 | 56 | 2,744 |

**Approval Actions:**
```
┌─────────────────────────────────────────────┐
│ [Approve]  [Reject]  [Request Clarification]│
│                                             │
│ Remarks (if reject): [________________]     │
└─────────────────────────────────────────────┘
```

#### 15.3.3 Comparison View

**Compare with Previous Month:**
| Metric | Current Month | Previous Month | Variance |
|--------|---------------|----------------|----------|
| Total Commission | 4,55,000 | 4,20,000 | +8.3% |
| Agent Count | 45 | 42 | +7.1% |
| Average per Agent | 10,111 | 10,000 | +1.1% |

### 15.4 Agent Self-Service Portal

**Purpose:** Enable agents to view their commission statements and payment status

**User Role:** Agent

#### 15.4.1 Dashboard

**Agent View:**
| Section | Content |
|---------|---------|
| Welcome | Agent Name, Code, Division |
| Summary | Current Month Commission, YTD Commission, Last Payment |
| Quick Links | View Statements, Download Certificate, Raise Query |

#### 15.4.2 Commission Statements

**Month-wise List:**
| Month/Year | Gross Commission | TDS | Net Payable | Status | Actions |
|------------|------------------|-----|-------------|--------|---------|
| Jan 2026 | 5,200 | 104 | 5,096 | Paid | [View] [Download] |
| Dec 2025 | 4,800 | 96 | 4,704 | Paid | [View] [Download] |

#### 15.4.3 Policy-wise Breakdown

**Detailed View:**
| Policy Number | Policy Type | Premium | Rate | Commission | Date |
|---------------|-------------|---------|------|------------|------|
| PLI123456789 | Whole Life | 12,000 | 10% | 1,200 | 05-Jan-26 |

**Features:**
- Filter by commission type (First Year/Renewal)
- Search by policy number
- Export to PDF

#### 15.4.4 Query/Dispute Management

**Raise Query Form:**
```
┌─────────────────────────────────────────────┐
│ Query Type: [Commission Discrepancy]        │
│                                             │
│ Policy Number: [_______________]            │
│ Month/Year: [_______] [_____]              │
│                                             │
│ Description:                                │
│ [_____________________________________]     │
│ [_____________________________________]     │
│                                             │
│ [Submit Query] [Cancel]                     │
└─────────────────────────────────────────────┘
```

### 15.5 Search and Filter Requirements

#### 15.5.1 Global Search

**Search Fields:**
- Agent Code
- Agent Name
- Policy Number
- Transaction ID
- PAN Number
- Mobile Number

**Features:**
- Auto-suggest functionality
- Recent searches history
- Search within results

#### 15.5.2 Advanced Filters

**Filter Options:**
| Filter | Type | Options |
|--------|------|---------|
| Date Range | Date Range | From/To picker |
| Amount Range | Number Range | Min/Max input |
| Status | Multi-select | Pending/Approved/Paid/Rejected/On Hold |
| Circle | Dropdown | All circles |
| Division | Dropdown | All divisions in selected circle |
| Agent Type | Multi-select | DE/FO/DA/GDS/Coordinator |
| Policy Type | Multi-select | PLI/RPLI/AEA/Non-AEA |
| Commission Type | Multi-select | First Year/Renewal/Bonus |

**Filter Presets:**
- Save frequently used filter combinations
- Quick access to saved presets
- Share presets with other users

### 15.6 Export Capabilities

#### 15.6.1 Supported Formats

| Format | Use Case | Features |
|--------|----------|----------|
| Excel (.xlsx) | Data analysis | Multiple sheets, formulas |
| PDF | Statements, certificates | Watermark, digital signature |
| CSV | Data interchange | Plain text, compatible |

#### 15.6.2 Export Options

**Scope Selection:**
- Current page data
- Selected records
- All filtered records (with warning for large datasets)

**Content Options:**
- Include filters and summary
- Include/exclude audit trail
- Include/exclude TDS details

#### 15.6.3 Security

**Export Security Features:**
- Watermark on exported PDFs (user ID, timestamp)
- Log all export activities with user ID and timestamp
- Limit bulk exports (max 50,000 records per export)
- Mask sensitive data for non-privileged users

---
**End of Section 15: User Interface Requirements**

**Previous Section**: [Exception Handling](#14-exception-handling)
**Next Section**: [Security and Access Control Details](#16-security-and-access-control-details)

---

# Incentive, Commission and Producer Management - Security and Access Control Details

## 16. Security and Access Control Details

### 16.1 Role Definitions and Permissions

#### 16.1.1 Role Hierarchy

```
PLI Directorate
    │
    ├─→ Circle Head
    │       │
    │       └─→ Regional Head
    │               │
    │               └─→ Divisional Head
    │                       │
    │                       ├─→ Development Officer
    │                       └─→ ASP(HQ)/Office Superintendent
    │
    └─→ CEPT (Technical)
        │
        └─→ Audit
```

#### 16.1.2 Detailed Permissions Matrix

| Role | View Data | Run Calculations | Approve | Configure | Reports | Edit Rights |
|------|-----------|------------------|---------|-----------|---------|-------------|
| **Divisional Head** | Division only | Yes | Level 1 | No | Division | Agent data only |
| **Circle Head** | Circle only | No | No | No | Circle | No |
| **Regional Head** | Region only | No | No | No | Region | No |
| **PLI Directorate** | All | No | No | No | All | No |
| **CEPT** | All | Yes | No | Yes | All | System configuration |
| **Audit** | All (read-only) | No | No | No | All | No |
| **Agent** | Own only | No | No | No | Own only | No |

### 16.2 Maker-Checker Requirements

#### 16.2.1 Critical Operations

The following operations require Maker-Checker workflow:

| Operation | Maker | Checker | Escalation |
|-----------|-------|---------|------------|
| Rate configuration changes | CEPT Admin | CEPT Supervisor | Director |
| Manual commission adjustments | Finance Officer | Finance Manager | CFO |
| Clawback entries | Finance Officer | Finance Manager | CFO |
| Payment file generation | Finance Officer | Finance Manager | CFO |
| User role assignments | CEPT Admin | CEPT Supervisor | Director |

#### 16.2.2 Workflow Process

```
MAKER initiates transaction
    │
    ├─→ System logs: maker_id, timestamp, ip_address
    │
    ├─→ Transaction moves to checker_queue
    │
    └─→ CHECKER reviews
            │
            ├─→ APPROVE
            │       ├─→ System logs: checker_id, timestamp
            │       └─→ Transaction executed
            │
            └─→ REJECT
                    ├─→ System logs: rejection_reason
                    └─→ Returns to maker with comments
```

#### 16.2.3 Segregation Rules

| Rule | Description |
|------|-------------|
| **Same User Prevention** | Same user cannot be maker and checker for same transaction |
| **Grade Requirement** | Checker must be of higher or equal grade than maker |
| **Time Limit** | Transactions must be reviewed within 24 hours |
| **Queue Visibility** | Makers cannot view their own pending transactions |

### 16.3 Data Privacy and Confidentiality

#### 16.3.1 PII Protection

**Personally Identifiable Information (PII) Masking:**

| Data Field | Masking Rule | Example (Original) | Example (Masked) |
|------------|--------------|--------------------|------------------|
| Bank Account Number | Show last 4 digits | 123456789012 | XXXXXXXX9012 |
| PAN Number | Show last 4 characters | ABCDE1234F | XXXXX1234F |
| Phone Number | Middle digits masked | 9876543210 | 98XXX3210 |
| Email Address | Partial masking | agent@email.com | aXXXn@email.com |

#### 16.3.2 Data Masking Rules

**User Access Levels:**

| Access Level | PAN | Bank Account | Phone | Email | Full Access |
|--------------|-----|--------------|-------|-------|-------------|
| Administrator | Yes | Yes | Yes | Yes | All fields |
| Divisional Head | Masked | Masked | Masked | Masked | No |
| Circle Head | Masked | Masked | Masked | Masked | No |
| Audit | Yes | Yes | Yes | Yes | All fields (read-only) |
| Agent | Own only | Own only | Own only | Own only | Own data only |

#### 16.3.3 Audit Logging

**PII Access Log:**
```
Every access to unmasked PII data is logged:
- User ID
- Timestamp
- Data accessed (which fields)
- Record identifier (agent code)
- Access reason (where applicable)
- IP address
- Session ID
```

### 16.4 Session Management

#### 16.4.1 Session Security

| Setting | Value | Description |
|---------|-------|-------------|
| Session Timeout | 5 minutes | Auto-logout after inactivity |
| Concurrent Sessions | 1 active | New login terminates previous |
| Force Logout | On role change | Immediate re-authentication required |
| IP Logging | Per session | Track IP for security monitoring |
| Session Encryption | TLS 1.3 | All traffic encrypted |

#### 16.4.2 Password Policy

| Requirement | Specification |
|-------------|---------------|
| Minimum Length | 8 characters |
| Complexity | 1 uppercase, 1 lowercase, 1 number, 1 special character |
| Expiry | 180 days |
| History | Cannot reuse last 5 passwords |
| Lockout | 3 failed attempts |
| Lockout Duration | 30 minutes or admin unlock |
| Reset | Email verification required |

### 16.5 Audit and Compliance

#### 16.5.1 Audit Trail Capture

**All User Actions Logged:**
| Field | Description | Example |
|-------|-------------|---------|
| User ID | Who performed action | user123 |
| Timestamp | When action occurred | 2026-01-26 10:30:45 |
| Action | What was done | COMMISSION_APPROVAL |
| Old Value | Value before change | PENDING |
| New Value | Value after change | APPROVED |
| IP Address | Where from | 192.168.1.100 |
| Session ID | Session identifier | sess_abc123 |
| Record ID | Affected record | TS-2026-01-001 |

#### 16.5.2 Audit Trail Retention

| Requirement | Specification |
|-------------|---------------|
| Minimum Retention | 10 years |
| Archive Strategy | Older than 1 year moved to archive storage |
| Retrieval SLA | < 5 minutes for recent, < 1 day for archived |
| Search Capability | Full-text search on all fields |

#### 16.5.3 Compliance Reporting

**Internal Audit:**
- Monthly exception reports
- Quarterly access review
- Annual compliance audit

**CAG Audit Support:**
- Data extraction capability
- Reconciliation reports
- Historical data retrieval

**RTI Compliance:**
- Mechanism to extract specific citizen data
- Response tracking
- Disclosure logging

### 16.6 Data Encryption

#### 16.6.1 Encryption Standards

| Data State | Encryption Standard | Key Management |
|------------|---------------------|----------------|
| At Rest | AES-256 | Hardware Security Module (HSM) |
| In Transit | TLS 1.3 | SSL/TLS certificates |
| Backup | AES-256 | Separate key management |

#### 16.6.2 Encrypted Fields

| Field | Encryption | Reason |
|-------|------------|--------|
| PAN | AES-256 | Tax identifier |
| Bank Account | AES-256 | Financial data |
| Phone Number | AES-256 | Personal contact |
| Email Address | AES-256 | Personal contact |
| Date of Birth | AES-256 | Personal information |

### 16.7 Access Control Examples

#### 16.7.1 Divisional Head Access

**Can View:**
- All agents in their division
- Commission data for their division
- Trial statements for approval
- Disbursement status

**Cannot View:**
- Other divisions' data
- System configuration
- Other users' activities

**Can Perform:**
- Approve/reject trial statements
- View agent profiles (division only)
- Generate division reports

#### 16.7.2 Agent Access

**Can View:**
- Own profile (masked for display)
- Own commission statements
- Own payment history
- Own policy-wise commission

**Cannot View:**
- Other agents' data
- Commission rates
- System reports

**Can Perform:**
- Download own statements
- View TDS certificate
- Raise query/dispute

---
**End of Section 16: Security and Access Control Details**

**Previous Section**: [User Interface Requirements](#15-user-interface-requirements)
**Next Section**: [Performance SLAs](#17-performance-slas)

---

# Incentive, Commission and Producer Management - Performance SLAs

## 17. Performance SLAs

### 17.1 Response Time Requirements

#### 17.1.1 User Interface Response Times

| Operation | Expected Response Time | Measurement Method |
|-----------|------------------------|-------------------|
| User Login | < 2 seconds | Time from credentials submit to dashboard load |
| Dashboard Load | < 3 seconds | Time from request to full render |
| Search Results | < 2 seconds | Time from query submit to results display |
| Filter Application | < 1 second | Time from filter change to updated view |
| Report Generation (< 1000 records) | < 5 seconds | Time from request to download ready |
| Report Generation (< 10,000 records) | < 30 seconds | Time from request to download ready |
| Export to Excel (< 5000 records) | < 10 seconds | Time from request to file download |

#### 17.1.2 Commission Processing Response Times

| Operation | Expected Response Time | Notes |
|-----------|------------------------|-------|
| Commission Calculation (monthly batch) | < 2 hours | For ~10,000 policies |
| Payment File Generation | < 15 minutes | For ~500 agents |
| Trial Statement Generation | < 30 minutes | Depends on agent count |
| Final Statement Generation | < 15 minutes | After trial approval |
| Disbursement Processing | < 1 hour | EFT file generation |

### 17.2 Batch Processing Windows

#### 17.2.1 Daily Processing

| Task | Window | Duration | Dependencies |
|------|--------|----------|--------------|
| Policy Data Sync | 02:00 - 04:00 | 2 hours | Source system availability |
| Agent Master Sync | 04:00 - 05:00 | 1 hour | HRMS availability |
| TDS/GST Calculation | 05:00 - 06:00 | 1 hour | Commission calculation |
| SLA Monitoring | Every hour | < 1 minute | Timer trigger |

#### 17.2.2 Monthly Processing

| Task | Schedule | Window | Duration |
|------|----------|--------|----------|
| Commission Calculation | First working day | 09:00 - 13:00 | 4 hours |
| Trial Statement Generation | First working day | 13:00 - 14:00 | 1 hour |
| Approval Period | First 3 days | - | 3 days |
| Final Statement Generation | After approval | Same day | 1 hour |
| Payment File Generation | After final statement | Next day | 1 hour |

### 17.3 Concurrent User Requirements

#### 17.3.1 Expected Load

| Metric | Value | Notes |
|--------|-------|-------|
| Peak Concurrent Users | 150 | Commission processing days |
| Normal Concurrent Users | 50 | Regular days |
| Self-Service Users | 200 | Month-end (agents viewing statements) |

#### 17.3.2 Stress Testing

| Test Scenario | Load | Duration | Acceptance Criteria |
|---------------|-------|----------|---------------------|
| Normal Load | 150 users | 4 hours | Response time within SLA |
| Peak Load | 225 users (150% of peak) | 2 hours | Response time < 150% of SLA |
| Sustained Load | 150 users | 8 hours | No memory leaks, stable response |

### 17.4 Availability Requirements

#### 17.4.1 Uptime SLA

| Metric | Target | Exclusions |
|--------|--------|------------|
| System Availability | 99.5% | Planned maintenance |
| Planned Maintenance | 48 hours/year | With advance notification |
| Unplanned Downtime | < 44 hours/year | Maximum allowable |

**Availability Calculation:**
- Total hours in year: 8,760
- 99.5% availability: 8,716 hours
- Maximum downtime: 44 hours/year

#### 17.4.2 Maintenance Windows

| Type | Frequency | Duration | Notification |
|------|-----------|----------|--------------|
| Scheduled Maintenance | Monthly | 2 hours | 48 hours advance notice |
| Emergency Maintenance | As needed | 1 hour | Immediate notification to critical users |
| Security Updates | Quarterly | 4 hours | 1 week advance notice |

### 17.5 Performance Monitoring

#### 17.5.1 Key Metrics to Monitor

| Metric | Alert Threshold | Critical Threshold |
|--------|-----------------|-------------------|
| Response Time (p95) | > 3 seconds | > 5 seconds |
| Error Rate | > 1% | > 5% |
| CPU Utilization | > 70% | > 90% |
| Memory Utilization | > 75% | > 90% |
| Disk I/O Wait | > 20% | > 40% |
| Database Connection Pool | > 80% used | > 95% used |

#### 17.5.2 Monitoring Dashboard

**Real-time Metrics Display:**
```
┌─────────────────────────────────────────────────────┐
│ SYSTEM HEALTH: [GREEN]                              │
│                                                     │
│ Response Time:     1.2s ████████░░ 60% of SLA       │
│ CPU Usage:         45% ██████████ 45%               │
│ Memory Usage:      62% ████████░░ 62%               │
│ Active Users:      87 / 150 ███████░░ 58%           │
│ Error Rate:        0.2% ██░░░░░░░░ 0.2%             │
│ Queue Depth:       12 tasks █░░░░░░░░░              │
└─────────────────────────────────────────────────────┘
```

### 17.6 Performance Penalties

#### 17.6.1 SLA Breach Consequences

| SLA | Penalty | Recipient |
|-----|---------|-----------|
| Commission Disbursement > 10 days | 8% interest on delayed amount | Agent |
| Commission Batch > 6 hours | Escalation to Director | Finance Head |
| System Downtime > 4 hours | Credit to users | CEPT |

**Interest Calculation Example:**
```
Delayed Amount = Rs. 50,000
Delay = 4 days
Interest Rate = 8% per annum
Interest = 50,000 x 8% x (4/365) = Rs. 43.84
```

### 17.7 Capacity Planning

#### 17.7.1 Growth Projections

| Year | Policies | Agents | Transactions | Storage |
|------|----------|--------|--------------|---------|
| 2026 | 100,000 | 2,000 | 1.2M | 50 GB |
| 2027 | 120,000 | 2,400 | 1.44M | 60 GB |
| 2028 | 144,000 | 2,880 | 1.73M | 72 GB |

#### 17.7.2 Scalability Requirements

| Component | Current Capacity | Target Capacity | Scaling Strategy |
|-----------|------------------|-----------------|------------------|
| Application Server | 150 concurrent users | 300 concurrent users | Horizontal scaling |
| Database Server | 100 GB | 500 GB | Vertical + partitioning |
| File Storage | 200 GB | 1 TB | NAS/SAN expansion |

---
**End of Section 17: Performance SLAs**

**Previous Section**: [Security and Access Control Details](#16-security-and-access-control-details)
**Next Section**: [Sample Calculations and Examples](#18-sample-calculations-and-examples)

---

# Incentive, Commission and Producer Management - Sample Calculations and Examples

## 18. Sample Calculations and Examples

### 18.1 Example 1: PLI Non-AEA Policy (Term 20 years)

**Policy Details:**
| Field | Value |
|-------|-------|
| Policy Type | PLI Whole Life |
| Policy Term | 20 years |
| First Year Premium | Rs. 12,000 |
| Agent Type | Direct Agent |
| Agent Code | DA001 |
| PAN | Available |

**Commission Calculation:**
```
Step 1: Determine Rate
  Policy Type: Other than AEA
  Policy Term: 20 years (15 < Term <= 25)
  Applicable Rate: 10%

Step 2: Calculate Gross Commission
  Gross Commission = Premium x Rate
  Gross Commission = 12,000 x 10% = Rs. 1,200

Step 3: Calculate TDS
  TDS = Gross Commission x 2%
  TDS = 1,200 x 2% = Rs. 24

Step 4: Calculate Net Payable
  Net Payable = Gross Commission - TDS
  Net Payable = 1,200 - 24 = Rs. 1,176

Step 5: Calculate GST (Department Liability)
  GST = Gross Commission x 18%
  GST = 1,200 x 18% = Rs. 216
```

**Summary:**
| Component | Amount |
|-----------|--------|
| Gross Commission | Rs. 1,200 |
| TDS Deducted | Rs. 24 |
| Net Payable to Agent | Rs. 1,176 |
| GST Payable by Department | Rs. 216 |

### 18.2 Example 2: PLI AEA Policy (Term 10 years)

**Policy Details:**
| Field | Value |
|-------|-------|
| Policy Type | PLI AEA |
| Policy Term | 10 years |
| First Year Premium | Rs. 8,000 |
| Agent Type | Field Officer |
| Agent Code | FO023 |

**Commission Calculation:**
```
Step 1: Determine Rate
  Policy Type: AEA
  Policy Term: 10 years (Term <= 15)
  Applicable Rate: 5%

Step 2: Calculate Gross Commission
  Gross Commission = 8,000 x 5% = Rs. 400

Step 3: Calculate TDS
  TDS = 400 x 2% = Rs. 8

Step 4: Calculate Net Payable
  Net Payable = 400 - 8 = Rs. 392

Step 5: Monitoring Staff Incentive
  Development Officer (0.8%): 8,000 x 0.8% = Rs. 64
  Divisional Head (0.2%): 8,000 x 0.2% = Rs. 16
```

**Summary:**
| Component | Agent | DO | DH |
|-----------|-------|-----|-----|
| Gross Commission | Rs. 400 | Rs. 64 | Rs. 16 |
| TDS @ 2% | Rs. 8 | Rs. 1.28 | Rs. 0.32 |
| Net Payable | Rs. 392 | Rs. 62.72 | Rs. 15.68 |

### 18.3 Example 3: RPLI Policy with Monitoring Staff

**Policy Details:**
| Field | Value |
|-------|-------|
| Policy Type | RPLI |
| First Year Premium | Rs. 5,000 |
| Agent Type | GDS |
| Agent Code | GDS045 |
| Reporting Hierarchy | Sub-Divisional Head, Mail Overseer, Divisional Head |

**Commission Calculation:**
```
Step 1: Agent Commission
  RPLI Rate: 10%
  Agent Gross = 5,000 x 10% = Rs. 500
  Agent TDS = 500 x 2% = Rs. 10
  Agent Net = 500 - 10 = Rs. 490

Step 2: Monitoring Staff Commission
  Sub-Divisional Head (0.6%): 5,000 x 0.6% = Rs. 30
  Mail Overseer (0.2%): 5,000 x 0.2% = Rs. 10
  Divisional Head (0.2%): 5,000 x 0.2% = Rs. 10

Step 3: TDS for Monitoring Staff
  TDS = (30 + 10 + 10) x 2% = Rs. 1

Total Monitoring Staff Net = 50 - 1 = Rs. 49
```

**Summary:**
| Recipient | Gross | TDS | Net |
|-----------|-------|-----|-----|
| Agent (GDS) | Rs. 500 | Rs. 10 | Rs. 490 |
| Sub-Divisional Head | Rs. 30 | Rs. 0.60 | Rs. 29.40 |
| Mail Overseer | Rs. 10 | Rs. 0.20 | Rs. 9.80 |
| Divisional Head | Rs. 10 | Rs. 0.20 | Rs. 9.80 |
| **TOTAL** | **Rs. 550** | **Rs. 11** | **Rs. 539** |

### 18.4 Example 4: Multiple Policies for Same Agent

**Agent Details:**
| Field | Value |
|-------|-------|
| Agent Code | DA12345 |
| Agent Name | Rajesh Kumar |
| PAN | ABCDE1234F |

**Policies in Current Month:**

| Policy # | Type | Term | Premium | Rate | Commission |
|----------|------|------|---------|------|------------|
| PLI001 | Non-AEA | 12 years | 10,000 | 4% | 400 |
| RPLI001 | RPLI | - | 8,000 | 10% | 800 |
| PLI002 | Non-AEA | 30 years | 15,000 | 20% | 3,000 |
| PLI003 | AEA | 10 years | 5,000 | 5% | 250 |

**Consolidated Calculation:**
```
Step 1: Sum Gross Commission
  Total Gross = 400 + 800 + 3,000 + 250 = Rs. 4,450

Step 2: Calculate TDS
  Total TDS = 4,450 x 2% = Rs. 89

Step 3: Calculate Net Payable
  Net Payable = 4,450 - 89 = Rs. 4,361

Step 4: Calculate GST Liability
  GST = 4,450 x 18% = Rs. 801
```

**Trial Statement Format:**
| Description | Amount |
|-------------|--------|
| **First Year Commission** | |
| PLI001 (12 years, 4%) | 400 |
| PLI002 (30 years, 20%) | 3,000 |
| PLI003 (AEA 10 years, 5%) | 250 |
| **RPLI First Year (10%)** | |
| RPLI001 | 800 |
| **Total Gross Commission** | **4,450** |
| Less: TDS @ 2% | (89) |
| **Net Payable** | **4,361** |

### 18.5 Example 5: Renewal Commission Calculation

**Policy Details:**
| Field | Value |
|-------|-------|
| Policy Number | RPLI987654321 |
| Policy Procured | 15.08.2021 |
| Current Renewal Premium | Rs. 4,000 |
| Agent Type | Field Officer |
| Procurement Period | On or after 01.07.2020 |

**Commission Calculation:**
```
Step 1: Verify Eligibility
  Policy Age: > 12 months (eligible for renewal)
  Agent Status: Active
  Premium Status: Realized

Step 2: Determine Rate
  Policy Type: RPLI
  Renewal Rate: 2.5%

Step 3: Calculate Gross Commission
  Gross Commission = 4,000 x 2.5% = Rs. 100

Step 4: Calculate TDS
  TDS = 100 x 2% = Rs. 2

Step 5: Calculate Net Payable
  Net Payable = 100 - 2 = Rs. 98

Note: No monitoring staff incentive on renewal commission
```

### 18.6 Example 6: Commission Clawback Scenario

**Original Commission:**
| Field | Value |
|-------|-------|
| Policy Number | PLI123456789 |
| Original Commission | Rs. 1,200 |
| TDS Deducted | Rs. 24 |
| Net Paid to Agent | Rs. 1,176 |
| Payment Date | 15-Jan-2026 |

**Clawback Trigger:**
| Field | Value |
|-------|-------|
| Event | Policy Lapsed (premium bounced) |
| Date | 01-Feb-2026 |
| Clawback Amount | Rs. 1,200 |

**Clawback Calculation:**
```
Step 1: Create Clawback Entry
  Clawback Amount = Original Commission = Rs. 1,200
  Reason: Policy lapsed due to premium bounce

Step 2: Adjust in Next Cycle (Feb 2026)
  Feb Commission = Rs. 3,000
  Clawback Deduction = Rs. 1,200
  Net Feb Commission = Rs. 1,800
  Feb TDS = 1,800 x 2% = Rs. 36

Step 3: If Insufficient Commission
  If Feb Commission < Rs. 1,200:
    Carry forward balance to next month
    Initiate recovery through establishment for remainder
```

### 18.7 Example 7: Partial Disbursement

**Trial Statement Details:**
| Field | Value |
|-------|-------|
| Agent Code | DA001 |
| Total Commission | Rs. 50,000 |
| TDS | Rs. 1,000 |
| Net Payable | Rs. 49,000 |

**Partial Disbursement Scenario:**
```
Finance approves 60% partial disbursement:

Step 1: Calculate Disbursement Amount
  Disbursement Percentage = 60%
  Disbursement Amount = 49,000 x 60% = Rs. 29,400

Step 2: Calculate Pending Amount
  Pending Amount = 49,000 - 29,400 = Rs. 19,600

Step 3: Update Status
  Status = PARTIALLY_DISBURSED
  Disbursed = Rs. 29,400
  Pending = Rs. 19,600
```

**Trial Statement after Partial Disbursement:**
| Component | Amount |
|-----------|--------|
| Original Net Payable | Rs. 49,000 |
| First Disbursement (60%) | (Rs. 29,400) |
| **Balance Pending** | **Rs. 19,600** |

---
**End of Section 18: Sample Calculations and Examples**

**Previous Section**: [Performance SLAs](#17-performance-slas)
**Next Section**: [Glossary and Definitions](#19-glossary-and-definitions)

---

# Incentive, Commission and Producer Management - Glossary and Definitions

## 19. Glossary and Definitions

### 19.1 Abbreviations

| Abbreviation | Full Form | Description |
|--------------|-----------|-------------|
| **AEA** | Anticipated Endowment Assurance | A type of PLI policy with anticipated bonuses |
| **PLI** | Postal Life Insurance | Life insurance for postal and government employees |
| **RPLI** | Rural Postal Life Insurance | Life insurance for rural population |
| **FO** | Field Officer | An agent type (retired postal official) |
| **DA** | Direct Agent | An agent type directly appointed |
| **DE** | Departmental Employee | Postal department employee who can sell policies |
| **GDS** | Gramin Dak Sevak | Rural postal employee who can sell policies |
| **ASP** | Assistant Superintendent of Post Offices | A supervisory role |
| **HQ** | Headquarters | Main office location |
| **OS** | Office Superintendent | Office supervisory role |
| **Dy. SP** | Deputy Superintendent | Deputy supervisory role |
| **DO** | Development Officer | Officer overseeing field operations |
| **SDH** | Sub-Divisional Head | Head of a sub-division |
| **DH** | Divisional Head | Head of a division |
| **CBS** | Core Banking System | Banking system for financial transactions |
| **TDS** | Tax Deducted at Source | Tax deducted at source of income |
| **GST** | Goods and Services Tax | Indirect tax on goods and services |
| **RCM** | Reverse Charge Mechanism | GST paid by recipient instead of supplier |
| **PAN** | Permanent Account Number | Unique tax identifier in India |
| **UAT** | User Acceptance Testing | Testing by users before go-live |
| **SLA** | Service Level Agreement | Agreement on service levels |
| **RTO** | Recovery Time Objective | Target time for recovery after failure |
| **RPO** | Recovery Point Objective | Maximum acceptable data loss |
| **PII** | Personally Identifiable Information | Personal data requiring protection |
| **IRDAI** | Insurance Regulatory and Development Authority of India | Insurance regulator |
| **CEPT** | Centre for Excellence in Postal Technology | Technical support for India Post |
| **CAG** | Comptroller and Auditor General | Supreme audit institution of India |
| **RTI** | Right to Information | Indian law for public information access |

### 19.2 Business Terms

#### 19.2.1 Free-Look Period
**Definition:** 15-day window from policy issuance during which the policyholder can cancel without penalty.

**Impact on Commission:**
- Commission is NOT payable until free-look period completes
- If policy lapses during free-look, no commission is payable
- Free-look completion triggers commission eligibility

#### 19.2.2 Premium Realization
**Definition:** Actual receipt and confirmation of premium payment in the Department's account.

**Types:**
- **Offline Realization:** Cash/Cheque/DD credited to Department account
- **Online Realization:** Net banking/UPI/Card payment + T+5 working days

**Impact on Commission:**
- Commission calculated only on realized premiums
- Online payments have T+5 buffer for chargeback protection

#### 19.2.3 First Year Premium
**Definition:** Total premium collected during the first policy year (from policy inception to first renewal due date).

**Usage:** Basis for first-year commission calculation at applicable rates.

#### 19.2.4 Renewal Premium
**Definition:** Premium collected in subsequent policy years (2nd year onwards) to keep the policy in-force.

**Usage:** Basis for renewal commission calculation at lower rates.

#### 19.2.5 Cash Policy
**Definition:** PLI policies where premium is paid in cash (not through salary deduction).

**Commission Impact:**
- Eligible for 1% renewal commission (same as other policies from 01.07.2020)

#### 19.2.6 Policy Term
**Definition:** Duration of the policy in years as specified at policy issuance.

**Impact on Commission:**
- Determines applicable commission rate (especially for PLI non-AEA)
- Longer terms = higher commission rates

#### 19.2.7 Procurement
**Definition:** Act of enrolling a new policy by an agent.

**Usage:** Used in context of monitoring staff incentive (procurement incentive).

#### 19.2.8 Clawback
**Definition:** Recovery of commission already paid due to policy cancellation, lapse, or premium reversal.

**Triggers:**
- Premium bounce/reversal
- Policy cancellation during free-look
- Fraudulent policy detection
- Policy lapsed due to non-payment

#### 19.2.9 T+5
**Definition:** Transaction date plus 5 working days.

**Usage:** Settlement period for online payments before commission eligibility.
- Day 0: Transaction date
- Day 1-5: Working days (excluding holidays)
- Day 6: Commission eligible

### 19.3 System Terms

#### 19.3.1 Batch Processing
**Definition:** Automated execution of commission calculation for all eligible policies in a defined period.

**Schedule:** First working day of each month

**Process:**
1. Fetch eligible policies
2. Calculate commission for each
3. Generate trial statements
4. Apply validations
5. Create approval queue

#### 19.3.2 Approval Workflow
**Definition:** Sequential process of review and authorization before payment disbursement.

**Levels:**
1. Trial Statement Generation (System)
2. Divisional Head Review (Level 1)
3. Final Statement Generation (System)
4. Payment Processing (System)

#### 19.3.3 Maker-Checker
**Definition:** Dual control mechanism where one user (maker) initiates a transaction and another user (checker) must approve it.

**Applies To:**
- Rate configuration changes
- Manual commission adjustments
- Clawback entries
- Payment file generation
- User role assignments

#### 19.3.4 Audit Trail
**Definition:** Chronological record of all system activities and changes with timestamps and user identification.

**Captured For:**
- All user actions
- Configuration changes
- Data modifications
- Payment transactions
- Access to sensitive data

#### 19.3.5 Role-Based Access Control (RBAC)
**Definition:** Permission management system where access rights are assigned based on user roles.

**Roles:**
- Divisional Head
- Circle Head
- Regional Head
- PLI Directorate
- CEPT
- Audit
- Agent

#### 19.3.6 Data Masking
**Definition:** Hiding sensitive data from unauthorized users while maintaining functionality.

**Masked Fields:**
- PAN (show last 4 characters)
- Bank Account (show last 4 digits)
- Phone (middle digits masked)
- Email (partial masking)

#### 19.3.7 Trial Statement
**Definition:** Preliminary commission statement generated for review before final approval.

**Features:**
- Shows all calculations
- Awaiting approval
- Can be modified
- Approval required for disbursement

#### 19.3.8 Final Statement
**Definition:** Approved commission statement that is locked and ready for disbursement.

**Features:**
- Locked from modifications
- Approved by authority
- Triggers disbursement
- Basis for payment

#### 19.3.9 Disbursement
**Definition:** Process of paying commission to agents.

**Modes:**
- **Cheque:** Physical cheque generated
- **EFT:** Electronic fund transfer via PFMS/Bank

#### 19.3.10 Suspense Account
**Definition:** Temporary holding account for commissions that cannot be immediately disbursed.

**Reasons for Suspense:**
- Payment failures
- Disputed policies
- Investigation pending
- Agent details incomplete

---
**End of Section 19: Glossary and Definitions**

**Previous Section**: [Sample Calculations and Examples](#18-sample-calculations-and-examples)
**Next Section**: [Commission Clawback and Suspense Details](#20-commission-clawback-and-suspense-details)

---

# Incentive, Commission and Producer Management - Commission Clawback and Suspense Details

## 20. Commission Clawback and Suspense Details

### 20.1 Commission Clawback Overview

**Definition:** Recovery of commission already paid due to policy cancellation, lapse, or premium reversal.

**Purpose:** Ensure department does not pay commission on policies that are no longer in-force.

### 20.2 Clawback Triggers

| Trigger | Description | Priority | Timeline |
|---------|-------------|----------|----------|
| **Premium Bounce/Reversal** | Premium payment failed after commission paid | High | Immediate |
| **Policy Cancellation** | Policy cancelled during free-look or early lapse | High | Immediate |
| **Data Entry Error** | Incorrect data led to wrong commission | Medium | Next cycle |
| **Fraudulent Policy** | Policy obtained through fraud | Critical | Immediate |
| **Agent Disciplinary** | Agent terminated due to misconduct | Medium | Next cycle |

### 20.3 Clawback Process Workflow

```
STEP 1: TRIGGER DETECTION
  Detect policy event (lapse, cancellation, bounce)
  Verify commission already paid
  Calculate clawback amount

STEP 2: CREATE CLAWBACK ENTRY
  Generate clawback record
  Link to original commission transaction
  Set status = PENDING_RECOVERY

STEP 3: ADJUST IN NEXT CYCLE
  Calculate next month commission
  Apply clawback as negative entry
  Net payable = current_commission - clawback_amount

STEP 4: HANDLE INSUFFICIENT COMMISSION
  IF current_commission < clawback_amount THEN
    Carry forward balance
    Initiate establishment recovery
  END

STEP 5: COMPLETE RECOVERY
  Mark clawback as COMPLETED
  Update agent balance
  Notify agent of recovery
```

### 20.4 Graduated Recovery Schedule

**When commission is insufficient for full clawback:**

| Remaining Balance | Recovery Action |
|-------------------|-----------------|
| <= Rs. 5,000 | Deduct from next month only |
| Rs. 5,001 - Rs. 25,000 | Max 50% of future commissions |
| Rs. 25,001 - Rs. 50,000 | Max 50% for 3 months, then escalation |
| > Rs. 50,000 | Escalate to establishment recovery |

**Recovery Cap:**
- Maximum 50% of monthly commission can be deducted
- Minimum 50% must be paid to agent (essential income protection)

### 20.5 Clawback Notification

**Email to Agent:**
```
Subject: Commission Recovery Notice - Policy {policy_number}

Dear {agent_name},

This is to inform you that a recovery of Rs. {clawback_amount} has been
initiated against your commission account for the following reason:

Policy Details:
  Policy Number: {policy_number}
  Policy Holder: {policy_holder_name}
  Original Commission: Rs. {original_amount}
  Recovery Amount: Rs. {clawback_amount}
  Reason: {recovery_reason}

Recovery Schedule:
  Amount to be deducted from next commission: Rs. {deduction_amount}
  Balance remaining (if any): Rs. {balance_amount}

For any queries, please contact the Finance Department.

Regards,
Incentive Management System
Postal Life Insurance
```

### 20.6 Commission Suspense Overview

**Definition:** Temporary holding of commission when disbursement cannot be completed immediately.

**Purpose:** Track commissions that are pending resolution of issues.

### 20.7 Suspense Account Categories

| Category | Description | Typical Duration | Resolution |
|----------|-------------|------------------|------------|
| **Payment Failure** | EFT payment failed | 3-7 days | Retry or manual payment |
| **Bank Details Invalid** | Account details incorrect | 30 days | Update details |
| **Policy Under Investigation** | Policy authenticity verification | 30-90 days | Release or forfeit |
| **Agent Inactive** | Agent status changed | Until resolved | Hold or release |
| **Overpayment** | Excess payment made | Until recovery | Recovery from future |

### 20.8 Suspense Account Status Flow

```
┌─────────────┐
│   CREATED   │ ← Suspense entry created
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   PENDING   │ ← Waiting for resolution
└──────┬──────┘
       │
       ├─→→→ [PAYMENT_RETRY]
       │       │
       │       ├─→ SUCCESS → RELEASED
       │       └─→ FAILURE → FAILED
       │
       ├─→→→ [INVESTIGATION]
       │       │
       │       ├─→ CLEARED → RELEASED
       │       ├─→ FRAUD → FORFEITED
       │       └─→ INCONCLUSIVE → EXTENDED
       │
       └─→→→ [RECOVERY]
               │
               └─→ RECOVERED → CLOSED
```

### 20.9 Payment Retry Logic with Exponential Backoff

**Retry Schedule:**

| Attempt | Timing | Action |
|---------|--------|--------|
| 1 | Immediate | First payment attempt |
| 2 | After 1 hour | Retry with same details |
| 3 | After 3 hours | Retry after verification |
| 4 | Manual intervention | Escalate to finance team |

**Exponential Backoff Calculation:**
```
Retry Delay = Base Delay × (2 ^ (Attempt - 1))

Where:
  Base Delay = 1 hour
  Attempt = Retry attempt number (1-based)

Example:
  Attempt 1: 1 × (2 ^ 0) = 1 hour
  Attempt 2: 1 × (2 ^ 1) = 2 hours
  Attempt 3: 1 × (2 ^ 2) = 4 hours
```

### 20.10 Suspense Aging Report

**Report Format:**

| Suspense ID | Agent | Amount | Age (Days) | Category | Status | Action Required |
|-------------|-------|--------|------------|----------|--------|-----------------|
| SUS-001 | DA001 | 5,000 | 5 | Payment Failure | Retrying | Monitor |
| SUS-002 | FO023 | 3,500 | 35 | Bank Invalid | Pending | Contact agent |
| SUS-003 | GDS045 | 2,000 | 60 | Investigation | Under Review | Escalate |
| SUS-004 | DA005 | 8,000 | 95 | Investigation | Forfeited | Write-off |

**Aging Buckets:**
| Age Range | Status | Action |
|-----------|--------|--------|
| 0-7 days | Fresh | Normal processing |
| 8-30 days | Aging | Monitor closely |
| 31-60 days | Old | Escalate to supervisor |
| 61-90 days | Critical | Management review |
| > 90 days | Stale | Write-off consideration |

### 20.11 Suspense Release Process

**When suspense is released:**

```
STEP 1: VERIFICATION
  Confirm issue is resolved
  Verify bank details (if applicable)
  Check agent active status

STEP 2: CALCULATE RELEASE AMOUNT
  Original Amount - Recovery (if any) = Release Amount

STEP 3: UPDATE SUSPENSE RECORD
  Set status = RELEASED
  Record release_date
  Link to payment transaction

STEP 4: PROCESS PAYMENT
  Add to next payment batch
  Or immediate processing if urgent

STEP 5: NOTIFY AGENT
  Send release confirmation
  Include payment details
```

### 20.12 Suspense Forfeiture

**Conditions for Forfeiture:**

| Condition | Period | Action |
|-----------|--------|--------|
| No response from agent | 90 days | Final notice |
| Bank details not provided | 120 days | Hold till provided |
| Fraud confirmed | Immediate | Forfeit + legal action |
| Policy cancelled | 30 days | Forfeit if overpaid |

**Forfeiture Process:**
```
STEP 1: Verify all attempts made
STEP 2: Get management approval
STEP 3: Mark suspense as FORFEITED
STEP 4: Update accounting records
STEP 5: Report to audit
```

### 20.13 Accounting Treatment

**Clawback Accounting:**
```
Dr. Commission Recoverable (Asset)
  Cr. Commission Paid (Expense reduction)

When recovered:
  Dr. Bank/Cash
  Cr. Commission Recoverable
```

**Suspense Accounting:**
```
When created:
  Dr. Commission Suspense (Asset)
    Cr. Commission Payable (Liability reduction)

When released:
  Dr. Commission Payable
    Cr. Bank/Cash
    Cr. Commission Suspense

When forfeited:
  Dr. Commission Suspense
    Cr. Other Income (Gain)
```

---
**End of Section 20: Commission Clawback and Suspense Details**

**Previous Section**: [Glossary and Definitions](#19-glossary-and-definitions)

---

## Appendix A: Quick Reference

### File Structure

```
detailed/
├── IC_00_MAIN.md                  # Executive Summary & TOC
├── IC_01_BUSINESS_RULES.md        # 30+ Business Rules
├── IC_02_FUNCTIONAL_REQUIREMENTS.md  # 32 Functional Requirements
├── IC_03_VALIDATION_RULES.md      # 45+ Validation Rules
├── IC_04_ERROR_CODES.md           # 16 Error Codes
├── IC_05_WORKFLOWS.md             # 8 Workflows
├── IC_06_DATA_ENTITIES.md         # 12 Data Entities with schemas
├── IC_07_INTEGRATION_POINTS.md    # 4 External Integrations
├── IC_08_TEMPORAL_WORKFLOWS.md    # 4 Temporal Workflows with Go code
└── IC_09_TRACEABILITY.md          # Complete Traceability Matrix
```

### Document Metrics

- **Total Pages** (if printed): ~150 pages
- **Total Words**: ~60,000 words
- **Code Examples**: 15+ Go code blocks
- **Tables**: 50+ tables
- **Cross-references**: 100+ links

### Key Statistics

| Category | Count |
|----------|-------|
| Business Rules | 30 |
| Functional Requirements | 32 |
| Validation Rules | 45 |
| Error Codes | 16 |
| Workflows | 8 |
| Temporal Workflows | 4 |
| Data Entities | 12 |
| Integration Points | 4 |
| Commission Types | 3 |
| Agent Types | 4 |
| Payment Modes | 2 |

---

**End of Document**

**Previous Section**: [Integration Points](#8-integration-points)
