package auth

import "errors"

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
