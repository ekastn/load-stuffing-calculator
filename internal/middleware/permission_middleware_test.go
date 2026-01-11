package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/cache"
	"github.com/ekastn/load-stuffing-calculator/internal/middleware"
	"github.com/ekastn/load-stuffing-calculator/internal/mocks"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockQuerier is an alias to the shared mock implementation.
type MockQuerier = mocks.MockQuerier

func TestPermissionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("authorized_by_permission", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"article:create", "article:read"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "editor")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("authorized_by_global_star", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"*"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", types.RoleAdmin.String())

		middleware.Permission(mockQ, permCache, "anything")(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("authorized_by_resource_star", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"plan:*"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", types.RolePlanner.String())

		middleware.Permission(mockQ, permCache, "plan:update")(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"article:read"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "viewer")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("db_error", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return nil, errors.New("db down")
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "editor")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("authorized_cached", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				t.Error("GetPermissionsByRole should not be called")
				return nil, nil
			},
		}
		permCache := cache.NewPermissionCache()

		// Pre-fill cache
		permCache.Set("editor", []string{"article:create"})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", "editor")

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("forbidden_when_other_resource_star", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"plan:*"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", types.RolePlanner.String())

		middleware.Permission(mockQ, permCache, "user:create")(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("permission_any_allows_if_any_match", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"plan:read"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", types.RoleOperator.String())

		middleware.PermissionAny(mockQ, permCache, "plan:update", "plan:read")(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("permission_all_requires_all", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"plan:read"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", types.RoleOperator.String())

		middleware.PermissionAll(mockQ, permCache, "plan:read", "plan:update")(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("permission_all_success_when_all_granted", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				return []string{"plan:read", "plan:update", "plan:delete"}, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", types.RolePlanner.String())

		middleware.PermissionAll(mockQ, permCache, "plan:read", "plan:update")(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("role_not_in_context", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				t.Error("GetPermissionsByRole should not be called")
				return nil, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		// Don't set role in context

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "User role not found in context")
	})

	t.Run("invalid_role_type", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				t.Error("GetPermissionsByRole should not be called")
				return nil, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		// Set role as int instead of string
		c.Set("role", 123)

		middleware.Permission(mockQ, permCache, "article:create")(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid role type")
	})

	t.Run("permission_all_role_not_in_context", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				t.Error("GetPermissionsByRole should not be called")
				return nil, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		// Don't set role in context

		middleware.PermissionAll(mockQ, permCache, "plan:read", "plan:update")(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "User role not found in context")
	})

	t.Run("permission_any_role_not_in_context", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetPermissionsByRoleFunc: func(ctx context.Context, name string) ([]string, error) {
				t.Error("GetPermissionsByRole should not be called")
				return nil, nil
			},
		}
		permCache := cache.NewPermissionCache()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		// Don't set role in context

		middleware.PermissionAny(mockQ, permCache, "plan:read", "plan:update")(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "User role not found in context")
	})
}
