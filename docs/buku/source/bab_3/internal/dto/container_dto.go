package dto

// CreateContainerRequest adalah DTO untuk request pembuatan container baru.
// Tag `json` mendefinisikan nama field di JSON.
// Tag `binding` untuk validasi Gin—field wajib diisi dan nilainya harus positif.
type CreateContainerRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	LengthMm    float64 `json:"length_mm" binding:"required,gt=0"`
	WidthMm     float64 `json:"width_mm" binding:"required,gt=0"`
	HeightMm    float64 `json:"height_mm" binding:"required,gt=0"`
	MaxWeightKg float64 `json:"max_weight_kg" binding:"required,gt=0"`
}

// ContainerResponse adalah DTO untuk response data container.
// Berbeda dengan struct database, response bisa menyembunyikan field tertentu
// atau menambahkan computed fields.
type ContainerResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	LengthMm    float64 `json:"length_mm"`
	WidthMm     float64 `json:"width_mm"`
	HeightMm    float64 `json:"height_mm"`
	MaxWeightKg float64 `json:"max_weight_kg"`
}

// UpdateContainerRequest adalah DTO untuk request update container
type UpdateContainerRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	LengthMm    float64 `json:"length_mm" binding:"required,gt=0"`
	WidthMm     float64 `json:"width_mm" binding:"required,gt=0"`
	HeightMm    float64 `json:"height_mm" binding:"required,gt=0"`
	MaxWeightKg float64 `json:"max_weight_kg" binding:"required,gt=0"`
}
