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

func (h *PlanHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid plan ID format")
		return
	}

	plan, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Plan not found")
		return
	}

	// Get items for this plan
	items, err := h.svc.GetItems(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get plan items")
		return
	}

	// Get placements if plan is completed
	var placements []dto.PlacementResponse
	if plan.Status == "completed" {
		pl, err := h.svc.GetPlacements(c.Request.Context(), id)
		if err == nil {
			for _, p := range pl {
				placements = append(placements, dto.PlacementResponse{
					ID:         p.ID.String(),
					ProductID:  p.ProductID.String(),
					Label:      p.Label,
					PosX:       p.PosX,
					PosY:       p.PosY,
					PosZ:       p.PosZ,
					Rotation:   int(p.Rotation),
					StepNumber: int(p.StepNumber),
				})
			}
		}
	}

	// Convert items to response
	itemResponses := make([]dto.PlanItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, dto.PlanItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Label:     item.Label,
			Quantity:  int(item.Quantity),
			LengthMm:  item.LengthMm,
			WidthMm:   item.WidthMm,
			HeightMm:  item.HeightMm,
			WeightKg:  item.WeightKg,
		})
	}

	resp := dto.PlanDetailResponse{
		ID:          plan.ID.String(),
		ContainerID: plan.ContainerID.String(),
		Status:      plan.Status,
		Items:       itemResponses,
		Placements:  placements,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *PlanHandler) List(c *gin.Context) {
	plans, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list plans")
		return
	}

	resp := make([]dto.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		resp = append(resp, dto.PlanResponse{
			ID:            plan.ID.String(),
			ContainerID:   plan.ContainerID.String(),
			ContainerName: plan.ContainerName,
			Status:        plan.Status,
		})
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *PlanHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid plan ID format")
		return
	}

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	containerID, _ := uuid.Parse(req.ContainerID)

	plan, err := h.svc.Update(c.Request.Context(), id, containerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.PlanResponse{
		ID:          plan.ID.String(),
		ContainerID: plan.ContainerID.String(),
		Status:      plan.Status,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *PlanHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid plan ID format")
		return
	}

	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

func (h *PlanHandler) AddItem(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid plan ID format")
		return
	}

	var req dto.AddPlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	productID, _ := uuid.Parse(req.ProductID)

	item, err := h.svc.AddItem(c.Request.Context(), planID, productID, req.Quantity)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp := gin.H{
		"id":         item.ID.String(),
		"plan_id":    item.PlanID.String(),
		"product_id": item.ProductID.String(),
		"quantity":   item.Quantity,
	}

	response.Success(c, http.StatusCreated, resp)
}

func (h *PlanHandler) UpdateItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid item ID format")
		return
	}

	var req dto.UpdatePlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.svc.UpdateItem(c.Request.Context(), itemID, req.Quantity)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := gin.H{
		"id":         item.ID.String(),
		"plan_id":    item.PlanID.String(),
		"product_id": item.ProductID.String(),
		"quantity":   item.Quantity,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *PlanHandler) DeleteItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid item ID format")
		return
	}

	err = h.svc.DeleteItem(c.Request.Context(), itemID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
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
