package dto

type CreatePlanRequest struct {
	ContainerID string `json:"container_id" binding:"required,uuid"`
}

type UpdatePlanRequest struct {
	ContainerID string `json:"container_id" binding:"required,uuid"`
}

type AddPlanItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type UpdatePlanItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type PlanResponse struct {
	ID            string `json:"id"`
	ContainerID   string `json:"container_id"`
	ContainerName string `json:"container_name,omitempty"`
	Status        string `json:"status"`
}

type PlanItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Label     string  `json:"label"`
	Quantity  int     `json:"quantity"`
	LengthMm  float64 `json:"length_mm"`
	WidthMm   float64 `json:"width_mm"`
	HeightMm  float64 `json:"height_mm"`
	WeightKg  float64 `json:"weight_kg"`
}

type PlacementResponse struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	Label      string  `json:"label"`
	PosX       float64 `json:"pos_x"`
	PosY       float64 `json:"pos_y"`
	PosZ       float64 `json:"pos_z"`
	Rotation   int     `json:"rotation"`
	StepNumber int     `json:"step_number"`
}

type PlanDetailResponse struct {
	ID            string              `json:"id"`
	ContainerID   string              `json:"container_id"`
	ContainerName string              `json:"container_name"`
	Status        string              `json:"status"`
	Items         []PlanItemResponse  `json:"items"`
	Placements    []PlacementResponse `json:"placements,omitempty"`
}
