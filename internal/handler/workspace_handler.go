package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	workspaceSvc service.WorkspaceService
}

func NewWorkspaceHandler(workspaceSvc service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceSvc: workspaceSvc}
}

// ListWorkspaces godoc
//
//	@Summary		List workspaces
//	@Description	Retrieves a paginated list of workspaces the user is a member of.
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"		default(1)
//	@Param			limit	query		int	false	"Items per page"	default(10)
//	@Success		200		{object}	response.APIResponse{data=[]dto.WorkspaceResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/workspaces [get]
func (h *WorkspaceHandler) ListWorkspaces(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.workspaceSvc.ListWorkspaces(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list workspaces: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// CreateWorkspace godoc
//
//	@Summary		Create workspace
//	@Description	Creates an organization workspace; the creator becomes the owner member.
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateWorkspaceRequest	true	"Workspace creation data"
//	@Success		201		{object}	response.APIResponse{data=dto.WorkspaceResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/workspaces [post]
func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var req dto.CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.workspaceSvc.CreateWorkspace(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to create workspace: "+err.Error())
		return
	}
	response.Success(c, http.StatusCreated, resp)
}

// UpdateWorkspace godoc
//
//	@Summary		Update workspace
//	@Description	Renames a workspace and/or transfers ownership.
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Workspace ID"
//	@Param			request	body		dto.UpdateWorkspaceRequest	true	"Workspace update data"
//	@Success		200		{object}	response.APIResponse{data=dto.WorkspaceResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/workspaces/{id} [patch]
func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Workspace ID is required")
		return
	}

	var req dto.UpdateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.workspaceSvc.UpdateWorkspace(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update workspace: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// DeleteWorkspace godoc
//
//	@Summary		Delete workspace
//	@Description	Deletes an organization workspace and cascades workspace-owned records.
//	@Tags			workspaces
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Workspace ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/workspaces/{id} [delete]
func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Workspace ID is required")
		return
	}

	if err := h.workspaceSvc.DeleteWorkspace(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to delete workspace: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, nil)
}
