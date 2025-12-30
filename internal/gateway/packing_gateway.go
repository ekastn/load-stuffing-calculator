package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// PackingGateway talks to the external packing microservice.
// It is intentionally slim (HTTP, JSON, errors) and does no domain mapping.
//
// The packing service is expected to expose:
// - POST /pack
// - GET  /health
//
// It returns positions/rotations in the same units as requested.
// In our integration we always send units="mm".
//
// Note: We keep the API stable; py3dbp options remain server-internal.
// AllowRotation is currently ignored.

type PackingGateway interface {
	Pack(ctx context.Context, req PackRequest) (*PackResponse, error)
}

type HTTPPackingGateway struct {
	baseURL string
	client  *http.Client
}

func NewHTTPPackingGateway(baseURL string, timeout time.Duration) *HTTPPackingGateway {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:5051"
	}

	return &HTTPPackingGateway{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type PackRequest struct {
	Units     string          `json:"units"`
	Container PackContainerIn `json:"container"`
	Items     []PackItemIn    `json:"items"`
	Options   PackOptionsIn   `json:"options"`
}

type PackContainerIn struct {
	Length    float64 `json:"length"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	MaxWeight float64 `json:"max_weight"`
}

type PackItemIn struct {
	ItemID   string  `json:"item_id"`
	Label    string  `json:"label"`
	Length   float64 `json:"length"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Weight   float64 `json:"weight"`
	Quantity int     `json:"quantity"`
}

type PackOptionsIn struct {
	FixPoint            bool    `json:"fix_point"`
	CheckStable         bool    `json:"check_stable"`
	SupportSurfaceRatio float64 `json:"support_surface_ratio"`
	BiggerFirst         bool    `json:"bigger_first"`
	PutType             int     `json:"put_type"`
}

type PackResponse struct {
	Success bool          `json:"success"`
	Data    *PackDataOut  `json:"data,omitempty"`
	Error   *PackErrorOut `json:"error,omitempty"`
}

type PackErrorOut struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

type PackDataOut struct {
	Units      string             `json:"units"`
	Placements []PackPlacementOut `json:"placements"`
	Unfitted   []PackUnfittedOut  `json:"unfitted"`
	Stats      PackStatsOut       `json:"stats"`
}

type PackPlacementOut struct {
	ItemID     string  `json:"item_id"`
	Label      string  `json:"label"`
	PosX       float64 `json:"pos_x"`
	PosY       float64 `json:"pos_y"`
	PosZ       float64 `json:"pos_z"`
	Rotation   int     `json:"rotation"`
	StepNumber int     `json:"step_number"`
}

type PackUnfittedOut struct {
	ItemID string `json:"item_id"`
	Label  string `json:"label"`
	Count  int    `json:"count"`
}

type PackStatsOut struct {
	ExpandedItems int `json:"expanded_items"`
	FittedCount   int `json:"fitted_count"`
	UnfittedCount int `json:"unfitted_count"`
	PackTimeMs    int `json:"pack_time_ms"`
	TotalTimeMs   int `json:"total_time_ms"`
}

func (g *HTTPPackingGateway) Pack(ctx context.Context, req PackRequest) (*PackResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal pack request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/pack", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build /pack request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call packing service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read packing response: %w", err)
	}

	var decoded PackResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, fmt.Errorf("decode packing response (status %d): %w", resp.StatusCode, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := fmt.Sprintf("packing service returned %d", resp.StatusCode)
		if decoded.Error != nil && decoded.Error.Message != "" {
			msg = decoded.Error.Message
		}
		return nil, fmt.Errorf("packing service error: %s", msg)
	}

	if !decoded.Success {
		msg := "packing service returned success=false"
		if decoded.Error != nil && decoded.Error.Message != "" {
			msg = decoded.Error.Message
		}
		return nil, fmt.Errorf("packing service error: %s", msg)
	}
	if decoded.Data == nil {
		return nil, fmt.Errorf("packing service error: missing data")
	}

	return &decoded, nil
}
