package middleware

import (
	"net/http"
	"strings"

	"load-stuffing-calculator/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Try get from Cookie
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			tokenString = cookie
		}

		// 2. If no cookie, try Authorization header
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token provided"})
			return
		}

		// 3. Verify token
		userID, err := authService.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token"})
			return
		}

		// 4. Set user ID in context
		c.Set("userID", userID)
		c.Next()
	}
}
