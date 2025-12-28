package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
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
		name       string
		req        dto.CreateProductRequest
		createErr  error
		createResp store.Product
		wantErr    bool
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
		{
			name:    "trial_no_workspace_forbidden",
			req:     dto.CreateProductRequest{Name: name},
			wantErr: true,
		},
		{
			name: "founder_no_override_creates_global_preset",
			req:  dto.CreateProductRequest{Name: name},
			createResp: store.Product{
				ProductID: uuid.New(),
				Name:      name,
			},
			wantErr: false,
		},
		{
			name: "founder_with_override_creates_scoped",
			req:  dto.CreateProductRequest{Name: name},
			createResp: store.Product{
				ProductID: uuid.New(),
				Name:      name,
			},
			wantErr: false,
		},
		{
			name:    "founder_invalid_override_errors",
			req:     dto.CreateProductRequest{Name: name},
			wantErr: true,
		},
		{
			name: "founder_override_requires_no_token_workspace",
			req:  dto.CreateProductRequest{Name: name},
			createResp: store.Product{
				ProductID: uuid.New(),
				Name:      name,
			},
			wantErr: false,
		},
		{
			name:    "non_founder_missing_workspace_errors",
			req:     dto.CreateProductRequest{Name: name},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sawWorkspaceID := (*uuid.UUID)(nil)
			mockQ := &MockQuerier{
				CreateProductFunc: func(ctx context.Context, arg store.CreateProductParams) (store.Product, error) {
					sawWorkspaceID = arg.WorkspaceID
					switch tt.name {
					case "trial_no_workspace_forbidden", "non_founder_missing_workspace_errors", "founder_invalid_override_errors":
						return store.Product{}, fmt.Errorf("unexpected db call")
					}
					return tt.createResp, tt.createErr
				},
			}

			s := service.NewProductService(mockQ)

			workspaceID := uuid.New()
			overrideWorkspaceID := uuid.New()

			ctx := ctxWithWorkspaceID(workspaceID)
			switch tt.name {
			case "trial_no_workspace_forbidden":
				ctx = context.Background()
			case "founder_no_override_creates_global_preset":
				ctx = ctxWithRole("founder")
			case "founder_with_override_creates_scoped":
				ctx = auth.WithWorkspaceOverrideID(ctxWithRoleAndWorkspace("founder", workspaceID), overrideWorkspaceID.String())
			case "founder_invalid_override_errors":
				ctx = auth.WithWorkspaceOverrideID(ctxWithRoleAndWorkspace("founder", workspaceID), "not-a-uuid")
			case "founder_override_requires_no_token_workspace":
				ctx = auth.WithWorkspaceOverrideID(ctxWithRole("founder"), overrideWorkspaceID.String())
			case "non_founder_missing_workspace_errors":
				ctx = ctxWithRole("admin")
			}

			resp, err := s.CreateProduct(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Ensure we don't call the DB for validation-only failures.
			switch tt.name {
			case "trial_no_workspace_forbidden", "non_founder_missing_workspace_errors", "founder_invalid_override_errors":
				if sawWorkspaceID != nil {
					t.Fatalf("unexpected db call")
				}
				return
			}

			// For success and db-level errors, validate what workspace was used.
			switch tt.name {
			case "founder_no_override_creates_global_preset":
				if sawWorkspaceID != nil {
					t.Fatalf("expected WorkspaceID to be nil, got %v", *sawWorkspaceID)
				}
			case "founder_with_override_creates_scoped", "founder_override_requires_no_token_workspace":
				if sawWorkspaceID == nil || *sawWorkspaceID != overrideWorkspaceID {
					t.Fatalf("expected WorkspaceID to be override %v, got %v", overrideWorkspaceID, sawWorkspaceID)
				}
			case "success", "db_error":
				if sawWorkspaceID == nil || *sawWorkspaceID != workspaceID {
					t.Fatalf("expected WorkspaceID to be token workspace %v, got %v", workspaceID, sawWorkspaceID)
				}
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
			workspaceID := uuid.New()
			ctx := ctxWithWorkspaceID(workspaceID)

			mockQ := &MockQuerier{
				GetProductFunc: func(ctx context.Context, arg store.GetProductParams) (store.Product, error) {
					if arg.ProductID.String() != tt.id {
						return store.Product{}, fmt.Errorf("id mismatch")
					}
					if arg.WorkspaceID == nil || *arg.WorkspaceID != workspaceID {
						return store.Product{}, fmt.Errorf("workspace mismatch")
					}
					return tt.getResp, tt.getErr
				},
			}

			s := service.NewProductService(mockQ)
			resp, err := s.GetProduct(ctx, tt.id)

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

	t.Run("founder_no_override_uses_any", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		called := false

		mockQ := &MockQuerier{
			GetProductAnyFunc: func(ctx context.Context, productID uuid.UUID) (store.Product, error) {
				called = true
				if productID != id {
					return store.Product{}, fmt.Errorf("id mismatch")
				}
				return store.Product{ProductID: id, Name: name}, nil
			},
			GetProductFunc: func(ctx context.Context, arg store.GetProductParams) (store.Product, error) {
				return store.Product{}, fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewProductService(mockQ)
		resp, err := s.GetProduct(ctx, id.String())
		if err != nil {
			t.Fatalf("GetProduct() error = %v", err)
		}
		if !called {
			t.Fatalf("expected GetProductAny to be called")
		}
		if resp.ID != id.String() {
			t.Fatalf("ID = %v, want %v", resp.ID, id.String())
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())

		mockQ := &MockQuerier{
			GetProductFunc: func(ctx context.Context, arg store.GetProductParams) (store.Product, error) {
				if arg.ProductID != id {
					return store.Product{}, fmt.Errorf("id mismatch")
				}
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return store.Product{}, fmt.Errorf("workspace mismatch")
				}
				return store.Product{ProductID: id, Name: name}, nil
			},
			GetProductAnyFunc: func(ctx context.Context, productID uuid.UUID) (store.Product, error) {
				return store.Product{}, fmt.Errorf("unexpected any call")
			},
		}

		s := service.NewProductService(mockQ)
		resp, err := s.GetProduct(ctx, id.String())
		if err != nil {
			t.Fatalf("GetProduct() error = %v", err)
		}
		if resp.ID != id.String() {
			t.Fatalf("ID = %v, want %v", resp.ID, id.String())
		}
	})
}

func TestProductService_ListProducts(t *testing.T) {
	tests := []struct {
		name        string
		page, limit int32
		listResp    []store.Product
		listErr     error
		wantErr     bool
		wantLen     int
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
			name: "trial_no_workspace_returns_global_presets",
			page: 1, limit: 10,
			listResp: []store.Product{
				{ProductID: uuid.New(), Name: "global1"},
			},
			wantLen: 1,
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
					if tt.name == "trial_no_workspace_returns_global_presets" && arg.WorkspaceID != nil {
						return nil, fmt.Errorf("expected nil workspace id")
					}
					return tt.listResp, tt.listErr
				},
			}

			s := service.NewProductService(mockQ)
			ctx := ctxWithWorkspaceID(uuid.New())
			if tt.name == "trial_no_workspace_returns_global_presets" {
				ctx = context.Background()
			}
			resp, err := s.ListProducts(ctx, tt.page, tt.limit)

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

	t.Run("founder_no_override_uses_all", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		page, limit := int32(2), int32(12)
		offset := (page - 1) * limit
		called := false

		mockQ := &MockQuerier{
			ListProductsAllFunc: func(ctx context.Context, arg store.ListProductsAllParams) ([]store.Product, error) {
				called = true
				if arg.Limit != limit || arg.Offset != offset {
					return nil, fmt.Errorf("limit/offset mismatch")
				}
				return []store.Product{{ProductID: uuid.New(), Name: "global"}}, nil
			},
			ListProductsFunc: func(ctx context.Context, arg store.ListProductsParams) ([]store.Product, error) {
				return nil, fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewProductService(mockQ)
		resp, err := s.ListProducts(ctx, page, limit)
		if err != nil {
			t.Fatalf("ListProducts() error = %v", err)
		}
		if !called {
			t.Fatalf("expected ListProductsAll to be called")
		}
		if len(resp) != 1 {
			t.Fatalf("resp len = %v, want %v", len(resp), 1)
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())
		page, limit := int32(1), int32(10)
		offset := (page - 1) * limit

		mockQ := &MockQuerier{
			ListProductsFunc: func(ctx context.Context, arg store.ListProductsParams) ([]store.Product, error) {
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return nil, fmt.Errorf("workspace mismatch")
				}
				if arg.Limit != limit || arg.Offset != offset {
					return nil, fmt.Errorf("limit/offset mismatch")
				}
				return []store.Product{{ProductID: uuid.New(), Name: "scoped"}}, nil
			},
			ListProductsAllFunc: func(ctx context.Context, arg store.ListProductsAllParams) ([]store.Product, error) {
				return nil, fmt.Errorf("unexpected all call")
			},
		}

		s := service.NewProductService(mockQ)
		resp, err := s.ListProducts(ctx, page, limit)
		if err != nil {
			t.Fatalf("ListProducts() error = %v", err)
		}
		if len(resp) != 1 {
			t.Fatalf("resp len = %v, want %v", len(resp), 1)
		}
	})
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
			name:    "trial_no_workspace_forbidden",
			id:      id.String(),
			req:     dto.UpdateProductRequest{Name: name},
			wantErr: true,
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
					if tt.name == "trial_no_workspace_forbidden" {
						return fmt.Errorf("unexpected db call")
					}
					return tt.updateErr
				},
			}

			s := service.NewProductService(mockQ)
			ctx := ctxWithWorkspaceID(uuid.New())
			if tt.name == "trial_no_workspace_forbidden" {
				ctx = context.Background()
			}
			err := s.UpdateProduct(ctx, tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("founder_no_override_uses_any", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		req := dto.UpdateProductRequest{Name: "founder_update"}
		called := false

		mockQ := &MockQuerier{
			UpdateProductAnyFunc: func(ctx context.Context, arg store.UpdateProductAnyParams) error {
				called = true
				if arg.ProductID != id {
					return fmt.Errorf("id mismatch")
				}
				if arg.Name != req.Name {
					return fmt.Errorf("name mismatch")
				}
				return nil
			},
			UpdateProductFunc: func(ctx context.Context, arg store.UpdateProductParams) error {
				return fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewProductService(mockQ)
		err := s.UpdateProduct(ctx, id.String(), req)
		if err != nil {
			t.Fatalf("UpdateProduct() error = %v", err)
		}
		if !called {
			t.Fatalf("expected UpdateProductAny to be called")
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())
		req := dto.UpdateProductRequest{Name: "founder_update_scoped"}

		mockQ := &MockQuerier{
			UpdateProductFunc: func(ctx context.Context, arg store.UpdateProductParams) error {
				if arg.ProductID != id {
					return fmt.Errorf("id mismatch")
				}
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return fmt.Errorf("workspace mismatch")
				}
				if arg.Name != req.Name {
					return fmt.Errorf("name mismatch")
				}
				return nil
			},
			UpdateProductAnyFunc: func(ctx context.Context, arg store.UpdateProductAnyParams) error {
				return fmt.Errorf("unexpected any call")
			},
		}

		s := service.NewProductService(mockQ)
		err := s.UpdateProduct(ctx, id.String(), req)
		if err != nil {
			t.Fatalf("UpdateProduct() error = %v", err)
		}
	})
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
			name:    "success",
			id:      id.String(),
			wantErr: false,
		},
		{
			name:    "trial_no_workspace_forbidden",
			id:      id.String(),
			wantErr: true,
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
			workspaceID := uuid.New()
			ctx := ctxWithWorkspaceID(workspaceID)

			mockQ := &MockQuerier{
				DeleteProductFunc: func(ctx context.Context, arg store.DeleteProductParams) error {
					if arg.ProductID.String() != tt.id {
						return fmt.Errorf("id mismatch")
					}
					if tt.name == "trial_no_workspace_forbidden" {
						return fmt.Errorf("unexpected db call")
					}
					if arg.WorkspaceID == nil || *arg.WorkspaceID != workspaceID {
						return fmt.Errorf("workspace mismatch")
					}
					return tt.deleteErr
				},
			}

			s := service.NewProductService(mockQ)
			if tt.name == "trial_no_workspace_forbidden" {
				ctx = context.Background()
			}
			err := s.DeleteProduct(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("founder_no_override_uses_any", func(t *testing.T) {
		ctx := ctxWithRole("founder")
		called := false

		mockQ := &MockQuerier{
			DeleteProductAnyFunc: func(ctx context.Context, productID uuid.UUID) error {
				called = true
				if productID != id {
					return fmt.Errorf("id mismatch")
				}
				return nil
			},
			DeleteProductFunc: func(ctx context.Context, arg store.DeleteProductParams) error {
				return fmt.Errorf("unexpected scoped call")
			},
		}

		s := service.NewProductService(mockQ)
		err := s.DeleteProduct(ctx, id.String())
		if err != nil {
			t.Fatalf("DeleteProduct() error = %v", err)
		}
		if !called {
			t.Fatalf("expected DeleteProductAny to be called")
		}
	})

	t.Run("founder_with_override_uses_scoped", func(t *testing.T) {
		workspaceID := uuid.New()
		overrideWorkspaceID := uuid.New()
		ctx := ctxWithRoleAndWorkspace("founder", workspaceID)
		ctx = auth.WithWorkspaceOverrideID(ctx, overrideWorkspaceID.String())

		mockQ := &MockQuerier{
			DeleteProductFunc: func(ctx context.Context, arg store.DeleteProductParams) error {
				if arg.ProductID != id {
					return fmt.Errorf("id mismatch")
				}
				if arg.WorkspaceID == nil || *arg.WorkspaceID != overrideWorkspaceID {
					return fmt.Errorf("workspace mismatch")
				}
				return nil
			},
			DeleteProductAnyFunc: func(ctx context.Context, productID uuid.UUID) error {
				return fmt.Errorf("unexpected any call")
			},
		}

		s := service.NewProductService(mockQ)
		err := s.DeleteProduct(ctx, id.String())
		if err != nil {
			t.Fatalf("DeleteProduct() error = %v", err)
		}
	})
}
