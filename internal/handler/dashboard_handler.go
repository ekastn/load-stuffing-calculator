package handler

import (
	"net/http"

	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardSvc service.DashboardService
}

func NewDashboardHandler(dashboardSvc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc}
}

// GetStats godoc
//
//	@Summary		Get Dashboard Stats
//	@Description	Returns aggregated statistics for the dashboard based on user role.
//	@Tags			dashboard
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.APIResponse{data=dto.DashboardStatsResponse}
//	@Failure		500	{object}	response.APIResponse
//	@Security		BearerAuth
//	@Router			/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		role = ""
	}

	stats, err := h.dashboardSvc.GetStats(c.Request.Context(), role.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch dashboard statistics")
		return
	}

	response.Success(c, http.StatusOK, stats)
}
