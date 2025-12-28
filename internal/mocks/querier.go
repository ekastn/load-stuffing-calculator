package mocks

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

// MockQuerier implements store.Querier for testing purposes.
type MockQuerier struct {
	GetUserByUsernameFunc  func(ctx context.Context, username string) (store.GetUserByUsernameRow, error)
	GetUserByEmailFunc     func(ctx context.Context, email string) (store.GetUserByEmailRow, error)
	CreateRefreshTokenFunc func(ctx context.Context, arg store.CreateRefreshTokenParams) error
	GetRoleByNameFunc      func(ctx context.Context, name string) (store.GetRoleByNameRow, error)

	CreateUserFunc                  func(ctx context.Context, arg store.CreateUserParams) (store.User, error)
	GetRefreshTokenFunc             func(ctx context.Context, token string) (store.GetRefreshTokenRow, error)
	UpdateRefreshTokenWorkspaceFunc func(ctx context.Context, arg store.UpdateRefreshTokenWorkspaceParams) error
	GetUserByIDFunc                 func(ctx context.Context, userID uuid.UUID) (store.GetUserByIDRow, error)
	ListUsersFunc                   func(ctx context.Context, arg store.ListUsersParams) ([]store.ListUsersRow, error)
	RevokeRefreshTokenFunc          func(ctx context.Context, token string) error
	UpdateUserFunc                  func(ctx context.Context, arg store.UpdateUserParams) error
	DeleteUserFunc                  func(ctx context.Context, userID uuid.UUID) error
	UpdateUserPasswordFunc          func(ctx context.Context, arg store.UpdateUserPasswordParams) error
	GetPermissionsByRoleFunc        func(ctx context.Context, name string) ([]string, error)
	CreateRoleFunc                  func(ctx context.Context, arg store.CreateRoleParams) (store.Role, error)
	GetRoleFunc                     func(ctx context.Context, id uuid.UUID) (store.Role, error)
	ListRolesFunc                   func(ctx context.Context, arg store.ListRolesParams) ([]store.Role, error)
	UpdateRoleFunc                  func(ctx context.Context, arg store.UpdateRoleParams) error
	DeleteRoleFunc                  func(ctx context.Context, id uuid.UUID) error
	CreatePermissionFunc            func(ctx context.Context, arg store.CreatePermissionParams) (store.Permission, error)
	GetPermissionFunc               func(ctx context.Context, id uuid.UUID) (store.Permission, error)
	ListPermissionsFunc             func(ctx context.Context, arg store.ListPermissionsParams) ([]store.Permission, error)
	UpdatePermissionFunc            func(ctx context.Context, arg store.UpdatePermissionParams) error
	DeletePermissionFunc            func(ctx context.Context, id uuid.UUID) error
	CreateContainerFunc             func(ctx context.Context, arg store.CreateContainerParams) (store.Container, error)
	GetContainerFunc                func(ctx context.Context, arg store.GetContainerParams) (store.Container, error)
	GetContainerAnyFunc             func(ctx context.Context, containerID uuid.UUID) (store.Container, error)
	ListContainersFunc              func(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error)
	ListContainersAllFunc           func(ctx context.Context, arg store.ListContainersAllParams) ([]store.Container, error)
	UpdateContainerFunc             func(ctx context.Context, arg store.UpdateContainerParams) error
	UpdateContainerAnyFunc          func(ctx context.Context, arg store.UpdateContainerAnyParams) error
	DeleteContainerFunc             func(ctx context.Context, arg store.DeleteContainerParams) error
	DeleteContainerAnyFunc          func(ctx context.Context, containerID uuid.UUID) error
	CreateProductFunc               func(ctx context.Context, arg store.CreateProductParams) (store.Product, error)
	GetProductFunc                  func(ctx context.Context, arg store.GetProductParams) (store.Product, error)
	GetProductAnyFunc               func(ctx context.Context, productID uuid.UUID) (store.Product, error)
	ListProductsFunc                func(ctx context.Context, arg store.ListProductsParams) ([]store.Product, error)
	ListProductsAllFunc             func(ctx context.Context, arg store.ListProductsAllParams) ([]store.Product, error)
	UpdateProductFunc               func(ctx context.Context, arg store.UpdateProductParams) error
	UpdateProductAnyFunc            func(ctx context.Context, arg store.UpdateProductAnyParams) error
	DeleteProductFunc               func(ctx context.Context, arg store.DeleteProductParams) error
	DeleteProductAnyFunc            func(ctx context.Context, productID uuid.UUID) error
	CreateLoadPlanFunc              func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error)
	AddLoadItemFunc                 func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error)
	GetLoadPlanFunc                 func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error)
	GetLoadPlanAnyFunc              func(ctx context.Context, planID uuid.UUID) (store.LoadPlan, error)
	GetLoadPlanForGuestFunc         func(ctx context.Context, arg store.GetLoadPlanForGuestParams) (store.LoadPlan, error)
	ListLoadPlansFunc               func(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error)
	ListLoadPlansAllFunc            func(ctx context.Context, arg store.ListLoadPlansAllParams) ([]store.LoadPlan, error)
	ListLoadPlansForGuestFunc       func(ctx context.Context, arg store.ListLoadPlansForGuestParams) ([]store.LoadPlan, error)
	UpdatePlanStatusFunc            func(ctx context.Context, arg store.UpdatePlanStatusParams) error
	UpdatePlanStatusAnyFunc         func(ctx context.Context, arg store.UpdatePlanStatusAnyParams) error
	ListLoadItemsFunc               func(ctx context.Context, planID *uuid.UUID) ([]store.LoadItem, error)
	GetLoadItemFunc                 func(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error)
	UpdateLoadItemFunc              func(ctx context.Context, arg store.UpdateLoadItemParams) error
	DeleteLoadItemFunc              func(ctx context.Context, arg store.DeleteLoadItemParams) error
	UpdateLoadPlanFunc              func(ctx context.Context, arg store.UpdateLoadPlanParams) error
	DeleteLoadPlanFunc              func(ctx context.Context, arg store.DeleteLoadPlanParams) error
	CreatePlanResultFunc            func(ctx context.Context, arg store.CreatePlanResultParams) (store.PlanResult, error)
	DeletePlanResultsFunc           func(ctx context.Context, planID *uuid.UUID) error
	CreatePlanPlacementFunc         func(ctx context.Context, arg []store.CreatePlanPlacementParams) (int64, error)
	GetPlanResultFunc               func(ctx context.Context, planID *uuid.UUID) (store.PlanResult, error)
	ListPlanPlacementsFunc          func(ctx context.Context, resultID *uuid.UUID) ([]store.PlanPlacement, error)
	CountPlansByCreatorFunc         func(ctx context.Context, arg store.CountPlansByCreatorParams) (int64, error)
	ClaimPlansFromGuestFunc         func(ctx context.Context, arg store.ClaimPlansFromGuestParams) error

	AddRolePermissionFunc         func(ctx context.Context, arg store.AddRolePermissionParams) error
	DeleteRolePermissionsFunc     func(ctx context.Context, roleID uuid.UUID) error
	GetRolePermissionsFunc        func(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error)
	CountActivePlansFunc          func(ctx context.Context) (int64, error)
	CountCompletedPlansFunc       func(ctx context.Context) (int64, error)
	CountCompletedPlansTodayFunc  func(ctx context.Context) (int64, error)
	CountContainersFunc           func(ctx context.Context) (int64, error)
	CountTotalItemsFunc           func(ctx context.Context) (int64, error)
	CountTotalUsersFunc           func(ctx context.Context) (int64, error)
	GetAvgVolumeUtilizationFunc   func(ctx context.Context) (float64, error)
	GetPlanStatusDistributionFunc func(ctx context.Context) ([]store.GetPlanStatusDistributionRow, error)

	AcceptInviteFunc                        func(ctx context.Context, arg store.AcceptInviteParams) error
	CreateInviteFunc                        func(ctx context.Context, arg store.CreateInviteParams) (store.Invite, error)
	GetInviteByTokenHashFunc                func(ctx context.Context, tokenHash string) (store.Invite, error)
	ListInvitesByWorkspaceFunc              func(ctx context.Context, arg store.ListInvitesByWorkspaceParams) ([]store.ListInvitesByWorkspaceRow, error)
	RevokeInviteFunc                        func(ctx context.Context, arg store.RevokeInviteParams) error
	CreateWorkspaceFunc                     func(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error)
	GetWorkspaceFunc                        func(ctx context.Context, workspaceID uuid.UUID) (store.Workspace, error)
	ListWorkspacesByOwnerFunc               func(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error)
	UpdateWorkspaceFunc                     func(ctx context.Context, arg store.UpdateWorkspaceParams) error
	DeleteWorkspaceFunc                     func(ctx context.Context, workspaceID uuid.UUID) error
	TransferWorkspaceOwnershipFunc          func(ctx context.Context, arg store.TransferWorkspaceOwnershipParams) error
	CreateMemberFunc                        func(ctx context.Context, arg store.CreateMemberParams) (store.Member, error)
	GetMemberFunc                           func(ctx context.Context, memberID uuid.UUID) (store.Member, error)
	GetMemberByWorkspaceAndUserFunc         func(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error)
	ListMembersByWorkspaceFunc              func(ctx context.Context, arg store.ListMembersByWorkspaceParams) ([]store.ListMembersByWorkspaceRow, error)
	UpdateMemberRoleFunc                    func(ctx context.Context, arg store.UpdateMemberRoleParams) error
	DeleteMemberFunc                        func(ctx context.Context, arg store.DeleteMemberParams) error
	GetPlatformRoleByUserIDFunc             func(ctx context.Context, userID uuid.UUID) (string, error)
	UpsertPlatformMemberFunc                func(ctx context.Context, arg store.UpsertPlatformMemberParams) error
	GetPersonalWorkspaceByOwnerFunc         func(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error)
	ListWorkspacesForUserFunc               func(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error)
	GetMemberRoleNameByWorkspaceAndUserFunc func(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error)
}

func (m *MockQuerier) UpdateUserPassword(ctx context.Context, arg store.UpdateUserPasswordParams) error {
	if m.UpdateUserPasswordFunc != nil {
		return m.UpdateUserPasswordFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateUserPassword not implemented")
}

func (m *MockQuerier) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, userID)
	}
	return fmt.Errorf("DeleteUser not implemented")
}

func (m *MockQuerier) ListPlanPlacements(ctx context.Context, resultID *uuid.UUID) ([]store.PlanPlacement, error) {
	if m.ListPlanPlacementsFunc != nil {
		return m.ListPlanPlacementsFunc(ctx, resultID)
	}
	return nil, fmt.Errorf("ListPlanPlacements not implemented")
}

func (m *MockQuerier) GetPlanResult(ctx context.Context, planID *uuid.UUID) (store.PlanResult, error) {
	if m.GetPlanResultFunc != nil {
		return m.GetPlanResultFunc(ctx, planID)
	}
	return store.PlanResult{}, fmt.Errorf("GetPlanResult not implemented")
}

func (m *MockQuerier) CreatePlanResult(ctx context.Context, arg store.CreatePlanResultParams) (store.PlanResult, error) {
	if m.CreatePlanResultFunc != nil {
		return m.CreatePlanResultFunc(ctx, arg)
	}
	return store.PlanResult{}, fmt.Errorf("CreatePlanResult not implemented")
}

func (m *MockQuerier) DeletePlanResults(ctx context.Context, planID *uuid.UUID) error {
	if m.DeletePlanResultsFunc != nil {
		return m.DeletePlanResultsFunc(ctx, planID)
	}
	return fmt.Errorf("DeletePlanResults not implemented")
}

func (m *MockQuerier) CreatePlanPlacement(ctx context.Context, arg []store.CreatePlanPlacementParams) (int64, error) {
	if m.CreatePlanPlacementFunc != nil {
		return m.CreatePlanPlacementFunc(ctx, arg)
	}
	return 0, fmt.Errorf("CreatePlanPlacement not implemented")
}

func (m *MockQuerier) CreateLoadPlan(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
	if m.CreateLoadPlanFunc != nil {
		return m.CreateLoadPlanFunc(ctx, arg)
	}
	return store.LoadPlan{}, fmt.Errorf("CreateLoadPlan not implemented")
}

func (m *MockQuerier) AddLoadItem(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
	if m.AddLoadItemFunc != nil {
		return m.AddLoadItemFunc(ctx, arg)
	}
	return store.LoadItem{}, fmt.Errorf("AddLoadItem not implemented")
}

func (m *MockQuerier) GetLoadPlan(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
	if m.GetLoadPlanFunc != nil {
		return m.GetLoadPlanFunc(ctx, arg)
	}
	return store.LoadPlan{}, fmt.Errorf("GetLoadPlan not implemented")
}

func (m *MockQuerier) GetLoadPlanAny(ctx context.Context, planID uuid.UUID) (store.LoadPlan, error) {
	if m.GetLoadPlanAnyFunc != nil {
		return m.GetLoadPlanAnyFunc(ctx, planID)
	}
	return store.LoadPlan{}, fmt.Errorf("GetLoadPlanAny not implemented")
}

func (m *MockQuerier) GetLoadPlanForGuest(ctx context.Context, arg store.GetLoadPlanForGuestParams) (store.LoadPlan, error) {
	if m.GetLoadPlanForGuestFunc != nil {
		return m.GetLoadPlanForGuestFunc(ctx, arg)
	}
	return store.LoadPlan{}, fmt.Errorf("GetLoadPlanForGuest not implemented")
}

func (m *MockQuerier) ListLoadPlans(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error) {
	if m.ListLoadPlansFunc != nil {
		return m.ListLoadPlansFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListLoadPlans not implemented")
}

func (m *MockQuerier) ListLoadPlansAll(ctx context.Context, arg store.ListLoadPlansAllParams) ([]store.LoadPlan, error) {
	if m.ListLoadPlansAllFunc != nil {
		return m.ListLoadPlansAllFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListLoadPlansAll not implemented")
}

func (m *MockQuerier) ListLoadPlansForGuest(ctx context.Context, arg store.ListLoadPlansForGuestParams) ([]store.LoadPlan, error) {
	if m.ListLoadPlansForGuestFunc != nil {
		return m.ListLoadPlansForGuestFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListLoadPlansForGuest not implemented")
}

func (m *MockQuerier) UpdatePlanStatus(ctx context.Context, arg store.UpdatePlanStatusParams) error {
	if m.UpdatePlanStatusFunc != nil {
		return m.UpdatePlanStatusFunc(ctx, arg)
	}
	return fmt.Errorf("UpdatePlanStatus not implemented")
}

func (m *MockQuerier) UpdatePlanStatusAny(ctx context.Context, arg store.UpdatePlanStatusAnyParams) error {
	if m.UpdatePlanStatusAnyFunc != nil {
		return m.UpdatePlanStatusAnyFunc(ctx, arg)
	}
	return fmt.Errorf("UpdatePlanStatusAny not implemented")
}

func (m *MockQuerier) ListLoadItems(ctx context.Context, planID *uuid.UUID) ([]store.LoadItem, error) {
	if m.ListLoadItemsFunc != nil {
		return m.ListLoadItemsFunc(ctx, planID)
	}
	return nil, fmt.Errorf("ListLoadItems not implemented")
}

func (m *MockQuerier) GetLoadItem(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error) {
	if m.GetLoadItemFunc != nil {
		return m.GetLoadItemFunc(ctx, arg)
	}
	return store.LoadItem{}, fmt.Errorf("GetLoadItem not implemented")
}

func (m *MockQuerier) UpdateLoadItem(ctx context.Context, arg store.UpdateLoadItemParams) error {
	if m.UpdateLoadItemFunc != nil {
		return m.UpdateLoadItemFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateLoadItem not implemented")
}

func (m *MockQuerier) DeleteLoadItem(ctx context.Context, arg store.DeleteLoadItemParams) error {
	if m.DeleteLoadItemFunc != nil {
		return m.DeleteLoadItemFunc(ctx, arg)
	}
	return fmt.Errorf("DeleteLoadItem not implemented")
}

func (m *MockQuerier) UpdateLoadPlan(ctx context.Context, arg store.UpdateLoadPlanParams) error {
	if m.UpdateLoadPlanFunc != nil {
		return m.UpdateLoadPlanFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateLoadPlan not implemented")
}

func (m *MockQuerier) DeleteLoadPlan(ctx context.Context, arg store.DeleteLoadPlanParams) error {
	if m.DeleteLoadPlanFunc != nil {
		return m.DeleteLoadPlanFunc(ctx, arg)
	}
	return fmt.Errorf("DeleteLoadPlan not implemented")
}

func (m *MockQuerier) CreateProduct(ctx context.Context, arg store.CreateProductParams) (store.Product, error) {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(ctx, arg)
	}
	return store.Product{}, fmt.Errorf("CreateProduct not implemented")
}

func (m *MockQuerier) GetProduct(ctx context.Context, arg store.GetProductParams) (store.Product, error) {
	if m.GetProductFunc != nil {
		return m.GetProductFunc(ctx, arg)
	}
	return store.Product{}, fmt.Errorf("GetProduct not implemented")
}

func (m *MockQuerier) GetProductAny(ctx context.Context, productID uuid.UUID) (store.Product, error) {
	if m.GetProductAnyFunc != nil {
		return m.GetProductAnyFunc(ctx, productID)
	}
	return store.Product{}, fmt.Errorf("GetProductAny not implemented")
}

func (m *MockQuerier) ListProducts(ctx context.Context, arg store.ListProductsParams) ([]store.Product, error) {
	if m.ListProductsFunc != nil {
		return m.ListProductsFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListProducts not implemented")
}

func (m *MockQuerier) ListProductsAll(ctx context.Context, arg store.ListProductsAllParams) ([]store.Product, error) {
	if m.ListProductsAllFunc != nil {
		return m.ListProductsAllFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListProductsAll not implemented")
}

func (m *MockQuerier) UpdateProduct(ctx context.Context, arg store.UpdateProductParams) error {
	if m.UpdateProductFunc != nil {
		return m.UpdateProductFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateProduct not implemented")
}

func (m *MockQuerier) UpdateProductAny(ctx context.Context, arg store.UpdateProductAnyParams) error {
	if m.UpdateProductAnyFunc != nil {
		return m.UpdateProductAnyFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateProductAny not implemented")
}

func (m *MockQuerier) DeleteProduct(ctx context.Context, arg store.DeleteProductParams) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(ctx, arg)
	}
	return fmt.Errorf("DeleteProduct not implemented")
}

func (m *MockQuerier) DeleteProductAny(ctx context.Context, productID uuid.UUID) error {
	if m.DeleteProductAnyFunc != nil {
		return m.DeleteProductAnyFunc(ctx, productID)
	}
	return fmt.Errorf("DeleteProductAny not implemented")
}

func (m *MockQuerier) CreateContainer(ctx context.Context, arg store.CreateContainerParams) (store.Container, error) {
	if m.CreateContainerFunc != nil {
		return m.CreateContainerFunc(ctx, arg)
	}
	return store.Container{}, fmt.Errorf("CreateContainer not implemented")
}

func (m *MockQuerier) GetContainer(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
	if m.GetContainerFunc != nil {
		return m.GetContainerFunc(ctx, arg)
	}
	return store.Container{}, fmt.Errorf("GetContainer not implemented")
}

func (m *MockQuerier) GetContainerAny(ctx context.Context, containerID uuid.UUID) (store.Container, error) {
	if m.GetContainerAnyFunc != nil {
		return m.GetContainerAnyFunc(ctx, containerID)
	}
	return store.Container{}, fmt.Errorf("GetContainerAny not implemented")
}

func (m *MockQuerier) ListContainers(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error) {
	if m.ListContainersFunc != nil {
		return m.ListContainersFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListContainers not implemented")
}

func (m *MockQuerier) ListContainersAll(ctx context.Context, arg store.ListContainersAllParams) ([]store.Container, error) {
	if m.ListContainersAllFunc != nil {
		return m.ListContainersAllFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListContainersAll not implemented")
}

func (m *MockQuerier) UpdateContainer(ctx context.Context, arg store.UpdateContainerParams) error {
	if m.UpdateContainerFunc != nil {
		return m.UpdateContainerFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateContainer not implemented")
}

func (m *MockQuerier) UpdateContainerAny(ctx context.Context, arg store.UpdateContainerAnyParams) error {
	if m.UpdateContainerAnyFunc != nil {
		return m.UpdateContainerAnyFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateContainerAny not implemented")
}

func (m *MockQuerier) DeleteContainer(ctx context.Context, arg store.DeleteContainerParams) error {
	if m.DeleteContainerFunc != nil {
		return m.DeleteContainerFunc(ctx, arg)
	}
	return fmt.Errorf("DeleteContainer not implemented")
}

func (m *MockQuerier) DeleteContainerAny(ctx context.Context, containerID uuid.UUID) error {
	if m.DeleteContainerAnyFunc != nil {
		return m.DeleteContainerAnyFunc(ctx, containerID)
	}
	return fmt.Errorf("DeleteContainerAny not implemented")
}

func (m *MockQuerier) CreatePermission(ctx context.Context, arg store.CreatePermissionParams) (store.Permission, error) {
	if m.CreatePermissionFunc != nil {
		return m.CreatePermissionFunc(ctx, arg)
	}
	return store.Permission{}, fmt.Errorf("CreatePermission not implemented")
}

func (m *MockQuerier) GetPermission(ctx context.Context, id uuid.UUID) (store.Permission, error) {
	if m.GetPermissionFunc != nil {
		return m.GetPermissionFunc(ctx, id)
	}
	return store.Permission{}, fmt.Errorf("GetPermission not implemented")
}

func (m *MockQuerier) ListPermissions(ctx context.Context, arg store.ListPermissionsParams) ([]store.Permission, error) {
	if m.ListPermissionsFunc != nil {
		return m.ListPermissionsFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListPermissions not implemented")
}

func (m *MockQuerier) UpdatePermission(ctx context.Context, arg store.UpdatePermissionParams) error {
	if m.UpdatePermissionFunc != nil {
		return m.UpdatePermissionFunc(ctx, arg)
	}
	return fmt.Errorf("UpdatePermission not implemented")
}

func (m *MockQuerier) DeletePermission(ctx context.Context, id uuid.UUID) error {
	if m.DeletePermissionFunc != nil {
		return m.DeletePermissionFunc(ctx, id)
	}
	return fmt.Errorf("DeletePermission not implemented")
}

func (m *MockQuerier) GetPermissionsByRole(ctx context.Context, name string) ([]string, error) {
	if m.GetPermissionsByRoleFunc != nil {
		return m.GetPermissionsByRoleFunc(ctx, name)
	}
	return nil, fmt.Errorf("GetPermissionsByRole not implemented")
}

func (m *MockQuerier) CreateRole(ctx context.Context, arg store.CreateRoleParams) (store.Role, error) {
	if m.CreateRoleFunc != nil {
		return m.CreateRoleFunc(ctx, arg)
	}
	return store.Role{}, fmt.Errorf("CreateRole not implemented")
}

func (m *MockQuerier) GetRole(ctx context.Context, id uuid.UUID) (store.Role, error) {
	if m.GetRoleFunc != nil {
		return m.GetRoleFunc(ctx, id)
	}
	return store.Role{}, fmt.Errorf("GetRole not implemented")
}

func (m *MockQuerier) ListRoles(ctx context.Context, arg store.ListRolesParams) ([]store.Role, error) {
	if m.ListRolesFunc != nil {
		return m.ListRolesFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListRoles not implemented")
}

func (m *MockQuerier) UpdateRole(ctx context.Context, arg store.UpdateRoleParams) error {
	if m.UpdateRoleFunc != nil {
		return m.UpdateRoleFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateRole not implemented")
}

func (m *MockQuerier) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if m.DeleteRoleFunc != nil {
		return m.DeleteRoleFunc(ctx, id)
	}
	return fmt.Errorf("DeleteRole not implemented")
}

func (m *MockQuerier) GetRoleByName(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
	if m.GetRoleByNameFunc != nil {
		return m.GetRoleByNameFunc(ctx, name)
	}
	return store.GetRoleByNameRow{}, fmt.Errorf("GetRoleByName not implemented")
}

func (m *MockQuerier) GetUserByUsername(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return store.GetUserByUsernameRow{}, fmt.Errorf("GetUserByUsername not implemented")
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (store.GetUserByEmailRow, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return store.GetUserByEmailRow{}, fmt.Errorf("GetUserByEmail not implemented")
}

func (m *MockQuerier) CreateRefreshToken(ctx context.Context, arg store.CreateRefreshTokenParams) error {
	if m.CreateRefreshTokenFunc != nil {
		return m.CreateRefreshTokenFunc(ctx, arg)
	}
	return fmt.Errorf("CreateRefreshToken not implemented")
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, arg)
	}
	return store.User{}, fmt.Errorf("CreateUser not implemented")
}

func (m *MockQuerier) GetRefreshToken(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
	if m.GetRefreshTokenFunc != nil {
		return m.GetRefreshTokenFunc(ctx, token)
	}
	return store.GetRefreshTokenRow{}, fmt.Errorf("GetRefreshToken not implemented")
}

func (m *MockQuerier) GetUserByID(ctx context.Context, userID uuid.UUID) (store.GetUserByIDRow, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, userID)
	}
	return store.GetUserByIDRow{}, fmt.Errorf("GetUserByID not implemented")
}

func (m *MockQuerier) ListUsers(ctx context.Context, arg store.ListUsersParams) ([]store.ListUsersRow, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListUsers not implemented")
}

func (m *MockQuerier) RevokeRefreshToken(ctx context.Context, token string) error {
	if m.RevokeRefreshTokenFunc != nil {
		return m.RevokeRefreshTokenFunc(ctx, token)
	}
	return fmt.Errorf("RevokeRefreshToken not implemented")
}

func (m *MockQuerier) UpdateUser(ctx context.Context, arg store.UpdateUserParams) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateUser not implemented")
}

func (m *MockQuerier) AddRolePermission(ctx context.Context, arg store.AddRolePermissionParams) error {
	if m.AddRolePermissionFunc != nil {
		return m.AddRolePermissionFunc(ctx, arg)
	}
	return fmt.Errorf("AddRolePermission not implemented")
}

func (m *MockQuerier) DeleteRolePermissions(ctx context.Context, roleID uuid.UUID) error {
	if m.DeleteRolePermissionsFunc != nil {
		return m.DeleteRolePermissionsFunc(ctx, roleID)
	}
	return fmt.Errorf("DeleteRolePermissions not implemented")
}

func (m *MockQuerier) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	if m.GetRolePermissionsFunc != nil {
		return m.GetRolePermissionsFunc(ctx, roleID)
	}
	return nil, fmt.Errorf("GetRolePermissions not implemented")
}

func (m *MockQuerier) CountActivePlans(ctx context.Context) (int64, error) {
	if m.CountActivePlansFunc != nil {
		return m.CountActivePlansFunc(ctx)
	}
	return 0, fmt.Errorf("CountActivePlans not implemented")
}

func (m *MockQuerier) CountCompletedPlans(ctx context.Context) (int64, error) {
	if m.CountCompletedPlansFunc != nil {
		return m.CountCompletedPlansFunc(ctx)
	}
	return 0, fmt.Errorf("CountCompletedPlans not implemented")
}

func (m *MockQuerier) CountCompletedPlansToday(ctx context.Context) (int64, error) {
	if m.CountCompletedPlansTodayFunc != nil {
		return m.CountCompletedPlansTodayFunc(ctx)
	}
	return 0, fmt.Errorf("CountCompletedPlansToday not implemented")
}

func (m *MockQuerier) CountContainers(ctx context.Context) (int64, error) {
	if m.CountContainersFunc != nil {
		return m.CountContainersFunc(ctx)
	}
	return 0, fmt.Errorf("CountContainers not implemented")
}

func (m *MockQuerier) CountTotalItems(ctx context.Context) (int64, error) {
	if m.CountTotalItemsFunc != nil {
		return m.CountTotalItemsFunc(ctx)
	}
	return 0, fmt.Errorf("CountTotalItems not implemented")
}

func (m *MockQuerier) CountTotalUsers(ctx context.Context) (int64, error) {
	if m.CountTotalUsersFunc != nil {
		return m.CountTotalUsersFunc(ctx)
	}
	return 0, fmt.Errorf("CountTotalUsers not implemented")
}

func (m *MockQuerier) CountPlansByCreator(ctx context.Context, arg store.CountPlansByCreatorParams) (int64, error) {
	if m.CountPlansByCreatorFunc != nil {
		return m.CountPlansByCreatorFunc(ctx, arg)
	}
	return 0, fmt.Errorf("CountPlansByCreator not implemented")
}

func (m *MockQuerier) ClaimPlansFromGuest(ctx context.Context, arg store.ClaimPlansFromGuestParams) error {
	if m.ClaimPlansFromGuestFunc != nil {
		return m.ClaimPlansFromGuestFunc(ctx, arg)
	}
	return fmt.Errorf("ClaimPlansFromGuest not implemented")
}

func (m *MockQuerier) GetAvgVolumeUtilization(ctx context.Context) (float64, error) {
	if m.GetAvgVolumeUtilizationFunc != nil {
		return m.GetAvgVolumeUtilizationFunc(ctx)
	}
	return 0, fmt.Errorf("GetAvgVolumeUtilization not implemented")
}

func (m *MockQuerier) GetPlanStatusDistribution(ctx context.Context) ([]store.GetPlanStatusDistributionRow, error) {
	if m.GetPlanStatusDistributionFunc != nil {
		return m.GetPlanStatusDistributionFunc(ctx)
	}
	return nil, fmt.Errorf("GetPlanStatusDistribution not implemented")
}

func (m *MockQuerier) AcceptInvite(ctx context.Context, arg store.AcceptInviteParams) error {
	if m.AcceptInviteFunc != nil {
		return m.AcceptInviteFunc(ctx, arg)
	}
	return fmt.Errorf("AcceptInvite not implemented")
}

func (m *MockQuerier) CreateInvite(ctx context.Context, arg store.CreateInviteParams) (store.Invite, error) {
	if m.CreateInviteFunc != nil {
		return m.CreateInviteFunc(ctx, arg)
	}
	return store.Invite{}, fmt.Errorf("CreateInvite not implemented")
}

func (m *MockQuerier) GetInviteByTokenHash(ctx context.Context, tokenHash string) (store.Invite, error) {
	if m.GetInviteByTokenHashFunc != nil {
		return m.GetInviteByTokenHashFunc(ctx, tokenHash)
	}
	return store.Invite{}, fmt.Errorf("GetInviteByTokenHash not implemented")
}

func (m *MockQuerier) ListInvitesByWorkspace(ctx context.Context, arg store.ListInvitesByWorkspaceParams) ([]store.ListInvitesByWorkspaceRow, error) {
	if m.ListInvitesByWorkspaceFunc != nil {
		return m.ListInvitesByWorkspaceFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListInvitesByWorkspace not implemented")
}

func (m *MockQuerier) RevokeInvite(ctx context.Context, arg store.RevokeInviteParams) error {
	if m.RevokeInviteFunc != nil {
		return m.RevokeInviteFunc(ctx, arg)
	}
	return fmt.Errorf("RevokeInvite not implemented")
}

func (m *MockQuerier) CreateWorkspace(ctx context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
	if m.CreateWorkspaceFunc != nil {
		return m.CreateWorkspaceFunc(ctx, arg)
	}
	return store.Workspace{}, fmt.Errorf("CreateWorkspace not implemented")
}

func (m *MockQuerier) GetWorkspace(ctx context.Context, workspaceID uuid.UUID) (store.Workspace, error) {
	if m.GetWorkspaceFunc != nil {
		return m.GetWorkspaceFunc(ctx, workspaceID)
	}
	return store.Workspace{}, fmt.Errorf("GetWorkspace not implemented")
}

func (m *MockQuerier) GetPersonalWorkspaceByOwner(ctx context.Context, ownerUserID uuid.UUID) (store.Workspace, error) {
	if m.GetPersonalWorkspaceByOwnerFunc != nil {
		return m.GetPersonalWorkspaceByOwnerFunc(ctx, ownerUserID)
	}
	return store.Workspace{}, fmt.Errorf("GetPersonalWorkspaceByOwner not implemented")
}

func (m *MockQuerier) ListWorkspacesForUser(ctx context.Context, arg store.ListWorkspacesForUserParams) ([]store.Workspace, error) {
	if m.ListWorkspacesForUserFunc != nil {
		return m.ListWorkspacesForUserFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListWorkspacesForUser not implemented")
}

func (m *MockQuerier) ListWorkspacesByOwner(ctx context.Context, arg store.ListWorkspacesByOwnerParams) ([]store.Workspace, error) {
	if m.ListWorkspacesByOwnerFunc != nil {
		return m.ListWorkspacesByOwnerFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListWorkspacesByOwner not implemented")
}

func (m *MockQuerier) UpdateWorkspace(ctx context.Context, arg store.UpdateWorkspaceParams) error {
	if m.UpdateWorkspaceFunc != nil {
		return m.UpdateWorkspaceFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateWorkspace not implemented")
}

func (m *MockQuerier) DeleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	if m.DeleteWorkspaceFunc != nil {
		return m.DeleteWorkspaceFunc(ctx, workspaceID)
	}
	return fmt.Errorf("DeleteWorkspace not implemented")
}

func (m *MockQuerier) TransferWorkspaceOwnership(ctx context.Context, arg store.TransferWorkspaceOwnershipParams) error {
	if m.TransferWorkspaceOwnershipFunc != nil {
		return m.TransferWorkspaceOwnershipFunc(ctx, arg)
	}
	return fmt.Errorf("TransferWorkspaceOwnership not implemented")
}

func (m *MockQuerier) CreateMember(ctx context.Context, arg store.CreateMemberParams) (store.Member, error) {
	if m.CreateMemberFunc != nil {
		return m.CreateMemberFunc(ctx, arg)
	}
	return store.Member{}, fmt.Errorf("CreateMember not implemented")
}

func (m *MockQuerier) GetMember(ctx context.Context, memberID uuid.UUID) (store.Member, error) {
	if m.GetMemberFunc != nil {
		return m.GetMemberFunc(ctx, memberID)
	}
	return store.Member{}, fmt.Errorf("GetMember not implemented")
}

func (m *MockQuerier) GetMemberByWorkspaceAndUser(ctx context.Context, arg store.GetMemberByWorkspaceAndUserParams) (store.Member, error) {
	if m.GetMemberByWorkspaceAndUserFunc != nil {
		return m.GetMemberByWorkspaceAndUserFunc(ctx, arg)
	}
	return store.Member{}, fmt.Errorf("GetMemberByWorkspaceAndUser not implemented")
}

func (m *MockQuerier) GetMemberRoleNameByWorkspaceAndUser(ctx context.Context, arg store.GetMemberRoleNameByWorkspaceAndUserParams) (string, error) {
	if m.GetMemberRoleNameByWorkspaceAndUserFunc != nil {
		return m.GetMemberRoleNameByWorkspaceAndUserFunc(ctx, arg)
	}
	return "", fmt.Errorf("GetMemberRoleNameByWorkspaceAndUser not implemented")
}

func (m *MockQuerier) ListMembersByWorkspace(ctx context.Context, arg store.ListMembersByWorkspaceParams) ([]store.ListMembersByWorkspaceRow, error) {
	if m.ListMembersByWorkspaceFunc != nil {
		return m.ListMembersByWorkspaceFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListMembersByWorkspace not implemented")
}

func (m *MockQuerier) UpdateMemberRole(ctx context.Context, arg store.UpdateMemberRoleParams) error {
	if m.UpdateMemberRoleFunc != nil {
		return m.UpdateMemberRoleFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateMemberRole not implemented")
}

func (m *MockQuerier) DeleteMember(ctx context.Context, arg store.DeleteMemberParams) error {
	if m.DeleteMemberFunc != nil {
		return m.DeleteMemberFunc(ctx, arg)
	}
	return fmt.Errorf("DeleteMember not implemented")
}

func (m *MockQuerier) GetPlatformRoleByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	if m.GetPlatformRoleByUserIDFunc != nil {
		return m.GetPlatformRoleByUserIDFunc(ctx, userID)
	}
	return "", fmt.Errorf("GetPlatformRoleByUserID not implemented")
}

func (m *MockQuerier) UpsertPlatformMember(ctx context.Context, arg store.UpsertPlatformMemberParams) error {
	if m.UpsertPlatformMemberFunc != nil {
		return m.UpsertPlatformMemberFunc(ctx, arg)
	}
	return fmt.Errorf("UpsertPlatformMember not implemented")
}

func (m *MockQuerier) UpdateRefreshTokenWorkspace(ctx context.Context, arg store.UpdateRefreshTokenWorkspaceParams) error {
	if m.UpdateRefreshTokenWorkspaceFunc != nil {
		return m.UpdateRefreshTokenWorkspaceFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateRefreshTokenWorkspace not implemented")
}

var _ store.Querier = (*MockQuerier)(nil)
