package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/ekastn/load-stuffing-calculator/internal/store"
	"github.com/google/uuid"
)

func TestProductService_CreateProduct(t *testing.T) {
	name := "Item 1"
	length := 100.0
	color := "#abcdef"

	tests := []struct {
		name        string
		req         dto.CreateProductRequest
		createErr   error
		createResp  store.Product
		wantErr     bool
	}{
		{
			name: "success",
			req: dto.CreateProductRequest{
				Name:     name,
				LengthMM: length,
				WidthMM:  50.0,
				HeightMM: 20.0,
				WeightKG: 1.5,
				ColorHex: &color,
			},
			createResp: store.Product{
				ProductID: uuid.New(),
				Name:      name,
				LengthMm:  toNumeric(length),
				WidthMm:   toNumeric(50.0),
				HeightMm:  toNumeric(20.0),
				WeightKg:  toNumeric(1.5),
				ColorHex:  &color,
			},
			wantErr: false,
		},
		{
			name: "db_error",
			req: dto.CreateProductRequest{
				Name: name,
			},
			createErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				CreateProductFunc: func(ctx context.Context, arg store.CreateProductParams) (store.Product, error) {
					return tt.createResp, tt.createErr
				},
			}

			s := service.NewProductService(mockQ)
			resp, err := s.CreateProduct(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.Name != tt.req.Name {
					t.Errorf("Name = %v, want %v", resp.Name, tt.req.Name)
				}
				if resp.LengthMM != tt.req.LengthMM {
					t.Errorf("Length mismatch")
				}
			}
		})
	}
}

func TestProductService_GetProduct(t *testing.T) {
	id := uuid.New()
	name := "Item 2"

	tests := []struct {
		name    string
		id      string
		getErr  error
		getResp store.Product
		wantErr bool
	}{
		{
			name: "success",
			id:   id.String(),
			getResp: store.Product{
				ProductID: id,
				Name:      name,
			},
			wantErr: false,
		},
		{
			name:    "not_found",
			id:      id.String(),
			getErr:  fmt.Errorf("not found"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				GetProductFunc: func(ctx context.Context, productID uuid.UUID) (store.Product, error) {
					if productID.String() != tt.id {
						return store.Product{}, fmt.Errorf("id mismatch")
					}
					return tt.getResp, tt.getErr
				},
			}

			s := service.NewProductService(mockQ)
			resp, err := s.GetProduct(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.ID != tt.id {
					t.Errorf("ID = %v, want %v", resp.ID, tt.id)
				}
			}
		})
	}
}

func TestProductService_ListProducts(t *testing.T) {
	tests := []struct {
		name         string
		page, limit  int32
		listResp     []store.Product
		listErr      error
		wantErr      bool
		wantLen      int
	}{
		{
			name: "success",
			page: 1, limit: 10,
			listResp: []store.Product{
				{ProductID: uuid.New(), Name: "p1"},
				{ProductID: uuid.New(), Name: "p2"},
			},
			wantLen: 2,
		},
		{
			name:    "db_error",
			listErr: fmt.Errorf("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				ListProductsFunc: func(ctx context.Context, arg store.ListProductsParams) ([]store.Product, error) {
					return tt.listResp, tt.listErr
				},
			}

			s := service.NewProductService(mockQ)
			resp, err := s.ListProducts(context.Background(), tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(resp) != tt.wantLen {
					t.Errorf("Resp len = %v, want %v", len(resp), tt.wantLen)
				}
			}
		})
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	id := uuid.New()
	name := "updated_p"

	tests := []struct {
		name      string
		id        string
		req       dto.UpdateProductRequest
		updateErr error
		wantErr   bool
	}{
		{
			name: "success",
			id:   id.String(),
			req: dto.UpdateProductRequest{
				Name:     name,
				LengthMM: 110,
				WidthMM:  60,
				HeightMM: 30,
				WeightKG: 2,
				ColorHex: stringPtr("#123456"),
			},
			wantErr: false,
		},
		{
			name:      "db_error",
			id:        id.String(),
			req:       dto.UpdateProductRequest{Name: name},
			updateErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				UpdateProductFunc: func(ctx context.Context, arg store.UpdateProductParams) error {
					return tt.updateErr
				},
			}

			s := service.NewProductService(mockQ)
			err := s.UpdateProduct(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name      string
		id        string
		deleteErr error
		wantErr   bool
	}{
		{
			name: "success",
			id:   id.String(),
			wantErr: false,
		},
		{
			name:      "db_error",
			id:        id.String(),
			deleteErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				DeleteProductFunc: func(ctx context.Context, productID uuid.UUID) error {
					return tt.deleteErr
				},
			}

			s := service.NewProductService(mockQ)
			err := s.DeleteProduct(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}