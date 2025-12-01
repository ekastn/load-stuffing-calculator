package service

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

type PermissionService interface {
	CreatePermission(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error)
	GetPermission(ctx context.Context, id string) (*dto.PermissionResponse, error)
	ListPermissions(ctx context.Context, page, limit int32) ([]dto.PermissionResponse, error)
	UpdatePermission(ctx context.Context, id string, req dto.UpdatePermissionRequest) error
	DeletePermission(ctx context.Context, id string) error
}

type permissionService struct {
	q store.Querier
}

func NewPermissionService(q store.Querier) PermissionService {
	return &permissionService{q: q}
}

func (s *permissionService) CreatePermission(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	perm, err := s.q.CreatePermission(ctx, store.CreatePermissionParams{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &dto.PermissionResponse{
		ID:          perm.PermissionID.String(),
		Name:        perm.Name,
		Description: perm.Description,
	}, nil
}

func (s *permissionService) GetPermission(ctx context.Context, id string) (*dto.PermissionResponse, error) {
	permID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid permission id: %w", err)
	}

	perm, err := s.q.GetPermission(ctx, permID)
	if err != nil {
		return nil, err
	}

	return &dto.PermissionResponse{
		ID:          perm.PermissionID.String(),
		Name:        perm.Name,
		Description: perm.Description,
	}, nil
}

func (s *permissionService) ListPermissions(ctx context.Context, page, limit int32) ([]dto.PermissionResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	perms, err := s.q.ListPermissions(ctx, store.ListPermissionsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var result []dto.PermissionResponse
	for _, p := range perms {
		result = append(result, dto.PermissionResponse{
			ID:          p.PermissionID.String(),
			Name:        p.Name,
			Description: p.Description,
		})
	}

	return result, nil
}

func (s *permissionService) UpdatePermission(ctx context.Context, id string, req dto.UpdatePermissionRequest) error {
	permID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid permission id: %w", err)
	}

	err = s.q.UpdatePermission(ctx, store.UpdatePermissionParams{
		PermissionID: permID,
		Name:         req.Name,
		Description:  req.Description,
	})
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}
	return nil
}

func (s *permissionService) DeletePermission(ctx context.Context, id string) error {
	permID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid permission id: %w", err)
	}

	err = s.q.DeletePermission(ctx, permID)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}
