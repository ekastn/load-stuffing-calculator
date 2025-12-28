package service

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

type ContainerService interface {
	CreateContainer(ctx context.Context, req dto.CreateContainerRequest) (*dto.ContainerResponse, error)
	GetContainer(ctx context.Context, id string) (*dto.ContainerResponse, error)
	ListContainers(ctx context.Context, page, limit int32) ([]dto.ContainerResponse, error)
	UpdateContainer(ctx context.Context, id string, req dto.UpdateContainerRequest) error
	DeleteContainer(ctx context.Context, id string) error
}

type containerService struct {
	q store.Querier
}

func NewContainerService(q store.Querier) ContainerService {
	return &containerService{q: q}
}

func (s *containerService) CreateContainer(ctx context.Context, req dto.CreateContainerRequest) (*dto.ContainerResponse, error) {
	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Founders can create either:
	// - Global presets (workspace_id = NULL) when no override is provided.
	// - Workspace-scoped records when ?workspace_id= is provided.
	if isFounder(ctx) {
		workspaceID = overrideWorkspaceID
	}

	if workspaceID == nil && !isFounder(ctx) {
		return nil, fmt.Errorf("workspace id is required")
	}

	container, err := s.q.CreateContainer(ctx, store.CreateContainerParams{
		WorkspaceID:   workspaceID,
		Name:          req.Name,
		InnerLengthMm: toNumeric(req.InnerLengthMM),
		InnerWidthMm:  toNumeric(req.InnerWidthMM),
		InnerHeightMm: toNumeric(req.InnerHeightMM),
		MaxWeightKg:   toNumeric(req.MaxWeightKG),
		Description:   req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	return mapContainerToResponse(container), nil
}

func (s *containerService) GetContainer(ctx context.Context, id string) (*dto.ContainerResponse, error) {
	containerID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid container id: %w", err)
	}

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		container, err := s.q.GetContainerAny(ctx, containerID)
		if err != nil {
			return nil, err
		}
		return mapContainerToResponse(container), nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}

	container, err := s.q.GetContainer(ctx, store.GetContainerParams{ContainerID: containerID, WorkspaceID: workspaceID})
	if err != nil {
		return nil, err
	}

	return mapContainerToResponse(container), nil
}

func (s *containerService) ListContainers(ctx context.Context, page, limit int32) ([]dto.ContainerResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		containers, err := s.q.ListContainersAll(ctx, store.ListContainersAllParams{Limit: limit, Offset: offset})
		if err != nil {
			return nil, err
		}
		var result []dto.ContainerResponse
		for _, c := range containers {
			result = append(result, *mapContainerToResponse(c))
		}
		return result, nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}

	containers, err := s.q.ListContainers(ctx, store.ListContainersParams{WorkspaceID: workspaceID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}

	var result []dto.ContainerResponse
	for _, c := range containers {
		result = append(result, *mapContainerToResponse(c))
	}

	return result, nil
}

func (s *containerService) UpdateContainer(ctx context.Context, id string, req dto.UpdateContainerRequest) error {
	containerID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid container id: %w", err)
	}

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		err = s.q.UpdateContainerAny(ctx, store.UpdateContainerAnyParams{
			ContainerID:   containerID,
			Name:          req.Name,
			InnerLengthMm: toNumeric(req.InnerLengthMM),
			InnerWidthMm:  toNumeric(req.InnerWidthMM),
			InnerHeightMm: toNumeric(req.InnerHeightMM),
			MaxWeightKg:   toNumeric(req.MaxWeightKG),
			Description:   req.Description,
		})
		if err != nil {
			return fmt.Errorf("failed to update container: %w", err)
		}
		return nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}
	if workspaceID == nil {
		return fmt.Errorf("workspace id is required")
	}

	err = s.q.UpdateContainer(ctx, store.UpdateContainerParams{
		ContainerID:   containerID,
		WorkspaceID:   workspaceID,
		Name:          req.Name,
		InnerLengthMm: toNumeric(req.InnerLengthMM),
		InnerWidthMm:  toNumeric(req.InnerWidthMM),
		InnerHeightMm: toNumeric(req.InnerHeightMM),
		MaxWeightKg:   toNumeric(req.MaxWeightKG),
		Description:   req.Description,
	})
	if err != nil {
		return fmt.Errorf("failed to update container: %w", err)
	}
	return nil
}

func (s *containerService) DeleteContainer(ctx context.Context, id string) error {
	containerID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid container id: %w", err)
	}

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		err = s.q.DeleteContainerAny(ctx, containerID)
		if err != nil {
			return fmt.Errorf("failed to delete container: %w", err)
		}
		return nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}
	if workspaceID == nil {
		return fmt.Errorf("workspace id is required")
	}

	err = s.q.DeleteContainer(ctx, store.DeleteContainerParams{ContainerID: containerID, WorkspaceID: workspaceID})
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}
	return nil
}

func mapContainerToResponse(c store.Container) *dto.ContainerResponse {
	return &dto.ContainerResponse{
		ID:            c.ContainerID.String(),
		Name:          c.Name,
		InnerLengthMM: toFloat(c.InnerLengthMm),
		InnerWidthMM:  toFloat(c.InnerWidthMm),
		InnerHeightMM: toFloat(c.InnerHeightMm),
		MaxWeightKG:   toFloat(c.MaxWeightKg),
		Description:   c.Description,
	}
}
