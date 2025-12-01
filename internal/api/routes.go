package api

import (
	"time"

	_ "github.com/ekastn/load-stuffing-calculator/internal/docs"
	"github.com/ekastn/load-stuffing-calculator/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *App) setupRoutes(r *gin.Engine) {
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/docs/index.html")
	})
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/doc.json")))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", a.HealthCheck)

		auth := v1.Group("/auth")
		{
			auth.POST("/login", a.authHandler.Login)
		}

		v1.Use(middleware.JWT(a.jwtSecret))

		users := v1.Group("/users", middleware.Role("admin"))
		{
			users.POST("", a.userHandler.CreateUser)
			users.GET("/:id", a.userHandler.GetUser)
			users.GET("", a.userHandler.ListUsers)
		}

		roles := v1.Group("/roles", middleware.Role("admin"))
		{
			roles.POST("", a.roleHandler.CreateRole)
			roles.GET("", a.roleHandler.ListRoles)
			roles.GET("/:id", a.roleHandler.GetRole)
			roles.PUT("/:id", a.roleHandler.UpdateRole)
			roles.DELETE("/:id", a.roleHandler.DeleteRole)
		}

		permissions := v1.Group("/permissions", middleware.Role("admin"))
		{
			permissions.POST("", a.permHandler.CreatePermission)
			permissions.GET("", a.permHandler.ListPermissions)
			permissions.GET("/:id", a.permHandler.GetPermission)
			permissions.PUT("/:id", a.permHandler.UpdatePermission)
			permissions.DELETE("/:id", a.permHandler.DeletePermission)
		}
	}
}

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the server is running and returns basic info.
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (a *App) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0", // TODO: Get from build info
	})
}
