package service_test

import (
	"context"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

func TestNewDashboardService(t *testing.T) {
	mockQ := &MockQuerier{}
	svc := service.NewDashboardService(mockQ)

	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestDashboardService_GetStats(t *testing.T) {
	testUUID := uuid.New()
	tests := []struct {
		name              string
		role              string
		workspaceID       *uuid.UUID
		mockSetup         func(*MockQuerier)
		wantAdminStats    bool
		wantPlannerStats  bool
		wantOperatorStats bool
		verifyAdmin       func(*testing.T, *dto.AdminStats)
		verifyPlanner     func(*testing.T, *dto.PlannerStats)
		verifyOperator    func(*testing.T, *dto.OperatorStats)
	}{
		{
			name:        "founder_returns_global_stats",
			role:        types.RoleFounder.String(),
			workspaceID: nil,
			mockSetup: func(m *MockQuerier) {
				m.CountGlobalUsersFunc = func(ctx context.Context) (int64, error) {
					return 150, nil
				}
				m.CountGlobalActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 25, nil
				}
				m.CountGlobalContainersFunc = func(ctx context.Context) (int64, error) {
					return 10, nil
				}
			},
			wantAdminStats:    true,
			wantPlannerStats:  false,
			wantOperatorStats: false,
			verifyAdmin: func(t *testing.T, stats *dto.AdminStats) {
				if stats.TotalUsers != 150 {
					t.Errorf("expected TotalUsers=150, got %d", stats.TotalUsers)
				}
				if stats.ActiveShipments != 25 {
					t.Errorf("expected ActiveShipments=25, got %d", stats.ActiveShipments)
				}
				if stats.ContainerTypes != 10 {
					t.Errorf("expected ContainerTypes=10, got %d", stats.ContainerTypes)
				}
			},
		},
		{
			name:        "workspace_admin_returns_workspace_stats",
			role:        types.RoleAdmin.String(),
			workspaceID: &testUUID,
			mockSetup: func(m *MockQuerier) {
				m.CountWorkspaceMembersFunc = func(ctx context.Context, wid uuid.UUID) (int64, error) {
					return 5, nil
				}
				m.CountWorkspaceActivePlansFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 3, nil
				}
				m.CountWorkspaceContainersFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 2, nil
				}
				// Also fetched for planner stats which admin sees
				m.CountWorkspaceCompletedPlansTodayFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 1, nil
				}
				m.GetWorkspaceAvgVolumeUtilizationFunc = func(ctx context.Context, wid *uuid.UUID) (float64, error) {
					return 50.0, nil
				}
				m.CountWorkspaceItemsFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 100, nil
				}
				// Also fetched for operator stats which admin sees
				m.CountWorkspaceCompletedPlansFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 10, nil
				}
			},
			wantAdminStats:    true,
			wantPlannerStats:  true,
			wantOperatorStats: true,
			verifyAdmin: func(t *testing.T, stats *dto.AdminStats) {
				if stats.TotalUsers != 5 {
					t.Errorf("expected TotalMembers=5, got %d", stats.TotalUsers)
				}
				if stats.ActiveShipments != 3 {
					t.Errorf("expected ActiveShipments=3, got %d", stats.ActiveShipments)
				}
			},
		},
		{
			name:        "planner_returns_only_planner_stats",
			role:        types.RolePlanner.String(),
			workspaceID: &testUUID,
			mockSetup: func(m *MockQuerier) {
				m.CountWorkspaceActivePlansFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 15, nil
				}
				m.CountWorkspaceCompletedPlansTodayFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 3, nil
				}
				m.GetWorkspaceAvgVolumeUtilizationFunc = func(ctx context.Context, wid *uuid.UUID) (float64, error) {
					return 82.3, nil
				}
				m.CountWorkspaceItemsFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 250, nil
				}
			},
			wantAdminStats:    false,
			wantPlannerStats:  true,
			wantOperatorStats: false,
			verifyPlanner: func(t *testing.T, stats *dto.PlannerStats) {
				if stats.PendingPlans != 15 {
					t.Errorf("expected PendingPlans=15, got %d", stats.PendingPlans)
				}
				if stats.CompletedToday != 3 {
					t.Errorf("expected CompletedToday=3, got %d", stats.CompletedToday)
				}
				if stats.AvgUtilization != 82.3 {
					t.Errorf("expected AvgUtilization=82.3, got %f", stats.AvgUtilization)
				}
				if stats.ItemsProcessed != 250 {
					t.Errorf("expected ItemsProcessed=250, got %d", stats.ItemsProcessed)
				}
			},
		},
		{
			name:        "operator_returns_only_operator_stats",
			role:        types.RoleOperator.String(),
			workspaceID: &testUUID,
			mockSetup: func(m *MockQuerier) {
				m.CountWorkspaceActivePlansFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 8, nil
				}
				m.CountWorkspaceCompletedPlansFunc = func(ctx context.Context, wid *uuid.UUID) (int64, error) {
					return 45, nil
				}
			},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: true,
			verifyOperator: func(t *testing.T, stats *dto.OperatorStats) {
				if stats.ActiveLoads != 8 {
					t.Errorf("expected ActiveLoads=8, got %d", stats.ActiveLoads)
				}
				if stats.Completed != 45 {
					t.Errorf("expected Completed=45, got %d", stats.Completed)
				}
			},
		},
		{
			name:              "unknown_role_returns_empty_stats",
			role:              "unknown",
			workspaceID:       &testUUID,
			mockSetup:         func(m *MockQuerier) {},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockQ := &MockQuerier{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockQ)
			}
			svc := service.NewDashboardService(mockQ)

			// Execute
			resp, err := svc.GetStats(context.Background(), tt.role, tt.workspaceID)

			// Assert no error (current implementation never errors)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify response structure
			if resp == nil {
				t.Fatal("expected non-nil response")
			}

			// Verify admin stats presence
			if tt.wantAdminStats && resp.Admin == nil {
				t.Error("expected admin stats but got nil")
			}
			if !tt.wantAdminStats && resp.Admin != nil {
				t.Error("expected no admin stats but got some")
			}

			// Verify planner stats presence
			if tt.wantPlannerStats && resp.Planner == nil {
				t.Error("expected planner stats but got nil")
			}
			if !tt.wantPlannerStats && resp.Planner != nil {
				t.Error("expected no planner stats but got some")
			}

			// Verify operator stats presence
			if tt.wantOperatorStats && resp.Operator == nil {
				t.Error("expected operator stats but got nil")
			}
			if !tt.wantOperatorStats && resp.Operator != nil {
				t.Error("expected no operator stats but got some")
			}

			// Detailed verification callbacks
			if tt.verifyAdmin != nil && resp.Admin != nil {
				tt.verifyAdmin(t, resp.Admin)
			}
			if tt.verifyPlanner != nil && resp.Planner != nil {
				tt.verifyPlanner(t, resp.Planner)
			}
			if tt.verifyOperator != nil && resp.Operator != nil {
				tt.verifyOperator(t, resp.Operator)
			}
		})
	}
}
