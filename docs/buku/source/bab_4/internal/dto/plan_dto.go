package dto

type CreatePlanRequest struct {
	ContainerID string `json:"container_id" binding:"required,uuid"`
}

type AddPlanItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type PlanResponse struct {
	ID          string `json:"id"`
	ContainerID string `json:"container_id"`
	Status      string `json:"status"`
}
