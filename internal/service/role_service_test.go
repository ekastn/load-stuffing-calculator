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
		{
			name:    "invalid_id",
			id:      "invalid-uuid",
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

func TestRoleService_UpdateRolePermissions(t *testing.T) {
	roleID := uuid.New()
	perm1 := uuid.New()
	perm2 := uuid.New()

	tests := []struct {
		name             string
		roleID           string
		permissionIDs    []string
		deletePermsErr   error
		addPermErr       error
		addPermCallCount int
		wantErr          bool
		errContains      string
	}{
		{
			name:             "success_with_two_permissions",
			roleID:           roleID.String(),
			permissionIDs:    []string{perm1.String(), perm2.String()},
			addPermCallCount: 2,
			wantErr:          false,
		},
		{
			name:             "success_with_no_permissions",
			roleID:           roleID.String(),
			permissionIDs:    []string{},
			addPermCallCount: 0,
			wantErr:          false,
		},
		{
			name:          "invalid_role_id",
			roleID:        "invalid-uuid",
			permissionIDs: []string{perm1.String()},
			wantErr:       true,
			errContains:   "invalid role id",
		},
		{
			name:          "invalid_permission_id",
			roleID:        roleID.String(),
			permissionIDs: []string{"invalid-uuid"},
			wantErr:       true,
			errContains:   "invalid permission id",
		},
		{
			name:           "delete_permissions_error",
			roleID:         roleID.String(),
			permissionIDs:  []string{perm1.String()},
			deletePermsErr: fmt.Errorf("delete error"),
			wantErr:        true,
			errContains:    "failed to clear existing permissions",
		},
		{
			name:          "add_permission_error",
			roleID:        roleID.String(),
			permissionIDs: []string{perm1.String()},
			addPermErr:    fmt.Errorf("add error"),
			wantErr:       true,
			errContains:   "failed to add permission",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addPermCallCount := 0

			mockQ := &MockQuerier{
				DeleteRolePermissionsFunc: func(ctx context.Context, roleID uuid.UUID) error {
					return tt.deletePermsErr
				},
				AddRolePermissionFunc: func(ctx context.Context, arg store.AddRolePermissionParams) error {
					addPermCallCount++
					return tt.addPermErr
				},
			}

			s := service.NewRoleService(mockQ)
			err := s.UpdateRolePermissions(context.Background(), tt.roleID, tt.permissionIDs)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRolePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateRolePermissions() error = %v, want error containing %q", err, tt.errContains)
				}
			}

			if !tt.wantErr && addPermCallCount != tt.addPermCallCount {
				t.Errorf("AddRolePermission called %d times, want %d", addPermCallCount, tt.addPermCallCount)
			}
		})
	}
}

func TestRoleService_GetRolePermissions(t *testing.T) {
	roleID := uuid.New()
	perm1 := uuid.New()
	perm2 := uuid.New()

	tests := []struct {
		name        string
		roleID      string
		mockPerms   []uuid.UUID
		mockErr     error
		wantErr     bool
		wantCount   int
		errContains string
	}{
		{
			name:      "success_with_permissions",
			roleID:    roleID.String(),
			mockPerms: []uuid.UUID{perm1, perm2},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:      "success_with_no_permissions",
			roleID:    roleID.String(),
			mockPerms: []uuid.UUID{},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:        "invalid_role_id",
			roleID:      "invalid-uuid",
			wantErr:     true,
			errContains: "invalid role id",
		},
		{
			name:    "db_error",
			roleID:  roleID.String(),
			mockErr: fmt.Errorf("database error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				GetRolePermissionsFunc: func(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
					return tt.mockPerms, tt.mockErr
				},
			}

			s := service.NewRoleService(mockQ)
			perms, err := s.GetRolePermissions(context.Background(), tt.roleID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRolePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("GetRolePermissions() error = %v, want error containing %q", err, tt.errContains)
				}
			}

			if !tt.wantErr {
				if len(perms) != tt.wantCount {
					t.Errorf("GetRolePermissions() returned %d permissions, want %d", len(perms), tt.wantCount)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
