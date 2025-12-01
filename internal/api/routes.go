package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (a *App) setupRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"time":    time.Now().Format(time.RFC3339),
			"version": "1.0.0", // TODO: Get from build info
		})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", a.authHandler.Login)
		}

		users := v1.Group("/users")
		{
			users.POST("", a.userHandler.CreateUser)
			users.GET("/:id", a.userHandler.GetUser)
			users.GET("", a.userHandler.ListUsers)
		}
	}
}
