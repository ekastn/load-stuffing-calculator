package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestMemberService_ListMembers(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	memberID := uuid.New()

	tests := []struct {
		name                string
		page                int32
		limit               int32
		overrideWorkspaceID *string
		ctx                 context.Context
		mockSetup           func(*MockQuerier)
		wantErr             bool
		wantLen             int
	}{
		{
			name:                "successful_list",
			page:                1,
			limit:               10,
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.ListMembersByWorkspaceFunc = func(ctx context.Context, arg store.ListMembersByWorkspaceParams) ([]store.ListMembersByWorkspaceRow, error) {
					return []store.ListMembersByWorkspaceRow{
						{MemberID: memberID, WorkspaceID: workspaceID, UserID: userID, RoleName: types.RoleUser.String(), Username: "test", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
					}, nil
				}
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:                "pagination_defaults",
			page:                0,
			limit:               0,
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.ListMembersByWorkspaceFunc = func(ctx context.Context, arg store.ListMembersByWorkspaceParams) ([]store.ListMembersByWorkspaceRow, error) {
					if arg.Limit != 10 {
						t.Errorf("expected default limit 10, got %d", arg.Limit)
					}
					if arg.Offset != 0 {
						t.Errorf("expected offset 0, got %d", arg.Offset)
					}
					return []store.ListMembersByWorkspaceRow{}, nil
				}
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:                "missing_workspace",
			page:                1,
			limit:               10,
			overrideWorkspaceID: nil,
			ctx:                 context.Background(),
			mockSetup:           func(mq *MockQuerier) {},
			wantErr:             true,
		},
		{
			name:                "forbidden_not_admin",
			page:                1,
			limit:               10,
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOperator.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name:                "permission_check_database_error",
			page:                1,
			limit:               10,
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return "", fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name:                "list_database_error",
			page:                1,
			limit:               10,
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.ListMembersByWorkspaceFunc = func(ctx context.Context, arg store.ListMembersByWorkspaceParams) ([]store.ListMembersByWorkspaceRow, error) {
					return nil, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			tt.mockSetup(mockQ)

			s := service.NewMemberService(mockQ)
			result, err := s.ListMembers(tt.ctx, tt.page, tt.limit, tt.overrideWorkspaceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("ListMembers() returned %d items, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestMemberService_AddMember(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	targetUserID := uuid.New()
	ownerID := uuid.New()
	roleID := uuid.New()
	memberID := uuid.New()

	tests := []struct {
		name                string
		req                 dto.AddMemberRequest
		overrideWorkspaceID *string
		ctx                 context.Context
		mockSetup           func(*MockQuerier)
		wantErr             bool
	}{
		{
			name: "successful_add_by_username",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByUsernameFunc = func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
					return store.GetUserByUsernameRow{UserID: targetUserID, Username: "testuser", Email: "test@example.com"}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: targetUserID, RoleID: roleID, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "successful_add_by_email",
			req: dto.AddMemberRequest{
				UserIdentifier: "test@example.com",
				Role:           types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByEmailFunc = func(ctx context.Context, email string) (store.GetUserByEmailRow, error) {
					return store.GetUserByEmailRow{UserID: targetUserID, Username: "testuser", Email: "test@example.com", PasswordHash: "hash"}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: targetUserID, RoleID: roleID, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "successful_add_by_uuid",
			req: dto.AddMemberRequest{
				UserIdentifier: targetUserID.String(),
				Role:           types.RolePlanner.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: targetUserID, Username: "testuser", Email: "test@example.com"}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: targetUserID, RoleID: roleID, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "missing_workspace",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 context.Background(),
			mockSetup:           func(mq *MockQuerier) {},
			wantErr:             true,
		},
		{
			name: "personal_workspace_rejection",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "personal"}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "forbidden_not_admin",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOperator.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name: "invalid_role_founder",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RoleFounder.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name: "user_not_found_by_username",
			req: dto.AddMemberRequest{
				UserIdentifier: "nonexistent",
				Role:           types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByUsernameFunc = func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
					return store.GetUserByUsernameRow{}, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name: "user_not_found_by_email",
			req: dto.AddMemberRequest{
				UserIdentifier: "nonexistent@example.com",
				Role:           types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByEmailFunc = func(ctx context.Context, email string) (store.GetUserByEmailRow, error) {
					return store.GetUserByEmailRow{}, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name: "user_not_found_by_uuid",
			req: dto.AddMemberRequest{
				UserIdentifier: uuid.New().String(),
				Role:           types.RolePlanner.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{}, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name: "already_a_member",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByUsernameFunc = func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
					return store.GetUserByUsernameRow{UserID: targetUserID, Username: "testuser", Email: "test@example.com"}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{}, &pgconn.PgError{Code: "23505"}
				}
			},
			wantErr: true,
		},
		{
			name: "create_member_database_error",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.GetUserByUsernameFunc = func(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
					return store.GetUserByUsernameRow{UserID: targetUserID, Username: "testuser", Email: "test@example.com"}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "get_role_database_error",
			req: dto.AddMemberRequest{
				UserIdentifier: "testuser",
				Role:           types.RolePlanner.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			tt.mockSetup(mockQ)

			s := service.NewMemberService(mockQ)
			_, err := s.AddMember(tt.ctx, tt.req, tt.overrideWorkspaceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemberService_UpdateMemberRole(t *testing.T) {
	workspaceID := uuid.New()
	memberID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()
	roleID := uuid.New()

	tests := []struct {
		name                string
		memberID            string
		req                 dto.UpdateMemberRoleRequest
		overrideWorkspaceID *string
		ctx                 context.Context
		mockSetup           func(*MockQuerier)
		wantErr             bool
	}{
		{
			name:     "successful_update",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: userID, RoleID: roleID, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.UpdateMemberRoleFunc = func(ctx context.Context, arg store.UpdateMemberRoleParams) error {
					return nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Username: "testuser", Email: "test@example.com"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:     "invalid_member_id",
			memberID: "invalid-uuid",
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name:     "missing_workspace",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 context.Background(),
			mockSetup:           func(mq *MockQuerier) {},
			wantErr:             true,
		},
		{
			name:     "personal_workspace_rejection",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "personal"}, nil
				}
			},
			wantErr: true,
		},
		{
			name:     "forbidden_not_admin",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOperator.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name:     "member_not_found",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{}, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name:     "member_workspace_mismatch",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				otherWorkspaceID := uuid.New()
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: otherWorkspaceID, UserID: userID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:     "cannot_change_owner_role",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: ownerID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:     "invalid_role_founder",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleFounder.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: userID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:     "get_role_database_error",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: userID}, nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name:     "update_member_role_database_error",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: userID}, nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.UpdateMemberRoleFunc = func(ctx context.Context, arg store.UpdateMemberRoleParams) error {
					return fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name:     "get_user_by_id_fails_after_update",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: userID}, nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.UpdateMemberRoleFunc = func(ctx context.Context, arg store.UpdateMemberRoleParams) error {
					return nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name:     "get_member_fails_after_update",
			memberID: memberID.String(),
			req: dto.UpdateMemberRoleRequest{
				Role: types.RoleAdmin.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				callCount := 0
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					callCount++
					if callCount == 1 {
						return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: userID}, nil
					}
					return store.Member{}, fmt.Errorf("database error")
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.UpdateMemberRoleFunc = func(ctx context.Context, arg store.UpdateMemberRoleParams) error {
					return nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Username: "testuser", Email: "test@example.com"}, nil
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			tt.mockSetup(mockQ)

			s := service.NewMemberService(mockQ)
			_, err := s.UpdateMemberRole(tt.ctx, tt.memberID, tt.req, tt.overrideWorkspaceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMemberRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemberService_DeleteMember(t *testing.T) {
	workspaceID := uuid.New()
	memberID := uuid.New()
	actingUserID := uuid.New()
	memberUserID := uuid.New()
	ownerID := uuid.New()

	tests := []struct {
		name                string
		memberID            string
		overrideWorkspaceID *string
		ctx                 context.Context
		mockSetup           func(*MockQuerier)
		wantErr             bool
	}{
		{
			name:                "successful_delete",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					if arg.UserID == memberUserID {
						return types.RoleUser.String(), nil
					}
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: memberUserID}, nil
				}
				mq.DeleteMemberFunc = func(ctx context.Context, arg store.DeleteMemberParams) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:                "invalid_member_id",
			memberID:            "invalid-uuid",
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name:                "missing_workspace",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 context.Background(),
			mockSetup:           func(mq *MockQuerier) {},
			wantErr:             true,
		},
		{
			name:                "personal_workspace_rejection",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "personal"}, nil
				}
			},
			wantErr: true,
		},
		{
			name:                "forbidden_not_admin",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOperator.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name:                "member_not_found",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{}, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name:                "member_workspace_mismatch",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				otherWorkspaceID := uuid.New()
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: otherWorkspaceID, UserID: memberUserID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:                "cannot_remove_owner_by_role",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					if arg.UserID == memberUserID {
						return types.RoleOwner.String(), nil
					}
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: memberUserID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:                "cannot_remove_owner_by_user_id",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					if arg.UserID == ownerID {
						return types.RoleUser.String(), nil
					}
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: ownerID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:                "get_member_role_no_rows",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					if arg.UserID == memberUserID {
						return "", sql.ErrNoRows
					}
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: memberUserID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:                "get_member_role_database_error",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					if arg.UserID == memberUserID {
						return "", fmt.Errorf("database error")
					}
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: memberUserID}, nil
				}
			},
			wantErr: true,
		},
		{
			name:                "delete_member_database_error",
			memberID:            memberID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, actingUserID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization", OwnerUserID: ownerID}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					if arg.UserID == memberUserID {
						return types.RoleUser.String(), nil
					}
					return types.RoleOwner.String(), nil
				}
				mq.GetMemberFunc = func(ctx context.Context, id uuid.UUID) (store.Member, error) {
					return store.Member{MemberID: memberID, WorkspaceID: workspaceID, UserID: memberUserID}, nil
				}
				mq.DeleteMemberFunc = func(ctx context.Context, arg store.DeleteMemberParams) error {
					return fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			tt.mockSetup(mockQ)

			s := service.NewMemberService(mockQ)
			err := s.DeleteMember(tt.ctx, tt.memberID, tt.overrideWorkspaceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
