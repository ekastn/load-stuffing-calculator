package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

func TestUserService_CreateUser(t *testing.T) {
	validRole := "planner"
	validRoleID := uuid.New()
	validPassword := "password123"

	tests := []struct {
		name           string
		req            dto.CreateUserRequest
		getRoleErr     error
		getRoleResp    store.GetRoleByNameRow
		createUserErr  error
		createUserResp store.User
		wantErr        bool
	}{
		{
			name: "success",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Email:    "new@example.com",
				Password: validPassword,
				Role:     validRole,
			},
			getRoleResp: store.GetRoleByNameRow{
				RoleID: validRoleID,
				Name:   validRole,
			},
			createUserResp: store.User{
				UserID:    uuid.New(),
				Username:  "newuser",
				Email:     "new@example.com",
				RoleID:    validRoleID,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "role_not_found",
			req: dto.CreateUserRequest{
				Role: "invalid_role",
			},
			getRoleErr: fmt.Errorf("role not found"),
			wantErr:    true,
		},
		{
			name: "create_user_db_error",
			req: dto.CreateUserRequest{
				Username: "newuser",
				Password: validPassword,
				Role:     validRole,
			},
			getRoleResp:   store.GetRoleByNameRow{RoleID: validRoleID},
			createUserErr: fmt.Errorf("db error"),
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				GetRoleByNameFunc: func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					if name != tt.req.Role {
						// For tests where we expect success/db error, role matches.
						// For role_not_found, we return error directly or check name.
						// Simplification: just return what test configures
					}
					return tt.getRoleResp, tt.getRoleErr
				},
				CreateUserFunc: func(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
					// Verify password hashing happened (simple check length)
					if len(arg.PasswordHash) == 0 {
						t.Error("Password was not hashed")
					}
					// Verify role ID
					if arg.RoleID != tt.getRoleResp.RoleID {
						t.Errorf("RoleID mismatch: got %v, want %v", arg.RoleID, tt.getRoleResp.RoleID)
					}
					return tt.createUserResp, tt.createUserErr
				},
			}

			s := service.NewUserService(mockQ)
			resp, err := s.CreateUser(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.ID != tt.createUserResp.UserID.String() {
					t.Errorf("Resp ID = %v, want %v", resp.ID, tt.createUserResp.UserID.String())
				}
				if resp.Role != tt.req.Role {
					t.Errorf("Resp Role = %v, want %v", resp.Role, tt.req.Role)
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	validID := uuid.New()
	validIDStr := validID.String()

	fullName := "John Doe"
	gender := "male"

	tests := []struct {
		name        string
		id          string
		getUserErr  error
		getUserResp store.GetUserByIDRow
		wantErr     bool
	}{
		{
			name: "success_with_profile",
			id:   validIDStr,
			getUserResp: store.GetUserByIDRow{
				UserID:   validID,
				Username: "john",
				Email:    "john@example.com",
				RoleName: "user",
				FullName: &fullName,
				Gender:   &gender,
			},
			wantErr: false,
		},
		{
			name: "success_no_profile",
			id:   validIDStr,
			getUserResp: store.GetUserByIDRow{
				UserID:   validID,
				Username: "john",
			},
			wantErr: false,
		},
		{
			name:       "not_found",
			id:         validIDStr,
			getUserErr: fmt.Errorf("not found"),
			wantErr:    true,
		},
		{
			name:    "invalid_uuid",
			id:      "invalid-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					if id.String() != tt.id {
						return store.GetUserByIDRow{}, fmt.Errorf("unexpected id")
					}
					return tt.getUserResp, tt.getUserErr
				},
			}

			s := service.NewUserService(mockQ)
			resp, err := s.GetUserByID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.ID != tt.getUserResp.UserID.String() {
					t.Errorf("Resp ID = %v, want %v", resp.ID, tt.getUserResp.UserID.String())
				}
				if tt.getUserResp.FullName != nil {
					if resp.Profile == nil {
						t.Error("Expected profile, got nil")
					} else if resp.Profile.FullName == nil || *resp.Profile.FullName != *tt.getUserResp.FullName {
						t.Errorf("Profile Name mismatch")
					}
				}
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	tests := []struct {
		name          string
		page, limit   int32
		listUsersResp []store.ListUsersRow
		listUsersErr  error
		wantErr       bool
		wantLen       int
	}{
		{
			name: "success_list",
			page: 1, limit: 10,
			listUsersResp: []store.ListUsersRow{
				{UserID: uuid.New(), Username: "u1"},
				{UserID: uuid.New(), Username: "u2"},
			},
			wantLen: 2,
		},
		{
			name: "db_error",
			page: 1, limit: 10,
			listUsersErr: fmt.Errorf("db error"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				ListUsersFunc: func(ctx context.Context, arg store.ListUsersParams) ([]store.ListUsersRow, error) {
					return tt.listUsersResp, tt.listUsersErr
				},
			}

			s := service.NewUserService(mockQ)
			resp, err := s.ListUsers(context.Background(), tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
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
