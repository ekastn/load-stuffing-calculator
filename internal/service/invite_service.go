package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

type InviteService interface {
	ListInvites(ctx context.Context, page, limit int32, overrideWorkspaceID *string) ([]dto.InviteResponse, error)
	CreateInvite(ctx context.Context, req dto.CreateInviteRequest, overrideWorkspaceID *string) (*dto.CreateInviteResponse, error)
	RevokeInvite(ctx context.Context, inviteID string, overrideWorkspaceID *string) error
	AcceptInvite(ctx context.Context, req dto.AcceptInviteRequest) (*dto.AcceptInviteResponse, error)
}

type inviteService struct {
	q         store.Querier
	jwtSecret string
}

func NewInviteService(q store.Querier, jwtSecret string) InviteService {
	return &inviteService{q: q, jwtSecret: jwtSecret}
}

func (s *inviteService) ListInvites(ctx context.Context, page, limit int32, overrideWorkspaceID *string) ([]dto.InviteResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	workspaceID, err := activeOrOverrideWorkspaceID(ctx, overrideWorkspaceID)
	if err != nil {
		return nil, err
	}
	if workspaceID == nil {
		return nil, fmt.Errorf("missing workspace")
	}

	if _, err := ensureWorkspaceAdminOrOwnerOrFounder(ctx, s.q, *workspaceID); err != nil {
		return nil, err
	}

	rows, err := s.q.ListInvitesByWorkspace(ctx, store.ListInvitesByWorkspaceParams{WorkspaceID: *workspaceID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, fmt.Errorf("failed to list invites: %w", err)
	}

	resp := make([]dto.InviteResponse, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, mapInviteRow(row))
	}
	return resp, nil
}

func (s *inviteService) CreateInvite(ctx context.Context, req dto.CreateInviteRequest, overrideWorkspaceID *string) (*dto.CreateInviteResponse, error) {
	workspaceID, err := activeOrOverrideWorkspaceID(ctx, overrideWorkspaceID)
	if err != nil {
		return nil, err
	}
	if workspaceID == nil {
		return nil, fmt.Errorf("missing workspace")
	}

	ws, err := ensureWorkspaceExists(ctx, s.q, *workspaceID)
	if err != nil {
		return nil, err
	}
	if err := ensureNotPersonalWorkspace(ws); err != nil {
		return nil, err
	}

	if _, err := ensureWorkspaceAdminOrOwnerOrFounder(ctx, s.q, *workspaceID); err != nil {
		return nil, err
	}

	if !types.IsAssignableWorkspaceRole(req.Role) {
		return nil, fmt.Errorf("invalid role")
	}
	roleID, err := lookupRoleID(ctx, s.q, req.Role)
	if err != nil {
		return nil, err
	}

	invitedByUserID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Create a random raw token; store only hash.
	rawToken, err := randomInviteToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	th := hashToken(rawToken)

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	invite, err := s.q.CreateInvite(ctx, store.CreateInviteParams{
		WorkspaceID:     *workspaceID,
		Email:           req.Email,
		RoleID:          roleID,
		TokenHash:       th,
		InvitedByUserID: invitedByUserID,
		ExpiresAt:       &expiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	// Best-effort for invited by username in response.
	invitedByUsername := ""
	if u, err := s.q.GetUserByID(ctx, invitedByUserID); err == nil {
		invitedByUsername = u.Username
	}

	resp := dto.InviteResponse{
		InviteID:          invite.InviteID.String(),
		WorkspaceID:       invite.WorkspaceID.String(),
		Email:             invite.Email,
		Role:              req.Role,
		InvitedByUserID:   invite.InvitedByUserID.String(),
		InvitedByUsername: invitedByUsername,
		ExpiresAt:         invite.ExpiresAt,
		AcceptedAt:        invite.AcceptedAt,
		RevokedAt:         invite.RevokedAt,
		CreatedAt:         invite.CreatedAt,
	}

	return &dto.CreateInviteResponse{Invite: resp, Token: rawToken}, nil
}

func (s *inviteService) RevokeInvite(ctx context.Context, inviteID string, overrideWorkspaceID *string) error {
	workspaceID, err := activeOrOverrideWorkspaceID(ctx, overrideWorkspaceID)
	if err != nil {
		return err
	}
	if workspaceID == nil {
		return fmt.Errorf("missing workspace")
	}

	ws, err := ensureWorkspaceExists(ctx, s.q, *workspaceID)
	if err != nil {
		return err
	}
	if err := ensureNotPersonalWorkspace(ws); err != nil {
		return err
	}

	if _, err := ensureWorkspaceAdminOrOwnerOrFounder(ctx, s.q, *workspaceID); err != nil {
		return err
	}

	iID, err := uuid.Parse(inviteID)
	if err != nil {
		return fmt.Errorf("invalid invite id")
	}

	if err := s.q.RevokeInvite(ctx, store.RevokeInviteParams{InviteID: iID, WorkspaceID: *workspaceID}); err != nil {
		return fmt.Errorf("failed to revoke invite: %w", err)
	}
	return nil
}

func (s *inviteService) AcceptInvite(ctx context.Context, req dto.AcceptInviteRequest) (*dto.AcceptInviteResponse, error) {
	// MVP requires login.
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := s.q.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	inv, err := s.q.GetInviteByTokenHash(ctx, hashToken(req.Token))
	if err != nil {
		return nil, fmt.Errorf("invalid invite")
	}

	if user.Email != inv.Email {
		return nil, fmt.Errorf("invite email mismatch")
	}

	ws, err := ensureWorkspaceExists(ctx, s.q, inv.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if err := ensureNotPersonalWorkspace(ws); err != nil {
		return nil, err
	}

	// Determine role name from invite role_id.
	roleName := ""
	roleRow, err := s.q.GetRole(ctx, inv.RoleID)
	if err != nil {
		return nil, fmt.Errorf("invalid invite role")
	}
	roleName = roleRow.Name
	if roleName == types.RoleFounder.String() || roleName == types.RoleOwner.String() {
		return nil, fmt.Errorf("invalid invite role")
	}

	// Ensure membership exists.
	_, err = s.q.GetMemberByWorkspaceAndUser(ctx, store.GetMemberByWorkspaceAndUserParams{WorkspaceID: inv.WorkspaceID, UserID: userID})
	if err != nil {
		_, createErr := s.q.CreateMember(ctx, store.CreateMemberParams{WorkspaceID: inv.WorkspaceID, UserID: userID, RoleID: inv.RoleID})
		if createErr != nil {
			if isUniqueViolation(createErr) {
				// already member
			} else {
				return nil, fmt.Errorf("failed to create membership: %w", createErr)
			}
		}
	}

	if err := s.q.AcceptInvite(ctx, store.AcceptInviteParams{InviteID: inv.InviteID, WorkspaceID: inv.WorkspaceID}); err != nil {
		return nil, fmt.Errorf("failed to accept invite: %w", err)
	}

	wsIDStr := inv.WorkspaceID.String()
	access, err := auth.GenerateAccessToken(userID.String(), roleName, &wsIDStr, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &dto.AcceptInviteResponse{AccessToken: access, ActiveWorkspaceID: wsIDStr, Role: roleName}, nil
}

func mapInviteRow(row store.ListInvitesByWorkspaceRow) dto.InviteResponse {
	return dto.InviteResponse{
		InviteID:          row.InviteID.String(),
		WorkspaceID:       row.WorkspaceID.String(),
		Email:             row.Email,
		Role:              row.RoleName,
		InvitedByUserID:   row.InvitedByUserID.String(),
		InvitedByUsername: row.InvitedByUsername,
		ExpiresAt:         row.ExpiresAt,
		AcceptedAt:        row.AcceptedAt,
		RevokedAt:         row.RevokedAt,
		CreatedAt:         row.CreatedAt,
	}
}

func randomInviteToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// hex string; 2 chars per byte.
	return fmt.Sprintf("%x", b), nil
}
