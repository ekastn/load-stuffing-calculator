package service

import (
	"context"
	"fmt"

	"load-stuffing-calculator/internal/gateway"
	"load-stuffing-calculator/internal/store"

	"github.com/google/uuid"
)

// PlanService mengelola lifecycle dari Plan:
// - Pembuatan plan baru
// - Penambahan items ke plan
// - Orkestrasi kalkulasi dengan Packing Service
type PlanService struct {
	store   store.Querier          // Dependency: akses database via generated interface
	gateway gateway.PackingGateway // Dependency: komunikasi dengan Packing Service
}

// NewPlanService membuat instance PlanService dengan dependencies yang di-inject.
// Ini memungkinkan substitusi dengan mock saat testing.
func NewPlanService(store store.Querier, gw gateway.PackingGateway) *PlanService {
	return &PlanService{
		store:   store,
		gateway: gw,
	}
}

func (s *PlanService) Create(ctx context.Context, containerID uuid.UUID) (*store.Plan, error) {
	// Validasi: container harus exist
	_, err := s.store.GetContainer(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("container not found: %w", err)
	}

	plan, err := s.store.CreatePlan(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("create plan: %w", err)
	}
	return &plan, nil
}

func (s *PlanService) AddItem(ctx context.Context, planID, productID uuid.UUID, quantity int) error {
	// Validasi
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	_, err := s.store.GetProduct(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	_, err = s.store.AddPlanItem(ctx, store.AddPlanItemParams{
		PlanID:    planID,
		ProductID: productID,
		Quantity:  int32(quantity),
	})
	return err
}

// CalculateResult berisi hasil kalkulasi packing
type CalculateResult struct {
	Placements    []gateway.PackPlacement
	UnfittedItems []gateway.PackUnfitted
	Statistics    gateway.PackStats
}

func (s *PlanService) Calculate(ctx context.Context, planID uuid.UUID) (*CalculateResult, error) {
	// 1. Ambil plan dan container
	plan, err := s.store.GetPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("get plan: %w", err)
	}

	container, err := s.store.GetContainer(ctx, plan.ContainerID)
	if err != nil {
		return nil, fmt.Errorf("get container: %w", err)
	}

	// 2. Ambil items dengan detail produk
	items, err := s.store.GetPlanItems(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("get plan items: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("plan has no items")
	}

	// 3. Update status ke "calculating"
	_, err = s.store.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
		ID:     planID,
		Status: "calculating",
	})
	if err != nil {
		return nil, fmt.Errorf("update status: %w", err)
	}

	// 4. Build request untuk Packing Service
	req := s.buildPackRequest(container, items)

	// 5. Panggil Packing Service
	resp, err := s.gateway.Pack(ctx, req)
	if err != nil {
		// Rollback status ke failed jika gagal
		s.store.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
			ID:     planID,
			Status: "failed",
		})
		return nil, fmt.Errorf("packing calculation failed: %w", err)
	}

	// 6. Simpan placements
	err = s.savePlacements(ctx, planID, resp.Data.Placements)
	if err != nil {
		return nil, fmt.Errorf("save placements: %w", err)
	}

	// 7. Update status ke "completed"
	_, err = s.store.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
		ID:     planID,
		Status: "completed",
	})

	return &CalculateResult{
		Placements:    resp.Data.Placements,
		UnfittedItems: resp.Data.Unfitted,
		Statistics:    resp.Data.Stats,
	}, nil
}

func (s *PlanService) buildPackRequest(container store.Container, items []store.GetPlanItemsRow) gateway.PackRequest {
	packItems := make([]gateway.PackItem, 0, len(items))
	for _, item := range items {
		packItems = append(packItems, gateway.PackItem{
			ItemID:   item.ProductID.String(),
			Label:    item.Label,
			Length:   item.LengthMm,
			Width:    item.WidthMm,
			Height:   item.HeightMm,
			Weight:   item.WeightKg,
			Quantity: int(item.Quantity),
		})
	}

	return gateway.PackRequest{
		Units: "mm",
		Container: gateway.PackContainer{
			Length:    container.LengthMm,
			Width:     container.WidthMm,
			Height:    container.HeightMm,
			MaxWeight: container.MaxWeightKg,
		},
		Items: packItems,
		Options: gateway.PackOptions{
			CheckStable: true,
			BiggerFirst: true,
		},
	}
}

func (s *PlanService) savePlacements(ctx context.Context, planID uuid.UUID, placements []gateway.PackPlacement) error {
	// Hapus placements lama terlebih dahulu
	err := s.store.DeletePlanPlacements(ctx, planID)
	if err != nil {
		return fmt.Errorf("delete old placements: %w", err)
	}

	// Simpan placements baru
	for _, p := range placements {
		productID, err := uuid.Parse(p.ItemID)
		if err != nil {
			continue
		}

		err = s.store.SavePlacement(ctx, store.SavePlacementParams{
			PlanID:     planID,
			ProductID:  productID,
			PosX:       p.PosX,
			PosY:       p.PosY,
			PosZ:       p.PosZ,
			Rotation:   int32(p.Rotation),
			StepNumber: int32(p.StepNumber),
		})
		if err != nil {
			return fmt.Errorf("save placement: %w", err)
		}
	}

	return nil
}
