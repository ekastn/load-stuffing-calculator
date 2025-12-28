package auth

import "context"

type contextKey string

const (
	ContextKeyUserID              contextKey = "user_id"
	ContextKeyRole                contextKey = "role"
	ContextKeyWorkspaceID         contextKey = "workspace_id"
	ContextKeyWorkspaceOverrideID contextKey = "workspace_override_id"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ContextKeyRole, role)
}

func WithWorkspaceID(ctx context.Context, workspaceID string) context.Context {
	return context.WithValue(ctx, ContextKeyWorkspaceID, workspaceID)
}

func WithWorkspaceOverrideID(ctx context.Context, workspaceID string) context.Context {
	return context.WithValue(ctx, ContextKeyWorkspaceOverrideID, workspaceID)
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

func WorkspaceIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ContextKeyWorkspaceID)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func WorkspaceOverrideIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ContextKeyWorkspaceOverrideID)
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
