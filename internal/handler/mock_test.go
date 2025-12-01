package handler_test

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
