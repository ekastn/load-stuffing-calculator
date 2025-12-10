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
	CreateRefreshTokenFunc func(ctx context.Context, arg store.CreateRefreshTokenParams) error
	GetRoleByNameFunc      func(ctx context.Context, name string) (store.GetRoleByNameRow, error)

	CreateUserFunc           func(ctx context.Context, arg store.CreateUserParams) (store.User, error)
	GetRefreshTokenFunc      func(ctx context.Context, token string) (store.GetRefreshTokenRow, error)
	GetUserByIDFunc          func(ctx context.Context, userID uuid.UUID) (store.GetUserByIDRow, error)
	ListUsersFunc            func(ctx context.Context, arg store.ListUsersParams) ([]store.ListUsersRow, error)
	RevokeRefreshTokenFunc   func(ctx context.Context, token string) error
	UpdateUserFunc           func(ctx context.Context, arg store.UpdateUserParams) error
	GetPermissionsByRoleFunc func(ctx context.Context, name string) ([]string, error)
	CreateRoleFunc           func(ctx context.Context, arg store.CreateRoleParams) (store.Role, error)
	GetRoleFunc              func(ctx context.Context, id uuid.UUID) (store.Role, error)
	ListRolesFunc            func(ctx context.Context, arg store.ListRolesParams) ([]store.Role, error)
	UpdateRoleFunc           func(ctx context.Context, arg store.UpdateRoleParams) error
	DeleteRoleFunc           func(ctx context.Context, id uuid.UUID) error
	CreatePermissionFunc     func(ctx context.Context, arg store.CreatePermissionParams) (store.Permission, error)
	GetPermissionFunc        func(ctx context.Context, id uuid.UUID) (store.Permission, error)
	ListPermissionsFunc      func(ctx context.Context, arg store.ListPermissionsParams) ([]store.Permission, error)
	UpdatePermissionFunc     func(ctx context.Context, arg store.UpdatePermissionParams) error
	DeletePermissionFunc     func(ctx context.Context, id uuid.UUID) error
	CreateContainerFunc      func(ctx context.Context, arg store.CreateContainerParams) (store.Container, error)
	GetContainerFunc         func(ctx context.Context, id uuid.UUID) (store.Container, error)
	ListContainersFunc       func(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error)
	UpdateContainerFunc      func(ctx context.Context, arg store.UpdateContainerParams) error
	DeleteContainerFunc      func(ctx context.Context, id uuid.UUID) error
	CreateProductFunc        func(ctx context.Context, arg store.CreateProductParams) (store.Product, error)
	GetProductFunc           func(ctx context.Context, id uuid.UUID) (store.Product, error)
	ListProductsFunc         func(ctx context.Context, arg store.ListProductsParams) ([]store.Product, error)
	UpdateProductFunc        func(ctx context.Context, arg store.UpdateProductParams) error
	DeleteProductFunc        func(ctx context.Context, id uuid.UUID) error
	CreateLoadPlanFunc       func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error)
	AddLoadItemFunc          func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error)
	GetLoadPlanFunc          func(ctx context.Context, planID uuid.UUID) (store.LoadPlan, error)
	ListLoadPlansFunc        func(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error)
	UpdatePlanStatusFunc     func(ctx context.Context, arg store.UpdatePlanStatusParams) error
	ListLoadItemsFunc        func(ctx context.Context, planID *uuid.UUID) ([]store.LoadItem, error)
	GetLoadItemFunc          func(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error)
	UpdateLoadItemFunc       func(ctx context.Context, arg store.UpdateLoadItemParams) error
	DeleteLoadItemFunc       func(ctx context.Context, arg store.DeleteLoadItemParams) error
	UpdateLoadPlanFunc       func(ctx context.Context, arg store.UpdateLoadPlanParams) error
	DeleteLoadPlanFunc       func(ctx context.Context, planID uuid.UUID) error
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

func (m *MockQuerier) GetLoadPlan(ctx context.Context, planID uuid.UUID) (store.LoadPlan, error) {
	if m.GetLoadPlanFunc != nil {
		return m.GetLoadPlanFunc(ctx, planID)
	}
	return store.LoadPlan{}, fmt.Errorf("GetLoadPlan not implemented")
}

func (m *MockQuerier) ListLoadPlans(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error) {
	if m.ListLoadPlansFunc != nil {
		return m.ListLoadPlansFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListLoadPlans not implemented")
}

func (m *MockQuerier) UpdatePlanStatus(ctx context.Context, arg store.UpdatePlanStatusParams) error {
	if m.UpdatePlanStatusFunc != nil {
		return m.UpdatePlanStatusFunc(ctx, arg)
	}
	return fmt.Errorf("UpdatePlanStatus not implemented")
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

func (m *MockQuerier) DeleteLoadPlan(ctx context.Context, planID uuid.UUID) error {
	if m.DeleteLoadPlanFunc != nil {
		return m.DeleteLoadPlanFunc(ctx, planID)
	}
	return fmt.Errorf("DeleteLoadPlan not implemented")
}

func (m *MockQuerier) CreateProduct(ctx context.Context, arg store.CreateProductParams) (store.Product, error) {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(ctx, arg)
	}
	return store.Product{}, fmt.Errorf("CreateProduct not implemented")
}

func (m *MockQuerier) GetProduct(ctx context.Context, id uuid.UUID) (store.Product, error) {
	if m.GetProductFunc != nil {
		return m.GetProductFunc(ctx, id)
	}
	return store.Product{}, fmt.Errorf("GetProduct not implemented")
}

func (m *MockQuerier) ListProducts(ctx context.Context, arg store.ListProductsParams) ([]store.Product, error) {
	if m.ListProductsFunc != nil {
		return m.ListProductsFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListProducts not implemented")
}

func (m *MockQuerier) UpdateProduct(ctx context.Context, arg store.UpdateProductParams) error {
	if m.UpdateProductFunc != nil {
		return m.UpdateProductFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateProduct not implemented")
}

func (m *MockQuerier) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(ctx, id)
	}
	return fmt.Errorf("DeleteProduct not implemented")
}

func (m *MockQuerier) CreateContainer(ctx context.Context, arg store.CreateContainerParams) (store.Container, error) {
	if m.CreateContainerFunc != nil {
		return m.CreateContainerFunc(ctx, arg)
	}
	return store.Container{}, fmt.Errorf("CreateContainer not implemented")
}

func (m *MockQuerier) GetContainer(ctx context.Context, id uuid.UUID) (store.Container, error) {
	if m.GetContainerFunc != nil {
		return m.GetContainerFunc(ctx, id)
	}
	return store.Container{}, fmt.Errorf("GetContainer not implemented")
}

func (m *MockQuerier) ListContainers(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error) {
	if m.ListContainersFunc != nil {
		return m.ListContainersFunc(ctx, arg)
	}
	return nil, fmt.Errorf("ListContainers not implemented")
}

func (m *MockQuerier) UpdateContainer(ctx context.Context, arg store.UpdateContainerParams) error {
	if m.UpdateContainerFunc != nil {
		return m.UpdateContainerFunc(ctx, arg)
	}
	return fmt.Errorf("UpdateContainer not implemented")
}

func (m *MockQuerier) DeleteContainer(ctx context.Context, id uuid.UUID) error {
	if m.DeleteContainerFunc != nil {
		return m.DeleteContainerFunc(ctx, id)
	}
	return fmt.Errorf("DeleteContainer not implemented")
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

var _ store.Querier = (*MockQuerier)(nil)
