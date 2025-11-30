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
