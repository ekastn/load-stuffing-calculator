package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanHandler struct {
	planSvc service.PlanService
}

func withFounderWorkspaceOverride(c *gin.Context) {
	workspaceID := c.Query("workspace_id")
	if workspaceID == "" {
		return
	}
	role, ok := auth.RoleFromContext(c.Request.Context())
	if !ok || role != types.RoleFounder.String() {
		return
	}
	ctx := auth.WithWorkspaceOverrideID(c.Request.Context(), workspaceID)
	c.Request = c.Request.WithContext(ctx)
}

func NewPlanHandler(planSvc service.PlanService) *PlanHandler {
	return &PlanHandler{planSvc: planSvc}
}

func respondPlanServiceError(c *gin.Context, err error, defaultStatus int, defaultMessage string) {
	switch {
	case errors.Is(err, service.ErrTrialLimitReached):
		response.Error(c, http.StatusTooManyRequests, "Trial limit reached")
	case errors.Is(err, service.ErrForbidden):
		response.Error(c, http.StatusForbidden, "Forbidden")
	default:
		response.Error(c, defaultStatus, defaultMessage+err.Error())
	}
}

// CreatePlan godoc
//
//	@Summary		Create a new load plan
//	@Description	Creates a new load plan with a container and items.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string					false	"Workspace override (founder only)"
//	@Param			request			body		dto.CreatePlanRequest	true	"Plan Creation Data"
//	@Success		201				{object}	response.APIResponse{data=dto.CreatePlanResponse}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans [post]
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	withFounderWorkspaceOverride(c)

	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.planSvc.CreateCompletePlan(c.Request.Context(), req)
	if err != nil {
		respondPlanServiceError(c, err, http.StatusInternalServerError, "Failed to create plan: ")
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetPlan godoc
//
//	@Summary		Get a plan by ID
//	@Description	Retrieves plan details including items and stats.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			id				path		string	true	"Plan ID"
//	@Success		200				{object}	response.APIResponse{data=dto.PlanDetailResponse}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		404				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id} [get]
func (h *PlanHandler) GetPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	withFounderWorkspaceOverride(c)

	resp, err := h.planSvc.GetPlan(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Plan not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListPlans godoc
//
//	@Summary		List plans
//	@Description	Retrieves a paginated list of plans.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			page			query		int		false	"Page number"		default(1)
//	@Param			limit			query		int		false	"Items per page"	default(10)
//	@Success		200				{object}	response.APIResponse{data=[]dto.PlanListItem}
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans [get]
func (h *PlanHandler) ListPlans(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	withFounderWorkspaceOverride(c)

	resp, err := h.planSvc.ListPlans(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list plans")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdatePlan godoc
//
//	@Summary		Update a plan
//	@Description	Updates an existing plan (status, container).
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string					false	"Workspace override (founder only)"
//	@Param			id				path		string					true	"Plan ID"
//	@Param			request			body		dto.UpdatePlanRequest	true	"Plan Update Data"
//	@Success		200				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id} [put]
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

	withFounderWorkspaceOverride(c)

	if err := h.planSvc.UpdatePlan(c.Request.Context(), id, req); err != nil {
		respondPlanServiceError(c, err, http.StatusInternalServerError, "Failed to update plan: ")
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeletePlan godoc
//
//	@Summary		Delete a plan
//	@Description	Deletes an existing plan.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			id				path		string	true	"Plan ID"
//	@Success		200				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id} [delete]
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	withFounderWorkspaceOverride(c)

	err := h.planSvc.DeletePlan(c.Request.Context(), id)
	if err != nil {
		respondPlanServiceError(c, err, http.StatusInternalServerError, "Failed to delete plan: ")
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// AddPlanItem godoc
//
//	@Summary		Add item to plan
//	@Description	Adds a new item to an existing plan.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string					false	"Workspace override (founder only)"
//	@Param			id				path		string					true	"Plan ID"
//	@Param			request			body		dto.AddPlanItemRequest	true	"Item Data"
//	@Success		201				{object}	response.APIResponse{data=dto.PlanItemDetail}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id}/items [post]
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

	withFounderWorkspaceOverride(c)

	resp, err := h.planSvc.AddPlanItem(c.Request.Context(), id, req)
	if err != nil {
		respondPlanServiceError(c, err, http.StatusInternalServerError, "Failed to add item: ")
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetPlanItem godoc
//
//	@Summary		Get plan item
//	@Description	Retrieves details of a specific item in a plan.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			id				path		string	true	"Plan ID"
//	@Param			itemId			path		string	true	"Item ID"
//	@Success		200				{object}	response.APIResponse{data=dto.PlanItemDetail}
//	@Failure		404				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id}/items/{itemId} [get]
func (h *PlanHandler) GetPlanItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")
	if id == "" || itemId == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID and Item ID are required")
		return
	}

	withFounderWorkspaceOverride(c)

	resp, err := h.planSvc.GetPlanItem(c.Request.Context(), id, itemId)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Item not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdatePlanItem godoc
//
//	@Summary		Update plan item
//	@Description	Updates a specific item in a plan.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string						false	"Workspace override (founder only)"
//	@Param			id				path		string						true	"Plan ID"
//	@Param			itemId			path		string						true	"Item ID"
//	@Param			request			body		dto.UpdatePlanItemRequest	true	"Update Data"
//	@Success		200				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id}/items/{itemId} [put]
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

	withFounderWorkspaceOverride(c)

	err := h.planSvc.UpdatePlanItem(c.Request.Context(), id, itemId, req)
	if err != nil {
		respondPlanServiceError(c, err, http.StatusInternalServerError, "Failed to update item: ")
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeletePlanItem godoc
//
//	@Summary		Delete plan item
//	@Description	Deletes a specific item from a plan.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			id				path		string	true	"Plan ID"
//	@Param			itemId			path		string	true	"Item ID"
//	@Success		200				{object}	response.APIResponse
//	@Failure		400				{object}	response.APIResponse
//	@Failure		500				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id}/items/{itemId} [delete]
func (h *PlanHandler) DeletePlanItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")
	if id == "" || itemId == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID and Item ID are required")
		return
	}

	withFounderWorkspaceOverride(c)

	err := h.planSvc.DeletePlanItem(c.Request.Context(), id, itemId)
	if err != nil {
		respondPlanServiceError(c, err, http.StatusInternalServerError, "Failed to delete item: ")
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// CalculatePlan godoc
//
//	@Summary		Calculate plan
//	@Description	Triggers the packing calculation for a plan.
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query	string						false	"Workspace override (founder only)"
//	@Param			id				path	string						true	"Plan ID"
//	@Param			request			body	dto.CalculatePlanRequest	false	"Calculation Options"
//
// @Success	200	{object}	response.APIResponse{data=dto.CalculationResult}
// @Failure	400	{object}	response.APIResponse
// @Failure	500	{object}	response.APIResponse
// @Security	BearerAuth
// @Router		/plans/{id}/calculate [post]
func (h *PlanHandler) CalculatePlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Plan ID is required")
		return
	}

	var req dto.CalculatePlanRequest
	// Empty body is allowed (defaults will apply).
	if err := c.ShouldBindJSON(&req); err != nil {
		// Gin returns io.EOF for empty request bodies; treat that as "no options".
		if !errors.Is(err, io.EOF) {
			response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
			return
		}
		req = dto.CalculatePlanRequest{}
	}

	withFounderWorkspaceOverride(c)

	resp, err := h.planSvc.CalculatePlan(c.Request.Context(), id, req)
	if err != nil {
		respondPlanServiceError(c, err, http.StatusInternalServerError, "Failed to calculate plan: ")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// GetPlanBarcodes returns generated barcodes for all placements in a plan
//
//	@Summary		Get plan barcodes
//	@Description	Returns generated barcodes for all placements in a plan
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string	false	"Workspace override (founder only)"
//	@Param			id				path		string	true	"Plan ID"
//	@Success		200				{object}	response.APIResponse{data=[]dto.BarcodeInfo}
//	@Failure		400				{object}	response.APIResponse
//	@Failure		404				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id}/barcodes [get]
func (h *PlanHandler) GetPlanBarcodes(c *gin.Context) {
	planID := c.Param("id")

	// Get plan with placements (reuse existing method)
	plan, err := h.planSvc.GetPlan(c.Request.Context(), planID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Plan not found")
		return
	}

	planUUID, err := uuid.Parse(plan.PlanID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Invalid plan ID in response")
		return
	}

	// Generate barcodes
	barcodes := make([]dto.BarcodeInfo, 0, len(plan.Calculation.Placements))

	for _, placement := range plan.Calculation.Placements {
		// Skip placements without step numbers (shouldn't happen for valid placements)
		if placement.StepNumber == 0 {
			continue
		}

		item := findItemByID(plan.Items, placement.ItemID)
		if item == nil {
			continue
		}

		itemUUID, err := uuid.Parse(placement.ItemID)
		if err != nil {
			continue
		}

		barcode := generateBarcode(planUUID, placement.StepNumber, itemUUID)

		barcodes = append(barcodes, dto.BarcodeInfo{
			StepNumber: placement.StepNumber,
			ItemID:     placement.ItemID,
			ItemLabel:  *item.Label,
			Barcode:    barcode,
			Position: dto.Position{
				X: placement.PositionX,
				Y: placement.PositionY,
				Z: placement.PositionZ,
			},
			Dimensions: dto.Dimensions{
				Length: item.LengthMM,
				Width:  item.WidthMM,
				Height: item.HeightMM,
			},
		})
	}

	// Sort by step number
	sort.Slice(barcodes, func(i, j int) bool {
		return barcodes[i].StepNumber < barcodes[j].StepNumber
	})

	response.Success(c, http.StatusOK, barcodes)
}

// ValidatePlanBarcode validates a scanned barcode against a plan
//
//	@Summary		Validate plan barcode
//	@Description	Validates a scanned barcode against a plan
//	@Tags			plans
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	query		string						false	"Workspace override (founder only)"
//	@Param			id				path		string						true	"Plan ID"
//	@Param			request			body		dto.ValidateBarcodeRequest	true	"Validation Data"
//	@Success		200				{object}	response.APIResponse{data=dto.ValidationResult}
//	@Failure		400				{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/plans/{id}/validations [post]
func (h *PlanHandler) ValidatePlanBarcode(c *gin.Context) {
	planID := c.Param("id")

	var req dto.ValidateBarcodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse scanned barcode
	parsed := parseBarcode(req.Barcode)
	if parsed == nil {
		response.Success(c, http.StatusOK, dto.ValidationResult{
			Valid:  false,
			Status: "INVALID_FORMAT",
			Error:  "Invalid barcode format",
		})
		return
	}

	// Verify barcode belongs to this plan
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	planShort := planUUID.String()[:8]
	if parsed.PlanID != planShort {
		response.Success(c, http.StatusOK, dto.ValidationResult{
			Valid:      false,
			Status:     "WRONG_PLAN",
			Error:      "Barcode belongs to different plan",
			PlanID:     parsed.PlanID,
			StepNumber: parsed.StepNumber,
			ItemID:     parsed.ItemID,
		})
		return
	}

	// Check if step matches expected
	valid := false
	if req.ExpectedStep != nil {
		valid = parsed.StepNumber == *req.ExpectedStep
	} else {
		// If no expected step provided, we just validate format and plan match
		// Ideally we should check if step exists in plan but that requires fetching plan
		valid = true
	}

	status := "MATCHED"
	if !valid {
		if req.ExpectedStep != nil {
			status = "OUT_OF_SEQUENCE"
		} else {
			status = "UNEXPECTED_STEP"
		}
	}

	response.Success(c, http.StatusOK, dto.ValidationResult{
		Valid:      valid,
		Status:     status,
		PlanID:     parsed.PlanID,
		StepNumber: parsed.StepNumber,
		ItemID:     parsed.ItemID,
		Barcode:    req.Barcode,
	})
}

// Helper: Generate deterministic barcode
func generateBarcode(planID uuid.UUID, stepNumber int, itemID uuid.UUID) string {
	planShort := planID.String()[:8]
	itemShort := itemID.String()[:8]
	return fmt.Sprintf("PLAN-%s-STEP-%03d-%s", planShort, stepNumber, itemShort)
}

// Helper: Parse scanned barcode
func parseBarcode(barcode string) *ParsedBarcode {
	parts := strings.Split(barcode, "-")
	if len(parts) != 5 || parts[0] != "PLAN" || parts[2] != "STEP" {
		return nil
	}

	stepNum, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil
	}

	return &ParsedBarcode{
		PlanID:     parts[1],
		StepNumber: stepNum,
		ItemID:     parts[4],
	}
}

// Helper: Find item by ID
func findItemByID(items []dto.PlanItemDetail, itemID string) *dto.PlanItemDetail {
	for i := range items {
		if items[i].ItemID == itemID {
			return &items[i]
		}
	}
	return nil
}

// Helper struct
type ParsedBarcode struct {
	PlanID     string
	StepNumber int
	ItemID     string
}
