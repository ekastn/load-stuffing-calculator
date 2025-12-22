package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/cache"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleSvc   service.RoleService
	permCache *cache.PermissionCache
}

func NewRoleHandler(roleSvc service.RoleService, permCache *cache.PermissionCache) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc, permCache: permCache}
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Creates a new role. Requires admin privileges.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateRoleRequest true "Role Creation Data"
// @Success      201  {object}  response.APIResponse{data=dto.RoleResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.roleSvc.CreateRole(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create role: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetRole godoc
// @Summary      Get a role by ID
// @Description  Retrieves role details by ID.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object}  response.APIResponse{data=dto.RoleResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Role ID is required")
		return
	}

	resp, err := h.roleSvc.GetRole(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Role not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListRoles godoc
// @Summary      List roles
// @Description  Retrieves a paginated list of roles.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Success      200  {object}  response.APIResponse{data=[]dto.RoleResponse}
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.roleSvc.ListRoles(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list roles")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdateRole godoc
// @Summary      Update a role
// @Description  Updates an existing role. Requires admin privileges.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Role ID"
// @Param        request body      dto.UpdateRoleRequest  true  "Role Update Data"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Role ID is required")
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.roleSvc.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update role: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Deletes a role by ID. Requires admin privileges.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Role ID is required")
		return
	}

	err := h.roleSvc.DeleteRole(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete role: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// GetRolePermissions godoc
// @Summary      Get permissions for a role
// @Description  Retrieves a list of permission IDs assigned to a role.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Role ID"
// @Success      200  {object}  response.APIResponse{data=[]string}
// @Failure      400  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Role ID is required")
		return
	}

	resp, err := h.roleSvc.GetRolePermissions(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get role permissions")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdateRolePermissions godoc
// @Summary      Update permissions for a role
// @Description  Replaces all permissions for a role. Requires admin privileges.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        id      path      string                          true  "Role ID"
// @Param        request body      dto.UpdateRolePermissionsRequest  true  "Role Permissions Data"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /roles/{id}/permissions [put]
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Role ID is required")
		return
	}

	var req dto.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.roleSvc.UpdateRolePermissions(c.Request.Context(), id, req.PermissionIDs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update role permissions: "+err.Error())
		return
	}

	h.permCache.Invalidate()
	response.Success(c, http.StatusOK, nil)
}
