package middleware

import (
	"net/http"

	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/gin-gonic/gin"
)

func Role(allowedRoles ...types.Role) gin.HandlerFunc {
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

		// Admin has access to everything
		if userRole == types.RoleAdmin.String() {
			c.Next()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed.String() {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Access denied")
		c.Abort()
	}
}