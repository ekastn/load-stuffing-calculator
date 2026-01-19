package service

import (
	"context"
	"fmt"

	"load-stuffing-calculator/internal/store"

	"github.com/google/uuid"
)

// ProductService menangani operasi bisnis terkait Product.
// Menerima store.Querier sebagai dependency untuk akses database.
type ProductService struct {
	store store.Querier
}

// NewProductService membuat instance baru ProductService.
func NewProductService(store store.Querier) *ProductService {
	return &ProductService{store: store}
}

// GetByID mengambil satu product berdasarkan ID.
func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*store.Product, error) {
	product, err := s.store.GetProduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product %s: %w", id, err)
	}
	return &product, nil
}

// List mengambil semua product yang tersedia.
func (s *ProductService) List(ctx context.Context) ([]store.Product, error) {
	products, err := s.store.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	return products, nil
}

// Create membuat product baru dengan validasi bisnis.
func (s *ProductService) Create(ctx context.Context, label, sku string, length, width, height, weight float64) (*store.Product, error) {
	// Validasi bisnis: dimensi harus positif
	if length <= 0 || width <= 0 || height <= 0 {
		return nil, fmt.Errorf("dimensions must be positive")
	}
	if weight <= 0 {
		return nil, fmt.Errorf("weight must be positive")
	}

	product, err := s.store.CreateProduct(ctx, store.CreateProductParams{
		Label:    label,
		Sku:      sku,
		LengthMm: length,
		WidthMm:  width,
		HeightMm: height,
		WeightKg: weight,
	})
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	return &product, nil
}

// Update memperbarui product yang sudah ada.
func (s *ProductService) Update(ctx context.Context, id uuid.UUID, label, sku string, length, width, height, weight float64) (*store.Product, error) {
	// Validasi bisnis
	if length <= 0 || width <= 0 || height <= 0 {
		return nil, fmt.Errorf("dimensions must be positive")
	}
	if weight <= 0 {
		return nil, fmt.Errorf("weight must be positive")
	}

	product, err := s.store.UpdateProduct(ctx, store.UpdateProductParams{
		ID:       id,
		Label:    label,
		Sku:      sku,
		LengthMm: length,
		WidthMm:  width,
		HeightMm: height,
		WeightKg: weight,
	})
	if err != nil {
		return nil, fmt.Errorf("update product %s: %w", id, err)
	}
	return &product, nil
}

// Delete menghapus product berdasarkan ID.
func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.store.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("delete product %s: %w", id, err)
	}
	return nil
}
