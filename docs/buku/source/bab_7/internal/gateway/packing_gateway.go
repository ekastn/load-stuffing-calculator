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

// HTTPPackingGateway adalah implementasi PackingGateway yang berkomunikasi via HTTP
type HTTPPackingGateway struct {
	baseURL string       // URL base Packing Service
	client  *http.Client // HTTP client dengan timeout
}

// NewHTTPPackingGateway membuat instance gateway dengan konfigurasi yang diberikan
func NewHTTPPackingGateway(baseURL string, timeout time.Duration) *HTTPPackingGateway {
	// Hapus trailing slash untuk konsistensi
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:5000" // Default jika tidak dikonfigurasi
	}

	return &HTTPPackingGateway{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout, // Timeout untuk seluruh request
		},
	}
}

// Pack mengirim request ke Packing Service dan mengembalikan hasil
func (g *HTTPPackingGateway) Pack(ctx context.Context, req PackRequest) (*PackResponse, error) {
	// 1. Serialize request ke JSON
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal pack request: %w", err)
	}

	// 2. Buat HTTP request dengan context
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		g.baseURL+"/pack",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("build /pack request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 3. Kirim request
	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call packing service: %w", err)
	}
	defer resp.Body.Close()

	// 4. Baca response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read packing response: %w", err)
	}

	// 5. Parse JSON response
	var decoded PackResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, fmt.Errorf("decode packing response (status %d): %w", resp.StatusCode, err)
	}

	// 6. Handle HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := fmt.Sprintf("packing service returned %d", resp.StatusCode)
		if decoded.Error != nil && decoded.Error.Message != "" {
			msg = decoded.Error.Message
		}
		return nil, fmt.Errorf("packing service error: %s", msg)
	}

	// 7. Handle application-level errors
	if !decoded.Success {
		msg := "packing service returned success=false"
		if decoded.Error != nil && decoded.Error.Message != "" {
			msg = decoded.Error.Message
		}
		return nil, fmt.Errorf("packing service error: %s", msg)
	}

	return &decoded, nil
}
