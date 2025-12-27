package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

type MemberService interface {
	ListMembers(ctx context.Context, page, limit int32, overrideWorkspaceID *string) ([]dto.MemberResponse, error)
	AddMember(ctx context.Context, req dto.AddMemberRequest, overrideWorkspaceID *string) (*dto.MemberResponse, error)
	UpdateMemberRole(ctx context.Context, memberID string, req dto.UpdateMemberRoleRequest, overrideWorkspaceID *string) (*dto.MemberResponse, error)
	DeleteMember(ctx context.Context, memberID string, overrideWorkspaceID *string) error
}

type memberService struct {
	q store.Querier
}

func NewMemberService(q store.Querier) MemberService {
	return &memberService{q: q}
}

func (s *memberService) ListMembers(ctx context.Context, page, limit int32, overrideWorkspaceID *string) ([]dto.MemberResponse, error) {
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

	rows, err := s.q.ListMembersByWorkspace(ctx, store.ListMembersByWorkspaceParams{WorkspaceID: *workspaceID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}

	resp := make([]dto.MemberResponse, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, dto.MemberResponse{
			MemberID:    row.MemberID.String(),
			WorkspaceID: row.WorkspaceID.String(),
			UserID:      row.UserID.String(),
			Role:        row.RoleName,
			Username:    row.Username,
			Email:       row.Email,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return resp, nil
}

func (s *memberService) AddMember(ctx context.Context, req dto.AddMemberRequest, overrideWorkspaceID *string) (*dto.MemberResponse, error) {
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

	id, username, email := parseUserIdentifier(req.UserIdentifier)
	var user store.GetUserByUsernameRow
	switch {
	case id != nil:
		row, err := s.q.GetUserByID(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}
		user = store.GetUserByUsernameRow{UserID: row.UserID, Username: row.Username, Email: row.Email, PasswordHash: "", RoleID: row.RoleID, RoleName: row.RoleName}
	case email != nil:
		row, err := s.q.GetUserByEmail(ctx, *email)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}
		user = store.GetUserByUsernameRow{UserID: row.UserID, Username: row.Username, Email: row.Email, PasswordHash: row.PasswordHash, RoleID: row.RoleID, RoleName: row.RoleName}
	case username != nil:
		row, err := s.q.GetUserByUsername(ctx, *username)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}
		user = row
	default:
		return nil, fmt.Errorf("invalid user identifier")
	}

	member, err := s.q.CreateMember(ctx, store.CreateMemberParams{WorkspaceID: *workspaceID, UserID: user.UserID, RoleID: roleID})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("already a member")
		}
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	resp := &dto.MemberResponse{
		MemberID:    member.MemberID.String(),
		WorkspaceID: member.WorkspaceID.String(),
		UserID:      member.UserID.String(),
		Role:        req.Role,
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
	}
	return resp, nil
}

func (s *memberService) UpdateMemberRole(ctx context.Context, memberID string, req dto.UpdateMemberRoleRequest, overrideWorkspaceID *string) (*dto.MemberResponse, error) {
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

	mID, err := uuid.Parse(memberID)
	if err != nil {
		return nil, fmt.Errorf("invalid member id")
	}

	member, err := s.q.GetMember(ctx, mID)
	if err != nil {
		return nil, fmt.Errorf("member not found")
	}
	if member.WorkspaceID != *workspaceID {
		return nil, fmt.Errorf("member not found")
	}

	if member.UserID == ws.OwnerUserID {
		return nil, fmt.Errorf("cannot change owner role")
	}

	if !types.IsAssignableWorkspaceRole(req.Role) {
		return nil, fmt.Errorf("invalid role")
	}
	roleID, err := lookupRoleID(ctx, s.q, req.Role)
	if err != nil {
		return nil, err
	}

	if err := s.q.UpdateMemberRole(ctx, store.UpdateMemberRoleParams{MemberID: member.MemberID, WorkspaceID: *workspaceID, RoleID: roleID}); err != nil {
		return nil, fmt.Errorf("failed to update member role: %w", err)
	}

	// Fetch user info for response.
	userRow, err := s.q.GetUserByID(ctx, member.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user")
	}

	member, err = s.q.GetMember(ctx, mID)
	if err != nil {
		return nil, fmt.Errorf("member not found")
	}

	resp := &dto.MemberResponse{
		MemberID:    member.MemberID.String(),
		WorkspaceID: member.WorkspaceID.String(),
		UserID:      member.UserID.String(),
		Role:        req.Role,
		Username:    userRow.Username,
		Email:       userRow.Email,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
	}
	return resp, nil
}

func (s *memberService) DeleteMember(ctx context.Context, memberID string, overrideWorkspaceID *string) error {
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

	mID, err := uuid.Parse(memberID)
	if err != nil {
		return fmt.Errorf("invalid member id")
	}

	member, err := s.q.GetMember(ctx, mID)
	if err != nil {
		return fmt.Errorf("member not found")
	}
	if member.WorkspaceID != *workspaceID {
		return fmt.Errorf("member not found")
	}

	roleName, err := s.q.GetMemberRoleNameByWorkspaceAndUser(ctx, store.GetMemberRoleNameByWorkspaceAndUserParams{WorkspaceID: *workspaceID, UserID: member.UserID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("member not found")
		}
		return fmt.Errorf("failed to resolve member role: %w", err)
	}
	if roleName == types.RoleOwner.String() || member.UserID == ws.OwnerUserID {
		return fmt.Errorf("cannot remove owner")
	}

	if err := s.q.DeleteMember(ctx, store.DeleteMemberParams{MemberID: member.MemberID, WorkspaceID: *workspaceID}); err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}
