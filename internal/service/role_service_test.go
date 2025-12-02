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

func TestRoleService_CreateRole(t *testing.T) {
	name := "new_role"
	desc := "A new role"

	tests := []struct {
		name       string
		req        dto.CreateRoleRequest
		createErr  error
		createResp store.Role
		wantErr    bool
	}{
		{
			name: "success",
			req: dto.CreateRoleRequest{
				Name:        name,
				Description: &desc,
			},
			createResp: store.Role{
				RoleID:      uuid.New(),
				Name:        name,
				Description: &desc,
			},
			wantErr: false,
		},
		{
			name: "db_error",
			req: dto.CreateRoleRequest{
				Name: name,
			},
			createErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				CreateRoleFunc: func(ctx context.Context, arg store.CreateRoleParams) (store.Role, error) {
					return tt.createResp, tt.createErr
				},
			}

			s := service.NewRoleService(mockQ)
			resp, err := s.CreateRole(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.Name != tt.req.Name {
					t.Errorf("Name = %v, want %v", resp.Name, tt.req.Name)
				}
				if resp.Description != tt.createResp.Description {
					t.Errorf("Description mismatch")
				}
			}
		})
	}
}

func TestRoleService_GetRole(t *testing.T) {
	id := uuid.New()
	name := "role"

	tests := []struct {
		name    string
		id      string
		getErr  error
		getResp store.Role
		wantErr bool
	}{
		{
			name: "success",
			id:   id.String(),
			getResp: store.Role{
				RoleID: id,
				Name:   name,
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
				GetRoleFunc: func(ctx context.Context, roleID uuid.UUID) (store.Role, error) {
					if roleID.String() != tt.id {
						return store.Role{}, fmt.Errorf("id mismatch")
					}
					return tt.getResp, tt.getErr
				},
			}

			s := service.NewRoleService(mockQ)
			resp, err := s.GetRole(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRoleService_ListRoles(t *testing.T) {
	tests := []struct {
		name          string
		page, limit   int32
		listRolesResp []store.Role
		listRolesErr  error
		wantErr       bool
		wantLen       int
	}{
		{
			name: "success",
			page: 1, limit: 10,
			listRolesResp: []store.Role{
				{RoleID: uuid.New(), Name: "role1"},
				{RoleID: uuid.New(), Name: "role2"},
			},
			wantLen: 2,
		},
		{
			name:         "db_error",
			listRolesErr: fmt.Errorf("db error"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				ListRolesFunc: func(ctx context.Context, arg store.ListRolesParams) ([]store.Role, error) {
					return tt.listRolesResp, tt.listRolesErr
				},
			}

			s := service.NewRoleService(mockQ)
			resp, err := s.ListRoles(context.Background(), tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListRoles() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRoleService_UpdateRole(t *testing.T) {
	id := uuid.New()
	name := "updated_role"

	tests := []struct {
		name      string
		id        string
		req       dto.UpdateRoleRequest
		updateErr error
		wantErr   bool
	}{
		{
			name:    "success",
			id:      id.String(),
			req:     dto.UpdateRoleRequest{Name: name},
			wantErr: false,
		},
		{
			name:      "db_error",
			id:        id.String(),
			req:       dto.UpdateRoleRequest{Name: name},
			updateErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
		{
			name:    "invalid_id",
			id:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				UpdateRoleFunc: func(ctx context.Context, arg store.UpdateRoleParams) error {
					return tt.updateErr
				},
			}

			s := service.NewRoleService(mockQ)
			err := s.UpdateRole(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoleService_DeleteRole(t *testing.T) {
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
		{
			name:    "invalid_id",
			id:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				DeleteRoleFunc: func(ctx context.Context, roleID uuid.UUID) error {
					return tt.deleteErr
				},
			}

			s := service.NewRoleService(mockQ)
			err := s.DeleteRole(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
