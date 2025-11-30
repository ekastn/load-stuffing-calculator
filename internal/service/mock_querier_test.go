package service_test

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

	// Unused methods of Querier interface, implemented to satisfy the interface
	CreateUserFunc         func(ctx context.Context, arg store.CreateUserParams) (store.User, error)
	GetRefreshTokenFunc    func(ctx context.Context, token string) (store.GetRefreshTokenRow, error)
	GetUserByIDFunc        func(ctx context.Context, userID uuid.UUID) (store.GetUserByIDRow, error)
	ListUsersFunc          func(ctx context.Context, arg store.ListUsersParams) ([]store.ListUsersRow, error)
	RevokeRefreshTokenFunc func(ctx context.Context, token string) error
	UpdateUserFunc         func(ctx context.Context, arg store.UpdateUserParams) error
}

// GetUserByUsername mocks the corresponding method from store.Querier.
func (m *MockQuerier) GetUserByUsername(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return store.GetUserByUsernameRow{}, fmt.Errorf("GetUserByUsername not implemented")
}

// CreateRefreshToken mocks the corresponding method from store.Querier.
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

// Ensure MockQuerier implements store.Querier (compile-time check)
var _ store.Querier = (*MockQuerier)(nil)
