package handler

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateBarcode(t *testing.T) {
	planID := uuid.MustParse("a3f2e8b1-c4d9-4f2a-b1c8-3d9e8f7a6b5c")
	itemID := uuid.MustParse("c4d9f2a3-8b1c-4e5f-9a2b-6d7e8f9a0b1c")
	step := 1

	expected := "PLAN-a3f2e8b1-STEP-001-c4d9f2a3"
	actual := generateBarcode(planID, step, itemID)

	assert.Equal(t, expected, actual)
}

func TestParseBarcode(t *testing.T) {
	tests := []struct {
		name    string
		barcode string
		want    *ParsedBarcode
	}{
		{
			name:    "valid barcode",
			barcode: "PLAN-a3f2e8b1-STEP-001-c4d9f2a3",
			want: &ParsedBarcode{
				PlanID:     "a3f2e8b1",
				StepNumber: 1,
				ItemID:     "c4d9f2a3",
			},
		},
		{
			name:    "invalid prefix",
			barcode: "ITEM-a3f2e8b1-STEP-001-c4d9f2a3",
			want:    nil,
		},
		{
			name:    "missing step part",
			barcode: "PLAN-a3f2e8b1-PART-001-c4d9f2a3",
			want:    nil,
		},
		{
			name:    "invalid length",
			barcode: "PLAN-a3f2e8b1-STEP-001",
			want:    nil,
		},
		{
			name:    "invalid step number",
			barcode: "PLAN-a3f2e8b1-STEP-ABC-c4d9f2a3",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBarcode(tt.barcode)
			assert.Equal(t, tt.want, got)
		})
	}
}
