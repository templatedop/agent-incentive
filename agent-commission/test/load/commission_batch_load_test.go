package load

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"agent-commission/core/domain"
	"agent-commission/repo/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommissionBatchLoadTest tests commission processing for 10,000 policies
// Performance requirements:
// - Process 10,000 policies within 6-hour SLA
// - Use batch optimization for database operations
// - Target: <100ms per policy calculation
func TestCommissionBatchLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	ctx := context.Background()

	// Setup
	db := setupLoadTestDB(t)
	defer cleanupLoadTestDB(t, db)

	commissionBatchRepo := postgres.NewCommissionBatchRepository(db)
	commissionRateRepo := postgres.NewCommissionRateRepository(db)

	// Configure load test parameters
	const (
		totalPolicies    = 10000
		batchSize        = 100 // Process in batches of 100
		concurrentWorkers = 10
	)

	t.Run("Sequential Processing - 10,000 Policies", func(t *testing.T) {
		t.Log("Starting sequential processing test...")
		startTime := time.Now()

		// Setup commission rate
		rate := setupCommissionRate(t, ctx, commissionRateRepo)

		// Create batch
		batch := createTestBatch(t, ctx, commissionBatchRepo, totalPolicies)

		// Generate policy data
		policies := generateLoadTestPolicies(totalPolicies)

		// Process policies sequentially
		successCount := 0
		failedCount := 0
		var totalGross, totalNet, totalTDS float64

		for i := 0; i < len(policies); i += batchSize {
			end := i + batchSize
			if end > len(policies) {
				end = len(policies)
			}

			policyBatch := policies[i:end]

			for _, policy := range policyBatch {
				transaction, err := calculateCommission(policy, rate, batch.BatchID)
				if err != nil {
					failedCount++
					continue
				}

				err = commissionBatchRepo.CreateCommissionTransaction(ctx, transaction)
				if err != nil {
					failedCount++
					continue
				}

				successCount++
				totalGross += transaction.GrossCommission
				totalNet += transaction.NetCommission
				totalTDS += transaction.TDSAmount
			}

			// Log progress
			t.Logf("Processed %d/%d policies...", end, totalPolicies)
		}

		duration := time.Since(startTime)
		avgTimePerPolicy := duration / time.Duration(totalPolicies)

		// Assertions
		assert.Equal(t, totalPolicies, successCount, "All policies should be processed")
		assert.Equal(t, 0, failedCount, "No policies should fail")
		assert.True(t, duration < 6*time.Hour, "Should complete within 6-hour SLA")

		t.Logf("✓ Sequential Processing Results:")
		t.Logf("  Total Policies: %d", totalPolicies)
		t.Logf("  Success: %d", successCount)
		t.Logf("  Failed: %d", failedCount)
		t.Logf("  Total Duration: %s", duration)
		t.Logf("  Avg Time Per Policy: %s", avgTimePerPolicy)
		t.Logf("  Total Gross Commission: %.2f", totalGross)
		t.Logf("  Total Net Commission: %.2f", totalNet)
		t.Logf("  Total TDS: %.2f", totalTDS)
	})

	t.Run("Concurrent Processing - 10,000 Policies", func(t *testing.T) {
		t.Log("Starting concurrent processing test...")
		startTime := time.Now()

		// Setup commission rate
		rate := setupCommissionRate(t, ctx, commissionRateRepo)

		// Create batch
		batch := createTestBatch(t, ctx, commissionBatchRepo, totalPolicies)

		// Generate policy data
		policies := generateLoadTestPolicies(totalPolicies)

		// Concurrent processing with worker pool
		var successCount, failedCount int32
		var totalGross, totalNet, totalTDS float64
		var mu sync.Mutex // Protect float additions

		// Create work channel
		workChan := make(chan loadTestPolicy, batchSize)
		var wg sync.WaitGroup

		// Start workers
		for w := 0; w < concurrentWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				for policy := range workChan {
					transaction, err := calculateCommission(policy, rate, batch.BatchID)
					if err != nil {
						atomic.AddInt32(&failedCount, 1)
						continue
					}

					err = commissionBatchRepo.CreateCommissionTransaction(ctx, transaction)
					if err != nil {
						atomic.AddInt32(&failedCount, 1)
						continue
					}

					atomic.AddInt32(&successCount, 1)

					// Update totals (with mutex for float operations)
					mu.Lock()
					totalGross += transaction.GrossCommission
					totalNet += transaction.NetCommission
					totalTDS += transaction.TDSAmount
					mu.Unlock()
				}
			}(w)
		}

		// Send work to workers
		for i, policy := range policies {
			workChan <- policy
			if (i+1)%1000 == 0 {
				t.Logf("Queued %d/%d policies...", i+1, totalPolicies)
			}
		}
		close(workChan)

		// Wait for all workers to complete
		wg.Wait()

		duration := time.Since(startTime)
		avgTimePerPolicy := duration / time.Duration(totalPolicies)
		throughput := float64(totalPolicies) / duration.Seconds()

		// Assertions
		assert.Equal(t, int32(totalPolicies), successCount, "All policies should be processed")
		assert.Equal(t, int32(0), failedCount, "No policies should fail")
		assert.True(t, duration < 6*time.Hour, "Should complete within 6-hour SLA")
		assert.True(t, avgTimePerPolicy < 100*time.Millisecond, "Avg time per policy should be < 100ms")

		t.Logf("✓ Concurrent Processing Results:")
		t.Logf("  Total Policies: %d", totalPolicies)
		t.Logf("  Concurrent Workers: %d", concurrentWorkers)
		t.Logf("  Success: %d", successCount)
		t.Logf("  Failed: %d", failedCount)
		t.Logf("  Total Duration: %s", duration)
		t.Logf("  Avg Time Per Policy: %s", avgTimePerPolicy)
		t.Logf("  Throughput: %.2f policies/second", throughput)
		t.Logf("  Total Gross Commission: %.2f", totalGross)
		t.Logf("  Total Net Commission: %.2f", totalNet)
		t.Logf("  Total TDS: %.2f", totalTDS)
	})

	t.Run("Batch-Optimized Processing - 10,000 Policies", func(t *testing.T) {
		t.Log("Starting batch-optimized processing test...")
		startTime := time.Now()

		// Setup commission rate
		rate := setupCommissionRate(t, ctx, commissionRateRepo)

		// Create batch
		batch := createTestBatch(t, ctx, commissionBatchRepo, totalPolicies)

		// Generate policy data
		policies := generateLoadTestPolicies(totalPolicies)

		// Process using pgx.Batch optimization
		successCount := 0
		var totalGross, totalNet, totalTDS float64

		for i := 0; i < len(policies); i += batchSize {
			end := i + batchSize
			if end > len(policies) {
				end = len(policies)
			}

			policyBatch := policies[i:end]

			// TODO: Use actual pgx.Batch implementation from repository
			// For now, simulate batch processing
			for _, policy := range policyBatch {
				transaction, err := calculateCommission(policy, rate, batch.BatchID)
				if err != nil {
					continue
				}

				// In real implementation, this would use:
				// dblib.QueueReturnRow(batch, insertQuery, scanFn, &transaction)
				err = commissionBatchRepo.CreateCommissionTransaction(ctx, transaction)
				if err != nil {
					continue
				}

				successCount++
				totalGross += transaction.GrossCommission
				totalNet += transaction.NetCommission
				totalTDS += transaction.TDSAmount
			}

			if (i+batchSize)%1000 == 0 {
				t.Logf("Batch-processed %d/%d policies...", i+batchSize, totalPolicies)
			}
		}

		duration := time.Since(startTime)
		avgTimePerPolicy := duration / time.Duration(totalPolicies)
		throughput := float64(totalPolicies) / duration.Seconds()

		// Assertions
		assert.Equal(t, totalPolicies, successCount, "All policies should be processed")
		assert.True(t, duration < 6*time.Hour, "Should complete within 6-hour SLA")

		t.Logf("✓ Batch-Optimized Processing Results:")
		t.Logf("  Total Policies: %d", totalPolicies)
		t.Logf("  Batch Size: %d", batchSize)
		t.Logf("  Success: %d", successCount)
		t.Logf("  Total Duration: %s", duration)
		t.Logf("  Avg Time Per Policy: %s", avgTimePerPolicy)
		t.Logf("  Throughput: %.2f policies/second", throughput)
		t.Logf("  Total Gross Commission: %.2f", totalGross)
		t.Logf("  Total Net Commission: %.2f", totalNet)
		t.Logf("  Total TDS: %.2f", totalTDS)
	})
}

// Helper types and functions

type loadTestPolicy struct {
	AgentID       string
	PolicyNumber  string
	ProductType   string
	PlanCode      string
	PremiumAmount float64
	PolicyType    string
}

func generateLoadTestPolicies(count int) []loadTestPolicy {
	policies := make([]loadTestPolicy, count)
	agents := 100 // Distribute across 100 agents

	for i := 0; i < count; i++ {
		agentIndex := i % agents
		policyType := "FIRST_YEAR"
		if i%5 == 0 { // 20% renewal policies
			policyType = "RENEWAL"
		}

		policies[i] = loadTestPolicy{
			AgentID:       fmt.Sprintf("AGT%05d", agentIndex+1),
			PolicyNumber:  fmt.Sprintf("POL%08d", i+1),
			ProductType:   "TERM_INSURANCE",
			PlanCode:      "TERM-001",
			PremiumAmount: float64(5000 + (i%100)*1000), // Varying premiums: 5000-105000
			PolicyType:    policyType,
		}
	}

	return policies
}

func setupCommissionRate(
	t *testing.T,
	ctx context.Context,
	repo *postgres.CommissionRateRepository,
) *domain.CommissionRate {
	rate := &domain.CommissionRate{
		ProductType:      "TERM_INSURANCE",
		AgentType:        "INDIVIDUAL",
		PlanCode:         "TERM-001",
		CommissionType:   domain.CommissionTypeFirstYear,
		RatePercentage:   40.0,
		MinCommission:    500.0,
		MaxCommission:    100000.0,
		EffectiveFrom:    time.Now().AddDate(0, -1, 0),
		EffectiveTo:      nil,
	}

	err := repo.CreateRate(ctx, rate)
	require.NoError(t, err, "Failed to create commission rate")
	return rate
}

func createTestBatch(
	t *testing.T,
	ctx context.Context,
	repo *postgres.CommissionBatchRepository,
	totalPolicies int,
) *domain.CommissionBatch {
	batch := &domain.CommissionBatch{
		BatchMonth:      time.Now().AddDate(0, -1, 0),
		BatchStatus:     domain.BatchStatusProcessing,
		TotalPolicies:   totalPolicies,
		ProcessedPolicies: 0,
		FailedPolicies:   0,
		CreatedBy:       "load_test",
	}

	err := repo.CreateBatch(ctx, batch)
	require.NoError(t, err, "Failed to create batch")
	return batch
}

func calculateCommission(
	policy loadTestPolicy,
	rate *domain.CommissionRate,
	batchID string,
) (*domain.CommissionTransaction, error) {
	// Calculate gross commission
	grossCommission := policy.PremiumAmount * (rate.RatePercentage / 100.0)

	// Apply min/max limits
	if grossCommission < rate.MinCommission {
		grossCommission = rate.MinCommission
	}
	if grossCommission > rate.MaxCommission {
		grossCommission = rate.MaxCommission
	}

	// Calculate TDS
	tdsPercentage := 5.0 // Assume PAN available
	if policy.AgentID[len(policy.AgentID)-1:] == "0" { // Simulate some agents without PAN
		tdsPercentage = 10.0
	}
	tdsAmount := grossCommission * (tdsPercentage / 100.0)
	netCommission := grossCommission - tdsAmount

	// Determine commission type
	commissionType := domain.CommissionTypeFirstYear
	if policy.PolicyType == "RENEWAL" {
		commissionType = domain.CommissionTypeRenewal
		// Renewal commission is typically lower
		grossCommission = grossCommission * 0.3 // 30% of first year
		tdsAmount = grossCommission * (tdsPercentage / 100.0)
		netCommission = grossCommission - tdsAmount
	}

	transaction := &domain.CommissionTransaction{
		BatchID:          batchID,
		AgentID:          policy.AgentID,
		PolicyNumber:     policy.PolicyNumber,
		CommissionType:   commissionType,
		GrossCommission:  grossCommission,
		TDSAmount:        tdsAmount,
		TDSPercentage:    tdsPercentage,
		NetCommission:    netCommission,
		CommissionStatus: domain.CommissionStatusCalculated,
		PremiumAmount:    policy.PremiumAmount,
		ProductType:      policy.ProductType,
		PlanCode:         policy.PlanCode,
	}

	return transaction, nil
}

// Test database setup for load tests
func setupLoadTestDB(t *testing.T) *postgres.DB {
	// TODO: Implement load test database setup
	// Should create a separate test database optimized for load testing
	t.Skip("Load test database setup not implemented")
	return nil
}

func cleanupLoadTestDB(t *testing.T, db *postgres.DB) {
	// TODO: Implement load test database cleanup
}
