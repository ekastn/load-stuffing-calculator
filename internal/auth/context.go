package auth

import "context"

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ContextKeyRole, role)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ContextKeyUserID)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ContextKeyRole)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
