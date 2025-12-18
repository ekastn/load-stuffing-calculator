package api

import (
	"time"

	_ "github.com/ekastn/load-stuffing-calculator/internal/docs"
	"github.com/ekastn/load-stuffing-calculator/internal/middleware"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
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

		users := v1.Group("/users", middleware.Role(types.RoleAdmin))
		{
			users.POST("", a.userHandler.CreateUser)
			users.GET("/:id", a.userHandler.GetUser)
			users.GET("", a.userHandler.ListUsers)
			users.PUT("/:id", a.userHandler.UpdateUser)
			users.DELETE("/:id", a.userHandler.DeleteUser)
			users.PUT("/:id/password", a.userHandler.ChangePassword)
		}

		roles := v1.Group("/roles", middleware.Role(types.RoleAdmin))
		{
			roles.POST("", a.roleHandler.CreateRole)
			roles.GET("", a.roleHandler.ListRoles)
			roles.GET("/:id", a.roleHandler.GetRole)
			roles.PUT("/:id", a.roleHandler.UpdateRole)
			roles.DELETE("/:id", a.roleHandler.DeleteRole)
		}

		permissions := v1.Group("/permissions", middleware.Role(types.RoleAdmin))
		{
			permissions.POST("", a.permHandler.CreatePermission)
			permissions.GET("", a.permHandler.ListPermissions)
			permissions.GET("/:id", a.permHandler.GetPermission)
			permissions.PUT("/:id", a.permHandler.UpdatePermission)
			permissions.DELETE("/:id", a.permHandler.DeletePermission)
		}

		containers := v1.Group("/containers", middleware.Role(types.RoleAdmin))
		{
			containers.POST("", a.containerHandler.CreateContainer)
			containers.GET("", a.containerHandler.ListContainers)
			containers.GET("/:id", a.containerHandler.GetContainer)
			containers.PUT("/:id", a.containerHandler.UpdateContainer)
			containers.DELETE("/:id", a.containerHandler.DeleteContainer)
		}

		products := v1.Group("/products", middleware.Role(types.RoleAdmin))
		{
			products.POST("", a.productHandler.CreateProduct)
			products.GET("", a.productHandler.ListProducts)
			products.GET("/:id", a.productHandler.GetProduct)
			products.PUT("/:id", a.productHandler.UpdateProduct)
			products.DELETE("/:id", a.productHandler.DeleteProduct)
		}

		plans := v1.Group("/plans")
		{
			// Read access: Planner + Operator (Admin implicit)
			plans.GET("", middleware.Role(types.RolePlanner, types.RoleOperator), a.planHandler.ListPlans)
			plans.GET("/:id", middleware.Role(types.RolePlanner, types.RoleOperator), a.planHandler.GetPlan)
			plans.GET("/:id/items/:itemId", middleware.Role(types.RolePlanner, types.RoleOperator), a.planHandler.GetPlanItem)

			// Write access: Planner only (Admin implicit)
			plans.POST("", middleware.Role(types.RolePlanner), a.planHandler.CreatePlan)
			plans.PUT("/:id", middleware.Role(types.RolePlanner), a.planHandler.UpdatePlan)
			plans.DELETE("/:id", middleware.Role(types.RolePlanner), a.planHandler.DeletePlan)

			plans.POST("/:id/items", middleware.Role(types.RolePlanner), a.planHandler.AddPlanItem)
			plans.PUT("/:id/items/:itemId", middleware.Role(types.RolePlanner), a.planHandler.UpdatePlanItem)
			plans.DELETE("/:id/items/:itemId", middleware.Role(types.RolePlanner), a.planHandler.DeletePlanItem)

			plans.POST("/:id/calculate", middleware.Role(types.RolePlanner, types.RoleOperator), a.planHandler.CalculatePlan)
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
