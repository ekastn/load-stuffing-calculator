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

	t.Run("admin_bypass", func(t *testing.T) {
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
		c.Set("role", "admin")

		middleware.Permission(mockQ, permCache, "anything")(c)

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
}
