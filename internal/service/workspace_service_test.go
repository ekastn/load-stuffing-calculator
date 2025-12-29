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
