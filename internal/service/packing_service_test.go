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

func TestPackingService_Pack_NilGateway(t *testing.T) {
	p := NewPackingService(nil)
	_, err := p.Pack(context.Background(), packer.ContainerInput{}, []packer.ItemInput{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "packing gateway is nil")
}

func TestPackingService_Pack_GatewayError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	gw := gateway.NewHTTPPackingGateway(srv.URL, 0)
	p := NewPackingService(gw)

	_, err := p.Pack(context.Background(), packer.ContainerInput{ID: "c", Length: 100, Width: 100, Height: 100}, []packer.ItemInput{
		{ID: "A", Length: 10, Width: 20, Height: 30, Weight: 1, Quantity: 1},
	})
	require.Error(t, err)
}

func TestPackingService_Pack_UnknownItemIDInPlacement(t *testing.T) {
	// Server returns a placement for an item_id we didn't send
	resp := gateway.PackResponse{
		Success: true,
		Data: &gateway.PackDataOut{
			Units: "mm",
			Placements: []gateway.PackPlacementOut{
				{ItemID: "UNKNOWN", Label: "Unknown", PosX: 0, PosY: 0, PosZ: 0, Rotation: 0, StepNumber: 1},
			},
			Unfitted: []gateway.PackUnfittedOut{},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	gw := gateway.NewHTTPPackingGateway(srv.URL, 0)
	p := NewPackingService(gw)

	_, err := p.Pack(context.Background(), packer.ContainerInput{ID: "c", Length: 100, Width: 100, Height: 100}, []packer.ItemInput{
		{ID: "A", Length: 10, Width: 20, Height: 30, Weight: 1, Quantity: 1},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown item_id")
}

func TestPackingService_Pack_UnknownItemIDInUnfitted(t *testing.T) {
	// Server returns an unfitted item_id we didn't send (best effort: should continue)
	resp := gateway.PackResponse{
		Success: true,
		Data: &gateway.PackDataOut{
			Units: "mm",
			Placements: []gateway.PackPlacementOut{
				{ItemID: "A", Label: "Box A", PosX: 0, PosY: 0, PosZ: 0, Rotation: 0, StepNumber: 1},
			},
			Unfitted: []gateway.PackUnfittedOut{
				{ItemID: "UNKNOWN", Count: 1},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	gw := gateway.NewHTTPPackingGateway(srv.URL, 0)
	p := NewPackingService(gw)

	res, err := p.Pack(context.Background(), packer.ContainerInput{ID: "c", Length: 100, Width: 100, Height: 100}, []packer.ItemInput{
		{ID: "A", Label: "Box A", Length: 10, Width: 20, Height: 30, Weight: 1, Quantity: 1},
	})
	require.NoError(t, err)
	// Should continue despite unknown unfitted item
	require.Len(t, res.PackedItems, 1)
	require.Len(t, res.UnfitItems, 0) // Unknown item skipped
}

func TestPackingService_applyRotation(t *testing.T) {
	tests := []struct {
		name     string
		l, w, h  float64
		rotation int
		wantL    float64
		wantW    float64
		wantH    float64
	}{
		{
			name: "rotation_0_no_change",
			l:    10, w: 20, h: 30,
			rotation: 0,
			wantL:    10, wantW: 20, wantH: 30,
		},
		{
			name: "rotation_1_swap_l_w",
			l:    10, w: 20, h: 30,
			rotation: 1,
			wantL:    20, wantW: 10, wantH: 30,
		},
		{
			name: "rotation_2",
			l:    10, w: 20, h: 30,
			rotation: 2,
			wantL:    20, wantW: 30, wantH: 10,
		},
		{
			name: "rotation_3",
			l:    10, w: 20, h: 30,
			rotation: 3,
			wantL:    30, wantW: 20, wantH: 10,
		},
		{
			name: "rotation_4",
			l:    10, w: 20, h: 30,
			rotation: 4,
			wantL:    30, wantW: 10, wantH: 20,
		},
		{
			name: "rotation_5",
			l:    10, w: 20, h: 30,
			rotation: 5,
			wantL:    10, wantW: 30, wantH: 20,
		},
		{
			name: "rotation_negative_returns_original",
			l:    10, w: 20, h: 30,
			rotation: -1,
			wantL:    10, wantW: 20, wantH: 30,
		},
		{
			name: "rotation_out_of_bounds_returns_original",
			l:    10, w: 20, h: 30,
			rotation: 6,
			wantL:    10, wantW: 20, wantH: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotL, gotW, gotH := applyRotation(tt.l, tt.w, tt.h, tt.rotation)
			require.Equal(t, tt.wantL, gotL, "length mismatch")
			require.Equal(t, tt.wantW, gotW, "width mismatch")
			require.Equal(t, tt.wantH, gotH, "height mismatch")
		})
	}
}
