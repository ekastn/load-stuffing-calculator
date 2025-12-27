package service_test

import (
	"context"
	"fmt"

	"github.com/ekastn/load-stuffing-calculator/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ctxWithWorkspaceID(workspaceID uuid.UUID) context.Context {
	return auth.WithWorkspaceID(context.Background(), workspaceID.String())
}

// toNumeric is a test helper to convert float64 to pgtype.Numeric.
func toNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
}

// stringPtr is a test helper to return a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// boolPtr is a test helper to return a pointer to a bool.
func boolPtr(b bool) *bool {
	return &b
}

// floatPtr is a test helper to return a pointer to a float64.
func floatPtr(f float64) *float64 {
	return &f
}
