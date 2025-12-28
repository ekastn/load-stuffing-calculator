package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

func TestContainerService_CreateContainer(t *testing.T) {
	name := "20ft"
	length := 6058.0
	desc := "Standard"

	tests := []struct {
		name       string
		req        dto.CreateContainerRequest
		createErr  error
		createResp store.Container
		wantErr    bool
	}{
		{
			name: "success",
			req: dto.CreateContainerRequest{
				Name:          name,
				InnerLengthMM: length,
				InnerWidthMM:  2438.0,
				InnerHeightMM: 2591.0,
				MaxWeightKG:   28000.0,
				Description:   &desc,
			},
			createResp: store.Container{
				ContainerID:   uuid.New(),
				Name:          name,
				InnerLengthMm: toNumeric(length),
				InnerWidthMm:  toNumeric(2438.0),
				InnerHeightMm: toNumeric(2591.0),
				MaxWeightKg:   toNumeric(28000.0),
				Description:   &desc,
			},
			wantErr: false,
		},
		{
			name: "db_error",
			req: dto.CreateContainerRequest{
				Name: name,
			},
			createErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
		{
			name:    "trial_no_workspace_forbidden",
			req:     dto.CreateContainerRequest{Name: name},
			wantErr: true,
		},
		{
			name: "founder_no_override_creates_global_preset",
			req:  dto.CreateContainerRequest{Name: name},
			createResp: store.Container{
				ContainerID: uuid.New(),
				Name:        name,
			},
			wantErr: false,
		},
		{
			name: "founder_with_override_creates_scoped",
			req:  dto.CreateContainerRequest{Name: name},
			createResp: store.Container{
				ContainerID: uuid.New(),
				Name:        name,
			},
			wantErr: false,
		},
		{
			name:    "founder_invalid_override_errors",
			req:     dto.CreateContainerRequest{Name: name},
			wantErr: true,
		},
		{
			name: "founder_override_requires_no_token_workspace",
			req:  dto.CreateContainerRequest{Name: name},
			createResp: store.Container{
				ContainerID: uuid.New(),
				Name:        name,
			},
			wantErr: false,
		},
		{
			name:    "non_founder_missing_workspace_errors",
			req:     dto.CreateContainerRequest{Name: name},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sawWorkspaceID := (*uuid.UUID)(nil)
			mockQ := &MockQuerier{
				CreateContainerFunc: func(ctx context.Context, arg store.CreateContainerParams) (store.Container, error) {
					sawWorkspaceID = arg.WorkspaceID
					switch tt.name {
					case "trial_no_workspace_forbidden", "non_founder_missing_workspace_errors", "founder_invalid_override_errors":
						return store.Container{}, fmt.Errorf("unexpected db call")
					}
					return tt.createResp, tt.createErr
				},
			}

			s := service.NewContainerService(mockQ)

			workspaceID := uuid.New()
			overrideWorkspaceID := uuid.New()

			ctx := ctxWithWorkspaceID(workspaceID)
			switch tt.name {
			case "trial_no_workspace_forbidden":
				ctx = context.Background()
			case "founder_no_override_creates_global_preset":
				ctx = ctxWithRole("founder")
			case "founder_with_override_creates_scoped":
				ctx = auth.WithWorkspaceOverrideID(ctxWithRoleAndWorkspace("founder", workspaceID), overrideWorkspaceID.String())
			case "founder_invalid_override_errors":
				ctx = auth.WithWorkspaceOverrideID(ctxWithRoleAndWorkspace("founder", workspaceID), "not-a-uuid")
			case "founder_override_requires_no_token_workspace":
				ctx = auth.WithWorkspaceOverrideID(ctxWithRole("founder"), overrideWorkspaceID.String())
			case "non_founder_missing_workspace_errors":
				ctx = ctxWithRole("admin")
			}

			resp, err := s.CreateContainer(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Ensure we don't call the DB for validation-only failures.
			switch tt.name {
			case "trial_no_workspace_forbidden", "non_founder_missing_workspace_errors", "founder_invalid_override_errors":
				if sawWorkspaceID != nil {
					t.Fatalf("unexpected db call")
				}
				return
			}

			// For success and db-level errors, validate what workspace was used.
			switch tt.name {
			case "founder_no_override_creates_global_preset":
				if sawWorkspaceID != nil {
					t.Fatalf("expected WorkspaceID to be nil, got %v", *sawWorkspaceID)
				}
			case "founder_with_override_creates_scoped", "founder_override_requires_no_token_workspace":
				if sawWorkspaceID == nil || *sawWorkspaceID != overrideWorkspaceID {
					t.Fatalf("expected WorkspaceID to be override %v, got %v", overrideWorkspaceID, sawWorkspaceID)
				}
			case "success", "db_error":
				if sawWorkspaceID == nil || *sawWorkspaceID != workspaceID {
					t.Fatalf("expected WorkspaceID to be token workspace %v, got %v", workspaceID, sawWorkspaceID)
				}
			}

			if !tt.wantErr {
				if resp.Name != tt.req.Name {
					t.Errorf("Name = %v, want %v", resp.Name, tt.req.Name)
				}
				if resp.InnerLengthMM != tt.req.InnerLengthMM {
					t.Errorf("Length mismatch")
				}
			}
		})
	}
}

func TestContainerService_GetContainer(t *testing.T) {
	id := uuid.New()
	name := "40ft"

	tests := []struct {
		name    string
		id      string
		getErr  error
		getResp store.Container
		wantErr bool
	}{
		{
			name: "success",
			id:   id.String(),
			getResp: store.Container{
				ContainerID: id,
				Name:        name,
			},
			wantErr: false,
		},
		{
			name:    "not_found",
			id:      id.String(),
			getErr:  fmt.Errorf("not found"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID := uuid.New()
			ctx := ctxWithWorkspaceID(workspaceID)

			mockQ := &MockQuerier{
				GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
					if arg.ContainerID.String() != tt.id {
						return store.Container{}, fmt.Errorf("id mismatch")
					}
					if arg.WorkspaceID == nil || *arg.WorkspaceID != workspaceID {
						return store.Container{}, fmt.Errorf("workspace mismatch")
					}
					return tt.getResp, tt.getErr
				},
			}

			s := service.NewContainerService(mockQ)
			resp, err := s.GetContainer(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.ID != tt.id {
					t.Errorf("ID = %v, want %v", resp.ID, tt.id)
				}
			}
		})
	}

	t.Run("founder_no_override_uses_any", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		called := false
		mockQ := &MockQuerier{
			GetContainerAnyFunc: func(ctx context.Context, containerID uuid.UUID) (store.Container, error) {
				called = true
				if containerID != id {
					return store.Container{}, fmt.Errorf("id mismatch")
				}
				return store.Container{ContainerID: id, Name: name}, nil
			},
			GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
				return store.Container{}, fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewContainerService(mockQ)
		resp, err := s.GetContainer(ctx, id.String())
		if err != nil {
			t.Fatalf("GetContainer() error = %v", err)
		}
		if !called {
			t.Fatalf("expected GetContainerAny to be called")
		}
		if resp.ID != id.String() {
			t.Fatalf("ID = %v, want %v", resp.ID, id.String())
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())

		mockQ := &MockQuerier{
			GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
				if arg.ContainerID != id {
					return store.Container{}, fmt.Errorf("id mismatch")
				}
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return store.Container{}, fmt.Errorf("workspace mismatch")
				}
				return store.Container{ContainerID: id, Name: name}, nil
			},
			GetContainerAnyFunc: func(ctx context.Context, containerID uuid.UUID) (store.Container, error) {
				return store.Container{}, fmt.Errorf("unexpected any call")
			},
		}

		s := service.NewContainerService(mockQ)
		resp, err := s.GetContainer(ctx, id.String())
		if err != nil {
			t.Fatalf("GetContainer() error = %v", err)
		}
		if resp.ID != id.String() {
			t.Fatalf("ID = %v, want %v", resp.ID, id.String())
		}
	})
}

func TestContainerService_ListContainers(t *testing.T) {
	tests := []struct {
		name        string
		page, limit int32
		listResp    []store.Container
		listErr     error
		wantErr     bool
		wantLen     int
	}{
		{
			name: "success",
			page: 1, limit: 10,
			listResp: []store.Container{
				{ContainerID: uuid.New(), Name: "c1"},
				{ContainerID: uuid.New(), Name: "c2"},
			},
			wantLen: 2,
		},
		{
			name: "trial_no_workspace_returns_global_presets",
			page: 1, limit: 10,
			listResp: []store.Container{
				{ContainerID: uuid.New(), Name: "global1"},
			},
			wantLen: 1,
		},
		{
			name:    "db_error",
			listErr: fmt.Errorf("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				ListContainersFunc: func(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error) {
					if tt.name == "trial_no_workspace_returns_global_presets" && arg.WorkspaceID != nil {
						return nil, fmt.Errorf("expected nil workspace id")
					}
					return tt.listResp, tt.listErr
				},
			}

			s := service.NewContainerService(mockQ)
			ctx := ctxWithWorkspaceID(uuid.New())
			if tt.name == "trial_no_workspace_returns_global_presets" {
				ctx = context.Background()
			}
			resp, err := s.ListContainers(ctx, tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListContainers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(resp) != tt.wantLen {
					t.Errorf("Resp len = %v, want %v", len(resp), tt.wantLen)
				}
			}
		})
	}

	t.Run("founder_no_override_uses_all", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		page, limit := int32(2), int32(15)
		offset := (page - 1) * limit
		called := false

		mockQ := &MockQuerier{
			ListContainersAllFunc: func(ctx context.Context, arg store.ListContainersAllParams) ([]store.Container, error) {
				called = true
				if arg.Limit != limit || arg.Offset != offset {
					return nil, fmt.Errorf("limit/offset mismatch")
				}
				return []store.Container{{ContainerID: uuid.New(), Name: "global"}}, nil
			},
			ListContainersFunc: func(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error) {
				return nil, fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewContainerService(mockQ)
		resp, err := s.ListContainers(ctx, page, limit)
		if err != nil {
			t.Fatalf("ListContainers() error = %v", err)
		}
		if !called {
			t.Fatalf("expected ListContainersAll to be called")
		}
		if len(resp) != 1 {
			t.Fatalf("resp len = %v, want %v", len(resp), 1)
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())
		page, limit := int32(1), int32(10)
		offset := (page - 1) * limit

		mockQ := &MockQuerier{
			ListContainersFunc: func(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error) {
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return nil, fmt.Errorf("workspace mismatch")
				}
				if arg.Limit != limit || arg.Offset != offset {
					return nil, fmt.Errorf("limit/offset mismatch")
				}
				return []store.Container{{ContainerID: uuid.New(), Name: "scoped"}}, nil
			},
			ListContainersAllFunc: func(ctx context.Context, arg store.ListContainersAllParams) ([]store.Container, error) {
				return nil, fmt.Errorf("unexpected all call")
			},
		}

		s := service.NewContainerService(mockQ)
		resp, err := s.ListContainers(ctx, page, limit)
		if err != nil {
			t.Fatalf("ListContainers() error = %v", err)
		}
		if len(resp) != 1 {
			t.Fatalf("resp len = %v, want %v", len(resp), 1)
		}
	})
}

func TestContainerService_UpdateContainer(t *testing.T) {
	id := uuid.New()
	name := "updated_c"

	tests := []struct {
		name      string
		id        string
		req       dto.UpdateContainerRequest
		updateErr error
		wantErr   bool
	}{
		{
			name:    "success",
			id:      id.String(),
			req:     dto.UpdateContainerRequest{Name: name},
			wantErr: false,
		},
		{
			name:      "db_error",
			id:        id.String(),
			req:       dto.UpdateContainerRequest{Name: name},
			updateErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
		{
			name:    "trial_no_workspace_forbidden",
			id:      id.String(),
			req:     dto.UpdateContainerRequest{Name: name},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				UpdateContainerFunc: func(ctx context.Context, arg store.UpdateContainerParams) error {
					if tt.name == "trial_no_workspace_forbidden" {
						return fmt.Errorf("unexpected db call")
					}
					return tt.updateErr
				},
			}

			s := service.NewContainerService(mockQ)
			ctx := ctxWithWorkspaceID(uuid.New())
			if tt.name == "trial_no_workspace_forbidden" {
				ctx = context.Background()
			}
			err := s.UpdateContainer(ctx, tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("founder_no_override_uses_any", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		req := dto.UpdateContainerRequest{Name: "founder_update"}
		called := false

		mockQ := &MockQuerier{
			UpdateContainerAnyFunc: func(ctx context.Context, arg store.UpdateContainerAnyParams) error {
				called = true
				if arg.ContainerID != id {
					return fmt.Errorf("id mismatch")
				}
				if arg.Name != req.Name {
					return fmt.Errorf("name mismatch")
				}
				return nil
			},
			UpdateContainerFunc: func(ctx context.Context, arg store.UpdateContainerParams) error {
				return fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewContainerService(mockQ)
		err := s.UpdateContainer(ctx, id.String(), req)
		if err != nil {
			t.Fatalf("UpdateContainer() error = %v", err)
		}
		if !called {
			t.Fatalf("expected UpdateContainerAny to be called")
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())
		req := dto.UpdateContainerRequest{Name: "founder_update_scoped"}

		mockQ := &MockQuerier{
			UpdateContainerFunc: func(ctx context.Context, arg store.UpdateContainerParams) error {
				if arg.ContainerID != id {
					return fmt.Errorf("id mismatch")
				}
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return fmt.Errorf("workspace mismatch")
				}
				if arg.Name != req.Name {
					return fmt.Errorf("name mismatch")
				}
				return nil
			},
			UpdateContainerAnyFunc: func(ctx context.Context, arg store.UpdateContainerAnyParams) error {
				return fmt.Errorf("unexpected any call")
			},
		}

		s := service.NewContainerService(mockQ)
		err := s.UpdateContainer(ctx, id.String(), req)
		if err != nil {
			t.Fatalf("UpdateContainer() error = %v", err)
		}
	})
}

func TestContainerService_DeleteContainer(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name      string
		id        string
		deleteErr error
		wantErr   bool
	}{
		{
			name:    "success",
			id:      id.String(),
			wantErr: false,
		},
		{
			name:    "trial_no_workspace_forbidden",
			id:      id.String(),
			wantErr: true,
		},
		{
			name:      "db_error",
			id:        id.String(),
			deleteErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				DeleteContainerFunc: func(ctx context.Context, arg store.DeleteContainerParams) error {
					if tt.name == "trial_no_workspace_forbidden" {
						return fmt.Errorf("unexpected db call")
					}
					return tt.deleteErr
				},
			}

			s := service.NewContainerService(mockQ)
			ctx := ctxWithWorkspaceID(uuid.New())
			if tt.name == "trial_no_workspace_forbidden" {
				ctx = context.Background()
			}
			err := s.DeleteContainer(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("founder_no_override_uses_any", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		called := false

		mockQ := &MockQuerier{
			DeleteContainerAnyFunc: func(ctx context.Context, containerID uuid.UUID) error {
				called = true
				if containerID != id {
					return fmt.Errorf("id mismatch")
				}
				return nil
			},
			DeleteContainerFunc: func(ctx context.Context, arg store.DeleteContainerParams) error {
				return fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewContainerService(mockQ)
		err := s.DeleteContainer(ctx, id.String())
		if err != nil {
			t.Fatalf("DeleteContainer() error = %v", err)
		}
		if !called {
			t.Fatalf("expected DeleteContainerAny to be called")
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())

		mockQ := &MockQuerier{
			DeleteContainerFunc: func(ctx context.Context, arg store.DeleteContainerParams) error {
				if arg.ContainerID != id {
					return fmt.Errorf("id mismatch")
				}
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return fmt.Errorf("workspace mismatch")
				}
				return nil
			},
			DeleteContainerAnyFunc: func(ctx context.Context, containerID uuid.UUID) error {
				return fmt.Errorf("unexpected any call")
			},
		}

		s := service.NewContainerService(mockQ)
		err := s.DeleteContainer(ctx, id.String())
		if err != nil {
			t.Fatalf("DeleteContainer() error = %v", err)
		}
	})
}
