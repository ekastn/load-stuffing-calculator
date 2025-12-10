package mocks

import (
	"context"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of service.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

// MockUserService is a mock implementation of service.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context, page, limit int32) ([]dto.UserResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.UserResponse), args.Error(1)
}

// MockRoleService is a mock implementation of service.RoleService
type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoleResponse), args.Error(1)
}

func (m *MockRoleService) GetRole(ctx context.Context, id string) (*dto.RoleResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RoleResponse), args.Error(1)
}

func (m *MockRoleService) ListRoles(ctx context.Context, page, limit int32) ([]dto.RoleResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.RoleResponse), args.Error(1)
}

func (m *MockRoleService) UpdateRole(ctx context.Context, id string, req dto.UpdateRoleRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockRoleService) DeleteRole(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockPermissionService is a mock implementation of service.PermissionService
type MockPermissionService struct {
	mock.Mock
}

func (m *MockPermissionService) CreatePermission(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) GetPermission(ctx context.Context, id string) (*dto.PermissionResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) ListPermissions(ctx context.Context, page, limit int32) ([]dto.PermissionResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.PermissionResponse), args.Error(1)
}

func (m *MockPermissionService) UpdatePermission(ctx context.Context, id string, req dto.UpdatePermissionRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockPermissionService) DeletePermission(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockContainerService is a mock implementation of service.ContainerService
type MockContainerService struct {
	mock.Mock
}

func (m *MockContainerService) CreateContainer(ctx context.Context, req dto.CreateContainerRequest) (*dto.ContainerResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ContainerResponse), args.Error(1)
}

func (m *MockContainerService) GetContainer(ctx context.Context, id string) (*dto.ContainerResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ContainerResponse), args.Error(1)
}

func (m *MockContainerService) ListContainers(ctx context.Context, page, limit int32) ([]dto.ContainerResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ContainerResponse), args.Error(1)
}

func (m *MockContainerService) UpdateContainer(ctx context.Context, id string, req dto.UpdateContainerRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockContainerService) DeleteContainer(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockProductService is a mock implementation of service.ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) ListProducts(ctx context.Context, page, limit int32) ([]dto.ProductResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockPlanService is a mock implementation of service.PlanService
type MockPlanService struct {
	mock.Mock
}

func (m *MockPlanService) CreateCompletePlan(ctx context.Context, req dto.CreatePlanRequest) (*dto.CreatePlanResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CreatePlanResponse), args.Error(1)
}

func (m *MockPlanService) GetPlan(ctx context.Context, id string) (*dto.PlanDetailResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PlanDetailResponse), args.Error(1)
}

func (m *MockPlanService) ListPlans(ctx context.Context, page, limit int32) ([]dto.PlanListItem, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.PlanListItem), args.Error(1)
}

func (m *MockPlanService) UpdatePlan(ctx context.Context, id string, req dto.UpdatePlanRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockPlanService) DeletePlan(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPlanService) AddPlanItem(ctx context.Context, planID string, req dto.AddPlanItemRequest) (*dto.PlanItemDetail, error) {
	args := m.Called(ctx, planID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PlanItemDetail), args.Error(1)
}

func (m *MockPlanService) GetPlanItem(ctx context.Context, planID, itemID string) (*dto.PlanItemDetail, error) {
	args := m.Called(ctx, planID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PlanItemDetail), args.Error(1)
}

func (m *MockPlanService) UpdatePlanItem(ctx context.Context, planID, itemID string, req dto.UpdatePlanItemRequest) error {
	args := m.Called(ctx, planID, itemID, req)
	return args.Error(0)
}

func (m *MockPlanService) DeletePlanItem(ctx context.Context, planID, itemID string) error {
	args := m.Called(ctx, planID, itemID)
	return args.Error(0)
}
