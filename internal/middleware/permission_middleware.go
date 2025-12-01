package middleware

import (
	"net/http"

	"github.com/ekastn/load-stuffing-calculator/internal/cache"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/gin-gonic/gin"
)

// Permission checks if the user has the required permission.
func Permission(q store.Querier, permCache *cache.PermissionCache, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User role not found in context")
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "Invalid role type")
			c.Abort()
			return
		}

		// Check if admin (optimization / fallback)
		if userRole == "admin" {
			c.Next()
			return
		}

		// Check cache
		permissions, found := permCache.Get(userRole)
		if !found {
			// Cache miss, query DB
			var err error
			permissions, err = q.GetPermissionsByRole(c.Request.Context(), userRole)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, "Failed to fetch permissions")
				c.Abort()
				return
			}
			// Update cache
			permCache.Set(userRole, permissions)
		}

		for _, p := range permissions {
			if p == requiredPermission {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Insufficient permissions")
		c.Abort()
	}
}
