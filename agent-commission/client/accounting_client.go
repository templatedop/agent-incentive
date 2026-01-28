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

// AccountingClient handles communication with Accounting System
// INT-IC-004: Accounting system integration for GL posting
type AccountingClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	companyCode string
}

// AccountingConfig holds accounting system configuration
type AccountingConfig struct {
	BaseURL     string
	APIKey      string
	CompanyCode string
	Timeout     time.Duration
}

// GLEntry represents a General Ledger entry
type GLEntry struct {
	AccountCode   string  `json:"account_code"`   // GL account code
	AccountName   string  `json:"account_name"`
	DebitAmount   float64 `json:"debit_amount"`
	CreditAmount  float64 `json:"credit_amount"`
	CostCenter    string  `json:"cost_center"`
	Description   string  `json:"description"`
}

// VoucherRequest represents a voucher posting request
type VoucherRequest struct {
	VoucherType     string    `json:"voucher_type"`     // JV (Journal Voucher), PV (Payment Voucher)
	VoucherDate     time.Time `json:"voucher_date"`
	ReferenceType   string    `json:"reference_type"`   // DISBURSEMENT, COMMISSION, CLAWBACK, SUSPENSE
	ReferenceNumber string    `json:"reference_number"` // Disbursement ID, Batch ID, etc.
	Narration       string    `json:"narration"`
	Entries         []GLEntry `json:"entries"`
	TotalDebit      float64   `json:"total_debit"`
	TotalCredit     float64   `json:"total_credit"`
	CompanyCode     string    `json:"company_code"`
	BranchCode      string    `json:"branch_code,omitempty"`
}

// VoucherResponse represents voucher posting response
type VoucherResponse struct {
	VoucherID       string    `json:"voucher_id"`       // System-generated voucher ID
	VoucherNumber   string    `json:"voucher_number"`   // Sequential voucher number
	Status          string    `json:"status"`           // POSTED, PENDING_APPROVAL, REJECTED
	PostedAt        time.Time `json:"posted_at"`
	ApprovedBy      string    `json:"approved_by,omitempty"`
	Message         string    `json:"message"`
	ErrorCode       string    `json:"error_code,omitempty"`
	ErrorMessage    string    `json:"error_message,omitempty"`
}

// VoucherStatusResponse represents voucher status query response
type VoucherStatusResponse struct {
	VoucherID       string    `json:"voucher_id"`
	VoucherNumber   string    `json:"voucher_number"`
	Status          string    `json:"status"`
	ReferenceType   string    `json:"reference_type"`
	ReferenceNumber string    `json:"reference_number"`
	PostedAt        time.Time `json:"posted_at"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	ApprovedBy      string    `json:"approved_by,omitempty"`
	RejectedReason  string    `json:"rejected_reason,omitempty"`
}

// NewAccountingClient creates a new accounting system client
func NewAccountingClient(config AccountingConfig) *AccountingClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &AccountingClient{
		baseURL:     config.BaseURL,
		apiKey:      config.APIKey,
		companyCode: config.CompanyCode,
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

// PostCommissionVoucher posts commission payment voucher
// Used after disbursement is approved
func (c *AccountingClient) PostCommissionVoucher(
	ctx context.Context,
	disbursementID string,
	agentID string,
	commissionAmount float64,
	tdsAmount float64,
	netAmount float64,
	voucherDate time.Time,
) (*VoucherResponse, error) {
	// Prepare GL entries for commission payment
	// Debit: Commission Expense Account
	// Credit: Agent Payable Account
	// Credit: TDS Payable Account (if TDS > 0)

	entries := []GLEntry{
		{
			AccountCode:  "6100001", // Commission Expense
			AccountName:  "Agent Commission Expense",
			DebitAmount:  commissionAmount,
			CreditAmount: 0,
			Description:  fmt.Sprintf("Commission to Agent %s", agentID),
		},
	}

	if tdsAmount > 0 {
		entries = append(entries, GLEntry{
			AccountCode:  "2100002", // TDS Payable
			AccountName:  "TDS Payable to Government",
			DebitAmount:  0,
			CreditAmount: tdsAmount,
			Description:  fmt.Sprintf("TDS on Commission - Agent %s", agentID),
		})
	}

	entries = append(entries, GLEntry{
		AccountCode:  "2100001", // Agent Payable
		AccountName:  "Agent Commission Payable",
		DebitAmount:  0,
		CreditAmount: netAmount,
		Description:  fmt.Sprintf("Net Commission Payable - Agent %s", agentID),
	})

	voucherReq := VoucherRequest{
		VoucherType:     "PV", // Payment Voucher
		VoucherDate:     voucherDate,
		ReferenceType:   "DISBURSEMENT",
		ReferenceNumber: disbursementID,
		Narration:       fmt.Sprintf("Commission payment to Agent %s - Disbursement %s", agentID, disbursementID),
		Entries:         entries,
		TotalDebit:      commissionAmount,
		TotalCredit:     commissionAmount,
		CompanyCode:     c.companyCode,
	}

	return c.postVoucher(ctx, &voucherReq)
}

// PostClawbackVoucher posts clawback recovery voucher
// Used when commission is clawed back
func (c *AccountingClient) PostClawbackVoucher(
	ctx context.Context,
	clawbackID int64,
	agentID string,
	recoveryAmount float64,
	voucherDate time.Time,
) (*VoucherResponse, error) {
	// Debit: Agent Receivable Account
	// Credit: Commission Expense Account (reversal)

	entries := []GLEntry{
		{
			AccountCode:  "1200001", // Agent Receivable
			AccountName:  "Agent Commission Receivable",
			DebitAmount:  recoveryAmount,
			CreditAmount: 0,
			Description:  fmt.Sprintf("Clawback from Agent %s", agentID),
		},
		{
			AccountCode:  "6100001", // Commission Expense (reversal)
			AccountName:  "Agent Commission Expense",
			DebitAmount:  0,
			CreditAmount: recoveryAmount,
			Description:  fmt.Sprintf("Commission Clawback Reversal - Agent %s", agentID),
		},
	}

	voucherReq := VoucherRequest{
		VoucherType:     "JV", // Journal Voucher
		VoucherDate:     voucherDate,
		ReferenceType:   "CLAWBACK",
		ReferenceNumber: fmt.Sprintf("CLAWBACK-%d", clawbackID),
		Narration:       fmt.Sprintf("Commission clawback recovery from Agent %s", agentID),
		Entries:         entries,
		TotalDebit:      recoveryAmount,
		TotalCredit:     recoveryAmount,
		CompanyCode:     c.companyCode,
	}

	return c.postVoucher(ctx, &voucherReq)
}

// PostSuspenseVoucher posts suspense account entry
// Used when payment cannot be processed
func (c *AccountingClient) PostSuspenseVoucher(
	ctx context.Context,
	suspenseID int64,
	agentID string,
	amount float64,
	reason string,
	voucherDate time.Time,
) (*VoucherResponse, error) {
	// Debit: Suspense Account
	// Credit: Agent Payable Account

	entries := []GLEntry{
		{
			AccountCode:  "1300001", // Suspense Account
			AccountName:  "Commission Suspense",
			DebitAmount:  amount,
			CreditAmount: 0,
			Description:  fmt.Sprintf("Suspense - %s - Agent %s", reason, agentID),
		},
		{
			AccountCode:  "2100001", // Agent Payable
			AccountName:  "Agent Commission Payable",
			DebitAmount:  0,
			CreditAmount: amount,
			Description:  fmt.Sprintf("Transfer to Suspense - Agent %s", agentID),
		},
	}

	voucherReq := VoucherRequest{
		VoucherType:     "JV", // Journal Voucher
		VoucherDate:     voucherDate,
		ReferenceType:   "SUSPENSE",
		ReferenceNumber: fmt.Sprintf("SUSPENSE-%d", suspenseID),
		Narration:       fmt.Sprintf("Commission suspense for Agent %s - Reason: %s", agentID, reason),
		Entries:         entries,
		TotalDebit:      amount,
		TotalCredit:     amount,
		CompanyCode:     c.companyCode,
	}

	return c.postVoucher(ctx, &voucherReq)
}

// PostSuspenseResolutionVoucher posts suspense resolution entry
// Used when suspense amount is released
func (c *AccountingClient) PostSuspenseResolutionVoucher(
	ctx context.Context,
	suspenseID int64,
	agentID string,
	amount float64,
	voucherDate time.Time,
) (*VoucherResponse, error) {
	// Debit: Agent Payable Account
	// Credit: Suspense Account

	entries := []GLEntry{
		{
			AccountCode:  "2100001", // Agent Payable
			AccountName:  "Agent Commission Payable",
			DebitAmount:  amount,
			CreditAmount: 0,
			Description:  fmt.Sprintf("Suspense Release - Agent %s", agentID),
		},
		{
			AccountCode:  "1300001", // Suspense Account
			AccountName:  "Commission Suspense",
			DebitAmount:  0,
			CreditAmount: amount,
			Description:  fmt.Sprintf("Suspense Resolution - Agent %s", agentID),
		},
	}

	voucherReq := VoucherRequest{
		VoucherType:     "JV", // Journal Voucher
		VoucherDate:     voucherDate,
		ReferenceType:   "SUSPENSE_RESOLUTION",
		ReferenceNumber: fmt.Sprintf("SUSPENSE-RES-%d", suspenseID),
		Narration:       fmt.Sprintf("Suspense resolution for Agent %s", agentID),
		Entries:         entries,
		TotalDebit:      amount,
		TotalCredit:     amount,
		CompanyCode:     c.companyCode,
	}

	return c.postVoucher(ctx, &voucherReq)
}

// postVoucher posts a voucher to the accounting system
func (c *AccountingClient) postVoucher(ctx context.Context, req *VoucherRequest) (*VoucherResponse, error) {
	url := fmt.Sprintf("%s/api/v1/vouchers/post", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal voucher request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)
	httpReq.Header.Set("X-Company-Code", c.companyCode)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("accounting system returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var voucherResp VoucherResponse
	if err := json.Unmarshal(bodyBytes, &voucherResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info(ctx, "Voucher posted: VoucherID=%s, VoucherNumber=%s, Type=%s, Reference=%s",
		voucherResp.VoucherID, voucherResp.VoucherNumber, req.VoucherType, req.ReferenceNumber)

	return &voucherResp, nil
}

// GetVoucherStatus queries the status of a voucher
func (c *AccountingClient) GetVoucherStatus(ctx context.Context, voucherID string) (*VoucherStatusResponse, error) {
	url := fmt.Sprintf("%s/api/v1/vouchers/%s/status", c.baseURL, voucherID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Company-Code", c.companyCode)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("voucher not found: %s", voucherID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("accounting system returned status %d: %s", resp.StatusCode, string(body))
	}

	var statusResp VoucherStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &statusResp, nil
}

// CancelVoucher cancels a pending voucher
func (c *AccountingClient) CancelVoucher(ctx context.Context, voucherID, reason string) error {
	url := fmt.Sprintf("%s/api/v1/vouchers/%s/cancel", c.baseURL, voucherID)

	payload := map[string]interface{}{
		"voucher_id": voucherID,
		"reason":     reason,
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
	req.Header.Set("X-Company-Code", c.companyCode)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("accounting system returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Info(ctx, "Voucher cancelled: %s, Reason: %s", voucherID, reason)
	return nil
}

// HealthCheck checks if accounting system is reachable
func (c *AccountingClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("accounting system health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("accounting system unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
