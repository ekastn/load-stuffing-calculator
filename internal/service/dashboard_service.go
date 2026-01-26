package service

import (
	"context"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

type DashboardService interface {
	GetStats(ctx context.Context, role string, workspaceID *uuid.UUID) (*dto.DashboardStatsResponse, error)
}

type dashboardService struct {
	q store.Querier
}

func NewDashboardService(q store.Querier) DashboardService {
	return &dashboardService{q: q}
}

func (s *dashboardService) GetStats(ctx context.Context, role string, workspaceID *uuid.UUID) (*dto.DashboardStatsResponse, error) {
	resp := &dto.DashboardStatsResponse{}

	// Helper to handle workspace ID requirement
	if (role != types.RoleFounder.String()) && workspaceID == nil {
		return resp, nil
	}

	// 1. Platform Admin (Founder) - Global Stats
	if role == types.RoleFounder.String() {
		totalUsers, _ := s.q.CountGlobalUsers(ctx)
		activePlans, _ := s.q.CountGlobalActivePlans(ctx)
		containerTypes, _ := s.q.CountGlobalContainers(ctx)
		resp.Admin = &dto.AdminStats{
			TotalUsers:      totalUsers,
			ActiveShipments: activePlans,
			ContainerTypes:  containerTypes,
			SuccessRate:     98.5, // Placeholder
		}
		return resp, nil
	}

	// Determine View Permissions
	isAdminView := role == types.RoleAdmin.String() || role == types.RoleOwner.String() || role == types.RolePersonal.String()
	isPlannerView := role == types.RolePlanner.String() || isAdminView
	isOperatorView := role == types.RoleOperator.String() || isAdminView

	// Shared Metrics (Optimization: Fetch once)
	var activePlans int64
	if isAdminView || isPlannerView || isOperatorView {
		activePlans, _ = s.q.CountWorkspaceActivePlans(ctx, workspaceID)
	}

	// 2. Workspace Admin Stats
	if isAdminView {
		totalMembers, _ := s.q.CountWorkspaceMembers(ctx, *workspaceID)
		containerTypes, _ := s.q.CountWorkspaceContainers(ctx, workspaceID) // includes global containers

		resp.Admin = &dto.AdminStats{
			TotalUsers:      totalMembers,
			ActiveShipments: activePlans,
			ContainerTypes:  containerTypes,
			SuccessRate:     98.5,
		}
	}

	// 3. Planner Stats
	if isPlannerView {
		completedToday, _ := s.q.CountWorkspaceCompletedPlansToday(ctx, workspaceID)
		avgUtil, _ := s.q.GetWorkspaceAvgVolumeUtilization(ctx, workspaceID)
		totalItems, _ := s.q.CountWorkspaceItems(ctx, workspaceID)

		resp.Planner = &dto.PlannerStats{
			PendingPlans:   activePlans,
			CompletedToday: completedToday,
			AvgUtilization: avgUtil,
			ItemsProcessed: totalItems,
		}
	}

	// 4. Operator Stats
	if isOperatorView {
		completed, _ := s.q.CountWorkspaceCompletedPlans(ctx, workspaceID)

		resp.Operator = &dto.OperatorStats{
			ActiveLoads:       activePlans,
			Completed:         completed,
			FailedValidations: 0,
			AvgTimePerLoad:    "24m",
		}
	}

	return resp, nil
}
