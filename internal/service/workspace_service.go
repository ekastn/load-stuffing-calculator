package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

type WorkspaceService interface {
	ListWorkspaces(ctx context.Context, page, limit int32) ([]dto.WorkspaceResponse, error)
	CreateWorkspace(ctx context.Context, req dto.CreateWorkspaceRequest) (*dto.WorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, id string, req dto.UpdateWorkspaceRequest) (*dto.WorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, id string) error
}

type workspaceService struct {
	q store.Querier
}

func NewWorkspaceService(q store.Querier) WorkspaceService {
	return &workspaceService{q: q}
}

func (s *workspaceService) ListWorkspaces(ctx context.Context, page, limit int32) ([]dto.WorkspaceResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	if isFounder(ctx) {
		rows, err := s.q.ListWorkspacesAll(ctx, store.ListWorkspacesAllParams{Limit: limit, Offset: offset})
		if err != nil {
			return nil, fmt.Errorf("failed to list workspaces: %w", err)
		}

		resp := make([]dto.WorkspaceResponse, 0, len(rows))
		for _, row := range rows {
			resp = append(resp, mapWorkspaceAllRow(row))
		}
		return resp, nil
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	workspaces, err := s.q.ListWorkspacesForUser(ctx, store.ListWorkspacesForUserParams{UserID: userID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}

	resp := make([]dto.WorkspaceResponse, 0, len(workspaces))
	for _, ws := range workspaces {
		resp = append(resp, mapWorkspace(ws))
	}
	return resp, nil
}

func (s *workspaceService) CreateWorkspace(ctx context.Context, req dto.CreateWorkspaceRequest) (*dto.WorkspaceResponse, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	workspaceType := "organization"
	ownerUserID := userID

	if isFounder(ctx) {
		if req.Type != nil {
			workspaceType = strings.ToLower(strings.TrimSpace(*req.Type))
			if workspaceType != "personal" && workspaceType != "organization" {
				return nil, fmt.Errorf("invalid workspace type")
			}
		}
		if req.OwnerUserID != nil {
			parsed, err := uuid.Parse(*req.OwnerUserID)
			if err != nil {
				return nil, fmt.Errorf("invalid owner_user_id")
			}
			ownerUserID = parsed
		}
	} else {
		// Non-platform users cannot escalate ownership or create personal workspaces.
		if req.Type != nil || req.OwnerUserID != nil {
			return nil, fmt.Errorf("forbidden")
		}
	}

	roleName := types.RoleOwner.String()
	if workspaceType == "personal" {
		roleName = types.RolePersonal.String()
	}
	roleID, err := lookupRoleID(ctx, s.q, roleName)
	if err != nil {
		return nil, err
	}

	ws, err := s.q.CreateWorkspace(ctx, store.CreateWorkspaceParams{Type: workspaceType, Name: req.Name, OwnerUserID: ownerUserID})
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Only add the workspace owner as a member.
	_, err = s.q.CreateMember(ctx, store.CreateMemberParams{WorkspaceID: ws.WorkspaceID, UserID: ownerUserID, RoleID: roleID})
	if err != nil {
		return nil, fmt.Errorf("failed to create owner membership: %w", err)
	}

	resp := mapWorkspace(ws)
	return &resp, nil
}

func (s *workspaceService) UpdateWorkspace(ctx context.Context, id string, req dto.UpdateWorkspaceRequest) (*dto.WorkspaceResponse, error) {
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id")
	}

	ws, err := ensureWorkspaceExists(ctx, s.q, workspaceID)
	if err != nil {
		return nil, err
	}

	// Only owner/founder can rename or transfer ownership.
	if err := ensureWorkspaceOwnerOrFounder(ctx, s.q, ws); err != nil {
		return nil, err
	}

	if req.Name != nil {
		if err := s.q.UpdateWorkspace(ctx, store.UpdateWorkspaceParams{WorkspaceID: ws.WorkspaceID, Name: *req.Name}); err != nil {
			return nil, fmt.Errorf("failed to update workspace: %w", err)
		}
	}

	if req.OwnerUserID != nil {
		newOwnerID, err := uuid.Parse(*req.OwnerUserID)
		if err != nil {
			return nil, fmt.Errorf("invalid owner_user_id")
		}

		// Ensure new owner is a member (or create as owner).
		_, memberErr := s.q.GetMemberByWorkspaceAndUser(ctx, store.GetMemberByWorkspaceAndUserParams{WorkspaceID: ws.WorkspaceID, UserID: newOwnerID})
		if memberErr != nil {
			roleName := types.RoleOwner.String()
			if ws.Type == "personal" {
				roleName = types.RolePersonal.String()
			}
			roleID, err := lookupRoleID(ctx, s.q, roleName)
			if err != nil {
				return nil, err
			}
			if _, err := s.q.CreateMember(ctx, store.CreateMemberParams{WorkspaceID: ws.WorkspaceID, UserID: newOwnerID, RoleID: roleID}); err != nil {
				return nil, fmt.Errorf("failed to add new owner as member: %w", err)
			}
		}

		if err := s.q.TransferWorkspaceOwnership(ctx, store.TransferWorkspaceOwnershipParams{WorkspaceID: ws.WorkspaceID, OwnerUserID: newOwnerID}); err != nil {
			return nil, fmt.Errorf("failed to transfer ownership: %w", err)
		}
	}

	ws, err = ensureWorkspaceExists(ctx, s.q, ws.WorkspaceID)
	if err != nil {
		return nil, err
	}

	resp := mapWorkspace(ws)
	return &resp, nil
}

func (s *workspaceService) DeleteWorkspace(ctx context.Context, id string) error {
	workspaceID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid workspace id")
	}

	ws, err := ensureWorkspaceExists(ctx, s.q, workspaceID)
	if err != nil {
		return err
	}

	if ws.Type == "personal" && !isFounder(ctx) {
		return fmt.Errorf("cannot delete personal workspace")
	}

	if err := ensureWorkspaceOwnerOrFounder(ctx, s.q, ws); err != nil {
		return err
	}

	if err := s.q.DeleteWorkspace(ctx, ws.WorkspaceID); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}
	return nil
}

func mapWorkspace(ws store.Workspace) dto.WorkspaceResponse {
	return dto.WorkspaceResponse{
		WorkspaceID: ws.WorkspaceID.String(),
		Type:        ws.Type,
		Name:        ws.Name,
		OwnerUserID: ws.OwnerUserID.String(),
		CreatedAt:   ws.CreatedAt,
		UpdatedAt:   ws.UpdatedAt,
	}
}

func mapWorkspaceAllRow(row store.ListWorkspacesAllRow) dto.WorkspaceResponse {
	ownerUsername := row.OwnerUsername
	ownerEmail := row.OwnerEmail

	return dto.WorkspaceResponse{
		WorkspaceID:   row.WorkspaceID.String(),
		Type:          row.Type,
		Name:          row.Name,
		OwnerUserID:   row.OwnerUserID.String(),
		OwnerUsername: &ownerUsername,
		OwnerEmail:    &ownerEmail,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
