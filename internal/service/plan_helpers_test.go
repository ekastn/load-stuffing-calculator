package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/mocks"
	"github.com/ekastn/load-stuffing-calculator/internal/packer"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockPacker is a simple mock for testing plan helpers
type mockPackerHelper struct{}

func (m *mockPackerHelper) Pack(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
	return packer.PackingResult{}, nil
}

// TestActorFromContext tests the private actorFromContext helper
func TestActorFromContext(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		ctx         context.Context
		expectError bool
		errorMsg    string
		expectedID  *uuid.UUID
	}{
		{
			name:        "missing_role_in_context",
			ctx:         context.Background(),
			expectError: true,
			errorMsg:    "missing role",
		},
		{
			name:        "empty_role_in_context",
			ctx:         auth.WithRole(context.Background(), ""),
			expectError: true,
			errorMsg:    "missing role",
		},
		{
			name:        "missing_user_id_in_context",
			ctx:         auth.WithRole(context.Background(), types.RolePlanner.String()),
			expectError: true,
			errorMsg:    "missing user id",
		},
		{
			name:        "empty_user_id_in_context",
			ctx:         auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), ""),
			expectError: true,
			errorMsg:    "missing user id",
		},
		{
			name:        "invalid_user_id_not_uuid",
			ctx:         auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), "not-a-uuid"),
			expectError: true,
			errorMsg:    "invalid user id",
		},
		{
			name:        "valid_context_with_role_and_user_id",
			ctx:         auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), userID.String()),
			expectError: false,
			expectedID:  &userID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor, err := actorFromContext(tt.ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, actor)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, actor)
				if tt.expectedID != nil {
					assert.Equal(t, *tt.expectedID, actor.id)
				}
			}
		})
	}
}

// TestResolvePlanScope tests the private resolvePlanScope helper
func TestResolvePlanScope(t *testing.T) {
	planID := uuid.New()
	userID := uuid.New()
	workspaceID := uuid.New()
	otherWorkspaceID := uuid.New()

	tests := []struct {
		name        string
		ctx         context.Context
		planID      uuid.UUID
		mockSetup   func(*mocks.MockQuerier)
		expectError bool
		errorMsg    string
	}{
		{
			name:        "actor_from_context_fails",
			ctx:         context.Background(),
			planID:      planID,
			mockSetup:   func(m *mocks.MockQuerier) {},
			expectError: true,
			errorMsg:    "missing role",
		},
		{
			name:   "trial_user_can_access_own_plan",
			ctx:    auth.WithUserID(auth.WithRole(context.Background(), types.RoleTrial.String()), userID.String()),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetLoadPlanForGuestFunc = func(ctx context.Context, params store.GetLoadPlanForGuestParams) (store.LoadPlan, error) {
					assert.Equal(t, planID, params.PlanID)
					assert.Equal(t, userID, params.CreatedByID)
					return store.LoadPlan{PlanID: planID, CreatedByID: userID}, nil
				}
			},
			expectError: false,
		},
		{
			name:   "trial_user_cannot_access_other_plan",
			ctx:    auth.WithUserID(auth.WithRole(context.Background(), types.RoleTrial.String()), userID.String()),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetLoadPlanForGuestFunc = func(ctx context.Context, params store.GetLoadPlanForGuestParams) (store.LoadPlan, error) {
					return store.LoadPlan{}, sql.ErrNoRows
				}
			},
			expectError: true,
			errorMsg:    "forbidden",
		},
		{
			name:   "founder_with_invalid_workspace_override",
			ctx:    auth.WithWorkspaceOverrideID(auth.WithUserID(auth.WithRole(context.Background(), types.RoleFounder.String()), userID.String()), "invalid-uuid"),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
			},
			expectError: true,
			errorMsg:    "invalid workspace id",
		},
		{
			name:   "founder_with_no_override_can_access_any_plan",
			ctx:    auth.WithUserID(auth.WithRole(context.Background(), types.RoleFounder.String()), userID.String()),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetLoadPlanAnyFunc = func(ctx context.Context, id uuid.UUID) (store.LoadPlan, error) {
					assert.Equal(t, planID, id)
					return store.LoadPlan{PlanID: planID}, nil
				}
			},
			expectError: false,
		},
		{
			name:   "founder_with_no_override_plan_not_found",
			ctx:    auth.WithUserID(auth.WithRole(context.Background(), types.RoleFounder.String()), userID.String()),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetLoadPlanAnyFunc = func(ctx context.Context, id uuid.UUID) (store.LoadPlan, error) {
					return store.LoadPlan{}, sql.ErrNoRows
				}
			},
			expectError: true,
			errorMsg:    "no rows",
		},
		{
			name:   "founder_with_valid_workspace_override",
			ctx:    auth.WithWorkspaceOverrideID(auth.WithWorkspaceID(auth.WithUserID(auth.WithRole(context.Background(), types.RoleFounder.String()), userID.String()), workspaceID.String()), otherWorkspaceID.String()),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetLoadPlanFunc = func(ctx context.Context, params store.GetLoadPlanParams) (store.LoadPlan, error) {
					assert.Equal(t, planID, params.PlanID)
					assert.Equal(t, otherWorkspaceID, *params.WorkspaceID)
					return store.LoadPlan{PlanID: planID, WorkspaceID: &otherWorkspaceID}, nil
				}
			},
			expectError: false,
		},
		{
			name:   "planner_can_access_plan_in_workspace",
			ctx:    auth.WithWorkspaceID(auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), userID.String()), workspaceID.String()),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetLoadPlanFunc = func(ctx context.Context, params store.GetLoadPlanParams) (store.LoadPlan, error) {
					assert.Equal(t, planID, params.PlanID)
					assert.Equal(t, workspaceID, *params.WorkspaceID)
					return store.LoadPlan{PlanID: planID, WorkspaceID: &workspaceID}, nil
				}
			},
			expectError: false,
		},
		{
			name:   "planner_plan_not_found_in_workspace",
			ctx:    auth.WithWorkspaceID(auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), userID.String()), workspaceID.String()),
			planID: planID,
			mockSetup: func(m *mocks.MockQuerier) {
				m.GetLoadPlanFunc = func(ctx context.Context, params store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{}, sql.ErrNoRows
				}
			},
			expectError: true,
			errorMsg:    "no rows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &mocks.MockQuerier{}
			tt.mockSetup(mockQ)
			mockPacker := &mockPackerHelper{}
			s := &planService{q: mockQ, p: mockPacker}

			scope, err := s.resolvePlanScope(tt.ctx, tt.planID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, scope)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, scope)
				assert.Equal(t, tt.planID, scope.plan.PlanID)
			}
		})
	}
}
