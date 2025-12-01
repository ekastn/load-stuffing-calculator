package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/cache"
	"github.com/ekastn/load-stuffing-calculator/internal/middleware"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuerier for middleware tests
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) GetPermissionsByRole(ctx context.Context, name string) ([]string, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// Implement other interface methods to satisfy store.Querier (stubs)
func (m *MockQuerier) CreateRefreshToken(ctx context.Context, arg store.CreateRefreshTokenParams) error {
	return nil
}
func (m *MockQuerier) CreateUser(ctx context.Context, arg store.CreateUserParams) (store.User, error) {
	return store.User{}, nil
}
func (m *MockQuerier) GetRefreshToken(ctx context.Context, token string) (store.GetRefreshTokenRow, error) {
	return store.GetRefreshTokenRow{}, nil
}
func (m *MockQuerier) GetRoleByName(ctx context.Context, name string) (store.GetRoleByNameRow, error) {
	return store.GetRoleByNameRow{}, nil
}
func (m *MockQuerier) GetUserByID(ctx context.Context, userID uuid.UUID) (store.GetUserByIDRow, error) {
	return store.GetUserByIDRow{}, nil
}
func (m *MockQuerier) GetUserByUsername(ctx context.Context, username string) (store.GetUserByUsernameRow, error) {
	return store.GetUserByUsernameRow{}, nil
}
func (m *MockQuerier) ListUsers(ctx context.Context, arg store.ListUsersParams) ([]store.ListUsersRow, error) {
	return nil, nil
}
func (m *MockQuerier) RevokeRefreshToken(ctx context.Context, token string) error       { return nil }
func (m *MockQuerier) UpdateUser(ctx context.Context, arg store.UpdateUserParams) error { return nil }

func TestPermissionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("authorized_by_permission", func(t *testing.T) {
		mockQ := new(MockQuerier)
		permCache := cache.NewPermissionCache()
		mockQ.On("GetPermissionsByRole", mock.Anything, "editor").Return([]string{"article:create", "article:read"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "editor")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusOK, w.Code) // Next() called (status 200 default)
		mockQ.AssertExpectations(t)
	})

	t.Run("admin_bypass", func(t *testing.T) {
		mockQ := new(MockQuerier)
		permCache := cache.NewPermissionCache()
		// Should NOT call GetPermissionsByRole

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "admin")

		middleware.Permission(mockQ, permCache, "anything")(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockQ.AssertNotCalled(t, "GetPermissionsByRole")
	})

	t.Run("forbidden", func(t *testing.T) {
		mockQ := new(MockQuerier)
		permCache := cache.NewPermissionCache()
		mockQ.On("GetPermissionsByRole", mock.Anything, "viewer").Return([]string{"article:read"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "viewer")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockQ.AssertExpectations(t)
	})

	t.Run("db_error", func(t *testing.T) {
		mockQ := new(MockQuerier)
		permCache := cache.NewPermissionCache()
		mockQ.On("GetPermissionsByRole", mock.Anything, "editor").Return(nil, errors.New("db down"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "editor")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("authorized_cached", func(t *testing.T) {
		mockQ := new(MockQuerier)
		permCache := cache.NewPermissionCache()

		// Pre-fill cache
		permCache.Set("editor", []string{"article:create"})

		// Should NOT call GetPermissionsByRole

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "editor")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockQ.AssertNotCalled(t, "GetPermissionsByRole")
	})
}
