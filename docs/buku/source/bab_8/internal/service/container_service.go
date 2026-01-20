package service

import (
	"context"
	"fmt"

	"load-stuffing-calculator/internal/store"

	"github.com/google/uuid"
)

// ContainerService menangani operasi bisnis terkait Container.
// Menerima store.Querier sebagai dependency untuk akses database.
type ContainerService struct {
	store store.Querier
}

// NewContainerService membuat instance baru ContainerService.
// Pattern ini disebut "constructor function" di Go.
func NewContainerService(store store.Querier) *ContainerService {
	return &ContainerService{store: store}
}

// GetByID mengambil satu container berdasarkan ID.
// Mengembalikan pointer ke Container atau error jika tidak ditemukan.
func (s *ContainerService) GetByID(ctx context.Context, id uuid.UUID) (*store.Container, error) {
	container, err := s.store.GetContainer(ctx, id)
	if err != nil {
		// Wrap error dengan konteks untuk debugging yang lebih mudah
		return nil, fmt.Errorf("get container %s: %w", id, err)
	}
	return &container, nil
}

// List mengambil semua container yang tersedia.
func (s *ContainerService) List(ctx context.Context) ([]store.Container, error) {
	containers, err := s.store.ListContainers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}
	return containers, nil
}

// Create membuat container baru dengan validasi bisnis.
// Validasi memastikan dimensi dan berat bernilai positif.
func (s *ContainerService) Create(ctx context.Context, name string, length, width, height, maxWeight float64) (*store.Container, error) {
	// Validasi bisnis: dimensi harus positif
	// Ini adalah aturan domain, bukan constraint database
	if length <= 0 || width <= 0 || height <= 0 {
		return nil, fmt.Errorf("dimensions must be positive")
	}
	if maxWeight <= 0 {
		return nil, fmt.Errorf("max weight must be positive")
	}

	// Panggil repository untuk menyimpan ke database
	container, err := s.store.CreateContainer(ctx, store.CreateContainerParams{
		Name:        name,
		LengthMm:    length,
		WidthMm:     width,
		HeightMm:    height,
		MaxWeightKg: maxWeight,
	})
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}
	return &container, nil
}

// Update memperbarui container yang sudah ada.
func (s *ContainerService) Update(ctx context.Context, id uuid.UUID, name string, length, width, height, maxWeight float64) (*store.Container, error) {
	// Validasi bisnis
	if length <= 0 || width <= 0 || height <= 0 {
		return nil, fmt.Errorf("dimensions must be positive")
	}
	if maxWeight <= 0 {
		return nil, fmt.Errorf("max weight must be positive")
	}

	container, err := s.store.UpdateContainer(ctx, store.UpdateContainerParams{
		ID:          id,
		Name:        name,
		LengthMm:    length,
		WidthMm:     width,
		HeightMm:    height,
		MaxWeightKg: maxWeight,
	})
	if err != nil {
		return nil, fmt.Errorf("update container %s: %w", id, err)
	}
	return &container, nil
}

// Delete menghapus container berdasarkan ID.
func (s *ContainerService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteContainer(ctx, id)
	if err != nil {
		return fmt.Errorf("delete container %s: %w", id, err)
	}
	return nil
}
