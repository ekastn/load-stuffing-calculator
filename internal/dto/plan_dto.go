package dto

import (
	"time"
)

type CreatePlanRequest struct {
	Title         string  `json:"title" binding:"required,max=150" example:"Ekspor Jepang - Desember 2025"`
	Notes         *string `json:"notes,omitempty" binding:"omitempty,max=500" example:"Barang elektronik di depan"`
	AutoCalculate *bool   `json:"auto_calculate" binding:"omitempty" example:"true"`

	Container CreatePlanContainer `json:"container" binding:"required"`
	Items     []CreatePlanItem    `json:"items" binding:"required,min=1,max=1000,dive"`
}

type CreatePlanContainer struct {
	// Preset container if null using custom container
	ContainerID *string `json:"container_id,omitempty" binding:"omitempty,uuid" example:"a1b2c3d4-..."`

	LengthMM    *float64 `json:"length_mm,omitempty" binding:"omitempty,gt=0" example:"12000"`
	WidthMM     *float64 `json:"width_mm,omitempty" binding:"omitempty,gt=0" example:"2350"`
	HeightMM    *float64 `json:"height_mm,omitempty" binding:"omitempty,gt=0" example:"2390"`
	MaxWeightKG *float64 `json:"max_weight_kg,omitempty" binding:"omitempty,gt=0" example:"28200"`
}

type CreatePlanItem struct {
	ProductSKU    *string `json:"product_sku,omitempty" binding:"omitempty,max=50" example:"TV55-001"`
	Label         *string `json:"label,omitempty" binding:"omitempty,max=100" example:"TV LED 55 inch"`
	LengthMM      float64 `json:"length_mm" binding:"required,gt=0" example:"1300"`
	WidthMM       float64 `json:"width_mm" binding:"required,gt=0" example:"800"`
	HeightMM      float64 `json:"height_mm" binding:"required,gt=0" example:"200"`
	WeightKG      float64 `json:"weight_kg" binding:"required,gt=0" example:"25.5"`
	Quantity      int     `json:"quantity" binding:"required,gt=0" example:"120"`
	AllowRotation *bool   `json:"allow_rotation,omitempty" binding:"-" example:"true"`
	ColorHex      *string `json:"color_hex,omitempty" binding:"omitempty,len=7,startswith=#" example:"#ff5733"`
}

type CreatePlanResponse struct {
	PlanID           string             `json:"plan_id" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	PlanCode         string             `json:"plan_code" example:"PLAN-20251209-104500"`
	Title            string             `json:"title" example:"Ekspor Jepang - Desember 2025"`
	Status           string             `json:"status" example:"DRAFT"` // DRAFT, IN_PROGRESS, COMPLETED, PARTIAL, FAILED, CANCELLED
	TotalItems       int                `json:"total_items" example:"160"`
	TotalWeightKG    float64            `json:"total_weight_kg" example:"8520.0"`
	TotalVolumeM3    float64            `json:"total_volume_m3" example:"58.2"`
	CalculationJobID *string            `json:"calculation_job_id,omitempty" example:"calc_8f9e2a1b"`
	Calculation      *CalculationResult `json:"calculation,omitempty"`
	CreatedAt        string             `json:"created_at" example:"2025-12-09T10:45:00Z"`
}

type PlanDetailResponse struct {
	PlanID      string             `json:"plan_id"`
	PlanCode    string             `json:"plan_code"`
	Title       string             `json:"title"`
	Notes       *string            `json:"notes,omitempty"`
	Status      string             `json:"status" example:"COMPLETED"` // DRAFT, IN_PROGRESS, COMPLETED, PARTIAL, FAILED, CANCELLED
	Container   PlanContainerInfo  `json:"container"`
	Stats       PlanStats          `json:"stats"`
	Items       []PlanItemDetail   `json:"items"`
	Calculation *CalculationResult `json:"calculation,omitempty"`
	CreatedBy   UserSummary        `json:"created_by"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
}

type PlanContainerInfo struct {
	ContainerID *string `json:"container_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	LengthMM    float64 `json:"length_mm"`
	WidthMM     float64 `json:"width_mm"`
	HeightMM    float64 `json:"height_mm"`
	MaxWeightKG float64 `json:"max_weight_kg"`
	VolumeM3    float64 `json:"volume_m3"`
}

type PlanStats struct {
	TotalItems           int     `json:"total_items"`
	TotalWeightKG        float64 `json:"total_weight_kg"`
	TotalVolumeM3        float64 `json:"total_volume_m3"`
	VolumeUtilizationPct float64 `json:"volume_utilization_pct"`
	WeightUtilizationPct float64 `json:"weight_utilization_pct"`
}

type PlanItemDetail struct {
	ItemID        string  `json:"item_id"`
	ProductSKU    *string `json:"product_sku,omitempty"`
	Label         *string `json:"label,omitempty"`
	LengthMM      float64 `json:"length_mm"`
	WidthMM       float64 `json:"width_mm"`
	HeightMM      float64 `json:"height_mm"`
	WeightKG      float64 `json:"weight_kg"`
	Quantity      int     `json:"quantity"`
	TotalWeightKG float64 `json:"total_weight_kg"`
	TotalVolumeM3 float64 `json:"total_volume_m3"`
	AllowRotation bool    `json:"allow_rotation"`
	StackingLimit int     `json:"stacking_limit"`
	ColorHex      *string `json:"color_hex,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type CalculationResult struct {
	JobID             string            `json:"job_id"`
	Status            string            `json:"status"` // queued | running | completed | failed
	Algorithm         string            `json:"algorithm" example:"maxrects-bssf"`
	CalculatedAt      *string           `json:"calculated_at,omitempty"`
	DurationMs        int64             `json:"duration_ms,omitempty"`
	EfficiencyScore   float64           `json:"efficiency_score,omitempty"`
	VolumeUtilization float64           `json:"volume_utilization_pct,omitempty"`
	VisualizationURL  string            `json:"visualization_url" example:"/visualizer?plan=f47ac10b-..."`
	Placements        []PlacementDetail `json:"placements,omitempty"`
}

type PlacementDetail struct {
	PlacementID string  `json:"placement_id"`
	ItemID      string  `json:"item_id"`
	PositionX   float64 `json:"pos_x"`
	PositionY   float64 `json:"pos_y"`
	PositionZ   float64 `json:"pos_z"`
	Rotation    int     `json:"rotation"`
	StepNumber  int     `json:"step_number"`
}

type PlanListItem struct {
	PlanID               string  `json:"plan_id"`
	PlanCode             string  `json:"plan_code"`
	Title                string  `json:"title"`
	Status               string  `json:"status" example:"DRAFT"`
	TotalItems           int     `json:"total_items"`
	TotalWeightKG        float64 `json:"total_weight_kg"`
	VolumeUtilizationPct float64 `json:"volume_utilization_pct,omitempty"`
	CreatedBy            string  `json:"created_by"`
	CreatedAt            string  `json:"created_at"`
}

type UpdatePlanRequest struct {
	Status    *string              `json:"status,omitempty" binding:"omitempty,oneof=DRAFT IN_PROGRESS COMPLETED PARTIAL FAILED CANCELLED"`
	Container *CreatePlanContainer `json:"container,omitempty"`
}

type AddPlanItemRequest struct {
	CreatePlanItem
}

type UpdatePlanItemRequest struct {
	Label         *string  `json:"label,omitempty" binding:"omitempty,max=100"`
	LengthMM      *float64 `json:"length_mm,omitempty" binding:"omitempty,gt=0"`
	WidthMM       *float64 `json:"width_mm,omitempty" binding:"omitempty,gt=0"`
	HeightMM      *float64 `json:"height_mm,omitempty" binding:"omitempty,gt=0"`
	WeightKG      *float64 `json:"weight_kg,omitempty" binding:"omitempty,gt=0"`
	Quantity      *int     `json:"quantity,omitempty" binding:"omitempty,gt=0"`
	AllowRotation *bool    `json:"allow_rotation,omitempty"`
	ColorHex      *string  `json:"color_hex,omitempty" binding:"omitempty,len=7,startswith=#"`
}
