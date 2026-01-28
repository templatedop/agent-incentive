package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "gitlab.cept.gov.in/it-2.0-common/api-log"
)

// PolicyClient handles communication with Policy Services
// INT-IC-002: Policy Services integration
type PolicyClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// PolicyServiceConfig holds policy service configuration
type PolicyServiceConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// PolicyData represents policy information from policy service
type PolicyData struct {
	PolicyNumber     string    `json:"policy_number"`
	AgentID          string    `json:"agent_id"`
	ProductType      string    `json:"product_type"`
	PlanCode         string    `json:"plan_code"`
	PolicyStatus     string    `json:"policy_status"`
	PremiumAmount    float64   `json:"premium_amount"`
	SumAssured       float64   `json:"sum_assured"`
	InceptionDate    time.Time `json:"inception_date"`
	MaturityDate     time.Time `json:"maturity_date"`
	PaymentMode      string    `json:"payment_mode"` // MONTHLY, QUARTERLY, YEARLY
	PolicyType       string    `json:"policy_type"`  // FIRST_YEAR, RENEWAL
	CommissionableYN bool      `json:"commissionable_yn"`
}

// PolicyListResponse represents the response from policy list API
type PolicyListResponse struct {
	Policies   []PolicyData `json:"policies"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
}

// PolicyStatusChangeEvent represents policy status change webhook payload
type PolicyStatusChangeEvent struct {
	PolicyNumber string    `json:"policy_number"`
	OldStatus    string    `json:"old_status"`
	NewStatus    string    `json:"new_status"`
	ChangeDate   time.Time `json:"change_date"`
	Reason       string    `json:"reason"`
	AgentID      string    `json:"agent_id"`
}

// NewPolicyClient creates a new policy service client
func NewPolicyClient(config PolicyServiceConfig) *PolicyClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &PolicyClient{
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// GetPoliciesForCommission fetches policies eligible for commission calculation
// Used in commission batch processing workflow
func (c *PolicyClient) GetPoliciesForCommission(
	ctx context.Context,
	batchMonth string, // Format: YYYY-MM
	page, limit int,
) (*PolicyListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/policies/commissionable?batch_month=%s&page=%d&limit=%d",
		c.baseURL, batchMonth, page, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("policy service returned status %d: %s", resp.StatusCode, string(body))
	}

	var policyResp PolicyListResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info(ctx, "Fetched %d policies for batch month %s (page %d)", len(policyResp.Policies), batchMonth, page)
	return &policyResp, nil
}

// GetPolicyByNumber fetches a specific policy by policy number
func (c *PolicyClient) GetPolicyByNumber(ctx context.Context, policyNumber string) (*PolicyData, error) {
	url := fmt.Sprintf("%s/api/v1/policies/%s", c.baseURL, policyNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("policy not found: %s", policyNumber)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("policy service returned status %d: %s", resp.StatusCode, string(body))
	}

	var policy PolicyData
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &policy, nil
}

// ValidatePolicyStatus validates if a policy is active and commissionable
func (c *PolicyClient) ValidatePolicyStatus(ctx context.Context, policyNumber string) (bool, error) {
	policy, err := c.GetPolicyByNumber(ctx, policyNumber)
	if err != nil {
		return false, err
	}

	// Policy must be active and commissionable
	if policy.PolicyStatus != "ACTIVE" && policy.PolicyStatus != "INFORCE" {
		log.Warn(ctx, "Policy %s is not active: %s", policyNumber, policy.PolicyStatus)
		return false, nil
	}

	if !policy.CommissionableYN {
		log.Warn(ctx, "Policy %s is not commissionable", policyNumber)
		return false, nil
	}

	return true, nil
}

// GetPoliciesByAgent fetches all policies for a specific agent
func (c *PolicyClient) GetPoliciesByAgent(
	ctx context.Context,
	agentID string,
	status string, // ACTIVE, LAPSED, SURRENDERED, etc.
	page, limit int,
) (*PolicyListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/policies?agent_id=%s&status=%s&page=%d&limit=%d",
		c.baseURL, agentID, status, page, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("policy service returned status %d: %s", resp.StatusCode, string(body))
	}

	var policyResp PolicyListResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &policyResp, nil
}

// NotifyCommissionProcessed sends notification to policy service about commission processing
// Used to update policy commission_paid_yn flag
func (c *PolicyClient) NotifyCommissionProcessed(
	ctx context.Context,
	policyNumber string,
	commissionAmount float64,
	batchID string,
) error {
	url := fmt.Sprintf("%s/api/v1/policies/%s/commission-processed", c.baseURL, policyNumber)

	payload := map[string]interface{}{
		"policy_number":     policyNumber,
		"commission_amount": commissionAmount,
		"batch_id":          batchID,
		"processed_at":      time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("policy service returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Info(ctx, "Notified policy service about commission processing for policy %s", policyNumber)
	return nil
}

// HealthCheck checks if policy service is reachable
func (c *PolicyClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("policy service health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("policy service unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
