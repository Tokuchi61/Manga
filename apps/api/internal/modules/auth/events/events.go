package events

const (
	EventLoginSucceeded          = "auth.login.succeeded"
	EventLoginFailed             = "auth.login.failed"
	EventSessionRevoked          = "auth.session.revoked"
	EventEmailVerificationSent   = "auth.email_verification.sent"
	EventSecuritySuspiciousLogin = "auth.security.suspicious_login"
)
