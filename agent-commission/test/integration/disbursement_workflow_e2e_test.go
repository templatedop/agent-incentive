package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"agent-commission/client"
	"agent-commission/core/domain"
	"agent-commission/repo/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDisbursementWorkflowE2E tests the complete disbursement workflow
// Covers: FR-IC-COM-009, FR-IC-COM-010, BR-IC-COM-008, BR-IC-COM-011, INT-IC-003, INT-IC-004
func TestDisbursementWorkflowE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Setup test database and clients
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	disbursementRepo := postgres.NewDisbursementRepository(db)

	// Mock external clients
	pfmsClient := setupMockPFMSClient(t)
	accountingClient := setupMockAccountingClient(t)

	t.Run("EFT Disbursement Flow", func(t *testing.T) {
		// Step 1: Create disbursement record
		t.Log("Step 1: Creating EFT disbursement record...")

		agentID := "AGT00001"
		netCommission := 95000.0
		grossCommission := 100000.0
		tdsAmount := 5000.0

		disbursement := &domain.Disbursement{
			AgentID:              agentID,
			StatementID:          12345,
			DisbursementMode:     domain.DisbursementModeEFT,
			DisbursementStatus:   domain.DisbursementStatusPending,
			GrossAmount:          grossCommission,
			TDSAmount:            tdsAmount,
			NetAmount:            netCommission,
			BankAccountNumber:    "1234567890",
			BankIFSC:             "SBIN0001234",
			BankName:             "State Bank of India",
			AccountHolderName:    "Test Agent",
			RequestedBy:          "test_user",
			RequestedAt:          time.Now(),
		}

		err := disbursementRepo.CreateDisbursement(ctx, disbursement)
		require.NoError(t, err, "Failed to create disbursement")
		assert.Greater(t, disbursement.DisbursementID, int64(0), "Disbursement ID should be set")

		// Step 2: Validate bank account
		t.Log("Step 2: Validating bank account...")

		isValid, err := pfmsClient.ValidateBankAccount(
			ctx,
			disbursement.BankAccountNumber,
			disbursement.BankIFSC,
			disbursement.AccountHolderName,
		)
		require.NoError(t, err, "Bank validation failed")
		assert.True(t, isValid, "Bank account should be valid")

		// Step 3: Post accounting voucher
		t.Log("Step 3: Posting accounting voucher...")

		voucherResp, err := accountingClient.PostCommissionVoucher(
			ctx,
			fmt.Sprintf("DISB-%d", disbursement.DisbursementID),
			agentID,
			grossCommission,
			tdsAmount,
			netCommission,
			time.Now(),
		)
		require.NoError(t, err, "Failed to post voucher")
		assert.NotEmpty(t, voucherResp.VoucherID, "Voucher ID should be set")
		assert.Equal(t, "POSTED", voucherResp.Status, "Voucher should be posted")

		disbursement.VoucherID = &voucherResp.VoucherID

		// Step 4: Initiate EFT payment
		t.Log("Step 4: Initiating EFT payment...")

		paymentRequest := &client.EFTPaymentRequest{
			RequestID:          fmt.Sprintf("EFT-%d-%d", disbursement.DisbursementID, time.Now().Unix()),
			PaymentType:        "NEFT",
			Amount:             netCommission,
			Currency:           "INR",
			BeneficiaryName:    disbursement.AccountHolderName,
			BeneficiaryAccount: disbursement.BankAccountNumber,
			BeneficiaryIFSC:    disbursement.BankIFSC,
			BeneficiaryBankName: disbursement.BankName,
			Purpose:            "Commission Payment",
			Remarks:            fmt.Sprintf("Commission payment for Agent %s", agentID),
			TransactionDate:    time.Now(),
			ReferenceType:      "DISBURSEMENT",
			ReferenceNumber:    fmt.Sprintf("DISB-%d", disbursement.DisbursementID),
		}

		paymentResp, err := pfmsClient.InitiateEFTPayment(ctx, paymentRequest)
		require.NoError(t, err, "Failed to initiate EFT payment")
		assert.NotEmpty(t, paymentResp.PaymentID, "Payment ID should be set")
		assert.Contains(t, []string{"INITIATED", "PENDING", "SUCCESS"}, paymentResp.Status)

		disbursement.PaymentID = &paymentResp.PaymentID
		disbursement.DisbursementStatus = domain.DisbursementStatusInProgress

		// Step 5: Update disbursement with payment details
		t.Log("Step 5: Updating disbursement with payment details...")

		err = disbursementRepo.UpdateDisbursementStatus(
			ctx,
			disbursement.DisbursementID,
			disbursement.DisbursementStatus,
		)
		require.NoError(t, err, "Failed to update disbursement status")

		// Step 6: Simulate payment status polling
		t.Log("Step 6: Polling payment status...")

		var paymentStatus *client.PaymentStatusResponse
		maxRetries := 5
		for i := 0; i < maxRetries; i++ {
			paymentStatus, err = pfmsClient.GetPaymentStatus(ctx, *disbursement.PaymentID)
			require.NoError(t, err, "Failed to get payment status")

			if paymentStatus.Status == "SUCCESS" || paymentStatus.Status == "FAILED" {
				break
			}

			time.Sleep(time.Second) // Simulate polling delay
		}

		// Step 7: Process payment confirmation
		t.Log("Step 7: Processing payment confirmation...")

		if paymentStatus.Status == "SUCCESS" {
			disbursement.DisbursementStatus = domain.DisbursementStatusCompleted
			disbursement.TransactionID = &paymentStatus.TransactionID
			completedAt := time.Now()
			disbursement.CompletedAt = &completedAt

			// Verify SLA compliance
			// BR-IC-COM-011: 10-working-day SLA
			slaDeadline := disbursement.RequestedAt.AddDate(0, 0, 10)
			assert.True(t, completedAt.Before(slaDeadline), "Disbursement should complete within 10-day SLA")
		} else {
			disbursement.DisbursementStatus = domain.DisbursementStatusFailed
			failureReason := paymentStatus.FailureReason
			disbursement.FailureReason = &failureReason
		}

		err = disbursementRepo.UpdateDisbursementStatus(
			ctx,
			disbursement.DisbursementID,
			disbursement.DisbursementStatus,
		)
		require.NoError(t, err, "Failed to update final status")

		// Step 8: Verify final state
		t.Log("Step 8: Verifying final state...")

		retrievedDisb, err := disbursementRepo.GetDisbursementByID(ctx, disbursement.DisbursementID)
		require.NoError(t, err, "Failed to retrieve disbursement")
		assert.Equal(t, domain.DisbursementStatusCompleted, retrievedDisb.DisbursementStatus)
		assert.NotNil(t, retrievedDisb.PaymentID, "Payment ID should be set")
		assert.NotNil(t, retrievedDisb.VoucherID, "Voucher ID should be set")
		assert.NotNil(t, retrievedDisb.TransactionID, "Transaction ID should be set")
		assert.NotNil(t, retrievedDisb.CompletedAt, "Completed timestamp should be set")

		t.Log("✓ EFT disbursement E2E test completed successfully")
	})

	t.Run("Cheque Disbursement Flow", func(t *testing.T) {
		// Step 1: Create cheque disbursement
		t.Log("Step 1: Creating cheque disbursement record...")

		agentID := "AGT00002"
		netCommission := 47500.0
		grossCommission := 50000.0
		tdsAmount := 2500.0

		disbursement := &domain.Disbursement{
			AgentID:            agentID,
			StatementID:        12346,
			DisbursementMode:   domain.DisbursementModeCheque,
			DisbursementStatus: domain.DisbursementStatusPending,
			GrossAmount:        grossCommission,
			TDSAmount:          tdsAmount,
			NetAmount:          netCommission,
			PayeeAddress:       "123 Test Street, Test City",
			RequestedBy:        "test_user",
			RequestedAt:        time.Now(),
		}

		err := disbursementRepo.CreateDisbursement(ctx, disbursement)
		require.NoError(t, err, "Failed to create disbursement")

		// Step 2: Post accounting voucher
		t.Log("Step 2: Posting accounting voucher...")

		voucherResp, err := accountingClient.PostCommissionVoucher(
			ctx,
			fmt.Sprintf("DISB-%d", disbursement.DisbursementID),
			agentID,
			grossCommission,
			tdsAmount,
			netCommission,
			time.Now(),
		)
		require.NoError(t, err, "Failed to post voucher")
		disbursement.VoucherID = &voucherResp.VoucherID

		// Step 3: Generate cheque
		t.Log("Step 3: Generating cheque...")

		chequeNumber := fmt.Sprintf("CHQ%010d", time.Now().Unix())
		chequeDate := time.Now()

		disbursement.ChequeNumber = &chequeNumber
		disbursement.ChequeDate = &chequeDate
		disbursement.DisbursementStatus = domain.DisbursementStatusInProgress

		err = disbursementRepo.UpdateDisbursementStatus(
			ctx,
			disbursement.DisbursementID,
			disbursement.DisbursementStatus,
		)
		require.NoError(t, err, "Failed to update status")

		// Step 4: Mark as dispatched
		t.Log("Step 4: Marking cheque as dispatched...")

		dispatchedAt := time.Now()
		disbursement.DispatchedAt = &dispatchedAt
		disbursement.DisbursementStatus = domain.DisbursementStatusDispatched

		err = disbursementRepo.UpdateDisbursementStatus(
			ctx,
			disbursement.DisbursementID,
			disbursement.DisbursementStatus,
		)
		require.NoError(t, err, "Failed to update to dispatched")

		// Step 5: Mark as completed (cheque cleared)
		t.Log("Step 5: Marking as completed...")

		completedAt := time.Now()
		disbursement.CompletedAt = &completedAt
		disbursement.DisbursementStatus = domain.DisbursementStatusCompleted

		err = disbursementRepo.UpdateDisbursementStatus(
			ctx,
			disbursement.DisbursementID,
			disbursement.DisbursementStatus,
		)
		require.NoError(t, err, "Failed to mark as completed")

		// Step 6: Verify final state
		t.Log("Step 6: Verifying final state...")

		retrievedDisb, err := disbursementRepo.GetDisbursementByID(ctx, disbursement.DisbursementID)
		require.NoError(t, err, "Failed to retrieve disbursement")
		assert.Equal(t, domain.DisbursementStatusCompleted, retrievedDisb.DisbursementStatus)
		assert.Equal(t, domain.DisbursementModeCheque, retrievedDisb.DisbursementMode)
		assert.NotNil(t, retrievedDisb.ChequeNumber, "Cheque number should be set")
		assert.NotNil(t, retrievedDisb.ChequeDate, "Cheque date should be set")
		assert.NotNil(t, retrievedDisb.DispatchedAt, "Dispatched timestamp should be set")

		t.Log("✓ Cheque disbursement E2E test completed successfully")
	})

	t.Run("Failed Disbursement with Retry", func(t *testing.T) {
		// Step 1: Create disbursement
		t.Log("Step 1: Creating disbursement that will fail...")

		disbursement := &domain.Disbursement{
			AgentID:              "AGT00003",
			StatementID:          12347,
			DisbursementMode:     domain.DisbursementModeEFT,
			DisbursementStatus:   domain.DisbursementStatusPending,
			GrossAmount:          10000.0,
			TDSAmount:            500.0,
			NetAmount:            9500.0,
			BankAccountNumber:    "INVALID123", // Invalid account
			BankIFSC:             "INVALID",
			BankName:             "Test Bank",
			AccountHolderName:    "Test Agent",
			RequestedBy:          "test_user",
			RequestedAt:          time.Now(),
		}

		err := disbursementRepo.CreateDisbursement(ctx, disbursement)
		require.NoError(t, err, "Failed to create disbursement")

		// Step 2: Attempt payment (will fail)
		t.Log("Step 2: Attempting payment (expected to fail)...")

		// Simulate validation failure
		isValid, _ := pfmsClient.ValidateBankAccount(
			ctx,
			disbursement.BankAccountNumber,
			disbursement.BankIFSC,
			disbursement.AccountHolderName,
		)
		assert.False(t, isValid, "Invalid bank account should fail validation")

		// Step 3: Mark as failed
		failureReason := "Invalid bank account details"
		disbursement.FailureReason = &failureReason
		disbursement.DisbursementStatus = domain.DisbursementStatusFailed

		err = disbursementRepo.UpdateDisbursementStatus(
			ctx,
			disbursement.DisbursementID,
			disbursement.DisbursementStatus,
		)
		require.NoError(t, err, "Failed to update status")

		// Step 4: Retry with corrected details
		t.Log("Step 4: Retrying with corrected bank details...")

		// Update bank details
		disbursement.BankAccountNumber = "9876543210"
		disbursement.BankIFSC = "SBIN0001234"
		disbursement.DisbursementStatus = domain.DisbursementStatusPending
		disbursement.FailureReason = nil

		err = disbursementRepo.UpdateDisbursementStatus(
			ctx,
			disbursement.DisbursementID,
			disbursement.DisbursementStatus,
		)
		require.NoError(t, err, "Failed to reset status for retry")

		t.Log("✓ Failed disbursement with retry test completed")
	})
}

// Mock client setup helpers
func setupMockPFMSClient(t *testing.T) *client.PFMSClient {
	// Return a mock PFMS client for testing
	// In real tests, this would use a test server or mock
	config := client.PFMSConfig{
		BaseURL:   "http://mock-pfms:8080",
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
		OrgCode:   "TEST-ORG",
		Timeout:   10 * time.Second,
	}
	return client.NewPFMSClient(config)
}

func setupMockAccountingClient(t *testing.T) *client.AccountingClient {
	// Return a mock accounting client for testing
	config := client.AccountingConfig{
		BaseURL:     "http://mock-accounting:8080",
		APIKey:      "test-api-key",
		CompanyCode: "TEST-CO",
		Timeout:     10 * time.Second,
	}
	return client.NewAccountingClient(config)
}
