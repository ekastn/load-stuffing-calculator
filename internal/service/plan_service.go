package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/packer"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/ekastn/load-stuffing-calculator/internal/types"
	"github.com/google/uuid"
)

type PlanService interface {
	CreateCompletePlan(ctx context.Context, req dto.CreatePlanRequest) (*dto.CreatePlanResponse, error)
	GetPlan(ctx context.Context, id string) (*dto.PlanDetailResponse, error)
	ListPlans(ctx context.Context, page, limit int32) ([]dto.PlanListItem, error)
	UpdatePlan(ctx context.Context, id string, req dto.UpdatePlanRequest) error
	DeletePlan(ctx context.Context, id string) error
	AddPlanItem(ctx context.Context, planID string, req dto.AddPlanItemRequest) (*dto.PlanItemDetail, error)
	GetPlanItem(ctx context.Context, planID, itemID string) (*dto.PlanItemDetail, error)
	UpdatePlanItem(ctx context.Context, planID, itemID string, req dto.UpdatePlanItemRequest) error
	DeletePlanItem(ctx context.Context, planID, itemID string) error
	CalculatePlan(ctx context.Context, planID string) (*dto.CalculationResult, error)
}

type planService struct {
	q store.Querier
	p packer.Packer
}

func NewPlanService(q store.Querier, p packer.Packer) PlanService {
	return &planService{q: q, p: p}
}

func (s *planService) CreateCompletePlan(ctx context.Context, req dto.CreatePlanRequest) (*dto.CreatePlanResponse, error) {
	autoCalc := true
	if req.AutoCalculate != nil {
		autoCalc = *req.AutoCalculate
	}

	var lengthMM, widthMM, heightMM, maxWeightKG float64
	var contLabel string = "Custom Container"

	if req.Container.ContainerID != nil {
		contUUID, err := uuid.Parse(*req.Container.ContainerID)
		if err != nil {
			return nil, fmt.Errorf("invalid container_id format")
		}
		cont, err := s.q.GetContainer(ctx, contUUID)
		if err != nil {
			return nil, fmt.Errorf("container not found: %w", err)
		}
		lengthMM = toFloat(cont.InnerLengthMm)
		widthMM = toFloat(cont.InnerWidthMm)
		heightMM = toFloat(cont.InnerHeightMm)
		maxWeightKG = toFloat(cont.MaxWeightKg)
		contLabel = cont.Name
	} else {
		if req.Container.LengthMM == nil || req.Container.WidthMM == nil ||
			req.Container.HeightMM == nil || req.Container.MaxWeightKG == nil {
			return nil, fmt.Errorf("custom container dimensions are required")
		}
		lengthMM = *req.Container.LengthMM
		widthMM = *req.Container.WidthMM
		heightMM = *req.Container.HeightMM
		maxWeightKG = *req.Container.MaxWeightKG
	}

	planCode := "PLD-" + time.Now().Format("20060102-150405")
	status := types.PlanStatusDraft.String()

	plan, err := s.q.CreateLoadPlan(ctx, store.CreateLoadPlanParams{
		PlanCode:    planCode,
		Status:      &status,
		ContLabel:   &contLabel,
		LengthMm:    toNumeric(lengthMM),
		WidthMm:     toNumeric(widthMM),
		HeightMm:    toNumeric(heightMM),
		MaxWeightKg: toNumeric(maxWeightKG),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	var totalQty int
	var totalWeight, totalVolume float64

	for _, item := range req.Items {
		allowRot := true
		if item.AllowRotation != nil {
			allowRot = *item.AllowRotation
		}
		color := "#3498db"
		if item.ColorHex != nil {
			color = *item.ColorHex
		}

		volumePerItem := item.LengthMM * item.WidthMM * item.HeightMM / 1_000_000_000.0
		totalVolItem := volumePerItem * float64(item.Quantity)
		totalWItem := item.WeightKG * float64(item.Quantity)

		_, err := s.q.AddLoadItem(ctx, store.AddLoadItemParams{
			PlanID:        &plan.PlanID,
			ItemLabel:     item.Label,
			LengthMm:      toNumeric(item.LengthMM),
			WidthMm:       toNumeric(item.WidthMM),
			HeightMm:      toNumeric(item.HeightMM),
			WeightKg:      toNumeric(item.WeightKG),
			Quantity:      int32(item.Quantity),
			AllowRotation: &allowRot,
			ColorHex:      &color,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add item: %w", err)
		}

		totalQty += item.Quantity
		totalWeight += totalWItem
		totalVolume += totalVolItem
	}

	var jobID *string
	if autoCalc {
		// Perform synchronous calculation
		calcRes, err := s.CalculatePlan(ctx, plan.PlanID.String())
		if err != nil {
			// If calculation fails, we log it (conceptually) and mark status as FAILED
			failStatus := types.PlanStatusFailed.String()
			_ = s.q.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
				PlanID: plan.PlanID,
				Status: &failStatus,
			})
			status = types.PlanStatusFailed.String()
			// We can return the error, or just continue with failed status.
			// Returning error might confuse client if Plan was created.
			// Let's just set status FAILED.
		} else {
			status = types.PlanStatusCompleted.String()
			if !strings.EqualFold(calcRes.Status, types.PlanStatusCompleted.String()) {
				status = types.PlanStatusPartial.String() // logic inside CalculatePlan handles DB update, here for response
			}
			jobID = &calcRes.JobID
		}
	} else {
		// If not auto-calc, leave as DRAFT
		status = types.PlanStatusDraft.String()
	}

	return &dto.CreatePlanResponse{
		PlanID:           plan.PlanID.String(),
		PlanCode:         plan.PlanCode,
		Title:            req.Title,
		Status:           status,
		TotalItems:       totalQty,
		TotalWeightKG:    totalWeight,
		TotalVolumeM3:    totalVolume,
		CalculationJobID: jobID,
		CreatedAt:        plan.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (s *planService) GetPlan(ctx context.Context, id string) (*dto.PlanDetailResponse, error) {
	planUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id")
	}

	plan, err := s.q.GetLoadPlan(ctx, planUUID)
	if err != nil {
		return nil, err
	}

	items, err := s.q.ListLoadItems(ctx, &plan.PlanID)
	if err != nil {
		return nil, err
	}

	var totalQty int
	var totalWeight, totalVolume float64
	var itemDetails []dto.PlanItemDetail

	for _, i := range items {
		l := toFloat(i.LengthMm)
		w := toFloat(i.WidthMm)
		h := toFloat(i.HeightMm)
		wg := toFloat(i.WeightKg)
		q := int(i.Quantity)

		vol := l * w * h / 1_000_000_000.0 * float64(q)
		tw := wg * float64(q)

		totalQty += q
		totalWeight += tw
		totalVolume += vol

		itemDetails = append(itemDetails, dto.PlanItemDetail{
			ItemID:        i.ItemID.String(),
			Label:         i.ItemLabel,
			LengthMM:      l,
			WidthMM:       w,
			HeightMM:      h,
			WeightKG:      wg,
			Quantity:      q,
			TotalWeightKG: tw,
			TotalVolumeM3: vol,
			AllowRotation: *i.AllowRotation,
			ColorHex:      i.ColorHex,
			CreatedAt:     "", // DB doesn't have created_at for item
		})
	}

	contL := toFloat(plan.LengthMm)
	contW := toFloat(plan.WidthMm)
	contH := toFloat(plan.HeightMm)
	contVol := contL * contW * contH / 1_000_000_000.0

	var calc *dto.CalculationResult
	res, err := s.q.GetPlanResult(ctx, &plan.PlanID)
	if err == nil {
		status := types.PlanStatusCompleted.String()
		if res.IsFeasible != nil && !*res.IsFeasible {
			status = types.PlanStatusPartial.String()
		}

		calc = &dto.CalculationResult{
			JobID:             res.ResultID.String(),
			Status:            status,
			Algorithm:         "BestFitDecreasing", // Default for now
			EfficiencyScore:   toFloat(res.VolumeUtilizationPct),
			VolumeUtilization: toFloat(res.VolumeUtilizationPct),
			VisualizationURL:  "/visualizer?plan=" + plan.PlanID.String(),
		}

		// Fetch placements
		placements, err := s.q.ListPlanPlacements(ctx, &res.ResultID)
		if err == nil {
			var plDetails []dto.PlacementDetail
			for _, pl := range placements {
				var iID string
				if pl.ItemID != nil {
					iID = pl.ItemID.String()
				}

				rot := 0
				if pl.RotationCode != nil {
					rot = int(*pl.RotationCode)
				}

				plDetails = append(plDetails, dto.PlacementDetail{
					PlacementID: pl.PlacementID.String(),
					ItemID:      iID,
					PositionX:   toFloat(pl.PosX),
					PositionY:   toFloat(pl.PosY),
					PositionZ:   toFloat(pl.PosZ),
					Rotation:    rot,
					StepNumber:  int(pl.StepNumber),
				})
			}
			calc.Placements = plDetails
		}
	}

	return &dto.PlanDetailResponse{
		PlanID:   plan.PlanID.String(),
		PlanCode: plan.PlanCode,
		Status:   getString(plan.Status),
		Container: dto.PlanContainerInfo{
			Name:        plan.ContLabel,
			LengthMM:    contL,
			WidthMM:     contW,
			HeightMM:    contH,
			MaxWeightKG: toFloat(plan.MaxWeightKg),
			VolumeM3:    contVol,
		},
		Stats: dto.PlanStats{
			TotalItems:    totalQty,
			TotalWeightKG: totalWeight,
			TotalVolumeM3: totalVolume,
		},
		Items:       itemDetails,
		Calculation: calc,
		CreatedAt:   plan.CreatedAt.Time,
	}, nil
}

func (s *planService) ListPlans(ctx context.Context, page, limit int32) ([]dto.PlanListItem, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	plans, err := s.q.ListLoadPlans(ctx, store.ListLoadPlansParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var result []dto.PlanListItem
	for _, p := range plans {
		result = append(result, dto.PlanListItem{
			PlanID:    p.PlanID.String(),
			PlanCode:  p.PlanCode,
			Status:    getString(p.Status),
			CreatedAt: p.CreatedAt.Time.Format(time.RFC3339),
			// TODO: Stats are 0 for now as they require aggregation
		})
	}
	return result, nil
}

func (s *planService) UpdatePlan(ctx context.Context, id string, req dto.UpdatePlanRequest) error {
	planUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid plan id")
	}

	plan, err := s.q.GetLoadPlan(ctx, planUUID)
	if err != nil {
		return err
	}

	params := store.UpdateLoadPlanParams{
		PlanID:      planUUID,
		PlanCode:    plan.PlanCode,
		ContLabel:   plan.ContLabel,
		LengthMm:    plan.LengthMm,
		WidthMm:     plan.WidthMm,
		HeightMm:    plan.HeightMm,
		MaxWeightKg: plan.MaxWeightKg,
		Status:      plan.Status,
	}

	if req.Status != nil {
		params.Status = req.Status
	}

	if req.Container != nil {
		if req.Container.ContainerID != nil {
			contUUID, err := uuid.Parse(*req.Container.ContainerID)
			if err != nil {
				return fmt.Errorf("invalid container_id format")
			}
			cont, err := s.q.GetContainer(ctx, contUUID)
			if err != nil {
				return fmt.Errorf("container not found: %w", err)
			}
			params.ContLabel = &cont.Name
			params.LengthMm = cont.InnerLengthMm
			params.WidthMm = cont.InnerWidthMm
			params.HeightMm = cont.InnerHeightMm
			params.MaxWeightKg = cont.MaxWeightKg
		} else {
			if req.Container.LengthMM != nil {
				params.LengthMm = toNumeric(*req.Container.LengthMM)
			}
			if req.Container.WidthMM != nil {
				params.WidthMm = toNumeric(*req.Container.WidthMM)
			}
			if req.Container.HeightMM != nil {
				params.HeightMm = toNumeric(*req.Container.HeightMM)
			}
			if req.Container.MaxWeightKG != nil {
				params.MaxWeightKg = toNumeric(*req.Container.MaxWeightKG)
			}
		}
	}

	return s.q.UpdateLoadPlan(ctx, params)
}

func (s *planService) DeletePlan(ctx context.Context, id string) error {
	planUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid plan id")
	}
	return s.q.DeleteLoadPlan(ctx, planUUID)
}

func (s *planService) AddPlanItem(ctx context.Context, planID string, req dto.AddPlanItemRequest) (*dto.PlanItemDetail, error) {
	pID, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id")
	}

	allowRot := true
	if req.AllowRotation != nil {
		allowRot = *req.AllowRotation
	}
	color := "#3498db"
	if req.ColorHex != nil {
		color = *req.ColorHex
	}

	item, err := s.q.AddLoadItem(ctx, store.AddLoadItemParams{
		PlanID:        &pID,
		ItemLabel:     req.Label,
		LengthMm:      toNumeric(req.LengthMM),
		WidthMm:       toNumeric(req.WidthMM),
		HeightMm:      toNumeric(req.HeightMM),
		WeightKg:      toNumeric(req.WeightKG),
		Quantity:      int32(req.Quantity),
		AllowRotation: &allowRot,
		ColorHex:      &color,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return mapLoadItemToDetail(item), nil
}

func (s *planService) GetPlanItem(ctx context.Context, planID, itemID string) (*dto.PlanItemDetail, error) {
	pID, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id")
	}
	iID, err := uuid.Parse(itemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item id")
	}

	item, err := s.q.GetLoadItem(ctx, store.GetLoadItemParams{PlanID: &pID, ItemID: iID})
	if err != nil {
		return nil, err
	}

	return mapLoadItemToDetail(item), nil
}

func (s *planService) UpdatePlanItem(ctx context.Context, planID, itemID string, req dto.UpdatePlanItemRequest) error {
	pID, err := uuid.Parse(planID)
	if err != nil {
		return fmt.Errorf("invalid plan id")
	}
	iID, err := uuid.Parse(itemID)
	if err != nil {
		return fmt.Errorf("invalid item id")
	}

	// Fetch existing to merge
	existing, err := s.q.GetLoadItem(ctx, store.GetLoadItemParams{PlanID: &pID, ItemID: iID})
	if err != nil {
		return err
	}

	params := store.UpdateLoadItemParams{
		PlanID:        &pID,
		ItemID:        iID,
		ItemLabel:     existing.ItemLabel,
		LengthMm:      existing.LengthMm,
		WidthMm:       existing.WidthMm,
		HeightMm:      existing.HeightMm,
		WeightKg:      existing.WeightKg,
		Quantity:      existing.Quantity,
		AllowRotation: existing.AllowRotation,
		ColorHex:      existing.ColorHex,
	}

	if req.Label != nil {
		params.ItemLabel = req.Label
	}
	if req.LengthMM != nil {
		params.LengthMm = toNumeric(*req.LengthMM)
	}
	if req.WidthMM != nil {
		params.WidthMm = toNumeric(*req.WidthMM)
	}
	if req.HeightMM != nil {
		params.HeightMm = toNumeric(*req.HeightMM)
	}
	if req.WeightKG != nil {
		params.WeightKg = toNumeric(*req.WeightKG)
	}
	if req.Quantity != nil {
		params.Quantity = int32(*req.Quantity)
	}
	if req.AllowRotation != nil {
		params.AllowRotation = req.AllowRotation
	}
	if req.ColorHex != nil {
		params.ColorHex = req.ColorHex
	}

	if err := s.q.UpdateLoadItem(ctx, params); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	return nil
}

func (s *planService) DeletePlanItem(ctx context.Context, planID, itemID string) error {
	pID, err := uuid.Parse(planID)
	if err != nil {
		return fmt.Errorf("invalid plan id")
	}
	iID, err := uuid.Parse(itemID)
	if err != nil {
		return fmt.Errorf("invalid item id")
	}

	if err := s.q.DeleteLoadItem(ctx, store.DeleteLoadItemParams{PlanID: &pID, ItemID: iID}); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (s *planService) CalculatePlan(ctx context.Context, planID string) (*dto.CalculationResult, error) {
	// 1. Fetch Plan & Items
	pID, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id")
	}

	plan, err := s.q.GetLoadPlan(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	items, err := s.q.ListLoadItems(ctx, &plan.PlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	// 2. Prepare Inputs
	contInput := packer.ContainerInput{
		ID:        plan.PlanID.String(),
		Length:    toFloat(plan.LengthMm),
		Width:     toFloat(plan.WidthMm),
		Height:    toFloat(plan.HeightMm),
		MaxWeight: toFloat(plan.MaxWeightKg),
	}

	var itemInputs []packer.ItemInput
	for _, item := range items {
		allowRot := true
		if item.AllowRotation != nil {
			allowRot = *item.AllowRotation
		}
		color := "#3498db"
		if item.ColorHex != nil {
			color = *item.ColorHex
		}

		itemInputs = append(itemInputs, packer.ItemInput{
			ID:            item.ItemID.String(),
			Label:         getString(item.ItemLabel),
			Length:        toFloat(item.LengthMm),
			Width:         toFloat(item.WidthMm),
			Height:        toFloat(item.HeightMm),
			Weight:        toFloat(item.WeightKg),
			Quantity:      int(item.Quantity),
			AllowRotation: allowRot,
			Color:         color,
		})
	}

	// 3. Run Packing
	res, err := s.p.Pack(ctx, contInput, itemInputs)
	if err != nil {
		return nil, fmt.Errorf("packing failed: %w", err)
	}

	// 4. Save Results
	// Delete old results first? Usually only one active result.
	_ = s.q.DeletePlanResults(ctx, &pID)

	savedRes, err := s.q.CreatePlanResult(ctx, store.CreatePlanResultParams{
		PlanID:               &pID,
		TotalLoadedWeightKg:  toNumeric(res.TotalWeightPackedKG),
		VolumeUtilizationPct: toNumeric(res.VolumeUtilisationPct),
		IsFeasible:           &res.IsFeasible,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save result: %w", err)
	}

	// 5. Save Placements (Bulk)
	var placements []store.CreatePlanPlacementParams
	for i, pItem := range res.PackedItems {
		itemID, _ := uuid.Parse(pItem.ItemID)

		rID := savedRes.ResultID
		iID := itemID
		rot := int32(pItem.RotationType)

		placements = append(placements, store.CreatePlanPlacementParams{
			ResultID:     &rID,
			ItemID:       &iID,
			PosX:         toNumeric(pItem.Position.X),
			PosY:         toNumeric(pItem.Position.Y),
			PosZ:         toNumeric(pItem.Position.Z),
			RotationCode: &rot,
			StepNumber:   int32(i + 1),
		})
	}

	if len(placements) > 0 {
		_, err := s.q.CreatePlanPlacement(ctx, placements)
		if err != nil {
			return nil, fmt.Errorf("failed to save placements: %w", err)
		}
	}

	// 6. Update Status
	newStatus := types.PlanStatusCompleted.String()
	if !res.IsFeasible {
		newStatus = types.PlanStatusPartial.String()
	}
	s.q.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
		PlanID: pID,
		Status: &newStatus,
	})

	// 7. Return DTO
	return &dto.CalculationResult{
		JobID:             savedRes.ResultID.String(),
		Status:            types.PlanStatusCompleted.String(), // Status from packer result, or from DB
		Algorithm:         res.Algorithm,
		EfficiencyScore:   res.VolumeUtilisationPct,
		VolumeUtilization: res.VolumeUtilisationPct,
		DurationMs:        res.DurationMs,
		VisualizationURL:  "/visualizer?plan=" + planID,
	}, nil
}

func mapLoadItemToDetail(i store.LoadItem) *dto.PlanItemDetail {
	l := toFloat(i.LengthMm)
	w := toFloat(i.WidthMm)
	h := toFloat(i.HeightMm)
	wg := toFloat(i.WeightKg)
	q := int(i.Quantity)

	vol := l * w * h / 1_000_000_000.0 * float64(q)
	tw := wg * float64(q)

	allowRot := true
	if i.AllowRotation != nil {
		allowRot = *i.AllowRotation
	}

	return &dto.PlanItemDetail{
		ItemID:        i.ItemID.String(),
		Label:         i.ItemLabel,
		LengthMM:      l,
		WidthMM:       w,
		HeightMM:      h,
		WeightKG:      wg,
		Quantity:      q,
		TotalWeightKG: tw,
		TotalVolumeM3: vol,
		AllowRotation: allowRot,
		ColorHex:      i.ColorHex,
	}
}
