package packer

import (
	"context"
	"fmt"
	"time"

	"github.com/bavix/boxpacker3"
)

// Packer defines the interface for 3D bin packing algorithms.
type Packer interface {
	Pack(ctx context.Context, container ContainerInput, items []ItemInput) (PackingResult, error)
}

// defaultPacker implements Packer using boxpacker3 library.
type packer struct {
}

// NewPacker creates a new instance of the default Packer.
func NewPacker() Packer {
	return &packer{}
}

func (p *packer) Pack(ctx context.Context, container ContainerInput, items []ItemInput) (PackingResult, error) {
	start := time.Now()

	boxes := []*boxpacker3.Box{p.toBox(container)}
	libItems, itemMap := p.toItems(items)

	// Configure Packer with BestFitDecreasing strategy
	bp := boxpacker3.NewPacker(
		boxpacker3.WithStrategy(boxpacker3.StrategyBestFitDecreasing),
	)

	// Run Packing
	packResult, err := bp.PackCtx(ctx, boxes, libItems)
	if err != nil {
		return PackingResult{}, fmt.Errorf("packing calculation failed: %w", err)
	}

	return p.buildResult(container, packResult, itemMap, time.Since(start)), nil
}

func (p *packer) toBox(c ContainerInput) *boxpacker3.Box {
	return boxpacker3.NewBox(
		c.ID,
		c.Length,
		c.Width,
		c.Height,
		c.MaxWeight*1000, // kg to grams
	)
}

func (p *packer) toItems(items []ItemInput) ([]*boxpacker3.Item, map[string]ItemInput) {
	var libItems []*boxpacker3.Item
	itemMap := make(map[string]ItemInput)

	for _, item := range items {
		itemMap[item.ID] = item
		for i := 0; i < item.Quantity; i++ {
			instanceID := fmt.Sprintf("%s:%d", item.ID, i)
			libItem := boxpacker3.NewItem(
				instanceID,
				item.Length,
				item.Width,
				item.Height,
				item.Weight*1000, // kg to grams
			)
			libItems = append(libItems, libItem)
		}
	}
	return libItems, itemMap
}

func (p *packer) buildResult(
	container ContainerInput,
	packResult *boxpacker3.Result,
	itemMap map[string]ItemInput,
	duration time.Duration,
) PackingResult {
	result := PackingResult{
		ContainerID: container.ID,
		Algorithm:   "BestFitDecreasing",
		DurationMs:  duration.Milliseconds(),
	}

	var usedBox *boxpacker3.Box
	if len(packResult.Boxes) > 0 {
		usedBox = packResult.Boxes[0]
	}

	if usedBox != nil {
		p.mapPackedItems(usedBox, itemMap, &result)
		p.calculateStats(container, &result)
	}

	p.mapUnfitItems(packResult.UnfitItems, itemMap, &result)

	result.IsFeasible = len(result.UnfitItems) == 0
	return result
}

func (p *packer) mapPackedItems(box *boxpacker3.Box, itemMap map[string]ItemInput, result *PackingResult) {
	for _, item := range box.GetItems() {
		originalItemID := p.parseOriginalID(item.GetID())
		originalInput := p.lookupItem(originalItemID, itemMap)

		pos := item.GetPosition()

		packedItem := PackedItem{
			ItemID:        originalItemID,
			InstanceID:    item.GetID(),
			Label:         originalInput.Label,
			ProductSKU:    originalInput.ProductSKU,
			RotatedLength: item.GetWidth(),
			RotatedWidth:  item.GetHeight(),
			RotatedHeight: item.GetDepth(),
			Position: Position{
				X: pos[0],
				Y: pos[1],
				Z: pos[2],
			},
			RotationType: 0,
		}

		result.PackedItems = append(result.PackedItems, packedItem)

		// Accumulate raw totals for stats
		result.TotalVolumePackedM3 += (packedItem.RotatedLength * packedItem.RotatedWidth * packedItem.RotatedHeight)
		result.TotalWeightPackedKG += item.GetWeight()
	}

	// Final unit conversions
	result.TotalVolumePackedM3 /= 1_000_000_000.0 // mm3 to m3
	result.TotalWeightPackedKG /= 1000.0          // grams to kg
	result.TotalPackedItems = len(result.PackedItems)
}

func (p *packer) mapUnfitItems(unfitItems []*boxpacker3.Item, itemMap map[string]ItemInput, result *PackingResult) {
	unfitCounts := make(map[string]int)
	for _, unfit := range unfitItems {
		id := p.parseOriginalID(unfit.GetID())
		unfitCounts[id]++
	}

	for id, count := range unfitCounts {
		if input, ok := itemMap[id]; ok {
			unfitInput := input
			unfitInput.Quantity = count
			result.UnfitItems = append(result.UnfitItems, unfitInput)
		}
	}
}

func (p *packer) calculateStats(container ContainerInput, result *PackingResult) {
	contVol := container.Length * container.Width * container.Height // mm3
	contVolM3 := contVol / 1_000_000_000.0

	if contVolM3 > 0 {
		result.VolumeUtilisationPct = (result.TotalVolumePackedM3 / contVolM3) * 100
	}
	if container.MaxWeight > 0 {
		result.WeightUtilisationPct = (result.TotalWeightPackedKG / container.MaxWeight) * 100
	}
}

func (p *packer) parseOriginalID(instanceID string) string {
	for i := len(instanceID) - 1; i >= 0; i-- {
		if instanceID[i] == ':' {
			return instanceID[:i]
		}
	}
	return instanceID
}

func (p *packer) lookupItem(id string, itemMap map[string]ItemInput) ItemInput {
	if input, ok := itemMap[id]; ok {
		return input
	}
	return ItemInput{ID: id, Label: "Unknown"}
}