package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	planSvc service.PlanService
}

func NewPlanHandler(planSvc service.PlanService) *PlanHandler {
	return &PlanHandler{planSvc: planSvc}
}

// CreatePlan godoc
// @Summary      Create a new load plan
// @Description  Creates a new load plan with a container and items.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePlanRequest true "Plan Creation Data"
// @Success      201  {object}  response.APIResponse{data=dto.CreatePlanResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans [post]
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.planSvc.CreateCompletePlan(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create plan: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetPlan godoc
// @Summary      Get a plan by ID
// @Description  Retrieves plan details including items and stats.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  response.APIResponse{data=dto.PlanDetailResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id} [get]
func (h *PlanHandler) GetPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	resp, err := h.planSvc.GetPlan(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Plan not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListPlans godoc
// @Summary      List plans
// @Description  Retrieves a paginated list of plans.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Success      200  {object}  response.APIResponse{data=[]dto.PlanListItem}
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans [get]
func (h *PlanHandler) ListPlans(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.planSvc.ListPlans(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list plans")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdatePlan godoc
// @Summary      Update a plan
// @Description  Updates an existing plan (status, container).
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Plan ID"
// @Param        request body      dto.UpdatePlanRequest  true  "Plan Update Data"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id} [put]
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.planSvc.UpdatePlan(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update plan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeletePlan godoc
// @Summary      Delete a plan
// @Description  Deletes a plan by ID.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id} [delete]
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	err := h.planSvc.DeletePlan(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete plan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// AddPlanItem godoc
// @Summary      Add item to plan
// @Description  Adds a new item to an existing plan.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Plan ID"
// @Param        request body      dto.AddPlanItemRequest true  "Item Data"
// @Success      201     {object}  response.APIResponse{data=dto.PlanItemDetail}
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id}/items [post]
func (h *PlanHandler) AddPlanItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	var req dto.AddPlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.planSvc.AddPlanItem(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to add item: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetPlanItem godoc
// @Summary      Get plan item
// @Description  Retrieves details of a specific item in a plan.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Plan ID"
// @Param        itemId  path      string  true  "Item ID"
// @Success      200     {object}  response.APIResponse{data=dto.PlanItemDetail}
// @Failure      404     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id}/items/{itemId} [get]
func (h *PlanHandler) GetPlanItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")
	if id == "" || itemId == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID and Item ID are required")
		return
	}

	resp, err := h.planSvc.GetPlanItem(c.Request.Context(), id, itemId)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Item not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdatePlanItem godoc
// @Summary      Update plan item
// @Description  Updates a specific item in a plan.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Plan ID"
// @Param        itemId  path      string                 true  "Item ID"
// @Param        request body      dto.UpdatePlanItemRequest true "Update Data"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id}/items/{itemId} [put]
func (h *PlanHandler) UpdatePlanItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")
	if id == "" || itemId == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID and Item ID are required")
		return
	}

	var req dto.UpdatePlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.planSvc.UpdatePlanItem(c.Request.Context(), id, itemId, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update item: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeletePlanItem godoc
// @Summary      Delete plan item
// @Description  Deletes a specific item from a plan.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Plan ID"
// @Param        itemId  path      string  true  "Item ID"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id}/items/{itemId} [delete]
func (h *PlanHandler) DeletePlanItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")
	if id == "" || itemId == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID and Item ID are required")
		return
	}

	err := h.planSvc.DeletePlanItem(c.Request.Context(), id, itemId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete item: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// CalculatePlan godoc
// @Summary      Calculate plan
// @Description  Triggers the packing calculation for a plan.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Plan ID"
// @Success      200  {object}  response.APIResponse{data=dto.CalculationResult}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /plans/{id}/calculate [post]
func (h *PlanHandler) CalculatePlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	resp, err := h.planSvc.CalculatePlan(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to calculate plan: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, resp)
}
