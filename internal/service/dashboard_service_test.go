package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
)

func TestNewDashboardService(t *testing.T) {
	mockQ := &MockQuerier{}
	svc := service.NewDashboardService(mockQ)

	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestDashboardService_GetStats(t *testing.T) {
	tests := []struct {
		name              string
		role              string
		mockSetup         func(*MockQuerier)
		wantAdminStats    bool
		wantPlannerStats  bool
		wantOperatorStats bool
		verifyAdmin       func(*testing.T, *dto.AdminStats)
		verifyPlanner     func(*testing.T, *dto.PlannerStats)
		verifyOperator    func(*testing.T, *dto.OperatorStats)
	}{
		{
			name: "admin_returns_all_stats",
			role: types.RoleAdmin.String(),
			mockSetup: func(m *MockQuerier) {
				m.CountTotalUsersFunc = func(ctx context.Context) (int64, error) {
					return 150, nil
				}
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 25, nil
				}
				m.CountContainersFunc = func(ctx context.Context) (int64, error) {
					return 10, nil
				}
				m.CountCompletedPlansTodayFunc = func(ctx context.Context) (int64, error) {
					return 5, nil
				}
				m.GetAvgVolumeUtilizationFunc = func(ctx context.Context) (float64, error) {
					return 78.5, nil
				}
				m.CountTotalItemsFunc = func(ctx context.Context) (int64, error) {
					return 500, nil
				}
				m.CountCompletedPlansFunc = func(ctx context.Context) (int64, error) {
					return 100, nil
				}
			},
			wantAdminStats:    true,
			wantPlannerStats:  true,
			wantOperatorStats: true,
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
				if stats.SuccessRate != 98.5 {
					t.Errorf("expected SuccessRate=98.5, got %f", stats.SuccessRate)
				}
			},
			verifyPlanner: func(t *testing.T, stats *dto.PlannerStats) {
				if stats.PendingPlans != 25 {
					t.Errorf("expected PendingPlans=25, got %d", stats.PendingPlans)
				}
				if stats.CompletedToday != 5 {
					t.Errorf("expected CompletedToday=5, got %d", stats.CompletedToday)
				}
				if stats.AvgUtilization != 78.5 {
					t.Errorf("expected AvgUtilization=78.5, got %f", stats.AvgUtilization)
				}
				if stats.ItemsProcessed != 500 {
					t.Errorf("expected ItemsProcessed=500, got %d", stats.ItemsProcessed)
				}
			},
			verifyOperator: func(t *testing.T, stats *dto.OperatorStats) {
				if stats.ActiveLoads != 25 {
					t.Errorf("expected ActiveLoads=25, got %d", stats.ActiveLoads)
				}
				if stats.Completed != 100 {
					t.Errorf("expected Completed=100, got %d", stats.Completed)
				}
				if stats.FailedValidations != 0 {
					t.Errorf("expected FailedValidations=0, got %d", stats.FailedValidations)
				}
				if stats.AvgTimePerLoad != "24m" {
					t.Errorf("expected AvgTimePerLoad='24m', got %s", stats.AvgTimePerLoad)
				}
			},
		},
		{
			name: "planner_returns_only_planner_stats",
			role: types.RolePlanner.String(),
			mockSetup: func(m *MockQuerier) {
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 15, nil
				}
				m.CountCompletedPlansTodayFunc = func(ctx context.Context) (int64, error) {
					return 3, nil
				}
				m.GetAvgVolumeUtilizationFunc = func(ctx context.Context) (float64, error) {
					return 82.3, nil
				}
				m.CountTotalItemsFunc = func(ctx context.Context) (int64, error) {
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
			name: "operator_returns_only_operator_stats",
			role: types.RoleOperator.String(),
			mockSetup: func(m *MockQuerier) {
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 8, nil
				}
				m.CountCompletedPlansFunc = func(ctx context.Context) (int64, error) {
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
				if stats.FailedValidations != 0 {
					t.Errorf("expected FailedValidations=0, got %d", stats.FailedValidations)
				}
				if stats.AvgTimePerLoad != "24m" {
					t.Errorf("expected AvgTimePerLoad='24m', got %s", stats.AvgTimePerLoad)
				}
			},
		},
		{
			name:              "unknown_role_returns_empty_stats",
			role:              "unknown",
			mockSetup:         func(m *MockQuerier) {},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: false,
		},
		{
			name:              "empty_role_returns_empty_stats",
			role:              "",
			mockSetup:         func(m *MockQuerier) {},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: false,
		},
		{
			name:              "uppercase_role_not_matched",
			role:              "ADMIN",
			mockSetup:         func(m *MockQuerier) {},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: false,
		},
		{
			name: "admin_with_zero_values",
			role: types.RoleAdmin.String(),
			mockSetup: func(m *MockQuerier) {
				m.CountTotalUsersFunc = func(ctx context.Context) (int64, error) {
					return 0, nil
				}
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 0, nil
				}
				m.CountContainersFunc = func(ctx context.Context) (int64, error) {
					return 0, nil
				}
				m.CountCompletedPlansTodayFunc = func(ctx context.Context) (int64, error) {
					return 0, nil
				}
				m.GetAvgVolumeUtilizationFunc = func(ctx context.Context) (float64, error) {
					return 0.0, nil
				}
				m.CountTotalItemsFunc = func(ctx context.Context) (int64, error) {
					return 0, nil
				}
				m.CountCompletedPlansFunc = func(ctx context.Context) (int64, error) {
					return 0, nil
				}
			},
			wantAdminStats:    true,
			wantPlannerStats:  true,
			wantOperatorStats: true,
			verifyAdmin: func(t *testing.T, stats *dto.AdminStats) {
				if stats.TotalUsers != 0 {
					t.Errorf("expected TotalUsers=0, got %d", stats.TotalUsers)
				}
				if stats.ActiveShipments != 0 {
					t.Errorf("expected ActiveShipments=0, got %d", stats.ActiveShipments)
				}
				if stats.ContainerTypes != 0 {
					t.Errorf("expected ContainerTypes=0, got %d", stats.ContainerTypes)
				}
				// SuccessRate is hardcoded to 98.5
				if stats.SuccessRate != 98.5 {
					t.Errorf("expected SuccessRate=98.5, got %f", stats.SuccessRate)
				}
			},
			verifyPlanner: func(t *testing.T, stats *dto.PlannerStats) {
				if stats.PendingPlans != 0 {
					t.Errorf("expected PendingPlans=0, got %d", stats.PendingPlans)
				}
				if stats.CompletedToday != 0 {
					t.Errorf("expected CompletedToday=0, got %d", stats.CompletedToday)
				}
			},
			verifyOperator: func(t *testing.T, stats *dto.OperatorStats) {
				if stats.ActiveLoads != 0 {
					t.Errorf("expected ActiveLoads=0, got %d", stats.ActiveLoads)
				}
				if stats.Completed != 0 {
					t.Errorf("expected Completed=0, got %d", stats.Completed)
				}
			},
		},
		{
			name:              "trial_role_returns_no_stats",
			role:              types.RoleTrial.String(),
			mockSetup:         func(m *MockQuerier) {},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: false,
		},
		{
			name:              "founder_role_returns_no_stats",
			role:              types.RoleFounder.String(),
			mockSetup:         func(m *MockQuerier) {},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: false,
		},
		{
			name:              "user_role_returns_no_stats",
			role:              types.RoleUser.String(),
			mockSetup:         func(m *MockQuerier) {},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: false,
		},
		{
			name: "admin_db_errors_ignored_returns_zero_values",
			role: types.RoleAdmin.String(),
			mockSetup: func(m *MockQuerier) {
				m.CountTotalUsersFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db connection error")
				}
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.CountContainersFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.CountCompletedPlansTodayFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.GetAvgVolumeUtilizationFunc = func(ctx context.Context) (float64, error) {
					return 0.0, fmt.Errorf("db error")
				}
				m.CountTotalItemsFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.CountCompletedPlansFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
			},
			wantAdminStats:    true,
			wantPlannerStats:  true,
			wantOperatorStats: true,
			verifyAdmin: func(t *testing.T, stats *dto.AdminStats) {
				// Should have zero values due to errors being ignored
				if stats.TotalUsers != 0 {
					t.Errorf("expected TotalUsers=0 (error ignored), got %d", stats.TotalUsers)
				}
				if stats.ActiveShipments != 0 {
					t.Errorf("expected ActiveShipments=0 (error ignored), got %d", stats.ActiveShipments)
				}
				if stats.ContainerTypes != 0 {
					t.Errorf("expected ContainerTypes=0 (error ignored), got %d", stats.ContainerTypes)
				}
			},
		},
		{
			name: "planner_db_errors_ignored_returns_zero_values",
			role: types.RolePlanner.String(),
			mockSetup: func(m *MockQuerier) {
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.CountCompletedPlansTodayFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.GetAvgVolumeUtilizationFunc = func(ctx context.Context) (float64, error) {
					return 0.0, fmt.Errorf("db error")
				}
				m.CountTotalItemsFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
			},
			wantAdminStats:    false,
			wantPlannerStats:  true,
			wantOperatorStats: false,
			verifyPlanner: func(t *testing.T, stats *dto.PlannerStats) {
				if stats.PendingPlans != 0 {
					t.Errorf("expected PendingPlans=0 (error ignored), got %d", stats.PendingPlans)
				}
				if stats.CompletedToday != 0 {
					t.Errorf("expected CompletedToday=0 (error ignored), got %d", stats.CompletedToday)
				}
			},
		},
		{
			name: "operator_db_errors_ignored_returns_zero_values",
			role: types.RoleOperator.String(),
			mockSetup: func(m *MockQuerier) {
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.CountCompletedPlansFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
			},
			wantAdminStats:    false,
			wantPlannerStats:  false,
			wantOperatorStats: true,
			verifyOperator: func(t *testing.T, stats *dto.OperatorStats) {
				if stats.ActiveLoads != 0 {
					t.Errorf("expected ActiveLoads=0 (error ignored), got %d", stats.ActiveLoads)
				}
				if stats.Completed != 0 {
					t.Errorf("expected Completed=0 (error ignored), got %d", stats.Completed)
				}
			},
		},
		{
			name: "admin_partial_success_some_queries_error",
			role: types.RoleAdmin.String(),
			mockSetup: func(m *MockQuerier) {
				// Some succeed
				m.CountTotalUsersFunc = func(ctx context.Context) (int64, error) {
					return 50, nil
				}
				m.CountActivePlansFunc = func(ctx context.Context) (int64, error) {
					return 10, nil
				}
				// Some fail
				m.CountContainersFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.CountCompletedPlansTodayFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.GetAvgVolumeUtilizationFunc = func(ctx context.Context) (float64, error) {
					return 0.0, fmt.Errorf("db error")
				}
				m.CountTotalItemsFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
				m.CountCompletedPlansFunc = func(ctx context.Context) (int64, error) {
					return 0, fmt.Errorf("db error")
				}
			},
			wantAdminStats:    true,
			wantPlannerStats:  true,
			wantOperatorStats: true,
			verifyAdmin: func(t *testing.T, stats *dto.AdminStats) {
				// Successful queries return values
				if stats.TotalUsers != 50 {
					t.Errorf("expected TotalUsers=50, got %d", stats.TotalUsers)
				}
				if stats.ActiveShipments != 10 {
					t.Errorf("expected ActiveShipments=10, got %d", stats.ActiveShipments)
				}
				// Failed queries return zero
				if stats.ContainerTypes != 0 {
					t.Errorf("expected ContainerTypes=0 (error ignored), got %d", stats.ContainerTypes)
				}
			},
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
			resp, err := svc.GetStats(context.Background(), tt.role)

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
