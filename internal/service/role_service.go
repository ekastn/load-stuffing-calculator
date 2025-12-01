package service

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

type RoleService interface {
	CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error)
	GetRole(ctx context.Context, id string) (*dto.RoleResponse, error)
	ListRoles(ctx context.Context, page, limit int32) ([]dto.RoleResponse, error)
	UpdateRole(ctx context.Context, id string, req dto.UpdateRoleRequest) error
	DeleteRole(ctx context.Context, id string) error
}

type roleService struct {
	q store.Querier
}

func NewRoleService(q store.Querier) RoleService {
	return &roleService{q: q}
}

func (s *roleService) CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.q.CreateRole(ctx, store.CreateRoleParams{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &dto.RoleResponse{
		ID:          role.RoleID.String(),
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (s *roleService) GetRole(ctx context.Context, id string) (*dto.RoleResponse, error) {
	roleID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid role id: %w", err)
	}

	role, err := s.q.GetRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return &dto.RoleResponse{
		ID:          role.RoleID.String(),
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (s *roleService) ListRoles(ctx context.Context, page, limit int32) ([]dto.RoleResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	roles, err := s.q.ListRoles(ctx, store.ListRolesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var result []dto.RoleResponse
	for _, r := range roles {
		result = append(result, dto.RoleResponse{
			ID:          r.RoleID.String(),
			Name:        r.Name,
			Description: r.Description,
		})
	}

	return result, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id string, req dto.UpdateRoleRequest) error {
	roleID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid role id: %w", err)
	}

	err = s.q.UpdateRole(ctx, store.UpdateRoleParams{
		RoleID:      roleID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (s *roleService) DeleteRole(ctx context.Context, id string) error {
	roleID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid role id: %w", err)
	}

	err = s.q.DeleteRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}
