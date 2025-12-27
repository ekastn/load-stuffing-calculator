package service

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
}

func workspaceIDFromContext(ctx context.Context) (*uuid.UUID, error) {
	workspaceID, ok := auth.WorkspaceIDFromContext(ctx)
	if !ok || workspaceID == "" {
		return nil, nil
	}
	wid, err := uuid.Parse(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id: %w", err)
	}
	return &wid, nil
}

func toFloat(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
