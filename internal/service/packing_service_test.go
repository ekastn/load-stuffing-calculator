package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/gateway"
	"github.com/ekastn/load-stuffing-calculator/internal/packer"
	"github.com/stretchr/testify/require"
)

func TestPackingService_Pack_SortsByStepAndComputesRotationDims(t *testing.T) {
	// Fake packing microservice response: steps deliberately out of order.
	resp := gateway.PackResponse{
		Success: true,
		Data: &gateway.PackDataOut{
			Units: "mm",
			Placements: []gateway.PackPlacementOut{
				{ItemID: "B", Label: "Box B", PosX: 0, PosY: 0, PosZ: 0, Rotation: 1, StepNumber: 2},
				{ItemID: "A", Label: "Box A", PosX: 10, PosY: 0, PosZ: 0, Rotation: 0, StepNumber: 1},
			},
			Unfitted: []gateway.PackUnfittedOut{},
			Stats:    gateway.PackStatsOut{TotalTimeMs: 12},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/pack", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	gw := gateway.NewHTTPPackingGateway(srv.URL, 0)
	p := NewPackingService(gw)

	res, err := p.Pack(context.Background(), packer.ContainerInput{ID: "c", Length: 100, Width: 100, Height: 100, MaxWeight: 100}, []packer.ItemInput{
		{ID: "A", Label: "Box A", Length: 10, Width: 20, Height: 30, Weight: 1, Quantity: 1},
		{ID: "B", Label: "Box B", Length: 11, Width: 22, Height: 33, Weight: 1, Quantity: 1},
	})
	require.NoError(t, err)

	require.Equal(t, int64(12), res.DurationMs)
	require.Len(t, res.PackedItems, 2)

	// Ensure sorted by step_number.
	require.Equal(t, "A", res.PackedItems[0].ItemID)
	require.Equal(t, "B", res.PackedItems[1].ItemID)

	// Rotation=1 swaps L and W.
	require.Equal(t, float64(22), res.PackedItems[1].RotatedLength)
	require.Equal(t, float64(11), res.PackedItems[1].RotatedWidth)
	require.Equal(t, float64(33), res.PackedItems[1].RotatedHeight)
}
