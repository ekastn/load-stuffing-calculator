package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ekastn/load-stuffing-calculator/internal/gateway"
	"github.com/ekastn/load-stuffing-calculator/internal/packer"
)

// packingService adapts the packing microservice (py3dbp) into the existing
// packer.Packer interface so we can replace the algorithm without changing
// plan_service (and without breaking current tests).
//
// API stability notes:
// - We always send units="mm".
// - dto.CalculatePlanRequest options (strategy/goal/gravity) are ignored.
// - AllowRotation is ignored for now.

type packingService struct {
	gw gateway.PackingGateway
}

func NewPackingService(gw gateway.PackingGateway) packer.Packer {
	return &packingService{gw: gw}
}

func (s *packingService) Pack(ctx context.Context, container packer.ContainerInput, items []packer.ItemInput) (packer.PackingResult, error) {
	if s.gw == nil {
		return packer.PackingResult{}, fmt.Errorf("packing gateway is nil")
	}

	// Server-internal defaults for py3dbp.
	// Keep these stable unless we explicitly decide to expose them.
	req := gateway.PackRequest{
		Units: "mm",
		Container: gateway.PackContainerIn{
			Length:    container.Length,
			Width:     container.Width,
			Height:    container.Height,
			MaxWeight: container.MaxWeight,
		},
		Options: gateway.PackOptionsIn{
			FixPoint:            true,
			CheckStable:         true,
			SupportSurfaceRatio: 0.75,
			BiggerFirst:         true,
			PutType:             1,
		},
	}

	itemByID := make(map[string]packer.ItemInput, len(items))
	for _, it := range items {
		itemByID[it.ID] = it
		req.Items = append(req.Items, gateway.PackItemIn{
			ItemID:   it.ID,
			Label:    it.Label,
			Length:   it.Length,
			Width:    it.Width,
			Height:   it.Height,
			Weight:   it.Weight,
			Quantity: it.Quantity,
		})
	}

	started := time.Now()
	resp, err := s.gw.Pack(ctx, req)
	if err != nil {
		return packer.PackingResult{}, err
	}

	// Ensure placements are in the intended order (py3dbp putOrder())
	placements := append([]gateway.PackPlacementOut(nil), resp.Data.Placements...)
	sort.SliceStable(placements, func(i, j int) bool {
		return placements[i].StepNumber < placements[j].StepNumber
	})

	result := packer.PackingResult{
		ContainerID: container.ID,
		Algorithm:   "py3dbp",
		DurationMs:  time.Since(started).Milliseconds(),
	}
	if resp.Data.Stats.TotalTimeMs > 0 {
		result.DurationMs = int64(resp.Data.Stats.TotalTimeMs)
	}

	// Build packed items.
	instanceCounter := make(map[string]int)
	for _, pl := range placements {
		in, ok := itemByID[pl.ItemID]
		if !ok {
			return packer.PackingResult{}, fmt.Errorf("packing service returned unknown item_id: %s", pl.ItemID)
		}

		rotL, rotW, rotH := applyRotation(in.Length, in.Width, in.Height, pl.Rotation)

		instanceCounter[pl.ItemID]++
		instanceID := fmt.Sprintf("%s:%d", pl.ItemID, instanceCounter[pl.ItemID])

		packed := packer.PackedItem{
			ItemID:        pl.ItemID,
			InstanceID:    instanceID,
			Label:         in.Label,
			ProductSKU:    in.ProductSKU,
			RotatedLength: rotL,
			RotatedWidth:  rotW,
			RotatedHeight: rotH,
			Position: packer.Position{
				X: pl.PosX,
				Y: pl.PosY,
				Z: pl.PosZ,
			},
			RotationType: pl.Rotation,
		}

		result.PackedItems = append(result.PackedItems, packed)
		result.TotalVolumePackedM3 += (rotL * rotW * rotH)
		result.TotalWeightPackedKG += in.Weight
	}

	// totals are in mm^3 and kg per item instance; convert volume to m^3
	result.TotalVolumePackedM3 /= 1_000_000_000.0
	result.TotalPackedItems = len(result.PackedItems)

	// Unfitted items
	for _, u := range resp.Data.Unfitted {
		in, ok := itemByID[u.ItemID]
		if !ok {
			// keep best effort; this shouldn’t happen but we don’t want to crash
			continue
		}
		unfit := in
		unfit.Quantity = u.Count
		result.UnfitItems = append(result.UnfitItems, unfit)
	}

	result.IsFeasible = len(result.UnfitItems) == 0

	// Utilisation
	containerVolumeM3 := (container.Length * container.Width * container.Height) / 1_000_000_000.0
	if containerVolumeM3 > 0 {
		result.VolumeUtilisationPct = (result.TotalVolumePackedM3 / containerVolumeM3) * 100
	}
	if container.MaxWeight > 0 {
		result.WeightUtilisationPct = (result.TotalWeightPackedKG / container.MaxWeight) * 100
	}

	return result, nil
}

func applyRotation(l, w, h float64, rotation int) (float64, float64, float64) {
	candidates := [6][3]float64{
		{l, w, h},
		{w, l, h},
		{w, h, l},
		{h, w, l},
		{h, l, w},
		{l, h, w},
	}
	if rotation < 0 || rotation >= len(candidates) {
		return l, w, h
	}
	rot := candidates[rotation]
	return rot[0], rot[1], rot[2]
}
