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
			auth.POST("/register", a.authHandler.Register)
			auth.POST("/guest", a.authHandler.GuestToken)
			auth.POST("/refresh", a.authHandler.RefreshToken)
		}

		v1.Use(middleware.JWT(a.jwtSecret))

		// Auth endpoints that require a valid access token.
		authAuthed := v1.Group("/auth")
		{
			authAuthed.POST("/switch-workspace", a.authHandler.SwitchWorkspace)
			authAuthed.GET("/me", a.authHandler.Me)
		}

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

		containers := v1.Group("/containers")
		{
			containers.POST("", perm.Require("container:create"), a.containerHandler.CreateContainer)
			containers.GET("", perm.Require("container:read"), a.containerHandler.ListContainers)
			containers.GET("/:id", perm.Require("container:read"), a.containerHandler.GetContainer)
			containers.PUT("/:id", perm.Require("container:update"), a.containerHandler.UpdateContainer)
			containers.DELETE("/:id", perm.Require("container:delete"), a.containerHandler.DeleteContainer)
		}

		products := v1.Group("/products")
		{
			products.POST("", perm.Require("product:create"), a.productHandler.CreateProduct)
			products.GET("", perm.Require("product:read"), a.productHandler.ListProducts)
			products.GET("/:id", perm.Require("product:read"), a.productHandler.GetProduct)
			products.PUT("/:id", perm.Require("product:update"), a.productHandler.UpdateProduct)
			products.DELETE("/:id", perm.Require("product:delete"), a.productHandler.DeleteProduct)
		}

		workspaces := v1.Group("/workspaces")
		{
			workspaces.GET("", perm.Require("workspace:read"), a.workspaceHandler.ListWorkspaces)
			workspaces.POST("", perm.Require("workspace:create"), a.workspaceHandler.CreateWorkspace)
			workspaces.PATCH("/:id", perm.Require("workspace:update"), a.workspaceHandler.UpdateWorkspace)
			workspaces.DELETE("/:id", perm.Require("workspace:delete"), a.workspaceHandler.DeleteWorkspace)
		}

		members := v1.Group("/members")
		{
			members.GET("", perm.Require("member:read"), a.memberHandler.ListMembers)
			members.POST("", perm.Require("member:create"), a.memberHandler.AddMember)
			members.PATCH("/:member_id", perm.Require("member:update"), a.memberHandler.UpdateMemberRole)
			members.DELETE("/:member_id", perm.Require("member:delete"), a.memberHandler.DeleteMember)
		}

		invites := v1.Group("/invites")
		{
			invites.GET("", perm.Require("invite:read"), a.inviteHandler.ListInvites)
			invites.POST("", perm.Require("invite:create"), a.inviteHandler.CreateInvite)
			invites.DELETE("/:invite_id", perm.Require("invite:delete"), a.inviteHandler.RevokeInvite)
			// Note: accept is intentionally not permission-gated; the user may not yet be a member.
			invites.POST("/accept", a.inviteHandler.AcceptInvite)
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

			plans.GET("/:id/barcodes", perm.Require("plan:read"), a.planHandler.GetPlanBarcodes)
			plans.POST("/:id/validations", perm.Require("plan:read"), a.planHandler.ValidatePlanBarcode)
		}
	}
}

// HealthCheck godoc
//
//	@Summary		Health Check
//	@Description	Checks if the server is running and returns basic info.
//	@Tags			system
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (a *App) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0", // TODO: Get from build info
	})
}
