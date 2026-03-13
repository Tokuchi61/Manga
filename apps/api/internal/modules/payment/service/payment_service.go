package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
	paymentvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

const defaultNonProdCallbackSigningSecret = "novascans-local-payment-callback-secret"

var (
	ErrValidation              = errors.New("payment_validation_failed")
	ErrNotFound                = errors.New("payment_not_found")
	ErrConflict                = errors.New("payment_conflict")
	ErrForbiddenAction         = errors.New("payment_forbidden_action")
	ErrManaPurchaseDisabled    = errors.New("payment_mana_purchase_disabled")
	ErrCheckoutDisabled        = errors.New("payment_checkout_disabled")
	ErrTransactionReadDisabled = errors.New("payment_transaction_read_disabled")
	ErrCallbackIntakePaused    = errors.New("payment_callback_intake_paused")
	ErrCallbackSignature       = errors.New("payment_callback_signature_invalid")
	ErrCallbackTimestamp       = errors.New("payment_callback_timestamp_invalid")
)

// CallbackEventPublisher publishes durable callback lifecycle events.
type CallbackEventPublisher interface {
	PublishPaymentCallback(ctx context.Context, event PaymentCallbackEvent) error
}

// PaymentCallbackEvent captures callback outcome for async consumers.
type PaymentCallbackEvent struct {
	ProviderEventID   string
	SessionID         string
	TransactionID     string
	TransactionStatus string
	RequestID         string
	CorrelationID     string
	ProcessedAt       time.Time
}

type Config struct {
	CallbackSigningSecret  string
	CallbackTimestampSkew  time.Duration
	CallbackEventPublisher CallbackEventPublisher
}

func (c Config) withDefaults() Config {
	if strings.TrimSpace(c.CallbackSigningSecret) == "" {
		env := strings.TrimSpace(strings.ToLower(os.Getenv("APP_ENV")))
		if env == "prod" || env == "production" {
			c.CallbackSigningSecret = ""
		} else {
			c.CallbackSigningSecret = defaultNonProdCallbackSigningSecret
		}
	}
	if c.CallbackTimestampSkew <= 0 {
		c.CallbackTimestampSkew = 5 * time.Minute
	}
	return c
}

// PaymentService owns stage-19 payment flows.
type PaymentService struct {
	store                  repository.Store
	validator              *validation.Validator
	now                    func() time.Time
	callbackSigningSecret  string
	callbackTimestampSkew  time.Duration
	callbackEventPublisher CallbackEventPublisher
}

func New(store repository.Store, validator *validation.Validator, cfg ...Config) *PaymentService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	resolvedCfg := Config{}.withDefaults()
	if len(cfg) > 0 {
		resolvedCfg = cfg[0].withDefaults()
	}
	return &PaymentService{
		store:                  store,
		validator:              validator,
		now:                    time.Now,
		callbackSigningSecret:  strings.TrimSpace(resolvedCfg.CallbackSigningSecret),
		callbackTimestampSkew:  resolvedCfg.CallbackTimestampSkew,
		callbackEventPublisher: resolvedCfg.CallbackEventPublisher,
	}
}

func (s *PaymentService) validateInput(payload any) error {
	if err := paymentvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *PaymentService) requireManaPurchaseEnabled(enabled bool) error {
	if !enabled {
		return ErrManaPurchaseDisabled
	}
	return nil
}

func (s *PaymentService) requireCheckoutEnabled(enabled bool) error {
	if !enabled {
		return ErrCheckoutDisabled
	}
	return nil
}

func (s *PaymentService) requireTransactionReadEnabled(enabled bool) error {
	if !enabled {
		return ErrTransactionReadDisabled
	}
	return nil
}

func (s *PaymentService) requireCallbackIntakeOpen(paused bool) error {
	if paused {
		return ErrCallbackIntakePaused
	}
	return nil
}

func (s *PaymentService) VerifyProviderCallbackSignature(payload []byte, signature string, timestamp string) error {
	secret := strings.TrimSpace(s.callbackSigningSecret)
	if secret == "" {
		return ErrCallbackSignature
	}
	timestamp = strings.TrimSpace(timestamp)
	signature = strings.TrimSpace(strings.ToLower(signature))
	if timestamp == "" || signature == "" {
		return ErrCallbackSignature
	}

	unixSeconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ErrCallbackTimestamp
	}
	issuedAt := time.Unix(unixSeconds, 0).UTC()
	now := s.now().UTC()
	if issuedAt.After(now.Add(s.callbackTimestampSkew)) || issuedAt.Before(now.Add(-s.callbackTimestampSkew)) {
		return ErrCallbackTimestamp
	}

	expected := signCallbackPayload(secret, timestamp, payload)
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrCallbackSignature
	}
	return nil
}

func (s *PaymentService) publishCallbackEvent(ctx context.Context, request dto.ProcessProviderCallbackRequest, tx entity.Transaction, processedAt time.Time) error {
	if s == nil || s.callbackEventPublisher == nil {
		return nil
	}
	return s.callbackEventPublisher.PublishPaymentCallback(ctx, PaymentCallbackEvent{
		ProviderEventID:   request.ProviderEventID,
		SessionID:         request.SessionID,
		TransactionID:     tx.TransactionID,
		TransactionStatus: tx.Status,
		RequestID:         request.ProviderEventID,
		CorrelationID:     request.ProviderReference,
		ProcessedAt:       processedAt.UTC(),
	})
}

func SignProviderCallback(secret string, timestamp string, payload []byte) string {
	return signCallbackPayload(strings.TrimSpace(secret), strings.TrimSpace(timestamp), payload)
}

func signCallbackPayload(secret string, timestamp string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp))
	_, _ = mac.Write([]byte("."))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func DefaultNonProdCallbackSigningSecret() string {
	return defaultNonProdCallbackSigningSecret
}
