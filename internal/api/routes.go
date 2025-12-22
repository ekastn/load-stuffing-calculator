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
			auth.POST("/refresh", a.authHandler.RefreshToken)
		}

		v1.Use(middleware.JWT(a.jwtSecret))

		perm := middleware.NewPermissionMiddleware(a.querier, a.permCache)

		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("", perm.Require("dashboard:read"), a.dashboardHandler.GetStats)
		}

		users := v1.Group("/users", perm.Require("user:*"))
		{
			users.POST("", a.userHandler.CreateUser)
			users.GET("/:id", a.userHandler.GetUser)
			users.GET("", a.userHandler.ListUsers)
			users.PUT("/:id", a.userHandler.UpdateUser)
			users.DELETE("/:id", a.userHandler.DeleteUser)
			users.PUT("/:id/password", a.userHandler.ChangePassword)
		}

		roles := v1.Group("/roles", perm.Require("role:*"))
		{
			roles.POST("", a.roleHandler.CreateRole)
			roles.GET("/:id", a.roleHandler.GetRole)
			roles.GET("", a.roleHandler.ListRoles)
			roles.PUT("/:id", a.roleHandler.UpdateRole)
			roles.DELETE("/:id", a.roleHandler.DeleteRole)
			roles.GET("/:id/permissions", a.roleHandler.GetRolePermissions)
			roles.PUT("/:id/permissions", a.roleHandler.UpdateRolePermissions)
		}
		permissions := v1.Group("/permissions", perm.Require("permission:*"))
		{
			permissions.POST("", a.permHandler.CreatePermission)
			permissions.GET("", a.permHandler.ListPermissions)
			permissions.GET("/:id", a.permHandler.GetPermission)
			permissions.PUT("/:id", a.permHandler.UpdatePermission)
			permissions.DELETE("/:id", a.permHandler.DeletePermission)
		}

		containers := v1.Group("/containers", perm.Require("container:*"))
		{
			containers.POST("", a.containerHandler.CreateContainer)
			containers.GET("", a.containerHandler.ListContainers)
			containers.GET("/:id", a.containerHandler.GetContainer)
			containers.PUT("/:id", a.containerHandler.UpdateContainer)
			containers.DELETE("/:id", a.containerHandler.DeleteContainer)
		}

		products := v1.Group("/products", perm.Require("product:*"))
		{
			products.POST("", a.productHandler.CreateProduct)
			products.GET("", a.productHandler.ListProducts)
			products.GET("/:id", a.productHandler.GetProduct)
			products.PUT("/:id", a.productHandler.UpdateProduct)
			products.DELETE("/:id", a.productHandler.DeleteProduct)
		}

		plans := v1.Group("/plans")
		{
			plans.GET("", perm.Require("plan:read"), a.planHandler.ListPlans)
			plans.GET("/:id", perm.Require("plan:read"), a.planHandler.GetPlan)
			plans.GET("/:id/items/:itemId", perm.Require("plan_item:read"), a.planHandler.GetPlanItem)

			plans.POST("", perm.Require("plan:create"), a.planHandler.CreatePlan)
			plans.PUT("/:id", perm.Require("plan:update"), a.planHandler.UpdatePlan)
			plans.DELETE("/:id", perm.Require("plan:delete"), a.planHandler.DeletePlan)

			plans.POST("/:id/items", perm.Require("plan_item:*"), a.planHandler.AddPlanItem)
			plans.PUT("/:id/items/:itemId", perm.Require("plan_item:*"), a.planHandler.UpdatePlanItem)
			plans.DELETE("/:id/items/:itemId", perm.Require("plan_item:*"), a.planHandler.DeletePlanItem)

			plans.POST("/:id/calculate", perm.Require("plan:calculate"), a.planHandler.CalculatePlan)
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
