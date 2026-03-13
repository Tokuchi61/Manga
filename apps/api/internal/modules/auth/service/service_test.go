package service

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/dto"
	authrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/stretchr/testify/require"
)

func newTestService(nowRef *time.Time) *AuthService {
	svc := New(authrepository.NewMemoryStore(), validation.New(), Config{
		FailedAttemptLimit:         2,
		LoginCooldown:              2 * time.Minute,
		VerificationResendCooldown: time.Minute,
		AccessTokenTTL:             15 * time.Minute,
		RefreshTokenTTL:            time.Hour,
		PasswordResetTokenTTL:      time.Hour,
		EmailVerificationTokenTTL:  time.Hour,
		ExposeSensitiveTokens:      true,
	})
	svc.now = func() time.Time { return nowRef.UTC() }
	return svc
}

func TestRegisterVerifyAndLoginFlow(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	registerRes, err := svc.Register(ctx, dto.RegisterRequest{
		Email:    "user@example.com",
		Password: "StrongPass123!",
	}, RequestMeta{RequestID: "req-1", CorrelationID: "corr-1"})
	require.NoError(t, err)
	require.False(t, registerRes.EmailVerified)
	require.NotEmpty(t, registerRes.VerificationToken)

	_, err = svc.Login(ctx, dto.LoginRequest{
		Email:    "user@example.com",
		Password: "StrongPass123!",
	}, RequestMeta{RequestID: "req-2", CorrelationID: "corr-2"})
	require.ErrorIs(t, err, ErrEmailNotVerified)

	_, err = svc.ConfirmEmailVerification(ctx, dto.ConfirmVerificationRequest{Token: registerRes.VerificationToken}, RequestMeta{RequestID: "req-3", CorrelationID: "corr-3"})
	require.NoError(t, err)

	loginRes, err := svc.Login(ctx, dto.LoginRequest{
		Email:    "user@example.com",
		Password: "StrongPass123!",
		Device:   "ios",
		IP:       "127.0.0.1",
	}, RequestMeta{RequestID: "req-4", CorrelationID: "corr-4"})
	require.NoError(t, err)
	require.NotEmpty(t, loginRes.SessionID)
	require.NotEmpty(t, loginRes.AccessToken)
	require.NotEmpty(t, loginRes.RefreshToken)
}

func TestLoginFailedLimitTriggersCooldown(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	registerRes, err := svc.Register(ctx, dto.RegisterRequest{Email: "cooldown@example.com", Password: "StrongPass123!"}, RequestMeta{RequestID: "req-1"})
	require.NoError(t, err)

	_, err = svc.ConfirmEmailVerification(ctx, dto.ConfirmVerificationRequest{Token: registerRes.VerificationToken}, RequestMeta{RequestID: "req-2"})
	require.NoError(t, err)

	_, err = svc.Login(ctx, dto.LoginRequest{Email: "cooldown@example.com", Password: "WrongPass123!"}, RequestMeta{RequestID: "req-3"})
	require.ErrorIs(t, err, ErrInvalidCredentials)

	_, err = svc.Login(ctx, dto.LoginRequest{Email: "cooldown@example.com", Password: "WrongPass123!"}, RequestMeta{RequestID: "req-4"})
	require.ErrorIs(t, err, ErrLoginCooldownActive)
}

func TestRefreshRotationAndSessionRevoke(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	registerRes, err := svc.Register(ctx, dto.RegisterRequest{Email: "refresh@example.com", Password: "StrongPass123!"}, RequestMeta{RequestID: "req-1"})
	require.NoError(t, err)
	_, err = svc.ConfirmEmailVerification(ctx, dto.ConfirmVerificationRequest{Token: registerRes.VerificationToken}, RequestMeta{RequestID: "req-2"})
	require.NoError(t, err)

	loginRes, err := svc.Login(ctx, dto.LoginRequest{Email: "refresh@example.com", Password: "StrongPass123!"}, RequestMeta{RequestID: "req-3"})
	require.NoError(t, err)

	refreshRes, err := svc.RefreshToken(ctx, dto.RefreshTokenRequest{RefreshToken: loginRes.RefreshToken}, RequestMeta{RequestID: "req-4"})
	require.NoError(t, err)
	require.NotEqual(t, loginRes.RefreshToken, refreshRes.RefreshToken)

	_, err = svc.RefreshToken(ctx, dto.RefreshTokenRequest{RefreshToken: loginRes.RefreshToken}, RequestMeta{RequestID: "req-5"})
	require.ErrorIs(t, err, ErrTokenInvalidOrExpired)

	_, err = svc.RevokeCurrentSession(ctx, dto.RevokeCurrentSessionRequest{CredentialID: loginRes.CredentialID, SessionID: loginRes.SessionID}, RequestMeta{RequestID: "req-6"})
	require.NoError(t, err)

	_, err = svc.RefreshToken(ctx, dto.RefreshTokenRequest{RefreshToken: refreshRes.RefreshToken}, RequestMeta{RequestID: "req-7"})
	require.ErrorIs(t, err, ErrSessionNotFound)
}

func TestForgotResetAndChangePassword(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	registerRes, err := svc.Register(ctx, dto.RegisterRequest{Email: "password@example.com", Password: "StrongPass123!"}, RequestMeta{RequestID: "req-1"})
	require.NoError(t, err)
	_, err = svc.ConfirmEmailVerification(ctx, dto.ConfirmVerificationRequest{Token: registerRes.VerificationToken}, RequestMeta{RequestID: "req-2"})
	require.NoError(t, err)

	forgotRes, err := svc.ForgotPassword(ctx, dto.ForgotPasswordRequest{Email: "password@example.com"}, RequestMeta{RequestID: "req-3"})
	require.NoError(t, err)
	require.NotEmpty(t, forgotRes.Token)

	_, err = svc.ResetPassword(ctx, dto.ResetPasswordRequest{Token: forgotRes.Token, NewPassword: "NewStrongPass123!"}, RequestMeta{RequestID: "req-4"})
	require.NoError(t, err)

	_, err = svc.Login(ctx, dto.LoginRequest{Email: "password@example.com", Password: "StrongPass123!"}, RequestMeta{RequestID: "req-5"})
	require.ErrorIs(t, err, ErrInvalidCredentials)

	loginRes, err := svc.Login(ctx, dto.LoginRequest{Email: "password@example.com", Password: "NewStrongPass123!"}, RequestMeta{RequestID: "req-6"})
	require.NoError(t, err)

	_, err = svc.ChangePassword(ctx, dto.ChangePasswordRequest{CredentialID: loginRes.CredentialID, OldPassword: "NewStrongPass123!", NewPassword: "FinalStrongPass123!"}, RequestMeta{RequestID: "req-7"})
	require.NoError(t, err)

	_, err = svc.Login(ctx, dto.LoginRequest{Email: "password@example.com", Password: "FinalStrongPass123!"}, RequestMeta{RequestID: "req-8"})
	require.NoError(t, err)
}

func TestVerificationResendCooldown(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	svc := newTestService(&now)
	ctx := context.Background()

	registerRes, err := svc.Register(ctx, dto.RegisterRequest{Email: "resend@example.com", Password: "StrongPass123!"}, RequestMeta{RequestID: "req-1"})
	require.NoError(t, err)
	require.NotEmpty(t, registerRes.VerificationToken)

	_, err = svc.SendEmailVerification(ctx, dto.SendVerificationRequest{CredentialID: registerRes.CredentialID}, RequestMeta{RequestID: "req-2"})
	require.ErrorIs(t, err, ErrVerificationResendCooldown)

	now = now.Add(2 * time.Minute)
	resendRes, err := svc.SendEmailVerification(ctx, dto.SendVerificationRequest{CredentialID: registerRes.CredentialID}, RequestMeta{RequestID: "req-3"})
	require.NoError(t, err)
	require.NotEmpty(t, resendRes.Token)
}
