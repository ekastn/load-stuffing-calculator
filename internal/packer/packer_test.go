package packer_test

import (
	"context"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/packer"
	"github.com/stretchr/testify/assert"
)

func TestPacker_Pack(t *testing.T) {
	p := packer.NewPacker()
	ctx := context.Background()

	// Standard Container (e.g. 100x100x100, 10kg)
	container := packer.ContainerInput{
		ID:        "CONT-001",
		Length:    1000,
		Width:     1000,
		Height:    1000,
		MaxWeight: 100, // 100 kg
	}

	t.Run("single_item_fits", func(t *testing.T) {
		items := []packer.ItemInput{
			{
				ID:       "ITEM-1",
				Label:    "Small Box",
				Length:   100,
				Width:    100,
				Height:   100,
				Weight:   1, // 1 kg
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.Equal(t, "BestFitDecreasing", res.Algorithm)
		assert.Equal(t, "CONT-001", res.ContainerID)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		assert.Equal(t, "ITEM-1", res.PackedItems[0].ItemID)
		assert.Empty(t, res.UnfitItems)
	})

	t.Run("rotation_required_to_fit", func(t *testing.T) {
		rotationContainer := packer.ContainerInput{
			ID:        "CONT-ROT",
			Length:    1000,
			Width:     600,
			Height:    400,
			MaxWeight: 100,
			Options: packer.PackOptions{
				Strategy: "bestfitdecreasing",
			},
		}

		items := []packer.ItemInput{
			{
				ID:       "ITEM-ROT",
				Label:    "Needs Rotation",
				Length:   500,
				Width:    800, // Too wide unless rotated
				Height:   300,
				Weight:   1,
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, rotationContainer, items)

		assert.NoError(t, err)
		assert.Equal(t, "BestFitDecreasing", res.Algorithm)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)

		pi := res.PackedItems[0]
		assert.Equal(t, "ITEM-ROT", pi.ItemID)

		// Only rotation 1 (swap length/width) can fit into 1000x600x400.
		assert.Equal(t, 1, pi.RotationType)
		assert.Equal(t, 800.0, pi.RotatedLength)
		assert.Equal(t, 500.0, pi.RotatedWidth)
		assert.Equal(t, 300.0, pi.RotatedHeight)

		assert.LessOrEqual(t, pi.Position.X+pi.RotatedLength, rotationContainer.Length)
		assert.LessOrEqual(t, pi.Position.Y+pi.RotatedWidth, rotationContainer.Width)
		assert.LessOrEqual(t, pi.Position.Z+pi.RotatedHeight, rotationContainer.Height)
	})

	t.Run("multiple_items_fit", func(t *testing.T) {
		items := []packer.ItemInput{
			{
				ID:       "ITEM-1",
				Length:   500,
				Width:    500,
				Height:   500,
				Weight:   10,
				Quantity: 8, // 8 * (500*500*500) = 8 * 0.125 m3 = 1 m3 (Exact fit volume-wise)
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Equal(t, 8, res.TotalPackedItems)
		assert.Len(t, res.PackedItems, 8)
	})

	t.Run("overflow_by_volume", func(t *testing.T) {
		items := []packer.ItemInput{
			{
				ID:       "ITEM-HUGE",
				Length:   2000, // Bigger than container
				Width:    2000,
				Height:   2000,
				Weight:   10,
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.False(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 0)
		assert.Len(t, res.UnfitItems, 1)
		assert.Equal(t, "ITEM-HUGE", res.UnfitItems[0].ID)
		assert.Equal(t, 1, res.UnfitItems[0].Quantity)
	})

	t.Run("overflow_by_weight", func(t *testing.T) {
		items := []packer.ItemInput{
			{
				ID:       "ITEM-HEAVY",
				Length:   100,
				Width:    100,
				Height:   100,
				Weight:   150, // 150 kg > 100 kg max
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.False(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 0)
		assert.Len(t, res.UnfitItems, 1)
	})

	t.Run("invalid_strategy_returns_error", func(t *testing.T) {
		badContainer := container
		badContainer.Options = packer.PackOptions{Strategy: "definitely-not-a-strategy"}

		items := []packer.ItemInput{
			{ID: "ITEM-1", Length: 100, Width: 100, Height: 100, Weight: 1, Quantity: 1},
		}

		_, err := p.Pack(ctx, badContainer, items)
		assert.Error(t, err)
	})

	// Gravity is implemented as a post-pass that settles items down to the
	// highest supporting surface beneath them (or the floor).
	t.Run("gravity_settles_items", func(t *testing.T) {
		gravityContainer := packer.ContainerInput{
			ID:        "CONT-GRAV",
			Length:    1000,
			Width:     1000,
			Height:    1000,
			MaxWeight: 100,
			Options: packer.PackOptions{
				Strategy: "bestfitdecreasing",
				Gravity:  true,
			},
		}

		items := []packer.ItemInput{
			{ID: "A", Length: 500, Width: 500, Height: 500, Weight: 1, Quantity: 1},
			{ID: "B", Length: 500, Width: 500, Height: 500, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(ctx, gravityContainer, items)
		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 2)

		// With gravity enabled, no packed item should "float": every item should
		// touch the floor (z=0) or another item's top surface.
		allowedSupports := map[float64]bool{0: true, 500: true}
		for _, pi := range res.PackedItems {
			assert.True(t, allowedSupports[pi.Position.Z], "unexpected z=%v", pi.Position.Z)
		}
	})

	// 1 fits, 1 doesn't
	t.Run("mixed_items", func(t *testing.T) {
		// 1 fits, 1 doesn't
		items := []packer.ItemInput{
			{
				ID:       "FIT",
				Length:   100,
				Width:    100,
				Height:   100,
				Weight:   10,
				Quantity: 1,
			},
			{
				ID:       "NO-FIT",
				Length:   2000,
				Width:    2000,
				Height:   2000,
				Weight:   10,
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.False(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		assert.Equal(t, "FIT", res.PackedItems[0].ItemID)
		assert.Len(t, res.UnfitItems, 1)
		assert.Equal(t, "NO-FIT", res.UnfitItems[0].ID)
	})

	t.Run("quantity_splitting", func(t *testing.T) {
		items := []packer.ItemInput{
			{
				ID:       "ITEM-SPLIT",
				Length:   500,
				Width:    1000, // Full width
				Height:   1000, // Full height
				Weight:   10,
				Quantity: 3, // Length 500 fits 2 times in 1000. 3rd should fail.
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.False(t, res.IsFeasible)
		assert.Equal(t, 2, res.TotalPackedItems)
		assert.Len(t, res.UnfitItems, 1)
		assert.Equal(t, "ITEM-SPLIT", res.UnfitItems[0].ID)
		assert.Equal(t, 1, res.UnfitItems[0].Quantity) // 1 failed
	})
}
