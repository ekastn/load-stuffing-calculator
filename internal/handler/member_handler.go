package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberSvc service.MemberService
}

func NewMemberHandler(memberSvc service.MemberService) *MemberHandler {
	return &MemberHandler{memberSvc: memberSvc}
}

// ListMembers godoc
//
//	@Summary		List members
//	@Description	Lists members for the active workspace (founder may override via workspace_id).
//	@Tags			members
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			page			query		int		false	"Page number"		default(1)
//	@Param			limit			query		int		false	"Items per page"	default(10)
//	@Success		200				{object}	response.APIResponse{data=[]dto.MemberResponse}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/members [get]
func (h *MemberHandler) ListMembers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	overrideWorkspaceID := c.Query("workspace_id")
	var override *string
	if overrideWorkspaceID != "" {
		override = &overrideWorkspaceID
	}

	resp, err := h.memberSvc.ListMembers(c.Request.Context(), int32(page), int32(limit), override)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to list members: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// AddMember godoc
//
//	@Summary		Add member
//	@Description	Adds an existing user to the active workspace (founder may override via workspace_id).
//	@Tags			members
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string					false	"Workspace override (founder only)"
//	@Param			request			body		dto.AddMemberRequest	true	"Member creation data"
//	@Success		201				{object}	response.APIResponse{data=dto.MemberResponse}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/members [post]
func (h *MemberHandler) AddMember(c *gin.Context) {
	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	overrideWorkspaceID := c.Query("workspace_id")
	var override *string
	if overrideWorkspaceID != "" {
		override = &overrideWorkspaceID
	}

	resp, err := h.memberSvc.AddMember(c.Request.Context(), req, override)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to add member: "+err.Error())
		return
	}
	response.Success(c, http.StatusCreated, resp)
}

// UpdateMemberRole godoc
//
//	@Summary		Update member role
//	@Description	Updates a member role in the active workspace (founder may override via workspace_id).
//	@Tags			members
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string						false	"Workspace override (founder only)"
//	@Param			member_id		path		string						true	"Member ID"
//	@Param			request			body		dto.UpdateMemberRoleRequest	true	"Member update data"
//	@Success		200				{object}	response.APIResponse{data=dto.MemberResponse}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/members/{member_id} [patch]
func (h *MemberHandler) UpdateMemberRole(c *gin.Context) {
	memberID := c.Param("member_id")
	if memberID == "" {
		response.Error(c, http.StatusBadRequest, "Member ID is required")
		return
	}

	var req dto.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	overrideWorkspaceID := c.Query("workspace_id")
	var override *string
	if overrideWorkspaceID != "" {
		override = &overrideWorkspaceID
	}

	resp, err := h.memberSvc.UpdateMemberRole(c.Request.Context(), memberID, req, override)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update member: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// DeleteMember godoc
//
//	@Summary		Delete member
//	@Description	Removes a member from the active workspace (founder may override via workspace_id). Owner membership cannot be removed.
//	@Tags			members
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			member_id		path		string	true	"Member ID"
//	@Success		200				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/members/{member_id} [delete]
func (h *MemberHandler) DeleteMember(c *gin.Context) {
	memberID := c.Param("member_id")
	if memberID == "" {
		response.Error(c, http.StatusBadRequest, "Member ID is required")
		return
	}

	overrideWorkspaceID := c.Query("workspace_id")
	var override *string
	if overrideWorkspaceID != "" {
		override = &overrideWorkspaceID
	}

	if err := h.memberSvc.DeleteMember(c.Request.Context(), memberID, override); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to delete member: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, nil)
}
