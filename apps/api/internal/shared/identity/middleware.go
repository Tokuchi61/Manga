package identity

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	HeaderActorUserID       = "X-Actor-User-ID"
	HeaderActorCredentialID = "X-Actor-Credential-ID"
	HeaderActorRoles        = "X-Actor-Roles"
)

func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := strings.TrimSpace(r.Header.Get(HeaderActorUserID))
		if _, err := uuid.Parse(userID); err != nil {
			writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_user_id")
			return
		}

		ctx := withUserID(r.Context(), userID)
		if roles := parseRoles(r.Header.Get(HeaderActorRoles)); len(roles) > 0 {
			ctx = withRoles(ctx, roles)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentialID := strings.TrimSpace(r.Header.Get(HeaderActorCredentialID))
		if _, err := uuid.Parse(credentialID); err != nil {
			writeError(w, http.StatusUnauthorized, "missing_or_invalid_actor_credential_id")
			return
		}

		ctx := withCredentialID(r.Context(), credentialID)
		if userID := strings.TrimSpace(r.Header.Get(HeaderActorUserID)); userID != "" {
			if _, err := uuid.Parse(userID); err == nil {
				ctx = withUserID(ctx, userID)
			}
		}
		if roles := parseRoles(r.Header.Get(HeaderActorRoles)); len(roles) > 0 {
			ctx = withRoles(ctx, roles)
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
			roles := parseRoles(r.Header.Get(HeaderActorRoles))
			if len(roles) == 0 {
				writeError(w, http.StatusForbidden, "missing_actor_roles")
				return
			}

			ctx := withRoles(r.Context(), roles)
			allowed := false
			for _, role := range normalized {
				if _, ok := roles[role]; ok {
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
