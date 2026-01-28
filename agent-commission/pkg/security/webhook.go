package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrInvalidSignature is returned when webhook signature verification fails
	ErrInvalidSignature = errors.New("invalid webhook signature")

	// ErrExpiredTimestamp is returned when webhook timestamp is too old
	ErrExpiredTimestamp = errors.New("webhook timestamp expired")

	// ErrMissingSignature is returned when signature header is missing
	ErrMissingSignature = errors.New("missing webhook signature")
)

// WebhookVerifier handles webhook signature verification
type WebhookVerifier struct {
	secretKey string
	tolerance time.Duration // Max age of webhook (e.g., 5 minutes)
}

// NewWebhookVerifier creates a new webhook verifier
func NewWebhookVerifier(secretKey string, tolerance time.Duration) *WebhookVerifier {
	if tolerance == 0 {
		tolerance = 5 * time.Minute // Default 5 minutes
	}

	return &WebhookVerifier{
		secretKey: secretKey,
		tolerance: tolerance,
	}
}

// VerifyHMACSHA256 verifies HMAC-SHA256 signature of webhook payload
// Used for PFMS, Policy Services webhooks
func (w *WebhookVerifier) VerifyHMACSHA256(payload []byte, signature string) error {
	if signature == "" {
		return ErrMissingSignature
	}

	expectedSignature := w.generateHMACSHA256(payload)

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return ErrInvalidSignature
	}

	return nil
}

// VerifyHMACSHA256WithTimestamp verifies signature with timestamp validation
// Format: "t=<timestamp>,v1=<signature>"
func (w *WebhookVerifier) VerifyHMACSHA256WithTimestamp(
	payload []byte,
	signatureHeader string,
	timestamp int64,
) error {
	if signatureHeader == "" {
		return ErrMissingSignature
	}

	// Verify timestamp is not too old
	webhookTime := time.Unix(timestamp, 0)
	if time.Since(webhookTime) > w.tolerance {
		return ErrExpiredTimestamp
	}

	// Create signed payload with timestamp
	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payload))
	expectedSignature := w.generateHMACSHA256([]byte(signedPayload))

	if !hmac.Equal([]byte(signatureHeader), []byte(expectedSignature)) {
		return ErrInvalidSignature
	}

	return nil
}

// generateHMACSHA256 generates HMAC-SHA256 signature (base64 encoded)
func (w *WebhookVerifier) generateHMACSHA256(payload []byte) string {
	h := hmac.New(sha256.New, []byte(w.secretKey))
	h.Write(payload)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GenerateHMACSHA256Hex generates HMAC-SHA256 signature (hex encoded)
func (w *WebhookVerifier) GenerateHMACSHA256Hex(payload []byte) string {
	h := hmac.New(sha256.New, []byte(w.secretKey))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyGitHubStyle verifies GitHub-style webhook signature
// Format: "sha256=<hex_signature>"
func (w *WebhookVerifier) VerifyGitHubStyle(payload []byte, signature string) error {
	if signature == "" {
		return ErrMissingSignature
	}

	// Remove "sha256=" prefix if present
	if len(signature) > 7 && signature[:7] == "sha256=" {
		signature = signature[7:]
	}

	expectedSignature := w.GenerateHMACSHA256Hex(payload)

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return ErrInvalidSignature
	}

	return nil
}

// IdempotencyKey represents a key for idempotency checks
type IdempotencyKey struct {
	RequestID string
	Timestamp time.Time
}

// IdempotencyStore interface for storing processed webhook IDs
type IdempotencyStore interface {
	Exists(requestID string) (bool, error)
	Store(requestID string, ttl time.Duration) error
}

// WebhookProcessor handles webhook processing with idempotency
type WebhookProcessor struct {
	verifier         *WebhookVerifier
	idempotencyStore IdempotencyStore
}

// NewWebhookProcessor creates a new webhook processor
func NewWebhookProcessor(verifier *WebhookVerifier, store IdempotencyStore) *WebhookProcessor {
	return &WebhookProcessor{
		verifier:         verifier,
		idempotencyStore: store,
	}
}

// ProcessWebhook verifies signature and checks idempotency
func (wp *WebhookProcessor) ProcessWebhook(
	payload []byte,
	signature string,
	requestID string,
) error {
	// Verify signature
	if err := wp.verifier.VerifyHMACSHA256(payload, signature); err != nil {
		return err
	}

	// Check idempotency
	if wp.idempotencyStore != nil {
		exists, err := wp.idempotencyStore.Exists(requestID)
		if err != nil {
			return fmt.Errorf("idempotency check failed: %w", err)
		}

		if exists {
			return fmt.Errorf("duplicate webhook request: %s", requestID)
		}

		// Store request ID (TTL: 24 hours)
		if err := wp.idempotencyStore.Store(requestID, 24*time.Hour); err != nil {
			return fmt.Errorf("failed to store idempotency key: %w", err)
		}
	}

	return nil
}
