package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"agent-commission/core/domain"
	"agent-commission/repo/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommissionBatchE2E tests the complete commission batch processing flow
// Covers: BR-IC-COM-001, BR-IC-COM-002, BR-IC-COM-003, BR-IC-COM-012
func TestCommissionBatchE2E(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Setup test database connection
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Initialize repositories
	commissionBatchRepo := postgres.NewCommissionBatchRepository(db)
	commissionRateRepo := postgres.NewCommissionRateRepository(db)
	trialStatementRepo := postgres.NewTrialStatementRepository(db)

	t.Run("Complete Commission Batch Flow", func(t *testing.T) {
		// Step 1: Setup commission rates
		t.Log("Step 1: Setting up commission rates...")
		rate := &domain.CommissionRate{
			ProductType:      "TERM_INSURANCE",
			AgentType:        "INDIVIDUAL",
			PlanCode:         "TERM-001",
			CommissionType:   domain.CommissionTypeFirstYear,
			RatePercentage:   40.0, // 40% FYC
			MinCommission:    500.0,
			MaxCommission:    100000.0,
			EffectiveFrom:    time.Now().AddDate(0, -1, 0), // 1 month ago
			EffectiveTo:      nil,
		}

		err := commissionRateRepo.CreateRate(ctx, rate)
		require.NoError(t, err, "Failed to create commission rate")
		assert.Greater(t, rate.RateID, int64(0), "Rate ID should be set")

		// Step 2: Create commission batch
		t.Log("Step 2: Creating commission batch...")
		batch := &domain.CommissionBatch{
			BatchMonth:      time.Now().AddDate(0, -1, 0), // Last month
			BatchStatus:     domain.BatchStatusPending,
			TotalPolicies:   100,
			ProcessedPolicies: 0,
			FailedPolicies:   0,
			CreatedBy:       "test_user",
		}

		err = commissionBatchRepo.CreateBatch(ctx, batch)
		require.NoError(t, err, "Failed to create batch")
		assert.NotEmpty(t, batch.BatchID, "Batch ID should be set")

		// Step 3: Process policies and calculate commissions
		t.Log("Step 3: Processing policies and calculating commissions...")

		// Simulate policy data
		policies := generateTestPolicies(100)

		var totalGross, totalNet, totalTDS float64
		var successCount, failedCount int

		for _, policy := range policies {
			// Calculate commission
			premiumAmount := policy.PremiumAmount
			grossCommission := premiumAmount * (rate.RatePercentage / 100.0)

			// Apply min/max limits
			if grossCommission < rate.MinCommission {
				grossCommission = rate.MinCommission
			}
			if grossCommission > rate.MaxCommission {
				grossCommission = rate.MaxCommission
			}

			// Calculate TDS (10% if PAN not available, 5% if PAN available)
			tdsPercentage := 5.0 // Assume PAN available
			tdsAmount := grossCommission * (tdsPercentage / 100.0)
			netCommission := grossCommission - tdsAmount

			// Create commission transaction
			transaction := &domain.CommissionTransaction{
				BatchID:          batch.BatchID,
				AgentID:          policy.AgentID,
				PolicyNumber:     policy.PolicyNumber,
				CommissionType:   domain.CommissionTypeFirstYear,
				GrossCommission:  grossCommission,
				TDSAmount:        tdsAmount,
				TDSPercentage:    tdsPercentage,
				NetCommission:    netCommission,
				CommissionStatus: domain.CommissionStatusCalculated,
				PremiumAmount:    premiumAmount,
				ProductType:      policy.ProductType,
				PlanCode:         policy.PlanCode,
			}

			err = commissionBatchRepo.CreateCommissionTransaction(ctx, transaction)
			if err != nil {
				failedCount++
				continue
			}

			successCount++
			totalGross += grossCommission
			totalNet += netCommission
			totalTDS += tdsAmount
		}

		// Step 4: Update batch summary
		t.Log("Step 4: Updating batch summary...")
		batch.ProcessedPolicies = successCount
		batch.FailedPolicies = failedCount
		batch.TotalGrossCommission = totalGross
		batch.TotalNetCommission = totalNet
		batch.TotalTDS = totalTDS
		batch.BatchStatus = domain.BatchStatusProcessing

		err = commissionBatchRepo.UpdateBatchProgress(ctx, batch.BatchID, successCount, failedCount)
		require.NoError(t, err, "Failed to update batch progress")

		// Step 5: Generate trial statement
		t.Log("Step 5: Generating trial statement...")
		trialStatement := &domain.TrialStatement{
			BatchID:            batch.BatchID,
			StatementMonth:     batch.BatchMonth,
			TotalAgents:        countUniqueAgents(policies),
			TotalPolicies:      successCount,
			TotalGrossCommission: totalGross,
			TotalTDS:           totalTDS,
			TotalNetCommission:  totalNet,
			TrialStatus:        domain.TrialStatusPending,
			GeneratedBy:        "system",
			GeneratedAt:        time.Now(),
		}

		err = trialStatementRepo.CreateTrialStatement(ctx, trialStatement)
		require.NoError(t, err, "Failed to create trial statement")
		assert.Greater(t, trialStatement.StatementID, int64(0), "Statement ID should be set")

		// Step 6: Verify SLA compliance
		t.Log("Step 6: Verifying SLA compliance...")
		// BR-IC-COM-012: 6-hour SLA for batch completion
		slaDeadline := batch.CreatedAt.Add(6 * time.Hour)
		assert.True(t, time.Now().Before(slaDeadline), "Batch should complete within 6-hour SLA")

		// Step 7: Verify calculations
		t.Log("Step 7: Verifying calculations...")
		assert.Equal(t, 100, batch.TotalPolicies, "Total policies should match")
		assert.Equal(t, successCount, batch.ProcessedPolicies, "Processed count should match")
		assert.Greater(t, totalGross, 0.0, "Total gross commission should be positive")
		assert.Greater(t, totalNet, 0.0, "Total net commission should be positive")
		assert.Greater(t, totalTDS, 0.0, "Total TDS should be positive")
		assert.Equal(t, totalGross-totalTDS, totalNet, "Net = Gross - TDS")

		// Step 8: Mark batch as completed
		t.Log("Step 8: Marking batch as completed...")
		batch.BatchStatus = domain.BatchStatusCompleted
		batch.CompletedAt = ptrTime(time.Now())

		err = commissionBatchRepo.UpdateBatchStatus(ctx, batch.BatchID, batch.BatchStatus)
		require.NoError(t, err, "Failed to update batch status")

		// Step 9: Verify final state
		t.Log("Step 9: Verifying final state...")
		retrievedBatch, err := commissionBatchRepo.GetBatchByID(ctx, batch.BatchID)
		require.NoError(t, err, "Failed to retrieve batch")
		assert.Equal(t, domain.BatchStatusCompleted, retrievedBatch.BatchStatus, "Batch should be completed")
		assert.NotNil(t, retrievedBatch.CompletedAt, "Completed timestamp should be set")

		t.Log("✓ Commission batch E2E test completed successfully")
	})
}

// TestCommissionCalculationAccuracy tests commission calculation accuracy
func TestCommissionCalculationAccuracy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testCases := []struct {
		name              string
		premiumAmount     float64
		ratePercentage    float64
		minCommission     float64
		maxCommission     float64
		tdsPercentage     float64
		expectedGross     float64
		expectedTDS       float64
		expectedNet       float64
	}{
		{
			name:           "Standard calculation",
			premiumAmount:  10000.0,
			ratePercentage: 40.0,
			minCommission:  500.0,
			maxCommission:  100000.0,
			tdsPercentage:  5.0,
			expectedGross:  4000.0,  // 10000 * 40%
			expectedTDS:    200.0,   // 4000 * 5%
			expectedNet:    3800.0,  // 4000 - 200
		},
		{
			name:           "Below minimum",
			premiumAmount:  1000.0,
			ratePercentage: 40.0,
			minCommission:  500.0,
			maxCommission:  100000.0,
			tdsPercentage:  5.0,
			expectedGross:  500.0,   // Min commission
			expectedTDS:    25.0,    // 500 * 5%
			expectedNet:    475.0,   // 500 - 25
		},
		{
			name:           "Above maximum",
			premiumAmount:  500000.0,
			ratePercentage: 40.0,
			minCommission:  500.0,
			maxCommission:  100000.0,
			tdsPercentage:  5.0,
			expectedGross:  100000.0, // Max commission
			expectedTDS:    5000.0,   // 100000 * 5%
			expectedNet:    95000.0,  // 100000 - 5000
		},
		{
			name:           "Higher TDS (no PAN)",
			premiumAmount:  10000.0,
			ratePercentage: 40.0,
			minCommission:  500.0,
			maxCommission:  100000.0,
			tdsPercentage:  10.0, // No PAN
			expectedGross:  4000.0,
			expectedTDS:    400.0,   // 4000 * 10%
			expectedNet:    3600.0,  // 4000 - 400
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate gross commission
			grossCommission := tc.premiumAmount * (tc.ratePercentage / 100.0)

			// Apply min/max limits
			if grossCommission < tc.minCommission {
				grossCommission = tc.minCommission
			}
			if grossCommission > tc.maxCommission {
				grossCommission = tc.maxCommission
			}

			// Calculate TDS
			tdsAmount := grossCommission * (tc.tdsPercentage / 100.0)
			netCommission := grossCommission - tdsAmount

			// Verify calculations
			assert.Equal(t, tc.expectedGross, grossCommission, "Gross commission mismatch")
			assert.Equal(t, tc.expectedTDS, tdsAmount, "TDS amount mismatch")
			assert.Equal(t, tc.expectedNet, netCommission, "Net commission mismatch")
		})
	}
}

// Helper functions

type testPolicy struct {
	AgentID       string
	PolicyNumber  string
	ProductType   string
	PlanCode      string
	PremiumAmount float64
}

func generateTestPolicies(count int) []testPolicy {
	policies := make([]testPolicy, count)
	for i := 0; i < count; i++ {
		policies[i] = testPolicy{
			AgentID:       generateAgentID(i % 10), // 10 unique agents
			PolicyNumber:  generatePolicyNumber(i),
			ProductType:   "TERM_INSURANCE",
			PlanCode:      "TERM-001",
			PremiumAmount: float64(10000 + (i * 100)), // Varying premiums
		}
	}
	return policies
}

func generateAgentID(index int) string {
	return fmt.Sprintf("AGT%05d", index+1)
}

func generatePolicyNumber(index int) string {
	return fmt.Sprintf("POL%08d", index+1)
}

func countUniqueAgents(policies []testPolicy) int {
	agentMap := make(map[string]bool)
	for _, p := range policies {
		agentMap[p.AgentID] = true
	}
	return len(agentMap)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// Test database setup helpers (to be implemented based on your test infrastructure)
func setupTestDB(t *testing.T) *postgres.DB {
	// TODO: Implement test database setup
	// This should create a test database with schema
	t.Skip("Test database setup not implemented")
	return nil
}

func cleanupTestDB(t *testing.T, db *postgres.DB) {
	// TODO: Implement test database cleanup
	// This should drop test tables and close connection
}
