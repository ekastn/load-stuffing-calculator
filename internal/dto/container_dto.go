package dto

type CreateContainerRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=100"`
	InnerLengthMM float64 `json:"inner_length_mm" binding:"required,gt=0"`
	InnerWidthMM  float64 `json:"inner_width_mm" binding:"required,gt=0"`
	InnerHeightMM float64 `json:"inner_height_mm" binding:"required,gt=0"`
	MaxWeightKG   float64 `json:"max_weight_kg" binding:"required,gt=0"`
	Description   *string `json:"description" binding:"omitempty,max=500"`
}

type UpdateContainerRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=100"`
	InnerLengthMM float64 `json:"inner_length_mm" binding:"required,gt=0"`
	InnerWidthMM  float64 `json:"inner_width_mm" binding:"required,gt=0"`
	InnerHeightMM float64 `json:"inner_height_mm" binding:"required,gt=0"`
	MaxWeightKG   float64 `json:"max_weight_kg" binding:"required,gt=0"`
	Description   *string `json:"description" binding:"omitempty,max=500"`
}

type ContainerResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	InnerLengthMM float64 `json:"inner_length_mm"`
	InnerWidthMM  float64 `json:"inner_width_mm"`
	InnerHeightMM float64 `json:"inner_height_mm"`
	MaxWeightKG   float64 `json:"max_weight_kg"`
	Description   *string `json:"description,omitempty"`
}
