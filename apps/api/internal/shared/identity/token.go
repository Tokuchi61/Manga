package identity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	defaultNonProdAccessTokenSecret = "novascans-local-access-token-secret"
	defaultTokenIssuer              = "novascans-api"
)

var (
	ErrAccessTokenMalformed = errors.New("access_token_malformed")
	ErrAccessTokenInvalid   = errors.New("access_token_invalid")
	ErrAccessTokenExpired   = errors.New("access_token_expired")
	ErrAccessTokenSecret    = errors.New("access_token_secret_missing")

	accessTokenSecretMu sync.RWMutex
	accessTokenSecret   string
)

type TokenClaims struct {
	UserID       string
	CredentialID string
	Roles        []string
	ExpiresAt    time.Time
	Issuer       string
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type tokenPayload struct {
	UserID       string   `json:"uid,omitempty"`
	CredentialID string   `json:"cid,omitempty"`
	Roles        []string `json:"roles,omitempty"`
	ExpiresAt    int64    `json:"exp"`
	Issuer       string   `json:"iss,omitempty"`
}

func SetAccessTokenSecret(secret string) {
	accessTokenSecretMu.Lock()
	accessTokenSecret = strings.TrimSpace(secret)
	accessTokenSecretMu.Unlock()
}

func AccessTokenSecret() string {
	return resolveSecret()
}

func IssueAccessToken(claims TokenClaims) (string, error) {
	return IssueAccessTokenWithSecret(resolveSecret(), claims)
}

func IssueAccessTokenWithSecret(secret string, claims TokenClaims) (string, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "", ErrAccessTokenSecret
	}
	if err := validateClaimsShape(claims); err != nil {
		return "", err
	}

	payload := tokenPayload{
		UserID:       strings.TrimSpace(claims.UserID),
		CredentialID: strings.TrimSpace(claims.CredentialID),
		Roles:        normalizeRoles(claims.Roles),
		ExpiresAt:    claims.ExpiresAt.UTC().Unix(),
		Issuer:       normalizeIssuer(claims.Issuer),
	}
	header := tokenHeader{Alg: "HS256", Typ: "NST"}

	headerRaw, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("%w: marshal header", ErrAccessTokenInvalid)
	}
	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("%w: marshal payload", ErrAccessTokenInvalid)
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerRaw)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadRaw)
	signingInput := headerEncoded + "." + payloadEncoded

	signature := computeSignature(signingInput, secret)
	return signingInput + "." + signature, nil
}

func ParseAccessToken(token string) (TokenClaims, error) {
	return ParseAccessTokenWithSecret(resolveSecret(), token)
}

func ParseAccessTokenWithSecret(secret string, token string) (TokenClaims, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return TokenClaims{}, ErrAccessTokenSecret
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return TokenClaims{}, ErrAccessTokenMalformed
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return TokenClaims{}, ErrAccessTokenMalformed
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := computeSignature(signingInput, secret)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return TokenClaims{}, ErrAccessTokenInvalid
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return TokenClaims{}, ErrAccessTokenMalformed
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		return TokenClaims{}, ErrAccessTokenMalformed
	}

	claims := TokenClaims{
		UserID:       strings.TrimSpace(payload.UserID),
		CredentialID: strings.TrimSpace(payload.CredentialID),
		Roles:        normalizeRoles(payload.Roles),
		ExpiresAt:    time.Unix(payload.ExpiresAt, 0).UTC(),
		Issuer:       normalizeIssuer(payload.Issuer),
	}
	if err := validateClaimsShape(claims); err != nil {
		return TokenClaims{}, ErrAccessTokenInvalid
	}
	if time.Now().UTC().After(claims.ExpiresAt.UTC()) {
		return TokenClaims{}, ErrAccessTokenExpired
	}

	return claims, nil
}

func validateClaimsShape(claims TokenClaims) error {
	userID := strings.TrimSpace(claims.UserID)
	credentialID := strings.TrimSpace(claims.CredentialID)

	if userID == "" && credentialID == "" {
		return ErrAccessTokenInvalid
	}
	if userID != "" {
		if _, err := uuid.Parse(userID); err != nil {
			return ErrAccessTokenInvalid
		}
	}
	if credentialID != "" {
		if _, err := uuid.Parse(credentialID); err != nil {
			return ErrAccessTokenInvalid
		}
	}
	if claims.ExpiresAt.IsZero() {
		return ErrAccessTokenInvalid
	}
	return nil
}

func normalizeRoles(input []string) []string {
	if len(input) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(input))
	roles := make([]string, 0, len(input))
	for _, role := range input {
		value := normalizeValue(role)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		roles = append(roles, value)
	}
	if len(roles) == 0 {
		return nil
	}
	return roles
}

func normalizeIssuer(issuer string) string {
	issuer = strings.TrimSpace(issuer)
	if issuer == "" {
		return defaultTokenIssuer
	}
	return issuer
}

func computeSignature(signingInput string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func resolveSecret() string {
	accessTokenSecretMu.RLock()
	current := strings.TrimSpace(accessTokenSecret)
	accessTokenSecretMu.RUnlock()
	if current != "" {
		return current
	}

	appEnv := normalizeValue(os.Getenv("APP_ENV"))
	if appEnv == "prod" || appEnv == "production" {
		return ""
	}
	return defaultNonProdAccessTokenSecret
}
