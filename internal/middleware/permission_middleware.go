package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ekastn/load-stuffing-calculator/internal/cache"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/gin-gonic/gin"
)

type permissionQuerier interface {
	GetPermissionsByRole(ctx context.Context, name string) ([]string, error)
}

type PermissionMiddleware interface {
	Require(requiredPermission string) gin.HandlerFunc
	Any(requiredPermissions ...string) gin.HandlerFunc
	All(requiredPermissions ...string) gin.HandlerFunc
}

type permissionMiddleware struct {
	q         permissionQuerier
	permCache *cache.PermissionCache
}

func NewPermissionMiddleware(q permissionQuerier, permCache *cache.PermissionCache) PermissionMiddleware {
	return &permissionMiddleware{q: q, permCache: permCache}
}

func permissionMatches(granted string, required string) bool {
	if granted == "*" {
		return true
	}
	if granted == required {
		return true
	}

	// Handle resource wildcards like "plan:*".
	if strings.HasSuffix(granted, ":*") {
		prefix := strings.TrimSuffix(granted, "*")
		return strings.HasPrefix(required, prefix)
	}

	return false
}

func (m *permissionMiddleware) getPermissionsForRequest(c *gin.Context) ([]string, bool) {
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User role not found in context")
		c.Abort()
		return nil, false
	}

	userRole, ok := role.(string)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid role type")
		c.Abort()
		return nil, false
	}

	permissions, found := m.permCache.Get(userRole)
	if found {
		return permissions, true
	}

	permissions, err := m.q.GetPermissionsByRole(c.Request.Context(), userRole)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch permissions")
		c.Abort()
		return nil, false
	}

	m.permCache.Set(userRole, permissions)
	return permissions, true
}

func (m *permissionMiddleware) Require(requiredPermission string) gin.HandlerFunc {
	return m.Any(requiredPermission)
}

func (m *permissionMiddleware) Any(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, ok := m.getPermissionsForRequest(c)
		if !ok {
			return
		}

		for _, required := range requiredPermissions {
			for _, granted := range permissions {
				if permissionMatches(granted, required) {
					c.Next()
					return
				}
			}
		}

		response.Error(c, http.StatusForbidden, "Insufficient permissions")
		c.Abort()
	}
}

func (m *permissionMiddleware) All(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, ok := m.getPermissionsForRequest(c)
		if !ok {
			return
		}

		for _, required := range requiredPermissions {
			hasRequired := false
			for _, granted := range permissions {
				if permissionMatches(granted, required) {
					hasRequired = true
					break
				}
			}
			if !hasRequired {
				response.Error(c, http.StatusForbidden, "Insufficient permissions")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// Permission checks if the user has the required permission.
func Permission(q permissionQuerier, permCache *cache.PermissionCache, requiredPermission string) gin.HandlerFunc {
	return NewPermissionMiddleware(q, permCache).Require(requiredPermission)
}

// PermissionAny allows request if user has at least one required permission.
func PermissionAny(q permissionQuerier, permCache *cache.PermissionCache, requiredPermissions ...string) gin.HandlerFunc {
	return NewPermissionMiddleware(q, permCache).Any(requiredPermissions...)
}

// PermissionAll allows request only if user has every required permission.
func PermissionAll(q permissionQuerier, permCache *cache.PermissionCache, requiredPermissions ...string) gin.HandlerFunc {
	return NewPermissionMiddleware(q, permCache).All(requiredPermissions...)
}
