package identity

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	HeaderAuthorization = "Authorization"
)

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := ensureIdentityContext(w, r)
		if !ok {
			return
		}
		if _, hasUser := UserID(ctx); !hasUser {
			writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_user_id")
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := ensureIdentityContext(w, r)
		if !ok {
			return
		}
		if _, hasCredential := CredentialID(ctx); !hasCredential {
			writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_credential_id")
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAnyRole(requiredRoles ...string) func(http.Handler) http.Handler {
	normalized := make([]string, 0, len(requiredRoles))
	for _, role := range requiredRoles {
		role = normalizeValue(role)
		if role == "" {
			continue
		}
		normalized = append(normalized, role)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, ok := ensureIdentityContext(w, r)
			if !ok {
				return
			}

			roles := rolesFromContext(ctx)
			if len(roles) == 0 {
				writeError(w, http.StatusForbidden, "missing_actor_roles")
				return
			}

			allowed := false
			for _, role := range normalized {
				if _, exists := roles[role]; exists {
					allowed = true
					break
				}
			}
			if !allowed {
				writeError(w, http.StatusForbidden, "insufficient_actor_role")
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ensureIdentityContext(w http.ResponseWriter, r *http.Request) (context.Context, bool) {
	ctx := r.Context()
	if hasIdentityContext(ctx) {
		return ctx, true
	}

	token, parseErr := bearerToken(r.Header.Get(HeaderAuthorization))
	if parseErr != "" {
		writeError(w, http.StatusUnauthorized, parseErr)
		return nil, false
	}

	claims, err := ParseAccessToken(token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_access_token")
		return nil, false
	}

	if claims.UserID != "" {
		ctx = withUserID(ctx, claims.UserID)
	}
	if claims.CredentialID != "" {
		ctx = withCredentialID(ctx, claims.CredentialID)
	}
	if roles := parseRoles(strings.Join(claims.Roles, ",")); len(roles) > 0 {
		ctx = withRoles(ctx, roles)
	}
	return ctx, true
}

func hasIdentityContext(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	if value, exists := ctx.Value(userIDKey).(string); exists && value != "" {
		return true
	}
	if value, exists := ctx.Value(credentialIDKey).(string); exists && value != "" {
		return true
	}
	if roles, exists := ctx.Value(rolesKey).(map[string]struct{}); exists && len(roles) > 0 {
		return true
	}
	return false
}

func rolesFromContext(ctx context.Context) map[string]struct{} {
	if ctx == nil {
		return nil
	}
	roles, _ := ctx.Value(rolesKey).(map[string]struct{})
	return roles
}

func bearerToken(value string) (string, string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", "missing_authorization_header"
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(strings.TrimSpace(parts[0]), "Bearer") {
		return "", "invalid_authorization_header"
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", "invalid_authorization_header"
	}
	return token, ""
}

func parseRoles(raw string) map[string]struct{} {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	roles := make(map[string]struct{})
	parts := strings.Split(raw, ",")
	for _, part := range parts {
		role := normalizeValue(part)
		if role == "" {
			continue
		}
		roles[role] = struct{}{}
	}
	if len(roles) == 0 {
		return nil
	}
	return roles
}

func normalizeValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
