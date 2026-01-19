package dto

// CreateProductRequest adalah DTO untuk request pembuatan product baru
type CreateProductRequest struct {
	Label    string  `json:"label" binding:"required,min=2,max=100"`
	SKU      string  `json:"sku" binding:"required,min=2,max=50"`
	LengthMm float64 `json:"length_mm" binding:"required,gt=0"`
	WidthMm  float64 `json:"width_mm" binding:"required,gt=0"`
	HeightMm float64 `json:"height_mm" binding:"required,gt=0"`
	WeightKg float64 `json:"weight_kg" binding:"required,gt=0"`
}

// UpdateProductRequest adalah DTO untuk request update product
type UpdateProductRequest struct {
	Label    string  `json:"label" binding:"required,min=2,max=100"`
	SKU      string  `json:"sku" binding:"required,min=2,max=50"`
	LengthMm float64 `json:"length_mm" binding:"required,gt=0"`
	WidthMm  float64 `json:"width_mm" binding:"required,gt=0"`
	HeightMm float64 `json:"height_mm" binding:"required,gt=0"`
	WeightKg float64 `json:"weight_kg" binding:"required,gt=0"`
}

// ProductResponse adalah DTO untuk response data product
type ProductResponse struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	SKU      string  `json:"sku"`
	LengthMm float64 `json:"length_mm"`
	WidthMm  float64 `json:"width_mm"`
	HeightMm float64 `json:"height_mm"`
	WeightKg float64 `json:"weight_kg"`
}
