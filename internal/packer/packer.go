package packer

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
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

	bp, algoName, err := p.newLibPacker(container.Options)
	if err != nil {
		return PackingResult{}, err
	}

	// Run Packing
	packResult, err := bp.PackCtx(ctx, boxes, libItems)
	if err != nil {
		return PackingResult{}, fmt.Errorf("packing calculation failed: %w", err)
	}

	result := p.buildResult(container, packResult, itemMap, time.Since(start))
	result.Algorithm = algoName

	if container.Options.Gravity {
		p.applyGravity(container, &result)
	}

	return result, nil
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
		Algorithm:   "",
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
		dim := item.GetDimension() // rotated dimensions in same axis-order as input (Length, Width, Height)

		rot := p.inferRotationType(
			[3]float64{originalInput.Length, originalInput.Width, originalInput.Height},
			[3]float64{dim[0], dim[1], dim[2]},
		)

		packedItem := PackedItem{
			ItemID:        originalItemID,
			InstanceID:    item.GetID(),
			Label:         originalInput.Label,
			ProductSKU:    originalInput.ProductSKU,
			RotatedLength: dim[0],
			RotatedWidth:  dim[1],
			RotatedHeight: dim[2],
			Position: Position{
				X: pos[0],
				Y: pos[1],
				Z: pos[2],
			},
			RotationType: rot,
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

func (p *packer) inferRotationType(originalLWH, rotatedLWH [3]float64) int {
	// The boxpacker3 rotation code is a permutation of dimensions in the same
	// order used when constructing the item: NewItem(..., w, h, d).
	// In our adapter we pass (Length, Width, Height), so we can treat the
	// rotated dimensions returned by GetDimension() as (Length, Width, Height)
	// in the corresponding container axes.
	perms := [6][3]int{
		{0, 1, 2},
		{1, 0, 2},
		{1, 2, 0},
		{2, 1, 0},
		{2, 0, 1},
		{0, 2, 1},
	}

	const eps = 1e-6
	isClose := func(a, b float64) bool {
		return math.Abs(a-b) <= eps
	}

	for i, perm := range perms {
		candidate := [3]float64{originalLWH[perm[0]], originalLWH[perm[1]], originalLWH[perm[2]]}
		if isClose(candidate[0], rotatedLWH[0]) && isClose(candidate[1], rotatedLWH[1]) && isClose(candidate[2], rotatedLWH[2]) {
			return i
		}
	}

	// Fallback (e.g., equal dimensions where multiple rotations match)
	return 0
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

type rect struct {
	x1, y1 float64
	x2, y2 float64
}

func (r rect) overlapsXY(o rect, eps float64) bool {
	return (r.x1+eps) < o.x2 && (o.x1+eps) < r.x2 && (r.y1+eps) < o.y2 && (o.y1+eps) < r.y2
}

func strategyName(strategy boxpacker3.PackingStrategy) string {
	switch strategy {
	case boxpacker3.StrategyMinimizeBoxes:
		return "MinimizeBoxes"
	case boxpacker3.StrategyGreedy:
		return "Greedy"
	case boxpacker3.StrategyBestFit:
		return "BestFit"
	case boxpacker3.StrategyBestFitDecreasing:
		return "BestFitDecreasing"
	case boxpacker3.StrategyNextFit:
		return "NextFit"
	case boxpacker3.StrategyWorstFit:
		return "WorstFit"
	case boxpacker3.StrategyAlmostWorstFit:
		return "AlmostWorstFit"
	default:
		return "MinimizeBoxes"
	}
}

func (p *packer) newLibPacker(opts PackOptions) (*boxpacker3.Packer, string, error) {
	strategy := strings.ToLower(strings.TrimSpace(opts.Strategy))
	goal := strings.ToLower(strings.TrimSpace(opts.Goal))

	switch strategy {
	case "", "bestfitdecreasing", "bfd":
		const strat = boxpacker3.StrategyBestFitDecreasing
		return boxpacker3.NewPacker(boxpacker3.WithStrategy(strat)), strategyName(strat), nil
	case "minimizeboxes", "ffd":
		const strat = boxpacker3.StrategyMinimizeBoxes
		return boxpacker3.NewPacker(boxpacker3.WithStrategy(strat)), strategyName(strat), nil
	case "greedy":
		const strat = boxpacker3.StrategyGreedy
		return boxpacker3.NewPacker(boxpacker3.WithStrategy(strat)), strategyName(strat), nil
	case "bestfit", "bf":
		const strat = boxpacker3.StrategyBestFit
		return boxpacker3.NewPacker(boxpacker3.WithStrategy(strat)), strategyName(strat), nil
	case "nextfit", "nf":
		const strat = boxpacker3.StrategyNextFit
		return boxpacker3.NewPacker(boxpacker3.WithStrategy(strat)), strategyName(strat), nil
	case "worstfit", "wf":
		const strat = boxpacker3.StrategyWorstFit
		return boxpacker3.NewPacker(boxpacker3.WithStrategy(strat)), strategyName(strat), nil
	case "almostworstfit", "awf":
		const strat = boxpacker3.StrategyAlmostWorstFit
		return boxpacker3.NewPacker(boxpacker3.WithStrategy(strat)), strategyName(strat), nil
	case "parallel", "auto":
		var comparator boxpacker3.ComparatorFunc
		switch goal {
		case "", "minimizeboxes":
			comparator = boxpacker3.MinimizeBoxesGoal
		case "maximizeitems":
			comparator = boxpacker3.MaximizeItemsGoal
		case "tightest":
			comparator = boxpacker3.TightestPackingGoal
		case "maxfill", "maxaveragefillrate":
			comparator = boxpacker3.MaxAverageFillRateGoal
		case "balanced":
			comparator = boxpacker3.BalancedPackingGoal
		default:
			return nil, "", fmt.Errorf("invalid pack goal: %q", opts.Goal)
		}

		ps := boxpacker3.NewParallelStrategy(
			boxpacker3.WithAlgorithms(
				boxpacker3.NewMinimizeBoxesStrategy(),
				boxpacker3.NewBestFitDecreasingStrategy(),
				boxpacker3.NewBestFitStrategy(),
				boxpacker3.NewGreedyStrategy(),
			),
			boxpacker3.WithGoal(comparator),
		)

		return boxpacker3.NewPacker(boxpacker3.WithAlgorithm(ps)), "Parallel(" + goal + ")", nil
	default:
		return nil, "", fmt.Errorf("invalid pack strategy: %q", opts.Strategy)
	}
}

func (p *packer) applyGravity(container ContainerInput, result *PackingResult) {
	if len(result.PackedItems) == 0 {
		return
	}

	packed := make([]*PackedItem, 0, len(result.PackedItems))
	for i := range result.PackedItems {
		packed = append(packed, &result.PackedItems[i])
	}

	sort.SliceStable(packed, func(i, j int) bool {
		if packed[i].Position.Z != packed[j].Position.Z {
			return packed[i].Position.Z < packed[j].Position.Z
		}
		if packed[i].Position.Y != packed[j].Position.Y {
			return packed[i].Position.Y < packed[j].Position.Y
		}
		return packed[i].Position.X < packed[j].Position.X
	})

	const eps = 1e-6
	for i, cur := range packed {
		curRect := rect{cur.Position.X, cur.Position.Y, cur.Position.X + cur.RotatedLength, cur.Position.Y + cur.RotatedWidth}

		bestZ := 0.0
		for j := 0; j < i; j++ {
			below := packed[j]
			belowRect := rect{below.Position.X, below.Position.Y, below.Position.X + below.RotatedLength, below.Position.Y + below.RotatedWidth}
			if !curRect.overlapsXY(belowRect, eps) {
				continue
			}

			supportZ := below.Position.Z + below.RotatedHeight
			if supportZ > bestZ {
				bestZ = supportZ
			}
		}

		if bestZ+cur.RotatedHeight <= container.Height+eps {
			cur.Position.Z = bestZ
		}
	}
}
