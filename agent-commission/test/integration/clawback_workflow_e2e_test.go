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

// TestClawbackWorkflowE2E tests the complete clawback workflow
// Covers: FR-IC-COM-011, BR-IC-CLAWBACK-001 (Graduated recovery), WF-IC-CLAWBACK-001
func TestClawbackWorkflowE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Setup test database and clients
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	clawbackRepo := postgres.NewClawbackRepository(db)
	commissionBatchRepo := postgres.NewCommissionBatchRepository(db)
	accountingClient := setupMockAccountingClient(t)
	policyClient := setupMockPolicyClient(t)

	t.Run("Graduated Clawback Recovery - Year 1", func(t *testing.T) {
		// Step 1: Setup - Create original commission transaction
		t.Log("Step 1: Creating original commission transaction...")

		policyNumber := "POL00000001"
		agentID := "AGT00001"
		originalCommission := 100000.0

		transaction := &domain.CommissionTransaction{
			BatchID:          "BATCH-2025-01",
			AgentID:          agentID,
			PolicyNumber:     policyNumber,
			CommissionType:   domain.CommissionTypeFirstYear,
			GrossCommission:  originalCommission,
			TDSAmount:        5000.0,
			TDSPercentage:    5.0,
			NetCommission:    95000.0,
			CommissionStatus: domain.CommissionStatusPaid,
			PremiumAmount:    250000.0,
			ProductType:      "TERM_INSURANCE",
			PlanCode:         "TERM-001",
		}

		err := commissionBatchRepo.CreateCommissionTransaction(ctx, transaction)
		require.NoError(t, err, "Failed to create commission transaction")

		// Step 2: Policy status change (triggers clawback)
		t.Log("Step 2: Simulating policy lapse (within 12 months)...")

		policyInceptionDate := time.Now().AddDate(0, -6, 0) // 6 months old policy
		lapseDate := time.Now()
		policyAgeMonths := 6

		// Step 3: Calculate clawback amount (100% for Year 1)
		t.Log("Step 3: Calculating clawback amount...")

		clawbackPercentage := domain.CalculateClawbackPercentage(policyAgeMonths)
		assert.Equal(t, 100.0, clawbackPercentage, "Year 1 should have 100% clawback")

		clawbackAmount := originalCommission * (clawbackPercentage / 100.0)
		assert.Equal(t, originalCommission, clawbackAmount, "Full amount should be clawed back")

		// Step 4: Create clawback record
		t.Log("Step 4: Creating clawback record...")

		clawback := &domain.Clawback{
			PolicyNumber:        policyNumber,
			AgentID:             agentID,
			OriginalCommission:  originalCommission,
			ClawbackPercentage:  clawbackPercentage,
			ClawbackAmount:      clawbackAmount,
			ClawbackReason:      domain.ClawbackReasonPolicyLapsed,
			ClawbackStatus:      domain.ClawbackStatusInitiated,
			PolicyInceptionDate: policyInceptionDate,
			LapseDate:           &lapseDate,
			InitiatedBy:         "system",
			InitiatedAt:         time.Now(),
		}

		err = clawbackRepo.CreateClawback(ctx, clawback)
		require.NoError(t, err, "Failed to create clawback")
		assert.Greater(t, clawback.ClawbackID, int64(0), "Clawback ID should be set")

		// Step 5: Post accounting entry for clawback
		t.Log("Step 5: Posting accounting entry...")

		voucherResp, err := accountingClient.PostClawbackVoucher(
			ctx,
			clawback.ClawbackID,
			agentID,
			clawbackAmount,
			time.Now(),
		)
		require.NoError(t, err, "Failed to post clawback voucher")
		assert.Equal(t, "POSTED", voucherResp.Status, "Voucher should be posted")

		clawback.VoucherID = &voucherResp.VoucherID

		// Step 6: Record recovery
		t.Log("Step 6: Recording recovery...")

		recovery := &domain.ClawbackRecovery{
			ClawbackID:      clawback.ClawbackID,
			RecoveryDate:    time.Now(),
			RecoveryAmount:  clawbackAmount,
			RecoveryMethod:  "DEDUCTION_FROM_FUTURE_COMMISSION",
			RecoveryStatus:  "COMPLETED",
		}

		err = clawbackRepo.RecordRecovery(ctx, recovery, true)
		require.NoError(t, err, "Failed to record recovery")

		// Step 7: Update clawback status
		clawback.ClawbackStatus = domain.ClawbackStatusRecovered
		completedAt := time.Now()
		clawback.CompletedAt = &completedAt

		err = clawbackRepo.UpdateClawbackStatus(ctx, clawback.ClawbackID, clawback.ClawbackStatus)
		require.NoError(t, err, "Failed to update clawback status")

		// Step 8: Verify final state
		t.Log("Step 8: Verifying clawback completion...")

		retrievedClawback, err := clawbackRepo.GetClawbackByID(ctx, clawback.ClawbackID)
		require.NoError(t, err, "Failed to retrieve clawback")
		assert.Equal(t, domain.ClawbackStatusRecovered, retrievedClawback.ClawbackStatus)
		assert.Equal(t, 100.0, retrievedClawback.ClawbackPercentage)
		assert.Equal(t, originalCommission, retrievedClawback.ClawbackAmount)

		t.Log("✓ Year 1 graduated clawback test completed")
	})

	t.Run("Graduated Clawback Recovery - Year 2", func(t *testing.T) {
		t.Log("Testing Year 2 clawback (75% recovery)...")

		policyAgeMonths := 18 // 1.5 years
		originalCommission := 80000.0

		clawbackPercentage := domain.CalculateClawbackPercentage(policyAgeMonths)
		assert.Equal(t, 75.0, clawbackPercentage, "Year 2 should have 75% clawback")

		clawbackAmount := originalCommission * (clawbackPercentage / 100.0)
		assert.Equal(t, 60000.0, clawbackAmount, "75% of commission should be clawed back")

		t.Log("✓ Year 2 clawback calculation verified")
	})

	t.Run("Graduated Clawback Recovery - Year 3", func(t *testing.T) {
		t.Log("Testing Year 3 clawback (50% recovery)...")

		policyAgeMonths := 30 // 2.5 years
		originalCommission := 60000.0

		clawbackPercentage := domain.CalculateClawbackPercentage(policyAgeMonths)
		assert.Equal(t, 50.0, clawbackPercentage, "Year 3 should have 50% clawback")

		clawbackAmount := originalCommission * (clawbackPercentage / 100.0)
		assert.Equal(t, 30000.0, clawbackAmount, "50% of commission should be clawed back")

		t.Log("✓ Year 3 clawback calculation verified")
	})

	t.Run("Graduated Clawback Recovery - Year 4", func(t *testing.T) {
		t.Log("Testing Year 4 clawback (25% recovery)...")

		policyAgeMonths := 42 // 3.5 years
		originalCommission := 50000.0

		clawbackPercentage := domain.CalculateClawbackPercentage(policyAgeMonths)
		assert.Equal(t, 25.0, clawbackPercentage, "Year 4 should have 25% clawback")

		clawbackAmount := originalCommission * (clawbackPercentage / 100.0)
		assert.Equal(t, 12500.0, clawbackAmount, "25% of commission should be clawed back")

		t.Log("✓ Year 4 clawback calculation verified")
	})

	t.Run("Graduated Clawback Recovery - Year 5+", func(t *testing.T) {
		t.Log("Testing Year 5+ clawback (0% recovery)...")

		policyAgeMonths := 60 // 5 years
		originalCommission := 40000.0

		clawbackPercentage := domain.CalculateClawbackPercentage(policyAgeMonths)
		assert.Equal(t, 0.0, clawbackPercentage, "Year 5+ should have 0% clawback")

		clawbackAmount := originalCommission * (clawbackPercentage / 100.0)
		assert.Equal(t, 0.0, clawbackAmount, "No clawback after 5 years")

		t.Log("✓ Year 5+ no clawback verified")
	})

	t.Run("Partial Recovery with Multiple Installments", func(t *testing.T) {
		// Step 1: Create clawback with large amount
		t.Log("Step 1: Creating clawback for partial recovery...")

		policyNumber := "POL00000002"
		agentID := "AGT00002"
		originalCommission := 200000.0
		policyAgeMonths := 3 // Year 1
		clawbackPercentage := domain.CalculateClawbackPercentage(policyAgeMonths)
		clawbackAmount := originalCommission * (clawbackPercentage / 100.0) // 200,000

		clawback := &domain.Clawback{
			PolicyNumber:       policyNumber,
			AgentID:            agentID,
			OriginalCommission: originalCommission,
			ClawbackPercentage: clawbackPercentage,
			ClawbackAmount:     clawbackAmount,
			ClawbackReason:     domain.ClawbackReasonPolicySurrendered,
			ClawbackStatus:     domain.ClawbackStatusInitiated,
			PolicyInceptionDate: time.Now().AddDate(0, -3, 0),
			InitiatedBy:        "system",
			InitiatedAt:        time.Now(),
		}

		err := clawbackRepo.CreateClawback(ctx, clawback)
		require.NoError(t, err, "Failed to create clawback")

		// Step 2: Record partial recoveries
		t.Log("Step 2: Recording partial recoveries...")

		recoveries := []float64{50000.0, 50000.0, 50000.0, 50000.0} // 4 installments
		totalRecovered := 0.0

		for i, amount := range recoveries {
			recovery := &domain.ClawbackRecovery{
				ClawbackID:     clawback.ClawbackID,
				RecoveryDate:   time.Now().AddDate(0, i, 0),
				RecoveryAmount: amount,
				RecoveryMethod: "DEDUCTION_FROM_FUTURE_COMMISSION",
				RecoveryStatus: "COMPLETED",
			}

			isFullyRecovered := (totalRecovered + amount >= clawbackAmount)
			err = clawbackRepo.RecordRecovery(ctx, recovery, isFullyRecovered)
			require.NoError(t, err, fmt.Sprintf("Failed to record recovery %d", i+1))

			totalRecovered += amount
			t.Logf("Installment %d: Recovered %.2f, Total: %.2f", i+1, amount, totalRecovered)
		}

		// Step 3: Verify full recovery
		assert.Equal(t, clawbackAmount, totalRecovered, "Total recovered should match clawback amount")

		// Retrieve and verify status
		retrievedClawback, err := clawbackRepo.GetClawbackByID(ctx, clawback.ClawbackID)
		require.NoError(t, err, "Failed to retrieve clawback")
		assert.Equal(t, domain.ClawbackStatusRecovered, retrievedClawback.ClawbackStatus)

		// Get all recoveries
		recoveryList, err := clawbackRepo.GetRecoveriesByClawbackID(ctx, clawback.ClawbackID)
		require.NoError(t, err, "Failed to get recoveries")
		assert.Equal(t, 4, len(recoveryList), "Should have 4 recovery records")

		t.Log("✓ Partial recovery test completed")
	})

	t.Run("Policy Revival - Reverse Clawback", func(t *testing.T) {
		t.Log("Testing policy revival scenario...")

		// Step 1: Create initial clawback
		policyNumber := "POL00000003"
		agentID := "AGT00003"
		originalCommission := 50000.0
		clawbackAmount := 50000.0 // 100% Year 1

		clawback := &domain.Clawback{
			PolicyNumber:        policyNumber,
			AgentID:             agentID,
			OriginalCommission:  originalCommission,
			ClawbackPercentage:  100.0,
			ClawbackAmount:      clawbackAmount,
			ClawbackReason:      domain.ClawbackReasonPolicyLapsed,
			ClawbackStatus:      domain.ClawbackStatusRecovered,
			PolicyInceptionDate: time.Now().AddDate(0, -6, 0),
			InitiatedBy:         "system",
			InitiatedAt:         time.Now().AddDate(0, 0, -30),
		}

		err := clawbackRepo.CreateClawback(ctx, clawback)
		require.NoError(t, err, "Failed to create clawback")

		// Step 2: Policy revived - reverse clawback
		t.Log("Policy revived, reversing clawback...")

		clawback.ClawbackStatus = domain.ClawbackStatusReversed
		reversedAt := time.Now()
		clawback.CompletedAt = &reversedAt

		err = clawbackRepo.UpdateClawbackStatus(ctx, clawback.ClawbackID, clawback.ClawbackStatus)
		require.NoError(t, err, "Failed to reverse clawback")

		// Verify reversal
		retrievedClawback, err := clawbackRepo.GetClawbackByID(ctx, clawback.ClawbackID)
		require.NoError(t, err, "Failed to retrieve clawback")
		assert.Equal(t, domain.ClawbackStatusReversed, retrievedClawback.ClawbackStatus)

		t.Log("✓ Policy revival and clawback reversal completed")
	})
}

// TestClawbackCalculationAccuracy tests clawback calculation accuracy
func TestClawbackCalculationAccuracy(t *testing.T) {
	testCases := []struct {
		name               string
		policyAgeMonths    int
		originalCommission float64
		expectedPercentage float64
		expectedAmount     float64
	}{
		{
			name:               "Month 1 - 100% recovery",
			policyAgeMonths:    1,
			originalCommission: 100000.0,
			expectedPercentage: 100.0,
			expectedAmount:     100000.0,
		},
		{
			name:               "Month 11 - 100% recovery",
			policyAgeMonths:    11,
			originalCommission: 80000.0,
			expectedPercentage: 100.0,
			expectedAmount:     80000.0,
		},
		{
			name:               "Month 13 - 75% recovery",
			policyAgeMonths:    13,
			originalCommission: 100000.0,
			expectedPercentage: 75.0,
			expectedAmount:     75000.0,
		},
		{
			name:               "Month 23 - 75% recovery",
			policyAgeMonths:    23,
			originalCommission: 60000.0,
			expectedPercentage: 75.0,
			expectedAmount:     45000.0,
		},
		{
			name:               "Month 25 - 50% recovery",
			policyAgeMonths:    25,
			originalCommission: 100000.0,
			expectedPercentage: 50.0,
			expectedAmount:     50000.0,
		},
		{
			name:               "Month 35 - 50% recovery",
			policyAgeMonths:    35,
			originalCommission: 80000.0,
			expectedPercentage: 50.0,
			expectedAmount:     40000.0,
		},
		{
			name:               "Month 37 - 25% recovery",
			policyAgeMonths:    37,
			originalCommission: 100000.0,
			expectedPercentage: 25.0,
			expectedAmount:     25000.0,
		},
		{
			name:               "Month 47 - 25% recovery",
			policyAgeMonths:    47,
			originalCommission: 60000.0,
			expectedPercentage: 25.0,
			expectedAmount:     15000.0,
		},
		{
			name:               "Month 48 - 0% recovery",
			policyAgeMonths:    48,
			originalCommission: 100000.0,
			expectedPercentage: 0.0,
			expectedAmount:     0.0,
		},
		{
			name:               "Month 60 - 0% recovery",
			policyAgeMonths:    60,
			originalCommission: 50000.0,
			expectedPercentage: 0.0,
			expectedAmount:     0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			percentage := domain.CalculateClawbackPercentage(tc.policyAgeMonths)
			amount := tc.originalCommission * (percentage / 100.0)

			assert.Equal(t, tc.expectedPercentage, percentage, "Percentage mismatch")
			assert.Equal(t, tc.expectedAmount, amount, "Amount mismatch")
		})
	}
}

// Mock policy client setup
func setupMockPolicyClient(t *testing.T) *client.PolicyClient {
	config := client.PolicyServiceConfig{
		BaseURL: "http://mock-policy:8080",
		APIKey:  "test-api-key",
		Timeout: 10 * time.Second,
	}
	return client.NewPolicyClient(config)
}
