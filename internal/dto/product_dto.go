package dto

type CreateProductRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=150"`
	SKU      *string `json:"sku"`
	LengthMM float64 `json:"length_mm" binding:"required,gt=0"`
	WidthMM  float64 `json:"width_mm" binding:"required,gt=0"`
	HeightMM float64 `json:"height_mm" binding:"required,gt=0"`
	WeightKG float64 `json:"weight_kg" binding:"required,gt=0"`
	ColorHex *string `json:"color_hex" binding:"omitempty,hexcolor"`
}

type UpdateProductRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=150"`
	SKU      *string `json:"sku"`
	LengthMM float64 `json:"length_mm" binding:"required,gt=0"`
	WidthMM  float64 `json:"width_mm" binding:"required,gt=0"`
	HeightMM float64 `json:"height_mm" binding:"required,gt=0"`
	WeightKG float64 `json:"weight_kg" binding:"required,gt=0"`
	ColorHex *string `json:"color_hex" binding:"omitempty,hexcolor"`
}

type ProductResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	SKU      *string `json:"sku,omitempty"`
	LengthMM float64 `json:"length_mm"`
	WidthMM  float64 `json:"width_mm"`
	HeightMM float64 `json:"height_mm"`
	WeightKG float64 `json:"weight_kg"`
	ColorHex *string `json:"color_hex,omitempty"`
}
