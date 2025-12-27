package middleware

import (
	"net/http"
	"strings"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/gin-gonic/gin"
)

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		if claims.WorkspaceID != nil {
			c.Set("workspace_id", *claims.WorkspaceID)
		}

		ctx := auth.WithUserID(c.Request.Context(), claims.UserID)
		ctx = auth.WithRole(ctx, claims.Role)
		if claims.WorkspaceID != nil {
			ctx = auth.WithWorkspaceID(ctx, *claims.WorkspaceID)
		}
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
