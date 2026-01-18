package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (a *App) setupRoutes(r *gin.Engine) {
	// Semua routes diawali dengan /api/v1 untuk versioning
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", a.HealthCheck)

		// Container routes: CRUD operasi untuk container
		containers := v1.Group("/containers")
		{
			containers.GET("", a.containerHandler.List)
			containers.GET("/:id", a.containerHandler.GetByID)
			containers.POST("", a.containerHandler.Create)
		}

		// Plan routes: pembuatan dan kalkulasi plan
		plans := v1.Group("/plans")
		{
			plans.POST("", a.planHandler.Create)
			plans.POST("/:id/calculate", a.planHandler.Calculate)
		}
	}
}

func (a *App) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	})
}
