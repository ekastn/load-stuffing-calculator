package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type InviteHandler struct {
	inviteSvc service.InviteService
}

func NewInviteHandler(inviteSvc service.InviteService) *InviteHandler {
	return &InviteHandler{inviteSvc: inviteSvc}
}

// ListInvites godoc
//
//	@Summary		List invites
//	@Description	Lists invites for the active workspace (founder may override via workspace_id).
//	@Tags			invites
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			page			query		int		false	"Page number"		default(1)
//	@Param			limit			query		int		false	"Items per page"	default(10)
//	@Success		200				{object}	response.APIResponse{data=[]dto.InviteResponse}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/invites [get]
func (h *InviteHandler) ListInvites(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	overrideWorkspaceID := c.Query("workspace_id")
	var override *string
	if overrideWorkspaceID != "" {
		override = &overrideWorkspaceID
	}

	resp, err := h.inviteSvc.ListInvites(c.Request.Context(), int32(page), int32(limit), override)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to list invites: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// CreateInvite godoc
//
//	@Summary		Create invite
//	@Description	Creates an invite for the active workspace and returns a raw token (shown only once).
//	@Tags			invites
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string					false	"Workspace override (founder only)"
//	@Param			request			body		dto.CreateInviteRequest	true	"Invite creation data"
//	@Success		201				{object}	response.APIResponse{data=dto.CreateInviteResponse}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/invites [post]
func (h *InviteHandler) CreateInvite(c *gin.Context) {
	var req dto.CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	overrideWorkspaceID := c.Query("workspace_id")
	var override *string
	if overrideWorkspaceID != "" {
		override = &overrideWorkspaceID
	}

	resp, err := h.inviteSvc.CreateInvite(c.Request.Context(), req, override)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to create invite: "+err.Error())
		return
	}
	response.Success(c, http.StatusCreated, resp)
}

// RevokeInvite godoc
//
//	@Summary		Revoke invite
//	@Description	Revokes a pending invite by ID.
//	@Tags			invites
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			invite_id		path		string	true	"Invite ID"
//	@Success		200				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/invites/{invite_id} [delete]
func (h *InviteHandler) RevokeInvite(c *gin.Context) {
	inviteID := c.Param("invite_id")
	if inviteID == "" {
		response.Error(c, http.StatusBadRequest, "Invite ID is required")
		return
	}

	overrideWorkspaceID := c.Query("workspace_id")
	var override *string
	if overrideWorkspaceID != "" {
		override = &overrideWorkspaceID
	}

	if err := h.inviteSvc.RevokeInvite(c.Request.Context(), inviteID, override); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to revoke invite: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, nil)
}

// AcceptInvite godoc
//
//	@Summary		Accept invite
//	@Description	Accepts an invite using a token and returns a new access token for the invite workspace.
//	@Tags			invites
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.AcceptInviteRequest	true	"Invite accept data"
//	@Success		200		{object}	response.APIResponse{data=dto.AcceptInviteResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/invites/accept [post]
func (h *InviteHandler) AcceptInvite(c *gin.Context) {
	var req dto.AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.inviteSvc.AcceptInvite(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to accept invite: "+err.Error())
		return
	}
	response.Success(c, http.StatusOK, resp)
}
