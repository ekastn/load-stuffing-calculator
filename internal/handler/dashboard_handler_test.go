package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDashboardHandler_GetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful_get_stats_admin", func(t *testing.T) {
		mockSvc := new(MockDashboardService)
		h := handler.NewDashboardHandler(mockSvc)

		expected := &dto.DashboardStatsResponse{
			Admin: &dto.AdminStats{
				TotalUsers:      10,
				ActiveShipments: 25,
				ContainerTypes:  5,
				SuccessRate:     95.5,
			},
		}

		mockSvc.On("GetStats", mock.Anything, "admin", mock.AnythingOfType("*uuid.UUID")).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/stats", nil)

		h.GetStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("successful_get_stats_planner", func(t *testing.T) {
		mockSvc := new(MockDashboardService)
		h := handler.NewDashboardHandler(mockSvc)

		expected := &dto.DashboardStatsResponse{
			Planner: &dto.PlannerStats{
				PendingPlans:   15,
				CompletedToday: 5,
				AvgUtilization: 85.5,
				ItemsProcessed: 1000,
			},
		}

		mockSvc.On("GetStats", mock.Anything, "planner", mock.AnythingOfType("*uuid.UUID")).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "planner")
		c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/stats", nil)

		h.GetStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("successful_get_stats_operator", func(t *testing.T) {
		mockSvc := new(MockDashboardService)
		h := handler.NewDashboardHandler(mockSvc)

		expected := &dto.DashboardStatsResponse{
			Operator: &dto.OperatorStats{
				ActiveLoads:       10,
				Completed:         50,
				FailedValidations: 2,
				AvgTimePerLoad:    "45m",
			},
		}

		mockSvc.On("GetStats", mock.Anything, "operator", mock.AnythingOfType("*uuid.UUID")).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "operator")
		c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/stats", nil)

		h.GetStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing_role_context", func(t *testing.T) {
		mockSvc := new(MockDashboardService)
		h := handler.NewDashboardHandler(mockSvc)

		expected := &dto.DashboardStatsResponse{}

		mockSvc.On("GetStats", mock.Anything, "", mock.AnythingOfType("*uuid.UUID")).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// No role set in context
		c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/stats", nil)

		h.GetStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockSvc := new(MockDashboardService)
		h := handler.NewDashboardHandler(mockSvc)

		mockSvc.On("GetStats", mock.Anything, "admin", mock.AnythingOfType("*uuid.UUID")).Return((*dto.DashboardStatsResponse)(nil), errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/stats", nil)

		h.GetStats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
