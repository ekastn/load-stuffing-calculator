package gateway

import (
	"context"
)

// PackingGateway mendefinisikan kontrak untuk komunikasi dengan Packing Service.
// Interface ini memungkinkan kita mengganti implementasi (HTTP, mock, dll)
// tanpa mengubah service layer.
type PackingGateway interface {
	Pack(ctx context.Context, req PackRequest) (*PackResponse, error)
}

// PackRequest adalah request yang dikirim ke Packing Service
type PackRequest struct {
	Units     string        `json:"units"`     // Satuan ukuran: "mm" atau "cm"
	Container PackContainer `json:"container"` // Dimensi container
	Items     []PackItem    `json:"items"`     // Items yang akan di-pack
	Options   PackOptions   `json:"options"`   // Opsi algoritma
}

// PackContainer mendeskripsikan container untuk packing
type PackContainer struct {
	Length    float64 `json:"length"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	MaxWeight float64 `json:"max_weight"`
}

// PackItem mendeskripsikan satu item yang akan di-pack
type PackItem struct {
	ItemID   string  `json:"item_id"` // ID unik untuk tracking
	Label    string  `json:"label"`   // Label untuk display
	Length   float64 `json:"length"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Weight   float64 `json:"weight"`
	Quantity int     `json:"quantity"` // Jumlah item dengan dimensi sama
}

// PackOptions berisi opsi untuk algoritma packing
type PackOptions struct {
	FixPoint            bool    `json:"fix_point"`
	CheckStable         bool    `json:"check_stable"`
	SupportSurfaceRatio float64 `json:"support_surface_ratio"`
	BiggerFirst         bool    `json:"bigger_first"`
}

// PackResponse adalah response dari Packing Service
type PackResponse struct {
	Success bool       `json:"success"`
	Data    *PackData  `json:"data,omitempty"`
	Error   *PackError `json:"error,omitempty"`
}

// PackData berisi hasil kalkulasi jika sukses
type PackData struct {
	Units      string          `json:"units"`
	Placements []PackPlacement `json:"placements"` // Posisi setiap item
	Unfitted   []PackUnfitted  `json:"unfitted"`   // Items yang tidak muat
	Stats      PackStats       `json:"stats"`      // Statistik packing
}

// PackPlacement mendeskripsikan posisi satu item dalam container
type PackPlacement struct {
	ItemID     string  `json:"item_id"`
	Label      string  `json:"label"`
	PosX       float64 `json:"pos_x"`       // Posisi X dalam container
	PosY       float64 `json:"pos_y"`       // Posisi Y dalam container
	PosZ       float64 `json:"pos_z"`       // Posisi Z (tinggi) dalam container
	Rotation   int     `json:"rotation"`    // Rotasi yang diterapkan
	StepNumber int     `json:"step_number"` // Urutan penempatan
}

// PackUnfitted mendeskripsikan item yang tidak muat dalam container
type PackUnfitted struct {
	ItemID string `json:"item_id"`
	Label  string `json:"label"`
	Count  int    `json:"count"`
}

// PackStats berisi statistik hasil packing
type PackStats struct {
	FittedCount   int `json:"fitted_count"`   // Jumlah item yang muat
	UnfittedCount int `json:"unfitted_count"` // Jumlah item yang tidak muat
	PackTimeMs    int `json:"pack_time_ms"`   // Waktu kalkulasi
}

// PackError berisi detail error jika gagal
type PackError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
