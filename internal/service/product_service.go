package service

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error)
	ListProducts(ctx context.Context, page, limit int32) ([]dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	q store.Querier
}

func NewProductService(q store.Querier) ProductService {
	return &productService{q: q}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Founders can create either:
	// - Global presets (workspace_id = NULL) when no override is provided.
	// - Workspace-scoped records when ?workspace_id= is provided.
	if isFounder(ctx) {
		workspaceID = overrideWorkspaceID
	}

	if workspaceID == nil && !isFounder(ctx) {
		return nil, fmt.Errorf("workspace id is required")
	}

	product, err := s.q.CreateProduct(ctx, store.CreateProductParams{
		WorkspaceID: workspaceID,
		Name:        req.Name,
		Sku:         req.SKU,
		LengthMm:    toNumeric(req.LengthMM),
		WidthMm:     toNumeric(req.WidthMM),
		HeightMm:    toNumeric(req.HeightMM),
		WeightKg:    toNumeric(req.WeightKG),
		ColorHex:    req.ColorHex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return mapProductToResponse(product), nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product id: %w", err)
	}

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		product, err := s.q.GetProductAny(ctx, productID)
		if err != nil {
			return nil, err
		}
		return mapProductToResponse(product), nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}

	product, err := s.q.GetProduct(ctx, store.GetProductParams{ProductID: productID, WorkspaceID: workspaceID})
	if err != nil {
		return nil, err
	}

	return mapProductToResponse(product), nil
}

func (s *productService) ListProducts(ctx context.Context, page, limit int32) ([]dto.ProductResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		products, err := s.q.ListProductsAll(ctx, store.ListProductsAllParams{Limit: limit, Offset: offset})
		if err != nil {
			return nil, err
		}
		var result []dto.ProductResponse
		for _, p := range products {
			result = append(result, *mapProductToResponse(p))
		}
		return result, nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}

	products, err := s.q.ListProducts(ctx, store.ListProductsParams{WorkspaceID: workspaceID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}

	var result []dto.ProductResponse
	for _, p := range products {
		result = append(result, *mapProductToResponse(p))
	}

	return result, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product id: %w", err)
	}

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		err = s.q.UpdateProductAny(ctx, store.UpdateProductAnyParams{
			ProductID: productID,
			Name:      req.Name,
			Sku:       req.SKU,
			LengthMm:  toNumeric(req.LengthMM),
			WidthMm:   toNumeric(req.WidthMM),
			HeightMm:  toNumeric(req.HeightMM),
			WeightKg:  toNumeric(req.WeightKG),
			ColorHex:  req.ColorHex,
		})
		if err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}
		return nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}
	if workspaceID == nil {
		return fmt.Errorf("workspace id is required")
	}

	err = s.q.UpdateProduct(ctx, store.UpdateProductParams{
		ProductID:   productID,
		WorkspaceID: workspaceID,
		Name:        req.Name,
		Sku:         req.SKU,
		LengthMm:    toNumeric(req.LengthMM),
		WidthMm:     toNumeric(req.WidthMM),
		HeightMm:    toNumeric(req.HeightMM),
		WeightKg:    toNumeric(req.WeightKG),
		ColorHex:    req.ColorHex,
	})
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product id: %w", err)
	}

	overrideWorkspaceID, err := workspaceOverrideIDFromContext(ctx)
	if err != nil {
		return err
	}

	if isFounder(ctx) && overrideWorkspaceID == nil {
		err = s.q.DeleteProductAny(ctx, productID)
		if err != nil {
			return fmt.Errorf("failed to delete product: %w", err)
		}
		return nil
	}

	workspaceID, err := workspaceIDFromContext(ctx)
	if err != nil {
		return err
	}
	if overrideWorkspaceID != nil {
		workspaceID = overrideWorkspaceID
	}
	if workspaceID == nil {
		return fmt.Errorf("workspace id is required")
	}

	err = s.q.DeleteProduct(ctx, store.DeleteProductParams{ProductID: productID, WorkspaceID: workspaceID})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func mapProductToResponse(p store.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:       p.ProductID.String(),
		Name:     p.Name,
		SKU:      p.Sku,
		LengthMM: toFloat(p.LengthMm),
		WidthMM:  toFloat(p.WidthMm),
		HeightMM: toFloat(p.HeightMm),
		WeightKG: toFloat(p.WeightKg),
		ColorHex: p.ColorHex,
	}
}
