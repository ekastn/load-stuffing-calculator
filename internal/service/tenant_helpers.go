package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func isFounder(ctx context.Context) bool {
	role, ok := auth.RoleFromContext(ctx)
	return ok && role == types.RoleFounder.String()
}

func workspaceOverrideIDFromContext(ctx context.Context) (*uuid.UUID, error) {
	wid, ok := auth.WorkspaceOverrideIDFromContext(ctx)
	if !ok || wid == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(wid)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id: %w", err)
	}
	return &parsed, nil
}

func userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := auth.UserIDFromContext(ctx)
	if !ok || userIDStr == "" {
		return uuid.Nil, fmt.Errorf("missing user id")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}
	return userID, nil
}

func parseOptionalWorkspaceID(workspaceID *string) (*uuid.UUID, error) {
	if workspaceID == nil || *workspaceID == "" {
		return nil, nil
	}
	wid, err := uuid.Parse(*workspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id")
	}
	return &wid, nil
}

func activeOrOverrideWorkspaceID(ctx context.Context, overrideWorkspaceID *string) (*uuid.UUID, error) {
	if isFounder(ctx) {
		if wid, err := parseOptionalWorkspaceID(overrideWorkspaceID); err != nil {
			return nil, err
		} else if wid != nil {
			return wid, nil
		}
	}
	return workspaceIDFromContext(ctx)
}

func ensureWorkspaceExists(ctx context.Context, q store.Querier, workspaceID uuid.UUID) (store.Workspace, error) {
	ws, err := q.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return store.Workspace{}, fmt.Errorf("workspace not found: %w", err)
	}
	return ws, nil
}

func ensureNotPersonalWorkspace(ws store.Workspace) error {
	if ws.Type == "personal" {
		return fmt.Errorf("operation not allowed for personal workspace")
	}
	return nil
}

func ensureWorkspaceOwnerOrFounder(ctx context.Context, q store.Querier, workspace store.Workspace) error {
	if isFounder(ctx) {
		return nil
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return err
	}
	if workspace.OwnerUserID != userID {
		return fmt.Errorf("forbidden")
	}
	return nil
}

func ensureWorkspaceAdminOrOwnerOrFounder(ctx context.Context, q store.Querier, workspaceID uuid.UUID) (string, error) {
	if isFounder(ctx) {
		return types.RoleFounder.String(), nil
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return "", err
	}
	roleName, err := q.GetMemberRoleNameByWorkspaceAndUser(ctx, store.GetMemberRoleNameByWorkspaceAndUserParams{WorkspaceID: workspaceID, UserID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("forbidden")
		}
		return "", fmt.Errorf("failed to resolve role: %w", err)
	}
	switch roleName {
	case types.RoleOwner.String(), types.RoleAdmin.String():
		return roleName, nil
	default:
		return "", fmt.Errorf("forbidden")
	}
}

func lookupRoleID(ctx context.Context, q store.Querier, roleName string) (uuid.UUID, error) {
	if !types.IsWorkspaceRole(roleName) {
		return uuid.Nil, fmt.Errorf("invalid role")
	}
	role, err := q.GetRoleByName(ctx, roleName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("role not found: %w", err)
	}
	return role.RoleID, nil
}

func parseUserIdentifier(identifier string) (id *uuid.UUID, username *string, email *string) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, nil, nil
	}
	if u, err := uuid.Parse(identifier); err == nil {
		return &u, nil, nil
	}
	if strings.Contains(identifier, "@") {
		return nil, nil, &identifier
	}
	return nil, &identifier, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
