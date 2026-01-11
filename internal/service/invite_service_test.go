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

func TestInviteService_ListInvites(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	inviteID := uuid.New()
	invitedByUserID := uuid.New()

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
				mq.ListInvitesByWorkspaceFunc = func(ctx context.Context, arg store.ListInvitesByWorkspaceParams) ([]store.ListInvitesByWorkspaceRow, error) {
					return []store.ListInvitesByWorkspaceRow{
						{
							InviteID:          inviteID,
							WorkspaceID:       workspaceID,
							Email:             "test@example.com",
							RoleName:          types.RoleOperator.String(),
							InvitedByUserID:   invitedByUserID,
							InvitedByUsername: "admin",
							CreatedAt:         time.Now(),
						},
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
				mq.ListInvitesByWorkspaceFunc = func(ctx context.Context, arg store.ListInvitesByWorkspaceParams) ([]store.ListInvitesByWorkspaceRow, error) {
					if arg.Limit != 10 || arg.Offset != 0 {
						return nil, fmt.Errorf("expected defaults: page=1, limit=10")
					}
					return []store.ListInvitesByWorkspaceRow{}, nil
				}
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:                "missing_workspace_context",
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
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleAdmin.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return "", fmt.Errorf("database connection error")
				}
			},
			wantErr: true,
		},
		{
			name:                "list_invites_database_error",
			page:                1,
			limit:               10,
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.ListInvitesByWorkspaceFunc = func(ctx context.Context, arg store.ListInvitesByWorkspaceParams) ([]store.ListInvitesByWorkspaceRow, error) {
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

			s := service.NewInviteService(mockQ, "test-secret")
			invites, err := s.ListInvites(tt.ctx, tt.page, tt.limit, tt.overrideWorkspaceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListInvites() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(invites) != tt.wantLen {
				t.Errorf("ListInvites() returned %d invites, want %d", len(invites), tt.wantLen)
			}
		})
	}
}

func TestInviteService_CreateInvite(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	roleID := uuid.New()
	inviteID := uuid.New()

	tests := []struct {
		name                string
		req                 dto.CreateInviteRequest
		overrideWorkspaceID *string
		ctx                 context.Context
		mockSetup           func(*MockQuerier)
		wantErr             bool
	}{
		{
			name: "successful_create",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateInviteFunc = func(ctx context.Context, arg store.CreateInviteParams) (store.Invite, error) {
					expiresAt := time.Now().Add(7 * 24 * time.Hour)
					return store.Invite{
						InviteID:        inviteID,
						WorkspaceID:     arg.WorkspaceID,
						Email:           arg.Email,
						RoleID:          arg.RoleID,
						TokenHash:       arg.TokenHash,
						InvitedByUserID: arg.InvitedByUserID,
						ExpiresAt:       &expiresAt,
						CreatedAt:       time.Now(),
					}, nil
				}
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Username: "admin", Email: "admin@example.com"}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "personal_workspace",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  types.RoleOperator.String(),
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
			name: "invalid_role",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  "invalid",
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name: "forbidden_not_admin",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOperator.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name: "missing_workspace",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 context.Background(),
			mockSetup:           func(mq *MockQuerier) {},
			wantErr:             true,
		},
		{
			name: "get_role_database_error",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{}, fmt.Errorf("role lookup error")
				}
			},
			wantErr: true,
		},
		{
			name: "database_error_create_invite",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateInviteFunc = func(ctx context.Context, arg store.CreateInviteParams) (store.Invite, error) {
					return store.Invite{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "successful_create_with_get_user_failure_best_effort",
			req: dto.CreateInviteRequest{
				Email: "invite@example.com",
				Role:  types.RoleOperator.String(),
			},
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateInviteFunc = func(ctx context.Context, arg store.CreateInviteParams) (store.Invite, error) {
					expiresAt := time.Now().Add(7 * 24 * time.Hour)
					return store.Invite{
						InviteID:        inviteID,
						WorkspaceID:     arg.WorkspaceID,
						Email:           arg.Email,
						RoleID:          arg.RoleID,
						TokenHash:       arg.TokenHash,
						InvitedByUserID: arg.InvitedByUserID,
						ExpiresAt:       &expiresAt,
						CreatedAt:       time.Now(),
					}, nil
				}
				// GetUserByID fails (best-effort), should not affect invite creation
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{}, fmt.Errorf("user fetch error")
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			tt.mockSetup(mockQ)

			s := service.NewInviteService(mockQ, "test-secret")
			resp, err := s.CreateInvite(tt.ctx, tt.req, tt.overrideWorkspaceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateInvite() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && resp == nil {
				t.Error("CreateInvite() returned nil response, expected valid response")
			}
			if !tt.wantErr && resp != nil && resp.Token == "" {
				t.Error("CreateInvite() returned empty token")
			}
		})
	}
}

func TestInviteService_RevokeInvite(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	inviteID := uuid.New()

	tests := []struct {
		name                string
		inviteID            string
		overrideWorkspaceID *string
		ctx                 context.Context
		mockSetup           func(*MockQuerier)
		wantErr             bool
	}{
		{
			name:                "successful_revoke",
			inviteID:            inviteID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.RevokeInviteFunc = func(ctx context.Context, arg store.RevokeInviteParams) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:                "invalid_invite_id",
			inviteID:            "invalid-uuid",
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name:                "personal_workspace",
			inviteID:            inviteID.String(),
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
			name:                "forbidden_not_admin",
			inviteID:            inviteID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOperator.String(), nil
				}
			},
			wantErr: true,
		},
		{
			name:                "missing_workspace_context",
			inviteID:            inviteID.String(),
			overrideWorkspaceID: nil,
			ctx:                 context.Background(),
			mockSetup:           func(mq *MockQuerier) {},
			wantErr:             true,
		},
		{
			name:                "get_workspace_database_error",
			inviteID:            inviteID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name:                "permission_check_database_error",
			inviteID:            inviteID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleAdmin.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return "", fmt.Errorf("database connection error")
				}
			},
			wantErr: true,
		},
		{
			name:                "database_error",
			inviteID:            inviteID.String(),
			overrideWorkspaceID: nil,
			ctx:                 ctxWithRoleWorkspaceAndUser(types.RoleOwner.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetMemberRoleNameByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
					return types.RoleOwner.String(), nil
				}
				mq.RevokeInviteFunc = func(ctx context.Context, arg store.RevokeInviteParams) error {
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

			s := service.NewInviteService(mockQ, "test-secret")
			err := s.RevokeInvite(tt.ctx, tt.inviteID, tt.overrideWorkspaceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RevokeInvite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInviteService_AcceptInvite(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	inviteID := uuid.New()
	roleID := uuid.New()
	tokenHash := "hashed_token"

	tests := []struct {
		name      string
		req       dto.AcceptInviteRequest
		ctx       context.Context
		mockSetup func(*MockQuerier)
		wantErr   bool
	}{
		{
			name: "successful_accept",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
						TokenHash:   tokenHash,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetRoleFunc = func(ctx context.Context, id uuid.UUID) (store.Role, error) {
					return store.Role{RoleID: roleID, Name: types.RoleOperator.String()}, nil
				}
				mq.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
					return store.Member{}, sql.ErrNoRows
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID}, nil
				}
				mq.AcceptInviteFunc = func(ctx context.Context, arg store.AcceptInviteParams) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "unauthorized_no_user_id",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx:       context.Background(),
			mockSetup: func(mq *MockQuerier) {},
			wantErr:   true,
		},
		{
			name: "invalid_invite_token",
			req: dto.AcceptInviteRequest{
				Token: "invalid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{}, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name: "email_mismatch",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID: userID,
						Email:  "user@example.com",
					}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "different@example.com",
						RoleID:      roleID,
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "personal_workspace",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "personal"}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "invalid_role_owner",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetRoleFunc = func(ctx context.Context, id uuid.UUID) (store.Role, error) {
					return store.Role{RoleID: roleID, Name: types.RoleOwner.String()}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "already_member",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetRoleFunc = func(ctx context.Context, id uuid.UUID) (store.Role, error) {
					return store.Role{RoleID: roleID, Name: types.RoleOperator.String()}, nil
				}
				mq.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
					return store.Member{}, sql.ErrNoRows
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{}, &pgconn.PgError{Code: "23505"}
				}
				mq.AcceptInviteFunc = func(ctx context.Context, arg store.AcceptInviteParams) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "get_user_by_id_fails",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{}, fmt.Errorf("user not found")
				}
			},
			wantErr: true,
		},
		{
			name: "get_role_fails",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetRoleFunc = func(ctx context.Context, id uuid.UUID) (store.Role, error) {
					return store.Role{}, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name: "invalid_role_founder",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetRoleFunc = func(ctx context.Context, id uuid.UUID) (store.Role, error) {
					return store.Role{RoleID: roleID, Name: types.RoleFounder.String()}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "create_member_database_error_non_unique",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetRoleFunc = func(ctx context.Context, id uuid.UUID) (store.Role, error) {
					return store.Role{RoleID: roleID, Name: types.RoleOperator.String()}, nil
				}
				mq.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
					return store.Member{}, sql.ErrNoRows
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{}, fmt.Errorf("database connection error")
				}
			},
			wantErr: true,
		},
		{
			name: "accept_invite_database_error",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{
						UserID:   userID,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
						TokenHash:   tokenHash,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{WorkspaceID: workspaceID, Type: "organization"}, nil
				}
				mq.GetRoleFunc = func(ctx context.Context, id uuid.UUID) (store.Role, error) {
					return store.Role{RoleID: roleID, Name: types.RoleOperator.String()}, nil
				}
				mq.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
					return store.Member{}, sql.ErrNoRows
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID}, nil
				}
				mq.AcceptInviteFunc = func(ctx context.Context, arg store.AcceptInviteParams) error {
					return fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "get_workspace_database_error",
			req: dto.AcceptInviteRequest{
				Token: "valid_token",
			},
			ctx: ctxWithRoleWorkspaceAndUser(types.RoleOperator.String(), workspaceID, userID),
			mockSetup: func(mq *MockQuerier) {
				mq.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (store.GetUserByIDRow, error) {
					return store.GetUserByIDRow{UserID: userID, Email: "test@example.com"}, nil
				}
				mq.GetInviteByTokenHashFunc = func(ctx context.Context, hash string) (store.Invite, error) {
					return store.Invite{
						InviteID:    inviteID,
						WorkspaceID: workspaceID,
						Email:       "test@example.com",
						RoleID:      roleID,
					}, nil
				}
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			tt.mockSetup(mockQ)

			s := service.NewInviteService(mockQ, "test-secret-key-for-jwt")
			resp, err := s.AcceptInvite(tt.ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("AcceptInvite() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && resp == nil {
				t.Error("AcceptInvite() returned nil response, expected valid response")
			}
			if !tt.wantErr && resp != nil && resp.AccessToken == "" {
				t.Error("AcceptInvite() returned empty access token")
			}
		})
	}
}
