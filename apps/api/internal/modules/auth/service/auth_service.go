package service

import (
	"errors"
	"fmt"
	"time"

	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	authvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
)

var (
	ErrValidation                 = errors.New("auth_validation_failed")
	ErrCredentialAlreadyExists    = errors.New("auth_credential_already_exists")
	ErrInvalidCredentials         = errors.New("auth_invalid_credentials")
	ErrEmailNotVerified           = errors.New("auth_email_not_verified")
	ErrCredentialSuspended        = errors.New("auth_credential_suspended")
	ErrCredentialBanned           = errors.New("auth_credential_banned")
	ErrLoginCooldownActive        = errors.New("auth_login_cooldown_active")
	ErrSessionNotFound            = errors.New("auth_session_not_found")
	ErrTokenInvalidOrExpired      = errors.New("auth_token_invalid_or_expired")
	ErrVerificationResendCooldown = errors.New("auth_verification_resend_cooldown")
)

// RequestMeta carries trace and source metadata.
type RequestMeta struct {
	RequestID     string
	CorrelationID string
	Device        string
	IP            string
}

// Config defines auth runtime and token settings.
type Config struct {
	FailedAttemptLimit         int
	LoginCooldown              time.Duration
	VerificationResendCooldown time.Duration
	AccessTokenTTL             time.Duration
	AccessTokenSecret          string
	AccessTokenIssuer          string
	RefreshTokenTTL            time.Duration
	PasswordResetTokenTTL      time.Duration
	EmailVerificationTokenTTL  time.Duration
	ExposeSensitiveTokens      bool
}

func (c Config) withDefaults() Config {
	if c.FailedAttemptLimit <= 0 {
		c.FailedAttemptLimit = 5
	}
	if c.LoginCooldown <= 0 {
		c.LoginCooldown = 5 * time.Minute
	}
	if c.VerificationResendCooldown <= 0 {
		c.VerificationResendCooldown = time.Minute
	}
	if c.AccessTokenTTL <= 0 {
		c.AccessTokenTTL = 15 * time.Minute
	}
	if c.AccessTokenSecret == "" {
		c.AccessTokenSecret = identity.AccessTokenSecret()
	}
	if c.AccessTokenIssuer == "" {
		c.AccessTokenIssuer = "novascans-api"
	}
	if c.RefreshTokenTTL <= 0 {
		c.RefreshTokenTTL = 30 * 24 * time.Hour
	}
	if c.PasswordResetTokenTTL <= 0 {
		c.PasswordResetTokenTTL = 30 * time.Minute
	}
	if c.EmailVerificationTokenTTL <= 0 {
		c.EmailVerificationTokenTTL = 24 * time.Hour
	}
	return c
}

// AuthService owns auth flows for stage-4 bootstrap.
type AuthService struct {
	store      authrepository.Store
	validator  *validation.Validator
	cfg        Config
	now        func() time.Time
	userLookup UserLookup
}

func New(store authrepository.Store, validator *validation.Validator, cfg Config) *AuthService {
	if store == nil {
		store = authrepository.NewMemoryStore()
	}
	return &AuthService{
		store:     store,
		validator: validator,
		cfg:       cfg.withDefaults(),
		now:       time.Now,
	}
}

func (s *AuthService) validateInput(payload any) error {
	if err := authvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *AuthService) shouldExposeSensitiveTokens() bool {
	if s == nil {
		return false
	}
	return s.cfg.ExposeSensitiveTokens
}
