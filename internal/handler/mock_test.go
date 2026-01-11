package handler_test

import (
	"github.com/ekastn/load-stuffing-calculator/internal/mocks"
	"github.com/stretchr/testify/mock"
)

// Aliases for shared mocks
type MockAuthService = mocks.MockAuthService
type MockUserService = mocks.MockUserService
type MockRoleService = mocks.MockRoleService
type MockPermissionService = mocks.MockPermissionService
type MockContainerService = mocks.MockContainerService
type MockProductService = mocks.MockProductService
type MockPlanService = mocks.MockPlanService
type MockInviteService = mocks.MockInviteService
type MockMemberService = mocks.MockMemberService
type MockDashboardService = mocks.MockDashboardService
type MockWorkspaceService = mocks.MockWorkspaceService

// MockPermCache is a mock for PermissionCache
type MockPermCache struct {
	mock.Mock
}

func (m *MockPermCache) Invalidate() {
	m.Called()
}
