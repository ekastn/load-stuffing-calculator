package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/packer"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func authedPlannerCtx() context.Context {
	ctx := auth.WithUserID(auth.WithRole(context.Background(), types.RolePlanner.String()), uuid.New().String())
	return auth.WithWorkspaceID(ctx, uuid.New().String())
}

func authedTrialCtxWithID(id uuid.UUID) context.Context {
	return auth.WithUserID(auth.WithRole(context.Background(), types.RoleTrial.String()), id.String())
}

func authedFounderCtx() context.Context {
	return auth.WithUserID(auth.WithRole(context.Background(), types.RoleFounder.String()), uuid.New().String())
}

func authedFounderCtxWithWorkspaceOverride(wsID uuid.UUID) context.Context {
	ctx := auth.WithUserID(auth.WithRole(context.Background(), types.RoleFounder.String()), uuid.New().String())
	ctx = auth.WithWorkspaceID(ctx, wsID.String())
	return auth.WithWorkspaceOverrideID(ctx, wsID.String())
}

func TestPlanService_CreateCompletePlan(t *testing.T) {
	planID := uuid.New()

	// Common item for tests
	itemReq := dto.CreatePlanItem{
		Label:    stringPtr("Item A"),
		LengthMM: 100,
		WidthMM:  100,
		HeightMM: 100,
		WeightKG: 10,
		Quantity: 1,
	}

	t.Run("success_with_preset_container", func(t *testing.T) {
		contID := uuid.New()
		contName := "20ft Standard"
		contLength := 5898.0
		contWidth := 2352.0
		contHeight := 2393.0
		contMaxWeight := 28200.0

		mockQ := &MockQuerier{
			GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
				assert.Equal(t, contID, arg.ContainerID)
				return store.Container{
					ContainerID:   contID,
					Name:          contName,
					InnerLengthMm: toNumeric(contLength),
					InnerWidthMm:  toNumeric(contWidth),
					InnerHeightMm: toNumeric(contHeight),
					MaxWeightKg:   toNumeric(contMaxWeight),
				}, nil
			},
			CreateLoadPlanFunc: func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
				assert.Equal(t, contName, *arg.ContLabel)
				return store.LoadPlan{
					PlanID:      planID,
					PlanCode:    arg.PlanCode,
					Status:      arg.Status,
					ContLabel:   arg.ContLabel,
					LengthMm:    arg.LengthMm,
					WidthMm:     arg.WidthMm,
					HeightMm:    arg.HeightMm,
					MaxWeightKg: arg.MaxWeightKg,
					CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				}, nil
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				assert.Equal(t, planID, *arg.PlanID)
				return store.LoadItem{ItemID: uuid.New()}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "Test Plan",
			Container: dto.CreatePlanContainer{
				ContainerID: stringPtr(contID.String()),
			},
			Items:         []dto.CreatePlanItem{itemReq},
			AutoCalculate: boolPtr(false),
		}
		resp, err := s.CreateCompletePlan(authedPlannerCtx(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, planID.String(), resp.PlanID)
		assert.Equal(t, types.PlanStatusDraft.String(), resp.Status)
	})

	t.Run("founder_with_override_uses_override_workspace", func(t *testing.T) {
		contID := uuid.New()
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()

		mockQ := &MockQuerier{
			GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
				assert.Equal(t, contID, arg.ContainerID)
				assert.Equal(t, &overrideWorkspaceID, arg.WorkspaceID)
				return store.Container{ContainerID: contID, Name: "Preset", InnerLengthMm: toNumeric(1), InnerWidthMm: toNumeric(1), InnerHeightMm: toNumeric(1), MaxWeightKg: toNumeric(1)}, nil
			},
			CreateLoadPlanFunc: func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
				assert.Equal(t, &overrideWorkspaceID, arg.WorkspaceID)
				return store.LoadPlan{PlanID: planID, PlanCode: arg.PlanCode, CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}}, nil
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{ItemID: uuid.New()}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		ctx := auth.WithWorkspaceOverrideID(ctxWithRoleAndWorkspace("founder", workspaceID), overrideWorkspaceID.String())
		ctx = auth.WithUserID(ctx, uuid.New().String())

		req := dto.CreatePlanRequest{
			Title: "Founder Override Plan",
			Container: dto.CreatePlanContainer{
				ContainerID: stringPtr(contID.String()),
			},
			Items:         []dto.CreatePlanItem{itemReq},
			AutoCalculate: boolPtr(false),
		}

		resp, err := s.CreateCompletePlan(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("success_with_custom_container", func(t *testing.T) {
		customLength := 1000.0
		customWidth := 500.0
		customHeight := 500.0
		customMaxWeight := 5000.0

		mockQ := &MockQuerier{
			CreateLoadPlanFunc: func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
				assert.Equal(t, "Custom Container", *arg.ContLabel)
				assert.Equal(t, toNumeric(customLength), arg.LengthMm)
				return store.LoadPlan{
					PlanID:    planID,
					PlanCode:  arg.PlanCode,
					CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
				}, nil
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{ItemID: uuid.New()}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "Custom Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    &customLength,
				WidthMM:     &customWidth,
				HeightMM:    &customHeight,
				MaxWeightKG: &customMaxWeight,
			},
			Items:         []dto.CreatePlanItem{itemReq},
			AutoCalculate: boolPtr(false), // Ensure auto-calc is off for this test too
		}
		resp, err := s.CreateCompletePlan(authedPlannerCtx(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("error_invalid_container_id", func(t *testing.T) {
		mockQ := &MockQuerier{}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "Invalid Container ID Plan",
			Container: dto.CreatePlanContainer{
				ContainerID: stringPtr("invalid-uuid"),
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		resp, err := s.CreateCompletePlan(authedPlannerCtx(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid container_id format")
	})

	t.Run("error_container_not_found", func(t *testing.T) {
		contID := uuid.New()
		mockQ := &MockQuerier{
			GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
				return store.Container{}, fmt.Errorf("container not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "Container Not Found Plan",
			Container: dto.CreatePlanContainer{
				ContainerID: stringPtr(contID.String()),
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		resp, err := s.CreateCompletePlan(authedPlannerCtx(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "container not found")
	})

	t.Run("error_custom_container_dims_missing", func(t *testing.T) {
		mockQ := &MockQuerier{}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title:     "Missing Custom Dims Plan",
			Container: dto.CreatePlanContainer{
				// Missing LengthMM, WidthMM, HeightMM, MaxWeightKG
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		resp, err := s.CreateCompletePlan(authedPlannerCtx(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "custom container dimensions are required")
	})

	t.Run("error_create_load_plan_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			CreateLoadPlanFunc: func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{}, fmt.Errorf("db error")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "DB Error Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    floatPtr(1000),
				WidthMM:     floatPtr(500),
				HeightMM:    floatPtr(500),
				MaxWeightKG: floatPtr(5000),
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		resp, err := s.CreateCompletePlan(authedPlannerCtx(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to create plan")
	})

	t.Run("error_add_load_item_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			CreateLoadPlanFunc: func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:      planID,
					PlanCode:    arg.PlanCode,
					Status:      arg.Status,
					ContLabel:   arg.ContLabel,
					LengthMm:    arg.LengthMm,
					WidthMm:     arg.WidthMm,
					HeightMm:    arg.HeightMm,
					MaxWeightKg: arg.MaxWeightKg,
					CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				}, nil
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{}, fmt.Errorf("item db error")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "Item Add Error Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    floatPtr(1000),
				WidthMM:     floatPtr(500),
				HeightMM:    floatPtr(500),
				MaxWeightKG: floatPtr(5000),
			},
			Items: []dto.CreatePlanItem{itemReq},
		}
		resp, err := s.CreateCompletePlan(authedPlannerCtx(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to add item")
	})

	t.Run("trial_limit_reached", func(t *testing.T) {
		guestID := uuid.New()
		createCalled := false

		mockQ := &MockQuerier{
			CountPlansByCreatorFunc: func(ctx context.Context, arg store.CountPlansByCreatorParams) (int64, error) {
				assert.Equal(t, "guest", arg.CreatedByType)
				assert.Equal(t, guestID, arg.CreatedByID)
				return 3, nil
			},
			CreateLoadPlanFunc: func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
				createCalled = true
				return store.LoadPlan{}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "Trial Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    floatPtr(1000),
				WidthMM:     floatPtr(500),
				HeightMM:    floatPtr(500),
				MaxWeightKG: floatPtr(5000),
			},
			Items:         []dto.CreatePlanItem{itemReq},
			AutoCalculate: boolPtr(false),
		}

		resp, err := s.CreateCompletePlan(authedTrialCtxWithID(guestID), req)
		assert.ErrorIs(t, err, service.ErrTrialLimitReached)
		assert.Nil(t, resp)
		assert.False(t, createCalled)
	})

	t.Run("trial_creates_with_guest_owner_fields", func(t *testing.T) {
		guestID := uuid.New()
		planID := uuid.New()

		mockQ := &MockQuerier{
			CountPlansByCreatorFunc: func(ctx context.Context, arg store.CountPlansByCreatorParams) (int64, error) {
				assert.Equal(t, "guest", arg.CreatedByType)
				assert.Equal(t, guestID, arg.CreatedByID)
				return 0, nil
			},
			CreateLoadPlanFunc: func(ctx context.Context, arg store.CreateLoadPlanParams) (store.LoadPlan, error) {
				assert.Equal(t, "guest", arg.CreatedByType)
				assert.Equal(t, guestID, arg.CreatedByID)
				return store.LoadPlan{PlanID: planID, PlanCode: arg.PlanCode, CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}}, nil
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{ItemID: uuid.New()}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.CreatePlanRequest{
			Title: "Trial Plan",
			Container: dto.CreatePlanContainer{
				LengthMM:    floatPtr(1000),
				WidthMM:     floatPtr(500),
				HeightMM:    floatPtr(500),
				MaxWeightKG: floatPtr(5000),
			},
			Items:         []dto.CreatePlanItem{itemReq},
			AutoCalculate: boolPtr(false),
		}

		resp, err := s.CreateCompletePlan(authedTrialCtxWithID(guestID), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, planID.String(), resp.PlanID)
	})
}

func TestPlanService_GetPlan(t *testing.T) {
	planID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				assert.Equal(t, planID, arg.PlanID)
				return store.LoadPlan{
					PlanID:      planID,
					PlanCode:    "CODE",
					ContLabel:   stringPtr("Test Container"),
					LengthMm:    toNumeric(1000),
					WidthMm:     toNumeric(1000),
					HeightMm:    toNumeric(1000),
					MaxWeightKg: toNumeric(1000),
					CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				}, nil
			},
			ListLoadItemsFunc: func(ctx context.Context, id *uuid.UUID) ([]store.LoadItem, error) {
				assert.Equal(t, planID, *id)
				return []store.LoadItem{
					{ItemID: uuid.New(), ItemLabel: stringPtr("Item1"), Quantity: 2, LengthMm: toNumeric(100), WidthMm: toNumeric(100), HeightMm: toNumeric(100), WeightKg: toNumeric(10), AllowRotation: boolPtr(true), ColorHex: stringPtr("#aabbcc")},
				}, nil
			},
			GetPlanResultFunc: func(ctx context.Context, id *uuid.UUID) (store.PlanResult, error) {
				return store.PlanResult{}, fmt.Errorf("no result")
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.GetPlan(authedPlannerCtx(), planID.String())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, planID.String(), resp.PlanID)
		assert.Len(t, resp.Items, 1)
		assert.Equal(t, 2, resp.Stats.TotalItems)
	})

	t.Run("trial_scoped", func(t *testing.T) {
		guestID := uuid.New()
		getGuestCalled := false
		getUserCalled := false

		mockQ := &MockQuerier{
			GetLoadPlanForGuestFunc: func(ctx context.Context, arg store.GetLoadPlanForGuestParams) (store.LoadPlan, error) {
				getGuestCalled = true
				assert.Equal(t, planID, arg.PlanID)
				assert.Equal(t, guestID, arg.CreatedByID)
				return store.LoadPlan{PlanID: planID, PlanCode: "CODE", CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}}, nil
			},
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				getUserCalled = true
				return store.LoadPlan{}, nil
			},
			ListLoadItemsFunc: func(ctx context.Context, id *uuid.UUID) ([]store.LoadItem, error) {
				return []store.LoadItem{}, nil
			},
			GetPlanResultFunc: func(ctx context.Context, id *uuid.UUID) (store.PlanResult, error) {
				return store.PlanResult{}, fmt.Errorf("no result")
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.GetPlan(authedTrialCtxWithID(guestID), planID.String())
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, getGuestCalled)
		assert.False(t, getUserCalled)
	})

	t.Run("error_plan_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{}, fmt.Errorf("plan not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.GetPlan(authedPlannerCtx(), uuid.New().String())

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("error_list_items_fails", func(t *testing.T) {
		planID := uuid.New()
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID}, nil
			},
			ListLoadItemsFunc: func(ctx context.Context, id *uuid.UUID) ([]store.LoadItem, error) {
				return nil, fmt.Errorf("list items db error")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.GetPlan(authedPlannerCtx(), planID.String())

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "list items db error")
	})

	t.Run("error_invalid_plan_id", func(t *testing.T) {
		mockQ := &MockQuerier{}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.GetPlan(authedPlannerCtx(), "invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid plan id")
	})

	t.Run("with_calculation_result_and_placements", func(t *testing.T) {
		planID := uuid.New()
		resultID := uuid.New()
		itemID := uuid.New()
		placementID := uuid.New()

		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:      planID,
					PlanCode:    "CODE",
					ContLabel:   stringPtr("Test Container"),
					LengthMm:    toNumeric(1000),
					WidthMm:     toNumeric(1000),
					HeightMm:    toNumeric(1000),
					MaxWeightKg: toNumeric(1000),
					CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
				}, nil
			},
			ListLoadItemsFunc: func(ctx context.Context, id *uuid.UUID) ([]store.LoadItem, error) {
				return []store.LoadItem{
					{ItemID: itemID, ItemLabel: stringPtr("Item1"), Quantity: 1, LengthMm: toNumeric(100), WidthMm: toNumeric(100), HeightMm: toNumeric(100), WeightKg: toNumeric(10), AllowRotation: boolPtr(true)},
				}, nil
			},
			GetPlanResultFunc: func(ctx context.Context, id *uuid.UUID) (store.PlanResult, error) {
				return store.PlanResult{
					ResultID:             resultID,
					PlanID:               &planID,
					TotalLoadedWeightKg:  toNumeric(10.0),
					VolumeUtilizationPct: toNumeric(75.5),
					IsFeasible:           boolPtr(true),
				}, nil
			},
			ListPlanPlacementsFunc: func(ctx context.Context, resID *uuid.UUID) ([]store.PlanPlacement, error) {
				assert.Equal(t, resultID, *resID)
				rotCode := int32(0)
				return []store.PlanPlacement{
					{
						PlacementID:  placementID,
						ResultID:     &resultID,
						ItemID:       &itemID,
						PosX:         toNumeric(0),
						PosY:         toNumeric(0),
						PosZ:         toNumeric(0),
						RotationCode: &rotCode,
						StepNumber:   1,
					},
				}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.GetPlan(authedPlannerCtx(), planID.String())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Calculation)
		assert.Equal(t, resultID.String(), resp.Calculation.JobID)
		assert.Equal(t, types.PlanStatusCompleted.String(), resp.Calculation.Status)
		assert.Equal(t, 75.5, resp.Calculation.VolumeUtilization)
		assert.Len(t, resp.Calculation.Placements, 1)
		assert.Equal(t, itemID.String(), resp.Calculation.Placements[0].ItemID)
	})

	t.Run("with_calculation_result_partial", func(t *testing.T) {
		planID := uuid.New()
		resultID := uuid.New()

		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:    planID,
					PlanCode:  "CODE",
					CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
				}, nil
			},
			ListLoadItemsFunc: func(ctx context.Context, id *uuid.UUID) ([]store.LoadItem, error) {
				return []store.LoadItem{}, nil
			},
			GetPlanResultFunc: func(ctx context.Context, id *uuid.UUID) (store.PlanResult, error) {
				return store.PlanResult{
					ResultID:             resultID,
					PlanID:               &planID,
					TotalLoadedWeightKg:  toNumeric(10.0),
					VolumeUtilizationPct: toNumeric(50.0),
					IsFeasible:           boolPtr(false), // Partial/infeasible
				}, nil
			},
			ListPlanPlacementsFunc: func(ctx context.Context, resID *uuid.UUID) ([]store.PlanPlacement, error) {
				return []store.PlanPlacement{}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.GetPlan(authedPlannerCtx(), planID.String())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Calculation)
		// When IsFeasible is false, status should be PARTIAL
		assert.Equal(t, types.PlanStatusPartial.String(), resp.Calculation.Status)
	})
}

func TestPlanService_ListPlans(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockQ := &MockQuerier{
			ListLoadPlansFunc: func(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error) {
				fmt.Println("DEBUG: Mock ListLoadPlansFunc called")
				return []store.LoadPlan{
					{PlanID: uuid.New(), PlanCode: "P1", CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}, Status: stringPtr("DRAFT"), ContLabel: stringPtr("Cont1"), LengthMm: toNumeric(100), WidthMm: toNumeric(100), HeightMm: toNumeric(100), MaxWeightKg: toNumeric(1000)},
				}, nil
			},
			// Need to mock ListLoadItems for stats aggregation
			ListLoadItemsFunc: func(ctx context.Context, id *uuid.UUID) ([]store.LoadItem, error) {
				return []store.LoadItem{
					{ItemID: uuid.New(), PlanID: id, LengthMm: toNumeric(100), WidthMm: toNumeric(100), HeightMm: toNumeric(100), WeightKg: toNumeric(10), Quantity: 1},
				}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.ListPlans(authedPlannerCtx(), 1, 10)

		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "P1", resp[0].PlanCode)
		assert.Equal(t, types.PlanStatusDraft.String(), resp[0].Status)
		assert.Equal(t, 0, resp[0].TotalItems)
		assert.Equal(t, 0.0, resp[0].TotalWeightKG)
		// Volume and Utilization would need calculation here
	})

	t.Run("trial_scoped", func(t *testing.T) {
		guestID := uuid.New()
		listGuestCalled := false
		listUserCalled := false

		mockQ := &MockQuerier{
			ListLoadPlansForGuestFunc: func(ctx context.Context, arg store.ListLoadPlansForGuestParams) ([]store.LoadPlan, error) {
				listGuestCalled = true
				assert.Equal(t, guestID, arg.CreatedByID)
				assert.Equal(t, int32(10), arg.Limit)
				assert.Equal(t, int32(0), arg.Offset)
				return []store.LoadPlan{{PlanID: uuid.New(), PlanCode: "P1", CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}, Status: stringPtr("DRAFT")}}, nil
			},
			ListLoadPlansFunc: func(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error) {
				listUserCalled = true
				return nil, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.ListPlans(authedTrialCtxWithID(guestID), 1, 10)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.True(t, listGuestCalled)
		assert.False(t, listUserCalled)
	})

	t.Run("db_error", func(t *testing.T) {
		mockQ := &MockQuerier{
			ListLoadPlansFunc: func(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error) {
				return nil, fmt.Errorf("db error")
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.ListPlans(authedPlannerCtx(), 1, 10)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("missing_workspace_context", func(t *testing.T) {
		mockQ := &MockQuerier{}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.ListPlans(context.Background(), 1, 10)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("founder_list_all", func(t *testing.T) {
		mockQ := &MockQuerier{
			ListLoadPlansAllFunc: func(ctx context.Context, arg store.ListLoadPlansAllParams) ([]store.LoadPlan, error) {
				return []store.LoadPlan{
					{PlanID: uuid.New(), PlanCode: "P1", CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}, Status: stringPtr("DRAFT")},
					{PlanID: uuid.New(), PlanCode: "P2", CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}, Status: stringPtr("DRAFT")},
				}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.ListPlans(authedFounderCtx(), 1, 10)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
	})

	t.Run("founder_with_workspace_override", func(t *testing.T) {
		workspaceID := uuid.New()
		mockQ := &MockQuerier{
			ListLoadPlansFunc: func(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error) {
				assert.Equal(t, &workspaceID, arg.WorkspaceID)
				return []store.LoadPlan{
					{PlanID: uuid.New(), PlanCode: "P1", CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true}, Status: stringPtr("DRAFT")},
				}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		ctx := authedFounderCtxWithWorkspaceOverride(workspaceID)
		resp, err := s.ListPlans(ctx, 1, 10)

		assert.NoError(t, err)
		assert.Len(t, resp, 1)
	})

	t.Run("trial_guest_db_error", func(t *testing.T) {
		guestID := uuid.New()
		mockQ := &MockQuerier{
			ListLoadPlansForGuestFunc: func(ctx context.Context, arg store.ListLoadPlansForGuestParams) ([]store.LoadPlan, error) {
				return nil, fmt.Errorf("guest plans db error")
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.ListPlans(authedTrialCtxWithID(guestID), 1, 10)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("pagination_defaults", func(t *testing.T) {
		mockQ := &MockQuerier{
			ListLoadPlansFunc: func(ctx context.Context, arg store.ListLoadPlansParams) ([]store.LoadPlan, error) {
				assert.Equal(t, int32(10), arg.Limit)
				assert.Equal(t, int32(0), arg.Offset)
				return []store.LoadPlan{}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		resp, err := s.ListPlans(authedPlannerCtx(), 0, 0)

		assert.NoError(t, err)
		assert.Len(t, resp, 0)
	})
}

func TestPlanService_UpdatePlan(t *testing.T) {
	planID := uuid.New()
	contID := uuid.New()
	statusNew := types.PlanStatusCompleted.String()

	t.Run("success_update_status", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:      planID,
					PlanCode:    "OLD_CODE",
					Status:      stringPtr(types.PlanStatusDraft.String()),
					ContLabel:   stringPtr("Old Cont"),
					LengthMm:    toNumeric(1000),
					WidthMm:     toNumeric(1000),
					HeightMm:    toNumeric(1000),
					MaxWeightKg: toNumeric(1000),
				}, nil
			},
			UpdateLoadPlanFunc: func(ctx context.Context, arg store.UpdateLoadPlanParams) error {
				assert.Equal(t, planID, arg.PlanID)
				assert.Equal(t, "OLD_CODE", arg.PlanCode)
				assert.Equal(t, statusNew, *arg.Status)
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{
			Status: stringPtr(statusNew),
		}
		err := s.UpdatePlan(authedPlannerCtx(), planID.String(), req)

		assert.NoError(t, err)
	})

	t.Run("success_update_container_preset", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:      planID,
					PlanCode:    "OLD_CODE",
					Status:      stringPtr(types.PlanStatusDraft.String()),
					ContLabel:   stringPtr("Old Cont"),
					LengthMm:    toNumeric(1000),
					WidthMm:     toNumeric(1000),
					HeightMm:    toNumeric(1000),
					MaxWeightKg: toNumeric(1000),
				}, nil
			},
			GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
				return store.Container{
					ContainerID:   contID,
					Name:          "New Cont",
					InnerLengthMm: toNumeric(2000),
					InnerWidthMm:  toNumeric(2000),
					InnerHeightMm: toNumeric(2000),
					MaxWeightKg:   toNumeric(2000),
				}, nil
			},
			UpdateLoadPlanFunc: func(ctx context.Context, arg store.UpdateLoadPlanParams) error {
				assert.Equal(t, "New Cont", *arg.ContLabel)
				assert.Equal(t, toNumeric(2000), arg.LengthMm)
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{
			Container: &dto.CreatePlanContainer{
				ContainerID: stringPtr(contID.String()),
			},
		}
		err := s.UpdatePlan(authedPlannerCtx(), planID.String(), req)
		assert.NoError(t, err)
	})

	t.Run("error_plan_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{}, fmt.Errorf("plan not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{
			Status: stringPtr(statusNew),
		}
		err := s.UpdatePlan(authedPlannerCtx(), uuid.New().String(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan not found")
	})

	t.Run("error_invalid_plan_id", func(t *testing.T) {
		mockQ := &MockQuerier{}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{Status: stringPtr(statusNew)}
		err := s.UpdatePlan(authedPlannerCtx(), "invalid-uuid", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid plan id")
	})

	t.Run("error_invalid_container_id", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, PlanCode: "CODE"}, nil
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{
			Container: &dto.CreatePlanContainer{
				ContainerID: stringPtr("invalid-uuid"),
			},
		}
		err := s.UpdatePlan(authedPlannerCtx(), planID.String(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid container_id format")
	})

	t.Run("error_container_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, PlanCode: "CODE"}, nil
			},
			GetContainerFunc: func(ctx context.Context, arg store.GetContainerParams) (store.Container, error) {
				return store.Container{}, fmt.Errorf("container not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{
			Container: &dto.CreatePlanContainer{
				ContainerID: stringPtr(contID.String()),
			},
		}
		err := s.UpdatePlan(authedPlannerCtx(), planID.String(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "container not found")
	})

	t.Run("success_update_custom_dimensions", func(t *testing.T) {
		lengthMM := 1500.0
		widthMM := 1600.0
		heightMM := 1700.0
		maxWeightKG := 1800.0

		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:   planID,
					PlanCode: "CODE",
				}, nil
			},
			UpdateLoadPlanFunc: func(ctx context.Context, arg store.UpdateLoadPlanParams) error {
				assert.Equal(t, toNumeric(lengthMM), arg.LengthMm)
				assert.Equal(t, toNumeric(widthMM), arg.WidthMm)
				assert.Equal(t, toNumeric(heightMM), arg.HeightMm)
				assert.Equal(t, toNumeric(maxWeightKG), arg.MaxWeightKg)
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{
			Container: &dto.CreatePlanContainer{
				LengthMM:    &lengthMM,
				WidthMM:     &widthMM,
				HeightMM:    &heightMM,
				MaxWeightKG: &maxWeightKG,
			},
		}
		err := s.UpdatePlan(authedPlannerCtx(), planID.String(), req)
		assert.NoError(t, err)
	})

	t.Run("error_update_database_error", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, PlanCode: "CODE"}, nil
			},
			UpdateLoadPlanFunc: func(ctx context.Context, arg store.UpdateLoadPlanParams) error {
				return fmt.Errorf("database update error")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanRequest{Status: stringPtr(statusNew)}
		err := s.UpdatePlan(authedPlannerCtx(), planID.String(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database update error")
	})
}

func TestPlanService_DeletePlan(t *testing.T) {
	planID := uuid.New()
	t.Run("success", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			DeleteLoadPlanFunc: func(ctx context.Context, arg store.DeleteLoadPlanParams) error {
				if arg.PlanID != planID {
					return assert.AnError
				}
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		err := s.DeletePlan(authedPlannerCtx(), planID.String())

		assert.NoError(t, err)
	})

	t.Run("trial_forbidden_when_plan_not_owned", func(t *testing.T) {
		guestID := uuid.New()
		deleteCalled := false

		mockQ := &MockQuerier{
			GetLoadPlanForGuestFunc: func(ctx context.Context, arg store.GetLoadPlanForGuestParams) (store.LoadPlan, error) {
				return store.LoadPlan{}, fmt.Errorf("not found")
			},
			DeleteLoadPlanFunc: func(ctx context.Context, arg store.DeleteLoadPlanParams) error {
				deleteCalled = true
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		err := s.DeletePlan(authedTrialCtxWithID(guestID), planID.String())

		assert.ErrorIs(t, err, service.ErrForbidden)
		assert.False(t, deleteCalled)
	})
}

func TestPlanService_AddPlanItem(t *testing.T) {
	planID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{
					ItemID:        uuid.New(),
					ItemLabel:     arg.ItemLabel,
					LengthMm:      arg.LengthMm,
					WidthMm:       arg.WidthMm,
					HeightMm:      arg.HeightMm,
					WeightKg:      arg.WeightKg,
					Quantity:      arg.Quantity,
					AllowRotation: arg.AllowRotation,
					ColorHex:      arg.ColorHex,
				}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.AddPlanItemRequest{}
		req.Label = stringPtr("Item")
		req.LengthMM = 100
		req.WidthMM = 100
		req.HeightMM = 100
		req.WeightKG = 10
		req.Quantity = 5
		req.AllowRotation = boolPtr(true)
		req.ColorHex = stringPtr("#abcdef")

		resp, err := s.AddPlanItem(authedPlannerCtx(), planID.String(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Item", *resp.Label)
		assert.Equal(t, 100.0, resp.LengthMM)
		assert.True(t, resp.AllowRotation)
	})

	t.Run("trial_forbidden_when_plan_not_owned", func(t *testing.T) {
		guestID := uuid.New()
		addCalled := false

		mockQ := &MockQuerier{
			GetLoadPlanForGuestFunc: func(ctx context.Context, arg store.GetLoadPlanForGuestParams) (store.LoadPlan, error) {
				return store.LoadPlan{}, fmt.Errorf("not found")
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				addCalled = true
				return store.LoadItem{}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.AddPlanItemRequest{}
		req.Label = stringPtr("Item")
		req.LengthMM = 1
		req.WidthMM = 1
		req.HeightMM = 1
		req.WeightKG = 1
		req.Quantity = 1

		resp, err := s.AddPlanItem(authedTrialCtxWithID(guestID), planID.String(), req)
		assert.ErrorIs(t, err, service.ErrForbidden)
		assert.Nil(t, resp)
		assert.False(t, addCalled)
	})

	t.Run("error_add_item_db_error", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			AddLoadItemFunc: func(ctx context.Context, arg store.AddLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{}, fmt.Errorf("db error")
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.AddPlanItemRequest{}
		req.Label = stringPtr("Item")
		req.LengthMM = 100
		req.WidthMM = 100
		req.HeightMM = 100
		req.WeightKG = 10
		req.Quantity = 5

		resp, err := s.AddPlanItem(authedPlannerCtx(), planID.String(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to add item")
	})
}

func TestPlanService_UpdatePlanItem(t *testing.T) {
	planID := uuid.New()
	itemID := uuid.New()

	initialItem := store.LoadItem{
		ItemID:        itemID,
		PlanID:        &planID,
		ItemLabel:     stringPtr("Old Label"),
		LengthMm:      toNumeric(500),
		WidthMm:       toNumeric(300),
		HeightMm:      toNumeric(200),
		WeightKg:      toNumeric(5),
		Quantity:      int32(10),
		AllowRotation: boolPtr(true),
		ColorHex:      stringPtr("#ff0000"),
	}

	t.Run("success_update_label", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			GetLoadItemFunc: func(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error) {
				return initialItem, nil
			},
			UpdateLoadItemFunc: func(ctx context.Context, arg store.UpdateLoadItemParams) error {
				assert.Equal(t, planID, *arg.PlanID)
				assert.Equal(t, itemID, arg.ItemID)
				assert.Equal(t, stringPtr("New Label"), arg.ItemLabel)
				assert.Equal(t, initialItem.LengthMm, arg.LengthMm)
				// ... assert other unchanged fields
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanItemRequest{
			Label: stringPtr("New Label"),
		}

		err := s.UpdatePlanItem(authedPlannerCtx(), planID.String(), itemID.String(), req)
		assert.NoError(t, err)
	})

	t.Run("success_update_all_fields", func(t *testing.T) {
		newLen := 600.0
		newQty := 20
		newColor := "#0000ff"
		newRotation := false

		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			GetLoadItemFunc: func(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error) {
				return initialItem, nil
			},
			UpdateLoadItemFunc: func(ctx context.Context, arg store.UpdateLoadItemParams) error {
				assert.Equal(t, toNumeric(newLen), arg.LengthMm)
				assert.Equal(t, int32(newQty), arg.Quantity)
				assert.Equal(t, &newRotation, arg.AllowRotation)
				assert.Equal(t, &newColor, arg.ColorHex)
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanItemRequest{
			LengthMM:      &newLen,
			Quantity:      &newQty,
			AllowRotation: &newRotation,
			ColorHex:      &newColor,
		}
		err := s.UpdatePlanItem(authedPlannerCtx(), planID.String(), itemID.String(), req)
		assert.NoError(t, err)
	})

	t.Run("error_plan_item_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			GetLoadItemFunc: func(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{}, fmt.Errorf("item not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		req := dto.UpdatePlanItemRequest{Label: stringPtr("New Label")}
		err := s.UpdatePlanItem(authedPlannerCtx(), planID.String(), uuid.New().String(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "item not found")
	})
}

func TestPlanService_DeletePlanItem(t *testing.T) {
	planID := uuid.New()
	itemID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			DeleteLoadItemFunc: func(ctx context.Context, arg store.DeleteLoadItemParams) error {
				assert.Equal(t, planID, *arg.PlanID)
				assert.Equal(t, itemID, arg.ItemID)
				return nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		err := s.DeletePlanItem(authedPlannerCtx(), planID.String(), itemID.String())
		assert.NoError(t, err)
	})

	t.Run("error_db_delete_fails", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
			DeleteLoadItemFunc: func(ctx context.Context, arg store.DeleteLoadItemParams) error {
				return fmt.Errorf("db error")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		err := s.DeletePlanItem(authedPlannerCtx(), planID.String(), itemID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete item")
	})

	t.Run("invalid_plan_id", func(t *testing.T) {
		mockQ := &MockQuerier{}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		err := s.DeletePlanItem(authedPlannerCtx(), "invalid-uuid", itemID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid plan id")
	})

	t.Run("invalid_item_id", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{PlanID: planID, WorkspaceID: arg.WorkspaceID}, nil
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		err := s.DeletePlanItem(authedPlannerCtx(), planID.String(), "invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid item id")
	})

	t.Run("plan_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{}, fmt.Errorf("plan not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		err := s.DeletePlanItem(authedPlannerCtx(), planID.String(), itemID.String())
		assert.Error(t, err)
	})
}
func TestPlanService_GetPlanItem(t *testing.T) {
	planID := uuid.New()
	itemID := uuid.New()
	workspaceID := uuid.New()
	label := "Test Item"
	colorHex := "#FF5733"

	t.Run("success", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:      planID,
					WorkspaceID: &workspaceID,
				}, nil
			},
			GetLoadItemFunc: func(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error) {
				assert.Equal(t, planID, *arg.PlanID)
				assert.Equal(t, itemID, arg.ItemID)
				return store.LoadItem{
					ItemID:        itemID,
					PlanID:        &planID,
					ItemLabel:     &label,
					LengthMm:      toNumeric(1000),
					WidthMm:       toNumeric(500),
					HeightMm:      toNumeric(800),
					WeightKg:      toNumeric(25.5),
					Quantity:      5,
					AllowRotation: boolPtr(true),
					ColorHex:      &colorHex,
				}, nil
			},
		}

		s := service.NewPlanService(mockQ, packer.NewPacker())
		result, err := s.GetPlanItem(authedPlannerCtx(), planID.String(), itemID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, itemID.String(), result.ItemID)
		assert.Equal(t, &label, result.Label)
		assert.Equal(t, 1000.0, result.LengthMM)
		assert.Equal(t, 500.0, result.WidthMM)
		assert.Equal(t, 800.0, result.HeightMM)
		assert.Equal(t, 25.5, result.WeightKG)
		assert.Equal(t, 5, result.Quantity)
		assert.Equal(t, &colorHex, result.ColorHex)
	})

	t.Run("error_invalid_plan_id", func(t *testing.T) {
		mockQ := &MockQuerier{}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		result, err := s.GetPlanItem(authedPlannerCtx(), "invalid-uuid", itemID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid plan id")
	})

	t.Run("error_invalid_item_id", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:      planID,
					WorkspaceID: &workspaceID,
				}, nil
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		result, err := s.GetPlanItem(authedPlannerCtx(), planID.String(), "invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid item id")
	})

	t.Run("error_plan_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{}, fmt.Errorf("plan not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		result, err := s.GetPlanItem(authedPlannerCtx(), planID.String(), itemID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("error_item_not_found", func(t *testing.T) {
		mockQ := &MockQuerier{
			GetLoadPlanFunc: func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
				return store.LoadPlan{
					PlanID:      planID,
					WorkspaceID: &workspaceID,
				}, nil
			},
			GetLoadItemFunc: func(ctx context.Context, arg store.GetLoadItemParams) (store.LoadItem, error) {
				return store.LoadItem{}, fmt.Errorf("item not found")
			},
		}
		s := service.NewPlanService(mockQ, packer.NewPacker())
		result, err := s.GetPlanItem(authedPlannerCtx(), planID.String(), itemID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// MockPacker is a mock implementation of the packer.Packer interface
type MockPacker struct {
	PackFunc func(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error)
}

func (m *MockPacker) Pack(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
	if m.PackFunc != nil {
		return m.PackFunc(ctx, container, items)
	}
	return packer.PackingResult{}, fmt.Errorf("PackFunc not implemented")
}

func TestPlanService_CalculatePlan(t *testing.T) {
	planID := uuid.New()
	workspaceID := uuid.New()
	itemID1 := uuid.New()
	resultID := uuid.New()

	tests := []struct {
		name       string
		planID     string
		opts       dto.CalculatePlanRequest
		ctx        context.Context
		mockSetup  func(*MockQuerier, *MockPacker)
		wantErr    bool
		assertFunc func(*testing.T, *dto.CalculationResult, error)
	}{
		{
			name:   "successful_calculation_feasible",
			planID: planID.String(),
			opts: dto.CalculatePlanRequest{
				Strategy: "best_fit",
				Goal:     "volume",
				Gravity:  boolPtr(false),
			},
			ctx: authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{
						PlanID:      planID,
						WorkspaceID: &workspaceID,
						LengthMm:    toNumeric(1000.0),
						WidthMm:     toNumeric(1000.0),
						HeightMm:    toNumeric(1000.0),
						MaxWeightKg: toNumeric(100.0),
					}, nil
				}
				mq.ListLoadItemsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) ([]store.LoadItem, error) {
					return []store.LoadItem{
						{
							ItemID:   itemID1,
							LengthMm: toNumeric(100.0),
							WidthMm:  toNumeric(100.0),
							HeightMm: toNumeric(100.0),
							WeightKg: toNumeric(10.0),
							Quantity: 2,
						},
					}, nil
				}
				mp.PackFunc = func(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
					return packer.PackingResult{
						IsFeasible:           true,
						TotalWeightPackedKG:  20.0,
						VolumeUtilisationPct: 15.5,
						Algorithm:            "test-algorithm",
						DurationMs:           100,
						PackedItems: []packer.PackedItem{
							{
								ItemID:       itemID1.String(),
								RotationType: 0,
								Position:     packer.Position{X: 0, Y: 0, Z: 0},
							},
							{
								ItemID:       itemID1.String(),
								RotationType: 0,
								Position:     packer.Position{X: 100, Y: 0, Z: 0},
							},
						},
						UnfitItems: []packer.ItemInput{}, // Empty for feasible case
					}, nil
				}
				mq.DeletePlanResultsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) error {
					return nil
				}
				mq.CreatePlanResultFunc = func(ctx context.Context, arg store.CreatePlanResultParams) (store.PlanResult, error) {
					return store.PlanResult{
						ResultID:             resultID,
						PlanID:               arg.PlanID,
						TotalLoadedWeightKg:  arg.TotalLoadedWeightKg,
						VolumeUtilizationPct: arg.VolumeUtilizationPct,
						IsFeasible:           arg.IsFeasible,
					}, nil
				}
				mq.CreatePlanPlacementFunc = func(ctx context.Context, arg []store.CreatePlanPlacementParams) (int64, error) {
					return int64(len(arg)), nil
				}
				mq.UpdatePlanStatusFunc = func(ctx context.Context, arg store.UpdatePlanStatusParams) error {
					assert.Equal(t, types.PlanStatusCompleted.String(), *arg.Status)
					return nil
				}
			},
			wantErr: false,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, resultID.String(), result.JobID)
				assert.Equal(t, types.PlanStatusCompleted.String(), result.Status)
				assert.Equal(t, "test-algorithm", result.Algorithm)
				assert.Equal(t, 15.5, result.VolumeUtilization)
				assert.Equal(t, 15.5, result.EfficiencyScore)
				assert.Equal(t, int64(100), result.DurationMs)
				assert.Equal(t, 2, len(result.Placements))
			},
		},
		{
			name:   "successful_calculation_partial",
			planID: planID.String(),
			opts: dto.CalculatePlanRequest{
				Strategy: "best_fit",
				Goal:     "volume",
			},
			ctx: authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{
						PlanID:      planID,
						WorkspaceID: &workspaceID,
						LengthMm:    toNumeric(1000.0),
						WidthMm:     toNumeric(1000.0),
						HeightMm:    toNumeric(1000.0),
						MaxWeightKg: toNumeric(100.0),
					}, nil
				}
				mq.ListLoadItemsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) ([]store.LoadItem, error) {
					return []store.LoadItem{
						{
							ItemID:   itemID1,
							LengthMm: toNumeric(100.0),
							WidthMm:  toNumeric(100.0),
							HeightMm: toNumeric(100.0),
							WeightKg: toNumeric(10.0),
							Quantity: 1,
						},
					}, nil
				}
				mp.PackFunc = func(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
					return packer.PackingResult{
						IsFeasible:           false,
						TotalWeightPackedKG:  10.0,
						VolumeUtilisationPct: 5.0,
						Algorithm:            "test-algorithm",
						PackedItems: []packer.PackedItem{
							{
								ItemID:       itemID1.String(),
								RotationType: 0,
								Position:     packer.Position{X: 0, Y: 0, Z: 0},
							},
						},
						UnfitItems: []packer.ItemInput{
							{ID: "unfit1"},
							{ID: "unfit2"},
							{ID: "unfit3"},
						}, // 3 unfit items for infeasible case
					}, nil
				}
				mq.DeletePlanResultsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) error {
					return nil
				}
				mq.CreatePlanResultFunc = func(ctx context.Context, arg store.CreatePlanResultParams) (store.PlanResult, error) {
					return store.PlanResult{
						ResultID:             resultID,
						PlanID:               arg.PlanID,
						TotalLoadedWeightKg:  arg.TotalLoadedWeightKg,
						VolumeUtilizationPct: arg.VolumeUtilizationPct,
						IsFeasible:           arg.IsFeasible,
					}, nil
				}
				mq.CreatePlanPlacementFunc = func(ctx context.Context, arg []store.CreatePlanPlacementParams) (int64, error) {
					return int64(len(arg)), nil
				}
				mq.UpdatePlanStatusFunc = func(ctx context.Context, arg store.UpdatePlanStatusParams) error {
					assert.Equal(t, types.PlanStatusPartial.String(), *arg.Status)
					return nil
				}
			},
			wantErr: false,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// When packing is infeasible, service returns completed status (not partial) in the DTO
				assert.Equal(t, types.PlanStatusCompleted.String(), result.Status)
				// Verify only 1 placement was created (from the 1 packed item)
				assert.Equal(t, 1, len(result.Placements))
			},
		},
		{
			name:      "invalid_plan_id",
			planID:    "invalid-uuid",
			opts:      dto.CalculatePlanRequest{},
			ctx:       authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {},
			wantErr:   true,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "invalid plan id")
			},
		},
		{
			name:   "plan_not_found",
			planID: planID.String(),
			opts:   dto.CalculatePlanRequest{},
			ctx:    authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{}, fmt.Errorf("plan not found")
				}
			},
			wantErr: true,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "plan not found")
			},
		},
		{
			name:   "database_error_list_items",
			planID: planID.String(),
			opts:   dto.CalculatePlanRequest{},
			ctx:    authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{
						PlanID:      planID,
						WorkspaceID: &workspaceID,
						LengthMm:    toNumeric(1000.0),
						WidthMm:     toNumeric(1000.0),
						HeightMm:    toNumeric(1000.0),
						MaxWeightKg: toNumeric(100.0),
					}, nil
				}
				mq.ListLoadItemsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) ([]store.LoadItem, error) {
					return nil, fmt.Errorf("database error")
				}
			},
			wantErr: true,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "failed to list items")
			},
		},
		{
			name:   "packing_algorithm_error",
			planID: planID.String(),
			opts:   dto.CalculatePlanRequest{},
			ctx:    authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{
						PlanID:      planID,
						WorkspaceID: &workspaceID,
						LengthMm:    toNumeric(1000.0),
						WidthMm:     toNumeric(1000.0),
						HeightMm:    toNumeric(1000.0),
						MaxWeightKg: toNumeric(100.0),
					}, nil
				}
				mq.ListLoadItemsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) ([]store.LoadItem, error) {
					return []store.LoadItem{{ItemID: itemID1}}, nil
				}
				mp.PackFunc = func(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
					return packer.PackingResult{}, fmt.Errorf("packing algorithm error")
				}
			},
			wantErr: true,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "packing failed")
			},
		},
		{
			name:   "database_error_save_result",
			planID: planID.String(),
			opts:   dto.CalculatePlanRequest{},
			ctx:    authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{
						PlanID:      planID,
						WorkspaceID: &workspaceID,
						LengthMm:    toNumeric(1000.0),
						WidthMm:     toNumeric(1000.0),
						HeightMm:    toNumeric(1000.0),
						MaxWeightKg: toNumeric(100.0),
					}, nil
				}
				mq.ListLoadItemsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) ([]store.LoadItem, error) {
					return []store.LoadItem{{ItemID: itemID1}}, nil
				}
				mp.PackFunc = func(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
					return packer.PackingResult{
						IsFeasible:  true,
						PackedItems: []packer.PackedItem{},
					}, nil
				}
				mq.DeletePlanResultsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) error {
					return nil
				}
				mq.CreatePlanResultFunc = func(ctx context.Context, arg store.CreatePlanResultParams) (store.PlanResult, error) {
					return store.PlanResult{}, fmt.Errorf("database error")
				}
			},
			wantErr: true,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "failed to save result")
			},
		},
		{
			name:   "database_error_save_placements",
			planID: planID.String(),
			opts:   dto.CalculatePlanRequest{},
			ctx:    authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{
						PlanID:      planID,
						WorkspaceID: &workspaceID,
						LengthMm:    toNumeric(1000.0),
						WidthMm:     toNumeric(1000.0),
						HeightMm:    toNumeric(1000.0),
						MaxWeightKg: toNumeric(100.0),
					}, nil
				}
				mq.ListLoadItemsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) ([]store.LoadItem, error) {
					return []store.LoadItem{{ItemID: itemID1}}, nil
				}
				mp.PackFunc = func(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
					return packer.PackingResult{
						IsFeasible: true,
						PackedItems: []packer.PackedItem{
							{
								ItemID:       itemID1.String(),
								RotationType: 0,
								Position:     packer.Position{X: 0, Y: 0, Z: 0},
							},
						},
					}, nil
				}
				mq.DeletePlanResultsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) error {
					return nil
				}
				mq.CreatePlanResultFunc = func(ctx context.Context, arg store.CreatePlanResultParams) (store.PlanResult, error) {
					return store.PlanResult{
						ResultID: resultID,
						PlanID:   arg.PlanID,
					}, nil
				}
				mq.CreatePlanPlacementFunc = func(ctx context.Context, arg []store.CreatePlanPlacementParams) (int64, error) {
					return 0, fmt.Errorf("placement insert error")
				}
			},
			wantErr: true,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "failed to save placements")
			},
		},
		{
			name:   "with_gravity_option",
			planID: planID.String(),
			opts: dto.CalculatePlanRequest{
				Strategy: "best_fit",
				Goal:     "volume",
				Gravity:  boolPtr(true),
			},
			ctx: authedPlannerCtx(),
			mockSetup: func(mq *MockQuerier, mp *MockPacker) {
				mq.GetLoadPlanFunc = func(ctx context.Context, arg store.GetLoadPlanParams) (store.LoadPlan, error) {
					return store.LoadPlan{
						PlanID:      planID,
						WorkspaceID: &workspaceID,
						LengthMm:    toNumeric(1000.0),
						WidthMm:     toNumeric(1000.0),
						HeightMm:    toNumeric(1000.0),
						MaxWeightKg: toNumeric(100.0),
					}, nil
				}
				mq.ListLoadItemsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) ([]store.LoadItem, error) {
					return []store.LoadItem{{ItemID: itemID1}}, nil
				}
				mp.PackFunc = func(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
					// Verify gravity option was passed
					assert.True(t, container.Options.Gravity)
					return packer.PackingResult{
						IsFeasible:  true,
						PackedItems: []packer.PackedItem{},
					}, nil
				}
				mq.DeletePlanResultsFunc = func(ctx context.Context, planIDPtr *uuid.UUID) error {
					return nil
				}
				mq.CreatePlanResultFunc = func(ctx context.Context, arg store.CreatePlanResultParams) (store.PlanResult, error) {
					return store.PlanResult{ResultID: resultID, PlanID: arg.PlanID}, nil
				}
				mq.CreatePlanPlacementFunc = func(ctx context.Context, arg []store.CreatePlanPlacementParams) (int64, error) {
					return 0, nil
				}
				mq.UpdatePlanStatusFunc = func(ctx context.Context, arg store.UpdatePlanStatusParams) error {
					return nil
				}
			},
			wantErr: false,
			assertFunc: func(t *testing.T, result *dto.CalculationResult, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{}
			mockP := &MockPacker{}
			tt.mockSetup(mockQ, mockP)

			s := service.NewPlanService(mockQ, mockP)
			result, err := s.CalculatePlan(tt.ctx, tt.planID, tt.opts)

			tt.assertFunc(t, result, err)
		})
	}
}
