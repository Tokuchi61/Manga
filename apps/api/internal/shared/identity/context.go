package identity

import "context"

type contextKey string

const (
	userIDKey       contextKey = "actor_user_id"
	credentialIDKey contextKey = "actor_credential_id"
	rolesKey        contextKey = "actor_roles"
)

func withUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func withCredentialID(ctx context.Context, credentialID string) context.Context {
	return context.WithValue(ctx, credentialIDKey, credentialID)
}

func withRoles(ctx context.Context, roles map[string]struct{}) context.Context {
	return context.WithValue(ctx, rolesKey, roles)
}

func UserID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	value, ok := ctx.Value(userIDKey).(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}

func CredentialID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	value, ok := ctx.Value(credentialIDKey).(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}

func HasRole(ctx context.Context, role string) bool {
	if ctx == nil || role == "" {
		return false
	}
	roleSet, ok := ctx.Value(rolesKey).(map[string]struct{})
	if !ok || len(roleSet) == 0 {
		return false
	}
	_, exists := roleSet[normalizeValue(role)]
	return exists
}

func HasAnyRole(ctx context.Context, roles ...string) bool {
	for _, role := range roles {
		if HasRole(ctx, role) {
			return true
		}
	}
	return false
}
