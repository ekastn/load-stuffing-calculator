package api

import (
	"time"

	"load-stuffing-calculator/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (a *App) setupRoutes(r *gin.Engine) {
	// Semua routes diawali dengan /api/v1 untuk versioning
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", a.HealthCheck)

		// Auth routes (Public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", a.authHandler.Register)
			auth.POST("/login", a.authHandler.Login)
		}

		// Protected Routes
		// Middleware akan memvalidasi token JWT sebelum request diproses
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(a.authService))
		{
			// Container routes: CRUD operasi untuk container
			containers := protected.Group("/containers")
			{
				containers.GET("", a.containerHandler.List)
				containers.GET("/:id", a.containerHandler.GetByID)
				containers.POST("", a.containerHandler.Create)
				containers.PUT("/:id", a.containerHandler.Update)
				containers.DELETE("/:id", a.containerHandler.Delete)
			}

			// Product routes: CRUD operasi untuk product
			products := protected.Group("/products")
			{
				products.GET("", a.productHandler.List)
				products.GET("/:id", a.productHandler.GetByID)
				products.POST("", a.productHandler.Create)
				products.PUT("/:id", a.productHandler.Update)
				products.DELETE("/:id", a.productHandler.Delete)
			}

			// Plan routes: CRUD dan kalkulasi
			plans := protected.Group("/plans")
			{
				plans.GET("", a.planHandler.List)
				plans.GET("/:id", a.planHandler.GetByID)
				plans.POST("", a.planHandler.Create)
				plans.PUT("/:id", a.planHandler.Update)
				plans.DELETE("/:id", a.planHandler.Delete)

				// Plan items management
				plans.POST("/:id/items", a.planHandler.AddItem)
				plans.PUT("/:id/items/:itemId", a.planHandler.UpdateItem)
				plans.DELETE("/:id/items/:itemId", a.planHandler.DeleteItem)

				// Calculation
				plans.POST("/:id/calculate", a.planHandler.Calculate)
			}

			// Dashboard
			protected.GET("/dashboard/stats", a.dashboardHandler.GetStats)
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
