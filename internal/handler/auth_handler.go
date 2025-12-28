package handler

import (
	"net/http"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc}
}

// Login godoc
//
//	@Summary		User Login
//	@Description	Authenticates a user and returns access and refresh tokens.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login Credentials"
//	@Success		200		{object}	response.APIResponse{data=dto.LoginResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.authSvc.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// Register godoc
//
//	@Summary		Register User
//	@Description	Creates a new user account. If `guest_token` is provided, guest-created plans are claimed into the new account.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"Registration data"
//	@Success		201		{object}	response.APIResponse{data=dto.RegisterResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.authSvc.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to register user: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GuestToken godoc
//
//	@Summary		Issue Guest Token
//	@Description	Issues a trial (guest) access token. Trial users can manage only their own plans and are limited to 3 total plans.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse{data=dto.GuestTokenResponse}
//	@Failure		500	{object}	response.APIResponse
//	@Router			/auth/guest [post]
func (h *AuthHandler) GuestToken(c *gin.Context) {
	resp, err := h.authSvc.GuestToken(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to issue guest token")
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// RefreshToken godoc
//
//	@Summary		Refresh Access Token
//	@Description	Rotates the refresh token and issues a new access token.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh Token"
//	@Success		200		{object}	response.APIResponse{data=dto.LoginResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// SwitchWorkspace godoc
//
//	@Summary		Switch active workspace
//	@Description	Updates the active workspace for the supplied refresh token and returns a new access token.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.SwitchWorkspaceRequest	true	"Switch workspace request"
//	@Success		200		{object}	response.APIResponse{data=dto.SwitchWorkspaceResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/auth/switch-workspace [post]
func (h *AuthHandler) SwitchWorkspace(c *gin.Context) {
	var req dto.SwitchWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.authSvc.SwitchWorkspace(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to switch workspace: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// Me godoc
//
//	@Summary		Get current session
//	@Description	Returns the current user, active workspace, permissions and platform membership.
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	response.APIResponse{data=dto.AuthMeResponse}
//	@Failure		401	{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	resp, err := h.authSvc.Me(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Failed to resolve session")
		return
	}

	response.Success(c, http.StatusOK, resp)
}
