package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "gitlab.cept.gov.in/it-2.0-common/api-log"
)

// PFMSClient handles communication with PFMS (Public Financial Management System)
// INT-IC-003: PFMS/Bank integration for EFT payments
type PFMSClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	secretKey  string // For HMAC signature
	orgCode    string // Organization code in PFMS
}

// PFMSConfig holds PFMS configuration
type PFMSConfig struct {
	BaseURL   string
	APIKey    string
	SecretKey string
	OrgCode   string
	Timeout   time.Duration
}

// EFTPaymentRequest represents EFT payment request to PFMS
type EFTPaymentRequest struct {
	RequestID          string    `json:"request_id"`      // Unique request ID for idempotency
	PaymentType        string    `json:"payment_type"`    // EFT, NEFT, RTGS
	Amount             float64   `json:"amount"`          // Payment amount
	Currency           string    `json:"currency"`        // INR
	BeneficiaryName    string    `json:"beneficiary_name"`
	BeneficiaryAccount string    `json:"beneficiary_account"`
	BeneficiaryIFSC    string    `json:"beneficiary_ifsc"`
	BeneficiaryBankName string   `json:"beneficiary_bank_name"`
	PayerAccount       string    `json:"payer_account"`   // Organization account
	PayerIFSC          string    `json:"payer_ifsc"`
	Purpose            string    `json:"purpose"`         // Commission Payment
	Remarks            string    `json:"remarks"`
	TransactionDate    time.Time `json:"transaction_date"`
	// Reference data
	ReferenceType   string `json:"reference_type"`   // DISBURSEMENT, COMMISSION, etc.
	ReferenceNumber string `json:"reference_number"` // Disbursement ID
}

// EFTPaymentResponse represents response from PFMS
type EFTPaymentResponse struct {
	PaymentID       string    `json:"payment_id"`       // PFMS payment ID
	RequestID       string    `json:"request_id"`
	Status          string    `json:"status"`           // INITIATED, PENDING, SUCCESS, FAILED
	TransactionID   string    `json:"transaction_id"`   // Bank transaction ID (UTR)
	ProcessedAt     time.Time `json:"processed_at"`
	Message         string    `json:"message"`
	ErrorCode       string    `json:"error_code,omitempty"`
	ErrorMessage    string    `json:"error_message,omitempty"`
}

// PaymentStatusResponse represents payment status query response
type PaymentStatusResponse struct {
	PaymentID       string    `json:"payment_id"`
	RequestID       string    `json:"request_id"`
	Status          string    `json:"status"`
	TransactionID   string    `json:"transaction_id"`
	Amount          float64   `json:"amount"`
	InitiatedAt     time.Time `json:"initiated_at"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	FailureReason   string    `json:"failure_reason,omitempty"`
	RetryCount      int       `json:"retry_count"`
	LastUpdated     time.Time `json:"last_updated"`
}

// PaymentConfirmationWebhook represents PFMS webhook payload
type PaymentConfirmationWebhook struct {
	PaymentID       string    `json:"payment_id"`
	RequestID       string    `json:"request_id"`
	TransactionID   string    `json:"transaction_id"`   // UTR number
	Status          string    `json:"status"`           // SUCCESS, FAILED, RETURNED
	Amount          float64   `json:"amount"`
	ProcessedAt     time.Time `json:"processed_at"`
	CompletedAt     time.Time `json:"completed_at"`
	FailureReason   string    `json:"failure_reason,omitempty"`
	Signature       string    `json:"signature"`        // HMAC signature
}

// NewPFMSClient creates a new PFMS client
func NewPFMSClient(config PFMSConfig) *PFMSClient {
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second // PFMS operations can be slow
	}

	return &PFMSClient{
		baseURL:   config.BaseURL,
		apiKey:    config.APIKey,
		secretKey: config.SecretKey,
		orgCode:   config.OrgCode,
		httpClient: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        50,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// InitiateEFTPayment initiates an EFT payment through PFMS
// Used in disbursement workflow for EFT mode
func (c *PFMSClient) InitiateEFTPayment(
	ctx context.Context,
	req *EFTPaymentRequest,
) (*EFTPaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments/eft/initiate", c.baseURL)

	// Add organization code
	payload := map[string]interface{}{
		"org_code":            c.orgCode,
		"request_id":          req.RequestID,
		"payment_type":        req.PaymentType,
		"amount":              req.Amount,
		"currency":            req.Currency,
		"beneficiary_name":    req.BeneficiaryName,
		"beneficiary_account": req.BeneficiaryAccount,
		"beneficiary_ifsc":    req.BeneficiaryIFSC,
		"beneficiary_bank":    req.BeneficiaryBankName,
		"payer_account":       req.PayerAccount,
		"payer_ifsc":          req.PayerIFSC,
		"purpose":             req.Purpose,
		"remarks":             req.Remarks,
		"transaction_date":    req.TransactionDate.Format(time.RFC3339),
		"reference_type":      req.ReferenceType,
		"reference_number":    req.ReferenceNumber,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Generate HMAC signature for request authentication
	signature := c.generateSignature(jsonData)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)
	httpReq.Header.Set("X-Org-Code", c.orgCode)
	httpReq.Header.Set("X-Signature", signature)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("PFMS returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var paymentResp EFTPaymentResponse
	if err := json.Unmarshal(bodyBytes, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info(ctx, "EFT payment initiated: PaymentID=%s, RequestID=%s, Status=%s",
		paymentResp.PaymentID, paymentResp.RequestID, paymentResp.Status)

	return &paymentResp, nil
}

// GetPaymentStatus queries the status of an EFT payment
func (c *PFMSClient) GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentStatusResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments/%s/status", c.baseURL, paymentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Org-Code", c.orgCode)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("payment not found: %s", paymentID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("PFMS returned status %d: %s", resp.StatusCode, string(body))
	}

	var statusResp PaymentStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &statusResp, nil
}

// CancelPayment cancels a pending EFT payment
// Only payments in INITIATED or PENDING status can be cancelled
func (c *PFMSClient) CancelPayment(ctx context.Context, paymentID, reason string) error {
	url := fmt.Sprintf("%s/api/v1/payments/%s/cancel", c.baseURL, paymentID)

	payload := map[string]interface{}{
		"payment_id": paymentID,
		"reason":     reason,
		"cancelled_by": "COMMISSION_SYSTEM",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	signature := c.generateSignature(jsonData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Org-Code", c.orgCode)
	req.Header.Set("X-Signature", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PFMS returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Info(ctx, "Payment cancelled: %s, Reason: %s", paymentID, reason)
	return nil
}

// RetryFailedPayment retries a failed EFT payment
// Creates a new payment with same details but new request ID
func (c *PFMSClient) RetryFailedPayment(
	ctx context.Context,
	originalRequestID string,
	newRequestID string,
) (*EFTPaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments/retry", c.baseURL)

	payload := map[string]interface{}{
		"original_request_id": originalRequestID,
		"new_request_id":      newRequestID,
		"org_code":            c.orgCode,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	signature := c.generateSignature(jsonData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Org-Code", c.orgCode)
	req.Header.Set("X-Signature", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("PFMS returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var paymentResp EFTPaymentResponse
	if err := json.Unmarshal(bodyBytes, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info(ctx, "Payment retry initiated: NewPaymentID=%s, OriginalRequestID=%s",
		paymentResp.PaymentID, originalRequestID)

	return &paymentResp, nil
}

// ValidateBankAccount validates bank account details with PFMS
// Performs penny drop or account validation check
func (c *PFMSClient) ValidateBankAccount(
	ctx context.Context,
	accountNumber, ifscCode, accountHolderName string,
) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/validate", c.baseURL)

	payload := map[string]interface{}{
		"account_number":      accountNumber,
		"ifsc_code":           ifscCode,
		"account_holder_name": accountHolderName,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Org-Code", c.orgCode)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Warn(ctx, "Bank account validation failed: %s", string(body))
		return false, nil
	}

	var validationResp struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&validationResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return validationResp.Valid, nil
}

// VerifyWebhookSignature verifies the HMAC signature of PFMS webhook
func (c *PFMSClient) VerifyWebhookSignature(payload []byte, signature string) bool {
	expectedSignature := c.generateSignature(payload)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// generateSignature generates HMAC-SHA256 signature for request authentication
func (c *PFMSClient) generateSignature(payload []byte) string {
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write(payload)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HealthCheck checks if PFMS is reachable
func (c *PFMSClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("PFMS health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PFMS unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
