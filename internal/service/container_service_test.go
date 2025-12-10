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

func TestContainerService_CreateContainer(t *testing.T) {
	name := "20ft"
	length := 6058.0
	desc := "Standard"

	tests := []struct {
		name       string
		req        dto.CreateContainerRequest
		createErr  error
		createResp store.Container
		wantErr    bool
	}{
		{
			name: "success",
			req: dto.CreateContainerRequest{
				Name:          name,
				InnerLengthMM: length,
				InnerWidthMM:  2438.0,
				InnerHeightMM: 2591.0,
				MaxWeightKG:   28000.0,
				Description:   &desc,
			},
			createResp: store.Container{
				ContainerID:   uuid.New(),
				Name:          name,
				InnerLengthMm: toNumeric(length),
				InnerWidthMm:  toNumeric(2438.0),
				InnerHeightMm: toNumeric(2591.0),
				MaxWeightKg:   toNumeric(28000.0),
				Description:   &desc,
			},
			wantErr: false,
		},
		{
			name: "db_error",
			req: dto.CreateContainerRequest{
				Name: name,
			},
			createErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				CreateContainerFunc: func(ctx context.Context, arg store.CreateContainerParams) (store.Container, error) {
					return tt.createResp, tt.createErr
				},
			}

			s := service.NewContainerService(mockQ)
			resp, err := s.CreateContainer(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.Name != tt.req.Name {
					t.Errorf("Name = %v, want %v", resp.Name, tt.req.Name)
				}
				if resp.InnerLengthMM != tt.req.InnerLengthMM {
					t.Errorf("Length mismatch")
				}
			}
		})
	}
}

func TestContainerService_GetContainer(t *testing.T) {
	id := uuid.New()
	name := "40ft"

	tests := []struct {
		name    string
		id      string
		getErr  error
		getResp store.Container
		wantErr bool
	}{
		{
			name: "success",
			id:   id.String(),
			getResp: store.Container{
				ContainerID: id,
				Name:        name,
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
				GetContainerFunc: func(ctx context.Context, containerID uuid.UUID) (store.Container, error) {
					if containerID.String() != tt.id {
						return store.Container{}, fmt.Errorf("id mismatch")
					}
					return tt.getResp, tt.getErr
				},
			}

			s := service.NewContainerService(mockQ)
			resp, err := s.GetContainer(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContainer() error = %v, wantErr %v", err, tt.wantErr)
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

func TestContainerService_ListContainers(t *testing.T) {
	tests := []struct {
		name        string
		page, limit int32
		listResp    []store.Container
		listErr     error
		wantErr     bool
		wantLen     int
	}{
		{
			name: "success",
			page: 1, limit: 10,
			listResp: []store.Container{
				{ContainerID: uuid.New(), Name: "c1"},
				{ContainerID: uuid.New(), Name: "c2"},
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
				ListContainersFunc: func(ctx context.Context, arg store.ListContainersParams) ([]store.Container, error) {
					return tt.listResp, tt.listErr
				},
			}

			s := service.NewContainerService(mockQ)
			resp, err := s.ListContainers(context.Background(), tt.page, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListContainers() error = %v, wantErr %v", err, tt.wantErr)
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

func TestContainerService_UpdateContainer(t *testing.T) {
	id := uuid.New()
	name := "updated_c"

	tests := []struct {
		name      string
		id        string
		req       dto.UpdateContainerRequest
		updateErr error
		wantErr   bool
	}{
		{
			name:    "success",
			id:      id.String(),
			req:     dto.UpdateContainerRequest{Name: name},
			wantErr: false,
		},
		{
			name:      "db_error",
			id:        id.String(),
			req:       dto.UpdateContainerRequest{Name: name},
			updateErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				UpdateContainerFunc: func(ctx context.Context, arg store.UpdateContainerParams) error {
					return tt.updateErr
				},
			}

			s := service.NewContainerService(mockQ)
			err := s.UpdateContainer(context.Background(), tt.id, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainerService_DeleteContainer(t *testing.T) {
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
			name:      "db_error",
			id:        id.String(),
			deleteErr: fmt.Errorf("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQ := &MockQuerier{
				DeleteContainerFunc: func(ctx context.Context, containerID uuid.UUID) error {
					return tt.deleteErr
				},
			}

			s := service.NewContainerService(mockQ)
			err := s.DeleteContainer(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}