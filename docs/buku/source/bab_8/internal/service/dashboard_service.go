package service

import (
	"context"
	"fmt"

	"load-stuffing-calculator/internal/store"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DashboardService struct {
	store store.Querier
}

func NewDashboardService(store store.Querier) *DashboardService {
	return &DashboardService{store: store}
}

func (s *DashboardService) GetStats(ctx context.Context, userID uuid.UUID) (*store.GetDashboardStatsRow, error) {
	stats, err := s.store.GetDashboardStats(ctx, pgtype.UUID{
		Bytes: userID,
		Valid: true,
	})
	if err != nil {
		return nil, fmt.Errorf("get dashboard stats: %w", err)
	}

	return &stats, nil
}
