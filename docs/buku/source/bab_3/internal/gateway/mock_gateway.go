package gateway

import (
	"context"
)

// MockPackingGateway adalah implementasi mock PackingGateway untuk testing dan demo.
// Gateway ini mensimulasikan response dari Packing Service tanpa perlu menjalankan
// server Python yang sebenarnya.
type MockPackingGateway struct{}

// NewMockPackingGateway membuat instance mock gateway
func NewMockPackingGateway() *MockPackingGateway {
	return &MockPackingGateway{}
}

// Pack mensimulasikan kalkulasi packing dengan mengembalikan response dummy
func (g *MockPackingGateway) Pack(ctx context.Context, req PackRequest) (*PackResponse, error) {
	// Simulasikan placements untuk setiap item
	placements := make([]PackPlacement, 0)
	unfitted := make([]PackUnfitted, 0)

	stepNumber := 1
	currentX := 0.0
	containerWidth := req.Container.Width

	for _, item := range req.Items {
		for i := 0; i < item.Quantity; i++ {
			// Simulasi sederhana: tumpuk item berurutan
			// Dalam implementasi nyata, algoritma akan menghitung posisi optimal
			if currentX+item.Length <= req.Container.Length {
				placements = append(placements, PackPlacement{
					ItemID:     item.ItemID,
					Label:      item.Label,
					PosX:       currentX,
					PosY:       0,
					PosZ:       0,
					Rotation:   0,
					StepNumber: stepNumber,
				})
				currentX += item.Length + 10 // spacing 10mm
				stepNumber++
			} else {
				// Item tidak muat
				unfitted = append(unfitted, PackUnfitted{
					ItemID: item.ItemID,
					Label:  item.Label,
					Count:  1,
				})
			}
		}
	}

	// Hitung fitted count
	fittedCount := len(placements)
	unfittedCount := 0
	for _, u := range unfitted {
		unfittedCount += u.Count
	}

	// Konsolidasi unfitted items yang sama
	unfittedMap := make(map[string]*PackUnfitted)
	for _, u := range unfitted {
		if existing, ok := unfittedMap[u.ItemID]; ok {
			existing.Count++
		} else {
			unfittedMap[u.ItemID] = &PackUnfitted{
				ItemID: u.ItemID,
				Label:  u.Label,
				Count:  1,
			}
		}
	}
	consolidatedUnfitted := make([]PackUnfitted, 0, len(unfittedMap))
	for _, u := range unfittedMap {
		consolidatedUnfitted = append(consolidatedUnfitted, *u)
	}
	// Untuk development future, containerWidth bisa digunakan untuk
	// simulasi yang lebih akurat
	_ = containerWidth

	return &PackResponse{
		Success: true,
		Data: &PackData{
			Units:      req.Units,
			Placements: placements,
			Unfitted:   consolidatedUnfitted,
			Stats: PackStats{
				FittedCount:   fittedCount,
				UnfittedCount: unfittedCount,
				PackTimeMs:    15, // Simulasi waktu kalkulasi
			},
		},
	}, nil
}
