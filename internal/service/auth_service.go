package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	GuestToken(ctx context.Context) (*dto.GuestTokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
	SwitchWorkspace(ctx context.Context, req dto.SwitchWorkspaceRequest) (*dto.SwitchWorkspaceResponse, error)
	Me(ctx context.Context) (*dto.AuthMeResponse, error)
}

type authService struct {
	q         store.Querier
	jwtSecret string
}

func NewAuthService(q store.Querier, jwtSecret string) AuthService {
	return &authService{
		q:         q,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) GuestToken(ctx context.Context) (*dto.GuestTokenResponse, error) {
	guestID := uuid.New()
	access, err := auth.GenerateAccessTokenWithTTL(guestID.String(), types.RoleTrial.String(), nil, s.jwtSecret, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	return &dto.GuestTokenResponse{AccessToken: access}, nil
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	plannerRole, err := s.q.GetRoleByName(ctx, types.RolePlanner.String())
	if err != nil {
		return nil, fmt.Errorf("planner role not found: %w", err)
	}
	ownerRole, err := s.q.GetRoleByName(ctx, types.RoleOwner.String())
	if err != nil {
		return nil, fmt.Errorf("owner role not found: %w", err)
	}
	personalRole, err := s.q.GetRoleByName(ctx, types.RolePersonal.String())
	if err != nil {
		return nil, fmt.Errorf("personal role not found: %w", err)
	}

	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(strings.ToLower(req.Email))

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.q.CreateUser(ctx, store.CreateUserParams{
		RoleID:       plannerRole.RoleID,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accountType := "personal"
	if req.AccountType != nil {
		candidate := strings.TrimSpace(strings.ToLower(*req.AccountType))
		if candidate != "" {
			accountType = candidate
		}
	}
	if accountType != "personal" && accountType != "organization" {
		return nil, fmt.Errorf("invalid account_type: must be 'personal' or 'organization'")
	}

	workspaceName := "my workspace"
	if req.WorkspaceName != nil {
		candidate := strings.TrimSpace(*req.WorkspaceName)
		if candidate != "" {
			workspaceName = candidate
		}
	}

	ws, err := s.q.CreateWorkspace(ctx, store.CreateWorkspaceParams{
		Type:        accountType,
		Name:        workspaceName,
		OwnerUserID: user.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	memberRole := ownerRole
	roleName := types.RoleOwner.String()
	if accountType == "personal" {
		memberRole = personalRole
		roleName = types.RolePersonal.String()
	}

	_, err = s.q.CreateMember(ctx, store.CreateMemberParams{
		WorkspaceID: ws.WorkspaceID,
		UserID:      user.UserID,
		RoleID:      memberRole.RoleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	workspaceID := ws.WorkspaceID.String()

	if req.GuestToken != nil && *req.GuestToken != "" {
		if err := s.claimGuestPlansWithWorkspace(ctx, *req.GuestToken, user.UserID.String(), &ws.WorkspaceID); err != nil {
			return nil, err
		}
	}

	access, err := auth.GenerateAccessToken(user.UserID.String(), roleName, &workspaceID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshExpires := time.Now().Add(30 * 24 * time.Hour)

	if err = s.q.CreateRefreshToken(ctx, store.CreateRefreshTokenParams{
		Token:       refresh,
		UserID:      user.UserID,
		WorkspaceID: &ws.WorkspaceID,
		ExpiresAt:   &refreshExpires,
	}); err != nil {
		return nil, fmt.Errorf("failed to create refresh token in store: %w", err)
	}

	return &dto.RegisterResponse{
		AccessToken:       access,
		RefreshToken:      refresh,
		ActiveWorkspaceID: &workspaceID,
		User: dto.UserSummary{
			ID:       user.UserID.String(),
			Username: user.Username,
			Role:     roleName,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.q.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	workspaceID, err := s.resolveDefaultWorkspaceID(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	roleName, err := s.resolveRoleName(ctx, user.UserID, workspaceID)
	if err != nil {
		return nil, err
	}

	workspaceIDStr := workspaceID.String()
	access, err := auth.GenerateAccessToken(user.UserID.String(), roleName, &workspaceIDStr, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	if err = s.q.CreateRefreshToken(ctx, store.CreateRefreshTokenParams{
		Token:       refresh,
		UserID:      user.UserID,
		WorkspaceID: &workspaceID,
		ExpiresAt:   &expiresAt,
	}); err != nil {
		return nil, fmt.Errorf("failed to create refresh token in store: %w", err)
	}

	if req.GuestToken != nil && *req.GuestToken != "" {
		if err := s.claimGuestPlansWithWorkspace(ctx, *req.GuestToken, user.UserID.String(), &workspaceID); err != nil {
			return nil, err
		}
	}

	return &dto.LoginResponse{
		AccessToken:       access,
		RefreshToken:      refresh,
		ActiveWorkspaceID: &workspaceIDStr,
		User: dto.UserSummary{
			ID:       user.UserID.String(),
			Username: user.Username,
			Role:     roleName,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, oldToken string) (*dto.LoginResponse, error) {
	tokenInfo, err := s.q.GetRefreshToken(ctx, oldToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !tokenInfo.RevokedAt.IsZero() {
		return nil, fmt.Errorf("refresh token revoked")
	}
	if tokenInfo.ExpiresAt != nil && tokenInfo.ExpiresAt.Before(time.Now()) {
		_ = s.q.RevokeRefreshToken(ctx, oldToken)
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.q.GetUserByID(ctx, tokenInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	activeWorkspaceID, err := s.resolveRefreshWorkspaceID(ctx, tokenInfo.UserID, tokenInfo.WorkspaceID)
	if err != nil {
		return nil, err
	}

	roleName, err := s.resolveRoleName(ctx, tokenInfo.UserID, activeWorkspaceID)
	if err != nil {
		return nil, err
	}

	if err := s.q.RevokeRefreshToken(ctx, oldToken); err != nil {
		return nil, fmt.Errorf("failed to revoke token: %w", err)
	}

	workspaceIDStr := activeWorkspaceID.String()
	newAccess, err := auth.GenerateAccessToken(user.UserID.String(), roleName, &workspaceIDStr, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	if err = s.q.CreateRefreshToken(ctx, store.CreateRefreshTokenParams{
		Token:       newRefresh,
		UserID:      user.UserID,
		WorkspaceID: &activeWorkspaceID,
		ExpiresAt:   &expiresAt,
	}); err != nil {
		return nil, fmt.Errorf("failed to create refresh token in store: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:       newAccess,
		RefreshToken:      newRefresh,
		ActiveWorkspaceID: &workspaceIDStr,
		User: dto.UserSummary{
			ID:       user.UserID.String(),
			Username: user.Username,
			Role:     roleName,
		},
	}, nil
}

func (s *authService) SwitchWorkspace(ctx context.Context, req dto.SwitchWorkspaceRequest) (*dto.SwitchWorkspaceResponse, error) {
	userIDStr, ok := auth.UserIDFromContext(ctx)
	if !ok || userIDStr == "" {
		return nil, fmt.Errorf("missing user id")
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	wsUUID, err := uuid.Parse(req.WorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id")
	}

	// Ensure this user can act in the target workspace.
	roleName, err := s.resolveRoleName(ctx, userUUID, wsUUID)
	if err != nil {
		return nil, err
	}

	// Ensure refresh token belongs to this user.
	tok, err := s.q.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	if tok.UserID != userUUID {
		return nil, fmt.Errorf("refresh token does not belong to user")
	}

	if err := s.q.UpdateRefreshTokenWorkspace(ctx, store.UpdateRefreshTokenWorkspaceParams{Token: req.RefreshToken, WorkspaceID: &wsUUID}); err != nil {
		return nil, fmt.Errorf("failed to update refresh token workspace: %w", err)
	}

	wsIDStr := wsUUID.String()
	newAccess, err := auth.GenerateAccessToken(userIDStr, roleName, &wsIDStr, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &dto.SwitchWorkspaceResponse{AccessToken: newAccess, ActiveWorkspaceID: wsIDStr}, nil
}

func (s *authService) resolveRefreshWorkspaceID(ctx context.Context, userID uuid.UUID, fromToken *uuid.UUID) (uuid.UUID, error) {
	if fromToken != nil {
		return *fromToken, nil
	}
	return s.resolveDefaultWorkspaceID(ctx, userID)
}

func (s *authService) resolveDefaultWorkspaceID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	ws, err := s.q.GetPersonalWorkspaceByOwner(ctx, userID)
	if err == nil {
		return ws.WorkspaceID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("failed to get personal workspace: %w", err)
	}

	owned, err := s.q.ListWorkspacesByOwner(ctx, store.ListWorkspacesByOwnerParams{OwnerUserID: userID, Limit: 1, Offset: 0})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to list owned workspaces: %w", err)
	}
	if len(owned) > 0 {
		return owned[0].WorkspaceID, nil
	}

	workspaces, err := s.q.ListWorkspacesForUser(ctx, store.ListWorkspacesForUserParams{UserID: userID, Limit: 1, Offset: 0})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to list workspaces: %w", err)
	}
	if len(workspaces) == 0 {
		return uuid.Nil, fmt.Errorf("no workspace for user")
	}
	return workspaces[0].WorkspaceID, nil
}

func (s *authService) resolveRoleName(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) (string, error) {
	platformRole, err := s.q.GetPlatformRoleByUserID(ctx, userID)
	if err == nil {
		return platformRole, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to get platform role: %w", err)
	}

	roleName, err := s.q.GetMemberRoleNameByWorkspaceAndUser(ctx, store.GetMemberRoleNameByWorkspaceAndUserParams{WorkspaceID: workspaceID, UserID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("not a workspace member")
		}
		return "", fmt.Errorf("failed to resolve member role: %w", err)
	}
	return roleName, nil
}

func (s *authService) Me(ctx context.Context) (*dto.AuthMeResponse, error) {
	userIDStr, ok := auth.UserIDFromContext(ctx)
	if !ok || userIDStr == "" {
		return nil, fmt.Errorf("missing user id")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	roleName, ok := auth.RoleFromContext(ctx)
	if !ok || roleName == "" {
		return nil, fmt.Errorf("missing role")
	}

	user, err := s.q.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	permissions, err := s.q.GetPermissionsByRole(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	isPlatformMember := false
	if _, err := s.q.GetPlatformRoleByUserID(ctx, userID); err == nil {
		isPlatformMember = true
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check platform membership: %w", err)
	}

	var workspaceID *string
	if wsID, ok := auth.WorkspaceIDFromContext(ctx); ok && wsID != "" {
		workspaceID = &wsID
	}

	return &dto.AuthMeResponse{
		User: dto.UserSummary{
			ID:       user.UserID.String(),
			Username: user.Username,
			Role:     roleName,
		},
		ActiveWorkspaceID: workspaceID,
		Permissions:       permissions,
		IsPlatformMember:  isPlatformMember,
	}, nil
}
