package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	service "github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

func TestMemberService_DeleteMember_CannotRemoveOwner(t *testing.T) {
	workspaceID := uuid.New()
	ownerID := uuid.New()
	memberID := uuid.New()

	ctx := auth.WithUserID(auth.WithRole(auth.WithWorkspaceID(context.Background(), workspaceID.String()), types.RoleAdmin.String()), uuid.New().String())

	mockQ := &MockQuerier{
		GetWorkspaceFunc: func(ctx context.Context, wid uuid.UUID) (store.Workspace, error) {
			return store.Workspace{WorkspaceID: wid, OwnerUserID: ownerID, Type: "organization", Name: "Org"}, nil
		},
		GetMemberFunc: func(ctx context.Context, mid uuid.UUID) (store.Member, error) {
			return store.Member{MemberID: mid, WorkspaceID: workspaceID, UserID: ownerID}, nil
		},
		GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
			return types.RoleOwner.String(), nil
		},
	}

	svc := service.NewMemberService(mockQ)
	if err := svc.DeleteMember(ctx, memberID.String(), nil); err == nil {
		t.Fatalf("expected error")
	}
}

func TestMemberService_AddMember_BlocksPersonalWorkspace(t *testing.T) {
	workspaceID := uuid.New()

	ctx := auth.WithUserID(auth.WithRole(auth.WithWorkspaceID(context.Background(), workspaceID.String()), types.RoleOwner.String()), uuid.New().String())

	mockQ := &MockQuerier{
		GetWorkspaceFunc: func(ctx context.Context, wid uuid.UUID) (store.Workspace, error) {
			return store.Workspace{WorkspaceID: wid, OwnerUserID: uuid.New(), Type: "personal", Name: "Personal"}, nil
		},
	}

	svc := service.NewMemberService(mockQ)
	_, err := svc.AddMember(ctx, dto.AddMemberRequest{UserIdentifier: "user@example.com", Role: types.RolePlanner.String()}, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestInviteService_CreateInvite_BlocksPersonalWorkspace(t *testing.T) {
	workspaceID := uuid.New()

	ctx := auth.WithUserID(auth.WithRole(auth.WithWorkspaceID(context.Background(), workspaceID.String()), types.RoleOwner.String()), uuid.New().String())

	mockQ := &MockQuerier{
		GetWorkspaceFunc: func(ctx context.Context, wid uuid.UUID) (store.Workspace, error) {
			return store.Workspace{WorkspaceID: wid, OwnerUserID: uuid.New(), Type: "personal", Name: "Personal"}, nil
		},
	}

	svc := service.NewInviteService(mockQ, "secret")
	_, err := svc.CreateInvite(ctx, dto.CreateInviteRequest{Email: "x@example.com", Role: types.RolePlanner.String()}, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMemberService_ListMembers_FounderOverrideWorkspaceID(t *testing.T) {
	workspaceID := uuid.New()
	otherWorkspaceID := uuid.New()
	override := otherWorkspaceID.String()

	ctx := auth.WithUserID(auth.WithRole(auth.WithWorkspaceID(context.Background(), workspaceID.String()), types.RoleFounder.String()), uuid.New().String())

	calledWorkspaceID := uuid.Nil
	mockQ := &MockQuerier{
		ListMembersByWorkspaceFunc: func(ctx context.Context, arg store.ListMembersByWorkspaceParams) ([]store.ListMembersByWorkspaceRow, error) {
			calledWorkspaceID = arg.WorkspaceID
			return []store.ListMembersByWorkspaceRow{}, nil
		},
		GetMemberRoleNameByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
			return types.RoleFounder.String(), nil
		},
	}

	svc := service.NewMemberService(mockQ)
	_, err := svc.ListMembers(ctx, 1, 10, &override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledWorkspaceID != otherWorkspaceID {
		t.Fatalf("expected override workspace id")
	}
}

func TestInviteService_AcceptInvite_ReturnsNewAccessToken(t *testing.T) {
	workspaceID := uuid.New()
	userID := uuid.New()
	inviteID := uuid.New()
	roleID := uuid.New()

	ctx := auth.WithUserID(auth.WithRole(auth.WithWorkspaceID(context.Background(), workspaceID.String()), types.RolePlanner.String()), userID.String())

	expires := time.Now().Add(24 * time.Hour)
	inv := store.Invite{InviteID: inviteID, WorkspaceID: workspaceID, Email: "x@example.com", RoleID: roleID, TokenHash: "", InvitedByUserID: uuid.New(), ExpiresAt: &expires}

	mockQ := &MockQuerier{
		GetUserByIDFunc: func(ctx context.Context, uid uuid.UUID) (store.GetUserByIDRow, error) {
			return store.GetUserByIDRow{UserID: uid, Username: "u", Email: "x@example.com", RoleID: uuid.New(), RoleName: types.RolePlanner.String(), CreatedAt: time.Now()}, nil
		},
		GetInviteByTokenHashFunc: func(ctx context.Context, tokenHash string) (store.Invite, error) {
			return inv, nil
		},
		GetWorkspaceFunc: func(ctx context.Context, wid uuid.UUID) (store.Workspace, error) {
			return store.Workspace{WorkspaceID: wid, OwnerUserID: uuid.New(), Type: "organization", Name: "Org"}, nil
		},
		GetRoleFunc: func(ctx context.Context, rid uuid.UUID) (store.Role, error) {
			return store.Role{RoleID: rid, Name: types.RolePlanner.String()}, nil
		},
		GetMemberByWorkspaceAndUserFunc: func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
			return store.Member{}, sql.ErrNoRows
		},
		CreateMemberFunc: func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
			return store.Member{MemberID: uuid.New(), WorkspaceID: arg.WorkspaceID, UserID: arg.UserID, RoleID: arg.RoleID, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
		AcceptInviteFunc: func(ctx context.Context, arg store.AcceptInviteParams) error {
			return nil
		},
	}

	svc := service.NewInviteService(mockQ, "secret")
	resp, err := svc.AcceptInvite(ctx, dto.AcceptInviteRequest{Token: "raw"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatalf("expected access token")
	}
	if resp.ActiveWorkspaceID != workspaceID.String() {
		t.Fatalf("expected active workspace id")
	}
}
