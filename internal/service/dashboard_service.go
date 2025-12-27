package service

import (
	"context"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
)

type DashboardService interface {
	GetStats(ctx context.Context, role string) (*dto.DashboardStatsResponse, error)
}

type dashboardService struct {
	q store.Querier
}

func NewDashboardService(q store.Querier) DashboardService {
	return &dashboardService{q: q}
}

func (s *dashboardService) GetStats(ctx context.Context, role string) (*dto.DashboardStatsResponse, error) {
	resp := &dto.DashboardStatsResponse{}

	if role == types.RoleAdmin.String() {
		totalUsers, _ := s.q.CountTotalUsers(ctx)
		activePlans, _ := s.q.CountActivePlans(ctx)
		containerTypes, _ := s.q.CountContainers(ctx)
		resp.Admin = &dto.AdminStats{
			TotalUsers:      totalUsers,
			ActiveShipments: activePlans,
			ContainerTypes:  containerTypes,
			SuccessRate:     98.5,
		}
	}

	if role == types.RolePlanner.String() || role == types.RoleAdmin.String() {
		pending, _ := s.q.CountActivePlans(ctx)
		completedToday, _ := s.q.CountCompletedPlansToday(ctx)
		avgUtil, _ := s.q.GetAvgVolumeUtilization(ctx)
		totalItems, _ := s.q.CountTotalItems(ctx)
		resp.Planner = &dto.PlannerStats{
			PendingPlans:   pending,
			CompletedToday: completedToday,
			AvgUtilization: avgUtil,
			ItemsProcessed: totalItems,
		}
	}

	if role == types.RoleOperator.String() || role == types.RoleAdmin.String() {
		active, _ := s.q.CountActivePlans(ctx)
		completed, _ := s.q.CountCompletedPlans(ctx)
		resp.Operator = &dto.OperatorStats{
			ActiveLoads:       active,
			Completed:         completed,
			FailedValidations: 0,
			AvgTimePerLoad:    "24m",
		}
	}

	return resp, nil
}
