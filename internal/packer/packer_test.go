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
		assert.Equal(t, "CONT-001", res.ContainerID)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		assert.Equal(t, "ITEM-1", res.PackedItems[0].ItemID)
		assert.Empty(t, res.UnfitItems)
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
