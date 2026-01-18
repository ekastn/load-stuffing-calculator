package handler

import (
	"net/http"

	"load-stuffing-calculator/internal/dto"
	"load-stuffing-calculator/internal/response"
	"load-stuffing-calculator/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContainerHandler struct {
	svc *service.ContainerService
}

func NewContainerHandler(svc *service.ContainerService) *ContainerHandler {
	return &ContainerHandler{svc: svc}
}

func (h *ContainerHandler) Create(c *gin.Context) {
	var req dto.CreateContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	container, err := h.svc.Create(
		c.Request.Context(),
		req.Name,
		req.LengthMm,
		req.WidthMm,
		req.HeightMm,
		req.MaxWeightKg,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert ke DTO response
	resp := dto.ContainerResponse{
		ID:          container.ID.String(),
		Name:        container.Name,
		LengthMm:    container.LengthMm,
		WidthMm:     container.WidthMm,
		HeightMm:    container.HeightMm,
		MaxWeightKg: container.MaxWeightKg,
	}

	response.Success(c, http.StatusCreated, resp)
}

func (h *ContainerHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "Container ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid container ID format")
		return
	}

	container, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Container not found")
		return
	}

	resp := dto.ContainerResponse{
		ID:          container.ID.String(),
		Name:        container.Name,
		LengthMm:    container.LengthMm,
		WidthMm:     container.WidthMm,
		HeightMm:    container.HeightMm,
		MaxWeightKg: container.MaxWeightKg,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *ContainerHandler) List(c *gin.Context) {
	containers, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list containers")
		return
	}

	// Convert slice ke DTO responses
	resp := make([]dto.ContainerResponse, 0, len(containers))
	for _, container := range containers {
		resp = append(resp, dto.ContainerResponse{
			ID:          container.ID.String(),
			Name:        container.Name,
			LengthMm:    container.LengthMm,
			WidthMm:     container.WidthMm,
			HeightMm:    container.HeightMm,
			MaxWeightKg: container.MaxWeightKg,
		})
	}

	response.Success(c, http.StatusOK, resp)
}
