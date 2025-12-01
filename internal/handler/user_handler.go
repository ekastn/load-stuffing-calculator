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
