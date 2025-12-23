package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Creates a new user with the specified role. Requires admin privileges.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateUserRequest	true	"User Creation Data"
//	@Success		201		{object}	response.APIResponse{data=dto.UserResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.userSvc.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetUser godoc
//
//	@Summary		Get a user by ID
//	@Description	Retrieves user details by their ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	response.APIResponse{data=dto.UserResponse}
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	resp, err := h.userSvc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListUsers godoc
//
//	@Summary		List users
//	@Description	Retrieves a paginated list of users.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"		default(1)
//	@Param			limit	query		int	false	"Items per page"	default(10)
//	@Success		200		{object}	response.APIResponse{data=[]dto.UserResponse}
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.userSvc.ListUsers(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list users")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Updates an existing user. Requires admin privileges.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"User ID"
//	@Param			request	body		dto.UpdateUserRequest	true	"User Update Data"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.userSvc.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Deletes a user by ID. Requires admin privileges.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	err := h.userSvc.DeleteUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// ChangePassword godoc
//
//	@Summary		Change user password
//	@Description	Changes the password for a specific user. Requires admin privileges.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User ID"
//	@Param			request	body		dto.ChangePasswordRequest	true	"Password Change Data"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if req.Password != req.ConfirmPassword {
		response.Error(c, http.StatusBadRequest, "Password and Confirm Password do not match")
		return
	}

	err := h.userSvc.ChangePassword(c.Request.Context(), id, req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to change password: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}
