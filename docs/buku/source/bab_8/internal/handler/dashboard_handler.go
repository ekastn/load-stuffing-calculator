package handler

import (
	"net/http"

	"load-stuffing-calculator/internal/response"
	"load-stuffing-calculator/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format in context")
		return
	}

	stats, err := h.svc.GetStats(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch dashboard stats")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"total_plans":         stats.TotalPlans,
		"total_containers":    stats.TotalContainers,
		"total_products":      stats.TotalProducts,
		"total_items_shipped": stats.TotalItemsShipped,
	})
}
