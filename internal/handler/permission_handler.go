package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permSvc service.PermissionService
}

func NewPermissionHandler(permSvc service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permSvc: permSvc}
}

// CreatePermission godoc
// @Summary      Create a new permission
// @Description  Creates a new permission. Requires admin privileges.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePermissionRequest true "Permission Creation Data"
// @Success      201  {object}  response.APIResponse{data=dto.PermissionResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.permSvc.CreatePermission(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create permission: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetPermission godoc
// @Summary      Get a permission by ID
// @Description  Retrieves permission details by ID.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Permission ID"
// @Success      200  {object}  response.APIResponse{data=dto.PermissionResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Permission ID is required")
		return
	}

	resp, err := h.permSvc.GetPermission(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Permission not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListPermissions godoc
// @Summary      List permissions
// @Description  Retrieves a paginated list of permissions.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Success      200  {object}  response.APIResponse{data=[]dto.PermissionResponse}
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.permSvc.ListPermissions(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list permissions")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdatePermission godoc
// @Summary      Update a permission
// @Description  Updates an existing permission. Requires admin privileges.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Permission ID"
// @Param        request body      dto.UpdatePermissionRequest  true  "Permission Update Data"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Permission ID is required")
		return
	}

	var req dto.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.permSvc.UpdatePermission(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update permission: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeletePermission godoc
// @Summary      Delete a permission
// @Description  Deletes a permission by ID. Requires admin privileges.
// @Tags         permissions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Permission ID"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Permission ID is required")
		return
	}

	err := h.permSvc.DeletePermission(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete permission: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}
