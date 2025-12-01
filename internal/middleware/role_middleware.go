package middleware

import (
	"net/http"

	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/gin-gonic/gin"
)

func Role(allowedRoles ...string) gin.HandlerFunc {
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

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Access denied")
		c.Abort()
	}
}
