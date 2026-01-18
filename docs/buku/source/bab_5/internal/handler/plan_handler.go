package handler

import (
	"net/http"

	"load-stuffing-calculator/internal/dto"
	"load-stuffing-calculator/internal/response"
	"load-stuffing-calculator/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanHandler struct {
	svc *service.PlanService
}

func NewPlanHandler(svc *service.PlanService) *PlanHandler {
	return &PlanHandler{svc: svc}
}

func (h *PlanHandler) Create(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	containerID, _ := uuid.Parse(req.ContainerID)

	plan, err := h.svc.Create(c.Request.Context(), containerID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := dto.PlanResponse{
		ID:          plan.ID.String(),
		ContainerID: plan.ContainerID.String(),
		Status:      plan.Status,
	}

	response.Success(c, http.StatusCreated, resp)
}

func (h *PlanHandler) Calculate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	result, err := h.svc.Calculate(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}
