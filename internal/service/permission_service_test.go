package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

func TestPermissionService_CreatePermission(t *testing.T) {
	name := "new_perm"
	desc := "A new perm"

	tests := []struct {
		name       string
		req        dto.CreatePermissionRequest
		createErr  error
		createResp store.Permission
		wantErr    bool
	}{
		{
			name: "success",
			req: dto.CreatePermissionRequest{
				Name:        name,
				Description: &desc,
			},
			createResp: store.Permission{
				PermissionID: uuid.New(),
				Name:         name,
				Description:  &desc,
			},
			wantErr: false,
		},
		{
			name: "db_error",
			req: dto.CreatePermissionRequest{
				Name: name,
			},
			createErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				CreatePermissionFunc: func(ctx context.Context, arg store.CreatePermissionParams) (store.Permission, error) {
					return tt.createResp, tt.createErr
				},
			}

			s := service.NewPermissionService(mockQ)
			resp, err := s.CreatePermission(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.Name != tt.req.Name {
					t.Errorf("Name = %v, want %v", resp.Name, tt.req.Name)
				}
			}
		})
	}
}

func TestPermissionService_GetPermission(t *testing.T) {
	id := uuid.New()
	name := "perm"

	tests := []struct {
		name    string
		id      string
		getErr  error
		getResp store.Permission
		wantErr bool
	}{
		{
			name: "success",
			id:   id.String(),
			getResp: store.Permission{
				PermissionID: id,
				Name:         name,
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
			mockQ := &MockQuerier{
				GetPermissionFunc: func(ctx context.Context, permID uuid.UUID) (store.Permission, error) {
					if permID.String() != tt.id {
						return store.Permission{}, fmt.Errorf("id mismatch")
					}
					return tt.getResp, tt.getErr
				},
			}

			s := service.NewPermissionService(mockQ)
			resp, err := s.GetPermission(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.ID != tt.id {
					t.Errorf("ID = %v, want %v", resp.ID, tt.id)
				}
			}
		})
	}
}

func TestPermissionService_ListPermissions(t *testing.T) {
	tests := []struct {
		name          string
		page, limit   int32
		listPermsResp []store.Permission
		listPermsErr  error
		wantErr       bool
		wantLen       int
	}{
		{
			name: "success",
			page: 1, limit: 10,
			listPermsResp: []store.Permission{
				{PermissionID: uuid.New(), Name: "perm1"},
				{PermissionID: uuid.New(), Name: "perm2"},
			},
			wantLen: 2,
		},
		{
			name:         "db_error",
			listPermsErr: fmt.Errorf("db error"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				ListPermissionsFunc: func(ctx context.Context, arg store.ListPermissionsParams) ([]store.Permission, error) {
					return tt.listPermsResp, tt.listPermsErr
				},
			}

			s := service.NewPermissionService(mockQ)
			resp, err := s.ListPermissions(context.Background(), tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(resp) != tt.wantLen {
					t.Errorf("Resp len = %v, want %v", len(resp), tt.wantLen)
				}
			}
		})
	}
}

func TestPermissionService_UpdatePermission(t *testing.T) {
	id := uuid.New()
	name := "updated_perm"

	tests := []struct {
		name      string
		id        string
		req       dto.UpdatePermissionRequest
		updateErr error
		wantErr   bool
	}{
		{
			name:    "success",
			id:      id.String(),
			req:     dto.UpdatePermissionRequest{Name: name},
			wantErr: false,
		},
		{
			name:      "db_error",
			id:        id.String(),
			req:       dto.UpdatePermissionRequest{Name: name},
			updateErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				UpdatePermissionFunc: func(ctx context.Context, arg store.UpdatePermissionParams) error {
					return tt.updateErr
				},
			}

			s := service.NewPermissionService(mockQ)
			err := s.UpdatePermission(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionService_DeletePermission(t *testing.T) {
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
			name:      "db_error",
			id:        id.String(),
			deleteErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				DeletePermissionFunc: func(ctx context.Context, permID uuid.UUID) error {
					return tt.deleteErr
				},
			}

			s := service.NewPermissionService(mockQ)
			err := s.DeletePermission(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
