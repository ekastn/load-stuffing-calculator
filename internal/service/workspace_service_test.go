package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/mocks"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

func ctxWithUserAndRole(role types.Role, userID uuid.UUID) context.Context {
	ctx := auth.WithRole(context.Background(), role.String())
	return auth.WithUserID(ctx, userID.String())
}

func TestWorkspaceService_ListWorkspaces_FounderUsesGlobalList(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	founderID := uuid.New()
	ctx := ctxWithUserAndRole(types.RoleFounder, founderID)

	calledAll := false
	mockQ.ListWorkspacesAllFunc = func(ctx context.Context, arg store.ListWorkspacesAllParams) ([]store.ListWorkspacesAllRow, error) {
		calledAll = true
		return []store.ListWorkspacesAllRow{
			{
				WorkspaceID:   uuid.New(),
				Type:          "organization",
				Name:          "Acme",
				OwnerUserID:   uuid.New(),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				OwnerUsername: "alice",
				OwnerEmail:    "alice@example.com",
			},
		}, nil
	}
	mockQ.ListWorkspacesForUserFunc = func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
		t.Fatalf("ListWorkspacesForUser should not be called for founder")
		return nil, nil
	}

	resp, err := svc.ListWorkspaces(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !calledAll {
		t.Fatalf("expected ListWorkspacesAll to be called")
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(resp))
	}
	if resp[0].OwnerUsername == nil || *resp[0].OwnerUsername != "alice" {
		t.Fatalf("expected owner_username to map")
	}
	if resp[0].OwnerEmail == nil || *resp[0].OwnerEmail != "alice@example.com" {
		t.Fatalf("expected owner_email to map")
	}
}

func TestWorkspaceService_ListWorkspaces_NonFounderUsesScopedList(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	userID := uuid.New()
	ctx := ctxWithUserAndRole(types.RoleOwner, userID)

	calledUser := false
	mockQ.ListWorkspacesForUserFunc = func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
		calledUser = true
		return []store.Workspace{{WorkspaceID: uuid.New(), Type: "organization", Name: "Org", OwnerUserID: userID}}, nil
	}
	mockQ.ListWorkspacesAllFunc = func(ctx context.Context, arg store.ListWorkspacesAllParams) ([]store.ListWorkspacesAllRow, error) {
		t.Fatalf("ListWorkspacesAll should not be called for non-founder")
		return nil, nil
	}

	_, err := svc.ListWorkspaces(ctx, 1, 10)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !calledUser {
		t.Fatalf("expected ListWorkspacesForUser to be called")
	}
}

func TestWorkspaceService_CreateWorkspace_FounderCanSetTypeAndOwner(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	founderID := uuid.New()
	newOwnerID := uuid.New()
	ctx := ctxWithUserAndRole(types.RoleFounder, founderID)

	mockQ.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
		if name != types.RolePersonal.String() {
			t.Fatalf("expected role lookup for personal, got %q", name)
		}
		return store.GetRoleByNameRow{RoleID: uuid.New(), Name: name}, nil
	}

	mockQ.CreateWorkspaceFunc = func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
		if arg.Type != "personal" {
			t.Fatalf("expected type=personal, got %q", arg.Type)
		}
		if arg.OwnerUserID != newOwnerID {
			t.Fatalf("expected owner_user_id=%v, got %v", newOwnerID, arg.OwnerUserID)
		}
		return store.Workspace{WorkspaceID: uuid.New(), Type: arg.Type, Name: arg.Name, OwnerUserID: arg.OwnerUserID}, nil
	}

	mockQ.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
		if arg.UserID != newOwnerID {
			t.Fatalf("expected membership created for owner")
		}
		return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID, RoleID: arg.RoleID}, nil
	}

	typ := "personal"
	ownerStr := newOwnerID.String()
	_, err := svc.CreateWorkspace(ctx, dto.CreateWorkspaceRequest{Name: "P", Type: &typ, OwnerUserID: &ownerStr})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWorkspaceService_CreateWorkspace_NonFounderCannotSetOwnerOrType(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	userID := uuid.New()
	ctx := ctxWithUserAndRole(types.RoleOwner, userID)

	typ := "personal"
	ownerStr := uuid.New().String()
	_, err := svc.CreateWorkspace(ctx, dto.CreateWorkspaceRequest{Name: "X", Type: &typ, OwnerUserID: &ownerStr})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestWorkspaceService_UpdateWorkspace_OwnershipTransferPersonalAddsPersonalMember(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	wsID := uuid.New()
	currentOwner := uuid.New()
	newOwner := uuid.New()

	ctx := ctxWithUserAndRole(types.RoleFounder, uuid.New())

	mockQ.GetWorkspaceFunc = func(ctx context.Context, workspaceID uuid.UUID) (store.Workspace, error) {
		if workspaceID != wsID {
			t.Fatalf("unexpected workspace id")
		}
		return store.Workspace{WorkspaceID: wsID, Type: "personal", Name: "P", OwnerUserID: currentOwner}, nil
	}

	mockQ.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
		return store.Member{}, errors.New("not found")
	}

	mockQ.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
		if name != types.RolePersonal.String() {
			t.Fatalf("expected role lookup personal, got %q", name)
		}
		return store.GetRoleByNameRow{RoleID: uuid.New(), Name: name}, nil
	}

	createdMember := false
	mockQ.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
		createdMember = true
		if arg.UserID != newOwner {
			t.Fatalf("expected member for new owner")
		}
		return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID, RoleID: arg.RoleID}, nil
	}

	transferred := false
	mockQ.TransferWorkspaceOwnershipFunc = func(ctx context.Context, arg store.TransferWorkspaceOwnershipParams) error {
		transferred = true
		if arg.OwnerUserID != newOwner {
			t.Fatalf("expected transfer to new owner")
		}
		return nil
	}

	mockQ.UpdateWorkspaceFunc = func(ctx context.Context, arg store.UpdateWorkspaceParams) error {
		return nil
	}

	newOwnerStr := newOwner.String()
	_, err := svc.UpdateWorkspace(ctx, wsID.String(), dto.UpdateWorkspaceRequest{OwnerUserID: &newOwnerStr})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !createdMember {
		t.Fatalf("expected CreateMember to be called")
	}
	if !transferred {
		t.Fatalf("expected TransferWorkspaceOwnership to be called")
	}
}

func TestWorkspaceService_DeleteWorkspace_FounderCanDeletePersonal(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	wsID := uuid.New()
	ctx := ctxWithUserAndRole(types.RoleFounder, uuid.New())

	mockQ.GetWorkspaceFunc = func(ctx context.Context, workspaceID uuid.UUID) (store.Workspace, error) {
		return store.Workspace{WorkspaceID: wsID, Type: "personal", Name: "P", OwnerUserID: uuid.New()}, nil
	}

	deleted := false
	mockQ.DeleteWorkspaceFunc = func(ctx context.Context, workspaceID uuid.UUID) error {
		deleted = true
		return nil
	}

	if err := svc.DeleteWorkspace(ctx, wsID.String()); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !deleted {
		t.Fatalf("expected delete to occur")
	}
}

func TestWorkspaceService_DeleteWorkspace_NonFounderCannotDeletePersonal(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	wsID := uuid.New()
	ownerID := uuid.New()
	ctx := ctxWithUserAndRole(types.RoleOwner, ownerID)

	mockQ.GetWorkspaceFunc = func(ctx context.Context, workspaceID uuid.UUID) (store.Workspace, error) {
		return store.Workspace{WorkspaceID: wsID, Type: "personal", Name: "P", OwnerUserID: ownerID}, nil
	}

	if err := svc.DeleteWorkspace(ctx, wsID.String()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestWorkspaceService_UpdateWorkspace_OwnershipTransfer_ExistingMemberDoesNotDuplicate(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	wsID := uuid.New()
	newOwner := uuid.New()

	ctx := ctxWithUserAndRole(types.RoleFounder, uuid.New())

	mockQ.GetWorkspaceFunc = func(ctx context.Context, workspaceID uuid.UUID) (store.Workspace, error) {
		return store.Workspace{WorkspaceID: wsID, Type: "organization", Name: "O", OwnerUserID: uuid.New()}, nil
	}

	mockQ.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
		return store.Member{MemberID: uuid.New()}, nil
	}

	mockQ.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
		t.Fatalf("did not expect CreateMember when member exists")
		return store.Member{}, nil
	}

	mockQ.UpdateWorkspaceFunc = func(ctx context.Context, arg store.UpdateWorkspaceParams) error {
		return nil
	}

	mockQ.TransferWorkspaceOwnershipFunc = func(ctx context.Context, arg store.TransferWorkspaceOwnershipParams) error {
		return nil
	}

	newOwnerStr := newOwner.String()
	_, err := svc.UpdateWorkspace(ctx, wsID.String(), dto.UpdateWorkspaceRequest{OwnerUserID: &newOwnerStr})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWorkspaceService_ListWorkspaces_GlobalListErrorPropagates(t *testing.T) {
	mockQ := &mocks.MockQuerier{}
	svc := service.NewWorkspaceService(mockQ)

	ctx := ctxWithUserAndRole(types.RoleFounder, uuid.New())
	expected := errors.New("boom")
	mockQ.ListWorkspacesAllFunc = func(ctx context.Context, arg store.ListWorkspacesAllParams) ([]store.ListWorkspacesAllRow, error) {
		return nil, expected
	}

	_, err := svc.ListWorkspaces(ctx, 1, 10)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, expected) {
		t.Fatalf("expected wrapped error")
	}
}

// Table-driven tests for UpdateWorkspace
func TestWorkspaceService_UpdateWorkspace(t *testing.T) {
	workspaceID := uuid.New()
	ownerID := uuid.New()
	newOwnerID := uuid.New()
	roleID := uuid.New()

	tests := []struct {
		name      string
		id        string
		req       dto.UpdateWorkspaceRequest
		ctx       context.Context
		mockSetup func(*mocks.MockQuerier)
		wantErr   bool
	}{
		{
			name: "successful_name_update",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				Name: stringPtr("New Name"),
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Old Name",
						OwnerUserID: ownerID,
					}, nil
				}
				mq.UpdateWorkspaceFunc = func(ctx context.Context, arg store.UpdateWorkspaceParams) error {
					if arg.Name != "New Name" {
						t.Errorf("expected name 'New Name', got %q", arg.Name)
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "invalid_workspace_id",
			id:   "invalid-uuid",
			req: dto.UpdateWorkspaceRequest{
				Name: stringPtr("Test"),
			},
			ctx:       ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {},
			wantErr:   true,
		},
		{
			name: "workspace_not_found",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				Name: stringPtr("Test"),
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{}, errors.New("not found")
				}
			},
			wantErr: true,
		},
		{
			name: "forbidden_non_owner",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				Name: stringPtr("Test"),
			},
			ctx: ctxWithUserAndRole(types.RoleAdmin, uuid.New()), // Different user, not owner
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test",
						OwnerUserID: ownerID, // Different from context user
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "invalid_new_owner_uuid",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				OwnerUserID: stringPtr("invalid-uuid"),
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test",
						OwnerUserID: ownerID,
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "database_error_on_update",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				Name: stringPtr("New Name"),
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Old Name",
						OwnerUserID: ownerID,
					}, nil
				}
				mq.UpdateWorkspaceFunc = func(ctx context.Context, arg store.UpdateWorkspaceParams) error {
					return errors.New("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "database_error_on_transfer",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				OwnerUserID: stringPtr(newOwnerID.String()),
			},
			ctx: ctxWithUserAndRole(types.RoleFounder, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test",
						OwnerUserID: ownerID,
					}, nil
				}
				mq.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
					return store.Member{}, errors.New("not found")
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID}, nil
				}
				mq.TransferWorkspaceOwnershipFunc = func(ctx context.Context, arg store.TransferWorkspaceOwnershipParams) error {
					return errors.New("transfer failed")
				}
			},
			wantErr: true,
		},
		{
			name: "role_lookup_error_on_transfer",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				OwnerUserID: stringPtr(newOwnerID.String()),
			},
			ctx: ctxWithUserAndRole(types.RoleFounder, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test",
						OwnerUserID: ownerID,
					}, nil
				}
				mq.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
					return store.Member{}, errors.New("not found")
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{}, errors.New("role not found")
				}
			},
			wantErr: true,
		},
		{
			name: "create_member_error_on_transfer",
			id:   workspaceID.String(),
			req: dto.UpdateWorkspaceRequest{
				OwnerUserID: stringPtr(newOwnerID.String()),
			},
			ctx: ctxWithUserAndRole(types.RoleFounder, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test",
						OwnerUserID: ownerID,
					}, nil
				}
				mq.GetMemberByWorkspaceAndUserFunc = func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
					return store.Member{}, errors.New("not found")
				}
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{}, errors.New("create member failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.mockSetup(mockQ)

			svc := service.NewWorkspaceService(mockQ)
			_, err := svc.UpdateWorkspace(tt.ctx, tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Table-driven tests for DeleteWorkspace
func TestWorkspaceService_DeleteWorkspace(t *testing.T) {
	workspaceID := uuid.New()
	ownerID := uuid.New()

	tests := []struct {
		name      string
		id        string
		ctx       context.Context
		mockSetup func(*mocks.MockQuerier)
		wantErr   bool
	}{
		{
			name: "successful_delete_organization",
			id:   workspaceID.String(),
			ctx:  ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test Org",
						OwnerUserID: ownerID,
					}, nil
				}
				mq.DeleteWorkspaceFunc = func(ctx context.Context, id uuid.UUID) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:      "invalid_workspace_id",
			id:        "invalid-uuid",
			ctx:       ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {},
			wantErr:   true,
		},
		{
			name: "workspace_not_found",
			id:   workspaceID.String(),
			ctx:  ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{}, errors.New("not found")
				}
			},
			wantErr: true,
		},
		{
			name: "forbidden_non_owner_organization",
			id:   workspaceID.String(),
			ctx:  ctxWithUserAndRole(types.RoleAdmin, uuid.New()), // Not the owner
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test",
						OwnerUserID: ownerID, // Different from context user
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "database_error_on_delete",
			id:   workspaceID.String(),
			ctx:  ctxWithUserAndRole(types.RoleOwner, ownerID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetWorkspaceFunc = func(ctx context.Context, id uuid.UUID) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: workspaceID,
						Type:        "organization",
						Name:        "Test",
						OwnerUserID: ownerID,
					}, nil
				}
				mq.DeleteWorkspaceFunc = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.mockSetup(mockQ)

			svc := service.NewWorkspaceService(mockQ)
			err := svc.DeleteWorkspace(tt.ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Table-driven tests for CreateWorkspace
func TestWorkspaceService_CreateWorkspace(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()

	tests := []struct {
		name      string
		req       dto.CreateWorkspaceRequest
		ctx       context.Context
		mockSetup func(*mocks.MockQuerier)
		wantErr   bool
	}{
		{
			name: "successful_create_organization",
			req: dto.CreateWorkspaceRequest{
				Name: "Test Org",
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, userID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					if name != types.RoleOwner.String() {
						t.Errorf("expected role %q, got %q", types.RoleOwner, name)
					}
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateWorkspaceFunc = func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
					if arg.Type != "organization" {
						t.Errorf("expected type organization, got %q", arg.Type)
					}
					return store.Workspace{
						WorkspaceID: uuid.New(),
						Type:        arg.Type,
						Name:        arg.Name,
						OwnerUserID: arg.OwnerUserID,
					}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{MemberID: uuid.New()}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "missing_user_id",
			req: dto.CreateWorkspaceRequest{
				Name: "Test",
			},
			ctx:       context.Background(), // No user ID in context
			mockSetup: func(mq *mocks.MockQuerier) {},
			wantErr:   true,
		},
		{
			name: "founder_invalid_workspace_type",
			req: dto.CreateWorkspaceRequest{
				Name: "Test",
				Type: stringPtr("invalid"),
			},
			ctx:       ctxWithUserAndRole(types.RoleFounder, userID),
			mockSetup: func(mq *mocks.MockQuerier) {},
			wantErr:   true,
		},
		{
			name: "founder_invalid_owner_uuid",
			req: dto.CreateWorkspaceRequest{
				Name:        "Test",
				OwnerUserID: stringPtr("invalid-uuid"),
			},
			ctx:       ctxWithUserAndRole(types.RoleFounder, userID),
			mockSetup: func(mq *mocks.MockQuerier) {},
			wantErr:   true,
		},
		{
			name: "role_lookup_error",
			req: dto.CreateWorkspaceRequest{
				Name: "Test",
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, userID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{}, errors.New("role not found")
				}
			},
			wantErr: true,
		},
		{
			name: "database_error_on_create_workspace",
			req: dto.CreateWorkspaceRequest{
				Name: "Test",
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, userID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateWorkspaceFunc = func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
					return store.Workspace{}, errors.New("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "database_error_on_create_member",
			req: dto.CreateWorkspaceRequest{
				Name: "Test",
			},
			ctx: ctxWithUserAndRole(types.RoleOwner, userID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.GetRoleByNameFunc = func(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
					return store.GetRoleByNameRow{RoleID: roleID, Name: name}, nil
				}
				mq.CreateWorkspaceFunc = func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
					return store.Workspace{
						WorkspaceID: uuid.New(),
						Type:        arg.Type,
						Name:        arg.Name,
						OwnerUserID: arg.OwnerUserID,
					}, nil
				}
				mq.CreateMemberFunc = func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
					return store.Member{}, errors.New("create member failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.mockSetup(mockQ)

			svc := service.NewWorkspaceService(mockQ)
			_, err := svc.CreateWorkspace(tt.ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Table-driven tests for ListWorkspaces
func TestWorkspaceService_ListWorkspaces(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		page      int32
		limit     int32
		ctx       context.Context
		mockSetup func(*mocks.MockQuerier)
		wantErr   bool
		wantLen   int
	}{
		{
			name:  "pagination_defaults",
			page:  0,
			limit: 0,
			ctx:   ctxWithUserAndRole(types.RoleOwner, userID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.ListWorkspacesForUserFunc = func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
					if arg.Limit != 10 {
						t.Errorf("expected default limit 10, got %d", arg.Limit)
					}
					if arg.Offset != 0 {
						t.Errorf("expected offset 0, got %d", arg.Offset)
					}
					return []store.Workspace{}, nil
				}
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:  "missing_user_id",
			page:  1,
			limit: 10,
			ctx:   context.Background(), // No user ID
			mockSetup: func(mq *mocks.MockQuerier) {
				// Founder check will fail, then userIDFromContext will fail
			},
			wantErr: true,
		},
		{
			name:  "database_error_for_user",
			page:  1,
			limit: 10,
			ctx:   ctxWithUserAndRole(types.RoleOwner, userID),
			mockSetup: func(mq *mocks.MockQuerier) {
				mq.ListWorkspacesForUserFunc = func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.mockSetup(mockQ)

			svc := service.NewWorkspaceService(mockQ)
			resp, err := svc.ListWorkspaces(tt.ctx, tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(resp) != tt.wantLen {
				t.Errorf("ListWorkspaces() returned %d workspaces, want %d", len(resp), tt.wantLen)
			}
		})
	}
}
