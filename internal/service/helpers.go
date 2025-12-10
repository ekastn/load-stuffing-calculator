package service

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func toNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
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
