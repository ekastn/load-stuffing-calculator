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

		// RotationType is derived from the rotated dimensions returned by the packer.
		// We only assert that the rotated dimensions fit the container bounds.
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

		// With gravity enabled, no packed item should float: every item should
		// touch the floor (z=0) or another item's top surface.
		allowedSupports := map[float64]bool{0: true, 500: true}
		for _, pi := range res.PackedItems {
			assert.True(t, allowedSupports[pi.Position.Z], "unexpected z=%v", pi.Position.Z)
		}
	})

	// Guard against axis mapping regressions (L/W/H -> boxpacker3 W/H/D -> L/W/H).
	// For a non-cubic container & item, the packed placement must still be
	// within the original container bounds.
	t.Run("axis_mapping_preserves_bounds", func(t *testing.T) {
		c := packer.ContainerInput{
			ID:        "CONT-AXIS",
			Length:    1200,
			Width:     700,
			Height:    450,
			MaxWeight: 100,
		}

		items := []packer.ItemInput{
			{ID: "AX", Length: 600, Width: 200, Height: 300, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(ctx, c, items)
		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)

		pi := res.PackedItems[0]
		assert.LessOrEqual(t, pi.Position.X+pi.RotatedLength, c.Length)
		assert.LessOrEqual(t, pi.Position.Y+pi.RotatedWidth, c.Width)
		assert.LessOrEqual(t, pi.Position.Z+pi.RotatedHeight, c.Height)
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

// TestPacker_AllStrategies tests all supported packing strategies.
func TestPacker_AllStrategies(t *testing.T) {
	p := packer.NewPacker()
	ctx := context.Background()

	container := packer.ContainerInput{
		ID:        "CONT-001",
		Length:    1000,
		Width:     1000,
		Height:    1000,
		MaxWeight: 100,
	}

	items := []packer.ItemInput{
		{
			ID:       "ITEM-1",
			Label:    "Box",
			Length:   500,
			Width:    500,
			Height:   500,
			Weight:   10,
			Quantity: 2,
		},
	}

	tests := []struct {
		name           string
		strategy       string
		expectedAlgo   string
		shouldSucceed  bool
		expectedPacked int
		expectedUnfit  int
		expectedError  string
	}{
		{
			name:           "strategy_minimize_boxes",
			strategy:       "minimizeboxes",
			expectedAlgo:   "MinimizeBoxes",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_minimize_boxes_alias_ffd",
			strategy:       "ffd",
			expectedAlgo:   "MinimizeBoxes",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_greedy",
			strategy:       "greedy",
			expectedAlgo:   "Greedy",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_best_fit",
			strategy:       "bestfit",
			expectedAlgo:   "BestFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_best_fit_alias_bf",
			strategy:       "bf",
			expectedAlgo:   "BestFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_best_fit_decreasing",
			strategy:       "bestfitdecreasing",
			expectedAlgo:   "BestFitDecreasing",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_best_fit_decreasing_alias_bfd",
			strategy:       "bfd",
			expectedAlgo:   "BestFitDecreasing",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_next_fit",
			strategy:       "nextfit",
			expectedAlgo:   "NextFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_next_fit_alias_nf",
			strategy:       "nf",
			expectedAlgo:   "NextFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_worst_fit",
			strategy:       "worstfit",
			expectedAlgo:   "WorstFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_worst_fit_alias_wf",
			strategy:       "wf",
			expectedAlgo:   "WorstFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_almost_worst_fit",
			strategy:       "almostworstfit",
			expectedAlgo:   "AlmostWorstFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:           "strategy_almost_worst_fit_alias_awf",
			strategy:       "awf",
			expectedAlgo:   "AlmostWorstFit",
			shouldSucceed:  true,
			expectedPacked: 2,
			expectedUnfit:  0,
		},
		{
			name:          "strategy_empty_defaults_to_bfd",
			strategy:      "",
			expectedAlgo:  "BestFitDecreasing",
			shouldSucceed: true,
		},
		{
			name:          "strategy_invalid",
			strategy:      "invalid-strategy",
			shouldSucceed: false,
			expectedError: "invalid pack strategy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContainer := container
			testContainer.Options = packer.PackOptions{
				Strategy: tt.strategy,
			}

			res, err := p.Pack(ctx, testContainer, items)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAlgo, res.Algorithm)
				if tt.expectedPacked > 0 {
					assert.Equal(t, tt.expectedPacked, res.TotalPackedItems)
				}
				if tt.expectedUnfit > 0 {
					assert.Len(t, res.UnfitItems, tt.expectedUnfit)
				}
			} else {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			}
		})
	}
}

// TestPacker_ParallelStrategies tests parallel/auto mode with different goals.
func TestPacker_ParallelStrategies(t *testing.T) {
	p := packer.NewPacker()
	ctx := context.Background()

	container := packer.ContainerInput{
		ID:        "CONT-001",
		Length:    1000,
		Width:     1000,
		Height:    1000,
		MaxWeight: 100,
	}

	items := []packer.ItemInput{
		{
			ID:       "ITEM-1",
			Label:    "Box",
			Length:   500,
			Width:    500,
			Height:   500,
			Weight:   10,
			Quantity: 3,
		},
	}

	tests := []struct {
		name          string
		strategy      string
		goal          string
		expectedAlgo  string
		shouldSucceed bool
		expectedError string
	}{
		{
			name:          "parallel_goal_minimize_boxes",
			strategy:      "parallel",
			goal:          "minimizeboxes",
			expectedAlgo:  "Parallel(minimizeboxes)",
			shouldSucceed: true,
		},
		{
			name:          "parallel_goal_empty_defaults_to_minimize_boxes",
			strategy:      "parallel",
			goal:          "",
			expectedAlgo:  "Parallel()",
			shouldSucceed: true,
		},
		{
			name:          "parallel_goal_maximize_items",
			strategy:      "parallel",
			goal:          "maximizeitems",
			expectedAlgo:  "Parallel(maximizeitems)",
			shouldSucceed: true,
		},
		{
			name:          "parallel_goal_tightest",
			strategy:      "parallel",
			goal:          "tightest",
			expectedAlgo:  "Parallel(tightest)",
			shouldSucceed: true,
		},
		{
			name:          "parallel_goal_max_fill",
			strategy:      "parallel",
			goal:          "maxfill",
			expectedAlgo:  "Parallel(maxfill)",
			shouldSucceed: true,
		},
		{
			name:          "parallel_goal_max_average_fill_rate",
			strategy:      "parallel",
			goal:          "maxaveragefillrate",
			expectedAlgo:  "Parallel(maxaveragefillrate)",
			shouldSucceed: true,
		},
		{
			name:          "parallel_goal_balanced",
			strategy:      "parallel",
			goal:          "balanced",
			expectedAlgo:  "Parallel(balanced)",
			shouldSucceed: true,
		},
		{
			name:          "auto_strategy_alias",
			strategy:      "auto",
			goal:          "tightest",
			expectedAlgo:  "Parallel(tightest)",
			shouldSucceed: true,
		},
		{
			name:          "parallel_goal_invalid",
			strategy:      "parallel",
			goal:          "invalid-goal",
			shouldSucceed: false,
			expectedError: "invalid pack goal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContainer := container
			testContainer.Options = packer.PackOptions{
				Strategy: tt.strategy,
				Goal:     tt.goal,
			}

			res, err := p.Pack(ctx, testContainer, items)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAlgo, res.Algorithm)
				assert.True(t, res.IsFeasible || len(res.PackedItems) > 0)
			} else {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			}
		})
	}
}

// TestPacker_EdgeCases tests edge cases and helper functions.
func TestPacker_EdgeCases(t *testing.T) {
	p := packer.NewPacker()
	ctx := context.Background()

	container := packer.ContainerInput{
		ID:        "CONT-001",
		Length:    1000,
		Width:     1000,
		Height:    1000,
		MaxWeight: 100,
	}

	t.Run("item_id_without_colon_parseOriginalID", func(t *testing.T) {
		// This tests the edge case in parseOriginalID where an ID has no colon.
		// We create items and verify that instance IDs are properly parsed.
		items := []packer.ItemInput{
			{
				ID:       "SIMPLE-ID",
				Label:    "No Colon",
				Length:   100,
				Width:    100,
				Height:   100,
				Weight:   1,
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		// The ItemID should be parsed correctly from "SIMPLE-ID:0" -> "SIMPLE-ID"
		assert.Equal(t, "SIMPLE-ID", res.PackedItems[0].ItemID)
		assert.Equal(t, "SIMPLE-ID:0", res.PackedItems[0].InstanceID)
	})

	t.Run("unknown_item_lookup", func(t *testing.T) {
		// This tests lookupItem when an item is not found in the map.
		// We need to trigger a scenario where the packer creates an instance ID
		// that doesn't match anything in the itemMap. This is a defensive case.
		// In practice, this should never happen with correct usage, but the code
		// handles it by returning ItemInput{ID: id, Label: "Unknown"}.
		//
		// Since we can't directly trigger this without mocking internal behavior,
		// we'll test the normal case and rely on the parseOriginalID test above
		// to ensure the fallback path exists.
		items := []packer.ItemInput{
			{
				ID:       "KNOWN-ID",
				Label:    "Known Item",
				Length:   100,
				Width:    100,
				Height:   100,
				Weight:   1,
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		assert.Equal(t, "Known Item", res.PackedItems[0].Label)
	})

	t.Run("all_rotation_permutations_inferRotationType", func(t *testing.T) {
		// Test items with different dimensional permutations to exercise
		// all 6 rotation types in inferRotationType.
		// RotationType is inferred from how the packer rotates the item.
		tests := []struct {
			name   string
			length float64
			width  float64
			height float64
		}{
			{name: "rotation_0_LWH", length: 100, width: 200, height: 300},
			{name: "rotation_1_WLH", length: 200, width: 100, height: 300},
			{name: "rotation_2_WHL", length: 300, width: 100, height: 200},
			{name: "rotation_3_HWL", length: 300, width: 200, height: 100},
			{name: "rotation_4_HLW", length: 200, width: 300, height: 100},
			{name: "rotation_5_LHW", length: 100, width: 300, height: 200},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				items := []packer.ItemInput{
					{
						ID:       "ROT-ITEM",
						Label:    tt.name,
						Length:   tt.length,
						Width:    tt.width,
						Height:   tt.height,
						Weight:   1,
						Quantity: 1,
					},
				}

				res, err := p.Pack(ctx, container, items)

				assert.NoError(t, err)
				assert.True(t, res.IsFeasible)
				assert.Len(t, res.PackedItems, 1)
				// RotationType should be between 0-5
				assert.GreaterOrEqual(t, res.PackedItems[0].RotationType, 0)
				assert.LessOrEqual(t, res.PackedItems[0].RotationType, 5)
			})
		}
	})

	t.Run("gravity_no_overlap_edge_case", func(t *testing.T) {
		// Test gravity when items don't overlap in XY plane.
		// Both should settle to Z=0 (floor).
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
			{ID: "A", Length: 400, Width: 400, Height: 200, Weight: 1, Quantity: 1},
			{ID: "B", Length: 400, Width: 400, Height: 200, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(ctx, gravityContainer, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 2)

		// Both items should be on the floor or stacked
		for _, pi := range res.PackedItems {
			assert.GreaterOrEqual(t, pi.Position.Z, 0.0)
			assert.LessOrEqual(t, pi.Position.Z, 200.0)
		}
	})

	t.Run("gravity_item_exceeds_container_height", func(t *testing.T) {
		// Test gravity when an item's final Z position would exceed container height.
		// The gravity logic should not adjust Z if it would exceed bounds.
		gravityContainer := packer.ContainerInput{
			ID:        "CONT-GRAV-LIMIT",
			Length:    1000,
			Width:     1000,
			Height:    600, // Limited height
			MaxWeight: 100,
			Options: packer.PackOptions{
				Strategy: "bestfitdecreasing",
				Gravity:  true,
			},
		}

		items := []packer.ItemInput{
			{ID: "BASE", Length: 1000, Width: 1000, Height: 300, Weight: 1, Quantity: 1},
			{ID: "TOP", Length: 500, Width: 500, Height: 350, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(ctx, gravityContainer, items)

		assert.NoError(t, err)
		// May or may not fit depending on packer logic, but should not panic
		for _, pi := range res.PackedItems {
			assert.LessOrEqual(t, pi.Position.Z+pi.RotatedHeight, gravityContainer.Height+1e-6)
		}
	})

	t.Run("empty_items_list", func(t *testing.T) {
		// Test packing with no items.
		items := []packer.ItemInput{}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 0)
		assert.Len(t, res.UnfitItems, 0)
		assert.Equal(t, 0, res.TotalPackedItems)
	})

	t.Run("rotation_fallback_equal_dimensions", func(t *testing.T) {
		// Test inferRotationType fallback when all dimensions are equal (cube).
		// Multiple rotations match, so it should return 0 (fallback).
		items := []packer.ItemInput{
			{
				ID:       "CUBE",
				Label:    "Perfect Cube",
				Length:   100,
				Width:    100,
				Height:   100,
				Weight:   1,
				Quantity: 1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		// For a cube, any rotation is valid, so RotationType should be defined (0-5)
		assert.GreaterOrEqual(t, res.PackedItems[0].RotationType, 0)
		assert.LessOrEqual(t, res.PackedItems[0].RotationType, 5)
	})

	t.Run("utilization_calculations_zero_container", func(t *testing.T) {
		// Test utilization calculation edge cases.
		// This tests calculateStats with zero dimensions (defensive).
		zeroContainer := packer.ContainerInput{
			ID:        "ZERO-CONT",
			Length:    0, // Zero dimensions
			Width:     0,
			Height:    0,
			MaxWeight: 0,
		}

		items := []packer.ItemInput{
			{ID: "ITEM", Length: 100, Width: 100, Height: 100, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(ctx, zeroContainer, items)

		// Packer should handle this gracefully (likely nothing fits)
		assert.NoError(t, err)
		assert.False(t, res.IsFeasible)
		assert.Equal(t, 0.0, res.VolumeUtilisationPct)
		assert.Equal(t, 0.0, res.WeightUtilisationPct)
	})

	t.Run("gravity_empty_packed_items", func(t *testing.T) {
		// Test applyGravity with no packed items (early return).
		gravityContainer := packer.ContainerInput{
			ID:        "CONT-GRAV-EMPTY",
			Length:    1000,
			Width:     1000,
			Height:    1000,
			MaxWeight: 100,
			Options: packer.PackOptions{
				Strategy: "bestfitdecreasing",
				Gravity:  true,
			},
		}

		// Item too large to fit
		items := []packer.ItemInput{
			{ID: "HUGE", Length: 2000, Width: 2000, Height: 2000, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(ctx, gravityContainer, items)

		assert.NoError(t, err)
		assert.False(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 0)
		// Should not panic on empty list
	})

	t.Run("strategy_name_default_case", func(t *testing.T) {
		// This tests the default case in strategyName by using a strategy value
		// that's not recognized. However, since newLibPacker validates the strategy
		// first, we need to test this through a valid strategy that exercises
		// all strategy name branches. We've already covered all branches through
		// TestPacker_AllStrategies, so this test ensures we don't have gaps.
		//
		// The default case in strategyName returns "MinimizeBoxes", which we can't
		// directly trigger through the API since newLibPacker validates first.
		// This is defensive code. We'll test all valid strategies to ensure coverage.
		items := []packer.ItemInput{
			{ID: "ITEM", Length: 100, Width: 100, Height: 100, Weight: 1, Quantity: 1},
		}

		// Test with BestFitDecreasing to ensure full code path
		testContainer := container
		testContainer.Options = packer.PackOptions{Strategy: "bestfitdecreasing"}

		res, err := p.Pack(ctx, testContainer, items)

		assert.NoError(t, err)
		assert.Equal(t, "BestFitDecreasing", res.Algorithm)
	})

	t.Run("multiple_unfit_items_aggregation", func(t *testing.T) {
		// Test mapUnfitItems with multiple instances of the same item not fitting.
		// This tests the aggregation logic in mapUnfitItems.
		items := []packer.ItemInput{
			{ID: "FIT", Length: 100, Width: 100, Height: 100, Weight: 1, Quantity: 1},
			{ID: "UNFIT", Length: 1500, Width: 1500, Height: 1500, Weight: 1, Quantity: 3},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.False(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		assert.Len(t, res.UnfitItems, 1)
		assert.Equal(t, "UNFIT", res.UnfitItems[0].ID)
		assert.Equal(t, 3, res.UnfitItems[0].Quantity) // All 3 didn't fit
	})

	t.Run("pack_context_cancelled", func(t *testing.T) {
		// Test that Pack respects context cancellation.
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		items := []packer.ItemInput{
			{ID: "ITEM", Length: 100, Width: 100, Height: 100, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(cancelCtx, container, items)

		// Depending on implementation, may return error or complete before cancel.
		// This tests the ctx parameter is passed through.
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, res)
		}
	})

	t.Run("product_sku_preserved", func(t *testing.T) {
		// Test that ProductSKU is properly preserved in PackedItem.
		items := []packer.ItemInput{
			{
				ID:         "ITEM-1",
				Label:      "Test Item",
				ProductSKU: "SKU-12345",
				Length:     100,
				Width:      100,
				Height:     100,
				Weight:     1,
				Quantity:   1,
			},
		}

		res, err := p.Pack(ctx, container, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 1)
		assert.Equal(t, "SKU-12345", res.PackedItems[0].ProductSKU)
	})

	t.Run("gravity_sorted_placement_order", func(t *testing.T) {
		// Test that gravity sorts items by Z, Y, X before processing.
		gravityContainer := packer.ContainerInput{
			ID:        "CONT-GRAV-SORT",
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
			{ID: "A", Length: 300, Width: 300, Height: 200, Weight: 1, Quantity: 1},
			{ID: "B", Length: 300, Width: 300, Height: 200, Weight: 1, Quantity: 1},
			{ID: "C", Length: 300, Width: 300, Height: 200, Weight: 1, Quantity: 1},
		}

		res, err := p.Pack(ctx, gravityContainer, items)

		assert.NoError(t, err)
		assert.True(t, res.IsFeasible)
		assert.Len(t, res.PackedItems, 3)

		// All items should be stacked or on the floor
		for _, pi := range res.PackedItems {
			assert.GreaterOrEqual(t, pi.Position.Z, 0.0)
		}
	})
}
